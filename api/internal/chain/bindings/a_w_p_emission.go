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

// AWPEmissionMetaData contains all meta data concerning the AWPEmission contract.
var AWPEmissionMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DECAY_PRECISION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_DECAY_FACTOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activeEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"appendAllocations\",\"inputs\":[{\"name\":\"packed_\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"effectiveEpoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"awpToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIAWPToken\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"baseEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"baseTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"currentDailyEmission\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"currentEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decayFactor\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"epochDuration\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"epochEmissionLocked\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"frozenEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEpochRecipientCount\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEpochTotalWeight\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEpochWeight\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRecipient\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRecipientCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalWeight\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWeight\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"guardian\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"awpToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"guardian_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"initialDailyEmission_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"genesisTime_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"epochDuration_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"treasury_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"maxRecipients\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"modifyAllocations\",\"inputs\":[{\"name\":\"patches_\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"newTotalWeight_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"effectiveEpoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseEpochUntil\",\"inputs\":[{\"name\":\"resumeTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pausedUntil\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setDecayFactor\",\"inputs\":[{\"name\":\"newDecayFactor\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setEpochDuration\",\"inputs\":[{\"name\":\"newDuration\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setGuardian\",\"inputs\":[{\"name\":\"g\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMaxRecipients\",\"inputs\":[{\"name\":\"newMax\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTreasury\",\"inputs\":[{\"name\":\"t\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"settleEpoch\",\"inputs\":[{\"name\":\"limit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"settleProgress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"settledEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"submitAllocations\",\"inputs\":[{\"name\":\"packed_\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"totalWeight_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"effectiveEpoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"treasury\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AllocationsAppended\",\"inputs\":[{\"name\":\"effectiveEpoch\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packed\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllocationsModified\",\"inputs\":[{\"name\":\"effectiveEpoch\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"patches\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"newTotalWeight\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AllocationsSubmitted\",\"inputs\":[{\"name\":\"effectiveEpoch\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"packed\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"totalWeight\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DecayFactorUpdated\",\"inputs\":[{\"name\":\"newDecayFactor\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochDurationUpdated\",\"inputs\":[{\"name\":\"oldDuration\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newDuration\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochPausedUntil\",\"inputs\":[{\"name\":\"resumeTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"frozenEpoch\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochSettled\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"totalEmission\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"recipientCount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GuardianUpdated\",\"inputs\":[{\"name\":\"newGuardian\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MaxRecipientsUpdated\",\"inputs\":[{\"name\":\"newMax\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RecipientAWPDistributed\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"awpAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TreasuryUpdated\",\"inputs\":[{\"name\":\"newTreasury\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EpochNotReady\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GenesisNotReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"IndexOutOfBounds\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidDecayFactor\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidResumeTime\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MiningComplete\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustBeFutureEpoch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotGuardian\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SettlementInProgress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TooManyRecipients\",\"inputs\":[{\"name\":\"count\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"max\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroEpochDuration\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroWeight\",\"inputs\":[]}]",
}

// AWPEmissionABI is the input ABI used to generate the binding from.
// Deprecated: Use AWPEmissionMetaData.ABI instead.
var AWPEmissionABI = AWPEmissionMetaData.ABI

// AWPEmission is an auto generated Go binding around an Ethereum contract.
type AWPEmission struct {
	AWPEmissionCaller     // Read-only binding to the contract
	AWPEmissionTransactor // Write-only binding to the contract
	AWPEmissionFilterer   // Log filterer for contract events
}

// AWPEmissionCaller is an auto generated read-only Go binding around an Ethereum contract.
type AWPEmissionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AWPEmissionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AWPEmissionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AWPEmissionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AWPEmissionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AWPEmissionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AWPEmissionSession struct {
	Contract     *AWPEmission      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AWPEmissionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AWPEmissionCallerSession struct {
	Contract *AWPEmissionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// AWPEmissionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AWPEmissionTransactorSession struct {
	Contract     *AWPEmissionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// AWPEmissionRaw is an auto generated low-level Go binding around an Ethereum contract.
type AWPEmissionRaw struct {
	Contract *AWPEmission // Generic contract binding to access the raw methods on
}

// AWPEmissionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AWPEmissionCallerRaw struct {
	Contract *AWPEmissionCaller // Generic read-only contract binding to access the raw methods on
}

// AWPEmissionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AWPEmissionTransactorRaw struct {
	Contract *AWPEmissionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAWPEmission creates a new instance of AWPEmission, bound to a specific deployed contract.
func NewAWPEmission(address common.Address, backend bind.ContractBackend) (*AWPEmission, error) {
	contract, err := bindAWPEmission(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AWPEmission{AWPEmissionCaller: AWPEmissionCaller{contract: contract}, AWPEmissionTransactor: AWPEmissionTransactor{contract: contract}, AWPEmissionFilterer: AWPEmissionFilterer{contract: contract}}, nil
}

// NewAWPEmissionCaller creates a new read-only instance of AWPEmission, bound to a specific deployed contract.
func NewAWPEmissionCaller(address common.Address, caller bind.ContractCaller) (*AWPEmissionCaller, error) {
	contract, err := bindAWPEmission(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AWPEmissionCaller{contract: contract}, nil
}

// NewAWPEmissionTransactor creates a new write-only instance of AWPEmission, bound to a specific deployed contract.
func NewAWPEmissionTransactor(address common.Address, transactor bind.ContractTransactor) (*AWPEmissionTransactor, error) {
	contract, err := bindAWPEmission(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AWPEmissionTransactor{contract: contract}, nil
}

// NewAWPEmissionFilterer creates a new log filterer instance of AWPEmission, bound to a specific deployed contract.
func NewAWPEmissionFilterer(address common.Address, filterer bind.ContractFilterer) (*AWPEmissionFilterer, error) {
	contract, err := bindAWPEmission(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AWPEmissionFilterer{contract: contract}, nil
}

// bindAWPEmission binds a generic wrapper to an already deployed contract.
func bindAWPEmission(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AWPEmissionMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AWPEmission *AWPEmissionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AWPEmission.Contract.AWPEmissionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AWPEmission *AWPEmissionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AWPEmission.Contract.AWPEmissionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AWPEmission *AWPEmissionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AWPEmission.Contract.AWPEmissionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AWPEmission *AWPEmissionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AWPEmission.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AWPEmission *AWPEmissionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AWPEmission.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AWPEmission *AWPEmissionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AWPEmission.Contract.contract.Transact(opts, method, params...)
}

// DECAYPRECISION is a free data retrieval call binding the contract method 0xef80c9f0.
//
// Solidity: function DECAY_PRECISION() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) DECAYPRECISION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "DECAY_PRECISION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DECAYPRECISION is a free data retrieval call binding the contract method 0xef80c9f0.
//
// Solidity: function DECAY_PRECISION() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) DECAYPRECISION() (*big.Int, error) {
	return _AWPEmission.Contract.DECAYPRECISION(&_AWPEmission.CallOpts)
}

// DECAYPRECISION is a free data retrieval call binding the contract method 0xef80c9f0.
//
// Solidity: function DECAY_PRECISION() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) DECAYPRECISION() (*big.Int, error) {
	return _AWPEmission.Contract.DECAYPRECISION(&_AWPEmission.CallOpts)
}

// MINDECAYFACTOR is a free data retrieval call binding the contract method 0x82a28de3.
//
// Solidity: function MIN_DECAY_FACTOR() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) MINDECAYFACTOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "MIN_DECAY_FACTOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINDECAYFACTOR is a free data retrieval call binding the contract method 0x82a28de3.
//
// Solidity: function MIN_DECAY_FACTOR() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) MINDECAYFACTOR() (*big.Int, error) {
	return _AWPEmission.Contract.MINDECAYFACTOR(&_AWPEmission.CallOpts)
}

// MINDECAYFACTOR is a free data retrieval call binding the contract method 0x82a28de3.
//
// Solidity: function MIN_DECAY_FACTOR() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) MINDECAYFACTOR() (*big.Int, error) {
	return _AWPEmission.Contract.MINDECAYFACTOR(&_AWPEmission.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_AWPEmission *AWPEmissionCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_AWPEmission *AWPEmissionSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _AWPEmission.Contract.UPGRADEINTERFACEVERSION(&_AWPEmission.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_AWPEmission *AWPEmissionCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _AWPEmission.Contract.UPGRADEINTERFACEVERSION(&_AWPEmission.CallOpts)
}

// ActiveEpoch is a free data retrieval call binding the contract method 0x9f6b4a3b.
//
// Solidity: function activeEpoch() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) ActiveEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "activeEpoch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ActiveEpoch is a free data retrieval call binding the contract method 0x9f6b4a3b.
//
// Solidity: function activeEpoch() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) ActiveEpoch() (*big.Int, error) {
	return _AWPEmission.Contract.ActiveEpoch(&_AWPEmission.CallOpts)
}

// ActiveEpoch is a free data retrieval call binding the contract method 0x9f6b4a3b.
//
// Solidity: function activeEpoch() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) ActiveEpoch() (*big.Int, error) {
	return _AWPEmission.Contract.ActiveEpoch(&_AWPEmission.CallOpts)
}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_AWPEmission *AWPEmissionCaller) AwpToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "awpToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_AWPEmission *AWPEmissionSession) AwpToken() (common.Address, error) {
	return _AWPEmission.Contract.AwpToken(&_AWPEmission.CallOpts)
}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_AWPEmission *AWPEmissionCallerSession) AwpToken() (common.Address, error) {
	return _AWPEmission.Contract.AwpToken(&_AWPEmission.CallOpts)
}

