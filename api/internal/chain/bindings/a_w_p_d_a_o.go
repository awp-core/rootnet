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

// AWPDAOMetaData contains all meta data concerning the AWPDAO contract.
var AWPDAOMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"stakeNFT_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"awpToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timelock_\",\"type\":\"address\",\"internalType\":\"contractTimelockController\"},{\"name\":\"votingDelay_\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"votingPeriod_\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumPercent_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"BALLOT_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"CLOCK_MODE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"COUNTING_MODE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"EXTENDED_BALLOT_TYPEHASH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"awpToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancel\",\"inputs\":[{\"name\":\"targets\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"values\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"calldatas\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"descriptionHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"castVote\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"castVoteBySig\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"castVoteWithReason\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"castVoteWithReasonAndParams\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"support\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"params\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"castVoteWithReasonAndParamsBySig\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"clock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"execute\",\"inputs\":[{\"name\":\"targets\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"values\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"calldatas\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"descriptionHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getVotes\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timepoint\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getVotesWithParams\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timepoint\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"params\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasVoted\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"hasVotedWithToken\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hashProposal\",\"inputs\":[{\"name\":\"targets\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"values\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"calldatas\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"descriptionHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"isSignalProposal\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nonces\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onERC1155BatchReceived\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onERC1155Received\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onERC721Received\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proposalCreatedAt\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposalDeadline\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposalEta\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposalNeedsQueuing\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposalProposer\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposalSnapshot\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposalThreshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposalTotalVotingPower\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proposalVotes\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"againstVotes\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"forVotes\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"abstainVotes\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"propose\",\"inputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"proposeWithTokens\",\"inputs\":[{\"name\":\"targets\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"values\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"calldatas\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"tokenIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"queue\",\"inputs\":[{\"name\":\"targets\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"values\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"calldatas\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"},{\"name\":\"descriptionHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"quorum\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quorumPercent\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"relay\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"setProposalThreshold\",\"inputs\":[{\"name\":\"newProposalThreshold\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setQuorumPercent\",\"inputs\":[{\"name\":\"newQuorumPercent\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setVotingDelay\",\"inputs\":[{\"name\":\"newVotingDelay\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setVotingPeriod\",\"inputs\":[{\"name\":\"newVotingPeriod\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"signalPropose\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"tokenIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakeNFT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIStakeNFT\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"state\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIGovernor.ProposalState\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"timelock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateTimelock\",\"inputs\":[{\"name\":\"newTimelock\",\"type\":\"address\",\"internalType\":\"contractTimelockController\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"version\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"votingDelay\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"votingPeriod\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProposalCanceled\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProposalCreated\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"proposer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"targets\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"values\",\"type\":\"uint256[]\",\"indexed\":false,\"internalType\":\"uint256[]\"},{\"name\":\"signatures\",\"type\":\"string[]\",\"indexed\":false,\"internalType\":\"string[]\"},{\"name\":\"calldatas\",\"type\":\"bytes[]\",\"indexed\":false,\"internalType\":\"bytes[]\"},{\"name\":\"voteStart\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"voteEnd\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"description\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProposalExecuted\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProposalQueued\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"etaSeconds\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ProposalThresholdSet\",\"inputs\":[{\"name\":\"oldProposalThreshold\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newProposalThreshold\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TimelockChange\",\"inputs\":[{\"name\":\"oldTimelock\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newTimelock\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VoteCast\",\"inputs\":[{\"name\":\"voter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"proposalId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"support\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"weight\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"reason\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VoteCastWithParams\",\"inputs\":[{\"name\":\"voter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"proposalId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"support\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"weight\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"reason\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"params\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VotingDelaySet\",\"inputs\":[{\"name\":\"oldVotingDelay\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newVotingDelay\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VotingPeriodSet\",\"inputs\":[{\"name\":\"oldVotingPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newVotingPeriod\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GovernorAlreadyCastVote\",\"inputs\":[{\"name\":\"voter\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"GovernorAlreadyQueuedProposal\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"GovernorDisabledDeposit\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GovernorInsufficientProposerVotes\",\"inputs\":[{\"name\":\"proposer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"votes\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"threshold\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"GovernorInvalidProposalLength\",\"inputs\":[{\"name\":\"targets\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"calldatas\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"values\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"GovernorInvalidSignature\",\"inputs\":[{\"name\":\"voter\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"GovernorInvalidVoteParams\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GovernorInvalidVoteType\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GovernorInvalidVotingPeriod\",\"inputs\":[{\"name\":\"votingPeriod\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"GovernorNonexistentProposal\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"GovernorNotQueuedProposal\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"GovernorOnlyExecutor\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"GovernorOnlyProposer\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"GovernorQueueNotImplemented\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GovernorRestrictedProposer\",\"inputs\":[{\"name\":\"proposer\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"GovernorUnexpectedProposalState\",\"inputs\":[{\"name\":\"proposalId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"current\",\"type\":\"uint8\",\"internalType\":\"enumIGovernor.ProposalState\"},{\"name\":\"expectedStates\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InsufficientVotingPower\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidAccountNonce\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"currentNonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidQuorumPercent\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidShortString\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LockExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MintedAfterProposal\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoTokens\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotTokenOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeCastOverflowedUintDowncast\",\"inputs\":[{\"name\":\"bits\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"StringTooLong\",\"inputs\":[{\"name\":\"str\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"TokenAlreadyVoted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UseProposeWithTokens\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UsecastVoteWithParams\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroTotalVotingPower\",\"inputs\":[]}]",
}

// AWPDAOABI is the input ABI used to generate the binding from.
// Deprecated: Use AWPDAOMetaData.ABI instead.
var AWPDAOABI = AWPDAOMetaData.ABI

// AWPDAO is an auto generated Go binding around an Ethereum contract.
type AWPDAO struct {
	AWPDAOCaller     // Read-only binding to the contract
	AWPDAOTransactor // Write-only binding to the contract
	AWPDAOFilterer   // Log filterer for contract events
}

// AWPDAOCaller is an auto generated read-only Go binding around an Ethereum contract.
type AWPDAOCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AWPDAOTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AWPDAOTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AWPDAOFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AWPDAOFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AWPDAOSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AWPDAOSession struct {
	Contract     *AWPDAO           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AWPDAOCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AWPDAOCallerSession struct {
	Contract *AWPDAOCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// AWPDAOTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AWPDAOTransactorSession struct {
	Contract     *AWPDAOTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AWPDAORaw is an auto generated low-level Go binding around an Ethereum contract.
type AWPDAORaw struct {
	Contract *AWPDAO // Generic contract binding to access the raw methods on
}

// AWPDAOCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AWPDAOCallerRaw struct {
	Contract *AWPDAOCaller // Generic read-only contract binding to access the raw methods on
}

// AWPDAOTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AWPDAOTransactorRaw struct {
	Contract *AWPDAOTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAWPDAO creates a new instance of AWPDAO, bound to a specific deployed contract.
func NewAWPDAO(address common.Address, backend bind.ContractBackend) (*AWPDAO, error) {
	contract, err := bindAWPDAO(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AWPDAO{AWPDAOCaller: AWPDAOCaller{contract: contract}, AWPDAOTransactor: AWPDAOTransactor{contract: contract}, AWPDAOFilterer: AWPDAOFilterer{contract: contract}}, nil
}

// NewAWPDAOCaller creates a new read-only instance of AWPDAO, bound to a specific deployed contract.
func NewAWPDAOCaller(address common.Address, caller bind.ContractCaller) (*AWPDAOCaller, error) {
	contract, err := bindAWPDAO(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AWPDAOCaller{contract: contract}, nil
}

// NewAWPDAOTransactor creates a new write-only instance of AWPDAO, bound to a specific deployed contract.
func NewAWPDAOTransactor(address common.Address, transactor bind.ContractTransactor) (*AWPDAOTransactor, error) {
	contract, err := bindAWPDAO(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AWPDAOTransactor{contract: contract}, nil
}

// NewAWPDAOFilterer creates a new log filterer instance of AWPDAO, bound to a specific deployed contract.
func NewAWPDAOFilterer(address common.Address, filterer bind.ContractFilterer) (*AWPDAOFilterer, error) {
	contract, err := bindAWPDAO(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AWPDAOFilterer{contract: contract}, nil
}

// bindAWPDAO binds a generic wrapper to an already deployed contract.
func bindAWPDAO(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AWPDAOMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AWPDAO *AWPDAORaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AWPDAO.Contract.AWPDAOCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AWPDAO *AWPDAORaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AWPDAO.Contract.AWPDAOTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AWPDAO *AWPDAORaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AWPDAO.Contract.AWPDAOTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AWPDAO *AWPDAOCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AWPDAO.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AWPDAO *AWPDAOTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AWPDAO.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AWPDAO *AWPDAOTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AWPDAO.Contract.contract.Transact(opts, method, params...)
}

// BALLOTTYPEHASH is a free data retrieval call binding the contract method 0xdeaaa7cc.
//
// Solidity: function BALLOT_TYPEHASH() view returns(bytes32)
func (_AWPDAO *AWPDAOCaller) BALLOTTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "BALLOT_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BALLOTTYPEHASH is a free data retrieval call binding the contract method 0xdeaaa7cc.
//
// Solidity: function BALLOT_TYPEHASH() view returns(bytes32)
func (_AWPDAO *AWPDAOSession) BALLOTTYPEHASH() ([32]byte, error) {
	return _AWPDAO.Contract.BALLOTTYPEHASH(&_AWPDAO.CallOpts)
}

// BALLOTTYPEHASH is a free data retrieval call binding the contract method 0xdeaaa7cc.
//
// Solidity: function BALLOT_TYPEHASH() view returns(bytes32)
func (_AWPDAO *AWPDAOCallerSession) BALLOTTYPEHASH() ([32]byte, error) {
	return _AWPDAO.Contract.BALLOTTYPEHASH(&_AWPDAO.CallOpts)
}

// CLOCKMODE is a free data retrieval call binding the contract method 0x4bf5d7e9.
//
// Solidity: function CLOCK_MODE() pure returns(string)
func (_AWPDAO *AWPDAOCaller) CLOCKMODE(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "CLOCK_MODE")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// CLOCKMODE is a free data retrieval call binding the contract method 0x4bf5d7e9.
//
// Solidity: function CLOCK_MODE() pure returns(string)
func (_AWPDAO *AWPDAOSession) CLOCKMODE() (string, error) {
	return _AWPDAO.Contract.CLOCKMODE(&_AWPDAO.CallOpts)
}

// CLOCKMODE is a free data retrieval call binding the contract method 0x4bf5d7e9.
//
// Solidity: function CLOCK_MODE() pure returns(string)
func (_AWPDAO *AWPDAOCallerSession) CLOCKMODE() (string, error) {
	return _AWPDAO.Contract.CLOCKMODE(&_AWPDAO.CallOpts)
}

// COUNTINGMODE is a free data retrieval call binding the contract method 0xdd4e2ba5.
//
// Solidity: function COUNTING_MODE() pure returns(string)
func (_AWPDAO *AWPDAOCaller) COUNTINGMODE(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "COUNTING_MODE")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// COUNTINGMODE is a free data retrieval call binding the contract method 0xdd4e2ba5.
//
// Solidity: function COUNTING_MODE() pure returns(string)
func (_AWPDAO *AWPDAOSession) COUNTINGMODE() (string, error) {
	return _AWPDAO.Contract.COUNTINGMODE(&_AWPDAO.CallOpts)
}

// COUNTINGMODE is a free data retrieval call binding the contract method 0xdd4e2ba5.
//
// Solidity: function COUNTING_MODE() pure returns(string)
func (_AWPDAO *AWPDAOCallerSession) COUNTINGMODE() (string, error) {
	return _AWPDAO.Contract.COUNTINGMODE(&_AWPDAO.CallOpts)
}

// EXTENDEDBALLOTTYPEHASH is a free data retrieval call binding the contract method 0x2fe3e261.
//
// Solidity: function EXTENDED_BALLOT_TYPEHASH() view returns(bytes32)
func (_AWPDAO *AWPDAOCaller) EXTENDEDBALLOTTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "EXTENDED_BALLOT_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EXTENDEDBALLOTTYPEHASH is a free data retrieval call binding the contract method 0x2fe3e261.
//
// Solidity: function EXTENDED_BALLOT_TYPEHASH() view returns(bytes32)
func (_AWPDAO *AWPDAOSession) EXTENDEDBALLOTTYPEHASH() ([32]byte, error) {
	return _AWPDAO.Contract.EXTENDEDBALLOTTYPEHASH(&_AWPDAO.CallOpts)
}

// EXTENDEDBALLOTTYPEHASH is a free data retrieval call binding the contract method 0x2fe3e261.
//
// Solidity: function EXTENDED_BALLOT_TYPEHASH() view returns(bytes32)
func (_AWPDAO *AWPDAOCallerSession) EXTENDEDBALLOTTYPEHASH() ([32]byte, error) {
	return _AWPDAO.Contract.EXTENDEDBALLOTTYPEHASH(&_AWPDAO.CallOpts)
}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_AWPDAO *AWPDAOCaller) AwpToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "awpToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_AWPDAO *AWPDAOSession) AwpToken() (common.Address, error) {
	return _AWPDAO.Contract.AwpToken(&_AWPDAO.CallOpts)
}

// AwpToken is a free data retrieval call binding the contract method 0x41a578cf.
//
// Solidity: function awpToken() view returns(address)
func (_AWPDAO *AWPDAOCallerSession) AwpToken() (common.Address, error) {
	return _AWPDAO.Contract.AwpToken(&_AWPDAO.CallOpts)
}

// CastVote is a free data retrieval call binding the contract method 0x56781388.
//
// Solidity: function castVote(uint256 , uint8 ) pure returns(uint256)
func (_AWPDAO *AWPDAOCaller) CastVote(opts *bind.CallOpts, arg0 *big.Int, arg1 uint8) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "castVote", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CastVote is a free data retrieval call binding the contract method 0x56781388.
//
// Solidity: function castVote(uint256 , uint8 ) pure returns(uint256)
func (_AWPDAO *AWPDAOSession) CastVote(arg0 *big.Int, arg1 uint8) (*big.Int, error) {
	return _AWPDAO.Contract.CastVote(&_AWPDAO.CallOpts, arg0, arg1)
}

// CastVote is a free data retrieval call binding the contract method 0x56781388.
//
// Solidity: function castVote(uint256 , uint8 ) pure returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) CastVote(arg0 *big.Int, arg1 uint8) (*big.Int, error) {
	return _AWPDAO.Contract.CastVote(&_AWPDAO.CallOpts, arg0, arg1)
}

// CastVoteBySig is a free data retrieval call binding the contract method 0x8ff262e3.
//
// Solidity: function castVoteBySig(uint256 , uint8 , address , bytes ) pure returns(uint256)
func (_AWPDAO *AWPDAOCaller) CastVoteBySig(opts *bind.CallOpts, arg0 *big.Int, arg1 uint8, arg2 common.Address, arg3 []byte) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "castVoteBySig", arg0, arg1, arg2, arg3)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CastVoteBySig is a free data retrieval call binding the contract method 0x8ff262e3.
//
// Solidity: function castVoteBySig(uint256 , uint8 , address , bytes ) pure returns(uint256)
func (_AWPDAO *AWPDAOSession) CastVoteBySig(arg0 *big.Int, arg1 uint8, arg2 common.Address, arg3 []byte) (*big.Int, error) {
	return _AWPDAO.Contract.CastVoteBySig(&_AWPDAO.CallOpts, arg0, arg1, arg2, arg3)
}

// CastVoteBySig is a free data retrieval call binding the contract method 0x8ff262e3.
//
// Solidity: function castVoteBySig(uint256 , uint8 , address , bytes ) pure returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) CastVoteBySig(arg0 *big.Int, arg1 uint8, arg2 common.Address, arg3 []byte) (*big.Int, error) {
	return _AWPDAO.Contract.CastVoteBySig(&_AWPDAO.CallOpts, arg0, arg1, arg2, arg3)
}

