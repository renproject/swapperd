// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package eth

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// EthSwapContractABI is the input ABI used to generate the binding from.
const EthSwapContractABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"initiate\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"initiatable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"}],\"name\":\"swapID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawBrokerFees\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"redeemable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"auditSecret\",\"outputs\":[{\"name\":\"secretKey\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refundable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_broker\",\"type\":\"address\"},{\"name\":\"_brokerFee\",\"type\":\"uint256\"},{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"initiateWithFees\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"redeemedAt\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"audit\",\"outputs\":[{\"name\":\"timelock\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"brokerFee\",\"type\":\"uint256\"},{\"name\":\"broker\",\"type\":\"address\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"secretLock\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_receiver\",\"type\":\"address\"},{\"name\":\"_secretKey\",\"type\":\"bytes32\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"brokerFees\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_VERSION\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_secretLock\",\"type\":\"bytes32\"}],\"name\":\"LogOpen\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"LogExpire\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_secretKey\",\"type\":\"bytes32\"}],\"name\":\"LogClose\",\"type\":\"event\"}]"

// EthSwapContractBin is the compiled bytecode used for deploying new contracts.
const EthSwapContractBin = `0x60806040523480156200001157600080fd5b50604051620011b0380380620011b0833981018060405260208110156200003757600080fd5b8101908080516401000000008111156200005057600080fd5b820160208101848111156200006457600080fd5b81516401000000008111828201871017156200007f57600080fd5b505080519093506200009b9250600091506020840190620000a3565b505062000148565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620000e657805160ff191683800117855562000116565b8280016001018555821562000116579182015b8281111562000116578251825591602001919060010190620000f9565b506200012492915062000128565b5090565b6200014591905b808211156200012457600081556001016200012f565b90565b61105880620001586000396000f3fe6080604052600436106100c9577c01000000000000000000000000000000000000000000000000000000006000350463027a257781146100ce57806309ece6181461010e5780634b2ac3fa1461014c5780634c6d37ff1461018e57806368f06b29146101b85780637249fbb6146101e2578063976d00f41461020c5780639fb3147514610236578063b8688e3f14610260578063bc4fcc4a146102ae578063c140635b146102d8578063c23b1a8514610348578063e1ec380c14610387578063ffa1ad74146103ba575b600080fd5b61010c600480360360a08110156100e457600080fd5b50803590600160a060020a036020820135169060408101359060608101359060800135610444565b005b34801561011a57600080fd5b506101386004803603602081101561013157600080fd5b5035610680565b604080519115158252519081900360200190f35b34801561015857600080fd5b5061017c6004803603604081101561016f57600080fd5b50803590602001356106a8565b60408051918252519081900360200190f35b34801561019a57600080fd5b5061010c600480360360208110156101b157600080fd5b50356106d4565b3480156101c457600080fd5b50610138600480360360208110156101db57600080fd5b5035610733565b3480156101ee57600080fd5b5061010c6004803603602081101561020557600080fd5b503561073c565b34801561021857600080fd5b5061017c6004803603602081101561022f57600080fd5b50356108bb565b34801561024257600080fd5b506101386004803603602081101561025957600080fd5b5035610949565b61010c600480360360e081101561027657600080fd5b50803590600160a060020a03602082013581169160408101359091169060608101359060808101359060a08101359060c0013561096f565b3480156102ba57600080fd5b5061017c600480360360208110156102d157600080fd5b5035610bbb565b3480156102e457600080fd5b50610302600480360360208110156102fb57600080fd5b5035610bcd565b604080519788526020880196909652600160a060020a03948516878701526060870193909352908316608086015290911660a084015260c0830152519081900360e00190f35b34801561035457600080fd5b5061010c6004803603606081101561036b57600080fd5b50803590600160a060020a036020820135169060400135610c75565b34801561039357600080fd5b5061017c600480360360208110156103aa57600080fd5b5035600160a060020a0316610f48565b3480156103c657600080fd5b506103cf610f5a565b6040805160208082528351818301528351919283929083019185019080838360005b838110156104095781810151838201526020016103f1565b50505050905090810190601f1680156104365780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b846000808281526002602052604090205460ff16600381111561046357fe5b146104b8576040805160e560020a62461bcd02815260206004820152601660248201527f73776170206f70656e65642070726576696f75736c7900000000000000000000604482015290519081900360640190fd5b3482146104c457600080fd5b6104cc610fe8565b61010060405190810160405280858152602001848152602001600081526020018681526020016000600102815260200133600160a060020a0316815260200187600160a060020a031681526020016000600160a060020a031681525090508060016000898152602001908152602001600020600082015181600001556020820151816001015560408201518160020155606082015181600301556080820151816004015560a08201518160050160006101000a815481600160a060020a030219169083600160a060020a0316021790555060c08201518160060160006101000a815481600160a060020a030219169083600160a060020a0316021790555060e08201518160070160006101000a815481600160a060020a030219169083600160a060020a0316021790555090505060016002600089815260200190815260200160002060006101000a81548160ff0219169083600381111561062a57fe5b021790555060408051888152600160a060020a038816602082015280820187905290517f497d46e9505eefe8b910d1a02e6b40d8769510023b0053c3ac4b9574b81c97bf9181900360600190a150505050505050565b6000805b60008381526002602052604090205460ff1660038111156106a157fe5b1492915050565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b336000908152600360205260409020548111156106f057600080fd5b33600081815260036020526040808220805485900390555183156108fc0291849190818181858888f1935050505015801561072f573d6000803e3d6000fd5b5050565b60006001610684565b80600160008281526002602052604090205460ff16600381111561075c57fe5b146107b1576040805160e560020a62461bcd02815260206004820152600d60248201527f73776170206e6f74206f70656e00000000000000000000000000000000000000604482015290519081900360640190fd5b6000828152600160205260409020548290421015610819576040805160e560020a62461bcd02815260206004820152601260248201527f73776170206e6f7420657870697261626c650000000000000000000000000000604482015290519081900360640190fd5b60008381526002602081815260408084208054600360ff199091161790556001918290528084206005810154938101549201549051600160a060020a0390931693910180156108fc02929091818181858888f19350505050158015610882573d6000803e3d6000fd5b506040805184815290517feb711459e1247171382f0da0387b86239b8e3ca60af3b15a9ff2f1eb3d7f6a1d9181900360200190a1505050565b600081600260008281526002602052604090205460ff1660038111156108dd57fe5b14610932576040805160e560020a62461bcd02815260206004820152601160248201527f73776170206e6f742072656465656d6564000000000000000000000000000000604482015290519081900360640190fd5b505060009081526001602052604090206004015490565b600081815260016020526040812054421080159061096957506001610684565b92915050565b866000808281526002602052604090205460ff16600381111561098e57fe5b146109e3576040805160e560020a62461bcd02815260206004820152601660248201527f73776170206f70656e65642070726576696f75736c7900000000000000000000604482015290519081900360640190fd5b34821480156109f25750848210155b15156109fd57600080fd5b610a05610fe8565b6101006040519081016040528085815260200187850381526020018781526020018681526020016000600102815260200133600160a060020a0316815260200189600160a060020a0316815260200188600160a060020a0316815250905080600160008b8152602001908152602001600020600082015181600001556020820151816001015560408201518160020155606082015181600301556080820151816004015560a08201518160050160006101000a815481600160a060020a030219169083600160a060020a0316021790555060c08201518160060160006101000a815481600160a060020a030219169083600160a060020a0316021790555060e08201518160070160006101000a815481600160a060020a030219169083600160a060020a031602179055509050506001600260008b815260200190815260200160002060006101000a81548160ff02191690836003811115610b6357fe5b0217905550604080518a8152600160a060020a038a16602082015280820187905290517f497d46e9505eefe8b910d1a02e6b40d8769510023b0053c3ac4b9574b81c97bf9181900360600190a1505050505050505050565b60046020526000908152604090205481565b6000806000806000806000610be0610fe8565b50505060009586525050600160208181526040958690208651610100810188528154808252938201549281018390526002820154978101889052600382015460608201819052600483015460808301526005830154600160a060020a0390811660a084018190526006850154821660c0850181905260079095015490911660e090930183905294999398929750919550935090565b82600160008281526002602052604090205460ff166003811115610c9557fe5b14610cea576040805160e560020a62461bcd02815260206004820152600d60248201527f73776170206e6f74206f70656e00000000000000000000000000000000000000604482015290519081900360640190fd5b8382600281604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b60208310610d3d5780518252601f199092019160209182019101610d1e565b51815160209384036101000a60001901801990921691161790526040519190930194509192505080830381855afa158015610d7c573d6000803e3d6000fd5b5050506040513d6020811015610d9157600080fd5b505160008381526001602052604090206003015414610dfa576040805160e560020a62461bcd02815260206004820152600e60248201527f696e76616c696420736563726574000000000000000000000000000000000000604482015290519081900360640190fd5b60008681526001602052604090206006015486903390600160a060020a03168114610e6f576040805160e560020a62461bcd02815260206004820152601460248201527f756e617574686f72697a6564207370656e646572000000000000000000000000604482015290519081900360640190fd5b600088815260016020818152604080842060048082018c90556002808552838720805460ff1916821790559084528286204290558101546007820154600160a060020a0390811687526003855283872080549092019091558d8652928490529092015491518a93918416926108fc81150292909190818181858888f19350505050158015610f01573d6000803e3d6000fd5b50604080518a81526020810189905281517f07da1fa25a1d885732677ce9c192cbec27051a4b69d391c9a64850f5a5112ba0929181900390910190a1505050505050505050565b60036020526000908152604090205481565b6000805460408051602060026001851615610100026000190190941693909304601f81018490048402820184019092528181529291830182828015610fe05780601f10610fb557610100808354040283529160200191610fe0565b820191906000526020600020905b815481529060010190602001808311610fc357829003601f168201915b505050505081565b6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e08101919091529056fea165627a7a7230582008358d83a236f8d96ccc5423a6ece0d39537ad954325abb637a4634d8dd146a00029`

