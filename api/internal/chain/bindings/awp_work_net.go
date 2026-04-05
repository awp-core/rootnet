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

// AWPWorkNetWorknetData is an auto generated low-level Go binding around an user-defined struct.
type AWPWorkNetWorknetData struct {
	Name           string
	Symbol         string
	WorknetManager common.Address
	WorknetToken   common.Address
	LpPool         [32]byte
	SkillsURI      string
	MinStake       *big.Int
	ImageURI       string
	MetadataURI    string
	Owner          common.Address
}

// AWPWorkNetWorknetIdentity is an auto generated low-level Go binding around an user-defined struct.
type AWPWorkNetWorknetIdentity struct {
	Name           string
	Symbol         string
	WorknetManager common.Address
	WorknetToken   common.Address
	LpPool         [32]byte
}

// AWPWorkNetWorknetMeta is an auto generated low-level Go binding around an user-defined struct.
type AWPWorkNetWorknetMeta struct {
	SkillsURI   string
	MinStake    *big.Int
	ImageURI    string
	MetadataURI string
}

// AWPWorkNetMetaData contains all meta data concerning the AWPWorkNet contract.
var AWPWorkNetMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"awpRegistry_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MAX_SKILLS_URI_LENGTH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_URI_LENGTH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"awpRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"contractURI\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getLPPool\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMinStake\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWorknetData\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structAWPWorkNet.WorknetData\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"worknetManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"worknetToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lpPool\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"imageURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"metadataURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWorknetIdentity\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structAWPWorkNet.WorknetIdentity\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"worknetManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"worknetToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lpPool\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWorknetManager\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWorknetMeta\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structAWPWorkNet.WorknetMeta\",\"components\":[{\"name\":\"skillsURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"imageURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"metadataURI\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getWorknetToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"guardian\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"guardian_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"worknetManager_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"worknetToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lpPool_\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"minStake_\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"skillsURI_\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBaseURI\",\"inputs\":[{\"name\":\"uri\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setContractURI\",\"inputs\":[{\"name\":\"uri\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setGuardian\",\"inputs\":[{\"name\":\"g\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setImageURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMetadataURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMinStake\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSkillsURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ContractURIUpdated\",\"inputs\":[{\"name\":\"uri\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GuardianUpdated\",\"inputs\":[{\"name\":\"newGuardian\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ImageURIUpdated\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"imageURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MetadataURIUpdated\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"metadataURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MetadataUpdate\",\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinStakeUpdated\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SkillsURIUpdated\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"JsonUnsafeCharacter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAWPRegistry\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAuthorized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotGuardian\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StringTooLong\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StringsInsufficientHexLength\",\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TokenNotExist\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
}

// AWPWorkNetABI is the input ABI used to generate the binding from.
// Deprecated: Use AWPWorkNetMetaData.ABI instead.
var AWPWorkNetABI = AWPWorkNetMetaData.ABI

// AWPWorkNet is an auto generated Go binding around an Ethereum contract.
type AWPWorkNet struct {
	AWPWorkNetCaller     // Read-only binding to the contract
	AWPWorkNetTransactor // Write-only binding to the contract
	AWPWorkNetFilterer   // Log filterer for contract events
}

// AWPWorkNetCaller is an auto generated read-only Go binding around an Ethereum contract.
type AWPWorkNetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AWPWorkNetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AWPWorkNetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AWPWorkNetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AWPWorkNetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AWPWorkNetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AWPWorkNetSession struct {
	Contract     *AWPWorkNet       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AWPWorkNetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AWPWorkNetCallerSession struct {
	Contract *AWPWorkNetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// AWPWorkNetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AWPWorkNetTransactorSession struct {
	Contract     *AWPWorkNetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// AWPWorkNetRaw is an auto generated low-level Go binding around an Ethereum contract.
type AWPWorkNetRaw struct {
	Contract *AWPWorkNet // Generic contract binding to access the raw methods on
}

// AWPWorkNetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AWPWorkNetCallerRaw struct {
	Contract *AWPWorkNetCaller // Generic read-only contract binding to access the raw methods on
}

// AWPWorkNetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AWPWorkNetTransactorRaw struct {
	Contract *AWPWorkNetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAWPWorkNet creates a new instance of AWPWorkNet, bound to a specific deployed contract.
func NewAWPWorkNet(address common.Address, backend bind.ContractBackend) (*AWPWorkNet, error) {
	contract, err := bindAWPWorkNet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AWPWorkNet{AWPWorkNetCaller: AWPWorkNetCaller{contract: contract}, AWPWorkNetTransactor: AWPWorkNetTransactor{contract: contract}, AWPWorkNetFilterer: AWPWorkNetFilterer{contract: contract}}, nil
}

// NewAWPWorkNetCaller creates a new read-only instance of AWPWorkNet, bound to a specific deployed contract.
func NewAWPWorkNetCaller(address common.Address, caller bind.ContractCaller) (*AWPWorkNetCaller, error) {
	contract, err := bindAWPWorkNet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AWPWorkNetCaller{contract: contract}, nil
}

// NewAWPWorkNetTransactor creates a new write-only instance of AWPWorkNet, bound to a specific deployed contract.
func NewAWPWorkNetTransactor(address common.Address, transactor bind.ContractTransactor) (*AWPWorkNetTransactor, error) {
	contract, err := bindAWPWorkNet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AWPWorkNetTransactor{contract: contract}, nil
}

// NewAWPWorkNetFilterer creates a new log filterer instance of AWPWorkNet, bound to a specific deployed contract.
func NewAWPWorkNetFilterer(address common.Address, filterer bind.ContractFilterer) (*AWPWorkNetFilterer, error) {
	contract, err := bindAWPWorkNet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AWPWorkNetFilterer{contract: contract}, nil
}

// bindAWPWorkNet binds a generic wrapper to an already deployed contract.
func bindAWPWorkNet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AWPWorkNetMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AWPWorkNet *AWPWorkNetRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AWPWorkNet.Contract.AWPWorkNetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AWPWorkNet *AWPWorkNetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.AWPWorkNetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AWPWorkNet *AWPWorkNetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.AWPWorkNetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AWPWorkNet *AWPWorkNetCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AWPWorkNet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AWPWorkNet *AWPWorkNetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AWPWorkNet *AWPWorkNetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.contract.Transact(opts, method, params...)
}

// MAXSKILLSURILENGTH is a free data retrieval call binding the contract method 0xdd2cac8e.
//
// Solidity: function MAX_SKILLS_URI_LENGTH() view returns(uint256)
func (_AWPWorkNet *AWPWorkNetCaller) MAXSKILLSURILENGTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "MAX_SKILLS_URI_LENGTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXSKILLSURILENGTH is a free data retrieval call binding the contract method 0xdd2cac8e.
//
// Solidity: function MAX_SKILLS_URI_LENGTH() view returns(uint256)
func (_AWPWorkNet *AWPWorkNetSession) MAXSKILLSURILENGTH() (*big.Int, error) {
	return _AWPWorkNet.Contract.MAXSKILLSURILENGTH(&_AWPWorkNet.CallOpts)
}

// MAXSKILLSURILENGTH is a free data retrieval call binding the contract method 0xdd2cac8e.
//
// Solidity: function MAX_SKILLS_URI_LENGTH() view returns(uint256)
func (_AWPWorkNet *AWPWorkNetCallerSession) MAXSKILLSURILENGTH() (*big.Int, error) {
	return _AWPWorkNet.Contract.MAXSKILLSURILENGTH(&_AWPWorkNet.CallOpts)
}

