package immediate

import (
	"fmt"

	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/republicprotocol/tau"
	"golang.org/x/crypto/sha3"
)

var ErrSwapExpired = fmt.Errorf("swap expired")
var ErrAuditPending = fmt.Errorf("audit pending")

type Contract interface {
	Initiate() error
	Audit() error
	Redeem([32]byte) error
	AuditSecret() ([32]byte, error)
	Refund() error
	Cost() blockchain.Cost
}

type ContractBuilder interface {
	BuildSwapContracts(request SwapRequest) (Contract, Contract, error)
}

type swapper struct {
	builder ContractBuilder
	swapMap map[swap.SwapID]SwapRequest
}

func New(cap int, builder ContractBuilder) tau.Task {
	return tau.New(tau.NewIO(cap), &swapper{
		builder: builder,
		swapMap: map[swap.SwapID]SwapRequest{},
	})
}

func (swapper *swapper) Reduce(msg tau.Message) tau.Message {
	switch msg := msg.(type) {
	case tau.Tick:
		return swapper.handleTick()
	case SwapRequest:
		return swapper.handleSwap(msg)
	default:
		return tau.NewError(fmt.Errorf("invalid message type in swapper: %T", msg))
	}
}

func (swapper *swapper) handleTick() tau.Message {
	msgs := []tau.Message{}
	for _, req := range swapper.swapMap {
		msgs = append(msgs, swapper.handleSwap(req))
	}
	return tau.NewMessageBatch(msgs)
}

func (swapper *swapper) handleSwap(req SwapRequest) tau.Message {
	native, foreign, err := swapper.builder.BuildSwapContracts(req)
	if err != nil {
		return tau.NewError(err)
	}
	if req.Blob.ShouldInitiateFirst {
		return swapper.initiate(req, native, foreign)
	}
	return swapper.respond(req, native, foreign)
}

func (swapper *swapper) initiate(req SwapRequest, native, foreign Contract) tau.Message {
	secret := sha3.Sum256(append([]byte(req.Blob.Password), []byte(req.Blob.ID)...))
	if err := native.Initiate(); err != nil {
		return swapper.handleResult(req, swap.Inactive, native, foreign, err, false)
	}
	if err := foreign.Audit(); err != nil {
		if err == ErrAuditPending {
			return swapper.handleResult(req, swap.AuditPending, native, foreign, nil, false)
		}
		if err != ErrSwapExpired {
			return swapper.handleResult(req, swap.AuditPending, native, foreign, err, false)
		}
		if err := native.Refund(); err != nil {
			return swapper.handleResult(req, swap.RefundFailed, native, foreign, err, false)
		}
		return swapper.handleResult(req, swap.Refunded, native, foreign, nil, true)
	}
	if err := foreign.Redeem(secret); err != nil {
		return swapper.handleResult(req, swap.Audited, native, foreign, err, false)
	}
	return swapper.handleResult(req, swap.Redeemed, native, foreign, nil, true)
}

func (swapper *swapper) respond(req SwapRequest, native, foreign Contract) tau.Message {
	if err := foreign.Audit(); err != nil {
		if err == ErrAuditPending {
			return swapper.handleResult(req, swap.AuditPending, native, foreign, nil, false)
		}
		if err == ErrSwapExpired {
			return swapper.handleResult(req, swap.Expired, native, foreign, err, true)
		}
		return swapper.handleResult(req, swap.AuditPending, native, foreign, err, false)
	}

	if err := native.Initiate(); err != nil {
		return swapper.handleResult(req, swap.Audited, native, foreign, err, false)
	}
	secret, err := native.AuditSecret()
	if err != nil {
		if err == ErrAuditPending {
			return swapper.handleResult(req, swap.AuditPending, native, foreign, nil, false)
		}
		if err != ErrSwapExpired {
			return swapper.handleResult(req, swap.Initiated, native, foreign, err, false)
		}
		if err := native.Refund(); err != nil {
			return swapper.handleResult(req, swap.RefundFailed, native, foreign, err, false)
		}
		return swapper.handleResult(req, swap.Refunded, native, foreign, nil, true)
	}
	if err := foreign.Redeem(secret); err != nil {
		return swapper.handleResult(req, swap.AuditedSecret, native, foreign, err, false)
	}
	return swapper.handleResult(req, swap.Redeemed, native, foreign, nil, true)
}

func (swapper *swapper) handleResult(req SwapRequest, status int, native, foreign Contract, err error, remove bool) tau.Message {
	messages := []tau.Message{}
	messages = append(messages, NewReceiptUpdate(req.Blob.ID, status, native, foreign))
	if err != nil {
		messages = append(messages, tau.NewError(err))
	}
	if remove {
		delete(swapper.swapMap, req.Blob.ID)
		return tau.NewMessageBatch(append(messages, DeleteSwap{req.Blob.ID}))
	}
	swapper.swapMap[req.Blob.ID] = req
	return tau.NewMessageBatch(messages)
}

type SwapRequest struct {
	Blob        swap.SwapBlob
	SendCost    blockchain.Cost
	ReceiveCost blockchain.Cost
}

func (msg SwapRequest) IsMessage() {
}

func NewSwapRequest(blob swap.SwapBlob, sendCost, receiveCost blockchain.Cost) SwapRequest {
	return SwapRequest{
		Blob:        blob,
		SendCost:    sendCost,
		ReceiveCost: receiveCost,
	}
}

type ReceiptUpdate swap.ReceiptUpdate

func (msg ReceiptUpdate) IsMessage() {
}

func NewReceiptUpdate(id swap.SwapID, status int, native, foreign Contract) ReceiptUpdate {
	return ReceiptUpdate(swap.NewReceiptUpdate(id, func(receipt *swap.SwapReceipt) {
		receipt.Status = status
		receipt.SendCost = blockchain.CostToCostBlob(native.Cost())
		receipt.ReceiveCost = blockchain.CostToCostBlob(foreign.Cost())
	}))
}

type DeleteSwap struct {
	ID swap.SwapID
}

func (msg DeleteSwap) IsMessage() {
}
