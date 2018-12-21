package status

import (
	"sync"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/sirupsen/logrus"
)

type ReceiptQuery struct {
	Responder chan<- map[swap.SwapID]swap.SwapReceipt
}

type Storage interface {
	Receipts() ([]swap.SwapReceipt, error)
	PutReceipt(receipt swap.SwapReceipt) error
	UpdateReceipt(id swap.SwapID, update func(receipt *swap.SwapReceipt)) error
}

type Statuses interface {
	Run(done <-chan struct{}, swaps <-chan swap.SwapReceipt, updates <-chan swap.ReceiptUpdate, queries <-chan ReceiptQuery)
}

type statuses struct {
	mu       *sync.RWMutex
	statuses map[swap.SwapID]swap.SwapReceipt
	storage  Storage
	logger   logrus.FieldLogger
}

func New(storage Storage, logger logrus.FieldLogger) Statuses {
	return &statuses{new(sync.RWMutex), map[swap.SwapID]swap.SwapReceipt{}, storage, logger}
}

func (statuses *statuses) Run(done <-chan struct{}, receipts <-chan swap.SwapReceipt, updates <-chan swap.ReceiptUpdate, queries <-chan ReceiptQuery) {
	// Loading historical swap receipts
	historicalReceipts, err := statuses.storage.Receipts()
	if err != nil {
		statuses.logger.Error(err)
	}
	co.ParForAll(historicalReceipts, func(i int) {
		statuses.set(historicalReceipts[i])
	})

	for {
		select {
		case <-done:
			return
		case receipt, ok := <-receipts:
			if !ok {
				return
			}
			statuses.set(receipt)
			go func() {
				if err := statuses.storage.PutReceipt(receipt); err != nil {
					statuses.logger.Error(err)
				}
			}()
		case update, ok := <-updates:
			if !ok {
				return
			}
			statuses.update(update)
			go func() {
				if err := statuses.storage.UpdateReceipt(update.ID, update.Update); err != nil {
					statuses.logger.Error(err)
				}
			}()
		case query, ok := <-queries:
			if !ok {
				return
			}
			go func() {
				query.Responder <- statuses.get()
			}()
		}
	}
}

func (statuses *statuses) get() map[swap.SwapID]swap.SwapReceipt {
	statuses.mu.RLock()
	defer statuses.mu.RUnlock()

	statusMap := make(map[swap.SwapID]swap.SwapReceipt, len(statuses.statuses))
	for id, status := range statuses.statuses {
		statusMap[id] = status
	}

	return statusMap
}

func (statuses *statuses) set(status swap.SwapReceipt) {
	statuses.mu.Lock()
	defer statuses.mu.Unlock()

	statuses.statuses[status.ID] = status
}

func (statuses *statuses) update(update swap.ReceiptUpdate) {
	statuses.mu.Lock()
	defer statuses.mu.Unlock()

	receipt := statuses.statuses[update.ID]
	update.Update(&receipt)
	statuses.statuses[update.ID] = receipt
}
