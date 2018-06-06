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

// AtomABI is the input ABI used to generate the binding from.
const AtomABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_withdrawTrader\",\"type\":\"address\"},{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"}],\"name\":\"initiate\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_orderID\",\"type\":\"bytes32\"},{\"name\":\"_swapDetails\",\"type\":\"bytes\"}],\"name\":\"submitDetails\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"auditSecret\",\"outputs\":[{\"name\":\"secretKey\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"swapDetails\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_secretKey\",\"type\":\"bytes32\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"audit\",\"outputs\":[{\"name\":\"timelock\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"secretLock\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_withdrawTrader\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_secretLock\",\"type\":\"bytes32\"}],\"name\":\"Open\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"Expire\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_secretKey\",\"type\":\"bytes\"}],\"name\":\"Close\",\"type\":\"event\"}]"

// AtomBin is the compiled bytecode used for deploying new contracts.
const AtomBin = `0x608060405234801561001057600080fd5b50610848806100206000396000f3006080604052600436106100825763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663412c0b5881146100875780637249fbb6146100a657806380ee79d4146100be578063976d00f41461011c578063b14631bb14610146578063b31597ad146101d3578063c140635b146101ee575b600080fd5b6100a4600435600160a060020a036024351660443560643561023b565b005b3480156100b257600080fd5b506100a4600435610324565b3480156100ca57600080fd5b5060408051602060046024803582810135601f81018590048502860185019096528585526100a495833595369560449491939091019190819084018382808284375094975061041f9650505050505050565b34801561012857600080fd5b50610134600435610443565b60408051918252519081900360200190f35b34801561015257600080fd5b5061015e6004356104e1565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610198578181015183820152602001610180565b50505050905090810190601f1680156101c55780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156101df57600080fd5b506100a460043560243561057a565b3480156101fa57600080fd5b506102066004356106c6565b604080519586526020860194909452600160a060020a0392831685850152911660608401526080830152519081900360a00190f35b61024361074c565b846000808281526001602052604090205460ff16600381111561026257fe5b1461026c57600080fd5b50506040805160c08101825291825234602080840191825233848401908152600160a060020a039687166060860190815260808601968752600060a087018181529981528084528581209651875593516001808801919091559151600287018054918a1673ffffffffffffffffffffffffffffffffffffffff1992831617905590516003870180549190991691161790965593516004840155945160059092019190915590829052909120805460ff19169091179055565b61032c61074c565b81600160008281526001602052604090205460ff16600381111561034c57fe5b1461035657600080fd5b600083815260208190526040902054839042101561037357600080fd5b600084815260208181526040808320815160c081018352815481526001808301548286019081526002840154600160a060020a03908116848701908152600380870154831660608701526004870154608087015260059096015460a08601528c895292909652848720805460ff1916909417909355519151925190975092169281156108fc029290818181858888f19350505050158015610418573d6000803e3d6000fd5b5050505050565b6000828152600260209081526040909120825161043e92840190610781565b505050565b600061044d61074c565b82600260008281526001602052604090205460ff16600381111561046d57fe5b1461047757600080fd5b50505060009081526020818152604091829020825160c081018452815481526001820154928101929092526002810154600160a060020a0390811693830193909352600381015490921660608201526004820154608082015260059091015460a090910181905290565b600260208181526000928352604092839020805484516001821615610100026000190190911693909304601f81018390048302840183019094528383529192908301828280156105725780601f1061054757610100808354040283529160200191610572565b820191906000526020600020905b81548152906001019060200180831161055557829003601f168201915b505050505081565b61058261074c565b82600160008281526001602052604090205460ff1660038111156105a257fe5b146105ac57600080fd5b6040805184815290518591859160029160208082019290918190038201816000865af11580156105e0573d6000803e3d6000fd5b5050506040513d60208110156105f557600080fd5b50516000838152602081905260409020600401541461061357600080fd5b600086815260208181526040808320815160c08101835281548152600180830154828601908152600280850154600160a060020a03908116858801526003860154811660608601908152600487015460808701526005909601805460a08701528f8a528e905592909652848720805460ff191690961790955591519351925190985092169281156108fc029290818181858888f193505050501580156106bd573d6000803e3d6000fd5b50505050505050565b60008060008060006106d661074c565b50505060009384525050506020818152604091829020825160c081018452815480825260018301549382018490526002830154600160a060020a039081169583018690526003840154166060830181905260048401546080840181905260059094015460a0909301929092529492939092909190565b6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a081019190915290565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106107c257805160ff19168380011785556107ef565b828001600101855582156107ef579182015b828111156107ef5782518255916020019190600101906107d4565b506107fb9291506107ff565b5090565b61081991905b808211156107fb5760008155600101610805565b905600a165627a7a72305820b885a5010d97b31797ada37060aebbc2d0e4701a9d0d828dfd760b8df2ac55810029`

