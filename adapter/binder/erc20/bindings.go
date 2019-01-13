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

// SwapperdERC20ABI is the input ABI used to generate the binding from.
const SwapperdERC20ABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"initiate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"initiatable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"TOKEN_ADDRESS\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"}],\"name\":\"swapID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawBrokerFees\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"redeemable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"auditSecret\",\"outputs\":[{\"name\":\"secretKey\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refundable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_secretKey\",\"type\":\"bytes32\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_broker\",\"type\":\"address\"},{\"name\":\"_brokerFee\",\"type\":\"uint256\"},{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"initiateWithFees\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"redeemedAt\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"audit\",\"outputs\":[{\"name\":\"timelock\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"brokerFee\",\"type\":\"uint256\"},{\"name\":\"broker\",\"type\":\"address\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"secretLock\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"brokerFees\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_VERSION\",\"type\":\"string\"},{\"name\":\"_TOKEN_ADDRESS\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_secretLock\",\"type\":\"bytes32\"}],\"name\":\"LogOpen\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"LogExpire\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_secretKey\",\"type\":\"bytes32\"}],\"name\":\"LogClose\",\"type\":\"event\"}]"

// SwapperdERC20Bin is the compiled bytecode used for deploying new contracts.
const SwapperdERC20Bin = `0x60806040523480156200001157600080fd5b50604051620013d3380380620013d3833981018060405260408110156200003757600080fd5b8101908080516401000000008111156200005057600080fd5b820160208101848111156200006457600080fd5b81516401000000008111828201871017156200007f57600080fd5b505060209182015181519194509250620000a09160009190850190620000c8565b5060018054600160a060020a031916600160a060020a0392909216919091179055506200016d565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200010b57805160ff19168380011785556200013b565b828001600101855582156200013b579182015b828111156200013b5782518255916020019190600101906200011e565b50620001499291506200014d565b5090565b6200016a91905b8082111562000149576000815560010162000154565b90565b611256806200017d6000396000f3fe6080604052600436106100d4577c01000000000000000000000000000000000000000000000000000000006000350463027a257781146100d957806309ece618146101265780630bdf5300146101645780634b2ac3fa146101955780634c6d37ff146101d757806368f06b29146102015780637249fbb61461022b578063976d00f4146102555780639fb314751461027f578063b31597ad146102a9578063b8688e3f146102d9578063bc4fcc4a14610334578063c140635b1461035e578063e1ec380c146103ce578063ffa1ad7414610401575b600080fd5b3480156100e557600080fd5b50610124600480360360a08110156100fc57600080fd5b50803590600160a060020a03602082013516906040810135906060810135906080013561048b565b005b34801561013257600080fd5b506101506004803603602081101561014957600080fd5b5035610747565b604080519115158252519081900360200190f35b34801561017057600080fd5b50610179610771565b60408051600160a060020a039092168252519081900360200190f35b3480156101a157600080fd5b506101c5600480360360408110156101b857600080fd5b5080359060200135610780565b60408051918252519081900360200190f35b3480156101e357600080fd5b50610124600480360360208110156101fa57600080fd5b50356107ac565b34801561020d57600080fd5b506101506004803603602081101561022457600080fd5b5035610863565b34801561023757600080fd5b506101246004803603602081101561024e57600080fd5b503561086c565b34801561026157600080fd5b506101c56004803603602081101561027857600080fd5b5035610a37565b34801561028b57600080fd5b50610150600480360360208110156102a257600080fd5b5035610ac7565b3480156102b557600080fd5b50610124600480360360408110156102cc57600080fd5b5080359060200135610aed565b3480156102e557600080fd5b50610124600480360360e08110156102fc57600080fd5b50803590600160a060020a03602082013581169160408101359091169060608101359060808101359060a08101359060c00135610dac565b34801561034057600080fd5b506101c56004803603602081101561035757600080fd5b503561108c565b34801561036a57600080fd5b506103886004803603602081101561038157600080fd5b503561109e565b604080519788526020880196909652600160a060020a03948516878701526060870193909352908316608086015290911660a084015260c0830152519081900360e00190f35b3480156103da57600080fd5b506101c5600480360360208110156103f157600080fd5b5035600160a060020a0316611146565b34801561040d57600080fd5b50610416611158565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610450578181015183820152602001610438565b50505050905090810190601f16801561047d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b846000808281526003602081905260409091205460ff16908111156104ac57fe5b14610501576040805160e560020a62461bcd02815260206004820152601660248201527f73776170206f70656e65642070726576696f75736c7900000000000000000000604482015290519081900360640190fd5b600154604080517f23b872dd000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018590529051600160a060020a03909216916323b872dd9160648082019260009290919082900301818387803b15801561057357600080fd5b505af1158015610587573d6000803e3d6000fd5b505050506105936111e6565b61010060405190810160405280858152602001848152602001600081526020018681526020016000600102815260200133600160a060020a0316815260200187600160a060020a031681526020016000600160a060020a031681525090508060026000898152602001908152602001600020600082015181600001556020820151816001015560408201518160020155606082015181600301556080820151816004015560a08201518160050160006101000a815481600160a060020a030219169083600160a060020a0316021790555060c08201518160060160006101000a815481600160a060020a030219169083600160a060020a0316021790555060e08201518160070160006101000a815481600160a060020a030219169083600160a060020a0316021790555090505060016003600089815260200190815260200160002060006101000a81548160ff021916908360038111156106f157fe5b021790555060408051888152600160a060020a038816602082015280820187905290517f497d46e9505eefe8b910d1a02e6b40d8769510023b0053c3ac4b9574b81c97bf9181900360600190a150505050505050565b6000805b60008381526003602081905260409091205460ff169081111561076a57fe5b1492915050565b600154600160a060020a031681565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b336000908152600460205260409020548111156107c857600080fd5b3360008181526004602081905260408083208054869003905560015481517fa9059cbb000000000000000000000000000000000000000000000000000000008152928301949094526024820185905251600160a060020a039093169263a9059cbb9260448084019391929182900301818387803b15801561084857600080fd5b505af115801561085c573d6000803e3d6000fd5b5050505050565b6000600161074b565b80600160008281526003602081905260409091205460ff169081111561088e57fe5b146108e3576040805160e560020a62461bcd02815260206004820152600d60248201527f73776170206e6f74206f70656e00000000000000000000000000000000000000604482015290519081900360640190fd5b600082815260026020526040902054829042101561094b576040805160e560020a62461bcd02815260206004820152601260248201527f73776170206e6f7420657870697261626c650000000000000000000000000000604482015290519081900360640190fd5b6000838152600360208181526040808420805460ff19169093179092556001805460029283905283852060058101549381015492015484517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a039485166004820152920160248301529251919092169263a9059cbb926044808201939182900301818387803b1580156109e657600080fd5b505af11580156109fa573d6000803e3d6000fd5b50506040805186815290517feb711459e1247171382f0da0387b86239b8e3ca60af3b15a9ff2f1eb3d7f6a1d9350908190036020019150a1505050565b600081600260008281526003602081905260409091205460ff1690811115610a5b57fe5b14610ab0576040805160e560020a62461bcd02815260206004820152601160248201527f73776170206e6f742072656465656d6564000000000000000000000000000000604482015290519081900360640190fd5b505060009081526002602052604090206004015490565b6000818152600260205260408120544210801590610ae75750600161074b565b92915050565b81600160008281526003602081905260409091205460ff1690811115610b0f57fe5b14610b64576040805160e560020a62461bcd02815260206004820152600d60248201527f73776170206e6f74206f70656e00000000000000000000000000000000000000604482015290519081900360640190fd5b8282600281604051602001808281526020019150506040516020818303038152906040526040518082805190602001908083835b60208310610bb75780518252601f199092019160209182019101610b98565b51815160209384036101000a60001901801990921691161790526040519190930194509192505080830381855afa158015610bf6573d6000803e3d6000fd5b5050506040513d6020811015610c0b57600080fd5b505160008381526002602052604090206003015414610c74576040805160e560020a62461bcd02815260206004820152600e60248201527f696e76616c696420736563726574000000000000000000000000000000000000604482015290519081900360640190fd5b600085815260026020818152604080842060048082018a905560038452828620805460ff1916861790556005845282862042905560018054959094526006820154939091015482517fa9059cbb000000000000000000000000000000000000000000000000000000008152600160a060020a039485169281019290925260248201529051919092169263a9059cbb926044808201939182900301818387803b158015610d1f57600080fd5b505af1158015610d33573d6000803e3d6000fd5b505050600086815260026020818152604080842092830154600790930154600160a060020a0316845260048252928390208054909201909155815188815290810187905281517f07da1fa25a1d885732677ce9c192cbec27051a4b69d391c9a64850f5a5112ba093509081900390910190a15050505050565b866000808281526003602081905260409091205460ff1690811115610dcd57fe5b14610e22576040805160e560020a62461bcd02815260206004820152601660248201527f73776170206f70656e65642070726576696f75736c7900000000000000000000604482015290519081900360640190fd5b600154604080517f23b872dd000000000000000000000000000000000000000000000000000000008152336004820152306024820152604481018590529051600160a060020a03909216916323b872dd9160648082019260009290919082900301818387803b158015610e9457600080fd5b505af1158015610ea8573d6000803e3d6000fd5b50505050600160a060020a03861615801590610ec357508415155b1515610ece57600080fd5b610ed66111e6565b6101006040519081016040528085815260200187850381526020018781526020018681526020016000600102815260200133600160a060020a0316815260200189600160a060020a0316815260200188600160a060020a0316815250905080600260008b8152602001908152602001600020600082015181600001556020820151816001015560408201518160020155606082015181600301556080820151816004015560a08201518160050160006101000a815481600160a060020a030219169083600160a060020a0316021790555060c08201518160060160006101000a815481600160a060020a030219169083600160a060020a0316021790555060e08201518160070160006101000a815481600160a060020a030219169083600160a060020a031602179055509050506001600360008b815260200190815260200160002060006101000a81548160ff0219169083600381111561103457fe5b0217905550604080518a8152600160a060020a038a16602082015280820187905290517f497d46e9505eefe8b910d1a02e6b40d8769510023b0053c3ac4b9574b81c97bf9181900360600190a1505050505050505050565b60056020526000908152604090205481565b60008060008060008060006110b16111e6565b50505060009586525050600260208181526040958690208651610100810188528154808252600183015493820184905293820154978101889052600382015460608201819052600483015460808301526005830154600160a060020a0390811660a084018190526006850154821660c0850181905260079095015490911660e090930183905294999398929750919550935090565b60046020526000908152604090205481565b6000805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156111de5780601f106111b3576101008083540402835291602001916111de565b820191906000526020600020905b8154815290600101906020018083116111c157829003601f168201915b505050505081565b6040805161010081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e08101919091529056fea165627a7a723058204a522c0dc985c585b8352a7a9a01b35220a1e531bfec1c32cfdf7c8791e500af0029`

