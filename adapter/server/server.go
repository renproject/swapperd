package server

import (
	"crypto/rand"
	"crypto/sha256"

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
	swap, err := decodePostSwap(swapReqRes)
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

func decodePostSwap(swapReqRes PostSwapRequestResponse) (foundation.Swap, error) {
	secret := [32]byte{}
	if swapReqRes.ShouldInitiateFirst {
		rand.Read(secret[:])
		hash := sha256.Sum256(secret[:])
		swapReqRes.SecretHash = MarshalSecretHash(hash)
	}
	swapID, err := UnmarshalSwapID(swapReqRes.ID)
	if err != nil {
		return foundation.Swap{}, nil
	}
	sendToken, err := UnmarshalToken(swapReqRes.SendToken)
	if err != nil {
		return foundation.Swap{}, nil
	}
	receiveToken, err := UnmarshalToken(swapReqRes.ReceiveToken)
	if err != nil {
		return foundation.Swap{}, nil
	}
	sendValue, err := UnmarshalAmount(swapReqRes.SendAmount)
	if err != nil {
		return foundation.Swap{}, nil
	}
	receiveValue, err := UnmarshalAmount(swapReqRes.ReceiveAmount)
	if err != nil {
		return foundation.Swap{}, nil
	}
	secretHash, err := UnmarshalSecretHash(swapReqRes.SecretHash)
	if err != nil {
		return foundation.Swap{}, nil
	}
	return foundation.Swap{
		ID:                 swapID,
		Secret:             secret,
		SecretHash:         secretHash,
		TimeLock:           swapReqRes.TimeLock,
		SendToAddress:      swapReqRes.SendTo,
		ReceiveFromAddress: swapReqRes.ReceiveFrom,
		SendValue:          sendValue,
		ReceiveValue:       receiveValue,
		SendToken:          sendToken,
		ReceiveToken:       receiveToken,
		IsFirst:            swapReqRes.ShouldInitiateFirst,
	}, nil
}
