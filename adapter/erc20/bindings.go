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

// CompatibleERC20FunctionsABI is the input ABI used to generate the binding from.
const CompatibleERC20FunctionsABI = "[]"

// CompatibleERC20FunctionsBin is the compiled bytecode used for deploying new contracts.
const CompatibleERC20FunctionsBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820dd6efa34d584fcd8e361c595eefdcf33cc322c6677015364bc1d8aba40f1fff30029`

// DeployCompatibleERC20Functions deploys a new Ethereum contract, binding an instance of CompatibleERC20Functions to it.
func DeployCompatibleERC20Functions(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CompatibleERC20Functions, error) {
	parsed, err := abi.JSON(strings.NewReader(CompatibleERC20FunctionsABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(CompatibleERC20FunctionsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CompatibleERC20Functions{CompatibleERC20FunctionsCaller: CompatibleERC20FunctionsCaller{contract: contract}, CompatibleERC20FunctionsTransactor: CompatibleERC20FunctionsTransactor{contract: contract}, CompatibleERC20FunctionsFilterer: CompatibleERC20FunctionsFilterer{contract: contract}}, nil
}

// CompatibleERC20Functions is an auto generated Go binding around an Ethereum contract.
type CompatibleERC20Functions struct {
	CompatibleERC20FunctionsCaller     // Read-only binding to the contract
	CompatibleERC20FunctionsTransactor // Write-only binding to the contract
	CompatibleERC20FunctionsFilterer   // Log filterer for contract events
}

// CompatibleERC20FunctionsCaller is an auto generated read-only Go binding around an Ethereum contract.
type CompatibleERC20FunctionsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompatibleERC20FunctionsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CompatibleERC20FunctionsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompatibleERC20FunctionsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CompatibleERC20FunctionsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CompatibleERC20FunctionsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CompatibleERC20FunctionsSession struct {
	Contract     *CompatibleERC20Functions // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// CompatibleERC20FunctionsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CompatibleERC20FunctionsCallerSession struct {
	Contract *CompatibleERC20FunctionsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// CompatibleERC20FunctionsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CompatibleERC20FunctionsTransactorSession struct {
	Contract     *CompatibleERC20FunctionsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// CompatibleERC20FunctionsRaw is an auto generated low-level Go binding around an Ethereum contract.
type CompatibleERC20FunctionsRaw struct {
	Contract *CompatibleERC20Functions // Generic contract binding to access the raw methods on
}

// CompatibleERC20FunctionsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CompatibleERC20FunctionsCallerRaw struct {
	Contract *CompatibleERC20FunctionsCaller // Generic read-only contract binding to access the raw methods on
}

// CompatibleERC20FunctionsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CompatibleERC20FunctionsTransactorRaw struct {
	Contract *CompatibleERC20FunctionsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCompatibleERC20Functions creates a new instance of CompatibleERC20Functions, bound to a specific deployed contract.
func NewCompatibleERC20Functions(address common.Address, backend bind.ContractBackend) (*CompatibleERC20Functions, error) {
	contract, err := bindCompatibleERC20Functions(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CompatibleERC20Functions{CompatibleERC20FunctionsCaller: CompatibleERC20FunctionsCaller{contract: contract}, CompatibleERC20FunctionsTransactor: CompatibleERC20FunctionsTransactor{contract: contract}, CompatibleERC20FunctionsFilterer: CompatibleERC20FunctionsFilterer{contract: contract}}, nil
}

// NewCompatibleERC20FunctionsCaller creates a new read-only instance of CompatibleERC20Functions, bound to a specific deployed contract.
func NewCompatibleERC20FunctionsCaller(address common.Address, caller bind.ContractCaller) (*CompatibleERC20FunctionsCaller, error) {
	contract, err := bindCompatibleERC20Functions(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CompatibleERC20FunctionsCaller{contract: contract}, nil
}

// NewCompatibleERC20FunctionsTransactor creates a new write-only instance of CompatibleERC20Functions, bound to a specific deployed contract.
func NewCompatibleERC20FunctionsTransactor(address common.Address, transactor bind.ContractTransactor) (*CompatibleERC20FunctionsTransactor, error) {
	contract, err := bindCompatibleERC20Functions(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CompatibleERC20FunctionsTransactor{contract: contract}, nil
}

// NewCompatibleERC20FunctionsFilterer creates a new log filterer instance of CompatibleERC20Functions, bound to a specific deployed contract.
func NewCompatibleERC20FunctionsFilterer(address common.Address, filterer bind.ContractFilterer) (*CompatibleERC20FunctionsFilterer, error) {
	contract, err := bindCompatibleERC20Functions(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CompatibleERC20FunctionsFilterer{contract: contract}, nil
}

// bindCompatibleERC20Functions binds a generic wrapper to an already deployed contract.
func bindCompatibleERC20Functions(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CompatibleERC20FunctionsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CompatibleERC20Functions *CompatibleERC20FunctionsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CompatibleERC20Functions.Contract.CompatibleERC20FunctionsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CompatibleERC20Functions *CompatibleERC20FunctionsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CompatibleERC20Functions.Contract.CompatibleERC20FunctionsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CompatibleERC20Functions *CompatibleERC20FunctionsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CompatibleERC20Functions.Contract.CompatibleERC20FunctionsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CompatibleERC20Functions *CompatibleERC20FunctionsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CompatibleERC20Functions.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CompatibleERC20Functions *CompatibleERC20FunctionsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CompatibleERC20Functions.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CompatibleERC20Functions *CompatibleERC20FunctionsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CompatibleERC20Functions.Contract.contract.Transact(opts, method, params...)
}

// MathABI is the input ABI used to generate the binding from.
const MathABI = "[]"

// MathBin is the compiled bytecode used for deploying new contracts.
const MathBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820a98957ed57f49c84fb32ce948c38879d13cf42d23419fb8bdff613b31b9f3e520029`

