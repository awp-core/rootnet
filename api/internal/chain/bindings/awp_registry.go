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

// AWPRegistryAgentInfo is an auto generated low-level Go binding around an user-defined struct.
type AWPRegistryAgentInfo struct {
	Root            common.Address
	IsValid         bool
	Stake           *big.Int
	RewardRecipient common.Address
}

// IAWPRegistrySubnetFullInfo is an auto generated low-level Go binding around an user-defined struct.
type IAWPRegistrySubnetFullInfo struct {
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

// IAWPRegistrySubnetInfo is an auto generated low-level Go binding around an user-defined struct.
type IAWPRegistrySubnetInfo struct {
	LpPool      [32]byte
	Status      uint8
	CreatedAt   uint64
	ActivatedAt uint64
}

// IAWPRegistrySubnetParams is an auto generated low-level Go binding around an user-defined struct.
type IAWPRegistrySubnetParams struct {
	Name          string
	Symbol        string
	SubnetManager common.Address
	Salt          [32]byte
	MinStake      *big.Int
	SkillsURI     string
}

// AWPRegistryMetaData contains all meta data concerning the AWPRegistry contract.
var AWPRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"deployer_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"treasury_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"guardian_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"INITIAL_ALPHA_MINT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_ACTIVE_SUBNETS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"activateSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"activateSubnetFor\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allocate\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allocateFor\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"alphaTokenFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"awpEmission\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"awpToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"banSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bind\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bindFor\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"boundTo\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deallocate\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deallocateFor\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultSubnetManagerImpl\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"delegates\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deregisterSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"dexConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveSubnetCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveSubnetIdAt\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAgentInfo\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structAWPRegistry.AgentInfo\",\"components\":[{\"name\":\"root\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isValid\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"stake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"rewardRecipient\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAgentsInfo\",\"inputs\":[{\"name\":\"agents\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structAWPRegistry.AgentInfo[]\",\"components\":[{\"name\":\"root\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isValid\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"stake\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"rewardRecipient\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIAWPRegistry.SubnetInfo\",\"components\":[{\"name\":\"lpPool\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumIAWPRegistry.SubnetStatus\"},{\"name\":\"createdAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"activatedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSubnetFull\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIAWPRegistry.SubnetFullInfo\",\"components\":[{\"name\":\"subnetManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"alphaToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lpPool\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumIAWPRegistry.SubnetStatus\"},{\"name\":\"createdAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"activatedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantDelegate\",\"inputs\":[{\"name\":\"delegate\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"guardian\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"immunityPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialAlphaPrice\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initializeRegistry\",\"inputs\":[{\"name\":\"awpToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subnetNFT_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"alphaTokenFactory_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"awpEmission_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lpManager_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakingVault_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakeNFT_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"defaultSubnetManagerImpl_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"dexConfig_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isRegistered\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isSubnetActive\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lpManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nextSubnetId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nonces\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pauseSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"reallocate\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromAgent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fromSubnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"toAgent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"toSubnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"recipient\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"register\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerSubnet\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIAWPRegistry.SubnetParams\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"subnetManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerSubnetFor\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIAWPRegistry.SubnetParams\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"subnetManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerSubnetForWithPermit\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIAWPRegistry.SubnetParams\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"subnetManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"permitV\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"permitR\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"permitS\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"registerV\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"registerR\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"registerS\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registeredCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registryInitialized\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"resolveRecipient\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"resumeSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeDelegate\",\"inputs\":[{\"name\":\"delegate\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAlphaTokenFactory\",\"inputs\":[{\"name\":\"factory\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDexConfig\",\"inputs\":[{\"name\":\"dexConfig_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setGuardian\",\"inputs\":[{\"name\":\"g\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setImmunityPeriod\",\"inputs\":[{\"name\":\"p\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setInitialAlphaPrice\",\"inputs\":[{\"name\":\"price\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRecipient\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRecipientFor\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"v\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSubnetManagerImpl\",\"inputs\":[{\"name\":\"impl\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakeNFT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stakingVault\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"subnetNFT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"subnets\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"lpPool\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumIAWPRegistry.SubnetStatus\"},{\"name\":\"createdAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"activatedAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"treasury\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unbanSubnet\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unbind\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Allocated\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AlphaTokenFactoryUpdated\",\"inputs\":[{\"name\":\"newFactory\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Bound\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"target\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Deallocated\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DefaultSubnetManagerImplUpdated\",\"inputs\":[{\"name\":\"newImpl\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DelegateGranted\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"delegate\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DelegateRevoked\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"delegate\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GuardianUpdated\",\"inputs\":[{\"name\":\"newGuardian\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ImmunityPeriodUpdated\",\"inputs\":[{\"name\":\"newPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InitialAlphaPriceUpdated\",\"inputs\":[{\"name\":\"newPrice\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LPCreated\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"poolId\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"awpAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"alphaAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Reallocated\",\"inputs\":[{\"name\":\"staker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fromAgent\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"fromSubnet\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"toAgent\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"toSubnet\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RecipientSet\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetActivated\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetBanned\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetDeregistered\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetPaused\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetRegistered\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"subnetManager\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"alphaToken\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"minStake\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"skillsURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetResumed\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubnetUnbanned\",\"inputs\":[{\"name\":\"subnetId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unbound\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UserRegistered\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotRevokeSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ChainTooLong\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CycleDetected\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpiredSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ImmunityNotExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidShortString\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSubnetParams\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSubnetStatus\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxActiveSubnetsReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAuthorized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotDeployer\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotGuardian\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotTimelock\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PriceTooHigh\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PriceTooLow\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"StringTooLong\",\"inputs\":[{\"name\":\"str\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"SubnetManagerRequired\",\"inputs\":[]}]",
}

// AWPRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use AWPRegistryMetaData.ABI instead.
var AWPRegistryABI = AWPRegistryMetaData.ABI

// AWPRegistry is an auto generated Go binding around an Ethereum contract.
type AWPRegistry struct {
	AWPRegistryCaller     // Read-only binding to the contract
	AWPRegistryTransactor // Write-only binding to the contract
	AWPRegistryFilterer   // Log filterer for contract events
}

// AWPRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type AWPRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AWPRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AWPRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AWPRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AWPRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AWPRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AWPRegistrySession struct {
	Contract     *AWPRegistry      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AWPRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AWPRegistryCallerSession struct {
	Contract *AWPRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// AWPRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AWPRegistryTransactorSession struct {
	Contract     *AWPRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// AWPRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type AWPRegistryRaw struct {
	Contract *AWPRegistry // Generic contract binding to access the raw methods on
}

// AWPRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AWPRegistryCallerRaw struct {
	Contract *AWPRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// AWPRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AWPRegistryTransactorRaw struct {
	Contract *AWPRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAWPRegistry creates a new instance of AWPRegistry, bound to a specific deployed contract.
func NewAWPRegistry(address common.Address, backend bind.ContractBackend) (*AWPRegistry, error) {
	contract, err := bindAWPRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AWPRegistry{AWPRegistryCaller: AWPRegistryCaller{contract: contract}, AWPRegistryTransactor: AWPRegistryTransactor{contract: contract}, AWPRegistryFilterer: AWPRegistryFilterer{contract: contract}}, nil
}

// NewAWPRegistryCaller creates a new read-only instance of AWPRegistry, bound to a specific deployed contract.
func NewAWPRegistryCaller(address common.Address, caller bind.ContractCaller) (*AWPRegistryCaller, error) {
	contract, err := bindAWPRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AWPRegistryCaller{contract: contract}, nil
}

// NewAWPRegistryTransactor creates a new write-only instance of AWPRegistry, bound to a specific deployed contract.
func NewAWPRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*AWPRegistryTransactor, error) {
	contract, err := bindAWPRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AWPRegistryTransactor{contract: contract}, nil
}

// NewAWPRegistryFilterer creates a new log filterer instance of AWPRegistry, bound to a specific deployed contract.
func NewAWPRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*AWPRegistryFilterer, error) {
	contract, err := bindAWPRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AWPRegistryFilterer{contract: contract}, nil
}

// bindAWPRegistry binds a generic wrapper to an already deployed contract.
func bindAWPRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AWPRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AWPRegistry *AWPRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AWPRegistry.Contract.AWPRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AWPRegistry *AWPRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AWPRegistry.Contract.AWPRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AWPRegistry *AWPRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AWPRegistry.Contract.AWPRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AWPRegistry *AWPRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AWPRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AWPRegistry *AWPRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AWPRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AWPRegistry *AWPRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AWPRegistry.Contract.contract.Transact(opts, method, params...)
}

// INITIALALPHAMINT is a free data retrieval call binding the contract method 0xb400555a.
//
// Solidity: function INITIAL_ALPHA_MINT() view returns(uint256)
func (_AWPRegistry *AWPRegistryCaller) INITIALALPHAMINT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "INITIAL_ALPHA_MINT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// INITIALALPHAMINT is a free data retrieval call binding the contract method 0xb400555a.
//
// Solidity: function INITIAL_ALPHA_MINT() view returns(uint256)
func (_AWPRegistry *AWPRegistrySession) INITIALALPHAMINT() (*big.Int, error) {
	return _AWPRegistry.Contract.INITIALALPHAMINT(&_AWPRegistry.CallOpts)
}

// INITIALALPHAMINT is a free data retrieval call binding the contract method 0xb400555a.
//
// Solidity: function INITIAL_ALPHA_MINT() view returns(uint256)
func (_AWPRegistry *AWPRegistryCallerSession) INITIALALPHAMINT() (*big.Int, error) {
	return _AWPRegistry.Contract.INITIALALPHAMINT(&_AWPRegistry.CallOpts)
}

// MAXACTIVESUBNETS is a free data retrieval call binding the contract method 0xbe65e4c2.
//
// Solidity: function MAX_ACTIVE_SUBNETS() view returns(uint128)
func (_AWPRegistry *AWPRegistryCaller) MAXACTIVESUBNETS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "MAX_ACTIVE_SUBNETS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXACTIVESUBNETS is a free data retrieval call binding the contract method 0xbe65e4c2.
//
// Solidity: function MAX_ACTIVE_SUBNETS() view returns(uint128)
func (_AWPRegistry *AWPRegistrySession) MAXACTIVESUBNETS() (*big.Int, error) {
	return _AWPRegistry.Contract.MAXACTIVESUBNETS(&_AWPRegistry.CallOpts)
}

// MAXACTIVESUBNETS is a free data retrieval call binding the contract method 0xbe65e4c2.
//
// Solidity: function MAX_ACTIVE_SUBNETS() view returns(uint128)
func (_AWPRegistry *AWPRegistryCallerSession) MAXACTIVESUBNETS() (*big.Int, error) {
	return _AWPRegistry.Contract.MAXACTIVESUBNETS(&_AWPRegistry.CallOpts)
}

// AlphaTokenFactory is a free data retrieval call binding the contract method 0xc1e0c9e7.
//
// Solidity: function alphaTokenFactory() view returns(address)
func (_AWPRegistry *AWPRegistryCaller) AlphaTokenFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "alphaTokenFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AlphaTokenFactory is a free data retrieval call binding the contract method 0xc1e0c9e7.
//
// Solidity: function alphaTokenFactory() view returns(address)
func (_AWPRegistry *AWPRegistrySession) AlphaTokenFactory() (common.Address, error) {
	return _AWPRegistry.Contract.AlphaTokenFactory(&_AWPRegistry.CallOpts)
}

// AlphaTokenFactory is a free data retrieval call binding the contract method 0xc1e0c9e7.
//
// Solidity: function alphaTokenFactory() view returns(address)
func (_AWPRegistry *AWPRegistryCallerSession) AlphaTokenFactory() (common.Address, error) {
	return _AWPRegistry.Contract.AlphaTokenFactory(&_AWPRegistry.CallOpts)
}

// AwpEmission is a free data retrieval call binding the contract method 0x67b26ba6.
//
// Solidity: function awpEmission() view returns(address)
func (_AWPRegistry *AWPRegistryCaller) AwpEmission(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "awpEmission")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpEmission is a free data retrieval call binding the contract method 0x67b26ba6.
//
// Solidity: function awpEmission() view returns(address)
func (_AWPRegistry *AWPRegistrySession) AwpEmission() (common.Address, error) {
	return _AWPRegistry.Contract.AwpEmission(&_AWPRegistry.CallOpts)
}

// AwpEmission is a free data retrieval call binding the contract method 0x67b26ba6.
//
// Solidity: function awpEmission() view returns(address)
func (_AWPRegistry *AWPRegistryCallerSession) AwpEmission() (common.Address, error) {
	return _AWPRegistry.Contract.AwpEmission(&_AWPRegistry.CallOpts)
}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_AWPRegistry *AWPRegistryCaller) AwpToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "awpToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_AWPRegistry *AWPRegistrySession) AwpToken() (common.Address, error) {
	return _AWPRegistry.Contract.AwpToken(&_AWPRegistry.CallOpts)
}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_AWPRegistry *AWPRegistryCallerSession) AwpToken() (common.Address, error) {
	return _AWPRegistry.Contract.AwpToken(&_AWPRegistry.CallOpts)
}

