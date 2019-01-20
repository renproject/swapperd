package swapper

import (
	"encoding/base64"
	"fmt"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/core/swapper/delayed"
	"github.com/republicprotocol/swapperd/core/swapper/immediate"
	"github.com/republicprotocol/swapperd/core/swapper/status"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/republicprotocol/tau"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	LoadCosts(id swap.SwapID) (blockchain.Cost, blockchain.Cost)
	Receipts() ([]swap.SwapReceipt, error)
	PutReceipt(receipt swap.SwapReceipt) error

	PutSwap(blob swap.SwapBlob) error
	DeletePendingSwap(swap.SwapID) error
	PendingSwaps() ([]swap.SwapBlob, error)
	UpdateReceipt(receiptUpdate swap.ReceiptUpdate) error
}

type core struct {
	delayedSwapper   tau.Task
	immediateSwapper tau.Task
	status           tau.Task
	storage          Storage
}

func New(cap int, storage Storage, delayedSwapperTask, immediateSwapperTask, statusTask tau.Task) tau.Task {
	return tau.New(tau.NewIO(cap), &core{delayedSwapperTask, immediateSwapperTask, statusTask, storage}, delayedSwapperTask, immediateSwapperTask, statusTask)
}

func (core *core) Reduce(msg tau.Message) tau.Message {
	switch msg := msg.(type) {
	case Bootload:
		return core.handleBootload(msg)
	case SwapRequest:
		return core.handleSwapRequest(msg)
	case immediate.ReceiptUpdate:
		return core.handleReceiptUpdate(swap.ReceiptUpdate(msg))
	case immediate.DeleteSwap:
		return core.handleDeleteSwap(msg.ID)
	case delayed.SwapRequest:
		return core.handleSwapRequest(SwapRequest(msg))
	case delayed.ReceiptUpdate:
		return core.handleReceiptUpdate(swap.ReceiptUpdate(msg))
	case delayed.DeleteSwap:
		return core.handleDeleteSwap(msg.ID)
	case status.ReceiptQuery:
		return core.handleReceiptQuery(msg)
	case tau.Error:
		return msg
	case tau.Tick:
		return core.handleTick(msg)
	default:
		return tau.NewError(fmt.Errorf("invalid message type in core: %T", msg))
	}
}

func (core *core) handleReceiptQuery(msg tau.Message) tau.Message {
	core.status.Send(msg)
	return nil
}

func (core *core) handleTick(msg tau.Message) tau.Message {
	core.status.Send(msg)
	core.immediateSwapper.Send(msg)
	core.delayedSwapper.Send(msg)
	return nil
}

func (core *core) handleReceiptUpdate(update swap.ReceiptUpdate) tau.Message {
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
		core.delayedSwapper.Send(delayed.DelayedSwapRequest(msg))
		return nil
	}

	sendCost, receiveCost := core.storage.LoadCosts(msg.ID)
	core.immediateSwapper.Send(immediate.NewSwapRequest(swap.SwapBlob(msg), sendCost, receiveCost))
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

	for _, pendingSwap := range pendingSwaps {
		hash, err := base64.StdEncoding.DecodeString(pendingSwap.PasswordHash)
		if pendingSwap.PasswordHash != "" && err != nil {
			continue
		}

		if pendingSwap.PasswordHash != "" && bcrypt.CompareHashAndPassword(hash, []byte(msg.Password)) != nil {
			continue
		}

		core.status.Send(status.ReceiptUpdate(swap.NewReceiptUpdate(pendingSwap.ID, func(receipt *swap.SwapReceipt) {
			receipt.Active = true
		})))

		pendingSwap.Password = msg.Password
		if pendingSwap.Delay {
			core.delayedSwapper.Send(delayed.DelayedSwapRequest(pendingSwap))
			continue
		}

		sendCost, receiveCost := core.storage.LoadCosts(pendingSwap.ID)
		core.immediateSwapper.Send(immediate.NewSwapRequest(pendingSwap, sendCost, receiveCost))
	}

	return nil
}

func (core *core) handleDeleteSwap(id swap.SwapID) tau.Message {
	if err := core.storage.DeletePendingSwap(id); err != nil {
		return tau.NewError(err)
	}
	return nil
}

type SwapRequest swap.SwapBlob

func (msg SwapRequest) IsMessage() {
}

type Bootload struct {
	Password string
}

func (msg Bootload) IsMessage() {
}
