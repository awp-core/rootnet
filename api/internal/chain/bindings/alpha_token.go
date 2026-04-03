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

// AlphaTokenMetaData contains all meta data concerning the AlphaToken contract.
var AlphaTokenMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MAX_SUPPLY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"admin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approveAndCall\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"burnFrom\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createdAt\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"currentMintableLimit\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grossMintedSinceLock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"worknetId_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"admin_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"minterPaused\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minters\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mintersLocked\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setMinterPaused\",\"inputs\":[{\"name\":\"minter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"paused\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWorknetMinter\",\"inputs\":[{\"name\":\"worknetManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supplyAtLock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferAndCall\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"worknetId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WorknetMinterSet\",\"inputs\":[{\"name\":\"worknetManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ExceedsMaxSupply\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExceedsMintableLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidCallback\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MinterPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MintersLocked\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotMinter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
}

// AlphaTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use AlphaTokenMetaData.ABI instead.
var AlphaTokenABI = AlphaTokenMetaData.ABI

// AlphaToken is an auto generated Go binding around an Ethereum contract.
type AlphaToken struct {
	AlphaTokenCaller     // Read-only binding to the contract
	AlphaTokenTransactor // Write-only binding to the contract
	AlphaTokenFilterer   // Log filterer for contract events
}

// AlphaTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type AlphaTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AlphaTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AlphaTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AlphaTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AlphaTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AlphaTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AlphaTokenSession struct {
	Contract     *AlphaToken       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AlphaTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AlphaTokenCallerSession struct {
	Contract *AlphaTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// AlphaTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AlphaTokenTransactorSession struct {
	Contract     *AlphaTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// AlphaTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type AlphaTokenRaw struct {
	Contract *AlphaToken // Generic contract binding to access the raw methods on
}

// AlphaTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AlphaTokenCallerRaw struct {
	Contract *AlphaTokenCaller // Generic read-only contract binding to access the raw methods on
}

// AlphaTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AlphaTokenTransactorRaw struct {
	Contract *AlphaTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAlphaToken creates a new instance of AlphaToken, bound to a specific deployed contract.
func NewAlphaToken(address common.Address, backend bind.ContractBackend) (*AlphaToken, error) {
	contract, err := bindAlphaToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AlphaToken{AlphaTokenCaller: AlphaTokenCaller{contract: contract}, AlphaTokenTransactor: AlphaTokenTransactor{contract: contract}, AlphaTokenFilterer: AlphaTokenFilterer{contract: contract}}, nil
}

// NewAlphaTokenCaller creates a new read-only instance of AlphaToken, bound to a specific deployed contract.
func NewAlphaTokenCaller(address common.Address, caller bind.ContractCaller) (*AlphaTokenCaller, error) {
	contract, err := bindAlphaToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AlphaTokenCaller{contract: contract}, nil
}

// NewAlphaTokenTransactor creates a new write-only instance of AlphaToken, bound to a specific deployed contract.
func NewAlphaTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*AlphaTokenTransactor, error) {
	contract, err := bindAlphaToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AlphaTokenTransactor{contract: contract}, nil
}

// NewAlphaTokenFilterer creates a new log filterer instance of AlphaToken, bound to a specific deployed contract.
func NewAlphaTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*AlphaTokenFilterer, error) {
	contract, err := bindAlphaToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AlphaTokenFilterer{contract: contract}, nil
}

