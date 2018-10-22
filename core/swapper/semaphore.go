package swapper

import (
	"sync"

	"github.com/republicprotocol/swapperd/foundation"
)

type Semaphore struct {
	statuses map[foundation.SwapID]bool
	mu       *sync.RWMutex
}

func NewSemaphore() Semaphore {
	return Semaphore{
		statuses: map[foundation.SwapID]bool{},
		mu:       new(sync.RWMutex),
	}
}

func (sem *Semaphore) TryWait(id foundation.SwapID) bool {
	sem.mu.Lock()
	defer sem.mu.Unlock()
	if sem.statuses[id] {
		return false
	}
	sem.statuses[id] = true
	return true
}

func (sem *Semaphore) Signal(id foundation.SwapID) {
	sem.mu.Lock()
	defer sem.mu.Unlock()
	sem.statuses[id] = false
}