// BaseEpoch is a free data retrieval call binding the contract method 0x67baf995.
//
// Solidity: function baseEpoch() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) BaseEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "baseEpoch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BaseEpoch is a free data retrieval call binding the contract method 0x67baf995.
//
// Solidity: function baseEpoch() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) BaseEpoch() (*big.Int, error) {
	return _AWPEmission.Contract.BaseEpoch(&_AWPEmission.CallOpts)
}

// BaseEpoch is a free data retrieval call binding the contract method 0x67baf995.
//
// Solidity: function baseEpoch() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) BaseEpoch() (*big.Int, error) {
	return _AWPEmission.Contract.BaseEpoch(&_AWPEmission.CallOpts)
}

// BaseTime is a free data retrieval call binding the contract method 0x08b096a0.
//
// Solidity: function baseTime() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) BaseTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "baseTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BaseTime is a free data retrieval call binding the contract method 0x08b096a0.
//
// Solidity: function baseTime() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) BaseTime() (*big.Int, error) {
	return _AWPEmission.Contract.BaseTime(&_AWPEmission.CallOpts)
}

// BaseTime is a free data retrieval call binding the contract method 0x08b096a0.
//
// Solidity: function baseTime() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) BaseTime() (*big.Int, error) {
	return _AWPEmission.Contract.BaseTime(&_AWPEmission.CallOpts)
}

// CurrentDailyEmission is a free data retrieval call binding the contract method 0x091075a9.
//
// Solidity: function currentDailyEmission() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) CurrentDailyEmission(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "currentDailyEmission")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentDailyEmission is a free data retrieval call binding the contract method 0x091075a9.
//
// Solidity: function currentDailyEmission() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) CurrentDailyEmission() (*big.Int, error) {
	return _AWPEmission.Contract.CurrentDailyEmission(&_AWPEmission.CallOpts)
}

// CurrentDailyEmission is a free data retrieval call binding the contract method 0x091075a9.
//
// Solidity: function currentDailyEmission() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) CurrentDailyEmission() (*big.Int, error) {
	return _AWPEmission.Contract.CurrentDailyEmission(&_AWPEmission.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) CurrentEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "currentEpoch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) CurrentEpoch() (*big.Int, error) {
	return _AWPEmission.Contract.CurrentEpoch(&_AWPEmission.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) CurrentEpoch() (*big.Int, error) {
	return _AWPEmission.Contract.CurrentEpoch(&_AWPEmission.CallOpts)
}

// DecayFactor is a free data retrieval call binding the contract method 0x20fb3016.
//
// Solidity: function decayFactor() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) DecayFactor(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "decayFactor")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DecayFactor is a free data retrieval call binding the contract method 0x20fb3016.
//
// Solidity: function decayFactor() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) DecayFactor() (*big.Int, error) {
	return _AWPEmission.Contract.DecayFactor(&_AWPEmission.CallOpts)
}

// DecayFactor is a free data retrieval call binding the contract method 0x20fb3016.
//
// Solidity: function decayFactor() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) DecayFactor() (*big.Int, error) {
	return _AWPEmission.Contract.DecayFactor(&_AWPEmission.CallOpts)
}

// EpochDuration is a free data retrieval call binding the contract method 0x4ff0876a.
//
// Solidity: function epochDuration() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) EpochDuration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "epochDuration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EpochDuration is a free data retrieval call binding the contract method 0x4ff0876a.
//
// Solidity: function epochDuration() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) EpochDuration() (*big.Int, error) {
	return _AWPEmission.Contract.EpochDuration(&_AWPEmission.CallOpts)
}

// EpochDuration is a free data retrieval call binding the contract method 0x4ff0876a.
//
// Solidity: function epochDuration() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) EpochDuration() (*big.Int, error) {
	return _AWPEmission.Contract.EpochDuration(&_AWPEmission.CallOpts)
}

// EpochEmissionLocked is a free data retrieval call binding the contract method 0x0c86ccd8.
//
// Solidity: function epochEmissionLocked() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) EpochEmissionLocked(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "epochEmissionLocked")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EpochEmissionLocked is a free data retrieval call binding the contract method 0x0c86ccd8.
//
// Solidity: function epochEmissionLocked() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) EpochEmissionLocked() (*big.Int, error) {
	return _AWPEmission.Contract.EpochEmissionLocked(&_AWPEmission.CallOpts)
}

// EpochEmissionLocked is a free data retrieval call binding the contract method 0x0c86ccd8.
//
// Solidity: function epochEmissionLocked() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) EpochEmissionLocked() (*big.Int, error) {
	return _AWPEmission.Contract.EpochEmissionLocked(&_AWPEmission.CallOpts)
}

// FrozenEpoch is a free data retrieval call binding the contract method 0x585db72a.
//
// Solidity: function frozenEpoch() view returns(uint64)
func (_AWPEmission *AWPEmissionCaller) FrozenEpoch(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "frozenEpoch")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// FrozenEpoch is a free data retrieval call binding the contract method 0x585db72a.
//
// Solidity: function frozenEpoch() view returns(uint64)
func (_AWPEmission *AWPEmissionSession) FrozenEpoch() (uint64, error) {
	return _AWPEmission.Contract.FrozenEpoch(&_AWPEmission.CallOpts)
}

// FrozenEpoch is a free data retrieval call binding the contract method 0x585db72a.
//
// Solidity: function frozenEpoch() view returns(uint64)
func (_AWPEmission *AWPEmissionCallerSession) FrozenEpoch() (uint64, error) {
	return _AWPEmission.Contract.FrozenEpoch(&_AWPEmission.CallOpts)
}

// GetEpochRecipientCount is a free data retrieval call binding the contract method 0x7b2c32fc.
//
// Solidity: function getEpochRecipientCount(uint256 epoch) view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) GetEpochRecipientCount(opts *bind.CallOpts, epoch *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "getEpochRecipientCount", epoch)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEpochRecipientCount is a free data retrieval call binding the contract method 0x7b2c32fc.
//
// Solidity: function getEpochRecipientCount(uint256 epoch) view returns(uint256)
func (_AWPEmission *AWPEmissionSession) GetEpochRecipientCount(epoch *big.Int) (*big.Int, error) {
	return _AWPEmission.Contract.GetEpochRecipientCount(&_AWPEmission.CallOpts, epoch)
}

// GetEpochRecipientCount is a free data retrieval call binding the contract method 0x7b2c32fc.
//
// Solidity: function getEpochRecipientCount(uint256 epoch) view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) GetEpochRecipientCount(epoch *big.Int) (*big.Int, error) {
	return _AWPEmission.Contract.GetEpochRecipientCount(&_AWPEmission.CallOpts, epoch)
}

// GetEpochTotalWeight is a free data retrieval call binding the contract method 0x3a895f5e.
//
// Solidity: function getEpochTotalWeight(uint256 epoch) view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) GetEpochTotalWeight(opts *bind.CallOpts, epoch *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "getEpochTotalWeight", epoch)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEpochTotalWeight is a free data retrieval call binding the contract method 0x3a895f5e.
//
// Solidity: function getEpochTotalWeight(uint256 epoch) view returns(uint256)
func (_AWPEmission *AWPEmissionSession) GetEpochTotalWeight(epoch *big.Int) (*big.Int, error) {
	return _AWPEmission.Contract.GetEpochTotalWeight(&_AWPEmission.CallOpts, epoch)
}

// GetEpochTotalWeight is a free data retrieval call binding the contract method 0x3a895f5e.
//
// Solidity: function getEpochTotalWeight(uint256 epoch) view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) GetEpochTotalWeight(epoch *big.Int) (*big.Int, error) {
	return _AWPEmission.Contract.GetEpochTotalWeight(&_AWPEmission.CallOpts, epoch)
}

// GetEpochWeight is a free data retrieval call binding the contract method 0xb1b381a2.
//
// Solidity: function getEpochWeight(uint256 epoch, address addr) view returns(uint96)
func (_AWPEmission *AWPEmissionCaller) GetEpochWeight(opts *bind.CallOpts, epoch *big.Int, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "getEpochWeight", epoch, addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEpochWeight is a free data retrieval call binding the contract method 0xb1b381a2.
//
// Solidity: function getEpochWeight(uint256 epoch, address addr) view returns(uint96)
func (_AWPEmission *AWPEmissionSession) GetEpochWeight(epoch *big.Int, addr common.Address) (*big.Int, error) {
	return _AWPEmission.Contract.GetEpochWeight(&_AWPEmission.CallOpts, epoch, addr)
}

// GetEpochWeight is a free data retrieval call binding the contract method 0xb1b381a2.
//
// Solidity: function getEpochWeight(uint256 epoch, address addr) view returns(uint96)
func (_AWPEmission *AWPEmissionCallerSession) GetEpochWeight(epoch *big.Int, addr common.Address) (*big.Int, error) {
	return _AWPEmission.Contract.GetEpochWeight(&_AWPEmission.CallOpts, epoch, addr)
}

// GetRecipient is a free data retrieval call binding the contract method 0x6d0cee75.
//
// Solidity: function getRecipient(uint256 index) view returns(address)
func (_AWPEmission *AWPEmissionCaller) GetRecipient(opts *bind.CallOpts, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "getRecipient", index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRecipient is a free data retrieval call binding the contract method 0x6d0cee75.
//
// Solidity: function getRecipient(uint256 index) view returns(address)
func (_AWPEmission *AWPEmissionSession) GetRecipient(index *big.Int) (common.Address, error) {
	return _AWPEmission.Contract.GetRecipient(&_AWPEmission.CallOpts, index)
}

// GetRecipient is a free data retrieval call binding the contract method 0x6d0cee75.
//
// Solidity: function getRecipient(uint256 index) view returns(address)
func (_AWPEmission *AWPEmissionCallerSession) GetRecipient(index *big.Int) (common.Address, error) {
	return _AWPEmission.Contract.GetRecipient(&_AWPEmission.CallOpts, index)
}

// GetRecipientCount is a free data retrieval call binding the contract method 0xaf99b63f.
//
// Solidity: function getRecipientCount() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) GetRecipientCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "getRecipientCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRecipientCount is a free data retrieval call binding the contract method 0xaf99b63f.
//
// Solidity: function getRecipientCount() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) GetRecipientCount() (*big.Int, error) {
	return _AWPEmission.Contract.GetRecipientCount(&_AWPEmission.CallOpts)
}

// GetRecipientCount is a free data retrieval call binding the contract method 0xaf99b63f.
//
// Solidity: function getRecipientCount() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) GetRecipientCount() (*big.Int, error) {
	return _AWPEmission.Contract.GetRecipientCount(&_AWPEmission.CallOpts)
}

