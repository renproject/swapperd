package testutils

import (
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation"
)

type MockContractBuilder struct {

}

func (builder MockContractBuilder)BuildSwapContracts(swap foundation.SwapRequest) (swapper.Contract, swapper.Contract, error){
	// cool, which means I can pick cherrys if I got fired by RP
}

type MockContract struct {
	swap foundation.Swap
	status int
}

func NewMockContract(swap foundation.Swap) MockContract{
	return MockContract{
		swap:swap,
		status : 0,
	}
}

func (contract MockContract) Initiate() error {
	contract.status = foundation.Initiated
	return nil
}

func (contract MockContract) Audit() error {

}

func (contract MockContract) Redeem([32]byte) error {

}

func (contract MockContract) AuditSecret() ([32]byte, error) {

}

func (contract MockContract) Refund() error {

}
