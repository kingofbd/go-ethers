// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Counter {
    uint256 public count;

    // 递增计数器（需交易，消耗Gas）
    function increment() external {
        count += 1;
    }

    // 递减计数器（需交易，消耗Gas）
    function decrement() external {
        require(count > 0, "Counter: cannot decrement below zero");
        count -= 1;
    }

    // 获取当前计数值（只读，无Gas消耗）
    function getCount() external view returns (uint256) {
        return count;          
    }
}