// DeployMath deploys a new Ethereum contract, binding an instance of Math to it.
func DeployMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Math, error) {
	parsed, err := abi.JSON(strings.NewReader(MathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Math{MathCaller: MathCaller{contract: contract}, MathTransactor: MathTransactor{contract: contract}, MathFilterer: MathFilterer{contract: contract}}, nil
}

// Math is an auto generated Go binding around an Ethereum contract.
type Math struct {
	MathCaller     // Read-only binding to the contract
	MathTransactor // Write-only binding to the contract
	MathFilterer   // Log filterer for contract events
}

// MathCaller is an auto generated read-only Go binding around an Ethereum contract.
type MathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MathSession struct {
	Contract     *Math             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MathCallerSession struct {
	Contract *MathCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// MathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MathTransactorSession struct {
	Contract     *MathTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MathRaw is an auto generated low-level Go binding around an Ethereum contract.
type MathRaw struct {
	Contract *Math // Generic contract binding to access the raw methods on
}

// MathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MathCallerRaw struct {
	Contract *MathCaller // Generic read-only contract binding to access the raw methods on
}

// MathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MathTransactorRaw struct {
	Contract *MathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMath creates a new instance of Math, bound to a specific deployed contract.
func NewMath(address common.Address, backend bind.ContractBackend) (*Math, error) {
	contract, err := bindMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Math{MathCaller: MathCaller{contract: contract}, MathTransactor: MathTransactor{contract: contract}, MathFilterer: MathFilterer{contract: contract}}, nil
}

// NewMathCaller creates a new read-only instance of Math, bound to a specific deployed contract.
func NewMathCaller(address common.Address, caller bind.ContractCaller) (*MathCaller, error) {
	contract, err := bindMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MathCaller{contract: contract}, nil
}

// NewMathTransactor creates a new write-only instance of Math, bound to a specific deployed contract.
func NewMathTransactor(address common.Address, transactor bind.ContractTransactor) (*MathTransactor, error) {
	contract, err := bindMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MathTransactor{contract: contract}, nil
}

// NewMathFilterer creates a new log filterer instance of Math, bound to a specific deployed contract.
func NewMathFilterer(address common.Address, filterer bind.ContractFilterer) (*MathFilterer, error) {
	contract, err := bindMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MathFilterer{contract: contract}, nil
}

// bindMath binds a generic wrapper to an already deployed contract.
func bindMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Math *MathRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Math.Contract.MathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Math *MathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Math.Contract.MathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Math *MathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Math.Contract.MathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Math *MathCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Math.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Math *MathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Math.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Math *MathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Math.Contract.contract.Transact(opts, method, params...)
}

// RenExAtomicSwapperABI is the input ABI used to generate the binding from.
const RenExAtomicSwapperABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_withdrawTrader\",\"type\":\"address\"},{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"initiate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"initiatable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TOKEN_ADDRESS\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"}],\"name\":\"swapID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"redeemable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"auditSecret\",\"outputs\":[{\"name\":\"secretKey\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refundable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_secretKey\",\"type\":\"bytes32\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"redeemedAt\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"audit\",\"outputs\":[{\"name\":\"timelock\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"secretLock\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_VERSION\",\"type\":\"string\"},{\"name\":\"_TOKEN_ADDRESS\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_withdrawTrader\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_secretLock\",\"type\":\"bytes32\"}],\"name\":\"LogOpen\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"LogExpire\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_secretKey\",\"type\":\"bytes32\"}],\"name\":\"LogClose\",\"type\":\"event\"}]"

