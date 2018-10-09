package executor

import (
	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/adapter/config"
	"github.com/republicprotocol/swapperd/adapter/keystore"
	"github.com/republicprotocol/swapperd/adapter/storage"
	"github.com/republicprotocol/swapperd/core/logger"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation"
)

type Executor struct {
	keystore  keystore.Keystore
	config    config.Config
	semaphore Semaphore
	storage.Storage
	logger.Logger
}

func New(keystore keystore.Keystore, config config.Config, storage storage.Storage, logger logger.Logger) Executor {
	return Executor{
		keystore,
		config,
		NewSemaphore(),
		storage,
		logger,
	}
}

func (executor *Executor) Run(swaps <-chan foundation.Swap, results chan<- foundation.SwapID) {
	done := make(chan foundation.SwapID)
	swapsToRetry := executor.LoadSwaps()
	go co.ForAll(swapsToRetry, func(i int) {
		swap := swapsToRetry[i]
		if executor.semaphore.TryWait(swap.ID) {
			return
		}
		defer executor.semaphore.Signal(swap.ID)
		native, foreign, err := executor.buildBinders(swap)
		if err != nil {
			return
		}
		swapper.Swap(native, foreign, swap, done)
		results <- swap.ID
	})

	for {
		select {
		case swap := <-swaps:
			executor.AddSwap(swap)
			native, foreign, err := executor.buildBinders(swap)
			if err != nil {
				continue
			}
			if executor.semaphore.TryWait(swap.ID) {
				continue
			}
			go func() {
				defer executor.semaphore.Signal(swap.ID)
				swapper.Swap(native, foreign, swap, done)
			}()
			results <- swap.ID
		case swapID := <-done:
			executor.DeleteSwap(swapID)
		}
	}
}

func (executor *Executor) buildBinders(swap foundation.Swap) (swapper.SwapContractBinder, swapper.SwapContractBinder, error) {

	return nil, nil, nil
}

func (executor *Executor) buildBinder() {

}
