package transfer

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/tau"
	"github.com/sirupsen/logrus"
)

type name struct {
}

type Blockchain interface {
	Transfer(password string, token blockchain.Token, to string, amount *big.Int) (string, error)
	Confirmations(txHash string) (int64, error)
}

type transfers struct {
	mu          *sync.RWMutex
	transferMap TransferReceiptMap
	logger      logrus.FieldLogger
	blockchain  Blockchain
}

func New(cap int, bc Blockchain, logger logrus.FieldLogger) tau.Task {
	return tau.New(tau.NewIO(cap), &transfers{new(sync.RWMutex), TransferReceiptMap{}, logger, bc})
}

func (transfers *transfers) Reduce(msg tau.Message) tau.Message {
	switch msg := msg.(type) {
	case tau.Tick:
		go transfers.update()
		return nil
	case TransferRequest:
		go func() {
			txHash, err := transfers.blockchain.Transfer(msg.Password, msg.Token, msg.To, msg.Amount)
			msg.Responder <- txHash
		}()
		return nil
	default:
		return tau.NewError(fmt.Errorf("invalid message type in transfers: %T", msg))
	}
}

func (transfers *transfers) update() {
	balanceMap, err := transfers.blockchain.Balances()
	if err != nil {
		transfers.logger.Errorf("cannot update transfers: %v", err)
		return
	}
	transfers.mu.Lock()
	defer transfers.mu.Unlock()
	transfers.balanceMap = balanceMap
}

func (transfers *transfers) read() BalanceMap {
	transfers.mu.RLock()
	defer transfers.mu.RUnlock()
	return transfers.balanceMap
}

type TransferRequest struct {
	Password string
	Token    blockchain.Token
	To       string
	Amount   *big.Int

	Responder chan<- string
}

func (request TransferRequest) IsMessage() {
}

type TransferReceiptMap map[string]TransferReceipt

func (request TransferReceiptMap) IsMessage() {
}

type TransferReceipt struct {
	Confirmations int64
	To            string
	From          string
	Token         blockchain.Token
	Value         *big.Int
	Fee           *big.Int
	Timestamp     int64
}
