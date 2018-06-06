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

// ArcABI is the input ABI used to generate the binding from.
const ArcABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_withdrawTrader\",\"type\":\"address\"},{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"}],\"name\":\"initiate\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"auditSecret\",\"outputs\":[{\"name\":\"secretKey\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_secretKey\",\"type\":\"bytes32\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"audit\",\"outputs\":[{\"name\":\"timelock\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"secretLock\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_withdrawTrader\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_secretLock\",\"type\":\"bytes32\"}],\"name\":\"Open\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"Expire\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_secretKey\",\"type\":\"bytes\"}],\"name\":\"Close\",\"type\":\"event\"}]"

// ArcBin is the compiled bytecode used for deploying new contracts.
const ArcBin = `0x608060405234801561001057600080fd5b506105ef806100206000396000f30060806040526004361061006c5763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663412c0b5881146100715780637249fbb614610090578063976d00f4146100a8578063b31597ad146100d2578063c140635b146100ed575b600080fd5b61008e600435600160a060020a036024351660443560643561013a565b005b34801561009c57600080fd5b5061008e600435610223565b3480156100b457600080fd5b506100c060043561031e565b60408051918252519081900360200190f35b3480156100de57600080fd5b5061008e6004356024356103bc565b3480156100f957600080fd5b50610105600435610508565b604080519586526020860194909452600160a060020a0392831685850152911660608401526080830152519081900360a00190f35b61014261058e565b846000808281526001602052604090205460ff16600381111561016157fe5b1461016b57600080fd5b50506040805160c08101825291825234602080840191825233848401908152600160a060020a039687166060860190815260808601968752600060a087018181529981528084528581209651875593516001808801919091559151600287018054918a1673ffffffffffffffffffffffffffffffffffffffff1992831617905590516003870180549190991691161790965593516004840155945160059092019190915590829052909120805460ff19169091179055565b61022b61058e565b81600160008281526001602052604090205460ff16600381111561024b57fe5b1461025557600080fd5b600083815260208190526040902054839042101561027257600080fd5b600084815260208181526040808320815160c081018352815481526001808301548286019081526002840154600160a060020a03908116848701908152600380870154831660608701526004870154608087015260059096015460a08601528c895292909652848720805460ff1916909417909355519151925190975092169281156108fc029290818181858888f19350505050158015610317573d6000803e3d6000fd5b5050505050565b600061032861058e565b82600260008281526001602052604090205460ff16600381111561034857fe5b1461035257600080fd5b50505060009081526020818152604091829020825160c081018452815481526001820154928101929092526002810154600160a060020a0390811693830193909352600381015490921660608201526004820154608082015260059091015460a090910181905290565b6103c461058e565b82600160008281526001602052604090205460ff1660038111156103e457fe5b146103ee57600080fd5b6040805184815290518591859160029160208082019290918190038201816000865af1158015610422573d6000803e3d6000fd5b5050506040513d602081101561043757600080fd5b50516000838152602081905260409020600401541461045557600080fd5b600086815260208181526040808320815160c08101835281548152600180830154828601908152600280850154600160a060020a03908116858801526003860154811660608601908152600487015460808701526005909601805460a08701528f8a528e905592909652848720805460ff191690961790955591519351925190985092169281156108fc029290818181858888f193505050501580156104ff573d6000803e3d6000fd5b50505050505050565b600080600080600061051861058e565b50505060009384525050506020818152604091829020825160c081018452815480825260018301549382018490526002830154600160a060020a039081169583018690526003840154166060830181905260048401546080840181905260059094015460a0909301929092529492939092909190565b6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a0810191909152905600a165627a7a7230582027456d6d9fa57b26f3b88353df661c6452f0f228cce3eec66b164e1d985fe8ed0029`

