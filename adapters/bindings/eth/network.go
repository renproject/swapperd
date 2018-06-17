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

// AtomNetworkABI is the input ABI used to generate the binding from.
const AtomNetworkABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"},{\"name\":\"_swapDetails\",\"type\":\"bytes\"}],\"name\":\"submitDetails\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"swapDetails\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// AtomNetworkBin is the compiled bytecode used for deploying new contracts.
const AtomNetworkBin = `0x608060405234801561001057600080fd5b506102c0806100206000396000f30060806040526004361061004b5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166380ee79d48114610050578063b14631bb146100b0575b600080fd5b34801561005c57600080fd5b5060408051602060046024803582810135601f81018590048502860185019096528585526100ae95833595369560449491939091019190819084018382808284375094975061013d9650505050505050565b005b3480156100bc57600080fd5b506100c860043561015f565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101025781810151838201526020016100ea565b50505050905090810190601f16801561012f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b600082815260208181526040909120825161015a928401906101f9565b505050565b600060208181529181526040908190208054825160026001831615610100026000190190921691909104601f8101859004850282018501909352828152929091908301828280156101f15780601f106101c6576101008083540402835291602001916101f1565b820191906000526020600020905b8154815290600101906020018083116101d457829003601f168201915b505050505081565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061023a57805160ff1916838001178555610267565b82800160010185558215610267579182015b8281111561026757825182559160200191906001019061024c565b50610273929150610277565b5090565b61029191905b80821115610273576000815560010161027d565b905600a165627a7a723058204de1123a16177b8620b558301590125ac0c8342cd6efe0e5956566ce99bb06070029`

// DeployAtomNetwork deploys a new Ethereum contract, binding an instance of AtomNetwork to it.
func DeployAtomNetwork(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AtomNetwork, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomNetworkABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AtomNetworkBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AtomNetwork{AtomNetworkCaller: AtomNetworkCaller{contract: contract}, AtomNetworkTransactor: AtomNetworkTransactor{contract: contract}, AtomNetworkFilterer: AtomNetworkFilterer{contract: contract}}, nil
}

// AtomNetwork is an auto generated Go binding around an Ethereum contract.
type AtomNetwork struct {
	AtomNetworkCaller     // Read-only binding to the contract
	AtomNetworkTransactor // Write-only binding to the contract
	AtomNetworkFilterer   // Log filterer for contract events
}

