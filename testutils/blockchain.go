package testutils

import (
	"errors"
	"sync"

	"github.com/republicprotocol/swapperd/foundation"
)

// MockBlockchain implements the `balance.Blockchain` interface.
type MockBlockchain struct {
	mu      *sync.Mutex
	balance map[foundation.TokenName]foundation.Balance
}

// NewMockBlockchain creates a new `MockBlockchain`.
func NewMockBlockchain(balance map[foundation.TokenName]foundation.Balance) *MockBlockchain {
	return &MockBlockchain{
		mu:      new(sync.Mutex),
		balance: copyBalanceMap(balance),
	}
}

// Balances implements the `balance.Blockchain` interface.
func (blockchain *MockBlockchain) Balances() (map[foundation.TokenName]foundation.Balance, error) {
	blockchain.mu.Lock()
	defer blockchain.mu.Unlock()

	return copyBalanceMap(blockchain.balance), nil
}

// UpdateBalance with given data.
func (blockchain *MockBlockchain) UpdateBalance(balance map[foundation.TokenName]foundation.Balance) {
	blockchain.mu.Lock()
	defer blockchain.mu.Unlock()

	blockchain.balance = copyBalanceMap(balance)
}

type FaultyBlockchain struct {
	balance map[foundation.TokenName]foundation.Balance
	counter int
}

func NewFaultyBlockchain(balance map[foundation.TokenName]foundation.Balance) *FaultyBlockchain {
	return &FaultyBlockchain{
		balance: balance,
		counter: 0,
	}
}

func (blockchain *FaultyBlockchain) Balances() (map[foundation.TokenName]foundation.Balance, error) {
	blockchain.counter++
	if blockchain.counter%2 != 0 {
		return blockchain.balance, nil
	}
	return nil, errors.New("cannot get the balance")
}

func copyBalanceMap(balance map[foundation.TokenName]foundation.Balance) map[foundation.TokenName]foundation.Balance {
	copied := map[foundation.TokenName]foundation.Balance{}
	for i, j := range balance {
		copied[i] = j
	}
	return copied
}
