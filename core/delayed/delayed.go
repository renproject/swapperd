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
	Run(done <-chan struct{}, delayedSwaps <-chan swap.SwapBlob, swaps chan<- swap.SwapBlob)
}

type Storage interface {
	UpdateReceipt(swapID swap.SwapID, update func(receipt *swap.SwapReceipt)) error
	DeletePendingSwap(swap.SwapID) error
}

type DelayCallback interface {
	DelayCallback(swap.SwapBlob) (swap.SwapBlob, error)
}

func New(delayCallback DelayCallback, storage Storage, logger logrus.FieldLogger) Callback {
	return &callback{delayCallback, storage, logger}
}

func (callback *callback) Run(done <-chan struct{}, delayedSwaps <-chan swap.SwapBlob, swaps chan<- swap.SwapBlob) {
	for {
		select {
		case <-done:
			return
		case blob, ok := <-delayedSwaps:
			if !ok {
				return
			}
			go callback.fill(blob, swaps)
		}
	}
}

func (callback *callback) fill(blob swap.SwapBlob, swaps chan<- swap.SwapBlob) {
	password := blob.Password
	blob.Password = ""
	for {
		filledBlob, err := callback.delayCallback.DelayCallback(blob)
		if err == nil {
			filledBlob.Password = password
			swaps <- filledBlob
			return
		}
		if err == ErrSwapCancelled {
			callback.handleRemoveSwap(blob.ID, swap.Cancelled)
			break
		}
		if time.Now().Unix() > blob.TimeLock {
			callback.handleRemoveSwap(blob.ID, swap.Expired)
			break
		}
		if err != ErrSwapDetailsUnavailable {
			callback.logger.Error(err)
		}
		time.Sleep(30 * time.Second)
	}
}

func (callback *callback) handleRemoveSwap(id swap.SwapID, status int) error {
	if err := callback.storage.UpdateReceipt(id, func(receipt *swap.SwapReceipt) {
		receipt.ID = id
	}); err != nil {
		callback.logger.Error(err)
		return err
	}
	if err := callback.storage.DeletePendingSwap(id); err != nil {
		callback.logger.Error(err)
		return err
	}
	return nil
}

func (callback *callback) handleUpdateSwap(blob swap.SwapBlob) error {
	if err := callback.storage.UpdateReceipt(blob.ID, func(receipt *swap.SwapReceipt) {
		receipt.ReceiveAmount = blob.ReceiveAmount
		receipt.SendAmount = blob.SendAmount
	}); err != nil {
		callback.logger.Error(err)
		return err
	}
	if err := callback.storage.DeletePendingSwap(blob.ID); err != nil {
		callback.logger.Error(err)
		return err
	}
	return nil
}