// DeploySwapperdERC20 deploys a new Ethereum contract, binding an instance of SwapperdERC20 to it.
func DeploySwapperdERC20(auth *bind.TransactOpts, backend bind.ContractBackend, _VERSION string, _TOKEN_ADDRESS common.Address) (common.Address, *types.Transaction, *SwapperdERC20, error) {
	parsed, err := abi.JSON(strings.NewReader(SwapperdERC20ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SwapperdERC20Bin), backend, _VERSION, _TOKEN_ADDRESS)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SwapperdERC20{SwapperdERC20Caller: SwapperdERC20Caller{contract: contract}, SwapperdERC20Transactor: SwapperdERC20Transactor{contract: contract}, SwapperdERC20Filterer: SwapperdERC20Filterer{contract: contract}}, nil
}

// SwapperdERC20 is an auto generated Go binding around an Ethereum contract.
type SwapperdERC20 struct {
	SwapperdERC20Caller     // Read-only binding to the contract
	SwapperdERC20Transactor // Write-only binding to the contract
	SwapperdERC20Filterer   // Log filterer for contract events
}

// SwapperdERC20Caller is an auto generated read-only Go binding around an Ethereum contract.
type SwapperdERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapperdERC20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type SwapperdERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapperdERC20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SwapperdERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwapperdERC20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SwapperdERC20Session struct {
	Contract     *SwapperdERC20    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SwapperdERC20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SwapperdERC20CallerSession struct {
	Contract *SwapperdERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// SwapperdERC20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SwapperdERC20TransactorSession struct {
	Contract     *SwapperdERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// SwapperdERC20Raw is an auto generated low-level Go binding around an Ethereum contract.
type SwapperdERC20Raw struct {
	Contract *SwapperdERC20 // Generic contract binding to access the raw methods on
}

// SwapperdERC20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SwapperdERC20CallerRaw struct {
	Contract *SwapperdERC20Caller // Generic read-only contract binding to access the raw methods on
}

// SwapperdERC20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SwapperdERC20TransactorRaw struct {
	Contract *SwapperdERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewSwapperdERC20 creates a new instance of SwapperdERC20, bound to a specific deployed contract.
func NewSwapperdERC20(address common.Address, backend bind.ContractBackend) (*SwapperdERC20, error) {
	contract, err := bindSwapperdERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SwapperdERC20{SwapperdERC20Caller: SwapperdERC20Caller{contract: contract}, SwapperdERC20Transactor: SwapperdERC20Transactor{contract: contract}, SwapperdERC20Filterer: SwapperdERC20Filterer{contract: contract}}, nil
}

// NewSwapperdERC20Caller creates a new read-only instance of SwapperdERC20, bound to a specific deployed contract.
func NewSwapperdERC20Caller(address common.Address, caller bind.ContractCaller) (*SwapperdERC20Caller, error) {
	contract, err := bindSwapperdERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SwapperdERC20Caller{contract: contract}, nil
}

// NewSwapperdERC20Transactor creates a new write-only instance of SwapperdERC20, bound to a specific deployed contract.
func NewSwapperdERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*SwapperdERC20Transactor, error) {
	contract, err := bindSwapperdERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SwapperdERC20Transactor{contract: contract}, nil
}

