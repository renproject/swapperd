package wallet

import (
	"fmt"

	"github.com/republicprotocol/swapperd/core/wallet/status"
	"github.com/republicprotocol/swapperd/core/wallet/swapper"
	"github.com/republicprotocol/swapperd/core/wallet/swapper/delayed"
	"github.com/republicprotocol/swapperd/core/wallet/swapper/immediate"
	"github.com/republicprotocol/swapperd/core/wallet/transfer"
	"github.com/republicprotocol/swapperd/foundation/swap"

	"github.com/republicprotocol/tau"
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
}

func New(cap int, storage Storage, bc transfer.Blockchain, builder immediate.ContractBuilder, callback delayed.DelayCallback) tau.Task {
	swapperTask := swapper.New(cap, storage, builder, callback)
	swapStatusTask := status.New(cap, storage)
	transferTask := transfer.New(cap, bc, storage)
	return tau.New(tau.NewIO(cap), &wallet{swapStatusTask, swapperTask, transferTask}, swapStatusTask, swapperTask, transferTask)
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
	wallet.swapperTask.Send(swapper.Bootload{msg.Password})
}

func (wallet *wallet) handleTick(msg tau.Tick) {
	wallet.swapperTask.Send(msg)
}

type Bootload struct {
	Password string
}

func (Bootload) IsMessage() {
}