// DeployArc deploys a new Ethereum contract, binding an instance of Arc to it.
func DeployArc(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Arc, error) {
	parsed, err := abi.JSON(strings.NewReader(ArcABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ArcBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Arc{ArcCaller: ArcCaller{contract: contract}, ArcTransactor: ArcTransactor{contract: contract}, ArcFilterer: ArcFilterer{contract: contract}}, nil
}

// Arc is an auto generated Go binding around an Ethereum contract.
type Arc struct {
	ArcCaller     // Read-only binding to the contract
	ArcTransactor // Write-only binding to the contract
	ArcFilterer   // Log filterer for contract events
}

// ArcCaller is an auto generated read-only Go binding around an Ethereum contract.
type ArcCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArcTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ArcTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArcFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ArcFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArcSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ArcSession struct {
	Contract     *Arc              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ArcCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ArcCallerSession struct {
	Contract *ArcCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ArcTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ArcTransactorSession struct {
	Contract     *ArcTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ArcRaw is an auto generated low-level Go binding around an Ethereum contract.
type ArcRaw struct {
	Contract *Arc // Generic contract binding to access the raw methods on
}

// ArcCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ArcCallerRaw struct {
	Contract *ArcCaller // Generic read-only contract binding to access the raw methods on
}

// ArcTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ArcTransactorRaw struct {
	Contract *ArcTransactor // Generic write-only contract binding to access the raw methods on
}

// NewArc creates a new instance of Arc, bound to a specific deployed contract.
func NewArc(address common.Address, backend bind.ContractBackend) (*Arc, error) {
	contract, err := bindArc(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Arc{ArcCaller: ArcCaller{contract: contract}, ArcTransactor: ArcTransactor{contract: contract}, ArcFilterer: ArcFilterer{contract: contract}}, nil
}

// NewArcCaller creates a new read-only instance of Arc, bound to a specific deployed contract.
func NewArcCaller(address common.Address, caller bind.ContractCaller) (*ArcCaller, error) {
	contract, err := bindArc(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ArcCaller{contract: contract}, nil
}

// NewArcTransactor creates a new write-only instance of Arc, bound to a specific deployed contract.
func NewArcTransactor(address common.Address, transactor bind.ContractTransactor) (*ArcTransactor, error) {
	contract, err := bindArc(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ArcTransactor{contract: contract}, nil
}

// NewArcFilterer creates a new log filterer instance of Arc, bound to a specific deployed contract.
func NewArcFilterer(address common.Address, filterer bind.ContractFilterer) (*ArcFilterer, error) {
	contract, err := bindArc(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ArcFilterer{contract: contract}, nil
}

// bindArc binds a generic wrapper to an already deployed contract.
func bindArc(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ArcABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Arc *ArcRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Arc.Contract.ArcCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Arc *ArcRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Arc.Contract.ArcTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Arc *ArcRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Arc.Contract.ArcTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Arc *ArcCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Arc.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Arc *ArcTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Arc.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Arc *ArcTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Arc.Contract.contract.Transact(opts, method, params...)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, from address, secretLock bytes32)
func (_Arc *ArcCaller) Audit(opts *bind.CallOpts, _swapID [32]byte) (struct {
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
	err := _Arc.contract.Call(opts, out, "audit", _swapID)
	return *ret, err
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, from address, secretLock bytes32)
func (_Arc *ArcSession) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	return _Arc.Contract.Audit(&_Arc.CallOpts, _swapID)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, from address, secretLock bytes32)
func (_Arc *ArcCallerSession) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	return _Arc.Contract.Audit(&_Arc.CallOpts, _swapID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_Arc *ArcCaller) AuditSecret(opts *bind.CallOpts, _swapID [32]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Arc.contract.Call(opts, out, "auditSecret", _swapID)
	return *ret0, err
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_Arc *ArcSession) AuditSecret(_swapID [32]byte) ([32]byte, error) {
	return _Arc.Contract.AuditSecret(&_Arc.CallOpts, _swapID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_Arc *ArcCallerSession) AuditSecret(_swapID [32]byte) ([32]byte, error) {
	return _Arc.Contract.AuditSecret(&_Arc.CallOpts, _swapID)
}

// Initiate is a paid mutator transaction binding the contract method 0x412c0b58.
//
// Solidity: function initiate(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_Arc *ArcTransactor) Initiate(opts *bind.TransactOpts, _swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _Arc.contract.Transact(opts, "initiate", _swapID, _withdrawTrader, _secretLock, _timelock)
}

// Initiate is a paid mutator transaction binding the contract method 0x412c0b58.
//
// Solidity: function initiate(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_Arc *ArcSession) Initiate(_swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _Arc.Contract.Initiate(&_Arc.TransactOpts, _swapID, _withdrawTrader, _secretLock, _timelock)
}

// Initiate is a paid mutator transaction binding the contract method 0x412c0b58.
//
// Solidity: function initiate(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_Arc *ArcTransactorSession) Initiate(_swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _Arc.Contract.Initiate(&_Arc.TransactOpts, _swapID, _withdrawTrader, _secretLock, _timelock)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_Arc *ArcTransactor) Redeem(opts *bind.TransactOpts, _swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _Arc.contract.Transact(opts, "redeem", _swapID, _secretKey)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_Arc *ArcSession) Redeem(_swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _Arc.Contract.Redeem(&_Arc.TransactOpts, _swapID, _secretKey)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_Arc *ArcTransactorSession) Redeem(_swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _Arc.Contract.Redeem(&_Arc.TransactOpts, _swapID, _secretKey)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_Arc *ArcTransactor) Refund(opts *bind.TransactOpts, _swapID [32]byte) (*types.Transaction, error) {
	return _Arc.contract.Transact(opts, "refund", _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_Arc *ArcSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _Arc.Contract.Refund(&_Arc.TransactOpts, _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_Arc *ArcTransactorSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _Arc.Contract.Refund(&_Arc.TransactOpts, _swapID)
}

// ArcCloseIterator is returned from FilterClose and is used to iterate over the raw logs and unpacked data for Close events raised by the Arc contract.
type ArcCloseIterator struct {
	Event *ArcClose // Event containing the contract specifics and raw log

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
func (it *ArcCloseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArcClose)
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
		it.Event = new(ArcClose)
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
func (it *ArcCloseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ArcCloseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ArcClose represents a Close event raised by the Arc contract.
type ArcClose struct {
	SwapID    [32]byte
	SecretKey []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterClose is a free log retrieval operation binding the contract event 0x692fd10a275135b9a2a2f5819db3d9965a5129ea2ad3640a0156dbce2fc81bdd.
//
// Solidity: e Close(_swapID bytes32, _secretKey bytes)
func (_Arc *ArcFilterer) FilterClose(opts *bind.FilterOpts) (*ArcCloseIterator, error) {

	logs, sub, err := _Arc.contract.FilterLogs(opts, "Close")
	if err != nil {
		return nil, err
	}
	return &ArcCloseIterator{contract: _Arc.contract, event: "Close", logs: logs, sub: sub}, nil
}

// WatchClose is a free log subscription operation binding the contract event 0x692fd10a275135b9a2a2f5819db3d9965a5129ea2ad3640a0156dbce2fc81bdd.
//
// Solidity: e Close(_swapID bytes32, _secretKey bytes)
func (_Arc *ArcFilterer) WatchClose(opts *bind.WatchOpts, sink chan<- *ArcClose) (event.Subscription, error) {

	logs, sub, err := _Arc.contract.WatchLogs(opts, "Close")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ArcClose)
				if err := _Arc.contract.UnpackLog(event, "Close", log); err != nil {
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

// ArcExpireIterator is returned from FilterExpire and is used to iterate over the raw logs and unpacked data for Expire events raised by the Arc contract.
type ArcExpireIterator struct {
	Event *ArcExpire // Event containing the contract specifics and raw log

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
func (it *ArcExpireIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArcExpire)
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
		it.Event = new(ArcExpire)
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
func (it *ArcExpireIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ArcExpireIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ArcExpire represents a Expire event raised by the Arc contract.
type ArcExpire struct {
	SwapID [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterExpire is a free log retrieval operation binding the contract event 0xbddd9b693ea862fad6ecf78fd51c065be26fda94d1f3cad3a7d691453a38a735.
//
// Solidity: e Expire(_swapID bytes32)
func (_Arc *ArcFilterer) FilterExpire(opts *bind.FilterOpts) (*ArcExpireIterator, error) {

	logs, sub, err := _Arc.contract.FilterLogs(opts, "Expire")
	if err != nil {
		return nil, err
	}
	return &ArcExpireIterator{contract: _Arc.contract, event: "Expire", logs: logs, sub: sub}, nil
}

// WatchExpire is a free log subscription operation binding the contract event 0xbddd9b693ea862fad6ecf78fd51c065be26fda94d1f3cad3a7d691453a38a735.
//
// Solidity: e Expire(_swapID bytes32)
func (_Arc *ArcFilterer) WatchExpire(opts *bind.WatchOpts, sink chan<- *ArcExpire) (event.Subscription, error) {

	logs, sub, err := _Arc.contract.WatchLogs(opts, "Expire")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ArcExpire)
				if err := _Arc.contract.UnpackLog(event, "Expire", log); err != nil {
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

// ArcOpenIterator is returned from FilterOpen and is used to iterate over the raw logs and unpacked data for Open events raised by the Arc contract.
type ArcOpenIterator struct {
	Event *ArcOpen // Event containing the contract specifics and raw log

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
func (it *ArcOpenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArcOpen)
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
		it.Event = new(ArcOpen)
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
func (it *ArcOpenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ArcOpenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ArcOpen represents a Open event raised by the Arc contract.
type ArcOpen struct {
	SwapID         [32]byte
	WithdrawTrader common.Address
	SecretLock     [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterOpen is a free log retrieval operation binding the contract event 0x6ed79a08bf5c8a7d4a330df315e4ac386627ecafbe5d2bfd6654237d967b24f3.
//
// Solidity: e Open(_swapID bytes32, _withdrawTrader address, _secretLock bytes32)
func (_Arc *ArcFilterer) FilterOpen(opts *bind.FilterOpts) (*ArcOpenIterator, error) {

	logs, sub, err := _Arc.contract.FilterLogs(opts, "Open")
	if err != nil {
		return nil, err
	}
	return &ArcOpenIterator{contract: _Arc.contract, event: "Open", logs: logs, sub: sub}, nil
}

// WatchOpen is a free log subscription operation binding the contract event 0x6ed79a08bf5c8a7d4a330df315e4ac386627ecafbe5d2bfd6654237d967b24f3.
//
// Solidity: e Open(_swapID bytes32, _withdrawTrader address, _secretLock bytes32)
func (_Arc *ArcFilterer) WatchOpen(opts *bind.WatchOpts, sink chan<- *ArcOpen) (event.Subscription, error) {

	logs, sub, err := _Arc.contract.WatchLogs(opts, "Open")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ArcOpen)
				if err := _Arc.contract.UnpackLog(event, "Open", log); err != nil {
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
