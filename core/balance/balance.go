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
	Response chan<- map[blockchain.TokenName]blockchain.Balance
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
	return &balances{
		mu:              new(sync.RWMutex),
		updateFrequency: updateFrequency,
		balanceMap:      map[blockchain.TokenName]blockchain.Balance{},
		logger:          logger,
		blockchain:      bc,
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

func (balances *balances) read() map[blockchain.TokenName]blockchain.Balance {
	balances.mu.RLock()
	defer balances.mu.RUnlock()

	return copyBalanceMap(balances.balanceMap)
}

func copyBalanceMap(balance map[blockchain.TokenName]blockchain.Balance) map[blockchain.TokenName]blockchain.Balance {
	copied := map[blockchain.TokenName]blockchain.Balance{}
	for i, j := range balance {
		copied[i] = j
	}
	return copied
}