// DeployAtom deploys a new Ethereum contract, binding an instance of Atom to it.
func DeployAtom(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Atom, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AtomBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Atom{AtomCaller: AtomCaller{contract: contract}, AtomTransactor: AtomTransactor{contract: contract}, AtomFilterer: AtomFilterer{contract: contract}}, nil
}

// Atom is an auto generated Go binding around an Ethereum contract.
type Atom struct {
	AtomCaller     // Read-only binding to the contract
	AtomTransactor // Write-only binding to the contract
	AtomFilterer   // Log filterer for contract events
}

// AtomCaller is an auto generated read-only Go binding around an Ethereum contract.
type AtomCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AtomTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AtomFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AtomSession struct {
	Contract     *Atom             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AtomCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AtomCallerSession struct {
	Contract *AtomCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// AtomTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AtomTransactorSession struct {
	Contract     *AtomTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AtomRaw is an auto generated low-level Go binding around an Ethereum contract.
type AtomRaw struct {
	Contract *Atom // Generic contract binding to access the raw methods on
}

// AtomCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AtomCallerRaw struct {
	Contract *AtomCaller // Generic read-only contract binding to access the raw methods on
}

// AtomTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AtomTransactorRaw struct {
	Contract *AtomTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAtom creates a new instance of Atom, bound to a specific deployed contract.
func NewAtom(address common.Address, backend bind.ContractBackend) (*Atom, error) {
	contract, err := bindAtom(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Atom{AtomCaller: AtomCaller{contract: contract}, AtomTransactor: AtomTransactor{contract: contract}, AtomFilterer: AtomFilterer{contract: contract}}, nil
}

// NewAtomCaller creates a new read-only instance of Atom, bound to a specific deployed contract.
func NewAtomCaller(address common.Address, caller bind.ContractCaller) (*AtomCaller, error) {
	contract, err := bindAtom(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AtomCaller{contract: contract}, nil
}

// NewAtomTransactor creates a new write-only instance of Atom, bound to a specific deployed contract.
func NewAtomTransactor(address common.Address, transactor bind.ContractTransactor) (*AtomTransactor, error) {
	contract, err := bindAtom(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AtomTransactor{contract: contract}, nil
}

// NewAtomFilterer creates a new log filterer instance of Atom, bound to a specific deployed contract.
func NewAtomFilterer(address common.Address, filterer bind.ContractFilterer) (*AtomFilterer, error) {
	contract, err := bindAtom(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AtomFilterer{contract: contract}, nil
}

// bindAtom binds a generic wrapper to an already deployed contract.
func bindAtom(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Atom *AtomRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Atom.Contract.AtomCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Atom *AtomRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Atom.Contract.AtomTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Atom *AtomRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Atom.Contract.AtomTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Atom *AtomCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Atom.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Atom *AtomTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Atom.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Atom *AtomTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Atom.Contract.contract.Transact(opts, method, params...)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, from address, secretLock bytes32)
func (_Atom *AtomCaller) Audit(opts *bind.CallOpts, _swapID [32]byte) (struct {
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
	err := _Atom.contract.Call(opts, out, "audit", _swapID)
	return *ret, err
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, from address, secretLock bytes32)
func (_Atom *AtomSession) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	return _Atom.Contract.Audit(&_Atom.CallOpts, _swapID)
}

// Audit is a free data retrieval call binding the contract method 0xc140635b.
//
// Solidity: function audit(_swapID bytes32) constant returns(timelock uint256, value uint256, to address, from address, secretLock bytes32)
func (_Atom *AtomCallerSession) Audit(_swapID [32]byte) (struct {
	Timelock   *big.Int
	Value      *big.Int
	To         common.Address
	From       common.Address
	SecretLock [32]byte
}, error) {
	return _Atom.Contract.Audit(&_Atom.CallOpts, _swapID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_Atom *AtomCaller) AuditSecret(opts *bind.CallOpts, _swapID [32]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Atom.contract.Call(opts, out, "auditSecret", _swapID)
	return *ret0, err
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_Atom *AtomSession) AuditSecret(_swapID [32]byte) ([32]byte, error) {
	return _Atom.Contract.AuditSecret(&_Atom.CallOpts, _swapID)
}

// AuditSecret is a free data retrieval call binding the contract method 0x976d00f4.
//
// Solidity: function auditSecret(_swapID bytes32) constant returns(secretKey bytes32)
func (_Atom *AtomCallerSession) AuditSecret(_swapID [32]byte) ([32]byte, error) {
	return _Atom.Contract.AuditSecret(&_Atom.CallOpts, _swapID)
}

// SwapDetails is a free data retrieval call binding the contract method 0xb14631bb.
//
// Solidity: function swapDetails( bytes32) constant returns(bytes)
func (_Atom *AtomCaller) SwapDetails(opts *bind.CallOpts, arg0 [32]byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _Atom.contract.Call(opts, out, "swapDetails", arg0)
	return *ret0, err
}

// SwapDetails is a free data retrieval call binding the contract method 0xb14631bb.
//
// Solidity: function swapDetails( bytes32) constant returns(bytes)
func (_Atom *AtomSession) SwapDetails(arg0 [32]byte) ([]byte, error) {
	return _Atom.Contract.SwapDetails(&_Atom.CallOpts, arg0)
}

// SwapDetails is a free data retrieval call binding the contract method 0xb14631bb.
//
// Solidity: function swapDetails( bytes32) constant returns(bytes)
func (_Atom *AtomCallerSession) SwapDetails(arg0 [32]byte) ([]byte, error) {
	return _Atom.Contract.SwapDetails(&_Atom.CallOpts, arg0)
}

// Initiate is a paid mutator transaction binding the contract method 0x412c0b58.
//
// Solidity: function initiate(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_Atom *AtomTransactor) Initiate(opts *bind.TransactOpts, _swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _Atom.contract.Transact(opts, "initiate", _swapID, _withdrawTrader, _secretLock, _timelock)
}

// Initiate is a paid mutator transaction binding the contract method 0x412c0b58.
//
// Solidity: function initiate(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_Atom *AtomSession) Initiate(_swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _Atom.Contract.Initiate(&_Atom.TransactOpts, _swapID, _withdrawTrader, _secretLock, _timelock)
}

// Initiate is a paid mutator transaction binding the contract method 0x412c0b58.
//
// Solidity: function initiate(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_Atom *AtomTransactorSession) Initiate(_swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _Atom.Contract.Initiate(&_Atom.TransactOpts, _swapID, _withdrawTrader, _secretLock, _timelock)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_Atom *AtomTransactor) Redeem(opts *bind.TransactOpts, _swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _Atom.contract.Transact(opts, "redeem", _swapID, _secretKey)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_Atom *AtomSession) Redeem(_swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _Atom.Contract.Redeem(&_Atom.TransactOpts, _swapID, _secretKey)
}

// Redeem is a paid mutator transaction binding the contract method 0xb31597ad.
//
// Solidity: function redeem(_swapID bytes32, _secretKey bytes32) returns()
func (_Atom *AtomTransactorSession) Redeem(_swapID [32]byte, _secretKey [32]byte) (*types.Transaction, error) {
	return _Atom.Contract.Redeem(&_Atom.TransactOpts, _swapID, _secretKey)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_Atom *AtomTransactor) Refund(opts *bind.TransactOpts, _swapID [32]byte) (*types.Transaction, error) {
	return _Atom.contract.Transact(opts, "refund", _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_Atom *AtomSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _Atom.Contract.Refund(&_Atom.TransactOpts, _swapID)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(_swapID bytes32) returns()
func (_Atom *AtomTransactorSession) Refund(_swapID [32]byte) (*types.Transaction, error) {
	return _Atom.Contract.Refund(&_Atom.TransactOpts, _swapID)
}

// SubmitDetails is a paid mutator transaction binding the contract method 0x80ee79d4.
//
// Solidity: function submitDetails(_orderID bytes32, _swapDetails bytes) returns()
func (_Atom *AtomTransactor) SubmitDetails(opts *bind.TransactOpts, _orderID [32]byte, _swapDetails []byte) (*types.Transaction, error) {
	return _Atom.contract.Transact(opts, "submitDetails", _orderID, _swapDetails)
}

// SubmitDetails is a paid mutator transaction binding the contract method 0x80ee79d4.
//
// Solidity: function submitDetails(_orderID bytes32, _swapDetails bytes) returns()
func (_Atom *AtomSession) SubmitDetails(_orderID [32]byte, _swapDetails []byte) (*types.Transaction, error) {
	return _Atom.Contract.SubmitDetails(&_Atom.TransactOpts, _orderID, _swapDetails)
}

// SubmitDetails is a paid mutator transaction binding the contract method 0x80ee79d4.
//
// Solidity: function submitDetails(_orderID bytes32, _swapDetails bytes) returns()
func (_Atom *AtomTransactorSession) SubmitDetails(_orderID [32]byte, _swapDetails []byte) (*types.Transaction, error) {
	return _Atom.Contract.SubmitDetails(&_Atom.TransactOpts, _orderID, _swapDetails)
}

// AtomCloseIterator is returned from FilterClose and is used to iterate over the raw logs and unpacked data for Close events raised by the Atom contract.
type AtomCloseIterator struct {
	Event *AtomClose // Event containing the contract specifics and raw log

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
func (it *AtomCloseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomClose)
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
		it.Event = new(AtomClose)
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
func (it *AtomCloseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomCloseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomClose represents a Close event raised by the Atom contract.
type AtomClose struct {
	SwapID    [32]byte
	SecretKey []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterClose is a free log retrieval operation binding the contract event 0x692fd10a275135b9a2a2f5819db3d9965a5129ea2ad3640a0156dbce2fc81bdd.
//
// Solidity: e Close(_swapID bytes32, _secretKey bytes)
func (_Atom *AtomFilterer) FilterClose(opts *bind.FilterOpts) (*AtomCloseIterator, error) {

	logs, sub, err := _Atom.contract.FilterLogs(opts, "Close")
	if err != nil {
		return nil, err
	}
	return &AtomCloseIterator{contract: _Atom.contract, event: "Close", logs: logs, sub: sub}, nil
}

// WatchClose is a free log subscription operation binding the contract event 0x692fd10a275135b9a2a2f5819db3d9965a5129ea2ad3640a0156dbce2fc81bdd.
//
// Solidity: e Close(_swapID bytes32, _secretKey bytes)
func (_Atom *AtomFilterer) WatchClose(opts *bind.WatchOpts, sink chan<- *AtomClose) (event.Subscription, error) {

	logs, sub, err := _Atom.contract.WatchLogs(opts, "Close")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomClose)
				if err := _Atom.contract.UnpackLog(event, "Close", log); err != nil {
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

// AtomExpireIterator is returned from FilterExpire and is used to iterate over the raw logs and unpacked data for Expire events raised by the Atom contract.
type AtomExpireIterator struct {
	Event *AtomExpire // Event containing the contract specifics and raw log

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
func (it *AtomExpireIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomExpire)
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
		it.Event = new(AtomExpire)
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
func (it *AtomExpireIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomExpireIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomExpire represents a Expire event raised by the Atom contract.
type AtomExpire struct {
	SwapID [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterExpire is a free log retrieval operation binding the contract event 0xbddd9b693ea862fad6ecf78fd51c065be26fda94d1f3cad3a7d691453a38a735.
//
// Solidity: e Expire(_swapID bytes32)
func (_Atom *AtomFilterer) FilterExpire(opts *bind.FilterOpts) (*AtomExpireIterator, error) {

	logs, sub, err := _Atom.contract.FilterLogs(opts, "Expire")
	if err != nil {
		return nil, err
	}
	return &AtomExpireIterator{contract: _Atom.contract, event: "Expire", logs: logs, sub: sub}, nil
}

// WatchExpire is a free log subscription operation binding the contract event 0xbddd9b693ea862fad6ecf78fd51c065be26fda94d1f3cad3a7d691453a38a735.
//
// Solidity: e Expire(_swapID bytes32)
func (_Atom *AtomFilterer) WatchExpire(opts *bind.WatchOpts, sink chan<- *AtomExpire) (event.Subscription, error) {

	logs, sub, err := _Atom.contract.WatchLogs(opts, "Expire")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomExpire)
				if err := _Atom.contract.UnpackLog(event, "Expire", log); err != nil {
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

// AtomOpenIterator is returned from FilterOpen and is used to iterate over the raw logs and unpacked data for Open events raised by the Atom contract.
type AtomOpenIterator struct {
	Event *AtomOpen // Event containing the contract specifics and raw log

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
func (it *AtomOpenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomOpen)
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
		it.Event = new(AtomOpen)
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
func (it *AtomOpenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomOpenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomOpen represents a Open event raised by the Atom contract.
type AtomOpen struct {
	SwapID         [32]byte
	WithdrawTrader common.Address
	SecretLock     [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterOpen is a free log retrieval operation binding the contract event 0x6ed79a08bf5c8a7d4a330df315e4ac386627ecafbe5d2bfd6654237d967b24f3.
//
// Solidity: e Open(_swapID bytes32, _withdrawTrader address, _secretLock bytes32)
func (_Atom *AtomFilterer) FilterOpen(opts *bind.FilterOpts) (*AtomOpenIterator, error) {

	logs, sub, err := _Atom.contract.FilterLogs(opts, "Open")
	if err != nil {
		return nil, err
	}
	return &AtomOpenIterator{contract: _Atom.contract, event: "Open", logs: logs, sub: sub}, nil
}

// WatchOpen is a free log subscription operation binding the contract event 0x6ed79a08bf5c8a7d4a330df315e4ac386627ecafbe5d2bfd6654237d967b24f3.
//
// Solidity: e Open(_swapID bytes32, _withdrawTrader address, _secretLock bytes32)
func (_Atom *AtomFilterer) WatchOpen(opts *bind.WatchOpts, sink chan<- *AtomOpen) (event.Subscription, error) {

	logs, sub, err := _Atom.contract.WatchLogs(opts, "Open")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomOpen)
				if err := _Atom.contract.UnpackLog(event, "Open", log); err != nil {
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
