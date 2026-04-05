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

// VeAWPMetaData contains all meta data concerning the VeAWP contract.
var VeAWPMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"awpToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"awpAllocator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"guardian_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MAX_WEIGHT_DURATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_LOCK_DURATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VOTE_WEIGHT_DIVISOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addToPosition\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"newLockEndTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"awpAllocator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"awpToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchWithdraw\",\"inputs\":[{\"name\":\"tokenIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lockDuration\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositWithPermit\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lockDuration\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPositionForVoting\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lockEndTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"createdAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"remainingSeconds\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"votingPower\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUserTotalStaked\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUserVotingPower\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"total\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getVotingPower\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"guardian\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"partialWithdraw\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"positions\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lockEndTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"createdAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"remainingTime\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rescueToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setGuardian\",\"inputs\":[{\"name\":\"g\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalStaked\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalVotingPower\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deposited\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"lockEndTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GuardianUpdated\",\"inputs\":[{\"name\":\"newGuardian\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MetadataUpdate\",\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PositionDecreased\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"withdrawnAmount\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"remainingAmount\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PositionIncreased\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"addedAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newLockEndTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawn\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CannotRescueStakedToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientUnallocated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LockCannotShorten\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LockMustExceedCurrentTime\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LockNotExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LockTooShort\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotGuardian\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotTokenOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NothingToUpdate\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PartialWithdrawExceedsBalance\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PositionExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
}

// VeAWPABI is the input ABI used to generate the binding from.
// Deprecated: Use VeAWPMetaData.ABI instead.
var VeAWPABI = VeAWPMetaData.ABI

// VeAWP is an auto generated Go binding around an Ethereum contract.
type VeAWP struct {
	VeAWPCaller     // Read-only binding to the contract
	VeAWPTransactor // Write-only binding to the contract
	VeAWPFilterer   // Log filterer for contract events
}

