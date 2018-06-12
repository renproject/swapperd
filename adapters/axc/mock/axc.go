package mock

import (
	"sync"

	"github.com/republicprotocol/atom-go/services/axc"
)

type MockAXC struct {
	mu     *sync.RWMutex
	owners map[[32]byte][]byte
}

func NewMockAXC() axc.AXC {
	return &MockAXC{
		mu:     &sync.RWMutex{},
		owners: make(map[[32]byte][]byte),
	}
}

func (axc *MockAXC) SetOwnerAddress(orderID [32]byte, address []byte) error {
	axc.mu.Lock()
	axc.owners[orderID] = address
	axc.mu.Unlock()
	return nil
}

func (axc *MockAXC) GetOwnerAddress(orderID [32]byte) ([]byte, error) {
	axc.mu.RLock()
	address := axc.owners[orderID]
	axc.mu.RUnlock()
	return address, nil
}
