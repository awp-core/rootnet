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

// WorknetTokenMetaData contains all meta data concerning the WorknetToken contract.
var WorknetTokenMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DOMAIN_SEPARATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_SUPPLY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approveAndCall\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approveAndCall\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"burnFrom\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createdAt\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"currentMintableLimit\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialized\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"minter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nonces\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"permit\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinter\",\"inputs\":[{\"name\":\"newMinter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supplyAtLock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferAndCall\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferAndCall\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFromAndCall\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFromAndCall\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"worknetId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinterSet\",\"inputs\":[{\"name\":\"newMinter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ERC1363ApproveFailed\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC1363InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1363InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1363TransferFailed\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC1363TransferFromFailed\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC2612ExpiredSignature\",\"inputs\":[{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC2612InvalidSigner\",\"inputs\":[{\"name\":\"signer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ExceedsMaxSupply\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExceedsMintableLimit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAccountNonce\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"currentNonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidShortString\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotMinter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StringTooLong\",\"inputs\":[{\"name\":\"str\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
}

// WorknetTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use WorknetTokenMetaData.ABI instead.
var WorknetTokenABI = WorknetTokenMetaData.ABI

// WorknetToken is an auto generated Go binding around an Ethereum contract.
type WorknetToken struct {
	WorknetTokenCaller     // Read-only binding to the contract
	WorknetTokenTransactor // Write-only binding to the contract
	WorknetTokenFilterer   // Log filterer for contract events
}

// WorknetTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type WorknetTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WorknetTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WorknetTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WorknetTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WorknetTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WorknetTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WorknetTokenSession struct {
	Contract     *WorknetToken     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WorknetTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WorknetTokenCallerSession struct {
	Contract *WorknetTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// WorknetTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WorknetTokenTransactorSession struct {
	Contract     *WorknetTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// WorknetTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type WorknetTokenRaw struct {
	Contract *WorknetToken // Generic contract binding to access the raw methods on
}

// WorknetTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WorknetTokenCallerRaw struct {
	Contract *WorknetTokenCaller // Generic read-only contract binding to access the raw methods on
}

// WorknetTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WorknetTokenTransactorRaw struct {
	Contract *WorknetTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWorknetToken creates a new instance of WorknetToken, bound to a specific deployed contract.
func NewWorknetToken(address common.Address, backend bind.ContractBackend) (*WorknetToken, error) {
	contract, err := bindWorknetToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WorknetToken{WorknetTokenCaller: WorknetTokenCaller{contract: contract}, WorknetTokenTransactor: WorknetTokenTransactor{contract: contract}, WorknetTokenFilterer: WorknetTokenFilterer{contract: contract}}, nil
}

// NewWorknetTokenCaller creates a new read-only instance of WorknetToken, bound to a specific deployed contract.
func NewWorknetTokenCaller(address common.Address, caller bind.ContractCaller) (*WorknetTokenCaller, error) {
	contract, err := bindWorknetToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WorknetTokenCaller{contract: contract}, nil
}

// NewWorknetTokenTransactor creates a new write-only instance of WorknetToken, bound to a specific deployed contract.
func NewWorknetTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*WorknetTokenTransactor, error) {
	contract, err := bindWorknetToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WorknetTokenTransactor{contract: contract}, nil
}

// NewWorknetTokenFilterer creates a new log filterer instance of WorknetToken, bound to a specific deployed contract.
func NewWorknetTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*WorknetTokenFilterer, error) {
	contract, err := bindWorknetToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WorknetTokenFilterer{contract: contract}, nil
}

