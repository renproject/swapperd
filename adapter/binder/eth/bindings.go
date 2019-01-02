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

// SwapperdEthABI is the input ABI used to generate the binding from.
const SwapperdEthABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"initiate\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"initiatable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"}],\"name\":\"swapID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawBrokerFees\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"redeemable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"auditSecret\",\"outputs\":[{\"name\":\"secretKey\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refundable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_secretKey\",\"type\":\"bytes32\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_broker\",\"type\":\"address\"},{\"name\":\"_brokerFee\",\"type\":\"uint256\"},{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"initiateWithFees\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"redeemedAt\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"audit\",\"outputs\":[{\"name\":\"timelock\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"brokerFee\",\"type\":\"uint256\"},{\"name\":\"broker\",\"type\":\"address\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"secretLock\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"brokerFees\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_VERSION\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_secretLock\",\"type\":\"bytes32\"}],\"name\":\"LogOpen\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"LogExpire\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_secretKey\",\"type\":\"bytes32\"}],\"name\":\"LogClose\",\"type\":\"event\"}]"

// SwapperdEthBin is the compiled bytecode used for deploying new contracts.
const SwapperdEthBin = `0x60806040523480156200001157600080fd5b506040516200113138038062001131833981018060405260208110156200003757600080fd5b8101908080516401000000008111156200005057600080fd5b820160208101848111156200006457600080fd5b81516401000000008111828201871017156200007f57600080fd5b505080519093506200009b9250600091506020840190620000a3565b505062000148565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620000e657805160ff191683800117855562000116565b8280016001018555821562000116579182015b8281111562000116578251825591602001919060010190620000f9565b506200012492915062000128565b5090565b6200014591905b808211156200012457600081556001016200012f565b90565b610fd980620001586000396000f3fe6080604052600436106100c9577c01000000000000000000000000000000000000000000000000000000006000350463027a257781146100ce57806309ece6181461010e5780634b2ac3fa1461014c5780634c6d37ff1461018e57806368f06b29146101b85780637249fbb6146101e2578063976d00f41461020c5780639fb3147514610236578063b31597ad14610260578063b8688e3f14610290578063bc4fcc4a146102de578063c140635b14610308578063e1ec380c14610378578063ffa1ad74146103ab575b600080fd5b61010c600480360360a08110156100e457600080fd5b50803590600160a060020a036020820135169060408101359060608101359060800135610435565b005b34801561011a57600080fd5b506101386004803603602081101561013157600080fd5b5035610671565b604080519115158252519081900360200190f35b34801561015857600080fd5b5061017c6004803603604081101561016f57600080fd5b5080359060200135610699565b60408051918252519081900360200190f35b34801561019a57600080fd5b5061010c600480360360208110156101b157600080fd5b50356106c5565b3480156101c457600080fd5b50610138600480360360208110156101db57600080fd5b5035610724565b3480156101ee57600080fd5b5061010c6004803603602081101561020557600080fd5b503561072d565b34801561021857600080fd5b5061017c6004803603602081101561022f57600080fd5b50356108ac565b34801561024257600080fd5b506101386004803603602081101561025957600080fd5b503561093a565b34801561026c57600080fd5b5061010c6004803603604081101561028357600080fd5b5080359060200135610960565b61010c600480360360e08110156102a657600080fd5b50803590600160a060020a03602082013581169160408101359091169060608101359060808101359060a08101359060c00135610bc3565b3480156102ea57600080fd5b5061017c6004803603602081101561030157600080fd5b5035610e0f565b34801561031457600080fd5b506103326004803603602081101561032b57600080fd5b5035610e21565b604080519788526020880196909652600160a060020a03948516878701526060870193909352908316608086015290911660a084015260c0830152519081900360e00190f35b34801561038457600080fd5b5061017c6004803603602081101561039b57600080fd5b5035600160a060020a0316610ec9565b3480156103b757600080fd5b506103c0610edb565b6040805160208082528351818301528351919283929083019185019080838360005b838110156103fa5781810151838201526020016103e2565b50505050905090810190601f1680156104275780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b846000808281526002602052604090205460ff16600381111561045457fe5b146104a9576040805160e560020a62461bcd02815260206004820152601660248201527f73776170206f70656e65642070726576696f75736c7900000000000000000000604482015290519081900360640190fd5b3482146104b557600080fd5b6104bd610f69565b61010060405190810160405280858152602001848152602001600081526020018681526020016000600102815260200133600160a060020a0316815260200187600160a060020a031681526020016000600160a060020a031681525090508060016000898152602001908152602001600020600082015181600001556020820151816001015560408201518160020155606082015181600301556080820151816004015560a08201518160050160006101000a815481600160a060020a030219169083600160a060020a0316021790555060c08201518160060160006101000a815481600160a060020a030219169083600160a060020a0316021790555060e08201518160070160006101000a815481600160a060020a030219169083600160a060020a0316021790555090505060016002600089815260200190815260200160002060006101000a81548160ff0219169083600381111561061b57fe5b021790555060408051888152600160a060020a038816602082015280820187905290517f497d46e9505eefe8b910d1a02e6b40d8769510023b0053c3ac4b9574b81c97bf9181900360600190a150505050505050565b6000805b60008381526002602052604090205460ff16600381111561069257fe5b1492915050565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b336000908152600360205260409020548111156106e157600080fd5b33600081815260036020526040808220805485900390555183156108fc0291849190818181858888f19350505050158015610720573d6000803e3d6000fd5b5050565b60006001610675565b80600160008281526002602052604090205460ff16600381111561074d57fe5b146107a2576040805160e560020a62461bcd02815260206004820152600d60248201527f73776170206e6f74206f70656e00000000000000000000000000000000000000604482015290519081900360640190fd5b600082815260016020526040902054829042101561080a576040805160e560020a62461bcd02815260206004820152601260248201527f73776170206e6f7420657870697261626c650000000000000000000000000000604482015290519081900360640190fd5b60008381526002602081815260408084208054600360ff199091161790556001918290528084206005810154938101549201549051600160a060020a0390931693910180156108fc02929091818181858888f19350505050158015610873573d6000803e3d6000fd5b506040805184815290517feb711459e1247171382f0da0387b86239b8e3ca60af3b15a9ff2f1eb3d7f6a1d9181900360200190a1505050565b600081600260008281526002602052604090205460ff1660038111156108ce57fe5b14610923576040805160e560020a62461bcd02815260206004820152601160248201527f73776170206e6f742072656465656d6564000000000000000000000000000000604482015290519081900360640190fd5b505060009081526001602052604090206004015490565b600081815260016020526040812054421080159061095a57506001610675565b92915050565b81600160008281526002602052604090205460ff16600381111561098057fe5b146109d5576040805160e560020a62461bcd02815260206004820152600d60248201527f73776170206e6f74206f70656e00000000000000000000000000000000000000604482015290519081900360640190fd5b8282600281604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b60208310610a285780518252601f199092019160209182019101610a09565b51815160209384036101000a60001901801990921691161790526040519190930194509192505080830381855afa158015610a67573d6000803e3d6000fd5b5050506040513d6020811015610a7c57600080fd5b505160008381526001602052604090206003015414610ae5576040805160e560020a62461bcd02815260206004820152600e60248201527f696e76616c696420736563726574000000000000000000000000000000000000604482015290519081900360640190fd5b600085815260016020818152604080842060048082018a90556002808552838720805460ff19169091179055835281852042905591839052600682015491909201549151600160a060020a03909116926108fc831502929190818181858888f19350505050158015610b5b573d6000803e3d6000fd5b5060008581526001602090815260408083206002015433845260038352928190208054909301909255815187815290810186905281517f07da1fa25a1d885732677ce9c192cbec27051a4b69d391c9a64850f5a5112ba0929181900390910190a15050505050565b866000808281526002602052604090205460ff166003811115610be257fe5b14610c37576040805160e560020a62461bcd02815260206004820152601660248201527f73776170206f70656e65642070726576696f75736c7900000000000000000000604482015290519081900360640190fd5b3482148015610c465750848210155b1515610c5157600080fd5b610c59610f69565b6101006040519081016040528085815260200187850381526020018781526020018681526020016000600102815260200133600160a060020a0316815260200189600160a060020a0316815260200188600160a060020a0316815250905080600160008b8152602001908152602001600020600082015181600001556020820151816001015560408201518160020155606082015181600301556080820151816004015560a08201518160050160006101000a815481600160a060020a030219169083600160a060020a0316021790555060c08201518160060160006101000a815481600160a060020a030219169083600160a060020a0316021790555060e08201518160070160006101000a815481600160a060020a030219169083600160a060020a031602179055509050506001600260008b815260200190815260200160002060006101000a81548160ff02191690836003811115610db757fe5b0217905550604080518a8152600160a060020a038a16602082015280820187905290517f497d46e9505eefe8b910d1a02e6b40d8769510023b0053c3ac4b9574b81c97bf9181900360600190a1505050505050505050565b60046020526000908152604090205481565b6000806000806000806000610e34610f69565b50505060009586525050600160208181526040958690208651610100810188528154808252938201549281018390526002820154978101889052600382015460608201819052600483015460808301526005830154600160a060020a0390811660a084018190526006850154821660c0850181905260079095015490911660e090930183905294999398929750919550935090565b60036020526000908152604090205481565b6000805460408051602060026001851615610100026000190190941693909304601f81018490048402820184019092528181529291830182828015610f615780601f10610f3657610100808354040283529160200191610f61565b820191906000526020600020905b815481529060010190602001808311610f4457829003601f168201915b505050505081565b6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e08101919091529056fea165627a7a723058203a2088fd385f49a13ab25717ba969e780724eae2b0d4efe9ac916ba560f701f90029`

