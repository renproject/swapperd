package delayed

import (
	"fmt"

	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/republicprotocol/tau"
)

var ErrSwapDetailsUnavailable = fmt.Errorf("swap details unavailable")
var ErrSwapCancelled = fmt.Errorf("swap cancelled")

type callback struct {
	delayCallback DelayCallback
	swapMap       map[swap.SwapID]DelayedSwapRequest
}

type DelayCallback interface {
	DelayCallback(swap.SwapBlob) (swap.SwapBlob, error)
}

func New(cap int, delayCallback DelayCallback) tau.Task {
	return tau.New(tau.NewIO(cap), &callback{delayCallback, map[swap.SwapID]DelayedSwapRequest{}})
}

func (callback *callback) Reduce(msg tau.Message) tau.Message {
	switch msg := msg.(type) {
	case DelayedSwapRequest:
		return callback.tryFill(msg)
	case tau.Tick:
		return callback.checkAgain()
	default:
		return tau.NewError(fmt.Errorf("invalid message type in delayed swapper: %T", msg))
	}
}

func (callback *callback) checkAgain() tau.Message {
	messages := []tau.Message{}
	for _, swap := range callback.swapMap {
		if msg := callback.tryFill(swap); msg != nil {
			messages = append(messages, msg)
		}
	}
	return tau.NewMessageBatch(messages)
}

func (callback *callback) tryFill(blob DelayedSwapRequest) tau.Message {
	password := blob.Password
	blob.Password = ""
	filledBlob, err := callback.delayCallback.DelayCallback(swap.SwapBlob(blob))
	if err == nil {
		filledBlob.Password = password
		return callback.handleUpdateSwap(SwapRequest(filledBlob))
	}
	if err == ErrSwapCancelled {
		return callback.handleCancelSwap(blob.ID)
	}
	blob.Password = password
	callback.swapMap[blob.ID] = blob
	if err != ErrSwapDetailsUnavailable {
		return tau.NewError(err)
	}
	return nil
}

func (callback *callback) handleCancelSwap(id swap.SwapID) tau.Message {
	update := ReceiptUpdate(swap.NewReceiptUpdate(id, func(receipt *swap.SwapReceipt) {
		receipt.ID = id
		receipt.Status = swap.Cancelled
	}))
	return tau.NewMessageBatch([]tau.Message{update, DeleteSwap{id}})
}

func (callback *callback) handleUpdateSwap(req SwapRequest) tau.Message {
	update := ReceiptUpdate(swap.NewReceiptUpdate(req.ID, func(receipt *swap.SwapReceipt) {
		receipt.ReceiveAmount = req.ReceiveAmount
		receipt.SendAmount = req.SendAmount
		receipt.TimeLock = req.TimeLock
	}))
	return tau.NewMessageBatch([]tau.Message{update, req})
}
