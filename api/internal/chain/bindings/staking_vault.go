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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"rootNet_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allocate\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deallocate\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"freezeAgentAllocations\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAgentStake\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAgentSubnets\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSubnetTotalStake\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reallocate\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromAgent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromSubnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"toAgent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"toSubnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rootNet\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setStakeNFT\",\"inputs\":[{\"name\":\"stakeNFT_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakeNFT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"subnetTotalStake\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"userTotalAllocated\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AgentAllocationsFrozen\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"totalFrozen\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadySet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientAllocation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientUnallocated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotRootNet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
	Bin: "0x60a03461008d57601f610d0938819003918201601f19168301916001600160401b038311848410176100915780849260209460405283398101031261008d57516001600160a01b038116810361008d57608052604051610c6390816100a68239608051818181610172015281816103300152818161050c0152818161065e015281816106ca01526107530152f35b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe6080604081815260049182361015610015575f80fd5b5f3560e01c9081630358ccb5146107df575080631ed2facb146107b957806332ffa4ce14610782578063405a0b061461073f57806348f069ec146106a45780636f9808131461063a578063716fb83d146104ec578063b48509e6146104c5578063d035a9a714610313578063d5d5278d14610129578063f1ad80c6146100cc5763f1e18405146100a3575f80fd5b346100c85760203660031901126100c85781602092355f528252805f20549051908152f35b5f80fd5b50346100c85760603660031901126100c8576020906100e961089a565b6100f16108b0565b9060018060a01b038091165f5260018452825f2091165f528252805f206044355f5282526001600160801b03815f2054169051908152f35b50346100c85760c03660031901126100c85761014361089a565b9061014c6108b0565b926001600160a01b0360643581811690604435908290036100c8576084359160a43596847f000000000000000000000000000000000000000000000000000000000000000016330361030457871580156102f4575b80156102ec575b80156102e4575b6102d4576001600160801b0385818a16921691825f5260209660018852895f209b169a8b5f528752885f20855f52875281895f2054169a818c106102c4576101fb8261029e9c9d610a66565b90845f52600189528a5f20815f5289528a5f20875f528952838b5f2092166001600160801b0319928184825416179055156102a2575b50835f5260018852895f20855f528852895f20875f528852895f209261025b845493828516610ab5565b1691161790555f5260028452855f20905f52835261027b82865f20610bd4565b505f52828252835f2061028f8782546108ff565b90555f52525f20918254610920565b9055005b845f52600289528a5f20905f5288526102bd868b5f20610af9565b505f610231565b89516337cb51dd60e21b81528990fd5b865163162908e360e11b81528690fd5b5083156101af565b5082156101a8565b506001600160801b0388116101a1565b86516257aacb60e11b81528690fd5b50346100c857610322366108c6565b9093926001600160a01b03917f0000000000000000000000000000000000000000000000000000000000000000831633036104b657801580156104a6575b801561049e575b61048e5782805f54169286519485916337178b4b60e11b83521693848a83015281602460209788935afa908115610484575f91610457575b50835f52600385526103b483885f2054610920565b116104475796610445976001600160801b0390845f5260018652875f20961695865f528552865f20885f528552865f20908154906103f6818616828416610ab5565b16906001600160801b031916179055825f5260038452855f2061041a838254610920565b9055865f52835261042f855f20918254610920565b90555f5260028152825f20915f52525f20610bd4565b005b855163d247d12160e01b81528890fd5b90508481813d831161047d575b61046e8183610a7f565b810103126100c857515f61039f565b503d610464565b87513d5f823e3d90fd5b845163162908e360e11b81528790fd5b508515610367565b506001600160801b038111610360565b84516257aacb60e11b81528790fd5b50346100c8575f3660031901126100c8575f5490516001600160a01b039091168152602090f35b50346100c8576104fb366108c6565b92949093926001600160a01b0391907f00000000000000000000000000000000000000000000000000000000000000008316330361062c578015801561061c575b61060d576001600160801b039183838316981697885f5260209460018652875f20961695865f528552865f20885f52855283875f2054168181106105fd579061058491610a66565b92885f5260018552865f20865f528552865f20885f528552865f20931692836001600160801b0319825416179055875f5260038452855f206105c78382546108ff565b9055865f5283526105dc855f209182546108ff565b9055156105e557005b610445945f5260028152825f20915f52525f20610af9565b87516337cb51dd60e21b81528390fd5b50835163162908e360e11b8152fd5b506001600160801b03811161053c565b5083516257aacb60e11b8152fd5b50346100c857806003193601126100c85761065361089a565b61065b6108b0565b917f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316330361069657610445838361092d565b516257aacb60e11b81528390fd5b5090346100c85760203660031901126100c8576106bf61089a565b906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000081163303610731575f549281841661072157169283156107145750506001600160a01b031916175f55005b5163d92e233d60e01b8152fd5b845163a741a04560e01b81528390fd5b5082516257aacb60e11b8152fd5b50346100c8575f3660031901126100c857517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b50346100c85760203660031901126100c8576020906001600160a01b036107a761089a565b165f5260038252805f20549051908152f35b50346100c85760203660031901126100c85781602092355f528252805f20549051908152f35b9050346100c857816003193601126100c857906107fa61089a565b6108026108b0565b60018060a01b038092165f5260209160028352835f2091165f528152815f209081548085528185019081935f52825f20905f5b818110610886575050508461084b910385610a7f565b825181815293518185018190528493840192915f5b82811061086f57505050500390f35b835185528695509381019392810192600101610860565b825484529284019260019283019201610835565b600435906001600160a01b03821682036100c857565b602435906001600160a01b03821682036100c857565b60809060031901126100c8576001600160a01b039060043582811681036100c8579160243590811681036100c857906044359060643590565b9190820391821161090c57565b634e487b7160e01b5f52601160045260245ffd5b9190820180921161090c57565b60018060a01b0380911690815f52602060028152604091825f20941693845f528152815f208054908115610a5e575f91805b6109ba575050835f5260038252825f2061097a8282546108ff565b905580610989575b5050505050565b7f823cbf835fe0c4834f90d65da805f8a714ff7df3c298a0de29e4547944378f4f9251908152a35f80808080610982565b5f19016109c78183610ad0565b90549060031b1c865f526001808652865f20895f528652865f20825f5286526001600160801b03875f2054169081610a0d575b5050610a07829184610af9565b5061095f565b610a0791849396610a57928b5f528952895f208c5f528952895f20885f528952895f206001600160801b0319815416905560048952895f20610a508382546108ff565b9055610920565b94916109fa565b505050505050565b6001600160801b03918216908216039190821161090c57565b90601f8019910116810190811067ffffffffffffffff821117610aa157604052565b634e487b7160e01b5f52604160045260245ffd5b9190916001600160801b038080941691160191821161090c57565b8054821015610ae5575f5260205f2001905f90565b634e487b7160e01b5f52603260045260245ffd5b906001820191815f528260205260405f2054908115155f14610bcc575f199180830181811161090c5782549084820191821161090c57818103610b81575b50505080548015610b6d57820191610b4f8383610ad0565b909182549160031b1b19169055555f526020525f6040812055600190565b634e487b7160e01b5f52603160045260245ffd5b610bb7610b91610ba19386610ad0565b90549060031b1c92839286610ad0565b819391549060031b91821b915f19901b19161790565b90555f528460205260405f20555f8080610b37565b505050505f90565b6001810190825f528160205260405f2054155f14610c2657805468010000000000000000811015610aa157610c13610ba1826001879401855584610ad0565b905554915f5260205260405f2055600190565b5050505f9056fea26469706673582212203180aa92bc40b102643815287b84896103b2840f3d357fba0c61ec19f32ef4a764736f6c63430008180033",
}

// StakingVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingVaultMetaData.ABI instead.
var StakingVaultABI = StakingVaultMetaData.ABI

// StakingVaultBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StakingVaultMetaData.Bin instead.
var StakingVaultBin = StakingVaultMetaData.Bin

// DeployStakingVault deploys a new Ethereum contract, binding an instance of StakingVault to it.
func DeployStakingVault(auth *bind.TransactOpts, backend bind.ContractBackend, rootNet_ common.Address) (common.Address, *types.Transaction, *StakingVault, error) {
	parsed, err := StakingVaultMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StakingVaultBin), backend, rootNet_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StakingVault{StakingVaultCaller: StakingVaultCaller{contract: contract}, StakingVaultTransactor: StakingVaultTransactor{contract: contract}, StakingVaultFilterer: StakingVaultFilterer{contract: contract}}, nil
}

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

// RootNet is a free data retrieval call binding the contract method 0x405a0b06.
//
// Solidity: function rootNet() view returns(address)
func (_StakingVault *StakingVaultCaller) RootNet(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakingVault.contract.Call(opts, &out, "rootNet")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RootNet is a free data retrieval call binding the contract method 0x405a0b06.
//
// Solidity: function rootNet() view returns(address)
func (_StakingVault *StakingVaultSession) RootNet() (common.Address, error) {
	return _StakingVault.Contract.RootNet(&_StakingVault.CallOpts)
}

// RootNet is a free data retrieval call binding the contract method 0x405a0b06.
//
// Solidity: function rootNet() view returns(address)
func (_StakingVault *StakingVaultCallerSession) RootNet() (common.Address, error) {
	return _StakingVault.Contract.RootNet(&_StakingVault.CallOpts)
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
