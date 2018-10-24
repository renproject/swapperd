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

// RenExAtomicSwapperABI is the input ABI used to generate the binding from.
const RenExAtomicSwapperABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"initiatable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_withdrawTrader\",\"type\":\"address\"},{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"}],\"name\":\"initiate\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_secretLock\",\"type\":\"bytes32\"},{\"name\":\"_timelock\",\"type\":\"uint256\"}],\"name\":\"swapID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"redeemable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"auditSecret\",\"outputs\":[{\"name\":\"secretKey\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"refundable\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"name\":\"_secretKey\",\"type\":\"bytes32\"}],\"name\":\"redeem\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"redeemedAt\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"audit\",\"outputs\":[{\"name\":\"timelock\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"secretLock\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_VERSION\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_withdrawTrader\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_secretLock\",\"type\":\"bytes32\"}],\"name\":\"LogOpen\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"}],\"name\":\"LogExpire\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"_swapID\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"_secretKey\",\"type\":\"bytes32\"}],\"name\":\"LogClose\",\"type\":\"event\"}]"

// RenExAtomicSwapperBin is the compiled bytecode used for deploying new contracts.
const RenExAtomicSwapperBin = `0x608060405234801561001057600080fd5b50604051610bab380380610bab83398101604052805101805161003a906000906020840190610041565b50506100dc565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061008257805160ff19168380011785556100af565b828001600101855582156100af579182015b828111156100af578251825591602001919060010190610094565b506100bb9291506100bf565b5090565b6100d991905b808211156100bb57600081556001016100c5565b90565b610ac0806100eb6000396000f3006080604052600436106100ae5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166309ece61881146100b3578063412c0b58146100df5780634b2ac3fa146100fe57806368f06b291461012b5780637249fbb614610143578063976d00f41461015b5780639fb3147514610173578063b31597ad1461018b578063bc4fcc4a146101a6578063c140635b146101be578063ffa1ad741461020b575b600080fd5b3480156100bf57600080fd5b506100cb600435610295565b604080519115158252519081900360200190f35b6100fc600435600160a060020a03602435166044356064356102bd565b005b34801561010a57600080fd5b50610119600435602435610446565b60408051918252519081900360200190f35b34801561013757600080fd5b506100cb6004356104cc565b34801561014f57600080fd5b506100fc6004356104d5565b34801561016757600080fd5b5061011960043561064b565b34801561017f57600080fd5b506100cb6004356106d9565b34801561019757600080fd5b506100fc6004356024356106ff565b3480156101b257600080fd5b5061011960043561093b565b3480156101ca57600080fd5b506101d660043561094d565b604080519586526020860194909452600160a060020a0392831685850152911660608401526080830152519081900360a00190f35b34801561021757600080fd5b506102206109d1565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561025a578181015183820152602001610242565b50505050905090810190601f1680156102875780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000805b60008381526002602052604090205460ff1660038111156102b657fe5b1492915050565b6102c5610a5f565b846000808281526002602052604090205460ff1660038111156102e457fe5b14610339576040805160e560020a62461bcd02815260206004820152601660248201527f73776170206f70656e65642070726576696f75736c7900000000000000000000604482015290519081900360640190fd5b6040805160c0810182528481523460208083019182528284018881526000606085018181523360808701908152600160a060020a038d811660a089019081528f855260018088528a86208a51815598518982015595516002808a019190915593516003890155915160048801805491831673ffffffffffffffffffffffffffffffffffffffff199283161790559151600590970180549790911696909116959095179094559290915292902080549194509060ff19168280021790555060408051878152600160a060020a038716602082015280820186905290517f497d46e9505eefe8b910d1a02e6b40d8769510023b0053c3ac4b9574b81c97bf9181900360600190a1505050505050565b6040805160208082018590528183018490528251808303840181526060909201928390528151600093918291908401908083835b602083106104995780518252601f19909201916020918201910161047a565b5181516020939093036101000a600019018019909116921691909117905260405192018290039091209695505050505050565b60006001610299565b80600160008281526002602052604090205460ff1660038111156104f557fe5b1461054a576040805160e560020a62461bcd02815260206004820152600d60248201527f73776170206e6f74206f70656e00000000000000000000000000000000000000604482015290519081900360640190fd5b60008281526001602052604090205482904210156105b2576040805160e560020a62461bcd02815260206004820152601260248201527f73776170206e6f7420657870697261626c650000000000000000000000000000604482015290519081900360640190fd5b6000838152600260209081526040808320805460ff1916600317905560019182905280832060048101549201549051600160a060020a03909216926108fc8215029290818181858888f19350505050158015610612573d6000803e3d6000fd5b506040805184815290517feb711459e1247171382f0da0387b86239b8e3ca60af3b15a9ff2f1eb3d7f6a1d9181900360200190a1505050565b600081600260008281526002602052604090205460ff16600381111561066d57fe5b146106c2576040805160e560020a62461bcd02815260206004820152601160248201527f73776170206e6f742072656465656d6564000000000000000000000000000000604482015290519081900360640190fd5b505060009081526001602052604090206003015490565b60008181526001602052604081205442108015906106f957506001610299565b92915050565b81600160008281526002602052604090205460ff16600381111561071f57fe5b14610774576040805160e560020a62461bcd02815260206004820152600d60248201527f73776170206e6f74206f70656e00000000000000000000000000000000000000604482015290519081900360640190fd5b60408051602080820185905282518083038201815291830192839052815186938693600293909282918401908083835b602083106107c35780518252601f1990920191602091820191016107a4565b51815160209384036101000a600019018019909216911617905260405191909301945091925050808303816000865af1158015610804573d6000803e3d6000fd5b5050506040513d602081101561081957600080fd5b505160008381526001602052604090206002015414610882576040805160e560020a62461bcd02815260206004820152600e60248201527f696e76616c696420736563726574000000000000000000000000000000000000604482015290519081900360640190fd5b600085815260016020818152604080842060038082018a90556002808552838720805460ff19169091179055835281852042905591839052600582015491909201549151600160a060020a03909116926108fc831502929190818181858888f193505050501580156108f8573d6000803e3d6000fd5b50604080518681526020810186905281517f07da1fa25a1d885732677ce9c192cbec27051a4b69d391c9a64850f5a5112ba0929181900390910190a15050505050565b60036020526000908152604090205481565b600080600080600061095d610a5f565b5050506000938452505060016020818152604093849020845160c0810186528154808252938201549281018390526002820154958101869052600382015460608201526004820154600160a060020a03908116608083018190526005909301541660a0909101819052929591949293509190565b6000805460408051602060026001851615610100026000190190941693909304601f81018490048402820184019092528181529291830182828015610a575780601f10610a2c57610100808354040283529160200191610a57565b820191906000526020600020905b815481529060010190602001808311610a3a57829003601f168201915b505050505081565b6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a0810191909152905600a165627a7a72305820ad30ee0c7213b04b85474ac92095e41ef5ae78a1a15db56c8989e720aa81197e0029`

