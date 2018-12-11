package swapper

import (
	"time"

	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
)

type Storage interface {
	UpdateStatus(update swap.StatusUpdate) error
	PendingSwap(swap.SwapID) (swap.SwapBlob, error)
	DeletePendingSwap(swap.SwapID) error
}

type Swapper interface {
	Run(done <-chan struct{}, swaps <-chan swap.SwapBlob, updates chan<- swap.StatusUpdate)
}

type Contract interface {
	Initiate() error
	Audit() error
	Redeem([32]byte) error
	AuditSecret() ([32]byte, error)
	Refund() error
}

type ContractBuilder interface {
	BuildSwapContracts(swap swap.SwapBlob) (Contract, Contract, error)
}

type DelayCallback interface {
	DelayCallback(swap.SwapBlob) (swap.SwapBlob, error)
}

type swapper struct {
	callback DelayCallback
	builder  ContractBuilder
	storage  Storage
	logger   logrus.FieldLogger
}

func New(callback DelayCallback, builder ContractBuilder, storage Storage, logger logrus.FieldLogger) Swapper {
	return &swapper{
		callback: callback,
		builder:  builder,
		storage:  storage,
		logger:   logger,
	}
}

func (swapper *swapper) Run(done <-chan struct{}, swaps <-chan swap.SwapBlob, updates chan<- swap.StatusUpdate) {
	for {
		select {
		case <-done:
			return
		case blob, ok := <-swaps:
			if !ok {
				return
			}
			go swapper.swap(blob, updates)
		}
	}
}

func (swapper *swapper) swap(blob swap.SwapBlob, updates chan<- swap.StatusUpdate) {
	logger := swapper.logger.WithField("SwapID", blob.ID)
	native, foreign, err := swapper.builder.BuildSwapContracts(blob)
	if err != nil {
		logger.Error(err)
		return
	}
	if blob.Delay {
		password := blob.Password
		blob.Password = ""
		filledSwap, err := swapper.callback.DelayCallback(blob)
		if err != nil {
			logger.Error(err)
			return
		}
		blob = filledSwap
		blob.Password = password
	}
	if blob.ShouldInitiateFirst {
		swapper.initiate(blob, native, foreign, updates)
		return
	}
	swapper.respond(blob, native, foreign, updates)
}

func (swapper *swapper) initiate(blob swap.SwapBlob, native, foreign Contract, updates chan<- swap.StatusUpdate) {
	var update = swap.NewStatusUpdate(blob.ID, swap.Inactive)
	defer func() {
		updates <- update
		if err := swapper.storage.UpdateStatus(update); err != nil {
			swapper.logger.Error(err)
		}
	}()

	secret := sha3.Sum256(append([]byte(blob.Password), []byte(blob.ID)...))
	logger := swapper.logger.WithField("SwapID", blob.ID)
	if err := native.Initiate(); err != nil {
		logger.Error(err)
		swapper.handleResult(blob, false, updates)
		return
	}
	update.Code = swap.Initiated
	if err := foreign.Audit(); err != nil {
		update.Code = swap.AuditFailed
		if err := native.Refund(); err != nil {
			logger.Error(err)
			swapper.handleResult(blob, false, updates)
			return
		}
		swapper.handleResult(blob, true, updates)
		update.Code = swap.Refunded
		return
	}
	update.Code = swap.Audited
	if err := foreign.Redeem(secret); err != nil {
		logger.Error(err)
		swapper.handleResult(blob, false, updates)
		return
	}
	update.Code = swap.Redeemed
	swapper.handleResult(blob, true, updates)
}

func (swapper *swapper) respond(blob swap.SwapBlob, native, foreign Contract, updates chan<- swap.StatusUpdate) {
	var update = swap.NewStatusUpdate(blob.ID, swap.Inactive)
	defer func() {
		updates <- update
		if err := swapper.storage.UpdateStatus(update); err != nil {
			swapper.logger.Error(err)
		}
	}()

	logger := swapper.logger.WithField("SwapID", blob.ID)
	if err := foreign.Audit(); err != nil {
		update.Code = swap.AuditFailed
		swapper.handleResult(blob, true, updates)
		return
	}

	update.Code = swap.Audited
	if err := native.Initiate(); err != nil {
		logger.Error(err)
		swapper.handleResult(blob, false, updates)
		return
	}

	update.Code = swap.Initiated
	secret, err := native.AuditSecret()
	if err != nil {
		if err := native.Refund(); err != nil {
			logger.Error(err)
			swapper.handleResult(blob, false, updates)
			return
		}
		update.Code = swap.Refunded
		swapper.handleResult(blob, true, updates)
		return
	}
	if err := foreign.Redeem(secret); err != nil {
		logger.Error(err)
		swapper.handleResult(blob, false, updates)
		return
	}
	update.Code = swap.Redeemed
	swapper.handleResult(blob, true, updates)
}

func (swapper *swapper) handleResult(blob swap.SwapBlob, remove bool, updates chan<- swap.StatusUpdate) {
	if remove {
		if err := swapper.storage.DeletePendingSwap(blob.ID); err != nil {
			swapper.logger.Error(err)
		}
		return
	}
	time.Sleep(5 * time.Minute)
	swapper.swap(blob, updates)
}
