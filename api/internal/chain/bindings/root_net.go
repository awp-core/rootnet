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

// IRootNetSubnetFullInfo is an auto generated low-level Go binding around an user-defined struct.
type IRootNetSubnetFullInfo struct {
	SubnetManager common.Address
	AlphaToken    common.Address
	LpPool        [32]byte
	Status        uint8
	CreatedAt     uint64
	ActivatedAt   uint64
	Name          string
	SkillsURI     string
	MinStake      *big.Int
	Owner         common.Address
}

// IRootNetSubnetInfo is an auto generated low-level Go binding around an user-defined struct.
type IRootNetSubnetInfo struct {
	LpPool      [32]byte
	Status      uint8
	CreatedAt   uint64
	ActivatedAt uint64
}

// IRootNetSubnetParams is an auto generated low-level Go binding around an user-defined struct.
type IRootNetSubnetParams struct {
	Name           string
	Symbol         string
	MetadataURI    string
	SubnetManager  common.Address
	CoordinatorURL string
	Salt           [32]byte
	MinStake       *big.Int
}

// RootNetAgentInfo is an auto generated low-level Go binding around an user-defined struct.
type RootNetAgentInfo struct {
	Owner           common.Address
	IsValid         bool
	Stake           *big.Int
	RewardRecipient common.Address
}

