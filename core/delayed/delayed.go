package delayed

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/republicprotocol/swapperd/foundation/swap"
)

var ErrSwapDetailsUnavailable = fmt.Errorf("swap details unavailable")
var ErrSwapCancelled = fmt.Errorf("swap cancelled")

type callback struct {
	delayCallback DelayCallback
	storage       Storage
	logger        logrus.FieldLogger
}

type Callback interface {
	Run(done <-chan struct{}, delayedSwaps <-chan swap.SwapBlob, swaps chan<- swap.SwapBlob, updates chan<- swap.ReceiptUpdate)
}

type Storage interface {
	UpdateReceipt(update swap.ReceiptUpdate) error
	DeletePendingSwap(swap.SwapID) error
}

type DelayCallback interface {
	DelayCallback(swap.SwapBlob) (swap.SwapBlob, error)
}

func New(delayCallback DelayCallback, storage Storage, logger logrus.FieldLogger) Callback {
	return &callback{delayCallback, storage, logger}
}

func (callback *callback) Run(done <-chan struct{}, delayedSwaps <-chan swap.SwapBlob, swaps chan<- swap.SwapBlob, updates chan<- swap.ReceiptUpdate) {
	for {
		select {
		case <-done:
			return
		case blob, ok := <-delayedSwaps:
			if !ok {
				return
			}
			go callback.fill(blob, swaps, updates)
		}
	}
}

func (callback *callback) fill(blob swap.SwapBlob, swaps chan<- swap.SwapBlob, updates chan<- swap.ReceiptUpdate) {
	password := blob.Password
	blob.Password = ""
	for {
		filledBlob, err := callback.delayCallback.DelayCallback(blob)
		if err == nil {
			filledBlob.Password = password
			callback.handleUpdateSwap(filledBlob, updates)
			swaps <- filledBlob
			return
		}
		if err == ErrSwapCancelled {
			callback.handleCancelSwap(blob.ID, updates)
			break
		}
		if err != ErrSwapDetailsUnavailable {
			callback.logger.Error(err)
		}
		time.Sleep(30 * time.Second)
	}
}

func (callback *callback) handleCancelSwap(id swap.SwapID, updates chan<- swap.ReceiptUpdate) {
	callback.logger.Infof("cancelled delayed swap (%s)", id)
	update := swap.NewReceiptUpdate(id, func(receipt *swap.SwapReceipt) {
		receipt.ID = id
		receipt.Status = swap.Cancelled
	})
	updates <- update
	if err := callback.storage.DeletePendingSwap(id); err != nil {
		callback.logger.Error(err)
	}
}

func (callback *callback) handleUpdateSwap(blob swap.SwapBlob, updates chan<- swap.ReceiptUpdate) {
	update := swap.NewReceiptUpdate(blob.ID, func(receipt *swap.SwapReceipt) {
		receipt.ReceiveAmount = blob.ReceiveAmount
		receipt.SendAmount = blob.SendAmount
		receipt.TimeLock = blob.TimeLock
	})
	updates <- update
}