// RenExAtomicSwapperBin is the compiled bytecode used for deploying new contracts.
const RenExAtomicSwapperBin = `0x608060405234801561001057600080fd5b50604051610d93380380610d9383398101604052805160208083015191909201805190926100439160009185019061006a565b5060018054600160a060020a031916600160a060020a039290921691909117905550610105565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100ab57805160ff19168380011785556100d8565b828001600101855582156100d8579182015b828111156100d85782518255916020019190600101906100bd565b506100e49291506100e8565b5090565b61010291905b808211156100e457600081556001016100ee565b90565b610c7f806101146000396000f3006080604052600436106100b95763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663027a257781146100be57806309ece618146100ed5780630bdf5300146101195780634b2ac3fa1461014a57806368f06b29146101775780637249fbb61461018f578063976d00f4146101a75780639fb31475146101bf578063b31597ad146101d7578063bc4fcc4a146101f2578063c140635b1461020a578063ffa1ad7414610257575b600080fd5b3480156100ca57600080fd5b506100eb600435600160a060020a03602435166044356064356084356102e1565b005b3480156100f957600080fd5b50610105600435610572565b604080519115158252519081900360200190f35b34801561012557600080fd5b5061012e61059c565b60408051600160a060020a039092168252519081900360200190f35b34801561015657600080fd5b506101656004356024356105ab565b60408051918252519081900360200190f35b34801561018357600080fd5b50610105600435610631565b34801561019b57600080fd5b506100eb60043561063a565b3480156101b357600080fd5b50610165600435610802565b3480156101cb57600080fd5b50610105600435610892565b3480156101e357600080fd5b506100eb6004356024356108b8565b3480156101fe57600080fd5b50610165600435610b42565b34801561021657600080fd5b50610222600435610b54565b604080519586526020860194909452600160a060020a0392831685850152911660608401526080830152519081900360a00190f35b34801561026357600080fd5b5061026c610b90565b6040805160208082528351818301528351919283929083019185019080838360005b838110156102a657818101518382015260200161028e565b50505050905090810190601f1680156102d35780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102e9610c1e565b856000808281526003602081905260409091205460ff169081111561030a57fe5b1461035f576040805160e560020a62461bcd02815260206004820152601660248201527f73776170206f70656e65642070726576696f75736c7900000000000000000000604482015290519081900360640190fd5b600154604080517f23b872dd000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018690529051600160a060020a03909216916323b872dd9160648082019260009290919082900301818387803b1580156103d157600080fd5b505af11580156103e5573d6000803e3d6000fd5b5050505060c06040519081016040528085815260200184815260200186600019168152602001600060010260001916815260200133600160a060020a0316815260200187600160a060020a03168152509150816002600089600019166000191681526020019081526020016000206000820151816000015560208201518160010155604082015181600201906000191690556060820151816003019060001916905560808201518160040160006101000a815481600160a060020a030219169083600160a060020a0316021790555060a08201518160050160006101000a815481600160a060020a030219169083600160a060020a03160217905550905050600160036000896000191660001916815260200190815260200160002060006101000a81548160ff0219169083600381111561051c57fe5b021790555060408051888152600160a060020a038816602082015280820187905290517f497d46e9505eefe8b910d1a02e6b40d8769510023b0053c3ac4b9574b81c97bf9181900360600190a150505050505050565b6000805b60008381526003602081905260409091205460ff169081111561059557fe5b1492915050565b600154600160a060020a031681565b6040805160208082018590528183018490528251808303840181526060909201928390528151600093918291908401908083835b602083106105fe5780518252601f1990920191602091820191016105df565b5181516020939093036101000a600019018019909116921691909117905260405192018290039091209695505050505050565b60006001610576565b80600160008281526003602081905260409091205460ff169081111561065c57fe5b146106b1576040805160e560020a62461bcd02815260206004820152600d60248201527f73776170206e6f74206f70656e00000000000000000000000000000000000000604482015290519081900360640190fd5b6000828152600260205260409020548290421015610719576040805160e560020a62461bcd02815260206004820152601260248201527f73776170206e6f7420657870697261626c650000000000000000000000000000604482015290519081900360640190fd5b6000838152600360208181526040808420805460ff1916909317909255600180546002909252828420600480820154919092015484517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a039283169381019390935260248301529251929091169263a9059cbb9260448084019382900301818387803b1580156107b157600080fd5b505af11580156107c5573d6000803e3d6000fd5b50506040805186815290517feb711459e1247171382f0da0387b86239b8e3ca60af3b15a9ff2f1eb3d7f6a1d9350908190036020019150a1505050565b600081600260008281526003602081905260409091205460ff169081111561082657fe5b1461087b576040805160e560020a62461bcd02815260206004820152601160248201527f73776170206e6f742072656465656d6564000000000000000000000000000000604482015290519081900360640190fd5b505060009081526002602052604090206003015490565b60008181526002602052604081205442108015906108b257506001610576565b92915050565b81600160008281526003602081905260409091205460ff16908111156108da57fe5b1461092f576040805160e560020a62461bcd02815260206004820152600d60248201527f73776170206e6f74206f70656e00000000000000000000000000000000000000604482015290519081900360640190fd5b60408051602080820185905282518083038201815291830192839052815186938693600293909282918401908083835b6020831061097e5780518252601f19909201916020918201910161095f565b51815160209384036101000a600019018019909216911617905260405191909301945091925050808303816000865af11580156109bf573d6000803e3d6000fd5b5050506040513d60208110156109d457600080fd5b50516000838152600260208190526040909120015414610a3e576040805160e560020a62461bcd02815260206004820152600e60248201527f696e76616c696420736563726574000000000000000000000000000000000000604482015290519081900360640190fd5b600085815260026020818152604080842060038082018a90558352818520805460ff191685179055600480845282862042905560018054959094526005820154939091015482517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a039485169281019290925260248201529051919092169263a9059cbb926044808201939182900301818387803b158015610ae857600080fd5b505af1158015610afc573d6000803e3d6000fd5b5050604080518881526020810188905281517f07da1fa25a1d885732677ce9c192cbec27051a4b69d391c9a64850f5a5112ba09450908190039091019150a15050505050565b60046020526000908152604090205481565b60009081526002602081905260409091208054600182015460058301546004840154939094015491949093600160a060020a0390811693169190565b6000805460408051602060026001851615610100026000190190941693909304601f81018490048402820184019092528181529291830182828015610c165780601f10610beb57610100808354040283529160200191610c16565b820191906000526020600020905b815481529060010190602001808311610bf957829003601f168201915b505050505081565b6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a0810191909152905600a165627a7a723058202bfb78f46175fe2f5ec08f0403f2b471bc4eb55873556f14fe52dacdfad834040029`

