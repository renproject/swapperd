package mock

import (
	"sync"

	"github.com/republicprotocol/atom-go/domains/match"
	"github.com/republicprotocol/atom-go/services/watch"
)

type MockWallet struct {
	mu      *sync.RWMutex
	matches map[[32]byte]match.Match
}

func NewMockWallet() watch.Wallet {
	matches := make(map[[32]byte]match.Match)
	return &MockWallet{
		mu:      &sync.RWMutex{},
		matches: matches,
	}
}

func (wallet *MockWallet) SetMatch(orderID [32]byte, m match.Match) error {
	wallet.mu.Lock()
	defer wallet.mu.Unlock()
	wallet.matches[orderID] = m
	return nil
}

func (wallet *MockWallet) GetMatch(orderID [32]byte) (match.Match, error) {
	wallet.mu.RLock()
	defer wallet.mu.RUnlock()
	return wallet.matches[orderID], nil
}
