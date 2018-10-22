package storage

import (
	"encoding/json"
	"sync"

	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation"
)

type Store interface {
	Read([]byte) ([]byte, error)
	Write([]byte, []byte) error
	Delete([]byte) error
}

type storage struct {
	mu    *sync.RWMutex
	store Store
}

type SwapStorage struct {
	Swaps        []foundation.Swap `json:"swaps"`
	PendingSwaps []foundation.Swap `json:"pendingSwaps"`
}

func New(store Store) swapper.Storage {
	return &storage{
		mu:    new(sync.RWMutex),
		store: store,
	}
}

func (storage *storage) AddSwap(swap foundation.Swap) error {
	swapStorageBytes, err := storage.store.Read([]byte("SwapStorage"))
	if err != nil {
		return err
	}

	swapStorage := SwapStorage{}
	if err := json.Unmarshal(swapStorageBytes, &swapStorage); err != nil {
		return err
	}

	swapStorage.Swaps = append(swapStorage.Swaps, swap)
	swapStorage.PendingSwaps = append(swapStorage.Swaps, swap)

	swapStorageBytes, err = json.Marshal(swapStorage)
	if err != nil {
		return err
	}

	return storage.store.Write([]byte("SwapStorage"), swapStorageBytes)
}

func (storage *storage) DeleteSwap(swapID foundation.SwapID) error {
	swapStorageBytes, err := storage.store.Read([]byte("SwapStorage"))
	if err != nil {
		return err
	}

	swapStorage := SwapStorage{}
	if err := json.Unmarshal(swapStorageBytes, &swapStorage); err != nil {
		return err
	}

	for i, swap := range swapStorage.PendingSwaps {
		if swap.ID == swapID {
			swapStorage.PendingSwaps = append(swapStorage.PendingSwaps[:i], swapStorage.PendingSwaps[:i+1]...)
		}
	}

	swapStorageBytes, err = json.Marshal(swapStorage)
	if err != nil {
		return err
	}

	return storage.store.Write([]byte("SwapStorage"), swapStorageBytes)
}

func (storage *storage) LoadPendingSwaps() []foundation.Swap {
	swapStorageBytes, err := storage.store.Read([]byte("SwapStorage"))
	if err != nil {
		return []foundation.Swap{}
	}

	swapStorage := SwapStorage{}
	if err := json.Unmarshal(swapStorageBytes, &swapStorage); err != nil {
		return []foundation.Swap{}
	}

	return swapStorage.PendingSwaps
}

func (storage *storage) LoadSwaps() []foundation.Swap {
	swapStorageBytes, err := storage.store.Read([]byte("SwapStorage"))
	if err != nil {
		return []foundation.Swap{}
	}

	swapStorage := SwapStorage{}
	if err := json.Unmarshal(swapStorageBytes, &swapStorage); err != nil {
		return []foundation.Swap{}
	}

	return swapStorage.Swaps
}
