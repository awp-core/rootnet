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

// LPManagerBaseMetaData contains all meta data concerning the LPManagerBase contract.
var LPManagerBaseMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"MAX_SQRT_RATIO\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint160\",\"internalType\":\"uint160\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_TICK\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_SQRT_RATIO\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint160\",\"internalType\":\"uint160\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_TICK\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"POOL_FEE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint24\",\"internalType\":\"uint24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"TICK_SPACING\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"int24\",\"internalType\":\"int24\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"awpRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"awpToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"compoundFees\",\"inputs\":[{\"name\":\"worknetToken\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createPoolAndAddLiquidity\",\"inputs\":[{\"name\":\"worknetToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"awpAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"worknetTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"poolId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"lpTokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"needsCompounding\",\"inputs\":[{\"name\":\"worknetToken\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"hasPool\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"permit2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"worknetTokenToPoolId\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"worknetTokenToTokenId\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"FeesCompounded\",\"inputs\":[{\"name\":\"worknetToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AmountExceedsPermit2Limit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoPool\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAWPRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotGuardian\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PoolAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
}

// LPManagerBaseABI is the input ABI used to generate the binding from.
// Deprecated: Use LPManagerBaseMetaData.ABI instead.
var LPManagerBaseABI = LPManagerBaseMetaData.ABI

// LPManagerBase is an auto generated Go binding around an Ethereum contract.
type LPManagerBase struct {
	LPManagerBaseCaller     // Read-only binding to the contract
	LPManagerBaseTransactor // Write-only binding to the contract
	LPManagerBaseFilterer   // Log filterer for contract events
}