// BoundTo is a free data retrieval call binding the contract method 0xf343e266.
//
// Solidity: function boundTo(address ) view returns(address)
func (_AWPRegistry *AWPRegistryCaller) BoundTo(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "boundTo", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BoundTo is a free data retrieval call binding the contract method 0xf343e266.
//
// Solidity: function boundTo(address ) view returns(address)
func (_AWPRegistry *AWPRegistrySession) BoundTo(arg0 common.Address) (common.Address, error) {
	return _AWPRegistry.Contract.BoundTo(&_AWPRegistry.CallOpts, arg0)
}

// BoundTo is a free data retrieval call binding the contract method 0xf343e266.
//
// Solidity: function boundTo(address ) view returns(address)
func (_AWPRegistry *AWPRegistryCallerSession) BoundTo(arg0 common.Address) (common.Address, error) {
	return _AWPRegistry.Contract.BoundTo(&_AWPRegistry.CallOpts, arg0)
}

// DefaultSubnetManagerImpl is a free data retrieval call binding the contract method 0xf4fda726.
//
// Solidity: function defaultSubnetManagerImpl() view returns(address)
func (_AWPRegistry *AWPRegistryCaller) DefaultSubnetManagerImpl(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "defaultSubnetManagerImpl")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DefaultSubnetManagerImpl is a free data retrieval call binding the contract method 0xf4fda726.
//
// Solidity: function defaultSubnetManagerImpl() view returns(address)
func (_AWPRegistry *AWPRegistrySession) DefaultSubnetManagerImpl() (common.Address, error) {
	return _AWPRegistry.Contract.DefaultSubnetManagerImpl(&_AWPRegistry.CallOpts)
}

// DefaultSubnetManagerImpl is a free data retrieval call binding the contract method 0xf4fda726.
//
// Solidity: function defaultSubnetManagerImpl() view returns(address)
func (_AWPRegistry *AWPRegistryCallerSession) DefaultSubnetManagerImpl() (common.Address, error) {
	return _AWPRegistry.Contract.DefaultSubnetManagerImpl(&_AWPRegistry.CallOpts)
}

// Delegates is a free data retrieval call binding the contract method 0xe5843242.
//
// Solidity: function delegates(address , address ) view returns(bool)
func (_AWPRegistry *AWPRegistryCaller) Delegates(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "delegates", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Delegates is a free data retrieval call binding the contract method 0xe5843242.
//
// Solidity: function delegates(address , address ) view returns(bool)
func (_AWPRegistry *AWPRegistrySession) Delegates(arg0 common.Address, arg1 common.Address) (bool, error) {
	return _AWPRegistry.Contract.Delegates(&_AWPRegistry.CallOpts, arg0, arg1)
}

// Delegates is a free data retrieval call binding the contract method 0xe5843242.
//
// Solidity: function delegates(address , address ) view returns(bool)
func (_AWPRegistry *AWPRegistryCallerSession) Delegates(arg0 common.Address, arg1 common.Address) (bool, error) {
	return _AWPRegistry.Contract.Delegates(&_AWPRegistry.CallOpts, arg0, arg1)
}

// DexConfig is a free data retrieval call binding the contract method 0x38d890d7.
//
// Solidity: function dexConfig() view returns(bytes)
func (_AWPRegistry *AWPRegistryCaller) DexConfig(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "dexConfig")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// DexConfig is a free data retrieval call binding the contract method 0x38d890d7.
//
// Solidity: function dexConfig() view returns(bytes)
func (_AWPRegistry *AWPRegistrySession) DexConfig() ([]byte, error) {
	return _AWPRegistry.Contract.DexConfig(&_AWPRegistry.CallOpts)
}

// DexConfig is a free data retrieval call binding the contract method 0x38d890d7.
//
// Solidity: function dexConfig() view returns(bytes)
func (_AWPRegistry *AWPRegistryCallerSession) DexConfig() ([]byte, error) {
	return _AWPRegistry.Contract.DexConfig(&_AWPRegistry.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_AWPRegistry *AWPRegistryCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "eip712Domain")

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
func (_AWPRegistry *AWPRegistrySession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _AWPRegistry.Contract.Eip712Domain(&_AWPRegistry.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_AWPRegistry *AWPRegistryCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _AWPRegistry.Contract.Eip712Domain(&_AWPRegistry.CallOpts)
}

// GetActiveSubnetCount is a free data retrieval call binding the contract method 0xc6a1a01a.
//
// Solidity: function getActiveSubnetCount() view returns(uint256)
func (_AWPRegistry *AWPRegistryCaller) GetActiveSubnetCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "getActiveSubnetCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActiveSubnetCount is a free data retrieval call binding the contract method 0xc6a1a01a.
//
// Solidity: function getActiveSubnetCount() view returns(uint256)
func (_AWPRegistry *AWPRegistrySession) GetActiveSubnetCount() (*big.Int, error) {
	return _AWPRegistry.Contract.GetActiveSubnetCount(&_AWPRegistry.CallOpts)
}

// GetActiveSubnetCount is a free data retrieval call binding the contract method 0xc6a1a01a.
//
// Solidity: function getActiveSubnetCount() view returns(uint256)
func (_AWPRegistry *AWPRegistryCallerSession) GetActiveSubnetCount() (*big.Int, error) {
	return _AWPRegistry.Contract.GetActiveSubnetCount(&_AWPRegistry.CallOpts)
}

// GetActiveSubnetIdAt is a free data retrieval call binding the contract method 0x38f48a89.
//
// Solidity: function getActiveSubnetIdAt(uint256 index) view returns(uint256)
func (_AWPRegistry *AWPRegistryCaller) GetActiveSubnetIdAt(opts *bind.CallOpts, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "getActiveSubnetIdAt", index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetActiveSubnetIdAt is a free data retrieval call binding the contract method 0x38f48a89.
//
// Solidity: function getActiveSubnetIdAt(uint256 index) view returns(uint256)
func (_AWPRegistry *AWPRegistrySession) GetActiveSubnetIdAt(index *big.Int) (*big.Int, error) {
	return _AWPRegistry.Contract.GetActiveSubnetIdAt(&_AWPRegistry.CallOpts, index)
}

// GetActiveSubnetIdAt is a free data retrieval call binding the contract method 0x38f48a89.
//
// Solidity: function getActiveSubnetIdAt(uint256 index) view returns(uint256)
func (_AWPRegistry *AWPRegistryCallerSession) GetActiveSubnetIdAt(index *big.Int) (*big.Int, error) {
	return _AWPRegistry.Contract.GetActiveSubnetIdAt(&_AWPRegistry.CallOpts, index)
}

// GetAgentInfo is a free data retrieval call binding the contract method 0x168f80f5.
//
// Solidity: function getAgentInfo(address agent, uint256 subnetId) view returns((address,bool,uint256,address))
func (_AWPRegistry *AWPRegistryCaller) GetAgentInfo(opts *bind.CallOpts, agent common.Address, subnetId *big.Int) (AWPRegistryAgentInfo, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "getAgentInfo", agent, subnetId)

	if err != nil {
		return *new(AWPRegistryAgentInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(AWPRegistryAgentInfo)).(*AWPRegistryAgentInfo)

	return out0, err

}

// GetAgentInfo is a free data retrieval call binding the contract method 0x168f80f5.
//
// Solidity: function getAgentInfo(address agent, uint256 subnetId) view returns((address,bool,uint256,address))
func (_AWPRegistry *AWPRegistrySession) GetAgentInfo(agent common.Address, subnetId *big.Int) (AWPRegistryAgentInfo, error) {
	return _AWPRegistry.Contract.GetAgentInfo(&_AWPRegistry.CallOpts, agent, subnetId)
}

// GetAgentInfo is a free data retrieval call binding the contract method 0x168f80f5.
//
// Solidity: function getAgentInfo(address agent, uint256 subnetId) view returns((address,bool,uint256,address))
func (_AWPRegistry *AWPRegistryCallerSession) GetAgentInfo(agent common.Address, subnetId *big.Int) (AWPRegistryAgentInfo, error) {
	return _AWPRegistry.Contract.GetAgentInfo(&_AWPRegistry.CallOpts, agent, subnetId)
}

// GetAgentsInfo is a free data retrieval call binding the contract method 0x4b6f6d67.
//
// Solidity: function getAgentsInfo(address[] agents, uint256 subnetId) view returns((address,bool,uint256,address)[])
func (_AWPRegistry *AWPRegistryCaller) GetAgentsInfo(opts *bind.CallOpts, agents []common.Address, subnetId *big.Int) ([]AWPRegistryAgentInfo, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "getAgentsInfo", agents, subnetId)

	if err != nil {
		return *new([]AWPRegistryAgentInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]AWPRegistryAgentInfo)).(*[]AWPRegistryAgentInfo)

	return out0, err

}

// GetAgentsInfo is a free data retrieval call binding the contract method 0x4b6f6d67.
//
// Solidity: function getAgentsInfo(address[] agents, uint256 subnetId) view returns((address,bool,uint256,address)[])
func (_AWPRegistry *AWPRegistrySession) GetAgentsInfo(agents []common.Address, subnetId *big.Int) ([]AWPRegistryAgentInfo, error) {
	return _AWPRegistry.Contract.GetAgentsInfo(&_AWPRegistry.CallOpts, agents, subnetId)
}

// GetAgentsInfo is a free data retrieval call binding the contract method 0x4b6f6d67.
//
// Solidity: function getAgentsInfo(address[] agents, uint256 subnetId) view returns((address,bool,uint256,address)[])
func (_AWPRegistry *AWPRegistryCallerSession) GetAgentsInfo(agents []common.Address, subnetId *big.Int) ([]AWPRegistryAgentInfo, error) {
	return _AWPRegistry.Contract.GetAgentsInfo(&_AWPRegistry.CallOpts, agents, subnetId)
}

// GetRegistry is a free data retrieval call binding the contract method 0x5ab1bd53.
//
// Solidity: function getRegistry() view returns(address, address, address, address, address, address, address, address, address)
func (_AWPRegistry *AWPRegistryCaller) GetRegistry(opts *bind.CallOpts) (common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "getRegistry")

	if err != nil {
		return *new(common.Address), *new(common.Address), *new(common.Address), *new(common.Address), *new(common.Address), *new(common.Address), *new(common.Address), *new(common.Address), *new(common.Address), err
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

	return out0, out1, out2, out3, out4, out5, out6, out7, out8, err

}

// GetRegistry is a free data retrieval call binding the contract method 0x5ab1bd53.
//
// Solidity: function getRegistry() view returns(address, address, address, address, address, address, address, address, address)
func (_AWPRegistry *AWPRegistrySession) GetRegistry() (common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, error) {
	return _AWPRegistry.Contract.GetRegistry(&_AWPRegistry.CallOpts)
}

// GetRegistry is a free data retrieval call binding the contract method 0x5ab1bd53.
//
// Solidity: function getRegistry() view returns(address, address, address, address, address, address, address, address, address)
func (_AWPRegistry *AWPRegistryCallerSession) GetRegistry() (common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, common.Address, error) {
	return _AWPRegistry.Contract.GetRegistry(&_AWPRegistry.CallOpts)
}

// GetSubnet is a free data retrieval call binding the contract method 0x58ca7504.
//
// Solidity: function getSubnet(uint256 subnetId) view returns((bytes32,uint8,uint64,uint64))
func (_AWPRegistry *AWPRegistryCaller) GetSubnet(opts *bind.CallOpts, subnetId *big.Int) (IAWPRegistrySubnetInfo, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "getSubnet", subnetId)

	if err != nil {
		return *new(IAWPRegistrySubnetInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IAWPRegistrySubnetInfo)).(*IAWPRegistrySubnetInfo)

	return out0, err

}

// GetSubnet is a free data retrieval call binding the contract method 0x58ca7504.
//
// Solidity: function getSubnet(uint256 subnetId) view returns((bytes32,uint8,uint64,uint64))
func (_AWPRegistry *AWPRegistrySession) GetSubnet(subnetId *big.Int) (IAWPRegistrySubnetInfo, error) {
	return _AWPRegistry.Contract.GetSubnet(&_AWPRegistry.CallOpts, subnetId)
}

// GetSubnet is a free data retrieval call binding the contract method 0x58ca7504.
//
// Solidity: function getSubnet(uint256 subnetId) view returns((bytes32,uint8,uint64,uint64))
func (_AWPRegistry *AWPRegistryCallerSession) GetSubnet(subnetId *big.Int) (IAWPRegistrySubnetInfo, error) {
	return _AWPRegistry.Contract.GetSubnet(&_AWPRegistry.CallOpts, subnetId)
}

// GetSubnetFull is a free data retrieval call binding the contract method 0x3de3b247.
//
// Solidity: function getSubnetFull(uint256 subnetId) view returns((address,address,bytes32,uint8,uint64,uint64,string,string,uint128,address))
func (_AWPRegistry *AWPRegistryCaller) GetSubnetFull(opts *bind.CallOpts, subnetId *big.Int) (IAWPRegistrySubnetFullInfo, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "getSubnetFull", subnetId)

	if err != nil {
		return *new(IAWPRegistrySubnetFullInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IAWPRegistrySubnetFullInfo)).(*IAWPRegistrySubnetFullInfo)

	return out0, err

}

// GetSubnetFull is a free data retrieval call binding the contract method 0x3de3b247.
//
// Solidity: function getSubnetFull(uint256 subnetId) view returns((address,address,bytes32,uint8,uint64,uint64,string,string,uint128,address))
func (_AWPRegistry *AWPRegistrySession) GetSubnetFull(subnetId *big.Int) (IAWPRegistrySubnetFullInfo, error) {
	return _AWPRegistry.Contract.GetSubnetFull(&_AWPRegistry.CallOpts, subnetId)
}

// GetSubnetFull is a free data retrieval call binding the contract method 0x3de3b247.
//
// Solidity: function getSubnetFull(uint256 subnetId) view returns((address,address,bytes32,uint8,uint64,uint64,string,string,uint128,address))
func (_AWPRegistry *AWPRegistryCallerSession) GetSubnetFull(subnetId *big.Int) (IAWPRegistrySubnetFullInfo, error) {
	return _AWPRegistry.Contract.GetSubnetFull(&_AWPRegistry.CallOpts, subnetId)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_AWPRegistry *AWPRegistryCaller) Guardian(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "guardian")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_AWPRegistry *AWPRegistrySession) Guardian() (common.Address, error) {
	return _AWPRegistry.Contract.Guardian(&_AWPRegistry.CallOpts)
}

// Guardian is a free data retrieval call binding the contract method 0x452a9320.
//
// Solidity: function guardian() view returns(address)
func (_AWPRegistry *AWPRegistryCallerSession) Guardian() (common.Address, error) {
	return _AWPRegistry.Contract.Guardian(&_AWPRegistry.CallOpts)
}

// ImmunityPeriod is a free data retrieval call binding the contract method 0x2672e1be.
//
// Solidity: function immunityPeriod() view returns(uint256)
func (_AWPRegistry *AWPRegistryCaller) ImmunityPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "immunityPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ImmunityPeriod is a free data retrieval call binding the contract method 0x2672e1be.
//
// Solidity: function immunityPeriod() view returns(uint256)
func (_AWPRegistry *AWPRegistrySession) ImmunityPeriod() (*big.Int, error) {
	return _AWPRegistry.Contract.ImmunityPeriod(&_AWPRegistry.CallOpts)
}

// ImmunityPeriod is a free data retrieval call binding the contract method 0x2672e1be.
//
// Solidity: function immunityPeriod() view returns(uint256)
func (_AWPRegistry *AWPRegistryCallerSession) ImmunityPeriod() (*big.Int, error) {
	return _AWPRegistry.Contract.ImmunityPeriod(&_AWPRegistry.CallOpts)
}

// InitialAlphaPrice is a free data retrieval call binding the contract method 0x6d345eea.
//
// Solidity: function initialAlphaPrice() view returns(uint256)
func (_AWPRegistry *AWPRegistryCaller) InitialAlphaPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "initialAlphaPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InitialAlphaPrice is a free data retrieval call binding the contract method 0x6d345eea.
//
// Solidity: function initialAlphaPrice() view returns(uint256)
func (_AWPRegistry *AWPRegistrySession) InitialAlphaPrice() (*big.Int, error) {
	return _AWPRegistry.Contract.InitialAlphaPrice(&_AWPRegistry.CallOpts)
}

// InitialAlphaPrice is a free data retrieval call binding the contract method 0x6d345eea.
//
// Solidity: function initialAlphaPrice() view returns(uint256)
func (_AWPRegistry *AWPRegistryCallerSession) InitialAlphaPrice() (*big.Int, error) {
	return _AWPRegistry.Contract.InitialAlphaPrice(&_AWPRegistry.CallOpts)
}

