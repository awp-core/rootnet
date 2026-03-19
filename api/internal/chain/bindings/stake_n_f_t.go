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

// StakeNFTMetaData contains all meta data concerning the StakeNFT contract.
var StakeNFTMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"awpToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakingVault_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"awpRegistry_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MAX_WEIGHT_DURATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_LOCK_DURATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SQRT_MAX_WEIGHT_FACTOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VOTE_WEIGHT_DIVISOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addToPosition\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newLockEndTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"awpRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"awpToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lockDuration\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositFor\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lockDuration\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositWithPermit\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lockDuration\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPositionForVoting\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lockEndTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"createdAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remainingSeconds\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"votingPower\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUserTotalStaked\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUserVotingPower\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"total\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getVotingPower\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"positions\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lockEndTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"createdAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"remainingTime\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakingVault\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalVotingPower\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposited\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"lockEndTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PositionIncreased\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"addedAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newLockEndTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawn\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientUnallocated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LockCannotShorten\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LockMustExceedCurrentTime\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LockNotExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LockTooShort\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAWPRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotTokenOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NothingToUpdate\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PositionExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// StakeNFTABI is the input ABI used to generate the binding from.
// Deprecated: Use StakeNFTMetaData.ABI instead.
var StakeNFTABI = StakeNFTMetaData.ABI

// StakeNFT is an auto generated Go binding around an Ethereum contract.
type StakeNFT struct {
	StakeNFTCaller     // Read-only binding to the contract
	StakeNFTTransactor // Write-only binding to the contract
	StakeNFTFilterer   // Log filterer for contract events
}

// StakeNFTCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakeNFTCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeNFTTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakeNFTTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeNFTFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakeNFTFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakeNFTSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakeNFTSession struct {
	Contract     *StakeNFT         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakeNFTCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakeNFTCallerSession struct {
	Contract *StakeNFTCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// StakeNFTTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakeNFTTransactorSession struct {
	Contract     *StakeNFTTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// StakeNFTRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakeNFTRaw struct {
	Contract *StakeNFT // Generic contract binding to access the raw methods on
}

// StakeNFTCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakeNFTCallerRaw struct {
	Contract *StakeNFTCaller // Generic read-only contract binding to access the raw methods on
}

// StakeNFTTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakeNFTTransactorRaw struct {
	Contract *StakeNFTTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakeNFT creates a new instance of StakeNFT, bound to a specific deployed contract.
