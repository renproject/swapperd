// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package eth

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// AtomWalletABI is the input ABI used to generate the binding from.
const AtomWalletABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_buyID\",\"type\":\"bytes32\"},{\"name\":\"_sellID\",\"type\":\"bytes32\"},{\"name\":\"_buyToken\",\"type\":\"uint32\"},{\"name\":\"_sellToken\",\"type\":\"uint32\"},{\"name\":\"_buyValue\",\"type\":\"uint256\"},{\"name\":\"_sellValue\",\"type\":\"uint256\"}],\"name\":\"setSettlementDetails\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"orderID\",\"type\":\"bytes32\"}],\"name\":\"getSettlementDetails\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"uint32\"},{\"name\":\"\",\"type\":\"uint32\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// AtomWalletBin is the compiled bytecode used for deploying new contracts.
const AtomWalletBin = `0x608060405234801561001057600080fd5b5061030e806100206000396000f30060806040526004361061004b5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416634150281d81146100505780636b195d3714610082575b600080fd5b34801561005c57600080fd5b5061008060043560243563ffffffff6044358116906064351660843560a4356100d6565b005b34801561008e57600080fd5b5061009a6004356102a0565b60408051968752602087019590955263ffffffff93841686860152919092166060850152608084019190915260a0830152519081900360c00190f35b60c06040519081016040528087600019168152602001866000191681526020018563ffffffff1681526020018463ffffffff168152602001838152602001828152506000808860001916600019168152602001908152602001600020600082015181600001906000191690556020820151816001019060001916905560408201518160020160006101000a81548163ffffffff021916908363ffffffff16021790555060608201518160020160046101000a81548163ffffffff021916908363ffffffff1602179055506080820151816003015560a0820151816004015590505060c06040519081016040528086600019168152602001876000191681526020018463ffffffff1681526020018563ffffffff168152602001828152602001838152506000808760001916600019168152602001908152602001600020600082015181600001906000191690556020820151816001019060001916905560408201518160020160006101000a81548163ffffffff021916908363ffffffff16021790555060608201518160020160046101000a81548163ffffffff021916908363ffffffff1602179055506080820151816003015560a08201518160040155905050505050505050565b600090815260208190526040902080546001820154600283015460038401546004909401549294919363ffffffff8083169464010000000090930416929091905600a165627a7a723058203f79d43d1c79a7ccef7bea2aa54a5f00aace7dcb9ac2cba4f53d0c5a53a299ae0029`

// DeployAtomWallet deploys a new Ethereum contract, binding an instance of AtomWallet to it.
func DeployAtomWallet(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AtomWallet, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomWalletABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AtomWalletBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AtomWallet{AtomWalletCaller: AtomWalletCaller{contract: contract}, AtomWalletTransactor: AtomWalletTransactor{contract: contract}, AtomWalletFilterer: AtomWalletFilterer{contract: contract}}, nil
}

// AtomWallet is an auto generated Go binding around an Ethereum contract.
type AtomWallet struct {
	AtomWalletCaller     // Read-only binding to the contract
	AtomWalletTransactor // Write-only binding to the contract
	AtomWalletFilterer   // Log filterer for contract events
}

// AtomWalletCaller is an auto generated read-only Go binding around an Ethereum contract.
type AtomWalletCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomWalletTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AtomWalletTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomWalletFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AtomWalletFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomWalletSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AtomWalletSession struct {
	Contract     *AtomWallet       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AtomWalletCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AtomWalletCallerSession struct {
	Contract *AtomWalletCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// AtomWalletTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AtomWalletTransactorSession struct {
	Contract     *AtomWalletTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// AtomWalletRaw is an auto generated low-level Go binding around an Ethereum contract.
type AtomWalletRaw struct {
	Contract *AtomWallet // Generic contract binding to access the raw methods on
}

// AtomWalletCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AtomWalletCallerRaw struct {
	Contract *AtomWalletCaller // Generic read-only contract binding to access the raw methods on
}

// AtomWalletTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AtomWalletTransactorRaw struct {
	Contract *AtomWalletTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAtomWallet creates a new instance of AtomWallet, bound to a specific deployed contract.
func NewAtomWallet(address common.Address, backend bind.ContractBackend) (*AtomWallet, error) {
	contract, err := bindAtomWallet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AtomWallet{AtomWalletCaller: AtomWalletCaller{contract: contract}, AtomWalletTransactor: AtomWalletTransactor{contract: contract}, AtomWalletFilterer: AtomWalletFilterer{contract: contract}}, nil
}

// NewAtomWalletCaller creates a new read-only instance of AtomWallet, bound to a specific deployed contract.
func NewAtomWalletCaller(address common.Address, caller bind.ContractCaller) (*AtomWalletCaller, error) {
	contract, err := bindAtomWallet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AtomWalletCaller{contract: contract}, nil
}

// NewAtomWalletTransactor creates a new write-only instance of AtomWallet, bound to a specific deployed contract.
func NewAtomWalletTransactor(address common.Address, transactor bind.ContractTransactor) (*AtomWalletTransactor, error) {
	contract, err := bindAtomWallet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AtomWalletTransactor{contract: contract}, nil
}

// NewAtomWalletFilterer creates a new log filterer instance of AtomWallet, bound to a specific deployed contract.
func NewAtomWalletFilterer(address common.Address, filterer bind.ContractFilterer) (*AtomWalletFilterer, error) {
	contract, err := bindAtomWallet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AtomWalletFilterer{contract: contract}, nil
}

// bindAtomWallet binds a generic wrapper to an already deployed contract.
func bindAtomWallet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomWalletABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomWallet *AtomWalletRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomWallet.Contract.AtomWalletCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomWallet *AtomWalletRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomWallet.Contract.AtomWalletTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomWallet *AtomWalletRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomWallet.Contract.AtomWalletTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomWallet *AtomWalletCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomWallet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomWallet *AtomWalletTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomWallet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomWallet *AtomWalletTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomWallet.Contract.contract.Transact(opts, method, params...)
}

