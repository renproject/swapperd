package server

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/republicprotocol/swapperd/adapter/wallet"
	"github.com/republicprotocol/swapperd/core/balance"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
)

type handler struct {
	passwordHash []byte
	wallet       wallet.Wallet
}

// The Handler for swapperd requests
type Handler interface {
	GetInfo() GetInfoResponse
	GetSwaps(chan<- swap.ReceiptQuery) GetSwapsResponse
	GetBalances(chan<- balance.BalanceQuery) GetBalancesResponse
	PostSwaps(PostSwapRequest, chan<- swap.SwapRequest) (PostSwapResponse, error)
	PostTransfers(PostTransfersRequest) (PostTransfersResponse, error)
}

func NewHandler(passwordHash []byte, wallet wallet.Wallet) Handler {
	return &handler{passwordHash, wallet}
}

func (handler *handler) GetInfo() GetInfoResponse {
	return GetInfoResponse{
		Version:              "0.2.0",
		SupportedBlockchains: handler.wallet.SupportedBlockchains(),
		SupportedTokens:      handler.wallet.SupportedTokens(),
	}
}

func (handler *handler) GetSwaps(statuses chan<- swap.ReceiptQuery) GetSwapsResponse {
	resp := GetSwapsResponse{}
	responder := make(chan map[swap.SwapID]swap.SwapReceipt)
	statuses <- swap.ReceiptQuery{Responder: responder}
	statusMap := <-responder
	for _, status := range statusMap {
		resp.Swaps = append(resp.Swaps, status)
	}
	return resp
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

func (handler *handler) PostSwaps(req PostSwapRequest, swaps chan<- swap.SwapRequest) (PostSwapResponse, error) {
	blob, secret, err := handler.patchSwap(req)
	if err != nil {
		return PostSwapResponse{}, err
	}
	swaps <- swap.NewSwapRequest(blob, secret, req.Password)
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

func (handler *handler) verifyPassword(password string) bool {
	return (bcrypt.CompareHashAndPassword(handler.passwordHash, []byte(password)) != nil)
}

func (handler *handler) patchSwap(req PostSwapRequest) (swap.SwapBlob, [32]byte, error) {
	if err := handler.validateTokenDetails(req.SwapBlob, req.Password); err != nil {
		return swap.SwapBlob{}, [32]byte{}, err
	}
	return patchSwapDetails(req.SwapBlob)
}

func (handler *handler) validateTokenDetails(swapBlob swap.SwapBlob, password string) error {
	// verify send details
	if err := handler.verifyTokenDetails(swapBlob.SendToken, swapBlob.SendTo, swapBlob.SendAmount); err != nil {
		return err
	}
	// verify receive details
	return handler.verifyTokenDetails(swapBlob.ReceiveToken, swapBlob.ReceiveFrom, swapBlob.ReceiveAmount)
}

func patchSwapDetails(swapBlob swap.SwapBlob) (swap.SwapBlob, [32]byte, error) {
	swapID := [32]byte{}
	rand.Read(swapID[:])
	swapBlob.ID = swap.SwapID(base64.StdEncoding.EncodeToString(swapID[:]))
	secret := [32]byte{}
	if swapBlob.ShouldInitiateFirst {
		swapBlob.TimeLock = time.Now().Unix() + 3*swap.ExpiryUnit
		rand.Read(secret[:])
		hash := sha256.Sum256(secret[:])
		swapBlob.SecretHash = base64.StdEncoding.EncodeToString(hash[:])
		return swapBlob, secret, nil
	}
	secretHash, err := base64.StdEncoding.DecodeString(swapBlob.SecretHash)
	if len(secretHash) != 32 || err != nil {
		return swapBlob, secret, fmt.Errorf("invalid secret hash")
	}
	if time.Now().Unix()+2*swap.ExpiryUnit > swapBlob.TimeLock {
		return swapBlob, secret, fmt.Errorf("not enough time to do the atomic swap")
	}
	return swapBlob, secret, nil
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
