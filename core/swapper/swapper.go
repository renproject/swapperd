package swapper

import (
	"time"

	co "github.com/republicprotocol/co-go"
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
	Swap     foundation.SwapRequest
	Password string
	Secret   [32]byte
}

type Logger interface {
	LogInfo(foundation.SwapID, string)
	LogDebug(foundation.SwapID, string)
	LogError(foundation.SwapID, error)
}

type Storage interface {
	LoadPendingQueries() []Query
	DeleteQuery(foundation.SwapID) error
	AddQuery(query Query) error
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

func NewQuery(swap foundation.SwapRequest, secret [32]byte, password string) Query {
	return Query{swap, password, secret}
}

func (swapper *swapper) Run(queries <-chan Query, swapStatuses chan<- foundation.SwapStatus, done <-chan struct{}) {
	results := make(chan result)
	for {
		select {
		case <-done:
			return
		case query := <-queries:
			swapper.AddQuery(query)
			native, foreign, err := swapper.BuildBinders(query)
			if err != nil {
				swapper.LogError(query.Swap.ID, err)
				continue
			}
			if swapper.TryWait(query.Swap.ID) {
				continue
			}
			go func() {
				defer swapper.Signal(query.Swap.ID)
				execute(native, foreign, swapper.Logger, query, swapStatuses, results)
			}()
		case swapResult := <-results:
			if swapResult.Success {
				swapper.DeleteQuery(swapResult.ID)
			}
		default:
			swapsToRetry := swapper.LoadPendingQueries()
			go co.ForAll(swapsToRetry, func(i int) {
				query := swapsToRetry[i]
				if swapper.TryWait(query.Swap.ID) {
					return
				}
				defer swapper.Signal(query.Swap.ID)
				native, foreign, err := swapper.BuildBinders(query)
				if err != nil {
					return
				}
				execute(native, foreign, swapper.Logger, query, swapStatuses, results)
			})
			time.Sleep(1 * time.Minute)
		}
	}
}

func execute(native, foreign SwapContractBinder, logger Logger, query Query, statuses chan<- foundation.SwapStatus, done chan<- result) {
	if query.Swap.ShouldInitiateFirst {
		initiate(native, foreign, logger, query, statuses, done)
	}
	respond(native, foreign, logger, query, statuses, done)
}

func initiate(native, foreign SwapContractBinder, logger Logger, query Query, statuses chan<- foundation.SwapStatus, done chan<- result) {
	if err := native.Initiate(); err != nil {
		logger.LogError(query.Swap.ID, err)
		done <- result{query.Swap.ID, false}
		return
	}
	statuses <- foundation.SwapStatus{query.Swap.ID, foundation.INITIATED}
	if err := foreign.Audit(); err != nil {
		statuses <- foundation.SwapStatus{query.Swap.ID, foundation.AUDIT_FAILED}
		if err := native.Refund(); err != nil {
			logger.LogError(query.Swap.ID, err)
			done <- result{query.Swap.ID, false}
			return
		}
		done <- result{query.Swap.ID, true}
		statuses <- foundation.SwapStatus{query.Swap.ID, foundation.REFUNDED}
		return
	}
	statuses <- foundation.SwapStatus{query.Swap.ID, foundation.AUDITED}
	if err := foreign.Redeem(query.Secret); err != nil {
		logger.LogError(query.Swap.ID, err)
		done <- result{query.Swap.ID, false}
		return
	}
	statuses <- foundation.SwapStatus{query.Swap.ID, foundation.REDEEMED}
	done <- result{query.Swap.ID, true}
	return
}

func respond(native, foreign SwapContractBinder, logger Logger, query Query, statuses chan<- foundation.SwapStatus, done chan<- result) {
	if err := foreign.Audit(); err != nil {
		statuses <- foundation.SwapStatus{query.Swap.ID, foundation.AUDIT_FAILED}
		done <- result{query.Swap.ID, true}
		return
	}
	statuses <- foundation.SwapStatus{query.Swap.ID, foundation.AUDITED}
	if err := native.Initiate(); err != nil {
		logger.LogError(query.Swap.ID, err)
		done <- result{query.Swap.ID, false}
		return
	}
	statuses <- foundation.SwapStatus{query.Swap.ID, foundation.INITIATED}
	secret, err := native.AuditSecret()
	if err != nil {
		if err := native.Refund(); err != nil {
			logger.LogError(query.Swap.ID, err)
			done <- result{query.Swap.ID, false}
			return
		}
		statuses <- foundation.SwapStatus{query.Swap.ID, foundation.REFUNDED}
		done <- result{query.Swap.ID, true}
	}
	if err := foreign.Redeem(secret); err != nil {
		logger.LogError(query.Swap.ID, err)
		done <- result{query.Swap.ID, false}
		return
	}
	statuses <- foundation.SwapStatus{query.Swap.ID, foundation.REDEEMED}
	done <- result{query.Swap.ID, true}
}
