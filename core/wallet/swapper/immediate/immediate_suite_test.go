package immediate_test

import (
	"fmt"
	"math/rand"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/republicprotocol/swapperd/core/wallet/swapper/immediate"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
)

func TestImmediate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Immediate Suite")
}

type MockContractBuilder struct {
}

func NewMockContractBuilder() *MockContractBuilder {
	return &MockContractBuilder{}
}

func (builder MockContractBuilder) BuildSwapContracts(request immediate.SwapRequest) (immediate.Contract, immediate.Contract, error) {
	if uint64(request.Blob.TimeLock)%36 == 8 {
		return nil, nil, fmt.Errorf("invalid swap object")
	}
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
	if uint64(contract.blob.TimeLock)%9 == 0 {
		return fmt.Errorf("Connection Failed")
	}
	return nil
}

func (contract *MockContract) Audit() error {
	contract.status = swap.Audited
	if uint64(contract.blob.TimeLock)%9 == 1 {
		return immediate.ErrSwapExpired
	}
	if uint64(contract.blob.TimeLock)%9 == 2 {
		return immediate.ErrAuditPending
	}
	if uint64(contract.blob.TimeLock)%9 == 3 {
		return fmt.Errorf("Connection Failed")
	}
	return nil
}

func (contract *MockContract) Redeem([32]byte) error {
	contract.status = swap.Redeemed
	if uint64(contract.blob.TimeLock)%9 == 4 {
		return fmt.Errorf("Connection Failed")
	}
	return nil
}

func (contract *MockContract) AuditSecret() ([32]byte, error) {
	if uint64(contract.blob.TimeLock)%9 == 5 {
		return [32]byte{}, fmt.Errorf("Connection Failed")
	}
	if uint64(contract.blob.TimeLock)%9 == 6 {
		return [32]byte{}, immediate.ErrAuditPending
	}
	if uint64(contract.blob.TimeLock)%9 == 7 {
		return [32]byte{}, immediate.ErrSwapExpired
	}
	return [32]byte{}, nil
}

func (contract *MockContract) Refund() error {
	contract.status = swap.Refunded
	if uint64(contract.blob.TimeLock)%18 == 7 || uint64(contract.blob.TimeLock)%18 == 1 {
		return fmt.Errorf("Connection Failed")
	}
	return nil
}

func (contract *MockContract) Cost() blockchain.Cost {
	return contract.cost
}
