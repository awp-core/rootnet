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
	Name          string
	Symbol        string
	SubnetManager common.Address
	Salt          [32]byte
	MinStake      *big.Int
	SkillsURI     string
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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"deployer_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"treasury_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"guardian_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"INITIAL_ALPHA_MINT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_ACTIVE_SUBNETS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"accessManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activateSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"activateSubnetFor\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allocate\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allocateFor\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"alphaTokenFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"awpEmission\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"awpToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"banSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bind\",\"inputs\":[{\"name\":\"principal\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bindFor\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"principal\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deallocate\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deallocateFor\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultSubnetManagerImpl\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deregisterSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveSubnetCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveSubnetIdAt\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAgentInfo\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRootNet.AgentInfo\",\"components\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isValid\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"stake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"rewardRecipient\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAgentsInfo\",\"inputs\":[{\"name\":\"agents\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structRootNet.AgentInfo[]\",\"components\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isValid\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"stake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"rewardRecipient\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIRootNet.SubnetInfo\",\"components\":[{\"name\":\"lpPool\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumIRootNet.SubnetStatus\"},{\"name\":\"createdAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"activatedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSubnetFull\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIRootNet.SubnetFullInfo\",\"components\":[{\"name\":\"subnetManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"alphaToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lpPool\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumIRootNet.SubnetStatus\"},{\"name\":\"createdAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"activatedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"guardian\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"immunityPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialAlphaPrice\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initializeRegistry\",\"inputs\":[{\"name\":\"awpToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetNFT_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"alphaTokenFactory_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"awpEmission_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lpManager_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"accessManager_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakingVault_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakeNFT_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"defaultSubnetManagerImpl_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isSubnetActive\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lpManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nextSubnetId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nonces\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reallocate\",\"inputs\":[{\"name\":\"fromAgent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromSubnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"toAgent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"toSubnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"register\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"register\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"depositAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lockDuration\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerAndStake\",\"inputs\":[{\"name\":\"depositAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lockDuration\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"allocateAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerFor\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerSubnet\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIRootNet.SubnetParams\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"subnetManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerSubnetFor\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIRootNet.SubnetParams\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"subnetManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerSubnetForWithPermit\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIRootNet.SubnetParams\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"subnetManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"permitV\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"permitR\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"permitS\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"registerV\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"registerR\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"registerS\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registryInitialized\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeAgent\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resumeSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAlphaTokenFactory\",\"inputs\":[{\"name\":\"factory\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDelegation\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_isManager\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setGuardian\",\"inputs\":[{\"name\":\"g\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setImmunityPeriod\",\"inputs\":[{\"name\":\"p\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setInitialAlphaPrice\",\"inputs\":[{\"name\":\"price\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRewardRecipient\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRewardRecipientFor\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSubnetManagerImpl\",\"inputs\":[{\"name\":\"impl\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakeNFT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stakingVault\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"subnetNFT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"subnets\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"lpPool\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumIRootNet.SubnetStatus\"},{\"name\":\"createdAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"activatedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"treasury\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unbanSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unbind\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AgentBound\",\"inputs\":[{\"name\":\"principal\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"oldPrincipal\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AgentRemoved\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AgentUnbound\",\"inputs\":[{\"name\":\"principal\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Allocated\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AlphaTokenFactoryUpdated\",\"inputs\":[{\"name\":\"newFactory\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deallocated\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultSubnetManagerImplUpdated\",\"inputs\":[{\"name\":\"newImpl\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DelegationUpdated\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"isManager\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GuardianUpdated\",\"inputs\":[{\"name\":\"newGuardian\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ImmunityPeriodUpdated\",\"inputs\":[{\"name\":\"newPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InitialAlphaPriceUpdated\",\"inputs\":[{\"name\":\"newPrice\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LPCreated\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"poolId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"awpAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"alphaAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Reallocated\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fromAgent\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"fromSubnet\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"toAgent\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"toSubnet\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardRecipientUpdated\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetActivated\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetBanned\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetDeregistered\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetPaused\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetRegistered\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"subnetManager\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"alphaToken\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetResumed\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetUnbanned\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UserRegistered\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpiredSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ImmunityNotExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientMinStake\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAgent\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidShortString\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSubnetParams\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSubnetStatus\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxActiveSubnetsReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotDeployer\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotGuardian\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotTimelock\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PriceTooHigh\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PriceTooLow\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"StringTooLong\",\"inputs\":[{\"name\":\"str\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"SubnetManagerRequired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnknownAddress\",\"inputs\":[]}]",
}