// VeAWPCaller is an auto generated read-only Go binding around an Ethereum contract.
type VeAWPCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VeAWPTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VeAWPTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VeAWPFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VeAWPFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VeAWPSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VeAWPSession struct {
	Contract     *VeAWP            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VeAWPCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VeAWPCallerSession struct {
	Contract *VeAWPCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// VeAWPTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VeAWPTransactorSession struct {
	Contract     *VeAWPTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VeAWPRaw is an auto generated low-level Go binding around an Ethereum contract.
type VeAWPRaw struct {
	Contract *VeAWP // Generic contract binding to access the raw methods on
}

// VeAWPCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VeAWPCallerRaw struct {
	Contract *VeAWPCaller // Generic read-only contract binding to access the raw methods on
}

// VeAWPTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VeAWPTransactorRaw struct {
	Contract *VeAWPTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVeAWP creates a new instance of VeAWP, bound to a specific deployed contract.
func NewVeAWP(address common.Address, backend bind.ContractBackend) (*VeAWP, error) {
	contract, err := bindVeAWP(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VeAWP{VeAWPCaller: VeAWPCaller{contract: contract}, VeAWPTransactor: VeAWPTransactor{contract: contract}, VeAWPFilterer: VeAWPFilterer{contract: contract}}, nil
}

// NewVeAWPCaller creates a new read-only instance of VeAWP, bound to a specific deployed contract.
func NewVeAWPCaller(address common.Address, caller bind.ContractCaller) (*VeAWPCaller, error) {
	contract, err := bindVeAWP(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VeAWPCaller{contract: contract}, nil
}

// NewVeAWPTransactor creates a new write-only instance of VeAWP, bound to a specific deployed contract.
func NewVeAWPTransactor(address common.Address, transactor bind.ContractTransactor) (*VeAWPTransactor, error) {
	contract, err := bindVeAWP(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VeAWPTransactor{contract: contract}, nil
}

// NewVeAWPFilterer creates a new log filterer instance of VeAWP, bound to a specific deployed contract.
func NewVeAWPFilterer(address common.Address, filterer bind.ContractFilterer) (*VeAWPFilterer, error) {
	contract, err := bindVeAWP(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VeAWPFilterer{contract: contract}, nil
}

// bindVeAWP binds a generic wrapper to an already deployed contract.
func bindVeAWP(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VeAWPMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VeAWP *VeAWPRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VeAWP.Contract.VeAWPCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VeAWP *VeAWPRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VeAWP.Contract.VeAWPTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VeAWP *VeAWPRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VeAWP.Contract.VeAWPTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VeAWP *VeAWPCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VeAWP.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VeAWP *VeAWPTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VeAWP.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VeAWP *VeAWPTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VeAWP.Contract.contract.Transact(opts, method, params...)
}

// MAXWEIGHTDURATION is a free data retrieval call binding the contract method 0x597f3d29.
//
// Solidity: function MAX_WEIGHT_DURATION() view returns(uint64)
func (_VeAWP *VeAWPCaller) MAXWEIGHTDURATION(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "MAX_WEIGHT_DURATION")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MAXWEIGHTDURATION is a free data retrieval call binding the contract method 0x597f3d29.
//
// Solidity: function MAX_WEIGHT_DURATION() view returns(uint64)
func (_VeAWP *VeAWPSession) MAXWEIGHTDURATION() (uint64, error) {
	return _VeAWP.Contract.MAXWEIGHTDURATION(&_VeAWP.CallOpts)
}

// MAXWEIGHTDURATION is a free data retrieval call binding the contract method 0x597f3d29.
//
// Solidity: function MAX_WEIGHT_DURATION() view returns(uint64)
func (_VeAWP *VeAWPCallerSession) MAXWEIGHTDURATION() (uint64, error) {
	return _VeAWP.Contract.MAXWEIGHTDURATION(&_VeAWP.CallOpts)
}

// MINLOCKDURATION is a free data retrieval call binding the contract method 0x78b4330f.
//
// Solidity: function MIN_LOCK_DURATION() view returns(uint64)
func (_VeAWP *VeAWPCaller) MINLOCKDURATION(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "MIN_LOCK_DURATION")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MINLOCKDURATION is a free data retrieval call binding the contract method 0x78b4330f.
//
// Solidity: function MIN_LOCK_DURATION() view returns(uint64)
func (_VeAWP *VeAWPSession) MINLOCKDURATION() (uint64, error) {
	return _VeAWP.Contract.MINLOCKDURATION(&_VeAWP.CallOpts)
}

// MINLOCKDURATION is a free data retrieval call binding the contract method 0x78b4330f.
//
// Solidity: function MIN_LOCK_DURATION() view returns(uint64)
func (_VeAWP *VeAWPCallerSession) MINLOCKDURATION() (uint64, error) {
	return _VeAWP.Contract.MINLOCKDURATION(&_VeAWP.CallOpts)
}

// VOTEWEIGHTDIVISOR is a free data retrieval call binding the contract method 0xe4cf87ce.
//
// Solidity: function VOTE_WEIGHT_DIVISOR() view returns(uint256)
func (_VeAWP *VeAWPCaller) VOTEWEIGHTDIVISOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "VOTE_WEIGHT_DIVISOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VOTEWEIGHTDIVISOR is a free data retrieval call binding the contract method 0xe4cf87ce.
//
// Solidity: function VOTE_WEIGHT_DIVISOR() view returns(uint256)
func (_VeAWP *VeAWPSession) VOTEWEIGHTDIVISOR() (*big.Int, error) {
	return _VeAWP.Contract.VOTEWEIGHTDIVISOR(&_VeAWP.CallOpts)
}

// VOTEWEIGHTDIVISOR is a free data retrieval call binding the contract method 0xe4cf87ce.
//
// Solidity: function VOTE_WEIGHT_DIVISOR() view returns(uint256)
func (_VeAWP *VeAWPCallerSession) VOTEWEIGHTDIVISOR() (*big.Int, error) {
	return _VeAWP.Contract.VOTEWEIGHTDIVISOR(&_VeAWP.CallOpts)
}

// AwpAllocator is a free data retrieval call binding the contract method 0xb304b059.
//
// Solidity: function awpAllocator() view returns(address)
func (_VeAWP *VeAWPCaller) AwpAllocator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "awpAllocator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpAllocator is a free data retrieval call binding the contract method 0xb304b059.
//
// Solidity: function awpAllocator() view returns(address)
func (_VeAWP *VeAWPSession) AwpAllocator() (common.Address, error) {
	return _VeAWP.Contract.AwpAllocator(&_VeAWP.CallOpts)
}

// AwpAllocator is a free data retrieval call binding the contract method 0xb304b059.
//
// Solidity: function awpAllocator() view returns(address)
func (_VeAWP *VeAWPCallerSession) AwpAllocator() (common.Address, error) {
	return _VeAWP.Contract.AwpAllocator(&_VeAWP.CallOpts)
}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_VeAWP *VeAWPCaller) AwpToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "awpToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_VeAWP *VeAWPSession) AwpToken() (common.Address, error) {
	return _VeAWP.Contract.AwpToken(&_VeAWP.CallOpts)
}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_VeAWP *VeAWPCallerSession) AwpToken() (common.Address, error) {
	return _VeAWP.Contract.AwpToken(&_VeAWP.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_VeAWP *VeAWPCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_VeAWP *VeAWPSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _VeAWP.Contract.BalanceOf(&_VeAWP.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_VeAWP *VeAWPCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _VeAWP.Contract.BalanceOf(&_VeAWP.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_VeAWP *VeAWPCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_VeAWP *VeAWPSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _VeAWP.Contract.GetApproved(&_VeAWP.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_VeAWP *VeAWPCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _VeAWP.Contract.GetApproved(&_VeAWP.CallOpts, tokenId)
}

// GetPositionForVoting is a free data retrieval call binding the contract method 0xada53c61.
//
// Solidity: function getPositionForVoting(uint256 tokenId) view returns(address owner, uint128 amount, uint64 lockEndTime, uint64 createdAt, uint64 remainingSeconds, uint256 votingPower)
func (_VeAWP *VeAWPCaller) GetPositionForVoting(opts *bind.CallOpts, tokenId *big.Int) (struct {
	Owner            common.Address
	Amount           *big.Int
	LockEndTime      uint64
	CreatedAt        uint64
	RemainingSeconds uint64
	VotingPower      *big.Int
}, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "getPositionForVoting", tokenId)

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
func (_VeAWP *VeAWPSession) GetPositionForVoting(tokenId *big.Int) (struct {
	Owner            common.Address
	Amount           *big.Int
	LockEndTime      uint64
	CreatedAt        uint64
	RemainingSeconds uint64
	VotingPower      *big.Int
}, error) {
	return _VeAWP.Contract.GetPositionForVoting(&_VeAWP.CallOpts, tokenId)
}

// GetPositionForVoting is a free data retrieval call binding the contract method 0xada53c61.
//
// Solidity: function getPositionForVoting(uint256 tokenId) view returns(address owner, uint128 amount, uint64 lockEndTime, uint64 createdAt, uint64 remainingSeconds, uint256 votingPower)
func (_VeAWP *VeAWPCallerSession) GetPositionForVoting(tokenId *big.Int) (struct {
	Owner            common.Address
	Amount           *big.Int
	LockEndTime      uint64
	CreatedAt        uint64
	RemainingSeconds uint64
	VotingPower      *big.Int
}, error) {
	return _VeAWP.Contract.GetPositionForVoting(&_VeAWP.CallOpts, tokenId)
}

// GetUserTotalStaked is a free data retrieval call binding the contract method 0x6e2f1696.
//
// Solidity: function getUserTotalStaked(address user) view returns(uint256)
func (_VeAWP *VeAWPCaller) GetUserTotalStaked(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "getUserTotalStaked", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserTotalStaked is a free data retrieval call binding the contract method 0x6e2f1696.
//
// Solidity: function getUserTotalStaked(address user) view returns(uint256)
func (_VeAWP *VeAWPSession) GetUserTotalStaked(user common.Address) (*big.Int, error) {
	return _VeAWP.Contract.GetUserTotalStaked(&_VeAWP.CallOpts, user)
}

// GetUserTotalStaked is a free data retrieval call binding the contract method 0x6e2f1696.
//
// Solidity: function getUserTotalStaked(address user) view returns(uint256)
func (_VeAWP *VeAWPCallerSession) GetUserTotalStaked(user common.Address) (*big.Int, error) {
	return _VeAWP.Contract.GetUserTotalStaked(&_VeAWP.CallOpts, user)
}

// GetUserVotingPower is a free data retrieval call binding the contract method 0xd033c3bb.
//
// Solidity: function getUserVotingPower(address user, uint256[] tokenIds) view returns(uint256 total)
func (_VeAWP *VeAWPCaller) GetUserVotingPower(opts *bind.CallOpts, user common.Address, tokenIds []*big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "getUserVotingPower", user, tokenIds)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserVotingPower is a free data retrieval call binding the contract method 0xd033c3bb.
//
// Solidity: function getUserVotingPower(address user, uint256[] tokenIds) view returns(uint256 total)
func (_VeAWP *VeAWPSession) GetUserVotingPower(user common.Address, tokenIds []*big.Int) (*big.Int, error) {
	return _VeAWP.Contract.GetUserVotingPower(&_VeAWP.CallOpts, user, tokenIds)
}

// GetUserVotingPower is a free data retrieval call binding the contract method 0xd033c3bb.
//
// Solidity: function getUserVotingPower(address user, uint256[] tokenIds) view returns(uint256 total)
func (_VeAWP *VeAWPCallerSession) GetUserVotingPower(user common.Address, tokenIds []*big.Int) (*big.Int, error) {
	return _VeAWP.Contract.GetUserVotingPower(&_VeAWP.CallOpts, user, tokenIds)
}

// GetVotingPower is a free data retrieval call binding the contract method 0x9c6d2976.
//
// Solidity: function getVotingPower(uint256 tokenId) view returns(uint256)
func (_VeAWP *VeAWPCaller) GetVotingPower(opts *bind.CallOpts, tokenId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "getVotingPower", tokenId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVotingPower is a free data retrieval call binding the contract method 0x9c6d2976.
//
// Solidity: function getVotingPower(uint256 tokenId) view returns(uint256)
func (_VeAWP *VeAWPSession) GetVotingPower(tokenId *big.Int) (*big.Int, error) {
	return _VeAWP.Contract.GetVotingPower(&_VeAWP.CallOpts, tokenId)
}

// GetVotingPower is a free data retrieval call binding the contract method 0x9c6d2976.
//
// Solidity: function getVotingPower(uint256 tokenId) view returns(uint256)
func (_VeAWP *VeAWPCallerSession) GetVotingPower(tokenId *big.Int) (*big.Int, error) {
	return _VeAWP.Contract.GetVotingPower(&_VeAWP.CallOpts, tokenId)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_VeAWP *VeAWPCaller) Guardian(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "guardian")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_VeAWP *VeAWPSession) Guardian() (common.Address, error) {
	return _VeAWP.Contract.Guardian(&_VeAWP.CallOpts)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_VeAWP *VeAWPCallerSession) Guardian() (common.Address, error) {
	return _VeAWP.Contract.Guardian(&_VeAWP.CallOpts)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_VeAWP *VeAWPCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_VeAWP *VeAWPSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _VeAWP.Contract.IsApprovedForAll(&_VeAWP.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_VeAWP *VeAWPCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _VeAWP.Contract.IsApprovedForAll(&_VeAWP.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_VeAWP *VeAWPCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_VeAWP *VeAWPSession) Name() (string, error) {
	return _VeAWP.Contract.Name(&_VeAWP.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_VeAWP *VeAWPCallerSession) Name() (string, error) {
	return _VeAWP.Contract.Name(&_VeAWP.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_VeAWP *VeAWPCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_VeAWP *VeAWPSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _VeAWP.Contract.OwnerOf(&_VeAWP.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_VeAWP *VeAWPCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _VeAWP.Contract.OwnerOf(&_VeAWP.CallOpts, tokenId)
}

// Positions is a free data retrieval call binding the contract method 0x99fbab88.
//
// Solidity: function positions(uint256 ) view returns(uint128 amount, uint64 lockEndTime, uint64 createdAt)
func (_VeAWP *VeAWPCaller) Positions(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Amount      *big.Int
	LockEndTime uint64
	CreatedAt   uint64
}, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "positions", arg0)

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
func (_VeAWP *VeAWPSession) Positions(arg0 *big.Int) (struct {
	Amount      *big.Int
	LockEndTime uint64
	CreatedAt   uint64
}, error) {
	return _VeAWP.Contract.Positions(&_VeAWP.CallOpts, arg0)
}

// Positions is a free data retrieval call binding the contract method 0x99fbab88.
//
// Solidity: function positions(uint256 ) view returns(uint128 amount, uint64 lockEndTime, uint64 createdAt)
func (_VeAWP *VeAWPCallerSession) Positions(arg0 *big.Int) (struct {
	Amount      *big.Int
	LockEndTime uint64
	CreatedAt   uint64
}, error) {
	return _VeAWP.Contract.Positions(&_VeAWP.CallOpts, arg0)
}

// RemainingTime is a free data retrieval call binding the contract method 0x0c64a7f2.
//
// Solidity: function remainingTime(uint256 tokenId) view returns(uint64)
func (_VeAWP *VeAWPCaller) RemainingTime(opts *bind.CallOpts, tokenId *big.Int) (uint64, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "remainingTime", tokenId)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// RemainingTime is a free data retrieval call binding the contract method 0x0c64a7f2.
//
// Solidity: function remainingTime(uint256 tokenId) view returns(uint64)
func (_VeAWP *VeAWPSession) RemainingTime(tokenId *big.Int) (uint64, error) {
	return _VeAWP.Contract.RemainingTime(&_VeAWP.CallOpts, tokenId)
}

// RemainingTime is a free data retrieval call binding the contract method 0x0c64a7f2.
//
// Solidity: function remainingTime(uint256 tokenId) view returns(uint64)
func (_VeAWP *VeAWPCallerSession) RemainingTime(tokenId *big.Int) (uint64, error) {
	return _VeAWP.Contract.RemainingTime(&_VeAWP.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_VeAWP *VeAWPCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_VeAWP *VeAWPSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _VeAWP.Contract.SupportsInterface(&_VeAWP.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_VeAWP *VeAWPCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _VeAWP.Contract.SupportsInterface(&_VeAWP.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_VeAWP *VeAWPCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_VeAWP *VeAWPSession) Symbol() (string, error) {
	return _VeAWP.Contract.Symbol(&_VeAWP.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_VeAWP *VeAWPCallerSession) Symbol() (string, error) {
	return _VeAWP.Contract.Symbol(&_VeAWP.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_VeAWP *VeAWPCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_VeAWP *VeAWPSession) TokenURI(tokenId *big.Int) (string, error) {
	return _VeAWP.Contract.TokenURI(&_VeAWP.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_VeAWP *VeAWPCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _VeAWP.Contract.TokenURI(&_VeAWP.CallOpts, tokenId)
}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() view returns(uint256)
func (_VeAWP *VeAWPCaller) TotalStaked(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "totalStaked")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() view returns(uint256)
func (_VeAWP *VeAWPSession) TotalStaked() (*big.Int, error) {
	return _VeAWP.Contract.TotalStaked(&_VeAWP.CallOpts)
}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() view returns(uint256)
func (_VeAWP *VeAWPCallerSession) TotalStaked() (*big.Int, error) {
	return _VeAWP.Contract.TotalStaked(&_VeAWP.CallOpts)
}

// TotalVotingPower is a free data retrieval call binding the contract method 0x671b3793.
//
// Solidity: function totalVotingPower() view returns(uint256)
func (_VeAWP *VeAWPCaller) TotalVotingPower(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeAWP.contract.Call(opts, &out, "totalVotingPower")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalVotingPower is a free data retrieval call binding the contract method 0x671b3793.
//
// Solidity: function totalVotingPower() view returns(uint256)
func (_VeAWP *VeAWPSession) TotalVotingPower() (*big.Int, error) {
	return _VeAWP.Contract.TotalVotingPower(&_VeAWP.CallOpts)
}

// TotalVotingPower is a free data retrieval call binding the contract method 0x671b3793.
//
// Solidity: function totalVotingPower() view returns(uint256)
func (_VeAWP *VeAWPCallerSession) TotalVotingPower() (*big.Int, error) {
	return _VeAWP.Contract.TotalVotingPower(&_VeAWP.CallOpts)
}

// AddToPosition is a paid mutator transaction binding the contract method 0xd2845e7d.
//
// Solidity: function addToPosition(uint256 tokenId, uint256 amount, uint64 newLockEndTime) returns()
func (_VeAWP *VeAWPTransactor) AddToPosition(opts *bind.TransactOpts, tokenId *big.Int, amount *big.Int, newLockEndTime uint64) (*types.Transaction, error) {
	return _VeAWP.contract.Transact(opts, "addToPosition", tokenId, amount, newLockEndTime)
}

// AddToPosition is a paid mutator transaction binding the contract method 0xd2845e7d.
//
// Solidity: function addToPosition(uint256 tokenId, uint256 amount, uint64 newLockEndTime) returns()
func (_VeAWP *VeAWPSession) AddToPosition(tokenId *big.Int, amount *big.Int, newLockEndTime uint64) (*types.Transaction, error) {
	return _VeAWP.Contract.AddToPosition(&_VeAWP.TransactOpts, tokenId, amount, newLockEndTime)
}

// AddToPosition is a paid mutator transaction binding the contract method 0xd2845e7d.
//
// Solidity: function addToPosition(uint256 tokenId, uint256 amount, uint64 newLockEndTime) returns()
func (_VeAWP *VeAWPTransactorSession) AddToPosition(tokenId *big.Int, amount *big.Int, newLockEndTime uint64) (*types.Transaction, error) {
	return _VeAWP.Contract.AddToPosition(&_VeAWP.TransactOpts, tokenId, amount, newLockEndTime)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_VeAWP *VeAWPTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _VeAWP.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_VeAWP *VeAWPSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _VeAWP.Contract.Approve(&_VeAWP.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_VeAWP *VeAWPTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _VeAWP.Contract.Approve(&_VeAWP.TransactOpts, to, tokenId)
}

// BatchWithdraw is a paid mutator transaction binding the contract method 0x72e55399.
//
// Solidity: function batchWithdraw(uint256[] tokenIds) returns()
func (_VeAWP *VeAWPTransactor) BatchWithdraw(opts *bind.TransactOpts, tokenIds []*big.Int) (*types.Transaction, error) {
	return _VeAWP.contract.Transact(opts, "batchWithdraw", tokenIds)
}

// BatchWithdraw is a paid mutator transaction binding the contract method 0x72e55399.
//
// Solidity: function batchWithdraw(uint256[] tokenIds) returns()
func (_VeAWP *VeAWPSession) BatchWithdraw(tokenIds []*big.Int) (*types.Transaction, error) {
	return _VeAWP.Contract.BatchWithdraw(&_VeAWP.TransactOpts, tokenIds)
}

// BatchWithdraw is a paid mutator transaction binding the contract method 0x72e55399.
//
// Solidity: function batchWithdraw(uint256[] tokenIds) returns()
func (_VeAWP *VeAWPTransactorSession) BatchWithdraw(tokenIds []*big.Int) (*types.Transaction, error) {
	return _VeAWP.Contract.BatchWithdraw(&_VeAWP.TransactOpts, tokenIds)
}

// Deposit is a paid mutator transaction binding the contract method 0x7d552ea6.
//
// Solidity: function deposit(uint256 amount, uint64 lockDuration) returns(uint256 tokenId)
func (_VeAWP *VeAWPTransactor) Deposit(opts *bind.TransactOpts, amount *big.Int, lockDuration uint64) (*types.Transaction, error) {
	return _VeAWP.contract.Transact(opts, "deposit", amount, lockDuration)
}

// Deposit is a paid mutator transaction binding the contract method 0x7d552ea6.
//
// Solidity: function deposit(uint256 amount, uint64 lockDuration) returns(uint256 tokenId)
func (_VeAWP *VeAWPSession) Deposit(amount *big.Int, lockDuration uint64) (*types.Transaction, error) {
	return _VeAWP.Contract.Deposit(&_VeAWP.TransactOpts, amount, lockDuration)
}

// Deposit is a paid mutator transaction binding the contract method 0x7d552ea6.
//
// Solidity: function deposit(uint256 amount, uint64 lockDuration) returns(uint256 tokenId)
func (_VeAWP *VeAWPTransactorSession) Deposit(amount *big.Int, lockDuration uint64) (*types.Transaction, error) {
	return _VeAWP.Contract.Deposit(&_VeAWP.TransactOpts, amount, lockDuration)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x8b5cf26b.
//
// Solidity: function depositWithPermit(uint256 amount, uint64 lockDuration, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256 tokenId)
func (_VeAWP *VeAWPTransactor) DepositWithPermit(opts *bind.TransactOpts, amount *big.Int, lockDuration uint64, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _VeAWP.contract.Transact(opts, "depositWithPermit", amount, lockDuration, deadline, v, r, s)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x8b5cf26b.
//
// Solidity: function depositWithPermit(uint256 amount, uint64 lockDuration, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256 tokenId)
func (_VeAWP *VeAWPSession) DepositWithPermit(amount *big.Int, lockDuration uint64, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _VeAWP.Contract.DepositWithPermit(&_VeAWP.TransactOpts, amount, lockDuration, deadline, v, r, s)
}

// DepositWithPermit is a paid mutator transaction binding the contract method 0x8b5cf26b.
//
// Solidity: function depositWithPermit(uint256 amount, uint64 lockDuration, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256 tokenId)
func (_VeAWP *VeAWPTransactorSession) DepositWithPermit(amount *big.Int, lockDuration uint64, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _VeAWP.Contract.DepositWithPermit(&_VeAWP.TransactOpts, amount, lockDuration, deadline, v, r, s)
}

// PartialWithdraw is a paid mutator transaction binding the contract method 0x808fe782.
//
// Solidity: function partialWithdraw(uint256 tokenId, uint128 amount) returns()
func (_VeAWP *VeAWPTransactor) PartialWithdraw(opts *bind.TransactOpts, tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VeAWP.contract.Transact(opts, "partialWithdraw", tokenId, amount)
}

// PartialWithdraw is a paid mutator transaction binding the contract method 0x808fe782.
//
// Solidity: function partialWithdraw(uint256 tokenId, uint128 amount) returns()
func (_VeAWP *VeAWPSession) PartialWithdraw(tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VeAWP.Contract.PartialWithdraw(&_VeAWP.TransactOpts, tokenId, amount)
}

// PartialWithdraw is a paid mutator transaction binding the contract method 0x808fe782.
//
// Solidity: function partialWithdraw(uint256 tokenId, uint128 amount) returns()
func (_VeAWP *VeAWPTransactorSession) PartialWithdraw(tokenId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _VeAWP.Contract.PartialWithdraw(&_VeAWP.TransactOpts, tokenId, amount)
}

// RescueToken is a paid mutator transaction binding the contract method 0xe5711e8b.
//
// Solidity: function rescueToken(address token, address to, uint256 amount) returns()
func (_VeAWP *VeAWPTransactor) RescueToken(opts *bind.TransactOpts, token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VeAWP.contract.Transact(opts, "rescueToken", token, to, amount)
}

// RescueToken is a paid mutator transaction binding the contract method 0xe5711e8b.
//
// Solidity: function rescueToken(address token, address to, uint256 amount) returns()
func (_VeAWP *VeAWPSession) RescueToken(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VeAWP.Contract.RescueToken(&_VeAWP.TransactOpts, token, to, amount)
}

// RescueToken is a paid mutator transaction binding the contract method 0xe5711e8b.
//
// Solidity: function rescueToken(address token, address to, uint256 amount) returns()
func (_VeAWP *VeAWPTransactorSession) RescueToken(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _VeAWP.Contract.RescueToken(&_VeAWP.TransactOpts, token, to, amount)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_VeAWP *VeAWPTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _VeAWP.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_VeAWP *VeAWPSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _VeAWP.Contract.SafeTransferFrom(&_VeAWP.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_VeAWP *VeAWPTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _VeAWP.Contract.SafeTransferFrom(&_VeAWP.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_VeAWP *VeAWPTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _VeAWP.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_VeAWP *VeAWPSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _VeAWP.Contract.SafeTransferFrom0(&_VeAWP.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_VeAWP *VeAWPTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _VeAWP.Contract.SafeTransferFrom0(&_VeAWP.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_VeAWP *VeAWPTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _VeAWP.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_VeAWP *VeAWPSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _VeAWP.Contract.SetApprovalForAll(&_VeAWP.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_VeAWP *VeAWPTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _VeAWP.Contract.SetApprovalForAll(&_VeAWP.TransactOpts, operator, approved)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_VeAWP *VeAWPTransactor) SetGuardian(opts *bind.TransactOpts, g common.Address) (*types.Transaction, error) {
	return _VeAWP.contract.Transact(opts, "setGuardian", g)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_VeAWP *VeAWPSession) SetGuardian(g common.Address) (*types.Transaction, error) {
	return _VeAWP.Contract.SetGuardian(&_VeAWP.TransactOpts, g)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_VeAWP *VeAWPTransactorSession) SetGuardian(g common.Address) (*types.Transaction, error) {
	return _VeAWP.Contract.SetGuardian(&_VeAWP.TransactOpts, g)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_VeAWP *VeAWPTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _VeAWP.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_VeAWP *VeAWPSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _VeAWP.Contract.TransferFrom(&_VeAWP.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_VeAWP *VeAWPTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _VeAWP.Contract.TransferFrom(&_VeAWP.TransactOpts, from, to, tokenId)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 tokenId) returns()
func (_VeAWP *VeAWPTransactor) Withdraw(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _VeAWP.contract.Transact(opts, "withdraw", tokenId)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 tokenId) returns()
func (_VeAWP *VeAWPSession) Withdraw(tokenId *big.Int) (*types.Transaction, error) {
	return _VeAWP.Contract.Withdraw(&_VeAWP.TransactOpts, tokenId)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 tokenId) returns()
func (_VeAWP *VeAWPTransactorSession) Withdraw(tokenId *big.Int) (*types.Transaction, error) {
	return _VeAWP.Contract.Withdraw(&_VeAWP.TransactOpts, tokenId)
}

// VeAWPApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the VeAWP contract.
type VeAWPApprovalIterator struct {
	Event *VeAWPApproval // Event containing the contract specifics and raw log

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
func (it *VeAWPApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VeAWPApproval)
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
		it.Event = new(VeAWPApproval)
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
func (it *VeAWPApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VeAWPApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VeAWPApproval represents a Approval event raised by the VeAWP contract.
type VeAWPApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_VeAWP *VeAWPFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*VeAWPApprovalIterator, error) {

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

	logs, sub, err := _VeAWP.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &VeAWPApprovalIterator{contract: _VeAWP.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_VeAWP *VeAWPFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *VeAWPApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _VeAWP.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VeAWPApproval)
				if err := _VeAWP.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_VeAWP *VeAWPFilterer) ParseApproval(log types.Log) (*VeAWPApproval, error) {
	event := new(VeAWPApproval)
	if err := _VeAWP.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VeAWPApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the VeAWP contract.
type VeAWPApprovalForAllIterator struct {
	Event *VeAWPApprovalForAll // Event containing the contract specifics and raw log

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
func (it *VeAWPApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VeAWPApprovalForAll)
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
		it.Event = new(VeAWPApprovalForAll)
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
func (it *VeAWPApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VeAWPApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VeAWPApprovalForAll represents a ApprovalForAll event raised by the VeAWP contract.
type VeAWPApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_VeAWP *VeAWPFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*VeAWPApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _VeAWP.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &VeAWPApprovalForAllIterator{contract: _VeAWP.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_VeAWP *VeAWPFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *VeAWPApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _VeAWP.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VeAWPApprovalForAll)
				if err := _VeAWP.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_VeAWP *VeAWPFilterer) ParseApprovalForAll(log types.Log) (*VeAWPApprovalForAll, error) {
	event := new(VeAWPApprovalForAll)
	if err := _VeAWP.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VeAWPDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the VeAWP contract.
type VeAWPDepositedIterator struct {
	Event *VeAWPDeposited // Event containing the contract specifics and raw log

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
func (it *VeAWPDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VeAWPDeposited)
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
		it.Event = new(VeAWPDeposited)
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
func (it *VeAWPDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VeAWPDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VeAWPDeposited represents a Deposited event raised by the VeAWP contract.
type VeAWPDeposited struct {
	User        common.Address
	TokenId     *big.Int
	Amount      *big.Int
	LockEndTime uint64
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x19e7166e374f41f05d851e7f5774e0d8424541e4b4353728a88a4c84fe7ba133.
//
// Solidity: event Deposited(address indexed user, uint256 indexed tokenId, uint256 amount, uint64 lockEndTime)
func (_VeAWP *VeAWPFilterer) FilterDeposited(opts *bind.FilterOpts, user []common.Address, tokenId []*big.Int) (*VeAWPDepositedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _VeAWP.contract.FilterLogs(opts, "Deposited", userRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &VeAWPDepositedIterator{contract: _VeAWP.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x19e7166e374f41f05d851e7f5774e0d8424541e4b4353728a88a4c84fe7ba133.
//
// Solidity: event Deposited(address indexed user, uint256 indexed tokenId, uint256 amount, uint64 lockEndTime)
func (_VeAWP *VeAWPFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *VeAWPDeposited, user []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _VeAWP.contract.WatchLogs(opts, "Deposited", userRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VeAWPDeposited)
				if err := _VeAWP.contract.UnpackLog(event, "Deposited", log); err != nil {
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
func (_VeAWP *VeAWPFilterer) ParseDeposited(log types.Log) (*VeAWPDeposited, error) {
	event := new(VeAWPDeposited)
	if err := _VeAWP.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VeAWPGuardianUpdatedIterator is returned from FilterGuardianUpdated and is used to iterate over the raw logs and unpacked data for GuardianUpdated events raised by the VeAWP contract.
type VeAWPGuardianUpdatedIterator struct {
	Event *VeAWPGuardianUpdated // Event containing the contract specifics and raw log

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
func (it *VeAWPGuardianUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VeAWPGuardianUpdated)
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
		it.Event = new(VeAWPGuardianUpdated)
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
func (it *VeAWPGuardianUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VeAWPGuardianUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VeAWPGuardianUpdated represents a GuardianUpdated event raised by the VeAWP contract.
type VeAWPGuardianUpdated struct {
	NewGuardian common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterGuardianUpdated is a free log retrieval operation binding the contract event 0x6bb7ff33e730289800c62ad882105a144a74010d2bdbb9a942544a3005ad55bf.
//
// Solidity: event GuardianUpdated(address indexed newGuardian)
func (_VeAWP *VeAWPFilterer) FilterGuardianUpdated(opts *bind.FilterOpts, newGuardian []common.Address) (*VeAWPGuardianUpdatedIterator, error) {

	var newGuardianRule []interface{}
	for _, newGuardianItem := range newGuardian {
		newGuardianRule = append(newGuardianRule, newGuardianItem)
	}

	logs, sub, err := _VeAWP.contract.FilterLogs(opts, "GuardianUpdated", newGuardianRule)
	if err != nil {
		return nil, err
	}
	return &VeAWPGuardianUpdatedIterator{contract: _VeAWP.contract, event: "GuardianUpdated", logs: logs, sub: sub}, nil
}

// WatchGuardianUpdated is a free log subscription operation binding the contract event 0x6bb7ff33e730289800c62ad882105a144a74010d2bdbb9a942544a3005ad55bf.
//
// Solidity: event GuardianUpdated(address indexed newGuardian)
func (_VeAWP *VeAWPFilterer) WatchGuardianUpdated(opts *bind.WatchOpts, sink chan<- *VeAWPGuardianUpdated, newGuardian []common.Address) (event.Subscription, error) {

	var newGuardianRule []interface{}
	for _, newGuardianItem := range newGuardian {
		newGuardianRule = append(newGuardianRule, newGuardianItem)
	}

	logs, sub, err := _VeAWP.contract.WatchLogs(opts, "GuardianUpdated", newGuardianRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VeAWPGuardianUpdated)
				if err := _VeAWP.contract.UnpackLog(event, "GuardianUpdated", log); err != nil {
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
func (_VeAWP *VeAWPFilterer) ParseGuardianUpdated(log types.Log) (*VeAWPGuardianUpdated, error) {
	event := new(VeAWPGuardianUpdated)
	if err := _VeAWP.contract.UnpackLog(event, "GuardianUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VeAWPMetadataUpdateIterator is returned from FilterMetadataUpdate and is used to iterate over the raw logs and unpacked data for MetadataUpdate events raised by the VeAWP contract.
type VeAWPMetadataUpdateIterator struct {
	Event *VeAWPMetadataUpdate // Event containing the contract specifics and raw log

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
func (it *VeAWPMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VeAWPMetadataUpdate)
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
		it.Event = new(VeAWPMetadataUpdate)
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
func (it *VeAWPMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VeAWPMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VeAWPMetadataUpdate represents a MetadataUpdate event raised by the VeAWP contract.
type VeAWPMetadataUpdate struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMetadataUpdate is a free log retrieval operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_VeAWP *VeAWPFilterer) FilterMetadataUpdate(opts *bind.FilterOpts) (*VeAWPMetadataUpdateIterator, error) {

	logs, sub, err := _VeAWP.contract.FilterLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &VeAWPMetadataUpdateIterator{contract: _VeAWP.contract, event: "MetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchMetadataUpdate is a free log subscription operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_VeAWP *VeAWPFilterer) WatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *VeAWPMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _VeAWP.contract.WatchLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VeAWPMetadataUpdate)
				if err := _VeAWP.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
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

// ParseMetadataUpdate is a log parse operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_VeAWP *VeAWPFilterer) ParseMetadataUpdate(log types.Log) (*VeAWPMetadataUpdate, error) {
	event := new(VeAWPMetadataUpdate)
	if err := _VeAWP.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VeAWPPositionDecreasedIterator is returned from FilterPositionDecreased and is used to iterate over the raw logs and unpacked data for PositionDecreased events raised by the VeAWP contract.
type VeAWPPositionDecreasedIterator struct {
	Event *VeAWPPositionDecreased // Event containing the contract specifics and raw log

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
func (it *VeAWPPositionDecreasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VeAWPPositionDecreased)
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
		it.Event = new(VeAWPPositionDecreased)
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
func (it *VeAWPPositionDecreasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VeAWPPositionDecreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VeAWPPositionDecreased represents a PositionDecreased event raised by the VeAWP contract.
type VeAWPPositionDecreased struct {
	TokenId         *big.Int
	WithdrawnAmount *big.Int
	RemainingAmount *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterPositionDecreased is a free log retrieval operation binding the contract event 0x88af8406386f907cac78cacec66023a69cb9bce0cf27d9c08190be97d041b324.
//
// Solidity: event PositionDecreased(uint256 indexed tokenId, uint128 withdrawnAmount, uint128 remainingAmount)
func (_VeAWP *VeAWPFilterer) FilterPositionDecreased(opts *bind.FilterOpts, tokenId []*big.Int) (*VeAWPPositionDecreasedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _VeAWP.contract.FilterLogs(opts, "PositionDecreased", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &VeAWPPositionDecreasedIterator{contract: _VeAWP.contract, event: "PositionDecreased", logs: logs, sub: sub}, nil
}

// WatchPositionDecreased is a free log subscription operation binding the contract event 0x88af8406386f907cac78cacec66023a69cb9bce0cf27d9c08190be97d041b324.
//
// Solidity: event PositionDecreased(uint256 indexed tokenId, uint128 withdrawnAmount, uint128 remainingAmount)
func (_VeAWP *VeAWPFilterer) WatchPositionDecreased(opts *bind.WatchOpts, sink chan<- *VeAWPPositionDecreased, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _VeAWP.contract.WatchLogs(opts, "PositionDecreased", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VeAWPPositionDecreased)
				if err := _VeAWP.contract.UnpackLog(event, "PositionDecreased", log); err != nil {
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

// ParsePositionDecreased is a log parse operation binding the contract event 0x88af8406386f907cac78cacec66023a69cb9bce0cf27d9c08190be97d041b324.
//
// Solidity: event PositionDecreased(uint256 indexed tokenId, uint128 withdrawnAmount, uint128 remainingAmount)
func (_VeAWP *VeAWPFilterer) ParsePositionDecreased(log types.Log) (*VeAWPPositionDecreased, error) {
	event := new(VeAWPPositionDecreased)
	if err := _VeAWP.contract.UnpackLog(event, "PositionDecreased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VeAWPPositionIncreasedIterator is returned from FilterPositionIncreased and is used to iterate over the raw logs and unpacked data for PositionIncreased events raised by the VeAWP contract.
type VeAWPPositionIncreasedIterator struct {
	Event *VeAWPPositionIncreased // Event containing the contract specifics and raw log

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
func (it *VeAWPPositionIncreasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VeAWPPositionIncreased)
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
		it.Event = new(VeAWPPositionIncreased)
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
func (it *VeAWPPositionIncreasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VeAWPPositionIncreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VeAWPPositionIncreased represents a PositionIncreased event raised by the VeAWP contract.
type VeAWPPositionIncreased struct {
	TokenId        *big.Int
	AddedAmount    *big.Int
	NewLockEndTime uint64
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterPositionIncreased is a free log retrieval operation binding the contract event 0x5f3a37f541da6e3c701efe8ac2d6f6d9070b50349f6aaf0ee403ba4206ad3614.
//
// Solidity: event PositionIncreased(uint256 indexed tokenId, uint256 addedAmount, uint64 newLockEndTime)
func (_VeAWP *VeAWPFilterer) FilterPositionIncreased(opts *bind.FilterOpts, tokenId []*big.Int) (*VeAWPPositionIncreasedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _VeAWP.contract.FilterLogs(opts, "PositionIncreased", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &VeAWPPositionIncreasedIterator{contract: _VeAWP.contract, event: "PositionIncreased", logs: logs, sub: sub}, nil
}

// WatchPositionIncreased is a free log subscription operation binding the contract event 0x5f3a37f541da6e3c701efe8ac2d6f6d9070b50349f6aaf0ee403ba4206ad3614.
//
// Solidity: event PositionIncreased(uint256 indexed tokenId, uint256 addedAmount, uint64 newLockEndTime)
func (_VeAWP *VeAWPFilterer) WatchPositionIncreased(opts *bind.WatchOpts, sink chan<- *VeAWPPositionIncreased, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _VeAWP.contract.WatchLogs(opts, "PositionIncreased", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VeAWPPositionIncreased)
				if err := _VeAWP.contract.UnpackLog(event, "PositionIncreased", log); err != nil {
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
func (_VeAWP *VeAWPFilterer) ParsePositionIncreased(log types.Log) (*VeAWPPositionIncreased, error) {
	event := new(VeAWPPositionIncreased)
	if err := _VeAWP.contract.UnpackLog(event, "PositionIncreased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VeAWPTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the VeAWP contract.
type VeAWPTransferIterator struct {
	Event *VeAWPTransfer // Event containing the contract specifics and raw log

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
func (it *VeAWPTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VeAWPTransfer)
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
		it.Event = new(VeAWPTransfer)
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
func (it *VeAWPTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VeAWPTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VeAWPTransfer represents a Transfer event raised by the VeAWP contract.
type VeAWPTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_VeAWP *VeAWPFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*VeAWPTransferIterator, error) {

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

	logs, sub, err := _VeAWP.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &VeAWPTransferIterator{contract: _VeAWP.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_VeAWP *VeAWPFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *VeAWPTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _VeAWP.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VeAWPTransfer)
				if err := _VeAWP.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_VeAWP *VeAWPFilterer) ParseTransfer(log types.Log) (*VeAWPTransfer, error) {
	event := new(VeAWPTransfer)
	if err := _VeAWP.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VeAWPWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the VeAWP contract.
type VeAWPWithdrawnIterator struct {
	Event *VeAWPWithdrawn // Event containing the contract specifics and raw log

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
func (it *VeAWPWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VeAWPWithdrawn)
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
		it.Event = new(VeAWPWithdrawn)
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
func (it *VeAWPWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VeAWPWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VeAWPWithdrawn represents a Withdrawn event raised by the VeAWP contract.
type VeAWPWithdrawn struct {
	User    common.Address
	TokenId *big.Int
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0x92ccf450a286a957af52509bc1c9939d1a6a481783e142e41e2499f0bb66ebc6.
//
// Solidity: event Withdrawn(address indexed user, uint256 indexed tokenId, uint256 amount)
func (_VeAWP *VeAWPFilterer) FilterWithdrawn(opts *bind.FilterOpts, user []common.Address, tokenId []*big.Int) (*VeAWPWithdrawnIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _VeAWP.contract.FilterLogs(opts, "Withdrawn", userRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &VeAWPWithdrawnIterator{contract: _VeAWP.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0x92ccf450a286a957af52509bc1c9939d1a6a481783e142e41e2499f0bb66ebc6.
//
// Solidity: event Withdrawn(address indexed user, uint256 indexed tokenId, uint256 amount)
func (_VeAWP *VeAWPFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *VeAWPWithdrawn, user []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _VeAWP.contract.WatchLogs(opts, "Withdrawn", userRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VeAWPWithdrawn)
				if err := _VeAWP.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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
func (_VeAWP *VeAWPFilterer) ParseWithdrawn(log types.Log) (*VeAWPWithdrawn, error) {
	event := new(VeAWPWithdrawn)
	if err := _VeAWP.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
