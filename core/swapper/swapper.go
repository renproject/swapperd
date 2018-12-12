package swapper

import (
	"time"

	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
)

type Storage interface {
	LoadCosts(id swap.SwapID) (blockchain.Cost, blockchain.Cost)
	DeletePendingSwap(swap.SwapID) error
}

type Swapper interface {
	Run(done <-chan struct{}, swaps <-chan swap.SwapBlob, updates chan<- swap.ReceiptUpdate)
}

type Contract interface {
	Initiate() error
	Audit() error
	Redeem([32]byte) error
	AuditSecret() ([32]byte, error)
	Refund() error
	Cost() blockchain.Cost
}

type ContractBuilder interface {
	BuildSwapContracts(swap swap.SwapBlob, sendCost, receiveCost blockchain.Cost) (Contract, Contract, error)
}

type swapper struct {
	builder ContractBuilder
	storage Storage
	logger  logrus.FieldLogger
}

func New(builder ContractBuilder, storage Storage, logger logrus.FieldLogger) Swapper {
	return &swapper{
		builder: builder,
		storage: storage,
		logger:  logger,
	}
}

type Bootload struct {
	Password string
}

func (swapper *swapper) Run(done <-chan struct{}, swaps <-chan swap.SwapBlob, updates chan<- swap.ReceiptUpdate) {
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

func (swapper *swapper) swap(blob swap.SwapBlob, updates chan<- swap.ReceiptUpdate) {
	logger := swapper.logger.WithField("SwapID", blob.ID)

	sendCost, receiveCost := swapper.storage.LoadCosts(blob.ID)
	native, foreign, err := swapper.builder.BuildSwapContracts(blob, sendCost, receiveCost)
	if err != nil {
		logger.Error(err)
		return
	}
	if blob.ShouldInitiateFirst {
		swapper.initiate(blob, native, foreign, updates)
		return
	}
	swapper.respond(blob, native, foreign, updates)
}

func (swapper *swapper) initiate(blob swap.SwapBlob, native, foreign Contract, updates chan<- swap.ReceiptUpdate) {
	swapStatus := swap.Inactive
	defer func() {
		updates <- swap.NewReceiptUpdate(blob.ID, func(receipt *swap.SwapReceipt) {
			receipt.Status = swapStatus
			receipt.SendCost = native.Cost()
			receipt.ReceiveCost = foreign.Cost()
		})
	}()

	secret := sha3.Sum256(append([]byte(blob.Password), []byte(blob.ID)...))
	logger := swapper.logger.WithField("SwapID", blob.ID)
	if err := native.Initiate(); err != nil {
		logger.Error(err)
		swapper.handleResult(blob, false, updates)
		return
	}

	swapStatus = swap.Initiated
	if err := foreign.Audit(); err != nil {
		swapStatus = swap.AuditFailed
		if err := native.Refund(); err != nil {
			logger.Error(err)
			swapper.handleResult(blob, false, updates)
			return
		}
		swapper.handleResult(blob, true, updates)
		swapStatus = swap.Refunded
		return
	}

	swapStatus = swap.Audited
	if err := foreign.Redeem(secret); err != nil {
		logger.Error(err)
		swapper.handleResult(blob, false, updates)
		return
	}

	swapStatus = swap.Redeemed
	swapper.handleResult(blob, true, updates)
}

func (swapper *swapper) respond(blob swap.SwapBlob, native, foreign Contract, updates chan<- swap.ReceiptUpdate) {
	swapStatus := swap.Inactive
	defer func() {
		updates <- swap.NewReceiptUpdate(blob.ID, func(receipt *swap.SwapReceipt) {
			receipt.Status = swapStatus
			receipt.SendCost = native.Cost()
			receipt.ReceiveCost = foreign.Cost()
		})
	}()

	logger := swapper.logger.WithField("SwapID", blob.ID)
	if err := foreign.Audit(); err != nil {
		swapStatus = swap.AuditFailed
		swapper.handleResult(blob, true, updates)
		return
	}

	swapStatus = swap.Audited
	if err := native.Initiate(); err != nil {
		logger.Error(err)
		swapper.handleResult(blob, false, updates)
		return
	}

	swapStatus = swap.Initiated
	secret, err := native.AuditSecret()
	if err != nil {
		if err := native.Refund(); err != nil {
			logger.Error(err)
			swapper.handleResult(blob, false, updates)
			return
		}
		swapStatus = swap.Refunded
		swapper.handleResult(blob, true, updates)
		return
	}
	if err := foreign.Redeem(secret); err != nil {
		logger.Error(err)
		swapper.handleResult(blob, false, updates)
		return
	}
	swapStatus = swap.Redeemed
	swapper.handleResult(blob, true, updates)
}

func (swapper *swapper) handleResult(blob swap.SwapBlob, remove bool, updates chan<- swap.ReceiptUpdate) {
	if remove {
		if err := swapper.storage.DeletePendingSwap(blob.ID); err != nil {
			swapper.logger.Error(err)
		}
		return
	}
	time.Sleep(5 * time.Minute)
	swapper.swap(blob, updates)
}
