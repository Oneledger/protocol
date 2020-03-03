pragma solidity 0.6.0;

contract LockRedeem {
    //Flag to pause and unpause contract
    bool ACTIVE = false;

    // numValidators holds the total number of validators
    uint public numValidators;
    //
    // struct ValidatorDetails {
    //     uint votingPower;
    //     uint validatorFee;
    //     uint8 index;
    // }
    mapping (address => uint8) public validators;
    address[] validatorList;

    // Require this amount of signatures to push a proposal through,Set only in constructor
    uint votingThreshold;
    uint activeThreshold;

    //the base cost of a transaction (21000 gas)
    //the cost of a contract deployment (32000 gas) [Not relevant for sign trasaction]
    //the cost for every zero byte of data or code for a transaction.
    //the cost of every non-zero byte of data or code for a transaction.
    uint trasactionCost = 23192 ;
    uint additionalGasCost = 37692 ;
    uint redeem_gas_charge = 10000000000000000;

    uint public migrationSignatures;
    mapping (address => bool) public migrationSigners;
    mapping (address => uint) migrationCount;
    address [] migrationAddress;

    // Default Voting power should be updated at one point
    int constant DEFAULT_VALIDATOR_POWER = 100;
    uint constant MIN_VALIDATORS = 0;
    //Approx
    // 270 blocks per hour
    // 5 blocks per min
    // 28800 old value
    uint256 LOCK_PERIOD;


    struct RedeemTX {
        address payable recipient;
        uint256 votes;
        uint256 amount;
        uint256 signature_count;
        uint256 until;
        uint256 redeemFee;
    }
    mapping (address => RedeemTX) redeemRequests;

    event RedeemRequest(
        address indexed recepient,
        uint256 amount_requested,
        uint256 redeemFeeCharged
    );

    event ValidatorSignedRedeem(
        address indexed recipient,
        address validator_addresss,
        uint256 amount,
        uint256 gasReturned
    );

    event Lock(
        address sender,
        uint256 amount_received
    );

    event ValidatorMigrated(
        address validator,
        address NewSmartContractAddress
    );

    event AddValidator(
        address indexed _address
    );

    modifier isActive() {
        require(ACTIVE);
        _;
    }

    modifier isInactive() {
        require(!ACTIVE);
        _;
    }

    modifier onlyValidator() {
        require(isValidator(msg.sender) ,"validator not present in List");
        _; // Continues control flow after this is validates
    }



    constructor(address[] memory initialValidators,uint _lock_period) public {
        // Require at least 4 validators
        require(initialValidators.length >= MIN_VALIDATORS, "insufficient validators passed to constructor");

        // Add the initial validators
        for (uint i = 0; i < initialValidators.length; i++) {
            // Ensure these validators are unique
            address v = initialValidators[i];
            require(validators[v] == 0, "found non-unique validator in initialValidators");
            addValidator(v);
        }
        ACTIVE = true ;
        LOCK_PERIOD = _lock_period;
        votingThreshold = (initialValidators.length * 2 / 3) + 1;
        activeThreshold = (initialValidators.length * 1 / 3) + 1;

    }

    function migrate (address newSmartContractAddress) public onlyValidator {
        require(migrationSigners[msg.sender]==false,"Validator Signed already");
        migrationSigners[msg.sender] = true ;
        (bool status,) = newSmartContractAddress.call(abi.encodePacked(bytes4(keccak256("MigrateFromOld()"))));
        require(status,"Unable to Migrate new Smart contract");
        migrationSignatures = migrationSignatures + 1;
        if(migrationCount[newSmartContractAddress]==0)
        {
            migrationAddress.push(newSmartContractAddress);
        }
        migrationCount[newSmartContractAddress] += 1;

        // Global flag ,needs to be set only once
        if (migrationSignatures == activeThreshold) {
            ACTIVE = false ;
        }
        // Trasfer needs to be done only once
        if (migrationSignatures == votingThreshold) {
            uint voteCount = 0 ;
            address maxVotedAddress ;
            for(uint i=0;i<migrationAddress.length;i++){
                if (migrationCount[migrationAddress[i]] > voteCount){
                    voteCount = migrationCount[migrationAddress[i]];
                    maxVotedAddress =migrationAddress[i];
                }
            }
            (bool success, ) = maxVotedAddress.call.value(address(this).balance)("");
            require(success, "Transfer failed");
        }
    }
    // function called by user
    function lock() payable public isActive {
        require(msg.value >= 0, "Must pay a balance more than 0");
        emit Lock(msg.sender,msg.value);
    }

    // function called by user
    function redeem(uint256 amount_)  payable public isActive {
        require(isredeemAvailable(msg.sender) ,"redeem to this address is not available yet");
        require(amount_ > 0, "amount should be bigger than 0");
        require(msg.value + redeemRequests[msg.sender].redeemFee >= redeem_gas_charge ,"Redeem fee not provided");

        // redeemRequests[msg.sender].isProcessing = true;
        redeemRequests[msg.sender].signature_count = uint256(0);
        redeemRequests[msg.sender].recipient = msg.sender;
        redeemRequests[msg.sender].amount = amount_ ;
        redeemRequests[msg.sender].until = block.number + LOCK_PERIOD;
        redeemRequests[msg.sender].redeemFee += msg.value;
        redeemRequests[msg.sender].votes = 0;
        emit RedeemRequest(redeemRequests[msg.sender].recipient,redeemRequests[msg.sender].amount,redeemRequests[msg.sender].redeemFee);
    }

    //function called by protocol
    function sign(uint amount_, address payable recipient_) public isActive onlyValidator {
        uint startGas = gasleft();
        require(isValidator(msg.sender),"validator not present in list");
        require(redeemRequests[recipient_].until > block.number, "redeem request is not available");
        require(redeemRequests[recipient_].amount == amount_,"redeem amount is different" );
        require((redeemRequests[recipient_].votes >> validators[msg.sender])% 2 == uint256(0), "validator already vote");

        // update votes
        redeemRequests[recipient_].votes = redeemRequests[recipient_].votes + (uint256(1) << validators[msg.sender]);
        redeemRequests[recipient_].signature_count += 1;

        // if reach threshold, transfer
        if (redeemRequests[recipient_].signature_count >= votingThreshold ) {
            (bool success, ) = redeemRequests[recipient_].recipient.call.value((redeemRequests[recipient_].amount))("");
            require(success, "Transfer failed.");
            redeemRequests[recipient_].amount = 0;
            redeemRequests[recipient_].until = block.number;

        }
        // Trasaction Cost is a an average trasaction cost .
        // additionalGasCost is the extra gas used by the lines of code after gas calculation is done.
        // This is an approximate calcualtion and actuall cost might vary sightly .
        uint gasUsed = startGas - gasleft() + trasactionCost + additionalGasCost ;
        //validators[msg.sender].validatorFee = gasUsed;
        uint gasFee = gasUsed * tx.gasprice;
        (bool success, ) = msg.sender.call.value(gasFee)("");
        // require(success, "Transfer back to validator failed");
        redeemRequests[recipient_].redeemFee -= gasFee;
        emit ValidatorSignedRedeem(recipient_, msg.sender, amount_,gasFee);
    }

    // function validatorFee() public view isActive onlyValidator returns(uint) {
    //     return validators[msg.sender].validatorFee;
    // }

    function collectUserFee() public isActive {
        require(verifyRedeem(msg.sender) > 0 , "request signing is still in progress");
        (bool success, ) = msg.sender.call.value(redeemRequests[msg.sender].redeemFee)("");
        require(success, "Transfer failed.");
    }

    function isValidator(address addr) public view returns(bool) {
        return validators[addr] > 0;
    }

    // function isValidator(address v) public view returns (bool) {
    //     return validators[v].votingPower > 0;
    // }

    function isredeemAvailable (address recepient_) public isActive view returns (bool)  {
        return redeemRequests[recepient_].until < block.number;
    }

    function hasValidatorSigned(address recipient_) public isActive view returns(bool) {
        return  ((redeemRequests[recipient_].votes >> validators[msg.sender]) % 2 == uint256(1));
    }

    //BeforeRedeem : -1  (amount == 0 and until = 0)
    //Ongoing : 0  (amount > 0 and until > block.number)
    //Success : 1  (amount = 0 and until < block.number)
    //Expired : 2  (amount > 0 and until < block.number)
    function verifyRedeem(address recipient_) public isActive view returns(int8){
        return (redeemRequests[recipient_].until > 0 && redeemRequests[recipient_].until < block.number) ? (redeemRequests[recipient_].amount > 0? int8(2) : int8(1) ): (redeemRequests[recipient_].amount == 0 ? int8(-1) : int8(0));
    }

    function getSignatureCount(address recipient_) public isActive view returns(uint256){
        return redeemRequests[recipient_].signature_count;
    }

    function getTotalEthBalance() public view returns(uint) {
        return address(this).balance ;
    }

    function getOLTEthAddress() public view returns(address){
        return address(this);
    }

    // internal functions
    function addValidator(address v) internal {
        validatorList.push(v);

        //validators[v].votingPower = DEFAULT_VALIDATOR_POWER;
        //validators[v].validatorFee = 0;
        validators[v] = uint8(validatorList.length);
        emit AddValidator(v);
    }

}
