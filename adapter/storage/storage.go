package storage

import (
	"github.com/republicprotocol/swapperd/foundation"
)

type Storage interface {
	AddSwap(swap foundation.Swap) error
	DeleteSwap(swapID foundation.SwapID) error
	LoadSwaps() []foundation.Swap
}

func NewStorage() Storage {
	return &storage{}
}

type storage struct {
}

func (storage *storage) AddSwap(swap foundation.Swap) error {
	return nil
}

func (storage *storage) DeleteSwap(orderID foundation.SwapID) error {
	return nil
}

func (storage *storage) LoadSwaps() []foundation.Swap {
	return []foundation.Swap{}
}

func (storage *storage) GetSwaps() []foundation.Swap {
	return []foundation.Swap{}
}
