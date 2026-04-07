// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {ERC20Burnable} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import {ERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";
import {ERC1363} from "@openzeppelin/contracts/token/ERC20/extensions/ERC1363.sol";
import {IWorknetTokenFactory} from "../interfaces/IWorknetTokenFactory.sol";

/// @title WorknetToken — Worknet token (CREATE2 deterministic deployment)
/// @notice ERC20 + ERC20Permit + ERC1363 + Burnable. Single-minter with linear time-based minting cap.
/// @dev Constructor reads name/symbol/worknetId from factory via callback (no constructor args).
///      This makes creationCode constant → CREATE2 address depends only on salt, not on token params.
///      Vanity salts are universal: mine once, reuse for any worknet.
///      Deploy + mint + setMinter happen atomically in AWPRegistry._activateWorknet.
///      After lock: mintable ≤ (MAX_SUPPLY − supplyAtLock) × elapsed / 365 days.
contract WorknetToken is ERC1363, ERC20Burnable, ERC20Permit {
    uint256 public constant MAX_SUPPLY = 10_000_000_000 * 1e18;

    uint256 public immutable worknetId;
    uint64 public immutable createdAt;

    /// @notice Permanent minter. address(0) = not yet initialized (open mint phase).
    address public minter;
    /// @notice Supply snapshot at lock time (excludes LP pre-mint from time cap budget).
    uint256 public supplyAtLock;

    event MinterSet(address indexed newMinter);

    error NotMinter();
    error AlreadyInitialized();
    error ZeroAddress();
    error ExceedsMaxSupply();
    error ExceedsMintableLimit();
    /// @dev Reads token params from factory via callback. No constructor args → constant creationCode.
    constructor()
        ERC20(
            IWorknetTokenFactory(msg.sender).pendingName(),
            IWorknetTokenFactory(msg.sender).pendingSymbol()
        )
        ERC20Permit(IWorknetTokenFactory(msg.sender).pendingName())
    {
        worknetId = IWorknetTokenFactory(msg.sender).pendingWorknetId();
        createdAt = uint64(block.timestamp);
    }

    function initialized() external view returns (bool) {
        return minter != address(0);
    }

    /// @notice Lock the permanent minter. One-time call.
    /// @dev No access control by design — adding immutable state would change creationCode
    /// and break universal vanity salt. The race window during _activateWorknet is safe
    /// because lpManager is trusted AWP code and the call chain is nonReentrant.
    function setMinter(address newMinter) external {
        if (minter != address(0)) revert AlreadyInitialized();
        if (newMinter == address(0)) revert ZeroAddress();
        minter = newMinter;
        supplyAtLock = totalSupply();
        emit MinterSet(newMinter);
    }

    /// @notice Mint tokens. Before init: open. After init: minter-only + time cap.
    function mint(address to, uint256 amount) external {
        uint256 supply = totalSupply();
        if (supply + amount > MAX_SUPPLY) revert ExceedsMaxSupply();
        address _minter = minter;
        if (_minter != address(0)) {
            if (msg.sender != _minter) revert NotMinter();
            uint256 lock = supplyAtLock;
            uint256 elapsed = block.timestamp - uint256(createdAt);
            if (elapsed == 0) elapsed = 1;
            uint256 budget = MAX_SUPPLY - lock;
            uint256 cap = budget * elapsed / 365 days;
            if (cap > budget) cap = budget;
            if (supply - lock + amount > cap) revert ExceedsMintableLimit();
        }
        _mint(to, amount);
    }

    /// @notice Remaining mintable headroom at current timestamp.
    function currentMintableLimit() external view returns (uint256) {
        if (minter == address(0)) return 0;
        uint256 lock = supplyAtLock;
        uint256 elapsed = block.timestamp - uint256(createdAt);
        if (elapsed == 0) elapsed = 1;
        uint256 budget = MAX_SUPPLY - lock;
        uint256 cap = budget * elapsed / 365 days;
        if (cap > budget) cap = budget;
        uint256 supply = totalSupply();
        uint256 used = supply > lock ? supply - lock : 0;
        return used >= cap ? 0 : cap - used;
    }
}
