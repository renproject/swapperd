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
const SwapperdEthABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"initiatable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_broker\",\"type\":\"address\"},{\"name\":\"_brokerFee\",\"type\":\"uint256\"},{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"}],\"name\":\"initiate\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"}],\"name\":\"swapID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"redeemable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"auditSecret\",\"outputs\":[{\"name\":\"secretKey\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refundable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_secretKey\",\"type\":\"bytes32\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"redeemedAt\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"audit\",\"outputs\":[{\"name\":\"timelock\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"secretLock\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_VERSION\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_secretLock\",\"type\":\"bytes32\"}],\"name\":\"LogOpen\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"LogExpire\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_secretKey\",\"type\":\"bytes32\"}],\"name\":\"LogClose\",\"type\":\"event\"}]"

// SwapperdEthBin is the compiled bytecode used for deploying new contracts.
const SwapperdEthBin = `0x608060405234801561001057600080fd5b50604051610dba380380610dba8339810180604052602081101561003357600080fd5b81019080805164010000000081111561004b57600080fd5b8201602081018481111561005e57600080fd5b815164010000000081118282018710171561007857600080fd5b505080519093506100929250600091506020840190610099565b5050610134565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100da57805160ff1916838001178555610107565b82800160010185558215610107579182015b828111156101075782518255916020019190600101906100ec565b50610113929150610117565b5090565b61013191905b80821115610113576000815560010161011d565b90565b610c77806101436000396000f3fe6080604052600436106100ae5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166309ece61881146100b3578063395de232146100f15780634b2ac3fa1461013b57806368f06b291461017d5780637249fbb6146101a7578063976d00f4146101d15780639fb31475146101fb578063b31597ad14610225578063bc4fcc4a14610255578063c140635b1461027f578063ffa1ad74146102de575b600080fd5b3480156100bf57600080fd5b506100dd600480360360208110156100d657600080fd5b5035610368565b604080519115158252519081900360200190f35b610139600480360360c081101561010757600080fd5b50803590600160a060020a03602082013581169160408101359091169060608101359060808101359060a00135610390565b005b34801561014757600080fd5b5061016b6004803603604081101561015e57600080fd5b50803590602001356105ce565b60408051918252519081900360200190f35b34801561018957600080fd5b506100dd600480360360208110156101a057600080fd5b50356105fa565b3480156101b357600080fd5b50610139600480360360208110156101ca57600080fd5b5035610603565b3480156101dd57600080fd5b5061016b600480360360208110156101f457600080fd5b5035610782565b34801561020757600080fd5b506100dd6004803603602081101561021e57600080fd5b5035610810565b34801561023157600080fd5b506101396004803603604081101561024857600080fd5b5080359060200135610836565b34801561026157600080fd5b5061016b6004803603602081101561027857600080fd5b5035610ac4565b34801561028b57600080fd5b506102a9600480360360208110156102a257600080fd5b5035610ad6565b604080519586526020860194909452600160a060020a0392831685850152911660608401526080830152519081900360a00190f35b3480156102ea57600080fd5b506102f3610b79565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561032d578181015183820152602001610315565b50505050905090810190601f16801561035a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000805b60008381526002602052604090205460ff16600381111561038957fe5b1492915050565b856000808281526002602052604090205460ff1660038111156103af57fe5b14610404576040805160e560020a62461bcd02815260206004820152601660248201527f73776170206f70656e65642070726576696f75736c7900000000000000000000604482015290519081900360640190fd5b3484111561041157600080fd5b610419610c07565b6101006040519081016040528084815260200186340381526020018681526020018581526020016000600102815260200133600160a060020a0316815260200188600160a060020a0316815260200187600160a060020a0316815250905080600160008a8152602001908152602001600020600082015181600001556020820151816001015560408201518160020155606082015181600301556080820151816004015560a08201518160050160006101000a815481600160a060020a030219169083600160a060020a0316021790555060c08201518160060160006101000a815481600160a060020a030219169083600160a060020a0316021790555060e08201518160070160006101000a815481600160a060020a030219169083600160a060020a031602179055509050506001600260008a815260200190815260200160002060006101000a81548160ff0219169083600381111561057757fe5b021790555060408051898152600160a060020a038916602082015280820186905290517f497d46e9505eefe8b910d1a02e6b40d8769510023b0053c3ac4b9574b81c97bf9181900360600190a15050505050505050565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b6000600161036c565b80600160008281526002602052604090205460ff16600381111561062357fe5b14610678576040805160e560020a62461bcd02815260206004820152600d60248201527f73776170206e6f74206f70656e00000000000000000000000000000000000000604482015290519081900360640190fd5b60008281526001602052604090205482904210156106e0576040805160e560020a62461bcd02815260206004820152601260248201527f73776170206e6f7420657870697261626c650000000000000000000000000000604482015290519081900360640190fd5b60008381526002602081815260408084208054600360ff199091161790556001918290528084206005810154938101549201549051600160a060020a0390931693910180156108fc02929091818181858888f19350505050158015610749573d6000803e3d6000fd5b506040805184815290517feb711459e1247171382f0da0387b86239b8e3ca60af3b15a9ff2f1eb3d7f6a1d9181900360200190a1505050565b600081600260008281526002602052604090205460ff1660038111156107a457fe5b146107f9576040805160e560020a62461bcd02815260206004820152601160248201527f73776170206e6f742072656465656d6564000000000000000000000000000000604482015290519081900360640190fd5b505060009081526001602052604090206004015490565b60008181526001602052604081205442108015906108305750600161036c565b92915050565b81600160008281526002602052604090205460ff16600381111561085657fe5b146108ab576040805160e560020a62461bcd02815260206004820152600d60248201527f73776170206e6f74206f70656e00000000000000000000000000000000000000604482015290519081900360640190fd5b8282600281604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b602083106108fe5780518252601f1990920191602091820191016108df565b51815160209384036101000a60001901801990921691161790526040519190930194509192505080830381855afa15801561093d573d6000803e3d6000fd5b5050506040513d602081101561095257600080fd5b5051600083815260016020526040902060030154146109bb576040805160e560020a62461bcd02815260206004820152600e60248201527f696e76616c696420736563726574000000000000000000000000000000000000604482015290519081900360640190fd5b6000858152600160208181526040808420600481018990556002808452828620805460ff191690911790556003835281852042905591839052600682015491909201549151600160a060020a03909116926108fc831502929190818181858888f19350505050158015610a32573d6000803e3d6000fd5b5060008581526001602052604080822060078101546002909101549151600160a060020a039091169282156108fc02929190818181858888f19350505050158015610a81573d6000803e3d6000fd5b50604080518681526020810186905281517f07da1fa25a1d885732677ce9c192cbec27051a4b69d391c9a64850f5a5112ba0929181900390910190a15050505050565b60036020526000908152604090205481565b6000806000806000610ae6610c07565b5050506000938452505060016020818152604093849020845161010081018652815480825293820154928101839052600282015495810195909552600381015460608601819052600482015460808701526005820154600160a060020a0390811660a088018190526006840154821660c0890181905260079094015490911660e090970196909652929591949093509190565b6000805460408051602060026001851615610100026000190190941693909304601f81018490048402820184019092528181529291830182828015610bff5780601f10610bd457610100808354040283529160200191610bff565b820191906000526020600020905b815481529060010190602001808311610be257829003601f168201915b505050505081565b6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e08101919091529056fea165627a7a723058209d43a63ae95aea463ff9db002104ac2a508f9efe02674baf0cdd653fe7e6aef50029`

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
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, from address, secretLock bytes32)
func (_SwapperdEth *SwapperdEthCaller) Audit(opts *bind.CallOpts, _swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	ret := new(struct {
		Timelock   *big.Int
		Value      *big.Int
		To         common.Address
		From       common.Address
		SecretLock [32]byte
	})
	out := ret
	err := _SwapperdEth.contract.Call(opts, out, "audit", _swapID)
	return *ret, err
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, from address, secretLock bytes32)
func (_SwapperdEth *SwapperdEthSession) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	return _SwapperdEth.Contract.Audit(&_SwapperdEth.CallOpts, _swapID)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, from address, secretLock bytes32)
func (_SwapperdEth *SwapperdEthCallerSession) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
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

