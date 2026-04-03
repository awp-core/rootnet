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

// WorknetNFTWorknetData is an auto generated low-level Go binding around an user-defined struct.
type WorknetNFTWorknetData struct {
	Name           string
	WorknetManager common.Address
	AlphaToken     common.Address
	SkillsURI      string
	MinStake       *big.Int
	Owner          common.Address
}

// WorknetNFTMetaData contains all meta data concerning the WorknetNFT contract.
var WorknetNFTMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"awpRegistry_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"awpRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAlphaToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinStake\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWorknetData\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structWorknetNFT.WorknetData\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"worknetManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"alphaToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWorknetManager\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"worknetManager_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"alphaToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"minStake_\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"skillsURI_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBaseURI\",\"inputs\":[{\"name\":\"uri\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMetadataURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"metadataURI_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinStake\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minStake_\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSkillsURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"skillsURI_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MetadataURIUpdated\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"metadataURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinStakeUpdated\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SkillsURIUpdated\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"JsonUnsafeCharacter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAWPRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotTokenOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StringsInsufficientHexLength\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TokenNotExist\",\"inputs\":[]}]",
}

// WorknetNFTABI is the input ABI used to generate the binding from.
// Deprecated: Use WorknetNFTMetaData.ABI instead.
var WorknetNFTABI = WorknetNFTMetaData.ABI

// WorknetNFT is an auto generated Go binding around an Ethereum contract.
type WorknetNFT struct {
	WorknetNFTCaller     // Read-only binding to the contract
	WorknetNFTTransactor // Write-only binding to the contract
	WorknetNFTFilterer   // Log filterer for contract events
}

