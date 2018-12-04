package balance

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/republicprotocol/swapperd/foundation"
)

type Blockchain interface {
	Balances() (map[foundation.TokenName]foundation.Balance, error)
}

type BalanceQuery struct {
	Response chan<- map[foundation.TokenName]foundation.Balance
}

type Balances interface {
	Run(done <-chan struct{}, queries <-chan BalanceQuery)
}

type balances struct {
	mu              *sync.RWMutex
	updateFrequency time.Duration
	balanceMap      map[foundation.TokenName]foundation.Balance
	logger          logrus.FieldLogger
	blockchain      Blockchain
}

func New(updateFrequency time.Duration, blockchain Blockchain, logger logrus.FieldLogger) Balances {
	return &balances{new(sync.RWMutex), updateFrequency, map[foundation.TokenName]foundation.Balance{}, logger, blockchain}
}

func (balances *balances) Run(done <-chan struct{}, queries <-chan BalanceQuery) {
	ticker := time.NewTicker(balances.updateFrequency)
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			go balances.update()
		case query, ok := <-queries:
			if !ok {
				return
			}
			query.Response <- balances.read()
		}
	}
}

func (balances *balances) update() {
	balanceMap, err := balances.blockchain.Balances()
	if err != nil {
		balances.logger.Errorf("cannot update balances: %v", err)
		return
	}
	balances.mu.Lock()
	defer balances.mu.Unlock()
	balances.balanceMap = balanceMap
}

func (balances *balances) read() map[foundation.TokenName]foundation.Balance {
	balances.mu.RLock()
	defer balances.mu.RUnlock()
	return balances.balanceMap
}
