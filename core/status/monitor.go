package status

import (
	"sync"

	"github.com/republicprotocol/swapperd/foundation"
)

type monitor struct {
	mu       *sync.RWMutex
	statuses map[foundation.SwapID]foundation.SwapStatus
}

func newMonitor() *monitor {
	return &monitor{
		mu:       new(sync.RWMutex),
		statuses: map[foundation.SwapID]foundation.SwapStatus{},
	}
}

func (monitor *monitor) get() map[foundation.SwapID]foundation.SwapStatus {
	monitor.mu.RLock()
	defer monitor.mu.RUnlock()
	statuses := make(map[foundation.SwapID]foundation.SwapStatus, len(monitor.statuses))
	for swapID, status := range monitor.statuses {
		statuses[swapID] = status
	}
	return statuses
}

func (monitor *monitor) set(id foundation.SwapID, status foundation.SwapStatus) {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()
	monitor.statuses[id] = status
}