// Initiate is a paid mutator transaction binding the contract method 0x395de232.
//
// Solidity: function initiate(_swapID bytes32, _spender address, _broker address, _brokerFee uint256, _secretLock bytes32, _timelock uint256) returns()
func (_SwapperdEth *SwapperdEthTransactor) Initiate(opts *bind.TransactOpts, _swapID [32]byte, _spender common.Address, _broker common.Address, _brokerFee *big.Int, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _SwapperdEth.contract.Transact(opts, "initiate", _swapID, _spender, _broker, _brokerFee, _secretLock, _timelock)
}

// Initiate is a paid mutator transaction binding the contract method 0x395de232.
//
// Solidity: function initiate(_swapID bytes32, _spender address, _broker address, _brokerFee uint256, _secretLock bytes32, _timelock uint256) returns()
func (_SwapperdEth *SwapperdEthSession) Initiate(_swapID [32]byte, _spender common.Address, _broker common.Address, _brokerFee *big.Int, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _SwapperdEth.Contract.Initiate(&_SwapperdEth.TransactOpts, _swapID, _spender, _broker, _brokerFee, _secretLock, _timelock)
}

// Initiate is a paid mutator transaction binding the contract method 0x395de232.
//
// Solidity: function initiate(_swapID bytes32, _spender address, _broker address, _brokerFee uint256, _secretLock bytes32, _timelock uint256) returns()
func (_SwapperdEth *SwapperdEthTransactorSession) Initiate(_swapID [32]byte, _spender common.Address, _broker common.Address, _brokerFee *big.Int, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _SwapperdEth.Contract.Initiate(&_SwapperdEth.TransactOpts, _swapID, _spender, _broker, _brokerFee, _secretLock, _timelock)
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