// GetSettlementDetails is a free data retrieval call binding the contract method 0x6b195d37.
//
// Solidity: function getSettlementDetails(orderID bytes32) constant returns(bytes32, bytes32, uint32, uint32, uint256, uint256)
func (_AtomWallet *AtomWalletCaller) GetSettlementDetails(opts *bind.CallOpts, orderID [32]byte) ([32]byte, [32]byte, uint32, uint32, *big.Int, *big.Int, error) {
	var (
		ret0 = new([32]byte)
		ret1 = new([32]byte)
		ret2 = new(uint32)
		ret3 = new(uint32)
		ret4 = new(*big.Int)
		ret5 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
		ret4,
		ret5,
	}
	err := _AtomWallet.contract.Call(opts, out, "getSettlementDetails", orderID)
	return *ret0, *ret1, *ret2, *ret3, *ret4, *ret5, err
}

// GetSettlementDetails is a free data retrieval call binding the contract method 0x6b195d37.
//
// Solidity: function getSettlementDetails(orderID bytes32) constant returns(bytes32, bytes32, uint32, uint32, uint256, uint256)
func (_AtomWallet *AtomWalletSession) GetSettlementDetails(orderID [32]byte) ([32]byte, [32]byte, uint32, uint32, *big.Int, *big.Int, error) {
	return _AtomWallet.Contract.GetSettlementDetails(&_AtomWallet.CallOpts, orderID)
}

// GetSettlementDetails is a free data retrieval call binding the contract method 0x6b195d37.
//
// Solidity: function getSettlementDetails(orderID bytes32) constant returns(bytes32, bytes32, uint32, uint32, uint256, uint256)
func (_AtomWallet *AtomWalletCallerSession) GetSettlementDetails(orderID [32]byte) ([32]byte, [32]byte, uint32, uint32, *big.Int, *big.Int, error) {
	return _AtomWallet.Contract.GetSettlementDetails(&_AtomWallet.CallOpts, orderID)
}

// SetSettlementDetails is a paid mutator transaction binding the contract method 0x4150281d.
//
// Solidity: function setSettlementDetails(_buyID bytes32, _sellID bytes32, _buyToken uint32, _sellToken uint32, _buyValue uint256, _sellValue uint256) returns()
func (_AtomWallet *AtomWalletTransactor) SetSettlementDetails(opts *bind.TransactOpts, _buyID [32]byte, _sellID [32]byte, _buyToken uint32, _sellToken uint32, _buyValue *big.Int, _sellValue *big.Int) (*types.Transaction, error) {
	return _AtomWallet.contract.Transact(opts, "setSettlementDetails", _buyID, _sellID, _buyToken, _sellToken, _buyValue, _sellValue)
}

// SetSettlementDetails is a paid mutator transaction binding the contract method 0x4150281d.
//
// Solidity: function setSettlementDetails(_buyID bytes32, _sellID bytes32, _buyToken uint32, _sellToken uint32, _buyValue uint256, _sellValue uint256) returns()
func (_AtomWallet *AtomWalletSession) SetSettlementDetails(_buyID [32]byte, _sellID [32]byte, _buyToken uint32, _sellToken uint32, _buyValue *big.Int, _sellValue *big.Int) (*types.Transaction, error) {
	return _AtomWallet.Contract.SetSettlementDetails(&_AtomWallet.TransactOpts, _buyID, _sellID, _buyToken, _sellToken, _buyValue, _sellValue)
}

// SetSettlementDetails is a paid mutator transaction binding the contract method 0x4150281d.
//
// Solidity: function setSettlementDetails(_buyID bytes32, _sellID bytes32, _buyToken uint32, _sellToken uint32, _buyValue uint256, _sellValue uint256) returns()
func (_AtomWallet *AtomWalletTransactorSession) SetSettlementDetails(_buyID [32]byte, _sellID [32]byte, _buyToken uint32, _sellToken uint32, _buyValue *big.Int, _sellValue *big.Int) (*types.Transaction, error) {
	return _AtomWallet.Contract.SetSettlementDetails(&_AtomWallet.TransactOpts, _buyID, _sellID, _buyToken, _sellToken, _buyValue, _sellValue)
}
