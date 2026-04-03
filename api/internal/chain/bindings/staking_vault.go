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

// StakingVaultMetaData contains all meta data concerning the StakingVault contract.
var StakingVaultMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allocate\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"worknetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allocateFor\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"worknetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"awpRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deallocate\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"worknetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deallocateFor\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"worknetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAgentStake\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"worknetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAgentWorknets\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWorknetTotalStake\",\"inputs\":[{\"name\":\"worknetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"guardian\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"awpRegistry_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"guardian_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"nonces\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reallocate\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromAgent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromWorknetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"toAgent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"toWorknetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setGuardian\",\"inputs\":[{\"name\":\"g\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setStakeNFT\",\"inputs\":[{\"name\":\"stakeNFT_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakeNFT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"userTotalAllocated\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"worknetTotalStake\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Allocated\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"worknetId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deallocated\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"worknetId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GuardianUpdated\",\"inputs\":[{\"name\":\"newGuardian\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Reallocated\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fromAgent\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"fromWorknetId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"toAgent\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"toWorknetId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakeNFTSet\",\"inputs\":[{\"name\":\"stakeNFT\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AllocationOverflow\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadySet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AmountExceedsUint128\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpiredSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientAllocation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientUnallocated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAWPRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAuthorized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotGuardian\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroWorknetId\",\"inputs\":[]}]",
}

// StakingVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingVaultMetaData.ABI instead.
var StakingVaultABI = StakingVaultMetaData.ABI

// StakingVault is an auto generated Go binding around an Ethereum contract.
type StakingVault struct {
	StakingVaultCaller     // Read-only binding to the contract
	StakingVaultTransactor // Write-only binding to the contract
	StakingVaultFilterer   // Log filterer for contract events
}

// StakingVaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingVaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingVaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingVaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingVaultSession struct {
	Contract     *StakingVault     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingVaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingVaultCallerSession struct {
	Contract *StakingVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// StakingVaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingVaultTransactorSession struct {
	Contract     *StakingVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// StakingVaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingVaultRaw struct {
	Contract *StakingVault // Generic contract binding to access the raw methods on
}

// StakingVaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingVaultCallerRaw struct {
	Contract *StakingVaultCaller // Generic read-only contract binding to access the raw methods on
}

// StakingVaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingVaultTransactorRaw struct {
	Contract *StakingVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakingVault creates a new instance of StakingVault, bound to a specific deployed contract.
func NewStakingVault(address common.Address, backend bind.ContractBackend) (*StakingVault, error) {
	contract, err := bindStakingVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakingVault{StakingVaultCaller: StakingVaultCaller{contract: contract}, StakingVaultTransactor: StakingVaultTransactor{contract: contract}, StakingVaultFilterer: StakingVaultFilterer{contract: contract}}, nil
}

// NewStakingVaultCaller creates a new read-only instance of StakingVault, bound to a specific deployed contract.
func NewStakingVaultCaller(address common.Address, caller bind.ContractCaller) (*StakingVaultCaller, error) {
	contract, err := bindStakingVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingVaultCaller{contract: contract}, nil
}

// NewStakingVaultTransactor creates a new write-only instance of StakingVault, bound to a specific deployed contract.
func NewStakingVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingVaultTransactor, error) {
	contract, err := bindStakingVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingVaultTransactor{contract: contract}, nil
}

// NewStakingVaultFilterer creates a new log filterer instance of StakingVault, bound to a specific deployed contract.
func NewStakingVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingVaultFilterer, error) {
	contract, err := bindStakingVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingVaultFilterer{contract: contract}, nil
}

