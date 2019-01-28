package wallet

import (
	"fmt"

	"github.com/republicprotocol/swapperd/core/wallet/swapper"
	"github.com/republicprotocol/swapperd/core/wallet/swapper/delayed"
	"github.com/republicprotocol/swapperd/core/wallet/swapper/immediate"
	"github.com/republicprotocol/swapperd/core/wallet/transfer"

	"github.com/republicprotocol/tau"
)

type Storage interface {
	transfer.Storage
	swapper.Storage
}

type wallet struct {
	swapperTask  tau.Task
	transferTask tau.Task
}

func New(cap int, storage Storage, bc transfer.Blockchain, builder immediate.ContractBuilder, callback delayed.DelayCallback) tau.Task {
	swapperTask := swapper.New(cap, storage, builder, callback)
	transferTask := transfer.New(cap, bc, storage)
	return tau.New(tau.NewIO(cap), &wallet{swapperTask, transferTask}, swapperTask, transferTask)
}

func (wallet *wallet) Reduce(msg tau.Message) tau.Message {
	switch msg := msg.(type) {
	case tau.Tick:
		wallet.swapperTask.Send(msg)
		wallet.transferTask.Send(msg)
	case SwapperRequest:
		wallet.swapperTask.Send(msg.Message)
	case TransferRequest:
		wallet.transferTask.Send(msg.Message)
	case tau.Error:
		return tau.MessageBatch{msg, AcceptRequest{}}
	default:
		return tau.MessageBatch{tau.NewError(fmt.Errorf("Unknown message type: %T in wallet task", msg)), AcceptRequest{}}
	}
	return AcceptRequest{}
}

type AcceptRequest struct {
}

func (AcceptRequest) IsMessage() {
}

type SwapperRequest struct {
	Message tau.Message
}

func (SwapperRequest) IsMessage() {
}

func NewSwapperRequest(msg tau.Message) SwapperRequest {
	return SwapperRequest{
		Message: msg,
	}
}

type TransferRequest struct {
	Message tau.Message
}

func (TransferRequest) IsMessage() {

}

func NewTransferRequest(msg tau.Message) TransferRequest {
	return TransferRequest{
		Message: msg,
	}
}
