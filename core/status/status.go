package status

import (
	"fmt"
	"sync"

	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/republicprotocol/tau"
)

type statuses struct {
	mu       *sync.RWMutex
	statuses map[swap.SwapID]swap.SwapReceipt
}

func New(cap int) tau.Task {
	return tau.New(tau.NewIO(cap), &statuses{new(sync.RWMutex), map[swap.SwapID]swap.SwapReceipt{}})
}

func (statuses *statuses) Reduce(msg tau.Message) tau.Message {
	switch msg := msg.(type) {
	case Receipt:
		statuses.set(swap.SwapReceipt(msg))
		return nil
	case ReceiptUpdate:
		statuses.update(swap.ReceiptUpdate(msg))
		return nil
	case ReceiptQuery:
		go func() {
			msg.Responder <- statuses.get()
		}()
		return nil
	default:
		return tau.NewError(fmt.Errorf("invalid message type in transfers: %T", msg))
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

type Bootload struct {
}

func (msg Bootload) IsMessage() {
}

type Receipt swap.SwapReceipt

func (msg Receipt) IsMessage() {
}

type ReceiptUpdate swap.ReceiptUpdate

func (msg ReceiptUpdate) IsMessage() {
}

type ReceiptQuery struct {
	Responder chan<- map[swap.SwapID]swap.SwapReceipt
}

func (msg ReceiptQuery) IsMessage() {
}
