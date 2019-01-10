package swapper

import (
	"fmt"

	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/republicprotocol/tau"
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
		return swapper.handleRetry()
	case SwapRequest:
		return swapper.handleSwap(msg)
	default:
		return tau.NewError(fmt.Errorf("invalid message type in swapper: %T", msg))
	}
}
