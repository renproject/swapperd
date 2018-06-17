// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package eth

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// AtomWalletABI is the input ABI used to generate the binding from.
const AtomWalletABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_id\",\"type\":\"bytes32\"},{\"name\":\"_orderType\",\"type\":\"uint8\"},{\"name\":\"_parity\",\"type\":\"uint8\"},{\"name\":\"_expiry\",\"type\":\"uint64\"},{\"name\":\"_tokens\",\"type\":\"uint64\"},{\"name\":\"_priceC\",\"type\":\"uint16\"},{\"name\":\"_priceQ\",\"type\":\"uint16\"},{\"name\":\"_volumeC\",\"type\":\"uint16\"},{\"name\":\"_volumeQ\",\"type\":\"uint16\"},{\"name\":\"_minimumVolumeC\",\"type\":\"uint16\"},{\"name\":\"_minimumVolumeQ\",\"type\":\"uint16\"},{\"name\":\"_nonceHash\",\"type\":\"uint256\"}],\"name\":\"submitOrder\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_buy\",\"type\":\"bytes32\"},{\"name\":\"_sell\",\"type\":\"bytes32\"}],\"name\":\"submitMatch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"orders\",\"outputs\":[{\"name\":\"parity\",\"type\":\"uint8\"},{\"name\":\"orderType\",\"type\":\"uint8\"},{\"name\":\"expiry\",\"type\":\"uint64\"},{\"name\":\"tokens\",\"type\":\"uint64\"},{\"name\":\"priceC\",\"type\":\"uint256\"},{\"name\":\"priceQ\",\"type\":\"uint256\"},{\"name\":\"volumeC\",\"type\":\"uint256\"},{\"name\":\"volumeQ\",\"type\":\"uint256\"},{\"name\":\"minimumVolumeC\",\"type\":\"uint256\"},{\"name\":\"minimumVolumeQ\",\"type\":\"uint256\"},{\"name\":\"nonceHash\",\"type\":\"uint256\"},{\"name\":\"trader\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"matches\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"bonds\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"}]"

// AtomWalletBin is the compiled bytecode used for deploying new contracts.
const AtomWalletBin = `0x608060405234801561001057600080fd5b50610545806100206000396000f3006080604052600436106100775763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663177d19c3811461008e5780632a337d30146100f35780632e1a7d4d1461010e5780639c3f1e90146101265780639fe9ada3146101c6578063fe10d774146101f0575b336000908152600260205260409020805434019055005b34801561009a57600080fd5b506100f160043560ff6024358116906044351667ffffffffffffffff6064358116906084351661ffff60a43581169060c43581169060e435811690610104358116906101243581169061014435166101643561021e565b005b3480156100ff57600080fd5b506100f1600435602435610416565b34801561011a57600080fd5b506100f160043561042f565b34801561013257600080fd5b5061013e600435610464565b6040805160ff9d8e1681529b909c1660208c015267ffffffffffffffff998a168b8d01529790981660608a0152608089019590955260a088019390935260c087019190915260e086015261010085015261012084015261014083019190915273ffffffffffffffffffffffffffffffffffffffff166101608201529051908190036101800190f35b3480156101d257600080fd5b506101de6004356104f5565b60408051918252519081900360200190f35b3480156101fc57600080fd5b506101de73ffffffffffffffffffffffffffffffffffffffff60043516610507565b610180604051908101604052808b60ff1681526020018c60ff1681526020018a67ffffffffffffffff1681526020018967ffffffffffffffff1681526020018861ffff1681526020018761ffff1681526020018661ffff1681526020018561ffff1681526020018461ffff1681526020018361ffff168152602001828152602001600073ffffffffffffffffffffffffffffffffffffffff168152506000808e6000191660001916815260200190815260200160002060008201518160000160006101000a81548160ff021916908360ff16021790555060208201518160000160016101000a81548160ff021916908360ff16021790555060408201518160000160026101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550606082015181600001600a6101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506080820151816001015560a0820151816002015560c0820151816003015560e082015181600401556101008201518160050155610120820151816006015561014082015181600701556101608201518160080160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550905050505050505050505050505050565b6000828152600160205260408082208390559181522055565b3360009081526002602052604090205481111561044b57600080fd5b3360009081526002602052604090208054919091039055565b60006020819052908152604090208054600182015460028301546003840154600485015460058601546006870154600788015460089098015460ff808916996101008a049091169867ffffffffffffffff6201000082048116996a010000000000000000000090920416979096909590949093909290919073ffffffffffffffffffffffffffffffffffffffff168c565b60016020526000908152604090205481565b600260205260009081526040902054815600a165627a7a72305820acdafd20a9a9fb97b6d823efefa00b06a540efd2543e62eeb6d6d5444324b4ef0029`

