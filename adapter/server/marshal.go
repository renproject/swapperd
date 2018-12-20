package server

import (
	"encoding/json"

	"github.com/republicprotocol/swapperd/core/wallet/transfer"

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

type GetIDResponse struct {
	PublicKey string `json:"publicKey"`
}

type PostSwapResponse struct {
	Swap      swap.SwapBlob `json:"swap"`
	Signature string        `json:"signature"`
}

type PostTransfersRequest struct {
	Token    string `json:"token"`
	To       string `json:"to"`
	Amount   string `json:"amount"`
	Password string `json:"password"`
}

type PostTransfersResponse transfer.TransferReceipt

type GetSignatureResponseJSON struct {
	Message   json.RawMessage `json:"message"`
	Signature string          `json:"signature"`
}

type GetSignatureResponseString struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

type GetTransfersResponse transfer.TransferReceiptMap
