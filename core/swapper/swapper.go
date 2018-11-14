package swapper

import (
	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/foundation"
)

type Swap struct {
	foundation.SwapBlob

	Secret   [32]byte `json:"secret"`
	Password string   `json:"password"`
}

func NewSwap(swapBlob foundation.SwapBlob, secret [32]byte, password string) Swap {
	return Swap{swapBlob, secret, password}
}

type Contract interface {
	Initiate() error
	Audit() error
	Redeem([32]byte) error
	AuditSecret() ([32]byte, error)
	Refund() error
}

type ContractBuilder interface {
	BuildSwapContracts(swap Swap) (Contract, Contract, error)
}

type SwapFiller interface {
	DelayCallback(foundation.SwapBlob) (foundation.SwapBlob, error)
}

type Storage interface {
	InsertSwap(swap Swap) error
	PendingSwap(foundation.SwapID) (Swap, error)
	DeletePendingSwap(foundation.SwapID) error
	PendingSwaps() ([]Swap, error)
}

type Logger interface {
	LogInfo(foundation.SwapID, string)
	LogDebug(foundation.SwapID, string)
	LogError(foundation.SwapID, error)
}

type Swapper interface {
	Run(done <-chan struct{}, swaps <-chan Swap, statuses chan<- foundation.SwapStatus)
}

type swapper struct {
	filler  SwapFiller
	builder ContractBuilder
	storage Storage
	logger  Logger
}

type result struct {
	id      foundation.SwapID
	success bool
}

func newResult(id foundation.SwapID, success bool) result {
	return result{id, success}
}

func New(filler SwapFiller, builder ContractBuilder, storage Storage, logger Logger) Swapper {
	return &swapper{
		filler:  filler,
		builder: builder,
		storage: storage,
		logger:  logger,
	}
}

func (swapper *swapper) Run(done <-chan struct{}, swaps <-chan Swap, statuses chan<- foundation.SwapStatus) {
	results := make(chan result)

	swapsToRetry, err := swapper.storage.PendingSwaps()
	if err != nil {
		return
	}
	co.ForAll(swapsToRetry, func(i int) {
		swap := swapsToRetry[i]
		go swapper.execute(results, statuses, swap)
	})

	for {
		select {
		case <-done:
			return

		case swap := <-swaps:
			swapper.storage.InsertSwap(swap)
			go swapper.execute(results, statuses, swap)

		case result := <-results:
			if result.success {
				if err := swapper.storage.DeletePendingSwap(result.id); err != nil {
					swapper.logger.LogError(result.id, err)
					continue
				}
				swapper.logger.LogInfo(result.id, "removed from pending swaps")
				continue
			}
			swap, err := swapper.storage.PendingSwap(result.id)
			if err != nil {
				swapper.logger.LogError(result.id, err)
				continue
			}
			go swapper.execute(results, statuses, swap)
		}
	}
}

func (swapper *swapper) execute(results chan<- result, statuses chan<- foundation.SwapStatus, swap Swap) {
	native, foreign, err := swapper.builder.BuildSwapContracts(swap)
	if err != nil {
		swapper.logger.LogError(swap.ID, err)
		results <- newResult(swap.ID, false)
		return
	}

	if swap.Delay {
		filledSwap, err := swapper.filler.DelayCallback(swap.SwapBlob)
		if err != nil {
			swapper.logger.LogError(swap.ID, err)
			results <- newResult(swap.ID, false)
			return
		}
		swap.SwapBlob = filledSwap
	}

	if swap.ShouldInitiateFirst {
		initiate(results, statuses, native, foreign, swap, swapper.logger)
		return
	}
	respond(results, statuses, native, foreign, swap, swapper.logger)
}

func initiate(results chan<- result, statuses chan<- foundation.SwapStatus, native, foreign Contract, swap Swap, logger Logger) {
	if err := native.Initiate(); err != nil {
		logger.LogError(swap.ID, err)
		results <- newResult(swap.ID, false)
		return
	}
	statuses <- foundation.NewSwapStatus(swap.ID, foundation.Initiated)
	if err := foreign.Audit(); err != nil {
		statuses <- foundation.NewSwapStatus(swap.ID, foundation.AuditFailed)
		if err := native.Refund(); err != nil {
			logger.LogError(swap.ID, err)
			results <- newResult(swap.ID, false)
			return
		}
		results <- newResult(swap.ID, true)
		statuses <- foundation.NewSwapStatus(swap.ID, foundation.Refunded)
		return
	}
	statuses <- foundation.NewSwapStatus(swap.ID, foundation.Audited)
	if err := foreign.Redeem(swap.Secret); err != nil {
		logger.LogError(swap.ID, err)
		results <- newResult(swap.ID, false)
		return
	}
	statuses <- foundation.NewSwapStatus(swap.ID, foundation.Redeemed)
	results <- newResult(swap.ID, true)
}

func respond(results chan<- result, statuses chan<- foundation.SwapStatus, native, foreign Contract, swap Swap, logger Logger) {
	if err := foreign.Audit(); err != nil {
		statuses <- foundation.NewSwapStatus(swap.ID, foundation.AuditFailed)
		results <- newResult(swap.ID, true)
		return
	}
	statuses <- foundation.NewSwapStatus(swap.ID, foundation.Audited)
	if err := native.Initiate(); err != nil {
		logger.LogError(swap.ID, err)
		results <- newResult(swap.ID, false)
		return
	}
	statuses <- foundation.NewSwapStatus(swap.ID, foundation.Initiated)
	secret, err := native.AuditSecret()
	if err != nil {
		if err := native.Refund(); err != nil {
			logger.LogError(swap.ID, err)
			results <- newResult(swap.ID, false)
			return
		}
		statuses <- foundation.NewSwapStatus(swap.ID, foundation.Refunded)
		results <- newResult(swap.ID, true)
		return
	}
	if err := foreign.Redeem(secret); err != nil {
		logger.LogError(swap.ID, err)
		results <- newResult(swap.ID, false)
		return
	}
	statuses <- foundation.NewSwapStatus(swap.ID, foundation.Redeemed)
	results <- newResult(swap.ID, true)
}
