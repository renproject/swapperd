package status

import (
	"sync"

	"github.com/republicprotocol/swapperd/foundation/swap"
)

type ReceiptQuery struct {
	Responder chan<- map[swap.SwapID]swap.SwapReceipt
}
type Statuses interface {
	Run(done <-chan struct{}, swaps <-chan swap.SwapReceipt, updates <-chan swap.StatusUpdate, queries <-chan ReceiptQuery)
}

type statuses struct {
	mu       *sync.RWMutex
	statuses map[swap.SwapID]swap.SwapReceipt
}

func New() Statuses {
	return &statuses{new(sync.RWMutex), map[swap.SwapID]swap.SwapReceipt{}}
}

func (statuses *statuses) Run(done <-chan struct{}, receipts <-chan swap.SwapReceipt, updates <-chan swap.StatusUpdate, queries <-chan ReceiptQuery) {
	for {
		select {
		case <-done:
			return
		case receipt, ok := <-receipts:
			if !ok {
				return
			}
			statuses.set(receipt)
		case update, ok := <-updates:
			if !ok {
				return
			}
			statuses.update(update)
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

func (statuses *statuses) update(status swap.StatusUpdate) {
	statuses.mu.Lock()
	defer statuses.mu.Unlock()
	statusObj := statuses.statuses[status.ID]
	statusObj.Status = status.Code
	statuses.statuses[status.ID] = statusObj
}
