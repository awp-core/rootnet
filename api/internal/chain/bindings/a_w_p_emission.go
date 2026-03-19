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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DECAY_FACTOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DECAY_PRECISION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"EMISSION_SPLIT_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activeEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allocationNonce\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"awpToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIAWPToken\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"currentDailyEmission\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"currentEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"emergencySetWeight\",\"inputs\":[{\"name\":\"epoch_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"weight\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"epochDuration\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"epochEmissionLocked\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"genesisTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEpochRecipientCount\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEpochTotalWeight\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEpochWeight\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOracleCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRecipient\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRecipientCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalWeight\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWeight\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"awpToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"treasury_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"initialDailyEmission_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"genesisTime_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"epochDuration_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"maxRecipients\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"oracleThreshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"oracles\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setOracleConfig\",\"inputs\":[{\"name\":\"oracles_\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"threshold_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"settleEpoch\",\"inputs\":[{\"name\":\"limit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"settleProgress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"settledEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"submitAllocations\",\"inputs\":[{\"name\":\"recipients_\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"weights_\",\"type\":\"uint96[]\",\"internalType\":\"uint96[]\"},{\"name\":\"signatures\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"effectiveEpoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"treasury\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AllocationsSubmitted\",\"inputs\":[{\"name\":\"nonce\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"recipients\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"weights\",\"type\":\"uint96[]\",\"indexed\":false,\"internalType\":\"uint96[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DAOMatchDistributed\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochSettled\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"totalEmission\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"recipientCount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GovernanceWeightUpdated\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"weight\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OracleConfigUpdated\",\"inputs\":[{\"name\":\"oracles\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"threshold\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RecipientAWPDistributed\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"awpAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ArrayLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateOracle\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DuplicateSigner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EpochNotReady\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidOracleConfig\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidParameter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignatureCount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MiningComplete\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustBeFutureEpoch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotTimelock\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OracleNotConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SettlementInProgress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UnknownOracle\",\"inputs\":[]}]",
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

// DECAYFACTOR is a free data retrieval call binding the contract method 0xe08eae5f.
//
// Solidity: function DECAY_FACTOR() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) DECAYFACTOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "DECAY_FACTOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DECAYFACTOR is a free data retrieval call binding the contract method 0xe08eae5f.
//
// Solidity: function DECAY_FACTOR() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) DECAYFACTOR() (*big.Int, error) {
	return _AWPEmission.Contract.DECAYFACTOR(&_AWPEmission.CallOpts)
}

// DECAYFACTOR is a free data retrieval call binding the contract method 0xe08eae5f.
//
// Solidity: function DECAY_FACTOR() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) DECAYFACTOR() (*big.Int, error) {
	return _AWPEmission.Contract.DECAYFACTOR(&_AWPEmission.CallOpts)
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

// EMISSIONSPLITBPS is a free data retrieval call binding the contract method 0x202a1204.
//
// Solidity: function EMISSION_SPLIT_BPS() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) EMISSIONSPLITBPS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "EMISSION_SPLIT_BPS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EMISSIONSPLITBPS is a free data retrieval call binding the contract method 0x202a1204.
//
// Solidity: function EMISSION_SPLIT_BPS() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) EMISSIONSPLITBPS() (*big.Int, error) {
	return _AWPEmission.Contract.EMISSIONSPLITBPS(&_AWPEmission.CallOpts)
}

// EMISSIONSPLITBPS is a free data retrieval call binding the contract method 0x202a1204.
//
// Solidity: function EMISSION_SPLIT_BPS() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) EMISSIONSPLITBPS() (*big.Int, error) {
	return _AWPEmission.Contract.EMISSIONSPLITBPS(&_AWPEmission.CallOpts)
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

// AllocationNonce is a free data retrieval call binding the contract method 0x5d77fd87.
//
// Solidity: function allocationNonce() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) AllocationNonce(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "allocationNonce")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AllocationNonce is a free data retrieval call binding the contract method 0x5d77fd87.
//
// Solidity: function allocationNonce() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) AllocationNonce() (*big.Int, error) {
	return _AWPEmission.Contract.AllocationNonce(&_AWPEmission.CallOpts)
}

