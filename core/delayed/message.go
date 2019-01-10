package delayed

import "github.com/republicprotocol/swapperd/foundation/swap"

type DelayedSwapRequest swap.SwapBlob

func (msg DelayedSwapRequest) IsMessage() {
}

type SwapRequest swap.SwapBlob

func (msg SwapRequest) IsMessage() {
}

type ReceiptUpdate swap.ReceiptUpdate

func (msg ReceiptUpdate) IsMessage() {
}

type DeleteSwap struct {
	ID swap.SwapID
}

func (msg DeleteSwap) IsMessage() {
}
