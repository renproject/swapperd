package testutils

import (
	"sync"

	"github.com/renproject/swapperd/foundation/swap"
)

type MockCallback struct {
	mu    *sync.RWMutex
	flags map[swap.SwapID]bool
	err   error
}

func NewMockCallback(err error) *MockCallback {
	return &MockCallback{
		mu:    new(sync.RWMutex),
		flags: map[swap.SwapID]bool{},
		err:   err,
	}
}

func (callback *MockCallback) DelayCallback(swap swap.SwapBlob) (swap.SwapBlob, error) {
	callback.mu.Lock()
	defer callback.mu.Unlock()
	if !callback.flags[swap.ID] {
		callback.flags[swap.ID] = true
		return swap, callback.err
	}
	return swap, nil
}