// AllocationNonce is a free data retrieval call binding the contract method 0x5d77fd87.
//
// Solidity: function allocationNonce() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) AllocationNonce() (*big.Int, error) {
	return _AWPEmission.Contract.AllocationNonce(&_AWPEmission.CallOpts)
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

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_AWPEmission *AWPEmissionCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "eip712Domain")

	outstruct := new(struct {
		Fields            [1]byte
		Name              string
		Version           string
		ChainId           *big.Int
		VerifyingContract common.Address
		Salt              [32]byte
		Extensions        []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_AWPEmission *AWPEmissionSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _AWPEmission.Contract.Eip712Domain(&_AWPEmission.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_AWPEmission *AWPEmissionCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _AWPEmission.Contract.Eip712Domain(&_AWPEmission.CallOpts)
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

// GenesisTime is a free data retrieval call binding the contract method 0x42c6498a.
//
// Solidity: function genesisTime() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) GenesisTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "genesisTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GenesisTime is a free data retrieval call binding the contract method 0x42c6498a.
//
// Solidity: function genesisTime() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) GenesisTime() (*big.Int, error) {
	return _AWPEmission.Contract.GenesisTime(&_AWPEmission.CallOpts)
}

// GenesisTime is a free data retrieval call binding the contract method 0x42c6498a.
//
// Solidity: function genesisTime() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) GenesisTime() (*big.Int, error) {
	return _AWPEmission.Contract.GenesisTime(&_AWPEmission.CallOpts)
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

// GetOracleCount is a free data retrieval call binding the contract method 0x3f4e4251.
//
// Solidity: function getOracleCount() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) GetOracleCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "getOracleCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOracleCount is a free data retrieval call binding the contract method 0x3f4e4251.
//
// Solidity: function getOracleCount() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) GetOracleCount() (*big.Int, error) {
	return _AWPEmission.Contract.GetOracleCount(&_AWPEmission.CallOpts)
}

// GetOracleCount is a free data retrieval call binding the contract method 0x3f4e4251.
//
// Solidity: function getOracleCount() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) GetOracleCount() (*big.Int, error) {
	return _AWPEmission.Contract.GetOracleCount(&_AWPEmission.CallOpts)
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

// OracleThreshold is a free data retrieval call binding the contract method 0xc379a75d.
//
// Solidity: function oracleThreshold() view returns(uint256)
func (_AWPEmission *AWPEmissionCaller) OracleThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "oracleThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OracleThreshold is a free data retrieval call binding the contract method 0xc379a75d.
//
// Solidity: function oracleThreshold() view returns(uint256)
func (_AWPEmission *AWPEmissionSession) OracleThreshold() (*big.Int, error) {
	return _AWPEmission.Contract.OracleThreshold(&_AWPEmission.CallOpts)
}

// OracleThreshold is a free data retrieval call binding the contract method 0xc379a75d.
//
// Solidity: function oracleThreshold() view returns(uint256)
func (_AWPEmission *AWPEmissionCallerSession) OracleThreshold() (*big.Int, error) {
	return _AWPEmission.Contract.OracleThreshold(&_AWPEmission.CallOpts)
}

// Oracles is a free data retrieval call binding the contract method 0x5b69a7d8.
//
// Solidity: function oracles(uint256 ) view returns(address)
func (_AWPEmission *AWPEmissionCaller) Oracles(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AWPEmission.contract.Call(opts, &out, "oracles", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Oracles is a free data retrieval call binding the contract method 0x5b69a7d8.
//
// Solidity: function oracles(uint256 ) view returns(address)
func (_AWPEmission *AWPEmissionSession) Oracles(arg0 *big.Int) (common.Address, error) {
	return _AWPEmission.Contract.Oracles(&_AWPEmission.CallOpts, arg0)
}

// Oracles is a free data retrieval call binding the contract method 0x5b69a7d8.
//
// Solidity: function oracles(uint256 ) view returns(address)
func (_AWPEmission *AWPEmissionCallerSession) Oracles(arg0 *big.Int) (common.Address, error) {
	return _AWPEmission.Contract.Oracles(&_AWPEmission.CallOpts, arg0)
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

// EmergencySetWeight is a paid mutator transaction binding the contract method 0xa7e5d57a.
//
// Solidity: function emergencySetWeight(uint256 epoch_, uint256 index, address addr, uint96 weight) returns()
func (_AWPEmission *AWPEmissionTransactor) EmergencySetWeight(opts *bind.TransactOpts, epoch_ *big.Int, index *big.Int, addr common.Address, weight *big.Int) (*types.Transaction, error) {
	return _AWPEmission.contract.Transact(opts, "emergencySetWeight", epoch_, index, addr, weight)
}

// EmergencySetWeight is a paid mutator transaction binding the contract method 0xa7e5d57a.
//
// Solidity: function emergencySetWeight(uint256 epoch_, uint256 index, address addr, uint96 weight) returns()
func (_AWPEmission *AWPEmissionSession) EmergencySetWeight(epoch_ *big.Int, index *big.Int, addr common.Address, weight *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.EmergencySetWeight(&_AWPEmission.TransactOpts, epoch_, index, addr, weight)
}

// EmergencySetWeight is a paid mutator transaction binding the contract method 0xa7e5d57a.
//
// Solidity: function emergencySetWeight(uint256 epoch_, uint256 index, address addr, uint96 weight) returns()
func (_AWPEmission *AWPEmissionTransactorSession) EmergencySetWeight(epoch_ *big.Int, index *big.Int, addr common.Address, weight *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.EmergencySetWeight(&_AWPEmission.TransactOpts, epoch_, index, addr, weight)
}

// Initialize is a paid mutator transaction binding the contract method 0xd13f90b4.
//
// Solidity: function initialize(address awpToken_, address treasury_, uint256 initialDailyEmission_, uint256 genesisTime_, uint256 epochDuration_) returns()
func (_AWPEmission *AWPEmissionTransactor) Initialize(opts *bind.TransactOpts, awpToken_ common.Address, treasury_ common.Address, initialDailyEmission_ *big.Int, genesisTime_ *big.Int, epochDuration_ *big.Int) (*types.Transaction, error) {
	return _AWPEmission.contract.Transact(opts, "initialize", awpToken_, treasury_, initialDailyEmission_, genesisTime_, epochDuration_)
}

// Initialize is a paid mutator transaction binding the contract method 0xd13f90b4.
//
// Solidity: function initialize(address awpToken_, address treasury_, uint256 initialDailyEmission_, uint256 genesisTime_, uint256 epochDuration_) returns()
func (_AWPEmission *AWPEmissionSession) Initialize(awpToken_ common.Address, treasury_ common.Address, initialDailyEmission_ *big.Int, genesisTime_ *big.Int, epochDuration_ *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.Initialize(&_AWPEmission.TransactOpts, awpToken_, treasury_, initialDailyEmission_, genesisTime_, epochDuration_)
}

// Initialize is a paid mutator transaction binding the contract method 0xd13f90b4.
//
// Solidity: function initialize(address awpToken_, address treasury_, uint256 initialDailyEmission_, uint256 genesisTime_, uint256 epochDuration_) returns()
func (_AWPEmission *AWPEmissionTransactorSession) Initialize(awpToken_ common.Address, treasury_ common.Address, initialDailyEmission_ *big.Int, genesisTime_ *big.Int, epochDuration_ *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.Initialize(&_AWPEmission.TransactOpts, awpToken_, treasury_, initialDailyEmission_, genesisTime_, epochDuration_)
}

// SetOracleConfig is a paid mutator transaction binding the contract method 0xa5ca27df.
//
// Solidity: function setOracleConfig(address[] oracles_, uint256 threshold_) returns()
func (_AWPEmission *AWPEmissionTransactor) SetOracleConfig(opts *bind.TransactOpts, oracles_ []common.Address, threshold_ *big.Int) (*types.Transaction, error) {
	return _AWPEmission.contract.Transact(opts, "setOracleConfig", oracles_, threshold_)
}

// SetOracleConfig is a paid mutator transaction binding the contract method 0xa5ca27df.
//
// Solidity: function setOracleConfig(address[] oracles_, uint256 threshold_) returns()
func (_AWPEmission *AWPEmissionSession) SetOracleConfig(oracles_ []common.Address, threshold_ *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.SetOracleConfig(&_AWPEmission.TransactOpts, oracles_, threshold_)
}

// SetOracleConfig is a paid mutator transaction binding the contract method 0xa5ca27df.
//
// Solidity: function setOracleConfig(address[] oracles_, uint256 threshold_) returns()
func (_AWPEmission *AWPEmissionTransactorSession) SetOracleConfig(oracles_ []common.Address, threshold_ *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.SetOracleConfig(&_AWPEmission.TransactOpts, oracles_, threshold_)
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

// SubmitAllocations is a paid mutator transaction binding the contract method 0xf7c70a46.
//
// Solidity: function submitAllocations(address[] recipients_, uint96[] weights_, bytes[] signatures, uint256 effectiveEpoch) returns()
func (_AWPEmission *AWPEmissionTransactor) SubmitAllocations(opts *bind.TransactOpts, recipients_ []common.Address, weights_ []*big.Int, signatures [][]byte, effectiveEpoch *big.Int) (*types.Transaction, error) {
	return _AWPEmission.contract.Transact(opts, "submitAllocations", recipients_, weights_, signatures, effectiveEpoch)
}

// SubmitAllocations is a paid mutator transaction binding the contract method 0xf7c70a46.
//
// Solidity: function submitAllocations(address[] recipients_, uint96[] weights_, bytes[] signatures, uint256 effectiveEpoch) returns()
func (_AWPEmission *AWPEmissionSession) SubmitAllocations(recipients_ []common.Address, weights_ []*big.Int, signatures [][]byte, effectiveEpoch *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.SubmitAllocations(&_AWPEmission.TransactOpts, recipients_, weights_, signatures, effectiveEpoch)
}

// SubmitAllocations is a paid mutator transaction binding the contract method 0xf7c70a46.
//
// Solidity: function submitAllocations(address[] recipients_, uint96[] weights_, bytes[] signatures, uint256 effectiveEpoch) returns()
func (_AWPEmission *AWPEmissionTransactorSession) SubmitAllocations(recipients_ []common.Address, weights_ []*big.Int, signatures [][]byte, effectiveEpoch *big.Int) (*types.Transaction, error) {
	return _AWPEmission.Contract.SubmitAllocations(&_AWPEmission.TransactOpts, recipients_, weights_, signatures, effectiveEpoch)
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
	Nonce      *big.Int
	Recipients []common.Address
	Weights    []*big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAllocationsSubmitted is a free log retrieval operation binding the contract event 0x5fc284d54073f14f094fecd84deb0ab9419a66f59bd75cdfe5259a90393415d3.
//
// Solidity: event AllocationsSubmitted(uint256 indexed nonce, address[] recipients, uint96[] weights)
func (_AWPEmission *AWPEmissionFilterer) FilterAllocationsSubmitted(opts *bind.FilterOpts, nonce []*big.Int) (*AWPEmissionAllocationsSubmittedIterator, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "AllocationsSubmitted", nonceRule)
	if err != nil {
		return nil, err
	}
	return &AWPEmissionAllocationsSubmittedIterator{contract: _AWPEmission.contract, event: "AllocationsSubmitted", logs: logs, sub: sub}, nil
}

// WatchAllocationsSubmitted is a free log subscription operation binding the contract event 0x5fc284d54073f14f094fecd84deb0ab9419a66f59bd75cdfe5259a90393415d3.
//
// Solidity: event AllocationsSubmitted(uint256 indexed nonce, address[] recipients, uint96[] weights)
func (_AWPEmission *AWPEmissionFilterer) WatchAllocationsSubmitted(opts *bind.WatchOpts, sink chan<- *AWPEmissionAllocationsSubmitted, nonce []*big.Int) (event.Subscription, error) {

	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "AllocationsSubmitted", nonceRule)
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

// ParseAllocationsSubmitted is a log parse operation binding the contract event 0x5fc284d54073f14f094fecd84deb0ab9419a66f59bd75cdfe5259a90393415d3.
//
// Solidity: event AllocationsSubmitted(uint256 indexed nonce, address[] recipients, uint96[] weights)
func (_AWPEmission *AWPEmissionFilterer) ParseAllocationsSubmitted(log types.Log) (*AWPEmissionAllocationsSubmitted, error) {
	event := new(AWPEmissionAllocationsSubmitted)
	if err := _AWPEmission.contract.UnpackLog(event, "AllocationsSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPEmissionDAOMatchDistributedIterator is returned from FilterDAOMatchDistributed and is used to iterate over the raw logs and unpacked data for DAOMatchDistributed events raised by the AWPEmission contract.
type AWPEmissionDAOMatchDistributedIterator struct {
	Event *AWPEmissionDAOMatchDistributed // Event containing the contract specifics and raw log

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
func (it *AWPEmissionDAOMatchDistributedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPEmissionDAOMatchDistributed)
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
		it.Event = new(AWPEmissionDAOMatchDistributed)
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
func (it *AWPEmissionDAOMatchDistributedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPEmissionDAOMatchDistributedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPEmissionDAOMatchDistributed represents a DAOMatchDistributed event raised by the AWPEmission contract.
type AWPEmissionDAOMatchDistributed struct {
	Epoch  *big.Int
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDAOMatchDistributed is a free log retrieval operation binding the contract event 0xcbb7180c1182752945c0986e8ee2b3c59c938ba01628512ec2b4c1509e0853aa.
//
// Solidity: event DAOMatchDistributed(uint256 indexed epoch, uint256 amount)
func (_AWPEmission *AWPEmissionFilterer) FilterDAOMatchDistributed(opts *bind.FilterOpts, epoch []*big.Int) (*AWPEmissionDAOMatchDistributedIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "DAOMatchDistributed", epochRule)
	if err != nil {
		return nil, err
	}
	return &AWPEmissionDAOMatchDistributedIterator{contract: _AWPEmission.contract, event: "DAOMatchDistributed", logs: logs, sub: sub}, nil
}

// WatchDAOMatchDistributed is a free log subscription operation binding the contract event 0xcbb7180c1182752945c0986e8ee2b3c59c938ba01628512ec2b4c1509e0853aa.
//
// Solidity: event DAOMatchDistributed(uint256 indexed epoch, uint256 amount)
func (_AWPEmission *AWPEmissionFilterer) WatchDAOMatchDistributed(opts *bind.WatchOpts, sink chan<- *AWPEmissionDAOMatchDistributed, epoch []*big.Int) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "DAOMatchDistributed", epochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPEmissionDAOMatchDistributed)
				if err := _AWPEmission.contract.UnpackLog(event, "DAOMatchDistributed", log); err != nil {
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

// ParseDAOMatchDistributed is a log parse operation binding the contract event 0xcbb7180c1182752945c0986e8ee2b3c59c938ba01628512ec2b4c1509e0853aa.
//
// Solidity: event DAOMatchDistributed(uint256 indexed epoch, uint256 amount)
func (_AWPEmission *AWPEmissionFilterer) ParseDAOMatchDistributed(log types.Log) (*AWPEmissionDAOMatchDistributed, error) {
	event := new(AWPEmissionDAOMatchDistributed)
	if err := _AWPEmission.contract.UnpackLog(event, "DAOMatchDistributed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPEmissionEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the AWPEmission contract.
type AWPEmissionEIP712DomainChangedIterator struct {
	Event *AWPEmissionEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *AWPEmissionEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPEmissionEIP712DomainChanged)
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
		it.Event = new(AWPEmissionEIP712DomainChanged)
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
func (it *AWPEmissionEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPEmissionEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPEmissionEIP712DomainChanged represents a EIP712DomainChanged event raised by the AWPEmission contract.
type AWPEmissionEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_AWPEmission *AWPEmissionFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*AWPEmissionEIP712DomainChangedIterator, error) {

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &AWPEmissionEIP712DomainChangedIterator{contract: _AWPEmission.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_AWPEmission *AWPEmissionFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *AWPEmissionEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPEmissionEIP712DomainChanged)
				if err := _AWPEmission.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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

// ParseEIP712DomainChanged is a log parse operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_AWPEmission *AWPEmissionFilterer) ParseEIP712DomainChanged(log types.Log) (*AWPEmissionEIP712DomainChanged, error) {
	event := new(AWPEmissionEIP712DomainChanged)
	if err := _AWPEmission.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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

// AWPEmissionGovernanceWeightUpdatedIterator is returned from FilterGovernanceWeightUpdated and is used to iterate over the raw logs and unpacked data for GovernanceWeightUpdated events raised by the AWPEmission contract.
type AWPEmissionGovernanceWeightUpdatedIterator struct {
	Event *AWPEmissionGovernanceWeightUpdated // Event containing the contract specifics and raw log

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
func (it *AWPEmissionGovernanceWeightUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPEmissionGovernanceWeightUpdated)
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
		it.Event = new(AWPEmissionGovernanceWeightUpdated)
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
func (it *AWPEmissionGovernanceWeightUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPEmissionGovernanceWeightUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPEmissionGovernanceWeightUpdated represents a GovernanceWeightUpdated event raised by the AWPEmission contract.
type AWPEmissionGovernanceWeightUpdated struct {
	Addr   common.Address
	Weight *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterGovernanceWeightUpdated is a free log retrieval operation binding the contract event 0x5e5ab064b47f32e3fe22c86c3b8b8df9d68b1e303fd4c3ef721554ae12466414.
//
// Solidity: event GovernanceWeightUpdated(address indexed addr, uint96 weight)
func (_AWPEmission *AWPEmissionFilterer) FilterGovernanceWeightUpdated(opts *bind.FilterOpts, addr []common.Address) (*AWPEmissionGovernanceWeightUpdatedIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "GovernanceWeightUpdated", addrRule)
	if err != nil {
		return nil, err
	}
	return &AWPEmissionGovernanceWeightUpdatedIterator{contract: _AWPEmission.contract, event: "GovernanceWeightUpdated", logs: logs, sub: sub}, nil
}

// WatchGovernanceWeightUpdated is a free log subscription operation binding the contract event 0x5e5ab064b47f32e3fe22c86c3b8b8df9d68b1e303fd4c3ef721554ae12466414.
//
// Solidity: event GovernanceWeightUpdated(address indexed addr, uint96 weight)
func (_AWPEmission *AWPEmissionFilterer) WatchGovernanceWeightUpdated(opts *bind.WatchOpts, sink chan<- *AWPEmissionGovernanceWeightUpdated, addr []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "GovernanceWeightUpdated", addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPEmissionGovernanceWeightUpdated)
				if err := _AWPEmission.contract.UnpackLog(event, "GovernanceWeightUpdated", log); err != nil {
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

// ParseGovernanceWeightUpdated is a log parse operation binding the contract event 0x5e5ab064b47f32e3fe22c86c3b8b8df9d68b1e303fd4c3ef721554ae12466414.
//
// Solidity: event GovernanceWeightUpdated(address indexed addr, uint96 weight)
func (_AWPEmission *AWPEmissionFilterer) ParseGovernanceWeightUpdated(log types.Log) (*AWPEmissionGovernanceWeightUpdated, error) {
	event := new(AWPEmissionGovernanceWeightUpdated)
	if err := _AWPEmission.contract.UnpackLog(event, "GovernanceWeightUpdated", log); err != nil {
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

// AWPEmissionOracleConfigUpdatedIterator is returned from FilterOracleConfigUpdated and is used to iterate over the raw logs and unpacked data for OracleConfigUpdated events raised by the AWPEmission contract.
type AWPEmissionOracleConfigUpdatedIterator struct {
	Event *AWPEmissionOracleConfigUpdated // Event containing the contract specifics and raw log

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
func (it *AWPEmissionOracleConfigUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPEmissionOracleConfigUpdated)
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
		it.Event = new(AWPEmissionOracleConfigUpdated)
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
func (it *AWPEmissionOracleConfigUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPEmissionOracleConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPEmissionOracleConfigUpdated represents a OracleConfigUpdated event raised by the AWPEmission contract.
type AWPEmissionOracleConfigUpdated struct {
	Oracles   []common.Address
	Threshold *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOracleConfigUpdated is a free log retrieval operation binding the contract event 0x094dfe9b10b430945fee27690a51f22b63e17293d41d0300c094c723f410c90f.
//
// Solidity: event OracleConfigUpdated(address[] oracles, uint256 threshold)
func (_AWPEmission *AWPEmissionFilterer) FilterOracleConfigUpdated(opts *bind.FilterOpts) (*AWPEmissionOracleConfigUpdatedIterator, error) {

	logs, sub, err := _AWPEmission.contract.FilterLogs(opts, "OracleConfigUpdated")
	if err != nil {
		return nil, err
	}
	return &AWPEmissionOracleConfigUpdatedIterator{contract: _AWPEmission.contract, event: "OracleConfigUpdated", logs: logs, sub: sub}, nil
}

// WatchOracleConfigUpdated is a free log subscription operation binding the contract event 0x094dfe9b10b430945fee27690a51f22b63e17293d41d0300c094c723f410c90f.
//
// Solidity: event OracleConfigUpdated(address[] oracles, uint256 threshold)
func (_AWPEmission *AWPEmissionFilterer) WatchOracleConfigUpdated(opts *bind.WatchOpts, sink chan<- *AWPEmissionOracleConfigUpdated) (event.Subscription, error) {

	logs, sub, err := _AWPEmission.contract.WatchLogs(opts, "OracleConfigUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPEmissionOracleConfigUpdated)
				if err := _AWPEmission.contract.UnpackLog(event, "OracleConfigUpdated", log); err != nil {
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

// ParseOracleConfigUpdated is a log parse operation binding the contract event 0x094dfe9b10b430945fee27690a51f22b63e17293d41d0300c094c723f410c90f.
//
// Solidity: event OracleConfigUpdated(address[] oracles, uint256 threshold)
func (_AWPEmission *AWPEmissionFilterer) ParseOracleConfigUpdated(log types.Log) (*AWPEmissionOracleConfigUpdated, error) {
	event := new(AWPEmissionOracleConfigUpdated)
	if err := _AWPEmission.contract.UnpackLog(event, "OracleConfigUpdated", log); err != nil {
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