// MAXURILENGTH is a free data retrieval call binding the contract method 0xaab5a877.
//
// Solidity: function MAX_URI_LENGTH() view returns(uint256)
func (_AWPWorkNet *AWPWorkNetCaller) MAXURILENGTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "MAX_URI_LENGTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXURILENGTH is a free data retrieval call binding the contract method 0xaab5a877.
//
// Solidity: function MAX_URI_LENGTH() view returns(uint256)
func (_AWPWorkNet *AWPWorkNetSession) MAXURILENGTH() (*big.Int, error) {
	return _AWPWorkNet.Contract.MAXURILENGTH(&_AWPWorkNet.CallOpts)
}

// MAXURILENGTH is a free data retrieval call binding the contract method 0xaab5a877.
//
// Solidity: function MAX_URI_LENGTH() view returns(uint256)
func (_AWPWorkNet *AWPWorkNetCallerSession) MAXURILENGTH() (*big.Int, error) {
	return _AWPWorkNet.Contract.MAXURILENGTH(&_AWPWorkNet.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_AWPWorkNet *AWPWorkNetCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_AWPWorkNet *AWPWorkNetSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _AWPWorkNet.Contract.UPGRADEINTERFACEVERSION(&_AWPWorkNet.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_AWPWorkNet *AWPWorkNetCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _AWPWorkNet.Contract.UPGRADEINTERFACEVERSION(&_AWPWorkNet.CallOpts)
}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_AWPWorkNet *AWPWorkNetCaller) AwpRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "awpRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_AWPWorkNet *AWPWorkNetSession) AwpRegistry() (common.Address, error) {
	return _AWPWorkNet.Contract.AwpRegistry(&_AWPWorkNet.CallOpts)
}

