package bootload

import (
	"github.com/republicprotocol/co-go"
	"github.com/sirupsen/logrus"

	"github.com/republicprotocol/swapperd/foundation/swap"
)

type Storage interface {
	PendingSwaps() ([]swap.SwapBlob, error)
	Swaps() ([]swap.SwapReceipt, error)
}

type Bootloader interface {
	Bootload(password string, receipts chan<- swap.SwapReceipt, swaps chan<- swap.SwapBlob)
}

type bootloader struct {
	storage Storage
	logger  logrus.FieldLogger
}

func New(storage Storage, logger logrus.FieldLogger) Bootloader {
	return &bootloader{
		storage: storage,
	}
}

func (bootloader *bootloader) Bootload(password string, receipts chan<- swap.SwapReceipt, swaps chan<- swap.SwapBlob) {
	bootloader.logger.Info("loading historical swap receipts")
	// Loading swap statuses from storage
	hostoricalReceipts, err := bootloader.storage.Swaps()
	if err != nil {
		bootloader.logger.Error(err)
		return
	}
	co.ParForAll(hostoricalReceipts, func(i int) {
		receipts <- hostoricalReceipts[i]
	})

	bootloader.logger.Info("loading pending atomic swaps")
	// Loading pending swaps from storage
	swapsToRetry, err := bootloader.storage.PendingSwaps()
	if err != nil {
		bootloader.logger.Error(err)
		return
	}
	co.ParForAll(swapsToRetry, func(i int) {
		blob := swapsToRetry[i]
		blob.Password = password
		swaps <- blob
	})
}