// GetTotalWeight is a free data retrieval call binding the contract method 0x06aba0e1.
//
// Solidity: function getTotalWeight() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) GetTotalWeight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "getTotalWeight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalWeight is a free data retrieval call binding the contract method 0x06aba0e1.
//
// Solidity: function getTotalWeight() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) GetTotalWeight() (*big.Int, error) {
	return _AWPEmission.Contract.GetTotalWeight(&_AWPEmission.CallOpts)
}

// GetTotalWeight is a free data retrieval call binding the contract method 0x06aba0e1.
//
// Solidity: function getTotalWeight() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) GetTotalWeight() (*big.Int, error) {
	return _AWPEmission.Contract.GetTotalWeight(&_AWPEmission.CallOpts)
}

// GetWeight is a free data retrieval call binding the contract method 0xac6c5251.
//
// Solidity: function getWeight(address addr) view returns(uint96)
func (_AWPEmission *AWPEmissionCaller) GetWeight(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "getWeight", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetWeight is a free data retrieval call binding the contract method 0xac6c5251.
//
// Solidity: function getWeight(address addr) view returns(uint96)
func (_AWPEmission *AWPEmissionSession) GetWeight(addr common.Address) (*big.Int, error) {
	return _AWPEmission.Contract.GetWeight(&_AWPEmission.CallOpts, addr)
}

// GetWeight is a free data retrieval call binding the contract method 0xac6c5251.
//
// Solidity: function getWeight(address addr) view returns(uint96)
func (_AWPEmission *AWPEmissionCallerSession) GetWeight(addr common.Address) (*big.Int, error) {
	return _AWPEmission.Contract.GetWeight(&_AWPEmission.CallOpts, addr)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_AWPEmission *AWPEmissionCaller) Guardian(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "guardian")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_AWPEmission *AWPEmissionSession) Guardian() (common.Address, error) {
	return _AWPEmission.Contract.Guardian(&_AWPEmission.CallOpts)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_AWPEmission *AWPEmissionCallerSession) Guardian() (common.Address, error) {
	return _AWPEmission.Contract.Guardian(&_AWPEmission.CallOpts)
}

// MaxRecipients is a free data retrieval call binding the contract method 0x88a13072.
//
// Solidity: function maxRecipients() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) MaxRecipients(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "maxRecipients")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxRecipients is a free data retrieval call binding the contract method 0x88a13072.
//
// Solidity: function maxRecipients() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) MaxRecipients() (*big.Int, error) {
	return _AWPEmission.Contract.MaxRecipients(&_AWPEmission.CallOpts)
}

// MaxRecipients is a free data retrieval call binding the contract method 0x88a13072.
//
// Solidity: function maxRecipients() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) MaxRecipients() (*big.Int, error) {
	return _AWPEmission.Contract.MaxRecipients(&_AWPEmission.CallOpts)
}

// PausedUntil is a free data retrieval call binding the contract method 0xda748b10.
//
// Solidity: function pausedUntil() view returns(uint64)
func (_AWPEmission *AWPEmissionCaller) PausedUntil(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "pausedUntil")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// PausedUntil is a free data retrieval call binding the contract method 0xda748b10.
//
// Solidity: function pausedUntil() view returns(uint64)
func (_AWPEmission *AWPEmissionSession) PausedUntil() (uint64, error) {
	return _AWPEmission.Contract.PausedUntil(&_AWPEmission.CallOpts)
}

// PausedUntil is a free data retrieval call binding the contract method 0xda748b10.
//
// Solidity: function pausedUntil() view returns(uint64)
func (_AWPEmission *AWPEmissionCallerSession) PausedUntil() (uint64, error) {
	return _AWPEmission.Contract.PausedUntil(&_AWPEmission.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_AWPEmission *AWPEmissionCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_AWPEmission *AWPEmissionSession) ProxiableUUID() ([32]byte, error) {
	return _AWPEmission.Contract.ProxiableUUID(&_AWPEmission.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_AWPEmission *AWPEmissionCallerSession) ProxiableUUID() ([32]byte, error) {
	return _AWPEmission.Contract.ProxiableUUID(&_AWPEmission.CallOpts)
}

// SettleProgress is a free data retrieval call binding the contract method 0x57200fc5.
//
// Solidity: function settleProgress() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) SettleProgress(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "settleProgress")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SettleProgress is a free data retrieval call binding the contract method 0x57200fc5.
//
// Solidity: function settleProgress() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) SettleProgress() (*big.Int, error) {
	return _AWPEmission.Contract.SettleProgress(&_AWPEmission.CallOpts)
}

// SettleProgress is a free data retrieval call binding the contract method 0x57200fc5.
//
// Solidity: function settleProgress() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) SettleProgress() (*big.Int, error) {
	return _AWPEmission.Contract.SettleProgress(&_AWPEmission.CallOpts)
}

// SettledEpoch is a free data retrieval call binding the contract method 0xb560d209.
//
// Solidity: function settledEpoch() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) SettledEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "settledEpoch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SettledEpoch is a free data retrieval call binding the contract method 0xb560d209.
//
// Solidity: function settledEpoch() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) SettledEpoch() (*big.Int, error) {
	return _AWPEmission.Contract.SettledEpoch(&_AWPEmission.CallOpts)
}

// SettledEpoch is a free data retrieval call binding the contract method 0xb560d209.
//
// Solidity: function settledEpoch() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) SettledEpoch() (*big.Int, error) {
	return _AWPEmission.Contract.SettledEpoch(&_AWPEmission.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_AWPEmission *AWPEmissionCaller) Treasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "treasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_AWPEmission *AWPEmissionSession) Treasury() (common.Address, error) {
	return _AWPEmission.Contract.Treasury(&_AWPEmission.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_AWPEmission *AWPEmissionCallerSession) Treasury() (common.Address, error) {
	return _AWPEmission.Contract.Treasury(&_AWPEmission.CallOpts)
}

// AppendAllocations is a paid mutator transaction binding the contract method 0xca2b5c60.
//
// Solidity: function appendAllocations(uint256[] packed_, uint256 effectiveEpoch) returns()
func (_AWPEmission *AWPEmissionTransactor) AppendAllocations(opts *bind.TransactOpts, packed_ []*big.Int, effectiveEpoch *big.Int) (*types.Transaction, error) {
	return _AWPEmission.contract.Transact(opts, "appendAllocations", packed_, effectiveEpoch)
}

// AppendAllocations is a paid mutator transaction binding the contract method 0xca2b5c60.
//
// Solidity: function appendAllocations(uint256[] packed_, uint256 effectiveEpoch) returns()
func (_AWPEmission *AWPEmissionSession) AppendAllocations(packed_ []*big.Int, effectiveEpoch *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.AppendAllocations(&_AWPEmission.TransactOpts, packed_, effectiveEpoch)
}

// AppendAllocations is a paid mutator transaction binding the contract method 0xca2b5c60.
//
// Solidity: function appendAllocations(uint256[] packed_, uint256 effectiveEpoch) returns()
func (_AWPEmission *AWPEmissionTransactorSession) AppendAllocations(packed_ []*big.Int, effectiveEpoch *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.AppendAllocations(&_AWPEmission.TransactOpts, packed_, effectiveEpoch)
}

// Initialize is a paid mutator transaction binding the contract method 0xb1a5d12d.
//
// Solidity: function initialize(address awpToken_, address guardian_, uint256 initialDailyEmission_, uint256 genesisTime_, uint256 epochDuration_, address treasury_) returns()
func (_AWPEmission *AWPEmissionTransactor) Initialize(opts *bind.TransactOpts, awpToken_ common.Address, guardian_ common.Address, initialDailyEmission_ *big.Int, genesisTime_ *big.Int, epochDuration_ *big.Int, treasury_ common.Address) (*types.Transaction, error) {
	return _AWPEmission.contract.Transact(opts, "initialize", awpToken_, guardian_, initialDailyEmission_, genesisTime_, epochDuration_, treasury_)
}