// DeployEthSwapContract deploys a new Ethereum contract, binding an instance of EthSwapContract to it.
func DeployEthSwapContract(auth *bind.TransactOpts, backend bind.ContractBackend, _VERSION string) (common.Address, *types.Transaction, *EthSwapContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EthSwapContractABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(EthSwapContractBin), backend, _VERSION)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EthSwapContract{EthSwapContractCaller: EthSwapContractCaller{contract: contract}, EthSwapContractTransactor: EthSwapContractTransactor{contract: contract}, EthSwapContractFilterer: EthSwapContractFilterer{contract: contract}}, nil
}

// EthSwapContract is an auto generated Go binding around an Ethereum contract.
type EthSwapContract struct {
	EthSwapContractCaller     // Read-only binding to the contract
	EthSwapContractTransactor // Write-only binding to the contract
	EthSwapContractFilterer   // Log filterer for contract events
}

// EthSwapContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthSwapContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthSwapContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthSwapContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthSwapContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthSwapContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthSwapContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthSwapContractSession struct {
	Contract     *EthSwapContract  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthSwapContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthSwapContractCallerSession struct {
	Contract *EthSwapContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// EthSwapContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthSwapContractTransactorSession struct {
	Contract     *EthSwapContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// EthSwapContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthSwapContractRaw struct {
	Contract *EthSwapContract // Generic contract binding to access the raw methods on
}

// EthSwapContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthSwapContractCallerRaw struct {
	Contract *EthSwapContractCaller // Generic read-only contract binding to access the raw methods on
}

// EthSwapContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthSwapContractTransactorRaw struct {
	Contract *EthSwapContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEthSwapContract creates a new instance of EthSwapContract, bound to a specific deployed contract.
func NewEthSwapContract(address common.Address, backend bind.ContractBackend) (*EthSwapContract, error) {
	contract, err := bindEthSwapContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EthSwapContract{EthSwapContractCaller: EthSwapContractCaller{contract: contract}, EthSwapContractTransactor: EthSwapContractTransactor{contract: contract}, EthSwapContractFilterer: EthSwapContractFilterer{contract: contract}}, nil
}

// NewEthSwapContractCaller creates a new read-only instance of EthSwapContract, bound to a specific deployed contract.
func NewEthSwapContractCaller(address common.Address, caller bind.ContractCaller) (*EthSwapContractCaller, error) {
	contract, err := bindEthSwapContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthSwapContractCaller{contract: contract}, nil
}

// NewEthSwapContractTransactor creates a new write-only instance of EthSwapContract, bound to a specific deployed contract.
func NewEthSwapContractTransactor(address common.Address, transactor bind.ContractTransactor) (*EthSwapContractTransactor, error) {
	contract, err := bindEthSwapContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthSwapContractTransactor{contract: contract}, nil
}

