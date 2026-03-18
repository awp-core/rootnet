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

// AccessManagerMetaData contains all meta data concerning the AccessManager contract.
var AccessManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"rootNet_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"agentOwner\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"bind\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"principal\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"oldPrincipal\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAgents\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getOwner\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRewardRecipient\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalUsers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isAgent\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isKnownAddress\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isManager\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isManagerAgent\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRegistered\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRegisteredAgent\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRegisteredUser\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"register\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registeredAt\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeAgent\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resolveCallerRole\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isUser\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isManager_\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rewardRecipients\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rootNet\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setManager\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_isManager\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRewardRecipient\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalUsers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unbind\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"oldPrincipal\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"error\",\"name\":\"AddressIsAgent\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AddressIsPrincipal\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AgentIsSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotRemoveSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotRevokeSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRecipient\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAgentOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotBound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotRootNet\",\"inputs\":[]}]",
	Bin: "0x60a03461009457601f610e6638819003918201601f19168301916001600160401b038311848410176100985780849260209460405283398101031261009457516001600160a01b038116810361009457608052604051610db990816100ad82396080518181816101fb0152818161026c015281816104ea01528181610610015281816106cc0152818161086a015261093d0152f35b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe6080604090808252600480361015610015575f80fd5b5f3560e01c91826301032b86146109fe575081631f5bdf5d1461037e57816321c965de146109a55781633ea4ec341461096c578163405a0b06146109295781634420e486146108435781634ebf59d2146107c55781637fb5d115146107a05781638ce6c37514610767578163902e16451461067c5781639be572f61461048d578163a1656952146105e3578163aaea3d9c146105a2578163afe24ab3146104ab578163bff1f9e11461048d578163c2a8702d146103b9578163c3c5a5471461037e578163c9b1983214610307578163cf5e7bd314610244578163d1ec9f14146101cd57508063e21b38d214610192578063f3ae2415146101565763fa5441611461011d575f80fd5b346101525760203660031901126101525760209061014161013c610a4a565b610bd5565b90516001600160a01b039091168152f35b5f80fd5b5034610152576020366003190112610152576020906001600160a01b0361017b610a4a565b165f526005825260ff815f20541690519015158152f35b5034610152576020366003190112610152576020906001600160a01b03806101b8610a4a565b165f5260038352815f20541615159051908152f35b82346101525780600319360112610152576101e6610a4a565b916101ef610a60565b926001600160a01b03917f000000000000000000000000000000000000000000000000000000000000000083163303610237575060209361022f91610a98565b915191168152f35b83516257aacb60e11b8152fd5b905034610152576020918260031936011261015257610261610a4a565b6001600160a01b03907f0000000000000000000000000000000000000000000000000000000000000000821633036102f957811690815f5260038552825f2054169283156102eb575f8281526003865283812080546001600160a01b031916905560058652838120805460ff191690558481529085528290206102e49190610c3b565b5051908152f35b825163179435b360e01b8152fd5b5050516257aacb60e11b8152fd5b823461015257602036600319011261015257606090610324610a4a565b9060018060a01b0380831692835f525f60205260ff835f20541690815f1461036c57935b5f52600560205260ff835f2054169183519416845215156020840152151590820152f35b50600360205281835f20541693610348565b8234610152576020366003190112610152576020906001600160a01b036103a3610a4a565b165f525f825260ff815f20541690519015158152f35b90503461015257602080600319360112610152576001600160a01b0392836103df610a4a565b165f52828252805f20918151918282855491828152019081955f52835f20905f5b858282106104795750505050839003601f01601f191683019467ffffffffffffffff861184871017610466575084815281855291518185018190529184019291905f5b8281106104505785850386f35b8351871685529381019392810192600101610443565b604190634e487b7160e01b5f525260245ffd5b835485529093019260019283019201610400565b8234610152575f366003190112610152576020906002549051908152f35b8234610152576060366003190112610152576104c5610a4a565b906104ce610a60565b6001600160a01b039290604435848116919082900361015257847f000000000000000000000000000000000000000000000000000000000000000016330361059357841693845f52600360205280845f2054169216809203610583578314610573579261057193835f526003602052825f206bffffffffffffffffffffffff60a01b81541690556005602052825f2060ff1981541690555f526020525f20610c3b565b005b8151635e03d55f60e01b81528490fd5b8251630e41dcbf60e21b81528590fd5b83516257aacb60e11b81528690fd5b8234610152576020366003190112610152576020906001600160a01b036105c7610a4a565b165f526001825267ffffffffffffffff815f2054169051908152f35b82346101525780600319360112610152576105fc610a4a565b610604610a60565b6001600160a01b0391907f00000000000000000000000000000000000000000000000000000000000000008316330361066d57821693841561065f5750165f90815260066020522080546001600160a01b0319169091179055005b8351634e46966960e11b8152fd5b505050516257aacb60e11b8152fd5b823461015257608036600319011261015257610696610a4a565b9161069f610a60565b90604435938415908115809603610152576001600160a01b039360643585811692919083900361015257857f000000000000000000000000000000000000000000000000000000000000000016330361075857851694855f52600360205280875f205416911603610748578161073e575b5061073057505f5260056020525f209060ff801983541691161790555f80f35b825163373d752960e01b8152fd5b9050821485610710565b8451630e41dcbf60e21b81528390fd5b86516257aacb60e11b81528590fd5b8234610152576020366003190112610152576020906001600160a01b038061078d610a4a565b165f5260038352815f2054169051908152f35b8234610152576020366003190112610152576020906001600160a01b0361017b610a4a565b82346101525780600319360112610152576020906107e1610a4a565b6001600160a01b039190826107f4610a60565b1692835f526003855280835f205416911680911492831561081a575b5050519015158152f35b8192935014908161082f575b50908380610810565b90505f525f825260ff815f20541683610826565b9050346101525760203660031901126101525761085e610a4a565b6001600160a01b0391907f00000000000000000000000000000000000000000000000000000000000000008316330361091b57821691825f525f60205260ff845f20541661090c57825f526003602052835f2054166108fe57505f525f602052805f20600160ff1982541617905560016020525f2067ffffffffffffffff421667ffffffffffffffff198254161790556108f9600254610a76565b600255005b825163212dbc5760e11b8152fd5b508251630ea075bf60e21b8152fd5b5082516257aacb60e11b8152fd5b8234610152575f36600319011261015257517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b8234610152576020366003190112610152576020906001600160a01b0380610992610a4a565b165f5260068352815f2054169051908152f35b8234610152576020366003190112610152576020906001600160a01b0390816109cc610a4a565b165f525f835260ff815f2054169182156109ea575b50519015158152f35b600384525f829020541615159150836109e1565b83903461015257602036600319011261015257602091610a1c610a4a565b6001600160a01b038181165f9081526006865292909220548216908115610a435750168152f35b9050168152f35b600435906001600160a01b038216820361015257565b602435906001600160a01b038216820361015257565b5f198114610a845760010190565b634e487b7160e01b5f52601160045260245ffd5b919060018060a01b0380931690815f5260205f815260409160ff835f205416610bc4578516808414610bb357805f526003825285835f205416610ba257805f525f8252825f20805460ff811615610b61575b5050835f526003825285835f20541695818714610b5a5791600491610b3f959493881680610b42575b50845f5260038252835f20816bffffffffffffffffffffffff60a01b8254161790555f52525f20610d16565b50565b5f52828252610b5385855f20610c3b565b505f610b13565b5050505050565b60019060ff191617905560018252825f2067ffffffffffffffff421667ffffffffffffffff19825416179055610b98600254610a76565b6002555f80610aea565b825163212dbc5760e11b8152600490fd5b82516319ad3f4d60e31b8152600490fd5b82516309cba1f360e41b8152600490fd5b6001600160a01b038181165f908152600360205260409020541680610c0d57505f60205260ff60405f205416610c0a57505f90565b90565b905090565b8054821015610c27575f5260205f2001905f90565b634e487b7160e01b5f52603260045260245ffd5b906001820191815f528260205260405f2054908115155f14610d0e575f1991808301818111610a8457825490848201918211610a8457818103610cc3575b50505080548015610caf57820191610c918383610c12565b909182549160031b1b19169055555f526020525f6040812055600190565b634e487b7160e01b5f52603160045260245ffd5b610cf9610cd3610ce39386610c12565b90549060031b1c92839286610c12565b819391549060031b91821b915f19901b19161790565b90555f528460205260405f20555f8080610c79565b505050505f90565b6001810190825f528160205260405f2054155f14610d7c57805468010000000000000000811015610d6857610d55610ce3826001879401855584610c12565b905554915f5260205260405f2055600190565b634e487b7160e01b5f52604160045260245ffd5b5050505f9056fea2646970667358221220fa34922c97779d1b381864f29450f19aee65b8233e35c0b5516d7b5b12edee9764736f6c63430008180033",
}

// AccessManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use AccessManagerMetaData.ABI instead.
var AccessManagerABI = AccessManagerMetaData.ABI

// AccessManagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AccessManagerMetaData.Bin instead.
var AccessManagerBin = AccessManagerMetaData.Bin

// DeployAccessManager deploys a new Ethereum contract, binding an instance of AccessManager to it.
func DeployAccessManager(auth *bind.TransactOpts, backend bind.ContractBackend, rootNet_ common.Address) (common.Address, *types.Transaction, *AccessManager, error) {
	parsed, err := AccessManagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AccessManagerBin), backend, rootNet_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AccessManager{AccessManagerCaller: AccessManagerCaller{contract: contract}, AccessManagerTransactor: AccessManagerTransactor{contract: contract}, AccessManagerFilterer: AccessManagerFilterer{contract: contract}}, nil
}

// AccessManager is an auto generated Go binding around an Ethereum contract.
type AccessManager struct {
	AccessManagerCaller     // Read-only binding to the contract
	AccessManagerTransactor // Write-only binding to the contract
	AccessManagerFilterer   // Log filterer for contract events
}

// AccessManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type AccessManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AccessManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AccessManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AccessManagerSession struct {
	Contract     *AccessManager    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AccessManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AccessManagerCallerSession struct {
	Contract *AccessManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// AccessManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AccessManagerTransactorSession struct {
	Contract     *AccessManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// AccessManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type AccessManagerRaw struct {
	Contract *AccessManager // Generic contract binding to access the raw methods on
}

// AccessManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AccessManagerCallerRaw struct {
	Contract *AccessManagerCaller // Generic read-only contract binding to access the raw methods on
}

// AccessManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AccessManagerTransactorRaw struct {
	Contract *AccessManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAccessManager creates a new instance of AccessManager, bound to a specific deployed contract.
func NewAccessManager(address common.Address, backend bind.ContractBackend) (*AccessManager, error) {
	contract, err := bindAccessManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccessManager{AccessManagerCaller: AccessManagerCaller{contract: contract}, AccessManagerTransactor: AccessManagerTransactor{contract: contract}, AccessManagerFilterer: AccessManagerFilterer{contract: contract}}, nil
}

// NewAccessManagerCaller creates a new read-only instance of AccessManager, bound to a specific deployed contract.
func NewAccessManagerCaller(address common.Address, caller bind.ContractCaller) (*AccessManagerCaller, error) {
	contract, err := bindAccessManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccessManagerCaller{contract: contract}, nil
}

// NewAccessManagerTransactor creates a new write-only instance of AccessManager, bound to a specific deployed contract.
func NewAccessManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*AccessManagerTransactor, error) {
	contract, err := bindAccessManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccessManagerTransactor{contract: contract}, nil
}

// NewAccessManagerFilterer creates a new log filterer instance of AccessManager, bound to a specific deployed contract.
func NewAccessManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*AccessManagerFilterer, error) {
	contract, err := bindAccessManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccessManagerFilterer{contract: contract}, nil
}

// bindAccessManager binds a generic wrapper to an already deployed contract.
func bindAccessManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AccessManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccessManager *AccessManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessManager.Contract.AccessManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccessManager *AccessManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessManager.Contract.AccessManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccessManager *AccessManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessManager.Contract.AccessManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccessManager *AccessManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccessManager *AccessManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccessManager *AccessManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessManager.Contract.contract.Transact(opts, method, params...)
}

