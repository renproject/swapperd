package server

import (
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
)

type GetInfoResponse struct {
	Version              string                  `json:"version"`
	Bootloaded           bool                    `json:"bootloaded"`
	SupportedBlockchains []blockchain.Blockchain `json:"supportedBlockchains"`
	SupportedTokens      []blockchain.Token      `json:"supportedTokens"`
}

type GetSwapsResponse struct {
	Swaps []swap.SwapReceipt `json:"swaps"`
}

type GetBalancesResponse map[blockchain.TokenName]blockchain.Balance

type GetAddressesResponse map[blockchain.TokenName]string

type PostSwapRequest swap.SwapBlob

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
