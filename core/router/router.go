package router

import (
	"time"

	"github.com/republicprotocol/swapperd/core/status"
	"github.com/republicprotocol/swapperd/core/swapper"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/core/request"
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
	Run(done <-chan struct{})
}

type router struct {
	swapper    swapper.Swapper
	statusBook status.Book
	storage    Storage
	logger     foundation.Logger
	listeners  []request.Listener
}

func New(swapper swapper.Swapper, statusBook status.Book, storage Storage, logger foundation.Logger, listeners ...request.Listener) Router {
	return &router{swapper, statusBook, storage, logger, listeners}
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
			co.ForAll(router.listeners, func(i int) {
				router.listeners[i].Run(done, swapRequests, statusQueries)
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
					router.storage.InsertSwap(swap)
					statuses <- foundation.NewSwapStatus(swap.SwapBlob)
					router.logger.LogInfo(swap.ID, "adding to the swap queue")
					swaps <- swap
				case update := <-updateRequests:
					if err := router.storage.UpdateStatus(update); err != nil {
						router.logger.LogError(update.ID, err)
					}
					updates <- update
				case result := <-results:
					if result.Success {
						if err := router.storage.DeletePendingSwap(result.ID); err != nil {
							router.logger.LogError(result.ID, err)
							continue
						}
						router.logger.LogInfo(result.ID, "removed from pending swaps")
						continue
					}
					swap, err := router.storage.PendingSwap(result.ID)
					if err != nil {
						router.logger.LogError(result.ID, err)
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