// NewEthSwapContractFilterer creates a new log filterer instance of EthSwapContract, bound to a specific deployed contract.
func NewEthSwapContractFilterer(address common.Address, filterer bind.ContractFilterer) (*EthSwapContractFilterer, error) {
	contract, err := bindEthSwapContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthSwapContractFilterer{contract: contract}, nil
}

// bindEthSwapContract binds a generic wrapper to an already deployed contract.
func bindEthSwapContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EthSwapContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthSwapContract *EthSwapContractRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _EthSwapContract.Contract.EthSwapContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthSwapContract *EthSwapContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthSwapContract.Contract.EthSwapContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthSwapContract *EthSwapContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthSwapContract.Contract.EthSwapContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthSwapContract *EthSwapContractCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _EthSwapContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthSwapContract *EthSwapContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthSwapContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthSwapContract *EthSwapContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthSwapContract.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_EthSwapContract *EthSwapContractCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _EthSwapContract.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_EthSwapContract *EthSwapContractSession) VERSION() (string, error) {
	return _EthSwapContract.Contract.VERSION(&_EthSwapContract.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_EthSwapContract *EthSwapContractCallerSession) VERSION() (string, error) {
	return _EthSwapContract.Contract.VERSION(&_EthSwapContract.CallOpts)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, brokerFee uint256, broker address, from address, secretLock bytes32)
func (_EthSwapContract *EthSwapContractCaller) Audit(opts *bind.CallOpts, _swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	BrokerFee  *big.Int
	Broker     common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	ret := new(struct {
		Timelock   *big.Int
		Value      *big.Int
		To         common.Address
		BrokerFee  *big.Int
		Broker     common.Address
		From       common.Address
		SecretLock [32]byte
	})
	out := ret
	err := _EthSwapContract.contract.Call(opts, out, "audit", _swapID)
	return *ret, err
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, brokerFee uint256, broker address, from address, secretLock bytes32)
func (_EthSwapContract *EthSwapContractSession) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	BrokerFee  *big.Int
	Broker     common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	return _EthSwapContract.Contract.Audit(&_EthSwapContract.CallOpts, _swapID)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, brokerFee uint256, broker address, from address, secretLock bytes32)
func (_EthSwapContract *EthSwapContractCallerSession) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	BrokerFee  *big.Int
	Broker     common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	return _EthSwapContract.Contract.Audit(&_EthSwapContract.CallOpts, _swapID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_EthSwapContract *EthSwapContractCaller) AuditSecret(opts *bind.CallOpts, _swapID [32]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _EthSwapContract.contract.Call(opts, out, "auditSecret", _swapID)
	return *ret0, err
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_EthSwapContract *EthSwapContractSession) AuditSecret(_swapID [32]byte) ([32]byte, error) {
	return _EthSwapContract.Contract.AuditSecret(&_EthSwapContract.CallOpts, _swapID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_EthSwapContract *EthSwapContractCallerSession) AuditSecret(_swapID [32]byte) ([32]byte, error) {
	return _EthSwapContract.Contract.AuditSecret(&_EthSwapContract.CallOpts, _swapID)
}

// BrokerFees is a free data retrieval call binding the contract method 0xe1ec380c.
//
// Solidity: function brokerFees( address) constant returns(uint256)
func (_EthSwapContract *EthSwapContractCaller) BrokerFees(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EthSwapContract.contract.Call(opts, out, "brokerFees", arg0)
	return *ret0, err
}

// BrokerFees is a free data retrieval call binding the contract method 0xe1ec380c.
//
// Solidity: function brokerFees( address) constant returns(uint256)
func (_EthSwapContract *EthSwapContractSession) BrokerFees(arg0 common.Address) (*big.Int, error) {
	return _EthSwapContract.Contract.BrokerFees(&_EthSwapContract.CallOpts, arg0)
}

// BrokerFees is a free data retrieval call binding the contract method 0xe1ec380c.
//
// Solidity: function brokerFees( address) constant returns(uint256)
func (_EthSwapContract *EthSwapContractCallerSession) BrokerFees(arg0 common.Address) (*big.Int, error) {
	return _EthSwapContract.Contract.BrokerFees(&_EthSwapContract.CallOpts, arg0)
}

// Initiatable is a free data retrieval call binding the contract method 0x09ece618.
//
// Solidity: function initiatable(_swapID bytes32) constant returns(bool)
func (_EthSwapContract *EthSwapContractCaller) Initiatable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _EthSwapContract.contract.Call(opts, out, "initiatable", _swapID)
	return *ret0, err
}

