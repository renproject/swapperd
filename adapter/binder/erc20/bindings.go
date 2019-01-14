// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package erc20

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

// CompatibleERC20ABI is the input ABI used to generate the binding from.
const CompatibleERC20ABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]"

// CompatibleERC20Bin is the compiled bytecode used for deploying new contracts.
const CompatibleERC20Bin = `0x`

// DeployCompatibleERC20 deploys a new Ethereum contract, binding an instance of CompatibleERC20 to it.
func DeployCompatibleERC20(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CompatibleERC20, error) {
	parsed, err := abi.JSON(strings.NewReader(CompatibleERC20ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(CompatibleERC20Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CompatibleERC20{CompatibleERC20Caller: CompatibleERC20Caller{contract: contract}, CompatibleERC20Transactor: CompatibleERC20Transactor{contract: contract}, CompatibleERC20Filterer: CompatibleERC20Filterer{contract: contract}}, nil
}

// CompatibleERC20 is an auto generated Go binding around an Ethereum contract.
type CompatibleERC20 struct {
	CompatibleERC20Caller     // Read-only binding to the contract
	CompatibleERC20Transactor // Write-only binding to the contract
	CompatibleERC20Filterer   // Log filterer for contract events
}

// CompatibleERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type CompatibleERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompatibleERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type CompatibleERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompatibleERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CompatibleERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompatibleERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CompatibleERC20Session struct {
	Contract     *CompatibleERC20  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CompatibleERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CompatibleERC20CallerSession struct {
	Contract *CompatibleERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// CompatibleERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CompatibleERC20TransactorSession struct {
	Contract     *CompatibleERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// CompatibleERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type CompatibleERC20Raw struct {
	Contract *CompatibleERC20 // Generic contract binding to access the raw methods on
}

// CompatibleERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CompatibleERC20CallerRaw struct {
	Contract *CompatibleERC20Caller // Generic read-only contract binding to access the raw methods on
}

// CompatibleERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CompatibleERC20TransactorRaw struct {
	Contract *CompatibleERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewCompatibleERC20 creates a new instance of CompatibleERC20, bound to a specific deployed contract.
func NewCompatibleERC20(address common.Address, backend bind.ContractBackend) (*CompatibleERC20, error) {
	contract, err := bindCompatibleERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CompatibleERC20{CompatibleERC20Caller: CompatibleERC20Caller{contract: contract}, CompatibleERC20Transactor: CompatibleERC20Transactor{contract: contract}, CompatibleERC20Filterer: CompatibleERC20Filterer{contract: contract}}, nil
}

// NewCompatibleERC20Caller creates a new read-only instance of CompatibleERC20, bound to a specific deployed contract.
func NewCompatibleERC20Caller(address common.Address, caller bind.ContractCaller) (*CompatibleERC20Caller, error) {
	contract, err := bindCompatibleERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CompatibleERC20Caller{contract: contract}, nil
}

// NewCompatibleERC20Transactor creates a new write-only instance of CompatibleERC20, bound to a specific deployed contract.
func NewCompatibleERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*CompatibleERC20Transactor, error) {
	contract, err := bindCompatibleERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CompatibleERC20Transactor{contract: contract}, nil
}

// NewCompatibleERC20Filterer creates a new log filterer instance of CompatibleERC20, bound to a specific deployed contract.
func NewCompatibleERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*CompatibleERC20Filterer, error) {
	contract, err := bindCompatibleERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CompatibleERC20Filterer{contract: contract}, nil
}