// AgentOwner is a free data retrieval call binding the contract method 0x8ce6c375.
//
// Solidity: function agentOwner(address ) view returns(address)
func (_AccessManager *AccessManagerCaller) AgentOwner(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _AccessManager.contract.Call(opts, &out, "agentOwner", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AgentOwner is a free data retrieval call binding the contract method 0x8ce6c375.
//
// Solidity: function agentOwner(address ) view returns(address)
func (_AccessManager *AccessManagerSession) AgentOwner(arg0 common.Address) (common.Address, error) {
	return _AccessManager.Contract.AgentOwner(&_AccessManager.CallOpts, arg0)
}

// AgentOwner is a free data retrieval call binding the contract method 0x8ce6c375.
//
// Solidity: function agentOwner(address ) view returns(address)
func (_AccessManager *AccessManagerCallerSession) AgentOwner(arg0 common.Address) (common.Address, error) {
	return _AccessManager.Contract.AgentOwner(&_AccessManager.CallOpts, arg0)
}

// GetAgents is a free data retrieval call binding the contract method 0xc2a8702d.
//
// Solidity: function getAgents(address user) view returns(address[])
func (_AccessManager *AccessManagerCaller) GetAgents(opts *bind.CallOpts, user common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _AccessManager.contract.Call(opts, &out, "getAgents", user)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetAgents is a free data retrieval call binding the contract method 0xc2a8702d.
//
// Solidity: function getAgents(address user) view returns(address[])
func (_AccessManager *AccessManagerSession) GetAgents(user common.Address) ([]common.Address, error) {
	return _AccessManager.Contract.GetAgents(&_AccessManager.CallOpts, user)
}

// GetAgents is a free data retrieval call binding the contract method 0xc2a8702d.
//
// Solidity: function getAgents(address user) view returns(address[])
func (_AccessManager *AccessManagerCallerSession) GetAgents(user common.Address) ([]common.Address, error) {
	return _AccessManager.Contract.GetAgents(&_AccessManager.CallOpts, user)
}

// GetOwner is a free data retrieval call binding the contract method 0xfa544161.
//
// Solidity: function getOwner(address addr) view returns(address)
func (_AccessManager *AccessManagerCaller) GetOwner(opts *bind.CallOpts, addr common.Address) (common.Address, error) {
	var out []interface{}
	err := _AccessManager.contract.Call(opts, &out, "getOwner", addr)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOwner is a free data retrieval call binding the contract method 0xfa544161.
//
// Solidity: function getOwner(address addr) view returns(address)
func (_AccessManager *AccessManagerSession) GetOwner(addr common.Address) (common.Address, error) {
	return _AccessManager.Contract.GetOwner(&_AccessManager.CallOpts, addr)
}

// GetOwner is a free data retrieval call binding the contract method 0xfa544161.
//
// Solidity: function getOwner(address addr) view returns(address)
func (_AccessManager *AccessManagerCallerSession) GetOwner(addr common.Address) (common.Address, error) {
	return _AccessManager.Contract.GetOwner(&_AccessManager.CallOpts, addr)
}

// GetRewardRecipient is a free data retrieval call binding the contract method 0x01032b86.
//
// Solidity: function getRewardRecipient(address user) view returns(address)
func (_AccessManager *AccessManagerCaller) GetRewardRecipient(opts *bind.CallOpts, user common.Address) (common.Address, error) {
	var out []interface{}
	err := _AccessManager.contract.Call(opts, &out, "getRewardRecipient", user)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRewardRecipient is a free data retrieval call binding the contract method 0x01032b86.
//
// Solidity: function getRewardRecipient(address user) view returns(address)
func (_AccessManager *AccessManagerSession) GetRewardRecipient(user common.Address) (common.Address, error) {
	return _AccessManager.Contract.GetRewardRecipient(&_AccessManager.CallOpts, user)
}

// GetRewardRecipient is a free data retrieval call binding the contract method 0x01032b86.
//
// Solidity: function getRewardRecipient(address user) view returns(address)
func (_AccessManager *AccessManagerCallerSession) GetRewardRecipient(user common.Address) (common.Address, error) {
	return _AccessManager.Contract.GetRewardRecipient(&_AccessManager.CallOpts, user)
}

// GetTotalUsers is a free data retrieval call binding the contract method 0x9be572f6.
//
// Solidity: function getTotalUsers() view returns(uint256)
func (_AccessManager *AccessManagerCaller) GetTotalUsers(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessManager.contract.Call(opts, &out, "getTotalUsers")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalUsers is a free data retrieval call binding the contract method 0x9be572f6.
//
// Solidity: function getTotalUsers() view returns(uint256)
func (_AccessManager *AccessManagerSession) GetTotalUsers() (*big.Int, error) {
	return _AccessManager.Contract.GetTotalUsers(&_AccessManager.CallOpts)
}

// GetTotalUsers is a free data retrieval call binding the contract method 0x9be572f6.
//
// Solidity: function getTotalUsers() view returns(uint256)
func (_AccessManager *AccessManagerCallerSession) GetTotalUsers() (*big.Int, error) {
	return _AccessManager.Contract.GetTotalUsers(&_AccessManager.CallOpts)
}

// IsAgent is a free data retrieval call binding the contract method 0x4ebf59d2.
//
// Solidity: function isAgent(address user, address agent) view returns(bool)
func (_AccessManager *AccessManagerCaller) IsAgent(opts *bind.CallOpts, user common.Address, agent common.Address) (bool, error) {
	var out []interface{}
	err := _AccessManager.contract.Call(opts, &out, "isAgent", user, agent)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAgent is a free data retrieval call binding the contract method 0x4ebf59d2.
//
// Solidity: function isAgent(address user, address agent) view returns(bool)
func (_AccessManager *AccessManagerSession) IsAgent(user common.Address, agent common.Address) (bool, error) {
	return _AccessManager.Contract.IsAgent(&_AccessManager.CallOpts, user, agent)
}

// IsAgent is a free data retrieval call binding the contract method 0x4ebf59d2.
//
// Solidity: function isAgent(address user, address agent) view returns(bool)
func (_AccessManager *AccessManagerCallerSession) IsAgent(user common.Address, agent common.Address) (bool, error) {
	return _AccessManager.Contract.IsAgent(&_AccessManager.CallOpts, user, agent)
}

// IsKnownAddress is a free data retrieval call binding the contract method 0x21c965de.
//
// Solidity: function isKnownAddress(address addr) view returns(bool)
func (_AccessManager *AccessManagerCaller) IsKnownAddress(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _AccessManager.contract.Call(opts, &out, "isKnownAddress", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsKnownAddress is a free data retrieval call binding the contract method 0x21c965de.
//
// Solidity: function isKnownAddress(address addr) view returns(bool)
func (_AccessManager *AccessManagerSession) IsKnownAddress(addr common.Address) (bool, error) {
	return _AccessManager.Contract.IsKnownAddress(&_AccessManager.CallOpts, addr)
}

// IsKnownAddress is a free data retrieval call binding the contract method 0x21c965de.
//
// Solidity: function isKnownAddress(address addr) view returns(bool)
func (_AccessManager *AccessManagerCallerSession) IsKnownAddress(addr common.Address) (bool, error) {
	return _AccessManager.Contract.IsKnownAddress(&_AccessManager.CallOpts, addr)
}

// IsManager is a free data retrieval call binding the contract method 0xf3ae2415.
//
// Solidity: function isManager(address ) view returns(bool)
func (_AccessManager *AccessManagerCaller) IsManager(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _AccessManager.contract.Call(opts, &out, "isManager", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsManager is a free data retrieval call binding the contract method 0xf3ae2415.
//
// Solidity: function isManager(address ) view returns(bool)
func (_AccessManager *AccessManagerSession) IsManager(arg0 common.Address) (bool, error) {
	return _AccessManager.Contract.IsManager(&_AccessManager.CallOpts, arg0)
}

// IsManager is a free data retrieval call binding the contract method 0xf3ae2415.
//
// Solidity: function isManager(address ) view returns(bool)
func (_AccessManager *AccessManagerCallerSession) IsManager(arg0 common.Address) (bool, error) {
	return _AccessManager.Contract.IsManager(&_AccessManager.CallOpts, arg0)
}

// IsManagerAgent is a free data retrieval call binding the contract method 0x7fb5d115.
//
// Solidity: function isManagerAgent(address agent) view returns(bool)
func (_AccessManager *AccessManagerCaller) IsManagerAgent(opts *bind.CallOpts, agent common.Address) (bool, error) {
	var out []interface{}
	err := _AccessManager.contract.Call(opts, &out, "isManagerAgent", agent)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsManagerAgent is a free data retrieval call binding the contract method 0x7fb5d115.
//
// Solidity: function isManagerAgent(address agent) view returns(bool)
func (_AccessManager *AccessManagerSession) IsManagerAgent(agent common.Address) (bool, error) {
	return _AccessManager.Contract.IsManagerAgent(&_AccessManager.CallOpts, agent)
}

// IsManagerAgent is a free data retrieval call binding the contract method 0x7fb5d115.
//
// Solidity: function isManagerAgent(address agent) view returns(bool)
func (_AccessManager *AccessManagerCallerSession) IsManagerAgent(agent common.Address) (bool, error) {
	return _AccessManager.Contract.IsManagerAgent(&_AccessManager.CallOpts, agent)
}

// IsRegistered is a free data retrieval call binding the contract method 0xc3c5a547.
//
// Solidity: function isRegistered(address ) view returns(bool)
func (_AccessManager *AccessManagerCaller) IsRegistered(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _AccessManager.contract.Call(opts, &out, "isRegistered", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsRegistered is a free data retrieval call binding the contract method 0xc3c5a547.
//
// Solidity: function isRegistered(address ) view returns(bool)
func (_AccessManager *AccessManagerSession) IsRegistered(arg0 common.Address) (bool, error) {
	return _AccessManager.Contract.IsRegistered(&_AccessManager.CallOpts, arg0)
}

// IsRegistered is a free data retrieval call binding the contract method 0xc3c5a547.
//
// Solidity: function isRegistered(address ) view returns(bool)
func (_AccessManager *AccessManagerCallerSession) IsRegistered(arg0 common.Address) (bool, error) {
	return _AccessManager.Contract.IsRegistered(&_AccessManager.CallOpts, arg0)
}

// IsRegisteredAgent is a free data retrieval call binding the contract method 0xe21b38d2.
//
// Solidity: function isRegisteredAgent(address addr) view returns(bool)
func (_AccessManager *AccessManagerCaller) IsRegisteredAgent(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _AccessManager.contract.Call(opts, &out, "isRegisteredAgent", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsRegisteredAgent is a free data retrieval call binding the contract method 0xe21b38d2.
//
// Solidity: function isRegisteredAgent(address addr) view returns(bool)
func (_AccessManager *AccessManagerSession) IsRegisteredAgent(addr common.Address) (bool, error) {
	return _AccessManager.Contract.IsRegisteredAgent(&_AccessManager.CallOpts, addr)
}

// IsRegisteredAgent is a free data retrieval call binding the contract method 0xe21b38d2.
//
// Solidity: function isRegisteredAgent(address addr) view returns(bool)
func (_AccessManager *AccessManagerCallerSession) IsRegisteredAgent(addr common.Address) (bool, error) {
	return _AccessManager.Contract.IsRegisteredAgent(&_AccessManager.CallOpts, addr)
}

// IsRegisteredUser is a free data retrieval call binding the contract method 0x1f5bdf5d.
//
// Solidity: function isRegisteredUser(address addr) view returns(bool)
func (_AccessManager *AccessManagerCaller) IsRegisteredUser(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _AccessManager.contract.Call(opts, &out, "isRegisteredUser", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsRegisteredUser is a free data retrieval call binding the contract method 0x1f5bdf5d.
//
// Solidity: function isRegisteredUser(address addr) view returns(bool)
func (_AccessManager *AccessManagerSession) IsRegisteredUser(addr common.Address) (bool, error) {
	return _AccessManager.Contract.IsRegisteredUser(&_AccessManager.CallOpts, addr)
}

// IsRegisteredUser is a free data retrieval call binding the contract method 0x1f5bdf5d.
//
// Solidity: function isRegisteredUser(address addr) view returns(bool)
func (_AccessManager *AccessManagerCallerSession) IsRegisteredUser(addr common.Address) (bool, error) {
	return _AccessManager.Contract.IsRegisteredUser(&_AccessManager.CallOpts, addr)
}

// RegisteredAt is a free data retrieval call binding the contract method 0xaaea3d9c.
//
// Solidity: function registeredAt(address ) view returns(uint64)
func (_AccessManager *AccessManagerCaller) RegisteredAt(opts *bind.CallOpts, arg0 common.Address) (uint64, error) {
	var out []interface{}
	err := _AccessManager.contract.Call(opts, &out, "registeredAt", arg0)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// RegisteredAt is a free data retrieval call binding the contract method 0xaaea3d9c.
//
// Solidity: function registeredAt(address ) view returns(uint64)
func (_AccessManager *AccessManagerSession) RegisteredAt(arg0 common.Address) (uint64, error) {
	return _AccessManager.Contract.RegisteredAt(&_AccessManager.CallOpts, arg0)
}

// RegisteredAt is a free data retrieval call binding the contract method 0xaaea3d9c.
//
// Solidity: function registeredAt(address ) view returns(uint64)
func (_AccessManager *AccessManagerCallerSession) RegisteredAt(arg0 common.Address) (uint64, error) {
	return _AccessManager.Contract.RegisteredAt(&_AccessManager.CallOpts, arg0)
}

// ResolveCallerRole is a free data retrieval call binding the contract method 0xc9b19832.
//
// Solidity: function resolveCallerRole(address addr) view returns(address owner, bool isUser, bool isManager_)
func (_AccessManager *AccessManagerCaller) ResolveCallerRole(opts *bind.CallOpts, addr common.Address) (struct {
	Owner     common.Address
	IsUser    bool
	IsManager bool
}, error) {
	var out []interface{}
	err := _AccessManager.contract.Call(opts, &out, "resolveCallerRole", addr)

	outstruct := new(struct {
		Owner     common.Address
		IsUser    bool
		IsManager bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Owner = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.IsUser = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.IsManager = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

// ResolveCallerRole is a free data retrieval call binding the contract method 0xc9b19832.
//
// Solidity: function resolveCallerRole(address addr) view returns(address owner, bool isUser, bool isManager_)
func (_AccessManager *AccessManagerSession) ResolveCallerRole(addr common.Address) (struct {
	Owner     common.Address
	IsUser    bool
	IsManager bool
}, error) {
	return _AccessManager.Contract.ResolveCallerRole(&_AccessManager.CallOpts, addr)
}

// ResolveCallerRole is a free data retrieval call binding the contract method 0xc9b19832.
//
// Solidity: function resolveCallerRole(address addr) view returns(address owner, bool isUser, bool isManager_)
func (_AccessManager *AccessManagerCallerSession) ResolveCallerRole(addr common.Address) (struct {
	Owner     common.Address
	IsUser    bool
	IsManager bool
}, error) {
	return _AccessManager.Contract.ResolveCallerRole(&_AccessManager.CallOpts, addr)
}

// RewardRecipients is a free data retrieval call binding the contract method 0x3ea4ec34.
//
// Solidity: function rewardRecipients(address ) view returns(address)
func (_AccessManager *AccessManagerCaller) RewardRecipients(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _AccessManager.contract.Call(opts, &out, "rewardRecipients", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RewardRecipients is a free data retrieval call binding the contract method 0x3ea4ec34.
//
// Solidity: function rewardRecipients(address ) view returns(address)
func (_AccessManager *AccessManagerSession) RewardRecipients(arg0 common.Address) (common.Address, error) {
	return _AccessManager.Contract.RewardRecipients(&_AccessManager.CallOpts, arg0)
}

// RewardRecipients is a free data retrieval call binding the contract method 0x3ea4ec34.
//
// Solidity: function rewardRecipients(address ) view returns(address)
func (_AccessManager *AccessManagerCallerSession) RewardRecipients(arg0 common.Address) (common.Address, error) {
	return _AccessManager.Contract.RewardRecipients(&_AccessManager.CallOpts, arg0)
}

// RootNet is a free data retrieval call binding the contract method 0x405a0b06.
//
// Solidity: function rootNet() view returns(address)
func (_AccessManager *AccessManagerCaller) RootNet(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessManager.contract.Call(opts, &out, "rootNet")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RootNet is a free data retrieval call binding the contract method 0x405a0b06.
//
// Solidity: function rootNet() view returns(address)
func (_AccessManager *AccessManagerSession) RootNet() (common.Address, error) {
	return _AccessManager.Contract.RootNet(&_AccessManager.CallOpts)
}

// RootNet is a free data retrieval call binding the contract method 0x405a0b06.
//
// Solidity: function rootNet() view returns(address)
func (_AccessManager *AccessManagerCallerSession) RootNet() (common.Address, error) {
	return _AccessManager.Contract.RootNet(&_AccessManager.CallOpts)
}

// TotalUsers is a free data retrieval call binding the contract method 0xbff1f9e1.
//
// Solidity: function totalUsers() view returns(uint256)
func (_AccessManager *AccessManagerCaller) TotalUsers(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccessManager.contract.Call(opts, &out, "totalUsers")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalUsers is a free data retrieval call binding the contract method 0xbff1f9e1.
//
// Solidity: function totalUsers() view returns(uint256)
func (_AccessManager *AccessManagerSession) TotalUsers() (*big.Int, error) {
	return _AccessManager.Contract.TotalUsers(&_AccessManager.CallOpts)
}

// TotalUsers is a free data retrieval call binding the contract method 0xbff1f9e1.
//
// Solidity: function totalUsers() view returns(uint256)
func (_AccessManager *AccessManagerCallerSession) TotalUsers() (*big.Int, error) {
	return _AccessManager.Contract.TotalUsers(&_AccessManager.CallOpts)
}

// Bind is a paid mutator transaction binding the contract method 0xd1ec9f14.
//
// Solidity: function bind(address agent, address principal) returns(address oldPrincipal)
func (_AccessManager *AccessManagerTransactor) Bind(opts *bind.TransactOpts, agent common.Address, principal common.Address) (*types.Transaction, error) {
	return _AccessManager.contract.Transact(opts, "bind", agent, principal)
}

// Bind is a paid mutator transaction binding the contract method 0xd1ec9f14.
//
// Solidity: function bind(address agent, address principal) returns(address oldPrincipal)
func (_AccessManager *AccessManagerSession) Bind(agent common.Address, principal common.Address) (*types.Transaction, error) {
	return _AccessManager.Contract.Bind(&_AccessManager.TransactOpts, agent, principal)
}

// Bind is a paid mutator transaction binding the contract method 0xd1ec9f14.
//
// Solidity: function bind(address agent, address principal) returns(address oldPrincipal)
func (_AccessManager *AccessManagerTransactorSession) Bind(agent common.Address, principal common.Address) (*types.Transaction, error) {
	return _AccessManager.Contract.Bind(&_AccessManager.TransactOpts, agent, principal)
}

// Register is a paid mutator transaction binding the contract method 0x4420e486.
//
// Solidity: function register(address user) returns()
func (_AccessManager *AccessManagerTransactor) Register(opts *bind.TransactOpts, user common.Address) (*types.Transaction, error) {
	return _AccessManager.contract.Transact(opts, "register", user)
}

// Register is a paid mutator transaction binding the contract method 0x4420e486.
//
// Solidity: function register(address user) returns()
func (_AccessManager *AccessManagerSession) Register(user common.Address) (*types.Transaction, error) {
	return _AccessManager.Contract.Register(&_AccessManager.TransactOpts, user)
}

// Register is a paid mutator transaction binding the contract method 0x4420e486.
//
// Solidity: function register(address user) returns()
func (_AccessManager *AccessManagerTransactorSession) Register(user common.Address) (*types.Transaction, error) {
	return _AccessManager.Contract.Register(&_AccessManager.TransactOpts, user)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0xafe24ab3.
//
// Solidity: function removeAgent(address user, address agent, address operator) returns()
func (_AccessManager *AccessManagerTransactor) RemoveAgent(opts *bind.TransactOpts, user common.Address, agent common.Address, operator common.Address) (*types.Transaction, error) {
	return _AccessManager.contract.Transact(opts, "removeAgent", user, agent, operator)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0xafe24ab3.
//
// Solidity: function removeAgent(address user, address agent, address operator) returns()
func (_AccessManager *AccessManagerSession) RemoveAgent(user common.Address, agent common.Address, operator common.Address) (*types.Transaction, error) {
	return _AccessManager.Contract.RemoveAgent(&_AccessManager.TransactOpts, user, agent, operator)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0xafe24ab3.
//
// Solidity: function removeAgent(address user, address agent, address operator) returns()
func (_AccessManager *AccessManagerTransactorSession) RemoveAgent(user common.Address, agent common.Address, operator common.Address) (*types.Transaction, error) {
	return _AccessManager.Contract.RemoveAgent(&_AccessManager.TransactOpts, user, agent, operator)
}

// SetManager is a paid mutator transaction binding the contract method 0x902e1645.
//
// Solidity: function setManager(address user, address agent, bool _isManager, address operator) returns()
func (_AccessManager *AccessManagerTransactor) SetManager(opts *bind.TransactOpts, user common.Address, agent common.Address, _isManager bool, operator common.Address) (*types.Transaction, error) {
	return _AccessManager.contract.Transact(opts, "setManager", user, agent, _isManager, operator)
}

// SetManager is a paid mutator transaction binding the contract method 0x902e1645.
//
// Solidity: function setManager(address user, address agent, bool _isManager, address operator) returns()
func (_AccessManager *AccessManagerSession) SetManager(user common.Address, agent common.Address, _isManager bool, operator common.Address) (*types.Transaction, error) {
	return _AccessManager.Contract.SetManager(&_AccessManager.TransactOpts, user, agent, _isManager, operator)
}

// SetManager is a paid mutator transaction binding the contract method 0x902e1645.
//
// Solidity: function setManager(address user, address agent, bool _isManager, address operator) returns()
func (_AccessManager *AccessManagerTransactorSession) SetManager(user common.Address, agent common.Address, _isManager bool, operator common.Address) (*types.Transaction, error) {
	return _AccessManager.Contract.SetManager(&_AccessManager.TransactOpts, user, agent, _isManager, operator)
}

// SetRewardRecipient is a paid mutator transaction binding the contract method 0xa1656952.
//
// Solidity: function setRewardRecipient(address user, address recipient) returns()
func (_AccessManager *AccessManagerTransactor) SetRewardRecipient(opts *bind.TransactOpts, user common.Address, recipient common.Address) (*types.Transaction, error) {
	return _AccessManager.contract.Transact(opts, "setRewardRecipient", user, recipient)
}

// SetRewardRecipient is a paid mutator transaction binding the contract method 0xa1656952.
//
// Solidity: function setRewardRecipient(address user, address recipient) returns()
func (_AccessManager *AccessManagerSession) SetRewardRecipient(user common.Address, recipient common.Address) (*types.Transaction, error) {
	return _AccessManager.Contract.SetRewardRecipient(&_AccessManager.TransactOpts, user, recipient)
}

// SetRewardRecipient is a paid mutator transaction binding the contract method 0xa1656952.
//
// Solidity: function setRewardRecipient(address user, address recipient) returns()
func (_AccessManager *AccessManagerTransactorSession) SetRewardRecipient(user common.Address, recipient common.Address) (*types.Transaction, error) {
	return _AccessManager.Contract.SetRewardRecipient(&_AccessManager.TransactOpts, user, recipient)
}

// Unbind is a paid mutator transaction binding the contract method 0xcf5e7bd3.
//
// Solidity: function unbind(address agent) returns(address oldPrincipal)
func (_AccessManager *AccessManagerTransactor) Unbind(opts *bind.TransactOpts, agent common.Address) (*types.Transaction, error) {
	return _AccessManager.contract.Transact(opts, "unbind", agent)
}

// Unbind is a paid mutator transaction binding the contract method 0xcf5e7bd3.
//
// Solidity: function unbind(address agent) returns(address oldPrincipal)
func (_AccessManager *AccessManagerSession) Unbind(agent common.Address) (*types.Transaction, error) {
	return _AccessManager.Contract.Unbind(&_AccessManager.TransactOpts, agent)
}

// Unbind is a paid mutator transaction binding the contract method 0xcf5e7bd3.
//
// Solidity: function unbind(address agent) returns(address oldPrincipal)
func (_AccessManager *AccessManagerTransactorSession) Unbind(agent common.Address) (*types.Transaction, error) {
	return _AccessManager.Contract.Unbind(&_AccessManager.TransactOpts, agent)
}