// RootNetABI is the input ABI used to generate the binding from.
// Deprecated: Use RootNetMetaData.ABI instead.
var RootNetABI = RootNetMetaData.ABI

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

// ActivateSubnetFor is a paid mutator transaction binding the contract method 0x08b55cff.
//
// Solidity: function activateSubnetFor(address user, uint256 subnetId, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetTransactor) ActivateSubnetFor(opts *bind.TransactOpts, user common.Address, subnetId *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "activateSubnetFor", user, subnetId, deadline, v, r, s)
}

// ActivateSubnetFor is a paid mutator transaction binding the contract method 0x08b55cff.
//
// Solidity: function activateSubnetFor(address user, uint256 subnetId, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetSession) ActivateSubnetFor(user common.Address, subnetId *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.ActivateSubnetFor(&_RootNet.TransactOpts, user, subnetId, deadline, v, r, s)
}

// ActivateSubnetFor is a paid mutator transaction binding the contract method 0x08b55cff.
//
// Solidity: function activateSubnetFor(address user, uint256 subnetId, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetTransactorSession) ActivateSubnetFor(user common.Address, subnetId *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.ActivateSubnetFor(&_RootNet.TransactOpts, user, subnetId, deadline, v, r, s)
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

// AllocateFor is a paid mutator transaction binding the contract method 0x7d66c5c5.
//
// Solidity: function allocateFor(address user, address agent, uint256 subnetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetTransactor) AllocateFor(opts *bind.TransactOpts, user common.Address, agent common.Address, subnetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "allocateFor", user, agent, subnetId, amount, deadline, v, r, s)
}

// AllocateFor is a paid mutator transaction binding the contract method 0x7d66c5c5.
//
// Solidity: function allocateFor(address user, address agent, uint256 subnetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetSession) AllocateFor(user common.Address, agent common.Address, subnetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.AllocateFor(&_RootNet.TransactOpts, user, agent, subnetId, amount, deadline, v, r, s)
}

// AllocateFor is a paid mutator transaction binding the contract method 0x7d66c5c5.
//
// Solidity: function allocateFor(address user, address agent, uint256 subnetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetTransactorSession) AllocateFor(user common.Address, agent common.Address, subnetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.AllocateFor(&_RootNet.TransactOpts, user, agent, subnetId, amount, deadline, v, r, s)
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

// DeallocateFor is a paid mutator transaction binding the contract method 0x10fe1208.
//
// Solidity: function deallocateFor(address user, address agent, uint256 subnetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetTransactor) DeallocateFor(opts *bind.TransactOpts, user common.Address, agent common.Address, subnetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "deallocateFor", user, agent, subnetId, amount, deadline, v, r, s)
}

// DeallocateFor is a paid mutator transaction binding the contract method 0x10fe1208.
//
// Solidity: function deallocateFor(address user, address agent, uint256 subnetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetSession) DeallocateFor(user common.Address, agent common.Address, subnetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.DeallocateFor(&_RootNet.TransactOpts, user, agent, subnetId, amount, deadline, v, r, s)
}

// DeallocateFor is a paid mutator transaction binding the contract method 0x10fe1208.
//
// Solidity: function deallocateFor(address user, address agent, uint256 subnetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetTransactorSession) DeallocateFor(user common.Address, agent common.Address, subnetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.DeallocateFor(&_RootNet.TransactOpts, user, agent, subnetId, amount, deadline, v, r, s)
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