// DeployRenExAtomicSwapper deploys a new Ethereum contract, binding an instance of RenExAtomicSwapper to it.
func DeployRenExAtomicSwapper(auth *bind.TransactOpts, backend bind.ContractBackend, _VERSION string, _TOKEN_ADDRESS common.Address) (common.Address, *types.Transaction, *RenExAtomicSwapper, error) {
	parsed, err := abi.JSON(strings.NewReader(RenExAtomicSwapperABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RenExAtomicSwapperBin), backend, _VERSION, _TOKEN_ADDRESS)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RenExAtomicSwapper{RenExAtomicSwapperCaller: RenExAtomicSwapperCaller{contract: contract}, RenExAtomicSwapperTransactor: RenExAtomicSwapperTransactor{contract: contract}, RenExAtomicSwapperFilterer: RenExAtomicSwapperFilterer{contract: contract}}, nil
}

// RenExAtomicSwapper is an auto generated Go binding around an Ethereum contract.
type RenExAtomicSwapper struct {
	RenExAtomicSwapperCaller     // Read-only binding to the contract
	RenExAtomicSwapperTransactor // Write-only binding to the contract
	RenExAtomicSwapperFilterer   // Log filterer for contract events
}

// RenExAtomicSwapperCaller is an auto generated read-only Go binding around an Ethereum contract.
type RenExAtomicSwapperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RenExAtomicSwapperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RenExAtomicSwapperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RenExAtomicSwapperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RenExAtomicSwapperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RenExAtomicSwapperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RenExAtomicSwapperSession struct {
	Contract     *RenExAtomicSwapper // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// RenExAtomicSwapperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RenExAtomicSwapperCallerSession struct {
	Contract *RenExAtomicSwapperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// RenExAtomicSwapperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RenExAtomicSwapperTransactorSession struct {
	Contract     *RenExAtomicSwapperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// RenExAtomicSwapperRaw is an auto generated low-level Go binding around an Ethereum contract.
type RenExAtomicSwapperRaw struct {
	Contract *RenExAtomicSwapper // Generic contract binding to access the raw methods on
}

// RenExAtomicSwapperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RenExAtomicSwapperCallerRaw struct {
	Contract *RenExAtomicSwapperCaller // Generic read-only contract binding to access the raw methods on
}

// RenExAtomicSwapperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RenExAtomicSwapperTransactorRaw struct {
	Contract *RenExAtomicSwapperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRenExAtomicSwapper creates a new instance of RenExAtomicSwapper, bound to a specific deployed contract.
func NewRenExAtomicSwapper(address common.Address, backend bind.ContractBackend) (*RenExAtomicSwapper, error) {
	contract, err := bindRenExAtomicSwapper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RenExAtomicSwapper{RenExAtomicSwapperCaller: RenExAtomicSwapperCaller{contract: contract}, RenExAtomicSwapperTransactor: RenExAtomicSwapperTransactor{contract: contract}, RenExAtomicSwapperFilterer: RenExAtomicSwapperFilterer{contract: contract}}, nil
}

// NewRenExAtomicSwapperCaller creates a new read-only instance of RenExAtomicSwapper, bound to a specific deployed contract.
func NewRenExAtomicSwapperCaller(address common.Address, caller bind.ContractCaller) (*RenExAtomicSwapperCaller, error) {
	contract, err := bindRenExAtomicSwapper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RenExAtomicSwapperCaller{contract: contract}, nil
}

// NewRenExAtomicSwapperTransactor creates a new write-only instance of RenExAtomicSwapper, bound to a specific deployed contract.
func NewRenExAtomicSwapperTransactor(address common.Address, transactor bind.ContractTransactor) (*RenExAtomicSwapperTransactor, error) {
	contract, err := bindRenExAtomicSwapper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RenExAtomicSwapperTransactor{contract: contract}, nil
}

// NewRenExAtomicSwapperFilterer creates a new log filterer instance of RenExAtomicSwapper, bound to a specific deployed contract.
func NewRenExAtomicSwapperFilterer(address common.Address, filterer bind.ContractFilterer) (*RenExAtomicSwapperFilterer, error) {
	contract, err := bindRenExAtomicSwapper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RenExAtomicSwapperFilterer{contract: contract}, nil
}

// bindRenExAtomicSwapper binds a generic wrapper to an already deployed contract.
func bindRenExAtomicSwapper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RenExAtomicSwapperABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RenExAtomicSwapper *RenExAtomicSwapperRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RenExAtomicSwapper.Contract.RenExAtomicSwapperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RenExAtomicSwapper *RenExAtomicSwapperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RenExAtomicSwapper.Contract.RenExAtomicSwapperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RenExAtomicSwapper *RenExAtomicSwapperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RenExAtomicSwapper.Contract.RenExAtomicSwapperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RenExAtomicSwapper *RenExAtomicSwapperCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RenExAtomicSwapper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RenExAtomicSwapper *RenExAtomicSwapperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RenExAtomicSwapper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RenExAtomicSwapper *RenExAtomicSwapperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RenExAtomicSwapper.Contract.contract.Transact(opts, method, params...)
}

