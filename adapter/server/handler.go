package server

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/republicprotocol/swapperd/adapter/wallet"
	"github.com/republicprotocol/swapperd/core/balance"
	"github.com/republicprotocol/swapperd/core/bootload"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
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
	bootloader   bootload.Bootloader
}

// The Handler for swapperd requests
type Handler interface {
	GetInfo() GetInfoResponse
	GetSwaps(chan<- swap.ReceiptQuery) (GetSwapsResponse, error)
	GetBalances(chan<- balance.BalanceQuery) GetBalancesResponse
	PostTransfers(PostTransfersRequest) (PostTransfersResponse, error)
	PostSwaps(PostSwapRequest, chan<- swap.SwapReceipt, chan<- swap.SwapBlob) (PostSwapResponse, error)
	PostBootload(password string, receipts chan<- swap.SwapReceipt, swaps chan<- swap.SwapBlob)
	VerifyPassword(password string) bool
}

func NewHandler(passwordHash []byte, wallet wallet.Wallet, storage Storage, bootloader bootload.Bootloader) Handler {
	return &handler{false, passwordHash, wallet, storage, bootloader}
}

func (handler *handler) GetInfo() GetInfoResponse {
	return GetInfoResponse{
		Version:              "0.2.0",
		SupportedBlockchains: handler.wallet.SupportedBlockchains(),
		SupportedTokens:      handler.wallet.SupportedTokens(),
	}
}

func (handler *handler) GetSwaps(statuses chan<- swap.ReceiptQuery) (GetSwapsResponse, error) {
	resp := GetSwapsResponse{}
	if !handler.bootloaded {
		return resp, NewErrBootloadRequired("get swaps")
	}
	responder := make(chan map[swap.SwapID]swap.SwapReceipt)
	statuses <- swap.ReceiptQuery{Responder: responder}
	statusMap := <-responder
	for _, status := range statusMap {
		resp.Swaps = append(resp.Swaps, status)
	}
	return resp, nil
}

func (handler *handler) GetBalances(balanceQueries chan<- balance.BalanceQuery) GetBalancesResponse {
	response := make(chan map[blockchain.TokenName]blockchain.Balance)
	query := balance.BalanceQuery{
		Response: response,
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
	blob.Password = ""
	handler.storage.InsertSwap(blob)
	blob.Password = password
	go func() {
		swaps <- blob
		receipts <- swap.NewSwapReceipt(blob)
	}()
	return PostSwapResponse{}, nil
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

func (handler *handler) PostBootload(password string, receipts chan<- swap.SwapReceipt, swaps chan<- swap.SwapBlob) {
	handler.bootloader.Bootload(password, receipts, swaps)
}

func (handler *handler) VerifyPassword(password string) bool {
	return (bcrypt.CompareHashAndPassword(handler.passwordHash, []byte(password)) != nil)
}

func (handler *handler) patchSwap(req swap.SwapBlob) (swap.SwapBlob, error) {
	if err := handler.validateTokenDetails(req); err != nil {
		return swap.SwapBlob{}, err
	}
	return patchSwapDetails(req)
}

func (handler *handler) validateTokenDetails(blob swap.SwapBlob) error {
	// verify send details
	if err := handler.verifyTokenDetails(blob.SendToken, blob.SendTo, blob.SendAmount); err != nil {
		return err
	}
	// verify receive details
	return handler.verifyTokenDetails(blob.ReceiveToken, blob.ReceiveFrom, blob.ReceiveAmount)
}

func patchSwapDetails(swapBlob swap.SwapBlob) (swap.SwapBlob, error) {
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

func (handler *handler) verifyTokenDetails(tokenString, addressString, amountString string) error {
	token, err := blockchain.PatchToken(tokenString)
	if err != nil {
		return err
	}
	amount, ok := big.NewInt(0).SetString(amountString, 10)
	if !ok {
		return fmt.Errorf("invalid amount %s", amountString)
	}
	if err := handler.wallet.VerifyAddress(token.Blockchain, addressString); err != nil {
		return err
	}
	return handler.wallet.VerifyBalance(token, amount)
}
