package swapper

import (
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/republicprotocol/tau"
	"golang.org/x/crypto/sha3"
)

func (swapper *swapper) handleRetry() tau.Message {
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
		return swapper.handleResult(req, swap.Refunded, native, foreign, err, true)
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
			return swapper.handleResult(req, swap.AuditFailed, native, foreign, err, true)
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
		return swapper.handleResult(req, swap.Refunded, native, foreign, err, true)
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
