package wallet

import (
	"fmt"
	"sync"

	"github.com/republicprotocol/swapperd/core/wallet/balance"
	"github.com/republicprotocol/swapperd/core/wallet/transfer"
	"github.com/republicprotocol/tau"
	"github.com/sirupsen/logrus"
)

type Blockchain interface {
	balance.Blockchain
	transfer.Blockchain
}

type Storage interface {
	transfer.Storage
}

type wallet struct {
	mu           *sync.RWMutex
	logger       logrus.FieldLogger
	transferTask tau.Task
	balanceTask  tau.Task
}

func New(cap int, bc Blockchain, storage Storage, logger logrus.FieldLogger) tau.Task {
	transferTask := transfer.New(cap, bc, storage, logger)
	balanceTask := balance.New(cap, bc, logger)
	return tau.New(tau.NewIO(cap), &wallet{new(sync.RWMutex), logger, transferTask, balanceTask}, transferTask, balanceTask)
}

func (wallet *wallet) Reduce(msg tau.Message) tau.Message {
	switch msg := msg.(type) {
	case Bootload:
		wallet.transferTask.Send(transfer.Bootload{})
		wallet.balanceTask.Send(balance.Bootload{})
	case transfer.TransferRequest, transfer.TransferReceiptRequest:
		wallet.transferTask.Send(msg)
	case balance.BalanceRequest:
		wallet.balanceTask.Send(msg)
	case tau.Tick:
		wallet.transferTask.Send(msg)
		wallet.balanceTask.Send(msg)
	default:
		return tau.NewError(fmt.Errorf("unsupported request"))
	}
	return nil
}

type Bootload struct {
}

func (request Bootload) IsMessage() {
}
