package server

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/republicprotocol/swapperd/adapter/funds"
	"github.com/republicprotocol/swapperd/core/auth"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation"
)

type server struct {
	authenticator  auth.Authenticator
	fundManager    funds.Manager
	swapperQueries chan<- swapper.Swap
	statusQueries  chan<- status.Query
}

func NewServer(authenticator auth.Authenticator, fundManager funds.Manager, swaps chan<- swapper.Swap, statusQueries chan<- status.Query) *server {
	return &server{authenticator, fundManager, swaps, statusQueries}
}

func (server *server) GetPing() GetPingResponse {
	return GetPingResponse{
		Version:         "0.1.0",
		SupportedTokens: server.fundManager.SupportedTokens(),
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

func (server *server) PostSwaps(swap foundation.SwapBlob, password string) foundation.SwapBlob {
	secret := [32]byte{}
	if swap.ShouldInitiateFirst {
		rand.Read(secret[:])
		hash := sha256.Sum256(secret[:])
		swap.SecretHash = MarshalSecretHash(hash)
	}
	server.swapperQueries <- swapper.NewSwap(swap, secret, password)
	return swap
}

func (server *server) GetBalances(password string) (GetBalanceResponse, error) {
	resp := GetBalanceResponse{}
	balanceMap, err := server.fundManager.Balances(password)
	if err != nil {
		return resp, err
	}
	for token, balance := range balanceMap {
		resp.Balances = append(resp.Balances, Balance{
			Token:   token.Name,
			Address: balance.Address,
			Amount:  balance.Amount.String(),
		})
	}
	return resp, nil
}

func (server *server) PostWithdraw(password string, postWithdrawRequest PostWithdrawRequest) error {
	token, err := UnmarshalToken(postWithdrawRequest.Token)
	if err != nil {
		return err
	}
	if postWithdrawRequest.Amount == "" {
		return server.fundManager.Withdraw(password, token, postWithdrawRequest.To, nil)
	}
	value, ok := big.NewInt(0).SetString(postWithdrawRequest.Amount, 10)
	if !ok {
		return fmt.Errorf("invalid amount")
	}
	return server.fundManager.Withdraw(password, token, postWithdrawRequest.To, value)
}