// IsRegistered is a free data retrieval call binding the contract method 0xc3c5a547.
//
// Solidity: function isRegistered(address addr) view returns(bool)
func (_AWPRegistry *AWPRegistryCaller) IsRegistered(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "isRegistered", addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsRegistered is a free data retrieval call binding the contract method 0xc3c5a547.
//
// Solidity: function isRegistered(address addr) view returns(bool)
func (_AWPRegistry *AWPRegistrySession) IsRegistered(addr common.Address) (bool, error) {
	return _AWPRegistry.Contract.IsRegistered(&_AWPRegistry.CallOpts, addr)
}

// IsRegistered is a free data retrieval call binding the contract method 0xc3c5a547.
//
// Solidity: function isRegistered(address addr) view returns(bool)
func (_AWPRegistry *AWPRegistryCallerSession) IsRegistered(addr common.Address) (bool, error) {
	return _AWPRegistry.Contract.IsRegistered(&_AWPRegistry.CallOpts, addr)
}

// IsSubnetActive is a free data retrieval call binding the contract method 0x7ab5e276.
//
// Solidity: function isSubnetActive(uint256 subnetId) view returns(bool)
func (_AWPRegistry *AWPRegistryCaller) IsSubnetActive(opts *bind.CallOpts, subnetId *big.Int) (bool, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "isSubnetActive", subnetId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsSubnetActive is a free data retrieval call binding the contract method 0x7ab5e276.
//
// Solidity: function isSubnetActive(uint256 subnetId) view returns(bool)
func (_AWPRegistry *AWPRegistrySession) IsSubnetActive(subnetId *big.Int) (bool, error) {
	return _AWPRegistry.Contract.IsSubnetActive(&_AWPRegistry.CallOpts, subnetId)
}

// IsSubnetActive is a free data retrieval call binding the contract method 0x7ab5e276.
//
// Solidity: function isSubnetActive(uint256 subnetId) view returns(bool)
func (_AWPRegistry *AWPRegistryCallerSession) IsSubnetActive(subnetId *big.Int) (bool, error) {
	return _AWPRegistry.Contract.IsSubnetActive(&_AWPRegistry.CallOpts, subnetId)
}

// LpManager is a free data retrieval call binding the contract method 0xb906f15a.
//
// Solidity: function lpManager() view returns(address)
func (_AWPRegistry *AWPRegistryCaller) LpManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "lpManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LpManager is a free data retrieval call binding the contract method 0xb906f15a.
//
// Solidity: function lpManager() view returns(address)
func (_AWPRegistry *AWPRegistrySession) LpManager() (common.Address, error) {
	return _AWPRegistry.Contract.LpManager(&_AWPRegistry.CallOpts)
}

// LpManager is a free data retrieval call binding the contract method 0xb906f15a.
//
// Solidity: function lpManager() view returns(address)
func (_AWPRegistry *AWPRegistryCallerSession) LpManager() (common.Address, error) {
	return _AWPRegistry.Contract.LpManager(&_AWPRegistry.CallOpts)
}

// NextSubnetId is a free data retrieval call binding the contract method 0xd929ff05.
//
// Solidity: function nextSubnetId() view returns(uint256)
func (_AWPRegistry *AWPRegistryCaller) NextSubnetId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "nextSubnetId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextSubnetId is a free data retrieval call binding the contract method 0xd929ff05.
//
// Solidity: function nextSubnetId() view returns(uint256)
func (_AWPRegistry *AWPRegistrySession) NextSubnetId() (*big.Int, error) {
	return _AWPRegistry.Contract.NextSubnetId(&_AWPRegistry.CallOpts)
}

// NextSubnetId is a free data retrieval call binding the contract method 0xd929ff05.
//
// Solidity: function nextSubnetId() view returns(uint256)
func (_AWPRegistry *AWPRegistryCallerSession) NextSubnetId() (*big.Int, error) {
	return _AWPRegistry.Contract.NextSubnetId(&_AWPRegistry.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_AWPRegistry *AWPRegistryCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "nonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_AWPRegistry *AWPRegistrySession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _AWPRegistry.Contract.Nonces(&_AWPRegistry.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_AWPRegistry *AWPRegistryCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _AWPRegistry.Contract.Nonces(&_AWPRegistry.CallOpts, arg0)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_AWPRegistry *AWPRegistryCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_AWPRegistry *AWPRegistrySession) Paused() (bool, error) {
	return _AWPRegistry.Contract.Paused(&_AWPRegistry.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_AWPRegistry *AWPRegistryCallerSession) Paused() (bool, error) {
	return _AWPRegistry.Contract.Paused(&_AWPRegistry.CallOpts)
}

// Recipient is a free data retrieval call binding the contract method 0xb3651eea.
//
// Solidity: function recipient(address ) view returns(address)
func (_AWPRegistry *AWPRegistryCaller) Recipient(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "recipient", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Recipient is a free data retrieval call binding the contract method 0xb3651eea.
//
// Solidity: function recipient(address ) view returns(address)
func (_AWPRegistry *AWPRegistrySession) Recipient(arg0 common.Address) (common.Address, error) {
	return _AWPRegistry.Contract.Recipient(&_AWPRegistry.CallOpts, arg0)
}

// Recipient is a free data retrieval call binding the contract method 0xb3651eea.
//
// Solidity: function recipient(address ) view returns(address)
func (_AWPRegistry *AWPRegistryCallerSession) Recipient(arg0 common.Address) (common.Address, error) {
	return _AWPRegistry.Contract.Recipient(&_AWPRegistry.CallOpts, arg0)
}

// RegisteredCount is a free data retrieval call binding the contract method 0x210ff9bb.
//
// Solidity: function registeredCount() view returns(uint256)
func (_AWPRegistry *AWPRegistryCaller) RegisteredCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "registeredCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RegisteredCount is a free data retrieval call binding the contract method 0x210ff9bb.
//
// Solidity: function registeredCount() view returns(uint256)
func (_AWPRegistry *AWPRegistrySession) RegisteredCount() (*big.Int, error) {
	return _AWPRegistry.Contract.RegisteredCount(&_AWPRegistry.CallOpts)
}

// RegisteredCount is a free data retrieval call binding the contract method 0x210ff9bb.
//
// Solidity: function registeredCount() view returns(uint256)
func (_AWPRegistry *AWPRegistryCallerSession) RegisteredCount() (*big.Int, error) {
	return _AWPRegistry.Contract.RegisteredCount(&_AWPRegistry.CallOpts)
}

// RegistryInitialized is a free data retrieval call binding the contract method 0x56354a24.
//
// Solidity: function registryInitialized() view returns(bool)
func (_AWPRegistry *AWPRegistryCaller) RegistryInitialized(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "registryInitialized")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// RegistryInitialized is a free data retrieval call binding the contract method 0x56354a24.
//
// Solidity: function registryInitialized() view returns(bool)
func (_AWPRegistry *AWPRegistrySession) RegistryInitialized() (bool, error) {
	return _AWPRegistry.Contract.RegistryInitialized(&_AWPRegistry.CallOpts)
}

// RegistryInitialized is a free data retrieval call binding the contract method 0x56354a24.
//
// Solidity: function registryInitialized() view returns(bool)
func (_AWPRegistry *AWPRegistryCallerSession) RegistryInitialized() (bool, error) {
	return _AWPRegistry.Contract.RegistryInitialized(&_AWPRegistry.CallOpts)
}

// ResolveRecipient is a free data retrieval call binding the contract method 0xfbea9d67.
//
// Solidity: function resolveRecipient(address addr) view returns(address)
func (_AWPRegistry *AWPRegistryCaller) ResolveRecipient(opts *bind.CallOpts, addr common.Address) (common.Address, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "resolveRecipient", addr)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ResolveRecipient is a free data retrieval call binding the contract method 0xfbea9d67.
//
// Solidity: function resolveRecipient(address addr) view returns(address)
func (_AWPRegistry *AWPRegistrySession) ResolveRecipient(addr common.Address) (common.Address, error) {
	return _AWPRegistry.Contract.ResolveRecipient(&_AWPRegistry.CallOpts, addr)
}

// ResolveRecipient is a free data retrieval call binding the contract method 0xfbea9d67.
//
// Solidity: function resolveRecipient(address addr) view returns(address)
func (_AWPRegistry *AWPRegistryCallerSession) ResolveRecipient(addr common.Address) (common.Address, error) {
	return _AWPRegistry.Contract.ResolveRecipient(&_AWPRegistry.CallOpts, addr)
}

// StakeNFT is a free data retrieval call binding the contract method 0xb48509e6.
//
// Solidity: function stakeNFT() view returns(address)
func (_AWPRegistry *AWPRegistryCaller) StakeNFT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "stakeNFT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakeNFT is a free data retrieval call binding the contract method 0xb48509e6.
//
// Solidity: function stakeNFT() view returns(address)
func (_AWPRegistry *AWPRegistrySession) StakeNFT() (common.Address, error) {
	return _AWPRegistry.Contract.StakeNFT(&_AWPRegistry.CallOpts)
}

// StakeNFT is a free data retrieval call binding the contract method 0xb48509e6.
//
// Solidity: function stakeNFT() view returns(address)
func (_AWPRegistry *AWPRegistryCallerSession) StakeNFT() (common.Address, error) {
	return _AWPRegistry.Contract.StakeNFT(&_AWPRegistry.CallOpts)
}

// StakingVault is a free data retrieval call binding the contract method 0x24e7964a.
//
// Solidity: function stakingVault() view returns(address)
func (_AWPRegistry *AWPRegistryCaller) StakingVault(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "stakingVault")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakingVault is a free data retrieval call binding the contract method 0x24e7964a.
//
// Solidity: function stakingVault() view returns(address)
func (_AWPRegistry *AWPRegistrySession) StakingVault() (common.Address, error) {
	return _AWPRegistry.Contract.StakingVault(&_AWPRegistry.CallOpts)
}

// StakingVault is a free data retrieval call binding the contract method 0x24e7964a.
//
// Solidity: function stakingVault() view returns(address)
func (_AWPRegistry *AWPRegistryCallerSession) StakingVault() (common.Address, error) {
	return _AWPRegistry.Contract.StakingVault(&_AWPRegistry.CallOpts)
}

// SubnetNFT is a free data retrieval call binding the contract method 0x11cba7e9.
//
// Solidity: function subnetNFT() view returns(address)
func (_AWPRegistry *AWPRegistryCaller) SubnetNFT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "subnetNFT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SubnetNFT is a free data retrieval call binding the contract method 0x11cba7e9.
//
// Solidity: function subnetNFT() view returns(address)
func (_AWPRegistry *AWPRegistrySession) SubnetNFT() (common.Address, error) {
	return _AWPRegistry.Contract.SubnetNFT(&_AWPRegistry.CallOpts)
}

// SubnetNFT is a free data retrieval call binding the contract method 0x11cba7e9.
//
// Solidity: function subnetNFT() view returns(address)
func (_AWPRegistry *AWPRegistryCallerSession) SubnetNFT() (common.Address, error) {
	return _AWPRegistry.Contract.SubnetNFT(&_AWPRegistry.CallOpts)
}

// Subnets is a free data retrieval call binding the contract method 0x475726f7.
//
// Solidity: function subnets(uint256 ) view returns(bytes32 lpPool, uint8 status, uint64 createdAt, uint64 activatedAt)
func (_AWPRegistry *AWPRegistryCaller) Subnets(opts *bind.CallOpts, arg0 *big.Int) (struct {
	LpPool      [32]byte
	Status      uint8
	CreatedAt   uint64
	ActivatedAt uint64
}, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "subnets", arg0)

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
func (_AWPRegistry *AWPRegistrySession) Subnets(arg0 *big.Int) (struct {
	LpPool      [32]byte
	Status      uint8
	CreatedAt   uint64
	ActivatedAt uint64
}, error) {
	return _AWPRegistry.Contract.Subnets(&_AWPRegistry.CallOpts, arg0)
}

// Subnets is a free data retrieval call binding the contract method 0x475726f7.
//
// Solidity: function subnets(uint256 ) view returns(bytes32 lpPool, uint8 status, uint64 createdAt, uint64 activatedAt)
func (_AWPRegistry *AWPRegistryCallerSession) Subnets(arg0 *big.Int) (struct {
	LpPool      [32]byte
	Status      uint8
	CreatedAt   uint64
	ActivatedAt uint64
}, error) {
	return _AWPRegistry.Contract.Subnets(&_AWPRegistry.CallOpts, arg0)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_AWPRegistry *AWPRegistryCaller) Treasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPRegistry.contract.Call(opts, &out, "treasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_AWPRegistry *AWPRegistrySession) Treasury() (common.Address, error) {
	return _AWPRegistry.Contract.Treasury(&_AWPRegistry.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_AWPRegistry *AWPRegistryCallerSession) Treasury() (common.Address, error) {
	return _AWPRegistry.Contract.Treasury(&_AWPRegistry.CallOpts)
}

// ActivateSubnet is a paid mutator transaction binding the contract method 0xcead1c96.
//
// Solidity: function activateSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistryTransactor) ActivateSubnet(opts *bind.TransactOpts, subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "activateSubnet", subnetId)
}

// ActivateSubnet is a paid mutator transaction binding the contract method 0xcead1c96.
//
// Solidity: function activateSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistrySession) ActivateSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.ActivateSubnet(&_AWPRegistry.TransactOpts, subnetId)
}

// ActivateSubnet is a paid mutator transaction binding the contract method 0xcead1c96.
//
// Solidity: function activateSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) ActivateSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.ActivateSubnet(&_AWPRegistry.TransactOpts, subnetId)
}

// ActivateSubnetFor is a paid mutator transaction binding the contract method 0x08b55cff.
//
// Solidity: function activateSubnetFor(address user, uint256 subnetId, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPRegistry *AWPRegistryTransactor) ActivateSubnetFor(opts *bind.TransactOpts, user common.Address, subnetId *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "activateSubnetFor", user, subnetId, deadline, v, r, s)
}

// ActivateSubnetFor is a paid mutator transaction binding the contract method 0x08b55cff.
//
// Solidity: function activateSubnetFor(address user, uint256 subnetId, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPRegistry *AWPRegistrySession) ActivateSubnetFor(user common.Address, subnetId *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.ActivateSubnetFor(&_AWPRegistry.TransactOpts, user, subnetId, deadline, v, r, s)
}

// ActivateSubnetFor is a paid mutator transaction binding the contract method 0x08b55cff.
//
// Solidity: function activateSubnetFor(address user, uint256 subnetId, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) ActivateSubnetFor(user common.Address, subnetId *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.ActivateSubnetFor(&_AWPRegistry.TransactOpts, user, subnetId, deadline, v, r, s)
}

// Allocate is a paid mutator transaction binding the contract method 0xd035a9a7.
//
// Solidity: function allocate(address staker, address agent, uint256 subnetId, uint256 amount) returns()
func (_AWPRegistry *AWPRegistryTransactor) Allocate(opts *bind.TransactOpts, staker common.Address, agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "allocate", staker, agent, subnetId, amount)
}

// Allocate is a paid mutator transaction binding the contract method 0xd035a9a7.
//
// Solidity: function allocate(address staker, address agent, uint256 subnetId, uint256 amount) returns()
func (_AWPRegistry *AWPRegistrySession) Allocate(staker common.Address, agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.Allocate(&_AWPRegistry.TransactOpts, staker, agent, subnetId, amount)
}

// Allocate is a paid mutator transaction binding the contract method 0xd035a9a7.
//
// Solidity: function allocate(address staker, address agent, uint256 subnetId, uint256 amount) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) Allocate(staker common.Address, agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.Allocate(&_AWPRegistry.TransactOpts, staker, agent, subnetId, amount)
}

// AllocateFor is a paid mutator transaction binding the contract method 0x7d66c5c5.
//
// Solidity: function allocateFor(address staker, address agent, uint256 subnetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPRegistry *AWPRegistryTransactor) AllocateFor(opts *bind.TransactOpts, staker common.Address, agent common.Address, subnetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "allocateFor", staker, agent, subnetId, amount, deadline, v, r, s)
}

// AllocateFor is a paid mutator transaction binding the contract method 0x7d66c5c5.
//
// Solidity: function allocateFor(address staker, address agent, uint256 subnetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPRegistry *AWPRegistrySession) AllocateFor(staker common.Address, agent common.Address, subnetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.AllocateFor(&_AWPRegistry.TransactOpts, staker, agent, subnetId, amount, deadline, v, r, s)
}

// AllocateFor is a paid mutator transaction binding the contract method 0x7d66c5c5.
//
// Solidity: function allocateFor(address staker, address agent, uint256 subnetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) AllocateFor(staker common.Address, agent common.Address, subnetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.AllocateFor(&_AWPRegistry.TransactOpts, staker, agent, subnetId, amount, deadline, v, r, s)
}

// BanSubnet is a paid mutator transaction binding the contract method 0xb79b7658.
//
// Solidity: function banSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistryTransactor) BanSubnet(opts *bind.TransactOpts, subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "banSubnet", subnetId)
}

// BanSubnet is a paid mutator transaction binding the contract method 0xb79b7658.
//
// Solidity: function banSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistrySession) BanSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.BanSubnet(&_AWPRegistry.TransactOpts, subnetId)
}

// BanSubnet is a paid mutator transaction binding the contract method 0xb79b7658.
//
// Solidity: function banSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) BanSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.BanSubnet(&_AWPRegistry.TransactOpts, subnetId)
}

// Bind is a paid mutator transaction binding the contract method 0x81bac14f.
//
// Solidity: function bind(address target) returns()
func (_AWPRegistry *AWPRegistryTransactor) Bind(opts *bind.TransactOpts, target common.Address) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "bind", target)
}

// Bind is a paid mutator transaction binding the contract method 0x81bac14f.
//
// Solidity: function bind(address target) returns()
func (_AWPRegistry *AWPRegistrySession) Bind(target common.Address) (*types.Transaction, error) {
	return _AWPRegistry.Contract.Bind(&_AWPRegistry.TransactOpts, target)
}

// Bind is a paid mutator transaction binding the contract method 0x81bac14f.
//
// Solidity: function bind(address target) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) Bind(target common.Address) (*types.Transaction, error) {
	return _AWPRegistry.Contract.Bind(&_AWPRegistry.TransactOpts, target)
}

// BindFor is a paid mutator transaction binding the contract method 0x7b234b81.
//
// Solidity: function bindFor(address agent, address target, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPRegistry *AWPRegistryTransactor) BindFor(opts *bind.TransactOpts, agent common.Address, target common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "bindFor", agent, target, deadline, v, r, s)
}

// BindFor is a paid mutator transaction binding the contract method 0x7b234b81.
//
// Solidity: function bindFor(address agent, address target, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPRegistry *AWPRegistrySession) BindFor(agent common.Address, target common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.BindFor(&_AWPRegistry.TransactOpts, agent, target, deadline, v, r, s)
}

// BindFor is a paid mutator transaction binding the contract method 0x7b234b81.
//
// Solidity: function bindFor(address agent, address target, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) BindFor(agent common.Address, target common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.BindFor(&_AWPRegistry.TransactOpts, agent, target, deadline, v, r, s)
}

// Deallocate is a paid mutator transaction binding the contract method 0x716fb83d.
//
// Solidity: function deallocate(address staker, address agent, uint256 subnetId, uint256 amount) returns()
func (_AWPRegistry *AWPRegistryTransactor) Deallocate(opts *bind.TransactOpts, staker common.Address, agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "deallocate", staker, agent, subnetId, amount)
}

// Deallocate is a paid mutator transaction binding the contract method 0x716fb83d.
//
// Solidity: function deallocate(address staker, address agent, uint256 subnetId, uint256 amount) returns()
func (_AWPRegistry *AWPRegistrySession) Deallocate(staker common.Address, agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.Deallocate(&_AWPRegistry.TransactOpts, staker, agent, subnetId, amount)
}

// Deallocate is a paid mutator transaction binding the contract method 0x716fb83d.
//
// Solidity: function deallocate(address staker, address agent, uint256 subnetId, uint256 amount) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) Deallocate(staker common.Address, agent common.Address, subnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.Deallocate(&_AWPRegistry.TransactOpts, staker, agent, subnetId, amount)
}

