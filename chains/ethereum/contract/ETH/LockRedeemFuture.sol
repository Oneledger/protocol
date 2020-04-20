pragma solidity >=0.5.0 <0.6.0;

contract LockRedeem {
    //Flag to pause and unpause contract
    bool ACTIVE = false;
    uint DEFAULT_VALIDATOR_POWER = 50;
    mapping (address => uint) public validators;
    address old_contract ;
    uint public numValidators = 0;
    uint signaturesrequiredformigration = 0;
    uint migrationsignaturecount = 0 ;
    address val;

    constructor(address _old_contract,uint noofValidatorsinold) public {
        old_contract = _old_contract;
        signaturesrequiredformigration = (noofValidatorsinold * 2 / 3 + 1);
    }
    function MigrateFromOld() public {
        require(msg.sender == old_contract);
        migrationsignaturecount = migrationsignaturecount + 1;
        addValidator(tx.origin);
    }
    function () external payable {
        require(migrationsignaturecount == signaturesrequiredformigration);
        ACTIVE = true;
    }
    function isActive ()public view returns (bool) {
        return ACTIVE;
    }

    function getMigrationCount() public view returns (uint) {
        return migrationsignaturecount;
    }
    function getTotalEthBalance() public view returns(uint) {
        return address(this).balance ;
    }
    function addValidator(address v) internal {
        validators[v] = DEFAULT_VALIDATOR_POWER;
        numValidators += 1;
    }


}