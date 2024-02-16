// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

address constant STAKING_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000001005;

IStaking constant STAKING_CONTRACT = IStaking(STAKING_PRECOMPILE_ADDRESS);

interface IStaking {
    // Transactions
    function delegate(string memory valAddress, uint256 amount) external returns (bool success);

    function redelegate(string memory srcAddress, string memory dstAddress, uint256 amount)
        external
        returns (bool success);

    function undelegate(string memory valAddress, uint256 amount) external returns (bool success);

    function getDelegation(address delegator, string memory valAddress) external view returns (uint256 shares);

    function getStakingPool(string memory valAddress)
        external
        view
        returns (uint256 totalShares, uint256 totalTokens);
}
