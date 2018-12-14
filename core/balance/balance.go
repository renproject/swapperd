package balance

import (
	"sync"
	"time"

	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/sirupsen/logrus"
)

type Blockchain interface {
	Balances() (map[blockchain.TokenName]blockchain.Balance, error)
}

type BalanceQuery struct {
	Responder chan<- map[blockchain.TokenName]blockchain.Balance
}

type Balances interface {
	Run(done <-chan struct{}, queries <-chan BalanceQuery)
}

type balances struct {
	mu              *sync.RWMutex
	updateFrequency time.Duration
	balanceMap      map[blockchain.TokenName]blockchain.Balance
	logger          logrus.FieldLogger
	blockchain      Blockchain
}

func New(updateFrequency time.Duration, bc Blockchain, logger logrus.FieldLogger) Balances {
	return &balances{new(sync.RWMutex), updateFrequency, map[blockchain.TokenName]blockchain.Balance{}, logger, bc}
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
			go func() {
				query.Responder <- balances.read()
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
	balances.balanceMap = balanceMap
}

func (balances *balances) read() map[blockchain.TokenName]blockchain.Balance {
	balances.mu.RLock()
	defer balances.mu.RUnlock()
	return balances.balanceMap
}
