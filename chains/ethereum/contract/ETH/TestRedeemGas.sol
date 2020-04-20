pragma solidity >=0.4.22 <0.6.0;
//pragma experimental ABIEncoderV2 ;

contract TestRedeemGas {
    mapping (address => uint) validators;
    mapping (address => bool) public isSigned;
    uint signature_count;
    uint threshold_signature;

    constructor(address[] memory validatorList) public {
        for (uint i = 0; i < validatorList.length; i++)
        {
            validators[validatorList[i]] = 1;
        }
        signature_count = 0;
        threshold_signature = validatorList.length;
    }


    function sign() public returns (bool) {
        require(isValidator(msg.sender),"validator not pressent in list");
        if(isSigned[msg.sender] == false)
        {
            isSigned[msg.sender] = true;
            signature_count ++ ;
        }
        return true;
    }


    function isValidator(address addr) public view returns(bool) {
        return validators[addr] > 0;
    }

    function processRedeem() public pure returns(bool){
        return true;
    }


}

