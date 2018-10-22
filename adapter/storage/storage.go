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
	Swaps          []foundation.Swap `json:"swaps"`
	PendingQueries []swapper.Query   `json:"pendingQueries"`
}

func New(store Store) swapper.Storage {
	return &storage{
		mu:    new(sync.RWMutex),
		store: store,
	}
}

func (storage *storage) AddQuery(query swapper.Query) error {
	swapStorageBytes, err := storage.store.Read([]byte("SwapStorage"))
	if err != nil {
		return err
	}

	swapStorage := SwapStorage{}
	if err := json.Unmarshal(swapStorageBytes, &swapStorage); err != nil {
		return err
	}

	swapStorage.Swaps = append(swapStorage.Swaps, query.Swap)
	swapStorage.PendingQueries = append(swapStorage.PendingQueries, query)

	swapStorageBytes, err = json.Marshal(swapStorage)
	if err != nil {
		return err
	}

	return storage.store.Write([]byte("SwapStorage"), swapStorageBytes)
}

func (storage *storage) DeleteQuery(swapID foundation.SwapID) error {
	swapStorageBytes, err := storage.store.Read([]byte("SwapStorage"))
	if err != nil {
		return err
	}

	swapStorage := SwapStorage{}
	if err := json.Unmarshal(swapStorageBytes, &swapStorage); err != nil {
		return err
	}

	for i, query := range swapStorage.PendingQueries {
		if query.Swap.ID == swapID {
			swapStorage.PendingQueries = append(swapStorage.PendingQueries[:i], swapStorage.PendingQueries[:i+1]...)
		}
	}

	swapStorageBytes, err = json.Marshal(swapStorage)
	if err != nil {
		return err
	}

	return storage.store.Write([]byte("SwapStorage"), swapStorageBytes)
}

func (storage *storage) LoadPendingQueries() []swapper.Query {
	swapStorageBytes, err := storage.store.Read([]byte("SwapStorage"))
	if err != nil {
		return []swapper.Query{}
	}

	swapStorage := SwapStorage{}
	if err := json.Unmarshal(swapStorageBytes, &swapStorage); err != nil {
		return []swapper.Query{}
	}

	return swapStorage.PendingQueries
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
