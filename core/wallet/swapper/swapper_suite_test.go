package swapper_test

import (
	"errors"
	"sync"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/renproject/swapperd/foundation/swap"
)

func TestSwapper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Swapper Suite")
}

type MockStorage struct {
	mu    *sync.RWMutex
	swaps map[swap.SwapID]swap.SwapBlob
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		mu:    new(sync.RWMutex),
		swaps: map[swap.SwapID]swap.SwapBlob{},
	}
}

func (store *MockStorage) PutSwap(swapBlob swap.SwapBlob) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	swapBlob.Password = ""
	store.swaps[swapBlob.ID] = swapBlob
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
	store.mu.RLock()
	defer store.mu.RUnlock()

	swaps := make([]swap.SwapBlob, 0, len(store.swaps))
	for _, swap := range store.swaps {
		swaps = append(swaps, store.swaps[swap.ID])
	}

	return swaps, nil
}

func (store *MockStorage) LoadCosts(id swap.SwapID) (blockchain.Cost, blockchain.Cost) {
	return blockchain.Cost{}, blockchain.Cost{}
}
