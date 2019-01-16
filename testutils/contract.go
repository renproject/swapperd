package testutils

import (
	"fmt"
	"math/rand"

	"github.com/republicprotocol/swapperd/core/swapper/immediate"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
)

type MockContractBuilder struct {
}

func NewMockContractBuilder() *MockContractBuilder {
	return &MockContractBuilder{}
}

func (builder MockContractBuilder) BuildSwapContracts(request immediate.SwapRequest) (immediate.Contract, immediate.Contract, error) {
	return NewMockContract(request.Blob, request.SendCost), NewMockContract(request.Blob, request.ReceiveCost), nil
}

type MockContract struct {
	rand   *rand.Rand
	blob   swap.SwapBlob
	cost   blockchain.Cost
	status int
}

func NewMockContract(blob swap.SwapBlob, cost blockchain.Cost) *MockContract {
	return &MockContract{
		blob:   blob,
		cost:   cost,
		status: 0,
	}
}

func (contract *MockContract) Initiate() error {
	contract.status = swap.Initiated
	if contract.blob.TimeLock%8 == 0 {
		return fmt.Errorf("Connection Failed")
	}
	return nil
}

func (contract *MockContract) Audit() error {
	contract.status = swap.Audited
	if contract.blob.TimeLock%8 == 1 {
		return immediate.ErrSwapExpired
	}
	if contract.blob.TimeLock%8 == 2 {
		return immediate.ErrAuditPending
	}
	if contract.blob.TimeLock%8 == 3 {
		return fmt.Errorf("Connection Failed")
	}
	return nil
}

func (contract *MockContract) Redeem([32]byte) error {
	contract.status = swap.Redeemed
	if contract.blob.TimeLock%8 == 4 {
		return fmt.Errorf("Connection Failed")
	}
	return nil
}

func (contract *MockContract) AuditSecret() ([32]byte, error) {
	if contract.blob.TimeLock%8 == 5 {
		return [32]byte{}, fmt.Errorf("Connection Failed")
	}
	if contract.blob.TimeLock%8 == 6 {
		return [32]byte{}, immediate.ErrAuditPending
	}
	if contract.blob.TimeLock%8 == 7 {
		return [32]byte{}, immediate.ErrSwapExpired
	}
	return [32]byte{}, nil
}

func (contract *MockContract) Refund() error {
	contract.status = swap.Refunded
	if contract.blob.TimeLock%16 == 7 || contract.blob.TimeLock%16 == 1 {
		return fmt.Errorf("Connection Failed")
	}
	return nil
}

func (contract *MockContract) Cost() blockchain.Cost {
	return contract.cost
}
