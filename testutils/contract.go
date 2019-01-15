package testutils

import (
	"github.com/republicprotocol/swapperd/core/swapper/immediate"
	"github.com/republicprotocol/swapperd/foundation/swap"
)

type MockContractBuilder struct {
}

func (builder MockContractBuilder) BuildSwapContracts(swapblob swap.SwapBlob) (immediate.Contract, immediate.Contract, error) {
	// cool, which means I can pick cherrys if I got fired by RP
	return nil, nil, nil
}

type MockContract struct {
	swap   swap.SwapBlob
	status int
}

func NewMockContract(swap swap.SwapBlob) MockContract {
	return MockContract{
		swap:   swap,
		status: 0,
	}
}

func (contract MockContract) Initiate() error {
	contract.status = swap.Initiated
	return nil
}

func (contract MockContract) Audit() error {
	contract.status = swap.Audited
	return nil
}

func (contract MockContract) Redeem([32]byte) error {
	contract.status = swap.Redeemed
	return nil
}

func (contract MockContract) AuditSecret() ([32]byte, error) {
	return [32]byte{}, nil
}

func (contract MockContract) Refund() error {
	contract.status = swap.Refunded
	return nil
}