// InitializeRegistry is a paid mutator transaction binding the contract method 0x0672c4be.
//
// Solidity: function initializeRegistry(address awpToken_, address subnetNFT_, address alphaTokenFactory_, address awpEmission_, address lpManager_, address accessManager_, address stakingVault_, address stakeNFT_, address defaultSubnetManagerImpl_) returns()
func (_RootNet *RootNetTransactor) InitializeRegistry(opts *bind.TransactOpts, awpToken_ common.Address, subnetNFT_ common.Address, alphaTokenFactory_ common.Address, awpEmission_ common.Address, lpManager_ common.Address, accessManager_ common.Address, stakingVault_ common.Address, stakeNFT_ common.Address, defaultSubnetManagerImpl_ common.Address) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "initializeRegistry", awpToken_, subnetNFT_, alphaTokenFactory_, awpEmission_, lpManager_, accessManager_, stakingVault_, stakeNFT_, defaultSubnetManagerImpl_)
}

// InitializeRegistry is a paid mutator transaction binding the contract method 0x0672c4be.
//
// Solidity: function initializeRegistry(address awpToken_, address subnetNFT_, address alphaTokenFactory_, address awpEmission_, address lpManager_, address accessManager_, address stakingVault_, address stakeNFT_, address defaultSubnetManagerImpl_) returns()
func (_RootNet *RootNetSession) InitializeRegistry(awpToken_ common.Address, subnetNFT_ common.Address, alphaTokenFactory_ common.Address, awpEmission_ common.Address, lpManager_ common.Address, accessManager_ common.Address, stakingVault_ common.Address, stakeNFT_ common.Address, defaultSubnetManagerImpl_ common.Address) (*types.Transaction, error) {
	return _RootNet.Contract.InitializeRegistry(&_RootNet.TransactOpts, awpToken_, subnetNFT_, alphaTokenFactory_, awpEmission_, lpManager_, accessManager_, stakingVault_, stakeNFT_, defaultSubnetManagerImpl_)
}

// InitializeRegistry is a paid mutator transaction binding the contract method 0x0672c4be.
//
// Solidity: function initializeRegistry(address awpToken_, address subnetNFT_, address alphaTokenFactory_, address awpEmission_, address lpManager_, address accessManager_, address stakingVault_, address stakeNFT_, address defaultSubnetManagerImpl_) returns()
func (_RootNet *RootNetTransactorSession) InitializeRegistry(awpToken_ common.Address, subnetNFT_ common.Address, alphaTokenFactory_ common.Address, awpEmission_ common.Address, lpManager_ common.Address, accessManager_ common.Address, stakingVault_ common.Address, stakeNFT_ common.Address, defaultSubnetManagerImpl_ common.Address) (*types.Transaction, error) {
	return _RootNet.Contract.InitializeRegistry(&_RootNet.TransactOpts, awpToken_, subnetNFT_, alphaTokenFactory_, awpEmission_, lpManager_, accessManager_, stakingVault_, stakeNFT_, defaultSubnetManagerImpl_)
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

// RegisterSubnet is a paid mutator transaction binding the contract method 0x5f24898d.
//
// Solidity: function registerSubnet((string,string,address,bytes32,uint128,string) params) returns(uint256)
func (_RootNet *RootNetTransactor) RegisterSubnet(opts *bind.TransactOpts, params IRootNetSubnetParams) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "registerSubnet", params)
}

// RegisterSubnet is a paid mutator transaction binding the contract method 0x5f24898d.
//
// Solidity: function registerSubnet((string,string,address,bytes32,uint128,string) params) returns(uint256)
func (_RootNet *RootNetSession) RegisterSubnet(params IRootNetSubnetParams) (*types.Transaction, error) {
	return _RootNet.Contract.RegisterSubnet(&_RootNet.TransactOpts, params)
}

// RegisterSubnet is a paid mutator transaction binding the contract method 0x5f24898d.
//
// Solidity: function registerSubnet((string,string,address,bytes32,uint128,string) params) returns(uint256)
func (_RootNet *RootNetTransactorSession) RegisterSubnet(params IRootNetSubnetParams) (*types.Transaction, error) {
	return _RootNet.Contract.RegisterSubnet(&_RootNet.TransactOpts, params)
}