// bindCompatibleERC20 binds a generic wrapper to an already deployed contract.
func bindCompatibleERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CompatibleERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CompatibleERC20 *CompatibleERC20Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CompatibleERC20.Contract.CompatibleERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CompatibleERC20 *CompatibleERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.CompatibleERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CompatibleERC20 *CompatibleERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.CompatibleERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CompatibleERC20 *CompatibleERC20CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CompatibleERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CompatibleERC20 *CompatibleERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CompatibleERC20 *CompatibleERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CompatibleERC20.contract.Call(opts, out, "allowance", owner, spender)
	return *ret0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _CompatibleERC20.Contract.Allowance(&_CompatibleERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(owner address, spender address) constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _CompatibleERC20.Contract.Allowance(&_CompatibleERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20Caller) BalanceOf(opts *bind.CallOpts, who common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CompatibleERC20.contract.Call(opts, out, "balanceOf", who)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20Session) BalanceOf(who common.Address) (*big.Int, error) {
	return _CompatibleERC20.Contract.BalanceOf(&_CompatibleERC20.CallOpts, who)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(who address) constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20CallerSession) BalanceOf(who common.Address) (*big.Int, error) {
	return _CompatibleERC20.Contract.BalanceOf(&_CompatibleERC20.CallOpts, who)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CompatibleERC20.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20Session) TotalSupply() (*big.Int, error) {
	return _CompatibleERC20.Contract.TotalSupply(&_CompatibleERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_CompatibleERC20 *CompatibleERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _CompatibleERC20.Contract.TotalSupply(&_CompatibleERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.Approve(&_CompatibleERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(spender address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.Approve(&_CompatibleERC20.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.Transfer(&_CompatibleERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(to address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.Transfer(&_CompatibleERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.TransferFrom(&_CompatibleERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(from address, to address, value uint256) returns()
func (_CompatibleERC20 *CompatibleERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _CompatibleERC20.Contract.TransferFrom(&_CompatibleERC20.TransactOpts, from, to, value)
}

// CompatibleERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the CompatibleERC20 contract.
type CompatibleERC20ApprovalIterator struct {
	Event *CompatibleERC20Approval // Event containing the contract specifics and raw log

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
func (it *CompatibleERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompatibleERC20Approval)
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
		it.Event = new(CompatibleERC20Approval)
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
func (it *CompatibleERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompatibleERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompatibleERC20Approval represents a Approval event raised by the CompatibleERC20 contract.
type CompatibleERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_CompatibleERC20 *CompatibleERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*CompatibleERC20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _CompatibleERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &CompatibleERC20ApprovalIterator{contract: _CompatibleERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: e Approval(owner indexed address, spender indexed address, value uint256)
func (_CompatibleERC20 *CompatibleERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *CompatibleERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _CompatibleERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompatibleERC20Approval)
				if err := _CompatibleERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

// CompatibleERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the CompatibleERC20 contract.
type CompatibleERC20TransferIterator struct {
	Event *CompatibleERC20Transfer // Event containing the contract specifics and raw log

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
func (it *CompatibleERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CompatibleERC20Transfer)
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
		it.Event = new(CompatibleERC20Transfer)
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
func (it *CompatibleERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CompatibleERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CompatibleERC20Transfer represents a Transfer event raised by the CompatibleERC20 contract.
type CompatibleERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_CompatibleERC20 *CompatibleERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*CompatibleERC20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CompatibleERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &CompatibleERC20TransferIterator{contract: _CompatibleERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: e Transfer(from indexed address, to indexed address, value uint256)
func (_CompatibleERC20 *CompatibleERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *CompatibleERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _CompatibleERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CompatibleERC20Transfer)
				if err := _CompatibleERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ERC20SwapContractABI is the input ABI used to generate the binding from.
const ERC20SwapContractABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"initiate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"initiatable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"}],\"name\":\"swapID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawBrokerFees\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"redeemable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"auditSecret\",\"outputs\":[{\"name\":\"secretKey\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refundable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_broker\",\"type\":\"address\"},{\"name\":\"_brokerFee\",\"type\":\"uint256\"},{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"initiateWithFees\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"audit\",\"outputs\":[{\"name\":\"timelock\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"brokerFee\",\"type\":\"uint256\"},{\"name\":\"broker\",\"type\":\"address\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"secretLock\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_receiver\",\"type\":\"address\"},{\"name\":\"_secretKey\",\"type\":\"bytes32\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ERC20SwapContractBin is the compiled bytecode used for deploying new contracts.
const ERC20SwapContractBin = `0x`

// DeployERC20SwapContract deploys a new Ethereum contract, binding an instance of ERC20SwapContract to it.
func DeployERC20SwapContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC20SwapContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20SwapContractABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ERC20SwapContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC20SwapContract{ERC20SwapContractCaller: ERC20SwapContractCaller{contract: contract}, ERC20SwapContractTransactor: ERC20SwapContractTransactor{contract: contract}, ERC20SwapContractFilterer: ERC20SwapContractFilterer{contract: contract}}, nil
}

// ERC20SwapContract is an auto generated Go binding around an Ethereum contract.
type ERC20SwapContract struct {
	ERC20SwapContractCaller     // Read-only binding to the contract
	ERC20SwapContractTransactor // Write-only binding to the contract
	ERC20SwapContractFilterer   // Log filterer for contract events
}

// ERC20SwapContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20SwapContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20SwapContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20SwapContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20SwapContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20SwapContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20SwapContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20SwapContractSession struct {
	Contract     *ERC20SwapContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ERC20SwapContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20SwapContractCallerSession struct {
	Contract *ERC20SwapContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ERC20SwapContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20SwapContractTransactorSession struct {
	Contract     *ERC20SwapContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ERC20SwapContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20SwapContractRaw struct {
	Contract *ERC20SwapContract // Generic contract binding to access the raw methods on
}

// ERC20SwapContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20SwapContractCallerRaw struct {
	Contract *ERC20SwapContractCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20SwapContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20SwapContractTransactorRaw struct {
	Contract *ERC20SwapContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20SwapContract creates a new instance of ERC20SwapContract, bound to a specific deployed contract.
func NewERC20SwapContract(address common.Address, backend bind.ContractBackend) (*ERC20SwapContract, error) {
	contract, err := bindERC20SwapContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20SwapContract{ERC20SwapContractCaller: ERC20SwapContractCaller{contract: contract}, ERC20SwapContractTransactor: ERC20SwapContractTransactor{contract: contract}, ERC20SwapContractFilterer: ERC20SwapContractFilterer{contract: contract}}, nil
}

// NewERC20SwapContractCaller creates a new read-only instance of ERC20SwapContract, bound to a specific deployed contract.
func NewERC20SwapContractCaller(address common.Address, caller bind.ContractCaller) (*ERC20SwapContractCaller, error) {
	contract, err := bindERC20SwapContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20SwapContractCaller{contract: contract}, nil
}

// NewERC20SwapContractTransactor creates a new write-only instance of ERC20SwapContract, bound to a specific deployed contract.
func NewERC20SwapContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20SwapContractTransactor, error) {
	contract, err := bindERC20SwapContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20SwapContractTransactor{contract: contract}, nil
}

// NewERC20SwapContractFilterer creates a new log filterer instance of ERC20SwapContract, bound to a specific deployed contract.
func NewERC20SwapContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20SwapContractFilterer, error) {
	contract, err := bindERC20SwapContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20SwapContractFilterer{contract: contract}, nil
}

// bindERC20SwapContract binds a generic wrapper to an already deployed contract.
func bindERC20SwapContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20SwapContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20SwapContract *ERC20SwapContractRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20SwapContract.Contract.ERC20SwapContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20SwapContract *ERC20SwapContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20SwapContract.Contract.ERC20SwapContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20SwapContract *ERC20SwapContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20SwapContract.Contract.ERC20SwapContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20SwapContract *ERC20SwapContractCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ERC20SwapContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20SwapContract *ERC20SwapContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20SwapContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20SwapContract *ERC20SwapContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20SwapContract.Contract.contract.Transact(opts, method, params...)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, brokerFee uint256, broker address, from address, secretLock bytes32)
func (_ERC20SwapContract *ERC20SwapContractCaller) Audit(opts *bind.CallOpts, _swapID [32]byte) (struct {
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
	err := _ERC20SwapContract.contract.Call(opts, out, "audit", _swapID)
	return *ret, err
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, brokerFee uint256, broker address, from address, secretLock bytes32)
func (_ERC20SwapContract *ERC20SwapContractSession) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	BrokerFee  *big.Int
	Broker     common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	return _ERC20SwapContract.Contract.Audit(&_ERC20SwapContract.CallOpts, _swapID)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, brokerFee uint256, broker address, from address, secretLock bytes32)
func (_ERC20SwapContract *ERC20SwapContractCallerSession) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	BrokerFee  *big.Int
	Broker     common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	return _ERC20SwapContract.Contract.Audit(&_ERC20SwapContract.CallOpts, _swapID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_ERC20SwapContract *ERC20SwapContractCaller) AuditSecret(opts *bind.CallOpts, _swapID [32]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _ERC20SwapContract.contract.Call(opts, out, "auditSecret", _swapID)
	return *ret0, err
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_ERC20SwapContract *ERC20SwapContractSession) AuditSecret(_swapID [32]byte) ([32]byte, error) {
	return _ERC20SwapContract.Contract.AuditSecret(&_ERC20SwapContract.CallOpts, _swapID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_ERC20SwapContract *ERC20SwapContractCallerSession) AuditSecret(_swapID [32]byte) ([32]byte, error) {
	return _ERC20SwapContract.Contract.AuditSecret(&_ERC20SwapContract.CallOpts, _swapID)
}

// Initiatable is a free data retrieval call binding the contract method 0x09ece618.
//
// Solidity: function initiatable(_swapID bytes32) constant returns(bool)
func (_ERC20SwapContract *ERC20SwapContractCaller) Initiatable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ERC20SwapContract.contract.Call(opts, out, "initiatable", _swapID)
	return *ret0, err
}

// Initiatable is a free data retrieval call binding the contract method 0x09ece618.
//
// Solidity: function initiatable(_swapID bytes32) constant returns(bool)
func (_ERC20SwapContract *ERC20SwapContractSession) Initiatable(_swapID [32]byte) (bool, error) {
	return _ERC20SwapContract.Contract.Initiatable(&_ERC20SwapContract.CallOpts, _swapID)
}

// Initiatable is a free data retrieval call binding the contract method 0x09ece618.
//
// Solidity: function initiatable(_swapID bytes32) constant returns(bool)
func (_ERC20SwapContract *ERC20SwapContractCallerSession) Initiatable(_swapID [32]byte) (bool, error) {
	return _ERC20SwapContract.Contract.Initiatable(&_ERC20SwapContract.CallOpts, _swapID)
}

// Redeemable is a free data retrieval call binding the contract method 0x68f06b29.
//
// Solidity: function redeemable(_swapID bytes32) constant returns(bool)
func (_ERC20SwapContract *ERC20SwapContractCaller) Redeemable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ERC20SwapContract.contract.Call(opts, out, "redeemable", _swapID)
	return *ret0, err
}

// Redeemable is a free data retrieval call binding the contract method 0x68f06b29.
//
// Solidity: function redeemable(_swapID bytes32) constant returns(bool)
func (_ERC20SwapContract *ERC20SwapContractSession) Redeemable(_swapID [32]byte) (bool, error) {
	return _ERC20SwapContract.Contract.Redeemable(&_ERC20SwapContract.CallOpts, _swapID)
}

// Redeemable is a free data retrieval call binding the contract method 0x68f06b29.
//
// Solidity: function redeemable(_swapID bytes32) constant returns(bool)
func (_ERC20SwapContract *ERC20SwapContractCallerSession) Redeemable(_swapID [32]byte) (bool, error) {
	return _ERC20SwapContract.Contract.Redeemable(&_ERC20SwapContract.CallOpts, _swapID)
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(_swapID bytes32) constant returns(bool)
func (_ERC20SwapContract *ERC20SwapContractCaller) Refundable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ERC20SwapContract.contract.Call(opts, out, "refundable", _swapID)
	return *ret0, err
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(_swapID bytes32) constant returns(bool)
func (_ERC20SwapContract *ERC20SwapContractSession) Refundable(_swapID [32]byte) (bool, error) {
	return _ERC20SwapContract.Contract.Refundable(&_ERC20SwapContract.CallOpts, _swapID)
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(_swapID bytes32) constant returns(bool)
func (_ERC20SwapContract *ERC20SwapContractCallerSession) Refundable(_swapID [32]byte) (bool, error) {
	return _ERC20SwapContract.Contract.Refundable(&_ERC20SwapContract.CallOpts, _swapID)
}

// SwapID is a free data retrieval call binding the contract method 0x4b2ac3fa.
//
// Solidity: function swapID(_secretLock bytes32, _timelock uint256) constant returns(bytes32)
func (_ERC20SwapContract *ERC20SwapContractCaller) SwapID(opts *bind.CallOpts, _secretLock [32]byte, _timelock *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _ERC20SwapContract.contract.Call(opts, out, "swapID", _secretLock, _timelock)
	return *ret0, err
}

// SwapID is a free data retrieval call binding the contract method 0x4b2ac3fa.
//
// Solidity: function swapID(_secretLock bytes32, _timelock uint256) constant returns(bytes32)
func (_ERC20SwapContract *ERC20SwapContractSession) SwapID(_secretLock [32]byte, _timelock *big.Int) ([32]byte, error) {
	return _ERC20SwapContract.Contract.SwapID(&_ERC20SwapContract.CallOpts, _secretLock, _timelock)
}

// SwapID is a free data retrieval call binding the contract method 0x4b2ac3fa.
//
// Solidity: function swapID(_secretLock bytes32, _timelock uint256) constant returns(bytes32)
func (_ERC20SwapContract *ERC20SwapContractCallerSession) SwapID(_secretLock [32]byte, _timelock *big.Int) ([32]byte, error) {
	return _ERC20SwapContract.Contract.SwapID(&_ERC20SwapContract.CallOpts, _secretLock, _timelock)
}

// Initiate is a paid mutator transaction binding the contract method 0x027a2577.
//
// Solidity: function initiate(_swapID bytes32, _spender address, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_ERC20SwapContract *ERC20SwapContractTransactor) Initiate(opts *bind.TransactOpts, _swapID [32]byte, _spender common.Address, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _ERC20SwapContract.contract.Transact(opts, "initiate", _swapID, _spender, _secretLock, _timelock, _value)
}

// Initiate is a paid mutator transaction binding the contract method 0x027a2577.
//
// Solidity: function initiate(_swapID bytes32, _spender address, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_ERC20SwapContract *ERC20SwapContractSession) Initiate(_swapID [32]byte, _spender common.Address, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _ERC20SwapContract.Contract.Initiate(&_ERC20SwapContract.TransactOpts, _swapID, _spender, _secretLock, _timelock, _value)
}

// Initiate is a paid mutator transaction binding the contract method 0x027a2577.
//
// Solidity: function initiate(_swapID bytes32, _spender address, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_ERC20SwapContract *ERC20SwapContractTransactorSession) Initiate(_swapID [32]byte, _spender common.Address, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _ERC20SwapContract.Contract.Initiate(&_ERC20SwapContract.TransactOpts, _swapID, _spender, _secretLock, _timelock, _value)
}

// InitiateWithFees is a paid mutator transaction binding the contract method 0xb8688e3f.
//
// Solidity: function initiateWithFees(_swapID bytes32, _spender address, _broker address, _brokerFee uint256, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_ERC20SwapContract *ERC20SwapContractTransactor) InitiateWithFees(opts *bind.TransactOpts, _swapID [32]byte, _spender common.Address, _broker common.Address, _brokerFee *big.Int, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _ERC20SwapContract.contract.Transact(opts, "initiateWithFees", _swapID, _spender, _broker, _brokerFee, _secretLock, _timelock, _value)
}

// InitiateWithFees is a paid mutator transaction binding the contract method 0xb8688e3f.
//
// Solidity: function initiateWithFees(_swapID bytes32, _spender address, _broker address, _brokerFee uint256, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_ERC20SwapContract *ERC20SwapContractSession) InitiateWithFees(_swapID [32]byte, _spender common.Address, _broker common.Address, _brokerFee *big.Int, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _ERC20SwapContract.Contract.InitiateWithFees(&_ERC20SwapContract.TransactOpts, _swapID, _spender, _broker, _brokerFee, _secretLock, _timelock, _value)
}

// InitiateWithFees is a paid mutator transaction binding the contract method 0xb8688e3f.
//
// Solidity: function initiateWithFees(_swapID bytes32, _spender address, _broker address, _brokerFee uint256, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_ERC20SwapContract *ERC20SwapContractTransactorSession) InitiateWithFees(_swapID [32]byte, _spender common.Address, _broker common.Address, _brokerFee *big.Int, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _ERC20SwapContract.Contract.InitiateWithFees(&_ERC20SwapContract.TransactOpts, _swapID, _spender, _broker, _brokerFee, _secretLock, _timelock, _value)
}

// Redeem is a paid mutator transaction binding the contract method 0xc23b1a85.
//
// Solidity: function redeem(_swapID bytes32, _receiver address, _secretKey bytes32) returns()
func (_ERC20SwapContract *ERC20SwapContractTransactor) Redeem(opts *bind.TransactOpts, _swapID [32]byte, _receiver common.Address, _secretKey [32]byte) (*types.Transaction, error) {
	return _ERC20SwapContract.contract.Transact(opts, "redeem", _swapID, _receiver, _secretKey)
}

// Redeem is a paid mutator transaction binding the contract method 0xc23b1a85.
//
// Solidity: function redeem(_swapID bytes32, _receiver address, _secretKey bytes32) returns()
func (_ERC20SwapContract *ERC20SwapContractSession) Redeem(_swapID [32]byte, _receiver common.Address, _secretKey [32]byte) (*types.Transaction, error) {
	return _ERC20SwapContract.Contract.Redeem(&_ERC20SwapContract.TransactOpts, _swapID, _receiver, _secretKey)
}

// Redeem is a paid mutator transaction binding the contract method 0xc23b1a85.
//
// Solidity: function redeem(_swapID bytes32, _receiver address, _secretKey bytes32) returns()
func (_ERC20SwapContract *ERC20SwapContractTransactorSession) Redeem(_swapID [32]byte, _receiver common.Address, _secretKey [32]byte) (*types.Transaction, error) {
	return _ERC20SwapContract.Contract.Redeem(&_ERC20SwapContract.TransactOpts, _swapID, _receiver, _secretKey)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_ERC20SwapContract *ERC20SwapContractTransactor) Refund(opts *bind.TransactOpts, _swapID [32]byte) (*types.Transaction, error) {
	return _ERC20SwapContract.contract.Transact(opts, "refund", _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_ERC20SwapContract *ERC20SwapContractSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _ERC20SwapContract.Contract.Refund(&_ERC20SwapContract.TransactOpts, _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_ERC20SwapContract *ERC20SwapContractTransactorSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _ERC20SwapContract.Contract.Refund(&_ERC20SwapContract.TransactOpts, _swapID)
}

// WithdrawBrokerFees is a paid mutator transaction binding the contract method 0x4c6d37ff.
//
// Solidity: function withdrawBrokerFees(_amount uint256) returns()
func (_ERC20SwapContract *ERC20SwapContractTransactor) WithdrawBrokerFees(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _ERC20SwapContract.contract.Transact(opts, "withdrawBrokerFees", _amount)
}

// WithdrawBrokerFees is a paid mutator transaction binding the contract method 0x4c6d37ff.
//
// Solidity: function withdrawBrokerFees(_amount uint256) returns()
func (_ERC20SwapContract *ERC20SwapContractSession) WithdrawBrokerFees(_amount *big.Int) (*types.Transaction, error) {
	return _ERC20SwapContract.Contract.WithdrawBrokerFees(&_ERC20SwapContract.TransactOpts, _amount)
}

// WithdrawBrokerFees is a paid mutator transaction binding the contract method 0x4c6d37ff.
//
// Solidity: function withdrawBrokerFees(_amount uint256) returns()
func (_ERC20SwapContract *ERC20SwapContractTransactorSession) WithdrawBrokerFees(_amount *big.Int) (*types.Transaction, error) {
	return _ERC20SwapContract.Contract.WithdrawBrokerFees(&_ERC20SwapContract.TransactOpts, _amount)
}
