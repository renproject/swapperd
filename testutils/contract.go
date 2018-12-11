package testutils

//
// import (
// 	"github.com/republicprotocol/swapperd/core/swapper"
// 	"github.com/republicprotocol/swapperd/foundation/swap"
// )
//
// type MockContractBuilder struct {
// }
//
// func (builder MockContractBuilder) BuildSwapContracts(swap swap.SwapRequest) (swapper.Contract, swapper.Contract, error) {
// 	// cool, which means I can pick cherrys if I got fired by RP
// 	return nil, nil, nil
// }
//
// type MockContract struct {
// 	swap   foundation.Swap
// 	status int
// }
//
// func NewMockContract(swap foundation.Swap) MockContract {
// 	return MockContract{
// 		swap:   swap,
// 		status: 0,
// 	}
// }
//
// func (contract MockContract) Initiate() error {
// 	return nil
// }
//
// func (contract MockContract) Audit() error {
// 	return nil
// }
//
// func (contract MockContract) Redeem([32]byte) error {
// 	return nil
// }
//
// func (contract MockContract) AuditSecret() ([32]byte, error) {
// 	return [32]byte{}, nil
// }
//
// func (contract MockContract) Refund() error {
// 	return nil
// }
