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

// PoolManagerReaderMetaData contains all meta data concerning the PoolManagerReader contract.
var PoolManagerReaderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"}],\"name\":\"getSlot0\",\"outputs\":[{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96\",\"type\":\"uint160\"},{\"internalType\":\"int24\",\"name\":\"tick\",\"type\":\"int24\"},{\"internalType\":\"uint24\",\"name\":\"protocolFee\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"lpFee\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// PoolManagerReaderABI is the input ABI used to generate the binding from.
// Deprecated: Use PoolManagerReaderMetaData.ABI instead.
var PoolManagerReaderABI = PoolManagerReaderMetaData.ABI

// PoolManagerReader is an auto generated Go binding around an Ethereum contract.
type PoolManagerReader struct {
	PoolManagerReaderCaller     // Read-only binding to the contract
	PoolManagerReaderTransactor // Write-only binding to the contract
	PoolManagerReaderFilterer   // Log filterer for contract events
}

// PoolManagerReaderCaller is an auto generated read-only Go binding around an Ethereum contract.
type PoolManagerReaderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolManagerReaderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PoolManagerReaderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolManagerReaderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PoolManagerReaderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolManagerReaderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PoolManagerReaderSession struct {
	Contract     *PoolManagerReader // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// PoolManagerReaderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PoolManagerReaderCallerSession struct {
	Contract *PoolManagerReaderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// PoolManagerReaderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PoolManagerReaderTransactorSession struct {
	Contract     *PoolManagerReaderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// PoolManagerReaderRaw is an auto generated low-level Go binding around an Ethereum contract.
type PoolManagerReaderRaw struct {
	Contract *PoolManagerReader // Generic contract binding to access the raw methods on
}

// PoolManagerReaderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PoolManagerReaderCallerRaw struct {
	Contract *PoolManagerReaderCaller // Generic read-only contract binding to access the raw methods on
}

// PoolManagerReaderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PoolManagerReaderTransactorRaw struct {
	Contract *PoolManagerReaderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPoolManagerReader creates a new instance of PoolManagerReader, bound to a specific deployed contract.
func NewPoolManagerReader(address common.Address, backend bind.ContractBackend) (*PoolManagerReader, error) {
	contract, err := bindPoolManagerReader(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PoolManagerReader{PoolManagerReaderCaller: PoolManagerReaderCaller{contract: contract}, PoolManagerReaderTransactor: PoolManagerReaderTransactor{contract: contract}, PoolManagerReaderFilterer: PoolManagerReaderFilterer{contract: contract}}, nil
}

// NewPoolManagerReaderCaller creates a new read-only instance of PoolManagerReader, bound to a specific deployed contract.
func NewPoolManagerReaderCaller(address common.Address, caller bind.ContractCaller) (*PoolManagerReaderCaller, error) {
	contract, err := bindPoolManagerReader(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PoolManagerReaderCaller{contract: contract}, nil
}

// NewPoolManagerReaderTransactor creates a new write-only instance of PoolManagerReader, bound to a specific deployed contract.
func NewPoolManagerReaderTransactor(address common.Address, transactor bind.ContractTransactor) (*PoolManagerReaderTransactor, error) {
	contract, err := bindPoolManagerReader(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PoolManagerReaderTransactor{contract: contract}, nil
}

// NewPoolManagerReaderFilterer creates a new log filterer instance of PoolManagerReader, bound to a specific deployed contract.
func NewPoolManagerReaderFilterer(address common.Address, filterer bind.ContractFilterer) (*PoolManagerReaderFilterer, error) {
	contract, err := bindPoolManagerReader(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PoolManagerReaderFilterer{contract: contract}, nil
}

// bindPoolManagerReader binds a generic wrapper to an already deployed contract.
func bindPoolManagerReader(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PoolManagerReaderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoolManagerReader *PoolManagerReaderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PoolManagerReader.Contract.PoolManagerReaderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoolManagerReader *PoolManagerReaderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolManagerReader.Contract.PoolManagerReaderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoolManagerReader *PoolManagerReaderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoolManagerReader.Contract.PoolManagerReaderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoolManagerReader *PoolManagerReaderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PoolManagerReader.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoolManagerReader *PoolManagerReaderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolManagerReader.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoolManagerReader *PoolManagerReaderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoolManagerReader.Contract.contract.Transact(opts, method, params...)
}

// GetSlot0 is a free data retrieval call binding the contract method 0xc815641c.
//
// Solidity: function getSlot0(bytes32 id) view returns(uint160 sqrtPriceX96, int24 tick, uint24 protocolFee, uint24 lpFee)
func (_PoolManagerReader *PoolManagerReaderCaller) GetSlot0(opts *bind.CallOpts, id [32]byte) (struct {
	SqrtPriceX96 *big.Int
	Tick         *big.Int
	ProtocolFee  *big.Int
	LpFee        *big.Int
}, error) {
	var out []interface{}
	err := _PoolManagerReader.contract.Call(opts, &out, "getSlot0", id)

	outstruct := new(struct {
		SqrtPriceX96 *big.Int
		Tick         *big.Int
		ProtocolFee  *big.Int
		LpFee        *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SqrtPriceX96 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Tick = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ProtocolFee = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.LpFee = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetSlot0 is a free data retrieval call binding the contract method 0xc815641c.
//
// Solidity: function getSlot0(bytes32 id) view returns(uint160 sqrtPriceX96, int24 tick, uint24 protocolFee, uint24 lpFee)
func (_PoolManagerReader *PoolManagerReaderSession) GetSlot0(id [32]byte) (struct {
	SqrtPriceX96 *big.Int
	Tick         *big.Int
	ProtocolFee  *big.Int
	LpFee        *big.Int
}, error) {
	return _PoolManagerReader.Contract.GetSlot0(&_PoolManagerReader.CallOpts, id)
}

// GetSlot0 is a free data retrieval call binding the contract method 0xc815641c.
//
// Solidity: function getSlot0(bytes32 id) view returns(uint160 sqrtPriceX96, int24 tick, uint24 protocolFee, uint24 lpFee)
func (_PoolManagerReader *PoolManagerReaderCallerSession) GetSlot0(id [32]byte) (struct {
	SqrtPriceX96 *big.Int
	Tick         *big.Int
	ProtocolFee  *big.Int
	LpFee        *big.Int
}, error) {
	return _PoolManagerReader.Contract.GetSlot0(&_PoolManagerReader.CallOpts, id)
}
