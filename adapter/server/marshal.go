package server

import (
	"encoding/json"

	"github.com/republicprotocol/swapperd/core/transfer"
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

type GetSwapResponse swap.SwapReceipt

type GetBalanceResponse blockchain.Balance
type GetBalancesResponse map[blockchain.TokenName]blockchain.Balance

type GetAddressesResponse map[blockchain.TokenName]string
type GetAddressResponse string

type PostSwapRequest swap.SwapBlob

type GetIDResponse struct {
	PublicKey string `json:"publicKey"`
}

type PostSwapResponse struct {
	ID        swap.SwapID   `json:"id"`
	Swap      swap.SwapBlob `json:"swap,omitempty"`
	Signature string        `json:"signature,omitempty"`
}

type PostRedeemSwapResponse struct {
	ID swap.SwapID `json:"id"`
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

type GetTransfersResponse struct {
	Transfers []transfer.TransferReceipt `json:"transfers"`
}

func MarshalGetTransfersResponse(receiptMap transfer.TransferReceiptMap) GetTransfersResponse {
	transfers := []transfer.TransferReceipt{}
	for _, receipt := range receiptMap {
		receipt.PasswordHash = ""
		transfers = append(transfers, receipt)
	}
	return GetTransfersResponse{
		Transfers: transfers,
	}
}
