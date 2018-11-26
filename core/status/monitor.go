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

	statuses := make(map[foundation.SwapID]foundation.SwapStatus)
	for id, status := range monitor.statuses {
		statuses[id] = status
	}
	return statuses
}

func (monitor *monitor) set(status foundation.SwapStatus) {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()

	monitor.statuses[status.ID] = status
}

func (monitor *monitor) update(status foundation.StatusUpdate) {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()
	
	statusObj := monitor.statuses[status.ID]
	statusObj.Status = status.Status
	monitor.statuses[status.ID] = statusObj
}
