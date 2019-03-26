package wallet

import (
	"fmt"

	"github.com/renproject/swapperd/core/wallet/status"
	"github.com/renproject/swapperd/core/wallet/swapper"
	"github.com/renproject/swapperd/core/wallet/swapper/delayed"
	"github.com/renproject/swapperd/core/wallet/swapper/immediate"
	"github.com/renproject/swapperd/core/wallet/transfer"
	"github.com/renproject/swapperd/foundation/swap"
	"github.com/republicprotocol/tau"
	"github.com/sirupsen/logrus"
)

type Storage interface {
	status.Storage
	transfer.Storage
	swapper.Storage
}

type wallet struct {
	swapStatusTask tau.Task
	swapperTask    tau.Task
	transferTask   tau.Task

	logger logrus.FieldLogger
}

type Wallet interface {
	transfer.Blockchain
	swapper.Wallet
}

func New(cap int, storage Storage, w Wallet, builder immediate.ContractBuilder, callback delayed.DelayCallback, logger logrus.FieldLogger) tau.Task {
	swapperTask := swapper.New(cap, storage, w, builder, callback, logger)
	swapStatusTask := status.New(cap, storage, logger)
	transferTask := transfer.New(cap, w, storage, logger)
	return tau.New(tau.NewIO(cap), &wallet{swapStatusTask, swapperTask, transferTask, logger}, swapStatusTask, swapperTask, transferTask)
}

func (wallet *wallet) Reduce(msg tau.Message) tau.Message {
	switch msg := msg.(type) {
	case Bootload:
		wallet.handleBootload(msg)
	case tau.Tick:
		wallet.handleTick(msg)
	case swapper.SwapRequest:
		wallet.handleSwapRequest(msg)
	case swapper.ReceiptUpdate:
		wallet.swapStatusTask.Send(status.ReceiptUpdate(msg))
	case transfer.TransferRequest:
		wallet.transferTask.Send(msg)
	case tau.Error:
		return msg
	default:
		return tau.NewError(fmt.Errorf("Unknown message type: %T in wallet task", msg))
	}
	return nil
}

func (wallet *wallet) handleSwapRequest(msg swapper.SwapRequest) {
	wallet.swapStatusTask.Send(status.NewReceipt(swap.SwapBlob(msg)))
	wallet.swapperTask.Send(msg)
}

func (wallet *wallet) handleBootload(msg Bootload) {
	wallet.swapperTask.Send(swapper.Bootload{Password: msg.Password})
}

func (wallet *wallet) handleTick(msg tau.Tick) {
	wallet.swapperTask.Send(msg)
}

type Bootload struct {
	Password string
}

func (Bootload) IsMessage() {
}
