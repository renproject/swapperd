package swapper

import (
	"github.com/republicprotocol/swapperd/foundation"
)

type swapper struct {
	Storage
	BinderBuilder
	Semaphore
	Logger
}

type result struct {
	ID      foundation.SwapID
	Success bool
}

type Query struct {
	Swap     foundation.Swap
	Password string
}

type Logger interface {
	LogInfo(foundation.SwapID, string)
	LogDebug(foundation.SwapID, string)
	LogError(foundation.SwapID, error)
}

type Storage interface {
	LoadSwaps() []foundation.Swap
	DeleteSwap(foundation.SwapID) error
	AddSwap(swap foundation.Swap) error
}

type SwapContractBinder interface {
	Initiate() error
	Audit() error
	Redeem([32]byte) error
	AuditSecret() ([32]byte, error)
	Refund() error
}

type BinderBuilder interface {
	BuildBinders(query Query) (SwapContractBinder, SwapContractBinder, error)
}

type Swapper interface {
	Run(swaps <-chan Query, swapStatuses chan<- foundation.SwapStatus, done <-chan struct{})
}

func New(storage Storage, builder BinderBuilder, logger Logger) Swapper {
	return &swapper{
		Storage:       storage,
		BinderBuilder: builder,
		Semaphore:     NewSemaphore(),
		Logger:        logger,
	}
}

func NewQuery(swap foundation.Swap, password string) Query {
	return Query{swap, password}
}

func (swapper *swapper) Run(queries <-chan Query, swapStatuses chan<- foundation.SwapStatus, done <-chan struct{}) {
	results := make(chan result)
	for {
		select {
		case <-done:
			return
		case query := <-queries:
			swapper.AddSwap(query.Swap)
			native, foreign, err := swapper.BuildBinders(query)
			if err != nil {
				swapper.LogError(query.Swap.ID, err)
				continue
			}
			// if swapper.TryWait(swap.ID) {
			// 	continue
			// }
			go func() {
				// defer swapper.Signal(swap.ID)
				execute(native, foreign, swapper.Logger, query.Swap, swapStatuses, results)
			}()
		case swapResult := <-results:
			if swapResult.Success {
				swapper.DeleteSwap(swapResult.ID)
			}
			// default:
			// 	swapsToRetry := swapper.LoadSwaps()
			// 	go co.ForAll(swapsToRetry, func(i int) {
			// 		swap := swapsToRetry[i]
			// 		if swapper.TryWait(swap.ID) {
			// 			return
			// 		}
			// 		defer swapper.Signal(swap.ID)
			// 		native, foreign, err := swapper.buildBinders(swap)
			// 		if err != nil {
			// 			return
			// 		}
			// 		execute(native, foreign, swapper.Logger, swap, swapStatuses, results)
			// 	})
			// 	time.Sleep(1 * time.Minute)
		}
	}
}

func execute(native, foreign SwapContractBinder, logger Logger, req foundation.Swap, statuses chan<- foundation.SwapStatus, done chan<- result) {
	if req.IsFirst {
		initiate(native, foreign, logger, req, statuses, done)
	}
	respond(native, foreign, logger, req, statuses, done)
}

func initiate(native, foreign SwapContractBinder, logger Logger, req foundation.Swap, statuses chan<- foundation.SwapStatus, done chan<- result) {
	if err := native.Initiate(); err != nil {
		logger.LogError(req.ID, err)
		done <- result{req.ID, false}
		return
	}
	statuses <- foundation.SwapStatus{req.ID, foundation.INITIATED}
	if err := foreign.Audit(); err != nil {
		statuses <- foundation.SwapStatus{req.ID, foundation.AUDIT_FAILED}
		if err := native.Refund(); err != nil {
			logger.LogError(req.ID, err)
			done <- result{req.ID, false}
			return
		}
		done <- result{req.ID, true}
		statuses <- foundation.SwapStatus{req.ID, foundation.REFUNDED}
		return
	}
	statuses <- foundation.SwapStatus{req.ID, foundation.AUDITED}
	if err := foreign.Redeem(req.Secret); err != nil {
		logger.LogError(req.ID, err)
		done <- result{req.ID, false}
		return
	}
	statuses <- foundation.SwapStatus{req.ID, foundation.REDEEMED}
	done <- result{req.ID, true}
	return
}

func respond(native, foreign SwapContractBinder, logger Logger, req foundation.Swap, statuses chan<- foundation.SwapStatus, done chan<- result) {
	if err := foreign.Audit(); err != nil {
		statuses <- foundation.SwapStatus{req.ID, foundation.AUDIT_FAILED}
		done <- result{req.ID, true}
		return
	}
	statuses <- foundation.SwapStatus{req.ID, foundation.AUDITED}
	if err := native.Initiate(); err != nil {
		logger.LogError(req.ID, err)
		done <- result{req.ID, false}
		return
	}
	statuses <- foundation.SwapStatus{req.ID, foundation.INITIATED}
	secret, err := native.AuditSecret()
	if err != nil {
		if err := native.Refund(); err != nil {
			logger.LogError(req.ID, err)
			done <- result{req.ID, false}
			return
		}
		statuses <- foundation.SwapStatus{req.ID, foundation.REFUNDED}
		done <- result{req.ID, true}
	}
	if err := foreign.Redeem(secret); err != nil {
		logger.LogError(req.ID, err)
		done <- result{req.ID, false}
		return
	}
	statuses <- foundation.SwapStatus{req.ID, foundation.REDEEMED}
	done <- result{req.ID, true}
}
