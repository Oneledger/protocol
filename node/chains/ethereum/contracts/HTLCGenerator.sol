pragma solidity ^0.4.0;

import "./HTLC.sol";

contract HTLCGenerator {
    using SafeMath for uint256;
    uint256 balances;

    function() external payable{
        createHTLC();
    }

    function createHTLC() public payable returns (HTLC newContract) {
        newContract = new HTLC(msg.sender);
    }
}
