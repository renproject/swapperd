package core

import (
	"encoding/base64"
	"fmt"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/core/delayed"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/republicprotocol/tau"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	LoadCosts(id swap.SwapID) (blockchain.Cost, blockchain.Cost)
	DeletePendingSwap(swap.SwapID) error
	Receipts() ([]swap.SwapReceipt, error)
	PutReceipt(receipt swap.SwapReceipt) error
	UpdateReceipt(receiptUpdate swap.ReceiptUpdate) error
	PutSwap(blob swap.SwapBlob) error
	PendingSwaps() ([]swap.SwapBlob, error)
}

type core struct {
	delayedSwapper tau.Task
	swapper        tau.Task
	status         tau.Task
	storage        Storage
}

func New(cap int, storage Storage, builder swapper.ContractBuilder, callback delayed.DelayCallback) tau.Task {
	delayedSwapperTask := delayed.New(cap, callback)
	swapperTask := swapper.New(cap, builder)
	statusTask := status.New(cap)
	return tau.New(tau.NewIO(cap), &core{delayedSwapperTask, swapperTask, statusTask, storage}, delayedSwapperTask, swapperTask, statusTask)
}

func (core *core) Reduce(msg tau.Message) tau.Message {
	switch msg := msg.(type) {
	case Bootload:
		return core.handleBootload(msg)
	case SwapRequest:
		return core.handleSwapRequest(msg)
	case swapper.ReceiptUpdate:
		return core.handleUpdateRequest(swap.ReceiptUpdate(msg))
	case swapper.DeleteSwap:
		return core.handleDeleteSwap(msg.ID)
	case delayed.SwapRequest:
		return core.handleSwapRequest(SwapRequest(msg))
	case delayed.ReceiptUpdate:
		return core.handleUpdateRequest(swap.ReceiptUpdate(msg))
	case delayed.DeleteSwap:
		return core.handleDeleteSwap(msg.ID)
	case status.ReceiptQuery:
		core.status.Send(msg)
		return nil
	case tau.Error:
		return msg
	case tau.Tick:
		core.status.Send(msg)
		core.swapper.Send(msg)
		core.delayedSwapper.Send(msg)
		return nil
	default:
		return tau.NewError(fmt.Errorf("invalid message type in core: %T", msg))
	}
}

func (core *core) handleUpdateRequest(update swap.ReceiptUpdate) tau.Message {
	core.status.Send(status.ReceiptUpdate(update))
	if err := core.storage.UpdateReceipt(swap.ReceiptUpdate(update)); err != nil {
		return tau.NewError(err)
	}
	return nil
}

func (core *core) handleSwapRequest(msg SwapRequest) tau.Message {
	if err := core.storage.PutSwap(swap.SwapBlob(msg)); err != nil {
		return tau.NewError(err)
	}

	receipt := swap.NewSwapReceipt(swap.SwapBlob(msg))
	core.status.Send(status.Receipt(receipt))
	if err := core.storage.PutReceipt(receipt); err != nil {
		return tau.NewError(err)
	}

	if msg.Delay {
		core.delayedSwapper.Send(delayed.SwapRequest(msg))
		return nil
	}

	sendCost, receiveCost := core.storage.LoadCosts(msg.ID)
	core.swapper.Send(swapper.NewSwapRequest(swap.SwapBlob(msg), sendCost, receiveCost))
	return nil
}

func (core *core) handleBootload(msg Bootload) tau.Message {
	return tau.NewMessageBatch([]tau.Message{core.handleSwapperBootload(msg), core.handleStatusBootload(msg)})
}

func (core *core) handleStatusBootload(msg Bootload) tau.Message {
	// Loading historical swap receipts
	historicalReceipts, err := core.storage.Receipts()
	if err != nil {
		return tau.NewError(err)
	}

	co.ParForAll(historicalReceipts, func(i int) {
		core.status.Send(status.Receipt(historicalReceipts[i]))
	})

	return nil
}

func (core *core) handleSwapperBootload(msg Bootload) tau.Message {
	pendingSwaps, err := core.storage.PendingSwaps()
	if err != nil {
		return tau.NewError(err)
	}

	co.ParForAll(pendingSwaps, func(i int) {
		pendingSwap := pendingSwaps[i]

		hash, err := base64.StdEncoding.DecodeString(pendingSwap.PasswordHash)
		if pendingSwap.PasswordHash != "" && err != nil {
			return
		}

		if pendingSwap.PasswordHash != "" && bcrypt.CompareHashAndPassword(hash, []byte(msg.Password)) != nil {
			return
		}

		core.status.Send(status.ReceiptUpdate(swap.NewReceiptUpdate(pendingSwap.ID, func(receipt *swap.SwapReceipt) {
			receipt.Active = true
		})))

		pendingSwap.Password = msg.Password
		if pendingSwap.Delay {
			core.delayedSwapper.Send(delayed.DelayedSwapRequest(pendingSwap))
			return
		}

		sendCost, receiveCost := core.storage.LoadCosts(pendingSwap.ID)
		core.swapper.Send(swapper.NewSwapRequest(pendingSwap, sendCost, receiveCost))
		return
	})
	return nil
}

func (core *core) handleDeleteSwap(id swap.SwapID) tau.Message {
	if err := core.storage.DeletePendingSwap(id); err != nil {
		return tau.NewError(err)
	}
	return nil
}