// RootNetMetaData contains all meta data concerning the RootNet contract.
var RootNetMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"deployer_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"treasury_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"guardian_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"INITIAL_ALPHA_MINT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_ACTIVE_SUBNETS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"accessManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activateSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allocate\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"alphaTokenFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"awpEmission\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"awpToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"banSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bind\",\"inputs\":[{\"name\":\"principal\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bindFor\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"principal\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deallocate\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultSubnetManagerImpl\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deregisterSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveSubnetCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveSubnetIdAt\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAgentInfo\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRootNet.AgentInfo\",\"components\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isValid\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"stake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"rewardRecipient\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAgentsInfo\",\"inputs\":[{\"name\":\"agents\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structRootNet.AgentInfo[]\",\"components\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isValid\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"stake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"rewardRecipient\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIRootNet.SubnetInfo\",\"components\":[{\"name\":\"lpPool\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumIRootNet.SubnetStatus\"},{\"name\":\"createdAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"activatedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSubnetFull\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIRootNet.SubnetFullInfo\",\"components\":[{\"name\":\"subnetManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"alphaToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lpPool\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumIRootNet.SubnetStatus\"},{\"name\":\"createdAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"activatedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"guardian\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"immunityPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialAlphaPrice\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initializeRegistry\",\"inputs\":[{\"name\":\"awpToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetNFT_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"alphaTokenFactory_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"awpEmission_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lpManager_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"accessManager_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakingVault_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakeNFT_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isSubnetActive\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lpManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nextSubnetId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nonces\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reallocate\",\"inputs\":[{\"name\":\"fromAgent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromSubnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"toAgent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"toSubnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"register\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"register\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"depositAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lockDuration\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerAndStake\",\"inputs\":[{\"name\":\"depositAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lockDuration\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"allocateAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerFor\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerSubnet\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIRootNet.SubnetParams\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"metadataURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"subnetManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"coordinatorURL\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerSubnetFor\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIRootNet.SubnetParams\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"metadataURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"subnetManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"coordinatorURL\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerSubnetForWithPermit\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIRootNet.SubnetParams\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"metadataURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"subnetManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"coordinatorURL\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"permitV\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"permitR\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"permitS\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"registerV\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"registerR\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"registerS\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registryInitialized\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeAgent\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resumeSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAlphaTokenFactory\",\"inputs\":[{\"name\":\"factory\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDelegation\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_isManager\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setGuardian\",\"inputs\":[{\"name\":\"g\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setImmunityPeriod\",\"inputs\":[{\"name\":\"p\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setInitialAlphaPrice\",\"inputs\":[{\"name\":\"price\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRewardRecipient\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSubnetManagerImpl\",\"inputs\":[{\"name\":\"impl\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakeNFT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stakingVault\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"subnetNFT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"subnets\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"lpPool\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumIRootNet.SubnetStatus\"},{\"name\":\"createdAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"activatedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"treasury\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unbanSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unbind\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMetadata\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"metadataURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"coordinatorURL\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AgentBound\",\"inputs\":[{\"name\":\"principal\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"oldPrincipal\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AgentRemoved\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AgentUnbound\",\"inputs\":[{\"name\":\"principal\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Allocated\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deallocated\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DelegationUpdated\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"isManager\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LPCreated\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"poolId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"awpAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"alphaAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MetadataUpdated\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"metadataURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"coordinatorURL\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Reallocated\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fromAgent\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"fromSubnet\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"toAgent\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"toSubnet\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardRecipientUpdated\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetActivated\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetBanned\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetDeregistered\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetPaused\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetRegistered\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"metadataURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"subnetManager\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"alphaToken\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"coordinatorURL\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetResumed\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetUnbanned\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UserRegistered\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpiredSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ImmunityNotExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientMinStake\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAgent\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidShortString\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSubnetParams\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSubnetStatus\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxActiveSubnetsReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotDeployer\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotGuardian\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotTimelock\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PriceTooHigh\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PriceTooLow\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"StringTooLong\",\"inputs\":[{\"name\":\"str\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"SubnetManagerRequired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnknownAddress\",\"inputs\":[]}]",
	Bin: "0x61016034620001fd57601f6200581838819003918201601f19168301926001600160401b0392909183851183861017620001e9578160609284926040978852833981010312620001fd5762000054816200021d565b6200006f8462000067602085016200021d565b93016200021d565b908451926200007e8462000201565b600a84526020840192691055d4149bdbdd13995d60b21b84528651620000a48162000201565b6001815260208101603160f81b815260ff195f54165f5560018055620000ca8762000232565b95610120968752620000dc8362000401565b97610140988952519020918260e05251902061010098818a524660a05280519160208301937f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f85528284015260608301524660808301523060a083015260a0825260c0820198828a10908a1117620001e9578890525190206080523060c0526001601155662386f26fc1000060145562278d0060155560018060a01b0380928160018060a01b0319951685600f541617600f551683600c541617600c551690600d541617600d5561526c9384620005ac853960805184614b37015260a05184614bf4015260c05184614aff015260e05184614b8601525183614bac015251826114ec015251816115180152f35b634e487b7160e01b5f52604160045260245ffd5b5f80fd5b604081019081106001600160401b03821117620001e957604052565b51906001600160a01b0382168203620001fd57565b805160209081811015620002cc5750601f8251116200026d57808251920151908083106200025f57501790565b825f19910360031b1b161790565b90604051809263305a27a960e01b82528060048301528251908160248401525f935b828510620002b2575050604492505f838284010152601f80199101168101030190fd5b84810182015186860160440152938101938593506200028f565b906001600160401b038211620001e957600254926001938481811c91168015620003f6575b83821014620003e257601f8111620003ab575b5081601f84116001146200034357509282939183925f9462000337575b50501b915f199060031b1c19161760025560ff90565b015192505f8062000321565b919083601f19811660025f52845f20945f905b8883831062000390575050501062000377575b505050811b0160025560ff90565b01515f1960f88460031b161c191690555f808062000369565b85870151885590960195948501948793509081019062000356565b60025f5284601f845f20920160051c820191601f860160051c015b828110620003d657505062000304565b5f8155018590620003c6565b634e487b7160e01b5f52602260045260245ffd5b90607f1690620002f1565b805160209190828110156200048f575090601f8251116200043057808251920151908083106200025f57501790565b90604051809263305a27a960e01b82528060048301528251908160248401525f935b82851062000475575050604492505f838284010152601f80199101168101030190fd5b848101820151868601604401529381019385935062000452565b6001600160401b038111620001e9576003928354926001938481811c91168015620005a0575b83821014620003e257601f81116200056a575b5081601f84116001146200050557509282939183925f94620004f9575b50501b915f1990841b1c191617905560ff90565b015192505f80620004e5565b919083601f198116875f52845f20945f905b888383106200054f575050501062000537575b505050811b01905560ff90565b01515f1983861b60f8161c191690555f80806200052a565b85870151885590960195948501948793509081019062000517565b855f5284601f845f20920160051c820191601f860160051c015b82811062000594575050620004c8565b5f815501859062000584565b90607f1690620004b556fe60a080604052600436101562000013575f80fd5b5f905f3560e01c90816303880623146200366c57508063080ec84114620035e35780630cf02c5e146200334257806311cba7e91462003317578063168f80f5146200306e5780631a46f4b81462002f035780631aa3a0081462002e655780631ddc304a1462002d8357806324e7964a1462002d585780632672e1be1462002d385780632bf1c05d1462002b4a57806333bbf0301462002afc5780633442656414620027ea57806338f48a8914620027825780633de3b24714620024605780633f4ba83a14620023e957806341a578cf14620023be57806344e047ca14620022c1578063452a93201462002296578063475726f714620022315780634b6f6d671462001eb55780635364944c1462001d9b57806356354a241462001d7357806358ca75041462001ccc5780635ab1bd531462001c405780635c975abb1462001c1c57806361d027b31462001bf1578063671a2a8a1462001a6957806367b26ba61462001a3e5780636d23f89514620017b65780636d345eea14620017965780637ab5e27614620017565780637b234b8114620016fa5780637ecebe0014620016bd57806381bac14f146200167a5780638456cb59146200160057806384b0196e14620014d35780638a0dac4a1462001481578063901a71e4146200142f57806397a6278e146200123c578063ab3f22d51462000fcf578063b400555a1462000fa7578063b48509e61462000f7c578063b6b257421462000e3c578063b79b76581462000c1d578063b906f15a1462000bf2578063be65e4c21462000bd3578063c1e0c9e71462000ba8578063c6a1a01a1462000b88578063c8c32a8414620009f4578063cd49dc031462000999578063cead1c961462000827578063d929ff051462000807578063e3da3f9b146200067c578063e521136f1462000531578063e7c1721214620004e6578063e7d89b71146200045d578063f4fda7261462000432578063fdcb606814620004075763fe427e9514620002ef575f80fd5b34620004045762000300366200388a565b6200030d92919262003e0b565b6200031762003e2f565b6200032162004820565b916200032e8184620048bd565b600a546001600160a01b039086908216803b15620004005760405163716fb83d60e01b81526001600160a01b0387811660048301528516602482015260448101889052606481018690529082908290608490829084905af18015620003f557620003d9575b5050604080519384523360208501529181169316917fd55bd7964253d1d9ce9187c8187b1c904274a3f374c9074f6de6fa77746bf34591819081015b0390a46001805580f35b620003e490620038b6565b620003f157855f62000393565b8580fd5b6040513d84823e3d90fd5b5080fd5b80fd5b503462000404578060031936011262000404576009546040516001600160a01b039091168152602090f35b50346200040457806003193601126200040457600e546040516001600160a01b039091168152602090f35b5034620004045760203660031901126200040457600c54600435906001600160a01b03163303620004d45764e8d4a510008110620004c2576c0c9f2c9cd04674edea400000008111620004b05760145580f35b60405163127f08c960e11b8152600490fd5b604051636dddf41160e11b8152600490fd5b60405163656a6d1560e11b8152600490fd5b503462000404576020366003190112620004045762000504620037cb565b600c546001600160a01b03919082163303620004d457166001600160601b0360a01b600e541617600e5580f35b50346200040457602036600319011262000404576200054f620037cb565b6009546040516364d8cc1960e11b81523360048201526001600160a01b0391821690606081602481855afa9081156200067157859162000638575b501562000626578084913b1562000400576040516350b2b4a960e11b81523360048201526001600160a01b03851660248201529082908290604490829084905af18015620003f5576200060a575b5050604051911681527fc8c11bb97ac2ffa10ce2e2a98f4c1fd8df84cfa2e1a15e013ed2383ab1f527ad60203392a280f35b6200061590620038b6565b6200062257825f620005d8565b8280fd5b60405163aba4733960e01b8152600490fd5b6200065f915060603d60601162000669575b6200065681836200393b565b810190620047ec565b5090505f6200058a565b503d6200064a565b6040513d87823e3d90fd5b50346200040457600319906101203683011262000404576200069d620037cb565b90602435926001600160401b038411620004005760e0908460040194360301126200040457604435620006cf62003810565b9160c4359260ff841684036200080357620006e962003e0b565b620006f362003e2f565b824211620007f157601454906a52b7d2dcc80cd2e40000009180830292830403620007dd576004546001600160a01b039190821690813b15620007d957878580949360ff60e494670de0b6b3a7640000604051998a98899763d505accf60e01b89521660048801523060248801520460448601528a606486015216608484015260843560a484015260a43560c48401525af18015620003f557620007c1575b6020620007b58787620007af6101043560e4358a8a868662003e97565b62004055565b60018055604051908152f35b620007cd8291620038b6565b62000404578062000792565b8480fd5b634e487b7160e01b83526011600452602483fd5b60405163df4cc36d60e01b8152600490fd5b5f80fd5b503462000404578060031936011262000404576020601154604051908152f35b503462000404576020806003193601126200040057600435906200084a62003e2f565b6005546040516331a9108f60e11b815260048101849052906001600160a01b039083908390602490829085165afa918215620006715785926200095a575b50339116036200094857601090828452526001604083200180549060ff8216620008b28162003821565b62000936576127106012541015620009245770ffffffffffffffff0000000000000000ff199091164260481b67ffffffffffffffff60481b16176001179055620008fc8162004a9f565b507f034804b969efac7a0df7757ada640ffdcc09f25dbcd4582c390f25d5622255c48280a280f35b60405163ab59c60f60e01b8152600490fd5b604051633671d60760e11b8152600490fd5b6040516330cd747160e01b8152600490fd5b9091508281813d831162000991575b6200097581836200393b565b81010312620007d95762000989906200395d565b905f62000888565b503d62000969565b50346200040457600319906020368301126200040457600435916001600160401b038311620004005760e090833603011262000404576020620007b583620009e062003e0b565b620009ea62003e2f565b6004013362004055565b50346200040457610100366003190112620004045762000a13620037cb565b62000a1d620037f9565b9062000a28620037e2565b906001600160a01b03606435818116908190036200080357608435828116809103620008035760a43591838316809303620008035760c43593808516809503620008035760e435968188168098036200080357600f54828116330362000b765760a01c60ff1662000b6557818a99816001600160601b0360a01b9916896004541617600455168760055416176005551685600654161760065584600754161760075583600854161760085582600954161760095583600a548284821617600a5584600b5494851617600b55161790813b1562000b60576024849291838093604051968795869463123c1a7b60e21b8652161760048401525af18015620003f55762000b48575b50600f80546001600160a81b031916600160a01b17905580f35b62000b5390620038b6565b6200040457805f62000b2e565b505050fd5b60405162dc149f60e41b8152600490fd5b604051638b906c9760e01b8152600490fd5b503462000404578060031936011262000404576020601254604051908152f35b503462000404578060031936011262000404576006546040516001600160a01b039091168152602090f35b5034620004045780600319360112620004045760206040516127108152f35b503462000404578060031936011262000404576008546040516001600160a01b039091168152602090f35b503462000404576020806003193601126200040057600c54600435916001600160a01b039182163303620004d45782845260108152600160408520019160ff8354169062000c6b8262003821565b600182149283158062000e25575b6200093657818792600554169160405163731865cb60e11b81528860048201528181602481875afa801562000671578391869162000de5575b5016928362000d16575b505050505062000ccc9062003821565b62000d04575b805460ff191660031790557fb887f21153957bddcf7211314cf42794076ccf98c6ae5cf8e2d065ec717c681b8280a280f35b62000d0f82620049bd565b5062000cd2565b81602491604051928380926363de476360e11b82528d60048301525afa9182156200067157859262000da3575b50501690813b1562000622578291604483926040519485938492635095056760e11b84526004840152600160248401525af18015620003f55762000d8b575b80808062000cbc565b62000d9690620038b6565b620007d957845f62000d82565b90809250813d831162000ddd575b62000dbd81836200393b565b8101031262000dd95762000dd1906200395d565b5f8062000d43565b8380fd5b503d62000db1565b809250838092503d831162000e1d575b62000e0181836200393b565b81010312620007d95762000e1683916200395d565b5f62000cb2565b503d62000df5565b5062000e318362003821565b600283141562000c79565b5034620004045780600319360112620004045762000e5962003e0b565b60095460405163cf5e7bd360e01b8152336004820152906001600160a01b039060209083906024908290879086165af191821562000f7157839262000f2d575b508281600a5416803b156200040057604051636f98081360e01b81526001600160a01b03851660048201523360248201529082908290604490829084905af18015620003f55762000f15575b50503391167f3e2d9d696fa5ddd5b13727a43861bb914938ca9d534d942f5c33725656c469b18380a36001805580f35b62000f2090620038b6565b6200062257825f62000ee5565b9091506020813d60201162000f68575b8162000f4c602093836200393b565b81010312620006225762000f60906200395d565b905f62000e99565b3d915062000f3d565b6040513d85823e3d90fd5b50346200040457806003193601126200040457600b546040516001600160a01b039091168152602090f35b5034620004045780600319360112620004045760206040516a52b7d2dcc80cd2e40000008152f35b5034620004045762000fe1366200388a565b62000fee92919262003e0b565b62000ff862003e2f565b6200100262004820565b91838552602060108152600160ff8160408920015416620010238162003821565b036200093657620010358285620048bd565b600a546001600160a01b03919087908316803b15620004005760405163d035a9a760e01b81526001600160a01b0388811660048301528616602482015260448101899052606481018790529082908290608490829084905af18015620003f55762001220575b50506024818360055416604051928380926373f231e760e01b82528b60048301525afa801562001215578890620011ce575b6001600160801b03915016908162001121575b5050604080519384523360208501529181169316917f655f98c7dae1bab3e2db10cdb4407717b9d219cf2e585adc1edba92d48af2b159181908101620003cf565b600a546040516378d6c06360e11b81526001600160a01b03888116600483015286166024820152604481018990529082908290606490829088165afa918215620011c35789926200118f575b5050106200117d575f80620010e0565b604051630f38cabd60e01b8152600490fd5b90809250813d8311620011bb575b620011a981836200393b565b810103126200080357515f806200116d565b503d6200119d565b6040513d8b823e3d90fd5b508181813d83116200120d575b620011e781836200393b565b810103126200120957620012036001600160801b039162003a21565b620010cd565b8780fd5b503d620011db565b6040513d8a823e3d90fd5b6200122b90620038b6565b6200123857865f6200109b565b8680fd5b5034620004045760208060031936011262000400576200125b620037cb565b6200126562003e0b565b6200126f62004820565b60095460405163275face960e11b81526001600160a01b03838116600483015284811660248301529294929183908290604490829086165afa90811562001424578691620013e7575b50158015620013d9575b620013c757808591600a5416803b156200062257604051636f98081360e01b81526001600160a01b038781166004830152861660248201529083908290604490829084905af190811562000f71578391620013af575b5050806009541694853b156200062257816064849283604051958694859363afe24ab360e01b855216998a600485015216998a60248401523360448401525af18015620003f55762001397575b50507f877ef5b4e3b78ab10b445521d0724510a2c3e98f0812879447b7e08785ca866e90604051338152a36001805580f35b620013a290620038b6565b62000dd957835f62001365565b620013ba90620038b6565b6200040057815f62001318565b6040516327f5ce6b60e01b8152600490fd5b5080841681841614620012c2565b90508281813d83116200141c575b6200140181836200393b565b81010312620003f1576200141590620039b8565b5f620012b8565b503d620013f5565b6040513d88823e3d90fd5b50346200040457602036600319011262000404576200144d620037cb565b600c546001600160a01b03919082163303620004d457168015620013c7576001600160601b0360a01b600654161760065580f35b50346200040457602036600319011262000404576200149f620037cb565b600c546001600160a01b03919082163303620004d457168015620013c7576001600160601b0360a01b600d541617600d5580f35b50346200040457806003193601126200040457620015117f000000000000000000000000000000000000000000000000000000000000000062004d37565b906200153d7f000000000000000000000000000000000000000000000000000000000000000062004e6c565b906040519060209060208301938385106001600160401b03861117620015ec579284926020620015a088966200159198604052858552604051988998600f60f81b8a5260e0858b015260e08a019062003863565b9088820360408a015262003863565b924660608801523060808801528460a088015286840360c088015251928381520193925b828110620015d457505050500390f35b835185528695509381019392810192600101620015c4565b634e487b7160e01b5f52604160045260245ffd5b50346200040457806003193601126200040457600d546001600160a01b0316330362001668576200163062003e2f565b600160ff198254161781557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2586020604051338152a180f35b6040516377b6878160e11b8152600490fd5b5034620004045760203660031901126200040457620016b66200169c620037cb565b620016a662003e0b565b620016b062003e2f565b62003cc1565b6001805580f35b50346200040457602036600319011262000404576020906040906001600160a01b03620016e9620037cb565b168152601683522054604051908152f35b5034620004045760c03660031901126200040457620016b66200171c620037cb565b62001726620037f9565b6200173062003810565b6200173a62003e0b565b6200174462003e2f565b60a43592608435926044359162003a87565b503462000404576020366003190112620004045760ff60016040602093600435815260108552200154166200178b8162003821565b600160405191148152f35b503462000404578060031936011262000404576020601454604051908152f35b5034620004045760603660031901126200040457620017d4620037cb565b6024356044356001600160401b038116908181036200080357620017f762003e0b565b6200180162003e2f565b600954604051631f5bdf5d60e01b81523360048201526020956001600160a01b039492909188918616908881602481855afa90811562000f7157839162001a01575b501562001999575b50848216918262001900575b5050508315159081620018f5575b5062001874575b846001805580f35b600b54604051630ed6579b60e11b815233600482015260248101949094526001600160401b039190911660448401528391839160649183918891165af1801562000f7157620018c6575b80806200186c565b813d8311620018ed575b620018dc81836200393b565b8101031262000803575f80620018be565b503d620018d0565b905015155f62001865565b8560095416803b1562000622576040516350b2b4a960e11b81523360048201526001600160a01b0392909216602483015282908290604490829084905af18015620003f55762001981575b50506040519081527fc8c11bb97ac2ffa10ce2e2a98f4c1fd8df84cfa2e1a15e013ed2383ab1f527ad863392a25f868162001857565b6200198c90620038b6565b6200123857865f6200194b565b803b156200040057818091602460405180948193632210724360e11b83523360048401525af18015620003f557620019e9575b5050335f80516020620052178339815191528880a2865f6200184b565b620019f490620038b6565b6200123857865f620019cc565b90508881813d831162001a36575b62001a1b81836200393b565b81010312620006225762001a2f90620039b8565b5f62001843565b503d62001a0f565b503462000404578060031936011262000404576007546040516001600160a01b039091168152602090f35b5034620004045760a0366003190112620004045762001a87620037cb565b6024359060443560ff81168103620008035762001aa362003e0b565b62001aad62003e2f565b824211620007f1576001600160a01b039182168085526016602052604085208054919491929062001ade8462003a78565b90556040519160208301937f1dc53c8f538cd1214fc81408a70baf66a4054f82e01699b136194e160ef3bce28552866040850152606084015260808301526080825260a082018281106001600160401b03821117620015ec578593859362001b6d9362001b639360405262001b5d60843593606435935190206200495f565b62004c1b565b9092919262004ca0565b160362001bdf57829060095416803b156200040057818091602460405180948193632210724360e11b83528860048401525af18015620003f55762001bc7575b50805f805160206200521783398151915291a26001805580f35b62001bd290620038b6565b6200040057815f62001bad565b604051638baa579f60e01b8152600490fd5b50346200040457806003193601126200040457600c546040516001600160a01b039091168152602090f35b5034620004045780600319360112620004045760ff60209154166040519015158152f35b5034620004045780600319360112620004045761014060018060a01b0380600454169080600554169080600654168160075416826008541683600954169084600a54169285600b54169486600c541696600d541697604051998a5260208a015260408901526060880152608087015260a086015260c085015260e0840152610100830152610120820152f35b5034620004045760203660031901126200040457604060809162001cef62003992565b50600435815260106020522060016040519162001d0c83620038ca565b805483520154602082019060ff811662001d268162003821565b82526001600160401b0380926040850190828460081c16825282606087019460481c168452604051955186525162001d5e8162003821565b60208601525116604084015251166060820152f35b50346200040457806003193601126200040457602060ff600f5460a01c166040519015158152f35b5034620004045760208060031936011262000400576004359062001dbe62003e2f565b6005546040516331a9108f60e11b815260048101849052906001600160a01b039083908390602490829085165afa9182156200067157859262001e76575b5033911603620009485760109082845252600160408320018054600260ff821662001e278162003821565b0362000936576127106012541015620009245760ff1916600117905562001e4e8162004a9f565b507ff1693a0767c0183c95caf97ea0be785bece8e3578b49ef89c9669b754c1ba9a08280a280f35b9091508281813d831162001ead575b62001e9181836200393b565b81010312620007d95762001ea5906200395d565b905f62001dfc565b503d62001e85565b50346200040457604036600319011262000404576001600160401b03806004351162000400573660236004350112156200040057600435600401351162000404573660246004356004013560051b600435010111620004045762001f1f6004356004013562003a36565b62001f2e60405191826200393b565b60048035013580825262001f429062003a36565b601f1901825b818110620022175750506009546001600160a01b0316825b60043560040135811062001fed578284604051918291602083016020845282518091526020604085019301915b81811062001f9c575050500390f35b91935091602060808262001fde60019488516060908160018060a01b039182815116855260208101511515602086015260408101516040860152015116910152565b01940191019184939262001f8d565b6200200260248260051b600435010162003a4e565b60405163fa54416160e01b81526001600160a01b03821660048201529190602083602481875afa92831562001424578693620021d3575b506001600160a01b03831615620021ca57600a546040516378d6c06360e11b81526001600160a01b0385811660048301529283166024808301919091523560448201529160209183916064918391165afa9081156200142457869162002192575b50915b6001600160a01b038116156200218757604051628195c360e11b81526001600160a01b038216600482015292602084602481885afa80156200217c5787906200213a575b60019450905b60405192620020f684620038ca565b60a086901b86900390811680855215156020850152604084019190915216606082015262002125828662003a63565b5262002132818562003a63565b500162001f60565b506020843d60201162002173575b8162002157602093836200393b565b8101031262001238576200216d6001946200395d565b620020e1565b3d915062002148565b6040513d89823e3d90fd5b6001928690620020e7565b90506020813d602011620021c1575b81620021b0602093836200393b565b810103126200080357515f6200209a565b3d9150620021a1565b5084916200209d565b9092506020813d6020116200220e575b81620021f2602093836200393b565b81010312620003f15762002206906200395d565b915f62002039565b3d9150620021e3565b6020906200222462003992565b8282860101520162001f48565b5034620004045760203660031901126200040457604060809160043581526010602052206001815491015460ff8116906001600160401b03916040519384526200227b8162003821565b6020840152818160081c16604084015260481c166060820152f35b50346200040457806003193601126200040457600d546040516001600160a01b039091168152602090f35b50346200040457602080600319360112620004005760043590602460018060a01b03828160055416604051938480926331a9108f60e11b82528860048301525afa918215620006715785926200237f575b5033911603620009485760109082845252600160408320018054600160ff82166200233d8162003821565b03620009365760ff191660021790556200235781620049bd565b507f789ca96cb827d1dcb6bfc7d9d084d0a574dadff90700e3602acedefb10f69afc8280a280f35b9091508281813d8311620023b6575b6200239a81836200393b565b81010312620007d957620023ae906200395d565b905f62002312565b503d6200238e565b503462000404578060031936011262000404576004546040516001600160a01b039091168152602090f35b50346200040457806003193601126200040457600c546001600160a01b03163303620004d457805460ff8116156200244e5760ff191681557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa6020604051338152a180f35b604051638dfc202b60e01b8152600490fd5b503462000404576020366003190112620004045760043590806101206040516200248a81620038e6565b8281528260208201528260408201528260608201528260808201528260a0820152606060c0820152606060e0820152826101008201520152818152601060205260408120908060018060a01b0360055416936024604051809681936342a3a26560e11b835260048301525afa928315620027755781936200268c575b505060018060a01b036020830151169160018060a01b0360408201511691600181549101546001600160401b038351916060850151936001600160801b036080870151169560a0600180821b03910151169660405198620025678a620038e6565b8952602089015260408801526200258160ff821662003821565b60ff81166060880152818160081c16608088015260481c1660a086015260c085015260e084015261010083015261012082015260405180916020825260018060a01b03815116602083015260018060a01b036020820151166040830152604081015160608301526060810151620025f88162003821565b60808301526001600160401b0360808201511660a08301526001600160401b0360a08201511660c083015260c08101516200265e62002646610140928360e087015261016086019062003863565b60e0840151858203601f190161010087015262003863565b6101008301516001600160801b031661012085810191909152909201516001600160a01b0316908301520390f35b909192503d8083833e620026a181836200393b565b8101602082820312620006225781516001600160401b0392838211620007d957019160c08383031262000dd9576040519360c0850185811083821117620015ec576040528351828111620004005783620026fd918601620039c6565b85526200270d602085016200395d565b602086015262002720604085016200395d565b60408601526060840151918211620004045750916200274860a09262002767948301620039c6565b60608501526200275b6080820162003a21565b6080850152016200395d565b60a0820152905f8062002506565b50604051903d90823e3d90fd5b5034620004045760203660031901126200040457600435601254811015620027d65760209060125f527fbb8a6a4669ba250d26cd7a459eca9d215f8307e33aebe50379bc5a3617ec34440154604051908152f35b634e487b7160e01b5f52603260045260245ffd5b5034620004045760a036600319011262000404576024356001600160401b038116600435818303620008035762002820620037e2565b9160643593608435926200283362003e0b565b6200283d62003e2f565b600954604051631f5bdf5d60e01b81523360048201526001600160a01b039460209491928a9087168685602481845afa8015620003f5578795839162002abe575b501562002a47575b5050831515908162002a3c575b50620029b9575b50505082151580620029ad575b80620029a3575b620028bc575b856001805580f35b60109085875252600160ff8160408820015416620028da8162003821565b036200093657620028ec8333620048bd565b8481600a5416803b15620004005760405163d035a9a760e01b81523360048201526001600160a01b038616602482015260448101879052606481018590529082908290608490829084905af18015620003f5576200298b575b50506040805192835233602084018190529390911692917f655f98c7dae1bab3e2db10cdb4407717b9d219cf2e585adc1edba92d48af2b159190a45f80808080620028b4565b6200299690620038b6565b620007d957845f62002945565b50841515620028ae565b508184161515620028a7565b600b54604051630ed6579b60e11b815233600482015260248101949094526001600160401b03919091166044840152829060649082908b9088165af19081156200217c57829162002a0d575b81906200289a565b813d831162002a34575b62002a2381836200393b565b810103126200080357805f62002a05565b503d62002a17565b905015155f62002893565b809194503b1562000dd957838091602460405180978193632210724360e11b83523360048401525af19384156200277557869462002a9d575b5050335f80516020620052178339815191528b80a2895f62002886565b62002aab91929450620038b6565b62002aba578391895f62002a80565b8880fd5b86819792503d831162002af4575b62002ad881836200393b565b81010312620004005762002aed8795620039b8565b5f6200287e565b503d62002acc565b5034620004045760203660031901126200040457600c54600435906001600160a01b03163303620004d45762093a80811062002b385760155580f35b6040516335354b8960e11b8152600490fd5b503462000404576020806003193601126200040057600c5460043591906001600160a01b039081163303620004d457828452601082526001604085200191600360ff84541662002b9a8162003821565b036200093657818592600554169160405163731865cb60e11b81528660048201528181602481875afa801562000671578391869162002cf8575b5016928362002c2e575b505050505061271060125410156200092457805460ff1916600117905562002c068162004a9f565b507fa04fe0f9f3200108443db1507380606e909a0f81c9eb84c0707c2651526684668280a280f35b81602491604051928380926363de476360e11b82528b60048301525afa9182156200067157859262002cba575b50501690813b1562000622578291604483926040519485938492635095056760e11b845260048401528160248401525af18015620003f55762002ca2575b80808062002bde565b62002cad90620038b6565b6200062257825f62002c99565b90809250813d831162002cf0575b62002cd481836200393b565b8101031262000dd95762002ce8906200395d565b5f8062002c5b565b503d62002cc8565b809250838092503d831162002d30575b62002d1481836200393b565b81010312620007d95762002d2983916200395d565b5f62002bd4565b503d62002d08565b503462000404578060031936011262000404576020601554604051908152f35b50346200040457806003193601126200040457600a546040516001600160a01b039091168152602090f35b503462000404576040366003190112620004045762002da1620037cb565b60243590811515809203620006225762002dba62004820565b6009546001600160a01b039081169391859190853b156200062257816084849283604051958694859363902e164560e01b855216998a600485015216998a60248401528860448401523360648401525af18015620003f55762002e4d575b505060407f34dbef79b9de038294b4a8f1789ad62e1b9ebaa23af56a3b75f375ce1185a9b1918151908152336020820152a380f35b62002e5890620038b6565b62000dd957835f62002e18565b5034620004045780600319360112620004045762002e8262003e0b565b62002e8c62003e2f565b60095481906001600160a01b0316803b1562002f0057818091602460405180948193632210724360e11b83523360048401525af18015620003f55762002ee8575b50335f80516020620052178339815191528280a26001805580f35b62002ef390620038b6565b6200040457805f62002ecd565b50fd5b5034620004045760a0366003190112620004045762002f21620037cb565b60243562002f2e620037e2565b60843560643562002f3e62003e0b565b62002f4862003e2f565b62002f5262004820565b938187526010602052600160ff8160408a2001541662002f728162003821565b03620009365762002f848686620048bd565b62002f908486620048bd565b600a546001600160a01b0390811695889190873b1562000622578160c484928383604051968795869463d5d5278d60e01b8652169e8f6004860152169b8c6024850152896044850152169b8c60648401528960848401528a60a48401525af18015620003f5576200304e575b5050907f726c93ba67bfe4c677e37114279f0ad9aab5ee9ffbd1158923be5d0fec3b1b459460c094939260405194855260208501526040840152606083015260808201523360a0820152a26001805580f35b6200305e909594939295620038b6565b6200123857909192865f62002ffc565b50346200040457604036600319011262000404576200308c620037cb565b906200309762003992565b5060095460405163fa54416160e01b81526001600160a01b03848116600483015292918316919060208082602481875afa91821562000f71578392620032d8575b5060405163275face960e11b81526001600160a01b03838116600483015287166024820152958187604481885afa9687156200328a57849762003295575b50600a546040516378d6c06360e11b81526001600160a01b0380861660048301529092166024808401919091523560448301529096908290889088168180606481015b03915afa9283156200328a578697859796979462003251575b5082906024604051809a8193628195c360e11b835216998a60048301525afa9384156200324557809462003203575b506080965060405195620031b587620038ca565b865215159085015260408401521660608201526200320160405180926060908160018060a01b039182815116855260208101511515602086015260408101516040860152015116910152565bf35b9093508187813d83116200323d575b6200321e81836200393b565b81010312620004045750620032356080966200395d565b925f620031a1565b503d62003212565b604051903d90823e3d90fd5b955092508185813d811162003282575b6200326d81836200393b565b81010312620008035781879551939062003172565b503d62003261565b6040513d86823e3d90fd5b96508187813d8311620032d0575b620032af81836200393b565b8101031262000dd95781620032c86200315998620039b8565b975062003116565b503d620032a3565b9080925081813d83116200330f575b620032f381836200393b565b81010312620006225762003307906200395d565b905f620030d8565b503d620032e7565b503462000404578060031936011262000404576005546040516001600160a01b039091168152602090f35b503462000803576020806003193601126200080357600c5460043591906001600160a01b039081163303620004d457825f5260108252600160405f2001546001600160401b03808260081c16918215620009365760481c16908115620035db57505b6015548101809111620035c757421115620035b55760055460405163731865cb60e11b815260048101859052929082168184602481845afa8015620035295786945f9162003574575b508316908162003490575b50506010906200340885620049bd565b5084845252816001604082205f8155015560055416803b156200040057818091602460405180948193630852cd8d60e31b83528860048401525af18015620003f55762003478575b50807f960c7566f4c9bb6958ff6e37a02b5ae69fa39dd75651fcc9b9a1029c01d0ff3291a280f35b6200348390620038b6565b6200040057815f62003450565b829394509160249192604051928380926363de476360e11b82528960048301525afa8015620035295784915f9162003534575b501690813b1562000803575f91604483926040519485938492635095056760e11b84526004840152600160248401525af1801562003529576200350c575b9081859392620033f8565b601094506200351e90929192620038b6565b5f9391909162003501565b6040513d5f823e3d90fd5b809250848092503d83116200356c575b6200355081836200393b565b8101031262000803576200356584916200395d565b5f620034c3565b503d62003544565b809550838092503d8311620035ad575b6200359081836200393b565b81010312620008035782620035a687956200395d565b90620033ed565b503d62003584565b6040516384e3b93f60e01b8152600490fd5b634e487b7160e01b5f52601160045260245ffd5b9050620033a4565b34620008035760031960c036820112620008035762003601620037cb565b602435916001600160401b038311620008035760e0908360040193360301126200080357604435916200363362003810565b916200363e62003e0b565b6200364862003e2f565b834211620007f157620007af620007b59360209560a4359160843591868662003e97565b34620008035760603660031901126200080357600435906001600160401b036024358181116200080357620036a69036906004016200379b565b90916044359081116200080357620036c39036906004016200379b565b909360018060a01b0360208260248184600554166331a9108f60e11b82528b60048301525afa91821562003529575f9262003757575b503391160362000948577f4bb348c8e52124f1a18e983f64ad1bc5d380a3ca43654fbb2b8f73c71f3054599362003752916200374360405195869560408752604087019162003972565b91848303602086015262003972565b0390a2005b9091506020813d60201162003792575b8162003776602093836200393b565b8101031262000803576200378a906200395d565b9087620036f9565b3d915062003767565b9181601f8401121562000803578235916001600160401b0383116200080357602083818601950101116200080357565b600435906001600160a01b03821682036200080357565b604435906001600160a01b03821682036200080357565b602435906001600160a01b03821682036200080357565b6064359060ff821682036200080357565b600411156200382c57565b634e487b7160e01b5f52602160045260245ffd5b5f5b838110620038525750505f910152565b818101518382015260200162003842565b906020916200387e8151809281855285808601910162003840565b601f01601f1916010190565b606090600319011262000803576004356001600160a01b03811681036200080357906024359060443590565b6001600160401b038111620015ec57604052565b608081019081106001600160401b03821117620015ec57604052565b61014081019081106001600160401b03821117620015ec57604052565b60c081019081106001600160401b03821117620015ec57604052565b604081019081106001600160401b03821117620015ec57604052565b90601f801991011681019081106001600160401b03821117620015ec57604052565b51906001600160a01b03821682036200080357565b908060209392818452848401375f828201840152601f01601f1916010190565b60405190620039a182620038ca565b5f6060838281528260208201528260408201520152565b519081151582036200080357565b81601f82011215620008035780516001600160401b038111620015ec5760405192620039fd601f8301601f1916602001856200393b565b81845260208284010111620008035762003a1e916020808501910162003840565b90565b51906001600160801b03821682036200080357565b6001600160401b038111620015ec5760051b60200190565b356001600160a01b0381168103620008035790565b8051821015620027d65760209160051b010190565b5f198114620035c75760010190565b9594936001600160a01b038281169493909290918515620013c757824211620007f1578362003b3662001b63828c16998a94855f5260209960168b5262001b5d8c8c60409b8c92835f209081549162003ae08362003a78565b90558c8551948501957fdb3dbac9db5f555ea8f43e6a1bea52126fc17fc7c5afa6b3caf705de1ba2895587528501526060840152608083015260a082015260a0815262003b2d8162003903565b5190206200495f565b160362003cb057600954825163347b27c560e21b81526001600160a01b038a811660048301529290921660248301528490829060449082905f9088165af190811562003ca6575f9162003c69575b508281169285841462003c5e578362003bc7575b50507f4ee4aa1bbc31e8b57dad2c2cffa4627ad65ac133b3cea2acb4870c44b5ea6b17939495965051908152a3565b600a541697883b1562000803578251636f98081360e01b81526001600160a01b03928316600482015291166024820152965f908890604490829084905af196871562003c54577f4ee4aa1bbc31e8b57dad2c2cffa4627ad65ac133b3cea2acb4870c44b5ea6b179495969762003c42575b8796959462003b98565b62003c4d90620038b6565b5f62003c38565b50513d5f823e3d90fd5b505050505050509050565b90508381813d831162003c9e575b62003c8381836200393b565b81010312620008035762003c97906200395d565b5f62003b84565b503d62003c77565b82513d5f823e3d90fd5b8151638baa579f60e01b8152600490fd5b6001600160a01b0390808216908115620013c75760095460405163347b27c560e21b81523360048201526001600160a01b03929092166024830152602090829060449082905f9088165af190811562003529575f9162003dc9575b508281169282841462003dc3578362003d60575b50506040519182527f4ee4aa1bbc31e8b57dad2c2cffa4627ad65ac133b3cea2acb4870c44b5ea6b1760203393a3565b600a5416803b156200080357604051636f98081360e01b81526001600160a01b039290921660048301523360248301525f908290604490829084905af18015620035295762003db1575b8062003d30565b62003dbc90620038b6565b5f62003daa565b50505050565b90506020813d60201162003e02575b8162003de7602093836200393b565b81010312620008035762003dfb906200395d565b5f62003d1c565b3d915062003dd8565b60026001541462003e1d576002600155565b604051633ee5aeb560e01b8152600490fd5b60ff5f541662003e3b57565b60405163d93c066560e01b8152600490fd5b903590601e19813603018212156200080357018035906001600160401b03821162000803576020019181360383136200080357565b356001600160801b0381168103620008035790565b92919394909462003ea9868062003e4d565b9490966020810162003ebc908262003e4d565b6040979197809a81850162003ed2908662003e4d565b95909a6060820162003ee49062003a4e565b9662003ef4608084018462003e4d565b9d909162003f0560c0860162003e82565b938751988997602089019b60e08d526101008a019062003f259262003972565b90601f19998a8a840301908a015262003f3e9262003972565b908787830301606088015262003f549262003972565b600160a01b600190039d8e809a166080870152868683030160a087015262003f7c9262003972565b9160a0013560c08401526001600160801b031660e083015203908101825262003fa690826200393b565b519020911695865f526016602052875f209081549162003fc68362003a78565b905588519160208301937f481d00d71014332dae55eae70bcf8183dfc8d771e5f91fa53819b22567e609bc8552898b8501526060840152608083015260a082015260a08152620040168162003903565b51902062004024906200495f565b92620040309362004c1b565b6200403b9162004ca0565b1603620040455750565b51638baa579f60e01b8152600490fd5b9162004062828062003e4d565b5f939150158015620047d4575b62002b3857602081019062004085828262003e4d565b9050158015620047bc575b62002b38576001600160a01b03620040ab6060830162003a4e565b1615918280620047a8575b620047965760145491826a52b7d2dcc80cd2e400000002926a52b7d2dcc80cd2e4000000840403620035c7576004546008546040516323b872dd60e01b60208083019182526001600160a01b038c811660248501529384166044840152670de0b6b3a76400008804606480850191909152835292909316925f916200413d6084826200393b565b519082855af11562003529575f513d6200478c5750803b155b620047745750620041cf9060115494620041708662003a78565b6011556020620041e262004185848062003e4d565b5f8a620041968a8996959662003e4d565b949060018060a01b0360065416956040519b8c98899788966332c6cee560e21b8852600488015260a0602488015260a487019162003972565b8481036003190160448601529162003972565b30606483015260a0890135608483015203925af192831562003529575f9362004730575b506008546001600160a01b039081169084163b1562000803576040516340c10f1960e01b815260048101919091526a52b7d2dcc80cd2e400000060248201525f81604481836001600160a01b0389165af18015620035295762004716575b50600854604080516321946c7360e21b81526001600160a01b038681166004830152670de0b6b3a7640000890460248301526a52b7d2dcc80cd2e400000060448301529a9b999a98999098909392899160649183918f91165af19687156200470b578a97620046d1575b5015620046bc57600454604051630a31ee5b60e41b60208201526001600160a01b0385811660248301529182166044820152606481018890529082166084808301919091528152620043208162003903565b600e5460405191906102c8808401916001600160a01b0316906001600160401b03831185841017620046a8576200436d9285949260409262004f4f87398152816020820152019062003863565b03908af08015620011c3576001600160a01b0316925b6001600160a01b0381163b15620046815760405163c3897a6760e01b81526001600160a01b0385811660048301528b9190829082906024908290849088165af18015620003f55762004690575b50506005546001600160a01b03168a620043eb858062003e4d565b9092620043fb60c0880162003e82565b813b1562000dd9576040516335c1b08160e01b81526001600160a01b0388166004820152602481018e905260c060448201529485938492869284926001600160801b03916200444f9160c486019162003972565b6001600160a01b03808f1660648601528b166084850152911660a483015203925af18015620046855762004669575b507ff1754be0b0979fb871647c40cf65eaca03d25e211ac7b3016ad3d70a49845173879560c062004645670de0b6b3a76400009760609a978d9e9f978e97620045957f0a28a1fd5e0909199ee082834df66cfaae2125e3503bf16d2dc46214278fc7ab9f620046149b8b6001600160401b0391604051936200450085620038ca565b84526001602085019482865260408082019386421685526060830195818752815260106020522090518155019351906200453a8262003821565b620045458262003821565b60ff68ffffffffffffffff008654925160081b1692169068ffffffffffffffffff191617178355511667ffffffffffffffff60481b82549160481b169067ffffffffffffffff60481b1916179055565b620045b1620045a5848062003e4d565b90608052978462003e4d565b9a9093620045d962004605620045cb604084018462003e4d565b929093608081019062003e4d565b97909e620045f66040519d8d8f9e8f908152019060805162003972565b8c810360208e01529162003972565b9189830360408b015262003972565b9360018060a01b03168e87015260018060a01b0316608086015284830360a086015260018060a01b03169762003972565b0390a36040519283520460208201526a52b7d2dcc80cd2e40000006040820152a290565b620046758b91620038b6565b62004681575f6200447e565b8980fd5b6040513d8d823e3d90fd5b6200469b90620038b6565b6200468157895f620043d0565b634e487b7160e01b8e52604160045260248efd5b620046ca6060830162003a4e565b9262004383565b9096506040813d60401162004702575b81620046f0604093836200393b565b81010312620046815751955f620042ce565b3d9150620046e1565b6040513d8c823e3d90fd5b620047259196979850620038b6565b5f9695945f62004264565b9092506020813d6020116200476b575b816200474f602093836200393b565b81010312620008035762004763906200395d565b915f62004206565b3d915062004740565b60249060405190635274afe760e01b82526004820152fd5b6001141562004156565b604051637991d6f760e01b8152600490fd5b50600e546001600160a01b031615620040b6565b506010620047cb838362003e4d565b90501162004090565b506040620047e3828062003e4d565b9050116200406f565b90816060910312620008035762004803816200395d565b9162003a1e60406200481860208501620039b8565b9301620039b8565b6009546040516364d8cc1960e11b8152336004820152906001600160a01b03906060908390602490829085165afa90811562003529575f925f905f9362004892575b506200488b57821615620013c75715620048795790565b60405163607e454560e11b8152600490fd5b5050503390565b91935050620048b3915060603d60601162000669576200065681836200393b565b9192905f62004862565b60095460405163275face960e11b81526001600160a01b039283166004820152928216602484015260209183916044918391165afa90811562003529575f916200491d575b50156200490b57565b60405163bebdc75760e01b8152600490fd5b90506020813d60201162004956575b816200493b602093836200393b565b8101031262000803576200494f90620039b8565b5f62004902565b3d91506200492c565b6042906200496c62004afc565b906040519161190160f01b8352600283015260228201522090565b601254811015620027d65760125f527fbb8a6a4669ba250d26cd7a459eca9d215f8307e33aebe50379bc5a3617ec344401905f90565b5f81815260136020526040902054801562004a99575f1990808201818111620035c75760125490838201918211620035c75781810362004a49575b505050601254801562004a355781019062004a138262004987565b909182549160031b1b191690556012555f5260136020525f6040812055600190565b634e487b7160e01b5f52603160045260245ffd5b62004a8262004a5c62004a6c9362004987565b90549060031b1c92839262004987565b819391549060031b91821b915f19901b19161790565b90555f52601360205260405f20555f8080620049f8565b50505f90565b805f52601360205260405f2054155f1462004af75760125468010000000000000000811015620015ec5762004ae062004a6c82600185940160125562004987565b9055601254905f52601360205260405f2055600190565b505f90565b307f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316148062004bf1575b1562004b59577f000000000000000000000000000000000000000000000000000000000000000090565b60405160208101907f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f82527f000000000000000000000000000000000000000000000000000000000000000060408201527f000000000000000000000000000000000000000000000000000000000000000060608201524660808201523060a082015260a0815262004beb8162003903565b51902090565b507f0000000000000000000000000000000000000000000000000000000000000000461462004b2f565b91907f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0841162004c95579160209360809260ff5f9560405194855216868401526040830152606082015282805260015afa1562003529575f516001600160a01b0381161562004c8b57905f905f90565b505f906001905f90565b5050505f9160039190565b62004cab8162003821565b8062004cb5575050565b62004cc08162003821565b6001810362004cdb5760405163f645eedf60e01b8152600490fd5b62004ce68162003821565b6002810362004d085760405163fce698f760e01b815260048101839052602490fd5b8062004d1660039262003821565b1462004d1f5750565b602490604051906335e2f38360e21b82526004820152fd5b60ff811462004d795760ff811690601f821162004d67576040519162004d5d836200391f565b8252602082015290565b604051632cd44ac360e21b8152600490fd5b506040515f600254906001908260011c6001841692831562004e61575b602094858310851462004e4d57828752869490811562004e2b575060011462004dca575b505062003a1e925003826200393b565b9093915060025f527f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace935f915b81831062004e1257505062003a1e93508201015f8062004dba565b8554878401850152948501948694509183019162004df7565b91505062003a1e94925060ff191682840152151560051b8201015f8062004dba565b634e487b7160e01b5f52602260045260245ffd5b90607f169062004d96565b60ff811462004e925760ff811690601f821162004d67576040519162004d5d836200391f565b506040515f600354906001908260011c6001841692831562004f43575b602094858310851462004e4d57828752869490811562004e2b575060011462004ee257505062003a1e925003826200393b565b9093915060035f527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b935f915b81831062004f2a57505062003a1e93508201015f8062004dba565b8554878401850152948501948694509183019162004f0f565b90607f169062004eaf56fe60806040526102c8803803806100148161018e565b92833981019060408183031261018a5780516001600160a01b03811680820361018a5760208381015190936001600160401b03821161018a570184601f8201121561018a5780519061006d610068836101c7565b61018e565b9582875285838301011161018a5784905f5b8381106101765750505f9186010152813b1561015e577f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80546001600160a01b03191682179055604051907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a283511561014057505f80848461012796519101845af4903d15610137573d610118610068826101c7565b9081525f81943d92013e6101e2565b505b604051608290816102468239f35b606092506101e2565b925050503461014f5750610129565b63b398979f60e01b8152600490fd5b60249060405190634c9c8ce360e01b82526004820152fd5b81810183015188820184015286920161007f565b5f80fd5b6040519190601f01601f191682016001600160401b038111838210176101b357604052565b634e487b7160e01b5f52604160045260245ffd5b6001600160401b0381116101b357601f01601f191660200190565b9061020957508051156101f757805190602001fd5b60405163d6bda27560e01b8152600490fd5b8151158061023c575b61021a575090565b604051639996b31560e01b81526001600160a01b039091166004820152602490fd5b50803b1561021256fe60806040527f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc545f9081906001600160a01b0316368280378136915af43d5f803e156048573d5ff35b3d5ffdfea2646970667358221220cd27d5af5088e1cdce52ec7ad4a5de5c98b49ccf39972645ad4e873925d11c1864736f6c6343000818003354db7a5cb4735e1aac1f53db512d3390390bb6637bd30ad4bf9fc98667d9b9b9a26469706673582212205ce86d2067a0b0f12a7df0151b2b1b39c40e6c36d7789666be0bface8744f5a764736f6c63430008180033",
}

