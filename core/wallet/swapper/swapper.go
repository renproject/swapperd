package swapper

import (
	"encoding/base64"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/renproject/swapperd/core/wallet/swapper/delayed"
	"github.com/renproject/swapperd/core/wallet/swapper/immediate"
	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/renproject/swapperd/foundation/swap"
	"github.com/renproject/tokens"
	"github.com/republicprotocol/tau"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	LoadCosts(id swap.SwapID) (blockchain.Cost, blockchain.Cost)
	PutSwap(blob swap.SwapBlob) error
	DeletePendingSwap(swap.SwapID) error
	PendingSwaps() ([]swap.SwapBlob, error)
}

type Wallet interface {
	LockBalance(token tokens.Name, value string) error
	immediate.Wallet
}

type swapper struct {
	delayedSwapper   tau.Task
	immediateSwapper tau.Task
	storage          Storage
	wallet           Wallet
	logger           logrus.FieldLogger
}

func New(cap int, storage Storage, wallet Wallet, builder immediate.ContractBuilder, callback delayed.DelayCallback, logger logrus.FieldLogger) tau.Task {
	delayedSwapperTask := delayed.New(cap, callback, logger)
	immediateSwapperTask := immediate.New(cap, builder, wallet, logger)
	return tau.New(tau.NewIO(cap), NewSwapper(delayedSwapperTask, immediateSwapperTask, wallet, storage, logger), delayedSwapperTask, immediateSwapperTask)
}

func NewSwapper(delayedSwapperTask, immediateSwapperTask tau.Task, wallet Wallet, storage Storage, logger logrus.FieldLogger) tau.Reducer {
	return &swapper{delayedSwapperTask, immediateSwapperTask, storage, wallet, logger}
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

	if err := swapper.wallet.LockBalance(msg.SendToken, msg.SendAmount); err != nil {
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

		pendingSwap.Password = msg.Password
		msgs = append(
			msgs,
			ReceiptUpdate(swap.NewReceiptUpdate(pendingSwap.ID, func(receipt *swap.SwapReceipt) {
				receipt.Active = true
			})),
			swapper.handleSwapRequest(SwapRequest(pendingSwap)),
		)
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