// Initialize is a paid mutator transaction binding the contract method 0xb1a5d12d.
//
// Solidity: function initialize(address awpToken_, address guardian_, uint256 initialDailyEmission_, uint256 genesisTime_, uint256 epochDuration_, address treasury_) returns()
func (_AWPEmission *AWPEmissionSession) Initialize(awpToken_ common.Address, guardian_ common.Address, initialDailyEmission_ *big.Int, genesisTime_ *big.Int, epochDuration_ *big.Int, treasury_ common.Address) (*types.Transaction, error) {
	return _AWPEmission.Contract.Initialize(&_AWPEmission.TransactOpts, awpToken_, guardian_, initialDailyEmission_, genesisTime_, epochDuration_, treasury_)
}

// Initialize is a paid mutator transaction binding the contract method 0xb1a5d12d.
//
// Solidity: function initialize(address awpToken_, address guardian_, uint256 initialDailyEmission_, uint256 genesisTime_, uint256 epochDuration_, address treasury_) returns()
func (_AWPEmission *AWPEmissionTransactorSession) Initialize(awpToken_ common.Address, guardian_ common.Address, initialDailyEmission_ *big.Int, genesisTime_ *big.Int, epochDuration_ *big.Int, treasury_ common.Address) (*types.Transaction, error) {
	return _AWPEmission.Contract.Initialize(&_AWPEmission.TransactOpts, awpToken_, guardian_, initialDailyEmission_, genesisTime_, epochDuration_, treasury_)
}

// ModifyAllocations is a paid mutator transaction binding the contract method 0x169ec78d.
//
// Solidity: function modifyAllocations(uint256[] patches_, uint256 newTotalWeight_, uint256 effectiveEpoch) returns()
func (_AWPEmission *AWPEmissionTransactor) ModifyAllocations(opts *bind.TransactOpts, patches_ []*big.Int, newTotalWeight_ *big.Int, effectiveEpoch *big.Int) (*types.Transaction, error) {
	return _AWPEmission.contract.Transact(opts, "modifyAllocations", patches_, newTotalWeight_, effectiveEpoch)
}

// ModifyAllocations is a paid mutator transaction binding the contract method 0x169ec78d.
//
// Solidity: function modifyAllocations(uint256[] patches_, uint256 newTotalWeight_, uint256 effectiveEpoch) returns()
func (_AWPEmission *AWPEmissionSession) ModifyAllocations(patches_ []*big.Int, newTotalWeight_ *big.Int, effectiveEpoch *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.ModifyAllocations(&_AWPEmission.TransactOpts, patches_, newTotalWeight_, effectiveEpoch)
}

// ModifyAllocations is a paid mutator transaction binding the contract method 0x169ec78d.
//
// Solidity: function modifyAllocations(uint256[] patches_, uint256 newTotalWeight_, uint256 effectiveEpoch) returns()
func (_AWPEmission *AWPEmissionTransactorSession) ModifyAllocations(patches_ []*big.Int, newTotalWeight_ *big.Int, effectiveEpoch *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.ModifyAllocations(&_AWPEmission.TransactOpts, patches_, newTotalWeight_, effectiveEpoch)
}

// PauseEpochUntil is a paid mutator transaction binding the contract method 0x108897cd.
//
// Solidity: function pauseEpochUntil(uint64 resumeTime) returns()
func (_AWPEmission *AWPEmissionTransactor) PauseEpochUntil(opts *bind.TransactOpts, resumeTime uint64) (*types.Transaction, error) {
	return _AWPEmission.contract.Transact(opts, "pauseEpochUntil", resumeTime)
}

// PauseEpochUntil is a paid mutator transaction binding the contract method 0x108897cd.
//
// Solidity: function pauseEpochUntil(uint64 resumeTime) returns()
func (_AWPEmission *AWPEmissionSession) PauseEpochUntil(resumeTime uint64) (*types.Transaction, error) {
	return _AWPEmission.Contract.PauseEpochUntil(&_AWPEmission.TransactOpts, resumeTime)
}

// PauseEpochUntil is a paid mutator transaction binding the contract method 0x108897cd.
//
// Solidity: function pauseEpochUntil(uint64 resumeTime) returns()
func (_AWPEmission *AWPEmissionTransactorSession) PauseEpochUntil(resumeTime uint64) (*types.Transaction, error) {
	return _AWPEmission.Contract.PauseEpochUntil(&_AWPEmission.TransactOpts, resumeTime)
}

// SetDecayFactor is a paid mutator transaction binding the contract method 0xb8c9059d.
//
// Solidity: function setDecayFactor(uint256 newDecayFactor) returns()
func (_AWPEmission *AWPEmissionTransactor) SetDecayFactor(opts *bind.TransactOpts, newDecayFactor *big.Int) (*types.Transaction, error) {
	return _AWPEmission.contract.Transact(opts, "setDecayFactor", newDecayFactor)
}

// SetDecayFactor is a paid mutator transaction binding the contract method 0xb8c9059d.
//
// Solidity: function setDecayFactor(uint256 newDecayFactor) returns()
func (_AWPEmission *AWPEmissionSession) SetDecayFactor(newDecayFactor *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.SetDecayFactor(&_AWPEmission.TransactOpts, newDecayFactor)
}

// SetDecayFactor is a paid mutator transaction binding the contract method 0xb8c9059d.
//
// Solidity: function setDecayFactor(uint256 newDecayFactor) returns()
func (_AWPEmission *AWPEmissionTransactorSession) SetDecayFactor(newDecayFactor *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.SetDecayFactor(&_AWPEmission.TransactOpts, newDecayFactor)
}

// SetEpochDuration is a paid mutator transaction binding the contract method 0x30024dfe.
//
// Solidity: function setEpochDuration(uint256 newDuration) returns()
func (_AWPEmission *AWPEmissionTransactor) SetEpochDuration(opts *bind.TransactOpts, newDuration *big.Int) (*types.Transaction, error) {
	return _AWPEmission.contract.Transact(opts, "setEpochDuration", newDuration)
}

// SetEpochDuration is a paid mutator transaction binding the contract method 0x30024dfe.
//
// Solidity: function setEpochDuration(uint256 newDuration) returns()
func (_AWPEmission *AWPEmissionSession) SetEpochDuration(newDuration *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.SetEpochDuration(&_AWPEmission.TransactOpts, newDuration)
}

// SetEpochDuration is a paid mutator transaction binding the contract method 0x30024dfe.
//
// Solidity: function setEpochDuration(uint256 newDuration) returns()
func (_AWPEmission *AWPEmissionTransactorSession) SetEpochDuration(newDuration *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.SetEpochDuration(&_AWPEmission.TransactOpts, newDuration)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_AWPEmission *AWPEmissionTransactor) SetGuardian(opts *bind.TransactOpts, g common.Address) (*types.Transaction, error) {
	return _AWPEmission.contract.Transact(opts, "setGuardian", g)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_AWPEmission *AWPEmissionSession) SetGuardian(g common.Address) (*types.Transaction, error) {
	return _AWPEmission.Contract.SetGuardian(&_AWPEmission.TransactOpts, g)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_AWPEmission *AWPEmissionTransactorSession) SetGuardian(g common.Address) (*types.Transaction, error) {
	return _AWPEmission.Contract.SetGuardian(&_AWPEmission.TransactOpts, g)
}

// SetMaxRecipients is a paid mutator transaction binding the contract method 0x148ec9ab.
//
// Solidity: function setMaxRecipients(uint256 newMax) returns()
func (_AWPEmission *AWPEmissionTransactor) SetMaxRecipients(opts *bind.TransactOpts, newMax *big.Int) (*types.Transaction, error) {
	return _AWPEmission.contract.Transact(opts, "setMaxRecipients", newMax)
}

// SetMaxRecipients is a paid mutator transaction binding the contract method 0x148ec9ab.
//
// Solidity: function setMaxRecipients(uint256 newMax) returns()
func (_AWPEmission *AWPEmissionSession) SetMaxRecipients(newMax *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.SetMaxRecipients(&_AWPEmission.TransactOpts, newMax)
}

// SetMaxRecipients is a paid mutator transaction binding the contract method 0x148ec9ab.
//
// Solidity: function setMaxRecipients(uint256 newMax) returns()
func (_AWPEmission *AWPEmissionTransactorSession) SetMaxRecipients(newMax *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.SetMaxRecipients(&_AWPEmission.TransactOpts, newMax)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address t) returns()
func (_AWPEmission *AWPEmissionTransactor) SetTreasury(opts *bind.TransactOpts, t common.Address) (*types.Transaction, error) {
	return _AWPEmission.contract.Transact(opts, "setTreasury", t)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address t) returns()
func (_AWPEmission *AWPEmissionSession) SetTreasury(t common.Address) (*types.Transaction, error) {
	return _AWPEmission.Contract.SetTreasury(&_AWPEmission.TransactOpts, t)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address t) returns()
func (_AWPEmission *AWPEmissionTransactorSession) SetTreasury(t common.Address) (*types.Transaction, error) {
	return _AWPEmission.Contract.SetTreasury(&_AWPEmission.TransactOpts, t)
}

// SettleEpoch is a paid mutator transaction binding the contract method 0x4d35fa7e.
//
// Solidity: function settleEpoch(uint256 limit) returns()
func (_AWPEmission *AWPEmissionTransactor) SettleEpoch(opts *bind.TransactOpts, limit *big.Int) (*types.Transaction, error) {
	return _AWPEmission.contract.Transact(opts, "settleEpoch", limit)
}

