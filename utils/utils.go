package utils

import "sync"

type SwapManager struct {
	statuses map[[32]byte]bool
	mu       *sync.RWMutex
}

func NewSwapManager() SwapManager {
	return SwapManager{
		statuses: map[[32]byte]bool{},
		mu:       new(sync.RWMutex),
	}
}

func (manager *SwapManager) Lock(id [32]byte) bool {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	if manager.statuses[id] {
		return false
	}
	manager.statuses[id] = true
	return true
}

func (manager *SwapManager) Unlock(id [32]byte) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.statuses[id] = false
}
