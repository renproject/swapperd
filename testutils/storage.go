package testutils

import (
	"errors"
	"sync"

	"github.com/renproject/swapperd/core/wallet/transfer"
	"github.com/renproject/swapperd/foundation/swap"
)

type MockStorage struct {
	mu        *sync.RWMutex
	receipts  map[swap.SwapID]swap.SwapReceipt
	transfers map[string]transfer.TransferReceipt
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

func (store *MockStorage) Receipts() ([]swap.SwapReceipt, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	swapReceipts := []swap.SwapReceipt{}
	for _, receipt := range store.receipts {
		swapReceipts = append(swapReceipts, receipt)
	}

	return swapReceipts, nil
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

func (store *MockStorage) PutTransfer(receipt transfer.TransferReceipt) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.transfers[receipt.TxHash] = receipt
	return nil
}

func (store *MockStorage) Transfers() ([]transfer.TransferReceipt, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	transfers := []transfer.TransferReceipt{}
	for _, transfer := range store.transfers {
		transfers = append(transfers, transfer)
	}

	return transfers, nil
}