// SettleEpoch is a paid mutator transaction binding the contract method 0x4d35fa7e.
//
// Solidity: function settleEpoch(uint256 limit) returns()
func (_AWPEmission *AWPEmissionSession) SettleEpoch(limit *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.SettleEpoch(&_AWPEmission.TransactOpts, limit)
}

// SettleEpoch is a paid mutator transaction binding the contract method 0x4d35fa7e.
//
// Solidity: function settleEpoch(uint256 limit) returns()
func (_AWPEmission *AWPEmissionTransactorSession) SettleEpoch(limit *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.SettleEpoch(&_AWPEmission.TransactOpts, limit)
}

// SubmitAllocations is a paid mutator transaction binding the contract method 0x0da70833.
//
// Solidity: function submitAllocations(uint256[] packed_, uint256 totalWeight_, uint256 effectiveEpoch) returns()
func (_AWPEmission *AWPEmissionTransactor) SubmitAllocations(opts *bind.TransactOpts, packed_ []*big.Int, totalWeight_ *big.Int, effectiveEpoch *big.Int) (*types.Transaction, error) {
	return _AWPEmission.contract.Transact(opts, "submitAllocations", packed_, totalWeight_, effectiveEpoch)
}

// SubmitAllocations is a paid mutator transaction binding the contract method 0x0da70833.
//
// Solidity: function submitAllocations(uint256[] packed_, uint256 totalWeight_, uint256 effectiveEpoch) returns()
func (_AWPEmission *AWPEmissionSession) SubmitAllocations(packed_ []*big.Int, totalWeight_ *big.Int, effectiveEpoch *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.SubmitAllocations(&_AWPEmission.TransactOpts, packed_, totalWeight_, effectiveEpoch)
}