// DeallocateFor is a paid mutator transaction binding the contract method 0x10fe1208.
//
// Solidity: function deallocateFor(address staker, address agent, uint256 subnetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPRegistry *AWPRegistryTransactor) DeallocateFor(opts *bind.TransactOpts, staker common.Address, agent common.Address, subnetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "deallocateFor", staker, agent, subnetId, amount, deadline, v, r, s)
}

// DeallocateFor is a paid mutator transaction binding the contract method 0x10fe1208.
//
// Solidity: function deallocateFor(address staker, address agent, uint256 subnetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPRegistry *AWPRegistrySession) DeallocateFor(staker common.Address, agent common.Address, subnetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.DeallocateFor(&_AWPRegistry.TransactOpts, staker, agent, subnetId, amount, deadline, v, r, s)
}

// DeallocateFor is a paid mutator transaction binding the contract method 0x10fe1208.
//
// Solidity: function deallocateFor(address staker, address agent, uint256 subnetId, uint256 amount, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) DeallocateFor(staker common.Address, agent common.Address, subnetId *big.Int, amount *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.DeallocateFor(&_AWPRegistry.TransactOpts, staker, agent, subnetId, amount, deadline, v, r, s)
}

// DeregisterSubnet is a paid mutator transaction binding the contract method 0x0cf02c5e.
//
// Solidity: function deregisterSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistryTransactor) DeregisterSubnet(opts *bind.TransactOpts, subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "deregisterSubnet", subnetId)
}

// DeregisterSubnet is a paid mutator transaction binding the contract method 0x0cf02c5e.
//
// Solidity: function deregisterSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistrySession) DeregisterSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.DeregisterSubnet(&_AWPRegistry.TransactOpts, subnetId)
}

// DeregisterSubnet is a paid mutator transaction binding the contract method 0x0cf02c5e.
//
// Solidity: function deregisterSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) DeregisterSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.DeregisterSubnet(&_AWPRegistry.TransactOpts, subnetId)
}

// GrantDelegate is a paid mutator transaction binding the contract method 0xa757acd9.
//
// Solidity: function grantDelegate(address delegate) returns()
func (_AWPRegistry *AWPRegistryTransactor) GrantDelegate(opts *bind.TransactOpts, delegate common.Address) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "grantDelegate", delegate)
}

// GrantDelegate is a paid mutator transaction binding the contract method 0xa757acd9.
//
// Solidity: function grantDelegate(address delegate) returns()
func (_AWPRegistry *AWPRegistrySession) GrantDelegate(delegate common.Address) (*types.Transaction, error) {
	return _AWPRegistry.Contract.GrantDelegate(&_AWPRegistry.TransactOpts, delegate)
}

// GrantDelegate is a paid mutator transaction binding the contract method 0xa757acd9.
//
// Solidity: function grantDelegate(address delegate) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) GrantDelegate(delegate common.Address) (*types.Transaction, error) {
	return _AWPRegistry.Contract.GrantDelegate(&_AWPRegistry.TransactOpts, delegate)
}

// InitializeRegistry is a paid mutator transaction binding the contract method 0x534489a0.
//
// Solidity: function initializeRegistry(address awpToken_, address subnetNFT_, address alphaTokenFactory_, address awpEmission_, address lpManager_, address stakingVault_, address stakeNFT_, address defaultSubnetManagerImpl_, bytes dexConfig_) returns()
func (_AWPRegistry *AWPRegistryTransactor) InitializeRegistry(opts *bind.TransactOpts, awpToken_ common.Address, subnetNFT_ common.Address, alphaTokenFactory_ common.Address, awpEmission_ common.Address, lpManager_ common.Address, stakingVault_ common.Address, stakeNFT_ common.Address, defaultSubnetManagerImpl_ common.Address, dexConfig_ []byte) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "initializeRegistry", awpToken_, subnetNFT_, alphaTokenFactory_, awpEmission_, lpManager_, stakingVault_, stakeNFT_, defaultSubnetManagerImpl_, dexConfig_)
}

// InitializeRegistry is a paid mutator transaction binding the contract method 0x534489a0.
//
// Solidity: function initializeRegistry(address awpToken_, address subnetNFT_, address alphaTokenFactory_, address awpEmission_, address lpManager_, address stakingVault_, address stakeNFT_, address defaultSubnetManagerImpl_, bytes dexConfig_) returns()
func (_AWPRegistry *AWPRegistrySession) InitializeRegistry(awpToken_ common.Address, subnetNFT_ common.Address, alphaTokenFactory_ common.Address, awpEmission_ common.Address, lpManager_ common.Address, stakingVault_ common.Address, stakeNFT_ common.Address, defaultSubnetManagerImpl_ common.Address, dexConfig_ []byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.InitializeRegistry(&_AWPRegistry.TransactOpts, awpToken_, subnetNFT_, alphaTokenFactory_, awpEmission_, lpManager_, stakingVault_, stakeNFT_, defaultSubnetManagerImpl_, dexConfig_)
}

// InitializeRegistry is a paid mutator transaction binding the contract method 0x534489a0.
//
// Solidity: function initializeRegistry(address awpToken_, address subnetNFT_, address alphaTokenFactory_, address awpEmission_, address lpManager_, address stakingVault_, address stakeNFT_, address defaultSubnetManagerImpl_, bytes dexConfig_) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) InitializeRegistry(awpToken_ common.Address, subnetNFT_ common.Address, alphaTokenFactory_ common.Address, awpEmission_ common.Address, lpManager_ common.Address, stakingVault_ common.Address, stakeNFT_ common.Address, defaultSubnetManagerImpl_ common.Address, dexConfig_ []byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.InitializeRegistry(&_AWPRegistry.TransactOpts, awpToken_, subnetNFT_, alphaTokenFactory_, awpEmission_, lpManager_, stakingVault_, stakeNFT_, defaultSubnetManagerImpl_, dexConfig_)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_AWPRegistry *AWPRegistryTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_AWPRegistry *AWPRegistrySession) Pause() (*types.Transaction, error) {
	return _AWPRegistry.Contract.Pause(&_AWPRegistry.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_AWPRegistry *AWPRegistryTransactorSession) Pause() (*types.Transaction, error) {
	return _AWPRegistry.Contract.Pause(&_AWPRegistry.TransactOpts)
}

// PauseSubnet is a paid mutator transaction binding the contract method 0x44e047ca.
//
// Solidity: function pauseSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistryTransactor) PauseSubnet(opts *bind.TransactOpts, subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "pauseSubnet", subnetId)
}

// PauseSubnet is a paid mutator transaction binding the contract method 0x44e047ca.
//
// Solidity: function pauseSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistrySession) PauseSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.PauseSubnet(&_AWPRegistry.TransactOpts, subnetId)
}

// PauseSubnet is a paid mutator transaction binding the contract method 0x44e047ca.
//
// Solidity: function pauseSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) PauseSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.PauseSubnet(&_AWPRegistry.TransactOpts, subnetId)
}

// Reallocate is a paid mutator transaction binding the contract method 0xd5d5278d.
//
// Solidity: function reallocate(address staker, address fromAgent, uint256 fromSubnetId, address toAgent, uint256 toSubnetId, uint256 amount) returns()
func (_AWPRegistry *AWPRegistryTransactor) Reallocate(opts *bind.TransactOpts, staker common.Address, fromAgent common.Address, fromSubnetId *big.Int, toAgent common.Address, toSubnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "reallocate", staker, fromAgent, fromSubnetId, toAgent, toSubnetId, amount)
}

// Reallocate is a paid mutator transaction binding the contract method 0xd5d5278d.
//
// Solidity: function reallocate(address staker, address fromAgent, uint256 fromSubnetId, address toAgent, uint256 toSubnetId, uint256 amount) returns()
func (_AWPRegistry *AWPRegistrySession) Reallocate(staker common.Address, fromAgent common.Address, fromSubnetId *big.Int, toAgent common.Address, toSubnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.Reallocate(&_AWPRegistry.TransactOpts, staker, fromAgent, fromSubnetId, toAgent, toSubnetId, amount)
}

// Reallocate is a paid mutator transaction binding the contract method 0xd5d5278d.
//
// Solidity: function reallocate(address staker, address fromAgent, uint256 fromSubnetId, address toAgent, uint256 toSubnetId, uint256 amount) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) Reallocate(staker common.Address, fromAgent common.Address, fromSubnetId *big.Int, toAgent common.Address, toSubnetId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.Reallocate(&_AWPRegistry.TransactOpts, staker, fromAgent, fromSubnetId, toAgent, toSubnetId, amount)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_AWPRegistry *AWPRegistryTransactor) Register(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "register")
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_AWPRegistry *AWPRegistrySession) Register() (*types.Transaction, error) {
	return _AWPRegistry.Contract.Register(&_AWPRegistry.TransactOpts)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_AWPRegistry *AWPRegistryTransactorSession) Register() (*types.Transaction, error) {
	return _AWPRegistry.Contract.Register(&_AWPRegistry.TransactOpts)
}

// RegisterSubnet is a paid mutator transaction binding the contract method 0x5f24898d.
//
// Solidity: function registerSubnet((string,string,address,bytes32,uint128,string) params) returns(uint256)
func (_AWPRegistry *AWPRegistryTransactor) RegisterSubnet(opts *bind.TransactOpts, params IAWPRegistrySubnetParams) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "registerSubnet", params)
}

// RegisterSubnet is a paid mutator transaction binding the contract method 0x5f24898d.
//
// Solidity: function registerSubnet((string,string,address,bytes32,uint128,string) params) returns(uint256)
func (_AWPRegistry *AWPRegistrySession) RegisterSubnet(params IAWPRegistrySubnetParams) (*types.Transaction, error) {
	return _AWPRegistry.Contract.RegisterSubnet(&_AWPRegistry.TransactOpts, params)
}

// RegisterSubnet is a paid mutator transaction binding the contract method 0x5f24898d.
//
// Solidity: function registerSubnet((string,string,address,bytes32,uint128,string) params) returns(uint256)
func (_AWPRegistry *AWPRegistryTransactorSession) RegisterSubnet(params IAWPRegistrySubnetParams) (*types.Transaction, error) {
	return _AWPRegistry.Contract.RegisterSubnet(&_AWPRegistry.TransactOpts, params)
}

// RegisterSubnetFor is a paid mutator transaction binding the contract method 0x1aa3ff5a.
//
// Solidity: function registerSubnetFor(address user, (string,string,address,bytes32,uint128,string) params, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_AWPRegistry *AWPRegistryTransactor) RegisterSubnetFor(opts *bind.TransactOpts, user common.Address, params IAWPRegistrySubnetParams, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "registerSubnetFor", user, params, deadline, v, r, s)
}

// RegisterSubnetFor is a paid mutator transaction binding the contract method 0x1aa3ff5a.
//
// Solidity: function registerSubnetFor(address user, (string,string,address,bytes32,uint128,string) params, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_AWPRegistry *AWPRegistrySession) RegisterSubnetFor(user common.Address, params IAWPRegistrySubnetParams, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.RegisterSubnetFor(&_AWPRegistry.TransactOpts, user, params, deadline, v, r, s)
}

// RegisterSubnetFor is a paid mutator transaction binding the contract method 0x1aa3ff5a.
//
// Solidity: function registerSubnetFor(address user, (string,string,address,bytes32,uint128,string) params, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_AWPRegistry *AWPRegistryTransactorSession) RegisterSubnetFor(user common.Address, params IAWPRegistrySubnetParams, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.RegisterSubnetFor(&_AWPRegistry.TransactOpts, user, params, deadline, v, r, s)
}

// RegisterSubnetForWithPermit is a paid mutator transaction binding the contract method 0xedf12231.
//
// Solidity: function registerSubnetForWithPermit(address user, (string,string,address,bytes32,uint128,string) params, uint256 deadline, uint8 permitV, bytes32 permitR, bytes32 permitS, uint8 registerV, bytes32 registerR, bytes32 registerS) returns(uint256)
func (_AWPRegistry *AWPRegistryTransactor) RegisterSubnetForWithPermit(opts *bind.TransactOpts, user common.Address, params IAWPRegistrySubnetParams, deadline *big.Int, permitV uint8, permitR [32]byte, permitS [32]byte, registerV uint8, registerR [32]byte, registerS [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "registerSubnetForWithPermit", user, params, deadline, permitV, permitR, permitS, registerV, registerR, registerS)
}

// RegisterSubnetForWithPermit is a paid mutator transaction binding the contract method 0xedf12231.
//
// Solidity: function registerSubnetForWithPermit(address user, (string,string,address,bytes32,uint128,string) params, uint256 deadline, uint8 permitV, bytes32 permitR, bytes32 permitS, uint8 registerV, bytes32 registerR, bytes32 registerS) returns(uint256)
func (_AWPRegistry *AWPRegistrySession) RegisterSubnetForWithPermit(user common.Address, params IAWPRegistrySubnetParams, deadline *big.Int, permitV uint8, permitR [32]byte, permitS [32]byte, registerV uint8, registerR [32]byte, registerS [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.RegisterSubnetForWithPermit(&_AWPRegistry.TransactOpts, user, params, deadline, permitV, permitR, permitS, registerV, registerR, registerS)
}

// RegisterSubnetForWithPermit is a paid mutator transaction binding the contract method 0xedf12231.
//
// Solidity: function registerSubnetForWithPermit(address user, (string,string,address,bytes32,uint128,string) params, uint256 deadline, uint8 permitV, bytes32 permitR, bytes32 permitS, uint8 registerV, bytes32 registerR, bytes32 registerS) returns(uint256)
func (_AWPRegistry *AWPRegistryTransactorSession) RegisterSubnetForWithPermit(user common.Address, params IAWPRegistrySubnetParams, deadline *big.Int, permitV uint8, permitR [32]byte, permitS [32]byte, registerV uint8, registerR [32]byte, registerS [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.RegisterSubnetForWithPermit(&_AWPRegistry.TransactOpts, user, params, deadline, permitV, permitR, permitS, registerV, registerR, registerS)
}

// ResumeSubnet is a paid mutator transaction binding the contract method 0x5364944c.
//
// Solidity: function resumeSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistryTransactor) ResumeSubnet(opts *bind.TransactOpts, subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "resumeSubnet", subnetId)
}

// ResumeSubnet is a paid mutator transaction binding the contract method 0x5364944c.
//
// Solidity: function resumeSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistrySession) ResumeSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.ResumeSubnet(&_AWPRegistry.TransactOpts, subnetId)
}

// ResumeSubnet is a paid mutator transaction binding the contract method 0x5364944c.
//
// Solidity: function resumeSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) ResumeSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.ResumeSubnet(&_AWPRegistry.TransactOpts, subnetId)
}

// RevokeDelegate is a paid mutator transaction binding the contract method 0xfa352c00.
//
// Solidity: function revokeDelegate(address delegate) returns()
func (_AWPRegistry *AWPRegistryTransactor) RevokeDelegate(opts *bind.TransactOpts, delegate common.Address) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "revokeDelegate", delegate)
}

// RevokeDelegate is a paid mutator transaction binding the contract method 0xfa352c00.
//
// Solidity: function revokeDelegate(address delegate) returns()
func (_AWPRegistry *AWPRegistrySession) RevokeDelegate(delegate common.Address) (*types.Transaction, error) {
	return _AWPRegistry.Contract.RevokeDelegate(&_AWPRegistry.TransactOpts, delegate)
}

// RevokeDelegate is a paid mutator transaction binding the contract method 0xfa352c00.
//
// Solidity: function revokeDelegate(address delegate) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) RevokeDelegate(delegate common.Address) (*types.Transaction, error) {
	return _AWPRegistry.Contract.RevokeDelegate(&_AWPRegistry.TransactOpts, delegate)
}

// SetAlphaTokenFactory is a paid mutator transaction binding the contract method 0x901a71e4.
//
// Solidity: function setAlphaTokenFactory(address factory) returns()
func (_AWPRegistry *AWPRegistryTransactor) SetAlphaTokenFactory(opts *bind.TransactOpts, factory common.Address) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "setAlphaTokenFactory", factory)
}

// SetAlphaTokenFactory is a paid mutator transaction binding the contract method 0x901a71e4.
//
// Solidity: function setAlphaTokenFactory(address factory) returns()
func (_AWPRegistry *AWPRegistrySession) SetAlphaTokenFactory(factory common.Address) (*types.Transaction, error) {
	return _AWPRegistry.Contract.SetAlphaTokenFactory(&_AWPRegistry.TransactOpts, factory)
}

