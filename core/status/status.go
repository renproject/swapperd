package status

import (
	"sync"

	"github.com/republicprotocol/swapperd/foundation"
)

type Query struct {
	Responder chan<- map[foundation.SwapID]foundation.Status
}

type monitor struct {
	mu       *sync.RWMutex
	statuses map[foundation.SwapID]foundation.Status
}

type StatusBook interface {
	Run(statuses <-chan foundation.SwapStatus, queries <-chan Query, done <-chan struct{})
}

func New() StatusBook {
	return &monitor{
		new(sync.RWMutex),
		make(map[foundation.SwapID]foundation.Status),
	}
}

func (monitor *monitor) Run(statuses <-chan foundation.SwapStatus, queries <-chan Query, done <-chan struct{}) {
	for {
		select {
		case <-done:
			return
		case status, ok := <-statuses:
			if !ok {
				return
			}
			monitor.write(status.ID, status.Status)
		case query, ok := <-queries:
			if !ok {
				return
			}
			query.Responder <- monitor.read()
		}
	}
}

func (monitor *monitor) read() map[foundation.SwapID]foundation.Status {
	monitor.mu.RLock()
	defer monitor.mu.RUnlock()
	statuses := make(map[foundation.SwapID]foundation.Status, len(monitor.statuses))
	for swapID, status := range monitor.statuses {
		statuses[swapID] = status
	}
	return statuses
}

func (monitor *monitor) write(id foundation.SwapID, status foundation.Status) {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()
	monitor.statuses[id] = status
}