// AwpRegistry is a free data retrieval call binding the contract method 0x38fb1eb4.
//
// Solidity: function awpRegistry() view returns(address)
func (_AWPWorkNet *AWPWorkNetCallerSession) AwpRegistry() (common.Address, error) {
	return _AWPWorkNet.Contract.AwpRegistry(&_AWPWorkNet.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_AWPWorkNet *AWPWorkNetCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_AWPWorkNet *AWPWorkNetSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _AWPWorkNet.Contract.BalanceOf(&_AWPWorkNet.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_AWPWorkNet *AWPWorkNetCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _AWPWorkNet.Contract.BalanceOf(&_AWPWorkNet.CallOpts, owner)
}

// ContractURI is a free data retrieval call binding the contract method 0xe8a3d485.
//
// Solidity: function contractURI() view returns(string)
func (_AWPWorkNet *AWPWorkNetCaller) ContractURI(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "contractURI")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ContractURI is a free data retrieval call binding the contract method 0xe8a3d485.
//
// Solidity: function contractURI() view returns(string)
func (_AWPWorkNet *AWPWorkNetSession) ContractURI() (string, error) {
	return _AWPWorkNet.Contract.ContractURI(&_AWPWorkNet.CallOpts)
}

// ContractURI is a free data retrieval call binding the contract method 0xe8a3d485.
//
// Solidity: function contractURI() view returns(string)
func (_AWPWorkNet *AWPWorkNetCallerSession) ContractURI() (string, error) {
	return _AWPWorkNet.Contract.ContractURI(&_AWPWorkNet.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_AWPWorkNet *AWPWorkNetCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_AWPWorkNet *AWPWorkNetSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _AWPWorkNet.Contract.GetApproved(&_AWPWorkNet.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_AWPWorkNet *AWPWorkNetCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _AWPWorkNet.Contract.GetApproved(&_AWPWorkNet.CallOpts, tokenId)
}

// GetLPPool is a free data retrieval call binding the contract method 0x35d69727.
//
// Solidity: function getLPPool(uint256 tokenId) view returns(bytes32)
func (_AWPWorkNet *AWPWorkNetCaller) GetLPPool(opts *bind.CallOpts, tokenId *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "getLPPool", tokenId)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetLPPool is a free data retrieval call binding the contract method 0x35d69727.
//
// Solidity: function getLPPool(uint256 tokenId) view returns(bytes32)
func (_AWPWorkNet *AWPWorkNetSession) GetLPPool(tokenId *big.Int) ([32]byte, error) {
	return _AWPWorkNet.Contract.GetLPPool(&_AWPWorkNet.CallOpts, tokenId)
}

// GetLPPool is a free data retrieval call binding the contract method 0x35d69727.
//
// Solidity: function getLPPool(uint256 tokenId) view returns(bytes32)
func (_AWPWorkNet *AWPWorkNetCallerSession) GetLPPool(tokenId *big.Int) ([32]byte, error) {
	return _AWPWorkNet.Contract.GetLPPool(&_AWPWorkNet.CallOpts, tokenId)
}

// GetMinStake is a free data retrieval call binding the contract method 0x73f231e7.
//
// Solidity: function getMinStake(uint256 tokenId) view returns(uint128)
func (_AWPWorkNet *AWPWorkNetCaller) GetMinStake(opts *bind.CallOpts, tokenId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "getMinStake", tokenId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinStake is a free data retrieval call binding the contract method 0x73f231e7.
//
// Solidity: function getMinStake(uint256 tokenId) view returns(uint128)
func (_AWPWorkNet *AWPWorkNetSession) GetMinStake(tokenId *big.Int) (*big.Int, error) {
	return _AWPWorkNet.Contract.GetMinStake(&_AWPWorkNet.CallOpts, tokenId)
}

// GetMinStake is a free data retrieval call binding the contract method 0x73f231e7.
//
// Solidity: function getMinStake(uint256 tokenId) view returns(uint128)
func (_AWPWorkNet *AWPWorkNetCallerSession) GetMinStake(tokenId *big.Int) (*big.Int, error) {
	return _AWPWorkNet.Contract.GetMinStake(&_AWPWorkNet.CallOpts, tokenId)
}

// GetWorknetData is a free data retrieval call binding the contract method 0x927979a0.
//
// Solidity: function getWorknetData(uint256 tokenId) view returns((string,string,address,address,bytes32,string,uint128,string,string,address))
func (_AWPWorkNet *AWPWorkNetCaller) GetWorknetData(opts *bind.CallOpts, tokenId *big.Int) (AWPWorkNetWorknetData, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "getWorknetData", tokenId)

	if err != nil {
		return *new(AWPWorkNetWorknetData), err
	}

	out0 := *abi.ConvertType(out[0], new(AWPWorkNetWorknetData)).(*AWPWorkNetWorknetData)

	return out0, err

}

// GetWorknetData is a free data retrieval call binding the contract method 0x927979a0.
//
// Solidity: function getWorknetData(uint256 tokenId) view returns((string,string,address,address,bytes32,string,uint128,string,string,address))
func (_AWPWorkNet *AWPWorkNetSession) GetWorknetData(tokenId *big.Int) (AWPWorkNetWorknetData, error) {
	return _AWPWorkNet.Contract.GetWorknetData(&_AWPWorkNet.CallOpts, tokenId)
}

// GetWorknetData is a free data retrieval call binding the contract method 0x927979a0.
//
// Solidity: function getWorknetData(uint256 tokenId) view returns((string,string,address,address,bytes32,string,uint128,string,string,address))
func (_AWPWorkNet *AWPWorkNetCallerSession) GetWorknetData(tokenId *big.Int) (AWPWorkNetWorknetData, error) {
	return _AWPWorkNet.Contract.GetWorknetData(&_AWPWorkNet.CallOpts, tokenId)
}

// GetWorknetIdentity is a free data retrieval call binding the contract method 0xace7cd11.
//
// Solidity: function getWorknetIdentity(uint256 tokenId) view returns((string,string,address,address,bytes32))
func (_AWPWorkNet *AWPWorkNetCaller) GetWorknetIdentity(opts *bind.CallOpts, tokenId *big.Int) (AWPWorkNetWorknetIdentity, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "getWorknetIdentity", tokenId)

	if err != nil {
		return *new(AWPWorkNetWorknetIdentity), err
	}

	out0 := *abi.ConvertType(out[0], new(AWPWorkNetWorknetIdentity)).(*AWPWorkNetWorknetIdentity)

	return out0, err

}

// GetWorknetIdentity is a free data retrieval call binding the contract method 0xace7cd11.
//
// Solidity: function getWorknetIdentity(uint256 tokenId) view returns((string,string,address,address,bytes32))
func (_AWPWorkNet *AWPWorkNetSession) GetWorknetIdentity(tokenId *big.Int) (AWPWorkNetWorknetIdentity, error) {
	return _AWPWorkNet.Contract.GetWorknetIdentity(&_AWPWorkNet.CallOpts, tokenId)
}

// GetWorknetIdentity is a free data retrieval call binding the contract method 0xace7cd11.
//
// Solidity: function getWorknetIdentity(uint256 tokenId) view returns((string,string,address,address,bytes32))
func (_AWPWorkNet *AWPWorkNetCallerSession) GetWorknetIdentity(tokenId *big.Int) (AWPWorkNetWorknetIdentity, error) {
	return _AWPWorkNet.Contract.GetWorknetIdentity(&_AWPWorkNet.CallOpts, tokenId)
}

// GetWorknetManager is a free data retrieval call binding the contract method 0xbc4d45c6.
//
// Solidity: function getWorknetManager(uint256 tokenId) view returns(address)
func (_AWPWorkNet *AWPWorkNetCaller) GetWorknetManager(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "getWorknetManager", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetWorknetManager is a free data retrieval call binding the contract method 0xbc4d45c6.
//
// Solidity: function getWorknetManager(uint256 tokenId) view returns(address)
func (_AWPWorkNet *AWPWorkNetSession) GetWorknetManager(tokenId *big.Int) (common.Address, error) {
	return _AWPWorkNet.Contract.GetWorknetManager(&_AWPWorkNet.CallOpts, tokenId)
}

// GetWorknetManager is a free data retrieval call binding the contract method 0xbc4d45c6.
//
// Solidity: function getWorknetManager(uint256 tokenId) view returns(address)
func (_AWPWorkNet *AWPWorkNetCallerSession) GetWorknetManager(tokenId *big.Int) (common.Address, error) {
	return _AWPWorkNet.Contract.GetWorknetManager(&_AWPWorkNet.CallOpts, tokenId)
}

// GetWorknetMeta is a free data retrieval call binding the contract method 0xfb012085.
//
// Solidity: function getWorknetMeta(uint256 tokenId) view returns((string,uint128,string,string))
func (_AWPWorkNet *AWPWorkNetCaller) GetWorknetMeta(opts *bind.CallOpts, tokenId *big.Int) (AWPWorkNetWorknetMeta, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "getWorknetMeta", tokenId)

	if err != nil {
		return *new(AWPWorkNetWorknetMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(AWPWorkNetWorknetMeta)).(*AWPWorkNetWorknetMeta)

	return out0, err

}

// GetWorknetMeta is a free data retrieval call binding the contract method 0xfb012085.
//
// Solidity: function getWorknetMeta(uint256 tokenId) view returns((string,uint128,string,string))
func (_AWPWorkNet *AWPWorkNetSession) GetWorknetMeta(tokenId *big.Int) (AWPWorkNetWorknetMeta, error) {
	return _AWPWorkNet.Contract.GetWorknetMeta(&_AWPWorkNet.CallOpts, tokenId)
}

// GetWorknetMeta is a free data retrieval call binding the contract method 0xfb012085.
//
// Solidity: function getWorknetMeta(uint256 tokenId) view returns((string,uint128,string,string))
func (_AWPWorkNet *AWPWorkNetCallerSession) GetWorknetMeta(tokenId *big.Int) (AWPWorkNetWorknetMeta, error) {
	return _AWPWorkNet.Contract.GetWorknetMeta(&_AWPWorkNet.CallOpts, tokenId)
}

// GetWorknetToken is a free data retrieval call binding the contract method 0x079d1141.
//
// Solidity: function getWorknetToken(uint256 tokenId) view returns(address)
func (_AWPWorkNet *AWPWorkNetCaller) GetWorknetToken(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "getWorknetToken", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetWorknetToken is a free data retrieval call binding the contract method 0x079d1141.
//
// Solidity: function getWorknetToken(uint256 tokenId) view returns(address)
func (_AWPWorkNet *AWPWorkNetSession) GetWorknetToken(tokenId *big.Int) (common.Address, error) {
	return _AWPWorkNet.Contract.GetWorknetToken(&_AWPWorkNet.CallOpts, tokenId)
}

// GetWorknetToken is a free data retrieval call binding the contract method 0x079d1141.
//
// Solidity: function getWorknetToken(uint256 tokenId) view returns(address)
func (_AWPWorkNet *AWPWorkNetCallerSession) GetWorknetToken(tokenId *big.Int) (common.Address, error) {
	return _AWPWorkNet.Contract.GetWorknetToken(&_AWPWorkNet.CallOpts, tokenId)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_AWPWorkNet *AWPWorkNetCaller) Guardian(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "guardian")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_AWPWorkNet *AWPWorkNetSession) Guardian() (common.Address, error) {
	return _AWPWorkNet.Contract.Guardian(&_AWPWorkNet.CallOpts)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_AWPWorkNet *AWPWorkNetCallerSession) Guardian() (common.Address, error) {
	return _AWPWorkNet.Contract.Guardian(&_AWPWorkNet.CallOpts)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_AWPWorkNet *AWPWorkNetCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_AWPWorkNet *AWPWorkNetSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _AWPWorkNet.Contract.IsApprovedForAll(&_AWPWorkNet.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_AWPWorkNet *AWPWorkNetCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _AWPWorkNet.Contract.IsApprovedForAll(&_AWPWorkNet.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AWPWorkNet *AWPWorkNetCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AWPWorkNet *AWPWorkNetSession) Name() (string, error) {
	return _AWPWorkNet.Contract.Name(&_AWPWorkNet.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AWPWorkNet *AWPWorkNetCallerSession) Name() (string, error) {
	return _AWPWorkNet.Contract.Name(&_AWPWorkNet.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_AWPWorkNet *AWPWorkNetCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_AWPWorkNet *AWPWorkNetSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _AWPWorkNet.Contract.OwnerOf(&_AWPWorkNet.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_AWPWorkNet *AWPWorkNetCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _AWPWorkNet.Contract.OwnerOf(&_AWPWorkNet.CallOpts, tokenId)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_AWPWorkNet *AWPWorkNetCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_AWPWorkNet *AWPWorkNetSession) ProxiableUUID() ([32]byte, error) {
	return _AWPWorkNet.Contract.ProxiableUUID(&_AWPWorkNet.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_AWPWorkNet *AWPWorkNetCallerSession) ProxiableUUID() ([32]byte, error) {
	return _AWPWorkNet.Contract.ProxiableUUID(&_AWPWorkNet.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AWPWorkNet *AWPWorkNetCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AWPWorkNet *AWPWorkNetSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AWPWorkNet.Contract.SupportsInterface(&_AWPWorkNet.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AWPWorkNet *AWPWorkNetCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AWPWorkNet.Contract.SupportsInterface(&_AWPWorkNet.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AWPWorkNet *AWPWorkNetCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AWPWorkNet *AWPWorkNetSession) Symbol() (string, error) {
	return _AWPWorkNet.Contract.Symbol(&_AWPWorkNet.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AWPWorkNet *AWPWorkNetCallerSession) Symbol() (string, error) {
	return _AWPWorkNet.Contract.Symbol(&_AWPWorkNet.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_AWPWorkNet *AWPWorkNetCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _AWPWorkNet.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_AWPWorkNet *AWPWorkNetSession) TokenURI(tokenId *big.Int) (string, error) {
	return _AWPWorkNet.Contract.TokenURI(&_AWPWorkNet.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_AWPWorkNet *AWPWorkNetCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _AWPWorkNet.Contract.TokenURI(&_AWPWorkNet.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_AWPWorkNet *AWPWorkNetTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AWPWorkNet.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_AWPWorkNet *AWPWorkNetSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.Approve(&_AWPWorkNet.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_AWPWorkNet *AWPWorkNetTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.Approve(&_AWPWorkNet.TransactOpts, to, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_AWPWorkNet *AWPWorkNetTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _AWPWorkNet.contract.Transact(opts, "burn", tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_AWPWorkNet *AWPWorkNetSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.Burn(&_AWPWorkNet.TransactOpts, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_AWPWorkNet *AWPWorkNetTransactorSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.Burn(&_AWPWorkNet.TransactOpts, tokenId)
}

// Initialize is a paid mutator transaction binding the contract method 0x077f224a.
//
// Solidity: function initialize(string name_, string symbol_, address guardian_) returns()
func (_AWPWorkNet *AWPWorkNetTransactor) Initialize(opts *bind.TransactOpts, name_ string, symbol_ string, guardian_ common.Address) (*types.Transaction, error) {
	return _AWPWorkNet.contract.Transact(opts, "initialize", name_, symbol_, guardian_)
}

// Initialize is a paid mutator transaction binding the contract method 0x077f224a.
//
// Solidity: function initialize(string name_, string symbol_, address guardian_) returns()
func (_AWPWorkNet *AWPWorkNetSession) Initialize(name_ string, symbol_ string, guardian_ common.Address) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.Initialize(&_AWPWorkNet.TransactOpts, name_, symbol_, guardian_)
}

// Initialize is a paid mutator transaction binding the contract method 0x077f224a.
//
// Solidity: function initialize(string name_, string symbol_, address guardian_) returns()
func (_AWPWorkNet *AWPWorkNetTransactorSession) Initialize(name_ string, symbol_ string, guardian_ common.Address) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.Initialize(&_AWPWorkNet.TransactOpts, name_, symbol_, guardian_)
}

// Mint is a paid mutator transaction binding the contract method 0x88a13b7c.
//
// Solidity: function mint(address to, uint256 tokenId, string name_, string symbol_, address worknetManager_, address worknetToken_, bytes32 lpPool_, uint128 minStake_, string skillsURI_) returns()
func (_AWPWorkNet *AWPWorkNetTransactor) Mint(opts *bind.TransactOpts, to common.Address, tokenId *big.Int, name_ string, symbol_ string, worknetManager_ common.Address, worknetToken_ common.Address, lpPool_ [32]byte, minStake_ *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _AWPWorkNet.contract.Transact(opts, "mint", to, tokenId, name_, symbol_, worknetManager_, worknetToken_, lpPool_, minStake_, skillsURI_)
}

// Mint is a paid mutator transaction binding the contract method 0x88a13b7c.
//
// Solidity: function mint(address to, uint256 tokenId, string name_, string symbol_, address worknetManager_, address worknetToken_, bytes32 lpPool_, uint128 minStake_, string skillsURI_) returns()
func (_AWPWorkNet *AWPWorkNetSession) Mint(to common.Address, tokenId *big.Int, name_ string, symbol_ string, worknetManager_ common.Address, worknetToken_ common.Address, lpPool_ [32]byte, minStake_ *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.Mint(&_AWPWorkNet.TransactOpts, to, tokenId, name_, symbol_, worknetManager_, worknetToken_, lpPool_, minStake_, skillsURI_)
}

// Mint is a paid mutator transaction binding the contract method 0x88a13b7c.
//
// Solidity: function mint(address to, uint256 tokenId, string name_, string symbol_, address worknetManager_, address worknetToken_, bytes32 lpPool_, uint128 minStake_, string skillsURI_) returns()
func (_AWPWorkNet *AWPWorkNetTransactorSession) Mint(to common.Address, tokenId *big.Int, name_ string, symbol_ string, worknetManager_ common.Address, worknetToken_ common.Address, lpPool_ [32]byte, minStake_ *big.Int, skillsURI_ string) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.Mint(&_AWPWorkNet.TransactOpts, to, tokenId, name_, symbol_, worknetManager_, worknetToken_, lpPool_, minStake_, skillsURI_)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_AWPWorkNet *AWPWorkNetTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AWPWorkNet.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_AWPWorkNet *AWPWorkNetSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SafeTransferFrom(&_AWPWorkNet.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_AWPWorkNet *AWPWorkNetTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SafeTransferFrom(&_AWPWorkNet.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_AWPWorkNet *AWPWorkNetTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _AWPWorkNet.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_AWPWorkNet *AWPWorkNetSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SafeTransferFrom0(&_AWPWorkNet.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_AWPWorkNet *AWPWorkNetTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SafeTransferFrom0(&_AWPWorkNet.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_AWPWorkNet *AWPWorkNetTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _AWPWorkNet.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_AWPWorkNet *AWPWorkNetSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SetApprovalForAll(&_AWPWorkNet.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_AWPWorkNet *AWPWorkNetTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SetApprovalForAll(&_AWPWorkNet.TransactOpts, operator, approved)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string uri) returns()
func (_AWPWorkNet *AWPWorkNetTransactor) SetBaseURI(opts *bind.TransactOpts, uri string) (*types.Transaction, error) {
	return _AWPWorkNet.contract.Transact(opts, "setBaseURI", uri)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string uri) returns()
func (_AWPWorkNet *AWPWorkNetSession) SetBaseURI(uri string) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SetBaseURI(&_AWPWorkNet.TransactOpts, uri)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x55f804b3.
//
// Solidity: function setBaseURI(string uri) returns()
func (_AWPWorkNet *AWPWorkNetTransactorSession) SetBaseURI(uri string) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SetBaseURI(&_AWPWorkNet.TransactOpts, uri)
}

// SetContractURI is a paid mutator transaction binding the contract method 0x938e3d7b.
//
// Solidity: function setContractURI(string uri) returns()
func (_AWPWorkNet *AWPWorkNetTransactor) SetContractURI(opts *bind.TransactOpts, uri string) (*types.Transaction, error) {
	return _AWPWorkNet.contract.Transact(opts, "setContractURI", uri)
}

// SetContractURI is a paid mutator transaction binding the contract method 0x938e3d7b.
//
// Solidity: function setContractURI(string uri) returns()
func (_AWPWorkNet *AWPWorkNetSession) SetContractURI(uri string) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SetContractURI(&_AWPWorkNet.TransactOpts, uri)
}

// SetContractURI is a paid mutator transaction binding the contract method 0x938e3d7b.
//
// Solidity: function setContractURI(string uri) returns()
func (_AWPWorkNet *AWPWorkNetTransactorSession) SetContractURI(uri string) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SetContractURI(&_AWPWorkNet.TransactOpts, uri)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_AWPWorkNet *AWPWorkNetTransactor) SetGuardian(opts *bind.TransactOpts, g common.Address) (*types.Transaction, error) {
	return _AWPWorkNet.contract.Transact(opts, "setGuardian", g)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_AWPWorkNet *AWPWorkNetSession) SetGuardian(g common.Address) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SetGuardian(&_AWPWorkNet.TransactOpts, g)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_AWPWorkNet *AWPWorkNetTransactorSession) SetGuardian(g common.Address) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SetGuardian(&_AWPWorkNet.TransactOpts, g)
}

// SetImageURI is a paid mutator transaction binding the contract method 0x029624e0.
//
// Solidity: function setImageURI(uint256 tokenId, string v) returns()
func (_AWPWorkNet *AWPWorkNetTransactor) SetImageURI(opts *bind.TransactOpts, tokenId *big.Int, v string) (*types.Transaction, error) {
	return _AWPWorkNet.contract.Transact(opts, "setImageURI", tokenId, v)
}

// SetImageURI is a paid mutator transaction binding the contract method 0x029624e0.
//
// Solidity: function setImageURI(uint256 tokenId, string v) returns()
func (_AWPWorkNet *AWPWorkNetSession) SetImageURI(tokenId *big.Int, v string) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SetImageURI(&_AWPWorkNet.TransactOpts, tokenId, v)
}

// SetImageURI is a paid mutator transaction binding the contract method 0x029624e0.
//
// Solidity: function setImageURI(uint256 tokenId, string v) returns()
func (_AWPWorkNet *AWPWorkNetTransactorSession) SetImageURI(tokenId *big.Int, v string) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SetImageURI(&_AWPWorkNet.TransactOpts, tokenId, v)
}

// SetMetadataURI is a paid mutator transaction binding the contract method 0x087dce94.
//
// Solidity: function setMetadataURI(uint256 tokenId, string v) returns()
func (_AWPWorkNet *AWPWorkNetTransactor) SetMetadataURI(opts *bind.TransactOpts, tokenId *big.Int, v string) (*types.Transaction, error) {
	return _AWPWorkNet.contract.Transact(opts, "setMetadataURI", tokenId, v)
}

// SetMetadataURI is a paid mutator transaction binding the contract method 0x087dce94.
//
// Solidity: function setMetadataURI(uint256 tokenId, string v) returns()
func (_AWPWorkNet *AWPWorkNetSession) SetMetadataURI(tokenId *big.Int, v string) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SetMetadataURI(&_AWPWorkNet.TransactOpts, tokenId, v)
}

// SetMetadataURI is a paid mutator transaction binding the contract method 0x087dce94.
//
// Solidity: function setMetadataURI(uint256 tokenId, string v) returns()
func (_AWPWorkNet *AWPWorkNetTransactorSession) SetMetadataURI(tokenId *big.Int, v string) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SetMetadataURI(&_AWPWorkNet.TransactOpts, tokenId, v)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x63a9bbe5.
//
// Solidity: function setMinStake(uint256 tokenId, uint128 v) returns()
func (_AWPWorkNet *AWPWorkNetTransactor) SetMinStake(opts *bind.TransactOpts, tokenId *big.Int, v *big.Int) (*types.Transaction, error) {
	return _AWPWorkNet.contract.Transact(opts, "setMinStake", tokenId, v)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x63a9bbe5.
//
// Solidity: function setMinStake(uint256 tokenId, uint128 v) returns()
func (_AWPWorkNet *AWPWorkNetSession) SetMinStake(tokenId *big.Int, v *big.Int) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SetMinStake(&_AWPWorkNet.TransactOpts, tokenId, v)
}

// SetMinStake is a paid mutator transaction binding the contract method 0x63a9bbe5.
//
// Solidity: function setMinStake(uint256 tokenId, uint128 v) returns()
func (_AWPWorkNet *AWPWorkNetTransactorSession) SetMinStake(tokenId *big.Int, v *big.Int) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SetMinStake(&_AWPWorkNet.TransactOpts, tokenId, v)
}

// SetSkillsURI is a paid mutator transaction binding the contract method 0x7c2f4cd6.
//
// Solidity: function setSkillsURI(uint256 tokenId, string v) returns()
func (_AWPWorkNet *AWPWorkNetTransactor) SetSkillsURI(opts *bind.TransactOpts, tokenId *big.Int, v string) (*types.Transaction, error) {
	return _AWPWorkNet.contract.Transact(opts, "setSkillsURI", tokenId, v)
}

// SetSkillsURI is a paid mutator transaction binding the contract method 0x7c2f4cd6.
//
// Solidity: function setSkillsURI(uint256 tokenId, string v) returns()
func (_AWPWorkNet *AWPWorkNetSession) SetSkillsURI(tokenId *big.Int, v string) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SetSkillsURI(&_AWPWorkNet.TransactOpts, tokenId, v)
}

// SetSkillsURI is a paid mutator transaction binding the contract method 0x7c2f4cd6.
//
// Solidity: function setSkillsURI(uint256 tokenId, string v) returns()
func (_AWPWorkNet *AWPWorkNetTransactorSession) SetSkillsURI(tokenId *big.Int, v string) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.SetSkillsURI(&_AWPWorkNet.TransactOpts, tokenId, v)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_AWPWorkNet *AWPWorkNetTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AWPWorkNet.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_AWPWorkNet *AWPWorkNetSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.TransferFrom(&_AWPWorkNet.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_AWPWorkNet *AWPWorkNetTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.TransferFrom(&_AWPWorkNet.TransactOpts, from, to, tokenId)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AWPWorkNet *AWPWorkNetTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AWPWorkNet.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AWPWorkNet *AWPWorkNetSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.UpgradeToAndCall(&_AWPWorkNet.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AWPWorkNet *AWPWorkNetTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AWPWorkNet.Contract.UpgradeToAndCall(&_AWPWorkNet.TransactOpts, newImplementation, data)
}

// AWPWorkNetApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the AWPWorkNet contract.
type AWPWorkNetApprovalIterator struct {
	Event *AWPWorkNetApproval // Event containing the contract specifics and raw log

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
func (it *AWPWorkNetApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPWorkNetApproval)
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
		it.Event = new(AWPWorkNetApproval)
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
func (it *AWPWorkNetApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPWorkNetApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPWorkNetApproval represents a Approval event raised by the AWPWorkNet contract.
type AWPWorkNetApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_AWPWorkNet *AWPWorkNetFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*AWPWorkNetApprovalIterator, error) {

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

	logs, sub, err := _AWPWorkNet.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &AWPWorkNetApprovalIterator{contract: _AWPWorkNet.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_AWPWorkNet *AWPWorkNetFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *AWPWorkNetApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _AWPWorkNet.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPWorkNetApproval)
				if err := _AWPWorkNet.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_AWPWorkNet *AWPWorkNetFilterer) ParseApproval(log types.Log) (*AWPWorkNetApproval, error) {
	event := new(AWPWorkNetApproval)
	if err := _AWPWorkNet.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPWorkNetApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the AWPWorkNet contract.
type AWPWorkNetApprovalForAllIterator struct {
	Event *AWPWorkNetApprovalForAll // Event containing the contract specifics and raw log

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
func (it *AWPWorkNetApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPWorkNetApprovalForAll)
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
		it.Event = new(AWPWorkNetApprovalForAll)
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
func (it *AWPWorkNetApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPWorkNetApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPWorkNetApprovalForAll represents a ApprovalForAll event raised by the AWPWorkNet contract.
type AWPWorkNetApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_AWPWorkNet *AWPWorkNetFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*AWPWorkNetApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _AWPWorkNet.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &AWPWorkNetApprovalForAllIterator{contract: _AWPWorkNet.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_AWPWorkNet *AWPWorkNetFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *AWPWorkNetApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _AWPWorkNet.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPWorkNetApprovalForAll)
				if err := _AWPWorkNet.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_AWPWorkNet *AWPWorkNetFilterer) ParseApprovalForAll(log types.Log) (*AWPWorkNetApprovalForAll, error) {
	event := new(AWPWorkNetApprovalForAll)
	if err := _AWPWorkNet.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPWorkNetContractURIUpdatedIterator is returned from FilterContractURIUpdated and is used to iterate over the raw logs and unpacked data for ContractURIUpdated events raised by the AWPWorkNet contract.
type AWPWorkNetContractURIUpdatedIterator struct {
	Event *AWPWorkNetContractURIUpdated // Event containing the contract specifics and raw log

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
func (it *AWPWorkNetContractURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPWorkNetContractURIUpdated)
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
		it.Event = new(AWPWorkNetContractURIUpdated)
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
func (it *AWPWorkNetContractURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPWorkNetContractURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPWorkNetContractURIUpdated represents a ContractURIUpdated event raised by the AWPWorkNet contract.
type AWPWorkNetContractURIUpdated struct {
	Uri string
	Raw types.Log // Blockchain specific contextual infos
}

// FilterContractURIUpdated is a free log retrieval operation binding the contract event 0x905d981207a7d0b6c62cc46ab0be2a076d0298e4a86d0ab79882dbd01ac37378.
//
// Solidity: event ContractURIUpdated(string uri)
func (_AWPWorkNet *AWPWorkNetFilterer) FilterContractURIUpdated(opts *bind.FilterOpts) (*AWPWorkNetContractURIUpdatedIterator, error) {

	logs, sub, err := _AWPWorkNet.contract.FilterLogs(opts, "ContractURIUpdated")
	if err != nil {
		return nil, err
	}
	return &AWPWorkNetContractURIUpdatedIterator{contract: _AWPWorkNet.contract, event: "ContractURIUpdated", logs: logs, sub: sub}, nil
}

// WatchContractURIUpdated is a free log subscription operation binding the contract event 0x905d981207a7d0b6c62cc46ab0be2a076d0298e4a86d0ab79882dbd01ac37378.
//
// Solidity: event ContractURIUpdated(string uri)
func (_AWPWorkNet *AWPWorkNetFilterer) WatchContractURIUpdated(opts *bind.WatchOpts, sink chan<- *AWPWorkNetContractURIUpdated) (event.Subscription, error) {

	logs, sub, err := _AWPWorkNet.contract.WatchLogs(opts, "ContractURIUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPWorkNetContractURIUpdated)
				if err := _AWPWorkNet.contract.UnpackLog(event, "ContractURIUpdated", log); err != nil {
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

// ParseContractURIUpdated is a log parse operation binding the contract event 0x905d981207a7d0b6c62cc46ab0be2a076d0298e4a86d0ab79882dbd01ac37378.
//
// Solidity: event ContractURIUpdated(string uri)
func (_AWPWorkNet *AWPWorkNetFilterer) ParseContractURIUpdated(log types.Log) (*AWPWorkNetContractURIUpdated, error) {
	event := new(AWPWorkNetContractURIUpdated)
	if err := _AWPWorkNet.contract.UnpackLog(event, "ContractURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPWorkNetGuardianUpdatedIterator is returned from FilterGuardianUpdated and is used to iterate over the raw logs and unpacked data for GuardianUpdated events raised by the AWPWorkNet contract.
type AWPWorkNetGuardianUpdatedIterator struct {
	Event *AWPWorkNetGuardianUpdated // Event containing the contract specifics and raw log

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
func (it *AWPWorkNetGuardianUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPWorkNetGuardianUpdated)
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
		it.Event = new(AWPWorkNetGuardianUpdated)
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
func (it *AWPWorkNetGuardianUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPWorkNetGuardianUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPWorkNetGuardianUpdated represents a GuardianUpdated event raised by the AWPWorkNet contract.
type AWPWorkNetGuardianUpdated struct {
	NewGuardian common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterGuardianUpdated is a free log retrieval operation binding the contract event 0x6bb7ff33e730289800c62ad882105a144a74010d2bdbb9a942544a3005ad55bf.
//
// Solidity: event GuardianUpdated(address indexed newGuardian)
func (_AWPWorkNet *AWPWorkNetFilterer) FilterGuardianUpdated(opts *bind.FilterOpts, newGuardian []common.Address) (*AWPWorkNetGuardianUpdatedIterator, error) {

	var newGuardianRule []interface{}
	for _, newGuardianItem := range newGuardian {
		newGuardianRule = append(newGuardianRule, newGuardianItem)
	}

	logs, sub, err := _AWPWorkNet.contract.FilterLogs(opts, "GuardianUpdated", newGuardianRule)
	if err != nil {
		return nil, err
	}
	return &AWPWorkNetGuardianUpdatedIterator{contract: _AWPWorkNet.contract, event: "GuardianUpdated", logs: logs, sub: sub}, nil
}

// WatchGuardianUpdated is a free log subscription operation binding the contract event 0x6bb7ff33e730289800c62ad882105a144a74010d2bdbb9a942544a3005ad55bf.
//
// Solidity: event GuardianUpdated(address indexed newGuardian)
func (_AWPWorkNet *AWPWorkNetFilterer) WatchGuardianUpdated(opts *bind.WatchOpts, sink chan<- *AWPWorkNetGuardianUpdated, newGuardian []common.Address) (event.Subscription, error) {

	var newGuardianRule []interface{}
	for _, newGuardianItem := range newGuardian {
		newGuardianRule = append(newGuardianRule, newGuardianItem)
	}

	logs, sub, err := _AWPWorkNet.contract.WatchLogs(opts, "GuardianUpdated", newGuardianRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPWorkNetGuardianUpdated)
				if err := _AWPWorkNet.contract.UnpackLog(event, "GuardianUpdated", log); err != nil {
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
func (_AWPWorkNet *AWPWorkNetFilterer) ParseGuardianUpdated(log types.Log) (*AWPWorkNetGuardianUpdated, error) {
	event := new(AWPWorkNetGuardianUpdated)
	if err := _AWPWorkNet.contract.UnpackLog(event, "GuardianUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPWorkNetImageURIUpdatedIterator is returned from FilterImageURIUpdated and is used to iterate over the raw logs and unpacked data for ImageURIUpdated events raised by the AWPWorkNet contract.
type AWPWorkNetImageURIUpdatedIterator struct {
	Event *AWPWorkNetImageURIUpdated // Event containing the contract specifics and raw log

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
func (it *AWPWorkNetImageURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPWorkNetImageURIUpdated)
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
		it.Event = new(AWPWorkNetImageURIUpdated)
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
func (it *AWPWorkNetImageURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPWorkNetImageURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPWorkNetImageURIUpdated represents a ImageURIUpdated event raised by the AWPWorkNet contract.
type AWPWorkNetImageURIUpdated struct {
	TokenId  *big.Int
	ImageURI string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterImageURIUpdated is a free log retrieval operation binding the contract event 0xc44d776786146efd7a19fd31681be2d271fa4a60fae7eb121f724808ba3d7021.
//
// Solidity: event ImageURIUpdated(uint256 indexed tokenId, string imageURI)
func (_AWPWorkNet *AWPWorkNetFilterer) FilterImageURIUpdated(opts *bind.FilterOpts, tokenId []*big.Int) (*AWPWorkNetImageURIUpdatedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _AWPWorkNet.contract.FilterLogs(opts, "ImageURIUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &AWPWorkNetImageURIUpdatedIterator{contract: _AWPWorkNet.contract, event: "ImageURIUpdated", logs: logs, sub: sub}, nil
}

// WatchImageURIUpdated is a free log subscription operation binding the contract event 0xc44d776786146efd7a19fd31681be2d271fa4a60fae7eb121f724808ba3d7021.
//
// Solidity: event ImageURIUpdated(uint256 indexed tokenId, string imageURI)
func (_AWPWorkNet *AWPWorkNetFilterer) WatchImageURIUpdated(opts *bind.WatchOpts, sink chan<- *AWPWorkNetImageURIUpdated, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _AWPWorkNet.contract.WatchLogs(opts, "ImageURIUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPWorkNetImageURIUpdated)
				if err := _AWPWorkNet.contract.UnpackLog(event, "ImageURIUpdated", log); err != nil {
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

// ParseImageURIUpdated is a log parse operation binding the contract event 0xc44d776786146efd7a19fd31681be2d271fa4a60fae7eb121f724808ba3d7021.
//
// Solidity: event ImageURIUpdated(uint256 indexed tokenId, string imageURI)
func (_AWPWorkNet *AWPWorkNetFilterer) ParseImageURIUpdated(log types.Log) (*AWPWorkNetImageURIUpdated, error) {
	event := new(AWPWorkNetImageURIUpdated)
	if err := _AWPWorkNet.contract.UnpackLog(event, "ImageURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPWorkNetInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the AWPWorkNet contract.
type AWPWorkNetInitializedIterator struct {
	Event *AWPWorkNetInitialized // Event containing the contract specifics and raw log

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
func (it *AWPWorkNetInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPWorkNetInitialized)
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
		it.Event = new(AWPWorkNetInitialized)
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
func (it *AWPWorkNetInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPWorkNetInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPWorkNetInitialized represents a Initialized event raised by the AWPWorkNet contract.
type AWPWorkNetInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AWPWorkNet *AWPWorkNetFilterer) FilterInitialized(opts *bind.FilterOpts) (*AWPWorkNetInitializedIterator, error) {

	logs, sub, err := _AWPWorkNet.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &AWPWorkNetInitializedIterator{contract: _AWPWorkNet.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AWPWorkNet *AWPWorkNetFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *AWPWorkNetInitialized) (event.Subscription, error) {

	logs, sub, err := _AWPWorkNet.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPWorkNetInitialized)
				if err := _AWPWorkNet.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_AWPWorkNet *AWPWorkNetFilterer) ParseInitialized(log types.Log) (*AWPWorkNetInitialized, error) {
	event := new(AWPWorkNetInitialized)
	if err := _AWPWorkNet.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPWorkNetMetadataURIUpdatedIterator is returned from FilterMetadataURIUpdated and is used to iterate over the raw logs and unpacked data for MetadataURIUpdated events raised by the AWPWorkNet contract.
type AWPWorkNetMetadataURIUpdatedIterator struct {
	Event *AWPWorkNetMetadataURIUpdated // Event containing the contract specifics and raw log

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
func (it *AWPWorkNetMetadataURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPWorkNetMetadataURIUpdated)
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
		it.Event = new(AWPWorkNetMetadataURIUpdated)
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
func (it *AWPWorkNetMetadataURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPWorkNetMetadataURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPWorkNetMetadataURIUpdated represents a MetadataURIUpdated event raised by the AWPWorkNet contract.
type AWPWorkNetMetadataURIUpdated struct {
	TokenId     *big.Int
	MetadataURI string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMetadataURIUpdated is a free log retrieval operation binding the contract event 0xbf65482a576bba07ddf407b0dd39c63d560c7765323c11cc051d4a9413881a61.
//
// Solidity: event MetadataURIUpdated(uint256 indexed tokenId, string metadataURI)
func (_AWPWorkNet *AWPWorkNetFilterer) FilterMetadataURIUpdated(opts *bind.FilterOpts, tokenId []*big.Int) (*AWPWorkNetMetadataURIUpdatedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _AWPWorkNet.contract.FilterLogs(opts, "MetadataURIUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &AWPWorkNetMetadataURIUpdatedIterator{contract: _AWPWorkNet.contract, event: "MetadataURIUpdated", logs: logs, sub: sub}, nil
}

// WatchMetadataURIUpdated is a free log subscription operation binding the contract event 0xbf65482a576bba07ddf407b0dd39c63d560c7765323c11cc051d4a9413881a61.
//
// Solidity: event MetadataURIUpdated(uint256 indexed tokenId, string metadataURI)
func (_AWPWorkNet *AWPWorkNetFilterer) WatchMetadataURIUpdated(opts *bind.WatchOpts, sink chan<- *AWPWorkNetMetadataURIUpdated, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _AWPWorkNet.contract.WatchLogs(opts, "MetadataURIUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPWorkNetMetadataURIUpdated)
				if err := _AWPWorkNet.contract.UnpackLog(event, "MetadataURIUpdated", log); err != nil {
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
func (_AWPWorkNet *AWPWorkNetFilterer) ParseMetadataURIUpdated(log types.Log) (*AWPWorkNetMetadataURIUpdated, error) {
	event := new(AWPWorkNetMetadataURIUpdated)
	if err := _AWPWorkNet.contract.UnpackLog(event, "MetadataURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPWorkNetMetadataUpdateIterator is returned from FilterMetadataUpdate and is used to iterate over the raw logs and unpacked data for MetadataUpdate events raised by the AWPWorkNet contract.
type AWPWorkNetMetadataUpdateIterator struct {
	Event *AWPWorkNetMetadataUpdate // Event containing the contract specifics and raw log

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
func (it *AWPWorkNetMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPWorkNetMetadataUpdate)
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
		it.Event = new(AWPWorkNetMetadataUpdate)
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
func (it *AWPWorkNetMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPWorkNetMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPWorkNetMetadataUpdate represents a MetadataUpdate event raised by the AWPWorkNet contract.
type AWPWorkNetMetadataUpdate struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMetadataUpdate is a free log retrieval operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_AWPWorkNet *AWPWorkNetFilterer) FilterMetadataUpdate(opts *bind.FilterOpts) (*AWPWorkNetMetadataUpdateIterator, error) {

	logs, sub, err := _AWPWorkNet.contract.FilterLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &AWPWorkNetMetadataUpdateIterator{contract: _AWPWorkNet.contract, event: "MetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchMetadataUpdate is a free log subscription operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_AWPWorkNet *AWPWorkNetFilterer) WatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *AWPWorkNetMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _AWPWorkNet.contract.WatchLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPWorkNetMetadataUpdate)
				if err := _AWPWorkNet.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
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
func (_AWPWorkNet *AWPWorkNetFilterer) ParseMetadataUpdate(log types.Log) (*AWPWorkNetMetadataUpdate, error) {
	event := new(AWPWorkNetMetadataUpdate)
	if err := _AWPWorkNet.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPWorkNetMinStakeUpdatedIterator is returned from FilterMinStakeUpdated and is used to iterate over the raw logs and unpacked data for MinStakeUpdated events raised by the AWPWorkNet contract.
type AWPWorkNetMinStakeUpdatedIterator struct {
	Event *AWPWorkNetMinStakeUpdated // Event containing the contract specifics and raw log

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
func (it *AWPWorkNetMinStakeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPWorkNetMinStakeUpdated)
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
		it.Event = new(AWPWorkNetMinStakeUpdated)
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
func (it *AWPWorkNetMinStakeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPWorkNetMinStakeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPWorkNetMinStakeUpdated represents a MinStakeUpdated event raised by the AWPWorkNet contract.
type AWPWorkNetMinStakeUpdated struct {
	TokenId  *big.Int
	MinStake *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMinStakeUpdated is a free log retrieval operation binding the contract event 0xd0b53d029b5624436a948f7e5e2d9854defd8058cb4a20ff51a0ff9599ad6de8.
//
// Solidity: event MinStakeUpdated(uint256 indexed tokenId, uint128 minStake)
func (_AWPWorkNet *AWPWorkNetFilterer) FilterMinStakeUpdated(opts *bind.FilterOpts, tokenId []*big.Int) (*AWPWorkNetMinStakeUpdatedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _AWPWorkNet.contract.FilterLogs(opts, "MinStakeUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &AWPWorkNetMinStakeUpdatedIterator{contract: _AWPWorkNet.contract, event: "MinStakeUpdated", logs: logs, sub: sub}, nil
}

// WatchMinStakeUpdated is a free log subscription operation binding the contract event 0xd0b53d029b5624436a948f7e5e2d9854defd8058cb4a20ff51a0ff9599ad6de8.
//
// Solidity: event MinStakeUpdated(uint256 indexed tokenId, uint128 minStake)
func (_AWPWorkNet *AWPWorkNetFilterer) WatchMinStakeUpdated(opts *bind.WatchOpts, sink chan<- *AWPWorkNetMinStakeUpdated, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _AWPWorkNet.contract.WatchLogs(opts, "MinStakeUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPWorkNetMinStakeUpdated)
				if err := _AWPWorkNet.contract.UnpackLog(event, "MinStakeUpdated", log); err != nil {
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
func (_AWPWorkNet *AWPWorkNetFilterer) ParseMinStakeUpdated(log types.Log) (*AWPWorkNetMinStakeUpdated, error) {
	event := new(AWPWorkNetMinStakeUpdated)
	if err := _AWPWorkNet.contract.UnpackLog(event, "MinStakeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPWorkNetSkillsURIUpdatedIterator is returned from FilterSkillsURIUpdated and is used to iterate over the raw logs and unpacked data for SkillsURIUpdated events raised by the AWPWorkNet contract.
type AWPWorkNetSkillsURIUpdatedIterator struct {
	Event *AWPWorkNetSkillsURIUpdated // Event containing the contract specifics and raw log

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
func (it *AWPWorkNetSkillsURIUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPWorkNetSkillsURIUpdated)
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
		it.Event = new(AWPWorkNetSkillsURIUpdated)
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
func (it *AWPWorkNetSkillsURIUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPWorkNetSkillsURIUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPWorkNetSkillsURIUpdated represents a SkillsURIUpdated event raised by the AWPWorkNet contract.
type AWPWorkNetSkillsURIUpdated struct {
	TokenId   *big.Int
	SkillsURI string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSkillsURIUpdated is a free log retrieval operation binding the contract event 0xd1332ed84c54e159e7a4245f8e021aff5b3389b685598c228394168ae30d1020.
//
// Solidity: event SkillsURIUpdated(uint256 indexed tokenId, string skillsURI)
func (_AWPWorkNet *AWPWorkNetFilterer) FilterSkillsURIUpdated(opts *bind.FilterOpts, tokenId []*big.Int) (*AWPWorkNetSkillsURIUpdatedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _AWPWorkNet.contract.FilterLogs(opts, "SkillsURIUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &AWPWorkNetSkillsURIUpdatedIterator{contract: _AWPWorkNet.contract, event: "SkillsURIUpdated", logs: logs, sub: sub}, nil
}

// WatchSkillsURIUpdated is a free log subscription operation binding the contract event 0xd1332ed84c54e159e7a4245f8e021aff5b3389b685598c228394168ae30d1020.
//
// Solidity: event SkillsURIUpdated(uint256 indexed tokenId, string skillsURI)
func (_AWPWorkNet *AWPWorkNetFilterer) WatchSkillsURIUpdated(opts *bind.WatchOpts, sink chan<- *AWPWorkNetSkillsURIUpdated, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _AWPWorkNet.contract.WatchLogs(opts, "SkillsURIUpdated", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPWorkNetSkillsURIUpdated)
				if err := _AWPWorkNet.contract.UnpackLog(event, "SkillsURIUpdated", log); err != nil {
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
func (_AWPWorkNet *AWPWorkNetFilterer) ParseSkillsURIUpdated(log types.Log) (*AWPWorkNetSkillsURIUpdated, error) {
	event := new(AWPWorkNetSkillsURIUpdated)
	if err := _AWPWorkNet.contract.UnpackLog(event, "SkillsURIUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPWorkNetTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the AWPWorkNet contract.
type AWPWorkNetTransferIterator struct {
	Event *AWPWorkNetTransfer // Event containing the contract specifics and raw log

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
func (it *AWPWorkNetTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPWorkNetTransfer)
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
		it.Event = new(AWPWorkNetTransfer)
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
func (it *AWPWorkNetTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPWorkNetTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPWorkNetTransfer represents a Transfer event raised by the AWPWorkNet contract.
type AWPWorkNetTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_AWPWorkNet *AWPWorkNetFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*AWPWorkNetTransferIterator, error) {

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

	logs, sub, err := _AWPWorkNet.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &AWPWorkNetTransferIterator{contract: _AWPWorkNet.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_AWPWorkNet *AWPWorkNetFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *AWPWorkNetTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _AWPWorkNet.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPWorkNetTransfer)
				if err := _AWPWorkNet.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_AWPWorkNet *AWPWorkNetFilterer) ParseTransfer(log types.Log) (*AWPWorkNetTransfer, error) {
	event := new(AWPWorkNetTransfer)
	if err := _AWPWorkNet.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPWorkNetUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the AWPWorkNet contract.
type AWPWorkNetUpgradedIterator struct {
	Event *AWPWorkNetUpgraded // Event containing the contract specifics and raw log

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
func (it *AWPWorkNetUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPWorkNetUpgraded)
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
		it.Event = new(AWPWorkNetUpgraded)
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
func (it *AWPWorkNetUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPWorkNetUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPWorkNetUpgraded represents a Upgraded event raised by the AWPWorkNet contract.
type AWPWorkNetUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AWPWorkNet *AWPWorkNetFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*AWPWorkNetUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AWPWorkNet.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &AWPWorkNetUpgradedIterator{contract: _AWPWorkNet.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AWPWorkNet *AWPWorkNetFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *AWPWorkNetUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AWPWorkNet.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPWorkNetUpgraded)
				if err := _AWPWorkNet.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_AWPWorkNet *AWPWorkNetFilterer) ParseUpgraded(log types.Log) (*AWPWorkNetUpgraded, error) {
	event := new(AWPWorkNetUpgraded)
	if err := _AWPWorkNet.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