// NewSwapperdERC20Filterer creates a new log filterer instance of SwapperdERC20, bound to a specific deployed contract.
func NewSwapperdERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*SwapperdERC20Filterer, error) {
	contract, err := bindSwapperdERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SwapperdERC20Filterer{contract: contract}, nil
}

// bindSwapperdERC20 binds a generic wrapper to an already deployed contract.
func bindSwapperdERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SwapperdERC20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SwapperdERC20 *SwapperdERC20Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SwapperdERC20.Contract.SwapperdERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SwapperdERC20 *SwapperdERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwapperdERC20.Contract.SwapperdERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SwapperdERC20 *SwapperdERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SwapperdERC20.Contract.SwapperdERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SwapperdERC20 *SwapperdERC20CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SwapperdERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SwapperdERC20 *SwapperdERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwapperdERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SwapperdERC20 *SwapperdERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SwapperdERC20.Contract.contract.Transact(opts, method, params...)
}

// TOKENADDRESS is a free data retrieval call binding the contract method 0x0bdf5300.
//
// Solidity: function TOKEN_ADDRESS() constant returns(address)
func (_SwapperdERC20 *SwapperdERC20Caller) TOKENADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _SwapperdERC20.contract.Call(opts, out, "TOKEN_ADDRESS")
	return *ret0, err
}

