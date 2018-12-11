package status

import (
	"sync"

	"github.com/republicprotocol/swapperd/foundation/swap"
)

type monitor struct {
	mu       *sync.RWMutex
	statuses map[swap.SwapID]swap.SwapReceipt
}

func newMonitor() *monitor {
	return &monitor{
		mu:       new(sync.RWMutex),
		statuses: map[swap.SwapID]swap.SwapReceipt{},
	}
}

func (monitor *monitor) get() map[swap.SwapID]swap.SwapReceipt {
	monitor.mu.RLock()
	defer monitor.mu.RUnlock()
	statuses := make(map[swap.SwapID]swap.SwapReceipt, len(monitor.statuses))
	for id, status := range monitor.statuses {
		statuses[id] = status
	}
	return statuses
}

func (monitor *monitor) set(status swap.SwapReceipt) {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()
	monitor.statuses[status.ID] = status
}

func (monitor *monitor) update(status swap.StatusUpdate) {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()
	statusObj := monitor.statuses[status.ID]
	statusObj.Status = status.Code
	monitor.statuses[status.ID] = statusObj
}
