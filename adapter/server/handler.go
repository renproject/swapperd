package server

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/republicprotocol/swapperd/adapter/fund"
	"github.com/republicprotocol/swapperd/foundation"
	"golang.org/x/crypto/sha3"
)

type handler struct {
	passwordHash [32]byte
	fundManager  fund.Manager
	swaps        chan<- foundation.SwapRequest
	statuses     chan<- foundation.StatusQuery
}

// The Handler for swapperd requests
type Handler interface {
	GetInfo() GetInfoResponse
	GetBalances() (GetBalancesResponse, error)
	GetSwaps() GetSwapsResponse
	PostSwaps(PostSwapRequest) (PostSwapResponse, error)
	PostTransfers(PostTransfersRequest) (PostTransfersResponse, error)
}

func NewHandler(passwordHash [32]byte, fundManager fund.Manager, swaps chan<- foundation.SwapRequest, statuses chan<- foundation.StatusQuery) Handler {
	return &handler{passwordHash, fundManager, swaps, statuses}
}

func (handler *handler) GetInfo() GetInfoResponse {
	return GetInfoResponse{
		Version:              "0.2.0",
		SupportedBlockchains: handler.fundManager.SupportedBlockchains(),
		SupportedTokens:      handler.fundManager.SupportedTokens(),
	}
}

func (handler *handler) GetSwaps() GetSwapsResponse {
	resp := GetSwapsResponse{}
	responder := make(chan map[foundation.SwapID]foundation.SwapStatus)
	handler.statuses <- foundation.StatusQuery{Responder: responder}
	statusMap := <-responder
	for _, status := range statusMap {
		resp.Swaps = append(resp.Swaps, status)
	}
	return resp
}

func (handler *handler) GetBalances() (GetBalancesResponse, error) {
	return handler.fundManager.Balances()
}

func (handler *handler) PostSwaps(req PostSwapRequest) (PostSwapResponse, error) {
	swap, secret, err := handler.patchSwap(req)
	if err != nil {
		return PostSwapResponse{}, err
	}
	handler.swaps <- foundation.NewSwapRequest(swap, secret, req.Password)
	return PostSwapResponse{}, nil
}

func (handler *handler) PostTransfers(req PostTransfersRequest) (PostTransfersResponse, error) {
	response := PostTransfersResponse{}
	token, err := foundation.PatchToken(req.Token)
	if err != nil {
		return response, err
	}

	if err := handler.fundManager.VerifyAddress(token.Blockchain, req.To); err != nil {
		return response, err
	}

	amount, ok := big.NewInt(0).SetString(req.Amount, 10)
	if !ok {
		return response, fmt.Errorf("invalid amount %s", req.Amount)
	}

	if err := handler.fundManager.VerifyBalance(token, amount); err != nil {
		return response, err
	}

	txHash, err := handler.fundManager.Transfer(req.Password, token, req.To, amount)
	if err != nil {
		return response, err
	}

	response.TxHash = txHash
	return response, nil
}

func (handler *handler) verifyPassword(password string) bool {
	passwordHash := sha3.Sum256([]byte(password))
	passwordOk := subtle.ConstantTimeCompare(passwordHash[:], handler.passwordHash[:]) == 1
	return passwordOk
}

func (handler *handler) patchSwap(req PostSwapRequest) (foundation.SwapBlob, [32]byte, error) {
	if err := handler.validateTokenDetails(req.SwapBlob, req.Password); err != nil {
		return foundation.SwapBlob{}, [32]byte{}, err
	}
	return patchSwapDetails(req.SwapBlob)
}

func (handler *handler) validateTokenDetails(swapBlob foundation.SwapBlob, password string) error {
	// verify send details
	if err := handler.verifyTokenDetails(swapBlob.SendToken, swapBlob.SendTo, swapBlob.SendAmount); err != nil {
		return err
	}
	// verify receive details
	return handler.verifyTokenDetails(swapBlob.ReceiveToken, swapBlob.ReceiveFrom, swapBlob.ReceiveAmount)
}

func patchSwapDetails(swapBlob foundation.SwapBlob) (foundation.SwapBlob, [32]byte, error) {
	swapID := [32]byte{}
	rand.Read(swapID[:])
	swapBlob.ID = foundation.SwapID(base64.StdEncoding.EncodeToString(swapID[:]))
	secret := [32]byte{}
	if swapBlob.ShouldInitiateFirst {
		swapBlob.TimeLock = time.Now().Unix() + 3*foundation.ExpiryUnit
		rand.Read(secret[:])
		hash := sha256.Sum256(secret[:])
		swapBlob.SecretHash = base64.StdEncoding.EncodeToString(hash[:])
		return swapBlob, secret, nil
	}
	secretHash, err := base64.StdEncoding.DecodeString(swapBlob.SecretHash)
	if len(secretHash) != 32 || err != nil {
		return swapBlob, secret, fmt.Errorf("invalid secret hash")
	}
	if time.Now().Unix()+2*foundation.ExpiryUnit > swapBlob.TimeLock {
		return swapBlob, secret, fmt.Errorf("not enough time to do the atomic swap")
	}
	return swapBlob, secret, nil
}

func (handler *handler) verifyTokenDetails(tokenString, addressString, amountString string) error {
	token, err := foundation.PatchToken(tokenString)
	if err != nil {
		return err
	}
	amount, ok := big.NewInt(0).SetString(amountString, 10)
	if !ok {
		return fmt.Errorf("invalid amount %s", amountString)
	}
	if err := handler.fundManager.VerifyAddress(token.Blockchain, addressString); err != nil {
		return err
	}
	return handler.fundManager.VerifyBalance(token, amount)
}
