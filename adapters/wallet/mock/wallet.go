package mock

import (
	"math/big"
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

func (wallet *MockWallet) SetMatch(m match.Match) error {
	wallet.mu.Lock()
	defer wallet.mu.Unlock()
	wallet.matches[m.PersonalOrderID()] = match.NewMatch(m.PersonalOrderID(), m.ForeignOrderID(), m.SendValue(), m.ReceiveValue().Sub(m.ReceiveValue(), big.NewInt(10000)), m.SendCurrency(), m.ReceiveCurrency())
	wallet.matches[m.ForeignOrderID()] = match.NewMatch(m.ForeignOrderID(), m.PersonalOrderID(), m.ReceiveValue(), m.SendValue().Sub(m.SendValue(), big.NewInt(10000)), m.ReceiveCurrency(), m.SendCurrency())
	return nil
}

func (wallet *MockWallet) GetMatch(orderID [32]byte) (match.Match, error) {
	wallet.mu.RLock()
	defer wallet.mu.RUnlock()
	return wallet.matches[orderID], nil
}
