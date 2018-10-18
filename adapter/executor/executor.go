package executor

import (
	"github.com/republicprotocol/libbtc-go"
	"github.com/republicprotocol/swapperd/adapter/btc"
	"github.com/republicprotocol/swapperd/adapter/eth/erc20"
	"github.com/republicprotocol/swapperd/adapter/eth/eth"
	"github.com/republicprotocol/swapperd/adapter/storage"
	"github.com/republicprotocol/swapperd/core"
	"github.com/republicprotocol/swapperd/foundation"
)

type Executor struct {
	accounts  account.Accounts
	semaphore Semaphore
	// storage.Storage
	core.Logger
}

func New(accounts account.Accounts, storage storage.Storage, logger core.Logger) Executor {
	return Executor{
		accounts,
		NewSemaphore(),
		logger,
	}
}

func (executor *Executor) Run(swaps <-chan foundation.Swap, results chan<- foundation.SwapID) {
	done := make(chan foundation.SwapID)
	// swapsToRetry := executor.LoadSwaps()
	// go co.ForAll(swapsToRetry, func(i int) {
	// 	swap := swapsToRetry[i]
	// 	if executor.semaphore.TryWait(swap.ID) {
	// 		return
	// 	}
	// 	defer executor.semaphore.Signal(swap.ID)
	// 	native, foreign, err := executor.buildBinders(swap)
	// 	if err != nil {
	// 		return
	// 	}
	// 	swapper.Swap(native, foreign, swap, done)
	// 	results <- swap.ID
	// })

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

func (executor *Executor) buildBinders(swap foundation.Swap) (core.SwapContractBinder, core.SwapContractBinder, error) {

	return nil, nil, nil
}

func (executor *Executor) buildBinder(token foundation.Token, swap foundation.Swap) (core.SwapContractBinder, error) {
	switch token {
	case foundation.TokenBTC:
		return btc.NewBitcoinAtom(executor.GetAccount(token).(libbtc.Account), executor.Logger, swap)
	case foundation.TokenETH:
		return eth.NewEthereumAtom(executor.GetAccount(token).(libbtc.Account), executor.Logger, swap)
	case foundation.TokenWBTC:
		return erc20.NewERC20Atom(executor.GetAccount(token).(libbtc.Account), executor.Logger, swap)
	default:
		return nil, foundation.NewErrUnsupportedToken(token.Name)
	}
}