// SetAlphaTokenFactory is a paid mutator transaction binding the contract method 0x901a71e4.
//
// Solidity: function setAlphaTokenFactory(address factory) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) SetAlphaTokenFactory(factory common.Address) (*types.Transaction, error) {
	return _AWPRegistry.Contract.SetAlphaTokenFactory(&_AWPRegistry.TransactOpts, factory)
}

// SetDexConfig is a paid mutator transaction binding the contract method 0x042fce70.
//
// Solidity: function setDexConfig(bytes dexConfig_) returns()
func (_AWPRegistry *AWPRegistryTransactor) SetDexConfig(opts *bind.TransactOpts, dexConfig_ []byte) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "setDexConfig", dexConfig_)
}

// SetDexConfig is a paid mutator transaction binding the contract method 0x042fce70.
//
// Solidity: function setDexConfig(bytes dexConfig_) returns()
func (_AWPRegistry *AWPRegistrySession) SetDexConfig(dexConfig_ []byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.SetDexConfig(&_AWPRegistry.TransactOpts, dexConfig_)
}

// SetDexConfig is a paid mutator transaction binding the contract method 0x042fce70.
//
// Solidity: function setDexConfig(bytes dexConfig_) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) SetDexConfig(dexConfig_ []byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.SetDexConfig(&_AWPRegistry.TransactOpts, dexConfig_)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_AWPRegistry *AWPRegistryTransactor) SetGuardian(opts *bind.TransactOpts, g common.Address) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "setGuardian", g)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_AWPRegistry *AWPRegistrySession) SetGuardian(g common.Address) (*types.Transaction, error) {
	return _AWPRegistry.Contract.SetGuardian(&_AWPRegistry.TransactOpts, g)
}

// SetGuardian is a paid mutator transaction binding the contract method 0x8a0dac4a.
//
// Solidity: function setGuardian(address g) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) SetGuardian(g common.Address) (*types.Transaction, error) {
	return _AWPRegistry.Contract.SetGuardian(&_AWPRegistry.TransactOpts, g)
}

// SetImmunityPeriod is a paid mutator transaction binding the contract method 0x33bbf030.
//
// Solidity: function setImmunityPeriod(uint256 p) returns()
func (_AWPRegistry *AWPRegistryTransactor) SetImmunityPeriod(opts *bind.TransactOpts, p *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "setImmunityPeriod", p)
}

// SetImmunityPeriod is a paid mutator transaction binding the contract method 0x33bbf030.
//
// Solidity: function setImmunityPeriod(uint256 p) returns()
func (_AWPRegistry *AWPRegistrySession) SetImmunityPeriod(p *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.SetImmunityPeriod(&_AWPRegistry.TransactOpts, p)
}

// SetImmunityPeriod is a paid mutator transaction binding the contract method 0x33bbf030.
//
// Solidity: function setImmunityPeriod(uint256 p) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) SetImmunityPeriod(p *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.SetImmunityPeriod(&_AWPRegistry.TransactOpts, p)
}

// SetInitialAlphaPrice is a paid mutator transaction binding the contract method 0xe7d89b71.
//
// Solidity: function setInitialAlphaPrice(uint256 price) returns()
func (_AWPRegistry *AWPRegistryTransactor) SetInitialAlphaPrice(opts *bind.TransactOpts, price *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "setInitialAlphaPrice", price)
}

// SetInitialAlphaPrice is a paid mutator transaction binding the contract method 0xe7d89b71.
//
// Solidity: function setInitialAlphaPrice(uint256 price) returns()
func (_AWPRegistry *AWPRegistrySession) SetInitialAlphaPrice(price *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.SetInitialAlphaPrice(&_AWPRegistry.TransactOpts, price)
}

// SetInitialAlphaPrice is a paid mutator transaction binding the contract method 0xe7d89b71.
//
// Solidity: function setInitialAlphaPrice(uint256 price) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) SetInitialAlphaPrice(price *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.SetInitialAlphaPrice(&_AWPRegistry.TransactOpts, price)
}

// SetRecipient is a paid mutator transaction binding the contract method 0x3bbed4a0.
//
// Solidity: function setRecipient(address addr) returns()
func (_AWPRegistry *AWPRegistryTransactor) SetRecipient(opts *bind.TransactOpts, addr common.Address) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "setRecipient", addr)
}

// SetRecipient is a paid mutator transaction binding the contract method 0x3bbed4a0.
//
// Solidity: function setRecipient(address addr) returns()
func (_AWPRegistry *AWPRegistrySession) SetRecipient(addr common.Address) (*types.Transaction, error) {
	return _AWPRegistry.Contract.SetRecipient(&_AWPRegistry.TransactOpts, addr)
}

// SetRecipient is a paid mutator transaction binding the contract method 0x3bbed4a0.
//
// Solidity: function setRecipient(address addr) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) SetRecipient(addr common.Address) (*types.Transaction, error) {
	return _AWPRegistry.Contract.SetRecipient(&_AWPRegistry.TransactOpts, addr)
}

// SetRecipientFor is a paid mutator transaction binding the contract method 0x0026a047.
//
// Solidity: function setRecipientFor(address user, address _recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPRegistry *AWPRegistryTransactor) SetRecipientFor(opts *bind.TransactOpts, user common.Address, _recipient common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "setRecipientFor", user, _recipient, deadline, v, r, s)
}

// SetRecipientFor is a paid mutator transaction binding the contract method 0x0026a047.
//
// Solidity: function setRecipientFor(address user, address _recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPRegistry *AWPRegistrySession) SetRecipientFor(user common.Address, _recipient common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.SetRecipientFor(&_AWPRegistry.TransactOpts, user, _recipient, deadline, v, r, s)
}

// SetRecipientFor is a paid mutator transaction binding the contract method 0x0026a047.
//
// Solidity: function setRecipientFor(address user, address _recipient, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) SetRecipientFor(user common.Address, _recipient common.Address, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _AWPRegistry.Contract.SetRecipientFor(&_AWPRegistry.TransactOpts, user, _recipient, deadline, v, r, s)
}

// SetSubnetManagerImpl is a paid mutator transaction binding the contract method 0xe7c17212.
//
// Solidity: function setSubnetManagerImpl(address impl) returns()
func (_AWPRegistry *AWPRegistryTransactor) SetSubnetManagerImpl(opts *bind.TransactOpts, impl common.Address) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "setSubnetManagerImpl", impl)
}

// SetSubnetManagerImpl is a paid mutator transaction binding the contract method 0xe7c17212.
//
// Solidity: function setSubnetManagerImpl(address impl) returns()
func (_AWPRegistry *AWPRegistrySession) SetSubnetManagerImpl(impl common.Address) (*types.Transaction, error) {
	return _AWPRegistry.Contract.SetSubnetManagerImpl(&_AWPRegistry.TransactOpts, impl)
}

// SetSubnetManagerImpl is a paid mutator transaction binding the contract method 0xe7c17212.
//
// Solidity: function setSubnetManagerImpl(address impl) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) SetSubnetManagerImpl(impl common.Address) (*types.Transaction, error) {
	return _AWPRegistry.Contract.SetSubnetManagerImpl(&_AWPRegistry.TransactOpts, impl)
}

// UnbanSubnet is a paid mutator transaction binding the contract method 0x2bf1c05d.
//
// Solidity: function unbanSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistryTransactor) UnbanSubnet(opts *bind.TransactOpts, subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "unbanSubnet", subnetId)
}

// UnbanSubnet is a paid mutator transaction binding the contract method 0x2bf1c05d.
//
// Solidity: function unbanSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistrySession) UnbanSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.UnbanSubnet(&_AWPRegistry.TransactOpts, subnetId)
}

// UnbanSubnet is a paid mutator transaction binding the contract method 0x2bf1c05d.
//
// Solidity: function unbanSubnet(uint256 subnetId) returns()
func (_AWPRegistry *AWPRegistryTransactorSession) UnbanSubnet(subnetId *big.Int) (*types.Transaction, error) {
	return _AWPRegistry.Contract.UnbanSubnet(&_AWPRegistry.TransactOpts, subnetId)
}

// Unbind is a paid mutator transaction binding the contract method 0xb6b25742.
//
// Solidity: function unbind() returns()
func (_AWPRegistry *AWPRegistryTransactor) Unbind(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "unbind")
}

// Unbind is a paid mutator transaction binding the contract method 0xb6b25742.
//
// Solidity: function unbind() returns()
func (_AWPRegistry *AWPRegistrySession) Unbind() (*types.Transaction, error) {
	return _AWPRegistry.Contract.Unbind(&_AWPRegistry.TransactOpts)
}

