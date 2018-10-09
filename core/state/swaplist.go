package state

import (
	"sync"

	"github.com/republicprotocol/swapperd/foundation"
)

type ActiveSwapList interface {
	AddSwap(swap foundation.Swap) error
	DeleteSwap(swapID foundation.SwapID) error
	LoadSwaps() []foundation.Swap
}

type ProtectedSwapMap struct {
	mu *sync.RWMutex
	SwapMap
}

func (state *state) AddSwap(swap foundation.Swap) error {

}

func (state *state) DeleteSwap(orderID [32]byte) error {

}

func (state *state) LoadSwaps() [][32]byte {

}
