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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"awpRegistry_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allocate\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"awpRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deallocate\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"freezeAgentAllocations\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAgentStake\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAgentSubnets\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSubnetTotalStake\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reallocate\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromAgent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromSubnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"toAgent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"toSubnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setStakeNFT\",\"inputs\":[{\"name\":\"stakeNFT_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakeNFT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"subnetTotalStake\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"userTotalAllocated\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AgentAllocationsFrozen\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"totalFrozen\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadySet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientAllocation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientUnallocated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAWPRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
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

// GetAgentStake is a free data retrieval call binding the contract method 0xf1ad80c6.
//
// Solidity: function getAgentStake(address user, address agent, uint256 subnetId) view returns(uint256)
func (_StakingVault *StakingVaultCaller) GetAgentStake(opts *bind.CallOpts, user common.Address, agent common.Address, subnetId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _StakingVault.contract.Call(opts, &out, "getAgentStake", user, agent, subnetId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAgentStake is a free data retrieval call binding the contract method 0xf1ad80c6.
//
// Solidity: function getAgentStake(address user, address agent, uint256 subnetId) view returns(uint256)
func (_StakingVault *StakingVaultSession) GetAgentStake(user common.Address, agent common.Address, subnetId *big.Int) (*big.Int, error) {
	return _StakingVault.Contract.GetAgentStake(&_StakingVault.CallOpts, user, agent, subnetId)
}

// GetAgentStake is a free data retrieval call binding the contract method 0xf1ad80c6.
//
// Solidity: function getAgentStake(address user, address agent, uint256 subnetId) view returns(uint256)
func (_StakingVault *StakingVaultCallerSession) GetAgentStake(user common.Address, agent common.Address, subnetId *big.Int) (*big.Int, error) {
	return _StakingVault.Contract.GetAgentStake(&_StakingVault.CallOpts, user, agent, subnetId)
}

// GetAgentSubnets is a free data retrieval call binding the contract method 0x0358ccb5.
//
// Solidity: function getAgentSubnets(address user, address agent) view returns(uint256[])
func (_StakingVault *StakingVaultCaller) GetAgentSubnets(opts *bind.CallOpts, user common.Address, agent common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _StakingVault.contract.Call(opts, &out, "getAgentSubnets", user, agent)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetAgentSubnets is a free data retrieval call binding the contract method 0x0358ccb5.
//
// Solidity: function getAgentSubnets(address user, address agent) view returns(uint256[])
func (_StakingVault *StakingVaultSession) GetAgentSubnets(user common.Address, agent common.Address) ([]*big.Int, error) {
	return _StakingVault.Contract.GetAgentSubnets(&_StakingVault.CallOpts, user, agent)
}

// GetAgentSubnets is a free data retrieval call binding the contract method 0x0358ccb5.
//
// Solidity: function getAgentSubnets(address user, address agent) view returns(uint256[])
func (_StakingVault *StakingVaultCallerSession) GetAgentSubnets(user common.Address, agent common.Address) ([]*big.Int, error) {
	return _StakingVault.Contract.GetAgentSubnets(&_StakingVault.CallOpts, user, agent)
}

// GetSubnetTotalStake is a free data retrieval call binding the contract method 0x1ed2facb.
//
// Solidity: function getSubnetTotalStake(uint256 subnetId) view returns(uint256)
func (_StakingVault *StakingVaultCaller) GetSubnetTotalStake(opts *bind.CallOpts, subnetId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _StakingVault.contract.Call(opts, &out, "getSubnetTotalStake", subnetId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetSubnetTotalStake is a free data retrieval call binding the contract method 0x1ed2facb.
//
// Solidity: function getSubnetTotalStake(uint256 subnetId) view returns(uint256)
func (_StakingVault *StakingVaultSession) GetSubnetTotalStake(subnetId *big.Int) (*big.Int, error) {
	return _StakingVault.Contract.GetSubnetTotalStake(&_StakingVault.CallOpts, subnetId)
}

// GetSubnetTotalStake is a free data retrieval call binding the contract method 0x1ed2facb.
//
// Solidity: function getSubnetTotalStake(uint256 subnetId) view returns(uint256)
func (_StakingVault *StakingVaultCallerSession) GetSubnetTotalStake(subnetId *big.Int) (*big.Int, error) {
	return _StakingVault.Contract.GetSubnetTotalStake(&_StakingVault.CallOpts, subnetId)
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

// SubnetTotalStake is a free data retrieval call binding the contract method 0xf1e18405.
//
// Solidity: function subnetTotalStake(uint256 ) view returns(uint256)
func (_StakingVault *StakingVaultCaller) SubnetTotalStake(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _StakingVault.contract.Call(opts, &out, "subnetTotalStake", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SubnetTotalStake is a free data retrieval call binding the contract method 0xf1e18405.
//
// Solidity: function subnetTotalStake(uint256 ) view returns(uint256)
func (_StakingVault *StakingVaultSession) SubnetTotalStake(arg0 *big.Int) (*big.Int, error) {
	return _StakingVault.Contract.SubnetTotalStake(&_StakingVault.CallOpts, arg0)
}

// SubnetTotalStake is a free data retrieval call binding the contract method 0xf1e18405.
//
// Solidity: function subnetTotalStake(uint256 ) view returns(uint256)
func (_StakingVault *StakingVaultCallerSession) SubnetTotalStake(arg0 *big.Int) (*big.Int, error) {
	return _StakingVault.Contract.SubnetTotalStake(&_StakingVault.CallOpts, arg0)
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

// Allocate is a paid mutator transaction binding the contract method 0xd035a9a7.
//
// Solidity: function allocate(address user, address agent, uint256 subnetId, uint256 amount) returns()
func (_StakingVault *StakingVaultTransactor) Allocate(opts *bind.TransactOpts, user common.Address, agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.contract.Transact(opts, "allocate", user, agent, subnetId, amount)
}

// Allocate is a paid mutator transaction binding the contract method 0xd035a9a7.
//
// Solidity: function allocate(address user, address agent, uint256 subnetId, uint256 amount) returns()
func (_StakingVault *StakingVaultSession) Allocate(user common.Address, agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.Contract.Allocate(&_StakingVault.TransactOpts, user, agent, subnetId, amount)
}

// Allocate is a paid mutator transaction binding the contract method 0xd035a9a7.
//
// Solidity: function allocate(address user, address agent, uint256 subnetId, uint256 amount) returns()
func (_StakingVault *StakingVaultTransactorSession) Allocate(user common.Address, agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.Contract.Allocate(&_StakingVault.TransactOpts, user, agent, subnetId, amount)
}

// Deallocate is a paid mutator transaction binding the contract method 0x716fb83d.
//
// Solidity: function deallocate(address user, address agent, uint256 subnetId, uint256 amount) returns()
func (_StakingVault *StakingVaultTransactor) Deallocate(opts *bind.TransactOpts, user common.Address, agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.contract.Transact(opts, "deallocate", user, agent, subnetId, amount)
}

// Deallocate is a paid mutator transaction binding the contract method 0x716fb83d.
//
// Solidity: function deallocate(address user, address agent, uint256 subnetId, uint256 amount) returns()
func (_StakingVault *StakingVaultSession) Deallocate(user common.Address, agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.Contract.Deallocate(&_StakingVault.TransactOpts, user, agent, subnetId, amount)
}

// Deallocate is a paid mutator transaction binding the contract method 0x716fb83d.
//
// Solidity: function deallocate(address user, address agent, uint256 subnetId, uint256 amount) returns()
func (_StakingVault *StakingVaultTransactorSession) Deallocate(user common.Address, agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.Contract.Deallocate(&_StakingVault.TransactOpts, user, agent, subnetId, amount)
}

// FreezeAgentAllocations is a paid mutator transaction binding the contract method 0x6f980813.
//
// Solidity: function freezeAgentAllocations(address user, address agent) returns()
func (_StakingVault *StakingVaultTransactor) FreezeAgentAllocations(opts *bind.TransactOpts, user common.Address, agent common.Address) (*types.Transaction, error) {
	return _StakingVault.contract.Transact(opts, "freezeAgentAllocations", user, agent)
}

// FreezeAgentAllocations is a paid mutator transaction binding the contract method 0x6f980813.
//
// Solidity: function freezeAgentAllocations(address user, address agent) returns()
func (_StakingVault *StakingVaultSession) FreezeAgentAllocations(user common.Address, agent common.Address) (*types.Transaction, error) {
	return _StakingVault.Contract.FreezeAgentAllocations(&_StakingVault.TransactOpts, user, agent)
}

// FreezeAgentAllocations is a paid mutator transaction binding the contract method 0x6f980813.
//
// Solidity: function freezeAgentAllocations(address user, address agent) returns()
func (_StakingVault *StakingVaultTransactorSession) FreezeAgentAllocations(user common.Address, agent common.Address) (*types.Transaction, error) {
	return _StakingVault.Contract.FreezeAgentAllocations(&_StakingVault.TransactOpts, user, agent)
}

// Reallocate is a paid mutator transaction binding the contract method 0xd5d5278d.
//
// Solidity: function reallocate(address user, address fromAgent, uint256 fromSubnetId, address toAgent, uint256 toSubnetId, uint256 amount) returns()
func (_StakingVault *StakingVaultTransactor) Reallocate(opts *bind.TransactOpts, user common.Address, fromAgent common.Address, fromSubnetId *big.Int, toAgent common.Address, toSubnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.contract.Transact(opts, "reallocate", user, fromAgent, fromSubnetId, toAgent, toSubnetId, amount)
}

// Reallocate is a paid mutator transaction binding the contract method 0xd5d5278d.
//
// Solidity: function reallocate(address user, address fromAgent, uint256 fromSubnetId, address toAgent, uint256 toSubnetId, uint256 amount) returns()
func (_StakingVault *StakingVaultSession) Reallocate(user common.Address, fromAgent common.Address, fromSubnetId *big.Int, toAgent common.Address, toSubnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.Contract.Reallocate(&_StakingVault.TransactOpts, user, fromAgent, fromSubnetId, toAgent, toSubnetId, amount)
}

// Reallocate is a paid mutator transaction binding the contract method 0xd5d5278d.
//
// Solidity: function reallocate(address user, address fromAgent, uint256 fromSubnetId, address toAgent, uint256 toSubnetId, uint256 amount) returns()
func (_StakingVault *StakingVaultTransactorSession) Reallocate(user common.Address, fromAgent common.Address, fromSubnetId *big.Int, toAgent common.Address, toSubnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _StakingVault.Contract.Reallocate(&_StakingVault.TransactOpts, user, fromAgent, fromSubnetId, toAgent, toSubnetId, amount)
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

// StakingVaultAgentAllocationsFrozenIterator is returned from FilterAgentAllocationsFrozen and is used to iterate over the raw logs and unpacked data for AgentAllocationsFrozen events raised by the StakingVault contract.
type StakingVaultAgentAllocationsFrozenIterator struct {
	Event *StakingVaultAgentAllocationsFrozen // Event containing the contract specifics and raw log

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
func (it *StakingVaultAgentAllocationsFrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingVaultAgentAllocationsFrozen)
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
		it.Event = new(StakingVaultAgentAllocationsFrozen)
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
func (it *StakingVaultAgentAllocationsFrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingVaultAgentAllocationsFrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingVaultAgentAllocationsFrozen represents a AgentAllocationsFrozen event raised by the StakingVault contract.
type StakingVaultAgentAllocationsFrozen struct {
	User        common.Address
	Agent       common.Address
	TotalFrozen *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAgentAllocationsFrozen is a free log retrieval operation binding the contract event 0x823cbf835fe0c4834f90d65da805f8a714ff7df3c298a0de29e4547944378f4f.
//
// Solidity: event AgentAllocationsFrozen(address indexed user, address indexed agent, uint256 totalFrozen)
func (_StakingVault *StakingVaultFilterer) FilterAgentAllocationsFrozen(opts *bind.FilterOpts, user []common.Address, agent []common.Address) (*StakingVaultAgentAllocationsFrozenIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _StakingVault.contract.FilterLogs(opts, "AgentAllocationsFrozen", userRule, agentRule)
	if err != nil {
		return nil, err
	}
	return &StakingVaultAgentAllocationsFrozenIterator{contract: _StakingVault.contract, event: "AgentAllocationsFrozen", logs: logs, sub: sub}, nil
}

// WatchAgentAllocationsFrozen is a free log subscription operation binding the contract event 0x823cbf835fe0c4834f90d65da805f8a714ff7df3c298a0de29e4547944378f4f.
//
// Solidity: event AgentAllocationsFrozen(address indexed user, address indexed agent, uint256 totalFrozen)
func (_StakingVault *StakingVaultFilterer) WatchAgentAllocationsFrozen(opts *bind.WatchOpts, sink chan<- *StakingVaultAgentAllocationsFrozen, user []common.Address, agent []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _StakingVault.contract.WatchLogs(opts, "AgentAllocationsFrozen", userRule, agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingVaultAgentAllocationsFrozen)
				if err := _StakingVault.contract.UnpackLog(event, "AgentAllocationsFrozen", log); err != nil {
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

// ParseAgentAllocationsFrozen is a log parse operation binding the contract event 0x823cbf835fe0c4834f90d65da805f8a714ff7df3c298a0de29e4547944378f4f.
//
// Solidity: event AgentAllocationsFrozen(address indexed user, address indexed agent, uint256 totalFrozen)
func (_StakingVault *StakingVaultFilterer) ParseAgentAllocationsFrozen(log types.Log) (*StakingVaultAgentAllocationsFrozen, error) {
	event := new(StakingVaultAgentAllocationsFrozen)
	if err := _StakingVault.contract.UnpackLog(event, "AgentAllocationsFrozen", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
