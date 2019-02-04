package status_test

import (
	"errors"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/swapperd/foundation/swap"

	"testing"
)

func TestStatus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Status Suite")
}

type MockStorage struct {
	mu       *sync.RWMutex
	receipts map[swap.SwapID]swap.SwapReceipt
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		mu:       new(sync.RWMutex),
		receipts: map[swap.SwapID]swap.SwapReceipt{},
	}
}

func (store *MockStorage) Receipt(id swap.SwapID) (swap.SwapReceipt, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	return store.receipts[id], nil
}

func (store *MockStorage) PutReceipt(receipt swap.SwapReceipt) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.receipts[receipt.ID] = receipt
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
