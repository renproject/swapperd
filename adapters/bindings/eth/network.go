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

// NetworkABI is the input ABI used to generate the binding from.
const NetworkABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"},{\"name\":\"_swapDetails\",\"type\":\"bytes\"}],\"name\":\"submitDetails\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"swapDetails\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// NetworkBin is the compiled bytecode used for deploying new contracts.
const NetworkBin = `0x608060405234801561001057600080fd5b506102c0806100206000396000f30060806040526004361061004b5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166380ee79d48114610050578063b14631bb146100b0575b600080fd5b34801561005c57600080fd5b5060408051602060046024803582810135601f81018590048502860185019096528585526100ae95833595369560449491939091019190819084018382808284375094975061013d9650505050505050565b005b3480156100bc57600080fd5b506100c860043561015f565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101025781810151838201526020016100ea565b50505050905090810190601f16801561012f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b600082815260208181526040909120825161015a928401906101f9565b505050565b600060208181529181526040908190208054825160026001831615610100026000190190921691909104601f8101859004850282018501909352828152929091908301828280156101f15780601f106101c6576101008083540402835291602001916101f1565b820191906000526020600020905b8154815290600101906020018083116101d457829003601f168201915b505050505081565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061023a57805160ff1916838001178555610267565b82800160010185558215610267579182015b8281111561026757825182559160200191906001019061024c565b50610273929150610277565b5090565b61029191905b80821115610273576000815560010161027d565b905600a165627a7a72305820316e7382e66ea72c5ac1f5a919f77e41c826c12d8990bb4b60a995b13ba08aa20029`

// DeployNetwork deploys a new Ethereum contract, binding an instance of Network to it.
func DeployNetwork(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Network, error) {
	parsed, err := abi.JSON(strings.NewReader(NetworkABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(NetworkBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Network{NetworkCaller: NetworkCaller{contract: contract}, NetworkTransactor: NetworkTransactor{contract: contract}, NetworkFilterer: NetworkFilterer{contract: contract}}, nil
}

// Network is an auto generated Go binding around an Ethereum contract.
type Network struct {
	NetworkCaller     // Read-only binding to the contract
	NetworkTransactor // Write-only binding to the contract
	NetworkFilterer   // Log filterer for contract events
}

// NetworkCaller is an auto generated read-only Go binding around an Ethereum contract.
type NetworkCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NetworkTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NetworkFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NetworkSession struct {
	Contract     *Network          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NetworkCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NetworkCallerSession struct {
	Contract *NetworkCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// NetworkTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NetworkTransactorSession struct {
	Contract     *NetworkTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// NetworkRaw is an auto generated low-level Go binding around an Ethereum contract.
type NetworkRaw struct {
	Contract *Network // Generic contract binding to access the raw methods on
}

// NetworkCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NetworkCallerRaw struct {
	Contract *NetworkCaller // Generic read-only contract binding to access the raw methods on
}

// NetworkTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NetworkTransactorRaw struct {
	Contract *NetworkTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNetwork creates a new instance of Network, bound to a specific deployed contract.
func NewNetwork(address common.Address, backend bind.ContractBackend) (*Network, error) {
	contract, err := bindNetwork(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Network{NetworkCaller: NetworkCaller{contract: contract}, NetworkTransactor: NetworkTransactor{contract: contract}, NetworkFilterer: NetworkFilterer{contract: contract}}, nil
}

// NewNetworkCaller creates a new read-only instance of Network, bound to a specific deployed contract.
func NewNetworkCaller(address common.Address, caller bind.ContractCaller) (*NetworkCaller, error) {
	contract, err := bindNetwork(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkCaller{contract: contract}, nil
}

// NewNetworkTransactor creates a new write-only instance of Network, bound to a specific deployed contract.
func NewNetworkTransactor(address common.Address, transactor bind.ContractTransactor) (*NetworkTransactor, error) {
	contract, err := bindNetwork(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkTransactor{contract: contract}, nil
}

// NewNetworkFilterer creates a new log filterer instance of Network, bound to a specific deployed contract.
func NewNetworkFilterer(address common.Address, filterer bind.ContractFilterer) (*NetworkFilterer, error) {
	contract, err := bindNetwork(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NetworkFilterer{contract: contract}, nil
}

// bindNetwork binds a generic wrapper to an already deployed contract.
func bindNetwork(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NetworkABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Network *NetworkRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Network.Contract.NetworkCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Network *NetworkRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Network.Contract.NetworkTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Network *NetworkRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Network.Contract.NetworkTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Network *NetworkCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Network.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Network *NetworkTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Network.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Network *NetworkTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Network.Contract.contract.Transact(opts, method, params...)
}

// SwapDetails is a free data retrieval call binding the contract method 0xb14631bb.
//
// Solidity: function swapDetails( bytes32) constant returns(bytes)
func (_Network *NetworkCaller) SwapDetails(opts *bind.CallOpts, arg0 [32]byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _Network.contract.Call(opts, out, "swapDetails", arg0)
	return *ret0, err
}

// SwapDetails is a free data retrieval call binding the contract method 0xb14631bb.
//
// Solidity: function swapDetails( bytes32) constant returns(bytes)
func (_Network *NetworkSession) SwapDetails(arg0 [32]byte) ([]byte, error) {
	return _Network.Contract.SwapDetails(&_Network.CallOpts, arg0)
}

// SwapDetails is a free data retrieval call binding the contract method 0xb14631bb.
//
// Solidity: function swapDetails( bytes32) constant returns(bytes)
func (_Network *NetworkCallerSession) SwapDetails(arg0 [32]byte) ([]byte, error) {
	return _Network.Contract.SwapDetails(&_Network.CallOpts, arg0)
}

// SubmitDetails is a paid mutator transaction binding the contract method 0x80ee79d4.
//
// Solidity: function submitDetails(_orderID bytes32, _swapDetails bytes) returns()
func (_Network *NetworkTransactor) SubmitDetails(opts *bind.TransactOpts, _orderID [32]byte, _swapDetails []byte) (*types.Transaction, error) {
	return _Network.contract.Transact(opts, "submitDetails", _orderID, _swapDetails)
}

// SubmitDetails is a paid mutator transaction binding the contract method 0x80ee79d4.
//
// Solidity: function submitDetails(_orderID bytes32, _swapDetails bytes) returns()
func (_Network *NetworkSession) SubmitDetails(_orderID [32]byte, _swapDetails []byte) (*types.Transaction, error) {
	return _Network.Contract.SubmitDetails(&_Network.TransactOpts, _orderID, _swapDetails)
}

// SubmitDetails is a paid mutator transaction binding the contract method 0x80ee79d4.
//
// Solidity: function submitDetails(_orderID bytes32, _swapDetails bytes) returns()
func (_Network *NetworkTransactorSession) SubmitDetails(_orderID [32]byte, _swapDetails []byte) (*types.Transaction, error) {
	return _Network.Contract.SubmitDetails(&_Network.TransactOpts, _orderID, _swapDetails)
}