// RootNetABI is the input ABI used to generate the binding from.
// Deprecated: Use RootNetMetaData.ABI instead.
var RootNetABI = RootNetMetaData.ABI

// RootNetBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RootNetMetaData.Bin instead.
var RootNetBin = RootNetMetaData.Bin

// DeployRootNet deploys a new Ethereum contract, binding an instance of RootNet to it.
func DeployRootNet(auth *bind.TransactOpts, backend bind.ContractBackend, deployer_ common.Address, treasury_ common.Address, guardian_ common.Address) (common.Address, *types.Transaction, *RootNet, error) {
	parsed, err := RootNetMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RootNetBin), backend, deployer_, treasury_, guardian_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &RootNet{RootNetCaller: RootNetCaller{contract: contract}, RootNetTransactor: RootNetTransactor{contract: contract}, RootNetFilterer: RootNetFilterer{contract: contract}}, nil
}

// RootNet is an auto generated Go binding around an Ethereum contract.
type RootNet struct {
	RootNetCaller     // Read-only binding to the contract
	RootNetTransactor // Write-only binding to the contract
	RootNetFilterer   // Log filterer for contract events
}

// RootNetCaller is an auto generated read-only Go binding around an Ethereum contract.
type RootNetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootNetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RootNetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootNetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RootNetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootNetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RootNetSession struct {
	Contract     *RootNet          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RootNetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RootNetCallerSession struct {
	Contract *RootNetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// RootNetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RootNetTransactorSession struct {
	Contract     *RootNetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// RootNetRaw is an auto generated low-level Go binding around an Ethereum contract.
type RootNetRaw struct {
	Contract *RootNet // Generic contract binding to access the raw methods on
}

// RootNetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RootNetCallerRaw struct {
	Contract *RootNetCaller // Generic read-only contract binding to access the raw methods on
}

// RootNetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RootNetTransactorRaw struct {
	Contract *RootNetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRootNet creates a new instance of RootNet, bound to a specific deployed contract.
func NewRootNet(address common.Address, backend bind.ContractBackend) (*RootNet, error) {
	contract, err := bindRootNet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RootNet{RootNetCaller: RootNetCaller{contract: contract}, RootNetTransactor: RootNetTransactor{contract: contract}, RootNetFilterer: RootNetFilterer{contract: contract}}, nil
}

// NewRootNetCaller creates a new read-only instance of RootNet, bound to a specific deployed contract.
func NewRootNetCaller(address common.Address, caller bind.ContractCaller) (*RootNetCaller, error) {
	contract, err := bindRootNet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RootNetCaller{contract: contract}, nil
}

// NewRootNetTransactor creates a new write-only instance of RootNet, bound to a specific deployed contract.
func NewRootNetTransactor(address common.Address, transactor bind.ContractTransactor) (*RootNetTransactor, error) {
	contract, err := bindRootNet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RootNetTransactor{contract: contract}, nil
}

// NewRootNetFilterer creates a new log filterer instance of RootNet, bound to a specific deployed contract.
func NewRootNetFilterer(address common.Address, filterer bind.ContractFilterer) (*RootNetFilterer, error) {
	contract, err := bindRootNet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RootNetFilterer{contract: contract}, nil
}

// bindRootNet binds a generic wrapper to an already deployed contract.
func bindRootNet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RootNetMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RootNet *RootNetRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RootNet.Contract.RootNetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RootNet *RootNetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootNet.Contract.RootNetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RootNet *RootNetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RootNet.Contract.RootNetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RootNet *RootNetCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RootNet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RootNet *RootNetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootNet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RootNet *RootNetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RootNet.Contract.contract.Transact(opts, method, params...)
}

// INITIALALPHAMINT is a free data retrieval call binding the contract method 0xb400555a.
//
// Solidity: function INITIAL_ALPHA_MINT() view returns(uint256)
func (_RootNet *RootNetCaller) INITIALALPHAMINT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "INITIAL_ALPHA_MINT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// INITIALALPHAMINT is a free data retrieval call binding the contract method 0xb400555a.
//
// Solidity: function INITIAL_ALPHA_MINT() view returns(uint256)
func (_RootNet *RootNetSession) INITIALALPHAMINT() (*big.Int, error) {
	return _RootNet.Contract.INITIALALPHAMINT(&_RootNet.CallOpts)
}

// INITIALALPHAMINT is a free data retrieval call binding the contract method 0xb400555a.
//
// Solidity: function INITIAL_ALPHA_MINT() view returns(uint256)
func (_RootNet *RootNetCallerSession) INITIALALPHAMINT() (*big.Int, error) {
	return _RootNet.Contract.INITIALALPHAMINT(&_RootNet.CallOpts)
}

// MAXACTIVESUBNETS is a free data retrieval call binding the contract method 0xbe65e4c2.
//
// Solidity: function MAX_ACTIVE_SUBNETS() view returns(uint128)
func (_RootNet *RootNetCaller) MAXACTIVESUBNETS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "MAX_ACTIVE_SUBNETS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXACTIVESUBNETS is a free data retrieval call binding the contract method 0xbe65e4c2.
//
// Solidity: function MAX_ACTIVE_SUBNETS() view returns(uint128)
func (_RootNet *RootNetSession) MAXACTIVESUBNETS() (*big.Int, error) {
	return _RootNet.Contract.MAXACTIVESUBNETS(&_RootNet.CallOpts)
}

// MAXACTIVESUBNETS is a free data retrieval call binding the contract method 0xbe65e4c2.
//
// Solidity: function MAX_ACTIVE_SUBNETS() view returns(uint128)
func (_RootNet *RootNetCallerSession) MAXACTIVESUBNETS() (*big.Int, error) {
	return _RootNet.Contract.MAXACTIVESUBNETS(&_RootNet.CallOpts)
}

// AccessManager is a free data retrieval call binding the contract method 0xfdcb6068.
//
// Solidity: function accessManager() view returns(address)
func (_RootNet *RootNetCaller) AccessManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "accessManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AccessManager is a free data retrieval call binding the contract method 0xfdcb6068.
//
// Solidity: function accessManager() view returns(address)
func (_RootNet *RootNetSession) AccessManager() (common.Address, error) {
	return _RootNet.Contract.AccessManager(&_RootNet.CallOpts)
}

// AccessManager is a free data retrieval call binding the contract method 0xfdcb6068.
//
// Solidity: function accessManager() view returns(address)
func (_RootNet *RootNetCallerSession) AccessManager() (common.Address, error) {
	return _RootNet.Contract.AccessManager(&_RootNet.CallOpts)
}

// AlphaTokenFactory is a free data retrieval call binding the contract method 0xc1e0c9e7.
//
// Solidity: function alphaTokenFactory() view returns(address)
func (_RootNet *RootNetCaller) AlphaTokenFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "alphaTokenFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AlphaTokenFactory is a free data retrieval call binding the contract method 0xc1e0c9e7.
//
// Solidity: function alphaTokenFactory() view returns(address)
func (_RootNet *RootNetSession) AlphaTokenFactory() (common.Address, error) {
	return _RootNet.Contract.AlphaTokenFactory(&_RootNet.CallOpts)
}

