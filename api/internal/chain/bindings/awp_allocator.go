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

// AWPAllocatorMetaData contains all meta data concerning the AWPAllocator contract.
var AWPAllocatorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"awpRegistry_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"veAWP_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allocate\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"worknetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allocateFor\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"worknetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"awpRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchAllocate\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agents\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"worknetIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"amounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"batchDeallocate\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agents\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"worknetIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"amounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deallocate\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"worknetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deallocateAll\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"worknetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deallocateFor\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"worknetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAgentStake\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"worknetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAgentWorknets\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWorknetTotalStake\",\"inputs\":[{\"name\":\"worknetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"guardian\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"awpRegistry_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"guardian_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"nonces\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reallocate\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromAgent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromWorknetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"toAgent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"toWorknetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setGuardian\",\"inputs\":[{\"name\":\"g\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"userTotalAllocated\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"veAWP\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"worknetTotalStake\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Allocated\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"worknetId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deallocated\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"worknetId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GuardianUpdated\",\"inputs\":[{\"name\":\"newGuardian\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Reallocated\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fromAgent\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"fromWorknetId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"toAgent\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"toWorknetId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ArrayLengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpiredSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientAllocation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientUnallocated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAuthorized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotGuardian\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroWorknetId\",\"inputs\":[]}]",
}

// AWPAllocatorABI is the input ABI used to generate the binding from.
// Deprecated: Use AWPAllocatorMetaData.ABI instead.
var AWPAllocatorABI = AWPAllocatorMetaData.ABI

// AWPAllocator is an auto generated Go binding around an Ethereum contract.
type AWPAllocator struct {
	AWPAllocatorCaller     // Read-only binding to the contract
	AWPAllocatorTransactor // Write-only binding to the contract
	AWPAllocatorFilterer   // Log filterer for contract events
}

// AWPAllocatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type AWPAllocatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AWPAllocatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AWPAllocatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AWPAllocatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AWPAllocatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AWPAllocatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AWPAllocatorSession struct {
	Contract     *AWPAllocator     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AWPAllocatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AWPAllocatorCallerSession struct {
	Contract *AWPAllocatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// AWPAllocatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AWPAllocatorTransactorSession struct {
	Contract     *AWPAllocatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// AWPAllocatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type AWPAllocatorRaw struct {
	Contract *AWPAllocator // Generic contract binding to access the raw methods on
}

// AWPAllocatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AWPAllocatorCallerRaw struct {
	Contract *AWPAllocatorCaller // Generic read-only contract binding to access the raw methods on
}

// AWPAllocatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AWPAllocatorTransactorRaw struct {
	Contract *AWPAllocatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAWPAllocator creates a new instance of AWPAllocator, bound to a specific deployed contract.
func NewAWPAllocator(address common.Address, backend bind.ContractBackend) (*AWPAllocator, error) {
	contract, err := bindAWPAllocator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AWPAllocator{AWPAllocatorCaller: AWPAllocatorCaller{contract: contract}, AWPAllocatorTransactor: AWPAllocatorTransactor{contract: contract}, AWPAllocatorFilterer: AWPAllocatorFilterer{contract: contract}}, nil
}

// NewAWPAllocatorCaller creates a new read-only instance of AWPAllocator, bound to a specific deployed contract.
func NewAWPAllocatorCaller(address common.Address, caller bind.ContractCaller) (*AWPAllocatorCaller, error) {
	contract, err := bindAWPAllocator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AWPAllocatorCaller{contract: contract}, nil
}

// NewAWPAllocatorTransactor creates a new write-only instance of AWPAllocator, bound to a specific deployed contract.
func NewAWPAllocatorTransactor(address common.Address, transactor bind.ContractTransactor) (*AWPAllocatorTransactor, error) {
	contract, err := bindAWPAllocator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AWPAllocatorTransactor{contract: contract}, nil
}

// NewAWPAllocatorFilterer creates a new log filterer instance of AWPAllocator, bound to a specific deployed contract.
func NewAWPAllocatorFilterer(address common.Address, filterer bind.ContractFilterer) (*AWPAllocatorFilterer, error) {
	contract, err := bindAWPAllocator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AWPAllocatorFilterer{contract: contract}, nil
}