// CastVoteWithReason is a free data retrieval call binding the contract method 0x7b3c71d3.
//
// Solidity: function castVoteWithReason(uint256 , uint8 , string ) pure returns(uint256)
func (_AWPDAO *AWPDAOCaller) CastVoteWithReason(opts *bind.CallOpts, arg0 *big.Int, arg1 uint8, arg2 string) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "castVoteWithReason", arg0, arg1, arg2)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CastVoteWithReason is a free data retrieval call binding the contract method 0x7b3c71d3.
//
// Solidity: function castVoteWithReason(uint256 , uint8 , string ) pure returns(uint256)
func (_AWPDAO *AWPDAOSession) CastVoteWithReason(arg0 *big.Int, arg1 uint8, arg2 string) (*big.Int, error) {
	return _AWPDAO.Contract.CastVoteWithReason(&_AWPDAO.CallOpts, arg0, arg1, arg2)
}

// CastVoteWithReason is a free data retrieval call binding the contract method 0x7b3c71d3.
//
// Solidity: function castVoteWithReason(uint256 , uint8 , string ) pure returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) CastVoteWithReason(arg0 *big.Int, arg1 uint8, arg2 string) (*big.Int, error) {
	return _AWPDAO.Contract.CastVoteWithReason(&_AWPDAO.CallOpts, arg0, arg1, arg2)
}

// CastVoteWithReasonAndParamsBySig is a free data retrieval call binding the contract method 0x5b8d0e0d.
//
// Solidity: function castVoteWithReasonAndParamsBySig(uint256 , uint8 , address , string , bytes , bytes ) pure returns(uint256)
func (_AWPDAO *AWPDAOCaller) CastVoteWithReasonAndParamsBySig(opts *bind.CallOpts, arg0 *big.Int, arg1 uint8, arg2 common.Address, arg3 string, arg4 []byte, arg5 []byte) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "castVoteWithReasonAndParamsBySig", arg0, arg1, arg2, arg3, arg4, arg5)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CastVoteWithReasonAndParamsBySig is a free data retrieval call binding the contract method 0x5b8d0e0d.
//
// Solidity: function castVoteWithReasonAndParamsBySig(uint256 , uint8 , address , string , bytes , bytes ) pure returns(uint256)
func (_AWPDAO *AWPDAOSession) CastVoteWithReasonAndParamsBySig(arg0 *big.Int, arg1 uint8, arg2 common.Address, arg3 string, arg4 []byte, arg5 []byte) (*big.Int, error) {
	return _AWPDAO.Contract.CastVoteWithReasonAndParamsBySig(&_AWPDAO.CallOpts, arg0, arg1, arg2, arg3, arg4, arg5)
}

// CastVoteWithReasonAndParamsBySig is a free data retrieval call binding the contract method 0x5b8d0e0d.
//
// Solidity: function castVoteWithReasonAndParamsBySig(uint256 , uint8 , address , string , bytes , bytes ) pure returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) CastVoteWithReasonAndParamsBySig(arg0 *big.Int, arg1 uint8, arg2 common.Address, arg3 string, arg4 []byte, arg5 []byte) (*big.Int, error) {
	return _AWPDAO.Contract.CastVoteWithReasonAndParamsBySig(&_AWPDAO.CallOpts, arg0, arg1, arg2, arg3, arg4, arg5)
}

// Clock is a free data retrieval call binding the contract method 0x91ddadf4.
//
// Solidity: function clock() view returns(uint48)
func (_AWPDAO *AWPDAOCaller) Clock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "clock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Clock is a free data retrieval call binding the contract method 0x91ddadf4.
//
// Solidity: function clock() view returns(uint48)
func (_AWPDAO *AWPDAOSession) Clock() (*big.Int, error) {
	return _AWPDAO.Contract.Clock(&_AWPDAO.CallOpts)
}

// Clock is a free data retrieval call binding the contract method 0x91ddadf4.
//
// Solidity: function clock() view returns(uint48)
func (_AWPDAO *AWPDAOCallerSession) Clock() (*big.Int, error) {
	return _AWPDAO.Contract.Clock(&_AWPDAO.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_AWPDAO *AWPDAOCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "eip712Domain")

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
func (_AWPDAO *AWPDAOSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _AWPDAO.Contract.Eip712Domain(&_AWPDAO.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_AWPDAO *AWPDAOCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _AWPDAO.Contract.Eip712Domain(&_AWPDAO.CallOpts)
}

// GetVotes is a free data retrieval call binding the contract method 0xeb9019d4.
//
// Solidity: function getVotes(address account, uint256 timepoint) view returns(uint256)
func (_AWPDAO *AWPDAOCaller) GetVotes(opts *bind.CallOpts, account common.Address, timepoint *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "getVotes", account, timepoint)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVotes is a free data retrieval call binding the contract method 0xeb9019d4.
//
// Solidity: function getVotes(address account, uint256 timepoint) view returns(uint256)
func (_AWPDAO *AWPDAOSession) GetVotes(account common.Address, timepoint *big.Int) (*big.Int, error) {
	return _AWPDAO.Contract.GetVotes(&_AWPDAO.CallOpts, account, timepoint)
}

// GetVotes is a free data retrieval call binding the contract method 0xeb9019d4.
//
// Solidity: function getVotes(address account, uint256 timepoint) view returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) GetVotes(account common.Address, timepoint *big.Int) (*big.Int, error) {
	return _AWPDAO.Contract.GetVotes(&_AWPDAO.CallOpts, account, timepoint)
}

// GetVotesWithParams is a free data retrieval call binding the contract method 0x9a802a6d.
//
// Solidity: function getVotesWithParams(address account, uint256 timepoint, bytes params) view returns(uint256)
func (_AWPDAO *AWPDAOCaller) GetVotesWithParams(opts *bind.CallOpts, account common.Address, timepoint *big.Int, params []byte) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "getVotesWithParams", account, timepoint, params)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVotesWithParams is a free data retrieval call binding the contract method 0x9a802a6d.
//
// Solidity: function getVotesWithParams(address account, uint256 timepoint, bytes params) view returns(uint256)
func (_AWPDAO *AWPDAOSession) GetVotesWithParams(account common.Address, timepoint *big.Int, params []byte) (*big.Int, error) {
	return _AWPDAO.Contract.GetVotesWithParams(&_AWPDAO.CallOpts, account, timepoint, params)
}

// GetVotesWithParams is a free data retrieval call binding the contract method 0x9a802a6d.
//
// Solidity: function getVotesWithParams(address account, uint256 timepoint, bytes params) view returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) GetVotesWithParams(account common.Address, timepoint *big.Int, params []byte) (*big.Int, error) {
	return _AWPDAO.Contract.GetVotesWithParams(&_AWPDAO.CallOpts, account, timepoint, params)
}