// AlphaTokenFactory is a free data retrieval call binding the contract method 0xc1e0c9e7.
//
// Solidity: function alphaTokenFactory() view returns(address)
func (_RootNet *RootNetCallerSession) AlphaTokenFactory() (common.Address, error) {
	return _RootNet.Contract.AlphaTokenFactory(&_RootNet.CallOpts)
}

// AwpEmission is a free data retrieval call binding the contract method 0x67b26ba6.
//
// Solidity: function awpEmission() view returns(address)
func (_RootNet *RootNetCaller) AwpEmission(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "awpEmission")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpEmission is a free data retrieval call binding the contract method 0x67b26ba6.
//
// Solidity: function awpEmission() view returns(address)
func (_RootNet *RootNetSession) AwpEmission() (common.Address, error) {
	return _RootNet.Contract.AwpEmission(&_RootNet.CallOpts)
}

// AwpEmission is a free data retrieval call binding the contract method 0x67b26ba6.
//
// Solidity: function awpEmission() view returns(address)
func (_RootNet *RootNetCallerSession) AwpEmission() (common.Address, error) {
	return _RootNet.Contract.AwpEmission(&_RootNet.CallOpts)
}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_RootNet *RootNetCaller) AwpToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "awpToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_RootNet *RootNetSession) AwpToken() (common.Address, error) {
	return _RootNet.Contract.AwpToken(&_RootNet.CallOpts)
}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_RootNet *RootNetCallerSession) AwpToken() (common.Address, error) {
	return _RootNet.Contract.AwpToken(&_RootNet.CallOpts)
}

// DefaultSubnetManagerImpl is a free data retrieval call binding the contract method 0xf4fda726.
//
// Solidity: function defaultSubnetManagerImpl() view returns(address)
func (_RootNet *RootNetCaller) DefaultSubnetManagerImpl(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "defaultSubnetManagerImpl")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DefaultSubnetManagerImpl is a free data retrieval call binding the contract method 0xf4fda726.
//
// Solidity: function defaultSubnetManagerImpl() view returns(address)
func (_RootNet *RootNetSession) DefaultSubnetManagerImpl() (common.Address, error) {
	return _RootNet.Contract.DefaultSubnetManagerImpl(&_RootNet.CallOpts)
}

// DefaultSubnetManagerImpl is a free data retrieval call binding the contract method 0xf4fda726.
//
// Solidity: function defaultSubnetManagerImpl() view returns(address)
func (_RootNet *RootNetCallerSession) DefaultSubnetManagerImpl() (common.Address, error) {
	return _RootNet.Contract.DefaultSubnetManagerImpl(&_RootNet.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_RootNet *RootNetCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "eip712Domain")

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
func (_RootNet *RootNetSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _RootNet.Contract.Eip712Domain(&_RootNet.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_RootNet *RootNetCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _RootNet.Contract.Eip712Domain(&_RootNet.CallOpts)
}

// GetActiveSubnetCount is a free data retrieval call binding the contract method 0xc6a1a01a.
//
// Solidity: function getActiveSubnetCount() view returns(uint256)
func (_RootNet *RootNetCaller) GetActiveSubnetCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "getActiveSubnetCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActiveSubnetCount is a free data retrieval call binding the contract method 0xc6a1a01a.
//
// Solidity: function getActiveSubnetCount() view returns(uint256)
func (_RootNet *RootNetSession) GetActiveSubnetCount() (*big.Int, error) {
	return _RootNet.Contract.GetActiveSubnetCount(&_RootNet.CallOpts)
}

// GetActiveSubnetCount is a free data retrieval call binding the contract method 0xc6a1a01a.
//
// Solidity: function getActiveSubnetCount() view returns(uint256)
func (_RootNet *RootNetCallerSession) GetActiveSubnetCount() (*big.Int, error) {
	return _RootNet.Contract.GetActiveSubnetCount(&_RootNet.CallOpts)
}

// GetActiveSubnetIdAt is a free data retrieval call binding the contract method 0x38f48a89.
//
// Solidity: function getActiveSubnetIdAt(uint256 index) view returns(uint256)
func (_RootNet *RootNetCaller) GetActiveSubnetIdAt(opts *bind.CallOpts, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "getActiveSubnetIdAt", index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActiveSubnetIdAt is a free data retrieval call binding the contract method 0x38f48a89.
//
// Solidity: function getActiveSubnetIdAt(uint256 index) view returns(uint256)
func (_RootNet *RootNetSession) GetActiveSubnetIdAt(index *big.Int) (*big.Int, error) {
	return _RootNet.Contract.GetActiveSubnetIdAt(&_RootNet.CallOpts, index)
}

// GetActiveSubnetIdAt is a free data retrieval call binding the contract method 0x38f48a89.
//
// Solidity: function getActiveSubnetIdAt(uint256 index) view returns(uint256)
func (_RootNet *RootNetCallerSession) GetActiveSubnetIdAt(index *big.Int) (*big.Int, error) {
	return _RootNet.Contract.GetActiveSubnetIdAt(&_RootNet.CallOpts, index)
}

// GetAgentInfo is a free data retrieval call binding the contract method 0x168f80f5.
//
// Solidity: function getAgentInfo(address agent, uint256 subnetId) view returns((address,bool,uint256,address))
func (_RootNet *RootNetCaller) GetAgentInfo(opts *bind.CallOpts, agent common.Address, subnetId *big.Int) (RootNetAgentInfo, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "getAgentInfo", agent, subnetId)

	if err != nil {
		return *new(RootNetAgentInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(RootNetAgentInfo)).(*RootNetAgentInfo)

	return out0, err

}

// GetAgentInfo is a free data retrieval call binding the contract method 0x168f80f5.
//
// Solidity: function getAgentInfo(address agent, uint256 subnetId) view returns((address,bool,uint256,address))
func (_RootNet *RootNetSession) GetAgentInfo(agent common.Address, subnetId *big.Int) (RootNetAgentInfo, error) {
	return _RootNet.Contract.GetAgentInfo(&_RootNet.CallOpts, agent, subnetId)
}

// GetAgentInfo is a free data retrieval call binding the contract method 0x168f80f5.
//
// Solidity: function getAgentInfo(address agent, uint256 subnetId) view returns((address,bool,uint256,address))
func (_RootNet *RootNetCallerSession) GetAgentInfo(agent common.Address, subnetId *big.Int) (RootNetAgentInfo, error) {
	return _RootNet.Contract.GetAgentInfo(&_RootNet.CallOpts, agent, subnetId)
}

// GetAgentsInfo is a free data retrieval call binding the contract method 0x4b6f6d67.
//
// Solidity: function getAgentsInfo(address[] agents, uint256 subnetId) view returns((address,bool,uint256,address)[])
func (_RootNet *RootNetCaller) GetAgentsInfo(opts *bind.CallOpts, agents []common.Address, subnetId *big.Int) ([]RootNetAgentInfo, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "getAgentsInfo", agents, subnetId)

	if err != nil {
		return *new([]RootNetAgentInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]RootNetAgentInfo)).(*[]RootNetAgentInfo)

	return out0, err

}

// GetAgentsInfo is a free data retrieval call binding the contract method 0x4b6f6d67.
//
// Solidity: function getAgentsInfo(address[] agents, uint256 subnetId) view returns((address,bool,uint256,address)[])
func (_RootNet *RootNetSession) GetAgentsInfo(agents []common.Address, subnetId *big.Int) ([]RootNetAgentInfo, error) {
	return _RootNet.Contract.GetAgentsInfo(&_RootNet.CallOpts, agents, subnetId)
}

// GetAgentsInfo is a free data retrieval call binding the contract method 0x4b6f6d67.
//
// Solidity: function getAgentsInfo(address[] agents, uint256 subnetId) view returns((address,bool,uint256,address)[])
func (_RootNet *RootNetCallerSession) GetAgentsInfo(agents []common.Address, subnetId *big.Int) ([]RootNetAgentInfo, error) {
	return _RootNet.Contract.GetAgentsInfo(&_RootNet.CallOpts, agents, subnetId)
}

// GetRegistry is a free data retrieval call binding the contract method 0x5ab1bd53.
//
// Solidity: function getRegistry() view returns(address, address, address, address, address, address, address, address, address, address)
func (_RootNet *RootNetCaller) GetRegistry(opts *bind.CallOpts) (common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "getRegistry")

	if err != nil {
		return *new(common.Address), *new(common.Address), *new(common.Address), *new(common.Address), *new(common.Address), *new(common.Address), *new(common.Address), *new(common.Address), *new(common.Address), *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	out2 := *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	out3 := *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	out4 := *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	out5 := *abi.ConvertType(out[5], new(common.Address)).(*common.Address)
	out6 := *abi.ConvertType(out[6], new(common.Address)).(*common.Address)
	out7 := *abi.ConvertType(out[7], new(common.Address)).(*common.Address)
	out8 := *abi.ConvertType(out[8], new(common.Address)).(*common.Address)
	out9 := *abi.ConvertType(out[9], new(common.Address)).(*common.Address)

	return out0, out1, out2, out3, out4, out5, out6, out7, out8, out9, err

}

// GetRegistry is a free data retrieval call binding the contract method 0x5ab1bd53.
//
// Solidity: function getRegistry() view returns(address, address, address, address, address, address, address, address, address, address)
func (_RootNet *RootNetSession) GetRegistry() (common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, error) {
	return _RootNet.Contract.GetRegistry(&_RootNet.CallOpts)
}

// GetRegistry is a free data retrieval call binding the contract method 0x5ab1bd53.
//
// Solidity: function getRegistry() view returns(address, address, address, address, address, address, address, address, address, address)
func (_RootNet *RootNetCallerSession) GetRegistry() (common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, error) {
	return _RootNet.Contract.GetRegistry(&_RootNet.CallOpts)
}

// GetSubnet is a free data retrieval call binding the contract method 0x58ca7504.
//
// Solidity: function getSubnet(uint256 subnetId) view returns((bytes32,uint8,uint64,uint64))
func (_RootNet *RootNetCaller) GetSubnet(opts *bind.CallOpts, subnetId *big.Int) (IRootNetSubnetInfo, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "getSubnet", subnetId)

	if err != nil {
		return *new(IRootNetSubnetInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IRootNetSubnetInfo)).(*IRootNetSubnetInfo)

	return out0, err

}

// GetSubnet is a free data retrieval call binding the contract method 0x58ca7504.
//
// Solidity: function getSubnet(uint256 subnetId) view returns((bytes32,uint8,uint64,uint64))
func (_RootNet *RootNetSession) GetSubnet(subnetId *big.Int) (IRootNetSubnetInfo, error) {
	return _RootNet.Contract.GetSubnet(&_RootNet.CallOpts, subnetId)
}

// GetSubnet is a free data retrieval call binding the contract method 0x58ca7504.
//
// Solidity: function getSubnet(uint256 subnetId) view returns((bytes32,uint8,uint64,uint64))
func (_RootNet *RootNetCallerSession) GetSubnet(subnetId *big.Int) (IRootNetSubnetInfo, error) {
	return _RootNet.Contract.GetSubnet(&_RootNet.CallOpts, subnetId)
}

// GetSubnetFull is a free data retrieval call binding the contract method 0x3de3b247.
//
// Solidity: function getSubnetFull(uint256 subnetId) view returns((address,address,bytes32,uint8,uint64,uint64,string,string,uint128,address))
func (_RootNet *RootNetCaller) GetSubnetFull(opts *bind.CallOpts, subnetId *big.Int) (IRootNetSubnetFullInfo, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "getSubnetFull", subnetId)

	if err != nil {
		return *new(IRootNetSubnetFullInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IRootNetSubnetFullInfo)).(*IRootNetSubnetFullInfo)

	return out0, err

}

// GetSubnetFull is a free data retrieval call binding the contract method 0x3de3b247.
//
// Solidity: function getSubnetFull(uint256 subnetId) view returns((address,address,bytes32,uint8,uint64,uint64,string,string,uint128,address))
func (_RootNet *RootNetSession) GetSubnetFull(subnetId *big.Int) (IRootNetSubnetFullInfo, error) {
	return _RootNet.Contract.GetSubnetFull(&_RootNet.CallOpts, subnetId)
}

// GetSubnetFull is a free data retrieval call binding the contract method 0x3de3b247.
//
// Solidity: function getSubnetFull(uint256 subnetId) view returns((address,address,bytes32,uint8,uint64,uint64,string,string,uint128,address))
func (_RootNet *RootNetCallerSession) GetSubnetFull(subnetId *big.Int) (IRootNetSubnetFullInfo, error) {
	return _RootNet.Contract.GetSubnetFull(&_RootNet.CallOpts, subnetId)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_RootNet *RootNetCaller) Guardian(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "guardian")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_RootNet *RootNetSession) Guardian() (common.Address, error) {
	return _RootNet.Contract.Guardian(&_RootNet.CallOpts)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_RootNet *RootNetCallerSession) Guardian() (common.Address, error) {
	return _RootNet.Contract.Guardian(&_RootNet.CallOpts)
}

// ImmunityPeriod is a free data retrieval call binding the contract method 0x2672e1be.
//
// Solidity: function immunityPeriod() view returns(uint256)
func (_RootNet *RootNetCaller) ImmunityPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "immunityPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ImmunityPeriod is a free data retrieval call binding the contract method 0x2672e1be.
//
// Solidity: function immunityPeriod() view returns(uint256)
func (_RootNet *RootNetSession) ImmunityPeriod() (*big.Int, error) {
	return _RootNet.Contract.ImmunityPeriod(&_RootNet.CallOpts)
}

// ImmunityPeriod is a free data retrieval call binding the contract method 0x2672e1be.
//
// Solidity: function immunityPeriod() view returns(uint256)
func (_RootNet *RootNetCallerSession) ImmunityPeriod() (*big.Int, error) {
	return _RootNet.Contract.ImmunityPeriod(&_RootNet.CallOpts)
}

// InitialAlphaPrice is a free data retrieval call binding the contract method 0x6d345eea.
//
// Solidity: function initialAlphaPrice() view returns(uint256)
func (_RootNet *RootNetCaller) InitialAlphaPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "initialAlphaPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InitialAlphaPrice is a free data retrieval call binding the contract method 0x6d345eea.
//
// Solidity: function initialAlphaPrice() view returns(uint256)
func (_RootNet *RootNetSession) InitialAlphaPrice() (*big.Int, error) {
	return _RootNet.Contract.InitialAlphaPrice(&_RootNet.CallOpts)
}

// InitialAlphaPrice is a free data retrieval call binding the contract method 0x6d345eea.
//
// Solidity: function initialAlphaPrice() view returns(uint256)
func (_RootNet *RootNetCallerSession) InitialAlphaPrice() (*big.Int, error) {
	return _RootNet.Contract.InitialAlphaPrice(&_RootNet.CallOpts)
}