// bindStakingVault binds a generic wrapper to an already deployed contract.
func bindStakingVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StakingVaultMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingVault *StakingVaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingVault.Contract.StakingVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingVault *StakingVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingVault.Contract.StakingVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingVault *StakingVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingVault.Contract.StakingVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingVault *StakingVaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingVault *StakingVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingVault *StakingVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingVault.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_StakingVault *StakingVaultCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _StakingVault.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_StakingVault *StakingVaultSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _StakingVault.Contract.UPGRADEINTERFACEVERSION(&_StakingVault.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_StakingVault *StakingVaultCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _StakingVault.Contract.UPGRADEINTERFACEVERSION(&_StakingVault.CallOpts)
}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_StakingVault *StakingVaultCaller) AwpRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakingVault.contract.Call(opts, &out, "awpRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_StakingVault *StakingVaultSession) AwpRegistry() (common.Address, error) {
	return _StakingVault.Contract.AwpRegistry(&_StakingVault.CallOpts)
}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_StakingVault *StakingVaultCallerSession) AwpRegistry() (common.Address, error) {
	return _StakingVault.Contract.AwpRegistry(&_StakingVault.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_StakingVault *StakingVaultCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _StakingVault.contract.Call(opts, &out, "eip712Domain")

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
func (_StakingVault *StakingVaultSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _StakingVault.Contract.Eip712Domain(&_StakingVault.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_StakingVault *StakingVaultCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _StakingVault.Contract.Eip712Domain(&_StakingVault.CallOpts)
}

// GetAgentStake is a free data retrieval call binding the contract method 0xf1ad80c6.
//
// Solidity: function getAgentStake(address user, address agent, uint256 worknetId) view returns(uint256)
func (_StakingVault *StakingVaultCaller) GetAgentStake(opts *bind.CallOpts, user common.Address, agent common.Address, worknetId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _StakingVault.contract.Call(opts, &out, "getAgentStake", user, agent, worknetId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAgentStake is a free data retrieval call binding the contract method 0xf1ad80c6.
//
// Solidity: function getAgentStake(address user, address agent, uint256 worknetId) view returns(uint256)
func (_StakingVault *StakingVaultSession) GetAgentStake(user common.Address, agent common.Address, worknetId *big.Int) (*big.Int, error) {
	return _StakingVault.Contract.GetAgentStake(&_StakingVault.CallOpts, user, agent, worknetId)
}

// GetAgentStake is a free data retrieval call binding the contract method 0xf1ad80c6.
//
// Solidity: function getAgentStake(address user, address agent, uint256 worknetId) view returns(uint256)
func (_StakingVault *StakingVaultCallerSession) GetAgentStake(user common.Address, agent common.Address, worknetId *big.Int) (*big.Int, error) {
	return _StakingVault.Contract.GetAgentStake(&_StakingVault.CallOpts, user, agent, worknetId)
}

// GetAgentWorknets is a free data retrieval call binding the contract method 0xfd4fd2e8.
//
// Solidity: function getAgentWorknets(address user, address agent) view returns(uint256[])
func (_StakingVault *StakingVaultCaller) GetAgentWorknets(opts *bind.CallOpts, user common.Address, agent common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _StakingVault.contract.Call(opts, &out, "getAgentWorknets", user, agent)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetAgentWorknets is a free data retrieval call binding the contract method 0xfd4fd2e8.
//
// Solidity: function getAgentWorknets(address user, address agent) view returns(uint256[])
func (_StakingVault *StakingVaultSession) GetAgentWorknets(user common.Address, agent common.Address) ([]*big.Int, error) {
	return _StakingVault.Contract.GetAgentWorknets(&_StakingVault.CallOpts, user, agent)
}

// GetAgentWorknets is a free data retrieval call binding the contract method 0xfd4fd2e8.
//
// Solidity: function getAgentWorknets(address user, address agent) view returns(uint256[])
func (_StakingVault *StakingVaultCallerSession) GetAgentWorknets(user common.Address, agent common.Address) ([]*big.Int, error) {
	return _StakingVault.Contract.GetAgentWorknets(&_StakingVault.CallOpts, user, agent)
}

// GetWorknetTotalStake is a free data retrieval call binding the contract method 0x5bfe237f.
//
// Solidity: function getWorknetTotalStake(uint256 worknetId) view returns(uint256)
func (_StakingVault *StakingVaultCaller) GetWorknetTotalStake(opts *bind.CallOpts, worknetId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _StakingVault.contract.Call(opts, &out, "getWorknetTotalStake", worknetId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetWorknetTotalStake is a free data retrieval call binding the contract method 0x5bfe237f.
//
// Solidity: function getWorknetTotalStake(uint256 worknetId) view returns(uint256)
func (_StakingVault *StakingVaultSession) GetWorknetTotalStake(worknetId *big.Int) (*big.Int, error) {
	return _StakingVault.Contract.GetWorknetTotalStake(&_StakingVault.CallOpts, worknetId)
}

// GetWorknetTotalStake is a free data retrieval call binding the contract method 0x5bfe237f.
//
// Solidity: function getWorknetTotalStake(uint256 worknetId) view returns(uint256)
func (_StakingVault *StakingVaultCallerSession) GetWorknetTotalStake(worknetId *big.Int) (*big.Int, error) {
	return _StakingVault.Contract.GetWorknetTotalStake(&_StakingVault.CallOpts, worknetId)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_StakingVault *StakingVaultCaller) Guardian(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakingVault.contract.Call(opts, &out, "guardian")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_StakingVault *StakingVaultSession) Guardian() (common.Address, error) {
	return _StakingVault.Contract.Guardian(&_StakingVault.CallOpts)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_StakingVault *StakingVaultCallerSession) Guardian() (common.Address, error) {
	return _StakingVault.Contract.Guardian(&_StakingVault.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_StakingVault *StakingVaultCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakingVault.contract.Call(opts, &out, "nonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_StakingVault *StakingVaultSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _StakingVault.Contract.Nonces(&_StakingVault.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_StakingVault *StakingVaultCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _StakingVault.Contract.Nonces(&_StakingVault.CallOpts, arg0)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_StakingVault *StakingVaultCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _StakingVault.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_StakingVault *StakingVaultSession) ProxiableUUID() ([32]byte, error) {
	return _StakingVault.Contract.ProxiableUUID(&_StakingVault.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_StakingVault *StakingVaultCallerSession) ProxiableUUID() ([32]byte, error) {
	return _StakingVault.Contract.ProxiableUUID(&_StakingVault.CallOpts)
}

// StakeNFT is a free data retrieval call binding the contract method 0xb48509e6.
//
// Solidity: function stakeNFT() view returns(address)
func (_StakingVault *StakingVaultCaller) StakeNFT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakingVault.contract.Call(opts, &out, "stakeNFT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakeNFT is a free data retrieval call binding the contract method 0xb48509e6.
//
// Solidity: function stakeNFT() view returns(address)
func (_StakingVault *StakingVaultSession) StakeNFT() (common.Address, error) {
	return _StakingVault.Contract.StakeNFT(&_StakingVault.CallOpts)
}

// StakeNFT is a free data retrieval call binding the contract method 0xb48509e6.
//
// Solidity: function stakeNFT() view returns(address)
func (_StakingVault *StakingVaultCallerSession) StakeNFT() (common.Address, error) {
	return _StakingVault.Contract.StakeNFT(&_StakingVault.CallOpts)
}

// UserTotalAllocated is a free data retrieval call binding the contract method 0x32ffa4ce.
//
// Solidity: function userTotalAllocated(address ) view returns(uint256)
func (_StakingVault *StakingVaultCaller) UserTotalAllocated(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakingVault.contract.Call(opts, &out, "userTotalAllocated", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UserTotalAllocated is a free data retrieval call binding the contract method 0x32ffa4ce.
//
// Solidity: function userTotalAllocated(address ) view returns(uint256)
func (_StakingVault *StakingVaultSession) UserTotalAllocated(arg0 common.Address) (*big.Int, error) {
	return _StakingVault.Contract.UserTotalAllocated(&_StakingVault.CallOpts, arg0)
}

// UserTotalAllocated is a free data retrieval call binding the contract method 0x32ffa4ce.
//
// Solidity: function userTotalAllocated(address ) view returns(uint256)
func (_StakingVault *StakingVaultCallerSession) UserTotalAllocated(arg0 common.Address) (*big.Int, error) {
	return _StakingVault.Contract.UserTotalAllocated(&_StakingVault.CallOpts, arg0)
}

// WorknetTotalStake is a free data retrieval call binding the contract method 0x8c00c09c.
//
// Solidity: function worknetTotalStake(uint256 ) view returns(uint256)
func (_StakingVault *StakingVaultCaller) WorknetTotalStake(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _StakingVault.contract.Call(opts, &out, "worknetTotalStake", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WorknetTotalStake is a free data retrieval call binding the contract method 0x8c00c09c.
//
// Solidity: function worknetTotalStake(uint256 ) view returns(uint256)
func (_StakingVault *StakingVaultSession) WorknetTotalStake(arg0 *big.Int) (*big.Int, error) {
	return _StakingVault.Contract.WorknetTotalStake(&_StakingVault.CallOpts, arg0)
}

// WorknetTotalStake is a free data retrieval call binding the contract method 0x8c00c09c.
//
// Solidity: function worknetTotalStake(uint256 ) view returns(uint256)
func (_StakingVault *StakingVaultCallerSession) WorknetTotalStake(arg0 *big.Int) (*big.Int, error) {
	return _StakingVault.Contract.WorknetTotalStake(&_StakingVault.CallOpts, arg0)
}

// Allocate is a paid mutator transaction binding the contract method 0xd035a9a7.
//
// Solidity: function allocate(address staker, address agent, uint256 worknetId, uint256 amount) returns()
func (_StakingVault *StakingVaultTransactor) Allocate(opts *bind.TransactOpts, staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.contract.Transact(opts, "allocate", staker, agent, worknetId, amount)
}

// Allocate is a paid mutator transaction binding the contract method 0xd035a9a7.
//
// Solidity: function allocate(address staker, address agent, uint256 worknetId, uint256 amount) returns()
func (_StakingVault *StakingVaultSession) Allocate(staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.Contract.Allocate(&_StakingVault.TransactOpts, staker, agent, worknetId, amount)
}

// Allocate is a paid mutator transaction binding the contract method 0xd035a9a7.
//
// Solidity: function allocate(address staker, address agent, uint256 worknetId, uint256 amount) returns()
func (_StakingVault *StakingVaultTransactorSession) Allocate(staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.Contract.Allocate(&_StakingVault.TransactOpts, staker, agent, worknetId, amount)
}

// AllocateFor is a paid mutator transaction binding the contract method 0x7d66c5c5.
//
// Solidity: function allocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_StakingVault *StakingVaultTransactor) AllocateFor(opts *bind.TransactOpts, staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _StakingVault.contract.Transact(opts, "allocateFor", staker, agent, worknetId, amount, deadline, v, r, s)
}

// AllocateFor is a paid mutator transaction binding the contract method 0x7d66c5c5.
//
// Solidity: function allocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_StakingVault *StakingVaultSession) AllocateFor(staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _StakingVault.Contract.AllocateFor(&_StakingVault.TransactOpts, staker, agent, worknetId, amount, deadline, v, r, s)
}

// AllocateFor is a paid mutator transaction binding the contract method 0x7d66c5c5.
//
// Solidity: function allocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_StakingVault *StakingVaultTransactorSession) AllocateFor(staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _StakingVault.Contract.AllocateFor(&_StakingVault.TransactOpts, staker, agent, worknetId, amount, deadline, v, r, s)
}

// Deallocate is a paid mutator transaction binding the contract method 0x716fb83d.
//
// Solidity: function deallocate(address staker, address agent, uint256 worknetId, uint256 amount) returns()
func (_StakingVault *StakingVaultTransactor) Deallocate(opts *bind.TransactOpts, staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.contract.Transact(opts, "deallocate", staker, agent, worknetId, amount)
}

// Deallocate is a paid mutator transaction binding the contract method 0x716fb83d.
//
// Solidity: function deallocate(address staker, address agent, uint256 worknetId, uint256 amount) returns()
func (_StakingVault *StakingVaultSession) Deallocate(staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.Contract.Deallocate(&_StakingVault.TransactOpts, staker, agent, worknetId, amount)
}

// Deallocate is a paid mutator transaction binding the contract method 0x716fb83d.
//
// Solidity: function deallocate(address staker, address agent, uint256 worknetId, uint256 amount) returns()
func (_StakingVault *StakingVaultTransactorSession) Deallocate(staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.Contract.Deallocate(&_StakingVault.TransactOpts, staker, agent, worknetId, amount)
}

// DeallocateFor is a paid mutator transaction binding the contract method 0x10fe1208.
//
// Solidity: function deallocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_StakingVault *StakingVaultTransactor) DeallocateFor(opts *bind.TransactOpts, staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _StakingVault.contract.Transact(opts, "deallocateFor", staker, agent, worknetId, amount, deadline, v, r, s)
}

// DeallocateFor is a paid mutator transaction binding the contract method 0x10fe1208.
//
// Solidity: function deallocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_StakingVault *StakingVaultSession) DeallocateFor(staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _StakingVault.Contract.DeallocateFor(&_StakingVault.TransactOpts, staker, agent, worknetId, amount, deadline, v, r, s)
}

// DeallocateFor is a paid mutator transaction binding the contract method 0x10fe1208.
//
// Solidity: function deallocateFor(address staker, address agent, uint256 worknetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_StakingVault *StakingVaultTransactorSession) DeallocateFor(staker common.Address, agent common.Address, worknetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _StakingVault.Contract.DeallocateFor(&_StakingVault.TransactOpts, staker, agent, worknetId, amount, deadline, v, r, s)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address awpRegistry_, address guardian_) returns()
func (_StakingVault *StakingVaultTransactor) Initialize(opts *bind.TransactOpts, awpRegistry_ common.Address, guardian_ common.Address) (*types.Transaction, error) {
	return _StakingVault.contract.Transact(opts, "initialize", awpRegistry_, guardian_)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address awpRegistry_, address guardian_) returns()
func (_StakingVault *StakingVaultSession) Initialize(awpRegistry_ common.Address, guardian_ common.Address) (*types.Transaction, error) {
	return _StakingVault.Contract.Initialize(&_StakingVault.TransactOpts, awpRegistry_, guardian_)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address awpRegistry_, address guardian_) returns()
func (_StakingVault *StakingVaultTransactorSession) Initialize(awpRegistry_ common.Address, guardian_ common.Address) (*types.Transaction, error) {
	return _StakingVault.Contract.Initialize(&_StakingVault.TransactOpts, awpRegistry_, guardian_)
}

// Reallocate is a paid mutator transaction binding the contract method 0xd5d5278d.
//
// Solidity: function reallocate(address staker, address fromAgent, uint256 fromWorknetId, address toAgent, uint256 toWorknetId, uint256 amount) returns()
func (_StakingVault *StakingVaultTransactor) Reallocate(opts *bind.TransactOpts, staker common.Address, fromAgent common.Address, fromWorknetId *big.Int, toAgent common.Address, toWorknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.contract.Transact(opts, "reallocate", staker, fromAgent, fromWorknetId, toAgent, toWorknetId, amount)
}

// Reallocate is a paid mutator transaction binding the contract method 0xd5d5278d.
//
// Solidity: function reallocate(address staker, address fromAgent, uint256 fromWorknetId, address toAgent, uint256 toWorknetId, uint256 amount) returns()
func (_StakingVault *StakingVaultSession) Reallocate(staker common.Address, fromAgent common.Address, fromWorknetId *big.Int, toAgent common.Address, toWorknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.Contract.Reallocate(&_StakingVault.TransactOpts, staker, fromAgent, fromWorknetId, toAgent, toWorknetId, amount)
}

// Reallocate is a paid mutator transaction binding the contract method 0xd5d5278d.
//
// Solidity: function reallocate(address staker, address fromAgent, uint256 fromWorknetId, address toAgent, uint256 toWorknetId, uint256 amount) returns()
func (_StakingVault *StakingVaultTransactorSession) Reallocate(staker common.Address, fromAgent common.Address, fromWorknetId *big.Int, toAgent common.Address, toWorknetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.Contract.Reallocate(&_StakingVault.TransactOpts, staker, fromAgent, fromWorknetId, toAgent, toWorknetId, amount)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_StakingVault *StakingVaultTransactor) SetGuardian(opts *bind.TransactOpts, g common.Address) (*types.Transaction, error) {
	return _StakingVault.contract.Transact(opts, "setGuardian", g)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_StakingVault *StakingVaultSession) SetGuardian(g common.Address) (*types.Transaction, error) {
	return _StakingVault.Contract.SetGuardian(&_StakingVault.TransactOpts, g)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_StakingVault *StakingVaultTransactorSession) SetGuardian(g common.Address) (*types.Transaction, error) {
	return _StakingVault.Contract.SetGuardian(&_StakingVault.TransactOpts, g)
}

// SetStakeNFT is a paid mutator transaction binding the contract method 0x48f069ec.
//
// Solidity: function setStakeNFT(address stakeNFT_) returns()
func (_StakingVault *StakingVaultTransactor) SetStakeNFT(opts *bind.TransactOpts, stakeNFT_ common.Address) (*types.Transaction, error) {
	return _StakingVault.contract.Transact(opts, "setStakeNFT", stakeNFT_)
}

// SetStakeNFT is a paid mutator transaction binding the contract method 0x48f069ec.
//
// Solidity: function setStakeNFT(address stakeNFT_) returns()
func (_StakingVault *StakingVaultSession) SetStakeNFT(stakeNFT_ common.Address) (*types.Transaction, error) {
	return _StakingVault.Contract.SetStakeNFT(&_StakingVault.TransactOpts, stakeNFT_)
}

// SetStakeNFT is a paid mutator transaction binding the contract method 0x48f069ec.
//
// Solidity: function setStakeNFT(address stakeNFT_) returns()
func (_StakingVault *StakingVaultTransactorSession) SetStakeNFT(stakeNFT_ common.Address) (*types.Transaction, error) {
	return _StakingVault.Contract.SetStakeNFT(&_StakingVault.TransactOpts, stakeNFT_)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_StakingVault *StakingVaultTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _StakingVault.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_StakingVault *StakingVaultSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _StakingVault.Contract.UpgradeToAndCall(&_StakingVault.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_StakingVault *StakingVaultTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _StakingVault.Contract.UpgradeToAndCall(&_StakingVault.TransactOpts, newImplementation, data)
}

// StakingVaultAllocatedIterator is returned from FilterAllocated and is used to iterate over the raw logs and unpacked data for Allocated events raised by the StakingVault contract.
type StakingVaultAllocatedIterator struct {
	Event *StakingVaultAllocated // Event containing the contract specifics and raw log

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
func (it *StakingVaultAllocatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingVaultAllocated)
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
		it.Event = new(StakingVaultAllocated)
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
func (it *StakingVaultAllocatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingVaultAllocatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingVaultAllocated represents a Allocated event raised by the StakingVault contract.
type StakingVaultAllocated struct {
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
func (_StakingVault *StakingVaultFilterer) FilterAllocated(opts *bind.FilterOpts, staker []common.Address, agent []common.Address) (*StakingVaultAllocatedIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _StakingVault.contract.FilterLogs(opts, "Allocated", stakerRule, agentRule)
	if err != nil {
		return nil, err
	}
	return &StakingVaultAllocatedIterator{contract: _StakingVault.contract, event: "Allocated", logs: logs, sub: sub}, nil
}

// WatchAllocated is a free log subscription operation binding the contract event 0x655f98c7dae1bab3e2db10cdb4407717b9d219cf2e585adc1edba92d48af2b15.
//
// Solidity: event Allocated(address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator)
func (_StakingVault *StakingVaultFilterer) WatchAllocated(opts *bind.WatchOpts, sink chan<- *StakingVaultAllocated, staker []common.Address, agent []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _StakingVault.contract.WatchLogs(opts, "Allocated", stakerRule, agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingVaultAllocated)
				if err := _StakingVault.contract.UnpackLog(event, "Allocated", log); err != nil {
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
func (_StakingVault *StakingVaultFilterer) ParseAllocated(log types.Log) (*StakingVaultAllocated, error) {
	event := new(StakingVaultAllocated)
	if err := _StakingVault.contract.UnpackLog(event, "Allocated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingVaultDeallocatedIterator is returned from FilterDeallocated and is used to iterate over the raw logs and unpacked data for Deallocated events raised by the StakingVault contract.
type StakingVaultDeallocatedIterator struct {
	Event *StakingVaultDeallocated // Event containing the contract specifics and raw log

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
func (it *StakingVaultDeallocatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingVaultDeallocated)
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
		it.Event = new(StakingVaultDeallocated)
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
func (it *StakingVaultDeallocatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingVaultDeallocatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingVaultDeallocated represents a Deallocated event raised by the StakingVault contract.
type StakingVaultDeallocated struct {
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
func (_StakingVault *StakingVaultFilterer) FilterDeallocated(opts *bind.FilterOpts, staker []common.Address, agent []common.Address) (*StakingVaultDeallocatedIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _StakingVault.contract.FilterLogs(opts, "Deallocated", stakerRule, agentRule)
	if err != nil {
		return nil, err
	}
	return &StakingVaultDeallocatedIterator{contract: _StakingVault.contract, event: "Deallocated", logs: logs, sub: sub}, nil
}

// WatchDeallocated is a free log subscription operation binding the contract event 0xd55bd7964253d1d9ce9187c8187b1c904274a3f374c9074f6de6fa77746bf345.
//
// Solidity: event Deallocated(address indexed staker, address indexed agent, uint256 worknetId, uint256 amount, address operator)
func (_StakingVault *StakingVaultFilterer) WatchDeallocated(opts *bind.WatchOpts, sink chan<- *StakingVaultDeallocated, staker []common.Address, agent []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _StakingVault.contract.WatchLogs(opts, "Deallocated", stakerRule, agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingVaultDeallocated)
				if err := _StakingVault.contract.UnpackLog(event, "Deallocated", log); err != nil {
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
func (_StakingVault *StakingVaultFilterer) ParseDeallocated(log types.Log) (*StakingVaultDeallocated, error) {
	event := new(StakingVaultDeallocated)
	if err := _StakingVault.contract.UnpackLog(event, "Deallocated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingVaultEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the StakingVault contract.
type StakingVaultEIP712DomainChangedIterator struct {
	Event *StakingVaultEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *StakingVaultEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingVaultEIP712DomainChanged)
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
		it.Event = new(StakingVaultEIP712DomainChanged)
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
func (it *StakingVaultEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingVaultEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingVaultEIP712DomainChanged represents a EIP712DomainChanged event raised by the StakingVault contract.
type StakingVaultEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_StakingVault *StakingVaultFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*StakingVaultEIP712DomainChangedIterator, error) {

	logs, sub, err := _StakingVault.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &StakingVaultEIP712DomainChangedIterator{contract: _StakingVault.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_StakingVault *StakingVaultFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *StakingVaultEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _StakingVault.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingVaultEIP712DomainChanged)
				if err := _StakingVault.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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
func (_StakingVault *StakingVaultFilterer) ParseEIP712DomainChanged(log types.Log) (*StakingVaultEIP712DomainChanged, error) {
	event := new(StakingVaultEIP712DomainChanged)
	if err := _StakingVault.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingVaultGuardianUpdatedIterator is returned from FilterGuardianUpdated and is used to iterate over the raw logs and unpacked data for GuardianUpdated events raised by the StakingVault contract.
type StakingVaultGuardianUpdatedIterator struct {
	Event *StakingVaultGuardianUpdated // Event containing the contract specifics and raw log

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
func (it *StakingVaultGuardianUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingVaultGuardianUpdated)
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
		it.Event = new(StakingVaultGuardianUpdated)
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
func (it *StakingVaultGuardianUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingVaultGuardianUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingVaultGuardianUpdated represents a GuardianUpdated event raised by the StakingVault contract.
type StakingVaultGuardianUpdated struct {
	NewGuardian common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterGuardianUpdated is a free log retrieval operation binding the contract event 0x6bb7ff33e730289800c62ad882105a144a74010d2bdbb9a942544a3005ad55bf.
//
// Solidity: event GuardianUpdated(address indexed newGuardian)
func (_StakingVault *StakingVaultFilterer) FilterGuardianUpdated(opts *bind.FilterOpts, newGuardian []common.Address) (*StakingVaultGuardianUpdatedIterator, error) {

	var newGuardianRule []interface{}
	for _, newGuardianItem := range newGuardian {
		newGuardianRule = append(newGuardianRule, newGuardianItem)
	}

	logs, sub, err := _StakingVault.contract.FilterLogs(opts, "GuardianUpdated", newGuardianRule)
	if err != nil {
		return nil, err
	}
	return &StakingVaultGuardianUpdatedIterator{contract: _StakingVault.contract, event: "GuardianUpdated", logs: logs, sub: sub}, nil
}

// WatchGuardianUpdated is a free log subscription operation binding the contract event 0x6bb7ff33e730289800c62ad882105a144a74010d2bdbb9a942544a3005ad55bf.
//
// Solidity: event GuardianUpdated(address indexed newGuardian)
func (_StakingVault *StakingVaultFilterer) WatchGuardianUpdated(opts *bind.WatchOpts, sink chan<- *StakingVaultGuardianUpdated, newGuardian []common.Address) (event.Subscription, error) {

	var newGuardianRule []interface{}
	for _, newGuardianItem := range newGuardian {
		newGuardianRule = append(newGuardianRule, newGuardianItem)
	}

	logs, sub, err := _StakingVault.contract.WatchLogs(opts, "GuardianUpdated", newGuardianRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingVaultGuardianUpdated)
				if err := _StakingVault.contract.UnpackLog(event, "GuardianUpdated", log); err != nil {
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
func (_StakingVault *StakingVaultFilterer) ParseGuardianUpdated(log types.Log) (*StakingVaultGuardianUpdated, error) {
	event := new(StakingVaultGuardianUpdated)
	if err := _StakingVault.contract.UnpackLog(event, "GuardianUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingVaultInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the StakingVault contract.
type StakingVaultInitializedIterator struct {
	Event *StakingVaultInitialized // Event containing the contract specifics and raw log

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
func (it *StakingVaultInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingVaultInitialized)
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
		it.Event = new(StakingVaultInitialized)
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
func (it *StakingVaultInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingVaultInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingVaultInitialized represents a Initialized event raised by the StakingVault contract.
type StakingVaultInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_StakingVault *StakingVaultFilterer) FilterInitialized(opts *bind.FilterOpts) (*StakingVaultInitializedIterator, error) {

	logs, sub, err := _StakingVault.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &StakingVaultInitializedIterator{contract: _StakingVault.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_StakingVault *StakingVaultFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *StakingVaultInitialized) (event.Subscription, error) {

	logs, sub, err := _StakingVault.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingVaultInitialized)
				if err := _StakingVault.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_StakingVault *StakingVaultFilterer) ParseInitialized(log types.Log) (*StakingVaultInitialized, error) {
	event := new(StakingVaultInitialized)
	if err := _StakingVault.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingVaultReallocatedIterator is returned from FilterReallocated and is used to iterate over the raw logs and unpacked data for Reallocated events raised by the StakingVault contract.
type StakingVaultReallocatedIterator struct {
	Event *StakingVaultReallocated // Event containing the contract specifics and raw log

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
func (it *StakingVaultReallocatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingVaultReallocated)
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
		it.Event = new(StakingVaultReallocated)
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
func (it *StakingVaultReallocatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingVaultReallocatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingVaultReallocated represents a Reallocated event raised by the StakingVault contract.
type StakingVaultReallocated struct {
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
func (_StakingVault *StakingVaultFilterer) FilterReallocated(opts *bind.FilterOpts, staker []common.Address) (*StakingVaultReallocatedIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}

	logs, sub, err := _StakingVault.contract.FilterLogs(opts, "Reallocated", stakerRule)
	if err != nil {
		return nil, err
	}
	return &StakingVaultReallocatedIterator{contract: _StakingVault.contract, event: "Reallocated", logs: logs, sub: sub}, nil
}

// WatchReallocated is a free log subscription operation binding the contract event 0x726c93ba67bfe4c677e37114279f0ad9aab5ee9ffbd1158923be5d0fec3b1b45.
//
// Solidity: event Reallocated(address indexed staker, address fromAgent, uint256 fromWorknetId, address toAgent, uint256 toWorknetId, uint256 amount, address operator)
func (_StakingVault *StakingVaultFilterer) WatchReallocated(opts *bind.WatchOpts, sink chan<- *StakingVaultReallocated, staker []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}

	logs, sub, err := _StakingVault.contract.WatchLogs(opts, "Reallocated", stakerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingVaultReallocated)
				if err := _StakingVault.contract.UnpackLog(event, "Reallocated", log); err != nil {
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
func (_StakingVault *StakingVaultFilterer) ParseReallocated(log types.Log) (*StakingVaultReallocated, error) {
	event := new(StakingVaultReallocated)
	if err := _StakingVault.contract.UnpackLog(event, "Reallocated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingVaultStakeNFTSetIterator is returned from FilterStakeNFTSet and is used to iterate over the raw logs and unpacked data for StakeNFTSet events raised by the StakingVault contract.
type StakingVaultStakeNFTSetIterator struct {
	Event *StakingVaultStakeNFTSet // Event containing the contract specifics and raw log

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
func (it *StakingVaultStakeNFTSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingVaultStakeNFTSet)
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
		it.Event = new(StakingVaultStakeNFTSet)
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
func (it *StakingVaultStakeNFTSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingVaultStakeNFTSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingVaultStakeNFTSet represents a StakeNFTSet event raised by the StakingVault contract.
type StakingVaultStakeNFTSet struct {
	StakeNFT common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterStakeNFTSet is a free log retrieval operation binding the contract event 0xa7dc7ffb248216564bae4b23cefffef79fb6c1d1809fa0e0064ad670c04d0a68.
//
// Solidity: event StakeNFTSet(address indexed stakeNFT)
func (_StakingVault *StakingVaultFilterer) FilterStakeNFTSet(opts *bind.FilterOpts, stakeNFT []common.Address) (*StakingVaultStakeNFTSetIterator, error) {

	var stakeNFTRule []interface{}
	for _, stakeNFTItem := range stakeNFT {
		stakeNFTRule = append(stakeNFTRule, stakeNFTItem)
	}

	logs, sub, err := _StakingVault.contract.FilterLogs(opts, "StakeNFTSet", stakeNFTRule)
	if err != nil {
		return nil, err
	}
	return &StakingVaultStakeNFTSetIterator{contract: _StakingVault.contract, event: "StakeNFTSet", logs: logs, sub: sub}, nil
}

// WatchStakeNFTSet is a free log subscription operation binding the contract event 0xa7dc7ffb248216564bae4b23cefffef79fb6c1d1809fa0e0064ad670c04d0a68.
//
// Solidity: event StakeNFTSet(address indexed stakeNFT)
func (_StakingVault *StakingVaultFilterer) WatchStakeNFTSet(opts *bind.WatchOpts, sink chan<- *StakingVaultStakeNFTSet, stakeNFT []common.Address) (event.Subscription, error) {

	var stakeNFTRule []interface{}
	for _, stakeNFTItem := range stakeNFT {
		stakeNFTRule = append(stakeNFTRule, stakeNFTItem)
	}

	logs, sub, err := _StakingVault.contract.WatchLogs(opts, "StakeNFTSet", stakeNFTRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingVaultStakeNFTSet)
				if err := _StakingVault.contract.UnpackLog(event, "StakeNFTSet", log); err != nil {
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

// ParseStakeNFTSet is a log parse operation binding the contract event 0xa7dc7ffb248216564bae4b23cefffef79fb6c1d1809fa0e0064ad670c04d0a68.
//
// Solidity: event StakeNFTSet(address indexed stakeNFT)
func (_StakingVault *StakingVaultFilterer) ParseStakeNFTSet(log types.Log) (*StakingVaultStakeNFTSet, error) {
	event := new(StakingVaultStakeNFTSet)
	if err := _StakingVault.contract.UnpackLog(event, "StakeNFTSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingVaultUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the StakingVault contract.
type StakingVaultUpgradedIterator struct {
	Event *StakingVaultUpgraded // Event containing the contract specifics and raw log

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
func (it *StakingVaultUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingVaultUpgraded)
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
		it.Event = new(StakingVaultUpgraded)
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
func (it *StakingVaultUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingVaultUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingVaultUpgraded represents a Upgraded event raised by the StakingVault contract.
type StakingVaultUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_StakingVault *StakingVaultFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*StakingVaultUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _StakingVault.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &StakingVaultUpgradedIterator{contract: _StakingVault.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_StakingVault *StakingVaultFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *StakingVaultUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _StakingVault.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingVaultUpgraded)
				if err := _StakingVault.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_StakingVault *StakingVaultFilterer) ParseUpgraded(log types.Log) (*StakingVaultUpgraded, error) {
	event := new(StakingVaultUpgraded)
	if err := _StakingVault.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