// RegisterSubnetFor is a paid mutator transaction binding the contract method 0x1aa3ff5a.
//
// Solidity: function registerSubnetFor(address user, (string,string,address,bytes32,uint128,string) params, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_RootNet *RootNetTransactor) RegisterSubnetFor(opts *bind.TransactOpts, user common.Address, params IRootNetSubnetParams, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "registerSubnetFor", user, params, deadline, v, r, s)
}

// RegisterSubnetFor is a paid mutator transaction binding the contract method 0x1aa3ff5a.
//
// Solidity: function registerSubnetFor(address user, (string,string,address,bytes32,uint128,string) params, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_RootNet *RootNetSession) RegisterSubnetFor(user common.Address, params IRootNetSubnetParams, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.RegisterSubnetFor(&_RootNet.TransactOpts, user, params, deadline, v, r, s)
}

// RegisterSubnetFor is a paid mutator transaction binding the contract method 0x1aa3ff5a.
//
// Solidity: function registerSubnetFor(address user, (string,string,address,bytes32,uint128,string) params, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_RootNet *RootNetTransactorSession) RegisterSubnetFor(user common.Address, params IRootNetSubnetParams, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.RegisterSubnetFor(&_RootNet.TransactOpts, user, params, deadline, v, r, s)
}

// RegisterSubnetForWithPermit is a paid mutator transaction binding the contract method 0xedf12231.
//
// Solidity: function registerSubnetForWithPermit(address user, (string,string,address,bytes32,uint128,string) params, uint256 deadline, uint8 permitV, bytes32 permitR, bytes32 permitS, uint8 registerV, bytes32 registerR, bytes32 registerS) returns(uint256)
func (_RootNet *RootNetTransactor) RegisterSubnetForWithPermit(opts *bind.TransactOpts, user common.Address, params IRootNetSubnetParams, deadline *big.Int, permitV uint8, permitR [32]byte, permitS [32]byte, registerV uint8, registerR [32]byte, registerS [32]byte) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "registerSubnetForWithPermit", user, params, deadline, permitV, permitR, permitS, registerV, registerR, registerS)
}

// RegisterSubnetForWithPermit is a paid mutator transaction binding the contract method 0xedf12231.
//
// Solidity: function registerSubnetForWithPermit(address user, (string,string,address,bytes32,uint128,string) params, uint256 deadline, uint8 permitV, bytes32 permitR, bytes32 permitS, uint8 registerV, bytes32 registerR, bytes32 registerS) returns(uint256)
func (_RootNet *RootNetSession) RegisterSubnetForWithPermit(user common.Address, params IRootNetSubnetParams, deadline *big.Int, permitV uint8, permitR [32]byte, permitS [32]byte, registerV uint8, registerR [32]byte, registerS [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.RegisterSubnetForWithPermit(&_RootNet.TransactOpts, user, params, deadline, permitV, permitR, permitS, registerV, registerR, registerS)
}

// RegisterSubnetForWithPermit is a paid mutator transaction binding the contract method 0xedf12231.
//
// Solidity: function registerSubnetForWithPermit(address user, (string,string,address,bytes32,uint128,string) params, uint256 deadline, uint8 permitV, bytes32 permitR, bytes32 permitS, uint8 registerV, bytes32 registerR, bytes32 registerS) returns(uint256)
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

// SetRewardRecipientFor is a paid mutator transaction binding the contract method 0xce1297db.
//
// Solidity: function setRewardRecipientFor(address user, address recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetTransactor) SetRewardRecipientFor(opts *bind.TransactOpts, user common.Address, recipient common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.contract.Transact(opts, "setRewardRecipientFor", user, recipient, deadline, v, r, s)
}

// SetRewardRecipientFor is a paid mutator transaction binding the contract method 0xce1297db.
//
// Solidity: function setRewardRecipientFor(address user, address recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetSession) SetRewardRecipientFor(user common.Address, recipient common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.SetRewardRecipientFor(&_RootNet.TransactOpts, user, recipient, deadline, v, r, s)
}

// SetRewardRecipientFor is a paid mutator transaction binding the contract method 0xce1297db.
//
// Solidity: function setRewardRecipientFor(address user, address recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_RootNet *RootNetTransactorSession) SetRewardRecipientFor(user common.Address, recipient common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _RootNet.Contract.SetRewardRecipientFor(&_RootNet.TransactOpts, user, recipient, deadline, v, r, s)
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

