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

// SubnetNFTSubnetData is an auto generated low-level Go binding around an user-defined struct.
type SubnetNFTSubnetData struct {
	Name          string
	SubnetManager common.Address
	AlphaToken    common.Address
	SkillsURI     string
	MinStake      *big.Int
	Owner         common.Address
}

// SubnetNFTMetaData contains all meta data concerning the SubnetNFT contract.
var SubnetNFTMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"awpRegistry_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"awpRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAlphaToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinStake\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSubnetData\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structSubnetNFT.SubnetData\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"subnetManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"alphaToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSubnetManager\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"subnetManager_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"alphaToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"minStake_\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"skillsURI_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBaseURI\",\"inputs\":[{\"name\":\"uri\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinStake\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minStake_\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSkillsURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"skillsURI_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinStakeUpdated\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SkillsURIUpdated\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotAWPRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotTokenOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TokenNotExist\",\"inputs\":[]}]",
}

// SubnetNFTABI is the input ABI used to generate the binding from.
// Deprecated: Use SubnetNFTMetaData.ABI instead.
var SubnetNFTABI = SubnetNFTMetaData.ABI

// SubnetNFT is an auto generated Go binding around an Ethereum contract.
type SubnetNFT struct {
	SubnetNFTCaller     // Read-only binding to the contract
	SubnetNFTTransactor // Write-only binding to the contract
	SubnetNFTFilterer   // Log filterer for contract events
}