// DeployAtomWallet deploys a new Ethereum contract, binding an instance of AtomWallet to it.
func DeployAtomWallet(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AtomWallet, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomWalletABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AtomWalletBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AtomWallet{AtomWalletCaller: AtomWalletCaller{contract: contract}, AtomWalletTransactor: AtomWalletTransactor{contract: contract}, AtomWalletFilterer: AtomWalletFilterer{contract: contract}}, nil
}

// AtomWallet is an auto generated Go binding around an Ethereum contract.
type AtomWallet struct {
	AtomWalletCaller     // Read-only binding to the contract
	AtomWalletTransactor // Write-only binding to the contract
	AtomWalletFilterer   // Log filterer for contract events
}

// AtomWalletCaller is an auto generated read-only Go binding around an Ethereum contract.
type AtomWalletCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomWalletTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AtomWalletTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomWalletFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AtomWalletFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomWalletSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AtomWalletSession struct {
	Contract     *AtomWallet       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AtomWalletCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AtomWalletCallerSession struct {
	Contract *AtomWalletCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// AtomWalletTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AtomWalletTransactorSession struct {
	Contract     *AtomWalletTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// AtomWalletRaw is an auto generated low-level Go binding around an Ethereum contract.
type AtomWalletRaw struct {
	Contract *AtomWallet // Generic contract binding to access the raw methods on
}

// AtomWalletCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AtomWalletCallerRaw struct {
	Contract *AtomWalletCaller // Generic read-only contract binding to access the raw methods on
}

// AtomWalletTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AtomWalletTransactorRaw struct {
	Contract *AtomWalletTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAtomWallet creates a new instance of AtomWallet, bound to a specific deployed contract.
func NewAtomWallet(address common.Address, backend bind.ContractBackend) (*AtomWallet, error) {
	contract, err := bindAtomWallet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AtomWallet{AtomWalletCaller: AtomWalletCaller{contract: contract}, AtomWalletTransactor: AtomWalletTransactor{contract: contract}, AtomWalletFilterer: AtomWalletFilterer{contract: contract}}, nil
}

// NewAtomWalletCaller creates a new read-only instance of AtomWallet, bound to a specific deployed contract.
func NewAtomWalletCaller(address common.Address, caller bind.ContractCaller) (*AtomWalletCaller, error) {
	contract, err := bindAtomWallet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AtomWalletCaller{contract: contract}, nil
}

// NewAtomWalletTransactor creates a new write-only instance of AtomWallet, bound to a specific deployed contract.
func NewAtomWalletTransactor(address common.Address, transactor bind.ContractTransactor) (*AtomWalletTransactor, error) {
	contract, err := bindAtomWallet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AtomWalletTransactor{contract: contract}, nil
}

// NewAtomWalletFilterer creates a new log filterer instance of AtomWallet, bound to a specific deployed contract.
func NewAtomWalletFilterer(address common.Address, filterer bind.ContractFilterer) (*AtomWalletFilterer, error) {
	contract, err := bindAtomWallet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AtomWalletFilterer{contract: contract}, nil
}

// bindAtomWallet binds a generic wrapper to an already deployed contract.
func bindAtomWallet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AtomWalletABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomWallet *AtomWalletRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomWallet.Contract.AtomWalletCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomWallet *AtomWalletRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomWallet.Contract.AtomWalletTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomWallet *AtomWalletRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomWallet.Contract.AtomWalletTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomWallet *AtomWalletCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AtomWallet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomWallet *AtomWalletTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomWallet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomWallet *AtomWalletTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomWallet.Contract.contract.Transact(opts, method, params...)
}

// Bonds is a free data retrieval call binding the contract method 0xfe10d774.
//
// Solidity: function bonds( address) constant returns(uint256)
func (_AtomWallet *AtomWalletCaller) Bonds(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AtomWallet.contract.Call(opts, out, "bonds", arg0)
	return *ret0, err
}

// Bonds is a free data retrieval call binding the contract method 0xfe10d774.
//
// Solidity: function bonds( address) constant returns(uint256)
func (_AtomWallet *AtomWalletSession) Bonds(arg0 common.Address) (*big.Int, error) {
	return _AtomWallet.Contract.Bonds(&_AtomWallet.CallOpts, arg0)
}

// Bonds is a free data retrieval call binding the contract method 0xfe10d774.
//
// Solidity: function bonds( address) constant returns(uint256)
func (_AtomWallet *AtomWalletCallerSession) Bonds(arg0 common.Address) (*big.Int, error) {
	return _AtomWallet.Contract.Bonds(&_AtomWallet.CallOpts, arg0)
}