// IsSubnetActive is a free data retrieval call binding the contract method 0x7ab5e276.
//
// Solidity: function isSubnetActive(uint256 subnetId) view returns(bool)
func (_RootNet *RootNetCaller) IsSubnetActive(opts *bind.CallOpts, subnetId *big.Int) (bool, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "isSubnetActive", subnetId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsSubnetActive is a free data retrieval call binding the contract method 0x7ab5e276.
//
// Solidity: function isSubnetActive(uint256 subnetId) view returns(bool)
func (_RootNet *RootNetSession) IsSubnetActive(subnetId *big.Int) (bool, error) {
	return _RootNet.Contract.IsSubnetActive(&_RootNet.CallOpts, subnetId)
}

// IsSubnetActive is a free data retrieval call binding the contract method 0x7ab5e276.
//
// Solidity: function isSubnetActive(uint256 subnetId) view returns(bool)
func (_RootNet *RootNetCallerSession) IsSubnetActive(subnetId *big.Int) (bool, error) {
	return _RootNet.Contract.IsSubnetActive(&_RootNet.CallOpts, subnetId)
}

// LpManager is a free data retrieval call binding the contract method 0xb906f15a.
//
// Solidity: function lpManager() view returns(address)
func (_RootNet *RootNetCaller) LpManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "lpManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LpManager is a free data retrieval call binding the contract method 0xb906f15a.
//
// Solidity: function lpManager() view returns(address)
func (_RootNet *RootNetSession) LpManager() (common.Address, error) {
	return _RootNet.Contract.LpManager(&_RootNet.CallOpts)
}

// LpManager is a free data retrieval call binding the contract method 0xb906f15a.
//
// Solidity: function lpManager() view returns(address)
func (_RootNet *RootNetCallerSession) LpManager() (common.Address, error) {
	return _RootNet.Contract.LpManager(&_RootNet.CallOpts)
}

// NextSubnetId is a free data retrieval call binding the contract method 0xd929ff05.
//
// Solidity: function nextSubnetId() view returns(uint256)
func (_RootNet *RootNetCaller) NextSubnetId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "nextSubnetId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextSubnetId is a free data retrieval call binding the contract method 0xd929ff05.
//
// Solidity: function nextSubnetId() view returns(uint256)
func (_RootNet *RootNetSession) NextSubnetId() (*big.Int, error) {
	return _RootNet.Contract.NextSubnetId(&_RootNet.CallOpts)
}

// NextSubnetId is a free data retrieval call binding the contract method 0xd929ff05.
//
// Solidity: function nextSubnetId() view returns(uint256)
func (_RootNet *RootNetCallerSession) NextSubnetId() (*big.Int, error) {
	return _RootNet.Contract.NextSubnetId(&_RootNet.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_RootNet *RootNetCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "nonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_RootNet *RootNetSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _RootNet.Contract.Nonces(&_RootNet.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_RootNet *RootNetCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _RootNet.Contract.Nonces(&_RootNet.CallOpts, arg0)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RootNet *RootNetCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RootNet *RootNetSession) Paused() (bool, error) {
	return _RootNet.Contract.Paused(&_RootNet.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_RootNet *RootNetCallerSession) Paused() (bool, error) {
	return _RootNet.Contract.Paused(&_RootNet.CallOpts)
}

// RegistryInitialized is a free data retrieval call binding the contract method 0x56354a24.
//
// Solidity: function registryInitialized() view returns(bool)
func (_RootNet *RootNetCaller) RegistryInitialized(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "registryInitialized")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// RegistryInitialized is a free data retrieval call binding the contract method 0x56354a24.
//
// Solidity: function registryInitialized() view returns(bool)
func (_RootNet *RootNetSession) RegistryInitialized() (bool, error) {
	return _RootNet.Contract.RegistryInitialized(&_RootNet.CallOpts)
}

// RegistryInitialized is a free data retrieval call binding the contract method 0x56354a24.
//
// Solidity: function registryInitialized() view returns(bool)
func (_RootNet *RootNetCallerSession) RegistryInitialized() (bool, error) {
	return _RootNet.Contract.RegistryInitialized(&_RootNet.CallOpts)
}

// StakeNFT is a free data retrieval call binding the contract method 0xb48509e6.
//
// Solidity: function stakeNFT() view returns(address)
func (_RootNet *RootNetCaller) StakeNFT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "stakeNFT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakeNFT is a free data retrieval call binding the contract method 0xb48509e6.
//
// Solidity: function stakeNFT() view returns(address)
func (_RootNet *RootNetSession) StakeNFT() (common.Address, error) {
	return _RootNet.Contract.StakeNFT(&_RootNet.CallOpts)
}

// StakeNFT is a free data retrieval call binding the contract method 0xb48509e6.
//
// Solidity: function stakeNFT() view returns(address)
func (_RootNet *RootNetCallerSession) StakeNFT() (common.Address, error) {
	return _RootNet.Contract.StakeNFT(&_RootNet.CallOpts)
}

// StakingVault is a free data retrieval call binding the contract method 0x24e7964a.
//
// Solidity: function stakingVault() view returns(address)
func (_RootNet *RootNetCaller) StakingVault(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "stakingVault")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakingVault is a free data retrieval call binding the contract method 0x24e7964a.
//
// Solidity: function stakingVault() view returns(address)
func (_RootNet *RootNetSession) StakingVault() (common.Address, error) {
	return _RootNet.Contract.StakingVault(&_RootNet.CallOpts)
}

// StakingVault is a free data retrieval call binding the contract method 0x24e7964a.
//
// Solidity: function stakingVault() view returns(address)
func (_RootNet *RootNetCallerSession) StakingVault() (common.Address, error) {
	return _RootNet.Contract.StakingVault(&_RootNet.CallOpts)
}

// SubnetNFT is a free data retrieval call binding the contract method 0x11cba7e9.
//
// Solidity: function subnetNFT() view returns(address)
func (_RootNet *RootNetCaller) SubnetNFT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "subnetNFT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SubnetNFT is a free data retrieval call binding the contract method 0x11cba7e9.
//
// Solidity: function subnetNFT() view returns(address)
func (_RootNet *RootNetSession) SubnetNFT() (common.Address, error) {
	return _RootNet.Contract.SubnetNFT(&_RootNet.CallOpts)
}

// SubnetNFT is a free data retrieval call binding the contract method 0x11cba7e9.
//
// Solidity: function subnetNFT() view returns(address)
func (_RootNet *RootNetCallerSession) SubnetNFT() (common.Address, error) {
	return _RootNet.Contract.SubnetNFT(&_RootNet.CallOpts)
}

// Subnets is a free data retrieval call binding the contract method 0x475726f7.
//
// Solidity: function subnets(uint256 ) view returns(bytes32 lpPool, uint8 status, uint64 createdAt, uint64 activatedAt)
func (_RootNet *RootNetCaller) Subnets(opts *bind.CallOpts, arg0 *big.Int) (struct {
	LpPool      [32]byte
	Status      uint8
	CreatedAt   uint64
	ActivatedAt uint64
}, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "subnets", arg0)

	outstruct := new(struct {
		LpPool      [32]byte
		Status      uint8
		CreatedAt   uint64
		ActivatedAt uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.LpPool = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Status = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.CreatedAt = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.ActivatedAt = *abi.ConvertType(out[3], new(uint64)).(*uint64)

	return *outstruct, err

}

// Subnets is a free data retrieval call binding the contract method 0x475726f7.
//
// Solidity: function subnets(uint256 ) view returns(bytes32 lpPool, uint8 status, uint64 createdAt, uint64 activatedAt)
func (_RootNet *RootNetSession) Subnets(arg0 *big.Int) (struct {
	LpPool      [32]byte
	Status      uint8
	CreatedAt   uint64
	ActivatedAt uint64
}, error) {
	return _RootNet.Contract.Subnets(&_RootNet.CallOpts, arg0)
}

// Subnets is a free data retrieval call binding the contract method 0x475726f7.
//
// Solidity: function subnets(uint256 ) view returns(bytes32 lpPool, uint8 status, uint64 createdAt, uint64 activatedAt)
func (_RootNet *RootNetCallerSession) Subnets(arg0 *big.Int) (struct {
	LpPool      [32]byte
	Status      uint8
	CreatedAt   uint64
	ActivatedAt uint64
}, error) {
	return _RootNet.Contract.Subnets(&_RootNet.CallOpts, arg0)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_RootNet *RootNetCaller) Treasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RootNet.contract.Call(opts, &out, "treasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_RootNet *RootNetSession) Treasury() (common.Address, error) {
	return _RootNet.Contract.Treasury(&_RootNet.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_RootNet *RootNetCallerSession) Treasury() (common.Address, error) {
	return _RootNet.Contract.Treasury(&_RootNet.CallOpts)
}

// ActivateSubnet is a paid mutator transaction binding the contract method 0xcead1c96.
//
// Solidity: function activateSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetTransactor) ActivateSubnet(opts *bind.TransactOpts, subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "activateSubnet", subnetId)
}

// ActivateSubnet is a paid mutator transaction binding the contract method 0xcead1c96.
//
// Solidity: function activateSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetSession) ActivateSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.ActivateSubnet(&_RootNet.TransactOpts, subnetId)
}

// ActivateSubnet is a paid mutator transaction binding the contract method 0xcead1c96.
//
// Solidity: function activateSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetTransactorSession) ActivateSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.ActivateSubnet(&_RootNet.TransactOpts, subnetId)
}

// Allocate is a paid mutator transaction binding the contract method 0xab3f22d5.
//
// Solidity: function allocate(address agent, uint256 subnetId, uint256 amount) returns()
func (_RootNet *RootNetTransactor) Allocate(opts *bind.TransactOpts, agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "allocate", agent, subnetId, amount)
}

// Allocate is a paid mutator transaction binding the contract method 0xab3f22d5.
//
// Solidity: function allocate(address agent, uint256 subnetId, uint256 amount) returns()
func (_RootNet *RootNetSession) Allocate(agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.Allocate(&_RootNet.TransactOpts, agent, subnetId, amount)
}

// Allocate is a paid mutator transaction binding the contract method 0xab3f22d5.
//
// Solidity: function allocate(address agent, uint256 subnetId, uint256 amount) returns()
func (_RootNet *RootNetTransactorSession) Allocate(agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.Allocate(&_RootNet.TransactOpts, agent, subnetId, amount)
}

// BanSubnet is a paid mutator transaction binding the contract method 0xb79b7658.
//
// Solidity: function banSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetTransactor) BanSubnet(opts *bind.TransactOpts, subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "banSubnet", subnetId)
}

// BanSubnet is a paid mutator transaction binding the contract method 0xb79b7658.
//
// Solidity: function banSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetSession) BanSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.BanSubnet(&_RootNet.TransactOpts, subnetId)
}

// BanSubnet is a paid mutator transaction binding the contract method 0xb79b7658.
//
// Solidity: function banSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetTransactorSession) BanSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.BanSubnet(&_RootNet.TransactOpts, subnetId)
}

// Bind is a paid mutator transaction binding the contract method 0x81bac14f.
//
// Solidity: function bind(address principal) returns()
func (_RootNet *RootNetTransactor) Bind(opts *bind.TransactOpts, principal common.Address) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "bind", principal)
}

// Bind is a paid mutator transaction binding the contract method 0x81bac14f.
//
// Solidity: function bind(address principal) returns()
func (_RootNet *RootNetSession) Bind(principal common.Address) (*types.Transaction, error) {
	return _RootNet.Contract.Bind(&_RootNet.TransactOpts, principal)
}

// Bind is a paid mutator transaction binding the contract method 0x81bac14f.
//
// Solidity: function bind(address principal) returns()
func (_RootNet *RootNetTransactorSession) Bind(principal common.Address) (*types.Transaction, error) {
	return _RootNet.Contract.Bind(&_RootNet.TransactOpts, principal)
}

// BindFor is a paid mutator transaction binding the contract method 0x7b234b81.
//
// Solidity: function bindFor(address agent, address principal, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetTransactor) BindFor(opts *bind.TransactOpts, agent common.Address, principal common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "bindFor", agent, principal, deadline, v, r, s)
}

// BindFor is a paid mutator transaction binding the contract method 0x7b234b81.
//
// Solidity: function bindFor(address agent, address principal, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetSession) BindFor(agent common.Address, principal common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.BindFor(&_RootNet.TransactOpts, agent, principal, deadline, v, r, s)
}

// BindFor is a paid mutator transaction binding the contract method 0x7b234b81.
//
// Solidity: function bindFor(address agent, address principal, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetTransactorSession) BindFor(agent common.Address, principal common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.BindFor(&_RootNet.TransactOpts, agent, principal, deadline, v, r, s)
}

// Deallocate is a paid mutator transaction binding the contract method 0xfe427e95.
//
// Solidity: function deallocate(address agent, uint256 subnetId, uint256 amount) returns()
func (_RootNet *RootNetTransactor) Deallocate(opts *bind.TransactOpts, agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "deallocate", agent, subnetId, amount)
}

// Deallocate is a paid mutator transaction binding the contract method 0xfe427e95.
//
// Solidity: function deallocate(address agent, uint256 subnetId, uint256 amount) returns()
func (_RootNet *RootNetSession) Deallocate(agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.Deallocate(&_RootNet.TransactOpts, agent, subnetId, amount)
}

// Deallocate is a paid mutator transaction binding the contract method 0xfe427e95.
//
// Solidity: function deallocate(address agent, uint256 subnetId, uint256 amount) returns()
func (_RootNet *RootNetTransactorSession) Deallocate(agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.Deallocate(&_RootNet.TransactOpts, agent, subnetId, amount)
}

// DeregisterSubnet is a paid mutator transaction binding the contract method 0x0cf02c5e.
//
// Solidity: function deregisterSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetTransactor) DeregisterSubnet(opts *bind.TransactOpts, subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "deregisterSubnet", subnetId)
}

// DeregisterSubnet is a paid mutator transaction binding the contract method 0x0cf02c5e.
//
// Solidity: function deregisterSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetSession) DeregisterSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.DeregisterSubnet(&_RootNet.TransactOpts, subnetId)
}

// DeregisterSubnet is a paid mutator transaction binding the contract method 0x0cf02c5e.
//
// Solidity: function deregisterSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetTransactorSession) DeregisterSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.DeregisterSubnet(&_RootNet.TransactOpts, subnetId)
}

// InitializeRegistry is a paid mutator transaction binding the contract method 0xc8c32a84.
//
// Solidity: function initializeRegistry(address awpToken_, address subnetNFT_, address alphaTokenFactory_, address awpEmission_, address lpManager_, address accessManager_, address stakingVault_, address stakeNFT_) returns()
func (_RootNet *RootNetTransactor) InitializeRegistry(opts *bind.TransactOpts, awpToken_ common.Address, subnetNFT_ common.Address, alphaTokenFactory_ common.Address, awpEmission_ common.Address, lpManager_ common.Address, accessManager_ common.Address, stakingVault_ common.Address, stakeNFT_ common.Address) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "initializeRegistry", awpToken_, subnetNFT_, alphaTokenFactory_, awpEmission_, lpManager_, accessManager_, stakingVault_, stakeNFT_)
}

// InitializeRegistry is a paid mutator transaction binding the contract method 0xc8c32a84.
//
// Solidity: function initializeRegistry(address awpToken_, address subnetNFT_, address alphaTokenFactory_, address awpEmission_, address lpManager_, address accessManager_, address stakingVault_, address stakeNFT_) returns()
func (_RootNet *RootNetSession) InitializeRegistry(awpToken_ common.Address, subnetNFT_ common.Address, alphaTokenFactory_ common.Address, awpEmission_ common.Address, lpManager_ common.Address, accessManager_ common.Address, stakingVault_ common.Address, stakeNFT_ common.Address) (*types.Transaction, error) {
	return _RootNet.Contract.InitializeRegistry(&_RootNet.TransactOpts, awpToken_, subnetNFT_, alphaTokenFactory_, awpEmission_, lpManager_, accessManager_, stakingVault_, stakeNFT_)
}

// InitializeRegistry is a paid mutator transaction binding the contract method 0xc8c32a84.
//
// Solidity: function initializeRegistry(address awpToken_, address subnetNFT_, address alphaTokenFactory_, address awpEmission_, address lpManager_, address accessManager_, address stakingVault_, address stakeNFT_) returns()
func (_RootNet *RootNetTransactorSession) InitializeRegistry(awpToken_ common.Address, subnetNFT_ common.Address, alphaTokenFactory_ common.Address, awpEmission_ common.Address, lpManager_ common.Address, accessManager_ common.Address, stakingVault_ common.Address, stakeNFT_ common.Address) (*types.Transaction, error) {
	return _RootNet.Contract.InitializeRegistry(&_RootNet.TransactOpts, awpToken_, subnetNFT_, alphaTokenFactory_, awpEmission_, lpManager_, accessManager_, stakingVault_, stakeNFT_)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RootNet *RootNetTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RootNet *RootNetSession) Pause() (*types.Transaction, error) {
	return _RootNet.Contract.Pause(&_RootNet.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_RootNet *RootNetTransactorSession) Pause() (*types.Transaction, error) {
	return _RootNet.Contract.Pause(&_RootNet.TransactOpts)
}

// PauseSubnet is a paid mutator transaction binding the contract method 0x44e047ca.
//
// Solidity: function pauseSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetTransactor) PauseSubnet(opts *bind.TransactOpts, subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "pauseSubnet", subnetId)
}

// PauseSubnet is a paid mutator transaction binding the contract method 0x44e047ca.
//
// Solidity: function pauseSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetSession) PauseSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.PauseSubnet(&_RootNet.TransactOpts, subnetId)
}

// PauseSubnet is a paid mutator transaction binding the contract method 0x44e047ca.
//
// Solidity: function pauseSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetTransactorSession) PauseSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.PauseSubnet(&_RootNet.TransactOpts, subnetId)
}

// Reallocate is a paid mutator transaction binding the contract method 0x1a46f4b8.
//
// Solidity: function reallocate(address fromAgent, uint256 fromSubnetId, address toAgent, uint256 toSubnetId, uint256 amount) returns()
func (_RootNet *RootNetTransactor) Reallocate(opts *bind.TransactOpts, fromAgent common.Address, fromSubnetId *big.Int, toAgent common.Address, toSubnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "reallocate", fromAgent, fromSubnetId, toAgent, toSubnetId, amount)
}

// Reallocate is a paid mutator transaction binding the contract method 0x1a46f4b8.
//
// Solidity: function reallocate(address fromAgent, uint256 fromSubnetId, address toAgent, uint256 toSubnetId, uint256 amount) returns()
func (_RootNet *RootNetSession) Reallocate(fromAgent common.Address, fromSubnetId *big.Int, toAgent common.Address, toSubnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.Reallocate(&_RootNet.TransactOpts, fromAgent, fromSubnetId, toAgent, toSubnetId, amount)
}

// Reallocate is a paid mutator transaction binding the contract method 0x1a46f4b8.
//
// Solidity: function reallocate(address fromAgent, uint256 fromSubnetId, address toAgent, uint256 toSubnetId, uint256 amount) returns()
func (_RootNet *RootNetTransactorSession) Reallocate(fromAgent common.Address, fromSubnetId *big.Int, toAgent common.Address, toSubnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.Reallocate(&_RootNet.TransactOpts, fromAgent, fromSubnetId, toAgent, toSubnetId, amount)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_RootNet *RootNetTransactor) Register(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "register")
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_RootNet *RootNetSession) Register() (*types.Transaction, error) {
	return _RootNet.Contract.Register(&_RootNet.TransactOpts)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_RootNet *RootNetTransactorSession) Register() (*types.Transaction, error) {
	return _RootNet.Contract.Register(&_RootNet.TransactOpts)
}

// Register0 is a paid mutator transaction binding the contract method 0x6d23f895.
//
// Solidity: function register(address recipient, uint256 depositAmount, uint64 lockDuration) returns()
func (_RootNet *RootNetTransactor) Register0(opts *bind.TransactOpts, recipient common.Address, depositAmount *big.Int, lockDuration uint64) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "register0", recipient, depositAmount, lockDuration)
}

// Register0 is a paid mutator transaction binding the contract method 0x6d23f895.
//
// Solidity: function register(address recipient, uint256 depositAmount, uint64 lockDuration) returns()
func (_RootNet *RootNetSession) Register0(recipient common.Address, depositAmount *big.Int, lockDuration uint64) (*types.Transaction, error) {
	return _RootNet.Contract.Register0(&_RootNet.TransactOpts, recipient, depositAmount, lockDuration)
}

// Register0 is a paid mutator transaction binding the contract method 0x6d23f895.
//
// Solidity: function register(address recipient, uint256 depositAmount, uint64 lockDuration) returns()
func (_RootNet *RootNetTransactorSession) Register0(recipient common.Address, depositAmount *big.Int, lockDuration uint64) (*types.Transaction, error) {
	return _RootNet.Contract.Register0(&_RootNet.TransactOpts, recipient, depositAmount, lockDuration)
}

// RegisterAndStake is a paid mutator transaction binding the contract method 0x34426564.
//
// Solidity: function registerAndStake(uint256 depositAmount, uint64 lockDuration, address agent, uint256 subnetId, uint256 allocateAmount) returns()
func (_RootNet *RootNetTransactor) RegisterAndStake(opts *bind.TransactOpts, depositAmount *big.Int, lockDuration uint64, agent common.Address, subnetId *big.Int, allocateAmount *big.Int) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "registerAndStake", depositAmount, lockDuration, agent, subnetId, allocateAmount)
}

// RegisterAndStake is a paid mutator transaction binding the contract method 0x34426564.
//
// Solidity: function registerAndStake(uint256 depositAmount, uint64 lockDuration, address agent, uint256 subnetId, uint256 allocateAmount) returns()
func (_RootNet *RootNetSession) RegisterAndStake(depositAmount *big.Int, lockDuration uint64, agent common.Address, subnetId *big.Int, allocateAmount *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.RegisterAndStake(&_RootNet.TransactOpts, depositAmount, lockDuration, agent, subnetId, allocateAmount)
}

// RegisterAndStake is a paid mutator transaction binding the contract method 0x34426564.
//
// Solidity: function registerAndStake(uint256 depositAmount, uint64 lockDuration, address agent, uint256 subnetId, uint256 allocateAmount) returns()
func (_RootNet *RootNetTransactorSession) RegisterAndStake(depositAmount *big.Int, lockDuration uint64, agent common.Address, subnetId *big.Int, allocateAmount *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.RegisterAndStake(&_RootNet.TransactOpts, depositAmount, lockDuration, agent, subnetId, allocateAmount)
}

// RegisterFor is a paid mutator transaction binding the contract method 0x671a2a8a.
//
// Solidity: function registerFor(address user, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetTransactor) RegisterFor(opts *bind.TransactOpts, user common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "registerFor", user, deadline, v, r, s)
}

// RegisterFor is a paid mutator transaction binding the contract method 0x671a2a8a.
//
// Solidity: function registerFor(address user, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetSession) RegisterFor(user common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.RegisterFor(&_RootNet.TransactOpts, user, deadline, v, r, s)
}

// RegisterFor is a paid mutator transaction binding the contract method 0x671a2a8a.
//
// Solidity: function registerFor(address user, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetTransactorSession) RegisterFor(user common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.RegisterFor(&_RootNet.TransactOpts, user, deadline, v, r, s)
}

// RegisterSubnet is a paid mutator transaction binding the contract method 0xcd49dc03.
//
// Solidity: function registerSubnet((string,string,string,address,string,bytes32,uint128) params) returns(uint256)
func (_RootNet *RootNetTransactor) RegisterSubnet(opts *bind.TransactOpts, params IRootNetSubnetParams) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "registerSubnet", params)
}

// RegisterSubnet is a paid mutator transaction binding the contract method 0xcd49dc03.
//
// Solidity: function registerSubnet((string,string,string,address,string,bytes32,uint128) params) returns(uint256)
func (_RootNet *RootNetSession) RegisterSubnet(params IRootNetSubnetParams) (*types.Transaction, error) {
	return _RootNet.Contract.RegisterSubnet(&_RootNet.TransactOpts, params)
}

// RegisterSubnet is a paid mutator transaction binding the contract method 0xcd49dc03.
//
// Solidity: function registerSubnet((string,string,string,address,string,bytes32,uint128) params) returns(uint256)
func (_RootNet *RootNetTransactorSession) RegisterSubnet(params IRootNetSubnetParams) (*types.Transaction, error) {
	return _RootNet.Contract.RegisterSubnet(&_RootNet.TransactOpts, params)
}

