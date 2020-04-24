pragma solidity >=0.5.0 <0.6.0;

contract LockRedeem {
    //Flag to pause and unpause contract
    bool ACTIVE = false;
    // numValidators holds the total number of validators
    uint public numValidators;
    // mapping to store the validators to there power.
    mapping (address => uint) public validators;
    // Require this amount of signatures to push a proposal through,Set only in constructor
    uint votingThreshold;
    uint activeThreshold;

    // numValidators holds the total number of validators
    uint public migrationSignatures;
    // mapping to store the validators to there power.
    mapping (address => bool) public migrationSigners;
    uint constant DEFAULT_VALIDATOR_POWER = 50;
    uint constant MIN_VALIDATORS = 0;
    uint256 LOCK_PERIOD = 28800;

    struct RedeemTX {
        address payable recipient;
        mapping (address => bool) votes;
        uint256 amount;
        uint256 signature_count;
        bool isCompleted ;
        uint256 until;
    }
    mapping (address => RedeemTX) redeemRequests;

    event RedeemRequest(
        address indexed recepient,
        uint256 amount_requested
    );

    event ValidatorSignedRedeem(
        address indexed recipient,
        address validator_addresss,
        uint256 amount
    );

    event Lock(
        address sender,
        uint256 amount_received
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
        require(validators[msg.sender] >0 ,"validator not present in List");
        _; // Continues control flow after this is validates
    }

    function isValidator(address v) public view returns (bool) {
        return validators[v] > 0;
    }
    constructor(address[] memory initialValidators) public {
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
        votingThreshold = (initialValidators.length * 2 / 3) + 1;
        activeThreshold = (initialValidators.length * 1 / 3) + 1;
    }


    function migrate (address newSmartContractAddress) public onlyValidator {
        require(migrationSigners[msg.sender]==false,"Validator Signed already");
        migrationSigners[msg.sender] = true ;
        migrationSignatures = migrationSignatures + 1;
        (bool status,) = newSmartContractAddress.call(abi.encodePacked(bytes4(keccak256("MigrateFromOld()"))));
        require(status,"Unable to Migrate new Smart contract");

        // Global flag ,needs to be set only once
        if (migrationSignatures == activeThreshold) {
            ACTIVE = false ;
        }
        // Trasfer needs to be set only once
        if (migrationSignatures == votingThreshold) {
            (bool success, ) = newSmartContractAddress.call.value(address(this).balance)("");
            require(success, "Transfer failed");
        }
    }



    // function isValidator(address addr) public view returns(bool) {
    //     return validators[addr] == true;
    // }

    // function called by user
    function lock() payable public isActive {
        require(msg.value >= 0, "Must pay a balance more than 0");
        emit Lock(msg.sender,msg.value);
    }

    function isredeemAvailable (address recepient_) public view returns (bool){
        return redeemRequests[recepient_].until == 0 || redeemRequests[recepient_].until < block.number;
    }

    // function called by user
    function redeem(uint256 amount_)  public isActive {
        require(redeemRequests[msg.sender].amount == uint256(0));
        require(amount_ > 0, "amount should be bigger than 0");
        require(redeemRequests[msg.sender].until < block.number, "request is locked, not available");

        redeemRequests[msg.sender].isCompleted = false;
        redeemRequests[msg.sender].signature_count = uint256(0);
        redeemRequests[msg.sender].recipient = msg.sender;
        redeemRequests[msg.sender].amount = amount_ ;
        redeemRequests[msg.sender].until = block.number + LOCK_PERIOD;
        emit RedeemRequest(redeemRequests[msg.sender].recipient,redeemRequests[msg.sender].amount);
    }

    //function called by Validators using protocol to sign on RedeemRequests
    function sign(uint amount_, address payable recipient_) public isActive onlyValidator {
        require(!redeemRequests[recipient_].isCompleted, "redeem request is completed");
        require(redeemRequests[recipient_].amount == amount_,"redeem amount Compromised" );
        require(!redeemRequests[recipient_].votes[msg.sender]);
        // update votes
        redeemRequests[recipient_].votes[msg.sender] = true;
        redeemRequests[recipient_].signature_count += 1;
        // if threshold is reached then transfer the amount
        if (redeemRequests[recipient_].signature_count >= votingThreshold ) {
            (bool success, ) = redeemRequests[recipient_].recipient.call.value((redeemRequests[recipient_].amount))("");
            require(success, "Transfer failed.");
            redeemRequests[recipient_].amount = 0;
            redeemRequests[recipient_].isCompleted = true;
        }
        emit ValidatorSignedRedeem(recipient_, msg.sender, amount_);
    }

    function hasValidatorSigned(address recipient_) public view returns(bool) {
        return redeemRequests[recipient_].votes[msg.sender];
    }

    function verifyRedeem(address recipient_) public view returns(bool){
        return redeemRequests[recipient_].isCompleted;
    }

    function getSignatureCount(address recipient_) public view returns(uint256){
        return redeemRequests[recipient_].signature_count;
    }

    function getTotalEthBalance() public view returns(uint) {
        return address(this).balance ;
    }

    function getOLTEthAddress() public view returns(address){
        return address(this);
    }

    function addValidator(address v) internal {
        validators[v] = DEFAULT_VALIDATOR_POWER;
        numValidators += 1;
        emit AddValidator(v);
    }


}
//["0x14723A09ACff6D2A60DcdF7aA4AFf308FDDC160C","0x4B0897b0513fdC7C541B6d9D7E929C4e5364D2dB","0x583031D1113aD414F02576BD6afaBfb302140225","0xdD870fA1b7C4700F2BD7f44238821C26f7392148"]
