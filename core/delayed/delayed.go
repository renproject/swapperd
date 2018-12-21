package delayed

import (
	"fmt"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/sirupsen/logrus"
	"time"
)

// ErrSwapDetailsUnavailable is returned when the swap details is not available
var ErrSwapDetailsUnavailable = fmt.Errorf("swap details unavailable")

// ErrSwapCancelled is returned when the swap has been cancelled
var ErrSwapCancelled = fmt.Errorf("swap cancelled")

// Callback handles the delayed atomic swap.
type Callback interface {

	// Run reads the delayed atomic swap from a read-only channel and keeps
	// polling for details to do the swap. When the swap details are finalized,
	// it write the full swap to a write-only channel and notify the status
	// change of the swap.
	Run(done <-chan struct{}, delayedSwaps <-chan swap.SwapBlob, swaps chan<- swap.SwapBlob, updates chan<- swap.ReceiptUpdate)
}

type Storage interface {
	UpdateReceipt(update swap.ReceiptUpdate) error
	DeletePendingSwap(swap.SwapID) error
}

// DelayCallback tries to fill the given delayed SwapBlob.
type DelayCallback interface {

	// DelayCallback will return the SwapBlob with full details if it's
	// available. It returns ErrSwapCancelled if the order has been canceled.
	// It returns ErrSwapDetailsUnavailable if the swap details are not available.
	DelayCallback(swap.SwapBlob) (swap.SwapBlob, error)
}

type callback struct {
	delayCallback DelayCallback
	storage       Storage
	logger        logrus.FieldLogger
}

// New returns a new callback with given components.
func New(delayCallback DelayCallback, storage Storage, logger logrus.FieldLogger) Callback {
	return &callback{delayCallback, storage, logger}
}

// Run implements the `Callback` interface.
func (callback *callback) Run(done <-chan struct{}, delayedSwaps <-chan swap.SwapBlob, swaps chan<- swap.SwapBlob, updates chan<- swap.ReceiptUpdate) {
	for {
		select {
		case <-done:
			return
		case blob, ok := <-delayedSwaps:
			if !ok {
				return
			}
			go callback.fill(blob, done, swaps, updates)
		}
	}
}

func (callback *callback) fill(blob swap.SwapBlob, done <-chan struct{}, swaps chan<- swap.SwapBlob, updates chan<- swap.ReceiptUpdate) {
	password := blob.Password
	blob.Password = ""
	for {
		filledBlob, err := callback.delayCallback.DelayCallback(blob)
		switch err {
		case nil:
			filledBlob.Password = password
			callback.handleUpdateSwap(filledBlob, done, updates)
			select {
			case <-done:
			case swaps <- filledBlob:
			}
			return
		case ErrSwapCancelled:
			callback.handleCancelSwap(blob.ID, done, updates)
			return
		case ErrSwapDetailsUnavailable:
		default:
			callback.logger.Error(err)
		}
		time.Sleep(30 * time.Second)
	}
}

func (callback *callback) handleCancelSwap(id swap.SwapID, done <-chan struct{}, updates chan<- swap.ReceiptUpdate) {
	callback.logger.Infof("cancelled delayed swap (%s)", id)
	update := swap.NewReceiptUpdate(id, func(receipt *swap.SwapReceipt) {
		receipt.ID = id
		receipt.Status = swap.Cancelled
	})
	select {
	case <-done:
		return
	case updates <- update:
	}
	if err := callback.storage.DeletePendingSwap(id); err != nil {
		callback.logger.Error(err)
	}
}

func (callback *callback) handleUpdateSwap(blob swap.SwapBlob, done <-chan struct{}, updates chan<- swap.ReceiptUpdate) {
	update := swap.NewReceiptUpdate(blob.ID, func(receipt *swap.SwapReceipt) {
		receipt.ReceiveAmount = blob.ReceiveAmount
		receipt.SendAmount = blob.SendAmount
	})
	select {
	case <-done:
	case updates <- update:
	}
}