// Matches is a free data retrieval call binding the contract method 0x9fe9ada3.
//
// Solidity: function matches( bytes32) constant returns(bytes32)
func (_AtomWallet *AtomWalletCaller) Matches(opts *bind.CallOpts, arg0 [32]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _AtomWallet.contract.Call(opts, out, "matches", arg0)
	return *ret0, err
}

// Matches is a free data retrieval call binding the contract method 0x9fe9ada3.
//
// Solidity: function matches( bytes32) constant returns(bytes32)
func (_AtomWallet *AtomWalletSession) Matches(arg0 [32]byte) ([32]byte, error) {
	return _AtomWallet.Contract.Matches(&_AtomWallet.CallOpts, arg0)
}

// Matches is a free data retrieval call binding the contract method 0x9fe9ada3.
//
// Solidity: function matches( bytes32) constant returns(bytes32)
func (_AtomWallet *AtomWalletCallerSession) Matches(arg0 [32]byte) ([32]byte, error) {
	return _AtomWallet.Contract.Matches(&_AtomWallet.CallOpts, arg0)
}

// Orders is a free data retrieval call binding the contract method 0x9c3f1e90.
//
// Solidity: function orders( bytes32) constant returns(parity uint8, orderType uint8, expiry uint64, tokens uint64, priceC uint256, priceQ uint256, volumeC uint256, volumeQ uint256, minimumVolumeC uint256, minimumVolumeQ uint256, nonceHash uint256, trader address)
func (_AtomWallet *AtomWalletCaller) Orders(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Parity         uint8
	OrderType      uint8
	Expiry         uint64
	Tokens         uint64
	PriceC         *big.Int
	PriceQ         *big.Int
	VolumeC        *big.Int
	VolumeQ        *big.Int
	MinimumVolumeC *big.Int
	MinimumVolumeQ *big.Int
	NonceHash      *big.Int
	Trader         common.Address
}, error) {
	ret := new(struct {
		Parity         uint8
		OrderType      uint8
		Expiry         uint64
		Tokens         uint64
		PriceC         *big.Int
		PriceQ         *big.Int
		VolumeC        *big.Int
		VolumeQ        *big.Int
		MinimumVolumeC *big.Int
		MinimumVolumeQ *big.Int
		NonceHash      *big.Int
		Trader         common.Address
	})
	out := ret
	err := _AtomWallet.contract.Call(opts, out, "orders", arg0)
	return *ret, err
}

// Orders is a free data retrieval call binding the contract method 0x9c3f1e90.
//
// Solidity: function orders( bytes32) constant returns(parity uint8, orderType uint8, expiry uint64, tokens uint64, priceC uint256, priceQ uint256, volumeC uint256, volumeQ uint256, minimumVolumeC uint256, minimumVolumeQ uint256, nonceHash uint256, trader address)
func (_AtomWallet *AtomWalletSession) Orders(arg0 [32]byte) (struct {
	Parity         uint8
	OrderType      uint8
	Expiry         uint64
	Tokens         uint64
	PriceC         *big.Int
	PriceQ         *big.Int
	VolumeC        *big.Int
	VolumeQ        *big.Int
	MinimumVolumeC *big.Int
	MinimumVolumeQ *big.Int
	NonceHash      *big.Int
	Trader         common.Address
}, error) {
	return _AtomWallet.Contract.Orders(&_AtomWallet.CallOpts, arg0)
}

// Orders is a free data retrieval call binding the contract method 0x9c3f1e90.
//
// Solidity: function orders( bytes32) constant returns(parity uint8, orderType uint8, expiry uint64, tokens uint64, priceC uint256, priceQ uint256, volumeC uint256, volumeQ uint256, minimumVolumeC uint256, minimumVolumeQ uint256, nonceHash uint256, trader address)
func (_AtomWallet *AtomWalletCallerSession) Orders(arg0 [32]byte) (struct {
	Parity         uint8
	OrderType      uint8
	Expiry         uint64
	Tokens         uint64
	PriceC         *big.Int
	PriceQ         *big.Int
	VolumeC        *big.Int
	VolumeQ        *big.Int
	MinimumVolumeC *big.Int
	MinimumVolumeQ *big.Int
	NonceHash      *big.Int
	Trader         common.Address
}, error) {
	return _AtomWallet.Contract.Orders(&_AtomWallet.CallOpts, arg0)
}

// SubmitMatch is a paid mutator transaction binding the contract method 0x2a337d30.
//
// Solidity: function submitMatch(_buy bytes32, _sell bytes32) returns()
func (_AtomWallet *AtomWalletTransactor) SubmitMatch(opts *bind.TransactOpts, _buy [32]byte, _sell [32]byte) (*types.Transaction, error) {
	return _AtomWallet.contract.Transact(opts, "submitMatch", _buy, _sell)
}

