package server

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/adapter/wallet"
	"github.com/republicprotocol/swapperd/core/balance"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

func NewErrBootloadRequired(msg string) error {
	return fmt.Errorf("please bootload before calling %s", msg)
}

type handler struct {
	bootloaded   bool
	passwordHash []byte
	wallet       wallet.Wallet
	storage      Storage
	logger       logrus.FieldLogger
}

// The Handler for swapperd requests
type Handler interface {
	GetInfo() GetInfoResponse
	GetSwaps(chan<- status.ReceiptQuery) (GetSwapsResponse, error)
	GetBalances(chan<- balance.BalanceQuery) GetBalancesResponse
	GetAddresses() (GetAddressesResponse, error)
	PostTransfers(PostTransfersRequest) (PostTransfersResponse, error)
	PostSwaps(PostSwapRequest, chan<- swap.SwapReceipt, chan<- swap.SwapBlob) (PostSwapResponse, error)
	PostDelayedSwaps(PostSwapRequest, chan<- swap.SwapReceipt, chan<- swap.SwapBlob) error
	PostBootload(password string, swaps, delayedSwaps chan<- swap.SwapBlob) error
	VerifyPassword(password string) bool
}

func NewHandler(passwordHash []byte, wallet wallet.Wallet, storage Storage, logger logrus.FieldLogger) Handler {
	return &handler{false, passwordHash, wallet, storage, logger}
}

func (handler *handler) GetInfo() GetInfoResponse {
	return GetInfoResponse{
		Version:              "0.2.0",
		Bootloaded:           handler.bootloaded,
		SupportedBlockchains: handler.wallet.SupportedBlockchains(),
		SupportedTokens:      handler.wallet.SupportedTokens(),
	}
}

func (handler *handler) GetAddresses() (GetAddressesResponse, error) {
	return handler.wallet.Addresses()
}

func (handler *handler) GetSwaps(statuses chan<- status.ReceiptQuery) (GetSwapsResponse, error) {
	resp := GetSwapsResponse{}
	if !handler.bootloaded {
		return resp, NewErrBootloadRequired("get swaps")
	}
	responder := make(chan map[swap.SwapID]swap.SwapReceipt)
	statuses <- status.ReceiptQuery{Responder: responder}
	statusMap := <-responder
	for _, status := range statusMap {
		resp.Swaps = append(resp.Swaps, status)
	}
	return resp, nil
}

func (handler *handler) GetBalances(balanceQueries chan<- balance.BalanceQuery) GetBalancesResponse {
	response := make(chan map[blockchain.TokenName]blockchain.Balance)
	query := balance.BalanceQuery{
		Responder: response,
	}
	balanceQueries <- query
	resp := <-response
	return resp
}

func (handler *handler) PostSwaps(swapReq PostSwapRequest, receipts chan<- swap.SwapReceipt, swaps chan<- swap.SwapBlob) (PostSwapResponse, error) {
	if !handler.bootloaded {
		return PostSwapResponse{}, NewErrBootloadRequired("get swaps")
	}
	password := swapReq.Password
	blob, err := handler.patchSwap(swap.SwapBlob(swapReq))
	if err != nil {
		return PostSwapResponse{}, err
	}

	receipt := handler.newSwapReceipt(blob)
	blob.Password = ""

	if err := handler.storage.PutSwap(blob); err != nil {
		return PostSwapResponse{}, err
	}

	blob.Password = password
	go func() {
		swaps <- blob
		receipts <- receipt
	}()
	return PostSwapResponse{}, nil
}

func (handler *handler) PostDelayedSwaps(swapReq PostSwapRequest, receipts chan<- swap.SwapReceipt, swaps chan<- swap.SwapBlob) error {
	if !handler.bootloaded {
		return NewErrBootloadRequired("get swaps")
	}
	password := swapReq.Password
	blob, err := handler.patchDelayedSwap(swap.SwapBlob(swapReq))
	if err != nil {
		return err
	}

	receipt := handler.newSwapReceipt(blob)
	blob.Password = ""
	if err := handler.storage.PutSwap(blob); err != nil {
		return err
	}
	blob.Password = password
	go func() {
		swaps <- blob
		receipts <- receipt
	}()

	return nil
}

func (handler *handler) PostTransfers(req PostTransfersRequest) (PostTransfersResponse, error) {
	response := PostTransfersResponse{}
	token, err := blockchain.PatchToken(req.Token)
	if err != nil {
		return response, err
	}

	if err := handler.wallet.VerifyAddress(token.Blockchain, req.To); err != nil {
		return response, err
	}

	amount, ok := big.NewInt(0).SetString(req.Amount, 10)
	if !ok {
		return response, fmt.Errorf("invalid amount %s", req.Amount)
	}

	if err := handler.wallet.VerifyBalance(token, amount); err != nil {
		return response, err
	}

	txHash, err := handler.wallet.Transfer(req.Password, token, req.To, amount)
	if err != nil {
		return response, err
	}

	response.TxHash = txHash
	return response, nil
}

