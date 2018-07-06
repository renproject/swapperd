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

// AtomicInfoABI is the input ABI used to generate the binding from.
const AtomicInfoABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"getOwnerAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"},{\"name\":\"_swapDetails\",\"type\":\"bytes\"}],\"name\":\"submitDetails\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"},{\"name\":\"_owner\",\"type\":\"bytes\"}],\"name\":\"setOwnerAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"swapDetails\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// AtomicInfoBin is the compiled bytecode used for deploying new contracts.
const AtomicInfoBin = `0x608060405234801561001057600080fd5b506103d1806100206000396000f3006080604052600436106100605763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041662d43ec6811461006557806380ee79d4146100f2578063951065c114610152578063b14631bb146101b0575b600080fd5b34801561007157600080fd5b5061007d6004356101c8565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100b757818101518382015260200161009f565b50505050905090810190601f1680156100e45780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156100fe57600080fd5b5060408051602060046024803582810135601f81018590048502860185019096528585526101509583359536956044949193909101919081908401838280828437509497506102629650505050505050565b005b34801561015e57600080fd5b5060408051602060046024803582810135601f81018590048502860185019096528585526101509583359536956044949193909101919081908401838280828437509497506102869650505050505050565b3480156101bc57600080fd5b5061007d6004356102a3565b600060208181529181526040908190208054825160026001831615610100026000190190921691909104601f81018590048502820185019093528281529290919083018282801561025a5780601f1061022f5761010080835404028352916020019161025a565b820191906000526020600020905b81548152906001019060200180831161023d57829003601f168201915b505050505081565b600082815260016020908152604090912082516102819284019061030a565b505050565b60008281526020818152604090912082516102819284019061030a565b60016020818152600092835260409283902080548451600294821615610100026000190190911693909304601f810183900483028401830190945283835291929083018282801561025a5780601f1061022f5761010080835404028352916020019161025a565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061034b57805160ff1916838001178555610378565b82800160010185558215610378579182015b8281111561037857825182559160200191906001019061035d565b50610384929150610388565b5090565b6103a291905b80821115610384576000815560010161038e565b905600a165627a7a723058209fe94c54f34920fe095610c388786f5cad5ca9d944ac81e9554213872bfb20b80029`

// DeployAtomicInfo deploys a new Ethereum contract, binding an instance of AtomicInfo to it.
func DeployAtomicInfo(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AtomicInfo, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomicInfoABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AtomicInfoBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AtomicInfo{AtomicInfoCaller: AtomicInfoCaller{contract: contract}, AtomicInfoTransactor: AtomicInfoTransactor{contract: contract}, AtomicInfoFilterer: AtomicInfoFilterer{contract: contract}}, nil
}

// AtomicInfo is an auto generated Go binding around an Ethereum contract.
type AtomicInfo struct {
	AtomicInfoCaller     // Read-only binding to the contract
	AtomicInfoTransactor // Write-only binding to the contract
	AtomicInfoFilterer   // Log filterer for contract events
}

// AtomicInfoCaller is an auto generated read-only Go binding around an Ethereum contract.
type AtomicInfoCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomicInfoTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AtomicInfoTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomicInfoFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AtomicInfoFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomicInfoSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AtomicInfoSession struct {
	Contract     *AtomicInfo       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AtomicInfoCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AtomicInfoCallerSession struct {
	Contract *AtomicInfoCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// AtomicInfoTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AtomicInfoTransactorSession struct {
	Contract     *AtomicInfoTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// AtomicInfoRaw is an auto generated low-level Go binding around an Ethereum contract.
type AtomicInfoRaw struct {
	Contract *AtomicInfo // Generic contract binding to access the raw methods on
}

// AtomicInfoCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AtomicInfoCallerRaw struct {
	Contract *AtomicInfoCaller // Generic read-only contract binding to access the raw methods on
}

// AtomicInfoTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AtomicInfoTransactorRaw struct {
	Contract *AtomicInfoTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAtomicInfo creates a new instance of AtomicInfo, bound to a specific deployed contract.
func NewAtomicInfo(address common.Address, backend bind.ContractBackend) (*AtomicInfo, error) {
	contract, err := bindAtomicInfo(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AtomicInfo{AtomicInfoCaller: AtomicInfoCaller{contract: contract}, AtomicInfoTransactor: AtomicInfoTransactor{contract: contract}, AtomicInfoFilterer: AtomicInfoFilterer{contract: contract}}, nil
}

// NewAtomicInfoCaller creates a new read-only instance of AtomicInfo, bound to a specific deployed contract.
func NewAtomicInfoCaller(address common.Address, caller bind.ContractCaller) (*AtomicInfoCaller, error) {
	contract, err := bindAtomicInfo(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AtomicInfoCaller{contract: contract}, nil
}

// NewAtomicInfoTransactor creates a new write-only instance of AtomicInfo, bound to a specific deployed contract.
func NewAtomicInfoTransactor(address common.Address, transactor bind.ContractTransactor) (*AtomicInfoTransactor, error) {
	contract, err := bindAtomicInfo(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AtomicInfoTransactor{contract: contract}, nil
}

// NewAtomicInfoFilterer creates a new log filterer instance of AtomicInfo, bound to a specific deployed contract.
func NewAtomicInfoFilterer(address common.Address, filterer bind.ContractFilterer) (*AtomicInfoFilterer, error) {
	contract, err := bindAtomicInfo(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AtomicInfoFilterer{contract: contract}, nil
}

// bindAtomicInfo binds a generic wrapper to an already deployed contract.
func bindAtomicInfo(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomicInfoABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomicInfo *AtomicInfoRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomicInfo.Contract.AtomicInfoCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomicInfo *AtomicInfoRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomicInfo.Contract.AtomicInfoTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomicInfo *AtomicInfoRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomicInfo.Contract.AtomicInfoTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomicInfo *AtomicInfoCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomicInfo.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomicInfo *AtomicInfoTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomicInfo.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomicInfo *AtomicInfoTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomicInfo.Contract.contract.Transact(opts, method, params...)
}

// GetOwnerAddress is a free data retrieval call binding the contract method 0x00d43ec6.
//
// Solidity: function getOwnerAddress( bytes32) constant returns(bytes)
func (_AtomicInfo *AtomicInfoCaller) GetOwnerAddress(opts *bind.CallOpts, arg0 [32]byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _AtomicInfo.contract.Call(opts, out, "getOwnerAddress", arg0)
	return *ret0, err
}

// GetOwnerAddress is a free data retrieval call binding the contract method 0x00d43ec6.
//
// Solidity: function getOwnerAddress( bytes32) constant returns(bytes)
func (_AtomicInfo *AtomicInfoSession) GetOwnerAddress(arg0 [32]byte) ([]byte, error) {
	return _AtomicInfo.Contract.GetOwnerAddress(&_AtomicInfo.CallOpts, arg0)
}

// GetOwnerAddress is a free data retrieval call binding the contract method 0x00d43ec6.
//
// Solidity: function getOwnerAddress( bytes32) constant returns(bytes)
func (_AtomicInfo *AtomicInfoCallerSession) GetOwnerAddress(arg0 [32]byte) ([]byte, error) {
	return _AtomicInfo.Contract.GetOwnerAddress(&_AtomicInfo.CallOpts, arg0)
}

// SwapDetails is a free data retrieval call binding the contract method 0xb14631bb.
//
// Solidity: function swapDetails( bytes32) constant returns(bytes)
func (_AtomicInfo *AtomicInfoCaller) SwapDetails(opts *bind.CallOpts, arg0 [32]byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _AtomicInfo.contract.Call(opts, out, "swapDetails", arg0)
	return *ret0, err
}

// SwapDetails is a free data retrieval call binding the contract method 0xb14631bb.
//
// Solidity: function swapDetails( bytes32) constant returns(bytes)
func (_AtomicInfo *AtomicInfoSession) SwapDetails(arg0 [32]byte) ([]byte, error) {
	return _AtomicInfo.Contract.SwapDetails(&_AtomicInfo.CallOpts, arg0)
}

// SwapDetails is a free data retrieval call binding the contract method 0xb14631bb.
//
// Solidity: function swapDetails( bytes32) constant returns(bytes)
func (_AtomicInfo *AtomicInfoCallerSession) SwapDetails(arg0 [32]byte) ([]byte, error) {
	return _AtomicInfo.Contract.SwapDetails(&_AtomicInfo.CallOpts, arg0)
}

// SetOwnerAddress is a paid mutator transaction binding the contract method 0x951065c1.
//
// Solidity: function setOwnerAddress(_orderID bytes32, _owner bytes) returns()
func (_AtomicInfo *AtomicInfoTransactor) SetOwnerAddress(opts *bind.TransactOpts, _orderID [32]byte, _owner []byte) (*types.Transaction, error) {
	return _AtomicInfo.contract.Transact(opts, "setOwnerAddress", _orderID, _owner)
}

// SetOwnerAddress is a paid mutator transaction binding the contract method 0x951065c1.
//
// Solidity: function setOwnerAddress(_orderID bytes32, _owner bytes) returns()
func (_AtomicInfo *AtomicInfoSession) SetOwnerAddress(_orderID [32]byte, _owner []byte) (*types.Transaction, error) {
	return _AtomicInfo.Contract.SetOwnerAddress(&_AtomicInfo.TransactOpts, _orderID, _owner)
}

// SetOwnerAddress is a paid mutator transaction binding the contract method 0x951065c1.
//
// Solidity: function setOwnerAddress(_orderID bytes32, _owner bytes) returns()
func (_AtomicInfo *AtomicInfoTransactorSession) SetOwnerAddress(_orderID [32]byte, _owner []byte) (*types.Transaction, error) {
	return _AtomicInfo.Contract.SetOwnerAddress(&_AtomicInfo.TransactOpts, _orderID, _owner)
}

// SubmitDetails is a paid mutator transaction binding the contract method 0x80ee79d4.
//
// Solidity: function submitDetails(_orderID bytes32, _swapDetails bytes) returns()
func (_AtomicInfo *AtomicInfoTransactor) SubmitDetails(opts *bind.TransactOpts, _orderID [32]byte, _swapDetails []byte) (*types.Transaction, error) {
	return _AtomicInfo.contract.Transact(opts, "submitDetails", _orderID, _swapDetails)
}

// SubmitDetails is a paid mutator transaction binding the contract method 0x80ee79d4.
//
// Solidity: function submitDetails(_orderID bytes32, _swapDetails bytes) returns()
func (_AtomicInfo *AtomicInfoSession) SubmitDetails(_orderID [32]byte, _swapDetails []byte) (*types.Transaction, error) {
	return _AtomicInfo.Contract.SubmitDetails(&_AtomicInfo.TransactOpts, _orderID, _swapDetails)
}

// SubmitDetails is a paid mutator transaction binding the contract method 0x80ee79d4.
//
// Solidity: function submitDetails(_orderID bytes32, _swapDetails bytes) returns()
func (_AtomicInfo *AtomicInfoTransactorSession) SubmitDetails(_orderID [32]byte, _swapDetails []byte) (*types.Transaction, error) {
	return _AtomicInfo.Contract.SubmitDetails(&_AtomicInfo.TransactOpts, _orderID, _swapDetails)
}