// WorknetNFTCaller is an auto generated read-only Go binding around an Ethereum contract.
type WorknetNFTCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WorknetNFTTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WorknetNFTTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WorknetNFTFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WorknetNFTFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WorknetNFTSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WorknetNFTSession struct {
	Contract     *WorknetNFT       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WorknetNFTCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WorknetNFTCallerSession struct {
	Contract *WorknetNFTCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// WorknetNFTTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WorknetNFTTransactorSession struct {
	Contract     *WorknetNFTTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// WorknetNFTRaw is an auto generated low-level Go binding around an Ethereum contract.
type WorknetNFTRaw struct {
	Contract *WorknetNFT // Generic contract binding to access the raw methods on
}

// WorknetNFTCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WorknetNFTCallerRaw struct {
	Contract *WorknetNFTCaller // Generic read-only contract binding to access the raw methods on
}

// WorknetNFTTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WorknetNFTTransactorRaw struct {
	Contract *WorknetNFTTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWorknetNFT creates a new instance of WorknetNFT, bound to a specific deployed contract.
func NewWorknetNFT(address common.Address, backend bind.ContractBackend) (*WorknetNFT, error) {
	contract, err := bindWorknetNFT(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WorknetNFT{WorknetNFTCaller: WorknetNFTCaller{contract: contract}, WorknetNFTTransactor: WorknetNFTTransactor{contract: contract}, WorknetNFTFilterer: WorknetNFTFilterer{contract: contract}}, nil
}

// NewWorknetNFTCaller creates a new read-only instance of WorknetNFT, bound to a specific deployed contract.
func NewWorknetNFTCaller(address common.Address, caller bind.ContractCaller) (*WorknetNFTCaller, error) {
	contract, err := bindWorknetNFT(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WorknetNFTCaller{contract: contract}, nil
}

// NewWorknetNFTTransactor creates a new write-only instance of WorknetNFT, bound to a specific deployed contract.
func NewWorknetNFTTransactor(address common.Address, transactor bind.ContractTransactor) (*WorknetNFTTransactor, error) {
	contract, err := bindWorknetNFT(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WorknetNFTTransactor{contract: contract}, nil
}

// NewWorknetNFTFilterer creates a new log filterer instance of WorknetNFT, bound to a specific deployed contract.
func NewWorknetNFTFilterer(address common.Address, filterer bind.ContractFilterer) (*WorknetNFTFilterer, error) {
	contract, err := bindWorknetNFT(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WorknetNFTFilterer{contract: contract}, nil
}

// bindWorknetNFT binds a generic wrapper to an already deployed contract.
func bindWorknetNFT(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WorknetNFTMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WorknetNFT *WorknetNFTRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WorknetNFT.Contract.WorknetNFTCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WorknetNFT *WorknetNFTRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WorknetNFT.Contract.WorknetNFTTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WorknetNFT *WorknetNFTRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WorknetNFT.Contract.WorknetNFTTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WorknetNFT *WorknetNFTCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WorknetNFT.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WorknetNFT *WorknetNFTTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WorknetNFT.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WorknetNFT *WorknetNFTTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WorknetNFT.Contract.contract.Transact(opts, method, params...)
}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_WorknetNFT *WorknetNFTCaller) AwpRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WorknetNFT.contract.Call(opts, &out, "awpRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_WorknetNFT *WorknetNFTSession) AwpRegistry() (common.Address, error) {
	return _WorknetNFT.Contract.AwpRegistry(&_WorknetNFT.CallOpts)
}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_WorknetNFT *WorknetNFTCallerSession) AwpRegistry() (common.Address, error) {
	return _WorknetNFT.Contract.AwpRegistry(&_WorknetNFT.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_WorknetNFT *WorknetNFTCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WorknetNFT.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_WorknetNFT *WorknetNFTSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _WorknetNFT.Contract.BalanceOf(&_WorknetNFT.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_WorknetNFT *WorknetNFTCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _WorknetNFT.Contract.BalanceOf(&_WorknetNFT.CallOpts, owner)
}

// GetAlphaToken is a free data retrieval call binding the contract method 0xc7bc8ec6.
//
// Solidity: function getAlphaToken(uint256 tokenId) view returns(address)
func (_WorknetNFT *WorknetNFTCaller) GetAlphaToken(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _WorknetNFT.contract.Call(opts, &out, "getAlphaToken", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAlphaToken is a free data retrieval call binding the contract method 0xc7bc8ec6.
//
// Solidity: function getAlphaToken(uint256 tokenId) view returns(address)
func (_WorknetNFT *WorknetNFTSession) GetAlphaToken(tokenId *big.Int) (common.Address, error) {
	return _WorknetNFT.Contract.GetAlphaToken(&_WorknetNFT.CallOpts, tokenId)
}

// GetAlphaToken is a free data retrieval call binding the contract method 0xc7bc8ec6.
//
// Solidity: function getAlphaToken(uint256 tokenId) view returns(address)
func (_WorknetNFT *WorknetNFTCallerSession) GetAlphaToken(tokenId *big.Int) (common.Address, error) {
	return _WorknetNFT.Contract.GetAlphaToken(&_WorknetNFT.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_WorknetNFT *WorknetNFTCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _WorknetNFT.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_WorknetNFT *WorknetNFTSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _WorknetNFT.Contract.GetApproved(&_WorknetNFT.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_WorknetNFT *WorknetNFTCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _WorknetNFT.Contract.GetApproved(&_WorknetNFT.CallOpts, tokenId)
}

// GetMinStake is a free data retrieval call binding the contract method 0x73f231e7.
//
// Solidity: function getMinStake(uint256 tokenId) view returns(uint128)
func (_WorknetNFT *WorknetNFTCaller) GetMinStake(opts *bind.CallOpts, tokenId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _WorknetNFT.contract.Call(opts, &out, "getMinStake", tokenId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinStake is a free data retrieval call binding the contract method 0x73f231e7.
//
// Solidity: function getMinStake(uint256 tokenId) view returns(uint128)
func (_WorknetNFT *WorknetNFTSession) GetMinStake(tokenId *big.Int) (*big.Int, error) {
	return _WorknetNFT.Contract.GetMinStake(&_WorknetNFT.CallOpts, tokenId)
}

// GetMinStake is a free data retrieval call binding the contract method 0x73f231e7.
//
// Solidity: function getMinStake(uint256 tokenId) view returns(uint128)
func (_WorknetNFT *WorknetNFTCallerSession) GetMinStake(tokenId *big.Int) (*big.Int, error) {
	return _WorknetNFT.Contract.GetMinStake(&_WorknetNFT.CallOpts, tokenId)
}

// GetWorknetData is a free data retrieval call binding the contract method 0x927979a0.
//
// Solidity: function getWorknetData(uint256 tokenId) view returns((string,address,address,string,uint128,address))
func (_WorknetNFT *WorknetNFTCaller) GetWorknetData(opts *bind.CallOpts, tokenId *big.Int) (WorknetNFTWorknetData, error) {
	var out []interface{}
	err := _WorknetNFT.contract.Call(opts, &out, "getWorknetData", tokenId)

	if err != nil {
		return *new(WorknetNFTWorknetData), err
	}

	out0 := *abi.ConvertType(out[0], new(WorknetNFTWorknetData)).(*WorknetNFTWorknetData)

	return out0, err

}

// GetWorknetData is a free data retrieval call binding the contract method 0x927979a0.
//
// Solidity: function getWorknetData(uint256 tokenId) view returns((string,address,address,string,uint128,address))
func (_WorknetNFT *WorknetNFTSession) GetWorknetData(tokenId *big.Int) (WorknetNFTWorknetData, error) {
	return _WorknetNFT.Contract.GetWorknetData(&_WorknetNFT.CallOpts, tokenId)
}

// GetWorknetData is a free data retrieval call binding the contract method 0x927979a0.
//
// Solidity: function getWorknetData(uint256 tokenId) view returns((string,address,address,string,uint128,address))
func (_WorknetNFT *WorknetNFTCallerSession) GetWorknetData(tokenId *big.Int) (WorknetNFTWorknetData, error) {
	return _WorknetNFT.Contract.GetWorknetData(&_WorknetNFT.CallOpts, tokenId)
}

// GetWorknetManager is a free data retrieval call binding the contract method 0xbc4d45c6.
//
// Solidity: function getWorknetManager(uint256 tokenId) view returns(address)
func (_WorknetNFT *WorknetNFTCaller) GetWorknetManager(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _WorknetNFT.contract.Call(opts, &out, "getWorknetManager", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetWorknetManager is a free data retrieval call binding the contract method 0xbc4d45c6.
//
// Solidity: function getWorknetManager(uint256 tokenId) view returns(address)
func (_WorknetNFT *WorknetNFTSession) GetWorknetManager(tokenId *big.Int) (common.Address, error) {
	return _WorknetNFT.Contract.GetWorknetManager(&_WorknetNFT.CallOpts, tokenId)
}

// GetWorknetManager is a free data retrieval call binding the contract method 0xbc4d45c6.
//
// Solidity: function getWorknetManager(uint256 tokenId) view returns(address)
func (_WorknetNFT *WorknetNFTCallerSession) GetWorknetManager(tokenId *big.Int) (common.Address, error) {
	return _WorknetNFT.Contract.GetWorknetManager(&_WorknetNFT.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_WorknetNFT *WorknetNFTCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _WorknetNFT.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_WorknetNFT *WorknetNFTSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _WorknetNFT.Contract.IsApprovedForAll(&_WorknetNFT.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_WorknetNFT *WorknetNFTCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _WorknetNFT.Contract.IsApprovedForAll(&_WorknetNFT.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WorknetNFT *WorknetNFTCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WorknetNFT.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WorknetNFT *WorknetNFTSession) Name() (string, error) {
	return _WorknetNFT.Contract.Name(&_WorknetNFT.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WorknetNFT *WorknetNFTCallerSession) Name() (string, error) {
	return _WorknetNFT.Contract.Name(&_WorknetNFT.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_WorknetNFT *WorknetNFTCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _WorknetNFT.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_WorknetNFT *WorknetNFTSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _WorknetNFT.Contract.OwnerOf(&_WorknetNFT.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_WorknetNFT *WorknetNFTCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _WorknetNFT.Contract.OwnerOf(&_WorknetNFT.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_WorknetNFT *WorknetNFTCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _WorknetNFT.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_WorknetNFT *WorknetNFTSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _WorknetNFT.Contract.SupportsInterface(&_WorknetNFT.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_WorknetNFT *WorknetNFTCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _WorknetNFT.Contract.SupportsInterface(&_WorknetNFT.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WorknetNFT *WorknetNFTCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WorknetNFT.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WorknetNFT *WorknetNFTSession) Symbol() (string, error) {
	return _WorknetNFT.Contract.Symbol(&_WorknetNFT.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WorknetNFT *WorknetNFTCallerSession) Symbol() (string, error) {
	return _WorknetNFT.Contract.Symbol(&_WorknetNFT.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_WorknetNFT *WorknetNFTCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _WorknetNFT.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_WorknetNFT *WorknetNFTSession) TokenURI(tokenId *big.Int) (string, error) {
	return _WorknetNFT.Contract.TokenURI(&_WorknetNFT.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_WorknetNFT *WorknetNFTCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _WorknetNFT.Contract.TokenURI(&_WorknetNFT.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_WorknetNFT *WorknetNFTTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _WorknetNFT.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_WorknetNFT *WorknetNFTSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _WorknetNFT.Contract.Approve(&_WorknetNFT.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_WorknetNFT *WorknetNFTTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _WorknetNFT.Contract.Approve(&_WorknetNFT.TransactOpts, to, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_WorknetNFT *WorknetNFTTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _WorknetNFT.contract.Transact(opts, "burn", tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_WorknetNFT *WorknetNFTSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _WorknetNFT.Contract.Burn(&_WorknetNFT.TransactOpts, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_WorknetNFT *WorknetNFTTransactorSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _WorknetNFT.Contract.Burn(&_WorknetNFT.TransactOpts, tokenId)
}

// Mint is a paid mutator transaction binding the contract method 0xa9584608.
//
// Solidity: function mint(address to, uint256 tokenId, string name_, address worknetManager_, address alphaToken_, uint128 minStake_, string skillsURI_) returns()
func (_WorknetNFT *WorknetNFTTransactor) Mint(opts *bind.TransactOpts, to common.Address, tokenId *big.Int, name_ string, worknetManager_ common.Address, alphaToken_ common.Address, minStake_ *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _WorknetNFT.contract.Transact(opts, "mint", to, tokenId, name_, worknetManager_, alphaToken_, minStake_, skillsURI_)
}

// Mint is a paid mutator transaction binding the contract method 0xa9584608.
//
// Solidity: function mint(address to, uint256 tokenId, string name_, address worknetManager_, address alphaToken_, uint128 minStake_, string skillsURI_) returns()
func (_WorknetNFT *WorknetNFTSession) Mint(to common.Address, tokenId *big.Int, name_ string, worknetManager_ common.Address, alphaToken_ common.Address, minStake_ *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _WorknetNFT.Contract.Mint(&_WorknetNFT.TransactOpts, to, tokenId, name_, worknetManager_, alphaToken_, minStake_, skillsURI_)
}

// Mint is a paid mutator transaction binding the contract method 0xa9584608.
//
// Solidity: function mint(address to, uint256 tokenId, string name_, address worknetManager_, address alphaToken_, uint128 minStake_, string skillsURI_) returns()
func (_WorknetNFT *WorknetNFTTransactorSession) Mint(to common.Address, tokenId *big.Int, name_ string, worknetManager_ common.Address, alphaToken_ common.Address, minStake_ *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _WorknetNFT.Contract.Mint(&_WorknetNFT.TransactOpts, to, tokenId, name_, worknetManager_, alphaToken_, minStake_, skillsURI_)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_WorknetNFT *WorknetNFTTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _WorknetNFT.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_WorknetNFT *WorknetNFTSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _WorknetNFT.Contract.SafeTransferFrom(&_WorknetNFT.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_WorknetNFT *WorknetNFTTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _WorknetNFT.Contract.SafeTransferFrom(&_WorknetNFT.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_WorknetNFT *WorknetNFTTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _WorknetNFT.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_WorknetNFT *WorknetNFTSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _WorknetNFT.Contract.SafeTransferFrom0(&_WorknetNFT.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_WorknetNFT *WorknetNFTTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _WorknetNFT.Contract.SafeTransferFrom0(&_WorknetNFT.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_WorknetNFT *WorknetNFTTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _WorknetNFT.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_WorknetNFT *WorknetNFTSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _WorknetNFT.Contract.SetApprovalForAll(&_WorknetNFT.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_WorknetNFT *WorknetNFTTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _WorknetNFT.Contract.SetApprovalForAll(&_WorknetNFT.TransactOpts, operator, approved)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string uri) returns()
func (_WorknetNFT *WorknetNFTTransactor) SetBaseURI(opts *bind.TransactOpts, uri string) (*types.Transaction, error) {
	return _WorknetNFT.contract.Transact(opts, "setBaseURI", uri)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string uri) returns()
func (_WorknetNFT *WorknetNFTSession) SetBaseURI(uri string) (*types.Transaction, error) {
	return _WorknetNFT.Contract.SetBaseURI(&_WorknetNFT.TransactOpts, uri)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string uri) returns()
func (_WorknetNFT *WorknetNFTTransactorSession) SetBaseURI(uri string) (*types.Transaction, error) {
	return _WorknetNFT.Contract.SetBaseURI(&_WorknetNFT.TransactOpts, uri)
}

// SetMetadataURI is a paid mutator transaction binding the contract method 0x087dce94.
//
// Solidity: function setMetadataURI(uint256 tokenId, string metadataURI_) returns()
func (_WorknetNFT *WorknetNFTTransactor) SetMetadataURI(opts *bind.TransactOpts, tokenId *big.Int, metadataURI_ string) (*types.Transaction, error) {
	return _WorknetNFT.contract.Transact(opts, "setMetadataURI", tokenId, metadataURI_)
}

// SetMetadataURI is a paid mutator transaction binding the contract method 0x087dce94.
//
// Solidity: function setMetadataURI(uint256 tokenId, string metadataURI_) returns()
func (_WorknetNFT *WorknetNFTSession) SetMetadataURI(tokenId *big.Int, metadataURI_ string) (*types.Transaction, error) {
	return _WorknetNFT.Contract.SetMetadataURI(&_WorknetNFT.TransactOpts, tokenId, metadataURI_)
}

// SetMetadataURI is a paid mutator transaction binding the contract method 0x087dce94.
//
// Solidity: function setMetadataURI(uint256 tokenId, string metadataURI_) returns()
func (_WorknetNFT *WorknetNFTTransactorSession) SetMetadataURI(tokenId *big.Int, metadataURI_ string) (*types.Transaction, error) {
	return _WorknetNFT.Contract.SetMetadataURI(&_WorknetNFT.TransactOpts, tokenId, metadataURI_)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x63a9bbe5.
//
// Solidity: function setMinStake(uint256 tokenId, uint128 minStake_) returns()
func (_WorknetNFT *WorknetNFTTransactor) SetMinStake(opts *bind.TransactOpts, tokenId *big.Int, minStake_ *big.Int) (*types.Transaction, error) {
	return _WorknetNFT.contract.Transact(opts, "setMinStake", tokenId, minStake_)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x63a9bbe5.
//
// Solidity: function setMinStake(uint256 tokenId, uint128 minStake_) returns()
func (_WorknetNFT *WorknetNFTSession) SetMinStake(tokenId *big.Int, minStake_ *big.Int) (*types.Transaction, error) {
	return _WorknetNFT.Contract.SetMinStake(&_WorknetNFT.TransactOpts, tokenId, minStake_)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x63a9bbe5.
//
// Solidity: function setMinStake(uint256 tokenId, uint128 minStake_) returns()
func (_WorknetNFT *WorknetNFTTransactorSession) SetMinStake(tokenId *big.Int, minStake_ *big.Int) (*types.Transaction, error) {
	return _WorknetNFT.Contract.SetMinStake(&_WorknetNFT.TransactOpts, tokenId, minStake_)
}

// SetSkillsURI is a paid mutator transaction binding the contract method 0x7c2f4cd6.
//
// Solidity: function setSkillsURI(uint256 tokenId, string skillsURI_) returns()
func (_WorknetNFT *WorknetNFTTransactor) SetSkillsURI(opts *bind.TransactOpts, tokenId *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _WorknetNFT.contract.Transact(opts, "setSkillsURI", tokenId, skillsURI_)
}

// SetSkillsURI is a paid mutator transaction binding the contract method 0x7c2f4cd6.
//
// Solidity: function setSkillsURI(uint256 tokenId, string skillsURI_) returns()
func (_WorknetNFT *WorknetNFTSession) SetSkillsURI(tokenId *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _WorknetNFT.Contract.SetSkillsURI(&_WorknetNFT.TransactOpts, tokenId, skillsURI_)
}

// SetSkillsURI is a paid mutator transaction binding the contract method 0x7c2f4cd6.
//
// Solidity: function setSkillsURI(uint256 tokenId, string skillsURI_) returns()
func (_WorknetNFT *WorknetNFTTransactorSession) SetSkillsURI(tokenId *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _WorknetNFT.Contract.SetSkillsURI(&_WorknetNFT.TransactOpts, tokenId, skillsURI_)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_WorknetNFT *WorknetNFTTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _WorknetNFT.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_WorknetNFT *WorknetNFTSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _WorknetNFT.Contract.TransferFrom(&_WorknetNFT.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_WorknetNFT *WorknetNFTTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _WorknetNFT.Contract.TransferFrom(&_WorknetNFT.TransactOpts, from, to, tokenId)
}

// WorknetNFTApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the WorknetNFT contract.
type WorknetNFTApprovalIterator struct {
	Event *WorknetNFTApproval // Event containing the contract specifics and raw log

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
func (it *WorknetNFTApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorknetNFTApproval)
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
		it.Event = new(WorknetNFTApproval)
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
func (it *WorknetNFTApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WorknetNFTApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WorknetNFTApproval represents a Approval event raised by the WorknetNFT contract.
type WorknetNFTApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_WorknetNFT *WorknetNFTFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*WorknetNFTApprovalIterator, error) {

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

	logs, sub, err := _WorknetNFT.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &WorknetNFTApprovalIterator{contract: _WorknetNFT.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_WorknetNFT *WorknetNFTFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *WorknetNFTApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _WorknetNFT.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WorknetNFTApproval)
				if err := _WorknetNFT.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_WorknetNFT *WorknetNFTFilterer) ParseApproval(log types.Log) (*WorknetNFTApproval, error) {
	event := new(WorknetNFTApproval)
	if err := _WorknetNFT.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WorknetNFTApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the WorknetNFT contract.
type WorknetNFTApprovalForAllIterator struct {
	Event *WorknetNFTApprovalForAll // Event containing the contract specifics and raw log

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
func (it *WorknetNFTApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorknetNFTApprovalForAll)
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
		it.Event = new(WorknetNFTApprovalForAll)
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
func (it *WorknetNFTApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WorknetNFTApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WorknetNFTApprovalForAll represents a ApprovalForAll event raised by the WorknetNFT contract.
type WorknetNFTApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_WorknetNFT *WorknetNFTFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*WorknetNFTApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _WorknetNFT.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &WorknetNFTApprovalForAllIterator{contract: _WorknetNFT.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_WorknetNFT *WorknetNFTFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *WorknetNFTApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _WorknetNFT.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WorknetNFTApprovalForAll)
				if err := _WorknetNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_WorknetNFT *WorknetNFTFilterer) ParseApprovalForAll(log types.Log) (*WorknetNFTApprovalForAll, error) {
	event := new(WorknetNFTApprovalForAll)
	if err := _WorknetNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WorknetNFTMetadataURIUpdatedIterator is returned from FilterMetadataURIUpdated and is used to iterate over the raw logs and unpacked data for MetadataURIUpdated events raised by the WorknetNFT contract.
type WorknetNFTMetadataURIUpdatedIterator struct {
	Event *WorknetNFTMetadataURIUpdated // Event containing the contract specifics and raw log

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
func (it *WorknetNFTMetadataURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorknetNFTMetadataURIUpdated)
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
		it.Event = new(WorknetNFTMetadataURIUpdated)
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
func (it *WorknetNFTMetadataURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WorknetNFTMetadataURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WorknetNFTMetadataURIUpdated represents a MetadataURIUpdated event raised by the WorknetNFT contract.
type WorknetNFTMetadataURIUpdated struct {
	TokenId     *big.Int
	MetadataURI string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMetadataURIUpdated is a free log retrieval operation binding the contract event 0xbf65482a576bba07ddf407b0dd39c63d560c7765323c11cc051d4a9413881a61.
//
// Solidity: event MetadataURIUpdated(uint256 indexed tokenId, string metadataURI)
func (_WorknetNFT *WorknetNFTFilterer) FilterMetadataURIUpdated(opts *bind.FilterOpts, tokenId []*big.Int) (*WorknetNFTMetadataURIUpdatedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _WorknetNFT.contract.FilterLogs(opts, "MetadataURIUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &WorknetNFTMetadataURIUpdatedIterator{contract: _WorknetNFT.contract, event: "MetadataURIUpdated", logs: logs, sub: sub}, nil
}

// WatchMetadataURIUpdated is a free log subscription operation binding the contract event 0xbf65482a576bba07ddf407b0dd39c63d560c7765323c11cc051d4a9413881a61.
//
// Solidity: event MetadataURIUpdated(uint256 indexed tokenId, string metadataURI)
func (_WorknetNFT *WorknetNFTFilterer) WatchMetadataURIUpdated(opts *bind.WatchOpts, sink chan<- *WorknetNFTMetadataURIUpdated, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _WorknetNFT.contract.WatchLogs(opts, "MetadataURIUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WorknetNFTMetadataURIUpdated)
				if err := _WorknetNFT.contract.UnpackLog(event, "MetadataURIUpdated", log); err != nil {
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

// ParseMetadataURIUpdated is a log parse operation binding the contract event 0xbf65482a576bba07ddf407b0dd39c63d560c7765323c11cc051d4a9413881a61.
//
// Solidity: event MetadataURIUpdated(uint256 indexed tokenId, string metadataURI)
func (_WorknetNFT *WorknetNFTFilterer) ParseMetadataURIUpdated(log types.Log) (*WorknetNFTMetadataURIUpdated, error) {
	event := new(WorknetNFTMetadataURIUpdated)
	if err := _WorknetNFT.contract.UnpackLog(event, "MetadataURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WorknetNFTMinStakeUpdatedIterator is returned from FilterMinStakeUpdated and is used to iterate over the raw logs and unpacked data for MinStakeUpdated events raised by the WorknetNFT contract.
type WorknetNFTMinStakeUpdatedIterator struct {
	Event *WorknetNFTMinStakeUpdated // Event containing the contract specifics and raw log

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
func (it *WorknetNFTMinStakeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorknetNFTMinStakeUpdated)
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
		it.Event = new(WorknetNFTMinStakeUpdated)
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
func (it *WorknetNFTMinStakeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WorknetNFTMinStakeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WorknetNFTMinStakeUpdated represents a MinStakeUpdated event raised by the WorknetNFT contract.
type WorknetNFTMinStakeUpdated struct {
	TokenId  *big.Int
	MinStake *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMinStakeUpdated is a free log retrieval operation binding the contract event 0xd0b53d029b5624436a948f7e5e2d9854defd8058cb4a20ff51a0ff9599ad6de8.
//
// Solidity: event MinStakeUpdated(uint256 indexed tokenId, uint128 minStake)
func (_WorknetNFT *WorknetNFTFilterer) FilterMinStakeUpdated(opts *bind.FilterOpts, tokenId []*big.Int) (*WorknetNFTMinStakeUpdatedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _WorknetNFT.contract.FilterLogs(opts, "MinStakeUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &WorknetNFTMinStakeUpdatedIterator{contract: _WorknetNFT.contract, event: "MinStakeUpdated", logs: logs, sub: sub}, nil
}

// WatchMinStakeUpdated is a free log subscription operation binding the contract event 0xd0b53d029b5624436a948f7e5e2d9854defd8058cb4a20ff51a0ff9599ad6de8.
//
// Solidity: event MinStakeUpdated(uint256 indexed tokenId, uint128 minStake)
func (_WorknetNFT *WorknetNFTFilterer) WatchMinStakeUpdated(opts *bind.WatchOpts, sink chan<- *WorknetNFTMinStakeUpdated, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _WorknetNFT.contract.WatchLogs(opts, "MinStakeUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WorknetNFTMinStakeUpdated)
				if err := _WorknetNFT.contract.UnpackLog(event, "MinStakeUpdated", log); err != nil {
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
func (_WorknetNFT *WorknetNFTFilterer) ParseMinStakeUpdated(log types.Log) (*WorknetNFTMinStakeUpdated, error) {
	event := new(WorknetNFTMinStakeUpdated)
	if err := _WorknetNFT.contract.UnpackLog(event, "MinStakeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WorknetNFTSkillsURIUpdatedIterator is returned from FilterSkillsURIUpdated and is used to iterate over the raw logs and unpacked data for SkillsURIUpdated events raised by the WorknetNFT contract.
type WorknetNFTSkillsURIUpdatedIterator struct {
	Event *WorknetNFTSkillsURIUpdated // Event containing the contract specifics and raw log

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
func (it *WorknetNFTSkillsURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorknetNFTSkillsURIUpdated)
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
		it.Event = new(WorknetNFTSkillsURIUpdated)
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
func (it *WorknetNFTSkillsURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WorknetNFTSkillsURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WorknetNFTSkillsURIUpdated represents a SkillsURIUpdated event raised by the WorknetNFT contract.
type WorknetNFTSkillsURIUpdated struct {
	TokenId   *big.Int
	SkillsURI string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSkillsURIUpdated is a free log retrieval operation binding the contract event 0xd1332ed84c54e159e7a4245f8e021aff5b3389b685598c228394168ae30d1020.
//
// Solidity: event SkillsURIUpdated(uint256 indexed tokenId, string skillsURI)
func (_WorknetNFT *WorknetNFTFilterer) FilterSkillsURIUpdated(opts *bind.FilterOpts, tokenId []*big.Int) (*WorknetNFTSkillsURIUpdatedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _WorknetNFT.contract.FilterLogs(opts, "SkillsURIUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &WorknetNFTSkillsURIUpdatedIterator{contract: _WorknetNFT.contract, event: "SkillsURIUpdated", logs: logs, sub: sub}, nil
}

// WatchSkillsURIUpdated is a free log subscription operation binding the contract event 0xd1332ed84c54e159e7a4245f8e021aff5b3389b685598c228394168ae30d1020.
//
// Solidity: event SkillsURIUpdated(uint256 indexed tokenId, string skillsURI)
func (_WorknetNFT *WorknetNFTFilterer) WatchSkillsURIUpdated(opts *bind.WatchOpts, sink chan<- *WorknetNFTSkillsURIUpdated, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _WorknetNFT.contract.WatchLogs(opts, "SkillsURIUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WorknetNFTSkillsURIUpdated)
				if err := _WorknetNFT.contract.UnpackLog(event, "SkillsURIUpdated", log); err != nil {
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
func (_WorknetNFT *WorknetNFTFilterer) ParseSkillsURIUpdated(log types.Log) (*WorknetNFTSkillsURIUpdated, error) {
	event := new(WorknetNFTSkillsURIUpdated)
	if err := _WorknetNFT.contract.UnpackLog(event, "SkillsURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WorknetNFTTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the WorknetNFT contract.
type WorknetNFTTransferIterator struct {
	Event *WorknetNFTTransfer // Event containing the contract specifics and raw log

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
func (it *WorknetNFTTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorknetNFTTransfer)
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
		it.Event = new(WorknetNFTTransfer)
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
func (it *WorknetNFTTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WorknetNFTTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WorknetNFTTransfer represents a Transfer event raised by the WorknetNFT contract.
type WorknetNFTTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_WorknetNFT *WorknetNFTFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*WorknetNFTTransferIterator, error) {

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

	logs, sub, err := _WorknetNFT.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &WorknetNFTTransferIterator{contract: _WorknetNFT.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_WorknetNFT *WorknetNFTFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *WorknetNFTTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _WorknetNFT.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WorknetNFTTransfer)
				if err := _WorknetNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_WorknetNFT *WorknetNFTFilterer) ParseTransfer(log types.Log) (*WorknetNFTTransfer, error) {
	event := new(WorknetNFTTransfer)
	if err := _WorknetNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
