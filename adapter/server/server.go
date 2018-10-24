package server

import (
	"crypto/rand"
	"crypto/sha256"

	"github.com/republicprotocol/swapperd/adapter/balance"
	"github.com/republicprotocol/swapperd/core/auth"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation"
)

type server struct {
	authenticator  auth.Authenticator
	swapperQueries chan<- swapper.Swap
	statusQueries  chan<- status.Query
	balanceQueries chan<- balance.Query
}

func NewServer(authenticator auth.Authenticator, swaps chan<- swapper.Swap, statusQueries chan<- status.Query, balanceQueries chan<- balance.Query) *server {
	return &server{authenticator, swaps, statusQueries, balanceQueries}
}

func (server *server) GetPing() GetPingResponse {
	return GetPingResponse{
		Version: "0.1.0",
		SupportedTokens: []foundation.Token{
			foundation.TokenBTC,
			foundation.TokenETH,
			foundation.TokenWBTC,
		},
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
	query, responder, errs := balance.NewQuery(password)
	server.balanceQueries <- query
	if err := <-errs; err != nil {
		return resp, err
	}
	balanceMap := <-responder
	for token, balance := range balanceMap {
		resp.Balances = append(resp.Balances, Balance{
			TokenName: token.Name,
			Address:   balance.Address,
			Amount:    balance.Amount.String(),
		})
	}
	return resp, nil
}