// TOKENADDRESS is a free data retrieval call binding the contract method 0x0bdf5300.
//
// Solidity: function TOKEN_ADDRESS() constant returns(address)
func (_RenExAtomicSwapper *RenExAtomicSwapperCaller) TOKENADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _RenExAtomicSwapper.contract.Call(opts, out, "TOKEN_ADDRESS")
	return *ret0, err
}

// TOKENADDRESS is a free data retrieval call binding the contract method 0x0bdf5300.
//
// Solidity: function TOKEN_ADDRESS() constant returns(address)
func (_RenExAtomicSwapper *RenExAtomicSwapperSession) TOKENADDRESS() (common.Address, error) {
	return _RenExAtomicSwapper.Contract.TOKENADDRESS(&_RenExAtomicSwapper.CallOpts)
}

// TOKENADDRESS is a free data retrieval call binding the contract method 0x0bdf5300.
//
// Solidity: function TOKEN_ADDRESS() constant returns(address)
func (_RenExAtomicSwapper *RenExAtomicSwapperCallerSession) TOKENADDRESS() (common.Address, error) {
	return _RenExAtomicSwapper.Contract.TOKENADDRESS(&_RenExAtomicSwapper.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_RenExAtomicSwapper *RenExAtomicSwapperCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _RenExAtomicSwapper.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_RenExAtomicSwapper *RenExAtomicSwapperSession) VERSION() (string, error) {
	return _RenExAtomicSwapper.Contract.VERSION(&_RenExAtomicSwapper.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_RenExAtomicSwapper *RenExAtomicSwapperCallerSession) VERSION() (string, error) {
	return _RenExAtomicSwapper.Contract.VERSION(&_RenExAtomicSwapper.CallOpts)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, from address, secretLock bytes32)
func (_RenExAtomicSwapper *RenExAtomicSwapperCaller) Audit(opts *bind.CallOpts, _swapID [32]byte) (struct {
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
	err := _RenExAtomicSwapper.contract.Call(opts, out, "audit", _swapID)
	return *ret, err
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, from address, secretLock bytes32)
func (_RenExAtomicSwapper *RenExAtomicSwapperSession) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	return _RenExAtomicSwapper.Contract.Audit(&_RenExAtomicSwapper.CallOpts, _swapID)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, from address, secretLock bytes32)
func (_RenExAtomicSwapper *RenExAtomicSwapperCallerSession) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	return _RenExAtomicSwapper.Contract.Audit(&_RenExAtomicSwapper.CallOpts, _swapID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_RenExAtomicSwapper *RenExAtomicSwapperCaller) AuditSecret(opts *bind.CallOpts, _swapID [32]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _RenExAtomicSwapper.contract.Call(opts, out, "auditSecret", _swapID)
	return *ret0, err
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_RenExAtomicSwapper *RenExAtomicSwapperSession) AuditSecret(_swapID [32]byte) ([32]byte, error) {
	return _RenExAtomicSwapper.Contract.AuditSecret(&_RenExAtomicSwapper.CallOpts, _swapID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_RenExAtomicSwapper *RenExAtomicSwapperCallerSession) AuditSecret(_swapID [32]byte) ([32]byte, error) {
	return _RenExAtomicSwapper.Contract.AuditSecret(&_RenExAtomicSwapper.CallOpts, _swapID)
}

// Initiatable is a free data retrieval call binding the contract method 0x09ece618.
//
// Solidity: function initiatable(_swapID bytes32) constant returns(bool)
func (_RenExAtomicSwapper *RenExAtomicSwapperCaller) Initiatable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _RenExAtomicSwapper.contract.Call(opts, out, "initiatable", _swapID)
	return *ret0, err
}

// Initiatable is a free data retrieval call binding the contract method 0x09ece618.
//
// Solidity: function initiatable(_swapID bytes32) constant returns(bool)
func (_RenExAtomicSwapper *RenExAtomicSwapperSession) Initiatable(_swapID [32]byte) (bool, error) {
	return _RenExAtomicSwapper.Contract.Initiatable(&_RenExAtomicSwapper.CallOpts, _swapID)
}

// Initiatable is a free data retrieval call binding the contract method 0x09ece618.
//
// Solidity: function initiatable(_swapID bytes32) constant returns(bool)
func (_RenExAtomicSwapper *RenExAtomicSwapperCallerSession) Initiatable(_swapID [32]byte) (bool, error) {
	return _RenExAtomicSwapper.Contract.Initiatable(&_RenExAtomicSwapper.CallOpts, _swapID)
}

// Redeemable is a free data retrieval call binding the contract method 0x68f06b29.
//
// Solidity: function redeemable(_swapID bytes32) constant returns(bool)
func (_RenExAtomicSwapper *RenExAtomicSwapperCaller) Redeemable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _RenExAtomicSwapper.contract.Call(opts, out, "redeemable", _swapID)
	return *ret0, err
}

// Redeemable is a free data retrieval call binding the contract method 0x68f06b29.
//
// Solidity: function redeemable(_swapID bytes32) constant returns(bool)
func (_RenExAtomicSwapper *RenExAtomicSwapperSession) Redeemable(_swapID [32]byte) (bool, error) {
	return _RenExAtomicSwapper.Contract.Redeemable(&_RenExAtomicSwapper.CallOpts, _swapID)
}

// Redeemable is a free data retrieval call binding the contract method 0x68f06b29.
//
// Solidity: function redeemable(_swapID bytes32) constant returns(bool)
func (_RenExAtomicSwapper *RenExAtomicSwapperCallerSession) Redeemable(_swapID [32]byte) (bool, error) {
	return _RenExAtomicSwapper.Contract.Redeemable(&_RenExAtomicSwapper.CallOpts, _swapID)
}

// RedeemedAt is a free data retrieval call binding the contract method 0xbc4fcc4a.
//
// Solidity: function redeemedAt( bytes32) constant returns(uint256)
func (_RenExAtomicSwapper *RenExAtomicSwapperCaller) RedeemedAt(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RenExAtomicSwapper.contract.Call(opts, out, "redeemedAt", arg0)
	return *ret0, err
}

// RedeemedAt is a free data retrieval call binding the contract method 0xbc4fcc4a.
//
// Solidity: function redeemedAt( bytes32) constant returns(uint256)
func (_RenExAtomicSwapper *RenExAtomicSwapperSession) RedeemedAt(arg0 [32]byte) (*big.Int, error) {
	return _RenExAtomicSwapper.Contract.RedeemedAt(&_RenExAtomicSwapper.CallOpts, arg0)
}

// RedeemedAt is a free data retrieval call binding the contract method 0xbc4fcc4a.
//
// Solidity: function redeemedAt( bytes32) constant returns(uint256)
func (_RenExAtomicSwapper *RenExAtomicSwapperCallerSession) RedeemedAt(arg0 [32]byte) (*big.Int, error) {
	return _RenExAtomicSwapper.Contract.RedeemedAt(&_RenExAtomicSwapper.CallOpts, arg0)
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(_swapID bytes32) constant returns(bool)
func (_RenExAtomicSwapper *RenExAtomicSwapperCaller) Refundable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _RenExAtomicSwapper.contract.Call(opts, out, "refundable", _swapID)
	return *ret0, err
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(_swapID bytes32) constant returns(bool)
func (_RenExAtomicSwapper *RenExAtomicSwapperSession) Refundable(_swapID [32]byte) (bool, error) {
	return _RenExAtomicSwapper.Contract.Refundable(&_RenExAtomicSwapper.CallOpts, _swapID)
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(_swapID bytes32) constant returns(bool)
func (_RenExAtomicSwapper *RenExAtomicSwapperCallerSession) Refundable(_swapID [32]byte) (bool, error) {
	return _RenExAtomicSwapper.Contract.Refundable(&_RenExAtomicSwapper.CallOpts, _swapID)
}

// SwapID is a free data retrieval call binding the contract method 0x4b2ac3fa.
//
// Solidity: function swapID(_secretLock bytes32, _timelock uint256) constant returns(bytes32)
func (_RenExAtomicSwapper *RenExAtomicSwapperCaller) SwapID(opts *bind.CallOpts, _secretLock [32]byte, _timelock *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _RenExAtomicSwapper.contract.Call(opts, out, "swapID", _secretLock, _timelock)
	return *ret0, err
}

// SwapID is a free data retrieval call binding the contract method 0x4b2ac3fa.
//
// Solidity: function swapID(_secretLock bytes32, _timelock uint256) constant returns(bytes32)
func (_RenExAtomicSwapper *RenExAtomicSwapperSession) SwapID(_secretLock [32]byte, _timelock *big.Int) ([32]byte, error) {
	return _RenExAtomicSwapper.Contract.SwapID(&_RenExAtomicSwapper.CallOpts, _secretLock, _timelock)
}

// SwapID is a free data retrieval call binding the contract method 0x4b2ac3fa.
//
// Solidity: function swapID(_secretLock bytes32, _timelock uint256) constant returns(bytes32)
func (_RenExAtomicSwapper *RenExAtomicSwapperCallerSession) SwapID(_secretLock [32]byte, _timelock *big.Int) ([32]byte, error) {
	return _RenExAtomicSwapper.Contract.SwapID(&_RenExAtomicSwapper.CallOpts, _secretLock, _timelock)
}

// Initiate is a paid mutator transaction binding the contract method 0x027a2577.
//
// Solidity: function initiate(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_RenExAtomicSwapper *RenExAtomicSwapperTransactor) Initiate(opts *bind.TransactOpts, _swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _RenExAtomicSwapper.contract.Transact(opts, "initiate", _swapID, _withdrawTrader, _secretLock, _timelock, _value)
}

// Initiate is a paid mutator transaction binding the contract method 0x027a2577.
//
// Solidity: function initiate(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_RenExAtomicSwapper *RenExAtomicSwapperSession) Initiate(_swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _RenExAtomicSwapper.Contract.Initiate(&_RenExAtomicSwapper.TransactOpts, _swapID, _withdrawTrader, _secretLock, _timelock, _value)
}

// Initiate is a paid mutator transaction binding the contract method 0x027a2577.
//
// Solidity: function initiate(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_RenExAtomicSwapper *RenExAtomicSwapperTransactorSession) Initiate(_swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _RenExAtomicSwapper.Contract.Initiate(&_RenExAtomicSwapper.TransactOpts, _swapID, _withdrawTrader, _secretLock, _timelock, _value)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_RenExAtomicSwapper *RenExAtomicSwapperTransactor) Redeem(opts *bind.TransactOpts, _swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _RenExAtomicSwapper.contract.Transact(opts, "redeem", _swapID, _secretKey)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_RenExAtomicSwapper *RenExAtomicSwapperSession) Redeem(_swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _RenExAtomicSwapper.Contract.Redeem(&_RenExAtomicSwapper.TransactOpts, _swapID, _secretKey)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_RenExAtomicSwapper *RenExAtomicSwapperTransactorSession) Redeem(_swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _RenExAtomicSwapper.Contract.Redeem(&_RenExAtomicSwapper.TransactOpts, _swapID, _secretKey)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_RenExAtomicSwapper *RenExAtomicSwapperTransactor) Refund(opts *bind.TransactOpts, _swapID [32]byte) (*types.Transaction, error) {
	return _RenExAtomicSwapper.contract.Transact(opts, "refund", _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_RenExAtomicSwapper *RenExAtomicSwapperSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _RenExAtomicSwapper.Contract.Refund(&_RenExAtomicSwapper.TransactOpts, _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_RenExAtomicSwapper *RenExAtomicSwapperTransactorSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _RenExAtomicSwapper.Contract.Refund(&_RenExAtomicSwapper.TransactOpts, _swapID)
}

// RenExAtomicSwapperLogCloseIterator is returned from FilterLogClose and is used to iterate over the raw logs and unpacked data for LogClose events raised by the RenExAtomicSwapper contract.
type RenExAtomicSwapperLogCloseIterator struct {
	Event *RenExAtomicSwapperLogClose // Event containing the contract specifics and raw log

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
func (it *RenExAtomicSwapperLogCloseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RenExAtomicSwapperLogClose)
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
		it.Event = new(RenExAtomicSwapperLogClose)
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
func (it *RenExAtomicSwapperLogCloseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RenExAtomicSwapperLogCloseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RenExAtomicSwapperLogClose represents a LogClose event raised by the RenExAtomicSwapper contract.
type RenExAtomicSwapperLogClose struct {
	SwapID    [32]byte
	SecretKey [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLogClose is a free log retrieval operation binding the contract event 0x07da1fa25a1d885732677ce9c192cbec27051a4b69d391c9a64850f5a5112ba0.
//
// Solidity: e LogClose(_swapID bytes32, _secretKey bytes32)
func (_RenExAtomicSwapper *RenExAtomicSwapperFilterer) FilterLogClose(opts *bind.FilterOpts) (*RenExAtomicSwapperLogCloseIterator, error) {

	logs, sub, err := _RenExAtomicSwapper.contract.FilterLogs(opts, "LogClose")
	if err != nil {
		return nil, err
	}
	return &RenExAtomicSwapperLogCloseIterator{contract: _RenExAtomicSwapper.contract, event: "LogClose", logs: logs, sub: sub}, nil
}

// WatchLogClose is a free log subscription operation binding the contract event 0x07da1fa25a1d885732677ce9c192cbec27051a4b69d391c9a64850f5a5112ba0.
//
// Solidity: e LogClose(_swapID bytes32, _secretKey bytes32)
func (_RenExAtomicSwapper *RenExAtomicSwapperFilterer) WatchLogClose(opts *bind.WatchOpts, sink chan<- *RenExAtomicSwapperLogClose) (event.Subscription, error) {

	logs, sub, err := _RenExAtomicSwapper.contract.WatchLogs(opts, "LogClose")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RenExAtomicSwapperLogClose)
				if err := _RenExAtomicSwapper.contract.UnpackLog(event, "LogClose", log); err != nil {
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

// RenExAtomicSwapperLogExpireIterator is returned from FilterLogExpire and is used to iterate over the raw logs and unpacked data for LogExpire events raised by the RenExAtomicSwapper contract.
type RenExAtomicSwapperLogExpireIterator struct {
	Event *RenExAtomicSwapperLogExpire // Event containing the contract specifics and raw log

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
func (it *RenExAtomicSwapperLogExpireIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RenExAtomicSwapperLogExpire)
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
		it.Event = new(RenExAtomicSwapperLogExpire)
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
func (it *RenExAtomicSwapperLogExpireIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RenExAtomicSwapperLogExpireIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RenExAtomicSwapperLogExpire represents a LogExpire event raised by the RenExAtomicSwapper contract.
type RenExAtomicSwapperLogExpire struct {
	SwapID [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogExpire is a free log retrieval operation binding the contract event 0xeb711459e1247171382f0da0387b86239b8e3ca60af3b15a9ff2f1eb3d7f6a1d.
//
// Solidity: e LogExpire(_swapID bytes32)
func (_RenExAtomicSwapper *RenExAtomicSwapperFilterer) FilterLogExpire(opts *bind.FilterOpts) (*RenExAtomicSwapperLogExpireIterator, error) {

	logs, sub, err := _RenExAtomicSwapper.contract.FilterLogs(opts, "LogExpire")
	if err != nil {
		return nil, err
	}
	return &RenExAtomicSwapperLogExpireIterator{contract: _RenExAtomicSwapper.contract, event: "LogExpire", logs: logs, sub: sub}, nil
}

// WatchLogExpire is a free log subscription operation binding the contract event 0xeb711459e1247171382f0da0387b86239b8e3ca60af3b15a9ff2f1eb3d7f6a1d.
//
// Solidity: e LogExpire(_swapID bytes32)
func (_RenExAtomicSwapper *RenExAtomicSwapperFilterer) WatchLogExpire(opts *bind.WatchOpts, sink chan<- *RenExAtomicSwapperLogExpire) (event.Subscription, error) {

	logs, sub, err := _RenExAtomicSwapper.contract.WatchLogs(opts, "LogExpire")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RenExAtomicSwapperLogExpire)
				if err := _RenExAtomicSwapper.contract.UnpackLog(event, "LogExpire", log); err != nil {
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

// RenExAtomicSwapperLogOpenIterator is returned from FilterLogOpen and is used to iterate over the raw logs and unpacked data for LogOpen events raised by the RenExAtomicSwapper contract.
type RenExAtomicSwapperLogOpenIterator struct {
	Event *RenExAtomicSwapperLogOpen // Event containing the contract specifics and raw log

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
func (it *RenExAtomicSwapperLogOpenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RenExAtomicSwapperLogOpen)
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
		it.Event = new(RenExAtomicSwapperLogOpen)
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
func (it *RenExAtomicSwapperLogOpenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RenExAtomicSwapperLogOpenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RenExAtomicSwapperLogOpen represents a LogOpen event raised by the RenExAtomicSwapper contract.
type RenExAtomicSwapperLogOpen struct {
	SwapID         [32]byte
	WithdrawTrader common.Address
	SecretLock     [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterLogOpen is a free log retrieval operation binding the contract event 0x497d46e9505eefe8b910d1a02e6b40d8769510023b0053c3ac4b9574b81c97bf.
//
// Solidity: e LogOpen(_swapID bytes32, _withdrawTrader address, _secretLock bytes32)
func (_RenExAtomicSwapper *RenExAtomicSwapperFilterer) FilterLogOpen(opts *bind.FilterOpts) (*RenExAtomicSwapperLogOpenIterator, error) {

	logs, sub, err := _RenExAtomicSwapper.contract.FilterLogs(opts, "LogOpen")
	if err != nil {
		return nil, err
	}
	return &RenExAtomicSwapperLogOpenIterator{contract: _RenExAtomicSwapper.contract, event: "LogOpen", logs: logs, sub: sub}, nil
}

// WatchLogOpen is a free log subscription operation binding the contract event 0x497d46e9505eefe8b910d1a02e6b40d8769510023b0053c3ac4b9574b81c97bf.
//
// Solidity: e LogOpen(_swapID bytes32, _withdrawTrader address, _secretLock bytes32)
func (_RenExAtomicSwapper *RenExAtomicSwapperFilterer) WatchLogOpen(opts *bind.WatchOpts, sink chan<- *RenExAtomicSwapperLogOpen) (event.Subscription, error) {

	logs, sub, err := _RenExAtomicSwapper.contract.WatchLogs(opts, "LogOpen")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RenExAtomicSwapperLogOpen)
				if err := _RenExAtomicSwapper.contract.UnpackLog(event, "LogOpen", log); err != nil {
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

// SafeMathABI is the input ABI used to generate the binding from.
const SafeMathABI = "[]"

// SafeMathBin is the compiled bytecode used for deploying new contracts.
const SafeMathBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a7230582088033d534f66266dd3cfc8a8f61e92367283c30e2bb591989a19561d224255280029`

// DeploySafeMath deploys a new Ethereum contract, binding an instance of SafeMath to it.
func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// SafeMath is an auto generated Go binding around an Ethereum contract.
type SafeMath struct {
	SafeMathCaller     // Read-only binding to the contract
	SafeMathTransactor // Write-only binding to the contract
	SafeMathFilterer   // Log filterer for contract events
}

// SafeMathCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMathSession struct {
	Contract     *SafeMath         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMathCallerSession struct {
	Contract *SafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMathTransactorSession struct {
	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeMathRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMathRaw struct {
	Contract *SafeMath // Generic contract binding to access the raw methods on
}

// SafeMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMathCallerRaw struct {
	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMathTransactorRaw struct {
	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
	contract, err := bindSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
	contract, err := bindSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathCaller{contract: contract}, nil
}

// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
	contract, err := bindSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathTransactor{contract: contract}, nil
}

// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
	contract, err := bindSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathFilterer{contract: contract}, nil
}

// bindSafeMath binds a generic wrapper to an already deployed contract.
func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transact(opts, method, params...)
}
