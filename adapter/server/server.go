package server

import (
	"github.com/republicprotocol/swapperd/foundation"
)

type server struct {
	swaps chan<- foundation.Swap
}

func NewServer(swaps chan<- foundation.Swap) *server {
	return &server{swaps}
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

func (server *server) GetSwaps() (GetSwapResponse, error) {
	// TODO: Implement the logic
	return GetSwapResponse{}, nil
}

func (server *server) PostSwaps(swapReqRes PostSwapRequestResponse) (PostSwapRequestResponse, error) {
	swap, err := UnmarshalSwapRequestResponse(swapReqRes)
	if err != nil {
		return PostSwapRequestResponse{}, err
	}
	server.swaps <- swap
	swapReqRes.SecretHash = MarshalSecretHash(swap.SecretHash)
	return swapReqRes, nil
}

func (server *server) GetBalances() (GetBalanceResponse, error) {
	// TODO: Implement the logic
	return GetBalanceResponse{}, nil
}
