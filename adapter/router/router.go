package router

import (
	"time"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation"
	"github.com/sirupsen/logrus"
)

type Storage interface {
	InsertSwap(swap foundation.SwapRequest) error
	PendingSwap(foundation.SwapID) (foundation.SwapRequest, error)
	DeletePendingSwap(foundation.SwapID) error
	PendingSwaps() ([]foundation.SwapRequest, error)
	UpdateStatus(update foundation.StatusUpdate) error
	Swaps() ([]foundation.SwapStatus, error)
}

type Server interface {
	Run(done <-chan struct{}, swapRequests chan<- foundation.SwapRequest, statusQueries chan<- foundation.StatusQuery)
}

type Router interface {
	Run(done <-chan struct{})
}

type router struct {
	swapper    swapper.Swapper
	statusBook status.Book
	storage    Storage
	logger     logrus.FieldLogger
	servers    []Server
}

func New(swapper swapper.Swapper, statusBook status.Book, storage Storage, logger logrus.FieldLogger, servers ...Server) Router {
	return &router{swapper, statusBook, storage, logger, servers}
}

func (router *router) Run(done <-chan struct{}) {
	swapRequests := make(chan foundation.SwapRequest)
	statusQueries := make(chan foundation.StatusQuery)
	statuses := make(chan foundation.SwapStatus)
	swaps := make(chan foundation.SwapRequest)
	updateRequests := make(chan foundation.StatusUpdate)
	updates := make(chan foundation.StatusUpdate)
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

	// starting the swapperd's router
	co.ParBegin(
		func() {
			co.ForAll(router.servers, func(i int) {
				router.servers[i].Run(done, swapRequests, statusQueries)
			})
		},
		func() {
			router.swapper.Run(done, swaps, results, updateRequests)
		},
		func() {
			router.statusBook.Run(done, statuses, updates, statusQueries)
		},
		func() {
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
		},
	)
}