// TOKENADDRESS is a free data retrieval call binding the contract method 0x0bdf5300.
//
// Solidity: function TOKEN_ADDRESS() constant returns(address)
func (_SwapperdERC20 *SwapperdERC20Session) TOKENADDRESS() (common.Address, error) {
	return _SwapperdERC20.Contract.TOKENADDRESS(&_SwapperdERC20.CallOpts)
}

// TOKENADDRESS is a free data retrieval call binding the contract method 0x0bdf5300.
//
// Solidity: function TOKEN_ADDRESS() constant returns(address)
func (_SwapperdERC20 *SwapperdERC20CallerSession) TOKENADDRESS() (common.Address, error) {
	return _SwapperdERC20.Contract.TOKENADDRESS(&_SwapperdERC20.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_SwapperdERC20 *SwapperdERC20Caller) VERSION(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _SwapperdERC20.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_SwapperdERC20 *SwapperdERC20Session) VERSION() (string, error) {
	return _SwapperdERC20.Contract.VERSION(&_SwapperdERC20.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_SwapperdERC20 *SwapperdERC20CallerSession) VERSION() (string, error) {
	return _SwapperdERC20.Contract.VERSION(&_SwapperdERC20.CallOpts)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, brokerFee uint256, broker address, from address, secretLock bytes32)
func (_SwapperdERC20 *SwapperdERC20Caller) Audit(opts *bind.CallOpts, _swapID [32]byte) (struct {
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
	err := _SwapperdERC20.contract.Call(opts, out, "audit", _swapID)
	return *ret, err
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, brokerFee uint256, broker address, from address, secretLock bytes32)
func (_SwapperdERC20 *SwapperdERC20Session) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	BrokerFee  *big.Int
	Broker     common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	return _SwapperdERC20.Contract.Audit(&_SwapperdERC20.CallOpts, _swapID)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, brokerFee uint256, broker address, from address, secretLock bytes32)
func (_SwapperdERC20 *SwapperdERC20CallerSession) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	BrokerFee  *big.Int
	Broker     common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	return _SwapperdERC20.Contract.Audit(&_SwapperdERC20.CallOpts, _swapID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_SwapperdERC20 *SwapperdERC20Caller) AuditSecret(opts *bind.CallOpts, _swapID [32]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _SwapperdERC20.contract.Call(opts, out, "auditSecret", _swapID)
	return *ret0, err
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_SwapperdERC20 *SwapperdERC20Session) AuditSecret(_swapID [32]byte) ([32]byte, error) {
	return _SwapperdERC20.Contract.AuditSecret(&_SwapperdERC20.CallOpts, _swapID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_SwapperdERC20 *SwapperdERC20CallerSession) AuditSecret(_swapID [32]byte) ([32]byte, error) {
	return _SwapperdERC20.Contract.AuditSecret(&_SwapperdERC20.CallOpts, _swapID)
}

// BrokerFees is a free data retrieval call binding the contract method 0xe1ec380c.
//
// Solidity: function brokerFees( address) constant returns(uint256)
func (_SwapperdERC20 *SwapperdERC20Caller) BrokerFees(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SwapperdERC20.contract.Call(opts, out, "brokerFees", arg0)
	return *ret0, err
}

// BrokerFees is a free data retrieval call binding the contract method 0xe1ec380c.
//
// Solidity: function brokerFees( address) constant returns(uint256)
func (_SwapperdERC20 *SwapperdERC20Session) BrokerFees(arg0 common.Address) (*big.Int, error) {
	return _SwapperdERC20.Contract.BrokerFees(&_SwapperdERC20.CallOpts, arg0)
}

// BrokerFees is a free data retrieval call binding the contract method 0xe1ec380c.
//
// Solidity: function brokerFees( address) constant returns(uint256)
func (_SwapperdERC20 *SwapperdERC20CallerSession) BrokerFees(arg0 common.Address) (*big.Int, error) {
	return _SwapperdERC20.Contract.BrokerFees(&_SwapperdERC20.CallOpts, arg0)
}

// Initiatable is a free data retrieval call binding the contract method 0x09ece618.
//
// Solidity: function initiatable(_swapID bytes32) constant returns(bool)
func (_SwapperdERC20 *SwapperdERC20Caller) Initiatable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SwapperdERC20.contract.Call(opts, out, "initiatable", _swapID)
	return *ret0, err
}

// Initiatable is a free data retrieval call binding the contract method 0x09ece618.
//
// Solidity: function initiatable(_swapID bytes32) constant returns(bool)
func (_SwapperdERC20 *SwapperdERC20Session) Initiatable(_swapID [32]byte) (bool, error) {
	return _SwapperdERC20.Contract.Initiatable(&_SwapperdERC20.CallOpts, _swapID)
}

// Initiatable is a free data retrieval call binding the contract method 0x09ece618.
//
// Solidity: function initiatable(_swapID bytes32) constant returns(bool)
func (_SwapperdERC20 *SwapperdERC20CallerSession) Initiatable(_swapID [32]byte) (bool, error) {
	return _SwapperdERC20.Contract.Initiatable(&_SwapperdERC20.CallOpts, _swapID)
}

// Redeemable is a free data retrieval call binding the contract method 0x68f06b29.
//
// Solidity: function redeemable(_swapID bytes32) constant returns(bool)
func (_SwapperdERC20 *SwapperdERC20Caller) Redeemable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SwapperdERC20.contract.Call(opts, out, "redeemable", _swapID)
	return *ret0, err
}

// Redeemable is a free data retrieval call binding the contract method 0x68f06b29.
//
// Solidity: function redeemable(_swapID bytes32) constant returns(bool)
func (_SwapperdERC20 *SwapperdERC20Session) Redeemable(_swapID [32]byte) (bool, error) {
	return _SwapperdERC20.Contract.Redeemable(&_SwapperdERC20.CallOpts, _swapID)
}

// Redeemable is a free data retrieval call binding the contract method 0x68f06b29.
//
// Solidity: function redeemable(_swapID bytes32) constant returns(bool)
func (_SwapperdERC20 *SwapperdERC20CallerSession) Redeemable(_swapID [32]byte) (bool, error) {
	return _SwapperdERC20.Contract.Redeemable(&_SwapperdERC20.CallOpts, _swapID)
}

// RedeemedAt is a free data retrieval call binding the contract method 0xbc4fcc4a.
//
// Solidity: function redeemedAt( bytes32) constant returns(uint256)
func (_SwapperdERC20 *SwapperdERC20Caller) RedeemedAt(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _SwapperdERC20.contract.Call(opts, out, "redeemedAt", arg0)
	return *ret0, err
}

// RedeemedAt is a free data retrieval call binding the contract method 0xbc4fcc4a.
//
// Solidity: function redeemedAt( bytes32) constant returns(uint256)
func (_SwapperdERC20 *SwapperdERC20Session) RedeemedAt(arg0 [32]byte) (*big.Int, error) {
	return _SwapperdERC20.Contract.RedeemedAt(&_SwapperdERC20.CallOpts, arg0)
}

// RedeemedAt is a free data retrieval call binding the contract method 0xbc4fcc4a.
//
// Solidity: function redeemedAt( bytes32) constant returns(uint256)
func (_SwapperdERC20 *SwapperdERC20CallerSession) RedeemedAt(arg0 [32]byte) (*big.Int, error) {
	return _SwapperdERC20.Contract.RedeemedAt(&_SwapperdERC20.CallOpts, arg0)
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(_swapID bytes32) constant returns(bool)
func (_SwapperdERC20 *SwapperdERC20Caller) Refundable(opts *bind.CallOpts, _swapID [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _SwapperdERC20.contract.Call(opts, out, "refundable", _swapID)
	return *ret0, err
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(_swapID bytes32) constant returns(bool)
func (_SwapperdERC20 *SwapperdERC20Session) Refundable(_swapID [32]byte) (bool, error) {
	return _SwapperdERC20.Contract.Refundable(&_SwapperdERC20.CallOpts, _swapID)
}

// Refundable is a free data retrieval call binding the contract method 0x9fb31475.
//
// Solidity: function refundable(_swapID bytes32) constant returns(bool)
func (_SwapperdERC20 *SwapperdERC20CallerSession) Refundable(_swapID [32]byte) (bool, error) {
	return _SwapperdERC20.Contract.Refundable(&_SwapperdERC20.CallOpts, _swapID)
}

// SwapID is a free data retrieval call binding the contract method 0x4b2ac3fa.
//
// Solidity: function swapID(_secretLock bytes32, _timelock uint256) constant returns(bytes32)
func (_SwapperdERC20 *SwapperdERC20Caller) SwapID(opts *bind.CallOpts, _secretLock [32]byte, _timelock *big.Int) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _SwapperdERC20.contract.Call(opts, out, "swapID", _secretLock, _timelock)
	return *ret0, err
}

// SwapID is a free data retrieval call binding the contract method 0x4b2ac3fa.
//
// Solidity: function swapID(_secretLock bytes32, _timelock uint256) constant returns(bytes32)
func (_SwapperdERC20 *SwapperdERC20Session) SwapID(_secretLock [32]byte, _timelock *big.Int) ([32]byte, error) {
	return _SwapperdERC20.Contract.SwapID(&_SwapperdERC20.CallOpts, _secretLock, _timelock)
}

// SwapID is a free data retrieval call binding the contract method 0x4b2ac3fa.
//
// Solidity: function swapID(_secretLock bytes32, _timelock uint256) constant returns(bytes32)
func (_SwapperdERC20 *SwapperdERC20CallerSession) SwapID(_secretLock [32]byte, _timelock *big.Int) ([32]byte, error) {
	return _SwapperdERC20.Contract.SwapID(&_SwapperdERC20.CallOpts, _secretLock, _timelock)
}

// Initiate is a paid mutator transaction binding the contract method 0x027a2577.
//
// Solidity: function initiate(_swapID bytes32, _spender address, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_SwapperdERC20 *SwapperdERC20Transactor) Initiate(opts *bind.TransactOpts, _swapID [32]byte, _spender common.Address, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _SwapperdERC20.contract.Transact(opts, "initiate", _swapID, _spender, _secretLock, _timelock, _value)
}

// Initiate is a paid mutator transaction binding the contract method 0x027a2577.
//
// Solidity: function initiate(_swapID bytes32, _spender address, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_SwapperdERC20 *SwapperdERC20Session) Initiate(_swapID [32]byte, _spender common.Address, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _SwapperdERC20.Contract.Initiate(&_SwapperdERC20.TransactOpts, _swapID, _spender, _secretLock, _timelock, _value)
}

// Initiate is a paid mutator transaction binding the contract method 0x027a2577.
//
// Solidity: function initiate(_swapID bytes32, _spender address, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_SwapperdERC20 *SwapperdERC20TransactorSession) Initiate(_swapID [32]byte, _spender common.Address, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _SwapperdERC20.Contract.Initiate(&_SwapperdERC20.TransactOpts, _swapID, _spender, _secretLock, _timelock, _value)
}

