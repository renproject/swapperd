package request

import "github.com/republicprotocol/swapperd/foundation"

type GetInfoResponse struct {
	Version              string                  `json:"version"`
	SupportedBlockchains []foundation.Blockchain `json:"supportedBlockchains"`
	SupportedTokens      []foundation.Token      `json:"supportedTokens"`
}

type GetSwapsResponse struct {
	Swaps []foundation.SwapStatus `json:"swaps"`
}

type GetBalancesResponse map[foundation.TokenName]foundation.Balance

type PostSwapRequest struct {
	foundation.SwapBlob
	Password string
}

type PostSwapResponse struct {
}

type PostTransfersRequest struct {
	Token    string `json:"token"`
	To       string `json:"to"`
	Amount   string `json:"amount"`
	Password string `json:"password"`
}

type PostTransfersResponse struct {
	TxHash string `json:"txHash"`
}