// AtomNetworkCaller is an auto generated read-only Go binding around an Ethereum contract.
type AtomNetworkCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomNetworkTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AtomNetworkTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomNetworkFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AtomNetworkFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomNetworkSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AtomNetworkSession struct {
	Contract     *AtomNetwork      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AtomNetworkCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AtomNetworkCallerSession struct {
	Contract *AtomNetworkCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// AtomNetworkTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AtomNetworkTransactorSession struct {
	Contract     *AtomNetworkTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// AtomNetworkRaw is an auto generated low-level Go binding around an Ethereum contract.
type AtomNetworkRaw struct {
	Contract *AtomNetwork // Generic contract binding to access the raw methods on
}

// AtomNetworkCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AtomNetworkCallerRaw struct {
	Contract *AtomNetworkCaller // Generic read-only contract binding to access the raw methods on
}

// AtomNetworkTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AtomNetworkTransactorRaw struct {
	Contract *AtomNetworkTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAtomNetwork creates a new instance of AtomNetwork, bound to a specific deployed contract.
func NewAtomNetwork(address common.Address, backend bind.ContractBackend) (*AtomNetwork, error) {
	contract, err := bindAtomNetwork(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AtomNetwork{AtomNetworkCaller: AtomNetworkCaller{contract: contract}, AtomNetworkTransactor: AtomNetworkTransactor{contract: contract}, AtomNetworkFilterer: AtomNetworkFilterer{contract: contract}}, nil
}

// NewAtomNetworkCaller creates a new read-only instance of AtomNetwork, bound to a specific deployed contract.
func NewAtomNetworkCaller(address common.Address, caller bind.ContractCaller) (*AtomNetworkCaller, error) {
	contract, err := bindAtomNetwork(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AtomNetworkCaller{contract: contract}, nil
}

// NewAtomNetworkTransactor creates a new write-only instance of AtomNetwork, bound to a specific deployed contract.
func NewAtomNetworkTransactor(address common.Address, transactor bind.ContractTransactor) (*AtomNetworkTransactor, error) {
	contract, err := bindAtomNetwork(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AtomNetworkTransactor{contract: contract}, nil
}

// NewAtomNetworkFilterer creates a new log filterer instance of AtomNetwork, bound to a specific deployed contract.
func NewAtomNetworkFilterer(address common.Address, filterer bind.ContractFilterer) (*AtomNetworkFilterer, error) {
	contract, err := bindAtomNetwork(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AtomNetworkFilterer{contract: contract}, nil
}

// bindAtomNetwork binds a generic wrapper to an already deployed contract.
func bindAtomNetwork(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomNetworkABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomNetwork *AtomNetworkRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomNetwork.Contract.AtomNetworkCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomNetwork *AtomNetworkRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomNetwork.Contract.AtomNetworkTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomNetwork *AtomNetworkRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomNetwork.Contract.AtomNetworkTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomNetwork *AtomNetworkCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomNetwork.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomNetwork *AtomNetworkTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomNetwork.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomNetwork *AtomNetworkTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomNetwork.Contract.contract.Transact(opts, method, params...)
}

// SwapDetails is a free data retrieval call binding the contract method 0xb14631bb.
//
// Solidity: function swapDetails( bytes32) constant returns(bytes)
func (_AtomNetwork *AtomNetworkCaller) SwapDetails(opts *bind.CallOpts, arg0 [32]byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _AtomNetwork.contract.Call(opts, out, "swapDetails", arg0)
	return *ret0, err
}

// SwapDetails is a free data retrieval call binding the contract method 0xb14631bb.
//
// Solidity: function swapDetails( bytes32) constant returns(bytes)
func (_AtomNetwork *AtomNetworkSession) SwapDetails(arg0 [32]byte) ([]byte, error) {
	return _AtomNetwork.Contract.SwapDetails(&_AtomNetwork.CallOpts, arg0)
}

// SwapDetails is a free data retrieval call binding the contract method 0xb14631bb.
//
// Solidity: function swapDetails( bytes32) constant returns(bytes)
func (_AtomNetwork *AtomNetworkCallerSession) SwapDetails(arg0 [32]byte) ([]byte, error) {
	return _AtomNetwork.Contract.SwapDetails(&_AtomNetwork.CallOpts, arg0)
}

// SubmitDetails is a paid mutator transaction binding the contract method 0x80ee79d4.
//
// Solidity: function submitDetails(_orderID bytes32, _swapDetails bytes) returns()
func (_AtomNetwork *AtomNetworkTransactor) SubmitDetails(opts *bind.TransactOpts, _orderID [32]byte, _swapDetails []byte) (*types.Transaction, error) {
	return _AtomNetwork.contract.Transact(opts, "submitDetails", _orderID, _swapDetails)
}

// SubmitDetails is a paid mutator transaction binding the contract method 0x80ee79d4.
//
// Solidity: function submitDetails(_orderID bytes32, _swapDetails bytes) returns()
func (_AtomNetwork *AtomNetworkSession) SubmitDetails(_orderID [32]byte, _swapDetails []byte) (*types.Transaction, error) {
	return _AtomNetwork.Contract.SubmitDetails(&_AtomNetwork.TransactOpts, _orderID, _swapDetails)
}

// SubmitDetails is a paid mutator transaction binding the contract method 0x80ee79d4.
//
// Solidity: function submitDetails(_orderID bytes32, _swapDetails bytes) returns()
func (_AtomNetwork *AtomNetworkTransactorSession) SubmitDetails(_orderID [32]byte, _swapDetails []byte) (*types.Transaction, error) {
	return _AtomNetwork.Contract.SubmitDetails(&_AtomNetwork.TransactOpts, _orderID, _swapDetails)
}