// RegisterSubnetFor is a paid mutator transaction binding the contract method 0x080ec841.
//
// Solidity: function registerSubnetFor(address user, (string,string,string,address,string,bytes32,uint128) params, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_RootNet *RootNetTransactor) RegisterSubnetFor(opts *bind.TransactOpts, user common.Address, params IRootNetSubnetParams, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "registerSubnetFor", user, params, deadline, v, r, s)
}

// RegisterSubnetFor is a paid mutator transaction binding the contract method 0x080ec841.
//
// Solidity: function registerSubnetFor(address user, (string,string,string,address,string,bytes32,uint128) params, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_RootNet *RootNetSession) RegisterSubnetFor(user common.Address, params IRootNetSubnetParams, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.RegisterSubnetFor(&_RootNet.TransactOpts, user, params, deadline, v, r, s)
}

// RegisterSubnetFor is a paid mutator transaction binding the contract method 0x080ec841.
//
// Solidity: function registerSubnetFor(address user, (string,string,string,address,string,bytes32,uint128) params, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_RootNet *RootNetTransactorSession) RegisterSubnetFor(user common.Address, params IRootNetSubnetParams, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.RegisterSubnetFor(&_RootNet.TransactOpts, user, params, deadline, v, r, s)
}

// RegisterSubnetForWithPermit is a paid mutator transaction binding the contract method 0xe3da3f9b.
//
// Solidity: function registerSubnetForWithPermit(address user, (string,string,string,address,string,bytes32,uint128) params, uint256 deadline, uint8 permitV, bytes32 permitR, bytes32 permitS, uint8 registerV, bytes32 registerR, bytes32 registerS) returns(uint256)
func (_RootNet *RootNetTransactor) RegisterSubnetForWithPermit(opts *bind.TransactOpts, user common.Address, params IRootNetSubnetParams, deadline *big.Int, permitV uint8, permitR [32]byte, permitS [32]byte, registerV uint8, registerR [32]byte, registerS [32]byte) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "registerSubnetForWithPermit", user, params, deadline, permitV, permitR, permitS, registerV, registerR, registerS)
}

// RegisterSubnetForWithPermit is a paid mutator transaction binding the contract method 0xe3da3f9b.
//
// Solidity: function registerSubnetForWithPermit(address user, (string,string,string,address,string,bytes32,uint128) params, uint256 deadline, uint8 permitV, bytes32 permitR, bytes32 permitS, uint8 registerV, bytes32 registerR, bytes32 registerS) returns(uint256)
func (_RootNet *RootNetSession) RegisterSubnetForWithPermit(user common.Address, params IRootNetSubnetParams, deadline *big.Int, permitV uint8, permitR [32]byte, permitS [32]byte, registerV uint8, registerR [32]byte, registerS [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.RegisterSubnetForWithPermit(&_RootNet.TransactOpts, user, params, deadline, permitV, permitR, permitS, registerV, registerR, registerS)
}

// RegisterSubnetForWithPermit is a paid mutator transaction binding the contract method 0xe3da3f9b.
//
// Solidity: function registerSubnetForWithPermit(address user, (string,string,string,address,string,bytes32,uint128) params, uint256 deadline, uint8 permitV, bytes32 permitR, bytes32 permitS, uint8 registerV, bytes32 registerR, bytes32 registerS) returns(uint256)
func (_RootNet *RootNetTransactorSession) RegisterSubnetForWithPermit(user common.Address, params IRootNetSubnetParams, deadline *big.Int, permitV uint8, permitR [32]byte, permitS [32]byte, registerV uint8, registerR [32]byte, registerS [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.RegisterSubnetForWithPermit(&_RootNet.TransactOpts, user, params, deadline, permitV, permitR, permitS, registerV, registerR, registerS)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x97a6278e.
//
// Solidity: function removeAgent(address agent) returns()
func (_RootNet *RootNetTransactor) RemoveAgent(opts *bind.TransactOpts, agent common.Address) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "removeAgent", agent)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x97a6278e.
//
// Solidity: function removeAgent(address agent) returns()
func (_RootNet *RootNetSession) RemoveAgent(agent common.Address) (*types.Transaction, error) {
	return _RootNet.Contract.RemoveAgent(&_RootNet.TransactOpts, agent)
}

// RemoveAgent is a paid mutator transaction binding the contract method 0x97a6278e.
//
// Solidity: function removeAgent(address agent) returns()
func (_RootNet *RootNetTransactorSession) RemoveAgent(agent common.Address) (*types.Transaction, error) {
	return _RootNet.Contract.RemoveAgent(&_RootNet.TransactOpts, agent)
}

// ResumeSubnet is a paid mutator transaction binding the contract method 0x5364944c.
//
// Solidity: function resumeSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetTransactor) ResumeSubnet(opts *bind.TransactOpts, subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "resumeSubnet", subnetId)
}

// ResumeSubnet is a paid mutator transaction binding the contract method 0x5364944c.
//
// Solidity: function resumeSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetSession) ResumeSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.ResumeSubnet(&_RootNet.TransactOpts, subnetId)
}

// ResumeSubnet is a paid mutator transaction binding the contract method 0x5364944c.
//
// Solidity: function resumeSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetTransactorSession) ResumeSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.ResumeSubnet(&_RootNet.TransactOpts, subnetId)
}

// SetAlphaTokenFactory is a paid mutator transaction binding the contract method 0x901a71e4.
//
// Solidity: function setAlphaTokenFactory(address factory) returns()
func (_RootNet *RootNetTransactor) SetAlphaTokenFactory(opts *bind.TransactOpts, factory common.Address) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "setAlphaTokenFactory", factory)
}

// SetAlphaTokenFactory is a paid mutator transaction binding the contract method 0x901a71e4.
//
// Solidity: function setAlphaTokenFactory(address factory) returns()
func (_RootNet *RootNetSession) SetAlphaTokenFactory(factory common.Address) (*types.Transaction, error) {
	return _RootNet.Contract.SetAlphaTokenFactory(&_RootNet.TransactOpts, factory)
}

// SetAlphaTokenFactory is a paid mutator transaction binding the contract method 0x901a71e4.
//
// Solidity: function setAlphaTokenFactory(address factory) returns()
func (_RootNet *RootNetTransactorSession) SetAlphaTokenFactory(factory common.Address) (*types.Transaction, error) {
	return _RootNet.Contract.SetAlphaTokenFactory(&_RootNet.TransactOpts, factory)
}

// SetDelegation is a paid mutator transaction binding the contract method 0x1ddc304a.
//
// Solidity: function setDelegation(address agent, bool _isManager) returns()
func (_RootNet *RootNetTransactor) SetDelegation(opts *bind.TransactOpts, agent common.Address, _isManager bool) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "setDelegation", agent, _isManager)
}

// SetDelegation is a paid mutator transaction binding the contract method 0x1ddc304a.
//
// Solidity: function setDelegation(address agent, bool _isManager) returns()
func (_RootNet *RootNetSession) SetDelegation(agent common.Address, _isManager bool) (*types.Transaction, error) {
	return _RootNet.Contract.SetDelegation(&_RootNet.TransactOpts, agent, _isManager)
}

// SetDelegation is a paid mutator transaction binding the contract method 0x1ddc304a.
//
// Solidity: function setDelegation(address agent, bool _isManager) returns()
func (_RootNet *RootNetTransactorSession) SetDelegation(agent common.Address, _isManager bool) (*types.Transaction, error) {
	return _RootNet.Contract.SetDelegation(&_RootNet.TransactOpts, agent, _isManager)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_RootNet *RootNetTransactor) SetGuardian(opts *bind.TransactOpts, g common.Address) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "setGuardian", g)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_RootNet *RootNetSession) SetGuardian(g common.Address) (*types.Transaction, error) {
	return _RootNet.Contract.SetGuardian(&_RootNet.TransactOpts, g)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_RootNet *RootNetTransactorSession) SetGuardian(g common.Address) (*types.Transaction, error) {
	return _RootNet.Contract.SetGuardian(&_RootNet.TransactOpts, g)
}

// SetImmunityPeriod is a paid mutator transaction binding the contract method 0x33bbf030.
//
// Solidity: function setImmunityPeriod(uint256 p) returns()
func (_RootNet *RootNetTransactor) SetImmunityPeriod(opts *bind.TransactOpts, p *big.Int) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "setImmunityPeriod", p)
}

// SetImmunityPeriod is a paid mutator transaction binding the contract method 0x33bbf030.
//
// Solidity: function setImmunityPeriod(uint256 p) returns()
func (_RootNet *RootNetSession) SetImmunityPeriod(p *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.SetImmunityPeriod(&_RootNet.TransactOpts, p)
}

// SetImmunityPeriod is a paid mutator transaction binding the contract method 0x33bbf030.
//
// Solidity: function setImmunityPeriod(uint256 p) returns()
func (_RootNet *RootNetTransactorSession) SetImmunityPeriod(p *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.SetImmunityPeriod(&_RootNet.TransactOpts, p)
}

// SetInitialAlphaPrice is a paid mutator transaction binding the contract method 0xe7d89b71.
//
// Solidity: function setInitialAlphaPrice(uint256 price) returns()
func (_RootNet *RootNetTransactor) SetInitialAlphaPrice(opts *bind.TransactOpts, price *big.Int) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "setInitialAlphaPrice", price)
}

// SetInitialAlphaPrice is a paid mutator transaction binding the contract method 0xe7d89b71.
//
// Solidity: function setInitialAlphaPrice(uint256 price) returns()
func (_RootNet *RootNetSession) SetInitialAlphaPrice(price *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.SetInitialAlphaPrice(&_RootNet.TransactOpts, price)
}

// SetInitialAlphaPrice is a paid mutator transaction binding the contract method 0xe7d89b71.
//
// Solidity: function setInitialAlphaPrice(uint256 price) returns()
func (_RootNet *RootNetTransactorSession) SetInitialAlphaPrice(price *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.SetInitialAlphaPrice(&_RootNet.TransactOpts, price)
}

// SetRewardRecipient is a paid mutator transaction binding the contract method 0xe521136f.
//
// Solidity: function setRewardRecipient(address recipient) returns()
func (_RootNet *RootNetTransactor) SetRewardRecipient(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "setRewardRecipient", recipient)
}

// SetRewardRecipient is a paid mutator transaction binding the contract method 0xe521136f.
//
// Solidity: function setRewardRecipient(address recipient) returns()
func (_RootNet *RootNetSession) SetRewardRecipient(recipient common.Address) (*types.Transaction, error) {
	return _RootNet.Contract.SetRewardRecipient(&_RootNet.TransactOpts, recipient)
}

// SetRewardRecipient is a paid mutator transaction binding the contract method 0xe521136f.
//
// Solidity: function setRewardRecipient(address recipient) returns()
func (_RootNet *RootNetTransactorSession) SetRewardRecipient(recipient common.Address) (*types.Transaction, error) {
	return _RootNet.Contract.SetRewardRecipient(&_RootNet.TransactOpts, recipient)
}

// SetSubnetManagerImpl is a paid mutator transaction binding the contract method 0xe7c17212.
//
// Solidity: function setSubnetManagerImpl(address impl) returns()
func (_RootNet *RootNetTransactor) SetSubnetManagerImpl(opts *bind.TransactOpts, impl common.Address) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "setSubnetManagerImpl", impl)
}

// SetSubnetManagerImpl is a paid mutator transaction binding the contract method 0xe7c17212.
//
// Solidity: function setSubnetManagerImpl(address impl) returns()
func (_RootNet *RootNetSession) SetSubnetManagerImpl(impl common.Address) (*types.Transaction, error) {
	return _RootNet.Contract.SetSubnetManagerImpl(&_RootNet.TransactOpts, impl)
}

// SetSubnetManagerImpl is a paid mutator transaction binding the contract method 0xe7c17212.
//
// Solidity: function setSubnetManagerImpl(address impl) returns()
func (_RootNet *RootNetTransactorSession) SetSubnetManagerImpl(impl common.Address) (*types.Transaction, error) {
	return _RootNet.Contract.SetSubnetManagerImpl(&_RootNet.TransactOpts, impl)
}

// UnbanSubnet is a paid mutator transaction binding the contract method 0x2bf1c05d.
//
// Solidity: function unbanSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetTransactor) UnbanSubnet(opts *bind.TransactOpts, subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "unbanSubnet", subnetId)
}

// UnbanSubnet is a paid mutator transaction binding the contract method 0x2bf1c05d.
//
// Solidity: function unbanSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetSession) UnbanSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.UnbanSubnet(&_RootNet.TransactOpts, subnetId)
}

// UnbanSubnet is a paid mutator transaction binding the contract method 0x2bf1c05d.
//
// Solidity: function unbanSubnet(uint256 subnetId) returns()
func (_RootNet *RootNetTransactorSession) UnbanSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _RootNet.Contract.UnbanSubnet(&_RootNet.TransactOpts, subnetId)
}

// Unbind is a paid mutator transaction binding the contract method 0xb6b25742.
//
// Solidity: function unbind() returns()
func (_RootNet *RootNetTransactor) Unbind(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "unbind")
}

// Unbind is a paid mutator transaction binding the contract method 0xb6b25742.
//
// Solidity: function unbind() returns()
func (_RootNet *RootNetSession) Unbind() (*types.Transaction, error) {
	return _RootNet.Contract.Unbind(&_RootNet.TransactOpts)
}

// Unbind is a paid mutator transaction binding the contract method 0xb6b25742.
//
// Solidity: function unbind() returns()
func (_RootNet *RootNetTransactorSession) Unbind() (*types.Transaction, error) {
	return _RootNet.Contract.Unbind(&_RootNet.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RootNet *RootNetTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RootNet *RootNetSession) Unpause() (*types.Transaction, error) {
	return _RootNet.Contract.Unpause(&_RootNet.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_RootNet *RootNetTransactorSession) Unpause() (*types.Transaction, error) {
	return _RootNet.Contract.Unpause(&_RootNet.TransactOpts)
}

// UpdateMetadata is a paid mutator transaction binding the contract method 0x03880623.
//
// Solidity: function updateMetadata(uint256 subnetId, string metadataURI, string coordinatorURL) returns()
func (_RootNet *RootNetTransactor) UpdateMetadata(opts *bind.TransactOpts, subnetId *big.Int, metadataURI string, coordinatorURL string) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "updateMetadata", subnetId, metadataURI, coordinatorURL)
}

// UpdateMetadata is a paid mutator transaction binding the contract method 0x03880623.
//
// Solidity: function updateMetadata(uint256 subnetId, string metadataURI, string coordinatorURL) returns()
func (_RootNet *RootNetSession) UpdateMetadata(subnetId *big.Int, metadataURI string, coordinatorURL string) (*types.Transaction, error) {
	return _RootNet.Contract.UpdateMetadata(&_RootNet.TransactOpts, subnetId, metadataURI, coordinatorURL)
}

// UpdateMetadata is a paid mutator transaction binding the contract method 0x03880623.
//
// Solidity: function updateMetadata(uint256 subnetId, string metadataURI, string coordinatorURL) returns()
func (_RootNet *RootNetTransactorSession) UpdateMetadata(subnetId *big.Int, metadataURI string, coordinatorURL string) (*types.Transaction, error) {
	return _RootNet.Contract.UpdateMetadata(&_RootNet.TransactOpts, subnetId, metadataURI, coordinatorURL)
}

