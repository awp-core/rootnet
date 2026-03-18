// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {TimelockController} from "@openzeppelin/contracts/governance/TimelockController.sol";

/// @title Treasury — Treasury contract based on OpenZeppelin TimelockController
/// @notice Standard OZ TimelockController implementation with zero custom code
/// @dev Design notes:
///   - All proposed operations are queued by AWPDAO and executed after minDelay
///   - proposers: list of addresses allowed to submit operations (typically AWPDAO)
///   - executors: list of addresses allowed to execute matured operations (typically address(0), meaning anyone can execute)
///   - admin: initial admin (typically renounced after deployment; managed by the governance process)
contract Treasury is TimelockController {
    /// @notice Constructor
    /// @param minDelay Minimum delay in seconds from operation queuing to execution
    /// @param proposers List of addresses allowed to submit operations
    /// @param executors List of addresses allowed to execute operations
    /// @param admin Initial admin address (can be address(0) for no additional admin)
    constructor(uint256 minDelay, address[] memory proposers, address[] memory executors, address admin)
        TimelockController(minDelay, proposers, executors, admin)
    {}
}
