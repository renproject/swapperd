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

// AtomSwapABI is the input ABI used to generate the binding from.
const AtomSwapABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_withdrawTrader\",\"type\":\"address\"},{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"}],\"name\":\"initiate\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"auditSecret\",\"outputs\":[{\"name\":\"secretKey\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"swapDetails\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_secretKey\",\"type\":\"bytes32\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"audit\",\"outputs\":[{\"name\":\"timelock\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"secretLock\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_withdrawTrader\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_secretLock\",\"type\":\"bytes32\"}],\"name\":\"Open\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"Expire\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_secretKey\",\"type\":\"bytes\"}],\"name\":\"Close\",\"type\":\"event\"}]"

// AtomSwapBin is the compiled bytecode used for deploying new contracts.
const AtomSwapBin = `0x608060405234801561001057600080fd5b50610720806100206000396000f3006080604052600436106100775763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663412c0b58811461007c5780637249fbb61461009b578063976d00f4146100b3578063b14631bb146100dd578063b31597ad1461016a578063c140635b14610185575b600080fd5b610099600435600160a060020a03602435166044356064356101d2565b005b3480156100a757600080fd5b506100996004356102bb565b3480156100bf57600080fd5b506100cb6004356103b6565b60408051918252519081900360200190f35b3480156100e957600080fd5b506100f5600435610454565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561012f578181015183820152602001610117565b50505050905090810190601f16801561015c5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561017657600080fd5b506100996004356024356104ed565b34801561019157600080fd5b5061019d600435610639565b604080519586526020860194909452600160a060020a0392831685850152911660608401526080830152519081900360a00190f35b6101da6106bf565b846000808281526001602052604090205460ff1660038111156101f957fe5b1461020357600080fd5b50506040805160c08101825291825234602080840191825233848401908152600160a060020a039687166060860190815260808601968752600060a087018181529981528084528581209651875593516001808801919091559151600287018054918a1673ffffffffffffffffffffffffffffffffffffffff1992831617905590516003870180549190991691161790965593516004840155945160059092019190915590829052909120805460ff19169091179055565b6102c36106bf565b81600160008281526001602052604090205460ff1660038111156102e357fe5b146102ed57600080fd5b600083815260208190526040902054839042101561030a57600080fd5b600084815260208181526040808320815160c081018352815481526001808301548286019081526002840154600160a060020a03908116848701908152600380870154831660608701526004870154608087015260059096015460a08601528c895292909652848720805460ff1916909417909355519151925190975092169281156108fc029290818181858888f193505050501580156103af573d6000803e3d6000fd5b5050505050565b60006103c06106bf565b82600260008281526001602052604090205460ff1660038111156103e057fe5b146103ea57600080fd5b50505060009081526020818152604091829020825160c081018452815481526001820154928101929092526002810154600160a060020a0390811693830193909352600381015490921660608201526004820154608082015260059091015460a090910181905290565b600260208181526000928352604092839020805484516001821615610100026000190190911693909304601f81018390048302840183019094528383529192908301828280156104e55780601f106104ba576101008083540402835291602001916104e5565b820191906000526020600020905b8154815290600101906020018083116104c857829003601f168201915b505050505081565b6104f56106bf565b82600160008281526001602052604090205460ff16600381111561051557fe5b1461051f57600080fd5b6040805184815290518591859160029160208082019290918190038201816000865af1158015610553573d6000803e3d6000fd5b5050506040513d602081101561056857600080fd5b50516000838152602081905260409020600401541461058657600080fd5b600086815260208181526040808320815160c08101835281548152600180830154828601908152600280850154600160a060020a03908116858801526003860154811660608601908152600487015460808701526005909601805460a08701528f8a528e905592909652848720805460ff191690961790955591519351925190985092169281156108fc029290818181858888f19350505050158015610630573d6000803e3d6000fd5b50505050505050565b60008060008060006106496106bf565b50505060009384525050506020818152604091829020825160c081018452815480825260018301549382018490526002830154600160a060020a039081169583018690526003840154166060830181905260048401546080840181905260059094015460a0909301929092529492939092909190565b6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a0810191909152905600a165627a7a723058206978e11ad984d6cfd355c27bf8426340efe2395e71d2997d337ae7de8029b7fb0029`

// DeployAtomSwap deploys a new Ethereum contract, binding an instance of AtomSwap to it.
func DeployAtomSwap(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AtomSwap, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomSwapABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AtomSwapBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AtomSwap{AtomSwapCaller: AtomSwapCaller{contract: contract}, AtomSwapTransactor: AtomSwapTransactor{contract: contract}, AtomSwapFilterer: AtomSwapFilterer{contract: contract}}, nil
}

