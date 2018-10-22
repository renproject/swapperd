package swapper

import (
	"sync"

	"github.com/republicprotocol/swapperd/foundation"
)

type binarySemaphore struct {
	swaps map[foundation.SwapID]bool
	mu    *sync.RWMutex
}

func newBinarySemaphore() *binarySemaphore {
	return &binarySemaphore{
		swaps: map[foundation.SwapID]bool{},
		mu:    new(sync.RWMutex),
	}
}

func (sem *binarySemaphore) TryWait(id foundation.SwapID) bool {
	sem.mu.Lock()
	defer sem.mu.Unlock()

	if sem.swaps[id] {
		return false
	}
	sem.swaps[id] = true
	return true
}

func (sem *binarySemaphore) Signal(id foundation.SwapID) {
	sem.mu.Lock()
	defer sem.mu.Unlock()

	sem.swaps[id] = false
}
