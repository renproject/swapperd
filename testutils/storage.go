package testutils

import (
	"errors"
	"sync"

	"github.com/republicprotocol/swapperd/foundation/swap"
)

type MockStorage struct {
	mu       *sync.RWMutex
	swaps    map[swap.SwapID]swap.SwapBlob
	receipts map[swap.SwapID]swap.SwapReceipt
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		mu:       new(sync.RWMutex),
		swaps:    map[swap.SwapID]swap.SwapBlob{},
		receipts: map[swap.SwapID]swap.SwapReceipt{},
	}
}

func (store *MockStorage) InsertSwap(swapBlob swap.SwapBlob) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.swaps[swapBlob.ID] = swapBlob
	store.receipts[swapBlob.ID] = swap.NewSwapReceipt(swapBlob)
	return nil
}

func (store *MockStorage) PendingSwap(id swap.SwapID) (swap.SwapBlob, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	swapblob, ok := store.swaps[id]
	if !ok {
		return swap.SwapBlob{}, errors.New("pending swap not found")
	}
	return swapblob, nil
}

func (store *MockStorage) DeletePendingSwap(id swap.SwapID) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	delete(store.swaps, id)
	return nil
}

func (store *MockStorage) PendingSwaps() ([]swap.SwapBlob, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	swaps := make([]swap.SwapBlob, 0, len(store.receipts))
	for _, receipt := range store.receipts {
		if receipt.Status == swap.Initiated {
			swaps = append(swaps, store.swaps[receipt.ID])
		}
	}

	return swaps, nil
}

func (store *MockStorage) UpdateStatus(update swap.StatusUpdate) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	receipt, ok := store.receipts[update.ID]
	if !ok {
		return errors.New("swap not found")
	}
	receipt.Status = update.Code
	store.receipts[update.ID] = receipt
	return nil
}

func (store *MockStorage) Swaps() ([]swap.SwapReceipt, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	receipts := make([]swap.SwapReceipt, 0, len(store.receipts))
	for _, receipt := range store.receipts {
		receipts = append(receipts, receipt)
	}

	return receipts, nil
}