// DeploySwapperdEth deploys a new Ethereum contract, binding an instance of SwapperdEth to it.
func DeploySwapperdEth(auth *bind.TransactOpts, backend bind.ContractBackend, _VERSION string) (common.Address, *types.Transaction, *SwapperdEth, error) {
	parsed, err := abi.JSON(strings.NewReader(SwapperdEthABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SwapperdEthBin), backend, _VERSION)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SwapperdEth{SwapperdEthCaller: SwapperdEthCaller{contract: contract}, SwapperdEthTransactor: SwapperdEthTransactor{contract: contract}, SwapperdEthFilterer: SwapperdEthFilterer{contract: contract}}, nil
}

// SwapperdEth is an auto generated Go binding around an Ethereum contract.
type SwapperdEth struct {
	SwapperdEthCaller     // Read-only binding to the contract
	SwapperdEthTransactor // Write-only binding to the contract
	SwapperdEthFilterer   // Log filterer for contract events
}

// SwapperdEthCaller is an auto generated read-only Go binding around an Ethereum contract.
type SwapperdEthCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapperdEthTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SwapperdEthTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapperdEthFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SwapperdEthFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapperdEthSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SwapperdEthSession struct {
	Contract     *SwapperdEth      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SwapperdEthCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SwapperdEthCallerSession struct {
	Contract *SwapperdEthCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// SwapperdEthTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SwapperdEthTransactorSession struct {
	Contract     *SwapperdEthTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// SwapperdEthRaw is an auto generated low-level Go binding around an Ethereum contract.
type SwapperdEthRaw struct {
	Contract *SwapperdEth // Generic contract binding to access the raw methods on
}

// SwapperdEthCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SwapperdEthCallerRaw struct {
	Contract *SwapperdEthCaller // Generic read-only contract binding to access the raw methods on
}

// SwapperdEthTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SwapperdEthTransactorRaw struct {
	Contract *SwapperdEthTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSwapperdEth creates a new instance of SwapperdEth, bound to a specific deployed contract.
func NewSwapperdEth(address common.Address, backend bind.ContractBackend) (*SwapperdEth, error) {
	contract, err := bindSwapperdEth(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SwapperdEth{SwapperdEthCaller: SwapperdEthCaller{contract: contract}, SwapperdEthTransactor: SwapperdEthTransactor{contract: contract}, SwapperdEthFilterer: SwapperdEthFilterer{contract: contract}}, nil
}

// NewSwapperdEthCaller creates a new read-only instance of SwapperdEth, bound to a specific deployed contract.
func NewSwapperdEthCaller(address common.Address, caller bind.ContractCaller) (*SwapperdEthCaller, error) {
	contract, err := bindSwapperdEth(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SwapperdEthCaller{contract: contract}, nil
}

// NewSwapperdEthTransactor creates a new write-only instance of SwapperdEth, bound to a specific deployed contract.
func NewSwapperdEthTransactor(address common.Address, transactor bind.ContractTransactor) (*SwapperdEthTransactor, error) {
	contract, err := bindSwapperdEth(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SwapperdEthTransactor{contract: contract}, nil
}

// NewSwapperdEthFilterer creates a new log filterer instance of SwapperdEth, bound to a specific deployed contract.
func NewSwapperdEthFilterer(address common.Address, filterer bind.ContractFilterer) (*SwapperdEthFilterer, error) {
	contract, err := bindSwapperdEth(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SwapperdEthFilterer{contract: contract}, nil
}

// bindSwapperdEth binds a generic wrapper to an already deployed contract.
func bindSwapperdEth(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SwapperdEthABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SwapperdEth *SwapperdEthRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SwapperdEth.Contract.SwapperdEthCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SwapperdEth *SwapperdEthRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwapperdEth.Contract.SwapperdEthTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SwapperdEth *SwapperdEthRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SwapperdEth.Contract.SwapperdEthTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SwapperdEth *SwapperdEthCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SwapperdEth.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SwapperdEth *SwapperdEthTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwapperdEth.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SwapperdEth *SwapperdEthTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SwapperdEth.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_SwapperdEth *SwapperdEthCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _SwapperdEth.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_SwapperdEth *SwapperdEthSession) VERSION() (string, error) {
	return _SwapperdEth.Contract.VERSION(&_SwapperdEth.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_SwapperdEth *SwapperdEthCallerSession) VERSION() (string, error) {
	return _SwapperdEth.Contract.VERSION(&_SwapperdEth.CallOpts)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, brokerFee uint256, broker address, from address, secretLock bytes32)
func (_SwapperdEth *SwapperdEthCaller) Audit(opts *bind.CallOpts, _swapID [32]byte) (struct {
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
	err := _SwapperdEth.contract.Call(opts, out, "audit", _swapID)
	return *ret, err
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, brokerFee uint256, broker address, from address, secretLock bytes32)
func (_SwapperdEth *SwapperdEthSession) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	BrokerFee  *big.Int
	Broker     common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	return _SwapperdEth.Contract.Audit(&_SwapperdEth.CallOpts, _swapID)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, brokerFee uint256, broker address, from address, secretLock bytes32)
func (_SwapperdEth *SwapperdEthCallerSession) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	BrokerFee  *big.Int
	Broker     common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	return _SwapperdEth.Contract.Audit(&_SwapperdEth.CallOpts, _swapID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_SwapperdEth *SwapperdEthCaller) AuditSecret(opts *bind.CallOpts, _swapID [32]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _SwapperdEth.contract.Call(opts, out, "auditSecret", _swapID)
	return *ret0, err
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_SwapperdEth *SwapperdEthSession) AuditSecret(_swapID [32]byte) ([32]byte, error) {
	return _SwapperdEth.Contract.AuditSecret(&_SwapperdEth.CallOpts, _swapID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_SwapperdEth *SwapperdEthCallerSession) AuditSecret(_swapID [32]byte) ([32]byte, error) {
	return _SwapperdEth.Contract.AuditSecret(&_SwapperdEth.CallOpts, _swapID)
}

// BrokerFees is a free data retrieval call binding the contract method 0xe1ec380c.
//
// Solidity: function brokerFees( address) constant returns(uint256)
func (_SwapperdEth *SwapperdEthCaller) BrokerFees(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SwapperdEth.contract.Call(opts, out, "brokerFees", arg0)
	return *ret0, err
}

// BrokerFees is a free data retrieval call binding the contract method 0xe1ec380c.
//
// Solidity: function brokerFees( address) constant returns(uint256)
func (_SwapperdEth *SwapperdEthSession) BrokerFees(arg0 common.Address) (*big.Int, error) {
	return _SwapperdEth.Contract.BrokerFees(&_SwapperdEth.CallOpts, arg0)
}

// BrokerFees is a free data retrieval call binding the contract method 0xe1ec380c.
//
// Solidity: function brokerFees( address) constant returns(uint256)
func (_SwapperdEth *SwapperdEthCallerSession) BrokerFees(arg0 common.Address) (*big.Int, error) {
	return _SwapperdEth.Contract.BrokerFees(&_SwapperdEth.CallOpts, arg0)
}

// Initiatable is a free data retrieval call binding the contract method 0x09ece618.
//
// Solidity: function initiatable(_swapID bytes32) constant returns(bool)
func (_SwapperdEth *SwapperdEthCaller) Initiatable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SwapperdEth.contract.Call(opts, out, "initiatable", _swapID)
	return *ret0, err
}

// Initiatable is a free data retrieval call binding the contract method 0x09ece618.
//
// Solidity: function initiatable(_swapID bytes32) constant returns(bool)
func (_SwapperdEth *SwapperdEthSession) Initiatable(_swapID [32]byte) (bool, error) {
	return _SwapperdEth.Contract.Initiatable(&_SwapperdEth.CallOpts, _swapID)
}

// Initiatable is a free data retrieval call binding the contract method 0x09ece618.
//
// Solidity: function initiatable(_swapID bytes32) constant returns(bool)
func (_SwapperdEth *SwapperdEthCallerSession) Initiatable(_swapID [32]byte) (bool, error) {
	return _SwapperdEth.Contract.Initiatable(&_SwapperdEth.CallOpts, _swapID)
}

// Redeemable is a free data retrieval call binding the contract method 0x68f06b29.
//
// Solidity: function redeemable(_swapID bytes32) constant returns(bool)
func (_SwapperdEth *SwapperdEthCaller) Redeemable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SwapperdEth.contract.Call(opts, out, "redeemable", _swapID)
	return *ret0, err
}

// Redeemable is a free data retrieval call binding the contract method 0x68f06b29.
//
// Solidity: function redeemable(_swapID bytes32) constant returns(bool)
func (_SwapperdEth *SwapperdEthSession) Redeemable(_swapID [32]byte) (bool, error) {
	return _SwapperdEth.Contract.Redeemable(&_SwapperdEth.CallOpts, _swapID)
}

// Redeemable is a free data retrieval call binding the contract method 0x68f06b29.
//
// Solidity: function redeemable(_swapID bytes32) constant returns(bool)
func (_SwapperdEth *SwapperdEthCallerSession) Redeemable(_swapID [32]byte) (bool, error) {
	return _SwapperdEth.Contract.Redeemable(&_SwapperdEth.CallOpts, _swapID)
}

// RedeemedAt is a free data retrieval call binding the contract method 0xbc4fcc4a.
//
// Solidity: function redeemedAt( bytes32) constant returns(uint256)
func (_SwapperdEth *SwapperdEthCaller) RedeemedAt(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SwapperdEth.contract.Call(opts, out, "redeemedAt", arg0)
	return *ret0, err
}

// RedeemedAt is a free data retrieval call binding the contract method 0xbc4fcc4a.
//
// Solidity: function redeemedAt( bytes32) constant returns(uint256)
func (_SwapperdEth *SwapperdEthSession) RedeemedAt(arg0 [32]byte) (*big.Int, error) {
	return _SwapperdEth.Contract.RedeemedAt(&_SwapperdEth.CallOpts, arg0)
}

// RedeemedAt is a free data retrieval call binding the contract method 0xbc4fcc4a.
//
// Solidity: function redeemedAt( bytes32) constant returns(uint256)
func (_SwapperdEth *SwapperdEthCallerSession) RedeemedAt(arg0 [32]byte) (*big.Int, error) {
	return _SwapperdEth.Contract.RedeemedAt(&_SwapperdEth.CallOpts, arg0)
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(_swapID bytes32) constant returns(bool)
func (_SwapperdEth *SwapperdEthCaller) Refundable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SwapperdEth.contract.Call(opts, out, "refundable", _swapID)
	return *ret0, err
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(_swapID bytes32) constant returns(bool)
func (_SwapperdEth *SwapperdEthSession) Refundable(_swapID [32]byte) (bool, error) {
	return _SwapperdEth.Contract.Refundable(&_SwapperdEth.CallOpts, _swapID)
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(_swapID bytes32) constant returns(bool)
func (_SwapperdEth *SwapperdEthCallerSession) Refundable(_swapID [32]byte) (bool, error) {
	return _SwapperdEth.Contract.Refundable(&_SwapperdEth.CallOpts, _swapID)
}

// SwapID is a free data retrieval call binding the contract method 0x4b2ac3fa.
//
// Solidity: function swapID(_secretLock bytes32, _timelock uint256) constant returns(bytes32)
func (_SwapperdEth *SwapperdEthCaller) SwapID(opts *bind.CallOpts, _secretLock [32]byte, _timelock *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _SwapperdEth.contract.Call(opts, out, "swapID", _secretLock, _timelock)
	return *ret0, err
}

// SwapID is a free data retrieval call binding the contract method 0x4b2ac3fa.
//
// Solidity: function swapID(_secretLock bytes32, _timelock uint256) constant returns(bytes32)
func (_SwapperdEth *SwapperdEthSession) SwapID(_secretLock [32]byte, _timelock *big.Int) ([32]byte, error) {
	return _SwapperdEth.Contract.SwapID(&_SwapperdEth.CallOpts, _secretLock, _timelock)
}

// SwapID is a free data retrieval call binding the contract method 0x4b2ac3fa.
//
// Solidity: function swapID(_secretLock bytes32, _timelock uint256) constant returns(bytes32)
func (_SwapperdEth *SwapperdEthCallerSession) SwapID(_secretLock [32]byte, _timelock *big.Int) ([32]byte, error) {
	return _SwapperdEth.Contract.SwapID(&_SwapperdEth.CallOpts, _secretLock, _timelock)
}

// Initiate is a paid mutator transaction binding the contract method 0x027a2577.
//
// Solidity: function initiate(_swapID bytes32, _spender address, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_SwapperdEth *SwapperdEthTransactor) Initiate(opts *bind.TransactOpts, _swapID [32]byte, _spender common.Address, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _SwapperdEth.contract.Transact(opts, "initiate", _swapID, _spender, _secretLock, _timelock, _value)
}

// Initiate is a paid mutator transaction binding the contract method 0x027a2577.
//
// Solidity: function initiate(_swapID bytes32, _spender address, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_SwapperdEth *SwapperdEthSession) Initiate(_swapID [32]byte, _spender common.Address, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _SwapperdEth.Contract.Initiate(&_SwapperdEth.TransactOpts, _swapID, _spender, _secretLock, _timelock, _value)
}

// Initiate is a paid mutator transaction binding the contract method 0x027a2577.
//
// Solidity: function initiate(_swapID bytes32, _spender address, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_SwapperdEth *SwapperdEthTransactorSession) Initiate(_swapID [32]byte, _spender common.Address, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _SwapperdEth.Contract.Initiate(&_SwapperdEth.TransactOpts, _swapID, _spender, _secretLock, _timelock, _value)
}

// InitiateWithFees is a paid mutator transaction binding the contract method 0xb8688e3f.
//
// Solidity: function initiateWithFees(_swapID bytes32, _spender address, _broker address, _brokerFee uint256, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_SwapperdEth *SwapperdEthTransactor) InitiateWithFees(opts *bind.TransactOpts, _swapID [32]byte, _spender common.Address, _broker common.Address, _brokerFee *big.Int, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _SwapperdEth.contract.Transact(opts, "initiateWithFees", _swapID, _spender, _broker, _brokerFee, _secretLock, _timelock, _value)
}

// InitiateWithFees is a paid mutator transaction binding the contract method 0xb8688e3f.
//
// Solidity: function initiateWithFees(_swapID bytes32, _spender address, _broker address, _brokerFee uint256, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_SwapperdEth *SwapperdEthSession) InitiateWithFees(_swapID [32]byte, _spender common.Address, _broker common.Address, _brokerFee *big.Int, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _SwapperdEth.Contract.InitiateWithFees(&_SwapperdEth.TransactOpts, _swapID, _spender, _broker, _brokerFee, _secretLock, _timelock, _value)
}

// InitiateWithFees is a paid mutator transaction binding the contract method 0xb8688e3f.
//
// Solidity: function initiateWithFees(_swapID bytes32, _spender address, _broker address, _brokerFee uint256, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_SwapperdEth *SwapperdEthTransactorSession) InitiateWithFees(_swapID [32]byte, _spender common.Address, _broker common.Address, _brokerFee *big.Int, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _SwapperdEth.Contract.InitiateWithFees(&_SwapperdEth.TransactOpts, _swapID, _spender, _broker, _brokerFee, _secretLock, _timelock, _value)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_SwapperdEth *SwapperdEthTransactor) Redeem(opts *bind.TransactOpts, _swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _SwapperdEth.contract.Transact(opts, "redeem", _swapID, _secretKey)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_SwapperdEth *SwapperdEthSession) Redeem(_swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _SwapperdEth.Contract.Redeem(&_SwapperdEth.TransactOpts, _swapID, _secretKey)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_SwapperdEth *SwapperdEthTransactorSession) Redeem(_swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _SwapperdEth.Contract.Redeem(&_SwapperdEth.TransactOpts, _swapID, _secretKey)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_SwapperdEth *SwapperdEthTransactor) Refund(opts *bind.TransactOpts, _swapID [32]byte) (*types.Transaction, error) {
	return _SwapperdEth.contract.Transact(opts, "refund", _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_SwapperdEth *SwapperdEthSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _SwapperdEth.Contract.Refund(&_SwapperdEth.TransactOpts, _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_SwapperdEth *SwapperdEthTransactorSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _SwapperdEth.Contract.Refund(&_SwapperdEth.TransactOpts, _swapID)
}

// WithdrawBrokerFees is a paid mutator transaction binding the contract method 0x4c6d37ff.
//
// Solidity: function withdrawBrokerFees(_amount uint256) returns()
func (_SwapperdEth *SwapperdEthTransactor) WithdrawBrokerFees(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _SwapperdEth.contract.Transact(opts, "withdrawBrokerFees", _amount)
}

// WithdrawBrokerFees is a paid mutator transaction binding the contract method 0x4c6d37ff.
//
// Solidity: function withdrawBrokerFees(_amount uint256) returns()
func (_SwapperdEth *SwapperdEthSession) WithdrawBrokerFees(_amount *big.Int) (*types.Transaction, error) {
	return _SwapperdEth.Contract.WithdrawBrokerFees(&_SwapperdEth.TransactOpts, _amount)
}

// WithdrawBrokerFees is a paid mutator transaction binding the contract method 0x4c6d37ff.
//
// Solidity: function withdrawBrokerFees(_amount uint256) returns()
func (_SwapperdEth *SwapperdEthTransactorSession) WithdrawBrokerFees(_amount *big.Int) (*types.Transaction, error) {
	return _SwapperdEth.Contract.WithdrawBrokerFees(&_SwapperdEth.TransactOpts, _amount)
}

// SwapperdEthLogCloseIterator is returned from FilterLogClose and is used to iterate over the raw logs and unpacked data for LogClose events raised by the SwapperdEth contract.
type SwapperdEthLogCloseIterator struct {
	Event *SwapperdEthLogClose // Event containing the contract specifics and raw log

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
func (it *SwapperdEthLogCloseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SwapperdEthLogClose)
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
		it.Event = new(SwapperdEthLogClose)
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
func (it *SwapperdEthLogCloseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SwapperdEthLogCloseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SwapperdEthLogClose represents a LogClose event raised by the SwapperdEth contract.
type SwapperdEthLogClose struct {
	SwapID    [32]byte
	SecretKey [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLogClose is a free log retrieval operation binding the contract event 0x07da1fa25a1d885732677ce9c192cbec27051a4b69d391c9a64850f5a5112ba0.
//
// Solidity: e LogClose(_swapID bytes32, _secretKey bytes32)
func (_SwapperdEth *SwapperdEthFilterer) FilterLogClose(opts *bind.FilterOpts) (*SwapperdEthLogCloseIterator, error) {

	logs, sub, err := _SwapperdEth.contract.FilterLogs(opts, "LogClose")
	if err != nil {
		return nil, err
	}
	return &SwapperdEthLogCloseIterator{contract: _SwapperdEth.contract, event: "LogClose", logs: logs, sub: sub}, nil
}

// WatchLogClose is a free log subscription operation binding the contract event 0x07da1fa25a1d885732677ce9c192cbec27051a4b69d391c9a64850f5a5112ba0.
//
// Solidity: e LogClose(_swapID bytes32, _secretKey bytes32)
func (_SwapperdEth *SwapperdEthFilterer) WatchLogClose(opts *bind.WatchOpts, sink chan<- *SwapperdEthLogClose) (event.Subscription, error) {

	logs, sub, err := _SwapperdEth.contract.WatchLogs(opts, "LogClose")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SwapperdEthLogClose)
				if err := _SwapperdEth.contract.UnpackLog(event, "LogClose", log); err != nil {
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

// SwapperdEthLogExpireIterator is returned from FilterLogExpire and is used to iterate over the raw logs and unpacked data for LogExpire events raised by the SwapperdEth contract.
type SwapperdEthLogExpireIterator struct {
	Event *SwapperdEthLogExpire // Event containing the contract specifics and raw log

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
func (it *SwapperdEthLogExpireIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SwapperdEthLogExpire)
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
		it.Event = new(SwapperdEthLogExpire)
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
func (it *SwapperdEthLogExpireIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SwapperdEthLogExpireIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SwapperdEthLogExpire represents a LogExpire event raised by the SwapperdEth contract.
type SwapperdEthLogExpire struct {
	SwapID [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogExpire is a free log retrieval operation binding the contract event 0xeb711459e1247171382f0da0387b86239b8e3ca60af3b15a9ff2f1eb3d7f6a1d.
//
// Solidity: e LogExpire(_swapID bytes32)
func (_SwapperdEth *SwapperdEthFilterer) FilterLogExpire(opts *bind.FilterOpts) (*SwapperdEthLogExpireIterator, error) {

	logs, sub, err := _SwapperdEth.contract.FilterLogs(opts, "LogExpire")
	if err != nil {
		return nil, err
	}
	return &SwapperdEthLogExpireIterator{contract: _SwapperdEth.contract, event: "LogExpire", logs: logs, sub: sub}, nil
}

// WatchLogExpire is a free log subscription operation binding the contract event 0xeb711459e1247171382f0da0387b86239b8e3ca60af3b15a9ff2f1eb3d7f6a1d.
//
// Solidity: e LogExpire(_swapID bytes32)
func (_SwapperdEth *SwapperdEthFilterer) WatchLogExpire(opts *bind.WatchOpts, sink chan<- *SwapperdEthLogExpire) (event.Subscription, error) {

	logs, sub, err := _SwapperdEth.contract.WatchLogs(opts, "LogExpire")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SwapperdEthLogExpire)
				if err := _SwapperdEth.contract.UnpackLog(event, "LogExpire", log); err != nil {
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

// SwapperdEthLogOpenIterator is returned from FilterLogOpen and is used to iterate over the raw logs and unpacked data for LogOpen events raised by the SwapperdEth contract.
type SwapperdEthLogOpenIterator struct {
	Event *SwapperdEthLogOpen // Event containing the contract specifics and raw log

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
func (it *SwapperdEthLogOpenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SwapperdEthLogOpen)
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
		it.Event = new(SwapperdEthLogOpen)
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
func (it *SwapperdEthLogOpenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SwapperdEthLogOpenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SwapperdEthLogOpen represents a LogOpen event raised by the SwapperdEth contract.
type SwapperdEthLogOpen struct {
	SwapID     [32]byte
	Spender    common.Address
	SecretLock [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogOpen is a free log retrieval operation binding the contract event 0x497d46e9505eefe8b910d1a02e6b40d8769510023b0053c3ac4b9574b81c97bf.
//
// Solidity: e LogOpen(_swapID bytes32, _spender address, _secretLock bytes32)
func (_SwapperdEth *SwapperdEthFilterer) FilterLogOpen(opts *bind.FilterOpts) (*SwapperdEthLogOpenIterator, error) {

	logs, sub, err := _SwapperdEth.contract.FilterLogs(opts, "LogOpen")
	if err != nil {
		return nil, err
	}
	return &SwapperdEthLogOpenIterator{contract: _SwapperdEth.contract, event: "LogOpen", logs: logs, sub: sub}, nil
}

// WatchLogOpen is a free log subscription operation binding the contract event 0x497d46e9505eefe8b910d1a02e6b40d8769510023b0053c3ac4b9574b81c97bf.
//
// Solidity: e LogOpen(_swapID bytes32, _spender address, _secretLock bytes32)
func (_SwapperdEth *SwapperdEthFilterer) WatchLogOpen(opts *bind.WatchOpts, sink chan<- *SwapperdEthLogOpen) (event.Subscription, error) {

	logs, sub, err := _SwapperdEth.contract.WatchLogs(opts, "LogOpen")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SwapperdEthLogOpen)
				if err := _SwapperdEth.contract.UnpackLog(event, "LogOpen", log); err != nil {
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