func NewStakeNFT(address common.Address, backend bind.ContractBackend) (*StakeNFT, error) {
	contract, err := bindStakeNFT(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakeNFT{StakeNFTCaller: StakeNFTCaller{contract: contract}, StakeNFTTransactor: StakeNFTTransactor{contract: contract}, StakeNFTFilterer: StakeNFTFilterer{contract: contract}}, nil
}

// NewStakeNFTCaller creates a new read-only instance of StakeNFT, bound to a specific deployed contract.
func NewStakeNFTCaller(address common.Address, caller bind.ContractCaller) (*StakeNFTCaller, error) {
	contract, err := bindStakeNFT(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakeNFTCaller{contract: contract}, nil
}

// NewStakeNFTTransactor creates a new write-only instance of StakeNFT, bound to a specific deployed contract.
func NewStakeNFTTransactor(address common.Address, transactor bind.ContractTransactor) (*StakeNFTTransactor, error) {
	contract, err := bindStakeNFT(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakeNFTTransactor{contract: contract}, nil
}

// NewStakeNFTFilterer creates a new log filterer instance of StakeNFT, bound to a specific deployed contract.
func NewStakeNFTFilterer(address common.Address, filterer bind.ContractFilterer) (*StakeNFTFilterer, error) {
	contract, err := bindStakeNFT(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakeNFTFilterer{contract: contract}, nil
}

// bindStakeNFT binds a generic wrapper to an already deployed contract.
func bindStakeNFT(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StakeNFTMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakeNFT *StakeNFTRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakeNFT.Contract.StakeNFTCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakeNFT *StakeNFTRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeNFT.Contract.StakeNFTTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakeNFT *StakeNFTRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakeNFT.Contract.StakeNFTTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakeNFT *StakeNFTCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakeNFT.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakeNFT *StakeNFTTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakeNFT.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakeNFT *StakeNFTTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakeNFT.Contract.contract.Transact(opts, method, params...)
}

// MAXWEIGHTDURATION is a free data retrieval call binding the contract method 0x597f3d29.
//
// Solidity: function MAX_WEIGHT_DURATION() view returns(uint64)
func (_StakeNFT *StakeNFTCaller) MAXWEIGHTDURATION(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "MAX_WEIGHT_DURATION")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MAXWEIGHTDURATION is a free data retrieval call binding the contract method 0x597f3d29.
//
// Solidity: function MAX_WEIGHT_DURATION() view returns(uint64)
func (_StakeNFT *StakeNFTSession) MAXWEIGHTDURATION() (uint64, error) {
	return _StakeNFT.Contract.MAXWEIGHTDURATION(&_StakeNFT.CallOpts)
}

// MAXWEIGHTDURATION is a free data retrieval call binding the contract method 0x597f3d29.
//
// Solidity: function MAX_WEIGHT_DURATION() view returns(uint64)
func (_StakeNFT *StakeNFTCallerSession) MAXWEIGHTDURATION() (uint64, error) {
	return _StakeNFT.Contract.MAXWEIGHTDURATION(&_StakeNFT.CallOpts)
}

// MINLOCKDURATION is a free data retrieval call binding the contract method 0x78b4330f.
//
// Solidity: function MIN_LOCK_DURATION() view returns(uint64)
func (_StakeNFT *StakeNFTCaller) MINLOCKDURATION(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "MIN_LOCK_DURATION")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MINLOCKDURATION is a free data retrieval call binding the contract method 0x78b4330f.
//
// Solidity: function MIN_LOCK_DURATION() view returns(uint64)
func (_StakeNFT *StakeNFTSession) MINLOCKDURATION() (uint64, error) {
	return _StakeNFT.Contract.MINLOCKDURATION(&_StakeNFT.CallOpts)
}

// MINLOCKDURATION is a free data retrieval call binding the contract method 0x78b4330f.
//
// Solidity: function MIN_LOCK_DURATION() view returns(uint64)
func (_StakeNFT *StakeNFTCallerSession) MINLOCKDURATION() (uint64, error) {
	return _StakeNFT.Contract.MINLOCKDURATION(&_StakeNFT.CallOpts)
}

// SQRTMAXWEIGHTFACTOR is a free data retrieval call binding the contract method 0xbbfcec68.
//
// Solidity: function SQRT_MAX_WEIGHT_FACTOR() view returns(uint256)
func (_StakeNFT *StakeNFTCaller) SQRTMAXWEIGHTFACTOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "SQRT_MAX_WEIGHT_FACTOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SQRTMAXWEIGHTFACTOR is a free data retrieval call binding the contract method 0xbbfcec68.
//
// Solidity: function SQRT_MAX_WEIGHT_FACTOR() view returns(uint256)
func (_StakeNFT *StakeNFTSession) SQRTMAXWEIGHTFACTOR() (*big.Int, error) {
	return _StakeNFT.Contract.SQRTMAXWEIGHTFACTOR(&_StakeNFT.CallOpts)
}

// SQRTMAXWEIGHTFACTOR is a free data retrieval call binding the contract method 0xbbfcec68.
//
// Solidity: function SQRT_MAX_WEIGHT_FACTOR() view returns(uint256)
func (_StakeNFT *StakeNFTCallerSession) SQRTMAXWEIGHTFACTOR() (*big.Int, error) {
	return _StakeNFT.Contract.SQRTMAXWEIGHTFACTOR(&_StakeNFT.CallOpts)
}

// VOTEWEIGHTDIVISOR is a free data retrieval call binding the contract method 0xe4cf87ce.
//
// Solidity: function VOTE_WEIGHT_DIVISOR() view returns(uint256)
func (_StakeNFT *StakeNFTCaller) VOTEWEIGHTDIVISOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "VOTE_WEIGHT_DIVISOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VOTEWEIGHTDIVISOR is a free data retrieval call binding the contract method 0xe4cf87ce.
//
// Solidity: function VOTE_WEIGHT_DIVISOR() view returns(uint256)
func (_StakeNFT *StakeNFTSession) VOTEWEIGHTDIVISOR() (*big.Int, error) {
	return _StakeNFT.Contract.VOTEWEIGHTDIVISOR(&_StakeNFT.CallOpts)
}

// VOTEWEIGHTDIVISOR is a free data retrieval call binding the contract method 0xe4cf87ce.
//
// Solidity: function VOTE_WEIGHT_DIVISOR() view returns(uint256)
func (_StakeNFT *StakeNFTCallerSession) VOTEWEIGHTDIVISOR() (*big.Int, error) {
	return _StakeNFT.Contract.VOTEWEIGHTDIVISOR(&_StakeNFT.CallOpts)
}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_StakeNFT *StakeNFTCaller) AwpRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "awpRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_StakeNFT *StakeNFTSession) AwpRegistry() (common.Address, error) {
	return _StakeNFT.Contract.AwpRegistry(&_StakeNFT.CallOpts)
}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_StakeNFT *StakeNFTCallerSession) AwpRegistry() (common.Address, error) {
	return _StakeNFT.Contract.AwpRegistry(&_StakeNFT.CallOpts)
}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_StakeNFT *StakeNFTCaller) AwpToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "awpToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_StakeNFT *StakeNFTSession) AwpToken() (common.Address, error) {
	return _StakeNFT.Contract.AwpToken(&_StakeNFT.CallOpts)
}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_StakeNFT *StakeNFTCallerSession) AwpToken() (common.Address, error) {
	return _StakeNFT.Contract.AwpToken(&_StakeNFT.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_StakeNFT *StakeNFTCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_StakeNFT *StakeNFTSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _StakeNFT.Contract.BalanceOf(&_StakeNFT.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_StakeNFT *StakeNFTCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _StakeNFT.Contract.BalanceOf(&_StakeNFT.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_StakeNFT *StakeNFTCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_StakeNFT *StakeNFTSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _StakeNFT.Contract.GetApproved(&_StakeNFT.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_StakeNFT *StakeNFTCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _StakeNFT.Contract.GetApproved(&_StakeNFT.CallOpts, tokenId)
}

// GetPositionForVoting is a free data retrieval call binding the contract method 0xada53c61.
//
// Solidity: function getPositionForVoting(uint256 tokenId) view returns(address owner, uint128 amount, uint64 lockEndTime, uint64 createdAt, uint64 remainingSeconds, uint256 votingPower)
func (_StakeNFT *StakeNFTCaller) GetPositionForVoting(opts *bind.CallOpts, tokenId *big.Int) (struct {
	Owner            common.Address
	Amount           *big.Int
	LockEndTime      uint64
	CreatedAt        uint64
	RemainingSeconds uint64
	VotingPower      *big.Int
}, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "getPositionForVoting", tokenId)

	outstruct := new(struct {
		Owner            common.Address
		Amount           *big.Int
		LockEndTime      uint64
		CreatedAt        uint64
		RemainingSeconds uint64
		VotingPower      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Owner = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.LockEndTime = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.CreatedAt = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.RemainingSeconds = *abi.ConvertType(out[4], new(uint64)).(*uint64)
	outstruct.VotingPower = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetPositionForVoting is a free data retrieval call binding the contract method 0xada53c61.
//
// Solidity: function getPositionForVoting(uint256 tokenId) view returns(address owner, uint128 amount, uint64 lockEndTime, uint64 createdAt, uint64 remainingSeconds, uint256 votingPower)
func (_StakeNFT *StakeNFTSession) GetPositionForVoting(tokenId *big.Int) (struct {
	Owner            common.Address
	Amount           *big.Int
	LockEndTime      uint64
	CreatedAt        uint64
	RemainingSeconds uint64
	VotingPower      *big.Int
}, error) {
	return _StakeNFT.Contract.GetPositionForVoting(&_StakeNFT.CallOpts, tokenId)
}

// GetPositionForVoting is a free data retrieval call binding the contract method 0xada53c61.
//
// Solidity: function getPositionForVoting(uint256 tokenId) view returns(address owner, uint128 amount, uint64 lockEndTime, uint64 createdAt, uint64 remainingSeconds, uint256 votingPower)
func (_StakeNFT *StakeNFTCallerSession) GetPositionForVoting(tokenId *big.Int) (struct {
	Owner            common.Address
	Amount           *big.Int
	LockEndTime      uint64
	CreatedAt        uint64
	RemainingSeconds uint64
	VotingPower      *big.Int
}, error) {
	return _StakeNFT.Contract.GetPositionForVoting(&_StakeNFT.CallOpts, tokenId)
}

// GetUserTotalStaked is a free data retrieval call binding the contract method 0x6e2f1696.
//
// Solidity: function getUserTotalStaked(address user) view returns(uint256)
func (_StakeNFT *StakeNFTCaller) GetUserTotalStaked(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "getUserTotalStaked", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserTotalStaked is a free data retrieval call binding the contract method 0x6e2f1696.
//
// Solidity: function getUserTotalStaked(address user) view returns(uint256)
func (_StakeNFT *StakeNFTSession) GetUserTotalStaked(user common.Address) (*big.Int, error) {
	return _StakeNFT.Contract.GetUserTotalStaked(&_StakeNFT.CallOpts, user)
}

// GetUserTotalStaked is a free data retrieval call binding the contract method 0x6e2f1696.
//
// Solidity: function getUserTotalStaked(address user) view returns(uint256)
func (_StakeNFT *StakeNFTCallerSession) GetUserTotalStaked(user common.Address) (*big.Int, error) {
	return _StakeNFT.Contract.GetUserTotalStaked(&_StakeNFT.CallOpts, user)
}

// GetUserVotingPower is a free data retrieval call binding the contract method 0xd033c3bb.
//
// Solidity: function getUserVotingPower(address user, uint256[] tokenIds) view returns(uint256 total)
func (_StakeNFT *StakeNFTCaller) GetUserVotingPower(opts *bind.CallOpts, user common.Address, tokenIds []*big.Int) (*big.Int, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "getUserVotingPower", user, tokenIds)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserVotingPower is a free data retrieval call binding the contract method 0xd033c3bb.
//
// Solidity: function getUserVotingPower(address user, uint256[] tokenIds) view returns(uint256 total)
func (_StakeNFT *StakeNFTSession) GetUserVotingPower(user common.Address, tokenIds []*big.Int) (*big.Int, error) {
	return _StakeNFT.Contract.GetUserVotingPower(&_StakeNFT.CallOpts, user, tokenIds)
}

// GetUserVotingPower is a free data retrieval call binding the contract method 0xd033c3bb.
//
// Solidity: function getUserVotingPower(address user, uint256[] tokenIds) view returns(uint256 total)
func (_StakeNFT *StakeNFTCallerSession) GetUserVotingPower(user common.Address, tokenIds []*big.Int) (*big.Int, error) {
	return _StakeNFT.Contract.GetUserVotingPower(&_StakeNFT.CallOpts, user, tokenIds)
}

// GetVotingPower is a free data retrieval call binding the contract method 0x9c6d2976.
//
// Solidity: function getVotingPower(uint256 tokenId) view returns(uint256)
func (_StakeNFT *StakeNFTCaller) GetVotingPower(opts *bind.CallOpts, tokenId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "getVotingPower", tokenId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVotingPower is a free data retrieval call binding the contract method 0x9c6d2976.
//
// Solidity: function getVotingPower(uint256 tokenId) view returns(uint256)
func (_StakeNFT *StakeNFTSession) GetVotingPower(tokenId *big.Int) (*big.Int, error) {
	return _StakeNFT.Contract.GetVotingPower(&_StakeNFT.CallOpts, tokenId)
}

// GetVotingPower is a free data retrieval call binding the contract method 0x9c6d2976.
//
// Solidity: function getVotingPower(uint256 tokenId) view returns(uint256)
func (_StakeNFT *StakeNFTCallerSession) GetVotingPower(tokenId *big.Int) (*big.Int, error) {
	return _StakeNFT.Contract.GetVotingPower(&_StakeNFT.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_StakeNFT *StakeNFTCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_StakeNFT *StakeNFTSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _StakeNFT.Contract.IsApprovedForAll(&_StakeNFT.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_StakeNFT *StakeNFTCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _StakeNFT.Contract.IsApprovedForAll(&_StakeNFT.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_StakeNFT *StakeNFTCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_StakeNFT *StakeNFTSession) Name() (string, error) {
	return _StakeNFT.Contract.Name(&_StakeNFT.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_StakeNFT *StakeNFTCallerSession) Name() (string, error) {
	return _StakeNFT.Contract.Name(&_StakeNFT.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_StakeNFT *StakeNFTCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_StakeNFT *StakeNFTSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _StakeNFT.Contract.OwnerOf(&_StakeNFT.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_StakeNFT *StakeNFTCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _StakeNFT.Contract.OwnerOf(&_StakeNFT.CallOpts, tokenId)
}

// Positions is a free data retrieval call binding the contract method 0x99fbab88.
//
// Solidity: function positions(uint256 ) view returns(uint128 amount, uint64 lockEndTime, uint64 createdAt)
func (_StakeNFT *StakeNFTCaller) Positions(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Amount      *big.Int
	LockEndTime uint64
	CreatedAt   uint64
}, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "positions", arg0)

	outstruct := new(struct {
		Amount      *big.Int
		LockEndTime uint64
		CreatedAt   uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Amount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.LockEndTime = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.CreatedAt = *abi.ConvertType(out[2], new(uint64)).(*uint64)

	return *outstruct, err

}

// Positions is a free data retrieval call binding the contract method 0x99fbab88.
//
// Solidity: function positions(uint256 ) view returns(uint128 amount, uint64 lockEndTime, uint64 createdAt)
func (_StakeNFT *StakeNFTSession) Positions(arg0 *big.Int) (struct {
	Amount      *big.Int
	LockEndTime uint64
	CreatedAt   uint64
}, error) {
	return _StakeNFT.Contract.Positions(&_StakeNFT.CallOpts, arg0)
}

// Positions is a free data retrieval call binding the contract method 0x99fbab88.
//
// Solidity: function positions(uint256 ) view returns(uint128 amount, uint64 lockEndTime, uint64 createdAt)
func (_StakeNFT *StakeNFTCallerSession) Positions(arg0 *big.Int) (struct {
	Amount      *big.Int
	LockEndTime uint64
	CreatedAt   uint64
}, error) {
	return _StakeNFT.Contract.Positions(&_StakeNFT.CallOpts, arg0)
}

// RemainingTime is a free data retrieval call binding the contract method 0x0c64a7f2.
//
// Solidity: function remainingTime(uint256 tokenId) view returns(uint64)
func (_StakeNFT *StakeNFTCaller) RemainingTime(opts *bind.CallOpts, tokenId *big.Int) (uint64, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "remainingTime", tokenId)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// RemainingTime is a free data retrieval call binding the contract method 0x0c64a7f2.
//
// Solidity: function remainingTime(uint256 tokenId) view returns(uint64)
func (_StakeNFT *StakeNFTSession) RemainingTime(tokenId *big.Int) (uint64, error) {
	return _StakeNFT.Contract.RemainingTime(&_StakeNFT.CallOpts, tokenId)
}

// RemainingTime is a free data retrieval call binding the contract method 0x0c64a7f2.
//
// Solidity: function remainingTime(uint256 tokenId) view returns(uint64)
func (_StakeNFT *StakeNFTCallerSession) RemainingTime(tokenId *big.Int) (uint64, error) {
	return _StakeNFT.Contract.RemainingTime(&_StakeNFT.CallOpts, tokenId)
}

// StakingVault is a free data retrieval call binding the contract method 0x24e7964a.
//
// Solidity: function stakingVault() view returns(address)
func (_StakeNFT *StakeNFTCaller) StakingVault(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "stakingVault")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakingVault is a free data retrieval call binding the contract method 0x24e7964a.
//
// Solidity: function stakingVault() view returns(address)
func (_StakeNFT *StakeNFTSession) StakingVault() (common.Address, error) {
	return _StakeNFT.Contract.StakingVault(&_StakeNFT.CallOpts)
}

// StakingVault is a free data retrieval call binding the contract method 0x24e7964a.
//
// Solidity: function stakingVault() view returns(address)
func (_StakeNFT *StakeNFTCallerSession) StakingVault() (common.Address, error) {
	return _StakeNFT.Contract.StakingVault(&_StakeNFT.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_StakeNFT *StakeNFTCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_StakeNFT *StakeNFTSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _StakeNFT.Contract.SupportsInterface(&_StakeNFT.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_StakeNFT *StakeNFTCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _StakeNFT.Contract.SupportsInterface(&_StakeNFT.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_StakeNFT *StakeNFTCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_StakeNFT *StakeNFTSession) Symbol() (string, error) {
	return _StakeNFT.Contract.Symbol(&_StakeNFT.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_StakeNFT *StakeNFTCallerSession) Symbol() (string, error) {
	return _StakeNFT.Contract.Symbol(&_StakeNFT.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_StakeNFT *StakeNFTCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_StakeNFT *StakeNFTSession) TokenURI(tokenId *big.Int) (string, error) {
	return _StakeNFT.Contract.TokenURI(&_StakeNFT.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_StakeNFT *StakeNFTCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _StakeNFT.Contract.TokenURI(&_StakeNFT.CallOpts, tokenId)
}

// TotalVotingPower is a free data retrieval call binding the contract method 0x671b3793.
//
// Solidity: function totalVotingPower() view returns(uint256)
func (_StakeNFT *StakeNFTCaller) TotalVotingPower(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakeNFT.contract.Call(opts, &out, "totalVotingPower")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalVotingPower is a free data retrieval call binding the contract method 0x671b3793.
//
// Solidity: function totalVotingPower() view returns(uint256)
func (_StakeNFT *StakeNFTSession) TotalVotingPower() (*big.Int, error) {
	return _StakeNFT.Contract.TotalVotingPower(&_StakeNFT.CallOpts)
}

// TotalVotingPower is a free data retrieval call binding the contract method 0x671b3793.
//
// Solidity: function totalVotingPower() view returns(uint256)
func (_StakeNFT *StakeNFTCallerSession) TotalVotingPower() (*big.Int, error) {
	return _StakeNFT.Contract.TotalVotingPower(&_StakeNFT.CallOpts)
}

// AddToPosition is a paid mutator transaction binding the contract method 0xd2845e7d.
//
// Solidity: function addToPosition(uint256 tokenId, uint256 amount, uint64 newLockEndTime) returns()
func (_StakeNFT *StakeNFTTransactor) AddToPosition(opts *bind.TransactOpts, tokenId *big.Int, amount *big.Int, newLockEndTime uint64) (*types.Transaction, error) {
	return _StakeNFT.contract.Transact(opts, "addToPosition", tokenId, amount, newLockEndTime)
}

// AddToPosition is a paid mutator transaction binding the contract method 0xd2845e7d.
//
// Solidity: function addToPosition(uint256 tokenId, uint256 amount, uint64 newLockEndTime) returns()
func (_StakeNFT *StakeNFTSession) AddToPosition(tokenId *big.Int, amount *big.Int, newLockEndTime uint64) (*types.Transaction, error) {
	return _StakeNFT.Contract.AddToPosition(&_StakeNFT.TransactOpts, tokenId, amount, newLockEndTime)
}

// AddToPosition is a paid mutator transaction binding the contract method 0xd2845e7d.
//
// Solidity: function addToPosition(uint256 tokenId, uint256 amount, uint64 newLockEndTime) returns()
func (_StakeNFT *StakeNFTTransactorSession) AddToPosition(tokenId *big.Int, amount *big.Int, newLockEndTime uint64) (*types.Transaction, error) {
	return _StakeNFT.Contract.AddToPosition(&_StakeNFT.TransactOpts, tokenId, amount, newLockEndTime)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_StakeNFT *StakeNFTTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _StakeNFT.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_StakeNFT *StakeNFTSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _StakeNFT.Contract.Approve(&_StakeNFT.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_StakeNFT *StakeNFTTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _StakeNFT.Contract.Approve(&_StakeNFT.TransactOpts, to, tokenId)
}

// Deposit is a paid mutator transaction binding the contract method 0x7d552ea6.
//
// Solidity: function deposit(uint256 amount, uint64 lockDuration) returns(uint256 tokenId)
func (_StakeNFT *StakeNFTTransactor) Deposit(opts *bind.TransactOpts, amount *big.Int, lockDuration uint64) (*types.Transaction, error) {
	return _StakeNFT.contract.Transact(opts, "deposit", amount, lockDuration)
}

// Deposit is a paid mutator transaction binding the contract method 0x7d552ea6.
//
// Solidity: function deposit(uint256 amount, uint64 lockDuration) returns(uint256 tokenId)
func (_StakeNFT *StakeNFTSession) Deposit(amount *big.Int, lockDuration uint64) (*types.Transaction, error) {
	return _StakeNFT.Contract.Deposit(&_StakeNFT.TransactOpts, amount, lockDuration)
}

// Deposit is a paid mutator transaction binding the contract method 0x7d552ea6.
//
// Solidity: function deposit(uint256 amount, uint64 lockDuration) returns(uint256 tokenId)
func (_StakeNFT *StakeNFTTransactorSession) Deposit(amount *big.Int, lockDuration uint64) (*types.Transaction, error) {
	return _StakeNFT.Contract.Deposit(&_StakeNFT.TransactOpts, amount, lockDuration)
}

// DepositFor is a paid mutator transaction binding the contract method 0x1dacaf36.
//
// Solidity: function depositFor(address user, uint256 amount, uint64 lockDuration) returns(uint256 tokenId)
func (_StakeNFT *StakeNFTTransactor) DepositFor(opts *bind.TransactOpts, user common.Address, amount *big.Int, lockDuration uint64) (*types.Transaction, error) {
	return _StakeNFT.contract.Transact(opts, "depositFor", user, amount, lockDuration)
}

// DepositFor is a paid mutator transaction binding the contract method 0x1dacaf36.
//
// Solidity: function depositFor(address user, uint256 amount, uint64 lockDuration) returns(uint256 tokenId)
func (_StakeNFT *StakeNFTSession) DepositFor(user common.Address, amount *big.Int, lockDuration uint64) (*types.Transaction, error) {
	return _StakeNFT.Contract.DepositFor(&_StakeNFT.TransactOpts, user, amount, lockDuration)
}

// DepositFor is a paid mutator transaction binding the contract method 0x1dacaf36.
//
// Solidity: function depositFor(address user, uint256 amount, uint64 lockDuration) returns(uint256 tokenId)
func (_StakeNFT *StakeNFTTransactorSession) DepositFor(user common.Address, amount *big.Int, lockDuration uint64) (*types.Transaction, error) {
	return _StakeNFT.Contract.DepositFor(&_StakeNFT.TransactOpts, user, amount, lockDuration)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x8b5cf26b.
//
// Solidity: function depositWithPermit(uint256 amount, uint64 lockDuration, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256 tokenId)
func (_StakeNFT *StakeNFTTransactor) DepositWithPermit(opts *bind.TransactOpts, amount *big.Int, lockDuration uint64, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _StakeNFT.contract.Transact(opts, "depositWithPermit", amount, lockDuration, deadline, v, r, s)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x8b5cf26b.
//
// Solidity: function depositWithPermit(uint256 amount, uint64 lockDuration, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256 tokenId)
func (_StakeNFT *StakeNFTSession) DepositWithPermit(amount *big.Int, lockDuration uint64, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _StakeNFT.Contract.DepositWithPermit(&_StakeNFT.TransactOpts, amount, lockDuration, deadline, v, r, s)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x8b5cf26b.
//
// Solidity: function depositWithPermit(uint256 amount, uint64 lockDuration, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256 tokenId)
func (_StakeNFT *StakeNFTTransactorSession) DepositWithPermit(amount *big.Int, lockDuration uint64, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _StakeNFT.Contract.DepositWithPermit(&_StakeNFT.TransactOpts, amount, lockDuration, deadline, v, r, s)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_StakeNFT *StakeNFTTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _StakeNFT.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_StakeNFT *StakeNFTSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _StakeNFT.Contract.SafeTransferFrom(&_StakeNFT.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_StakeNFT *StakeNFTTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _StakeNFT.Contract.SafeTransferFrom(&_StakeNFT.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_StakeNFT *StakeNFTTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _StakeNFT.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_StakeNFT *StakeNFTSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _StakeNFT.Contract.SafeTransferFrom0(&_StakeNFT.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_StakeNFT *StakeNFTTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _StakeNFT.Contract.SafeTransferFrom0(&_StakeNFT.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_StakeNFT *StakeNFTTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _StakeNFT.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_StakeNFT *StakeNFTSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _StakeNFT.Contract.SetApprovalForAll(&_StakeNFT.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_StakeNFT *StakeNFTTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _StakeNFT.Contract.SetApprovalForAll(&_StakeNFT.TransactOpts, operator, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_StakeNFT *StakeNFTTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _StakeNFT.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_StakeNFT *StakeNFTSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _StakeNFT.Contract.TransferFrom(&_StakeNFT.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_StakeNFT *StakeNFTTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _StakeNFT.Contract.TransferFrom(&_StakeNFT.TransactOpts, from, to, tokenId)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 tokenId) returns()
func (_StakeNFT *StakeNFTTransactor) Withdraw(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _StakeNFT.contract.Transact(opts, "withdraw", tokenId)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 tokenId) returns()
func (_StakeNFT *StakeNFTSession) Withdraw(tokenId *big.Int) (*types.Transaction, error) {
	return _StakeNFT.Contract.Withdraw(&_StakeNFT.TransactOpts, tokenId)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 tokenId) returns()
func (_StakeNFT *StakeNFTTransactorSession) Withdraw(tokenId *big.Int) (*types.Transaction, error) {
	return _StakeNFT.Contract.Withdraw(&_StakeNFT.TransactOpts, tokenId)
}

// StakeNFTApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the StakeNFT contract.
type StakeNFTApprovalIterator struct {
	Event *StakeNFTApproval // Event containing the contract specifics and raw log

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
func (it *StakeNFTApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeNFTApproval)
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
		it.Event = new(StakeNFTApproval)
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
func (it *StakeNFTApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeNFTApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeNFTApproval represents a Approval event raised by the StakeNFT contract.
type StakeNFTApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_StakeNFT *StakeNFTFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*StakeNFTApprovalIterator, error) {

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

	logs, sub, err := _StakeNFT.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &StakeNFTApprovalIterator{contract: _StakeNFT.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_StakeNFT *StakeNFTFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *StakeNFTApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _StakeNFT.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeNFTApproval)
				if err := _StakeNFT.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_StakeNFT *StakeNFTFilterer) ParseApproval(log types.Log) (*StakeNFTApproval, error) {
	event := new(StakeNFTApproval)
	if err := _StakeNFT.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeNFTApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the StakeNFT contract.
type StakeNFTApprovalForAllIterator struct {
	Event *StakeNFTApprovalForAll // Event containing the contract specifics and raw log

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
func (it *StakeNFTApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeNFTApprovalForAll)
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
		it.Event = new(StakeNFTApprovalForAll)
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
func (it *StakeNFTApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeNFTApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeNFTApprovalForAll represents a ApprovalForAll event raised by the StakeNFT contract.
type StakeNFTApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_StakeNFT *StakeNFTFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*StakeNFTApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _StakeNFT.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &StakeNFTApprovalForAllIterator{contract: _StakeNFT.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_StakeNFT *StakeNFTFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *StakeNFTApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _StakeNFT.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeNFTApprovalForAll)
				if err := _StakeNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_StakeNFT *StakeNFTFilterer) ParseApprovalForAll(log types.Log) (*StakeNFTApprovalForAll, error) {
	event := new(StakeNFTApprovalForAll)
	if err := _StakeNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeNFTDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the StakeNFT contract.
type StakeNFTDepositedIterator struct {
	Event *StakeNFTDeposited // Event containing the contract specifics and raw log

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
func (it *StakeNFTDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeNFTDeposited)
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
		it.Event = new(StakeNFTDeposited)
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
func (it *StakeNFTDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeNFTDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeNFTDeposited represents a Deposited event raised by the StakeNFT contract.
type StakeNFTDeposited struct {
	User        common.Address
	TokenId     *big.Int
	Amount      *big.Int
	LockEndTime uint64
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x19e7166e374f41f05d851e7f5774e0d8424541e4b4353728a88a4c84fe7ba133.
//
// Solidity: event Deposited(address indexed user, uint256 indexed tokenId, uint256 amount, uint64 lockEndTime)
func (_StakeNFT *StakeNFTFilterer) FilterDeposited(opts *bind.FilterOpts, user []common.Address, tokenId []*big.Int) (*StakeNFTDepositedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _StakeNFT.contract.FilterLogs(opts, "Deposited", userRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &StakeNFTDepositedIterator{contract: _StakeNFT.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x19e7166e374f41f05d851e7f5774e0d8424541e4b4353728a88a4c84fe7ba133.
//
// Solidity: event Deposited(address indexed user, uint256 indexed tokenId, uint256 amount, uint64 lockEndTime)
func (_StakeNFT *StakeNFTFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *StakeNFTDeposited, user []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _StakeNFT.contract.WatchLogs(opts, "Deposited", userRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeNFTDeposited)
				if err := _StakeNFT.contract.UnpackLog(event, "Deposited", log); err != nil {
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

// ParseDeposited is a log parse operation binding the contract event 0x19e7166e374f41f05d851e7f5774e0d8424541e4b4353728a88a4c84fe7ba133.
//
// Solidity: event Deposited(address indexed user, uint256 indexed tokenId, uint256 amount, uint64 lockEndTime)
func (_StakeNFT *StakeNFTFilterer) ParseDeposited(log types.Log) (*StakeNFTDeposited, error) {
	event := new(StakeNFTDeposited)
	if err := _StakeNFT.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeNFTPositionIncreasedIterator is returned from FilterPositionIncreased and is used to iterate over the raw logs and unpacked data for PositionIncreased events raised by the StakeNFT contract.
type StakeNFTPositionIncreasedIterator struct {
	Event *StakeNFTPositionIncreased // Event containing the contract specifics and raw log

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
func (it *StakeNFTPositionIncreasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeNFTPositionIncreased)
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
		it.Event = new(StakeNFTPositionIncreased)
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
func (it *StakeNFTPositionIncreasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeNFTPositionIncreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeNFTPositionIncreased represents a PositionIncreased event raised by the StakeNFT contract.
type StakeNFTPositionIncreased struct {
	TokenId        *big.Int
	AddedAmount    *big.Int
	NewLockEndTime uint64
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterPositionIncreased is a free log retrieval operation binding the contract event 0x5f3a37f541da6e3c701efe8ac2d6f6d9070b50349f6aaf0ee403ba4206ad3614.
//
// Solidity: event PositionIncreased(uint256 indexed tokenId, uint256 addedAmount, uint64 newLockEndTime)
func (_StakeNFT *StakeNFTFilterer) FilterPositionIncreased(opts *bind.FilterOpts, tokenId []*big.Int) (*StakeNFTPositionIncreasedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _StakeNFT.contract.FilterLogs(opts, "PositionIncreased", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &StakeNFTPositionIncreasedIterator{contract: _StakeNFT.contract, event: "PositionIncreased", logs: logs, sub: sub}, nil
}

// WatchPositionIncreased is a free log subscription operation binding the contract event 0x5f3a37f541da6e3c701efe8ac2d6f6d9070b50349f6aaf0ee403ba4206ad3614.
//
// Solidity: event PositionIncreased(uint256 indexed tokenId, uint256 addedAmount, uint64 newLockEndTime)
func (_StakeNFT *StakeNFTFilterer) WatchPositionIncreased(opts *bind.WatchOpts, sink chan<- *StakeNFTPositionIncreased, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _StakeNFT.contract.WatchLogs(opts, "PositionIncreased", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeNFTPositionIncreased)
				if err := _StakeNFT.contract.UnpackLog(event, "PositionIncreased", log); err != nil {
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

// ParsePositionIncreased is a log parse operation binding the contract event 0x5f3a37f541da6e3c701efe8ac2d6f6d9070b50349f6aaf0ee403ba4206ad3614.
//
// Solidity: event PositionIncreased(uint256 indexed tokenId, uint256 addedAmount, uint64 newLockEndTime)
func (_StakeNFT *StakeNFTFilterer) ParsePositionIncreased(log types.Log) (*StakeNFTPositionIncreased, error) {
	event := new(StakeNFTPositionIncreased)
	if err := _StakeNFT.contract.UnpackLog(event, "PositionIncreased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeNFTTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the StakeNFT contract.
type StakeNFTTransferIterator struct {
	Event *StakeNFTTransfer // Event containing the contract specifics and raw log

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
func (it *StakeNFTTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeNFTTransfer)
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
		it.Event = new(StakeNFTTransfer)
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
func (it *StakeNFTTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeNFTTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeNFTTransfer represents a Transfer event raised by the StakeNFT contract.
type StakeNFTTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_StakeNFT *StakeNFTFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*StakeNFTTransferIterator, error) {

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

	logs, sub, err := _StakeNFT.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &StakeNFTTransferIterator{contract: _StakeNFT.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_StakeNFT *StakeNFTFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *StakeNFTTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _StakeNFT.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeNFTTransfer)
				if err := _StakeNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_StakeNFT *StakeNFTFilterer) ParseTransfer(log types.Log) (*StakeNFTTransfer, error) {
	event := new(StakeNFTTransfer)
	if err := _StakeNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakeNFTWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the StakeNFT contract.
type StakeNFTWithdrawnIterator struct {
	Event *StakeNFTWithdrawn // Event containing the contract specifics and raw log

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
func (it *StakeNFTWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakeNFTWithdrawn)
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
		it.Event = new(StakeNFTWithdrawn)
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
func (it *StakeNFTWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakeNFTWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakeNFTWithdrawn represents a Withdrawn event raised by the StakeNFT contract.
type StakeNFTWithdrawn struct {
	User    common.Address
	TokenId *big.Int
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0x92ccf450a286a957af52509bc1c9939d1a6a481783e142e41e2499f0bb66ebc6.
//
// Solidity: event Withdrawn(address indexed user, uint256 indexed tokenId, uint256 amount)
func (_StakeNFT *StakeNFTFilterer) FilterWithdrawn(opts *bind.FilterOpts, user []common.Address, tokenId []*big.Int) (*StakeNFTWithdrawnIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _StakeNFT.contract.FilterLogs(opts, "Withdrawn", userRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &StakeNFTWithdrawnIterator{contract: _StakeNFT.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0x92ccf450a286a957af52509bc1c9939d1a6a481783e142e41e2499f0bb66ebc6.
//
// Solidity: event Withdrawn(address indexed user, uint256 indexed tokenId, uint256 amount)
func (_StakeNFT *StakeNFTFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *StakeNFTWithdrawn, user []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _StakeNFT.contract.WatchLogs(opts, "Withdrawn", userRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakeNFTWithdrawn)
				if err := _StakeNFT.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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

// ParseWithdrawn is a log parse operation binding the contract event 0x92ccf450a286a957af52509bc1c9939d1a6a481783e142e41e2499f0bb66ebc6.
//
// Solidity: event Withdrawn(address indexed user, uint256 indexed tokenId, uint256 amount)
func (_StakeNFT *StakeNFTFilterer) ParseWithdrawn(log types.Log) (*StakeNFTWithdrawn, error) {
	event := new(StakeNFTWithdrawn)
	if err := _StakeNFT.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
