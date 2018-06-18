pragma solidity ^0.4.0;

import "./HTLC.sol";

contract HTLCGenerator {

    uint256 balances;
    mapping (address => address) public newHTLC;

    function() external payable{
        createHTLC();
    }

    function createHTLC() public payable  {
        newHTLC[msg.sender] = new HTLC(msg.sender);
    }

    function getAddress() public view returns (address) {
        return newHTLC[msg.sender];
    }
}