// SubmitMatch is a paid mutator transaction binding the contract method 0x2a337d30.
//
// Solidity: function submitMatch(_buy bytes32, _sell bytes32) returns()
func (_AtomWallet *AtomWalletSession) SubmitMatch(_buy [32]byte, _sell [32]byte) (*types.Transaction, error) {
	return _AtomWallet.Contract.SubmitMatch(&_AtomWallet.TransactOpts, _buy, _sell)
}

// SubmitMatch is a paid mutator transaction binding the contract method 0x2a337d30.
//
// Solidity: function submitMatch(_buy bytes32, _sell bytes32) returns()
func (_AtomWallet *AtomWalletTransactorSession) SubmitMatch(_buy [32]byte, _sell [32]byte) (*types.Transaction, error) {
	return _AtomWallet.Contract.SubmitMatch(&_AtomWallet.TransactOpts, _buy, _sell)
}

// SubmitOrder is a paid mutator transaction binding the contract method 0x177d19c3.
//
// Solidity: function submitOrder(_id bytes32, _orderType uint8, _parity uint8, _expiry uint64, _tokens uint64, _priceC uint16, _priceQ uint16, _volumeC uint16, _volumeQ uint16, _minimumVolumeC uint16, _minimumVolumeQ uint16, _nonceHash uint256) returns()
func (_AtomWallet *AtomWalletTransactor) SubmitOrder(opts *bind.TransactOpts, _id [32]byte, _orderType uint8, _parity uint8, _expiry uint64, _tokens uint64, _priceC uint16, _priceQ uint16, _volumeC uint16, _volumeQ uint16, _minimumVolumeC uint16, _minimumVolumeQ uint16, _nonceHash *big.Int) (*types.Transaction, error) {
	return _AtomWallet.contract.Transact(opts, "submitOrder", _id, _orderType, _parity, _expiry, _tokens, _priceC, _priceQ, _volumeC, _volumeQ, _minimumVolumeC, _minimumVolumeQ, _nonceHash)
}

// SubmitOrder is a paid mutator transaction binding the contract method 0x177d19c3.
//
// Solidity: function submitOrder(_id bytes32, _orderType uint8, _parity uint8, _expiry uint64, _tokens uint64, _priceC uint16, _priceQ uint16, _volumeC uint16, _volumeQ uint16, _minimumVolumeC uint16, _minimumVolumeQ uint16, _nonceHash uint256) returns()
func (_AtomWallet *AtomWalletSession) SubmitOrder(_id [32]byte, _orderType uint8, _parity uint8, _expiry uint64, _tokens uint64, _priceC uint16, _priceQ uint16, _volumeC uint16, _volumeQ uint16, _minimumVolumeC uint16, _minimumVolumeQ uint16, _nonceHash *big.Int) (*types.Transaction, error) {
	return _AtomWallet.Contract.SubmitOrder(&_AtomWallet.TransactOpts, _id, _orderType, _parity, _expiry, _tokens, _priceC, _priceQ, _volumeC, _volumeQ, _minimumVolumeC, _minimumVolumeQ, _nonceHash)
}

// SubmitOrder is a paid mutator transaction binding the contract method 0x177d19c3.
//
// Solidity: function submitOrder(_id bytes32, _orderType uint8, _parity uint8, _expiry uint64, _tokens uint64, _priceC uint16, _priceQ uint16, _volumeC uint16, _volumeQ uint16, _minimumVolumeC uint16, _minimumVolumeQ uint16, _nonceHash uint256) returns()
func (_AtomWallet *AtomWalletTransactorSession) SubmitOrder(_id [32]byte, _orderType uint8, _parity uint8, _expiry uint64, _tokens uint64, _priceC uint16, _priceQ uint16, _volumeC uint16, _volumeQ uint16, _minimumVolumeC uint16, _minimumVolumeQ uint16, _nonceHash *big.Int) (*types.Transaction, error) {
	return _AtomWallet.Contract.SubmitOrder(&_AtomWallet.TransactOpts, _id, _orderType, _parity, _expiry, _tokens, _priceC, _priceQ, _volumeC, _volumeQ, _minimumVolumeC, _minimumVolumeQ, _nonceHash)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(_amount uint256) returns()
func (_AtomWallet *AtomWalletTransactor) Withdraw(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _AtomWallet.contract.Transact(opts, "withdraw", _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(_amount uint256) returns()
func (_AtomWallet *AtomWalletSession) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _AtomWallet.Contract.Withdraw(&_AtomWallet.TransactOpts, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(_amount uint256) returns()
func (_AtomWallet *AtomWalletTransactorSession) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _AtomWallet.Contract.Withdraw(&_AtomWallet.TransactOpts, _amount)
}
