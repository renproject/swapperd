package swapper

import (
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/sirupsen/logrus"
)

type Swapper interface {
	Run(done <-chan struct{}, swaps <-chan swap.SwapRequest, results chan<- swap.SwapResult, updates chan<- swap.StatusUpdate)
}

type Contract interface {
	Initiate() error
	Audit() error
	Redeem([32]byte) error
	AuditSecret() ([32]byte, error)
	Refund() error
}

type ContractBuilder interface {
	BuildSwapContracts(swap swap.SwapRequest) (Contract, Contract, error)
}

type DelayCallback interface {
	DelayCallback(swap.SwapBlob) (swap.SwapBlob, error)
}

type swapper struct {
	callback DelayCallback
	builder  ContractBuilder
	logger   logrus.FieldLogger
}

func New(callback DelayCallback, builder ContractBuilder, logger logrus.FieldLogger) Swapper {
	return &swapper{
		callback: callback,
		builder:  builder,
		logger:   logger,
	}
}

func (swapper *swapper) Run(done <-chan struct{}, swaps <-chan swap.SwapRequest, results chan<- swap.SwapResult, updates chan<- swap.StatusUpdate) {
	for {
		select {
		case <-done:
			return
		case swapRequest, ok := <-swaps:
			if !ok {
				return
			}
			logger := swapper.logger.WithField("SwapID", swapRequest.ID)
			native, foreign, err := swapper.builder.BuildSwapContracts(swapRequest)
			if err != nil {
				logger.Error(err)
				results <- swap.NewSwapResult(swapRequest.ID, false)
				continue
			}
			if swapRequest.Delay {
				filledSwap, err := swapper.callback.DelayCallback(swapRequest.SwapBlob)
				if err != nil {
					logger.Error(err)
					results <- swap.NewSwapResult(swapRequest.ID, false)
					continue
				}
				swapRequest.SwapBlob = filledSwap
			}
			if swapRequest.ShouldInitiateFirst {
				go swapper.initiate(results, updates, native, foreign, swapRequest)
				continue
			}
			go swapper.respond(results, updates, native, foreign, swapRequest)
		}
	}
}

func (swapper *swapper) initiate(results chan<- swap.SwapResult, updates chan<- swap.StatusUpdate, native, foreign Contract, swapRequest swap.SwapRequest) {
	logger := swapper.logger.WithField("SwapID", swapRequest.ID)
	if err := native.Initiate(); err != nil {
		logger.Error(err)
		results <- swap.NewSwapResult(swapRequest.ID, false)
		return
	}
	updates <- swap.NewStatusUpdate(swapRequest.ID, swap.Initiated)
	if err := foreign.Audit(); err != nil {
		updates <- swap.NewStatusUpdate(swapRequest.ID, swap.AuditFailed)
		if err := native.Refund(); err != nil {
			logger.Error(err)
			results <- swap.NewSwapResult(swapRequest.ID, false)
			return
		}
		results <- swap.NewSwapResult(swapRequest.ID, true)
		updates <- swap.NewStatusUpdate(swapRequest.ID, swap.Refunded)
		return
	}
	updates <- swap.NewStatusUpdate(swapRequest.ID, swap.Audited)
	if err := foreign.Redeem(swapRequest.Secret); err != nil {
		logger.Error(err)
		results <- swap.NewSwapResult(swapRequest.ID, false)
		return
	}
	updates <- swap.NewStatusUpdate(swapRequest.ID, swap.Redeemed)
	results <- swap.NewSwapResult(swapRequest.ID, true)
}

func (swapper *swapper) respond(results chan<- swap.SwapResult, updates chan<- swap.StatusUpdate, native, foreign Contract, swapRequest swap.SwapRequest) {
	logger := swapper.logger.WithField("SwapID", swapRequest.ID)
	if err := foreign.Audit(); err != nil {
		updates <- swap.NewStatusUpdate(swapRequest.ID, swap.AuditFailed)
		results <- swap.NewSwapResult(swapRequest.ID, true)
		return
	}
	updates <- swap.NewStatusUpdate(swapRequest.ID, swap.Audited)
	if err := native.Initiate(); err != nil {
		logger.Error(err)
		results <- swap.NewSwapResult(swapRequest.ID, false)
		return
	}
	updates <- swap.NewStatusUpdate(swapRequest.ID, swap.Initiated)
	secret, err := native.AuditSecret()
	if err != nil {
		if err := native.Refund(); err != nil {
			logger.Error(err)
			results <- swap.NewSwapResult(swapRequest.ID, false)
			return
		}
		updates <- swap.NewStatusUpdate(swapRequest.ID, swap.Refunded)
		results <- swap.NewSwapResult(swapRequest.ID, true)
		return
	}
	if err := foreign.Redeem(secret); err != nil {
		logger.Error(err)
		results <- swap.NewSwapResult(swapRequest.ID, false)
		return
	}
	updates <- swap.NewStatusUpdate(swapRequest.ID, swap.Redeemed)
	results <- swap.NewSwapResult(swapRequest.ID, true)
}
