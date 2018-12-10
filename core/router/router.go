package router

import (
	"time"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/sirupsen/logrus"
)

type Storage interface {
	InsertSwap(swap swap.SwapRequest) error
	PendingSwap(swap.SwapID) (swap.SwapRequest, error)
	DeletePendingSwap(swap.SwapID) error
	PendingSwaps() ([]swap.SwapRequest, error)
	UpdateStatus(update swap.StatusUpdate) error
	Swaps() ([]swap.SwapReceipt, error)
}

type Router interface {
	Run(done <-chan struct{}, swapRequests <-chan swap.SwapRequest, updateRequests <-chan swap.StatusUpdate, swaps chan<- swap.SwapRequest, updates chan<- swap.StatusUpdate, statuses chan<- swap.SwapReceipt)
}

type router struct {
	storage Storage
	logger  logrus.FieldLogger
}

func New(storage Storage, logger logrus.FieldLogger) Router {
	return &router{storage, logger}
}

func (router *router) Run(done <-chan struct{}, swapRequests <-chan swap.SwapRequest, updateRequests <-chan swap.StatusUpdate, swaps chan<- swap.SwapRequest, updates chan<- swap.StatusUpdate, statuses chan<- swap.SwapReceipt) {
	results := make(chan swap.SwapResult)

	// Loading swap statuses from storage
	SwapReceiptes, err := router.storage.Swaps()
	if err != nil {
		return
	}
	go co.ParForAll(SwapReceiptes, func(i int) {
		statuses <- SwapReceiptes[i]
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
		case swapReq := <-swapRequests:
			logger := router.logger.WithField("SwapID", swapReq.ID)
			router.storage.InsertSwap(swapReq)
			statuses <- swap.NewSwapReceipt(swapReq.SwapBlob)
			logger.Info("adding to the swap queue")
			swaps <- swapReq
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
