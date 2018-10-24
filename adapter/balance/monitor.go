package balance

import (
	"math/big"
	"sync"

	"github.com/republicprotocol/swapperd/foundation"
)

type monitor struct {
	mu       *sync.RWMutex
	balances map[foundation.Token]*big.Int
}

func newMonitor() *monitor {
	return &monitor{
		mu:       new(sync.RWMutex),
		balances: map[foundation.Token]*big.Int{},
	}
}

func (monitor *monitor) get() map[foundation.Token]*big.Int {
	monitor.mu.RLock()
	defer monitor.mu.RUnlock()
	balances := make(map[foundation.Token]*big.Int, len(monitor.balances))
	for token, balance := range monitor.balances {
		balances[token] = balance
	}
	return balances
}

func (monitor *monitor) set(token foundation.Token, balance *big.Int) {
	monitor.mu.Lock()
	defer monitor.mu.Unlock()
	monitor.balances[token] = balance
}
