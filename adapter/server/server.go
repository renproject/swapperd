package server

import (
	"github.com/republicprotocol/swapperd/core/auth"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation"
)

type server struct {
	authenticator  auth.Authenticator
	swapperQueries chan<- swapper.Query
	statusQueries  chan<- status.Query
}

func NewServer(authenticator auth.Authenticator, swaps chan<- swapper.Query, statusQueries chan<- status.Query) *server {
	return &server{authenticator, swaps, statusQueries}
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
	responder := make(chan map[foundation.SwapID]foundation.Status)
	server.statusQueries <- status.Query{Responder: responder}
	swapMap := <-responder
	for id, status := range swapMap {
		resp.Swaps = append(resp.Swaps, SwapStatus{
			ID:     MarshalSwapID(id),
			Status: int64(status),
		})
	}
	return resp
}

func (server *server) PostSwaps(swapReqRes PostSwapRequestResponse) (PostSwapRequestResponse, error) {
	swap, err := UnmarshalSwapRequestResponse(swapReqRes)
	if err != nil {
		return PostSwapRequestResponse{}, err
	}
	server.swapperQueries <- swapper.NewQuery(swap, "")
	return swapReqRes, nil
}

func (server *server) GetBalances() (GetBalanceResponse, error) {
	// TODO: Implement the logic
	return GetBalanceResponse{}, nil
}
