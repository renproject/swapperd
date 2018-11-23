package swapper

import (
	"github.com/republicprotocol/swapperd/foundation"
)

type Swapper interface {
	Run(done <-chan struct{}, swaps <-chan foundation.SwapRequest, results chan<- foundation.SwapResult, updates chan<- foundation.StatusUpdate)
}

type Contract interface {
	Initiate() error
	Audit() error
	Redeem([32]byte) error
	AuditSecret() ([32]byte, error)
	Refund() error
}

type ContractBuilder interface {
	BuildSwapContracts(swap foundation.SwapRequest) (Contract, Contract, error)
}

type DelayCallback interface {
	DelayCallback(foundation.SwapBlob) (foundation.SwapBlob, error)
}

type swapper struct {
	callback DelayCallback
	builder  ContractBuilder
	logger   foundation.Logger
}

func New(callback DelayCallback, builder ContractBuilder, logger foundation.Logger) Swapper {
	return &swapper{
		callback: callback,
		builder:  builder,
		logger:   logger,
	}
}

func (swapper *swapper) Run(done <-chan struct{}, swaps <-chan foundation.SwapRequest, results chan<- foundation.SwapResult, updates chan<- foundation.StatusUpdate) {
	for {
		select {
		case <-done:
			return
		case swap, ok := <-swaps:
			if !ok {
				return
			}
			native, foreign, err := swapper.builder.BuildSwapContracts(swap)
			if err != nil {
				swapper.logger.LogError(swap.ID, err)
				results <- foundation.NewSwapResult(swap.ID, false)
				continue
			}
			if swap.Delay {
				filledSwap, err := swapper.callback.DelayCallback(swap.SwapBlob)
				if err != nil {
					swapper.logger.LogError(swap.ID, err)
					results <- foundation.NewSwapResult(swap.ID, false)
					continue
				}
				swap.SwapBlob = filledSwap
			}
			if swap.ShouldInitiateFirst {
				go swapper.initiate(results, updates, native, foreign, swap)
				continue
			}
			go swapper.respond(results, updates, native, foreign, swap)
		}
	}
}

func (swapper *swapper) initiate(results chan<- foundation.SwapResult, updates chan<- foundation.StatusUpdate, native, foreign Contract, swap foundation.SwapRequest) {
	if err := native.Initiate(); err != nil {
		swapper.logger.LogError(swap.ID, err)
		results <- foundation.NewSwapResult(swap.ID, false)
		return
	}
	updates <- foundation.NewStatusUpdate(swap.ID, foundation.Initiated)
	if err := foreign.Audit(); err != nil {
		updates <- foundation.NewStatusUpdate(swap.ID, foundation.AuditFailed)
		if err := native.Refund(); err != nil {
			swapper.logger.LogError(swap.ID, err)
			results <- foundation.NewSwapResult(swap.ID, false)
			return
		}
		results <- foundation.NewSwapResult(swap.ID, true)
		updates <- foundation.NewStatusUpdate(swap.ID, foundation.Refunded)
		return
	}
	updates <- foundation.NewStatusUpdate(swap.ID, foundation.Audited)
	if err := foreign.Redeem(swap.Secret); err != nil {
		swapper.logger.LogError(swap.ID, err)
		results <- foundation.NewSwapResult(swap.ID, false)
		return
	}
	updates <- foundation.NewStatusUpdate(swap.ID, foundation.Redeemed)
	results <- foundation.NewSwapResult(swap.ID, true)
}

func (swapper *swapper) respond(results chan<- foundation.SwapResult, updates chan<- foundation.StatusUpdate, native, foreign Contract, swap foundation.SwapRequest) {
	if err := foreign.Audit(); err != nil {
		updates <- foundation.NewStatusUpdate(swap.ID, foundation.AuditFailed)
		results <- foundation.NewSwapResult(swap.ID, true)
		return
	}
	updates <- foundation.NewStatusUpdate(swap.ID, foundation.Audited)
	if err := native.Initiate(); err != nil {
		swapper.logger.LogError(swap.ID, err)
		results <- foundation.NewSwapResult(swap.ID, false)
		return
	}
	updates <- foundation.NewStatusUpdate(swap.ID, foundation.Initiated)
	secret, err := native.AuditSecret()
	if err != nil {
		if err := native.Refund(); err != nil {
			swapper.logger.LogError(swap.ID, err)
			results <- foundation.NewSwapResult(swap.ID, false)
			return
		}
		updates <- foundation.NewStatusUpdate(swap.ID, foundation.Refunded)
		results <- foundation.NewSwapResult(swap.ID, true)
		return
	}
	if err := foreign.Redeem(secret); err != nil {
		swapper.logger.LogError(swap.ID, err)
		results <- foundation.NewSwapResult(swap.ID, false)
		return
	}
	updates <- foundation.NewStatusUpdate(swap.ID, foundation.Redeemed)
	results <- foundation.NewSwapResult(swap.ID, true)
}
