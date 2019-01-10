package core

import "github.com/republicprotocol/swapperd/foundation/swap"

type SwapRequest swap.SwapBlob

func (msg SwapRequest) IsMessage() {
}

type Bootload struct {
	Password string
}

func (msg Bootload) IsMessage() {
}
