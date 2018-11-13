package server

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/republicprotocol/swapperd/adapter/fund"
	"github.com/republicprotocol/swapperd/core/auth"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation"
)

type server struct {
	authenticator  auth.Authenticator
	fundManager    fund.Manager
	swapperQueries chan<- swapper.Swap
	statusQueries  chan<- status.Query
}

func NewServer(authenticator auth.Authenticator, fundManager fund.Manager, swaps chan<- swapper.Swap, statusQueries chan<- status.Query) *server {
	return &server{authenticator, fundManager, swaps, statusQueries}
}

func (server *server) GetWhoAmI() GetWhoAmIResponse {
	return GetWhoAmIResponse{
		Version:              "0.2.0",
		SupportedBlockchains: server.fundManager.SupportedBlockchains(),
		SupportedTokens:      server.fundManager.SupportedTokens(),
	}
}

func (server *server) GetSwaps() GetSwapsResponse {
	resp := GetSwapsResponse{}
	responder := make(chan map[foundation.SwapID]foundation.SwapStatus)
	server.statusQueries <- status.Query{Responder: responder}
	swapMap := <-responder
	for _, status := range swapMap {
		resp.Swaps = append(resp.Swaps, SwapStatus{
			ID:     string(status.ID),
			Status: status.Status,
		})
	}
	return resp
}

func (server *server) PostSwaps(swap foundation.SwapBlob, password string) (foundation.SwapBlob, error) {
	swap, secret, err := server.patchSwap(swap, password)
	if err != nil {
		return swap, err
	}
	server.swapperQueries <- swapper.NewSwap(swap, secret, password)
	return swap, nil
}

func (server *server) GetBalances() (GetBalanceResponse, error) {
	resp := GetBalanceResponse{}
	balanceBook, err := server.fundManager.Balances()
	if err != nil {
		return resp, err
	}
	for token, balance := range balanceBook {
		resp.Balances = append(resp.Balances, Balance{
			Token:   token.Name,
			Address: balance.Address,
			Amount:  balance.Amount.String(),
		})
	}
	return resp, nil
}

func (server *server) PostWithdraw(password string, postWithdrawRequest PostWithdrawalsRequest) (PostWithdrawalsResponse, error) {
	response := PostWithdrawalsResponse{}
	token, err := UnmarshalToken(postWithdrawRequest.Token)
	if err != nil {
		return response, err
	}
	if postWithdrawRequest.Amount == "" {
		txHash, err := server.fundManager.Withdraw(password, token, postWithdrawRequest.To, nil)
		if err != nil {
			return response, err
		}
		response.TxHash = txHash
	}
	value, ok := big.NewInt(0).SetString(postWithdrawRequest.Amount, 10)
	if !ok {
		return response, fmt.Errorf("invalid amount")
	}
	txHash, err := server.fundManager.Withdraw(password, token, postWithdrawRequest.To, value)
	if err != nil {
		return response, err
	}
	response.TxHash = txHash
	return response, nil
}

func (server *server) patchSwap(swapBlob foundation.SwapBlob, password string) (foundation.SwapBlob, [32]byte, error) {
	if err := server.validateTokenDetails(swapBlob, password); err != nil {
		return foundation.SwapBlob{}, [32]byte{}, err
	}
	return patchSwapDetails(swapBlob)
}

func (server *server) validateTokenDetails(swapBlob foundation.SwapBlob, password string) error {
	balanceBook, err := server.fundManager.Balances()
	if err != nil {
		return err
	}

	sendToken, err := UnmarshalToken(swapBlob.SendToken)
	if err != nil {
		return err
	}
	amount, err := UnmarshalAmount(swapBlob.SendAmount)
	if err != nil {
		return err
	}
	if balanceBook[sendToken].Amount.Cmp(amount) < 0 {
		return fmt.Errorf("insufficient balance required: %v current: %v", amount, balanceBook[sendToken].Amount)
	}

	switch sendToken.Blockchain {
	case foundation.Ethereum:
		minVal, ok := big.NewInt(0).SetString("5000000000000000", 10) // 0.005 eth
		if !ok {
			return fmt.Errorf("invalid minimum value")
		}
		if balanceBook[foundation.TokenETH].Amount.Cmp(minVal) < 0 {
			return fmt.Errorf("minimum balance required to start an atomic swap on ethereum blockchain is 0.005 eth (to cover the transaction fees)")
		}
	case foundation.Bitcoin:
		if amount.Cmp(big.NewInt(10000)) < 0 {
			return fmt.Errorf("minimum send amount for bitcoin is 10000 sat")
		}
	}
	return nil
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
		swapBlob.SecretHash = MarshalSecretHash(hash)
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
