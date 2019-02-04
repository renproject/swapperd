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

func (store *MockStorage) PutSwap(swapBlob swap.SwapBlob) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	swapBlob.Password = ""
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
	store.mu.RLock()
	defer store.mu.RUnlock()

	swaps := make([]swap.SwapBlob, 0, len(store.swaps))
	for _, swap := range store.swaps {
		swaps = append(swaps, store.swaps[swap.ID])
	}

	return swaps, nil
}

func (store *MockStorage) Receipts() ([]swap.SwapReceipt, error) {
	return []swap.SwapReceipt{}, nil
}

func (store *MockStorage) PutReceipt(receipt swap.SwapReceipt) error {
	return nil
}

func (store *MockStorage) UpdateReceipt(update swap.ReceiptUpdate) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	receipt, ok := store.receipts[update.ID]
	if !ok {
		return errors.New("swap not found")
	}

	update.Update(&receipt)
	store.receipts[update.ID] = receipt
	return nil
}

func (store *MockStorage) LoadCosts(id swap.SwapID) (blockchain.Cost, blockchain.Cost) {
	store.mu.Lock()
	defer store.mu.Unlock()
	receipt, ok := store.receipts[id]
	if !ok {
		return blockchain.Cost{}, blockchain.Cost{}
	}
	return blockchain.CostBlobToCost(receipt.SendCost), blockchain.CostBlobToCost(receipt.ReceiveCost)
}
