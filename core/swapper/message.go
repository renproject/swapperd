package swapper

import (
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
)

type SwapRequest struct {
	Blob        swap.SwapBlob
	SendCost    blockchain.Cost
	ReceiveCost blockchain.Cost
}

func (msg SwapRequest) IsMessage() {
}

func NewSwapRequest(blob swap.SwapBlob, sendCost, receiveCost blockchain.Cost) SwapRequest {
	return SwapRequest{
		Blob:        blob,
		SendCost:    sendCost,
		ReceiveCost: receiveCost,
	}
}

type ReceiptUpdate swap.ReceiptUpdate

func (msg ReceiptUpdate) IsMessage() {
}

func NewReceiptUpdate(id swap.SwapID, status int, native, foreign Contract) ReceiptUpdate {
	return ReceiptUpdate(swap.NewReceiptUpdate(id, func(receipt *swap.SwapReceipt) {
		receipt.Status = status
		receipt.SendCost = blockchain.CostToCostBlob(native.Cost())
		receipt.ReceiveCost = blockchain.CostToCostBlob(foreign.Cost())
	}))
}

type DeleteSwap struct {
	ID swap.SwapID
}

func (msg DeleteSwap) IsMessage() {
}
