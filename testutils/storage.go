package testutils

import (
	"errors"
	"sync"

	"github.com/republicprotocol/swapperd/foundation"
)

type MockStorage struct {
	mu     *sync.RWMutex
	swaps  map[foundation.SwapID]foundation.SwapRequest
	status map[foundation.SwapID]foundation.SwapStatus
}

func NewMockStorage() MockStorage {
	return MockStorage{
		mu:    new(sync.RWMutex),
		swaps: map[foundation.SwapID]foundation.SwapRequest{},
	}
}

func (store *MockStorage) InsertSwap(swap foundation.SwapRequest) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.swaps[swap.ID] = swap
	store.status[swap.ID] = foundation.NewSwapStatus(swap.SwapBlob)
	return nil
}

func (store *MockStorage) PendingSwap(id foundation.SwapID) (foundation.SwapRequest, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	swap, ok := store.swaps[id]
	if !ok {
		return foundation.SwapRequest{}, errors.New("pending swap not found")
	}
	return swap, nil
}

func (store *MockStorage) DeletePendingSwap(id foundation.SwapID) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	delete(store.swaps, id)
	return nil
}

func (store *MockStorage) PendingSwaps() ([]foundation.SwapRequest, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	swaps := make([]foundation.SwapRequest, 0, len(store.status))
	for _, status := range store.status {
		if status.Status == foundation.Initiated {
			swaps = append(swaps, store.swaps[status.ID])
		}
	}

	return swaps, nil
}

func (store *MockStorage) UpdateStatus(update foundation.StatusUpdate) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	status, ok := store.status[update.ID]
	if !ok {
		return errors.New("swap not found")
	}
	status.Status = update.Status
	store.status[update.ID] = status
	return nil
}

func (store *MockStorage) Swaps() ([]foundation.SwapStatus, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	statuses := make([]foundation.SwapStatus, 0, len(store.status))
	for _, status := range store.status {
		statuses = append(statuses, status)
	}

	return statuses, nil
}