// Initiatable is a free data retrieval call binding the contract method 0x09ece618.
//
// Solidity: function initiatable(_swapID bytes32) constant returns(bool)
func (_EthSwapContract *EthSwapContractSession) Initiatable(_swapID [32]byte) (bool, error) {
	return _EthSwapContract.Contract.Initiatable(&_EthSwapContract.CallOpts, _swapID)
}

// Initiatable is a free data retrieval call binding the contract method 0x09ece618.
//
// Solidity: function initiatable(_swapID bytes32) constant returns(bool)
func (_EthSwapContract *EthSwapContractCallerSession) Initiatable(_swapID [32]byte) (bool, error) {
	return _EthSwapContract.Contract.Initiatable(&_EthSwapContract.CallOpts, _swapID)
}

// Redeemable is a free data retrieval call binding the contract method 0x68f06b29.
//
// Solidity: function redeemable(_swapID bytes32) constant returns(bool)
func (_EthSwapContract *EthSwapContractCaller) Redeemable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _EthSwapContract.contract.Call(opts, out, "redeemable", _swapID)
	return *ret0, err
}

// Redeemable is a free data retrieval call binding the contract method 0x68f06b29.
//
// Solidity: function redeemable(_swapID bytes32) constant returns(bool)
func (_EthSwapContract *EthSwapContractSession) Redeemable(_swapID [32]byte) (bool, error) {
	return _EthSwapContract.Contract.Redeemable(&_EthSwapContract.CallOpts, _swapID)
}

// Redeemable is a free data retrieval call binding the contract method 0x68f06b29.
//
// Solidity: function redeemable(_swapID bytes32) constant returns(bool)
func (_EthSwapContract *EthSwapContractCallerSession) Redeemable(_swapID [32]byte) (bool, error) {
	return _EthSwapContract.Contract.Redeemable(&_EthSwapContract.CallOpts, _swapID)
}

// RedeemedAt is a free data retrieval call binding the contract method 0xbc4fcc4a.
//
// Solidity: function redeemedAt( bytes32) constant returns(uint256)
func (_EthSwapContract *EthSwapContractCaller) RedeemedAt(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _EthSwapContract.contract.Call(opts, out, "redeemedAt", arg0)
	return *ret0, err
}

// RedeemedAt is a free data retrieval call binding the contract method 0xbc4fcc4a.
//
// Solidity: function redeemedAt( bytes32) constant returns(uint256)
func (_EthSwapContract *EthSwapContractSession) RedeemedAt(arg0 [32]byte) (*big.Int, error) {
	return _EthSwapContract.Contract.RedeemedAt(&_EthSwapContract.CallOpts, arg0)
}