func (handler *handler) PostBootload(password string, swaps, delayedSwaps chan<- swap.SwapBlob) error {
	if handler.bootloaded {
		return fmt.Errorf("already bootloaded")
	}

	pendingSwaps, err := handler.storage.PendingSwaps()
	if err != nil {
		return err
	}

	handler.logger.Infof("loading %d pending atomic swaps", len(pendingSwaps))

	co.ParForAll(pendingSwaps, func(i int) {
		handler.logger.Info(pendingSwaps[i])
		swap := pendingSwaps[i]
		swap.Password = password
		if swap.Delay {
			delayedSwaps <- swap
			return
		}
		swaps <- swap
	})

	handler.bootloaded = true
	return nil
}

func (handler *handler) VerifyPassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword(handler.passwordHash, []byte(password)); err != nil {
		handler.logger.Info("password length", len(handler.passwordHash))
		handler.logger.Error(err)
	}
	return (bcrypt.CompareHashAndPassword(handler.passwordHash, []byte(password)) == nil)
}

func (handler *handler) patchSwap(swapBlob swap.SwapBlob) (swap.SwapBlob, error) {
	// verify send details
	if err := handler.verifyTokenDetails(swapBlob.SendToken, swapBlob.SendTo, swapBlob.SendAmount, true); err != nil {
		return swapBlob, err
	}
	// verify receive details
	if err := handler.verifyTokenDetails(swapBlob.ReceiveToken, swapBlob.ReceiveFrom, swapBlob.ReceiveAmount, true); err != nil {
		return swapBlob, err
	}

	swapID := [32]byte{}
	rand.Read(swapID[:])
	swapBlob.ID = swap.SwapID(base64.StdEncoding.EncodeToString(swapID[:]))
	secret := [32]byte{}
	if swapBlob.ShouldInitiateFirst {
		swapBlob.TimeLock = time.Now().Unix() + 3*swap.ExpiryUnit
		secret = sha3.Sum256(append([]byte(swapBlob.Password), []byte(swapBlob.ID)...))
		hash := sha256.Sum256(secret[:])
		swapBlob.SecretHash = base64.StdEncoding.EncodeToString(hash[:])
		return swapBlob, nil
	}
	secretHash, err := base64.StdEncoding.DecodeString(swapBlob.SecretHash)
	if len(secretHash) != 32 || err != nil {
		return swapBlob, fmt.Errorf("invalid secret hash")
	}
	if time.Now().Unix()+2*swap.ExpiryUnit > swapBlob.TimeLock {
		return swapBlob, fmt.Errorf("not enough time to do the atomic swap")
	}
	return swapBlob, nil
}

func (handler *handler) patchDelayedSwap(blob swap.SwapBlob) (swap.SwapBlob, error) {
	if blob.DelayCallbackURL == "" {
		return blob, fmt.Errorf("delay url cannot be empty")
	}

	swapID := [32]byte{}
	rand.Read(swapID[:])
	blob.ID = swap.SwapID(base64.StdEncoding.EncodeToString(swapID[:]))

	if err := handler.verifyTokenDetails(blob.SendToken, blob.SendTo, blob.SendAmount, false); err != nil {
		return blob, err
	}

	if _, err := blockchain.PatchToken(blob.ReceiveToken); err != nil {
		return blob, err
	}

	secret := sha3.Sum256(append([]byte(blob.ID), []byte(blob.Password)...))
	secretHash := sha256.Sum256(secret[:])
	blob.SecretHash = base64.StdEncoding.EncodeToString(secretHash[:])
	blob.TimeLock = time.Now().Unix() + 3*swap.ExpiryUnit
	return blob, nil
}

func (handler *handler) verifyTokenDetails(tokenString, addressString, amountString string, verifyAddress bool) error {
	token, err := blockchain.PatchToken(tokenString)
	if err != nil {
		return err
	}
	amount, ok := big.NewInt(0).SetString(amountString, 10)
	if !ok {
		return fmt.Errorf("invalid amount %s", amountString)
	}

	if verifyAddress {
		if err := handler.wallet.VerifyAddress(token.Blockchain, addressString); err != nil {
			return err
		}
	}

	return handler.wallet.VerifyBalance(token, amount)
}

func (handler *handler) newSwapReceipt(blob swap.SwapBlob) swap.SwapReceipt {
	return swap.SwapReceipt{blob.ID, blob.SendToken, blob.ReceiveToken, blob.SendAmount, blob.ReceiveAmount, blockchain.Cost{}, blockchain.Cost{}, time.Now().Unix(), 1, blob.Delay, blob.DelayInfo}
}