// LPManagerBaseCaller is an auto generated read-only Go binding around an Ethereum contract.
type LPManagerBaseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LPManagerBaseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LPManagerBaseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LPManagerBaseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LPManagerBaseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LPManagerBaseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LPManagerBaseSession struct {
	Contract     *LPManagerBase    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LPManagerBaseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LPManagerBaseCallerSession struct {
	Contract *LPManagerBaseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// LPManagerBaseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LPManagerBaseTransactorSession struct {
	Contract     *LPManagerBaseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// LPManagerBaseRaw is an auto generated low-level Go binding around an Ethereum contract.
type LPManagerBaseRaw struct {
	Contract *LPManagerBase // Generic contract binding to access the raw methods on
}

// LPManagerBaseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LPManagerBaseCallerRaw struct {
	Contract *LPManagerBaseCaller // Generic read-only contract binding to access the raw methods on
}

// LPManagerBaseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LPManagerBaseTransactorRaw struct {
	Contract *LPManagerBaseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLPManagerBase creates a new instance of LPManagerBase, bound to a specific deployed contract.
func NewLPManagerBase(address common.Address, backend bind.ContractBackend) (*LPManagerBase, error) {
	contract, err := bindLPManagerBase(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LPManagerBase{LPManagerBaseCaller: LPManagerBaseCaller{contract: contract}, LPManagerBaseTransactor: LPManagerBaseTransactor{contract: contract}, LPManagerBaseFilterer: LPManagerBaseFilterer{contract: contract}}, nil
}

// NewLPManagerBaseCaller creates a new read-only instance of LPManagerBase, bound to a specific deployed contract.
func NewLPManagerBaseCaller(address common.Address, caller bind.ContractCaller) (*LPManagerBaseCaller, error) {
	contract, err := bindLPManagerBase(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LPManagerBaseCaller{contract: contract}, nil
}

// NewLPManagerBaseTransactor creates a new write-only instance of LPManagerBase, bound to a specific deployed contract.
func NewLPManagerBaseTransactor(address common.Address, transactor bind.ContractTransactor) (*LPManagerBaseTransactor, error) {
	contract, err := bindLPManagerBase(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LPManagerBaseTransactor{contract: contract}, nil
}

// NewLPManagerBaseFilterer creates a new log filterer instance of LPManagerBase, bound to a specific deployed contract.
func NewLPManagerBaseFilterer(address common.Address, filterer bind.ContractFilterer) (*LPManagerBaseFilterer, error) {
	contract, err := bindLPManagerBase(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LPManagerBaseFilterer{contract: contract}, nil
}

// bindLPManagerBase binds a generic wrapper to an already deployed contract.
func bindLPManagerBase(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LPManagerBaseMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LPManagerBase *LPManagerBaseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LPManagerBase.Contract.LPManagerBaseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LPManagerBase *LPManagerBaseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LPManagerBase.Contract.LPManagerBaseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LPManagerBase *LPManagerBaseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LPManagerBase.Contract.LPManagerBaseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LPManagerBase *LPManagerBaseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LPManagerBase.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LPManagerBase *LPManagerBaseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LPManagerBase.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LPManagerBase *LPManagerBaseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LPManagerBase.Contract.contract.Transact(opts, method, params...)
}

// MAXSQRTRATIO is a free data retrieval call binding the contract method 0x6d2cc304.
//
// Solidity: function MAX_SQRT_RATIO() view returns(uint160)
func (_LPManagerBase *LPManagerBaseCaller) MAXSQRTRATIO(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LPManagerBase.contract.Call(opts, &out, "MAX_SQRT_RATIO")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXSQRTRATIO is a free data retrieval call binding the contract method 0x6d2cc304.
//
// Solidity: function MAX_SQRT_RATIO() view returns(uint160)
func (_LPManagerBase *LPManagerBaseSession) MAXSQRTRATIO() (*big.Int, error) {
	return _LPManagerBase.Contract.MAXSQRTRATIO(&_LPManagerBase.CallOpts)
}

// MAXSQRTRATIO is a free data retrieval call binding the contract method 0x6d2cc304.
//
// Solidity: function MAX_SQRT_RATIO() view returns(uint160)
func (_LPManagerBase *LPManagerBaseCallerSession) MAXSQRTRATIO() (*big.Int, error) {
	return _LPManagerBase.Contract.MAXSQRTRATIO(&_LPManagerBase.CallOpts)
}

// MAXTICK is a free data retrieval call binding the contract method 0x6882a888.
//
// Solidity: function MAX_TICK() view returns(int24)
func (_LPManagerBase *LPManagerBaseCaller) MAXTICK(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LPManagerBase.contract.Call(opts, &out, "MAX_TICK")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXTICK is a free data retrieval call binding the contract method 0x6882a888.
//
// Solidity: function MAX_TICK() view returns(int24)
func (_LPManagerBase *LPManagerBaseSession) MAXTICK() (*big.Int, error) {
	return _LPManagerBase.Contract.MAXTICK(&_LPManagerBase.CallOpts)
}

// MAXTICK is a free data retrieval call binding the contract method 0x6882a888.
//
// Solidity: function MAX_TICK() view returns(int24)
func (_LPManagerBase *LPManagerBaseCallerSession) MAXTICK() (*big.Int, error) {
	return _LPManagerBase.Contract.MAXTICK(&_LPManagerBase.CallOpts)
}

// MINSQRTRATIO is a free data retrieval call binding the contract method 0xee8847ff.
//
// Solidity: function MIN_SQRT_RATIO() view returns(uint160)
func (_LPManagerBase *LPManagerBaseCaller) MINSQRTRATIO(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LPManagerBase.contract.Call(opts, &out, "MIN_SQRT_RATIO")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINSQRTRATIO is a free data retrieval call binding the contract method 0xee8847ff.
//
// Solidity: function MIN_SQRT_RATIO() view returns(uint160)
func (_LPManagerBase *LPManagerBaseSession) MINSQRTRATIO() (*big.Int, error) {
	return _LPManagerBase.Contract.MINSQRTRATIO(&_LPManagerBase.CallOpts)
}

// MINSQRTRATIO is a free data retrieval call binding the contract method 0xee8847ff.
//
// Solidity: function MIN_SQRT_RATIO() view returns(uint160)
func (_LPManagerBase *LPManagerBaseCallerSession) MINSQRTRATIO() (*big.Int, error) {
	return _LPManagerBase.Contract.MINSQRTRATIO(&_LPManagerBase.CallOpts)
}

// MINTICK is a free data retrieval call binding the contract method 0xa1634b14.
//
// Solidity: function MIN_TICK() view returns(int24)
func (_LPManagerBase *LPManagerBaseCaller) MINTICK(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LPManagerBase.contract.Call(opts, &out, "MIN_TICK")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINTICK is a free data retrieval call binding the contract method 0xa1634b14.
//
// Solidity: function MIN_TICK() view returns(int24)
func (_LPManagerBase *LPManagerBaseSession) MINTICK() (*big.Int, error) {
	return _LPManagerBase.Contract.MINTICK(&_LPManagerBase.CallOpts)
}

// MINTICK is a free data retrieval call binding the contract method 0xa1634b14.
//
// Solidity: function MIN_TICK() view returns(int24)
func (_LPManagerBase *LPManagerBaseCallerSession) MINTICK() (*big.Int, error) {
	return _LPManagerBase.Contract.MINTICK(&_LPManagerBase.CallOpts)
}

// POOLFEE is a free data retrieval call binding the contract method 0xdd1b9c4a.
//
// Solidity: function POOL_FEE() view returns(uint24)
func (_LPManagerBase *LPManagerBaseCaller) POOLFEE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LPManagerBase.contract.Call(opts, &out, "POOL_FEE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// POOLFEE is a free data retrieval call binding the contract method 0xdd1b9c4a.
//
// Solidity: function POOL_FEE() view returns(uint24)
func (_LPManagerBase *LPManagerBaseSession) POOLFEE() (*big.Int, error) {
	return _LPManagerBase.Contract.POOLFEE(&_LPManagerBase.CallOpts)
}

// POOLFEE is a free data retrieval call binding the contract method 0xdd1b9c4a.
//
// Solidity: function POOL_FEE() view returns(uint24)
func (_LPManagerBase *LPManagerBaseCallerSession) POOLFEE() (*big.Int, error) {
	return _LPManagerBase.Contract.POOLFEE(&_LPManagerBase.CallOpts)
}

// TICKSPACING is a free data retrieval call binding the contract method 0x46ca626b.
//
// Solidity: function TICK_SPACING() view returns(int24)
func (_LPManagerBase *LPManagerBaseCaller) TICKSPACING(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LPManagerBase.contract.Call(opts, &out, "TICK_SPACING")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TICKSPACING is a free data retrieval call binding the contract method 0x46ca626b.
//
// Solidity: function TICK_SPACING() view returns(int24)
func (_LPManagerBase *LPManagerBaseSession) TICKSPACING() (*big.Int, error) {
	return _LPManagerBase.Contract.TICKSPACING(&_LPManagerBase.CallOpts)
}

// TICKSPACING is a free data retrieval call binding the contract method 0x46ca626b.
//
// Solidity: function TICK_SPACING() view returns(int24)
func (_LPManagerBase *LPManagerBaseCallerSession) TICKSPACING() (*big.Int, error) {
	return _LPManagerBase.Contract.TICKSPACING(&_LPManagerBase.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_LPManagerBase *LPManagerBaseCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LPManagerBase.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_LPManagerBase *LPManagerBaseSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _LPManagerBase.Contract.UPGRADEINTERFACEVERSION(&_LPManagerBase.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_LPManagerBase *LPManagerBaseCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _LPManagerBase.Contract.UPGRADEINTERFACEVERSION(&_LPManagerBase.CallOpts)
}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_LPManagerBase *LPManagerBaseCaller) AwpRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LPManagerBase.contract.Call(opts, &out, "awpRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_LPManagerBase *LPManagerBaseSession) AwpRegistry() (common.Address, error) {
	return _LPManagerBase.Contract.AwpRegistry(&_LPManagerBase.CallOpts)
}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_LPManagerBase *LPManagerBaseCallerSession) AwpRegistry() (common.Address, error) {
	return _LPManagerBase.Contract.AwpRegistry(&_LPManagerBase.CallOpts)
}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_LPManagerBase *LPManagerBaseCaller) AwpToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LPManagerBase.contract.Call(opts, &out, "awpToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_LPManagerBase *LPManagerBaseSession) AwpToken() (common.Address, error) {
	return _LPManagerBase.Contract.AwpToken(&_LPManagerBase.CallOpts)
}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_LPManagerBase *LPManagerBaseCallerSession) AwpToken() (common.Address, error) {
	return _LPManagerBase.Contract.AwpToken(&_LPManagerBase.CallOpts)
}

// NeedsCompounding is a free data retrieval call binding the contract method 0x382c1706.
//
// Solidity: function needsCompounding(address worknetToken) view returns(bool hasPool, uint256 tokenId)
func (_LPManagerBase *LPManagerBaseCaller) NeedsCompounding(opts *bind.CallOpts, worknetToken common.Address) (struct {
	HasPool bool
	TokenId *big.Int
}, error) {
	var out []interface{}
	err := _LPManagerBase.contract.Call(opts, &out, "needsCompounding", worknetToken)

	outstruct := new(struct {
		HasPool bool
		TokenId *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.HasPool = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.TokenId = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// NeedsCompounding is a free data retrieval call binding the contract method 0x382c1706.
//
// Solidity: function needsCompounding(address worknetToken) view returns(bool hasPool, uint256 tokenId)
func (_LPManagerBase *LPManagerBaseSession) NeedsCompounding(worknetToken common.Address) (struct {
	HasPool bool
	TokenId *big.Int
}, error) {
	return _LPManagerBase.Contract.NeedsCompounding(&_LPManagerBase.CallOpts, worknetToken)
}

// NeedsCompounding is a free data retrieval call binding the contract method 0x382c1706.
//
// Solidity: function needsCompounding(address worknetToken) view returns(bool hasPool, uint256 tokenId)
func (_LPManagerBase *LPManagerBaseCallerSession) NeedsCompounding(worknetToken common.Address) (struct {
	HasPool bool
	TokenId *big.Int
}, error) {
	return _LPManagerBase.Contract.NeedsCompounding(&_LPManagerBase.CallOpts, worknetToken)
}

// Permit2 is a free data retrieval call binding the contract method 0x12261ee7.
//
// Solidity: function permit2() view returns(address)
func (_LPManagerBase *LPManagerBaseCaller) Permit2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LPManagerBase.contract.Call(opts, &out, "permit2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Permit2 is a free data retrieval call binding the contract method 0x12261ee7.
//
// Solidity: function permit2() view returns(address)
func (_LPManagerBase *LPManagerBaseSession) Permit2() (common.Address, error) {
	return _LPManagerBase.Contract.Permit2(&_LPManagerBase.CallOpts)
}

// Permit2 is a free data retrieval call binding the contract method 0x12261ee7.
//
// Solidity: function permit2() view returns(address)
func (_LPManagerBase *LPManagerBaseCallerSession) Permit2() (common.Address, error) {
	return _LPManagerBase.Contract.Permit2(&_LPManagerBase.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_LPManagerBase *LPManagerBaseCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LPManagerBase.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_LPManagerBase *LPManagerBaseSession) ProxiableUUID() ([32]byte, error) {
	return _LPManagerBase.Contract.ProxiableUUID(&_LPManagerBase.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_LPManagerBase *LPManagerBaseCallerSession) ProxiableUUID() ([32]byte, error) {
	return _LPManagerBase.Contract.ProxiableUUID(&_LPManagerBase.CallOpts)
}

// WorknetTokenToPoolId is a free data retrieval call binding the contract method 0x39a5e566.
//
// Solidity: function worknetTokenToPoolId(address ) view returns(bytes32)
func (_LPManagerBase *LPManagerBaseCaller) WorknetTokenToPoolId(opts *bind.CallOpts, arg0 common.Address) ([32]byte, error) {
	var out []interface{}
	err := _LPManagerBase.contract.Call(opts, &out, "worknetTokenToPoolId", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// WorknetTokenToPoolId is a free data retrieval call binding the contract method 0x39a5e566.
//
// Solidity: function worknetTokenToPoolId(address ) view returns(bytes32)
func (_LPManagerBase *LPManagerBaseSession) WorknetTokenToPoolId(arg0 common.Address) ([32]byte, error) {
	return _LPManagerBase.Contract.WorknetTokenToPoolId(&_LPManagerBase.CallOpts, arg0)
}

// WorknetTokenToPoolId is a free data retrieval call binding the contract method 0x39a5e566.
//
// Solidity: function worknetTokenToPoolId(address ) view returns(bytes32)
func (_LPManagerBase *LPManagerBaseCallerSession) WorknetTokenToPoolId(arg0 common.Address) ([32]byte, error) {
	return _LPManagerBase.Contract.WorknetTokenToPoolId(&_LPManagerBase.CallOpts, arg0)
}

// WorknetTokenToTokenId is a free data retrieval call binding the contract method 0xbba47471.
//
// Solidity: function worknetTokenToTokenId(address ) view returns(uint256)
func (_LPManagerBase *LPManagerBaseCaller) WorknetTokenToTokenId(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LPManagerBase.contract.Call(opts, &out, "worknetTokenToTokenId", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WorknetTokenToTokenId is a free data retrieval call binding the contract method 0xbba47471.
//
// Solidity: function worknetTokenToTokenId(address ) view returns(uint256)
func (_LPManagerBase *LPManagerBaseSession) WorknetTokenToTokenId(arg0 common.Address) (*big.Int, error) {
	return _LPManagerBase.Contract.WorknetTokenToTokenId(&_LPManagerBase.CallOpts, arg0)
}

// WorknetTokenToTokenId is a free data retrieval call binding the contract method 0xbba47471.
//
// Solidity: function worknetTokenToTokenId(address ) view returns(uint256)
func (_LPManagerBase *LPManagerBaseCallerSession) WorknetTokenToTokenId(arg0 common.Address) (*big.Int, error) {
	return _LPManagerBase.Contract.WorknetTokenToTokenId(&_LPManagerBase.CallOpts, arg0)
}

// CompoundFees is a paid mutator transaction binding the contract method 0x683b8b61.
//
// Solidity: function compoundFees(address worknetToken) returns()
func (_LPManagerBase *LPManagerBaseTransactor) CompoundFees(opts *bind.TransactOpts, worknetToken common.Address) (*types.Transaction, error) {
	return _LPManagerBase.contract.Transact(opts, "compoundFees", worknetToken)
}

// CompoundFees is a paid mutator transaction binding the contract method 0x683b8b61.
//
// Solidity: function compoundFees(address worknetToken) returns()
func (_LPManagerBase *LPManagerBaseSession) CompoundFees(worknetToken common.Address) (*types.Transaction, error) {
	return _LPManagerBase.Contract.CompoundFees(&_LPManagerBase.TransactOpts, worknetToken)
}

// CompoundFees is a paid mutator transaction binding the contract method 0x683b8b61.
//
// Solidity: function compoundFees(address worknetToken) returns()
func (_LPManagerBase *LPManagerBaseTransactorSession) CompoundFees(worknetToken common.Address) (*types.Transaction, error) {
	return _LPManagerBase.Contract.CompoundFees(&_LPManagerBase.TransactOpts, worknetToken)
}

// CreatePoolAndAddLiquidity is a paid mutator transaction binding the contract method 0x8651b1cc.
//
// Solidity: function createPoolAndAddLiquidity(address worknetToken, uint256 awpAmount, uint256 worknetTokenAmount) returns(bytes32 poolId, uint256 lpTokenId)
func (_LPManagerBase *LPManagerBaseTransactor) CreatePoolAndAddLiquidity(opts *bind.TransactOpts, worknetToken common.Address, awpAmount *big.Int, worknetTokenAmount *big.Int) (*types.Transaction, error) {
	return _LPManagerBase.contract.Transact(opts, "createPoolAndAddLiquidity", worknetToken, awpAmount, worknetTokenAmount)
}

// CreatePoolAndAddLiquidity is a paid mutator transaction binding the contract method 0x8651b1cc.
//
// Solidity: function createPoolAndAddLiquidity(address worknetToken, uint256 awpAmount, uint256 worknetTokenAmount) returns(bytes32 poolId, uint256 lpTokenId)
func (_LPManagerBase *LPManagerBaseSession) CreatePoolAndAddLiquidity(worknetToken common.Address, awpAmount *big.Int, worknetTokenAmount *big.Int) (*types.Transaction, error) {
	return _LPManagerBase.Contract.CreatePoolAndAddLiquidity(&_LPManagerBase.TransactOpts, worknetToken, awpAmount, worknetTokenAmount)
}

// CreatePoolAndAddLiquidity is a paid mutator transaction binding the contract method 0x8651b1cc.
//
// Solidity: function createPoolAndAddLiquidity(address worknetToken, uint256 awpAmount, uint256 worknetTokenAmount) returns(bytes32 poolId, uint256 lpTokenId)
func (_LPManagerBase *LPManagerBaseTransactorSession) CreatePoolAndAddLiquidity(worknetToken common.Address, awpAmount *big.Int, worknetTokenAmount *big.Int) (*types.Transaction, error) {
	return _LPManagerBase.Contract.CreatePoolAndAddLiquidity(&_LPManagerBase.TransactOpts, worknetToken, awpAmount, worknetTokenAmount)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_LPManagerBase *LPManagerBaseTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LPManagerBase.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_LPManagerBase *LPManagerBaseSession) Initialize() (*types.Transaction, error) {
	return _LPManagerBase.Contract.Initialize(&_LPManagerBase.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_LPManagerBase *LPManagerBaseTransactorSession) Initialize() (*types.Transaction, error) {
	return _LPManagerBase.Contract.Initialize(&_LPManagerBase.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_LPManagerBase *LPManagerBaseTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _LPManagerBase.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_LPManagerBase *LPManagerBaseSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _LPManagerBase.Contract.UpgradeToAndCall(&_LPManagerBase.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_LPManagerBase *LPManagerBaseTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _LPManagerBase.Contract.UpgradeToAndCall(&_LPManagerBase.TransactOpts, newImplementation, data)
}

// LPManagerBaseFeesCompoundedIterator is returned from FilterFeesCompounded and is used to iterate over the raw logs and unpacked data for FeesCompounded events raised by the LPManagerBase contract.
type LPManagerBaseFeesCompoundedIterator struct {
	Event *LPManagerBaseFeesCompounded // Event containing the contract specifics and raw log

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
func (it *LPManagerBaseFeesCompoundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LPManagerBaseFeesCompounded)
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
		it.Event = new(LPManagerBaseFeesCompounded)
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
func (it *LPManagerBaseFeesCompoundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LPManagerBaseFeesCompoundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LPManagerBaseFeesCompounded represents a FeesCompounded event raised by the LPManagerBase contract.
type LPManagerBaseFeesCompounded struct {
	WorknetToken common.Address
	TokenId      *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterFeesCompounded is a free log retrieval operation binding the contract event 0xbb85ff3bb01b1e5e9b75cd5e5ac87c83c843e76ca01a1565d7e050a4c8dbd223.
//
// Solidity: event FeesCompounded(address indexed worknetToken, uint256 tokenId)
func (_LPManagerBase *LPManagerBaseFilterer) FilterFeesCompounded(opts *bind.FilterOpts, worknetToken []common.Address) (*LPManagerBaseFeesCompoundedIterator, error) {

	var worknetTokenRule []interface{}
	for _, worknetTokenItem := range worknetToken {
		worknetTokenRule = append(worknetTokenRule, worknetTokenItem)
	}

	logs, sub, err := _LPManagerBase.contract.FilterLogs(opts, "FeesCompounded", worknetTokenRule)
	if err != nil {
		return nil, err
	}
	return &LPManagerBaseFeesCompoundedIterator{contract: _LPManagerBase.contract, event: "FeesCompounded", logs: logs, sub: sub}, nil
}

// WatchFeesCompounded is a free log subscription operation binding the contract event 0xbb85ff3bb01b1e5e9b75cd5e5ac87c83c843e76ca01a1565d7e050a4c8dbd223.
//
// Solidity: event FeesCompounded(address indexed worknetToken, uint256 tokenId)
func (_LPManagerBase *LPManagerBaseFilterer) WatchFeesCompounded(opts *bind.WatchOpts, sink chan<- *LPManagerBaseFeesCompounded, worknetToken []common.Address) (event.Subscription, error) {

	var worknetTokenRule []interface{}
	for _, worknetTokenItem := range worknetToken {
		worknetTokenRule = append(worknetTokenRule, worknetTokenItem)
	}

	logs, sub, err := _LPManagerBase.contract.WatchLogs(opts, "FeesCompounded", worknetTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LPManagerBaseFeesCompounded)
				if err := _LPManagerBase.contract.UnpackLog(event, "FeesCompounded", log); err != nil {
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

// ParseFeesCompounded is a log parse operation binding the contract event 0xbb85ff3bb01b1e5e9b75cd5e5ac87c83c843e76ca01a1565d7e050a4c8dbd223.
//
// Solidity: event FeesCompounded(address indexed worknetToken, uint256 tokenId)
func (_LPManagerBase *LPManagerBaseFilterer) ParseFeesCompounded(log types.Log) (*LPManagerBaseFeesCompounded, error) {
	event := new(LPManagerBaseFeesCompounded)
	if err := _LPManagerBase.contract.UnpackLog(event, "FeesCompounded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LPManagerBaseInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the LPManagerBase contract.
type LPManagerBaseInitializedIterator struct {
	Event *LPManagerBaseInitialized // Event containing the contract specifics and raw log

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
func (it *LPManagerBaseInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LPManagerBaseInitialized)
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
		it.Event = new(LPManagerBaseInitialized)
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
func (it *LPManagerBaseInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LPManagerBaseInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LPManagerBaseInitialized represents a Initialized event raised by the LPManagerBase contract.
type LPManagerBaseInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_LPManagerBase *LPManagerBaseFilterer) FilterInitialized(opts *bind.FilterOpts) (*LPManagerBaseInitializedIterator, error) {

	logs, sub, err := _LPManagerBase.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &LPManagerBaseInitializedIterator{contract: _LPManagerBase.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_LPManagerBase *LPManagerBaseFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *LPManagerBaseInitialized) (event.Subscription, error) {

	logs, sub, err := _LPManagerBase.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LPManagerBaseInitialized)
				if err := _LPManagerBase.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_LPManagerBase *LPManagerBaseFilterer) ParseInitialized(log types.Log) (*LPManagerBaseInitialized, error) {
	event := new(LPManagerBaseInitialized)
	if err := _LPManagerBase.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LPManagerBaseUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the LPManagerBase contract.
type LPManagerBaseUpgradedIterator struct {
	Event *LPManagerBaseUpgraded // Event containing the contract specifics and raw log

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
func (it *LPManagerBaseUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LPManagerBaseUpgraded)
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
		it.Event = new(LPManagerBaseUpgraded)
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
func (it *LPManagerBaseUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LPManagerBaseUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LPManagerBaseUpgraded represents a Upgraded event raised by the LPManagerBase contract.
type LPManagerBaseUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_LPManagerBase *LPManagerBaseFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*LPManagerBaseUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _LPManagerBase.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &LPManagerBaseUpgradedIterator{contract: _LPManagerBase.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_LPManagerBase *LPManagerBaseFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *LPManagerBaseUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _LPManagerBase.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LPManagerBaseUpgraded)
				if err := _LPManagerBase.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_LPManagerBase *LPManagerBaseFilterer) ParseUpgraded(log types.Log) (*LPManagerBaseUpgraded, error) {
	event := new(LPManagerBaseUpgraded)
	if err := _LPManagerBase.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