// RedeemedAt is a free data retrieval call binding the contract method 0xbc4fcc4a.
//
// Solidity: function redeemedAt( bytes32) constant returns(uint256)
func (_EthSwapContract *EthSwapContractCallerSession) RedeemedAt(arg0 [32]byte) (*big.Int, error) {
	return _EthSwapContract.Contract.RedeemedAt(&_EthSwapContract.CallOpts, arg0)
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(_swapID bytes32) constant returns(bool)
func (_EthSwapContract *EthSwapContractCaller) Refundable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _EthSwapContract.contract.Call(opts, out, "refundable", _swapID)
	return *ret0, err
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(_swapID bytes32) constant returns(bool)
func (_EthSwapContract *EthSwapContractSession) Refundable(_swapID [32]byte) (bool, error) {
	return _EthSwapContract.Contract.Refundable(&_EthSwapContract.CallOpts, _swapID)
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(_swapID bytes32) constant returns(bool)
func (_EthSwapContract *EthSwapContractCallerSession) Refundable(_swapID [32]byte) (bool, error) {
	return _EthSwapContract.Contract.Refundable(&_EthSwapContract.CallOpts, _swapID)
}

// SwapID is a free data retrieval call binding the contract method 0x4b2ac3fa.
//
// Solidity: function swapID(_secretLock bytes32, _timelock uint256) constant returns(bytes32)
func (_EthSwapContract *EthSwapContractCaller) SwapID(opts *bind.CallOpts, _secretLock [32]byte, _timelock *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _EthSwapContract.contract.Call(opts, out, "swapID", _secretLock, _timelock)
	return *ret0, err
}

// SwapID is a free data retrieval call binding the contract method 0x4b2ac3fa.
//
// Solidity: function swapID(_secretLock bytes32, _timelock uint256) constant returns(bytes32)
func (_EthSwapContract *EthSwapContractSession) SwapID(_secretLock [32]byte, _timelock *big.Int) ([32]byte, error) {
	return _EthSwapContract.Contract.SwapID(&_EthSwapContract.CallOpts, _secretLock, _timelock)
}

// SwapID is a free data retrieval call binding the contract method 0x4b2ac3fa.
//
// Solidity: function swapID(_secretLock bytes32, _timelock uint256) constant returns(bytes32)
func (_EthSwapContract *EthSwapContractCallerSession) SwapID(_secretLock [32]byte, _timelock *big.Int) ([32]byte, error) {
	return _EthSwapContract.Contract.SwapID(&_EthSwapContract.CallOpts, _secretLock, _timelock)
}

// Initiate is a paid mutator transaction binding the contract method 0x027a2577.
//
// Solidity: function initiate(_swapID bytes32, _spender address, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_EthSwapContract *EthSwapContractTransactor) Initiate(opts *bind.TransactOpts, _swapID [32]byte, _spender common.Address, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _EthSwapContract.contract.Transact(opts, "initiate", _swapID, _spender, _secretLock, _timelock, _value)
}

// Initiate is a paid mutator transaction binding the contract method 0x027a2577.
//
// Solidity: function initiate(_swapID bytes32, _spender address, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_EthSwapContract *EthSwapContractSession) Initiate(_swapID [32]byte, _spender common.Address, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _EthSwapContract.Contract.Initiate(&_EthSwapContract.TransactOpts, _swapID, _spender, _secretLock, _timelock, _value)
}

// Initiate is a paid mutator transaction binding the contract method 0x027a2577.
//
// Solidity: function initiate(_swapID bytes32, _spender address, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_EthSwapContract *EthSwapContractTransactorSession) Initiate(_swapID [32]byte, _spender common.Address, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _EthSwapContract.Contract.Initiate(&_EthSwapContract.TransactOpts, _swapID, _spender, _secretLock, _timelock, _value)
}

// InitiateWithFees is a paid mutator transaction binding the contract method 0xb8688e3f.
//
// Solidity: function initiateWithFees(_swapID bytes32, _spender address, _broker address, _brokerFee uint256, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_EthSwapContract *EthSwapContractTransactor) InitiateWithFees(opts *bind.TransactOpts, _swapID [32]byte, _spender common.Address, _broker common.Address, _brokerFee *big.Int, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _EthSwapContract.contract.Transact(opts, "initiateWithFees", _swapID, _spender, _broker, _brokerFee, _secretLock, _timelock, _value)
}

// InitiateWithFees is a paid mutator transaction binding the contract method 0xb8688e3f.
//
// Solidity: function initiateWithFees(_swapID bytes32, _spender address, _broker address, _brokerFee uint256, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_EthSwapContract *EthSwapContractSession) InitiateWithFees(_swapID [32]byte, _spender common.Address, _broker common.Address, _brokerFee *big.Int, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _EthSwapContract.Contract.InitiateWithFees(&_EthSwapContract.TransactOpts, _swapID, _spender, _broker, _brokerFee, _secretLock, _timelock, _value)
}

// InitiateWithFees is a paid mutator transaction binding the contract method 0xb8688e3f.
//
// Solidity: function initiateWithFees(_swapID bytes32, _spender address, _broker address, _brokerFee uint256, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_EthSwapContract *EthSwapContractTransactorSession) InitiateWithFees(_swapID [32]byte, _spender common.Address, _broker common.Address, _brokerFee *big.Int, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _EthSwapContract.Contract.InitiateWithFees(&_EthSwapContract.TransactOpts, _swapID, _spender, _broker, _brokerFee, _secretLock, _timelock, _value)
}

// Redeem is a paid mutator transaction binding the contract method 0xc23b1a85.
//
// Solidity: function redeem(_swapID bytes32, _receiver address, _secretKey bytes32) returns()
func (_EthSwapContract *EthSwapContractTransactor) Redeem(opts *bind.TransactOpts, _swapID [32]byte, _receiver common.Address, _secretKey [32]byte) (*types.Transaction, error) {
	return _EthSwapContract.contract.Transact(opts, "redeem", _swapID, _receiver, _secretKey)
}