// bindAlphaToken binds a generic wrapper to an already deployed contract.
func bindAlphaToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AlphaTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AlphaToken *AlphaTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AlphaToken.Contract.AlphaTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AlphaToken *AlphaTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AlphaToken.Contract.AlphaTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AlphaToken *AlphaTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AlphaToken.Contract.AlphaTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AlphaToken *AlphaTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AlphaToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AlphaToken *AlphaTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AlphaToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AlphaToken *AlphaTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AlphaToken.Contract.contract.Transact(opts, method, params...)
}

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_AlphaToken *AlphaTokenCaller) MAXSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AlphaToken.contract.Call(opts, &out, "MAX_SUPPLY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_AlphaToken *AlphaTokenSession) MAXSUPPLY() (*big.Int, error) {
	return _AlphaToken.Contract.MAXSUPPLY(&_AlphaToken.CallOpts)
}

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_AlphaToken *AlphaTokenCallerSession) MAXSUPPLY() (*big.Int, error) {
	return _AlphaToken.Contract.MAXSUPPLY(&_AlphaToken.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_AlphaToken *AlphaTokenCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AlphaToken.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_AlphaToken *AlphaTokenSession) Admin() (common.Address, error) {
	return _AlphaToken.Contract.Admin(&_AlphaToken.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_AlphaToken *AlphaTokenCallerSession) Admin() (common.Address, error) {
	return _AlphaToken.Contract.Admin(&_AlphaToken.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_AlphaToken *AlphaTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AlphaToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_AlphaToken *AlphaTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _AlphaToken.Contract.Allowance(&_AlphaToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_AlphaToken *AlphaTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _AlphaToken.Contract.Allowance(&_AlphaToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_AlphaToken *AlphaTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AlphaToken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_AlphaToken *AlphaTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _AlphaToken.Contract.BalanceOf(&_AlphaToken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_AlphaToken *AlphaTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _AlphaToken.Contract.BalanceOf(&_AlphaToken.CallOpts, account)
}

// CreatedAt is a free data retrieval call binding the contract method 0xcf09e0d0.
//
// Solidity: function createdAt() view returns(uint64)
func (_AlphaToken *AlphaTokenCaller) CreatedAt(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _AlphaToken.contract.Call(opts, &out, "createdAt")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// CreatedAt is a free data retrieval call binding the contract method 0xcf09e0d0.
//
// Solidity: function createdAt() view returns(uint64)
func (_AlphaToken *AlphaTokenSession) CreatedAt() (uint64, error) {
	return _AlphaToken.Contract.CreatedAt(&_AlphaToken.CallOpts)
}

// CreatedAt is a free data retrieval call binding the contract method 0xcf09e0d0.
//
// Solidity: function createdAt() view returns(uint64)
func (_AlphaToken *AlphaTokenCallerSession) CreatedAt() (uint64, error) {
	return _AlphaToken.Contract.CreatedAt(&_AlphaToken.CallOpts)
}

// CurrentMintableLimit is a free data retrieval call binding the contract method 0xfd3088ae.
//
// Solidity: function currentMintableLimit() view returns(uint256)
func (_AlphaToken *AlphaTokenCaller) CurrentMintableLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AlphaToken.contract.Call(opts, &out, "currentMintableLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentMintableLimit is a free data retrieval call binding the contract method 0xfd3088ae.
//
// Solidity: function currentMintableLimit() view returns(uint256)
func (_AlphaToken *AlphaTokenSession) CurrentMintableLimit() (*big.Int, error) {
	return _AlphaToken.Contract.CurrentMintableLimit(&_AlphaToken.CallOpts)
}

// CurrentMintableLimit is a free data retrieval call binding the contract method 0xfd3088ae.
//
// Solidity: function currentMintableLimit() view returns(uint256)
func (_AlphaToken *AlphaTokenCallerSession) CurrentMintableLimit() (*big.Int, error) {
	return _AlphaToken.Contract.CurrentMintableLimit(&_AlphaToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AlphaToken *AlphaTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AlphaToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AlphaToken *AlphaTokenSession) Decimals() (uint8, error) {
	return _AlphaToken.Contract.Decimals(&_AlphaToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AlphaToken *AlphaTokenCallerSession) Decimals() (uint8, error) {
	return _AlphaToken.Contract.Decimals(&_AlphaToken.CallOpts)
}

// GrossMintedSinceLock is a free data retrieval call binding the contract method 0x70afc24b.
//
// Solidity: function grossMintedSinceLock() view returns(uint256)
func (_AlphaToken *AlphaTokenCaller) GrossMintedSinceLock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AlphaToken.contract.Call(opts, &out, "grossMintedSinceLock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GrossMintedSinceLock is a free data retrieval call binding the contract method 0x70afc24b.
//
// Solidity: function grossMintedSinceLock() view returns(uint256)
func (_AlphaToken *AlphaTokenSession) GrossMintedSinceLock() (*big.Int, error) {
	return _AlphaToken.Contract.GrossMintedSinceLock(&_AlphaToken.CallOpts)
}

// GrossMintedSinceLock is a free data retrieval call binding the contract method 0x70afc24b.
//
// Solidity: function grossMintedSinceLock() view returns(uint256)
func (_AlphaToken *AlphaTokenCallerSession) GrossMintedSinceLock() (*big.Int, error) {
	return _AlphaToken.Contract.GrossMintedSinceLock(&_AlphaToken.CallOpts)
}

// MinterPaused is a free data retrieval call binding the contract method 0xdb44d6da.
//
// Solidity: function minterPaused(address ) view returns(bool)
func (_AlphaToken *AlphaTokenCaller) MinterPaused(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _AlphaToken.contract.Call(opts, &out, "minterPaused", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// MinterPaused is a free data retrieval call binding the contract method 0xdb44d6da.
//
// Solidity: function minterPaused(address ) view returns(bool)
func (_AlphaToken *AlphaTokenSession) MinterPaused(arg0 common.Address) (bool, error) {
	return _AlphaToken.Contract.MinterPaused(&_AlphaToken.CallOpts, arg0)
}

// MinterPaused is a free data retrieval call binding the contract method 0xdb44d6da.
//
// Solidity: function minterPaused(address ) view returns(bool)
func (_AlphaToken *AlphaTokenCallerSession) MinterPaused(arg0 common.Address) (bool, error) {
	return _AlphaToken.Contract.MinterPaused(&_AlphaToken.CallOpts, arg0)
}

// Minters is a free data retrieval call binding the contract method 0xf46eccc4.
//
// Solidity: function minters(address ) view returns(bool)
func (_AlphaToken *AlphaTokenCaller) Minters(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _AlphaToken.contract.Call(opts, &out, "minters", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Minters is a free data retrieval call binding the contract method 0xf46eccc4.
//
// Solidity: function minters(address ) view returns(bool)
func (_AlphaToken *AlphaTokenSession) Minters(arg0 common.Address) (bool, error) {
	return _AlphaToken.Contract.Minters(&_AlphaToken.CallOpts, arg0)
}

// Minters is a free data retrieval call binding the contract method 0xf46eccc4.
//
// Solidity: function minters(address ) view returns(bool)
func (_AlphaToken *AlphaTokenCallerSession) Minters(arg0 common.Address) (bool, error) {
	return _AlphaToken.Contract.Minters(&_AlphaToken.CallOpts, arg0)
}

// MintersLocked is a free data retrieval call binding the contract method 0x7fe290f7.
//
// Solidity: function mintersLocked() view returns(bool)
func (_AlphaToken *AlphaTokenCaller) MintersLocked(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AlphaToken.contract.Call(opts, &out, "mintersLocked")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// MintersLocked is a free data retrieval call binding the contract method 0x7fe290f7.
//
// Solidity: function mintersLocked() view returns(bool)
func (_AlphaToken *AlphaTokenSession) MintersLocked() (bool, error) {
	return _AlphaToken.Contract.MintersLocked(&_AlphaToken.CallOpts)
}

// MintersLocked is a free data retrieval call binding the contract method 0x7fe290f7.
//
// Solidity: function mintersLocked() view returns(bool)
func (_AlphaToken *AlphaTokenCallerSession) MintersLocked() (bool, error) {
	return _AlphaToken.Contract.MintersLocked(&_AlphaToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AlphaToken *AlphaTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AlphaToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AlphaToken *AlphaTokenSession) Name() (string, error) {
	return _AlphaToken.Contract.Name(&_AlphaToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AlphaToken *AlphaTokenCallerSession) Name() (string, error) {
	return _AlphaToken.Contract.Name(&_AlphaToken.CallOpts)
}

// SupplyAtLock is a free data retrieval call binding the contract method 0xea32345b.
//
// Solidity: function supplyAtLock() view returns(uint256)
func (_AlphaToken *AlphaTokenCaller) SupplyAtLock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AlphaToken.contract.Call(opts, &out, "supplyAtLock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SupplyAtLock is a free data retrieval call binding the contract method 0xea32345b.
//
// Solidity: function supplyAtLock() view returns(uint256)
func (_AlphaToken *AlphaTokenSession) SupplyAtLock() (*big.Int, error) {
	return _AlphaToken.Contract.SupplyAtLock(&_AlphaToken.CallOpts)
}

// SupplyAtLock is a free data retrieval call binding the contract method 0xea32345b.
//
// Solidity: function supplyAtLock() view returns(uint256)
func (_AlphaToken *AlphaTokenCallerSession) SupplyAtLock() (*big.Int, error) {
	return _AlphaToken.Contract.SupplyAtLock(&_AlphaToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AlphaToken *AlphaTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AlphaToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AlphaToken *AlphaTokenSession) Symbol() (string, error) {
	return _AlphaToken.Contract.Symbol(&_AlphaToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AlphaToken *AlphaTokenCallerSession) Symbol() (string, error) {
	return _AlphaToken.Contract.Symbol(&_AlphaToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_AlphaToken *AlphaTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AlphaToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_AlphaToken *AlphaTokenSession) TotalSupply() (*big.Int, error) {
	return _AlphaToken.Contract.TotalSupply(&_AlphaToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_AlphaToken *AlphaTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _AlphaToken.Contract.TotalSupply(&_AlphaToken.CallOpts)
}

// WorknetId is a free data retrieval call binding the contract method 0x08e5d167.
//
// Solidity: function worknetId() view returns(uint256)
func (_AlphaToken *AlphaTokenCaller) WorknetId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AlphaToken.contract.Call(opts, &out, "worknetId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WorknetId is a free data retrieval call binding the contract method 0x08e5d167.
//
// Solidity: function worknetId() view returns(uint256)
func (_AlphaToken *AlphaTokenSession) WorknetId() (*big.Int, error) {
	return _AlphaToken.Contract.WorknetId(&_AlphaToken.CallOpts)
}

// WorknetId is a free data retrieval call binding the contract method 0x08e5d167.
//
// Solidity: function worknetId() view returns(uint256)
func (_AlphaToken *AlphaTokenCallerSession) WorknetId() (*big.Int, error) {
	return _AlphaToken.Contract.WorknetId(&_AlphaToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_AlphaToken *AlphaTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _AlphaToken.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_AlphaToken *AlphaTokenSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _AlphaToken.Contract.Approve(&_AlphaToken.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_AlphaToken *AlphaTokenTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _AlphaToken.Contract.Approve(&_AlphaToken.TransactOpts, spender, value)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(address spender, uint256 amount, bytes data) returns(bool)
func (_AlphaToken *AlphaTokenTransactor) ApproveAndCall(opts *bind.TransactOpts, spender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _AlphaToken.contract.Transact(opts, "approveAndCall", spender, amount, data)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(address spender, uint256 amount, bytes data) returns(bool)
func (_AlphaToken *AlphaTokenSession) ApproveAndCall(spender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _AlphaToken.Contract.ApproveAndCall(&_AlphaToken.TransactOpts, spender, amount, data)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(address spender, uint256 amount, bytes data) returns(bool)
func (_AlphaToken *AlphaTokenTransactorSession) ApproveAndCall(spender common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _AlphaToken.Contract.ApproveAndCall(&_AlphaToken.TransactOpts, spender, amount, data)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_AlphaToken *AlphaTokenTransactor) Burn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _AlphaToken.contract.Transact(opts, "burn", value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_AlphaToken *AlphaTokenSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _AlphaToken.Contract.Burn(&_AlphaToken.TransactOpts, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_AlphaToken *AlphaTokenTransactorSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _AlphaToken.Contract.Burn(&_AlphaToken.TransactOpts, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_AlphaToken *AlphaTokenTransactor) BurnFrom(opts *bind.TransactOpts, account common.Address, value *big.Int) (*types.Transaction, error) {
	return _AlphaToken.contract.Transact(opts, "burnFrom", account, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_AlphaToken *AlphaTokenSession) BurnFrom(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _AlphaToken.Contract.BurnFrom(&_AlphaToken.TransactOpts, account, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_AlphaToken *AlphaTokenTransactorSession) BurnFrom(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _AlphaToken.Contract.BurnFrom(&_AlphaToken.TransactOpts, account, value)
}

// Initialize is a paid mutator transaction binding the contract method 0xbd3a13f6.
//
// Solidity: function initialize(string name_, string symbol_, uint256 worknetId_, address admin_) returns()
func (_AlphaToken *AlphaTokenTransactor) Initialize(opts *bind.TransactOpts, name_ string, symbol_ string, worknetId_ *big.Int, admin_ common.Address) (*types.Transaction, error) {
	return _AlphaToken.contract.Transact(opts, "initialize", name_, symbol_, worknetId_, admin_)
}

// Initialize is a paid mutator transaction binding the contract method 0xbd3a13f6.
//
// Solidity: function initialize(string name_, string symbol_, uint256 worknetId_, address admin_) returns()
func (_AlphaToken *AlphaTokenSession) Initialize(name_ string, symbol_ string, worknetId_ *big.Int, admin_ common.Address) (*types.Transaction, error) {
	return _AlphaToken.Contract.Initialize(&_AlphaToken.TransactOpts, name_, symbol_, worknetId_, admin_)
}

// Initialize is a paid mutator transaction binding the contract method 0xbd3a13f6.
//
// Solidity: function initialize(string name_, string symbol_, uint256 worknetId_, address admin_) returns()
func (_AlphaToken *AlphaTokenTransactorSession) Initialize(name_ string, symbol_ string, worknetId_ *big.Int, admin_ common.Address) (*types.Transaction, error) {
	return _AlphaToken.Contract.Initialize(&_AlphaToken.TransactOpts, name_, symbol_, worknetId_, admin_)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_AlphaToken *AlphaTokenTransactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AlphaToken.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_AlphaToken *AlphaTokenSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AlphaToken.Contract.Mint(&_AlphaToken.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_AlphaToken *AlphaTokenTransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AlphaToken.Contract.Mint(&_AlphaToken.TransactOpts, to, amount)
}

// SetMinterPaused is a paid mutator transaction binding the contract method 0xa12a0ace.
//
// Solidity: function setMinterPaused(address minter, bool paused) returns()
func (_AlphaToken *AlphaTokenTransactor) SetMinterPaused(opts *bind.TransactOpts, minter common.Address, paused bool) (*types.Transaction, error) {
	return _AlphaToken.contract.Transact(opts, "setMinterPaused", minter, paused)
}

// SetMinterPaused is a paid mutator transaction binding the contract method 0xa12a0ace.
//
// Solidity: function setMinterPaused(address minter, bool paused) returns()
func (_AlphaToken *AlphaTokenSession) SetMinterPaused(minter common.Address, paused bool) (*types.Transaction, error) {
	return _AlphaToken.Contract.SetMinterPaused(&_AlphaToken.TransactOpts, minter, paused)
}

// SetMinterPaused is a paid mutator transaction binding the contract method 0xa12a0ace.
//
// Solidity: function setMinterPaused(address minter, bool paused) returns()
func (_AlphaToken *AlphaTokenTransactorSession) SetMinterPaused(minter common.Address, paused bool) (*types.Transaction, error) {
	return _AlphaToken.Contract.SetMinterPaused(&_AlphaToken.TransactOpts, minter, paused)
}

// SetWorknetMinter is a paid mutator transaction binding the contract method 0xb218a8ee.
//
// Solidity: function setWorknetMinter(address worknetManager) returns()
func (_AlphaToken *AlphaTokenTransactor) SetWorknetMinter(opts *bind.TransactOpts, worknetManager common.Address) (*types.Transaction, error) {
	return _AlphaToken.contract.Transact(opts, "setWorknetMinter", worknetManager)
}

// SetWorknetMinter is a paid mutator transaction binding the contract method 0xb218a8ee.
//
// Solidity: function setWorknetMinter(address worknetManager) returns()
func (_AlphaToken *AlphaTokenSession) SetWorknetMinter(worknetManager common.Address) (*types.Transaction, error) {
	return _AlphaToken.Contract.SetWorknetMinter(&_AlphaToken.TransactOpts, worknetManager)
}

// SetWorknetMinter is a paid mutator transaction binding the contract method 0xb218a8ee.
//
// Solidity: function setWorknetMinter(address worknetManager) returns()
func (_AlphaToken *AlphaTokenTransactorSession) SetWorknetMinter(worknetManager common.Address) (*types.Transaction, error) {
	return _AlphaToken.Contract.SetWorknetMinter(&_AlphaToken.TransactOpts, worknetManager)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_AlphaToken *AlphaTokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AlphaToken.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_AlphaToken *AlphaTokenSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AlphaToken.Contract.Transfer(&_AlphaToken.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_AlphaToken *AlphaTokenTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AlphaToken.Contract.Transfer(&_AlphaToken.TransactOpts, to, value)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 amount, bytes data) returns(bool)
func (_AlphaToken *AlphaTokenTransactor) TransferAndCall(opts *bind.TransactOpts, to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _AlphaToken.contract.Transact(opts, "transferAndCall", to, amount, data)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 amount, bytes data) returns(bool)
func (_AlphaToken *AlphaTokenSession) TransferAndCall(to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _AlphaToken.Contract.TransferAndCall(&_AlphaToken.TransactOpts, to, amount, data)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 amount, bytes data) returns(bool)
func (_AlphaToken *AlphaTokenTransactorSession) TransferAndCall(to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _AlphaToken.Contract.TransferAndCall(&_AlphaToken.TransactOpts, to, amount, data)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_AlphaToken *AlphaTokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AlphaToken.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_AlphaToken *AlphaTokenSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AlphaToken.Contract.TransferFrom(&_AlphaToken.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_AlphaToken *AlphaTokenTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AlphaToken.Contract.TransferFrom(&_AlphaToken.TransactOpts, from, to, value)
}

// AlphaTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the AlphaToken contract.
type AlphaTokenApprovalIterator struct {
	Event *AlphaTokenApproval // Event containing the contract specifics and raw log

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
func (it *AlphaTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AlphaTokenApproval)
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
		it.Event = new(AlphaTokenApproval)
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
func (it *AlphaTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AlphaTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AlphaTokenApproval represents a Approval event raised by the AlphaToken contract.
type AlphaTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_AlphaToken *AlphaTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*AlphaTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _AlphaToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &AlphaTokenApprovalIterator{contract: _AlphaToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_AlphaToken *AlphaTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *AlphaTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _AlphaToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AlphaTokenApproval)
				if err := _AlphaToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_AlphaToken *AlphaTokenFilterer) ParseApproval(log types.Log) (*AlphaTokenApproval, error) {
	event := new(AlphaTokenApproval)
	if err := _AlphaToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AlphaTokenInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the AlphaToken contract.
type AlphaTokenInitializedIterator struct {
	Event *AlphaTokenInitialized // Event containing the contract specifics and raw log

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
func (it *AlphaTokenInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AlphaTokenInitialized)
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
		it.Event = new(AlphaTokenInitialized)
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
func (it *AlphaTokenInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AlphaTokenInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AlphaTokenInitialized represents a Initialized event raised by the AlphaToken contract.
type AlphaTokenInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AlphaToken *AlphaTokenFilterer) FilterInitialized(opts *bind.FilterOpts) (*AlphaTokenInitializedIterator, error) {

	logs, sub, err := _AlphaToken.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &AlphaTokenInitializedIterator{contract: _AlphaToken.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AlphaToken *AlphaTokenFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *AlphaTokenInitialized) (event.Subscription, error) {

	logs, sub, err := _AlphaToken.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AlphaTokenInitialized)
				if err := _AlphaToken.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AlphaToken *AlphaTokenFilterer) ParseInitialized(log types.Log) (*AlphaTokenInitialized, error) {
	event := new(AlphaTokenInitialized)
	if err := _AlphaToken.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AlphaTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the AlphaToken contract.
type AlphaTokenTransferIterator struct {
	Event *AlphaTokenTransfer // Event containing the contract specifics and raw log

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
func (it *AlphaTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AlphaTokenTransfer)
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
		it.Event = new(AlphaTokenTransfer)
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
func (it *AlphaTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AlphaTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AlphaTokenTransfer represents a Transfer event raised by the AlphaToken contract.
type AlphaTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_AlphaToken *AlphaTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AlphaTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AlphaToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AlphaTokenTransferIterator{contract: _AlphaToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_AlphaToken *AlphaTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *AlphaTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AlphaToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AlphaTokenTransfer)
				if err := _AlphaToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_AlphaToken *AlphaTokenFilterer) ParseTransfer(log types.Log) (*AlphaTokenTransfer, error) {
	event := new(AlphaTokenTransfer)
	if err := _AlphaToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AlphaTokenWorknetMinterSetIterator is returned from FilterWorknetMinterSet and is used to iterate over the raw logs and unpacked data for WorknetMinterSet events raised by the AlphaToken contract.
type AlphaTokenWorknetMinterSetIterator struct {
	Event *AlphaTokenWorknetMinterSet // Event containing the contract specifics and raw log

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
func (it *AlphaTokenWorknetMinterSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AlphaTokenWorknetMinterSet)
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
		it.Event = new(AlphaTokenWorknetMinterSet)
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
func (it *AlphaTokenWorknetMinterSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AlphaTokenWorknetMinterSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AlphaTokenWorknetMinterSet represents a WorknetMinterSet event raised by the AlphaToken contract.
type AlphaTokenWorknetMinterSet struct {
	WorknetManager common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterWorknetMinterSet is a free log retrieval operation binding the contract event 0x21adb39f6834aac66de042684d84e00f5ecd3b57a6e022b84b8531a113c8b1c0.
//
// Solidity: event WorknetMinterSet(address indexed worknetManager)
func (_AlphaToken *AlphaTokenFilterer) FilterWorknetMinterSet(opts *bind.FilterOpts, worknetManager []common.Address) (*AlphaTokenWorknetMinterSetIterator, error) {

	var worknetManagerRule []interface{}
	for _, worknetManagerItem := range worknetManager {
		worknetManagerRule = append(worknetManagerRule, worknetManagerItem)
	}

	logs, sub, err := _AlphaToken.contract.FilterLogs(opts, "WorknetMinterSet", worknetManagerRule)
	if err != nil {
		return nil, err
	}
	return &AlphaTokenWorknetMinterSetIterator{contract: _AlphaToken.contract, event: "WorknetMinterSet", logs: logs, sub: sub}, nil
}

// WatchWorknetMinterSet is a free log subscription operation binding the contract event 0x21adb39f6834aac66de042684d84e00f5ecd3b57a6e022b84b8531a113c8b1c0.
//
// Solidity: event WorknetMinterSet(address indexed worknetManager)
func (_AlphaToken *AlphaTokenFilterer) WatchWorknetMinterSet(opts *bind.WatchOpts, sink chan<- *AlphaTokenWorknetMinterSet, worknetManager []common.Address) (event.Subscription, error) {

	var worknetManagerRule []interface{}
	for _, worknetManagerItem := range worknetManager {
		worknetManagerRule = append(worknetManagerRule, worknetManagerItem)
	}

	logs, sub, err := _AlphaToken.contract.WatchLogs(opts, "WorknetMinterSet", worknetManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AlphaTokenWorknetMinterSet)
				if err := _AlphaToken.contract.UnpackLog(event, "WorknetMinterSet", log); err != nil {
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

// ParseWorknetMinterSet is a log parse operation binding the contract event 0x21adb39f6834aac66de042684d84e00f5ecd3b57a6e022b84b8531a113c8b1c0.
//
// Solidity: event WorknetMinterSet(address indexed worknetManager)
func (_AlphaToken *AlphaTokenFilterer) ParseWorknetMinterSet(log types.Log) (*AlphaTokenWorknetMinterSet, error) {
	event := new(AlphaTokenWorknetMinterSet)
	if err := _AlphaToken.contract.UnpackLog(event, "WorknetMinterSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