// SubmitAllocations is a paid mutator transaction binding the contract method 0x0da70833.
//
// Solidity: function submitAllocations(uint256[] packed_, uint256 totalWeight_, uint256 effectiveEpoch) returns()
func (_AWPEmission *AWPEmissionTransactorSession) SubmitAllocations(packed_ []*big.Int, totalWeight_ *big.Int, effectiveEpoch *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.SubmitAllocations(&_AWPEmission.TransactOpts, packed_, totalWeight_, effectiveEpoch)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AWPEmission *AWPEmissionTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AWPEmission.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AWPEmission *AWPEmissionSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AWPEmission.Contract.UpgradeToAndCall(&_AWPEmission.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AWPEmission *AWPEmissionTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AWPEmission.Contract.UpgradeToAndCall(&_AWPEmission.TransactOpts, newImplementation, data)
}

// AWPEmissionAllocationsAppendedIterator is returned from FilterAllocationsAppended and is used to iterate over the raw logs and unpacked data for AllocationsAppended events raised by the AWPEmission contract.
type AWPEmissionAllocationsAppendedIterator struct {
	Event *AWPEmissionAllocationsAppended // Event containing the contract specifics and raw log

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
func (it *AWPEmissionAllocationsAppendedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPEmissionAllocationsAppended)
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
		it.Event = new(AWPEmissionAllocationsAppended)
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
func (it *AWPEmissionAllocationsAppendedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPEmissionAllocationsAppendedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPEmissionAllocationsAppended represents a AllocationsAppended event raised by the AWPEmission contract.
type AWPEmissionAllocationsAppended struct {
	EffectiveEpoch *big.Int
	Packed         []*big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterAllocationsAppended is a free log retrieval operation binding the contract event 0xced6171830a1afd9dec6f4cb9dcd9ddde35d7bbd1bfff440f0c5b49c9cb01bd5.
//
// Solidity: event AllocationsAppended(uint256 indexed effectiveEpoch, uint256[] packed)
func (_AWPEmission *AWPEmissionFilterer) FilterAllocationsAppended(opts *bind.FilterOpts, effectiveEpoch []*big.Int) (*AWPEmissionAllocationsAppendedIterator, error) {

	var effectiveEpochRule []interface{}
	for _, effectiveEpochItem := range effectiveEpoch {
		effectiveEpochRule = append(effectiveEpochRule, effectiveEpochItem)
	}

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "AllocationsAppended", effectiveEpochRule)
	if err != nil {
		return nil, err
	}
	return &AWPEmissionAllocationsAppendedIterator{contract: _AWPEmission.contract, event: "AllocationsAppended", logs: logs, sub: sub}, nil
}

// WatchAllocationsAppended is a free log subscription operation binding the contract event 0xced6171830a1afd9dec6f4cb9dcd9ddde35d7bbd1bfff440f0c5b49c9cb01bd5.
//
// Solidity: event AllocationsAppended(uint256 indexed effectiveEpoch, uint256[] packed)
func (_AWPEmission *AWPEmissionFilterer) WatchAllocationsAppended(opts *bind.WatchOpts, sink chan<- *AWPEmissionAllocationsAppended, effectiveEpoch []*big.Int) (event.Subscription, error) {

	var effectiveEpochRule []interface{}
	for _, effectiveEpochItem := range effectiveEpoch {
		effectiveEpochRule = append(effectiveEpochRule, effectiveEpochItem)
	}

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "AllocationsAppended", effectiveEpochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPEmissionAllocationsAppended)
				if err := _AWPEmission.contract.UnpackLog(event, "AllocationsAppended", log); err != nil {
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

// ParseAllocationsAppended is a log parse operation binding the contract event 0xced6171830a1afd9dec6f4cb9dcd9ddde35d7bbd1bfff440f0c5b49c9cb01bd5.
//
// Solidity: event AllocationsAppended(uint256 indexed effectiveEpoch, uint256[] packed)
func (_AWPEmission *AWPEmissionFilterer) ParseAllocationsAppended(log types.Log) (*AWPEmissionAllocationsAppended, error) {
	event := new(AWPEmissionAllocationsAppended)
	if err := _AWPEmission.contract.UnpackLog(event, "AllocationsAppended", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPEmissionAllocationsModifiedIterator is returned from FilterAllocationsModified and is used to iterate over the raw logs and unpacked data for AllocationsModified events raised by the AWPEmission contract.
type AWPEmissionAllocationsModifiedIterator struct {
	Event *AWPEmissionAllocationsModified // Event containing the contract specifics and raw log

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
func (it *AWPEmissionAllocationsModifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPEmissionAllocationsModified)
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
		it.Event = new(AWPEmissionAllocationsModified)
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
func (it *AWPEmissionAllocationsModifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPEmissionAllocationsModifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPEmissionAllocationsModified represents a AllocationsModified event raised by the AWPEmission contract.
type AWPEmissionAllocationsModified struct {
	EffectiveEpoch *big.Int
	Patches        []*big.Int
	NewTotalWeight *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterAllocationsModified is a free log retrieval operation binding the contract event 0x23937aac5de5861678c5cbc4fef6ffe02ccda8484361ec9772d729d959930783.
//
// Solidity: event AllocationsModified(uint256 indexed effectiveEpoch, uint256[] patches, uint256 newTotalWeight)
func (_AWPEmission *AWPEmissionFilterer) FilterAllocationsModified(opts *bind.FilterOpts, effectiveEpoch []*big.Int) (*AWPEmissionAllocationsModifiedIterator, error) {

	var effectiveEpochRule []interface{}
	for _, effectiveEpochItem := range effectiveEpoch {
		effectiveEpochRule = append(effectiveEpochRule, effectiveEpochItem)
	}

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "AllocationsModified", effectiveEpochRule)
	if err != nil {
		return nil, err
	}
	return &AWPEmissionAllocationsModifiedIterator{contract: _AWPEmission.contract, event: "AllocationsModified", logs: logs, sub: sub}, nil
}

// WatchAllocationsModified is a free log subscription operation binding the contract event 0x23937aac5de5861678c5cbc4fef6ffe02ccda8484361ec9772d729d959930783.
//
// Solidity: event AllocationsModified(uint256 indexed effectiveEpoch, uint256[] patches, uint256 newTotalWeight)
func (_AWPEmission *AWPEmissionFilterer) WatchAllocationsModified(opts *bind.WatchOpts, sink chan<- *AWPEmissionAllocationsModified, effectiveEpoch []*big.Int) (event.Subscription, error) {

	var effectiveEpochRule []interface{}
	for _, effectiveEpochItem := range effectiveEpoch {
		effectiveEpochRule = append(effectiveEpochRule, effectiveEpochItem)
	}

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "AllocationsModified", effectiveEpochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPEmissionAllocationsModified)
				if err := _AWPEmission.contract.UnpackLog(event, "AllocationsModified", log); err != nil {
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

// ParseAllocationsModified is a log parse operation binding the contract event 0x23937aac5de5861678c5cbc4fef6ffe02ccda8484361ec9772d729d959930783.
//
// Solidity: event AllocationsModified(uint256 indexed effectiveEpoch, uint256[] patches, uint256 newTotalWeight)
func (_AWPEmission *AWPEmissionFilterer) ParseAllocationsModified(log types.Log) (*AWPEmissionAllocationsModified, error) {
	event := new(AWPEmissionAllocationsModified)
	if err := _AWPEmission.contract.UnpackLog(event, "AllocationsModified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPEmissionAllocationsSubmittedIterator is returned from FilterAllocationsSubmitted and is used to iterate over the raw logs and unpacked data for AllocationsSubmitted events raised by the AWPEmission contract.
type AWPEmissionAllocationsSubmittedIterator struct {
	Event *AWPEmissionAllocationsSubmitted // Event containing the contract specifics and raw log

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
func (it *AWPEmissionAllocationsSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPEmissionAllocationsSubmitted)
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
		it.Event = new(AWPEmissionAllocationsSubmitted)
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
func (it *AWPEmissionAllocationsSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPEmissionAllocationsSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPEmissionAllocationsSubmitted represents a AllocationsSubmitted event raised by the AWPEmission contract.
type AWPEmissionAllocationsSubmitted struct {
	EffectiveEpoch *big.Int
	Packed         []*big.Int
	TotalWeight    *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterAllocationsSubmitted is a free log retrieval operation binding the contract event 0x1974add8d276f2a35c6cdd57bbda88ad340725bb6becf62aa0325d62fa3a6372.
//
// Solidity: event AllocationsSubmitted(uint256 indexed effectiveEpoch, uint256[] packed, uint256 totalWeight)
func (_AWPEmission *AWPEmissionFilterer) FilterAllocationsSubmitted(opts *bind.FilterOpts, effectiveEpoch []*big.Int) (*AWPEmissionAllocationsSubmittedIterator, error) {

	var effectiveEpochRule []interface{}
	for _, effectiveEpochItem := range effectiveEpoch {
		effectiveEpochRule = append(effectiveEpochRule, effectiveEpochItem)
	}

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "AllocationsSubmitted", effectiveEpochRule)
	if err != nil {
		return nil, err
	}
	return &AWPEmissionAllocationsSubmittedIterator{contract: _AWPEmission.contract, event: "AllocationsSubmitted", logs: logs, sub: sub}, nil
}

// WatchAllocationsSubmitted is a free log subscription operation binding the contract event 0x1974add8d276f2a35c6cdd57bbda88ad340725bb6becf62aa0325d62fa3a6372.
//
// Solidity: event AllocationsSubmitted(uint256 indexed effectiveEpoch, uint256[] packed, uint256 totalWeight)
func (_AWPEmission *AWPEmissionFilterer) WatchAllocationsSubmitted(opts *bind.WatchOpts, sink chan<- *AWPEmissionAllocationsSubmitted, effectiveEpoch []*big.Int) (event.Subscription, error) {

	var effectiveEpochRule []interface{}
	for _, effectiveEpochItem := range effectiveEpoch {
		effectiveEpochRule = append(effectiveEpochRule, effectiveEpochItem)
	}

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "AllocationsSubmitted", effectiveEpochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPEmissionAllocationsSubmitted)
				if err := _AWPEmission.contract.UnpackLog(event, "AllocationsSubmitted", log); err != nil {
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

// ParseAllocationsSubmitted is a log parse operation binding the contract event 0x1974add8d276f2a35c6cdd57bbda88ad340725bb6becf62aa0325d62fa3a6372.
//
// Solidity: event AllocationsSubmitted(uint256 indexed effectiveEpoch, uint256[] packed, uint256 totalWeight)
func (_AWPEmission *AWPEmissionFilterer) ParseAllocationsSubmitted(log types.Log) (*AWPEmissionAllocationsSubmitted, error) {
	event := new(AWPEmissionAllocationsSubmitted)
	if err := _AWPEmission.contract.UnpackLog(event, "AllocationsSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPEmissionDecayFactorUpdatedIterator is returned from FilterDecayFactorUpdated and is used to iterate over the raw logs and unpacked data for DecayFactorUpdated events raised by the AWPEmission contract.
type AWPEmissionDecayFactorUpdatedIterator struct {
	Event *AWPEmissionDecayFactorUpdated // Event containing the contract specifics and raw log

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
func (it *AWPEmissionDecayFactorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPEmissionDecayFactorUpdated)
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
		it.Event = new(AWPEmissionDecayFactorUpdated)
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
func (it *AWPEmissionDecayFactorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPEmissionDecayFactorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPEmissionDecayFactorUpdated represents a DecayFactorUpdated event raised by the AWPEmission contract.
type AWPEmissionDecayFactorUpdated struct {
	NewDecayFactor *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDecayFactorUpdated is a free log retrieval operation binding the contract event 0x020897283b668d79f63f4a336a8d20473f83e97516d34d97db33234af6da821f.
//
// Solidity: event DecayFactorUpdated(uint256 newDecayFactor)
func (_AWPEmission *AWPEmissionFilterer) FilterDecayFactorUpdated(opts *bind.FilterOpts) (*AWPEmissionDecayFactorUpdatedIterator, error) {

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "DecayFactorUpdated")
	if err != nil {
		return nil, err
	}
	return &AWPEmissionDecayFactorUpdatedIterator{contract: _AWPEmission.contract, event: "DecayFactorUpdated", logs: logs, sub: sub}, nil
}

// WatchDecayFactorUpdated is a free log subscription operation binding the contract event 0x020897283b668d79f63f4a336a8d20473f83e97516d34d97db33234af6da821f.
//
// Solidity: event DecayFactorUpdated(uint256 newDecayFactor)
func (_AWPEmission *AWPEmissionFilterer) WatchDecayFactorUpdated(opts *bind.WatchOpts, sink chan<- *AWPEmissionDecayFactorUpdated) (event.Subscription, error) {

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "DecayFactorUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPEmissionDecayFactorUpdated)
				if err := _AWPEmission.contract.UnpackLog(event, "DecayFactorUpdated", log); err != nil {
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

// ParseDecayFactorUpdated is a log parse operation binding the contract event 0x020897283b668d79f63f4a336a8d20473f83e97516d34d97db33234af6da821f.
//
// Solidity: event DecayFactorUpdated(uint256 newDecayFactor)
func (_AWPEmission *AWPEmissionFilterer) ParseDecayFactorUpdated(log types.Log) (*AWPEmissionDecayFactorUpdated, error) {
	event := new(AWPEmissionDecayFactorUpdated)
	if err := _AWPEmission.contract.UnpackLog(event, "DecayFactorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPEmissionEpochDurationUpdatedIterator is returned from FilterEpochDurationUpdated and is used to iterate over the raw logs and unpacked data for EpochDurationUpdated events raised by the AWPEmission contract.
type AWPEmissionEpochDurationUpdatedIterator struct {
	Event *AWPEmissionEpochDurationUpdated // Event containing the contract specifics and raw log

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
func (it *AWPEmissionEpochDurationUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPEmissionEpochDurationUpdated)
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
		it.Event = new(AWPEmissionEpochDurationUpdated)
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
func (it *AWPEmissionEpochDurationUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPEmissionEpochDurationUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPEmissionEpochDurationUpdated represents a EpochDurationUpdated event raised by the AWPEmission contract.
type AWPEmissionEpochDurationUpdated struct {
	OldDuration *big.Int
	NewDuration *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterEpochDurationUpdated is a free log retrieval operation binding the contract event 0xda8ff87475657c76bff7b4e017c535d5fbf29958aedc8095a45c56d9fab528f6.
//
// Solidity: event EpochDurationUpdated(uint256 oldDuration, uint256 newDuration)
func (_AWPEmission *AWPEmissionFilterer) FilterEpochDurationUpdated(opts *bind.FilterOpts) (*AWPEmissionEpochDurationUpdatedIterator, error) {

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "EpochDurationUpdated")
	if err != nil {
		return nil, err
	}
	return &AWPEmissionEpochDurationUpdatedIterator{contract: _AWPEmission.contract, event: "EpochDurationUpdated", logs: logs, sub: sub}, nil
}

// WatchEpochDurationUpdated is a free log subscription operation binding the contract event 0xda8ff87475657c76bff7b4e017c535d5fbf29958aedc8095a45c56d9fab528f6.
//
// Solidity: event EpochDurationUpdated(uint256 oldDuration, uint256 newDuration)
func (_AWPEmission *AWPEmissionFilterer) WatchEpochDurationUpdated(opts *bind.WatchOpts, sink chan<- *AWPEmissionEpochDurationUpdated) (event.Subscription, error) {

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "EpochDurationUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPEmissionEpochDurationUpdated)
				if err := _AWPEmission.contract.UnpackLog(event, "EpochDurationUpdated", log); err != nil {
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

// ParseEpochDurationUpdated is a log parse operation binding the contract event 0xda8ff87475657c76bff7b4e017c535d5fbf29958aedc8095a45c56d9fab528f6.
//
// Solidity: event EpochDurationUpdated(uint256 oldDuration, uint256 newDuration)
func (_AWPEmission *AWPEmissionFilterer) ParseEpochDurationUpdated(log types.Log) (*AWPEmissionEpochDurationUpdated, error) {
	event := new(AWPEmissionEpochDurationUpdated)
	if err := _AWPEmission.contract.UnpackLog(event, "EpochDurationUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPEmissionEpochPausedUntilIterator is returned from FilterEpochPausedUntil and is used to iterate over the raw logs and unpacked data for EpochPausedUntil events raised by the AWPEmission contract.
type AWPEmissionEpochPausedUntilIterator struct {
	Event *AWPEmissionEpochPausedUntil // Event containing the contract specifics and raw log

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
func (it *AWPEmissionEpochPausedUntilIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPEmissionEpochPausedUntil)
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
		it.Event = new(AWPEmissionEpochPausedUntil)
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
func (it *AWPEmissionEpochPausedUntilIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPEmissionEpochPausedUntilIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPEmissionEpochPausedUntil represents a EpochPausedUntil event raised by the AWPEmission contract.
type AWPEmissionEpochPausedUntil struct {
	ResumeTime  uint64
	FrozenEpoch uint64
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterEpochPausedUntil is a free log retrieval operation binding the contract event 0x28f262d609401962a40e1b4fb7254066679926b444f86947cf8108c86ecd80e0.
//
// Solidity: event EpochPausedUntil(uint64 resumeTime, uint64 frozenEpoch)
func (_AWPEmission *AWPEmissionFilterer) FilterEpochPausedUntil(opts *bind.FilterOpts) (*AWPEmissionEpochPausedUntilIterator, error) {

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "EpochPausedUntil")
	if err != nil {
		return nil, err
	}
	return &AWPEmissionEpochPausedUntilIterator{contract: _AWPEmission.contract, event: "EpochPausedUntil", logs: logs, sub: sub}, nil
}

// WatchEpochPausedUntil is a free log subscription operation binding the contract event 0x28f262d609401962a40e1b4fb7254066679926b444f86947cf8108c86ecd80e0.
//
// Solidity: event EpochPausedUntil(uint64 resumeTime, uint64 frozenEpoch)
func (_AWPEmission *AWPEmissionFilterer) WatchEpochPausedUntil(opts *bind.WatchOpts, sink chan<- *AWPEmissionEpochPausedUntil) (event.Subscription, error) {

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "EpochPausedUntil")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPEmissionEpochPausedUntil)
				if err := _AWPEmission.contract.UnpackLog(event, "EpochPausedUntil", log); err != nil {
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

// ParseEpochPausedUntil is a log parse operation binding the contract event 0x28f262d609401962a40e1b4fb7254066679926b444f86947cf8108c86ecd80e0.
//
// Solidity: event EpochPausedUntil(uint64 resumeTime, uint64 frozenEpoch)
func (_AWPEmission *AWPEmissionFilterer) ParseEpochPausedUntil(log types.Log) (*AWPEmissionEpochPausedUntil, error) {
	event := new(AWPEmissionEpochPausedUntil)
	if err := _AWPEmission.contract.UnpackLog(event, "EpochPausedUntil", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPEmissionEpochSettledIterator is returned from FilterEpochSettled and is used to iterate over the raw logs and unpacked data for EpochSettled events raised by the AWPEmission contract.
type AWPEmissionEpochSettledIterator struct {
	Event *AWPEmissionEpochSettled // Event containing the contract specifics and raw log

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
func (it *AWPEmissionEpochSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPEmissionEpochSettled)
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
		it.Event = new(AWPEmissionEpochSettled)
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
func (it *AWPEmissionEpochSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPEmissionEpochSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPEmissionEpochSettled represents a EpochSettled event raised by the AWPEmission contract.
type AWPEmissionEpochSettled struct {
	Epoch          *big.Int
	TotalEmission  *big.Int
	RecipientCount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterEpochSettled is a free log retrieval operation binding the contract event 0xe6dfe5a4e48226c4ece8d7eb3d8e0b37cd466ccb8e8b30ac5f4cfd81b928f07b.
//
// Solidity: event EpochSettled(uint256 indexed epoch, uint256 totalEmission, uint256 recipientCount)
func (_AWPEmission *AWPEmissionFilterer) FilterEpochSettled(opts *bind.FilterOpts, epoch []*big.Int) (*AWPEmissionEpochSettledIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "EpochSettled", epochRule)
	if err != nil {
		return nil, err
	}
	return &AWPEmissionEpochSettledIterator{contract: _AWPEmission.contract, event: "EpochSettled", logs: logs, sub: sub}, nil
}

// WatchEpochSettled is a free log subscription operation binding the contract event 0xe6dfe5a4e48226c4ece8d7eb3d8e0b37cd466ccb8e8b30ac5f4cfd81b928f07b.
//
// Solidity: event EpochSettled(uint256 indexed epoch, uint256 totalEmission, uint256 recipientCount)
func (_AWPEmission *AWPEmissionFilterer) WatchEpochSettled(opts *bind.WatchOpts, sink chan<- *AWPEmissionEpochSettled, epoch []*big.Int) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "EpochSettled", epochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPEmissionEpochSettled)
				if err := _AWPEmission.contract.UnpackLog(event, "EpochSettled", log); err != nil {
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

// ParseEpochSettled is a log parse operation binding the contract event 0xe6dfe5a4e48226c4ece8d7eb3d8e0b37cd466ccb8e8b30ac5f4cfd81b928f07b.
//
// Solidity: event EpochSettled(uint256 indexed epoch, uint256 totalEmission, uint256 recipientCount)
func (_AWPEmission *AWPEmissionFilterer) ParseEpochSettled(log types.Log) (*AWPEmissionEpochSettled, error) {
	event := new(AWPEmissionEpochSettled)
	if err := _AWPEmission.contract.UnpackLog(event, "EpochSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPEmissionGuardianUpdatedIterator is returned from FilterGuardianUpdated and is used to iterate over the raw logs and unpacked data for GuardianUpdated events raised by the AWPEmission contract.
type AWPEmissionGuardianUpdatedIterator struct {
	Event *AWPEmissionGuardianUpdated // Event containing the contract specifics and raw log

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
func (it *AWPEmissionGuardianUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPEmissionGuardianUpdated)
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
		it.Event = new(AWPEmissionGuardianUpdated)
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
func (it *AWPEmissionGuardianUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPEmissionGuardianUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPEmissionGuardianUpdated represents a GuardianUpdated event raised by the AWPEmission contract.
type AWPEmissionGuardianUpdated struct {
	NewGuardian common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterGuardianUpdated is a free log retrieval operation binding the contract event 0x6bb7ff33e730289800c62ad882105a144a74010d2bdbb9a942544a3005ad55bf.
//
// Solidity: event GuardianUpdated(address indexed newGuardian)
func (_AWPEmission *AWPEmissionFilterer) FilterGuardianUpdated(opts *bind.FilterOpts, newGuardian []common.Address) (*AWPEmissionGuardianUpdatedIterator, error) {

	var newGuardianRule []interface{}
	for _, newGuardianItem := range newGuardian {
		newGuardianRule = append(newGuardianRule, newGuardianItem)
	}

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "GuardianUpdated", newGuardianRule)
	if err != nil {
		return nil, err
	}
	return &AWPEmissionGuardianUpdatedIterator{contract: _AWPEmission.contract, event: "GuardianUpdated", logs: logs, sub: sub}, nil
}

// WatchGuardianUpdated is a free log subscription operation binding the contract event 0x6bb7ff33e730289800c62ad882105a144a74010d2bdbb9a942544a3005ad55bf.
//
// Solidity: event GuardianUpdated(address indexed newGuardian)
func (_AWPEmission *AWPEmissionFilterer) WatchGuardianUpdated(opts *bind.WatchOpts, sink chan<- *AWPEmissionGuardianUpdated, newGuardian []common.Address) (event.Subscription, error) {

	var newGuardianRule []interface{}
	for _, newGuardianItem := range newGuardian {
		newGuardianRule = append(newGuardianRule, newGuardianItem)
	}

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "GuardianUpdated", newGuardianRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPEmissionGuardianUpdated)
				if err := _AWPEmission.contract.UnpackLog(event, "GuardianUpdated", log); err != nil {
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

// ParseGuardianUpdated is a log parse operation binding the contract event 0x6bb7ff33e730289800c62ad882105a144a74010d2bdbb9a942544a3005ad55bf.
//
// Solidity: event GuardianUpdated(address indexed newGuardian)
func (_AWPEmission *AWPEmissionFilterer) ParseGuardianUpdated(log types.Log) (*AWPEmissionGuardianUpdated, error) {
	event := new(AWPEmissionGuardianUpdated)
	if err := _AWPEmission.contract.UnpackLog(event, "GuardianUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPEmissionInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the AWPEmission contract.
type AWPEmissionInitializedIterator struct {
	Event *AWPEmissionInitialized // Event containing the contract specifics and raw log

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
func (it *AWPEmissionInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPEmissionInitialized)
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
		it.Event = new(AWPEmissionInitialized)
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
func (it *AWPEmissionInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPEmissionInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPEmissionInitialized represents a Initialized event raised by the AWPEmission contract.
type AWPEmissionInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AWPEmission *AWPEmissionFilterer) FilterInitialized(opts *bind.FilterOpts) (*AWPEmissionInitializedIterator, error) {

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &AWPEmissionInitializedIterator{contract: _AWPEmission.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AWPEmission *AWPEmissionFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *AWPEmissionInitialized) (event.Subscription, error) {

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPEmissionInitialized)
				if err := _AWPEmission.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_AWPEmission *AWPEmissionFilterer) ParseInitialized(log types.Log) (*AWPEmissionInitialized, error) {
	event := new(AWPEmissionInitialized)
	if err := _AWPEmission.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPEmissionMaxRecipientsUpdatedIterator is returned from FilterMaxRecipientsUpdated and is used to iterate over the raw logs and unpacked data for MaxRecipientsUpdated events raised by the AWPEmission contract.
type AWPEmissionMaxRecipientsUpdatedIterator struct {
	Event *AWPEmissionMaxRecipientsUpdated // Event containing the contract specifics and raw log

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
func (it *AWPEmissionMaxRecipientsUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPEmissionMaxRecipientsUpdated)
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
		it.Event = new(AWPEmissionMaxRecipientsUpdated)
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
func (it *AWPEmissionMaxRecipientsUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPEmissionMaxRecipientsUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPEmissionMaxRecipientsUpdated represents a MaxRecipientsUpdated event raised by the AWPEmission contract.
type AWPEmissionMaxRecipientsUpdated struct {
	NewMax *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMaxRecipientsUpdated is a free log retrieval operation binding the contract event 0x625e33d0247c41d7d9e8e3333a170ede206ce082fa83f578f61ff578f5a4fc0b.
//
// Solidity: event MaxRecipientsUpdated(uint256 newMax)
func (_AWPEmission *AWPEmissionFilterer) FilterMaxRecipientsUpdated(opts *bind.FilterOpts) (*AWPEmissionMaxRecipientsUpdatedIterator, error) {

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "MaxRecipientsUpdated")
	if err != nil {
		return nil, err
	}
	return &AWPEmissionMaxRecipientsUpdatedIterator{contract: _AWPEmission.contract, event: "MaxRecipientsUpdated", logs: logs, sub: sub}, nil
}

// WatchMaxRecipientsUpdated is a free log subscription operation binding the contract event 0x625e33d0247c41d7d9e8e3333a170ede206ce082fa83f578f61ff578f5a4fc0b.
//
// Solidity: event MaxRecipientsUpdated(uint256 newMax)
func (_AWPEmission *AWPEmissionFilterer) WatchMaxRecipientsUpdated(opts *bind.WatchOpts, sink chan<- *AWPEmissionMaxRecipientsUpdated) (event.Subscription, error) {

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "MaxRecipientsUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPEmissionMaxRecipientsUpdated)
				if err := _AWPEmission.contract.UnpackLog(event, "MaxRecipientsUpdated", log); err != nil {
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

// ParseMaxRecipientsUpdated is a log parse operation binding the contract event 0x625e33d0247c41d7d9e8e3333a170ede206ce082fa83f578f61ff578f5a4fc0b.
//
// Solidity: event MaxRecipientsUpdated(uint256 newMax)
func (_AWPEmission *AWPEmissionFilterer) ParseMaxRecipientsUpdated(log types.Log) (*AWPEmissionMaxRecipientsUpdated, error) {
	event := new(AWPEmissionMaxRecipientsUpdated)
	if err := _AWPEmission.contract.UnpackLog(event, "MaxRecipientsUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPEmissionRecipientAWPDistributedIterator is returned from FilterRecipientAWPDistributed and is used to iterate over the raw logs and unpacked data for RecipientAWPDistributed events raised by the AWPEmission contract.
type AWPEmissionRecipientAWPDistributedIterator struct {
	Event *AWPEmissionRecipientAWPDistributed // Event containing the contract specifics and raw log

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
func (it *AWPEmissionRecipientAWPDistributedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPEmissionRecipientAWPDistributed)
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
		it.Event = new(AWPEmissionRecipientAWPDistributed)
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
func (it *AWPEmissionRecipientAWPDistributedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPEmissionRecipientAWPDistributedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPEmissionRecipientAWPDistributed represents a RecipientAWPDistributed event raised by the AWPEmission contract.
type AWPEmissionRecipientAWPDistributed struct {
	Epoch     *big.Int
	Recipient common.Address
	AwpAmount *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRecipientAWPDistributed is a free log retrieval operation binding the contract event 0x48ef71e94da9b216d39856e2ba2d0b674a04efb63df0a414ab31e8ce0b57e594.
//
// Solidity: event RecipientAWPDistributed(uint256 indexed epoch, address indexed recipient, uint256 awpAmount)
func (_AWPEmission *AWPEmissionFilterer) FilterRecipientAWPDistributed(opts *bind.FilterOpts, epoch []*big.Int, recipient []common.Address) (*AWPEmissionRecipientAWPDistributedIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "RecipientAWPDistributed", epochRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &AWPEmissionRecipientAWPDistributedIterator{contract: _AWPEmission.contract, event: "RecipientAWPDistributed", logs: logs, sub: sub}, nil
}

// WatchRecipientAWPDistributed is a free log subscription operation binding the contract event 0x48ef71e94da9b216d39856e2ba2d0b674a04efb63df0a414ab31e8ce0b57e594.
//
// Solidity: event RecipientAWPDistributed(uint256 indexed epoch, address indexed recipient, uint256 awpAmount)
func (_AWPEmission *AWPEmissionFilterer) WatchRecipientAWPDistributed(opts *bind.WatchOpts, sink chan<- *AWPEmissionRecipientAWPDistributed, epoch []*big.Int, recipient []common.Address) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "RecipientAWPDistributed", epochRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPEmissionRecipientAWPDistributed)
				if err := _AWPEmission.contract.UnpackLog(event, "RecipientAWPDistributed", log); err != nil {
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

// ParseRecipientAWPDistributed is a log parse operation binding the contract event 0x48ef71e94da9b216d39856e2ba2d0b674a04efb63df0a414ab31e8ce0b57e594.
//
// Solidity: event RecipientAWPDistributed(uint256 indexed epoch, address indexed recipient, uint256 awpAmount)
func (_AWPEmission *AWPEmissionFilterer) ParseRecipientAWPDistributed(log types.Log) (*AWPEmissionRecipientAWPDistributed, error) {
	event := new(AWPEmissionRecipientAWPDistributed)
	if err := _AWPEmission.contract.UnpackLog(event, "RecipientAWPDistributed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPEmissionTreasuryUpdatedIterator is returned from FilterTreasuryUpdated and is used to iterate over the raw logs and unpacked data for TreasuryUpdated events raised by the AWPEmission contract.
type AWPEmissionTreasuryUpdatedIterator struct {
	Event *AWPEmissionTreasuryUpdated // Event containing the contract specifics and raw log

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
func (it *AWPEmissionTreasuryUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPEmissionTreasuryUpdated)
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
		it.Event = new(AWPEmissionTreasuryUpdated)
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
func (it *AWPEmissionTreasuryUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPEmissionTreasuryUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPEmissionTreasuryUpdated represents a TreasuryUpdated event raised by the AWPEmission contract.
type AWPEmissionTreasuryUpdated struct {
	NewTreasury common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTreasuryUpdated is a free log retrieval operation binding the contract event 0x7dae230f18360d76a040c81f050aa14eb9d6dc7901b20fc5d855e2a20fe814d1.
//
// Solidity: event TreasuryUpdated(address indexed newTreasury)
func (_AWPEmission *AWPEmissionFilterer) FilterTreasuryUpdated(opts *bind.FilterOpts, newTreasury []common.Address) (*AWPEmissionTreasuryUpdatedIterator, error) {

	var newTreasuryRule []interface{}
	for _, newTreasuryItem := range newTreasury {
		newTreasuryRule = append(newTreasuryRule, newTreasuryItem)
	}

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "TreasuryUpdated", newTreasuryRule)
	if err != nil {
		return nil, err
	}
	return &AWPEmissionTreasuryUpdatedIterator{contract: _AWPEmission.contract, event: "TreasuryUpdated", logs: logs, sub: sub}, nil
}

// WatchTreasuryUpdated is a free log subscription operation binding the contract event 0x7dae230f18360d76a040c81f050aa14eb9d6dc7901b20fc5d855e2a20fe814d1.
//
// Solidity: event TreasuryUpdated(address indexed newTreasury)
func (_AWPEmission *AWPEmissionFilterer) WatchTreasuryUpdated(opts *bind.WatchOpts, sink chan<- *AWPEmissionTreasuryUpdated, newTreasury []common.Address) (event.Subscription, error) {

	var newTreasuryRule []interface{}
	for _, newTreasuryItem := range newTreasury {
		newTreasuryRule = append(newTreasuryRule, newTreasuryItem)
	}

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "TreasuryUpdated", newTreasuryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPEmissionTreasuryUpdated)
				if err := _AWPEmission.contract.UnpackLog(event, "TreasuryUpdated", log); err != nil {
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

// ParseTreasuryUpdated is a log parse operation binding the contract event 0x7dae230f18360d76a040c81f050aa14eb9d6dc7901b20fc5d855e2a20fe814d1.
//
// Solidity: event TreasuryUpdated(address indexed newTreasury)
func (_AWPEmission *AWPEmissionFilterer) ParseTreasuryUpdated(log types.Log) (*AWPEmissionTreasuryUpdated, error) {
	event := new(AWPEmissionTreasuryUpdated)
	if err := _AWPEmission.contract.UnpackLog(event, "TreasuryUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPEmissionUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the AWPEmission contract.
type AWPEmissionUpgradedIterator struct {
	Event *AWPEmissionUpgraded // Event containing the contract specifics and raw log

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
func (it *AWPEmissionUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPEmissionUpgraded)
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
		it.Event = new(AWPEmissionUpgraded)
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
func (it *AWPEmissionUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPEmissionUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPEmissionUpgraded represents a Upgraded event raised by the AWPEmission contract.
type AWPEmissionUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AWPEmission *AWPEmissionFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*AWPEmissionUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &AWPEmissionUpgradedIterator{contract: _AWPEmission.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AWPEmission *AWPEmissionFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *AWPEmissionUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPEmissionUpgraded)
				if err := _AWPEmission.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_AWPEmission *AWPEmissionFilterer) ParseUpgraded(log types.Log) (*AWPEmissionUpgraded, error) {
	event := new(AWPEmissionUpgraded)
	if err := _AWPEmission.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
