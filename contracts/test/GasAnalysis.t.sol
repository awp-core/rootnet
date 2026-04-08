// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;
import {Test, console} from "forge-std/Test.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IERC20Permit} from "@openzeppelin/contracts/token/ERC20/extensions/IERC20Permit.sol";
import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import {veAWP} from "../src/core/veAWP.sol";
import {VeAWPHelper} from "../src/core/VeAWPHelper.sol";

contract GasAnalysis is Test {
    address constant AWP = 0x0000A1050AcF9DEA8af9c2E74f0D7CF43f1000A1;
    address constant VEAWP = 0x0000b534C63D78212f1BDCc315165852793A00A8;
    
    uint256 constant PK = 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80;
    address user;
    VeAWPHelper helper;

    function setUp() public {
        vm.createSelectFork(vm.envString("BASE_RPC_URL"));
        helper = new VeAWPHelper(AWP, VEAWP);
        user = vm.addr(PK);
        deal(AWP, user, 100_000 ether);
    }

    function test_gas_breakdown() public {
        uint256 amount = 10_000 ether;
        uint64 lockDuration = 90 days;
        uint256 deadline = block.timestamp + 1 hours;

        // Sign permit
        bytes32 ds = _getDomainSep();
        uint256 nonce = _getNonce(user);
        bytes32 structHash = keccak256(abi.encode(
            keccak256("Permit(address owner,address spender,uint256 value,uint256 nonce,uint256 deadline)"),
            user, address(helper), amount, nonce, deadline
        ));
        (uint8 v, bytes32 r, bytes32 s) = vm.sign(PK, keccak256(abi.encodePacked("\x19\x01", ds, structHash)));

        // === Step-by-step gas measurement ===
        
        // Step 1: Permit alone
        uint256 g0 = gasleft();
        IERC20Permit(AWP).permit(user, address(helper), amount, deadline, v, r, s);
        uint256 permitGas = g0 - gasleft();
        console.log("1. permit:                  ", permitGas);

        // Step 2: transferFrom(user -> helper)
        g0 = gasleft();
        vm.prank(address(helper));
        IERC20(AWP).transferFrom(user, address(helper), amount);
        uint256 tf1Gas = g0 - gasleft();
        console.log("2. transferFrom(user->help): ", tf1Gas);

        // Step 3: veAWP.deposit (includes internal transferFrom + mint + position write)
        g0 = gasleft();
        vm.prank(address(helper));
        uint256 tokenId = veAWP(VEAWP).deposit(amount, lockDuration);
        uint256 depositGas = g0 - gasleft();
        console.log("3. veAWP.deposit:           ", depositGas);

        // Step 4: NFT transferFrom(helper -> user)
        g0 = gasleft();
        vm.prank(address(helper));
        IERC721(VEAWP).transferFrom(address(helper), user, tokenId);
        uint256 nftGas = g0 - gasleft();
        console.log("4. NFT transfer(help->user): ", nftGas);

        console.log("----------------------------");
        console.log("TOTAL:                       ", permitGas + tf1Gas + depositGas + nftGas);
        console.log("(overhead vs direct deposit): ", permitGas + tf1Gas + nftGas);
    }

    function _getDomainSep() internal view returns (bytes32) {
        (bool ok, bytes memory d) = AWP.staticcall(abi.encodeWithSignature("DOMAIN_SEPARATOR()"));
        require(ok);
        return abi.decode(d, (bytes32));
    }
    function _getNonce(address o) internal view returns (uint256) {
        (bool ok, bytes memory d) = AWP.staticcall(abi.encodeWithSignature("nonces(address)", o));
        require(ok);
        return abi.decode(d, (uint256));
    }
}
