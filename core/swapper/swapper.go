package swapper

import (
	"time"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/swapperd/foundation"
)

type Swap struct {
	foundation.SwapBlob

	Secret   [32]byte
	Password string
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

type Storage interface {
	InsertSwap(swap Swap) error
	DeleteSwap(foundation.SwapID) error
	PendingSwaps() []Swap
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
	builder ContractBuilder
	storage Storage
	logger  Logger

	sem *binarySemaphore
}

type result struct {
	id      foundation.SwapID
	success bool
}

func newResult(id foundation.SwapID, success bool) result {
	return result{id, success}
}

func New(builder ContractBuilder, storage Storage, logger Logger) Swapper {
	return &swapper{
		builder: builder,
		storage: storage,
		logger:  logger,

		sem: newBinarySemaphore(),
	}
}

func (swapper *swapper) Run(done <-chan struct{}, swaps <-chan Swap, statuses chan<- foundation.SwapStatus) {
	results := make(chan result)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {

		case <-done:
			return

		case swap := <-swaps:
			swapper.storage.InsertSwap(swap)
			native, foreign, err := swapper.builder.BuildSwapContracts(swap)
			if err != nil {
				swapper.logger.LogError(swap.ID, err)
				continue
			}
			if !swapper.sem.TryWait(swap.ID) {
				continue
			}
			go func() {
				defer swapper.sem.Signal(swap.ID)
				execute(results, statuses, native, foreign, swap, swapper.logger)
			}()

		case result := <-results:
			if result.success {
				swapper.storage.DeleteSwap(result.id)
			}

		case <-ticker.C:
			swapsToRetry := swapper.storage.PendingSwaps()
			go co.ForAll(swapsToRetry, func(i int) {
				swap := swapsToRetry[i]
				if swapper.sem.TryWait(swap.ID) {
					return
				}
				defer swapper.sem.Signal(swap.ID)
				native, foreign, err := swapper.builder.BuildSwapContracts(swap)
				if err != nil {
					return
				}
				execute(results, statuses, native, foreign, swap, swapper.logger)
			})
		}
	}
}

func execute(results chan<- result, statuses chan<- foundation.SwapStatus, native, foreign Contract, swap Swap, logger Logger) {
	if swap.ShouldInitiateFirst {
		initiate(results, statuses, native, foreign, swap, logger)
	}
	respond(results, statuses, native, foreign, swap, logger)
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
	return
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
	}
	if err := foreign.Redeem(secret); err != nil {
		logger.LogError(swap.ID, err)
		results <- newResult(swap.ID, false)
		return
	}
	statuses <- foundation.NewSwapStatus(swap.ID, foundation.Redeemed)
	results <- newResult(swap.ID, true)
}