// bindWorknetToken binds a generic wrapper to an already deployed contract.
func bindWorknetToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WorknetTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WorknetToken *WorknetTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WorknetToken.Contract.WorknetTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WorknetToken *WorknetTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WorknetToken.Contract.WorknetTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WorknetToken *WorknetTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WorknetToken.Contract.WorknetTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WorknetToken *WorknetTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WorknetToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WorknetToken *WorknetTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WorknetToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WorknetToken *WorknetTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WorknetToken.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_WorknetToken *WorknetTokenCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _WorknetToken.contract.Call(opts, &out, "DOMAIN_SEPARATOR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_WorknetToken *WorknetTokenSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _WorknetToken.Contract.DOMAINSEPARATOR(&_WorknetToken.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_WorknetToken *WorknetTokenCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _WorknetToken.Contract.DOMAINSEPARATOR(&_WorknetToken.CallOpts)
}

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_WorknetToken *WorknetTokenCaller) MAXSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WorknetToken.contract.Call(opts, &out, "MAX_SUPPLY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_WorknetToken *WorknetTokenSession) MAXSUPPLY() (*big.Int, error) {
	return _WorknetToken.Contract.MAXSUPPLY(&_WorknetToken.CallOpts)
}

// MAXSUPPLY is a free data retrieval call binding the contract method 0x32cb6b0c.
//
// Solidity: function MAX_SUPPLY() view returns(uint256)
func (_WorknetToken *WorknetTokenCallerSession) MAXSUPPLY() (*big.Int, error) {
	return _WorknetToken.Contract.MAXSUPPLY(&_WorknetToken.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WorknetToken *WorknetTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WorknetToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WorknetToken *WorknetTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _WorknetToken.Contract.Allowance(&_WorknetToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_WorknetToken *WorknetTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _WorknetToken.Contract.Allowance(&_WorknetToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WorknetToken *WorknetTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WorknetToken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WorknetToken *WorknetTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _WorknetToken.Contract.BalanceOf(&_WorknetToken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_WorknetToken *WorknetTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _WorknetToken.Contract.BalanceOf(&_WorknetToken.CallOpts, account)
}

// CreatedAt is a free data retrieval call binding the contract method 0xcf09e0d0.
//
// Solidity: function createdAt() view returns(uint64)
func (_WorknetToken *WorknetTokenCaller) CreatedAt(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _WorknetToken.contract.Call(opts, &out, "createdAt")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// CreatedAt is a free data retrieval call binding the contract method 0xcf09e0d0.
//
// Solidity: function createdAt() view returns(uint64)
func (_WorknetToken *WorknetTokenSession) CreatedAt() (uint64, error) {
	return _WorknetToken.Contract.CreatedAt(&_WorknetToken.CallOpts)
}

// CreatedAt is a free data retrieval call binding the contract method 0xcf09e0d0.
//
// Solidity: function createdAt() view returns(uint64)
func (_WorknetToken *WorknetTokenCallerSession) CreatedAt() (uint64, error) {
	return _WorknetToken.Contract.CreatedAt(&_WorknetToken.CallOpts)
}

// CurrentMintableLimit is a free data retrieval call binding the contract method 0xfd3088ae.
//
// Solidity: function currentMintableLimit() view returns(uint256)
func (_WorknetToken *WorknetTokenCaller) CurrentMintableLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WorknetToken.contract.Call(opts, &out, "currentMintableLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentMintableLimit is a free data retrieval call binding the contract method 0xfd3088ae.
//
// Solidity: function currentMintableLimit() view returns(uint256)
func (_WorknetToken *WorknetTokenSession) CurrentMintableLimit() (*big.Int, error) {
	return _WorknetToken.Contract.CurrentMintableLimit(&_WorknetToken.CallOpts)
}

// CurrentMintableLimit is a free data retrieval call binding the contract method 0xfd3088ae.
//
// Solidity: function currentMintableLimit() view returns(uint256)
func (_WorknetToken *WorknetTokenCallerSession) CurrentMintableLimit() (*big.Int, error) {
	return _WorknetToken.Contract.CurrentMintableLimit(&_WorknetToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WorknetToken *WorknetTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _WorknetToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WorknetToken *WorknetTokenSession) Decimals() (uint8, error) {
	return _WorknetToken.Contract.Decimals(&_WorknetToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WorknetToken *WorknetTokenCallerSession) Decimals() (uint8, error) {
	return _WorknetToken.Contract.Decimals(&_WorknetToken.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_WorknetToken *WorknetTokenCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _WorknetToken.contract.Call(opts, &out, "eip712Domain")

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
func (_WorknetToken *WorknetTokenSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _WorknetToken.Contract.Eip712Domain(&_WorknetToken.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_WorknetToken *WorknetTokenCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _WorknetToken.Contract.Eip712Domain(&_WorknetToken.CallOpts)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_WorknetToken *WorknetTokenCaller) Initialized(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _WorknetToken.contract.Call(opts, &out, "initialized")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_WorknetToken *WorknetTokenSession) Initialized() (bool, error) {
	return _WorknetToken.Contract.Initialized(&_WorknetToken.CallOpts)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_WorknetToken *WorknetTokenCallerSession) Initialized() (bool, error) {
	return _WorknetToken.Contract.Initialized(&_WorknetToken.CallOpts)
}

// Minter is a free data retrieval call binding the contract method 0x07546172.
//
// Solidity: function minter() view returns(address)
func (_WorknetToken *WorknetTokenCaller) Minter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _WorknetToken.contract.Call(opts, &out, "minter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Minter is a free data retrieval call binding the contract method 0x07546172.
//
// Solidity: function minter() view returns(address)
func (_WorknetToken *WorknetTokenSession) Minter() (common.Address, error) {
	return _WorknetToken.Contract.Minter(&_WorknetToken.CallOpts)
}

// Minter is a free data retrieval call binding the contract method 0x07546172.
//
// Solidity: function minter() view returns(address)
func (_WorknetToken *WorknetTokenCallerSession) Minter() (common.Address, error) {
	return _WorknetToken.Contract.Minter(&_WorknetToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WorknetToken *WorknetTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WorknetToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WorknetToken *WorknetTokenSession) Name() (string, error) {
	return _WorknetToken.Contract.Name(&_WorknetToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WorknetToken *WorknetTokenCallerSession) Name() (string, error) {
	return _WorknetToken.Contract.Name(&_WorknetToken.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_WorknetToken *WorknetTokenCaller) Nonces(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WorknetToken.contract.Call(opts, &out, "nonces", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_WorknetToken *WorknetTokenSession) Nonces(owner common.Address) (*big.Int, error) {
	return _WorknetToken.Contract.Nonces(&_WorknetToken.CallOpts, owner)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_WorknetToken *WorknetTokenCallerSession) Nonces(owner common.Address) (*big.Int, error) {
	return _WorknetToken.Contract.Nonces(&_WorknetToken.CallOpts, owner)
}

// SupplyAtLock is a free data retrieval call binding the contract method 0xea32345b.
//
// Solidity: function supplyAtLock() view returns(uint256)
func (_WorknetToken *WorknetTokenCaller) SupplyAtLock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WorknetToken.contract.Call(opts, &out, "supplyAtLock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SupplyAtLock is a free data retrieval call binding the contract method 0xea32345b.
//
// Solidity: function supplyAtLock() view returns(uint256)
func (_WorknetToken *WorknetTokenSession) SupplyAtLock() (*big.Int, error) {
	return _WorknetToken.Contract.SupplyAtLock(&_WorknetToken.CallOpts)
}

// SupplyAtLock is a free data retrieval call binding the contract method 0xea32345b.
//
// Solidity: function supplyAtLock() view returns(uint256)
func (_WorknetToken *WorknetTokenCallerSession) SupplyAtLock() (*big.Int, error) {
	return _WorknetToken.Contract.SupplyAtLock(&_WorknetToken.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_WorknetToken *WorknetTokenCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _WorknetToken.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_WorknetToken *WorknetTokenSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _WorknetToken.Contract.SupportsInterface(&_WorknetToken.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_WorknetToken *WorknetTokenCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _WorknetToken.Contract.SupportsInterface(&_WorknetToken.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WorknetToken *WorknetTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WorknetToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WorknetToken *WorknetTokenSession) Symbol() (string, error) {
	return _WorknetToken.Contract.Symbol(&_WorknetToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WorknetToken *WorknetTokenCallerSession) Symbol() (string, error) {
	return _WorknetToken.Contract.Symbol(&_WorknetToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WorknetToken *WorknetTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WorknetToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WorknetToken *WorknetTokenSession) TotalSupply() (*big.Int, error) {
	return _WorknetToken.Contract.TotalSupply(&_WorknetToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WorknetToken *WorknetTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _WorknetToken.Contract.TotalSupply(&_WorknetToken.CallOpts)
}

// WorknetId is a free data retrieval call binding the contract method 0x08e5d167.
//
// Solidity: function worknetId() view returns(uint256)
func (_WorknetToken *WorknetTokenCaller) WorknetId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WorknetToken.contract.Call(opts, &out, "worknetId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WorknetId is a free data retrieval call binding the contract method 0x08e5d167.
//
// Solidity: function worknetId() view returns(uint256)
func (_WorknetToken *WorknetTokenSession) WorknetId() (*big.Int, error) {
	return _WorknetToken.Contract.WorknetId(&_WorknetToken.CallOpts)
}

// WorknetId is a free data retrieval call binding the contract method 0x08e5d167.
//
// Solidity: function worknetId() view returns(uint256)
func (_WorknetToken *WorknetTokenCallerSession) WorknetId() (*big.Int, error) {
	return _WorknetToken.Contract.WorknetId(&_WorknetToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.Approve(&_WorknetToken.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.Approve(&_WorknetToken.TransactOpts, spender, value)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0x3177029f.
//
// Solidity: function approveAndCall(address spender, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenTransactor) ApproveAndCall(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.contract.Transact(opts, "approveAndCall", spender, value)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0x3177029f.
//
// Solidity: function approveAndCall(address spender, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenSession) ApproveAndCall(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.ApproveAndCall(&_WorknetToken.TransactOpts, spender, value)
}

// ApproveAndCall is a paid mutator transaction binding the contract method 0x3177029f.
//
// Solidity: function approveAndCall(address spender, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenTransactorSession) ApproveAndCall(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.ApproveAndCall(&_WorknetToken.TransactOpts, spender, value)
}

// ApproveAndCall0 is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(address spender, uint256 value, bytes data) returns(bool)
func (_WorknetToken *WorknetTokenTransactor) ApproveAndCall0(opts *bind.TransactOpts, spender common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _WorknetToken.contract.Transact(opts, "approveAndCall0", spender, value, data)
}

// ApproveAndCall0 is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(address spender, uint256 value, bytes data) returns(bool)
func (_WorknetToken *WorknetTokenSession) ApproveAndCall0(spender common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _WorknetToken.Contract.ApproveAndCall0(&_WorknetToken.TransactOpts, spender, value, data)
}

// ApproveAndCall0 is a paid mutator transaction binding the contract method 0xcae9ca51.
//
// Solidity: function approveAndCall(address spender, uint256 value, bytes data) returns(bool)
func (_WorknetToken *WorknetTokenTransactorSession) ApproveAndCall0(spender common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _WorknetToken.Contract.ApproveAndCall0(&_WorknetToken.TransactOpts, spender, value, data)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_WorknetToken *WorknetTokenTransactor) Burn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.contract.Transact(opts, "burn", value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_WorknetToken *WorknetTokenSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.Burn(&_WorknetToken.TransactOpts, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns()
func (_WorknetToken *WorknetTokenTransactorSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.Burn(&_WorknetToken.TransactOpts, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_WorknetToken *WorknetTokenTransactor) BurnFrom(opts *bind.TransactOpts, account common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.contract.Transact(opts, "burnFrom", account, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_WorknetToken *WorknetTokenSession) BurnFrom(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.BurnFrom(&_WorknetToken.TransactOpts, account, value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 value) returns()
func (_WorknetToken *WorknetTokenTransactorSession) BurnFrom(account common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.BurnFrom(&_WorknetToken.TransactOpts, account, value)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_WorknetToken *WorknetTokenTransactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WorknetToken.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_WorknetToken *WorknetTokenSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.Mint(&_WorknetToken.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_WorknetToken *WorknetTokenTransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.Mint(&_WorknetToken.TransactOpts, to, amount)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_WorknetToken *WorknetTokenTransactor) Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _WorknetToken.contract.Transact(opts, "permit", owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_WorknetToken *WorknetTokenSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _WorknetToken.Contract.Permit(&_WorknetToken.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_WorknetToken *WorknetTokenTransactorSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _WorknetToken.Contract.Permit(&_WorknetToken.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// SetMinter is a paid mutator transaction binding the contract method 0xfca3b5aa.
//
// Solidity: function setMinter(address newMinter) returns()
func (_WorknetToken *WorknetTokenTransactor) SetMinter(opts *bind.TransactOpts, newMinter common.Address) (*types.Transaction, error) {
	return _WorknetToken.contract.Transact(opts, "setMinter", newMinter)
}

// SetMinter is a paid mutator transaction binding the contract method 0xfca3b5aa.
//
// Solidity: function setMinter(address newMinter) returns()
func (_WorknetToken *WorknetTokenSession) SetMinter(newMinter common.Address) (*types.Transaction, error) {
	return _WorknetToken.Contract.SetMinter(&_WorknetToken.TransactOpts, newMinter)
}

// SetMinter is a paid mutator transaction binding the contract method 0xfca3b5aa.
//
// Solidity: function setMinter(address newMinter) returns()
func (_WorknetToken *WorknetTokenTransactorSession) SetMinter(newMinter common.Address) (*types.Transaction, error) {
	return _WorknetToken.Contract.SetMinter(&_WorknetToken.TransactOpts, newMinter)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.Transfer(&_WorknetToken.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.Transfer(&_WorknetToken.TransactOpts, to, value)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x1296ee62.
//
// Solidity: function transferAndCall(address to, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenTransactor) TransferAndCall(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.contract.Transact(opts, "transferAndCall", to, value)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x1296ee62.
//
// Solidity: function transferAndCall(address to, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenSession) TransferAndCall(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.TransferAndCall(&_WorknetToken.TransactOpts, to, value)
}

// TransferAndCall is a paid mutator transaction binding the contract method 0x1296ee62.
//
// Solidity: function transferAndCall(address to, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenTransactorSession) TransferAndCall(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.TransferAndCall(&_WorknetToken.TransactOpts, to, value)
}

// TransferAndCall0 is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 value, bytes data) returns(bool)
func (_WorknetToken *WorknetTokenTransactor) TransferAndCall0(opts *bind.TransactOpts, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _WorknetToken.contract.Transact(opts, "transferAndCall0", to, value, data)
}

// TransferAndCall0 is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 value, bytes data) returns(bool)
func (_WorknetToken *WorknetTokenSession) TransferAndCall0(to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _WorknetToken.Contract.TransferAndCall0(&_WorknetToken.TransactOpts, to, value, data)
}

// TransferAndCall0 is a paid mutator transaction binding the contract method 0x4000aea0.
//
// Solidity: function transferAndCall(address to, uint256 value, bytes data) returns(bool)
func (_WorknetToken *WorknetTokenTransactorSession) TransferAndCall0(to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _WorknetToken.Contract.TransferAndCall0(&_WorknetToken.TransactOpts, to, value, data)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.TransferFrom(&_WorknetToken.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.TransferFrom(&_WorknetToken.TransactOpts, from, to, value)
}

// TransferFromAndCall is a paid mutator transaction binding the contract method 0xc1d34b89.
//
// Solidity: function transferFromAndCall(address from, address to, uint256 value, bytes data) returns(bool)
func (_WorknetToken *WorknetTokenTransactor) TransferFromAndCall(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _WorknetToken.contract.Transact(opts, "transferFromAndCall", from, to, value, data)
}

// TransferFromAndCall is a paid mutator transaction binding the contract method 0xc1d34b89.
//
// Solidity: function transferFromAndCall(address from, address to, uint256 value, bytes data) returns(bool)
func (_WorknetToken *WorknetTokenSession) TransferFromAndCall(from common.Address, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _WorknetToken.Contract.TransferFromAndCall(&_WorknetToken.TransactOpts, from, to, value, data)
}

// TransferFromAndCall is a paid mutator transaction binding the contract method 0xc1d34b89.
//
// Solidity: function transferFromAndCall(address from, address to, uint256 value, bytes data) returns(bool)
func (_WorknetToken *WorknetTokenTransactorSession) TransferFromAndCall(from common.Address, to common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _WorknetToken.Contract.TransferFromAndCall(&_WorknetToken.TransactOpts, from, to, value, data)
}

// TransferFromAndCall0 is a paid mutator transaction binding the contract method 0xd8fbe994.
//
// Solidity: function transferFromAndCall(address from, address to, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenTransactor) TransferFromAndCall0(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.contract.Transact(opts, "transferFromAndCall0", from, to, value)
}

// TransferFromAndCall0 is a paid mutator transaction binding the contract method 0xd8fbe994.
//
// Solidity: function transferFromAndCall(address from, address to, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenSession) TransferFromAndCall0(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.TransferFromAndCall0(&_WorknetToken.TransactOpts, from, to, value)
}

// TransferFromAndCall0 is a paid mutator transaction binding the contract method 0xd8fbe994.
//
// Solidity: function transferFromAndCall(address from, address to, uint256 value) returns(bool)
func (_WorknetToken *WorknetTokenTransactorSession) TransferFromAndCall0(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _WorknetToken.Contract.TransferFromAndCall0(&_WorknetToken.TransactOpts, from, to, value)
}

// WorknetTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the WorknetToken contract.
type WorknetTokenApprovalIterator struct {
	Event *WorknetTokenApproval // Event containing the contract specifics and raw log

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
func (it *WorknetTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorknetTokenApproval)
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
		it.Event = new(WorknetTokenApproval)
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
func (it *WorknetTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WorknetTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WorknetTokenApproval represents a Approval event raised by the WorknetToken contract.
type WorknetTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_WorknetToken *WorknetTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*WorknetTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _WorknetToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &WorknetTokenApprovalIterator{contract: _WorknetToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_WorknetToken *WorknetTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *WorknetTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _WorknetToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WorknetTokenApproval)
				if err := _WorknetToken.contract.UnpackLog(event, "Approval", log); err != nil {
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
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_WorknetToken *WorknetTokenFilterer) ParseApproval(log types.Log) (*WorknetTokenApproval, error) {
	event := new(WorknetTokenApproval)
	if err := _WorknetToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WorknetTokenEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the WorknetToken contract.
type WorknetTokenEIP712DomainChangedIterator struct {
	Event *WorknetTokenEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *WorknetTokenEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorknetTokenEIP712DomainChanged)
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
		it.Event = new(WorknetTokenEIP712DomainChanged)
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
func (it *WorknetTokenEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WorknetTokenEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WorknetTokenEIP712DomainChanged represents a EIP712DomainChanged event raised by the WorknetToken contract.
type WorknetTokenEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_WorknetToken *WorknetTokenFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*WorknetTokenEIP712DomainChangedIterator, error) {

	logs, sub, err := _WorknetToken.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &WorknetTokenEIP712DomainChangedIterator{contract: _WorknetToken.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_WorknetToken *WorknetTokenFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *WorknetTokenEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _WorknetToken.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WorknetTokenEIP712DomainChanged)
				if err := _WorknetToken.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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
func (_WorknetToken *WorknetTokenFilterer) ParseEIP712DomainChanged(log types.Log) (*WorknetTokenEIP712DomainChanged, error) {
	event := new(WorknetTokenEIP712DomainChanged)
	if err := _WorknetToken.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WorknetTokenMinterSetIterator is returned from FilterMinterSet and is used to iterate over the raw logs and unpacked data for MinterSet events raised by the WorknetToken contract.
type WorknetTokenMinterSetIterator struct {
	Event *WorknetTokenMinterSet // Event containing the contract specifics and raw log

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
func (it *WorknetTokenMinterSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorknetTokenMinterSet)
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
		it.Event = new(WorknetTokenMinterSet)
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
func (it *WorknetTokenMinterSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WorknetTokenMinterSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WorknetTokenMinterSet represents a MinterSet event raised by the WorknetToken contract.
type WorknetTokenMinterSet struct {
	NewMinter common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterMinterSet is a free log retrieval operation binding the contract event 0x726b590ef91a8c76ad05bbe91a57ef84605276528f49cd47d787f558a4e755b6.
//
// Solidity: event MinterSet(address indexed newMinter)
func (_WorknetToken *WorknetTokenFilterer) FilterMinterSet(opts *bind.FilterOpts, newMinter []common.Address) (*WorknetTokenMinterSetIterator, error) {

	var newMinterRule []interface{}
	for _, newMinterItem := range newMinter {
		newMinterRule = append(newMinterRule, newMinterItem)
	}

	logs, sub, err := _WorknetToken.contract.FilterLogs(opts, "MinterSet", newMinterRule)
	if err != nil {
		return nil, err
	}
	return &WorknetTokenMinterSetIterator{contract: _WorknetToken.contract, event: "MinterSet", logs: logs, sub: sub}, nil
}

// WatchMinterSet is a free log subscription operation binding the contract event 0x726b590ef91a8c76ad05bbe91a57ef84605276528f49cd47d787f558a4e755b6.
//
// Solidity: event MinterSet(address indexed newMinter)
func (_WorknetToken *WorknetTokenFilterer) WatchMinterSet(opts *bind.WatchOpts, sink chan<- *WorknetTokenMinterSet, newMinter []common.Address) (event.Subscription, error) {

	var newMinterRule []interface{}
	for _, newMinterItem := range newMinter {
		newMinterRule = append(newMinterRule, newMinterItem)
	}

	logs, sub, err := _WorknetToken.contract.WatchLogs(opts, "MinterSet", newMinterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WorknetTokenMinterSet)
				if err := _WorknetToken.contract.UnpackLog(event, "MinterSet", log); err != nil {
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

// ParseMinterSet is a log parse operation binding the contract event 0x726b590ef91a8c76ad05bbe91a57ef84605276528f49cd47d787f558a4e755b6.
//
// Solidity: event MinterSet(address indexed newMinter)
func (_WorknetToken *WorknetTokenFilterer) ParseMinterSet(log types.Log) (*WorknetTokenMinterSet, error) {
	event := new(WorknetTokenMinterSet)
	if err := _WorknetToken.contract.UnpackLog(event, "MinterSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WorknetTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the WorknetToken contract.
type WorknetTokenTransferIterator struct {
	Event *WorknetTokenTransfer // Event containing the contract specifics and raw log

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
func (it *WorknetTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WorknetTokenTransfer)
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
		it.Event = new(WorknetTokenTransfer)
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
func (it *WorknetTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WorknetTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WorknetTokenTransfer represents a Transfer event raised by the WorknetToken contract.
type WorknetTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_WorknetToken *WorknetTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*WorknetTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WorknetToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &WorknetTokenTransferIterator{contract: _WorknetToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_WorknetToken *WorknetTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *WorknetTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _WorknetToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WorknetTokenTransfer)
				if err := _WorknetToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_WorknetToken *WorknetTokenFilterer) ParseTransfer(log types.Log) (*WorknetTokenTransfer, error) {
	event := new(WorknetTokenTransfer)
	if err := _WorknetToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