// AtomSwap is an auto generated Go binding around an Ethereum contract.
type AtomSwap struct {
	AtomSwapCaller     // Read-only binding to the contract
	AtomSwapTransactor // Write-only binding to the contract
	AtomSwapFilterer   // Log filterer for contract events
}

// AtomSwapCaller is an auto generated read-only Go binding around an Ethereum contract.
type AtomSwapCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomSwapTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AtomSwapTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomSwapFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AtomSwapFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomSwapSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AtomSwapSession struct {
	Contract     *AtomSwap         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AtomSwapCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AtomSwapCallerSession struct {
	Contract *AtomSwapCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// AtomSwapTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AtomSwapTransactorSession struct {
	Contract     *AtomSwapTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// AtomSwapRaw is an auto generated low-level Go binding around an Ethereum contract.
type AtomSwapRaw struct {
	Contract *AtomSwap // Generic contract binding to access the raw methods on
}

// AtomSwapCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AtomSwapCallerRaw struct {
	Contract *AtomSwapCaller // Generic read-only contract binding to access the raw methods on
}

// AtomSwapTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AtomSwapTransactorRaw struct {
	Contract *AtomSwapTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAtomSwap creates a new instance of AtomSwap, bound to a specific deployed contract.
func NewAtomSwap(address common.Address, backend bind.ContractBackend) (*AtomSwap, error) {
	contract, err := bindAtomSwap(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AtomSwap{AtomSwapCaller: AtomSwapCaller{contract: contract}, AtomSwapTransactor: AtomSwapTransactor{contract: contract}, AtomSwapFilterer: AtomSwapFilterer{contract: contract}}, nil
}

// NewAtomSwapCaller creates a new read-only instance of AtomSwap, bound to a specific deployed contract.
func NewAtomSwapCaller(address common.Address, caller bind.ContractCaller) (*AtomSwapCaller, error) {
	contract, err := bindAtomSwap(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AtomSwapCaller{contract: contract}, nil
}

// NewAtomSwapTransactor creates a new write-only instance of AtomSwap, bound to a specific deployed contract.
func NewAtomSwapTransactor(address common.Address, transactor bind.ContractTransactor) (*AtomSwapTransactor, error) {
	contract, err := bindAtomSwap(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AtomSwapTransactor{contract: contract}, nil
}

// NewAtomSwapFilterer creates a new log filterer instance of AtomSwap, bound to a specific deployed contract.
func NewAtomSwapFilterer(address common.Address, filterer bind.ContractFilterer) (*AtomSwapFilterer, error) {
	contract, err := bindAtomSwap(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AtomSwapFilterer{contract: contract}, nil
}

// bindAtomSwap binds a generic wrapper to an already deployed contract.
func bindAtomSwap(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomSwapABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomSwap *AtomSwapRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomSwap.Contract.AtomSwapCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomSwap *AtomSwapRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomSwap.Contract.AtomSwapTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomSwap *AtomSwapRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomSwap.Contract.AtomSwapTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomSwap *AtomSwapCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomSwap.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomSwap *AtomSwapTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomSwap.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomSwap *AtomSwapTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomSwap.Contract.contract.Transact(opts, method, params...)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, from address, secretLock bytes32)
func (_AtomSwap *AtomSwapCaller) Audit(opts *bind.CallOpts, _swapID [32]byte) (struct {
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
	err := _AtomSwap.contract.Call(opts, out, "audit", _swapID)
	return *ret, err
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, from address, secretLock bytes32)
func (_AtomSwap *AtomSwapSession) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	return _AtomSwap.Contract.Audit(&_AtomSwap.CallOpts, _swapID)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, from address, secretLock bytes32)
func (_AtomSwap *AtomSwapCallerSession) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	return _AtomSwap.Contract.Audit(&_AtomSwap.CallOpts, _swapID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_AtomSwap *AtomSwapCaller) AuditSecret(opts *bind.CallOpts, _swapID [32]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _AtomSwap.contract.Call(opts, out, "auditSecret", _swapID)
	return *ret0, err
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_AtomSwap *AtomSwapSession) AuditSecret(_swapID [32]byte) ([32]byte, error) {
	return _AtomSwap.Contract.AuditSecret(&_AtomSwap.CallOpts, _swapID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_AtomSwap *AtomSwapCallerSession) AuditSecret(_swapID [32]byte) ([32]byte, error) {
	return _AtomSwap.Contract.AuditSecret(&_AtomSwap.CallOpts, _swapID)
}

// SwapDetails is a free data retrieval call binding the contract method 0xb14631bb.
//
// Solidity: function swapDetails( bytes32) constant returns(bytes)
func (_AtomSwap *AtomSwapCaller) SwapDetails(opts *bind.CallOpts, arg0 [32]byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _AtomSwap.contract.Call(opts, out, "swapDetails", arg0)
	return *ret0, err
}

// SwapDetails is a free data retrieval call binding the contract method 0xb14631bb.
//
// Solidity: function swapDetails( bytes32) constant returns(bytes)
func (_AtomSwap *AtomSwapSession) SwapDetails(arg0 [32]byte) ([]byte, error) {
	return _AtomSwap.Contract.SwapDetails(&_AtomSwap.CallOpts, arg0)
}

// SwapDetails is a free data retrieval call binding the contract method 0xb14631bb.
//
// Solidity: function swapDetails( bytes32) constant returns(bytes)
func (_AtomSwap *AtomSwapCallerSession) SwapDetails(arg0 [32]byte) ([]byte, error) {
	return _AtomSwap.Contract.SwapDetails(&_AtomSwap.CallOpts, arg0)
}

// Initiate is a paid mutator transaction binding the contract method 0x412c0b58.
//
// Solidity: function initiate(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_AtomSwap *AtomSwapTransactor) Initiate(opts *bind.TransactOpts, _swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _AtomSwap.contract.Transact(opts, "initiate", _swapID, _withdrawTrader, _secretLock, _timelock)
}

// Initiate is a paid mutator transaction binding the contract method 0x412c0b58.
//
// Solidity: function initiate(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_AtomSwap *AtomSwapSession) Initiate(_swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _AtomSwap.Contract.Initiate(&_AtomSwap.TransactOpts, _swapID, _withdrawTrader, _secretLock, _timelock)
}

// Initiate is a paid mutator transaction binding the contract method 0x412c0b58.
//
// Solidity: function initiate(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_AtomSwap *AtomSwapTransactorSession) Initiate(_swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _AtomSwap.Contract.Initiate(&_AtomSwap.TransactOpts, _swapID, _withdrawTrader, _secretLock, _timelock)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_AtomSwap *AtomSwapTransactor) Redeem(opts *bind.TransactOpts, _swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _AtomSwap.contract.Transact(opts, "redeem", _swapID, _secretKey)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_AtomSwap *AtomSwapSession) Redeem(_swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _AtomSwap.Contract.Redeem(&_AtomSwap.TransactOpts, _swapID, _secretKey)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_AtomSwap *AtomSwapTransactorSession) Redeem(_swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _AtomSwap.Contract.Redeem(&_AtomSwap.TransactOpts, _swapID, _secretKey)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_AtomSwap *AtomSwapTransactor) Refund(opts *bind.TransactOpts, _swapID [32]byte) (*types.Transaction, error) {
	return _AtomSwap.contract.Transact(opts, "refund", _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_AtomSwap *AtomSwapSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _AtomSwap.Contract.Refund(&_AtomSwap.TransactOpts, _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_AtomSwap *AtomSwapTransactorSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _AtomSwap.Contract.Refund(&_AtomSwap.TransactOpts, _swapID)
}

// AtomSwapCloseIterator is returned from FilterClose and is used to iterate over the raw logs and unpacked data for Close events raised by the AtomSwap contract.
type AtomSwapCloseIterator struct {
	Event *AtomSwapClose // Event containing the contract specifics and raw log

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
func (it *AtomSwapCloseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomSwapClose)
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
		it.Event = new(AtomSwapClose)
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
func (it *AtomSwapCloseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomSwapCloseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomSwapClose represents a Close event raised by the AtomSwap contract.
type AtomSwapClose struct {
	SwapID    [32]byte
	SecretKey []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterClose is a free log retrieval operation binding the contract event 0x692fd10a275135b9a2a2f5819db3d9965a5129ea2ad3640a0156dbce2fc81bdd.
//
// Solidity: e Close(_swapID bytes32, _secretKey bytes)
func (_AtomSwap *AtomSwapFilterer) FilterClose(opts *bind.FilterOpts) (*AtomSwapCloseIterator, error) {

	logs, sub, err := _AtomSwap.contract.FilterLogs(opts, "Close")
	if err != nil {
		return nil, err
	}
	return &AtomSwapCloseIterator{contract: _AtomSwap.contract, event: "Close", logs: logs, sub: sub}, nil
}

// WatchClose is a free log subscription operation binding the contract event 0x692fd10a275135b9a2a2f5819db3d9965a5129ea2ad3640a0156dbce2fc81bdd.
//
// Solidity: e Close(_swapID bytes32, _secretKey bytes)
func (_AtomSwap *AtomSwapFilterer) WatchClose(opts *bind.WatchOpts, sink chan<- *AtomSwapClose) (event.Subscription, error) {

	logs, sub, err := _AtomSwap.contract.WatchLogs(opts, "Close")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomSwapClose)
				if err := _AtomSwap.contract.UnpackLog(event, "Close", log); err != nil {
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

// AtomSwapExpireIterator is returned from FilterExpire and is used to iterate over the raw logs and unpacked data for Expire events raised by the AtomSwap contract.
type AtomSwapExpireIterator struct {
	Event *AtomSwapExpire // Event containing the contract specifics and raw log

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
func (it *AtomSwapExpireIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomSwapExpire)
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
		it.Event = new(AtomSwapExpire)
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
func (it *AtomSwapExpireIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomSwapExpireIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomSwapExpire represents a Expire event raised by the AtomSwap contract.
type AtomSwapExpire struct {
	SwapID [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterExpire is a free log retrieval operation binding the contract event 0xbddd9b693ea862fad6ecf78fd51c065be26fda94d1f3cad3a7d691453a38a735.
//
// Solidity: e Expire(_swapID bytes32)
func (_AtomSwap *AtomSwapFilterer) FilterExpire(opts *bind.FilterOpts) (*AtomSwapExpireIterator, error) {

	logs, sub, err := _AtomSwap.contract.FilterLogs(opts, "Expire")
	if err != nil {
		return nil, err
	}
	return &AtomSwapExpireIterator{contract: _AtomSwap.contract, event: "Expire", logs: logs, sub: sub}, nil
}

// WatchExpire is a free log subscription operation binding the contract event 0xbddd9b693ea862fad6ecf78fd51c065be26fda94d1f3cad3a7d691453a38a735.
//
// Solidity: e Expire(_swapID bytes32)
func (_AtomSwap *AtomSwapFilterer) WatchExpire(opts *bind.WatchOpts, sink chan<- *AtomSwapExpire) (event.Subscription, error) {

	logs, sub, err := _AtomSwap.contract.WatchLogs(opts, "Expire")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomSwapExpire)
				if err := _AtomSwap.contract.UnpackLog(event, "Expire", log); err != nil {
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

// AtomSwapOpenIterator is returned from FilterOpen and is used to iterate over the raw logs and unpacked data for Open events raised by the AtomSwap contract.
type AtomSwapOpenIterator struct {
	Event *AtomSwapOpen // Event containing the contract specifics and raw log

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
func (it *AtomSwapOpenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomSwapOpen)
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
		it.Event = new(AtomSwapOpen)
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
func (it *AtomSwapOpenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomSwapOpenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomSwapOpen represents a Open event raised by the AtomSwap contract.
type AtomSwapOpen struct {
	SwapID         [32]byte
	WithdrawTrader common.Address
	SecretLock     [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterOpen is a free log retrieval operation binding the contract event 0x6ed79a08bf5c8a7d4a330df315e4ac386627ecafbe5d2bfd6654237d967b24f3.
//
// Solidity: e Open(_swapID bytes32, _withdrawTrader address, _secretLock bytes32)
func (_AtomSwap *AtomSwapFilterer) FilterOpen(opts *bind.FilterOpts) (*AtomSwapOpenIterator, error) {

	logs, sub, err := _AtomSwap.contract.FilterLogs(opts, "Open")
	if err != nil {
		return nil, err
	}
	return &AtomSwapOpenIterator{contract: _AtomSwap.contract, event: "Open", logs: logs, sub: sub}, nil
}

// WatchOpen is a free log subscription operation binding the contract event 0x6ed79a08bf5c8a7d4a330df315e4ac386627ecafbe5d2bfd6654237d967b24f3.
//
// Solidity: e Open(_swapID bytes32, _withdrawTrader address, _secretLock bytes32)
func (_AtomSwap *AtomSwapFilterer) WatchOpen(opts *bind.WatchOpts, sink chan<- *AtomSwapOpen) (event.Subscription, error) {

	logs, sub, err := _AtomSwap.contract.WatchLogs(opts, "Open")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomSwapOpen)
				if err := _AtomSwap.contract.UnpackLog(event, "Open", log); err != nil {
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