// bindAWPAllocator binds a generic wrapper to an already deployed contract.
func bindAWPAllocator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AWPAllocatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AWPAllocator *AWPAllocatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AWPAllocator.Contract.AWPAllocatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AWPAllocator *AWPAllocatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AWPAllocator.Contract.AWPAllocatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AWPAllocator *AWPAllocatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AWPAllocator.Contract.AWPAllocatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AWPAllocator *AWPAllocatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AWPAllocator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AWPAllocator *AWPAllocatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AWPAllocator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AWPAllocator *AWPAllocatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AWPAllocator.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_AWPAllocator *AWPAllocatorCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AWPAllocator.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_AWPAllocator *AWPAllocatorSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _AWPAllocator.Contract.UPGRADEINTERFACEVERSION(&_AWPAllocator.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_AWPAllocator *AWPAllocatorCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _AWPAllocator.Contract.UPGRADEINTERFACEVERSION(&_AWPAllocator.CallOpts)
}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_AWPAllocator *AWPAllocatorCaller) AwpRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPAllocator.contract.Call(opts, &out, "awpRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_AWPAllocator *AWPAllocatorSession) AwpRegistry() (common.Address, error) {
	return _AWPAllocator.Contract.AwpRegistry(&_AWPAllocator.CallOpts)
}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_AWPAllocator *AWPAllocatorCallerSession) AwpRegistry() (common.Address, error) {
	return _AWPAllocator.Contract.AwpRegistry(&_AWPAllocator.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_AWPAllocator *AWPAllocatorCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _AWPAllocator.contract.Call(opts, &out, "eip712Domain")

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
func (_AWPAllocator *AWPAllocatorSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _AWPAllocator.Contract.Eip712Domain(&_AWPAllocator.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_AWPAllocator *AWPAllocatorCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _AWPAllocator.Contract.Eip712Domain(&_AWPAllocator.CallOpts)
}

// GetAgentStake is a free data retrieval call binding the contract method 0xf1ad80c6.
//
// Solidity: function getAgentStake(address user, address agent, uint256 worknetId) view returns(uint256)
func (_AWPAllocator *AWPAllocatorCaller) GetAgentStake(opts *bind.CallOpts, user common.Address, agent common.Address, worknetId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AWPAllocator.contract.Call(opts, &out, "getAgentStake", user, agent, worknetId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAgentStake is a free data retrieval call binding the contract method 0xf1ad80c6.
//
// Solidity: function getAgentStake(address user, address agent, uint256 worknetId) view returns(uint256)
func (_AWPAllocator *AWPAllocatorSession) GetAgentStake(user common.Address, agent common.Address, worknetId *big.Int) (*big.Int, error) {
	return _AWPAllocator.Contract.GetAgentStake(&_AWPAllocator.CallOpts, user, agent, worknetId)
}

// GetAgentStake is a free data retrieval call binding the contract method 0xf1ad80c6.
//
// Solidity: function getAgentStake(address user, address agent, uint256 worknetId) view returns(uint256)
func (_AWPAllocator *AWPAllocatorCallerSession) GetAgentStake(user common.Address, agent common.Address, worknetId *big.Int) (*big.Int, error) {
	return _AWPAllocator.Contract.GetAgentStake(&_AWPAllocator.CallOpts, user, agent, worknetId)
}

// GetAgentWorknets is a free data retrieval call binding the contract method 0xfd4fd2e8.
//
// Solidity: function getAgentWorknets(address user, address agent) view returns(uint256[])
func (_AWPAllocator *AWPAllocatorCaller) GetAgentWorknets(opts *bind.CallOpts, user common.Address, agent common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _AWPAllocator.contract.Call(opts, &out, "getAgentWorknets", user, agent)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetAgentWorknets is a free data retrieval call binding the contract method 0xfd4fd2e8.
//
// Solidity: function getAgentWorknets(address user, address agent) view returns(uint256[])
func (_AWPAllocator *AWPAllocatorSession) GetAgentWorknets(user common.Address, agent common.Address) ([]*big.Int, error) {
	return _AWPAllocator.Contract.GetAgentWorknets(&_AWPAllocator.CallOpts, user, agent)
}

// GetAgentWorknets is a free data retrieval call binding the contract method 0xfd4fd2e8.
//
// Solidity: function getAgentWorknets(address user, address agent) view returns(uint256[])
func (_AWPAllocator *AWPAllocatorCallerSession) GetAgentWorknets(user common.Address, agent common.Address) ([]*big.Int, error) {
	return _AWPAllocator.Contract.GetAgentWorknets(&_AWPAllocator.CallOpts, user, agent)
}

// GetWorknetTotalStake is a free data retrieval call binding the contract method 0x5bfe237f.
//
// Solidity: function getWorknetTotalStake(uint256 worknetId) view returns(uint256)
func (_AWPAllocator *AWPAllocatorCaller) GetWorknetTotalStake(opts *bind.CallOpts, worknetId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AWPAllocator.contract.Call(opts, &out, "getWorknetTotalStake", worknetId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetWorknetTotalStake is a free data retrieval call binding the contract method 0x5bfe237f.
//
// Solidity: function getWorknetTotalStake(uint256 worknetId) view returns(uint256)
func (_AWPAllocator *AWPAllocatorSession) GetWorknetTotalStake(worknetId *big.Int) (*big.Int, error) {
	return _AWPAllocator.Contract.GetWorknetTotalStake(&_AWPAllocator.CallOpts, worknetId)
}

// GetWorknetTotalStake is a free data retrieval call binding the contract method 0x5bfe237f.
//
// Solidity: function getWorknetTotalStake(uint256 worknetId) view returns(uint256)
func (_AWPAllocator *AWPAllocatorCallerSession) GetWorknetTotalStake(worknetId *big.Int) (*big.Int, error) {
	return _AWPAllocator.Contract.GetWorknetTotalStake(&_AWPAllocator.CallOpts, worknetId)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_AWPAllocator *AWPAllocatorCaller) Guardian(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPAllocator.contract.Call(opts, &out, "guardian")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_AWPAllocator *AWPAllocatorSession) Guardian() (common.Address, error) {
	return _AWPAllocator.Contract.Guardian(&_AWPAllocator.CallOpts)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_AWPAllocator *AWPAllocatorCallerSession) Guardian() (common.Address, error) {
	return _AWPAllocator.Contract.Guardian(&_AWPAllocator.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_AWPAllocator *AWPAllocatorCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AWPAllocator.contract.Call(opts, &out, "nonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_AWPAllocator *AWPAllocatorSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _AWPAllocator.Contract.Nonces(&_AWPAllocator.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_AWPAllocator *AWPAllocatorCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _AWPAllocator.Contract.Nonces(&_AWPAllocator.CallOpts, arg0)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_AWPAllocator *AWPAllocatorCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AWPAllocator.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_AWPAllocator *AWPAllocatorSession) ProxiableUUID() ([32]byte, error) {
	return _AWPAllocator.Contract.ProxiableUUID(&_AWPAllocator.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_AWPAllocator *AWPAllocatorCallerSession) ProxiableUUID() ([32]byte, error) {
	return _AWPAllocator.Contract.ProxiableUUID(&_AWPAllocator.CallOpts)
}

// UserTotalAllocated is a free data retrieval call binding the contract method 0x32ffa4ce.
//
// Solidity: function userTotalAllocated(address ) view returns(uint256)
func (_AWPAllocator *AWPAllocatorCaller) UserTotalAllocated(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AWPAllocator.contract.Call(opts, &out, "userTotalAllocated", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UserTotalAllocated is a free data retrieval call binding the contract method 0x32ffa4ce.
//
// Solidity: function userTotalAllocated(address ) view returns(uint256)
func (_AWPAllocator *AWPAllocatorSession) UserTotalAllocated(arg0 common.Address) (*big.Int, error) {
	return _AWPAllocator.Contract.UserTotalAllocated(&_AWPAllocator.CallOpts, arg0)
}

// UserTotalAllocated is a free data retrieval call binding the contract method 0x32ffa4ce.
//
// Solidity: function userTotalAllocated(address ) view returns(uint256)
func (_AWPAllocator *AWPAllocatorCallerSession) UserTotalAllocated(arg0 common.Address) (*big.Int, error) {
	return _AWPAllocator.Contract.UserTotalAllocated(&_AWPAllocator.CallOpts, arg0)
}

// VeAWP is a free data retrieval call binding the contract method 0x7bb8431f.
//
// Solidity: function veAWP() view returns(address)
func (_AWPAllocator *AWPAllocatorCaller) VeAWP(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPAllocator.contract.Call(opts, &out, "veAWP")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VeAWP is a free data retrieval call binding the contract method 0x7bb8431f.
//
// Solidity: function veAWP() view returns(address)
func (_AWPAllocator *AWPAllocatorSession) VeAWP() (common.Address, error) {
	return _AWPAllocator.Contract.VeAWP(&_AWPAllocator.CallOpts)
}

// VeAWP is a free data retrieval call binding the contract method 0x7bb8431f.
//
// Solidity: function veAWP() view returns(address)
func (_AWPAllocator *AWPAllocatorCallerSession) VeAWP() (common.Address, error) {
	return _AWPAllocator.Contract.VeAWP(&_AWPAllocator.CallOpts)
}

// WorknetTotalStake is a free data retrieval call binding the contract method 0x8c00c09c.
//
// Solidity: function worknetTotalStake(uint256 ) view returns(uint256)
func (_AWPAllocator *AWPAllocatorCaller) WorknetTotalStake(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AWPAllocator.contract.Call(opts, &out, "worknetTotalStake", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WorknetTotalStake is a free data retrieval call binding the contract method 0x8c00c09c.
//
// Solidity: function worknetTotalStake(uint256 ) view returns(uint256)
func (_AWPAllocator *AWPAllocatorSession) WorknetTotalStake(arg0 *big.Int) (*big.Int, error) {
	return _AWPAllocator.Contract.WorknetTotalStake(&_AWPAllocator.CallOpts, arg0)
}

// WorknetTotalStake is a free data retrieval call binding the contract method 0x8c00c09c.
//
// Solidity: function worknetTotalStake(uint256 ) view returns(uint256)
func (_AWPAllocator *AWPAllocatorCallerSession) WorknetTotalStake(arg0 *big.Int) (*big.Int, error) {
	return _AWPAllocator.Contract.WorknetTotalStake(&_AWPAllocator.CallOpts, arg0)
}

// Allocate is a paid mutator transaction binding the contract method 0xd035a9a7.
//
// Solidity: function allocate(address staker, address agent, uint256 worknetId, uint256 amount) returns()
func (_AWPAllocator *AWPAllocatorTransactor) Allocate(opts *bind.TransactOpts, staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPAllocator.contract.Transact(opts, "allocate", staker, agent, worknetId, amount)
}

// Allocate is a paid mutator transaction binding the contract method 0xd035a9a7.
//
// Solidity: function allocate(address staker, address agent, uint256 worknetId, uint256 amount) returns()
func (_AWPAllocator *AWPAllocatorSession) Allocate(staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPAllocator.Contract.Allocate(&_AWPAllocator.TransactOpts, staker, agent, worknetId, amount)
}

// Allocate is a paid mutator transaction binding the contract method 0xd035a9a7.
//
// Solidity: function allocate(address staker, address agent, uint256 worknetId, uint256 amount) returns()
func (_AWPAllocator *AWPAllocatorTransactorSession) Allocate(staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPAllocator.Contract.Allocate(&_AWPAllocator.TransactOpts, staker, agent, worknetId, amount)
}

// AllocateFor is a paid mutator transaction binding the contract method 0x7d66c5c5.
//
// Solidity: function allocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPAllocator *AWPAllocatorTransactor) AllocateFor(opts *bind.TransactOpts, staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPAllocator.contract.Transact(opts, "allocateFor", staker, agent, worknetId, amount, deadline, v, r, s)
}

// AllocateFor is a paid mutator transaction binding the contract method 0x7d66c5c5.
//
// Solidity: function allocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPAllocator *AWPAllocatorSession) AllocateFor(staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPAllocator.Contract.AllocateFor(&_AWPAllocator.TransactOpts, staker, agent, worknetId, amount, deadline, v, r, s)
}

// AllocateFor is a paid mutator transaction binding the contract method 0x7d66c5c5.
//
// Solidity: function allocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPAllocator *AWPAllocatorTransactorSession) AllocateFor(staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPAllocator.Contract.AllocateFor(&_AWPAllocator.TransactOpts, staker, agent, worknetId, amount, deadline, v, r, s)
}

// BatchAllocate is a paid mutator transaction binding the contract method 0x25ad39ee.
//
// Solidity: function batchAllocate(address staker, address[] agents, uint256[] worknetIds, uint256[] amounts) returns()
func (_AWPAllocator *AWPAllocatorTransactor) BatchAllocate(opts *bind.TransactOpts, staker common.Address, agents []common.Address, worknetIds []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _AWPAllocator.contract.Transact(opts, "batchAllocate", staker, agents, worknetIds, amounts)
}

// BatchAllocate is a paid mutator transaction binding the contract method 0x25ad39ee.
//
// Solidity: function batchAllocate(address staker, address[] agents, uint256[] worknetIds, uint256[] amounts) returns()
func (_AWPAllocator *AWPAllocatorSession) BatchAllocate(staker common.Address, agents []common.Address, worknetIds []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _AWPAllocator.Contract.BatchAllocate(&_AWPAllocator.TransactOpts, staker, agents, worknetIds, amounts)
}

// BatchAllocate is a paid mutator transaction binding the contract method 0x25ad39ee.
//
// Solidity: function batchAllocate(address staker, address[] agents, uint256[] worknetIds, uint256[] amounts) returns()
func (_AWPAllocator *AWPAllocatorTransactorSession) BatchAllocate(staker common.Address, agents []common.Address, worknetIds []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _AWPAllocator.Contract.BatchAllocate(&_AWPAllocator.TransactOpts, staker, agents, worknetIds, amounts)
}

// BatchDeallocate is a paid mutator transaction binding the contract method 0x2e17f308.
//
// Solidity: function batchDeallocate(address staker, address[] agents, uint256[] worknetIds, uint256[] amounts) returns()
func (_AWPAllocator *AWPAllocatorTransactor) BatchDeallocate(opts *bind.TransactOpts, staker common.Address, agents []common.Address, worknetIds []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _AWPAllocator.contract.Transact(opts, "batchDeallocate", staker, agents, worknetIds, amounts)
}

// BatchDeallocate is a paid mutator transaction binding the contract method 0x2e17f308.
//
// Solidity: function batchDeallocate(address staker, address[] agents, uint256[] worknetIds, uint256[] amounts) returns()
func (_AWPAllocator *AWPAllocatorSession) BatchDeallocate(staker common.Address, agents []common.Address, worknetIds []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _AWPAllocator.Contract.BatchDeallocate(&_AWPAllocator.TransactOpts, staker, agents, worknetIds, amounts)
}

// BatchDeallocate is a paid mutator transaction binding the contract method 0x2e17f308.
//
// Solidity: function batchDeallocate(address staker, address[] agents, uint256[] worknetIds, uint256[] amounts) returns()
func (_AWPAllocator *AWPAllocatorTransactorSession) BatchDeallocate(staker common.Address, agents []common.Address, worknetIds []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _AWPAllocator.Contract.BatchDeallocate(&_AWPAllocator.TransactOpts, staker, agents, worknetIds, amounts)
}

// Deallocate is a paid mutator transaction binding the contract method 0x716fb83d.
//
// Solidity: function deallocate(address staker, address agent, uint256 worknetId, uint256 amount) returns()
func (_AWPAllocator *AWPAllocatorTransactor) Deallocate(opts *bind.TransactOpts, staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPAllocator.contract.Transact(opts, "deallocate", staker, agent, worknetId, amount)
}

// Deallocate is a paid mutator transaction binding the contract method 0x716fb83d.
//
// Solidity: function deallocate(address staker, address agent, uint256 worknetId, uint256 amount) returns()
func (_AWPAllocator *AWPAllocatorSession) Deallocate(staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPAllocator.Contract.Deallocate(&_AWPAllocator.TransactOpts, staker, agent, worknetId, amount)
}

// Deallocate is a paid mutator transaction binding the contract method 0x716fb83d.
//
// Solidity: function deallocate(address staker, address agent, uint256 worknetId, uint256 amount) returns()
func (_AWPAllocator *AWPAllocatorTransactorSession) Deallocate(staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPAllocator.Contract.Deallocate(&_AWPAllocator.TransactOpts, staker, agent, worknetId, amount)
}

// DeallocateAll is a paid mutator transaction binding the contract method 0x586ac6b3.
//
// Solidity: function deallocateAll(address staker, address agent, uint256 worknetId) returns()
func (_AWPAllocator *AWPAllocatorTransactor) DeallocateAll(opts *bind.TransactOpts, staker common.Address, agent common.Address, worknetId *big.Int) (*types.Transaction, error) {
	return _AWPAllocator.contract.Transact(opts, "deallocateAll", staker, agent, worknetId)
}

// DeallocateAll is a paid mutator transaction binding the contract method 0x586ac6b3.
//
// Solidity: function deallocateAll(address staker, address agent, uint256 worknetId) returns()
func (_AWPAllocator *AWPAllocatorSession) DeallocateAll(staker common.Address, agent common.Address, worknetId *big.Int) (*types.Transaction, error) {
	return _AWPAllocator.Contract.DeallocateAll(&_AWPAllocator.TransactOpts, staker, agent, worknetId)
}

// DeallocateAll is a paid mutator transaction binding the contract method 0x586ac6b3.
//
// Solidity: function deallocateAll(address staker, address agent, uint256 worknetId) returns()
func (_AWPAllocator *AWPAllocatorTransactorSession) DeallocateAll(staker common.Address, agent common.Address, worknetId *big.Int) (*types.Transaction, error) {
	return _AWPAllocator.Contract.DeallocateAll(&_AWPAllocator.TransactOpts, staker, agent, worknetId)
}

// DeallocateFor is a paid mutator transaction binding the contract method 0x10fe1208.
//
// Solidity: function deallocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPAllocator *AWPAllocatorTransactor) DeallocateFor(opts *bind.TransactOpts, staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPAllocator.contract.Transact(opts, "deallocateFor", staker, agent, worknetId, amount, deadline, v, r, s)
}

// DeallocateFor is a paid mutator transaction binding the contract method 0x10fe1208.
//
// Solidity: function deallocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPAllocator *AWPAllocatorSession) DeallocateFor(staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPAllocator.Contract.DeallocateFor(&_AWPAllocator.TransactOpts, staker, agent, worknetId, amount, deadline, v, r, s)
}

// DeallocateFor is a paid mutator transaction binding the contract method 0x10fe1208.
//
// Solidity: function deallocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPAllocator *AWPAllocatorTransactorSession) DeallocateFor(staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPAllocator.Contract.DeallocateFor(&_AWPAllocator.TransactOpts, staker, agent, worknetId, amount, deadline, v, r, s)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address awpRegistry_, address guardian_) returns()
func (_AWPAllocator *AWPAllocatorTransactor) Initialize(opts *bind.TransactOpts, awpRegistry_ common.Address, guardian_ common.Address) (*types.Transaction, error) {
	return _AWPAllocator.contract.Transact(opts, "initialize", awpRegistry_, guardian_)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address awpRegistry_, address guardian_) returns()
func (_AWPAllocator *AWPAllocatorSession) Initialize(awpRegistry_ common.Address, guardian_ common.Address) (*types.Transaction, error) {
	return _AWPAllocator.Contract.Initialize(&_AWPAllocator.TransactOpts, awpRegistry_, guardian_)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address awpRegistry_, address guardian_) returns()
func (_AWPAllocator *AWPAllocatorTransactorSession) Initialize(awpRegistry_ common.Address, guardian_ common.Address) (*types.Transaction, error) {
	return _AWPAllocator.Contract.Initialize(&_AWPAllocator.TransactOpts, awpRegistry_, guardian_)
}

// Reallocate is a paid mutator transaction binding the contract method 0xd5d5278d.
//
// Solidity: function reallocate(address staker, address fromAgent, uint256 fromWorknetId, address toAgent, uint256 toWorknetId, uint256 amount) returns()
func (_AWPAllocator *AWPAllocatorTransactor) Reallocate(opts *bind.TransactOpts, staker common.Address, fromAgent common.Address, fromWorknetId *big.Int, toAgent common.Address, toWorknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPAllocator.contract.Transact(opts, "reallocate", staker, fromAgent, fromWorknetId, toAgent, toWorknetId, amount)
}

// Reallocate is a paid mutator transaction binding the contract method 0xd5d5278d.
//
// Solidity: function reallocate(address staker, address fromAgent, uint256 fromWorknetId, address toAgent, uint256 toWorknetId, uint256 amount) returns()
func (_AWPAllocator *AWPAllocatorSession) Reallocate(staker common.Address, fromAgent common.Address, fromWorknetId *big.Int, toAgent common.Address, toWorknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPAllocator.Contract.Reallocate(&_AWPAllocator.TransactOpts, staker, fromAgent, fromWorknetId, toAgent, toWorknetId, amount)
}

// Reallocate is a paid mutator transaction binding the contract method 0xd5d5278d.
//
// Solidity: function reallocate(address staker, address fromAgent, uint256 fromWorknetId, address toAgent, uint256 toWorknetId, uint256 amount) returns()
func (_AWPAllocator *AWPAllocatorTransactorSession) Reallocate(staker common.Address, fromAgent common.Address, fromWorknetId *big.Int, toAgent common.Address, toWorknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPAllocator.Contract.Reallocate(&_AWPAllocator.TransactOpts, staker, fromAgent, fromWorknetId, toAgent, toWorknetId, amount)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_AWPAllocator *AWPAllocatorTransactor) SetGuardian(opts *bind.TransactOpts, g common.Address) (*types.Transaction, error) {
	return _AWPAllocator.contract.Transact(opts, "setGuardian", g)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_AWPAllocator *AWPAllocatorSession) SetGuardian(g common.Address) (*types.Transaction, error) {
	return _AWPAllocator.Contract.SetGuardian(&_AWPAllocator.TransactOpts, g)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_AWPAllocator *AWPAllocatorTransactorSession) SetGuardian(g common.Address) (*types.Transaction, error) {
	return _AWPAllocator.Contract.SetGuardian(&_AWPAllocator.TransactOpts, g)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AWPAllocator *AWPAllocatorTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AWPAllocator.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AWPAllocator *AWPAllocatorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AWPAllocator.Contract.UpgradeToAndCall(&_AWPAllocator.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AWPAllocator *AWPAllocatorTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AWPAllocator.Contract.UpgradeToAndCall(&_AWPAllocator.TransactOpts, newImplementation, data)
}

// AWPAllocatorAllocatedIterator is returned from FilterAllocated and is used to iterate over the raw logs and unpacked data for Allocated events raised by the AWPAllocator contract.
type AWPAllocatorAllocatedIterator struct {
	Event *AWPAllocatorAllocated // Event containing the contract specifics and raw log

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
func (it *AWPAllocatorAllocatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPAllocatorAllocated)
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
		it.Event = new(AWPAllocatorAllocated)
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
func (it *AWPAllocatorAllocatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPAllocatorAllocatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPAllocatorAllocated represents a Allocated event raised by the AWPAllocator contract.
type AWPAllocatorAllocated struct {
	Staker    common.Address
	Agent     common.Address
	WorknetId *big.Int
	Amount    *big.Int
	Operator  common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAllocated is a free log retrieval operation binding the contract event 0x655f98c7dae1bab3e2db10cdb4407717b9d219cf2e585adc1edba92d48af2b15.
//
// Solidity: event Allocated(address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator)
func (_AWPAllocator *AWPAllocatorFilterer) FilterAllocated(opts *bind.FilterOpts, staker []common.Address, agent []common.Address) (*AWPAllocatorAllocatedIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _AWPAllocator.contract.FilterLogs(opts, "Allocated", stakerRule, agentRule)
	if err != nil {
		return nil, err
	}
	return &AWPAllocatorAllocatedIterator{contract: _AWPAllocator.contract, event: "Allocated", logs: logs, sub: sub}, nil
}

// WatchAllocated is a free log subscription operation binding the contract event 0x655f98c7dae1bab3e2db10cdb4407717b9d219cf2e585adc1edba92d48af2b15.
//
// Solidity: event Allocated(address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator)
func (_AWPAllocator *AWPAllocatorFilterer) WatchAllocated(opts *bind.WatchOpts, sink chan<- *AWPAllocatorAllocated, staker []common.Address, agent []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _AWPAllocator.contract.WatchLogs(opts, "Allocated", stakerRule, agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPAllocatorAllocated)
				if err := _AWPAllocator.contract.UnpackLog(event, "Allocated", log); err != nil {
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

// ParseAllocated is a log parse operation binding the contract event 0x655f98c7dae1bab3e2db10cdb4407717b9d219cf2e585adc1edba92d48af2b15.
//
// Solidity: event Allocated(address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator)
func (_AWPAllocator *AWPAllocatorFilterer) ParseAllocated(log types.Log) (*AWPAllocatorAllocated, error) {
	event := new(AWPAllocatorAllocated)
	if err := _AWPAllocator.contract.UnpackLog(event, "Allocated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPAllocatorDeallocatedIterator is returned from FilterDeallocated and is used to iterate over the raw logs and unpacked data for Deallocated events raised by the AWPAllocator contract.
type AWPAllocatorDeallocatedIterator struct {
	Event *AWPAllocatorDeallocated // Event containing the contract specifics and raw log

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
func (it *AWPAllocatorDeallocatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPAllocatorDeallocated)
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
		it.Event = new(AWPAllocatorDeallocated)
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
func (it *AWPAllocatorDeallocatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPAllocatorDeallocatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPAllocatorDeallocated represents a Deallocated event raised by the AWPAllocator contract.
type AWPAllocatorDeallocated struct {
	Staker    common.Address
	Agent     common.Address
	WorknetId *big.Int
	Amount    *big.Int
	Operator  common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDeallocated is a free log retrieval operation binding the contract event 0xd55bd7964253d1d9ce9187c8187b1c904274a3f374c9074f6de6fa77746bf345.
//
// Solidity: event Deallocated(address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator)
func (_AWPAllocator *AWPAllocatorFilterer) FilterDeallocated(opts *bind.FilterOpts, staker []common.Address, agent []common.Address) (*AWPAllocatorDeallocatedIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _AWPAllocator.contract.FilterLogs(opts, "Deallocated", stakerRule, agentRule)
	if err != nil {
		return nil, err
	}
	return &AWPAllocatorDeallocatedIterator{contract: _AWPAllocator.contract, event: "Deallocated", logs: logs, sub: sub}, nil
}

// WatchDeallocated is a free log subscription operation binding the contract event 0xd55bd7964253d1d9ce9187c8187b1c904274a3f374c9074f6de6fa77746bf345.
//
// Solidity: event Deallocated(address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator)
func (_AWPAllocator *AWPAllocatorFilterer) WatchDeallocated(opts *bind.WatchOpts, sink chan<- *AWPAllocatorDeallocated, staker []common.Address, agent []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _AWPAllocator.contract.WatchLogs(opts, "Deallocated", stakerRule, agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPAllocatorDeallocated)
				if err := _AWPAllocator.contract.UnpackLog(event, "Deallocated", log); err != nil {
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

// ParseDeallocated is a log parse operation binding the contract event 0xd55bd7964253d1d9ce9187c8187b1c904274a3f374c9074f6de6fa77746bf345.
//
// Solidity: event Deallocated(address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator)
func (_AWPAllocator *AWPAllocatorFilterer) ParseDeallocated(log types.Log) (*AWPAllocatorDeallocated, error) {
	event := new(AWPAllocatorDeallocated)
	if err := _AWPAllocator.contract.UnpackLog(event, "Deallocated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPAllocatorEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the AWPAllocator contract.
type AWPAllocatorEIP712DomainChangedIterator struct {
	Event *AWPAllocatorEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *AWPAllocatorEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPAllocatorEIP712DomainChanged)
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
		it.Event = new(AWPAllocatorEIP712DomainChanged)
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
func (it *AWPAllocatorEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPAllocatorEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPAllocatorEIP712DomainChanged represents a EIP712DomainChanged event raised by the AWPAllocator contract.
type AWPAllocatorEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_AWPAllocator *AWPAllocatorFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*AWPAllocatorEIP712DomainChangedIterator, error) {

	logs, sub, err := _AWPAllocator.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &AWPAllocatorEIP712DomainChangedIterator{contract: _AWPAllocator.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_AWPAllocator *AWPAllocatorFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *AWPAllocatorEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _AWPAllocator.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPAllocatorEIP712DomainChanged)
				if err := _AWPAllocator.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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
func (_AWPAllocator *AWPAllocatorFilterer) ParseEIP712DomainChanged(log types.Log) (*AWPAllocatorEIP712DomainChanged, error) {
	event := new(AWPAllocatorEIP712DomainChanged)
	if err := _AWPAllocator.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPAllocatorGuardianUpdatedIterator is returned from FilterGuardianUpdated and is used to iterate over the raw logs and unpacked data for GuardianUpdated events raised by the AWPAllocator contract.
type AWPAllocatorGuardianUpdatedIterator struct {
	Event *AWPAllocatorGuardianUpdated // Event containing the contract specifics and raw log

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
func (it *AWPAllocatorGuardianUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPAllocatorGuardianUpdated)
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
		it.Event = new(AWPAllocatorGuardianUpdated)
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
func (it *AWPAllocatorGuardianUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPAllocatorGuardianUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPAllocatorGuardianUpdated represents a GuardianUpdated event raised by the AWPAllocator contract.
type AWPAllocatorGuardianUpdated struct {
	NewGuardian common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterGuardianUpdated is a free log retrieval operation binding the contract event 0x6bb7ff33e730289800c62ad882105a144a74010d2bdbb9a942544a3005ad55bf.
//
// Solidity: event GuardianUpdated(address indexed newGuardian)
func (_AWPAllocator *AWPAllocatorFilterer) FilterGuardianUpdated(opts *bind.FilterOpts, newGuardian []common.Address) (*AWPAllocatorGuardianUpdatedIterator, error) {

	var newGuardianRule []interface{}
	for _, newGuardianItem := range newGuardian {
		newGuardianRule = append(newGuardianRule, newGuardianItem)
	}

	logs, sub, err := _AWPAllocator.contract.FilterLogs(opts, "GuardianUpdated", newGuardianRule)
	if err != nil {
		return nil, err
	}
	return &AWPAllocatorGuardianUpdatedIterator{contract: _AWPAllocator.contract, event: "GuardianUpdated", logs: logs, sub: sub}, nil
}

// WatchGuardianUpdated is a free log subscription operation binding the contract event 0x6bb7ff33e730289800c62ad882105a144a74010d2bdbb9a942544a3005ad55bf.
//
// Solidity: event GuardianUpdated(address indexed newGuardian)
func (_AWPAllocator *AWPAllocatorFilterer) WatchGuardianUpdated(opts *bind.WatchOpts, sink chan<- *AWPAllocatorGuardianUpdated, newGuardian []common.Address) (event.Subscription, error) {

	var newGuardianRule []interface{}
	for _, newGuardianItem := range newGuardian {
		newGuardianRule = append(newGuardianRule, newGuardianItem)
	}

	logs, sub, err := _AWPAllocator.contract.WatchLogs(opts, "GuardianUpdated", newGuardianRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPAllocatorGuardianUpdated)
				if err := _AWPAllocator.contract.UnpackLog(event, "GuardianUpdated", log); err != nil {
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
func (_AWPAllocator *AWPAllocatorFilterer) ParseGuardianUpdated(log types.Log) (*AWPAllocatorGuardianUpdated, error) {
	event := new(AWPAllocatorGuardianUpdated)
	if err := _AWPAllocator.contract.UnpackLog(event, "GuardianUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPAllocatorInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the AWPAllocator contract.
type AWPAllocatorInitializedIterator struct {
	Event *AWPAllocatorInitialized // Event containing the contract specifics and raw log

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
func (it *AWPAllocatorInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPAllocatorInitialized)
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
		it.Event = new(AWPAllocatorInitialized)
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
func (it *AWPAllocatorInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPAllocatorInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPAllocatorInitialized represents a Initialized event raised by the AWPAllocator contract.
type AWPAllocatorInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AWPAllocator *AWPAllocatorFilterer) FilterInitialized(opts *bind.FilterOpts) (*AWPAllocatorInitializedIterator, error) {

	logs, sub, err := _AWPAllocator.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &AWPAllocatorInitializedIterator{contract: _AWPAllocator.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AWPAllocator *AWPAllocatorFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *AWPAllocatorInitialized) (event.Subscription, error) {

	logs, sub, err := _AWPAllocator.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPAllocatorInitialized)
				if err := _AWPAllocator.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_AWPAllocator *AWPAllocatorFilterer) ParseInitialized(log types.Log) (*AWPAllocatorInitialized, error) {
	event := new(AWPAllocatorInitialized)
	if err := _AWPAllocator.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPAllocatorReallocatedIterator is returned from FilterReallocated and is used to iterate over the raw logs and unpacked data for Reallocated events raised by the AWPAllocator contract.
type AWPAllocatorReallocatedIterator struct {
	Event *AWPAllocatorReallocated // Event containing the contract specifics and raw log

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
func (it *AWPAllocatorReallocatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPAllocatorReallocated)
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
		it.Event = new(AWPAllocatorReallocated)
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
func (it *AWPAllocatorReallocatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPAllocatorReallocatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPAllocatorReallocated represents a Reallocated event raised by the AWPAllocator contract.
type AWPAllocatorReallocated struct {
	Staker        common.Address
	FromAgent     common.Address
	FromWorknetId *big.Int
	ToAgent       common.Address
	ToWorknetId   *big.Int
	Amount        *big.Int
	Operator      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterReallocated is a free log retrieval operation binding the contract event 0x726c93ba67bfe4c677e37114279f0ad9aab5ee9ffbd1158923be5d0fec3b1b45.
//
// Solidity: event Reallocated(address indexed staker, address fromAgent, uint256 fromWorknetId, address toAgent, uint256 toWorknetId, uint256 amount, address operator)
func (_AWPAllocator *AWPAllocatorFilterer) FilterReallocated(opts *bind.FilterOpts, staker []common.Address) (*AWPAllocatorReallocatedIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}

	logs, sub, err := _AWPAllocator.contract.FilterLogs(opts, "Reallocated", stakerRule)
	if err != nil {
		return nil, err
	}
	return &AWPAllocatorReallocatedIterator{contract: _AWPAllocator.contract, event: "Reallocated", logs: logs, sub: sub}, nil
}

// WatchReallocated is a free log subscription operation binding the contract event 0x726c93ba67bfe4c677e37114279f0ad9aab5ee9ffbd1158923be5d0fec3b1b45.
//
// Solidity: event Reallocated(address indexed staker, address fromAgent, uint256 fromWorknetId, address toAgent, uint256 toWorknetId, uint256 amount, address operator)
func (_AWPAllocator *AWPAllocatorFilterer) WatchReallocated(opts *bind.WatchOpts, sink chan<- *AWPAllocatorReallocated, staker []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}

	logs, sub, err := _AWPAllocator.contract.WatchLogs(opts, "Reallocated", stakerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPAllocatorReallocated)
				if err := _AWPAllocator.contract.UnpackLog(event, "Reallocated", log); err != nil {
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

// ParseReallocated is a log parse operation binding the contract event 0x726c93ba67bfe4c677e37114279f0ad9aab5ee9ffbd1158923be5d0fec3b1b45.
//
// Solidity: event Reallocated(address indexed staker, address fromAgent, uint256 fromWorknetId, address toAgent, uint256 toWorknetId, uint256 amount, address operator)
func (_AWPAllocator *AWPAllocatorFilterer) ParseReallocated(log types.Log) (*AWPAllocatorReallocated, error) {
	event := new(AWPAllocatorReallocated)
	if err := _AWPAllocator.contract.UnpackLog(event, "Reallocated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPAllocatorUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the AWPAllocator contract.
type AWPAllocatorUpgradedIterator struct {
	Event *AWPAllocatorUpgraded // Event containing the contract specifics and raw log

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
func (it *AWPAllocatorUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPAllocatorUpgraded)
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
		it.Event = new(AWPAllocatorUpgraded)
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
func (it *AWPAllocatorUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPAllocatorUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPAllocatorUpgraded represents a Upgraded event raised by the AWPAllocator contract.
type AWPAllocatorUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AWPAllocator *AWPAllocatorFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*AWPAllocatorUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AWPAllocator.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &AWPAllocatorUpgradedIterator{contract: _AWPAllocator.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AWPAllocator *AWPAllocatorFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *AWPAllocatorUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AWPAllocator.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPAllocatorUpgraded)
				if err := _AWPAllocator.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_AWPAllocator *AWPAllocatorFilterer) ParseUpgraded(log types.Log) (*AWPAllocatorUpgraded, error) {
	event := new(AWPAllocatorUpgraded)
	if err := _AWPAllocator.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
