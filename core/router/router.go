package router

import (
	"time"

	"github.com/republicprotocol/co-go"
	"github.com/sirupsen/logrus"

	"github.com/republicprotocol/swapperd/foundation"
)

type Storage interface {
	InsertSwap(swap foundation.SwapRequest) error

	PendingSwap(foundation.SwapID) (foundation.SwapRequest, error)

	DeletePendingSwap(foundation.SwapID) error

	PendingSwaps() ([]foundation.SwapRequest, error)

	UpdateStatus(update foundation.StatusUpdate) error

	Swaps() ([]foundation.SwapStatus, error)
}

type Router interface {
	Run(done <-chan struct{}, swapRequests <-chan foundation.SwapRequest, updateRequests <-chan foundation.StatusUpdate, swaps chan<- foundation.SwapRequest, updates chan<- foundation.StatusUpdate, statuses chan<- foundation.SwapStatus)
}

type router struct {
	storage Storage
	logger  logrus.FieldLogger
}

func New(storage Storage, logger logrus.FieldLogger) Router {
	return &router{storage, logger}
}

func (router *router) Run(done <-chan struct{}, swapRequests <-chan foundation.SwapRequest, updateRequests <-chan foundation.StatusUpdate, swaps chan<- foundation.SwapRequest, updates chan<- foundation.StatusUpdate, statuses chan<- foundation.SwapStatus) {
	results := make(chan foundation.SwapResult)

	// Loading swap statuses from storage
	SwapStatuses, err := router.storage.Swaps()
	if err != nil {
		return
	}
	go co.ParForAll(SwapStatuses, func(i int) {
		statuses <- SwapStatuses[i]
	})

	// Loading pending swaps from storage
	swapsToRetry, err := router.storage.PendingSwaps()
	if err != nil {
		return
	}
	go co.ParForAll(swapsToRetry, func(i int) {
		swaps <- swapsToRetry[i]
	})

	// fault tolerant middleware
	for {
		select {
		case <-done:
			return
		case swap := <-swapRequests:
			logger := router.logger.WithField("SwapID", swap.ID)
			router.storage.InsertSwap(swap)
			statuses <- foundation.NewSwapStatus(swap.SwapBlob)
			logger.Info("adding to the swap queue")
			swaps <- swap
		case update := <-updateRequests:
			logger := router.logger.WithField("SwapID", update.ID)
			if err := router.storage.UpdateStatus(update); err != nil {
				logger.Error(err)
				continue
			}
			updates <- update
		case result := <-results:
			logger := router.logger.WithField("SwapID", result.ID)
			if result.Success {
				if err := router.storage.DeletePendingSwap(result.ID); err != nil {
					logger.Error(err)
					continue
				}
				logger.Info("removed from pending swaps")
				continue
			}
			swap, err := router.storage.PendingSwap(result.ID)
			if err != nil {
				logger.Error(err)
				continue
			}
			go func() {
				time.Sleep(5 * time.Minute)
				swaps <- swap
			}()
		}
	}
}