// Unbind is a paid mutator transaction binding the contract method 0xb6b25742.
//
// Solidity: function unbind() returns()
func (_AWPRegistry *AWPRegistryTransactorSession) Unbind() (*types.Transaction, error) {
	return _AWPRegistry.Contract.Unbind(&_AWPRegistry.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_AWPRegistry *AWPRegistryTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AWPRegistry.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_AWPRegistry *AWPRegistrySession) Unpause() (*types.Transaction, error) {
	return _AWPRegistry.Contract.Unpause(&_AWPRegistry.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_AWPRegistry *AWPRegistryTransactorSession) Unpause() (*types.Transaction, error) {
	return _AWPRegistry.Contract.Unpause(&_AWPRegistry.TransactOpts)
}

// AWPRegistryAllocatedIterator is returned from FilterAllocated and is used to iterate over the raw logs and unpacked data for Allocated events raised by the AWPRegistry contract.
type AWPRegistryAllocatedIterator struct {
	Event *AWPRegistryAllocated // Event containing the contract specifics and raw log

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
func (it *AWPRegistryAllocatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryAllocated)
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
		it.Event = new(AWPRegistryAllocated)
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
func (it *AWPRegistryAllocatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryAllocatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryAllocated represents a Allocated event raised by the AWPRegistry contract.
type AWPRegistryAllocated struct {
	Staker   common.Address
	Agent    common.Address
	SubnetId *big.Int
	Amount   *big.Int
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAllocated is a free log retrieval operation binding the contract event 0x655f98c7dae1bab3e2db10cdb4407717b9d219cf2e585adc1edba92d48af2b15.
//
// Solidity: event Allocated(address indexed staker, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator)
func (_AWPRegistry *AWPRegistryFilterer) FilterAllocated(opts *bind.FilterOpts, staker []common.Address, agent []common.Address, subnetId []*big.Int) (*AWPRegistryAllocatedIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "Allocated", stakerRule, agentRule, subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistryAllocatedIterator{contract: _AWPRegistry.contract, event: "Allocated", logs: logs, sub: sub}, nil
}

// WatchAllocated is a free log subscription operation binding the contract event 0x655f98c7dae1bab3e2db10cdb4407717b9d219cf2e585adc1edba92d48af2b15.
//
// Solidity: event Allocated(address indexed staker, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator)
func (_AWPRegistry *AWPRegistryFilterer) WatchAllocated(opts *bind.WatchOpts, sink chan<- *AWPRegistryAllocated, staker []common.Address, agent []common.Address, subnetId []*big.Int) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "Allocated", stakerRule, agentRule, subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryAllocated)
				if err := _AWPRegistry.contract.UnpackLog(event, "Allocated", log); err != nil {
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
// Solidity: event Allocated(address indexed staker, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator)
func (_AWPRegistry *AWPRegistryFilterer) ParseAllocated(log types.Log) (*AWPRegistryAllocated, error) {
	event := new(AWPRegistryAllocated)
	if err := _AWPRegistry.contract.UnpackLog(event, "Allocated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistryAlphaTokenFactoryUpdatedIterator is returned from FilterAlphaTokenFactoryUpdated and is used to iterate over the raw logs and unpacked data for AlphaTokenFactoryUpdated events raised by the AWPRegistry contract.
type AWPRegistryAlphaTokenFactoryUpdatedIterator struct {
	Event *AWPRegistryAlphaTokenFactoryUpdated // Event containing the contract specifics and raw log

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
func (it *AWPRegistryAlphaTokenFactoryUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryAlphaTokenFactoryUpdated)
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
		it.Event = new(AWPRegistryAlphaTokenFactoryUpdated)
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
func (it *AWPRegistryAlphaTokenFactoryUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryAlphaTokenFactoryUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryAlphaTokenFactoryUpdated represents a AlphaTokenFactoryUpdated event raised by the AWPRegistry contract.
type AWPRegistryAlphaTokenFactoryUpdated struct {
	NewFactory common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAlphaTokenFactoryUpdated is a free log retrieval operation binding the contract event 0xca3b5054bdfbf81973dd36029b7ef8c5479d0739433700df6b2e6d690ead4a3e.
//
// Solidity: event AlphaTokenFactoryUpdated(address indexed newFactory)
func (_AWPRegistry *AWPRegistryFilterer) FilterAlphaTokenFactoryUpdated(opts *bind.FilterOpts, newFactory []common.Address) (*AWPRegistryAlphaTokenFactoryUpdatedIterator, error) {

	var newFactoryRule []interface{}
	for _, newFactoryItem := range newFactory {
		newFactoryRule = append(newFactoryRule, newFactoryItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "AlphaTokenFactoryUpdated", newFactoryRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistryAlphaTokenFactoryUpdatedIterator{contract: _AWPRegistry.contract, event: "AlphaTokenFactoryUpdated", logs: logs, sub: sub}, nil
}

// WatchAlphaTokenFactoryUpdated is a free log subscription operation binding the contract event 0xca3b5054bdfbf81973dd36029b7ef8c5479d0739433700df6b2e6d690ead4a3e.
//
// Solidity: event AlphaTokenFactoryUpdated(address indexed newFactory)
func (_AWPRegistry *AWPRegistryFilterer) WatchAlphaTokenFactoryUpdated(opts *bind.WatchOpts, sink chan<- *AWPRegistryAlphaTokenFactoryUpdated, newFactory []common.Address) (event.Subscription, error) {

	var newFactoryRule []interface{}
	for _, newFactoryItem := range newFactory {
		newFactoryRule = append(newFactoryRule, newFactoryItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "AlphaTokenFactoryUpdated", newFactoryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryAlphaTokenFactoryUpdated)
				if err := _AWPRegistry.contract.UnpackLog(event, "AlphaTokenFactoryUpdated", log); err != nil {
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
func (_AWPRegistry *AWPRegistryFilterer) ParseAlphaTokenFactoryUpdated(log types.Log) (*AWPRegistryAlphaTokenFactoryUpdated, error) {
	event := new(AWPRegistryAlphaTokenFactoryUpdated)
	if err := _AWPRegistry.contract.UnpackLog(event, "AlphaTokenFactoryUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistryBoundIterator is returned from FilterBound and is used to iterate over the raw logs and unpacked data for Bound events raised by the AWPRegistry contract.
type AWPRegistryBoundIterator struct {
	Event *AWPRegistryBound // Event containing the contract specifics and raw log

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
func (it *AWPRegistryBoundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryBound)
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
		it.Event = new(AWPRegistryBound)
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
func (it *AWPRegistryBoundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryBoundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryBound represents a Bound event raised by the AWPRegistry contract.
type AWPRegistryBound struct {
	Addr   common.Address
	Target common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBound is a free log retrieval operation binding the contract event 0x0d128562eaa47ab89086803e64a0f96847c0ed3cc63c26251f29ba1aede09d4e.
//
// Solidity: event Bound(address indexed addr, address indexed target)
func (_AWPRegistry *AWPRegistryFilterer) FilterBound(opts *bind.FilterOpts, addr []common.Address, target []common.Address) (*AWPRegistryBoundIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "Bound", addrRule, targetRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistryBoundIterator{contract: _AWPRegistry.contract, event: "Bound", logs: logs, sub: sub}, nil
}

// WatchBound is a free log subscription operation binding the contract event 0x0d128562eaa47ab89086803e64a0f96847c0ed3cc63c26251f29ba1aede09d4e.
//
// Solidity: event Bound(address indexed addr, address indexed target)
func (_AWPRegistry *AWPRegistryFilterer) WatchBound(opts *bind.WatchOpts, sink chan<- *AWPRegistryBound, addr []common.Address, target []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "Bound", addrRule, targetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryBound)
				if err := _AWPRegistry.contract.UnpackLog(event, "Bound", log); err != nil {
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

// ParseBound is a log parse operation binding the contract event 0x0d128562eaa47ab89086803e64a0f96847c0ed3cc63c26251f29ba1aede09d4e.
//
// Solidity: event Bound(address indexed addr, address indexed target)
func (_AWPRegistry *AWPRegistryFilterer) ParseBound(log types.Log) (*AWPRegistryBound, error) {
	event := new(AWPRegistryBound)
	if err := _AWPRegistry.contract.UnpackLog(event, "Bound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistryDeallocatedIterator is returned from FilterDeallocated and is used to iterate over the raw logs and unpacked data for Deallocated events raised by the AWPRegistry contract.
type AWPRegistryDeallocatedIterator struct {
	Event *AWPRegistryDeallocated // Event containing the contract specifics and raw log

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
func (it *AWPRegistryDeallocatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryDeallocated)
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
		it.Event = new(AWPRegistryDeallocated)
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
func (it *AWPRegistryDeallocatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryDeallocatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryDeallocated represents a Deallocated event raised by the AWPRegistry contract.
type AWPRegistryDeallocated struct {
	Staker   common.Address
	Agent    common.Address
	SubnetId *big.Int
	Amount   *big.Int
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDeallocated is a free log retrieval operation binding the contract event 0xd55bd7964253d1d9ce9187c8187b1c904274a3f374c9074f6de6fa77746bf345.
//
// Solidity: event Deallocated(address indexed staker, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator)
func (_AWPRegistry *AWPRegistryFilterer) FilterDeallocated(opts *bind.FilterOpts, staker []common.Address, agent []common.Address, subnetId []*big.Int) (*AWPRegistryDeallocatedIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "Deallocated", stakerRule, agentRule, subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistryDeallocatedIterator{contract: _AWPRegistry.contract, event: "Deallocated", logs: logs, sub: sub}, nil
}

// WatchDeallocated is a free log subscription operation binding the contract event 0xd55bd7964253d1d9ce9187c8187b1c904274a3f374c9074f6de6fa77746bf345.
//
// Solidity: event Deallocated(address indexed staker, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator)
func (_AWPRegistry *AWPRegistryFilterer) WatchDeallocated(opts *bind.WatchOpts, sink chan<- *AWPRegistryDeallocated, staker []common.Address, agent []common.Address, subnetId []*big.Int) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "Deallocated", stakerRule, agentRule, subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryDeallocated)
				if err := _AWPRegistry.contract.UnpackLog(event, "Deallocated", log); err != nil {
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
// Solidity: event Deallocated(address indexed staker, address indexed agent, uint256 indexed subnetId, uint256 amount, address operator)
func (_AWPRegistry *AWPRegistryFilterer) ParseDeallocated(log types.Log) (*AWPRegistryDeallocated, error) {
	event := new(AWPRegistryDeallocated)
	if err := _AWPRegistry.contract.UnpackLog(event, "Deallocated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistryDefaultSubnetManagerImplUpdatedIterator is returned from FilterDefaultSubnetManagerImplUpdated and is used to iterate over the raw logs and unpacked data for DefaultSubnetManagerImplUpdated events raised by the AWPRegistry contract.
type AWPRegistryDefaultSubnetManagerImplUpdatedIterator struct {
	Event *AWPRegistryDefaultSubnetManagerImplUpdated // Event containing the contract specifics and raw log

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
func (it *AWPRegistryDefaultSubnetManagerImplUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryDefaultSubnetManagerImplUpdated)
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
		it.Event = new(AWPRegistryDefaultSubnetManagerImplUpdated)
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
func (it *AWPRegistryDefaultSubnetManagerImplUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryDefaultSubnetManagerImplUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryDefaultSubnetManagerImplUpdated represents a DefaultSubnetManagerImplUpdated event raised by the AWPRegistry contract.
type AWPRegistryDefaultSubnetManagerImplUpdated struct {
	NewImpl common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterDefaultSubnetManagerImplUpdated is a free log retrieval operation binding the contract event 0xa37cb79f631c6bb2a11d965d06cce40e3c936eba1649879b8ffa233c0219f949.
//
// Solidity: event DefaultSubnetManagerImplUpdated(address indexed newImpl)
func (_AWPRegistry *AWPRegistryFilterer) FilterDefaultSubnetManagerImplUpdated(opts *bind.FilterOpts, newImpl []common.Address) (*AWPRegistryDefaultSubnetManagerImplUpdatedIterator, error) {

	var newImplRule []interface{}
	for _, newImplItem := range newImpl {
		newImplRule = append(newImplRule, newImplItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "DefaultSubnetManagerImplUpdated", newImplRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistryDefaultSubnetManagerImplUpdatedIterator{contract: _AWPRegistry.contract, event: "DefaultSubnetManagerImplUpdated", logs: logs, sub: sub}, nil
}

// WatchDefaultSubnetManagerImplUpdated is a free log subscription operation binding the contract event 0xa37cb79f631c6bb2a11d965d06cce40e3c936eba1649879b8ffa233c0219f949.
//
// Solidity: event DefaultSubnetManagerImplUpdated(address indexed newImpl)
func (_AWPRegistry *AWPRegistryFilterer) WatchDefaultSubnetManagerImplUpdated(opts *bind.WatchOpts, sink chan<- *AWPRegistryDefaultSubnetManagerImplUpdated, newImpl []common.Address) (event.Subscription, error) {

	var newImplRule []interface{}
	for _, newImplItem := range newImpl {
		newImplRule = append(newImplRule, newImplItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "DefaultSubnetManagerImplUpdated", newImplRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryDefaultSubnetManagerImplUpdated)
				if err := _AWPRegistry.contract.UnpackLog(event, "DefaultSubnetManagerImplUpdated", log); err != nil {
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
func (_AWPRegistry *AWPRegistryFilterer) ParseDefaultSubnetManagerImplUpdated(log types.Log) (*AWPRegistryDefaultSubnetManagerImplUpdated, error) {
	event := new(AWPRegistryDefaultSubnetManagerImplUpdated)
	if err := _AWPRegistry.contract.UnpackLog(event, "DefaultSubnetManagerImplUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistryDelegateGrantedIterator is returned from FilterDelegateGranted and is used to iterate over the raw logs and unpacked data for DelegateGranted events raised by the AWPRegistry contract.
type AWPRegistryDelegateGrantedIterator struct {
	Event *AWPRegistryDelegateGranted // Event containing the contract specifics and raw log

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
func (it *AWPRegistryDelegateGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryDelegateGranted)
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
		it.Event = new(AWPRegistryDelegateGranted)
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
func (it *AWPRegistryDelegateGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryDelegateGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryDelegateGranted represents a DelegateGranted event raised by the AWPRegistry contract.
type AWPRegistryDelegateGranted struct {
	Staker   common.Address
	Delegate common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDelegateGranted is a free log retrieval operation binding the contract event 0x0cd335986c24e121f32e8e0fd34f998524f62b9de25106243d284f86394bc2e9.
//
// Solidity: event DelegateGranted(address indexed staker, address indexed delegate)
func (_AWPRegistry *AWPRegistryFilterer) FilterDelegateGranted(opts *bind.FilterOpts, staker []common.Address, delegate []common.Address) (*AWPRegistryDelegateGrantedIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var delegateRule []interface{}
	for _, delegateItem := range delegate {
		delegateRule = append(delegateRule, delegateItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "DelegateGranted", stakerRule, delegateRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistryDelegateGrantedIterator{contract: _AWPRegistry.contract, event: "DelegateGranted", logs: logs, sub: sub}, nil
}

// WatchDelegateGranted is a free log subscription operation binding the contract event 0x0cd335986c24e121f32e8e0fd34f998524f62b9de25106243d284f86394bc2e9.
//
// Solidity: event DelegateGranted(address indexed staker, address indexed delegate)
func (_AWPRegistry *AWPRegistryFilterer) WatchDelegateGranted(opts *bind.WatchOpts, sink chan<- *AWPRegistryDelegateGranted, staker []common.Address, delegate []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var delegateRule []interface{}
	for _, delegateItem := range delegate {
		delegateRule = append(delegateRule, delegateItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "DelegateGranted", stakerRule, delegateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryDelegateGranted)
				if err := _AWPRegistry.contract.UnpackLog(event, "DelegateGranted", log); err != nil {
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

// ParseDelegateGranted is a log parse operation binding the contract event 0x0cd335986c24e121f32e8e0fd34f998524f62b9de25106243d284f86394bc2e9.
//
// Solidity: event DelegateGranted(address indexed staker, address indexed delegate)
func (_AWPRegistry *AWPRegistryFilterer) ParseDelegateGranted(log types.Log) (*AWPRegistryDelegateGranted, error) {
	event := new(AWPRegistryDelegateGranted)
	if err := _AWPRegistry.contract.UnpackLog(event, "DelegateGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistryDelegateRevokedIterator is returned from FilterDelegateRevoked and is used to iterate over the raw logs and unpacked data for DelegateRevoked events raised by the AWPRegistry contract.
type AWPRegistryDelegateRevokedIterator struct {
	Event *AWPRegistryDelegateRevoked // Event containing the contract specifics and raw log

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
func (it *AWPRegistryDelegateRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryDelegateRevoked)
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
		it.Event = new(AWPRegistryDelegateRevoked)
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
func (it *AWPRegistryDelegateRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryDelegateRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryDelegateRevoked represents a DelegateRevoked event raised by the AWPRegistry contract.
type AWPRegistryDelegateRevoked struct {
	Staker   common.Address
	Delegate common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDelegateRevoked is a free log retrieval operation binding the contract event 0x76e6646868d096078ac7f3f1401c3aaa55dc84890ec74b99c699e4754714b18e.
//
// Solidity: event DelegateRevoked(address indexed staker, address indexed delegate)
func (_AWPRegistry *AWPRegistryFilterer) FilterDelegateRevoked(opts *bind.FilterOpts, staker []common.Address, delegate []common.Address) (*AWPRegistryDelegateRevokedIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var delegateRule []interface{}
	for _, delegateItem := range delegate {
		delegateRule = append(delegateRule, delegateItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "DelegateRevoked", stakerRule, delegateRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistryDelegateRevokedIterator{contract: _AWPRegistry.contract, event: "DelegateRevoked", logs: logs, sub: sub}, nil
}

// WatchDelegateRevoked is a free log subscription operation binding the contract event 0x76e6646868d096078ac7f3f1401c3aaa55dc84890ec74b99c699e4754714b18e.
//
// Solidity: event DelegateRevoked(address indexed staker, address indexed delegate)
func (_AWPRegistry *AWPRegistryFilterer) WatchDelegateRevoked(opts *bind.WatchOpts, sink chan<- *AWPRegistryDelegateRevoked, staker []common.Address, delegate []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}
	var delegateRule []interface{}
	for _, delegateItem := range delegate {
		delegateRule = append(delegateRule, delegateItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "DelegateRevoked", stakerRule, delegateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryDelegateRevoked)
				if err := _AWPRegistry.contract.UnpackLog(event, "DelegateRevoked", log); err != nil {
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

// ParseDelegateRevoked is a log parse operation binding the contract event 0x76e6646868d096078ac7f3f1401c3aaa55dc84890ec74b99c699e4754714b18e.
//
// Solidity: event DelegateRevoked(address indexed staker, address indexed delegate)
func (_AWPRegistry *AWPRegistryFilterer) ParseDelegateRevoked(log types.Log) (*AWPRegistryDelegateRevoked, error) {
	event := new(AWPRegistryDelegateRevoked)
	if err := _AWPRegistry.contract.UnpackLog(event, "DelegateRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistryEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the AWPRegistry contract.
type AWPRegistryEIP712DomainChangedIterator struct {
	Event *AWPRegistryEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *AWPRegistryEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryEIP712DomainChanged)
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
		it.Event = new(AWPRegistryEIP712DomainChanged)
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
func (it *AWPRegistryEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryEIP712DomainChanged represents a EIP712DomainChanged event raised by the AWPRegistry contract.
type AWPRegistryEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_AWPRegistry *AWPRegistryFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*AWPRegistryEIP712DomainChangedIterator, error) {

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &AWPRegistryEIP712DomainChangedIterator{contract: _AWPRegistry.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_AWPRegistry *AWPRegistryFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *AWPRegistryEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryEIP712DomainChanged)
				if err := _AWPRegistry.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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
func (_AWPRegistry *AWPRegistryFilterer) ParseEIP712DomainChanged(log types.Log) (*AWPRegistryEIP712DomainChanged, error) {
	event := new(AWPRegistryEIP712DomainChanged)
	if err := _AWPRegistry.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistryGuardianUpdatedIterator is returned from FilterGuardianUpdated and is used to iterate over the raw logs and unpacked data for GuardianUpdated events raised by the AWPRegistry contract.
type AWPRegistryGuardianUpdatedIterator struct {
	Event *AWPRegistryGuardianUpdated // Event containing the contract specifics and raw log

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
func (it *AWPRegistryGuardianUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryGuardianUpdated)
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
		it.Event = new(AWPRegistryGuardianUpdated)
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
func (it *AWPRegistryGuardianUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryGuardianUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryGuardianUpdated represents a GuardianUpdated event raised by the AWPRegistry contract.
type AWPRegistryGuardianUpdated struct {
	NewGuardian common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterGuardianUpdated is a free log retrieval operation binding the contract event 0x6bb7ff33e730289800c62ad882105a144a74010d2bdbb9a942544a3005ad55bf.
//
// Solidity: event GuardianUpdated(address indexed newGuardian)
func (_AWPRegistry *AWPRegistryFilterer) FilterGuardianUpdated(opts *bind.FilterOpts, newGuardian []common.Address) (*AWPRegistryGuardianUpdatedIterator, error) {

	var newGuardianRule []interface{}
	for _, newGuardianItem := range newGuardian {
		newGuardianRule = append(newGuardianRule, newGuardianItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "GuardianUpdated", newGuardianRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistryGuardianUpdatedIterator{contract: _AWPRegistry.contract, event: "GuardianUpdated", logs: logs, sub: sub}, nil
}

// WatchGuardianUpdated is a free log subscription operation binding the contract event 0x6bb7ff33e730289800c62ad882105a144a74010d2bdbb9a942544a3005ad55bf.
//
// Solidity: event GuardianUpdated(address indexed newGuardian)
func (_AWPRegistry *AWPRegistryFilterer) WatchGuardianUpdated(opts *bind.WatchOpts, sink chan<- *AWPRegistryGuardianUpdated, newGuardian []common.Address) (event.Subscription, error) {

	var newGuardianRule []interface{}
	for _, newGuardianItem := range newGuardian {
		newGuardianRule = append(newGuardianRule, newGuardianItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "GuardianUpdated", newGuardianRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryGuardianUpdated)
				if err := _AWPRegistry.contract.UnpackLog(event, "GuardianUpdated", log); err != nil {
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
func (_AWPRegistry *AWPRegistryFilterer) ParseGuardianUpdated(log types.Log) (*AWPRegistryGuardianUpdated, error) {
	event := new(AWPRegistryGuardianUpdated)
	if err := _AWPRegistry.contract.UnpackLog(event, "GuardianUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistryImmunityPeriodUpdatedIterator is returned from FilterImmunityPeriodUpdated and is used to iterate over the raw logs and unpacked data for ImmunityPeriodUpdated events raised by the AWPRegistry contract.
type AWPRegistryImmunityPeriodUpdatedIterator struct {
	Event *AWPRegistryImmunityPeriodUpdated // Event containing the contract specifics and raw log

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
func (it *AWPRegistryImmunityPeriodUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryImmunityPeriodUpdated)
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
		it.Event = new(AWPRegistryImmunityPeriodUpdated)
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
func (it *AWPRegistryImmunityPeriodUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryImmunityPeriodUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryImmunityPeriodUpdated represents a ImmunityPeriodUpdated event raised by the AWPRegistry contract.
type AWPRegistryImmunityPeriodUpdated struct {
	NewPeriod *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterImmunityPeriodUpdated is a free log retrieval operation binding the contract event 0x49b186851943e5bbcefec9411c3238262c6e102e4000142f8f060143d1b8724c.
//
// Solidity: event ImmunityPeriodUpdated(uint256 newPeriod)
func (_AWPRegistry *AWPRegistryFilterer) FilterImmunityPeriodUpdated(opts *bind.FilterOpts) (*AWPRegistryImmunityPeriodUpdatedIterator, error) {

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "ImmunityPeriodUpdated")
	if err != nil {
		return nil, err
	}
	return &AWPRegistryImmunityPeriodUpdatedIterator{contract: _AWPRegistry.contract, event: "ImmunityPeriodUpdated", logs: logs, sub: sub}, nil
}

// WatchImmunityPeriodUpdated is a free log subscription operation binding the contract event 0x49b186851943e5bbcefec9411c3238262c6e102e4000142f8f060143d1b8724c.
//
// Solidity: event ImmunityPeriodUpdated(uint256 newPeriod)
func (_AWPRegistry *AWPRegistryFilterer) WatchImmunityPeriodUpdated(opts *bind.WatchOpts, sink chan<- *AWPRegistryImmunityPeriodUpdated) (event.Subscription, error) {

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "ImmunityPeriodUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryImmunityPeriodUpdated)
				if err := _AWPRegistry.contract.UnpackLog(event, "ImmunityPeriodUpdated", log); err != nil {
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
func (_AWPRegistry *AWPRegistryFilterer) ParseImmunityPeriodUpdated(log types.Log) (*AWPRegistryImmunityPeriodUpdated, error) {
	event := new(AWPRegistryImmunityPeriodUpdated)
	if err := _AWPRegistry.contract.UnpackLog(event, "ImmunityPeriodUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistryInitialAlphaPriceUpdatedIterator is returned from FilterInitialAlphaPriceUpdated and is used to iterate over the raw logs and unpacked data for InitialAlphaPriceUpdated events raised by the AWPRegistry contract.
type AWPRegistryInitialAlphaPriceUpdatedIterator struct {
	Event *AWPRegistryInitialAlphaPriceUpdated // Event containing the contract specifics and raw log

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
func (it *AWPRegistryInitialAlphaPriceUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryInitialAlphaPriceUpdated)
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
		it.Event = new(AWPRegistryInitialAlphaPriceUpdated)
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
func (it *AWPRegistryInitialAlphaPriceUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryInitialAlphaPriceUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryInitialAlphaPriceUpdated represents a InitialAlphaPriceUpdated event raised by the AWPRegistry contract.
type AWPRegistryInitialAlphaPriceUpdated struct {
	NewPrice *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterInitialAlphaPriceUpdated is a free log retrieval operation binding the contract event 0xab7ee876750d22d253d0b38988caea5f6285a832697e4889d9beb36515dde34e.
//
// Solidity: event InitialAlphaPriceUpdated(uint256 newPrice)
func (_AWPRegistry *AWPRegistryFilterer) FilterInitialAlphaPriceUpdated(opts *bind.FilterOpts) (*AWPRegistryInitialAlphaPriceUpdatedIterator, error) {

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "InitialAlphaPriceUpdated")
	if err != nil {
		return nil, err
	}
	return &AWPRegistryInitialAlphaPriceUpdatedIterator{contract: _AWPRegistry.contract, event: "InitialAlphaPriceUpdated", logs: logs, sub: sub}, nil
}

// WatchInitialAlphaPriceUpdated is a free log subscription operation binding the contract event 0xab7ee876750d22d253d0b38988caea5f6285a832697e4889d9beb36515dde34e.
//
// Solidity: event InitialAlphaPriceUpdated(uint256 newPrice)
func (_AWPRegistry *AWPRegistryFilterer) WatchInitialAlphaPriceUpdated(opts *bind.WatchOpts, sink chan<- *AWPRegistryInitialAlphaPriceUpdated) (event.Subscription, error) {

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "InitialAlphaPriceUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryInitialAlphaPriceUpdated)
				if err := _AWPRegistry.contract.UnpackLog(event, "InitialAlphaPriceUpdated", log); err != nil {
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
func (_AWPRegistry *AWPRegistryFilterer) ParseInitialAlphaPriceUpdated(log types.Log) (*AWPRegistryInitialAlphaPriceUpdated, error) {
	event := new(AWPRegistryInitialAlphaPriceUpdated)
	if err := _AWPRegistry.contract.UnpackLog(event, "InitialAlphaPriceUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistryLPCreatedIterator is returned from FilterLPCreated and is used to iterate over the raw logs and unpacked data for LPCreated events raised by the AWPRegistry contract.
type AWPRegistryLPCreatedIterator struct {
	Event *AWPRegistryLPCreated // Event containing the contract specifics and raw log

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
func (it *AWPRegistryLPCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryLPCreated)
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
		it.Event = new(AWPRegistryLPCreated)
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
func (it *AWPRegistryLPCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryLPCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryLPCreated represents a LPCreated event raised by the AWPRegistry contract.
type AWPRegistryLPCreated struct {
	SubnetId    *big.Int
	PoolId      [32]byte
	AwpAmount   *big.Int
	AlphaAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterLPCreated is a free log retrieval operation binding the contract event 0x0a28a1fd5e0909199ee082834df66cfaae2125e3503bf16d2dc46214278fc7ab.
//
// Solidity: event LPCreated(uint256 indexed subnetId, bytes32 poolId, uint256 awpAmount, uint256 alphaAmount)
func (_AWPRegistry *AWPRegistryFilterer) FilterLPCreated(opts *bind.FilterOpts, subnetId []*big.Int) (*AWPRegistryLPCreatedIterator, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "LPCreated", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistryLPCreatedIterator{contract: _AWPRegistry.contract, event: "LPCreated", logs: logs, sub: sub}, nil
}

// WatchLPCreated is a free log subscription operation binding the contract event 0x0a28a1fd5e0909199ee082834df66cfaae2125e3503bf16d2dc46214278fc7ab.
//
// Solidity: event LPCreated(uint256 indexed subnetId, bytes32 poolId, uint256 awpAmount, uint256 alphaAmount)
func (_AWPRegistry *AWPRegistryFilterer) WatchLPCreated(opts *bind.WatchOpts, sink chan<- *AWPRegistryLPCreated, subnetId []*big.Int) (event.Subscription, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "LPCreated", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryLPCreated)
				if err := _AWPRegistry.contract.UnpackLog(event, "LPCreated", log); err != nil {
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
func (_AWPRegistry *AWPRegistryFilterer) ParseLPCreated(log types.Log) (*AWPRegistryLPCreated, error) {
	event := new(AWPRegistryLPCreated)
	if err := _AWPRegistry.contract.UnpackLog(event, "LPCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistryPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the AWPRegistry contract.
type AWPRegistryPausedIterator struct {
	Event *AWPRegistryPaused // Event containing the contract specifics and raw log

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
func (it *AWPRegistryPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryPaused)
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
		it.Event = new(AWPRegistryPaused)
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
func (it *AWPRegistryPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryPaused represents a Paused event raised by the AWPRegistry contract.
type AWPRegistryPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_AWPRegistry *AWPRegistryFilterer) FilterPaused(opts *bind.FilterOpts) (*AWPRegistryPausedIterator, error) {

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &AWPRegistryPausedIterator{contract: _AWPRegistry.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_AWPRegistry *AWPRegistryFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *AWPRegistryPaused) (event.Subscription, error) {

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryPaused)
				if err := _AWPRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_AWPRegistry *AWPRegistryFilterer) ParsePaused(log types.Log) (*AWPRegistryPaused, error) {
	event := new(AWPRegistryPaused)
	if err := _AWPRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistryReallocatedIterator is returned from FilterReallocated and is used to iterate over the raw logs and unpacked data for Reallocated events raised by the AWPRegistry contract.
type AWPRegistryReallocatedIterator struct {
	Event *AWPRegistryReallocated // Event containing the contract specifics and raw log

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
func (it *AWPRegistryReallocatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryReallocated)
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
		it.Event = new(AWPRegistryReallocated)
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
func (it *AWPRegistryReallocatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryReallocatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryReallocated represents a Reallocated event raised by the AWPRegistry contract.
type AWPRegistryReallocated struct {
	Staker     common.Address
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
// Solidity: event Reallocated(address indexed staker, address fromAgent, uint256 fromSubnet, address toAgent, uint256 toSubnet, uint256 amount, address operator)
func (_AWPRegistry *AWPRegistryFilterer) FilterReallocated(opts *bind.FilterOpts, staker []common.Address) (*AWPRegistryReallocatedIterator, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "Reallocated", stakerRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistryReallocatedIterator{contract: _AWPRegistry.contract, event: "Reallocated", logs: logs, sub: sub}, nil
}

// WatchReallocated is a free log subscription operation binding the contract event 0x726c93ba67bfe4c677e37114279f0ad9aab5ee9ffbd1158923be5d0fec3b1b45.
//
// Solidity: event Reallocated(address indexed staker, address fromAgent, uint256 fromSubnet, address toAgent, uint256 toSubnet, uint256 amount, address operator)
func (_AWPRegistry *AWPRegistryFilterer) WatchReallocated(opts *bind.WatchOpts, sink chan<- *AWPRegistryReallocated, staker []common.Address) (event.Subscription, error) {

	var stakerRule []interface{}
	for _, stakerItem := range staker {
		stakerRule = append(stakerRule, stakerItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "Reallocated", stakerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryReallocated)
				if err := _AWPRegistry.contract.UnpackLog(event, "Reallocated", log); err != nil {
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
// Solidity: event Reallocated(address indexed staker, address fromAgent, uint256 fromSubnet, address toAgent, uint256 toSubnet, uint256 amount, address operator)
func (_AWPRegistry *AWPRegistryFilterer) ParseReallocated(log types.Log) (*AWPRegistryReallocated, error) {
	event := new(AWPRegistryReallocated)
	if err := _AWPRegistry.contract.UnpackLog(event, "Reallocated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistryRecipientSetIterator is returned from FilterRecipientSet and is used to iterate over the raw logs and unpacked data for RecipientSet events raised by the AWPRegistry contract.
type AWPRegistryRecipientSetIterator struct {
	Event *AWPRegistryRecipientSet // Event containing the contract specifics and raw log

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
func (it *AWPRegistryRecipientSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryRecipientSet)
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
		it.Event = new(AWPRegistryRecipientSet)
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
func (it *AWPRegistryRecipientSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryRecipientSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryRecipientSet represents a RecipientSet event raised by the AWPRegistry contract.
type AWPRegistryRecipientSet struct {
	Addr      common.Address
	Recipient common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRecipientSet is a free log retrieval operation binding the contract event 0xc1416b5cdab50a9fbc872236e1aa54566c6deb40024e63a4b1737ecacf09d6f9.
//
// Solidity: event RecipientSet(address indexed addr, address recipient)
func (_AWPRegistry *AWPRegistryFilterer) FilterRecipientSet(opts *bind.FilterOpts, addr []common.Address) (*AWPRegistryRecipientSetIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "RecipientSet", addrRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistryRecipientSetIterator{contract: _AWPRegistry.contract, event: "RecipientSet", logs: logs, sub: sub}, nil
}

// WatchRecipientSet is a free log subscription operation binding the contract event 0xc1416b5cdab50a9fbc872236e1aa54566c6deb40024e63a4b1737ecacf09d6f9.
//
// Solidity: event RecipientSet(address indexed addr, address recipient)
func (_AWPRegistry *AWPRegistryFilterer) WatchRecipientSet(opts *bind.WatchOpts, sink chan<- *AWPRegistryRecipientSet, addr []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "RecipientSet", addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryRecipientSet)
				if err := _AWPRegistry.contract.UnpackLog(event, "RecipientSet", log); err != nil {
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

// ParseRecipientSet is a log parse operation binding the contract event 0xc1416b5cdab50a9fbc872236e1aa54566c6deb40024e63a4b1737ecacf09d6f9.
//
// Solidity: event RecipientSet(address indexed addr, address recipient)
func (_AWPRegistry *AWPRegistryFilterer) ParseRecipientSet(log types.Log) (*AWPRegistryRecipientSet, error) {
	event := new(AWPRegistryRecipientSet)
	if err := _AWPRegistry.contract.UnpackLog(event, "RecipientSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistrySubnetActivatedIterator is returned from FilterSubnetActivated and is used to iterate over the raw logs and unpacked data for SubnetActivated events raised by the AWPRegistry contract.
type AWPRegistrySubnetActivatedIterator struct {
	Event *AWPRegistrySubnetActivated // Event containing the contract specifics and raw log

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
func (it *AWPRegistrySubnetActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistrySubnetActivated)
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
		it.Event = new(AWPRegistrySubnetActivated)
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
func (it *AWPRegistrySubnetActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistrySubnetActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistrySubnetActivated represents a SubnetActivated event raised by the AWPRegistry contract.
type AWPRegistrySubnetActivated struct {
	SubnetId *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSubnetActivated is a free log retrieval operation binding the contract event 0x034804b969efac7a0df7757ada640ffdcc09f25dbcd4582c390f25d5622255c4.
//
// Solidity: event SubnetActivated(uint256 indexed subnetId)
func (_AWPRegistry *AWPRegistryFilterer) FilterSubnetActivated(opts *bind.FilterOpts, subnetId []*big.Int) (*AWPRegistrySubnetActivatedIterator, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "SubnetActivated", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistrySubnetActivatedIterator{contract: _AWPRegistry.contract, event: "SubnetActivated", logs: logs, sub: sub}, nil
}

// WatchSubnetActivated is a free log subscription operation binding the contract event 0x034804b969efac7a0df7757ada640ffdcc09f25dbcd4582c390f25d5622255c4.
//
// Solidity: event SubnetActivated(uint256 indexed subnetId)
func (_AWPRegistry *AWPRegistryFilterer) WatchSubnetActivated(opts *bind.WatchOpts, sink chan<- *AWPRegistrySubnetActivated, subnetId []*big.Int) (event.Subscription, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "SubnetActivated", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistrySubnetActivated)
				if err := _AWPRegistry.contract.UnpackLog(event, "SubnetActivated", log); err != nil {
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
func (_AWPRegistry *AWPRegistryFilterer) ParseSubnetActivated(log types.Log) (*AWPRegistrySubnetActivated, error) {
	event := new(AWPRegistrySubnetActivated)
	if err := _AWPRegistry.contract.UnpackLog(event, "SubnetActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistrySubnetBannedIterator is returned from FilterSubnetBanned and is used to iterate over the raw logs and unpacked data for SubnetBanned events raised by the AWPRegistry contract.
type AWPRegistrySubnetBannedIterator struct {
	Event *AWPRegistrySubnetBanned // Event containing the contract specifics and raw log

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
func (it *AWPRegistrySubnetBannedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistrySubnetBanned)
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
		it.Event = new(AWPRegistrySubnetBanned)
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
func (it *AWPRegistrySubnetBannedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistrySubnetBannedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistrySubnetBanned represents a SubnetBanned event raised by the AWPRegistry contract.
type AWPRegistrySubnetBanned struct {
	SubnetId *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSubnetBanned is a free log retrieval operation binding the contract event 0xb887f21153957bddcf7211314cf42794076ccf98c6ae5cf8e2d065ec717c681b.
//
// Solidity: event SubnetBanned(uint256 indexed subnetId)
func (_AWPRegistry *AWPRegistryFilterer) FilterSubnetBanned(opts *bind.FilterOpts, subnetId []*big.Int) (*AWPRegistrySubnetBannedIterator, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "SubnetBanned", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistrySubnetBannedIterator{contract: _AWPRegistry.contract, event: "SubnetBanned", logs: logs, sub: sub}, nil
}

// WatchSubnetBanned is a free log subscription operation binding the contract event 0xb887f21153957bddcf7211314cf42794076ccf98c6ae5cf8e2d065ec717c681b.
//
// Solidity: event SubnetBanned(uint256 indexed subnetId)
func (_AWPRegistry *AWPRegistryFilterer) WatchSubnetBanned(opts *bind.WatchOpts, sink chan<- *AWPRegistrySubnetBanned, subnetId []*big.Int) (event.Subscription, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "SubnetBanned", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistrySubnetBanned)
				if err := _AWPRegistry.contract.UnpackLog(event, "SubnetBanned", log); err != nil {
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
func (_AWPRegistry *AWPRegistryFilterer) ParseSubnetBanned(log types.Log) (*AWPRegistrySubnetBanned, error) {
	event := new(AWPRegistrySubnetBanned)
	if err := _AWPRegistry.contract.UnpackLog(event, "SubnetBanned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistrySubnetDeregisteredIterator is returned from FilterSubnetDeregistered and is used to iterate over the raw logs and unpacked data for SubnetDeregistered events raised by the AWPRegistry contract.
type AWPRegistrySubnetDeregisteredIterator struct {
	Event *AWPRegistrySubnetDeregistered // Event containing the contract specifics and raw log

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
func (it *AWPRegistrySubnetDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistrySubnetDeregistered)
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
		it.Event = new(AWPRegistrySubnetDeregistered)
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
func (it *AWPRegistrySubnetDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistrySubnetDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistrySubnetDeregistered represents a SubnetDeregistered event raised by the AWPRegistry contract.
type AWPRegistrySubnetDeregistered struct {
	SubnetId *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSubnetDeregistered is a free log retrieval operation binding the contract event 0x960c7566f4c9bb6958ff6e37a02b5ae69fa39dd75651fcc9b9a1029c01d0ff32.
//
// Solidity: event SubnetDeregistered(uint256 indexed subnetId)
func (_AWPRegistry *AWPRegistryFilterer) FilterSubnetDeregistered(opts *bind.FilterOpts, subnetId []*big.Int) (*AWPRegistrySubnetDeregisteredIterator, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "SubnetDeregistered", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistrySubnetDeregisteredIterator{contract: _AWPRegistry.contract, event: "SubnetDeregistered", logs: logs, sub: sub}, nil
}

// WatchSubnetDeregistered is a free log subscription operation binding the contract event 0x960c7566f4c9bb6958ff6e37a02b5ae69fa39dd75651fcc9b9a1029c01d0ff32.
//
// Solidity: event SubnetDeregistered(uint256 indexed subnetId)
func (_AWPRegistry *AWPRegistryFilterer) WatchSubnetDeregistered(opts *bind.WatchOpts, sink chan<- *AWPRegistrySubnetDeregistered, subnetId []*big.Int) (event.Subscription, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "SubnetDeregistered", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistrySubnetDeregistered)
				if err := _AWPRegistry.contract.UnpackLog(event, "SubnetDeregistered", log); err != nil {
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
func (_AWPRegistry *AWPRegistryFilterer) ParseSubnetDeregistered(log types.Log) (*AWPRegistrySubnetDeregistered, error) {
	event := new(AWPRegistrySubnetDeregistered)
	if err := _AWPRegistry.contract.UnpackLog(event, "SubnetDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistrySubnetPausedIterator is returned from FilterSubnetPaused and is used to iterate over the raw logs and unpacked data for SubnetPaused events raised by the AWPRegistry contract.
type AWPRegistrySubnetPausedIterator struct {
	Event *AWPRegistrySubnetPaused // Event containing the contract specifics and raw log

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
func (it *AWPRegistrySubnetPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistrySubnetPaused)
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
		it.Event = new(AWPRegistrySubnetPaused)
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
func (it *AWPRegistrySubnetPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistrySubnetPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistrySubnetPaused represents a SubnetPaused event raised by the AWPRegistry contract.
type AWPRegistrySubnetPaused struct {
	SubnetId *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSubnetPaused is a free log retrieval operation binding the contract event 0x789ca96cb827d1dcb6bfc7d9d084d0a574dadff90700e3602acedefb10f69afc.
//
// Solidity: event SubnetPaused(uint256 indexed subnetId)
func (_AWPRegistry *AWPRegistryFilterer) FilterSubnetPaused(opts *bind.FilterOpts, subnetId []*big.Int) (*AWPRegistrySubnetPausedIterator, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "SubnetPaused", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistrySubnetPausedIterator{contract: _AWPRegistry.contract, event: "SubnetPaused", logs: logs, sub: sub}, nil
}

// WatchSubnetPaused is a free log subscription operation binding the contract event 0x789ca96cb827d1dcb6bfc7d9d084d0a574dadff90700e3602acedefb10f69afc.
//
// Solidity: event SubnetPaused(uint256 indexed subnetId)
func (_AWPRegistry *AWPRegistryFilterer) WatchSubnetPaused(opts *bind.WatchOpts, sink chan<- *AWPRegistrySubnetPaused, subnetId []*big.Int) (event.Subscription, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "SubnetPaused", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistrySubnetPaused)
				if err := _AWPRegistry.contract.UnpackLog(event, "SubnetPaused", log); err != nil {
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
func (_AWPRegistry *AWPRegistryFilterer) ParseSubnetPaused(log types.Log) (*AWPRegistrySubnetPaused, error) {
	event := new(AWPRegistrySubnetPaused)
	if err := _AWPRegistry.contract.UnpackLog(event, "SubnetPaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistrySubnetRegisteredIterator is returned from FilterSubnetRegistered and is used to iterate over the raw logs and unpacked data for SubnetRegistered events raised by the AWPRegistry contract.
type AWPRegistrySubnetRegisteredIterator struct {
	Event *AWPRegistrySubnetRegistered // Event containing the contract specifics and raw log

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
func (it *AWPRegistrySubnetRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistrySubnetRegistered)
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
		it.Event = new(AWPRegistrySubnetRegistered)
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
func (it *AWPRegistrySubnetRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistrySubnetRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistrySubnetRegistered represents a SubnetRegistered event raised by the AWPRegistry contract.
type AWPRegistrySubnetRegistered struct {
	SubnetId      *big.Int
	Owner         common.Address
	Name          string
	Symbol        string
	SubnetManager common.Address
	AlphaToken    common.Address
	MinStake      *big.Int
	SkillsURI     string
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSubnetRegistered is a free log retrieval operation binding the contract event 0x57db51b7e6590b5f0b8b7dcdd5e3d1a0fd57b37aa8095cf382d699a7d0af58fe.
//
// Solidity: event SubnetRegistered(uint256 indexed subnetId, address indexed owner, string name, string symbol, address subnetManager, address alphaToken, uint128 minStake, string skillsURI)
func (_AWPRegistry *AWPRegistryFilterer) FilterSubnetRegistered(opts *bind.FilterOpts, subnetId []*big.Int, owner []common.Address) (*AWPRegistrySubnetRegisteredIterator, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "SubnetRegistered", subnetIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistrySubnetRegisteredIterator{contract: _AWPRegistry.contract, event: "SubnetRegistered", logs: logs, sub: sub}, nil
}

// WatchSubnetRegistered is a free log subscription operation binding the contract event 0x57db51b7e6590b5f0b8b7dcdd5e3d1a0fd57b37aa8095cf382d699a7d0af58fe.
//
// Solidity: event SubnetRegistered(uint256 indexed subnetId, address indexed owner, string name, string symbol, address subnetManager, address alphaToken, uint128 minStake, string skillsURI)
func (_AWPRegistry *AWPRegistryFilterer) WatchSubnetRegistered(opts *bind.WatchOpts, sink chan<- *AWPRegistrySubnetRegistered, subnetId []*big.Int, owner []common.Address) (event.Subscription, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "SubnetRegistered", subnetIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistrySubnetRegistered)
				if err := _AWPRegistry.contract.UnpackLog(event, "SubnetRegistered", log); err != nil {
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

// ParseSubnetRegistered is a log parse operation binding the contract event 0x57db51b7e6590b5f0b8b7dcdd5e3d1a0fd57b37aa8095cf382d699a7d0af58fe.
//
// Solidity: event SubnetRegistered(uint256 indexed subnetId, address indexed owner, string name, string symbol, address subnetManager, address alphaToken, uint128 minStake, string skillsURI)
func (_AWPRegistry *AWPRegistryFilterer) ParseSubnetRegistered(log types.Log) (*AWPRegistrySubnetRegistered, error) {
	event := new(AWPRegistrySubnetRegistered)
	if err := _AWPRegistry.contract.UnpackLog(event, "SubnetRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistrySubnetResumedIterator is returned from FilterSubnetResumed and is used to iterate over the raw logs and unpacked data for SubnetResumed events raised by the AWPRegistry contract.
type AWPRegistrySubnetResumedIterator struct {
	Event *AWPRegistrySubnetResumed // Event containing the contract specifics and raw log

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
func (it *AWPRegistrySubnetResumedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistrySubnetResumed)
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
		it.Event = new(AWPRegistrySubnetResumed)
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
func (it *AWPRegistrySubnetResumedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistrySubnetResumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistrySubnetResumed represents a SubnetResumed event raised by the AWPRegistry contract.
type AWPRegistrySubnetResumed struct {
	SubnetId *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSubnetResumed is a free log retrieval operation binding the contract event 0xf1693a0767c0183c95caf97ea0be785bece8e3578b49ef89c9669b754c1ba9a0.
//
// Solidity: event SubnetResumed(uint256 indexed subnetId)
func (_AWPRegistry *AWPRegistryFilterer) FilterSubnetResumed(opts *bind.FilterOpts, subnetId []*big.Int) (*AWPRegistrySubnetResumedIterator, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "SubnetResumed", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistrySubnetResumedIterator{contract: _AWPRegistry.contract, event: "SubnetResumed", logs: logs, sub: sub}, nil
}

// WatchSubnetResumed is a free log subscription operation binding the contract event 0xf1693a0767c0183c95caf97ea0be785bece8e3578b49ef89c9669b754c1ba9a0.
//
// Solidity: event SubnetResumed(uint256 indexed subnetId)
func (_AWPRegistry *AWPRegistryFilterer) WatchSubnetResumed(opts *bind.WatchOpts, sink chan<- *AWPRegistrySubnetResumed, subnetId []*big.Int) (event.Subscription, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "SubnetResumed", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistrySubnetResumed)
				if err := _AWPRegistry.contract.UnpackLog(event, "SubnetResumed", log); err != nil {
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
func (_AWPRegistry *AWPRegistryFilterer) ParseSubnetResumed(log types.Log) (*AWPRegistrySubnetResumed, error) {
	event := new(AWPRegistrySubnetResumed)
	if err := _AWPRegistry.contract.UnpackLog(event, "SubnetResumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistrySubnetUnbannedIterator is returned from FilterSubnetUnbanned and is used to iterate over the raw logs and unpacked data for SubnetUnbanned events raised by the AWPRegistry contract.
type AWPRegistrySubnetUnbannedIterator struct {
	Event *AWPRegistrySubnetUnbanned // Event containing the contract specifics and raw log

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
func (it *AWPRegistrySubnetUnbannedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistrySubnetUnbanned)
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
		it.Event = new(AWPRegistrySubnetUnbanned)
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
func (it *AWPRegistrySubnetUnbannedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistrySubnetUnbannedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistrySubnetUnbanned represents a SubnetUnbanned event raised by the AWPRegistry contract.
type AWPRegistrySubnetUnbanned struct {
	SubnetId *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSubnetUnbanned is a free log retrieval operation binding the contract event 0xa04fe0f9f3200108443db1507380606e909a0f81c9eb84c0707c265152668466.
//
// Solidity: event SubnetUnbanned(uint256 indexed subnetId)
func (_AWPRegistry *AWPRegistryFilterer) FilterSubnetUnbanned(opts *bind.FilterOpts, subnetId []*big.Int) (*AWPRegistrySubnetUnbannedIterator, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "SubnetUnbanned", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistrySubnetUnbannedIterator{contract: _AWPRegistry.contract, event: "SubnetUnbanned", logs: logs, sub: sub}, nil
}

// WatchSubnetUnbanned is a free log subscription operation binding the contract event 0xa04fe0f9f3200108443db1507380606e909a0f81c9eb84c0707c265152668466.
//
// Solidity: event SubnetUnbanned(uint256 indexed subnetId)
func (_AWPRegistry *AWPRegistryFilterer) WatchSubnetUnbanned(opts *bind.WatchOpts, sink chan<- *AWPRegistrySubnetUnbanned, subnetId []*big.Int) (event.Subscription, error) {

	var subnetIdRule []interface{}
	for _, subnetIdItem := range subnetId {
		subnetIdRule = append(subnetIdRule, subnetIdItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "SubnetUnbanned", subnetIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistrySubnetUnbanned)
				if err := _AWPRegistry.contract.UnpackLog(event, "SubnetUnbanned", log); err != nil {
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
func (_AWPRegistry *AWPRegistryFilterer) ParseSubnetUnbanned(log types.Log) (*AWPRegistrySubnetUnbanned, error) {
	event := new(AWPRegistrySubnetUnbanned)
	if err := _AWPRegistry.contract.UnpackLog(event, "SubnetUnbanned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistryUnboundIterator is returned from FilterUnbound and is used to iterate over the raw logs and unpacked data for Unbound events raised by the AWPRegistry contract.
type AWPRegistryUnboundIterator struct {
	Event *AWPRegistryUnbound // Event containing the contract specifics and raw log

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
func (it *AWPRegistryUnboundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryUnbound)
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
		it.Event = new(AWPRegistryUnbound)
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
func (it *AWPRegistryUnboundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryUnboundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryUnbound represents a Unbound event raised by the AWPRegistry contract.
type AWPRegistryUnbound struct {
	Addr common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUnbound is a free log retrieval operation binding the contract event 0x075b57b3f4efe82dc79cb35e807bfca4feaf0b4def20db9f9a9b821cba49d425.
//
// Solidity: event Unbound(address indexed addr)
func (_AWPRegistry *AWPRegistryFilterer) FilterUnbound(opts *bind.FilterOpts, addr []common.Address) (*AWPRegistryUnboundIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "Unbound", addrRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistryUnboundIterator{contract: _AWPRegistry.contract, event: "Unbound", logs: logs, sub: sub}, nil
}

// WatchUnbound is a free log subscription operation binding the contract event 0x075b57b3f4efe82dc79cb35e807bfca4feaf0b4def20db9f9a9b821cba49d425.
//
// Solidity: event Unbound(address indexed addr)
func (_AWPRegistry *AWPRegistryFilterer) WatchUnbound(opts *bind.WatchOpts, sink chan<- *AWPRegistryUnbound, addr []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "Unbound", addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryUnbound)
				if err := _AWPRegistry.contract.UnpackLog(event, "Unbound", log); err != nil {
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

// ParseUnbound is a log parse operation binding the contract event 0x075b57b3f4efe82dc79cb35e807bfca4feaf0b4def20db9f9a9b821cba49d425.
//
// Solidity: event Unbound(address indexed addr)
func (_AWPRegistry *AWPRegistryFilterer) ParseUnbound(log types.Log) (*AWPRegistryUnbound, error) {
	event := new(AWPRegistryUnbound)
	if err := _AWPRegistry.contract.UnpackLog(event, "Unbound", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistryUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the AWPRegistry contract.
type AWPRegistryUnpausedIterator struct {
	Event *AWPRegistryUnpaused // Event containing the contract specifics and raw log

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
func (it *AWPRegistryUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryUnpaused)
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
		it.Event = new(AWPRegistryUnpaused)
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
func (it *AWPRegistryUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryUnpaused represents a Unpaused event raised by the AWPRegistry contract.
type AWPRegistryUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_AWPRegistry *AWPRegistryFilterer) FilterUnpaused(opts *bind.FilterOpts) (*AWPRegistryUnpausedIterator, error) {

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &AWPRegistryUnpausedIterator{contract: _AWPRegistry.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_AWPRegistry *AWPRegistryFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AWPRegistryUnpaused) (event.Subscription, error) {

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryUnpaused)
				if err := _AWPRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_AWPRegistry *AWPRegistryFilterer) ParseUnpaused(log types.Log) (*AWPRegistryUnpaused, error) {
	event := new(AWPRegistryUnpaused)
	if err := _AWPRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPRegistryUserRegisteredIterator is returned from FilterUserRegistered and is used to iterate over the raw logs and unpacked data for UserRegistered events raised by the AWPRegistry contract.
type AWPRegistryUserRegisteredIterator struct {
	Event *AWPRegistryUserRegistered // Event containing the contract specifics and raw log

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
func (it *AWPRegistryUserRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPRegistryUserRegistered)
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
		it.Event = new(AWPRegistryUserRegistered)
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
func (it *AWPRegistryUserRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPRegistryUserRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPRegistryUserRegistered represents a UserRegistered event raised by the AWPRegistry contract.
type AWPRegistryUserRegistered struct {
	User common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterUserRegistered is a free log retrieval operation binding the contract event 0x54db7a5cb4735e1aac1f53db512d3390390bb6637bd30ad4bf9fc98667d9b9b9.
//
// Solidity: event UserRegistered(address indexed user)
func (_AWPRegistry *AWPRegistryFilterer) FilterUserRegistered(opts *bind.FilterOpts, user []common.Address) (*AWPRegistryUserRegisteredIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _AWPRegistry.contract.FilterLogs(opts, "UserRegistered", userRule)
	if err != nil {
		return nil, err
	}
	return &AWPRegistryUserRegisteredIterator{contract: _AWPRegistry.contract, event: "UserRegistered", logs: logs, sub: sub}, nil
}

// WatchUserRegistered is a free log subscription operation binding the contract event 0x54db7a5cb4735e1aac1f53db512d3390390bb6637bd30ad4bf9fc98667d9b9b9.
//
// Solidity: event UserRegistered(address indexed user)
func (_AWPRegistry *AWPRegistryFilterer) WatchUserRegistered(opts *bind.WatchOpts, sink chan<- *AWPRegistryUserRegistered, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _AWPRegistry.contract.WatchLogs(opts, "UserRegistered", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPRegistryUserRegistered)
				if err := _AWPRegistry.contract.UnpackLog(event, "UserRegistered", log); err != nil {
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
func (_AWPRegistry *AWPRegistryFilterer) ParseUserRegistered(log types.Log) (*AWPRegistryUserRegistered, error) {
	event := new(AWPRegistryUserRegistered)
	if err := _AWPRegistry.contract.UnpackLog(event, "UserRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