// HasVoted is a free data retrieval call binding the contract method 0x43859632.
//
// Solidity: function hasVoted(uint256 , address ) pure returns(bool)
func (_AWPDAO *AWPDAOCaller) HasVoted(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "hasVoted", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasVoted is a free data retrieval call binding the contract method 0x43859632.
//
// Solidity: function hasVoted(uint256 , address ) pure returns(bool)
func (_AWPDAO *AWPDAOSession) HasVoted(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _AWPDAO.Contract.HasVoted(&_AWPDAO.CallOpts, arg0, arg1)
}

// HasVoted is a free data retrieval call binding the contract method 0x43859632.
//
// Solidity: function hasVoted(uint256 , address ) pure returns(bool)
func (_AWPDAO *AWPDAOCallerSession) HasVoted(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _AWPDAO.Contract.HasVoted(&_AWPDAO.CallOpts, arg0, arg1)
}

// HasVotedWithToken is a free data retrieval call binding the contract method 0x014f0e84.
//
// Solidity: function hasVotedWithToken(uint256 proposalId, uint256 tokenId) view returns(bool)
func (_AWPDAO *AWPDAOCaller) HasVotedWithToken(opts *bind.CallOpts, proposalId *big.Int, tokenId *big.Int) (bool, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "hasVotedWithToken", proposalId, tokenId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasVotedWithToken is a free data retrieval call binding the contract method 0x014f0e84.
//
// Solidity: function hasVotedWithToken(uint256 proposalId, uint256 tokenId) view returns(bool)
func (_AWPDAO *AWPDAOSession) HasVotedWithToken(proposalId *big.Int, tokenId *big.Int) (bool, error) {
	return _AWPDAO.Contract.HasVotedWithToken(&_AWPDAO.CallOpts, proposalId, tokenId)
}

// HasVotedWithToken is a free data retrieval call binding the contract method 0x014f0e84.
//
// Solidity: function hasVotedWithToken(uint256 proposalId, uint256 tokenId) view returns(bool)
func (_AWPDAO *AWPDAOCallerSession) HasVotedWithToken(proposalId *big.Int, tokenId *big.Int) (bool, error) {
	return _AWPDAO.Contract.HasVotedWithToken(&_AWPDAO.CallOpts, proposalId, tokenId)
}

// HashProposal is a free data retrieval call binding the contract method 0xc59057e4.
//
// Solidity: function hashProposal(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) pure returns(uint256)
func (_AWPDAO *AWPDAOCaller) HashProposal(opts *bind.CallOpts, targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "hashProposal", targets, values, calldatas, descriptionHash)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// HashProposal is a free data retrieval call binding the contract method 0xc59057e4.
//
// Solidity: function hashProposal(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) pure returns(uint256)
func (_AWPDAO *AWPDAOSession) HashProposal(targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*big.Int, error) {
	return _AWPDAO.Contract.HashProposal(&_AWPDAO.CallOpts, targets, values, calldatas, descriptionHash)
}

// HashProposal is a free data retrieval call binding the contract method 0xc59057e4.
//
// Solidity: function hashProposal(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) pure returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) HashProposal(targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*big.Int, error) {
	return _AWPDAO.Contract.HashProposal(&_AWPDAO.CallOpts, targets, values, calldatas, descriptionHash)
}

// IsSignalProposal is a free data retrieval call binding the contract method 0x4a5fa5be.
//
// Solidity: function isSignalProposal(uint256 proposalId) view returns(bool)
func (_AWPDAO *AWPDAOCaller) IsSignalProposal(opts *bind.CallOpts, proposalId *big.Int) (bool, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "isSignalProposal", proposalId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsSignalProposal is a free data retrieval call binding the contract method 0x4a5fa5be.
//
// Solidity: function isSignalProposal(uint256 proposalId) view returns(bool)
func (_AWPDAO *AWPDAOSession) IsSignalProposal(proposalId *big.Int) (bool, error) {
	return _AWPDAO.Contract.IsSignalProposal(&_AWPDAO.CallOpts, proposalId)
}

// IsSignalProposal is a free data retrieval call binding the contract method 0x4a5fa5be.
//
// Solidity: function isSignalProposal(uint256 proposalId) view returns(bool)
func (_AWPDAO *AWPDAOCallerSession) IsSignalProposal(proposalId *big.Int) (bool, error) {
	return _AWPDAO.Contract.IsSignalProposal(&_AWPDAO.CallOpts, proposalId)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AWPDAO *AWPDAOCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AWPDAO *AWPDAOSession) Name() (string, error) {
	return _AWPDAO.Contract.Name(&_AWPDAO.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AWPDAO *AWPDAOCallerSession) Name() (string, error) {
	return _AWPDAO.Contract.Name(&_AWPDAO.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_AWPDAO *AWPDAOCaller) Nonces(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "nonces", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_AWPDAO *AWPDAOSession) Nonces(owner common.Address) (*big.Int, error) {
	return _AWPDAO.Contract.Nonces(&_AWPDAO.CallOpts, owner)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) Nonces(owner common.Address) (*big.Int, error) {
	return _AWPDAO.Contract.Nonces(&_AWPDAO.CallOpts, owner)
}

// ProposalCreatedAt is a free data retrieval call binding the contract method 0x5f9103b2.
//
// Solidity: function proposalCreatedAt(uint256 proposalId) view returns(uint256)
func (_AWPDAO *AWPDAOCaller) ProposalCreatedAt(opts *bind.CallOpts, proposalId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "proposalCreatedAt", proposalId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProposalCreatedAt is a free data retrieval call binding the contract method 0x5f9103b2.
//
// Solidity: function proposalCreatedAt(uint256 proposalId) view returns(uint256)
func (_AWPDAO *AWPDAOSession) ProposalCreatedAt(proposalId *big.Int) (*big.Int, error) {
	return _AWPDAO.Contract.ProposalCreatedAt(&_AWPDAO.CallOpts, proposalId)
}

// ProposalCreatedAt is a free data retrieval call binding the contract method 0x5f9103b2.
//
// Solidity: function proposalCreatedAt(uint256 proposalId) view returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) ProposalCreatedAt(proposalId *big.Int) (*big.Int, error) {
	return _AWPDAO.Contract.ProposalCreatedAt(&_AWPDAO.CallOpts, proposalId)
}

// ProposalDeadline is a free data retrieval call binding the contract method 0xc01f9e37.
//
// Solidity: function proposalDeadline(uint256 proposalId) view returns(uint256)
func (_AWPDAO *AWPDAOCaller) ProposalDeadline(opts *bind.CallOpts, proposalId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "proposalDeadline", proposalId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProposalDeadline is a free data retrieval call binding the contract method 0xc01f9e37.
//
// Solidity: function proposalDeadline(uint256 proposalId) view returns(uint256)
func (_AWPDAO *AWPDAOSession) ProposalDeadline(proposalId *big.Int) (*big.Int, error) {
	return _AWPDAO.Contract.ProposalDeadline(&_AWPDAO.CallOpts, proposalId)
}

// ProposalDeadline is a free data retrieval call binding the contract method 0xc01f9e37.
//
// Solidity: function proposalDeadline(uint256 proposalId) view returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) ProposalDeadline(proposalId *big.Int) (*big.Int, error) {
	return _AWPDAO.Contract.ProposalDeadline(&_AWPDAO.CallOpts, proposalId)
}

// ProposalEta is a free data retrieval call binding the contract method 0xab58fb8e.
//
// Solidity: function proposalEta(uint256 proposalId) view returns(uint256)
func (_AWPDAO *AWPDAOCaller) ProposalEta(opts *bind.CallOpts, proposalId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "proposalEta", proposalId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProposalEta is a free data retrieval call binding the contract method 0xab58fb8e.
//
// Solidity: function proposalEta(uint256 proposalId) view returns(uint256)
func (_AWPDAO *AWPDAOSession) ProposalEta(proposalId *big.Int) (*big.Int, error) {
	return _AWPDAO.Contract.ProposalEta(&_AWPDAO.CallOpts, proposalId)
}

// ProposalEta is a free data retrieval call binding the contract method 0xab58fb8e.
//
// Solidity: function proposalEta(uint256 proposalId) view returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) ProposalEta(proposalId *big.Int) (*big.Int, error) {
	return _AWPDAO.Contract.ProposalEta(&_AWPDAO.CallOpts, proposalId)
}

// ProposalNeedsQueuing is a free data retrieval call binding the contract method 0xa9a95294.
//
// Solidity: function proposalNeedsQueuing(uint256 proposalId) view returns(bool)
func (_AWPDAO *AWPDAOCaller) ProposalNeedsQueuing(opts *bind.CallOpts, proposalId *big.Int) (bool, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "proposalNeedsQueuing", proposalId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ProposalNeedsQueuing is a free data retrieval call binding the contract method 0xa9a95294.
//
// Solidity: function proposalNeedsQueuing(uint256 proposalId) view returns(bool)
func (_AWPDAO *AWPDAOSession) ProposalNeedsQueuing(proposalId *big.Int) (bool, error) {
	return _AWPDAO.Contract.ProposalNeedsQueuing(&_AWPDAO.CallOpts, proposalId)
}

// ProposalNeedsQueuing is a free data retrieval call binding the contract method 0xa9a95294.
//
// Solidity: function proposalNeedsQueuing(uint256 proposalId) view returns(bool)
func (_AWPDAO *AWPDAOCallerSession) ProposalNeedsQueuing(proposalId *big.Int) (bool, error) {
	return _AWPDAO.Contract.ProposalNeedsQueuing(&_AWPDAO.CallOpts, proposalId)
}

// ProposalProposer is a free data retrieval call binding the contract method 0x143489d0.
//
// Solidity: function proposalProposer(uint256 proposalId) view returns(address)
func (_AWPDAO *AWPDAOCaller) ProposalProposer(opts *bind.CallOpts, proposalId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "proposalProposer", proposalId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProposalProposer is a free data retrieval call binding the contract method 0x143489d0.
//
// Solidity: function proposalProposer(uint256 proposalId) view returns(address)
func (_AWPDAO *AWPDAOSession) ProposalProposer(proposalId *big.Int) (common.Address, error) {
	return _AWPDAO.Contract.ProposalProposer(&_AWPDAO.CallOpts, proposalId)
}

// ProposalProposer is a free data retrieval call binding the contract method 0x143489d0.
//
// Solidity: function proposalProposer(uint256 proposalId) view returns(address)
func (_AWPDAO *AWPDAOCallerSession) ProposalProposer(proposalId *big.Int) (common.Address, error) {
	return _AWPDAO.Contract.ProposalProposer(&_AWPDAO.CallOpts, proposalId)
}

// ProposalSnapshot is a free data retrieval call binding the contract method 0x2d63f693.
//
// Solidity: function proposalSnapshot(uint256 proposalId) view returns(uint256)
func (_AWPDAO *AWPDAOCaller) ProposalSnapshot(opts *bind.CallOpts, proposalId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "proposalSnapshot", proposalId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProposalSnapshot is a free data retrieval call binding the contract method 0x2d63f693.
//
// Solidity: function proposalSnapshot(uint256 proposalId) view returns(uint256)
func (_AWPDAO *AWPDAOSession) ProposalSnapshot(proposalId *big.Int) (*big.Int, error) {
	return _AWPDAO.Contract.ProposalSnapshot(&_AWPDAO.CallOpts, proposalId)
}

// ProposalSnapshot is a free data retrieval call binding the contract method 0x2d63f693.
//
// Solidity: function proposalSnapshot(uint256 proposalId) view returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) ProposalSnapshot(proposalId *big.Int) (*big.Int, error) {
	return _AWPDAO.Contract.ProposalSnapshot(&_AWPDAO.CallOpts, proposalId)
}

// ProposalThreshold is a free data retrieval call binding the contract method 0xb58131b0.
//
// Solidity: function proposalThreshold() view returns(uint256)
func (_AWPDAO *AWPDAOCaller) ProposalThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "proposalThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProposalThreshold is a free data retrieval call binding the contract method 0xb58131b0.
//
// Solidity: function proposalThreshold() view returns(uint256)
func (_AWPDAO *AWPDAOSession) ProposalThreshold() (*big.Int, error) {
	return _AWPDAO.Contract.ProposalThreshold(&_AWPDAO.CallOpts)
}

// ProposalThreshold is a free data retrieval call binding the contract method 0xb58131b0.
//
// Solidity: function proposalThreshold() view returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) ProposalThreshold() (*big.Int, error) {
	return _AWPDAO.Contract.ProposalThreshold(&_AWPDAO.CallOpts)
}

// ProposalTotalVotingPower is a free data retrieval call binding the contract method 0xd6c7d925.
//
// Solidity: function proposalTotalVotingPower(uint256 proposalId) view returns(uint256)
func (_AWPDAO *AWPDAOCaller) ProposalTotalVotingPower(opts *bind.CallOpts, proposalId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "proposalTotalVotingPower", proposalId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProposalTotalVotingPower is a free data retrieval call binding the contract method 0xd6c7d925.
//
// Solidity: function proposalTotalVotingPower(uint256 proposalId) view returns(uint256)
func (_AWPDAO *AWPDAOSession) ProposalTotalVotingPower(proposalId *big.Int) (*big.Int, error) {
	return _AWPDAO.Contract.ProposalTotalVotingPower(&_AWPDAO.CallOpts, proposalId)
}

// ProposalTotalVotingPower is a free data retrieval call binding the contract method 0xd6c7d925.
//
// Solidity: function proposalTotalVotingPower(uint256 proposalId) view returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) ProposalTotalVotingPower(proposalId *big.Int) (*big.Int, error) {
	return _AWPDAO.Contract.ProposalTotalVotingPower(&_AWPDAO.CallOpts, proposalId)
}

// ProposalVotes is a free data retrieval call binding the contract method 0x544ffc9c.
//
// Solidity: function proposalVotes(uint256 proposalId) view returns(uint256 againstVotes, uint256 forVotes, uint256 abstainVotes)
func (_AWPDAO *AWPDAOCaller) ProposalVotes(opts *bind.CallOpts, proposalId *big.Int) (struct {
	AgainstVotes *big.Int
	ForVotes     *big.Int
	AbstainVotes *big.Int
}, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "proposalVotes", proposalId)

	outstruct := new(struct {
		AgainstVotes *big.Int
		ForVotes     *big.Int
		AbstainVotes *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AgainstVotes = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.ForVotes = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.AbstainVotes = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ProposalVotes is a free data retrieval call binding the contract method 0x544ffc9c.
//
// Solidity: function proposalVotes(uint256 proposalId) view returns(uint256 againstVotes, uint256 forVotes, uint256 abstainVotes)
func (_AWPDAO *AWPDAOSession) ProposalVotes(proposalId *big.Int) (struct {
	AgainstVotes *big.Int
	ForVotes     *big.Int
	AbstainVotes *big.Int
}, error) {
	return _AWPDAO.Contract.ProposalVotes(&_AWPDAO.CallOpts, proposalId)
}

// ProposalVotes is a free data retrieval call binding the contract method 0x544ffc9c.
//
// Solidity: function proposalVotes(uint256 proposalId) view returns(uint256 againstVotes, uint256 forVotes, uint256 abstainVotes)
func (_AWPDAO *AWPDAOCallerSession) ProposalVotes(proposalId *big.Int) (struct {
	AgainstVotes *big.Int
	ForVotes     *big.Int
	AbstainVotes *big.Int
}, error) {
	return _AWPDAO.Contract.ProposalVotes(&_AWPDAO.CallOpts, proposalId)
}

// Propose is a free data retrieval call binding the contract method 0x7d5e81e2.
//
// Solidity: function propose(address[] , uint256[] , bytes[] , string ) pure returns(uint256)
func (_AWPDAO *AWPDAOCaller) Propose(opts *bind.CallOpts, arg0 []common.Address, arg1 []*big.Int, arg2 [][]byte, arg3 string) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "propose", arg0, arg1, arg2, arg3)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Propose is a free data retrieval call binding the contract method 0x7d5e81e2.
//
// Solidity: function propose(address[] , uint256[] , bytes[] , string ) pure returns(uint256)
func (_AWPDAO *AWPDAOSession) Propose(arg0 []common.Address, arg1 []*big.Int, arg2 [][]byte, arg3 string) (*big.Int, error) {
	return _AWPDAO.Contract.Propose(&_AWPDAO.CallOpts, arg0, arg1, arg2, arg3)
}

// Propose is a free data retrieval call binding the contract method 0x7d5e81e2.
//
// Solidity: function propose(address[] , uint256[] , bytes[] , string ) pure returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) Propose(arg0 []common.Address, arg1 []*big.Int, arg2 [][]byte, arg3 string) (*big.Int, error) {
	return _AWPDAO.Contract.Propose(&_AWPDAO.CallOpts, arg0, arg1, arg2, arg3)
}

// Quorum is a free data retrieval call binding the contract method 0xf8ce560a.
//
// Solidity: function quorum(uint256 ) view returns(uint256)
func (_AWPDAO *AWPDAOCaller) Quorum(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "quorum", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Quorum is a free data retrieval call binding the contract method 0xf8ce560a.
//
// Solidity: function quorum(uint256 ) view returns(uint256)
func (_AWPDAO *AWPDAOSession) Quorum(arg0 *big.Int) (*big.Int, error) {
	return _AWPDAO.Contract.Quorum(&_AWPDAO.CallOpts, arg0)
}

// Quorum is a free data retrieval call binding the contract method 0xf8ce560a.
//
// Solidity: function quorum(uint256 ) view returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) Quorum(arg0 *big.Int) (*big.Int, error) {
	return _AWPDAO.Contract.Quorum(&_AWPDAO.CallOpts, arg0)
}

// QuorumPercent is a free data retrieval call binding the contract method 0xf81cbd26.
//
// Solidity: function quorumPercent() view returns(uint256)
func (_AWPDAO *AWPDAOCaller) QuorumPercent(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "quorumPercent")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// QuorumPercent is a free data retrieval call binding the contract method 0xf81cbd26.
//
// Solidity: function quorumPercent() view returns(uint256)
func (_AWPDAO *AWPDAOSession) QuorumPercent() (*big.Int, error) {
	return _AWPDAO.Contract.QuorumPercent(&_AWPDAO.CallOpts)
}

// QuorumPercent is a free data retrieval call binding the contract method 0xf81cbd26.
//
// Solidity: function quorumPercent() view returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) QuorumPercent() (*big.Int, error) {
	return _AWPDAO.Contract.QuorumPercent(&_AWPDAO.CallOpts)
}

// StakeNFT is a free data retrieval call binding the contract method 0xb48509e6.
//
// Solidity: function stakeNFT() view returns(address)
func (_AWPDAO *AWPDAOCaller) StakeNFT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "stakeNFT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakeNFT is a free data retrieval call binding the contract method 0xb48509e6.
//
// Solidity: function stakeNFT() view returns(address)
func (_AWPDAO *AWPDAOSession) StakeNFT() (common.Address, error) {
	return _AWPDAO.Contract.StakeNFT(&_AWPDAO.CallOpts)
}

// StakeNFT is a free data retrieval call binding the contract method 0xb48509e6.
//
// Solidity: function stakeNFT() view returns(address)
func (_AWPDAO *AWPDAOCallerSession) StakeNFT() (common.Address, error) {
	return _AWPDAO.Contract.StakeNFT(&_AWPDAO.CallOpts)
}

// State is a free data retrieval call binding the contract method 0x3e4f49e6.
//
// Solidity: function state(uint256 proposalId) view returns(uint8)
func (_AWPDAO *AWPDAOCaller) State(opts *bind.CallOpts, proposalId *big.Int) (uint8, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "state", proposalId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// State is a free data retrieval call binding the contract method 0x3e4f49e6.
//
// Solidity: function state(uint256 proposalId) view returns(uint8)
func (_AWPDAO *AWPDAOSession) State(proposalId *big.Int) (uint8, error) {
	return _AWPDAO.Contract.State(&_AWPDAO.CallOpts, proposalId)
}

// State is a free data retrieval call binding the contract method 0x3e4f49e6.
//
// Solidity: function state(uint256 proposalId) view returns(uint8)
func (_AWPDAO *AWPDAOCallerSession) State(proposalId *big.Int) (uint8, error) {
	return _AWPDAO.Contract.State(&_AWPDAO.CallOpts, proposalId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AWPDAO *AWPDAOCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AWPDAO *AWPDAOSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AWPDAO.Contract.SupportsInterface(&_AWPDAO.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AWPDAO *AWPDAOCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AWPDAO.Contract.SupportsInterface(&_AWPDAO.CallOpts, interfaceId)
}

// Timelock is a free data retrieval call binding the contract method 0xd33219b4.
//
// Solidity: function timelock() view returns(address)
func (_AWPDAO *AWPDAOCaller) Timelock(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "timelock")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Timelock is a free data retrieval call binding the contract method 0xd33219b4.
//
// Solidity: function timelock() view returns(address)
func (_AWPDAO *AWPDAOSession) Timelock() (common.Address, error) {
	return _AWPDAO.Contract.Timelock(&_AWPDAO.CallOpts)
}

// Timelock is a free data retrieval call binding the contract method 0xd33219b4.
//
// Solidity: function timelock() view returns(address)
func (_AWPDAO *AWPDAOCallerSession) Timelock() (common.Address, error) {
	return _AWPDAO.Contract.Timelock(&_AWPDAO.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_AWPDAO *AWPDAOCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_AWPDAO *AWPDAOSession) Version() (string, error) {
	return _AWPDAO.Contract.Version(&_AWPDAO.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_AWPDAO *AWPDAOCallerSession) Version() (string, error) {
	return _AWPDAO.Contract.Version(&_AWPDAO.CallOpts)
}

// VotingDelay is a free data retrieval call binding the contract method 0x3932abb1.
//
// Solidity: function votingDelay() view returns(uint256)
func (_AWPDAO *AWPDAOCaller) VotingDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "votingDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VotingDelay is a free data retrieval call binding the contract method 0x3932abb1.
//
// Solidity: function votingDelay() view returns(uint256)
func (_AWPDAO *AWPDAOSession) VotingDelay() (*big.Int, error) {
	return _AWPDAO.Contract.VotingDelay(&_AWPDAO.CallOpts)
}

// VotingDelay is a free data retrieval call binding the contract method 0x3932abb1.
//
// Solidity: function votingDelay() view returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) VotingDelay() (*big.Int, error) {
	return _AWPDAO.Contract.VotingDelay(&_AWPDAO.CallOpts)
}

// VotingPeriod is a free data retrieval call binding the contract method 0x02a251a3.
//
// Solidity: function votingPeriod() view returns(uint256)
func (_AWPDAO *AWPDAOCaller) VotingPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AWPDAO.contract.Call(opts, &out, "votingPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VotingPeriod is a free data retrieval call binding the contract method 0x02a251a3.
//
// Solidity: function votingPeriod() view returns(uint256)
func (_AWPDAO *AWPDAOSession) VotingPeriod() (*big.Int, error) {
	return _AWPDAO.Contract.VotingPeriod(&_AWPDAO.CallOpts)
}

// VotingPeriod is a free data retrieval call binding the contract method 0x02a251a3.
//
// Solidity: function votingPeriod() view returns(uint256)
func (_AWPDAO *AWPDAOCallerSession) VotingPeriod() (*big.Int, error) {
	return _AWPDAO.Contract.VotingPeriod(&_AWPDAO.CallOpts)
}

// Cancel is a paid mutator transaction binding the contract method 0x452115d6.
//
// Solidity: function cancel(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) returns(uint256)
func (_AWPDAO *AWPDAOTransactor) Cancel(opts *bind.TransactOpts, targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*types.Transaction, error) {
	return _AWPDAO.contract.Transact(opts, "cancel", targets, values, calldatas, descriptionHash)
}

// Cancel is a paid mutator transaction binding the contract method 0x452115d6.
//
// Solidity: function cancel(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) returns(uint256)
func (_AWPDAO *AWPDAOSession) Cancel(targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*types.Transaction, error) {
	return _AWPDAO.Contract.Cancel(&_AWPDAO.TransactOpts, targets, values, calldatas, descriptionHash)
}

// Cancel is a paid mutator transaction binding the contract method 0x452115d6.
//
// Solidity: function cancel(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) returns(uint256)
func (_AWPDAO *AWPDAOTransactorSession) Cancel(targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*types.Transaction, error) {
	return _AWPDAO.Contract.Cancel(&_AWPDAO.TransactOpts, targets, values, calldatas, descriptionHash)
}

// CastVoteWithReasonAndParams is a paid mutator transaction binding the contract method 0x5f398a14.
//
// Solidity: function castVoteWithReasonAndParams(uint256 proposalId, uint8 support, string reason, bytes params) returns(uint256)
func (_AWPDAO *AWPDAOTransactor) CastVoteWithReasonAndParams(opts *bind.TransactOpts, proposalId *big.Int, support uint8, reason string, params []byte) (*types.Transaction, error) {
	return _AWPDAO.contract.Transact(opts, "castVoteWithReasonAndParams", proposalId, support, reason, params)
}

// CastVoteWithReasonAndParams is a paid mutator transaction binding the contract method 0x5f398a14.
//
// Solidity: function castVoteWithReasonAndParams(uint256 proposalId, uint8 support, string reason, bytes params) returns(uint256)
func (_AWPDAO *AWPDAOSession) CastVoteWithReasonAndParams(proposalId *big.Int, support uint8, reason string, params []byte) (*types.Transaction, error) {
	return _AWPDAO.Contract.CastVoteWithReasonAndParams(&_AWPDAO.TransactOpts, proposalId, support, reason, params)
}

// CastVoteWithReasonAndParams is a paid mutator transaction binding the contract method 0x5f398a14.
//
// Solidity: function castVoteWithReasonAndParams(uint256 proposalId, uint8 support, string reason, bytes params) returns(uint256)
func (_AWPDAO *AWPDAOTransactorSession) CastVoteWithReasonAndParams(proposalId *big.Int, support uint8, reason string, params []byte) (*types.Transaction, error) {
	return _AWPDAO.Contract.CastVoteWithReasonAndParams(&_AWPDAO.TransactOpts, proposalId, support, reason, params)
}

// Execute is a paid mutator transaction binding the contract method 0x2656227d.
//
// Solidity: function execute(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) payable returns(uint256)
func (_AWPDAO *AWPDAOTransactor) Execute(opts *bind.TransactOpts, targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*types.Transaction, error) {
	return _AWPDAO.contract.Transact(opts, "execute", targets, values, calldatas, descriptionHash)
}

// Execute is a paid mutator transaction binding the contract method 0x2656227d.
//
// Solidity: function execute(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) payable returns(uint256)
func (_AWPDAO *AWPDAOSession) Execute(targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*types.Transaction, error) {
	return _AWPDAO.Contract.Execute(&_AWPDAO.TransactOpts, targets, values, calldatas, descriptionHash)
}

// Execute is a paid mutator transaction binding the contract method 0x2656227d.
//
// Solidity: function execute(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) payable returns(uint256)
func (_AWPDAO *AWPDAOTransactorSession) Execute(targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*types.Transaction, error) {
	return _AWPDAO.Contract.Execute(&_AWPDAO.TransactOpts, targets, values, calldatas, descriptionHash)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_AWPDAO *AWPDAOTransactor) OnERC1155BatchReceived(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _AWPDAO.contract.Transact(opts, "onERC1155BatchReceived", arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_AWPDAO *AWPDAOSession) OnERC1155BatchReceived(arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _AWPDAO.Contract.OnERC1155BatchReceived(&_AWPDAO.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_AWPDAO *AWPDAOTransactorSession) OnERC1155BatchReceived(arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _AWPDAO.Contract.OnERC1155BatchReceived(&_AWPDAO.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_AWPDAO *AWPDAOTransactor) OnERC1155Received(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _AWPDAO.contract.Transact(opts, "onERC1155Received", arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_AWPDAO *AWPDAOSession) OnERC1155Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _AWPDAO.Contract.OnERC1155Received(&_AWPDAO.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_AWPDAO *AWPDAOTransactorSession) OnERC1155Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _AWPDAO.Contract.OnERC1155Received(&_AWPDAO.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_AWPDAO *AWPDAOTransactor) OnERC721Received(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _AWPDAO.contract.Transact(opts, "onERC721Received", arg0, arg1, arg2, arg3)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_AWPDAO *AWPDAOSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _AWPDAO.Contract.OnERC721Received(&_AWPDAO.TransactOpts, arg0, arg1, arg2, arg3)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_AWPDAO *AWPDAOTransactorSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _AWPDAO.Contract.OnERC721Received(&_AWPDAO.TransactOpts, arg0, arg1, arg2, arg3)
}

// ProposeWithTokens is a paid mutator transaction binding the contract method 0xb407dd87.
//
// Solidity: function proposeWithTokens(address[] targets, uint256[] values, bytes[] calldatas, string description, uint256[] tokenIds) returns(uint256)
func (_AWPDAO *AWPDAOTransactor) ProposeWithTokens(opts *bind.TransactOpts, targets []common.Address, values []*big.Int, calldatas [][]byte, description string, tokenIds []*big.Int) (*types.Transaction, error) {
	return _AWPDAO.contract.Transact(opts, "proposeWithTokens", targets, values, calldatas, description, tokenIds)
}

// ProposeWithTokens is a paid mutator transaction binding the contract method 0xb407dd87.
//
// Solidity: function proposeWithTokens(address[] targets, uint256[] values, bytes[] calldatas, string description, uint256[] tokenIds) returns(uint256)
func (_AWPDAO *AWPDAOSession) ProposeWithTokens(targets []common.Address, values []*big.Int, calldatas [][]byte, description string, tokenIds []*big.Int) (*types.Transaction, error) {
	return _AWPDAO.Contract.ProposeWithTokens(&_AWPDAO.TransactOpts, targets, values, calldatas, description, tokenIds)
}

// ProposeWithTokens is a paid mutator transaction binding the contract method 0xb407dd87.
//
// Solidity: function proposeWithTokens(address[] targets, uint256[] values, bytes[] calldatas, string description, uint256[] tokenIds) returns(uint256)
func (_AWPDAO *AWPDAOTransactorSession) ProposeWithTokens(targets []common.Address, values []*big.Int, calldatas [][]byte, description string, tokenIds []*big.Int) (*types.Transaction, error) {
	return _AWPDAO.Contract.ProposeWithTokens(&_AWPDAO.TransactOpts, targets, values, calldatas, description, tokenIds)
}

// Queue is a paid mutator transaction binding the contract method 0x160cbed7.
//
// Solidity: function queue(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) returns(uint256)
func (_AWPDAO *AWPDAOTransactor) Queue(opts *bind.TransactOpts, targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*types.Transaction, error) {
	return _AWPDAO.contract.Transact(opts, "queue", targets, values, calldatas, descriptionHash)
}

// Queue is a paid mutator transaction binding the contract method 0x160cbed7.
//
// Solidity: function queue(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) returns(uint256)
func (_AWPDAO *AWPDAOSession) Queue(targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*types.Transaction, error) {
	return _AWPDAO.Contract.Queue(&_AWPDAO.TransactOpts, targets, values, calldatas, descriptionHash)
}

// Queue is a paid mutator transaction binding the contract method 0x160cbed7.
//
// Solidity: function queue(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) returns(uint256)
func (_AWPDAO *AWPDAOTransactorSession) Queue(targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*types.Transaction, error) {
	return _AWPDAO.Contract.Queue(&_AWPDAO.TransactOpts, targets, values, calldatas, descriptionHash)
}

// Relay is a paid mutator transaction binding the contract method 0xc28bc2fa.
//
// Solidity: function relay(address target, uint256 value, bytes data) payable returns()
func (_AWPDAO *AWPDAOTransactor) Relay(opts *bind.TransactOpts, target common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _AWPDAO.contract.Transact(opts, "relay", target, value, data)
}

// Relay is a paid mutator transaction binding the contract method 0xc28bc2fa.
//
// Solidity: function relay(address target, uint256 value, bytes data) payable returns()
func (_AWPDAO *AWPDAOSession) Relay(target common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _AWPDAO.Contract.Relay(&_AWPDAO.TransactOpts, target, value, data)
}

// Relay is a paid mutator transaction binding the contract method 0xc28bc2fa.
//
// Solidity: function relay(address target, uint256 value, bytes data) payable returns()
func (_AWPDAO *AWPDAOTransactorSession) Relay(target common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _AWPDAO.Contract.Relay(&_AWPDAO.TransactOpts, target, value, data)
}

// SetProposalThreshold is a paid mutator transaction binding the contract method 0xece40cc1.
//
// Solidity: function setProposalThreshold(uint256 newProposalThreshold) returns()
func (_AWPDAO *AWPDAOTransactor) SetProposalThreshold(opts *bind.TransactOpts, newProposalThreshold *big.Int) (*types.Transaction, error) {
	return _AWPDAO.contract.Transact(opts, "setProposalThreshold", newProposalThreshold)
}

// SetProposalThreshold is a paid mutator transaction binding the contract method 0xece40cc1.
//
// Solidity: function setProposalThreshold(uint256 newProposalThreshold) returns()
func (_AWPDAO *AWPDAOSession) SetProposalThreshold(newProposalThreshold *big.Int) (*types.Transaction, error) {
	return _AWPDAO.Contract.SetProposalThreshold(&_AWPDAO.TransactOpts, newProposalThreshold)
}

// SetProposalThreshold is a paid mutator transaction binding the contract method 0xece40cc1.
//
// Solidity: function setProposalThreshold(uint256 newProposalThreshold) returns()
func (_AWPDAO *AWPDAOTransactorSession) SetProposalThreshold(newProposalThreshold *big.Int) (*types.Transaction, error) {
	return _AWPDAO.Contract.SetProposalThreshold(&_AWPDAO.TransactOpts, newProposalThreshold)
}

// SetQuorumPercent is a paid mutator transaction binding the contract method 0x797294ac.
//
// Solidity: function setQuorumPercent(uint256 newQuorumPercent) returns()
func (_AWPDAO *AWPDAOTransactor) SetQuorumPercent(opts *bind.TransactOpts, newQuorumPercent *big.Int) (*types.Transaction, error) {
	return _AWPDAO.contract.Transact(opts, "setQuorumPercent", newQuorumPercent)
}

// SetQuorumPercent is a paid mutator transaction binding the contract method 0x797294ac.
//
// Solidity: function setQuorumPercent(uint256 newQuorumPercent) returns()
func (_AWPDAO *AWPDAOSession) SetQuorumPercent(newQuorumPercent *big.Int) (*types.Transaction, error) {
	return _AWPDAO.Contract.SetQuorumPercent(&_AWPDAO.TransactOpts, newQuorumPercent)
}

// SetQuorumPercent is a paid mutator transaction binding the contract method 0x797294ac.
//
// Solidity: function setQuorumPercent(uint256 newQuorumPercent) returns()
func (_AWPDAO *AWPDAOTransactorSession) SetQuorumPercent(newQuorumPercent *big.Int) (*types.Transaction, error) {
	return _AWPDAO.Contract.SetQuorumPercent(&_AWPDAO.TransactOpts, newQuorumPercent)
}

// SetVotingDelay is a paid mutator transaction binding the contract method 0x79051887.
//
// Solidity: function setVotingDelay(uint48 newVotingDelay) returns()
func (_AWPDAO *AWPDAOTransactor) SetVotingDelay(opts *bind.TransactOpts, newVotingDelay *big.Int) (*types.Transaction, error) {
	return _AWPDAO.contract.Transact(opts, "setVotingDelay", newVotingDelay)
}

// SetVotingDelay is a paid mutator transaction binding the contract method 0x79051887.
//
// Solidity: function setVotingDelay(uint48 newVotingDelay) returns()
func (_AWPDAO *AWPDAOSession) SetVotingDelay(newVotingDelay *big.Int) (*types.Transaction, error) {
	return _AWPDAO.Contract.SetVotingDelay(&_AWPDAO.TransactOpts, newVotingDelay)
}

// SetVotingDelay is a paid mutator transaction binding the contract method 0x79051887.
//
// Solidity: function setVotingDelay(uint48 newVotingDelay) returns()
func (_AWPDAO *AWPDAOTransactorSession) SetVotingDelay(newVotingDelay *big.Int) (*types.Transaction, error) {
	return _AWPDAO.Contract.SetVotingDelay(&_AWPDAO.TransactOpts, newVotingDelay)
}

// SetVotingPeriod is a paid mutator transaction binding the contract method 0xe540d01d.
//
// Solidity: function setVotingPeriod(uint32 newVotingPeriod) returns()
func (_AWPDAO *AWPDAOTransactor) SetVotingPeriod(opts *bind.TransactOpts, newVotingPeriod uint32) (*types.Transaction, error) {
	return _AWPDAO.contract.Transact(opts, "setVotingPeriod", newVotingPeriod)
}

// SetVotingPeriod is a paid mutator transaction binding the contract method 0xe540d01d.
//
// Solidity: function setVotingPeriod(uint32 newVotingPeriod) returns()
func (_AWPDAO *AWPDAOSession) SetVotingPeriod(newVotingPeriod uint32) (*types.Transaction, error) {
	return _AWPDAO.Contract.SetVotingPeriod(&_AWPDAO.TransactOpts, newVotingPeriod)
}

// SetVotingPeriod is a paid mutator transaction binding the contract method 0xe540d01d.
//
// Solidity: function setVotingPeriod(uint32 newVotingPeriod) returns()
func (_AWPDAO *AWPDAOTransactorSession) SetVotingPeriod(newVotingPeriod uint32) (*types.Transaction, error) {
	return _AWPDAO.Contract.SetVotingPeriod(&_AWPDAO.TransactOpts, newVotingPeriod)
}

// SignalPropose is a paid mutator transaction binding the contract method 0xb1b5d01d.
//
// Solidity: function signalPropose(string description, uint256[] tokenIds) returns(uint256)
func (_AWPDAO *AWPDAOTransactor) SignalPropose(opts *bind.TransactOpts, description string, tokenIds []*big.Int) (*types.Transaction, error) {
	return _AWPDAO.contract.Transact(opts, "signalPropose", description, tokenIds)
}

// SignalPropose is a paid mutator transaction binding the contract method 0xb1b5d01d.
//
// Solidity: function signalPropose(string description, uint256[] tokenIds) returns(uint256)
func (_AWPDAO *AWPDAOSession) SignalPropose(description string, tokenIds []*big.Int) (*types.Transaction, error) {
	return _AWPDAO.Contract.SignalPropose(&_AWPDAO.TransactOpts, description, tokenIds)
}

// SignalPropose is a paid mutator transaction binding the contract method 0xb1b5d01d.
//
// Solidity: function signalPropose(string description, uint256[] tokenIds) returns(uint256)
func (_AWPDAO *AWPDAOTransactorSession) SignalPropose(description string, tokenIds []*big.Int) (*types.Transaction, error) {
	return _AWPDAO.Contract.SignalPropose(&_AWPDAO.TransactOpts, description, tokenIds)
}

// UpdateTimelock is a paid mutator transaction binding the contract method 0xa890c910.
//
// Solidity: function updateTimelock(address newTimelock) returns()
func (_AWPDAO *AWPDAOTransactor) UpdateTimelock(opts *bind.TransactOpts, newTimelock common.Address) (*types.Transaction, error) {
	return _AWPDAO.contract.Transact(opts, "updateTimelock", newTimelock)
}

// UpdateTimelock is a paid mutator transaction binding the contract method 0xa890c910.
//
// Solidity: function updateTimelock(address newTimelock) returns()
func (_AWPDAO *AWPDAOSession) UpdateTimelock(newTimelock common.Address) (*types.Transaction, error) {
	return _AWPDAO.Contract.UpdateTimelock(&_AWPDAO.TransactOpts, newTimelock)
}

// UpdateTimelock is a paid mutator transaction binding the contract method 0xa890c910.
//
// Solidity: function updateTimelock(address newTimelock) returns()
func (_AWPDAO *AWPDAOTransactorSession) UpdateTimelock(newTimelock common.Address) (*types.Transaction, error) {
	return _AWPDAO.Contract.UpdateTimelock(&_AWPDAO.TransactOpts, newTimelock)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_AWPDAO *AWPDAOTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AWPDAO.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_AWPDAO *AWPDAOSession) Receive() (*types.Transaction, error) {
	return _AWPDAO.Contract.Receive(&_AWPDAO.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_AWPDAO *AWPDAOTransactorSession) Receive() (*types.Transaction, error) {
	return _AWPDAO.Contract.Receive(&_AWPDAO.TransactOpts)
}

// AWPDAOEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the AWPDAO contract.
type AWPDAOEIP712DomainChangedIterator struct {
	Event *AWPDAOEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *AWPDAOEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPDAOEIP712DomainChanged)
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
		it.Event = new(AWPDAOEIP712DomainChanged)
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
func (it *AWPDAOEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPDAOEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPDAOEIP712DomainChanged represents a EIP712DomainChanged event raised by the AWPDAO contract.
type AWPDAOEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_AWPDAO *AWPDAOFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*AWPDAOEIP712DomainChangedIterator, error) {

	logs, sub, err := _AWPDAO.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &AWPDAOEIP712DomainChangedIterator{contract: _AWPDAO.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_AWPDAO *AWPDAOFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *AWPDAOEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _AWPDAO.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPDAOEIP712DomainChanged)
				if err := _AWPDAO.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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
func (_AWPDAO *AWPDAOFilterer) ParseEIP712DomainChanged(log types.Log) (*AWPDAOEIP712DomainChanged, error) {
	event := new(AWPDAOEIP712DomainChanged)
	if err := _AWPDAO.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPDAOProposalCanceledIterator is returned from FilterProposalCanceled and is used to iterate over the raw logs and unpacked data for ProposalCanceled events raised by the AWPDAO contract.
type AWPDAOProposalCanceledIterator struct {
	Event *AWPDAOProposalCanceled // Event containing the contract specifics and raw log

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
func (it *AWPDAOProposalCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPDAOProposalCanceled)
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
		it.Event = new(AWPDAOProposalCanceled)
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
func (it *AWPDAOProposalCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPDAOProposalCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPDAOProposalCanceled represents a ProposalCanceled event raised by the AWPDAO contract.
type AWPDAOProposalCanceled struct {
	ProposalId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalCanceled is a free log retrieval operation binding the contract event 0x789cf55be980739dad1d0699b93b58e806b51c9d96619bfa8fe0a28abaa7b30c.
//
// Solidity: event ProposalCanceled(uint256 proposalId)
func (_AWPDAO *AWPDAOFilterer) FilterProposalCanceled(opts *bind.FilterOpts) (*AWPDAOProposalCanceledIterator, error) {

	logs, sub, err := _AWPDAO.contract.FilterLogs(opts, "ProposalCanceled")
	if err != nil {
		return nil, err
	}
	return &AWPDAOProposalCanceledIterator{contract: _AWPDAO.contract, event: "ProposalCanceled", logs: logs, sub: sub}, nil
}

// WatchProposalCanceled is a free log subscription operation binding the contract event 0x789cf55be980739dad1d0699b93b58e806b51c9d96619bfa8fe0a28abaa7b30c.
//
// Solidity: event ProposalCanceled(uint256 proposalId)
func (_AWPDAO *AWPDAOFilterer) WatchProposalCanceled(opts *bind.WatchOpts, sink chan<- *AWPDAOProposalCanceled) (event.Subscription, error) {

	logs, sub, err := _AWPDAO.contract.WatchLogs(opts, "ProposalCanceled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPDAOProposalCanceled)
				if err := _AWPDAO.contract.UnpackLog(event, "ProposalCanceled", log); err != nil {
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

// ParseProposalCanceled is a log parse operation binding the contract event 0x789cf55be980739dad1d0699b93b58e806b51c9d96619bfa8fe0a28abaa7b30c.
//
// Solidity: event ProposalCanceled(uint256 proposalId)
func (_AWPDAO *AWPDAOFilterer) ParseProposalCanceled(log types.Log) (*AWPDAOProposalCanceled, error) {
	event := new(AWPDAOProposalCanceled)
	if err := _AWPDAO.contract.UnpackLog(event, "ProposalCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPDAOProposalCreatedIterator is returned from FilterProposalCreated and is used to iterate over the raw logs and unpacked data for ProposalCreated events raised by the AWPDAO contract.
type AWPDAOProposalCreatedIterator struct {
	Event *AWPDAOProposalCreated // Event containing the contract specifics and raw log

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
func (it *AWPDAOProposalCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPDAOProposalCreated)
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
		it.Event = new(AWPDAOProposalCreated)
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
func (it *AWPDAOProposalCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPDAOProposalCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPDAOProposalCreated represents a ProposalCreated event raised by the AWPDAO contract.
type AWPDAOProposalCreated struct {
	ProposalId  *big.Int
	Proposer    common.Address
	Targets     []common.Address
	Values      []*big.Int
	Signatures  []string
	Calldatas   [][]byte
	VoteStart   *big.Int
	VoteEnd     *big.Int
	Description string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterProposalCreated is a free log retrieval operation binding the contract event 0x7d84a6263ae0d98d3329bd7b46bb4e8d6f98cd35a7adb45c274c8b7fd5ebd5e0.
//
// Solidity: event ProposalCreated(uint256 proposalId, address proposer, address[] targets, uint256[] values, string[] signatures, bytes[] calldatas, uint256 voteStart, uint256 voteEnd, string description)
func (_AWPDAO *AWPDAOFilterer) FilterProposalCreated(opts *bind.FilterOpts) (*AWPDAOProposalCreatedIterator, error) {

	logs, sub, err := _AWPDAO.contract.FilterLogs(opts, "ProposalCreated")
	if err != nil {
		return nil, err
	}
	return &AWPDAOProposalCreatedIterator{contract: _AWPDAO.contract, event: "ProposalCreated", logs: logs, sub: sub}, nil
}

// WatchProposalCreated is a free log subscription operation binding the contract event 0x7d84a6263ae0d98d3329bd7b46bb4e8d6f98cd35a7adb45c274c8b7fd5ebd5e0.
//
// Solidity: event ProposalCreated(uint256 proposalId, address proposer, address[] targets, uint256[] values, string[] signatures, bytes[] calldatas, uint256 voteStart, uint256 voteEnd, string description)
func (_AWPDAO *AWPDAOFilterer) WatchProposalCreated(opts *bind.WatchOpts, sink chan<- *AWPDAOProposalCreated) (event.Subscription, error) {

	logs, sub, err := _AWPDAO.contract.WatchLogs(opts, "ProposalCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPDAOProposalCreated)
				if err := _AWPDAO.contract.UnpackLog(event, "ProposalCreated", log); err != nil {
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

// ParseProposalCreated is a log parse operation binding the contract event 0x7d84a6263ae0d98d3329bd7b46bb4e8d6f98cd35a7adb45c274c8b7fd5ebd5e0.
//
// Solidity: event ProposalCreated(uint256 proposalId, address proposer, address[] targets, uint256[] values, string[] signatures, bytes[] calldatas, uint256 voteStart, uint256 voteEnd, string description)
func (_AWPDAO *AWPDAOFilterer) ParseProposalCreated(log types.Log) (*AWPDAOProposalCreated, error) {
	event := new(AWPDAOProposalCreated)
	if err := _AWPDAO.contract.UnpackLog(event, "ProposalCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPDAOProposalExecutedIterator is returned from FilterProposalExecuted and is used to iterate over the raw logs and unpacked data for ProposalExecuted events raised by the AWPDAO contract.
type AWPDAOProposalExecutedIterator struct {
	Event *AWPDAOProposalExecuted // Event containing the contract specifics and raw log

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
func (it *AWPDAOProposalExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPDAOProposalExecuted)
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
		it.Event = new(AWPDAOProposalExecuted)
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
func (it *AWPDAOProposalExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPDAOProposalExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPDAOProposalExecuted represents a ProposalExecuted event raised by the AWPDAO contract.
type AWPDAOProposalExecuted struct {
	ProposalId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalExecuted is a free log retrieval operation binding the contract event 0x712ae1383f79ac853f8d882153778e0260ef8f03b504e2866e0593e04d2b291f.
//
// Solidity: event ProposalExecuted(uint256 proposalId)
func (_AWPDAO *AWPDAOFilterer) FilterProposalExecuted(opts *bind.FilterOpts) (*AWPDAOProposalExecutedIterator, error) {

	logs, sub, err := _AWPDAO.contract.FilterLogs(opts, "ProposalExecuted")
	if err != nil {
		return nil, err
	}
	return &AWPDAOProposalExecutedIterator{contract: _AWPDAO.contract, event: "ProposalExecuted", logs: logs, sub: sub}, nil
}

// WatchProposalExecuted is a free log subscription operation binding the contract event 0x712ae1383f79ac853f8d882153778e0260ef8f03b504e2866e0593e04d2b291f.
//
// Solidity: event ProposalExecuted(uint256 proposalId)
func (_AWPDAO *AWPDAOFilterer) WatchProposalExecuted(opts *bind.WatchOpts, sink chan<- *AWPDAOProposalExecuted) (event.Subscription, error) {

	logs, sub, err := _AWPDAO.contract.WatchLogs(opts, "ProposalExecuted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPDAOProposalExecuted)
				if err := _AWPDAO.contract.UnpackLog(event, "ProposalExecuted", log); err != nil {
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

// ParseProposalExecuted is a log parse operation binding the contract event 0x712ae1383f79ac853f8d882153778e0260ef8f03b504e2866e0593e04d2b291f.
//
// Solidity: event ProposalExecuted(uint256 proposalId)
func (_AWPDAO *AWPDAOFilterer) ParseProposalExecuted(log types.Log) (*AWPDAOProposalExecuted, error) {
	event := new(AWPDAOProposalExecuted)
	if err := _AWPDAO.contract.UnpackLog(event, "ProposalExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPDAOProposalQueuedIterator is returned from FilterProposalQueued and is used to iterate over the raw logs and unpacked data for ProposalQueued events raised by the AWPDAO contract.
type AWPDAOProposalQueuedIterator struct {
	Event *AWPDAOProposalQueued // Event containing the contract specifics and raw log

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
func (it *AWPDAOProposalQueuedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPDAOProposalQueued)
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
		it.Event = new(AWPDAOProposalQueued)
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
func (it *AWPDAOProposalQueuedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPDAOProposalQueuedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPDAOProposalQueued represents a ProposalQueued event raised by the AWPDAO contract.
type AWPDAOProposalQueued struct {
	ProposalId *big.Int
	EtaSeconds *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalQueued is a free log retrieval operation binding the contract event 0x9a2e42fd6722813d69113e7d0079d3d940171428df7373df9c7f7617cfda2892.
//
// Solidity: event ProposalQueued(uint256 proposalId, uint256 etaSeconds)
func (_AWPDAO *AWPDAOFilterer) FilterProposalQueued(opts *bind.FilterOpts) (*AWPDAOProposalQueuedIterator, error) {

	logs, sub, err := _AWPDAO.contract.FilterLogs(opts, "ProposalQueued")
	if err != nil {
		return nil, err
	}
	return &AWPDAOProposalQueuedIterator{contract: _AWPDAO.contract, event: "ProposalQueued", logs: logs, sub: sub}, nil
}

// WatchProposalQueued is a free log subscription operation binding the contract event 0x9a2e42fd6722813d69113e7d0079d3d940171428df7373df9c7f7617cfda2892.
//
// Solidity: event ProposalQueued(uint256 proposalId, uint256 etaSeconds)
func (_AWPDAO *AWPDAOFilterer) WatchProposalQueued(opts *bind.WatchOpts, sink chan<- *AWPDAOProposalQueued) (event.Subscription, error) {

	logs, sub, err := _AWPDAO.contract.WatchLogs(opts, "ProposalQueued")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPDAOProposalQueued)
				if err := _AWPDAO.contract.UnpackLog(event, "ProposalQueued", log); err != nil {
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

// ParseProposalQueued is a log parse operation binding the contract event 0x9a2e42fd6722813d69113e7d0079d3d940171428df7373df9c7f7617cfda2892.
//
// Solidity: event ProposalQueued(uint256 proposalId, uint256 etaSeconds)
func (_AWPDAO *AWPDAOFilterer) ParseProposalQueued(log types.Log) (*AWPDAOProposalQueued, error) {
	event := new(AWPDAOProposalQueued)
	if err := _AWPDAO.contract.UnpackLog(event, "ProposalQueued", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPDAOProposalThresholdSetIterator is returned from FilterProposalThresholdSet and is used to iterate over the raw logs and unpacked data for ProposalThresholdSet events raised by the AWPDAO contract.
type AWPDAOProposalThresholdSetIterator struct {
	Event *AWPDAOProposalThresholdSet // Event containing the contract specifics and raw log

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
func (it *AWPDAOProposalThresholdSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPDAOProposalThresholdSet)
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
		it.Event = new(AWPDAOProposalThresholdSet)
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
func (it *AWPDAOProposalThresholdSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPDAOProposalThresholdSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPDAOProposalThresholdSet represents a ProposalThresholdSet event raised by the AWPDAO contract.
type AWPDAOProposalThresholdSet struct {
	OldProposalThreshold *big.Int
	NewProposalThreshold *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterProposalThresholdSet is a free log retrieval operation binding the contract event 0xccb45da8d5717e6c4544694297c4ba5cf151d455c9bb0ed4fc7a38411bc05461.
//
// Solidity: event ProposalThresholdSet(uint256 oldProposalThreshold, uint256 newProposalThreshold)
func (_AWPDAO *AWPDAOFilterer) FilterProposalThresholdSet(opts *bind.FilterOpts) (*AWPDAOProposalThresholdSetIterator, error) {

	logs, sub, err := _AWPDAO.contract.FilterLogs(opts, "ProposalThresholdSet")
	if err != nil {
		return nil, err
	}
	return &AWPDAOProposalThresholdSetIterator{contract: _AWPDAO.contract, event: "ProposalThresholdSet", logs: logs, sub: sub}, nil
}

// WatchProposalThresholdSet is a free log subscription operation binding the contract event 0xccb45da8d5717e6c4544694297c4ba5cf151d455c9bb0ed4fc7a38411bc05461.
//
// Solidity: event ProposalThresholdSet(uint256 oldProposalThreshold, uint256 newProposalThreshold)
func (_AWPDAO *AWPDAOFilterer) WatchProposalThresholdSet(opts *bind.WatchOpts, sink chan<- *AWPDAOProposalThresholdSet) (event.Subscription, error) {

	logs, sub, err := _AWPDAO.contract.WatchLogs(opts, "ProposalThresholdSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPDAOProposalThresholdSet)
				if err := _AWPDAO.contract.UnpackLog(event, "ProposalThresholdSet", log); err != nil {
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

// ParseProposalThresholdSet is a log parse operation binding the contract event 0xccb45da8d5717e6c4544694297c4ba5cf151d455c9bb0ed4fc7a38411bc05461.
//
// Solidity: event ProposalThresholdSet(uint256 oldProposalThreshold, uint256 newProposalThreshold)
func (_AWPDAO *AWPDAOFilterer) ParseProposalThresholdSet(log types.Log) (*AWPDAOProposalThresholdSet, error) {
	event := new(AWPDAOProposalThresholdSet)
	if err := _AWPDAO.contract.UnpackLog(event, "ProposalThresholdSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPDAOTimelockChangeIterator is returned from FilterTimelockChange and is used to iterate over the raw logs and unpacked data for TimelockChange events raised by the AWPDAO contract.
type AWPDAOTimelockChangeIterator struct {
	Event *AWPDAOTimelockChange // Event containing the contract specifics and raw log

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
func (it *AWPDAOTimelockChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPDAOTimelockChange)
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
		it.Event = new(AWPDAOTimelockChange)
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
func (it *AWPDAOTimelockChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPDAOTimelockChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPDAOTimelockChange represents a TimelockChange event raised by the AWPDAO contract.
type AWPDAOTimelockChange struct {
	OldTimelock common.Address
	NewTimelock common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTimelockChange is a free log retrieval operation binding the contract event 0x08f74ea46ef7894f65eabfb5e6e695de773a000b47c529ab559178069b226401.
//
// Solidity: event TimelockChange(address oldTimelock, address newTimelock)
func (_AWPDAO *AWPDAOFilterer) FilterTimelockChange(opts *bind.FilterOpts) (*AWPDAOTimelockChangeIterator, error) {

	logs, sub, err := _AWPDAO.contract.FilterLogs(opts, "TimelockChange")
	if err != nil {
		return nil, err
	}
	return &AWPDAOTimelockChangeIterator{contract: _AWPDAO.contract, event: "TimelockChange", logs: logs, sub: sub}, nil
}

// WatchTimelockChange is a free log subscription operation binding the contract event 0x08f74ea46ef7894f65eabfb5e6e695de773a000b47c529ab559178069b226401.
//
// Solidity: event TimelockChange(address oldTimelock, address newTimelock)
func (_AWPDAO *AWPDAOFilterer) WatchTimelockChange(opts *bind.WatchOpts, sink chan<- *AWPDAOTimelockChange) (event.Subscription, error) {

	logs, sub, err := _AWPDAO.contract.WatchLogs(opts, "TimelockChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPDAOTimelockChange)
				if err := _AWPDAO.contract.UnpackLog(event, "TimelockChange", log); err != nil {
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

// ParseTimelockChange is a log parse operation binding the contract event 0x08f74ea46ef7894f65eabfb5e6e695de773a000b47c529ab559178069b226401.
//
// Solidity: event TimelockChange(address oldTimelock, address newTimelock)
func (_AWPDAO *AWPDAOFilterer) ParseTimelockChange(log types.Log) (*AWPDAOTimelockChange, error) {
	event := new(AWPDAOTimelockChange)
	if err := _AWPDAO.contract.UnpackLog(event, "TimelockChange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPDAOVoteCastIterator is returned from FilterVoteCast and is used to iterate over the raw logs and unpacked data for VoteCast events raised by the AWPDAO contract.
type AWPDAOVoteCastIterator struct {
	Event *AWPDAOVoteCast // Event containing the contract specifics and raw log

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
func (it *AWPDAOVoteCastIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPDAOVoteCast)
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
		it.Event = new(AWPDAOVoteCast)
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
func (it *AWPDAOVoteCastIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPDAOVoteCastIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPDAOVoteCast represents a VoteCast event raised by the AWPDAO contract.
type AWPDAOVoteCast struct {
	Voter      common.Address
	ProposalId *big.Int
	Support    uint8
	Weight     *big.Int
	Reason     string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVoteCast is a free log retrieval operation binding the contract event 0xb8e138887d0aa13bab447e82de9d5c1777041ecd21ca36ba824ff1e6c07ddda4.
//
// Solidity: event VoteCast(address indexed voter, uint256 proposalId, uint8 support, uint256 weight, string reason)
func (_AWPDAO *AWPDAOFilterer) FilterVoteCast(opts *bind.FilterOpts, voter []common.Address) (*AWPDAOVoteCastIterator, error) {

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _AWPDAO.contract.FilterLogs(opts, "VoteCast", voterRule)
	if err != nil {
		return nil, err
	}
	return &AWPDAOVoteCastIterator{contract: _AWPDAO.contract, event: "VoteCast", logs: logs, sub: sub}, nil
}

// WatchVoteCast is a free log subscription operation binding the contract event 0xb8e138887d0aa13bab447e82de9d5c1777041ecd21ca36ba824ff1e6c07ddda4.
//
// Solidity: event VoteCast(address indexed voter, uint256 proposalId, uint8 support, uint256 weight, string reason)
func (_AWPDAO *AWPDAOFilterer) WatchVoteCast(opts *bind.WatchOpts, sink chan<- *AWPDAOVoteCast, voter []common.Address) (event.Subscription, error) {

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _AWPDAO.contract.WatchLogs(opts, "VoteCast", voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPDAOVoteCast)
				if err := _AWPDAO.contract.UnpackLog(event, "VoteCast", log); err != nil {
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

// ParseVoteCast is a log parse operation binding the contract event 0xb8e138887d0aa13bab447e82de9d5c1777041ecd21ca36ba824ff1e6c07ddda4.
//
// Solidity: event VoteCast(address indexed voter, uint256 proposalId, uint8 support, uint256 weight, string reason)
func (_AWPDAO *AWPDAOFilterer) ParseVoteCast(log types.Log) (*AWPDAOVoteCast, error) {
	event := new(AWPDAOVoteCast)
	if err := _AWPDAO.contract.UnpackLog(event, "VoteCast", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPDAOVoteCastWithParamsIterator is returned from FilterVoteCastWithParams and is used to iterate over the raw logs and unpacked data for VoteCastWithParams events raised by the AWPDAO contract.
type AWPDAOVoteCastWithParamsIterator struct {
	Event *AWPDAOVoteCastWithParams // Event containing the contract specifics and raw log

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
func (it *AWPDAOVoteCastWithParamsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPDAOVoteCastWithParams)
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
		it.Event = new(AWPDAOVoteCastWithParams)
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
func (it *AWPDAOVoteCastWithParamsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPDAOVoteCastWithParamsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPDAOVoteCastWithParams represents a VoteCastWithParams event raised by the AWPDAO contract.
type AWPDAOVoteCastWithParams struct {
	Voter      common.Address
	ProposalId *big.Int
	Support    uint8
	Weight     *big.Int
	Reason     string
	Params     []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVoteCastWithParams is a free log retrieval operation binding the contract event 0xe2babfbac5889a709b63bb7f598b324e08bc5a4fb9ec647fb3cbc9ec07eb8712.
//
// Solidity: event VoteCastWithParams(address indexed voter, uint256 proposalId, uint8 support, uint256 weight, string reason, bytes params)
func (_AWPDAO *AWPDAOFilterer) FilterVoteCastWithParams(opts *bind.FilterOpts, voter []common.Address) (*AWPDAOVoteCastWithParamsIterator, error) {

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _AWPDAO.contract.FilterLogs(opts, "VoteCastWithParams", voterRule)
	if err != nil {
		return nil, err
	}
	return &AWPDAOVoteCastWithParamsIterator{contract: _AWPDAO.contract, event: "VoteCastWithParams", logs: logs, sub: sub}, nil
}

// WatchVoteCastWithParams is a free log subscription operation binding the contract event 0xe2babfbac5889a709b63bb7f598b324e08bc5a4fb9ec647fb3cbc9ec07eb8712.
//
// Solidity: event VoteCastWithParams(address indexed voter, uint256 proposalId, uint8 support, uint256 weight, string reason, bytes params)
func (_AWPDAO *AWPDAOFilterer) WatchVoteCastWithParams(opts *bind.WatchOpts, sink chan<- *AWPDAOVoteCastWithParams, voter []common.Address) (event.Subscription, error) {

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _AWPDAO.contract.WatchLogs(opts, "VoteCastWithParams", voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPDAOVoteCastWithParams)
				if err := _AWPDAO.contract.UnpackLog(event, "VoteCastWithParams", log); err != nil {
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

// ParseVoteCastWithParams is a log parse operation binding the contract event 0xe2babfbac5889a709b63bb7f598b324e08bc5a4fb9ec647fb3cbc9ec07eb8712.
//
// Solidity: event VoteCastWithParams(address indexed voter, uint256 proposalId, uint8 support, uint256 weight, string reason, bytes params)
func (_AWPDAO *AWPDAOFilterer) ParseVoteCastWithParams(log types.Log) (*AWPDAOVoteCastWithParams, error) {
	event := new(AWPDAOVoteCastWithParams)
	if err := _AWPDAO.contract.UnpackLog(event, "VoteCastWithParams", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPDAOVotingDelaySetIterator is returned from FilterVotingDelaySet and is used to iterate over the raw logs and unpacked data for VotingDelaySet events raised by the AWPDAO contract.
type AWPDAOVotingDelaySetIterator struct {
	Event *AWPDAOVotingDelaySet // Event containing the contract specifics and raw log

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
func (it *AWPDAOVotingDelaySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPDAOVotingDelaySet)
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
		it.Event = new(AWPDAOVotingDelaySet)
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
func (it *AWPDAOVotingDelaySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPDAOVotingDelaySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPDAOVotingDelaySet represents a VotingDelaySet event raised by the AWPDAO contract.
type AWPDAOVotingDelaySet struct {
	OldVotingDelay *big.Int
	NewVotingDelay *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterVotingDelaySet is a free log retrieval operation binding the contract event 0xc565b045403dc03c2eea82b81a0465edad9e2e7fc4d97e11421c209da93d7a93.
//
// Solidity: event VotingDelaySet(uint256 oldVotingDelay, uint256 newVotingDelay)
func (_AWPDAO *AWPDAOFilterer) FilterVotingDelaySet(opts *bind.FilterOpts) (*AWPDAOVotingDelaySetIterator, error) {

	logs, sub, err := _AWPDAO.contract.FilterLogs(opts, "VotingDelaySet")
	if err != nil {
		return nil, err
	}
	return &AWPDAOVotingDelaySetIterator{contract: _AWPDAO.contract, event: "VotingDelaySet", logs: logs, sub: sub}, nil
}

// WatchVotingDelaySet is a free log subscription operation binding the contract event 0xc565b045403dc03c2eea82b81a0465edad9e2e7fc4d97e11421c209da93d7a93.
//
// Solidity: event VotingDelaySet(uint256 oldVotingDelay, uint256 newVotingDelay)
func (_AWPDAO *AWPDAOFilterer) WatchVotingDelaySet(opts *bind.WatchOpts, sink chan<- *AWPDAOVotingDelaySet) (event.Subscription, error) {

	logs, sub, err := _AWPDAO.contract.WatchLogs(opts, "VotingDelaySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPDAOVotingDelaySet)
				if err := _AWPDAO.contract.UnpackLog(event, "VotingDelaySet", log); err != nil {
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

// ParseVotingDelaySet is a log parse operation binding the contract event 0xc565b045403dc03c2eea82b81a0465edad9e2e7fc4d97e11421c209da93d7a93.
//
// Solidity: event VotingDelaySet(uint256 oldVotingDelay, uint256 newVotingDelay)
func (_AWPDAO *AWPDAOFilterer) ParseVotingDelaySet(log types.Log) (*AWPDAOVotingDelaySet, error) {
	event := new(AWPDAOVotingDelaySet)
	if err := _AWPDAO.contract.UnpackLog(event, "VotingDelaySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AWPDAOVotingPeriodSetIterator is returned from FilterVotingPeriodSet and is used to iterate over the raw logs and unpacked data for VotingPeriodSet events raised by the AWPDAO contract.
type AWPDAOVotingPeriodSetIterator struct {
	Event *AWPDAOVotingPeriodSet // Event containing the contract specifics and raw log

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
func (it *AWPDAOVotingPeriodSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AWPDAOVotingPeriodSet)
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
		it.Event = new(AWPDAOVotingPeriodSet)
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
func (it *AWPDAOVotingPeriodSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AWPDAOVotingPeriodSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AWPDAOVotingPeriodSet represents a VotingPeriodSet event raised by the AWPDAO contract.
type AWPDAOVotingPeriodSet struct {
	OldVotingPeriod *big.Int
	NewVotingPeriod *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterVotingPeriodSet is a free log retrieval operation binding the contract event 0x7e3f7f0708a84de9203036abaa450dccc85ad5ff52f78c170f3edb55cf5e8828.
//
// Solidity: event VotingPeriodSet(uint256 oldVotingPeriod, uint256 newVotingPeriod)
func (_AWPDAO *AWPDAOFilterer) FilterVotingPeriodSet(opts *bind.FilterOpts) (*AWPDAOVotingPeriodSetIterator, error) {

	logs, sub, err := _AWPDAO.contract.FilterLogs(opts, "VotingPeriodSet")
	if err != nil {
		return nil, err
	}
	return &AWPDAOVotingPeriodSetIterator{contract: _AWPDAO.contract, event: "VotingPeriodSet", logs: logs, sub: sub}, nil
}

// WatchVotingPeriodSet is a free log subscription operation binding the contract event 0x7e3f7f0708a84de9203036abaa450dccc85ad5ff52f78c170f3edb55cf5e8828.
//
// Solidity: event VotingPeriodSet(uint256 oldVotingPeriod, uint256 newVotingPeriod)
func (_AWPDAO *AWPDAOFilterer) WatchVotingPeriodSet(opts *bind.WatchOpts, sink chan<- *AWPDAOVotingPeriodSet) (event.Subscription, error) {

	logs, sub, err := _AWPDAO.contract.WatchLogs(opts, "VotingPeriodSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AWPDAOVotingPeriodSet)
				if err := _AWPDAO.contract.UnpackLog(event, "VotingPeriodSet", log); err != nil {
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

// ParseVotingPeriodSet is a log parse operation binding the contract event 0x7e3f7f0708a84de9203036abaa450dccc85ad5ff52f78c170f3edb55cf5e8828.
//
// Solidity: event VotingPeriodSet(uint256 oldVotingPeriod, uint256 newVotingPeriod)
func (_AWPDAO *AWPDAOFilterer) ParseVotingPeriodSet(log types.Log) (*AWPDAOVotingPeriodSet, error) {
	event := new(AWPDAOVotingPeriodSet)
	if err := _AWPDAO.contract.UnpackLog(event, "VotingPeriodSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