// Redeem is a paid mutator transaction binding the contract method 0xc23b1a85.
//
// Solidity: function redeem(_swapID bytes32, _receiver address, _secretKey bytes32) returns()
func (_EthSwapContract *EthSwapContractSession) Redeem(_swapID [32]byte, _receiver common.Address, _secretKey [32]byte) (*types.Transaction, error) {
	return _EthSwapContract.Contract.Redeem(&_EthSwapContract.TransactOpts, _swapID, _receiver, _secretKey)
}

// Redeem is a paid mutator transaction binding the contract method 0xc23b1a85.
//
// Solidity: function redeem(_swapID bytes32, _receiver address, _secretKey bytes32) returns()
func (_EthSwapContract *EthSwapContractTransactorSession) Redeem(_swapID [32]byte, _receiver common.Address, _secretKey [32]byte) (*types.Transaction, error) {
	return _EthSwapContract.Contract.Redeem(&_EthSwapContract.TransactOpts, _swapID, _receiver, _secretKey)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_EthSwapContract *EthSwapContractTransactor) Refund(opts *bind.TransactOpts, _swapID [32]byte) (*types.Transaction, error) {
	return _EthSwapContract.contract.Transact(opts, "refund", _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_EthSwapContract *EthSwapContractSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _EthSwapContract.Contract.Refund(&_EthSwapContract.TransactOpts, _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_EthSwapContract *EthSwapContractTransactorSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _EthSwapContract.Contract.Refund(&_EthSwapContract.TransactOpts, _swapID)
}

// WithdrawBrokerFees is a paid mutator transaction binding the contract method 0x4c6d37ff.
//
// Solidity: function withdrawBrokerFees(_amount uint256) returns()
func (_EthSwapContract *EthSwapContractTransactor) WithdrawBrokerFees(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _EthSwapContract.contract.Transact(opts, "withdrawBrokerFees", _amount)
}

// WithdrawBrokerFees is a paid mutator transaction binding the contract method 0x4c6d37ff.
//
// Solidity: function withdrawBrokerFees(_amount uint256) returns()
func (_EthSwapContract *EthSwapContractSession) WithdrawBrokerFees(_amount *big.Int) (*types.Transaction, error) {
	return _EthSwapContract.Contract.WithdrawBrokerFees(&_EthSwapContract.TransactOpts, _amount)
}

// WithdrawBrokerFees is a paid mutator transaction binding the contract method 0x4c6d37ff.
//
// Solidity: function withdrawBrokerFees(_amount uint256) returns()
func (_EthSwapContract *EthSwapContractTransactorSession) WithdrawBrokerFees(_amount *big.Int) (*types.Transaction, error) {
	return _EthSwapContract.Contract.WithdrawBrokerFees(&_EthSwapContract.TransactOpts, _amount)
}

// EthSwapContractLogCloseIterator is returned from FilterLogClose and is used to iterate over the raw logs and unpacked data for LogClose events raised by the EthSwapContract contract.
type EthSwapContractLogCloseIterator struct {
	Event *EthSwapContractLogClose // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EthSwapContractLogCloseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthSwapContractLogClose)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EthSwapContractLogClose)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EthSwapContractLogCloseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthSwapContractLogCloseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthSwapContractLogClose represents a LogClose event raised by the EthSwapContract contract.
type EthSwapContractLogClose struct {
	SwapID    [32]byte
	SecretKey [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLogClose is a free log retrieval operation binding the contract event 0x07da1fa25a1d885732677ce9c192cbec27051a4b69d391c9a64850f5a5112ba0.
//
// Solidity: e LogClose(_swapID bytes32, _secretKey bytes32)
func (_EthSwapContract *EthSwapContractFilterer) FilterLogClose(opts *bind.FilterOpts) (*EthSwapContractLogCloseIterator, error) {

	logs, sub, err := _EthSwapContract.contract.FilterLogs(opts, "LogClose")
	if err != nil {
		return nil, err
	}
	return &EthSwapContractLogCloseIterator{contract: _EthSwapContract.contract, event: "LogClose", logs: logs, sub: sub}, nil
}

// WatchLogClose is a free log subscription operation binding the contract event 0x07da1fa25a1d885732677ce9c192cbec27051a4b69d391c9a64850f5a5112ba0.
//
// Solidity: e LogClose(_swapID bytes32, _secretKey bytes32)
func (_EthSwapContract *EthSwapContractFilterer) WatchLogClose(opts *bind.WatchOpts, sink chan<- *EthSwapContractLogClose) (event.Subscription, error) {

	logs, sub, err := _EthSwapContract.contract.WatchLogs(opts, "LogClose")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthSwapContractLogClose)
				if err := _EthSwapContract.contract.UnpackLog(event, "LogClose", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// EthSwapContractLogExpireIterator is returned from FilterLogExpire and is used to iterate over the raw logs and unpacked data for LogExpire events raised by the EthSwapContract contract.
type EthSwapContractLogExpireIterator struct {
	Event *EthSwapContractLogExpire // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EthSwapContractLogExpireIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthSwapContractLogExpire)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EthSwapContractLogExpire)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EthSwapContractLogExpireIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthSwapContractLogExpireIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthSwapContractLogExpire represents a LogExpire event raised by the EthSwapContract contract.
type EthSwapContractLogExpire struct {
	SwapID [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogExpire is a free log retrieval operation binding the contract event 0xeb711459e1247171382f0da0387b86239b8e3ca60af3b15a9ff2f1eb3d7f6a1d.
//
// Solidity: e LogExpire(_swapID bytes32)
func (_EthSwapContract *EthSwapContractFilterer) FilterLogExpire(opts *bind.FilterOpts) (*EthSwapContractLogExpireIterator, error) {

	logs, sub, err := _EthSwapContract.contract.FilterLogs(opts, "LogExpire")
	if err != nil {
		return nil, err
	}
	return &EthSwapContractLogExpireIterator{contract: _EthSwapContract.contract, event: "LogExpire", logs: logs, sub: sub}, nil
}

// WatchLogExpire is a free log subscription operation binding the contract event 0xeb711459e1247171382f0da0387b86239b8e3ca60af3b15a9ff2f1eb3d7f6a1d.
//
// Solidity: e LogExpire(_swapID bytes32)
func (_EthSwapContract *EthSwapContractFilterer) WatchLogExpire(opts *bind.WatchOpts, sink chan<- *EthSwapContractLogExpire) (event.Subscription, error) {

	logs, sub, err := _EthSwapContract.contract.WatchLogs(opts, "LogExpire")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthSwapContractLogExpire)
				if err := _EthSwapContract.contract.UnpackLog(event, "LogExpire", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// EthSwapContractLogOpenIterator is returned from FilterLogOpen and is used to iterate over the raw logs and unpacked data for LogOpen events raised by the EthSwapContract contract.
type EthSwapContractLogOpenIterator struct {
	Event *EthSwapContractLogOpen // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EthSwapContractLogOpenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthSwapContractLogOpen)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EthSwapContractLogOpen)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EthSwapContractLogOpenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthSwapContractLogOpenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthSwapContractLogOpen represents a LogOpen event raised by the EthSwapContract contract.
type EthSwapContractLogOpen struct {
	SwapID     [32]byte
	Spender    common.Address
	SecretLock [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogOpen is a free log retrieval operation binding the contract event 0x497d46e9505eefe8b910d1a02e6b40d8769510023b0053c3ac4b9574b81c97bf.
//
// Solidity: e LogOpen(_swapID bytes32, _spender address, _secretLock bytes32)
func (_EthSwapContract *EthSwapContractFilterer) FilterLogOpen(opts *bind.FilterOpts) (*EthSwapContractLogOpenIterator, error) {

	logs, sub, err := _EthSwapContract.contract.FilterLogs(opts, "LogOpen")
	if err != nil {
		return nil, err
	}
	return &EthSwapContractLogOpenIterator{contract: _EthSwapContract.contract, event: "LogOpen", logs: logs, sub: sub}, nil
}

// WatchLogOpen is a free log subscription operation binding the contract event 0x497d46e9505eefe8b910d1a02e6b40d8769510023b0053c3ac4b9574b81c97bf.
//
// Solidity: e LogOpen(_swapID bytes32, _spender address, _secretLock bytes32)
func (_EthSwapContract *EthSwapContractFilterer) WatchLogOpen(opts *bind.WatchOpts, sink chan<- *EthSwapContractLogOpen) (event.Subscription, error) {

	logs, sub, err := _EthSwapContract.contract.WatchLogs(opts, "LogOpen")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthSwapContractLogOpen)
				if err := _EthSwapContract.contract.UnpackLog(event, "LogOpen", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
