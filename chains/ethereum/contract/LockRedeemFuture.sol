pragma solidity >=0.5.0 <0.6.0;

contract LockRedeem {
    //Flag to pause and unpause contract
    bool ACTIVE = false;
    address old_contract ;
    uint signaturesrequiredformigration = 0;
    uint migrationsignaturecount = 0 ;
    constructor(address _old_contract,uint noofValidatorsinold) public {
        old_contract = _old_contract;
        signaturesrequiredformigration = (noofValidatorsinold * 2 / 3 + 1);
    }
    function MigrateFromOld() public {
        require(msg.sender == old_contract);
        migrationsignaturecount = migrationsignaturecount + 1;
    }
    function () external payable {
        require(migrationsignaturecount == signaturesrequiredformigration);
        ACTIVE = true;
    }
    function isActive ()public view returns (bool) {
        return ACTIVE;
    }
    function getTotalEthBalance() public view returns(uint) {
        return address(this).balance ;
    }


}