// DeployRenExAtomicSwapper deploys a new Ethereum contract, binding an instance of RenExAtomicSwapper to it.
func DeployRenExAtomicSwapper(auth *bind.TransactOpts, backend bind.ContractBackend, _VERSION string) (common.Address, *types.Transaction, *RenExAtomicSwapper, error) {
	parsed, err := abi.JSON(strings.NewReader(RenExAtomicSwapperABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RenExAtomicSwapperBin), backend, _VERSION)
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

// Initiate is a paid mutator transaction binding the contract method 0x412c0b58.
//
// Solidity: function initiate(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_RenExAtomicSwapper *RenExAtomicSwapperTransactor) Initiate(opts *bind.TransactOpts, _swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _RenExAtomicSwapper.contract.Transact(opts, "initiate", _swapID, _withdrawTrader, _secretLock, _timelock)
}

// Initiate is a paid mutator transaction binding the contract method 0x412c0b58.
//
// Solidity: function initiate(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_RenExAtomicSwapper *RenExAtomicSwapperSession) Initiate(_swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _RenExAtomicSwapper.Contract.Initiate(&_RenExAtomicSwapper.TransactOpts, _swapID, _withdrawTrader, _secretLock, _timelock)
}

// Initiate is a paid mutator transaction binding the contract method 0x412c0b58.
//
// Solidity: function initiate(_swapID bytes32, _withdrawTrader address, _secretLock bytes32, _timelock uint256) returns()
func (_RenExAtomicSwapper *RenExAtomicSwapperTransactorSession) Initiate(_swapID [32]byte, _withdrawTrader common.Address, _secretLock [32]byte, _timelock *big.Int) (*types.Transaction, error) {
	return _RenExAtomicSwapper.Contract.Initiate(&_RenExAtomicSwapper.TransactOpts, _swapID, _withdrawTrader, _secretLock, _timelock)
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
