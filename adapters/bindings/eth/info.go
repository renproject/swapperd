// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package eth

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// AtomInfoABI is the input ABI used to generate the binding from.
const AtomInfoABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"getOwnerAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"},{\"name\":\"_owner\",\"type\":\"bytes\"}],\"name\":\"setOwnerAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// AtomInfoBin is the compiled bytecode used for deploying new contracts.
const AtomInfoBin = `0x608060405234801561001057600080fd5b506102bf806100206000396000f30060806040526004361061004a5763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041662d43ec6811461004f578063951065c1146100dc575b600080fd5b34801561005b57600080fd5b5061006760043561013c565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100a1578181015183820152602001610089565b50505050905090810190601f1680156100ce5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156100e857600080fd5b5060408051602060046024803582810135601f810185900485028601850190965285855261013a9583359536956044949193909101919081908401838280828437509497506101d69650505050505050565b005b600060208181529181526040908190208054825160026001831615610100026000190190921691909104601f8101859004850282018501909352828152929091908301828280156101ce5780601f106101a3576101008083540402835291602001916101ce565b820191906000526020600020905b8154815290600101906020018083116101b157829003601f168201915b505050505081565b60008281526020818152604090912082516101f3928401906101f8565b505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061023957805160ff1916838001178555610266565b82800160010185558215610266579182015b8281111561026657825182559160200191906001019061024b565b50610272929150610276565b5090565b61029091905b80821115610272576000815560010161027c565b905600a165627a7a723058200723537f531c166cb51f4beed2745abddc9f09584000ede213faa496b360d0ca0029`

// DeployAtomInfo deploys a new Ethereum contract, binding an instance of AtomInfo to it.
func DeployAtomInfo(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AtomInfo, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomInfoABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AtomInfoBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AtomInfo{AtomInfoCaller: AtomInfoCaller{contract: contract}, AtomInfoTransactor: AtomInfoTransactor{contract: contract}, AtomInfoFilterer: AtomInfoFilterer{contract: contract}}, nil
}

// AtomInfo is an auto generated Go binding around an Ethereum contract.
type AtomInfo struct {
	AtomInfoCaller     // Read-only binding to the contract
	AtomInfoTransactor // Write-only binding to the contract
	AtomInfoFilterer   // Log filterer for contract events
}

// AtomInfoCaller is an auto generated read-only Go binding around an Ethereum contract.
type AtomInfoCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomInfoTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AtomInfoTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomInfoFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AtomInfoFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomInfoSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AtomInfoSession struct {
	Contract     *AtomInfo         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AtomInfoCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AtomInfoCallerSession struct {
	Contract *AtomInfoCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// AtomInfoTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AtomInfoTransactorSession struct {
	Contract     *AtomInfoTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// AtomInfoRaw is an auto generated low-level Go binding around an Ethereum contract.
type AtomInfoRaw struct {
	Contract *AtomInfo // Generic contract binding to access the raw methods on
}

// AtomInfoCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AtomInfoCallerRaw struct {
	Contract *AtomInfoCaller // Generic read-only contract binding to access the raw methods on
}

// AtomInfoTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AtomInfoTransactorRaw struct {
	Contract *AtomInfoTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAtomInfo creates a new instance of AtomInfo, bound to a specific deployed contract.
func NewAtomInfo(address common.Address, backend bind.ContractBackend) (*AtomInfo, error) {
	contract, err := bindAtomInfo(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AtomInfo{AtomInfoCaller: AtomInfoCaller{contract: contract}, AtomInfoTransactor: AtomInfoTransactor{contract: contract}, AtomInfoFilterer: AtomInfoFilterer{contract: contract}}, nil
}

// NewAtomInfoCaller creates a new read-only instance of AtomInfo, bound to a specific deployed contract.
func NewAtomInfoCaller(address common.Address, caller bind.ContractCaller) (*AtomInfoCaller, error) {
	contract, err := bindAtomInfo(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AtomInfoCaller{contract: contract}, nil
}

// NewAtomInfoTransactor creates a new write-only instance of AtomInfo, bound to a specific deployed contract.
func NewAtomInfoTransactor(address common.Address, transactor bind.ContractTransactor) (*AtomInfoTransactor, error) {
	contract, err := bindAtomInfo(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AtomInfoTransactor{contract: contract}, nil
}

// NewAtomInfoFilterer creates a new log filterer instance of AtomInfo, bound to a specific deployed contract.
func NewAtomInfoFilterer(address common.Address, filterer bind.ContractFilterer) (*AtomInfoFilterer, error) {
	contract, err := bindAtomInfo(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AtomInfoFilterer{contract: contract}, nil
}

// bindAtomInfo binds a generic wrapper to an already deployed contract.
func bindAtomInfo(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomInfoABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomInfo *AtomInfoRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomInfo.Contract.AtomInfoCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomInfo *AtomInfoRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomInfo.Contract.AtomInfoTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomInfo *AtomInfoRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomInfo.Contract.AtomInfoTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomInfo *AtomInfoCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomInfo.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomInfo *AtomInfoTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomInfo.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomInfo *AtomInfoTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomInfo.Contract.contract.Transact(opts, method, params...)
}

// GetOwnerAddress is a free data retrieval call binding the contract method 0x00d43ec6.
//
// Solidity: function getOwnerAddress( bytes32) constant returns(bytes)
func (_AtomInfo *AtomInfoCaller) GetOwnerAddress(opts *bind.CallOpts, arg0 [32]byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _AtomInfo.contract.Call(opts, out, "getOwnerAddress", arg0)
	return *ret0, err
}

// GetOwnerAddress is a free data retrieval call binding the contract method 0x00d43ec6.
//
// Solidity: function getOwnerAddress( bytes32) constant returns(bytes)
func (_AtomInfo *AtomInfoSession) GetOwnerAddress(arg0 [32]byte) ([]byte, error) {
	return _AtomInfo.Contract.GetOwnerAddress(&_AtomInfo.CallOpts, arg0)
}

// GetOwnerAddress is a free data retrieval call binding the contract method 0x00d43ec6.
//
// Solidity: function getOwnerAddress( bytes32) constant returns(bytes)
func (_AtomInfo *AtomInfoCallerSession) GetOwnerAddress(arg0 [32]byte) ([]byte, error) {
	return _AtomInfo.Contract.GetOwnerAddress(&_AtomInfo.CallOpts, arg0)
}

// SetOwnerAddress is a paid mutator transaction binding the contract method 0x951065c1.
//
// Solidity: function setOwnerAddress(_orderID bytes32, _owner bytes) returns()
func (_AtomInfo *AtomInfoTransactor) SetOwnerAddress(opts *bind.TransactOpts, _orderID [32]byte, _owner []byte) (*types.Transaction, error) {
	return _AtomInfo.contract.Transact(opts, "setOwnerAddress", _orderID, _owner)
}

// SetOwnerAddress is a paid mutator transaction binding the contract method 0x951065c1.
//
// Solidity: function setOwnerAddress(_orderID bytes32, _owner bytes) returns()
func (_AtomInfo *AtomInfoSession) SetOwnerAddress(_orderID [32]byte, _owner []byte) (*types.Transaction, error) {
	return _AtomInfo.Contract.SetOwnerAddress(&_AtomInfo.TransactOpts, _orderID, _owner)
}

// SetOwnerAddress is a paid mutator transaction binding the contract method 0x951065c1.
//
// Solidity: function setOwnerAddress(_orderID bytes32, _owner bytes) returns()
func (_AtomInfo *AtomInfoTransactorSession) SetOwnerAddress(_orderID [32]byte, _owner []byte) (*types.Transaction, error) {
	return _AtomInfo.Contract.SetOwnerAddress(&_AtomInfo.TransactOpts, _orderID, _owner)
}
