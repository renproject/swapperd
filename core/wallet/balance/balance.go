package balance

import (
	"fmt"
	"sync"

	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/tau"
	"github.com/sirupsen/logrus"
)

type Blockchain interface {
	Balances() (BalanceMap, error)
}

type Balances interface {
}

type balances struct {
	mu         *sync.RWMutex
	balanceMap BalanceMap
	logger     logrus.FieldLogger
	blockchain Blockchain
}

func New(cap int, bc Blockchain, logger logrus.FieldLogger) tau.Task {
	return tau.New(tau.NewIO(cap), &balances{new(sync.RWMutex), BalanceMap{}, logger, bc})
}

func (balances *balances) Reduce(msg tau.Message) tau.Message {
	switch msg := msg.(type) {
	case tau.Tick:
		go balances.update()
		return nil
	case BalanceRequest:
		go func() {
			msg.Responder <- balances.read()
		}()
		return nil
	default:
		return tau.NewError(fmt.Errorf("invalid message type in balances: %T", msg))
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

func (balances *balances) read() BalanceMap {
	balances.mu.RLock()
	defer balances.mu.RUnlock()
	return balances.balanceMap
}

type BalanceRequest struct {
	Responder chan<- BalanceMap
}

func (request BalanceRequest) IsMessage() {
}

type BalanceMap map[blockchain.TokenName]blockchain.Balance

func (request BalanceMap) IsMessage() {
}
