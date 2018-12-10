package balance

import (
	"sync"
	"time"

	"github.com/republicprotocol/swapperd/foundation"

	"github.com/sirupsen/logrus"
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
	return &balances{
		mu:              new(sync.RWMutex),
		updateFrequency: updateFrequency,
		balanceMap:      map[foundation.TokenName]foundation.Balance{},
		logger:          logger,
		blockchain:      blockchain,
	}
}

func (balances *balances) Run(done <-chan struct{}, queries <-chan BalanceQuery) {
	ticker := time.NewTicker(balances.updateFrequency)
	defer ticker.Stop()
	go balances.update()

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
			go func() {
				select {
				case query.Response <- balances.read():
				case <-done:
				}
			}()
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

	balances.balanceMap = copyBalanceMap(balanceMap)
}

func (balances *balances) read() map[foundation.TokenName]foundation.Balance {
	balances.mu.RLock()
	defer balances.mu.RUnlock()

	return copyBalanceMap(balances.balanceMap)
}

func copyBalanceMap(balance map[foundation.TokenName]foundation.Balance) map[foundation.TokenName]foundation.Balance {
	copied := map[foundation.TokenName]foundation.Balance{}
	for i, j := range balance {
		copied[i] = j
	}
	return copied
}