// InitiateWithFees is a paid mutator transaction binding the contract method 0xb8688e3f.
//
// Solidity: function initiateWithFees(_swapID bytes32, _spender address, _broker address, _brokerFee uint256, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_SwapperdERC20 *SwapperdERC20Transactor) InitiateWithFees(opts *bind.TransactOpts, _swapID [32]byte, _spender common.Address, _broker common.Address, _brokerFee *big.Int, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _SwapperdERC20.contract.Transact(opts, "initiateWithFees", _swapID, _spender, _broker, _brokerFee, _secretLock, _timelock, _value)
}

// InitiateWithFees is a paid mutator transaction binding the contract method 0xb8688e3f.
//
// Solidity: function initiateWithFees(_swapID bytes32, _spender address, _broker address, _brokerFee uint256, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_SwapperdERC20 *SwapperdERC20Session) InitiateWithFees(_swapID [32]byte, _spender common.Address, _broker common.Address, _brokerFee *big.Int, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _SwapperdERC20.Contract.InitiateWithFees(&_SwapperdERC20.TransactOpts, _swapID, _spender, _broker, _brokerFee, _secretLock, _timelock, _value)
}

// InitiateWithFees is a paid mutator transaction binding the contract method 0xb8688e3f.
//
// Solidity: function initiateWithFees(_swapID bytes32, _spender address, _broker address, _brokerFee uint256, _secretLock bytes32, _timelock uint256, _value uint256) returns()
func (_SwapperdERC20 *SwapperdERC20TransactorSession) InitiateWithFees(_swapID [32]byte, _spender common.Address, _broker common.Address, _brokerFee *big.Int, _secretLock [32]byte, _timelock *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _SwapperdERC20.Contract.InitiateWithFees(&_SwapperdERC20.TransactOpts, _swapID, _spender, _broker, _brokerFee, _secretLock, _timelock, _value)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_SwapperdERC20 *SwapperdERC20Transactor) Redeem(opts *bind.TransactOpts, _swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _SwapperdERC20.contract.Transact(opts, "redeem", _swapID, _secretKey)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_SwapperdERC20 *SwapperdERC20Session) Redeem(_swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _SwapperdERC20.Contract.Redeem(&_SwapperdERC20.TransactOpts, _swapID, _secretKey)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_SwapperdERC20 *SwapperdERC20TransactorSession) Redeem(_swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _SwapperdERC20.Contract.Redeem(&_SwapperdERC20.TransactOpts, _swapID, _secretKey)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_SwapperdERC20 *SwapperdERC20Transactor) Refund(opts *bind.TransactOpts, _swapID [32]byte) (*types.Transaction, error) {
	return _SwapperdERC20.contract.Transact(opts, "refund", _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_SwapperdERC20 *SwapperdERC20Session) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _SwapperdERC20.Contract.Refund(&_SwapperdERC20.TransactOpts, _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_SwapperdERC20 *SwapperdERC20TransactorSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _SwapperdERC20.Contract.Refund(&_SwapperdERC20.TransactOpts, _swapID)
}

// WithdrawBrokerFees is a paid mutator transaction binding the contract method 0x4c6d37ff.
//
// Solidity: function withdrawBrokerFees(_amount uint256) returns()
func (_SwapperdERC20 *SwapperdERC20Transactor) WithdrawBrokerFees(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _SwapperdERC20.contract.Transact(opts, "withdrawBrokerFees", _amount)
}

// WithdrawBrokerFees is a paid mutator transaction binding the contract method 0x4c6d37ff.
//
// Solidity: function withdrawBrokerFees(_amount uint256) returns()
func (_SwapperdERC20 *SwapperdERC20Session) WithdrawBrokerFees(_amount *big.Int) (*types.Transaction, error) {
	return _SwapperdERC20.Contract.WithdrawBrokerFees(&_SwapperdERC20.TransactOpts, _amount)
}

// WithdrawBrokerFees is a paid mutator transaction binding the contract method 0x4c6d37ff.
//
// Solidity: function withdrawBrokerFees(_amount uint256) returns()
func (_SwapperdERC20 *SwapperdERC20TransactorSession) WithdrawBrokerFees(_amount *big.Int) (*types.Transaction, error) {
	return _SwapperdERC20.Contract.WithdrawBrokerFees(&_SwapperdERC20.TransactOpts, _amount)
}

// SwapperdERC20LogCloseIterator is returned from FilterLogClose and is used to iterate over the raw logs and unpacked data for LogClose events raised by the SwapperdERC20 contract.
type SwapperdERC20LogCloseIterator struct {
	Event *SwapperdERC20LogClose // Event containing the contract specifics and raw log

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
func (it *SwapperdERC20LogCloseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SwapperdERC20LogClose)
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
		it.Event = new(SwapperdERC20LogClose)
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
func (it *SwapperdERC20LogCloseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SwapperdERC20LogCloseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SwapperdERC20LogClose represents a LogClose event raised by the SwapperdERC20 contract.
type SwapperdERC20LogClose struct {
	SwapID    [32]byte
	SecretKey [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLogClose is a free log retrieval operation binding the contract event 0x07da1fa25a1d885732677ce9c192cbec27051a4b69d391c9a64850f5a5112ba0.
//
// Solidity: e LogClose(_swapID bytes32, _secretKey bytes32)
func (_SwapperdERC20 *SwapperdERC20Filterer) FilterLogClose(opts *bind.FilterOpts) (*SwapperdERC20LogCloseIterator, error) {

	logs, sub, err := _SwapperdERC20.contract.FilterLogs(opts, "LogClose")
	if err != nil {
		return nil, err
	}
	return &SwapperdERC20LogCloseIterator{contract: _SwapperdERC20.contract, event: "LogClose", logs: logs, sub: sub}, nil
}

// WatchLogClose is a free log subscription operation binding the contract event 0x07da1fa25a1d885732677ce9c192cbec27051a4b69d391c9a64850f5a5112ba0.
//
// Solidity: e LogClose(_swapID bytes32, _secretKey bytes32)
func (_SwapperdERC20 *SwapperdERC20Filterer) WatchLogClose(opts *bind.WatchOpts, sink chan<- *SwapperdERC20LogClose) (event.Subscription, error) {

	logs, sub, err := _SwapperdERC20.contract.WatchLogs(opts, "LogClose")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SwapperdERC20LogClose)
				if err := _SwapperdERC20.contract.UnpackLog(event, "LogClose", log); err != nil {
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

// SwapperdERC20LogExpireIterator is returned from FilterLogExpire and is used to iterate over the raw logs and unpacked data for LogExpire events raised by the SwapperdERC20 contract.
type SwapperdERC20LogExpireIterator struct {
	Event *SwapperdERC20LogExpire // Event containing the contract specifics and raw log

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
func (it *SwapperdERC20LogExpireIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SwapperdERC20LogExpire)
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
		it.Event = new(SwapperdERC20LogExpire)
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
func (it *SwapperdERC20LogExpireIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SwapperdERC20LogExpireIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SwapperdERC20LogExpire represents a LogExpire event raised by the SwapperdERC20 contract.
type SwapperdERC20LogExpire struct {
	SwapID [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterLogExpire is a free log retrieval operation binding the contract event 0xeb711459e1247171382f0da0387b86239b8e3ca60af3b15a9ff2f1eb3d7f6a1d.
//
// Solidity: e LogExpire(_swapID bytes32)
func (_SwapperdERC20 *SwapperdERC20Filterer) FilterLogExpire(opts *bind.FilterOpts) (*SwapperdERC20LogExpireIterator, error) {

	logs, sub, err := _SwapperdERC20.contract.FilterLogs(opts, "LogExpire")
	if err != nil {
		return nil, err
	}
	return &SwapperdERC20LogExpireIterator{contract: _SwapperdERC20.contract, event: "LogExpire", logs: logs, sub: sub}, nil
}

// WatchLogExpire is a free log subscription operation binding the contract event 0xeb711459e1247171382f0da0387b86239b8e3ca60af3b15a9ff2f1eb3d7f6a1d.
//
// Solidity: e LogExpire(_swapID bytes32)
func (_SwapperdERC20 *SwapperdERC20Filterer) WatchLogExpire(opts *bind.WatchOpts, sink chan<- *SwapperdERC20LogExpire) (event.Subscription, error) {

	logs, sub, err := _SwapperdERC20.contract.WatchLogs(opts, "LogExpire")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SwapperdERC20LogExpire)
				if err := _SwapperdERC20.contract.UnpackLog(event, "LogExpire", log); err != nil {
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

// SwapperdERC20LogOpenIterator is returned from FilterLogOpen and is used to iterate over the raw logs and unpacked data for LogOpen events raised by the SwapperdERC20 contract.
type SwapperdERC20LogOpenIterator struct {
	Event *SwapperdERC20LogOpen // Event containing the contract specifics and raw log

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
func (it *SwapperdERC20LogOpenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SwapperdERC20LogOpen)
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
		it.Event = new(SwapperdERC20LogOpen)
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
func (it *SwapperdERC20LogOpenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SwapperdERC20LogOpenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SwapperdERC20LogOpen represents a LogOpen event raised by the SwapperdERC20 contract.
type SwapperdERC20LogOpen struct {
	SwapID     [32]byte
	Spender    common.Address
	SecretLock [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogOpen is a free log retrieval operation binding the contract event 0x497d46e9505eefe8b910d1a02e6b40d8769510023b0053c3ac4b9574b81c97bf.
//
// Solidity: e LogOpen(_swapID bytes32, _spender address, _secretLock bytes32)
func (_SwapperdERC20 *SwapperdERC20Filterer) FilterLogOpen(opts *bind.FilterOpts) (*SwapperdERC20LogOpenIterator, error) {

	logs, sub, err := _SwapperdERC20.contract.FilterLogs(opts, "LogOpen")
	if err != nil {
		return nil, err
	}
	return &SwapperdERC20LogOpenIterator{contract: _SwapperdERC20.contract, event: "LogOpen", logs: logs, sub: sub}, nil
}

// WatchLogOpen is a free log subscription operation binding the contract event 0x497d46e9505eefe8b910d1a02e6b40d8769510023b0053c3ac4b9574b81c97bf.
//
// Solidity: e LogOpen(_swapID bytes32, _spender address, _secretLock bytes32)
func (_SwapperdERC20 *SwapperdERC20Filterer) WatchLogOpen(opts *bind.WatchOpts, sink chan<- *SwapperdERC20LogOpen) (event.Subscription, error) {

	logs, sub, err := _SwapperdERC20.contract.WatchLogs(opts, "LogOpen")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SwapperdERC20LogOpen)
				if err := _SwapperdERC20.contract.UnpackLog(event, "LogOpen", log); err != nil {
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