// RootNetAgentBoundIterator is returned from FilterAgentBound and is used to iterate over the raw logs and unpacked data for AgentBound events raised by the RootNet contract.
type RootNetAgentBoundIterator struct {
	Event *RootNetAgentBound // Event containing the contract specifics and raw log

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
func (it *RootNetAgentBoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetAgentBound)
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
		it.Event = new(RootNetAgentBound)
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
func (it *RootNetAgentBoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetAgentBoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetAgentBound represents a AgentBound event raised by the RootNet contract.
type RootNetAgentBound struct {
	Principal    common.Address
	Agent        common.Address
	OldPrincipal common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterAgentBound is a free log retrieval operation binding the contract event 0x4ee4aa1bbc31e8b57dad2c2cffa4627ad65ac133b3cea2acb4870c44b5ea6b17.
//
// Solidity: event AgentBound(address indexed principal, address indexed agent, address oldPrincipal)
func (_RootNet *RootNetFilterer) FilterAgentBound(opts *bind.FilterOpts, principal []common.Address, agent []common.Address) (*RootNetAgentBoundIterator, error) {

	var principalRule []interface{}
	for _, principalItem := range principal {
		principalRule = append(principalRule, principalItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "AgentBound", principalRule, agentRule)
	if err != nil {
		return nil, err
	}
	return &RootNetAgentBoundIterator{contract: _RootNet.contract, event: "AgentBound", logs: logs, sub: sub}, nil
}

// WatchAgentBound is a free log subscription operation binding the contract event 0x4ee4aa1bbc31e8b57dad2c2cffa4627ad65ac133b3cea2acb4870c44b5ea6b17.
//
// Solidity: event AgentBound(address indexed principal, address indexed agent, address oldPrincipal)
func (_RootNet *RootNetFilterer) WatchAgentBound(opts *bind.WatchOpts, sink chan<- *RootNetAgentBound, principal []common.Address, agent []common.Address) (event.Subscription, error) {

	var principalRule []interface{}
	for _, principalItem := range principal {
		principalRule = append(principalRule, principalItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "AgentBound", principalRule, agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetAgentBound)
				if err := _RootNet.contract.UnpackLog(event, "AgentBound", log); err != nil {
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

// ParseAgentBound is a log parse operation binding the contract event 0x4ee4aa1bbc31e8b57dad2c2cffa4627ad65ac133b3cea2acb4870c44b5ea6b17.
//
// Solidity: event AgentBound(address indexed principal, address indexed agent, address oldPrincipal)
func (_RootNet *RootNetFilterer) ParseAgentBound(log types.Log) (*RootNetAgentBound, error) {
	event := new(RootNetAgentBound)
	if err := _RootNet.contract.UnpackLog(event, "AgentBound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetAgentRemovedIterator is returned from FilterAgentRemoved and is used to iterate over the raw logs and unpacked data for AgentRemoved events raised by the RootNet contract.
type RootNetAgentRemovedIterator struct {
	Event *RootNetAgentRemoved // Event containing the contract specifics and raw log

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
func (it *RootNetAgentRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetAgentRemoved)
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
		it.Event = new(RootNetAgentRemoved)
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
func (it *RootNetAgentRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetAgentRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetAgentRemoved represents a AgentRemoved event raised by the RootNet contract.
type RootNetAgentRemoved struct {
	User     common.Address
	Agent    common.Address
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAgentRemoved is a free log retrieval operation binding the contract event 0x877ef5b4e3b78ab10b445521d0724510a2c3e98f0812879447b7e08785ca866e.
//
// Solidity: event AgentRemoved(address indexed user, address indexed agent, address operator)
func (_RootNet *RootNetFilterer) FilterAgentRemoved(opts *bind.FilterOpts, user []common.Address, agent []common.Address) (*RootNetAgentRemovedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "AgentRemoved", userRule, agentRule)
	if err != nil {
		return nil, err
	}
	return &RootNetAgentRemovedIterator{contract: _RootNet.contract, event: "AgentRemoved", logs: logs, sub: sub}, nil
}

// WatchAgentRemoved is a free log subscription operation binding the contract event 0x877ef5b4e3b78ab10b445521d0724510a2c3e98f0812879447b7e08785ca866e.
//
// Solidity: event AgentRemoved(address indexed user, address indexed agent, address operator)
func (_RootNet *RootNetFilterer) WatchAgentRemoved(opts *bind.WatchOpts, sink chan<- *RootNetAgentRemoved, user []common.Address, agent []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "AgentRemoved", userRule, agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetAgentRemoved)
				if err := _RootNet.contract.UnpackLog(event, "AgentRemoved", log); err != nil {
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

// ParseAgentRemoved is a log parse operation binding the contract event 0x877ef5b4e3b78ab10b445521d0724510a2c3e98f0812879447b7e08785ca866e.
//
// Solidity: event AgentRemoved(address indexed user, address indexed agent, address operator)
func (_RootNet *RootNetFilterer) ParseAgentRemoved(log types.Log) (*RootNetAgentRemoved, error) {
	event := new(RootNetAgentRemoved)
	if err := _RootNet.contract.UnpackLog(event, "AgentRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetAgentUnboundIterator is returned from FilterAgentUnbound and is used to iterate over the raw logs and unpacked data for AgentUnbound events raised by the RootNet contract.
type RootNetAgentUnboundIterator struct {
	Event *RootNetAgentUnbound // Event containing the contract specifics and raw log

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
func (it *RootNetAgentUnboundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetAgentUnbound)
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
		it.Event = new(RootNetAgentUnbound)
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
func (it *RootNetAgentUnboundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetAgentUnboundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetAgentUnbound represents a AgentUnbound event raised by the RootNet contract.
type RootNetAgentUnbound struct {
	Principal common.Address
	Agent     common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAgentUnbound is a free log retrieval operation binding the contract event 0x3e2d9d696fa5ddd5b13727a43861bb914938ca9d534d942f5c33725656c469b1.
//
// Solidity: event AgentUnbound(address indexed principal, address indexed agent)
func (_RootNet *RootNetFilterer) FilterAgentUnbound(opts *bind.FilterOpts, principal []common.Address, agent []common.Address) (*RootNetAgentUnboundIterator, error) {

	var principalRule []interface{}
	for _, principalItem := range principal {
		principalRule = append(principalRule, principalItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "AgentUnbound", principalRule, agentRule)
	if err != nil {
		return nil, err
	}
	return &RootNetAgentUnboundIterator{contract: _RootNet.contract, event: "AgentUnbound", logs: logs, sub: sub}, nil
}

// WatchAgentUnbound is a free log subscription operation binding the contract event 0x3e2d9d696fa5ddd5b13727a43861bb914938ca9d534d942f5c33725656c469b1.
//
// Solidity: event AgentUnbound(address indexed principal, address indexed agent)
func (_RootNet *RootNetFilterer) WatchAgentUnbound(opts *bind.WatchOpts, sink chan<- *RootNetAgentUnbound, principal []common.Address, agent []common.Address) (event.Subscription, error) {

	var principalRule []interface{}
	for _, principalItem := range principal {
		principalRule = append(principalRule, principalItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "AgentUnbound", principalRule, agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetAgentUnbound)
				if err := _RootNet.contract.UnpackLog(event, "AgentUnbound", log); err != nil {
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

// ParseAgentUnbound is a log parse operation binding the contract event 0x3e2d9d696fa5ddd5b13727a43861bb914938ca9d534d942f5c33725656c469b1.
//
// Solidity: event AgentUnbound(address indexed principal, address indexed agent)
func (_RootNet *RootNetFilterer) ParseAgentUnbound(log types.Log) (*RootNetAgentUnbound, error) {
	event := new(RootNetAgentUnbound)
	if err := _RootNet.contract.UnpackLog(event, "AgentUnbound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetAllocatedIterator is returned from FilterAllocated and is used to iterate over the raw logs and unpacked data for Allocated events raised by the RootNet contract.
type RootNetAllocatedIterator struct {
	Event *RootNetAllocated // Event containing the contract specifics and raw log

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
func (it *RootNetAllocatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetAllocated)
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
		it.Event = new(RootNetAllocated)
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
func (it *RootNetAllocatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetAllocatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetAllocated represents a Allocated event raised by the RootNet contract.
type RootNetAllocated struct {
	User     common.Address
	Agent    common.Address
	SubnetId *big.Int
	Amount   *big.Int
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAllocated is a free log retrieval operation binding the contract event 0x655f98c7dae1bab3e2db10cdb4407717b9d219cf2e585adc1edba92d48af2b15.
//
// Solidity: event Allocated(address indexed user, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator)
func (_RootNet *RootNetFilterer) FilterAllocated(opts *bind.FilterOpts, user []common.Address, agent []common.Address, subnetId []*big.Int) (*RootNetAllocatedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "Allocated", userRule, agentRule, subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &RootNetAllocatedIterator{contract: _RootNet.contract, event: "Allocated", logs: logs, sub: sub}, nil
}

// WatchAllocated is a free log subscription operation binding the contract event 0x655f98c7dae1bab3e2db10cdb4407717b9d219cf2e585adc1edba92d48af2b15.
//
// Solidity: event Allocated(address indexed user, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator)
func (_RootNet *RootNetFilterer) WatchAllocated(opts *bind.WatchOpts, sink chan<- *RootNetAllocated, user []common.Address, agent []common.Address, subnetId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "Allocated", userRule, agentRule, subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetAllocated)
				if err := _RootNet.contract.UnpackLog(event, "Allocated", log); err != nil {
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
// Solidity: event Allocated(address indexed user, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator)
func (_RootNet *RootNetFilterer) ParseAllocated(log types.Log) (*RootNetAllocated, error) {
	event := new(RootNetAllocated)
	if err := _RootNet.contract.UnpackLog(event, "Allocated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetDeallocatedIterator is returned from FilterDeallocated and is used to iterate over the raw logs and unpacked data for Deallocated events raised by the RootNet contract.
type RootNetDeallocatedIterator struct {
	Event *RootNetDeallocated // Event containing the contract specifics and raw log

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
func (it *RootNetDeallocatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetDeallocated)
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
		it.Event = new(RootNetDeallocated)
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
func (it *RootNetDeallocatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetDeallocatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetDeallocated represents a Deallocated event raised by the RootNet contract.
type RootNetDeallocated struct {
	User     common.Address
	Agent    common.Address
	SubnetId *big.Int
	Amount   *big.Int
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDeallocated is a free log retrieval operation binding the contract event 0xd55bd7964253d1d9ce9187c8187b1c904274a3f374c9074f6de6fa77746bf345.
//
// Solidity: event Deallocated(address indexed user, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator)
func (_RootNet *RootNetFilterer) FilterDeallocated(opts *bind.FilterOpts, user []common.Address, agent []common.Address, subnetId []*big.Int) (*RootNetDeallocatedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "Deallocated", userRule, agentRule, subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &RootNetDeallocatedIterator{contract: _RootNet.contract, event: "Deallocated", logs: logs, sub: sub}, nil
}

// WatchDeallocated is a free log subscription operation binding the contract event 0xd55bd7964253d1d9ce9187c8187b1c904274a3f374c9074f6de6fa77746bf345.
//
// Solidity: event Deallocated(address indexed user, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator)
func (_RootNet *RootNetFilterer) WatchDeallocated(opts *bind.WatchOpts, sink chan<- *RootNetDeallocated, user []common.Address, agent []common.Address, subnetId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "Deallocated", userRule, agentRule, subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetDeallocated)
				if err := _RootNet.contract.UnpackLog(event, "Deallocated", log); err != nil {
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
// Solidity: event Deallocated(address indexed user, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator)
func (_RootNet *RootNetFilterer) ParseDeallocated(log types.Log) (*RootNetDeallocated, error) {
	event := new(RootNetDeallocated)
	if err := _RootNet.contract.UnpackLog(event, "Deallocated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetDelegationUpdatedIterator is returned from FilterDelegationUpdated and is used to iterate over the raw logs and unpacked data for DelegationUpdated events raised by the RootNet contract.
type RootNetDelegationUpdatedIterator struct {
	Event *RootNetDelegationUpdated // Event containing the contract specifics and raw log

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
func (it *RootNetDelegationUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetDelegationUpdated)
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
		it.Event = new(RootNetDelegationUpdated)
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
func (it *RootNetDelegationUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetDelegationUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetDelegationUpdated represents a DelegationUpdated event raised by the RootNet contract.
type RootNetDelegationUpdated struct {
	User      common.Address
	Agent     common.Address
	IsManager bool
	Operator  common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDelegationUpdated is a free log retrieval operation binding the contract event 0x34dbef79b9de038294b4a8f1789ad62e1b9ebaa23af56a3b75f375ce1185a9b1.
//
// Solidity: event DelegationUpdated(address indexed user, address indexed agent, bool isManager, address operator)
func (_RootNet *RootNetFilterer) FilterDelegationUpdated(opts *bind.FilterOpts, user []common.Address, agent []common.Address) (*RootNetDelegationUpdatedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "DelegationUpdated", userRule, agentRule)
	if err != nil {
		return nil, err
	}
	return &RootNetDelegationUpdatedIterator{contract: _RootNet.contract, event: "DelegationUpdated", logs: logs, sub: sub}, nil
}

// WatchDelegationUpdated is a free log subscription operation binding the contract event 0x34dbef79b9de038294b4a8f1789ad62e1b9ebaa23af56a3b75f375ce1185a9b1.
//
// Solidity: event DelegationUpdated(address indexed user, address indexed agent, bool isManager, address operator)
func (_RootNet *RootNetFilterer) WatchDelegationUpdated(opts *bind.WatchOpts, sink chan<- *RootNetDelegationUpdated, user []common.Address, agent []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "DelegationUpdated", userRule, agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetDelegationUpdated)
				if err := _RootNet.contract.UnpackLog(event, "DelegationUpdated", log); err != nil {
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

// ParseDelegationUpdated is a log parse operation binding the contract event 0x34dbef79b9de038294b4a8f1789ad62e1b9ebaa23af56a3b75f375ce1185a9b1.
//
// Solidity: event DelegationUpdated(address indexed user, address indexed agent, bool isManager, address operator)
func (_RootNet *RootNetFilterer) ParseDelegationUpdated(log types.Log) (*RootNetDelegationUpdated, error) {
	event := new(RootNetDelegationUpdated)
	if err := _RootNet.contract.UnpackLog(event, "DelegationUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the RootNet contract.
type RootNetEIP712DomainChangedIterator struct {
	Event *RootNetEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *RootNetEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetEIP712DomainChanged)
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
		it.Event = new(RootNetEIP712DomainChanged)
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
func (it *RootNetEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetEIP712DomainChanged represents a EIP712DomainChanged event raised by the RootNet contract.
type RootNetEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_RootNet *RootNetFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*RootNetEIP712DomainChangedIterator, error) {

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &RootNetEIP712DomainChangedIterator{contract: _RootNet.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_RootNet *RootNetFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *RootNetEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetEIP712DomainChanged)
				if err := _RootNet.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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
func (_RootNet *RootNetFilterer) ParseEIP712DomainChanged(log types.Log) (*RootNetEIP712DomainChanged, error) {
	event := new(RootNetEIP712DomainChanged)
	if err := _RootNet.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetLPCreatedIterator is returned from FilterLPCreated and is used to iterate over the raw logs and unpacked data for LPCreated events raised by the RootNet contract.
type RootNetLPCreatedIterator struct {
	Event *RootNetLPCreated // Event containing the contract specifics and raw log

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
func (it *RootNetLPCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetLPCreated)
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
		it.Event = new(RootNetLPCreated)
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
func (it *RootNetLPCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetLPCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetLPCreated represents a LPCreated event raised by the RootNet contract.
type RootNetLPCreated struct {
	SubnetId    *big.Int
	PoolId      [32]byte
	AwpAmount   *big.Int
	AlphaAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterLPCreated is a free log retrieval operation binding the contract event 0x0a28a1fd5e0909199ee082834df66cfaae2125e3503bf16d2dc46214278fc7ab.
//
// Solidity: event LPCreated(uint256 indexed subnetId, bytes32 poolId, uint256 awpAmount, uint256 alphaAmount)
func (_RootNet *RootNetFilterer) FilterLPCreated(opts *bind.FilterOpts, subnetId []*big.Int) (*RootNetLPCreatedIterator, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "LPCreated", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &RootNetLPCreatedIterator{contract: _RootNet.contract, event: "LPCreated", logs: logs, sub: sub}, nil
}

// WatchLPCreated is a free log subscription operation binding the contract event 0x0a28a1fd5e0909199ee082834df66cfaae2125e3503bf16d2dc46214278fc7ab.
//
// Solidity: event LPCreated(uint256 indexed subnetId, bytes32 poolId, uint256 awpAmount, uint256 alphaAmount)
func (_RootNet *RootNetFilterer) WatchLPCreated(opts *bind.WatchOpts, sink chan<- *RootNetLPCreated, subnetId []*big.Int) (event.Subscription, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "LPCreated", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetLPCreated)
				if err := _RootNet.contract.UnpackLog(event, "LPCreated", log); err != nil {
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

// ParseLPCreated is a log parse operation binding the contract event 0x0a28a1fd5e0909199ee082834df66cfaae2125e3503bf16d2dc46214278fc7ab.
//
// Solidity: event LPCreated(uint256 indexed subnetId, bytes32 poolId, uint256 awpAmount, uint256 alphaAmount)
func (_RootNet *RootNetFilterer) ParseLPCreated(log types.Log) (*RootNetLPCreated, error) {
	event := new(RootNetLPCreated)
	if err := _RootNet.contract.UnpackLog(event, "LPCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetMetadataUpdatedIterator is returned from FilterMetadataUpdated and is used to iterate over the raw logs and unpacked data for MetadataUpdated events raised by the RootNet contract.
type RootNetMetadataUpdatedIterator struct {
	Event *RootNetMetadataUpdated // Event containing the contract specifics and raw log

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
func (it *RootNetMetadataUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetMetadataUpdated)
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
		it.Event = new(RootNetMetadataUpdated)
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
func (it *RootNetMetadataUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetMetadataUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetMetadataUpdated represents a MetadataUpdated event raised by the RootNet contract.
type RootNetMetadataUpdated struct {
	SubnetId       *big.Int
	MetadataURI    string
	CoordinatorURL string
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterMetadataUpdated is a free log retrieval operation binding the contract event 0x4bb348c8e52124f1a18e983f64ad1bc5d380a3ca43654fbb2b8f73c71f305459.
//
// Solidity: event MetadataUpdated(uint256 indexed subnetId, string metadataURI, string coordinatorURL)
func (_RootNet *RootNetFilterer) FilterMetadataUpdated(opts *bind.FilterOpts, subnetId []*big.Int) (*RootNetMetadataUpdatedIterator, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "MetadataUpdated", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &RootNetMetadataUpdatedIterator{contract: _RootNet.contract, event: "MetadataUpdated", logs: logs, sub: sub}, nil
}

// WatchMetadataUpdated is a free log subscription operation binding the contract event 0x4bb348c8e52124f1a18e983f64ad1bc5d380a3ca43654fbb2b8f73c71f305459.
//
// Solidity: event MetadataUpdated(uint256 indexed subnetId, string metadataURI, string coordinatorURL)
func (_RootNet *RootNetFilterer) WatchMetadataUpdated(opts *bind.WatchOpts, sink chan<- *RootNetMetadataUpdated, subnetId []*big.Int) (event.Subscription, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "MetadataUpdated", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetMetadataUpdated)
				if err := _RootNet.contract.UnpackLog(event, "MetadataUpdated", log); err != nil {
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

// ParseMetadataUpdated is a log parse operation binding the contract event 0x4bb348c8e52124f1a18e983f64ad1bc5d380a3ca43654fbb2b8f73c71f305459.
//
// Solidity: event MetadataUpdated(uint256 indexed subnetId, string metadataURI, string coordinatorURL)
func (_RootNet *RootNetFilterer) ParseMetadataUpdated(log types.Log) (*RootNetMetadataUpdated, error) {
	event := new(RootNetMetadataUpdated)
	if err := _RootNet.contract.UnpackLog(event, "MetadataUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the RootNet contract.
type RootNetPausedIterator struct {
	Event *RootNetPaused // Event containing the contract specifics and raw log

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
func (it *RootNetPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetPaused)
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
		it.Event = new(RootNetPaused)
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
func (it *RootNetPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetPaused represents a Paused event raised by the RootNet contract.
type RootNetPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RootNet *RootNetFilterer) FilterPaused(opts *bind.FilterOpts) (*RootNetPausedIterator, error) {

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &RootNetPausedIterator{contract: _RootNet.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RootNet *RootNetFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *RootNetPaused) (event.Subscription, error) {

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetPaused)
				if err := _RootNet.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_RootNet *RootNetFilterer) ParsePaused(log types.Log) (*RootNetPaused, error) {
	event := new(RootNetPaused)
	if err := _RootNet.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetReallocatedIterator is returned from FilterReallocated and is used to iterate over the raw logs and unpacked data for Reallocated events raised by the RootNet contract.
type RootNetReallocatedIterator struct {
	Event *RootNetReallocated // Event containing the contract specifics and raw log

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
func (it *RootNetReallocatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetReallocated)
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
		it.Event = new(RootNetReallocated)
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
func (it *RootNetReallocatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetReallocatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetReallocated represents a Reallocated event raised by the RootNet contract.
type RootNetReallocated struct {
	User       common.Address
	FromAgent  common.Address
	FromSubnet *big.Int
	ToAgent    common.Address
	ToSubnet   *big.Int
	Amount     *big.Int
	Operator   common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterReallocated is a free log retrieval operation binding the contract event 0x726c93ba67bfe4c677e37114279f0ad9aab5ee9ffbd1158923be5d0fec3b1b45.
//
// Solidity: event Reallocated(address indexed user, address fromAgent, uint256 fromSubnet, address toAgent, uint256 toSubnet, uint256 amount, address operator)
func (_RootNet *RootNetFilterer) FilterReallocated(opts *bind.FilterOpts, user []common.Address) (*RootNetReallocatedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "Reallocated", userRule)
	if err != nil {
		return nil, err
	}
	return &RootNetReallocatedIterator{contract: _RootNet.contract, event: "Reallocated", logs: logs, sub: sub}, nil
}

// WatchReallocated is a free log subscription operation binding the contract event 0x726c93ba67bfe4c677e37114279f0ad9aab5ee9ffbd1158923be5d0fec3b1b45.
//
// Solidity: event Reallocated(address indexed user, address fromAgent, uint256 fromSubnet, address toAgent, uint256 toSubnet, uint256 amount, address operator)
func (_RootNet *RootNetFilterer) WatchReallocated(opts *bind.WatchOpts, sink chan<- *RootNetReallocated, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "Reallocated", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetReallocated)
				if err := _RootNet.contract.UnpackLog(event, "Reallocated", log); err != nil {
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
// Solidity: event Reallocated(address indexed user, address fromAgent, uint256 fromSubnet, address toAgent, uint256 toSubnet, uint256 amount, address operator)
func (_RootNet *RootNetFilterer) ParseReallocated(log types.Log) (*RootNetReallocated, error) {
	event := new(RootNetReallocated)
	if err := _RootNet.contract.UnpackLog(event, "Reallocated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetRewardRecipientUpdatedIterator is returned from FilterRewardRecipientUpdated and is used to iterate over the raw logs and unpacked data for RewardRecipientUpdated events raised by the RootNet contract.
type RootNetRewardRecipientUpdatedIterator struct {
	Event *RootNetRewardRecipientUpdated // Event containing the contract specifics and raw log

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
func (it *RootNetRewardRecipientUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetRewardRecipientUpdated)
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
		it.Event = new(RootNetRewardRecipientUpdated)
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
func (it *RootNetRewardRecipientUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetRewardRecipientUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetRewardRecipientUpdated represents a RewardRecipientUpdated event raised by the RootNet contract.
type RootNetRewardRecipientUpdated struct {
	User      common.Address
	Recipient common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRewardRecipientUpdated is a free log retrieval operation binding the contract event 0xc8c11bb97ac2ffa10ce2e2a98f4c1fd8df84cfa2e1a15e013ed2383ab1f527ad.
//
// Solidity: event RewardRecipientUpdated(address indexed user, address recipient)
func (_RootNet *RootNetFilterer) FilterRewardRecipientUpdated(opts *bind.FilterOpts, user []common.Address) (*RootNetRewardRecipientUpdatedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "RewardRecipientUpdated", userRule)
	if err != nil {
		return nil, err
	}
	return &RootNetRewardRecipientUpdatedIterator{contract: _RootNet.contract, event: "RewardRecipientUpdated", logs: logs, sub: sub}, nil
}

// WatchRewardRecipientUpdated is a free log subscription operation binding the contract event 0xc8c11bb97ac2ffa10ce2e2a98f4c1fd8df84cfa2e1a15e013ed2383ab1f527ad.
//
// Solidity: event RewardRecipientUpdated(address indexed user, address recipient)
func (_RootNet *RootNetFilterer) WatchRewardRecipientUpdated(opts *bind.WatchOpts, sink chan<- *RootNetRewardRecipientUpdated, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "RewardRecipientUpdated", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetRewardRecipientUpdated)
				if err := _RootNet.contract.UnpackLog(event, "RewardRecipientUpdated", log); err != nil {
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

// ParseRewardRecipientUpdated is a log parse operation binding the contract event 0xc8c11bb97ac2ffa10ce2e2a98f4c1fd8df84cfa2e1a15e013ed2383ab1f527ad.
//
// Solidity: event RewardRecipientUpdated(address indexed user, address recipient)
func (_RootNet *RootNetFilterer) ParseRewardRecipientUpdated(log types.Log) (*RootNetRewardRecipientUpdated, error) {
	event := new(RootNetRewardRecipientUpdated)
	if err := _RootNet.contract.UnpackLog(event, "RewardRecipientUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetSubnetActivatedIterator is returned from FilterSubnetActivated and is used to iterate over the raw logs and unpacked data for SubnetActivated events raised by the RootNet contract.
type RootNetSubnetActivatedIterator struct {
	Event *RootNetSubnetActivated // Event containing the contract specifics and raw log

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
func (it *RootNetSubnetActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetSubnetActivated)
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
		it.Event = new(RootNetSubnetActivated)
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
func (it *RootNetSubnetActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetSubnetActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetSubnetActivated represents a SubnetActivated event raised by the RootNet contract.
type RootNetSubnetActivated struct {
	SubnetId *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSubnetActivated is a free log retrieval operation binding the contract event 0x034804b969efac7a0df7757ada640ffdcc09f25dbcd4582c390f25d5622255c4.
//
// Solidity: event SubnetActivated(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) FilterSubnetActivated(opts *bind.FilterOpts, subnetId []*big.Int) (*RootNetSubnetActivatedIterator, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "SubnetActivated", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &RootNetSubnetActivatedIterator{contract: _RootNet.contract, event: "SubnetActivated", logs: logs, sub: sub}, nil
}

// WatchSubnetActivated is a free log subscription operation binding the contract event 0x034804b969efac7a0df7757ada640ffdcc09f25dbcd4582c390f25d5622255c4.
//
// Solidity: event SubnetActivated(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) WatchSubnetActivated(opts *bind.WatchOpts, sink chan<- *RootNetSubnetActivated, subnetId []*big.Int) (event.Subscription, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "SubnetActivated", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetSubnetActivated)
				if err := _RootNet.contract.UnpackLog(event, "SubnetActivated", log); err != nil {
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

// ParseSubnetActivated is a log parse operation binding the contract event 0x034804b969efac7a0df7757ada640ffdcc09f25dbcd4582c390f25d5622255c4.
//
// Solidity: event SubnetActivated(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) ParseSubnetActivated(log types.Log) (*RootNetSubnetActivated, error) {
	event := new(RootNetSubnetActivated)
	if err := _RootNet.contract.UnpackLog(event, "SubnetActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetSubnetBannedIterator is returned from FilterSubnetBanned and is used to iterate over the raw logs and unpacked data for SubnetBanned events raised by the RootNet contract.
type RootNetSubnetBannedIterator struct {
	Event *RootNetSubnetBanned // Event containing the contract specifics and raw log

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
func (it *RootNetSubnetBannedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetSubnetBanned)
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
		it.Event = new(RootNetSubnetBanned)
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
func (it *RootNetSubnetBannedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetSubnetBannedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetSubnetBanned represents a SubnetBanned event raised by the RootNet contract.
type RootNetSubnetBanned struct {
	SubnetId *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSubnetBanned is a free log retrieval operation binding the contract event 0xb887f21153957bddcf7211314cf42794076ccf98c6ae5cf8e2d065ec717c681b.
//
// Solidity: event SubnetBanned(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) FilterSubnetBanned(opts *bind.FilterOpts, subnetId []*big.Int) (*RootNetSubnetBannedIterator, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "SubnetBanned", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &RootNetSubnetBannedIterator{contract: _RootNet.contract, event: "SubnetBanned", logs: logs, sub: sub}, nil
}

// WatchSubnetBanned is a free log subscription operation binding the contract event 0xb887f21153957bddcf7211314cf42794076ccf98c6ae5cf8e2d065ec717c681b.
//
// Solidity: event SubnetBanned(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) WatchSubnetBanned(opts *bind.WatchOpts, sink chan<- *RootNetSubnetBanned, subnetId []*big.Int) (event.Subscription, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "SubnetBanned", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetSubnetBanned)
				if err := _RootNet.contract.UnpackLog(event, "SubnetBanned", log); err != nil {
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

// ParseSubnetBanned is a log parse operation binding the contract event 0xb887f21153957bddcf7211314cf42794076ccf98c6ae5cf8e2d065ec717c681b.
//
// Solidity: event SubnetBanned(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) ParseSubnetBanned(log types.Log) (*RootNetSubnetBanned, error) {
	event := new(RootNetSubnetBanned)
	if err := _RootNet.contract.UnpackLog(event, "SubnetBanned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetSubnetDeregisteredIterator is returned from FilterSubnetDeregistered and is used to iterate over the raw logs and unpacked data for SubnetDeregistered events raised by the RootNet contract.
type RootNetSubnetDeregisteredIterator struct {
	Event *RootNetSubnetDeregistered // Event containing the contract specifics and raw log

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
func (it *RootNetSubnetDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetSubnetDeregistered)
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
		it.Event = new(RootNetSubnetDeregistered)
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
func (it *RootNetSubnetDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetSubnetDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetSubnetDeregistered represents a SubnetDeregistered event raised by the RootNet contract.
type RootNetSubnetDeregistered struct {
	SubnetId *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSubnetDeregistered is a free log retrieval operation binding the contract event 0x960c7566f4c9bb6958ff6e37a02b5ae69fa39dd75651fcc9b9a1029c01d0ff32.
//
// Solidity: event SubnetDeregistered(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) FilterSubnetDeregistered(opts *bind.FilterOpts, subnetId []*big.Int) (*RootNetSubnetDeregisteredIterator, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "SubnetDeregistered", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &RootNetSubnetDeregisteredIterator{contract: _RootNet.contract, event: "SubnetDeregistered", logs: logs, sub: sub}, nil
}

// WatchSubnetDeregistered is a free log subscription operation binding the contract event 0x960c7566f4c9bb6958ff6e37a02b5ae69fa39dd75651fcc9b9a1029c01d0ff32.
//
// Solidity: event SubnetDeregistered(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) WatchSubnetDeregistered(opts *bind.WatchOpts, sink chan<- *RootNetSubnetDeregistered, subnetId []*big.Int) (event.Subscription, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "SubnetDeregistered", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetSubnetDeregistered)
				if err := _RootNet.contract.UnpackLog(event, "SubnetDeregistered", log); err != nil {
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

// ParseSubnetDeregistered is a log parse operation binding the contract event 0x960c7566f4c9bb6958ff6e37a02b5ae69fa39dd75651fcc9b9a1029c01d0ff32.
//
// Solidity: event SubnetDeregistered(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) ParseSubnetDeregistered(log types.Log) (*RootNetSubnetDeregistered, error) {
	event := new(RootNetSubnetDeregistered)
	if err := _RootNet.contract.UnpackLog(event, "SubnetDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetSubnetPausedIterator is returned from FilterSubnetPaused and is used to iterate over the raw logs and unpacked data for SubnetPaused events raised by the RootNet contract.
type RootNetSubnetPausedIterator struct {
	Event *RootNetSubnetPaused // Event containing the contract specifics and raw log

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
func (it *RootNetSubnetPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetSubnetPaused)
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
		it.Event = new(RootNetSubnetPaused)
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
func (it *RootNetSubnetPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetSubnetPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetSubnetPaused represents a SubnetPaused event raised by the RootNet contract.
type RootNetSubnetPaused struct {
	SubnetId *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSubnetPaused is a free log retrieval operation binding the contract event 0x789ca96cb827d1dcb6bfc7d9d084d0a574dadff90700e3602acedefb10f69afc.
//
// Solidity: event SubnetPaused(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) FilterSubnetPaused(opts *bind.FilterOpts, subnetId []*big.Int) (*RootNetSubnetPausedIterator, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "SubnetPaused", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &RootNetSubnetPausedIterator{contract: _RootNet.contract, event: "SubnetPaused", logs: logs, sub: sub}, nil
}

// WatchSubnetPaused is a free log subscription operation binding the contract event 0x789ca96cb827d1dcb6bfc7d9d084d0a574dadff90700e3602acedefb10f69afc.
//
// Solidity: event SubnetPaused(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) WatchSubnetPaused(opts *bind.WatchOpts, sink chan<- *RootNetSubnetPaused, subnetId []*big.Int) (event.Subscription, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "SubnetPaused", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetSubnetPaused)
				if err := _RootNet.contract.UnpackLog(event, "SubnetPaused", log); err != nil {
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

// ParseSubnetPaused is a log parse operation binding the contract event 0x789ca96cb827d1dcb6bfc7d9d084d0a574dadff90700e3602acedefb10f69afc.
//
// Solidity: event SubnetPaused(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) ParseSubnetPaused(log types.Log) (*RootNetSubnetPaused, error) {
	event := new(RootNetSubnetPaused)
	if err := _RootNet.contract.UnpackLog(event, "SubnetPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetSubnetRegisteredIterator is returned from FilterSubnetRegistered and is used to iterate over the raw logs and unpacked data for SubnetRegistered events raised by the RootNet contract.
type RootNetSubnetRegisteredIterator struct {
	Event *RootNetSubnetRegistered // Event containing the contract specifics and raw log

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
func (it *RootNetSubnetRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetSubnetRegistered)
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
		it.Event = new(RootNetSubnetRegistered)
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
func (it *RootNetSubnetRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetSubnetRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetSubnetRegistered represents a SubnetRegistered event raised by the RootNet contract.
type RootNetSubnetRegistered struct {
	SubnetId       *big.Int
	Owner          common.Address
	Name           string
	Symbol         string
	MetadataURI    string
	SubnetManager  common.Address
	AlphaToken     common.Address
	CoordinatorURL string
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterSubnetRegistered is a free log retrieval operation binding the contract event 0xf1754be0b0979fb871647c40cf65eaca03d25e211ac7b3016ad3d70a49845173.
//
// Solidity: event SubnetRegistered(uint256 indexed subnetId, address indexed owner, string name, string symbol, string metadataURI, address subnetManager, address alphaToken, string coordinatorURL)
func (_RootNet *RootNetFilterer) FilterSubnetRegistered(opts *bind.FilterOpts, subnetId []*big.Int, owner []common.Address) (*RootNetSubnetRegisteredIterator, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "SubnetRegistered", subnetIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &RootNetSubnetRegisteredIterator{contract: _RootNet.contract, event: "SubnetRegistered", logs: logs, sub: sub}, nil
}

// WatchSubnetRegistered is a free log subscription operation binding the contract event 0xf1754be0b0979fb871647c40cf65eaca03d25e211ac7b3016ad3d70a49845173.
//
// Solidity: event SubnetRegistered(uint256 indexed subnetId, address indexed owner, string name, string symbol, string metadataURI, address subnetManager, address alphaToken, string coordinatorURL)
func (_RootNet *RootNetFilterer) WatchSubnetRegistered(opts *bind.WatchOpts, sink chan<- *RootNetSubnetRegistered, subnetId []*big.Int, owner []common.Address) (event.Subscription, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "SubnetRegistered", subnetIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetSubnetRegistered)
				if err := _RootNet.contract.UnpackLog(event, "SubnetRegistered", log); err != nil {
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

// ParseSubnetRegistered is a log parse operation binding the contract event 0xf1754be0b0979fb871647c40cf65eaca03d25e211ac7b3016ad3d70a49845173.
//
// Solidity: event SubnetRegistered(uint256 indexed subnetId, address indexed owner, string name, string symbol, string metadataURI, address subnetManager, address alphaToken, string coordinatorURL)
func (_RootNet *RootNetFilterer) ParseSubnetRegistered(log types.Log) (*RootNetSubnetRegistered, error) {
	event := new(RootNetSubnetRegistered)
	if err := _RootNet.contract.UnpackLog(event, "SubnetRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetSubnetResumedIterator is returned from FilterSubnetResumed and is used to iterate over the raw logs and unpacked data for SubnetResumed events raised by the RootNet contract.
type RootNetSubnetResumedIterator struct {
	Event *RootNetSubnetResumed // Event containing the contract specifics and raw log

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
func (it *RootNetSubnetResumedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetSubnetResumed)
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
		it.Event = new(RootNetSubnetResumed)
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
func (it *RootNetSubnetResumedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetSubnetResumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetSubnetResumed represents a SubnetResumed event raised by the RootNet contract.
type RootNetSubnetResumed struct {
	SubnetId *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSubnetResumed is a free log retrieval operation binding the contract event 0xf1693a0767c0183c95caf97ea0be785bece8e3578b49ef89c9669b754c1ba9a0.
//
// Solidity: event SubnetResumed(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) FilterSubnetResumed(opts *bind.FilterOpts, subnetId []*big.Int) (*RootNetSubnetResumedIterator, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "SubnetResumed", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &RootNetSubnetResumedIterator{contract: _RootNet.contract, event: "SubnetResumed", logs: logs, sub: sub}, nil
}

// WatchSubnetResumed is a free log subscription operation binding the contract event 0xf1693a0767c0183c95caf97ea0be785bece8e3578b49ef89c9669b754c1ba9a0.
//
// Solidity: event SubnetResumed(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) WatchSubnetResumed(opts *bind.WatchOpts, sink chan<- *RootNetSubnetResumed, subnetId []*big.Int) (event.Subscription, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "SubnetResumed", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetSubnetResumed)
				if err := _RootNet.contract.UnpackLog(event, "SubnetResumed", log); err != nil {
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

// ParseSubnetResumed is a log parse operation binding the contract event 0xf1693a0767c0183c95caf97ea0be785bece8e3578b49ef89c9669b754c1ba9a0.
//
// Solidity: event SubnetResumed(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) ParseSubnetResumed(log types.Log) (*RootNetSubnetResumed, error) {
	event := new(RootNetSubnetResumed)
	if err := _RootNet.contract.UnpackLog(event, "SubnetResumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetSubnetUnbannedIterator is returned from FilterSubnetUnbanned and is used to iterate over the raw logs and unpacked data for SubnetUnbanned events raised by the RootNet contract.
type RootNetSubnetUnbannedIterator struct {
	Event *RootNetSubnetUnbanned // Event containing the contract specifics and raw log

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
func (it *RootNetSubnetUnbannedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetSubnetUnbanned)
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
		it.Event = new(RootNetSubnetUnbanned)
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
func (it *RootNetSubnetUnbannedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetSubnetUnbannedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetSubnetUnbanned represents a SubnetUnbanned event raised by the RootNet contract.
type RootNetSubnetUnbanned struct {
	SubnetId *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSubnetUnbanned is a free log retrieval operation binding the contract event 0xa04fe0f9f3200108443db1507380606e909a0f81c9eb84c0707c265152668466.
//
// Solidity: event SubnetUnbanned(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) FilterSubnetUnbanned(opts *bind.FilterOpts, subnetId []*big.Int) (*RootNetSubnetUnbannedIterator, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "SubnetUnbanned", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &RootNetSubnetUnbannedIterator{contract: _RootNet.contract, event: "SubnetUnbanned", logs: logs, sub: sub}, nil
}

// WatchSubnetUnbanned is a free log subscription operation binding the contract event 0xa04fe0f9f3200108443db1507380606e909a0f81c9eb84c0707c265152668466.
//
// Solidity: event SubnetUnbanned(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) WatchSubnetUnbanned(opts *bind.WatchOpts, sink chan<- *RootNetSubnetUnbanned, subnetId []*big.Int) (event.Subscription, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "SubnetUnbanned", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetSubnetUnbanned)
				if err := _RootNet.contract.UnpackLog(event, "SubnetUnbanned", log); err != nil {
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

// ParseSubnetUnbanned is a log parse operation binding the contract event 0xa04fe0f9f3200108443db1507380606e909a0f81c9eb84c0707c265152668466.
//
// Solidity: event SubnetUnbanned(uint256 indexed subnetId)
func (_RootNet *RootNetFilterer) ParseSubnetUnbanned(log types.Log) (*RootNetSubnetUnbanned, error) {
	event := new(RootNetSubnetUnbanned)
	if err := _RootNet.contract.UnpackLog(event, "SubnetUnbanned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the RootNet contract.
type RootNetUnpausedIterator struct {
	Event *RootNetUnpaused // Event containing the contract specifics and raw log

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
func (it *RootNetUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetUnpaused)
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
		it.Event = new(RootNetUnpaused)
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
func (it *RootNetUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetUnpaused represents a Unpaused event raised by the RootNet contract.
type RootNetUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RootNet *RootNetFilterer) FilterUnpaused(opts *bind.FilterOpts) (*RootNetUnpausedIterator, error) {

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &RootNetUnpausedIterator{contract: _RootNet.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RootNet *RootNetFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *RootNetUnpaused) (event.Subscription, error) {

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetUnpaused)
				if err := _RootNet.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_RootNet *RootNetFilterer) ParseUnpaused(log types.Log) (*RootNetUnpaused, error) {
	event := new(RootNetUnpaused)
	if err := _RootNet.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetUserRegisteredIterator is returned from FilterUserRegistered and is used to iterate over the raw logs and unpacked data for UserRegistered events raised by the RootNet contract.
type RootNetUserRegisteredIterator struct {
	Event *RootNetUserRegistered // Event containing the contract specifics and raw log

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
func (it *RootNetUserRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetUserRegistered)
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
		it.Event = new(RootNetUserRegistered)
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
func (it *RootNetUserRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetUserRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetUserRegistered represents a UserRegistered event raised by the RootNet contract.
type RootNetUserRegistered struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUserRegistered is a free log retrieval operation binding the contract event 0x54db7a5cb4735e1aac1f53db512d3390390bb6637bd30ad4bf9fc98667d9b9b9.
//
// Solidity: event UserRegistered(address indexed user)
func (_RootNet *RootNetFilterer) FilterUserRegistered(opts *bind.FilterOpts, user []common.Address) (*RootNetUserRegisteredIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "UserRegistered", userRule)
	if err != nil {
		return nil, err
	}
	return &RootNetUserRegisteredIterator{contract: _RootNet.contract, event: "UserRegistered", logs: logs, sub: sub}, nil
}

// WatchUserRegistered is a free log subscription operation binding the contract event 0x54db7a5cb4735e1aac1f53db512d3390390bb6637bd30ad4bf9fc98667d9b9b9.
//
// Solidity: event UserRegistered(address indexed user)
func (_RootNet *RootNetFilterer) WatchUserRegistered(opts *bind.WatchOpts, sink chan<- *RootNetUserRegistered, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "UserRegistered", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetUserRegistered)
				if err := _RootNet.contract.UnpackLog(event, "UserRegistered", log); err != nil {
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

// ParseUserRegistered is a log parse operation binding the contract event 0x54db7a5cb4735e1aac1f53db512d3390390bb6637bd30ad4bf9fc98667d9b9b9.
//
// Solidity: event UserRegistered(address indexed user)
func (_RootNet *RootNetFilterer) ParseUserRegistered(log types.Log) (*RootNetUserRegistered, error) {
	event := new(RootNetUserRegistered)
	if err := _RootNet.contract.UnpackLog(event, "UserRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