// SubnetNFTCaller is an auto generated read-only Go binding around an Ethereum contract.
type SubnetNFTCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubnetNFTTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SubnetNFTTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubnetNFTFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SubnetNFTFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubnetNFTSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SubnetNFTSession struct {
	Contract     *SubnetNFT        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SubnetNFTCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SubnetNFTCallerSession struct {
	Contract *SubnetNFTCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// SubnetNFTTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SubnetNFTTransactorSession struct {
	Contract     *SubnetNFTTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// SubnetNFTRaw is an auto generated low-level Go binding around an Ethereum contract.
type SubnetNFTRaw struct {
	Contract *SubnetNFT // Generic contract binding to access the raw methods on
}

// SubnetNFTCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SubnetNFTCallerRaw struct {
	Contract *SubnetNFTCaller // Generic read-only contract binding to access the raw methods on
}

// SubnetNFTTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SubnetNFTTransactorRaw struct {
	Contract *SubnetNFTTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSubnetNFT creates a new instance of SubnetNFT, bound to a specific deployed contract.
func NewSubnetNFT(address common.Address, backend bind.ContractBackend) (*SubnetNFT, error) {
	contract, err := bindSubnetNFT(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SubnetNFT{SubnetNFTCaller: SubnetNFTCaller{contract: contract}, SubnetNFTTransactor: SubnetNFTTransactor{contract: contract}, SubnetNFTFilterer: SubnetNFTFilterer{contract: contract}}, nil
}

// NewSubnetNFTCaller creates a new read-only instance of SubnetNFT, bound to a specific deployed contract.
func NewSubnetNFTCaller(address common.Address, caller bind.ContractCaller) (*SubnetNFTCaller, error) {
	contract, err := bindSubnetNFT(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SubnetNFTCaller{contract: contract}, nil
}

// NewSubnetNFTTransactor creates a new write-only instance of SubnetNFT, bound to a specific deployed contract.
func NewSubnetNFTTransactor(address common.Address, transactor bind.ContractTransactor) (*SubnetNFTTransactor, error) {
	contract, err := bindSubnetNFT(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SubnetNFTTransactor{contract: contract}, nil
}

// NewSubnetNFTFilterer creates a new log filterer instance of SubnetNFT, bound to a specific deployed contract.
func NewSubnetNFTFilterer(address common.Address, filterer bind.ContractFilterer) (*SubnetNFTFilterer, error) {
	contract, err := bindSubnetNFT(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SubnetNFTFilterer{contract: contract}, nil
}

// bindSubnetNFT binds a generic wrapper to an already deployed contract.
func bindSubnetNFT(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SubnetNFTMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SubnetNFT *SubnetNFTRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SubnetNFT.Contract.SubnetNFTCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SubnetNFT *SubnetNFTRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SubnetNFTTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SubnetNFT *SubnetNFTRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SubnetNFTTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SubnetNFT *SubnetNFTCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SubnetNFT.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SubnetNFT *SubnetNFTTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubnetNFT.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SubnetNFT *SubnetNFTTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubnetNFT.Contract.contract.Transact(opts, method, params...)
}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_SubnetNFT *SubnetNFTCaller) AwpRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "awpRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_SubnetNFT *SubnetNFTSession) AwpRegistry() (common.Address, error) {
	return _SubnetNFT.Contract.AwpRegistry(&_SubnetNFT.CallOpts)
}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_SubnetNFT *SubnetNFTCallerSession) AwpRegistry() (common.Address, error) {
	return _SubnetNFT.Contract.AwpRegistry(&_SubnetNFT.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_SubnetNFT *SubnetNFTCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_SubnetNFT *SubnetNFTSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _SubnetNFT.Contract.BalanceOf(&_SubnetNFT.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_SubnetNFT *SubnetNFTCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _SubnetNFT.Contract.BalanceOf(&_SubnetNFT.CallOpts, owner)
}

// GetAlphaToken is a free data retrieval call binding the contract method 0xc7bc8ec6.
//
// Solidity: function getAlphaToken(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTCaller) GetAlphaToken(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "getAlphaToken", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAlphaToken is a free data retrieval call binding the contract method 0xc7bc8ec6.
//
// Solidity: function getAlphaToken(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTSession) GetAlphaToken(tokenId *big.Int) (common.Address, error) {
	return _SubnetNFT.Contract.GetAlphaToken(&_SubnetNFT.CallOpts, tokenId)
}

// GetAlphaToken is a free data retrieval call binding the contract method 0xc7bc8ec6.
//
// Solidity: function getAlphaToken(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTCallerSession) GetAlphaToken(tokenId *big.Int) (common.Address, error) {
	return _SubnetNFT.Contract.GetAlphaToken(&_SubnetNFT.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _SubnetNFT.Contract.GetApproved(&_SubnetNFT.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _SubnetNFT.Contract.GetApproved(&_SubnetNFT.CallOpts, tokenId)
}

// GetMinStake is a free data retrieval call binding the contract method 0x73f231e7.
//
// Solidity: function getMinStake(uint256 tokenId) view returns(uint128)
func (_SubnetNFT *SubnetNFTCaller) GetMinStake(opts *bind.CallOpts, tokenId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "getMinStake", tokenId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinStake is a free data retrieval call binding the contract method 0x73f231e7.
//
// Solidity: function getMinStake(uint256 tokenId) view returns(uint128)
func (_SubnetNFT *SubnetNFTSession) GetMinStake(tokenId *big.Int) (*big.Int, error) {
	return _SubnetNFT.Contract.GetMinStake(&_SubnetNFT.CallOpts, tokenId)
}

// GetMinStake is a free data retrieval call binding the contract method 0x73f231e7.
//
// Solidity: function getMinStake(uint256 tokenId) view returns(uint128)
func (_SubnetNFT *SubnetNFTCallerSession) GetMinStake(tokenId *big.Int) (*big.Int, error) {
	return _SubnetNFT.Contract.GetMinStake(&_SubnetNFT.CallOpts, tokenId)
}

// GetSubnetData is a free data retrieval call binding the contract method 0x854744ca.
//
// Solidity: function getSubnetData(uint256 tokenId) view returns((string,address,address,string,uint128,address))
func (_SubnetNFT *SubnetNFTCaller) GetSubnetData(opts *bind.CallOpts, tokenId *big.Int) (SubnetNFTSubnetData, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "getSubnetData", tokenId)

	if err != nil {
		return *new(SubnetNFTSubnetData), err
	}

	out0 := *abi.ConvertType(out[0], new(SubnetNFTSubnetData)).(*SubnetNFTSubnetData)

	return out0, err

}

// GetSubnetData is a free data retrieval call binding the contract method 0x854744ca.
//
// Solidity: function getSubnetData(uint256 tokenId) view returns((string,address,address,string,uint128,address))
func (_SubnetNFT *SubnetNFTSession) GetSubnetData(tokenId *big.Int) (SubnetNFTSubnetData, error) {
	return _SubnetNFT.Contract.GetSubnetData(&_SubnetNFT.CallOpts, tokenId)
}

// GetSubnetData is a free data retrieval call binding the contract method 0x854744ca.
//
// Solidity: function getSubnetData(uint256 tokenId) view returns((string,address,address,string,uint128,address))
func (_SubnetNFT *SubnetNFTCallerSession) GetSubnetData(tokenId *big.Int) (SubnetNFTSubnetData, error) {
	return _SubnetNFT.Contract.GetSubnetData(&_SubnetNFT.CallOpts, tokenId)
}

// GetSubnetManager is a free data retrieval call binding the contract method 0xe630cb96.
//
// Solidity: function getSubnetManager(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTCaller) GetSubnetManager(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "getSubnetManager", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetSubnetManager is a free data retrieval call binding the contract method 0xe630cb96.
//
// Solidity: function getSubnetManager(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTSession) GetSubnetManager(tokenId *big.Int) (common.Address, error) {
	return _SubnetNFT.Contract.GetSubnetManager(&_SubnetNFT.CallOpts, tokenId)
}

// GetSubnetManager is a free data retrieval call binding the contract method 0xe630cb96.
//
// Solidity: function getSubnetManager(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTCallerSession) GetSubnetManager(tokenId *big.Int) (common.Address, error) {
	return _SubnetNFT.Contract.GetSubnetManager(&_SubnetNFT.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_SubnetNFT *SubnetNFTCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_SubnetNFT *SubnetNFTSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _SubnetNFT.Contract.IsApprovedForAll(&_SubnetNFT.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_SubnetNFT *SubnetNFTCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _SubnetNFT.Contract.IsApprovedForAll(&_SubnetNFT.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_SubnetNFT *SubnetNFTCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_SubnetNFT *SubnetNFTSession) Name() (string, error) {
	return _SubnetNFT.Contract.Name(&_SubnetNFT.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_SubnetNFT *SubnetNFTCallerSession) Name() (string, error) {
	return _SubnetNFT.Contract.Name(&_SubnetNFT.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _SubnetNFT.Contract.OwnerOf(&_SubnetNFT.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_SubnetNFT *SubnetNFTCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _SubnetNFT.Contract.OwnerOf(&_SubnetNFT.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_SubnetNFT *SubnetNFTCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_SubnetNFT *SubnetNFTSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _SubnetNFT.Contract.SupportsInterface(&_SubnetNFT.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_SubnetNFT *SubnetNFTCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _SubnetNFT.Contract.SupportsInterface(&_SubnetNFT.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_SubnetNFT *SubnetNFTCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_SubnetNFT *SubnetNFTSession) Symbol() (string, error) {
	return _SubnetNFT.Contract.Symbol(&_SubnetNFT.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_SubnetNFT *SubnetNFTCallerSession) Symbol() (string, error) {
	return _SubnetNFT.Contract.Symbol(&_SubnetNFT.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_SubnetNFT *SubnetNFTCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _SubnetNFT.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_SubnetNFT *SubnetNFTSession) TokenURI(tokenId *big.Int) (string, error) {
	return _SubnetNFT.Contract.TokenURI(&_SubnetNFT.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_SubnetNFT *SubnetNFTCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _SubnetNFT.Contract.TokenURI(&_SubnetNFT.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.Approve(&_SubnetNFT.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.Approve(&_SubnetNFT.TransactOpts, to, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "burn", tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.Burn(&_SubnetNFT.TransactOpts, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.Burn(&_SubnetNFT.TransactOpts, tokenId)
}

// Mint is a paid mutator transaction binding the contract method 0xa9584608.
//
// Solidity: function mint(address to, uint256 tokenId, string name_, address subnetManager_, address alphaToken_, uint128 minStake_, string skillsURI_) returns()
func (_SubnetNFT *SubnetNFTTransactor) Mint(opts *bind.TransactOpts, to common.Address, tokenId *big.Int, name_ string, subnetManager_ common.Address, alphaToken_ common.Address, minStake_ *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "mint", to, tokenId, name_, subnetManager_, alphaToken_, minStake_, skillsURI_)
}

// Mint is a paid mutator transaction binding the contract method 0xa9584608.
//
// Solidity: function mint(address to, uint256 tokenId, string name_, address subnetManager_, address alphaToken_, uint128 minStake_, string skillsURI_) returns()
func (_SubnetNFT *SubnetNFTSession) Mint(to common.Address, tokenId *big.Int, name_ string, subnetManager_ common.Address, alphaToken_ common.Address, minStake_ *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _SubnetNFT.Contract.Mint(&_SubnetNFT.TransactOpts, to, tokenId, name_, subnetManager_, alphaToken_, minStake_, skillsURI_)
}

// Mint is a paid mutator transaction binding the contract method 0xa9584608.
//
// Solidity: function mint(address to, uint256 tokenId, string name_, address subnetManager_, address alphaToken_, uint128 minStake_, string skillsURI_) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) Mint(to common.Address, tokenId *big.Int, name_ string, subnetManager_ common.Address, alphaToken_ common.Address, minStake_ *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _SubnetNFT.Contract.Mint(&_SubnetNFT.TransactOpts, to, tokenId, name_, subnetManager_, alphaToken_, minStake_, skillsURI_)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SafeTransferFrom(&_SubnetNFT.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SafeTransferFrom(&_SubnetNFT.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_SubnetNFT *SubnetNFTTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_SubnetNFT *SubnetNFTSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SafeTransferFrom0(&_SubnetNFT.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SafeTransferFrom0(&_SubnetNFT.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_SubnetNFT *SubnetNFTTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_SubnetNFT *SubnetNFTSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SetApprovalForAll(&_SubnetNFT.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SetApprovalForAll(&_SubnetNFT.TransactOpts, operator, approved)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string uri) returns()
func (_SubnetNFT *SubnetNFTTransactor) SetBaseURI(opts *bind.TransactOpts, uri string) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "setBaseURI", uri)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string uri) returns()
func (_SubnetNFT *SubnetNFTSession) SetBaseURI(uri string) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SetBaseURI(&_SubnetNFT.TransactOpts, uri)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string uri) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) SetBaseURI(uri string) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SetBaseURI(&_SubnetNFT.TransactOpts, uri)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x63a9bbe5.
//
// Solidity: function setMinStake(uint256 tokenId, uint128 minStake_) returns()
func (_SubnetNFT *SubnetNFTTransactor) SetMinStake(opts *bind.TransactOpts, tokenId *big.Int, minStake_ *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "setMinStake", tokenId, minStake_)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x63a9bbe5.
//
// Solidity: function setMinStake(uint256 tokenId, uint128 minStake_) returns()
func (_SubnetNFT *SubnetNFTSession) SetMinStake(tokenId *big.Int, minStake_ *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SetMinStake(&_SubnetNFT.TransactOpts, tokenId, minStake_)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x63a9bbe5.
//
// Solidity: function setMinStake(uint256 tokenId, uint128 minStake_) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) SetMinStake(tokenId *big.Int, minStake_ *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SetMinStake(&_SubnetNFT.TransactOpts, tokenId, minStake_)
}

// SetSkillsURI is a paid mutator transaction binding the contract method 0x7c2f4cd6.
//
// Solidity: function setSkillsURI(uint256 tokenId, string skillsURI_) returns()
func (_SubnetNFT *SubnetNFTTransactor) SetSkillsURI(opts *bind.TransactOpts, tokenId *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "setSkillsURI", tokenId, skillsURI_)
}

// SetSkillsURI is a paid mutator transaction binding the contract method 0x7c2f4cd6.
//
// Solidity: function setSkillsURI(uint256 tokenId, string skillsURI_) returns()
func (_SubnetNFT *SubnetNFTSession) SetSkillsURI(tokenId *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SetSkillsURI(&_SubnetNFT.TransactOpts, tokenId, skillsURI_)
}

// SetSkillsURI is a paid mutator transaction binding the contract method 0x7c2f4cd6.
//
// Solidity: function setSkillsURI(uint256 tokenId, string skillsURI_) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) SetSkillsURI(tokenId *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _SubnetNFT.Contract.SetSkillsURI(&_SubnetNFT.TransactOpts, tokenId, skillsURI_)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.TransferFrom(&_SubnetNFT.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_SubnetNFT *SubnetNFTTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _SubnetNFT.Contract.TransferFrom(&_SubnetNFT.TransactOpts, from, to, tokenId)
}

// SubnetNFTApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the SubnetNFT contract.
type SubnetNFTApprovalIterator struct {
	Event *SubnetNFTApproval // Event containing the contract specifics and raw log

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
func (it *SubnetNFTApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubnetNFTApproval)
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
		it.Event = new(SubnetNFTApproval)
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
func (it *SubnetNFTApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SubnetNFTApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SubnetNFTApproval represents a Approval event raised by the SubnetNFT contract.
type SubnetNFTApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_SubnetNFT *SubnetNFTFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*SubnetNFTApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _SubnetNFT.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &SubnetNFTApprovalIterator{contract: _SubnetNFT.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_SubnetNFT *SubnetNFTFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *SubnetNFTApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _SubnetNFT.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SubnetNFTApproval)
				if err := _SubnetNFT.contract.UnpackLog(event, "Approval", log); err != nil {
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
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_SubnetNFT *SubnetNFTFilterer) ParseApproval(log types.Log) (*SubnetNFTApproval, error) {
	event := new(SubnetNFTApproval)
	if err := _SubnetNFT.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SubnetNFTApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the SubnetNFT contract.
type SubnetNFTApprovalForAllIterator struct {
	Event *SubnetNFTApprovalForAll // Event containing the contract specifics and raw log

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
func (it *SubnetNFTApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubnetNFTApprovalForAll)
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
		it.Event = new(SubnetNFTApprovalForAll)
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
func (it *SubnetNFTApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SubnetNFTApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SubnetNFTApprovalForAll represents a ApprovalForAll event raised by the SubnetNFT contract.
type SubnetNFTApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_SubnetNFT *SubnetNFTFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*SubnetNFTApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _SubnetNFT.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &SubnetNFTApprovalForAllIterator{contract: _SubnetNFT.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_SubnetNFT *SubnetNFTFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *SubnetNFTApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _SubnetNFT.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SubnetNFTApprovalForAll)
				if err := _SubnetNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_SubnetNFT *SubnetNFTFilterer) ParseApprovalForAll(log types.Log) (*SubnetNFTApprovalForAll, error) {
	event := new(SubnetNFTApprovalForAll)
	if err := _SubnetNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SubnetNFTMinStakeUpdatedIterator is returned from FilterMinStakeUpdated and is used to iterate over the raw logs and unpacked data for MinStakeUpdated events raised by the SubnetNFT contract.
type SubnetNFTMinStakeUpdatedIterator struct {
	Event *SubnetNFTMinStakeUpdated // Event containing the contract specifics and raw log

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
func (it *SubnetNFTMinStakeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubnetNFTMinStakeUpdated)
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
		it.Event = new(SubnetNFTMinStakeUpdated)
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
func (it *SubnetNFTMinStakeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SubnetNFTMinStakeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SubnetNFTMinStakeUpdated represents a MinStakeUpdated event raised by the SubnetNFT contract.
type SubnetNFTMinStakeUpdated struct {
	TokenId  *big.Int
	MinStake *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMinStakeUpdated is a free log retrieval operation binding the contract event 0xd0b53d029b5624436a948f7e5e2d9854defd8058cb4a20ff51a0ff9599ad6de8.
//
// Solidity: event MinStakeUpdated(uint256 indexed tokenId, uint128 minStake)
func (_SubnetNFT *SubnetNFTFilterer) FilterMinStakeUpdated(opts *bind.FilterOpts, tokenId []*big.Int) (*SubnetNFTMinStakeUpdatedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _SubnetNFT.contract.FilterLogs(opts, "MinStakeUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &SubnetNFTMinStakeUpdatedIterator{contract: _SubnetNFT.contract, event: "MinStakeUpdated", logs: logs, sub: sub}, nil
}

// WatchMinStakeUpdated is a free log subscription operation binding the contract event 0xd0b53d029b5624436a948f7e5e2d9854defd8058cb4a20ff51a0ff9599ad6de8.
//
// Solidity: event MinStakeUpdated(uint256 indexed tokenId, uint128 minStake)
func (_SubnetNFT *SubnetNFTFilterer) WatchMinStakeUpdated(opts *bind.WatchOpts, sink chan<- *SubnetNFTMinStakeUpdated, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _SubnetNFT.contract.WatchLogs(opts, "MinStakeUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SubnetNFTMinStakeUpdated)
				if err := _SubnetNFT.contract.UnpackLog(event, "MinStakeUpdated", log); err != nil {
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

// ParseMinStakeUpdated is a log parse operation binding the contract event 0xd0b53d029b5624436a948f7e5e2d9854defd8058cb4a20ff51a0ff9599ad6de8.
//
// Solidity: event MinStakeUpdated(uint256 indexed tokenId, uint128 minStake)
func (_SubnetNFT *SubnetNFTFilterer) ParseMinStakeUpdated(log types.Log) (*SubnetNFTMinStakeUpdated, error) {
	event := new(SubnetNFTMinStakeUpdated)
	if err := _SubnetNFT.contract.UnpackLog(event, "MinStakeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SubnetNFTSkillsURIUpdatedIterator is returned from FilterSkillsURIUpdated and is used to iterate over the raw logs and unpacked data for SkillsURIUpdated events raised by the SubnetNFT contract.
type SubnetNFTSkillsURIUpdatedIterator struct {
	Event *SubnetNFTSkillsURIUpdated // Event containing the contract specifics and raw log

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
func (it *SubnetNFTSkillsURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubnetNFTSkillsURIUpdated)
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
		it.Event = new(SubnetNFTSkillsURIUpdated)
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
func (it *SubnetNFTSkillsURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SubnetNFTSkillsURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SubnetNFTSkillsURIUpdated represents a SkillsURIUpdated event raised by the SubnetNFT contract.
type SubnetNFTSkillsURIUpdated struct {
	TokenId   *big.Int
	SkillsURI string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSkillsURIUpdated is a free log retrieval operation binding the contract event 0xd1332ed84c54e159e7a4245f8e021aff5b3389b685598c228394168ae30d1020.
//
// Solidity: event SkillsURIUpdated(uint256 indexed tokenId, string skillsURI)
func (_SubnetNFT *SubnetNFTFilterer) FilterSkillsURIUpdated(opts *bind.FilterOpts, tokenId []*big.Int) (*SubnetNFTSkillsURIUpdatedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _SubnetNFT.contract.FilterLogs(opts, "SkillsURIUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &SubnetNFTSkillsURIUpdatedIterator{contract: _SubnetNFT.contract, event: "SkillsURIUpdated", logs: logs, sub: sub}, nil
}

// WatchSkillsURIUpdated is a free log subscription operation binding the contract event 0xd1332ed84c54e159e7a4245f8e021aff5b3389b685598c228394168ae30d1020.
//
// Solidity: event SkillsURIUpdated(uint256 indexed tokenId, string skillsURI)
func (_SubnetNFT *SubnetNFTFilterer) WatchSkillsURIUpdated(opts *bind.WatchOpts, sink chan<- *SubnetNFTSkillsURIUpdated, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _SubnetNFT.contract.WatchLogs(opts, "SkillsURIUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SubnetNFTSkillsURIUpdated)
				if err := _SubnetNFT.contract.UnpackLog(event, "SkillsURIUpdated", log); err != nil {
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

// ParseSkillsURIUpdated is a log parse operation binding the contract event 0xd1332ed84c54e159e7a4245f8e021aff5b3389b685598c228394168ae30d1020.
//
// Solidity: event SkillsURIUpdated(uint256 indexed tokenId, string skillsURI)
func (_SubnetNFT *SubnetNFTFilterer) ParseSkillsURIUpdated(log types.Log) (*SubnetNFTSkillsURIUpdated, error) {
	event := new(SubnetNFTSkillsURIUpdated)
	if err := _SubnetNFT.contract.UnpackLog(event, "SkillsURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SubnetNFTTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the SubnetNFT contract.
type SubnetNFTTransferIterator struct {
	Event *SubnetNFTTransfer // Event containing the contract specifics and raw log

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
func (it *SubnetNFTTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubnetNFTTransfer)
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
		it.Event = new(SubnetNFTTransfer)
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
func (it *SubnetNFTTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SubnetNFTTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SubnetNFTTransfer represents a Transfer event raised by the SubnetNFT contract.
type SubnetNFTTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_SubnetNFT *SubnetNFTFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*SubnetNFTTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _SubnetNFT.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &SubnetNFTTransferIterator{contract: _SubnetNFT.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_SubnetNFT *SubnetNFTFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *SubnetNFTTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _SubnetNFT.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SubnetNFTTransfer)
				if err := _SubnetNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
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
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_SubnetNFT *SubnetNFTFilterer) ParseTransfer(log types.Log) (*SubnetNFTTransfer, error) {
	event := new(SubnetNFTTransfer)
	if err := _SubnetNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