// RootNetAlphaTokenFactoryUpdatedIterator is returned from FilterAlphaTokenFactoryUpdated and is used to iterate over the raw logs and unpacked data for AlphaTokenFactoryUpdated events raised by the RootNet contract.
type RootNetAlphaTokenFactoryUpdatedIterator struct {
	Event *RootNetAlphaTokenFactoryUpdated // Event containing the contract specifics and raw log

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
func (it *RootNetAlphaTokenFactoryUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetAlphaTokenFactoryUpdated)
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
		it.Event = new(RootNetAlphaTokenFactoryUpdated)
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
func (it *RootNetAlphaTokenFactoryUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetAlphaTokenFactoryUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetAlphaTokenFactoryUpdated represents a AlphaTokenFactoryUpdated event raised by the RootNet contract.
type RootNetAlphaTokenFactoryUpdated struct {
	NewFactory common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAlphaTokenFactoryUpdated is a free log retrieval operation binding the contract event 0xca3b5054bdfbf81973dd36029b7ef8c5479d0739433700df6b2e6d690ead4a3e.
//
// Solidity: event AlphaTokenFactoryUpdated(address indexed newFactory)
func (_RootNet *RootNetFilterer) FilterAlphaTokenFactoryUpdated(opts *bind.FilterOpts, newFactory []common.Address) (*RootNetAlphaTokenFactoryUpdatedIterator, error) {

	var newFactoryRule []interface{}
	for _, newFactoryItem := range newFactory {
		newFactoryRule = append(newFactoryRule, newFactoryItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "AlphaTokenFactoryUpdated", newFactoryRule)
	if err != nil {
		return nil, err
	}
	return &RootNetAlphaTokenFactoryUpdatedIterator{contract: _RootNet.contract, event: "AlphaTokenFactoryUpdated", logs: logs, sub: sub}, nil
}

// WatchAlphaTokenFactoryUpdated is a free log subscription operation binding the contract event 0xca3b5054bdfbf81973dd36029b7ef8c5479d0739433700df6b2e6d690ead4a3e.
//
// Solidity: event AlphaTokenFactoryUpdated(address indexed newFactory)
func (_RootNet *RootNetFilterer) WatchAlphaTokenFactoryUpdated(opts *bind.WatchOpts, sink chan<- *RootNetAlphaTokenFactoryUpdated, newFactory []common.Address) (event.Subscription, error) {

	var newFactoryRule []interface{}
	for _, newFactoryItem := range newFactory {
		newFactoryRule = append(newFactoryRule, newFactoryItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "AlphaTokenFactoryUpdated", newFactoryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetAlphaTokenFactoryUpdated)
				if err := _RootNet.contract.UnpackLog(event, "AlphaTokenFactoryUpdated", log); err != nil {
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

// ParseAlphaTokenFactoryUpdated is a log parse operation binding the contract event 0xca3b5054bdfbf81973dd36029b7ef8c5479d0739433700df6b2e6d690ead4a3e.
//
// Solidity: event AlphaTokenFactoryUpdated(address indexed newFactory)
func (_RootNet *RootNetFilterer) ParseAlphaTokenFactoryUpdated(log types.Log) (*RootNetAlphaTokenFactoryUpdated, error) {
	event := new(RootNetAlphaTokenFactoryUpdated)
	if err := _RootNet.contract.UnpackLog(event, "AlphaTokenFactoryUpdated", log); err != nil {
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

// RootNetDefaultSubnetManagerImplUpdatedIterator is returned from FilterDefaultSubnetManagerImplUpdated and is used to iterate over the raw logs and unpacked data for DefaultSubnetManagerImplUpdated events raised by the RootNet contract.
type RootNetDefaultSubnetManagerImplUpdatedIterator struct {
	Event *RootNetDefaultSubnetManagerImplUpdated // Event containing the contract specifics and raw log

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
func (it *RootNetDefaultSubnetManagerImplUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetDefaultSubnetManagerImplUpdated)
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
		it.Event = new(RootNetDefaultSubnetManagerImplUpdated)
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
func (it *RootNetDefaultSubnetManagerImplUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetDefaultSubnetManagerImplUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetDefaultSubnetManagerImplUpdated represents a DefaultSubnetManagerImplUpdated event raised by the RootNet contract.
type RootNetDefaultSubnetManagerImplUpdated struct {
	NewImpl common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterDefaultSubnetManagerImplUpdated is a free log retrieval operation binding the contract event 0xa37cb79f631c6bb2a11d965d06cce40e3c936eba1649879b8ffa233c0219f949.
//
// Solidity: event DefaultSubnetManagerImplUpdated(address indexed newImpl)
func (_RootNet *RootNetFilterer) FilterDefaultSubnetManagerImplUpdated(opts *bind.FilterOpts, newImpl []common.Address) (*RootNetDefaultSubnetManagerImplUpdatedIterator, error) {

	var newImplRule []interface{}
	for _, newImplItem := range newImpl {
		newImplRule = append(newImplRule, newImplItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "DefaultSubnetManagerImplUpdated", newImplRule)
	if err != nil {
		return nil, err
	}
	return &RootNetDefaultSubnetManagerImplUpdatedIterator{contract: _RootNet.contract, event: "DefaultSubnetManagerImplUpdated", logs: logs, sub: sub}, nil
}

// WatchDefaultSubnetManagerImplUpdated is a free log subscription operation binding the contract event 0xa37cb79f631c6bb2a11d965d06cce40e3c936eba1649879b8ffa233c0219f949.
//
// Solidity: event DefaultSubnetManagerImplUpdated(address indexed newImpl)
func (_RootNet *RootNetFilterer) WatchDefaultSubnetManagerImplUpdated(opts *bind.WatchOpts, sink chan<- *RootNetDefaultSubnetManagerImplUpdated, newImpl []common.Address) (event.Subscription, error) {

	var newImplRule []interface{}
	for _, newImplItem := range newImpl {
		newImplRule = append(newImplRule, newImplItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "DefaultSubnetManagerImplUpdated", newImplRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetDefaultSubnetManagerImplUpdated)
				if err := _RootNet.contract.UnpackLog(event, "DefaultSubnetManagerImplUpdated", log); err != nil {
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

// ParseDefaultSubnetManagerImplUpdated is a log parse operation binding the contract event 0xa37cb79f631c6bb2a11d965d06cce40e3c936eba1649879b8ffa233c0219f949.
//
// Solidity: event DefaultSubnetManagerImplUpdated(address indexed newImpl)
func (_RootNet *RootNetFilterer) ParseDefaultSubnetManagerImplUpdated(log types.Log) (*RootNetDefaultSubnetManagerImplUpdated, error) {
	event := new(RootNetDefaultSubnetManagerImplUpdated)
	if err := _RootNet.contract.UnpackLog(event, "DefaultSubnetManagerImplUpdated", log); err != nil {
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

// RootNetGuardianUpdatedIterator is returned from FilterGuardianUpdated and is used to iterate over the raw logs and unpacked data for GuardianUpdated events raised by the RootNet contract.
type RootNetGuardianUpdatedIterator struct {
	Event *RootNetGuardianUpdated // Event containing the contract specifics and raw log

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
func (it *RootNetGuardianUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetGuardianUpdated)
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
		it.Event = new(RootNetGuardianUpdated)
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
func (it *RootNetGuardianUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetGuardianUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetGuardianUpdated represents a GuardianUpdated event raised by the RootNet contract.
type RootNetGuardianUpdated struct {
	NewGuardian common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterGuardianUpdated is a free log retrieval operation binding the contract event 0x6bb7ff33e730289800c62ad882105a144a74010d2bdbb9a942544a3005ad55bf.
//
// Solidity: event GuardianUpdated(address indexed newGuardian)
func (_RootNet *RootNetFilterer) FilterGuardianUpdated(opts *bind.FilterOpts, newGuardian []common.Address) (*RootNetGuardianUpdatedIterator, error) {

	var newGuardianRule []interface{}
	for _, newGuardianItem := range newGuardian {
		newGuardianRule = append(newGuardianRule, newGuardianItem)
	}

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "GuardianUpdated", newGuardianRule)
	if err != nil {
		return nil, err
	}
	return &RootNetGuardianUpdatedIterator{contract: _RootNet.contract, event: "GuardianUpdated", logs: logs, sub: sub}, nil
}

// WatchGuardianUpdated is a free log subscription operation binding the contract event 0x6bb7ff33e730289800c62ad882105a144a74010d2bdbb9a942544a3005ad55bf.
//
// Solidity: event GuardianUpdated(address indexed newGuardian)
func (_RootNet *RootNetFilterer) WatchGuardianUpdated(opts *bind.WatchOpts, sink chan<- *RootNetGuardianUpdated, newGuardian []common.Address) (event.Subscription, error) {

	var newGuardianRule []interface{}
	for _, newGuardianItem := range newGuardian {
		newGuardianRule = append(newGuardianRule, newGuardianItem)
	}

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "GuardianUpdated", newGuardianRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetGuardianUpdated)
				if err := _RootNet.contract.UnpackLog(event, "GuardianUpdated", log); err != nil {
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
func (_RootNet *RootNetFilterer) ParseGuardianUpdated(log types.Log) (*RootNetGuardianUpdated, error) {
	event := new(RootNetGuardianUpdated)
	if err := _RootNet.contract.UnpackLog(event, "GuardianUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetImmunityPeriodUpdatedIterator is returned from FilterImmunityPeriodUpdated and is used to iterate over the raw logs and unpacked data for ImmunityPeriodUpdated events raised by the RootNet contract.
type RootNetImmunityPeriodUpdatedIterator struct {
	Event *RootNetImmunityPeriodUpdated // Event containing the contract specifics and raw log

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
func (it *RootNetImmunityPeriodUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetImmunityPeriodUpdated)
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
		it.Event = new(RootNetImmunityPeriodUpdated)
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
func (it *RootNetImmunityPeriodUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetImmunityPeriodUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetImmunityPeriodUpdated represents a ImmunityPeriodUpdated event raised by the RootNet contract.
type RootNetImmunityPeriodUpdated struct {
	NewPeriod *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterImmunityPeriodUpdated is a free log retrieval operation binding the contract event 0x49b186851943e5bbcefec9411c3238262c6e102e4000142f8f060143d1b8724c.
//
// Solidity: event ImmunityPeriodUpdated(uint256 newPeriod)
func (_RootNet *RootNetFilterer) FilterImmunityPeriodUpdated(opts *bind.FilterOpts) (*RootNetImmunityPeriodUpdatedIterator, error) {

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "ImmunityPeriodUpdated")
	if err != nil {
		return nil, err
	}
	return &RootNetImmunityPeriodUpdatedIterator{contract: _RootNet.contract, event: "ImmunityPeriodUpdated", logs: logs, sub: sub}, nil
}

// WatchImmunityPeriodUpdated is a free log subscription operation binding the contract event 0x49b186851943e5bbcefec9411c3238262c6e102e4000142f8f060143d1b8724c.
//
// Solidity: event ImmunityPeriodUpdated(uint256 newPeriod)
func (_RootNet *RootNetFilterer) WatchImmunityPeriodUpdated(opts *bind.WatchOpts, sink chan<- *RootNetImmunityPeriodUpdated) (event.Subscription, error) {

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "ImmunityPeriodUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetImmunityPeriodUpdated)
				if err := _RootNet.contract.UnpackLog(event, "ImmunityPeriodUpdated", log); err != nil {
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

// ParseImmunityPeriodUpdated is a log parse operation binding the contract event 0x49b186851943e5bbcefec9411c3238262c6e102e4000142f8f060143d1b8724c.
//
// Solidity: event ImmunityPeriodUpdated(uint256 newPeriod)
func (_RootNet *RootNetFilterer) ParseImmunityPeriodUpdated(log types.Log) (*RootNetImmunityPeriodUpdated, error) {
	event := new(RootNetImmunityPeriodUpdated)
	if err := _RootNet.contract.UnpackLog(event, "ImmunityPeriodUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RootNetInitialAlphaPriceUpdatedIterator is returned from FilterInitialAlphaPriceUpdated and is used to iterate over the raw logs and unpacked data for InitialAlphaPriceUpdated events raised by the RootNet contract.
type RootNetInitialAlphaPriceUpdatedIterator struct {
	Event *RootNetInitialAlphaPriceUpdated // Event containing the contract specifics and raw log

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
func (it *RootNetInitialAlphaPriceUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootNetInitialAlphaPriceUpdated)
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
		it.Event = new(RootNetInitialAlphaPriceUpdated)
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
func (it *RootNetInitialAlphaPriceUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootNetInitialAlphaPriceUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootNetInitialAlphaPriceUpdated represents a InitialAlphaPriceUpdated event raised by the RootNet contract.
type RootNetInitialAlphaPriceUpdated struct {
	NewPrice *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterInitialAlphaPriceUpdated is a free log retrieval operation binding the contract event 0xab7ee876750d22d253d0b38988caea5f6285a832697e4889d9beb36515dde34e.
//
// Solidity: event InitialAlphaPriceUpdated(uint256 newPrice)
func (_RootNet *RootNetFilterer) FilterInitialAlphaPriceUpdated(opts *bind.FilterOpts) (*RootNetInitialAlphaPriceUpdatedIterator, error) {

	logs, sub, err := _RootNet.contract.FilterLogs(opts, "InitialAlphaPriceUpdated")
	if err != nil {
		return nil, err
	}
	return &RootNetInitialAlphaPriceUpdatedIterator{contract: _RootNet.contract, event: "InitialAlphaPriceUpdated", logs: logs, sub: sub}, nil
}

// WatchInitialAlphaPriceUpdated is a free log subscription operation binding the contract event 0xab7ee876750d22d253d0b38988caea5f6285a832697e4889d9beb36515dde34e.
//
// Solidity: event InitialAlphaPriceUpdated(uint256 newPrice)
func (_RootNet *RootNetFilterer) WatchInitialAlphaPriceUpdated(opts *bind.WatchOpts, sink chan<- *RootNetInitialAlphaPriceUpdated) (event.Subscription, error) {

	logs, sub, err := _RootNet.contract.WatchLogs(opts, "InitialAlphaPriceUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootNetInitialAlphaPriceUpdated)
				if err := _RootNet.contract.UnpackLog(event, "InitialAlphaPriceUpdated", log); err != nil {
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

// ParseInitialAlphaPriceUpdated is a log parse operation binding the contract event 0xab7ee876750d22d253d0b38988caea5f6285a832697e4889d9beb36515dde34e.
//
// Solidity: event InitialAlphaPriceUpdated(uint256 newPrice)
func (_RootNet *RootNetFilterer) ParseInitialAlphaPriceUpdated(log types.Log) (*RootNetInitialAlphaPriceUpdated, error) {
	event := new(RootNetInitialAlphaPriceUpdated)
	if err := _RootNet.contract.UnpackLog(event, "InitialAlphaPriceUpdated", log); err != nil {
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
	SubnetId      *big.Int
	Owner         common.Address
	Name          string
	Symbol        string
	SubnetManager common.Address
	AlphaToken    common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSubnetRegistered is a free log retrieval operation binding the contract event 0x8bf2da7b4bb5b09423a8727058489d29da8c78ca3d15f540facdf1ad5dbd09d1.
//
// Solidity: event SubnetRegistered(uint256 indexed subnetId, address indexed owner, string name, string symbol, address subnetManager, address alphaToken)
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

// WatchSubnetRegistered is a free log subscription operation binding the contract event 0x8bf2da7b4bb5b09423a8727058489d29da8c78ca3d15f540facdf1ad5dbd09d1.
//
// Solidity: event SubnetRegistered(uint256 indexed subnetId, address indexed owner, string name, string symbol, address subnetManager, address alphaToken)
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

// ParseSubnetRegistered is a log parse operation binding the contract event 0x8bf2da7b4bb5b09423a8727058489d29da8c78ca3d15f540facdf1ad5dbd09d1.
//
// Solidity: event SubnetRegistered(uint256 indexed subnetId, address indexed owner, string name, string symbol, address subnetManager, address alphaToken)
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
