// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// VeAWPHelperMetaData contains all meta data concerning the VeAWPHelper contract.
var VeAWPHelperMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"awpToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"veAWP_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"awpToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"depositFor\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lockDuration\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"veAWP\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"error\",\"name\":\"ApproveFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidUser\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
}

// VeAWPHelperABI is the input ABI used to generate the binding from.
// Deprecated: Use VeAWPHelperMetaData.ABI instead.
var VeAWPHelperABI = VeAWPHelperMetaData.ABI

// VeAWPHelper is an auto generated Go binding around an Ethereum contract.
type VeAWPHelper struct {
	VeAWPHelperCaller     // Read-only binding to the contract
	VeAWPHelperTransactor // Write-only binding to the contract
	VeAWPHelperFilterer   // Log filterer for contract events
}

// VeAWPHelperCaller is an auto generated read-only Go binding around an Ethereum contract.
type VeAWPHelperCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VeAWPHelperTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VeAWPHelperTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VeAWPHelperFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VeAWPHelperFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VeAWPHelperSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VeAWPHelperSession struct {
	Contract     *VeAWPHelper      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VeAWPHelperCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VeAWPHelperCallerSession struct {
	Contract *VeAWPHelperCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// VeAWPHelperTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VeAWPHelperTransactorSession struct {
	Contract     *VeAWPHelperTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// VeAWPHelperRaw is an auto generated low-level Go binding around an Ethereum contract.
type VeAWPHelperRaw struct {
	Contract *VeAWPHelper // Generic contract binding to access the raw methods on
}

// VeAWPHelperCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VeAWPHelperCallerRaw struct {
	Contract *VeAWPHelperCaller // Generic read-only contract binding to access the raw methods on
}

// VeAWPHelperTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VeAWPHelperTransactorRaw struct {
	Contract *VeAWPHelperTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVeAWPHelper creates a new instance of VeAWPHelper, bound to a specific deployed contract.
func NewVeAWPHelper(address common.Address, backend bind.ContractBackend) (*VeAWPHelper, error) {
	contract, err := bindVeAWPHelper(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VeAWPHelper{VeAWPHelperCaller: VeAWPHelperCaller{contract: contract}, VeAWPHelperTransactor: VeAWPHelperTransactor{contract: contract}, VeAWPHelperFilterer: VeAWPHelperFilterer{contract: contract}}, nil
}

// NewVeAWPHelperCaller creates a new read-only instance of VeAWPHelper, bound to a specific deployed contract.
func NewVeAWPHelperCaller(address common.Address, caller bind.ContractCaller) (*VeAWPHelperCaller, error) {
	contract, err := bindVeAWPHelper(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VeAWPHelperCaller{contract: contract}, nil
}

// NewVeAWPHelperTransactor creates a new write-only instance of VeAWPHelper, bound to a specific deployed contract.
func NewVeAWPHelperTransactor(address common.Address, transactor bind.ContractTransactor) (*VeAWPHelperTransactor, error) {
	contract, err := bindVeAWPHelper(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VeAWPHelperTransactor{contract: contract}, nil
}

// NewVeAWPHelperFilterer creates a new log filterer instance of VeAWPHelper, bound to a specific deployed contract.
func NewVeAWPHelperFilterer(address common.Address, filterer bind.ContractFilterer) (*VeAWPHelperFilterer, error) {
	contract, err := bindVeAWPHelper(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VeAWPHelperFilterer{contract: contract}, nil
}

// bindVeAWPHelper binds a generic wrapper to an already deployed contract.
func bindVeAWPHelper(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VeAWPHelperMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VeAWPHelper *VeAWPHelperRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VeAWPHelper.Contract.VeAWPHelperCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VeAWPHelper *VeAWPHelperRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VeAWPHelper.Contract.VeAWPHelperTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VeAWPHelper *VeAWPHelperRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VeAWPHelper.Contract.VeAWPHelperTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VeAWPHelper *VeAWPHelperCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VeAWPHelper.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VeAWPHelper *VeAWPHelperTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VeAWPHelper.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VeAWPHelper *VeAWPHelperTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VeAWPHelper.Contract.contract.Transact(opts, method, params...)
}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_VeAWPHelper *VeAWPHelperCaller) AwpToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VeAWPHelper.contract.Call(opts, &out, "awpToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_VeAWPHelper *VeAWPHelperSession) AwpToken() (common.Address, error) {
	return _VeAWPHelper.Contract.AwpToken(&_VeAWPHelper.CallOpts)
}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_VeAWPHelper *VeAWPHelperCallerSession) AwpToken() (common.Address, error) {
	return _VeAWPHelper.Contract.AwpToken(&_VeAWPHelper.CallOpts)
}

// VeAWP is a free data retrieval call binding the contract method 0x7bb8431f.
//
// Solidity: function veAWP() view returns(address)
func (_VeAWPHelper *VeAWPHelperCaller) VeAWP(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VeAWPHelper.contract.Call(opts, &out, "veAWP")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VeAWP is a free data retrieval call binding the contract method 0x7bb8431f.
//
// Solidity: function veAWP() view returns(address)
func (_VeAWPHelper *VeAWPHelperSession) VeAWP() (common.Address, error) {
	return _VeAWPHelper.Contract.VeAWP(&_VeAWPHelper.CallOpts)
}

// VeAWP is a free data retrieval call binding the contract method 0x7bb8431f.
//
// Solidity: function veAWP() view returns(address)
func (_VeAWPHelper *VeAWPHelperCallerSession) VeAWP() (common.Address, error) {
	return _VeAWPHelper.Contract.VeAWP(&_VeAWPHelper.CallOpts)
}

// DepositFor is a paid mutator transaction binding the contract method 0x59133324.
//
// Solidity: function depositFor(address user, uint256 amount, uint64 lockDuration, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256 tokenId)
func (_VeAWPHelper *VeAWPHelperTransactor) DepositFor(opts *bind.TransactOpts, user common.Address, amount *big.Int, lockDuration uint64, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _VeAWPHelper.contract.Transact(opts, "depositFor", user, amount, lockDuration, deadline, v, r, s)
}

// DepositFor is a paid mutator transaction binding the contract method 0x59133324.
//
// Solidity: function depositFor(address user, uint256 amount, uint64 lockDuration, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256 tokenId)
func (_VeAWPHelper *VeAWPHelperSession) DepositFor(user common.Address, amount *big.Int, lockDuration uint64, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _VeAWPHelper.Contract.DepositFor(&_VeAWPHelper.TransactOpts, user, amount, lockDuration, deadline, v, r, s)
}

// DepositFor is a paid mutator transaction binding the contract method 0x59133324.
//
// Solidity: function depositFor(address user, uint256 amount, uint64 lockDuration, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256 tokenId)
func (_VeAWPHelper *VeAWPHelperTransactorSession) DepositFor(user common.Address, amount *big.Int, lockDuration uint64, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _VeAWPHelper.Contract.DepositFor(&_VeAWPHelper.TransactOpts, user, amount, lockDuration, deadline, v, r, s)
}
