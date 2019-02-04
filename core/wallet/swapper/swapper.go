package swapper

import (
	"encoding/base64"
	"fmt"

	"github.com/republicprotocol/swapperd/core/wallet/swapper/delayed"
	"github.com/republicprotocol/swapperd/core/wallet/swapper/immediate"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/republicprotocol/tau"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	LoadCosts(id swap.SwapID) (blockchain.Cost, blockchain.Cost)

	PutSwap(blob swap.SwapBlob) error
	DeletePendingSwap(swap.SwapID) error
	PendingSwaps() ([]swap.SwapBlob, error)
}

type swapper struct {
	delayedSwapper   tau.Task
	immediateSwapper tau.Task
	storage          Storage
}

func New(cap int, storage Storage, builder immediate.ContractBuilder, callback delayed.DelayCallback) tau.Task {
	delayedSwapperTask := delayed.New(cap, callback)
	immediateSwapperTask := immediate.New(cap, builder)
	return tau.New(tau.NewIO(cap), NewSwapper(delayedSwapperTask, immediateSwapperTask, storage), delayedSwapperTask, immediateSwapperTask)
}

func NewSwapper(delayedSwapperTask, immediateSwapperTask tau.Task, storage Storage) tau.Reducer {
	return &swapper{delayedSwapperTask, immediateSwapperTask, storage}
}

func (swapper *swapper) Reduce(msg tau.Message) tau.Message {
	switch msg := msg.(type) {
	case Bootload:
		return swapper.handleBootload(msg)
	case SwapRequest:
		return swapper.handleSwapRequest(msg)
	case immediate.ReceiptUpdate:
		return ReceiptUpdate(msg)
	case immediate.DeleteSwap:
		return swapper.handleDeleteSwap(msg.ID)
	case delayed.SwapRequest:
		return swapper.handleSwapRequest(SwapRequest(msg))
	case delayed.ReceiptUpdate:
		return ReceiptUpdate(msg)
	case delayed.DeleteSwap:
		return swapper.handleDeleteSwap(msg.ID)
	case tau.Error:
		return msg
	case tau.Tick:
		return swapper.handleTick(msg)
	default:
		return tau.NewError(fmt.Errorf("invalid message type in swapper: %T", msg))
	}
}

func (swapper *swapper) handleTick(msg tau.Message) tau.Message {
	swapper.immediateSwapper.Send(msg)
	swapper.delayedSwapper.Send(msg)
	return nil
}

func (swapper *swapper) handleSwapRequest(msg SwapRequest) tau.Message {
	if err := swapper.storage.PutSwap(swap.SwapBlob(msg)); err != nil {
		return tau.NewError(err)
	}

	if msg.Delay {
		swapper.delayedSwapper.Send(delayed.DelayedSwapRequest(msg))
		return nil
	}

	sendCost, receiveCost := swapper.storage.LoadCosts(msg.ID)
	swapper.immediateSwapper.Send(immediate.NewSwapRequest(swap.SwapBlob(msg), sendCost, receiveCost))
	return nil
}

func (swapper *swapper) handleBootload(msg Bootload) tau.Message {
	pendingSwaps, err := swapper.storage.PendingSwaps()
	if err != nil {
		return tau.NewError(err)
	}

	msgs := []tau.Message{}
	for _, pendingSwap := range pendingSwaps {
		hash, err := base64.StdEncoding.DecodeString(pendingSwap.PasswordHash)
		if pendingSwap.PasswordHash != "" && err != nil {
			continue
		}

		if pendingSwap.PasswordHash != "" && bcrypt.CompareHashAndPassword(hash, []byte(msg.Password)) != nil {
			continue
		}

		msgs = append(msgs, ReceiptUpdate(swap.NewReceiptUpdate(pendingSwap.ID, func(receipt *swap.SwapReceipt) {
			receipt.Active = true
		})))

		pendingSwap.Password = msg.Password
		if pendingSwap.Delay {
			swapper.delayedSwapper.Send(delayed.DelayedSwapRequest(pendingSwap))
			continue
		}

		sendCost, receiveCost := swapper.storage.LoadCosts(pendingSwap.ID)
		swapper.immediateSwapper.Send(immediate.NewSwapRequest(pendingSwap, sendCost, receiveCost))
	}

	return tau.NewMessageBatch(msgs)
}

func (swapper *swapper) handleDeleteSwap(id swap.SwapID) tau.Message {
	if err := swapper.storage.DeletePendingSwap(id); err != nil {
		return tau.NewError(err)
	}
	return nil
}

type SwapRequest swap.SwapBlob

func (SwapRequest) IsMessage() {
}

type Bootload struct {
	Password string
}

func (Bootload) IsMessage() {
}

type ReceiptUpdate swap.ReceiptUpdate

func (ReceiptUpdate) IsMessage() {
}
