pragma solidity >=0.5.0 <0.6.0;

contract LockRedeem {
    //Flag to pause and unpause contract
    bool ACTIVE = false;
    // numValidators holds the total number of validators
    uint public numValidators;
    // mapping to store the validators to the ether they  have earne from signing trasactions
    struct ValidatorDetails {
        uint votingPower;
        uint validatorFee;
    }
    mapping (address => ValidatorDetails) public validators;
    address [] initialValidatorList;
    // Require this amount of signatures to push a proposal through,Set only in constructor
    uint votingThreshold;
    uint activeThreshold;
    //the base cost of a transaction (21000 gas)
    //the cost of a contract deployment (32000 gas) [Not relevant for sign trasaction]
    //the cost for every zero byte of data or code for a transaction.
    //the cost of every non-zero byte of data or code for a transaction.
    uint trasactionCost = 23192 ;
    uint additionalGasCost = 37692 ;
    uint redeem_gas_charge = 1000000000000000000;
    // numValidators holds the total number of validators
    uint public migrationSignatures;
    // mapping to store the validators to a bool ,which accounts to the fact they have signed.
    mapping (address => bool) public migrationSigners;
    mapping (address => uint) migrationAddressMap;
    address [] migrationAddressList;

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
        uint redeemFee;
    }
    mapping (address => RedeemTX) redeemRequests;

    event RedeemRequest(
        address indexed recepient,
        uint256 amount_requested
    );

    event ValidatorSignedRedeem(
        address indexed recipient,
        address validator_addresss,
        uint256 amount,
        uint gasReturned
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
        require(validators[msg.sender].votingPower >0 ,"validator not present in List");
        _; // Continues control flow after this is validates
    }

    function isValidator(address v) public view returns (bool) {
        return validators[v].votingPower > 0;
    }
//    function migrate (address newSmartContractAddress) public onlyValidator {
//        require(migrationSigners[msg.sender]==false,"Validator Signed already");
//        migrationSigners[msg.sender] = true ;
//        //MigratefromOl can be used to migrate validator to new contract address.
//        (bool status,) = newSmartContractAddress.call(abi.encodePacked(bytes4(keccak256("MigrateFromOld()"))));
//        require(status,"Unable to Migrate Validator new Smart contract");
//        emit ValidatorMigrated(msg.sender,newSmartContractAddress);
//        //migrationSignatures = migrationSignatures + 1;
//        if(migrationAddressMap[newSmartContractAddress]==0)
//        {
//            migrationAddressList.push(newSmartContractAddress);
//        }
//        migrationAddressMap[newSmartContractAddress] += 1;
//
//        if (migrationAddressMap[newSmartContractAddress] > activeThreshold) {
//            ACTIVE = false ;
//        }
//        if (migrationAddressMap[newSmartContractAddress] > votingThreshold) {
//            (bool success, ) = maxVotedAddress.call.value(address(this).balance)("");
//        }

        // Global flag ,needs to be set only once
//        if (migrationSignatures == activeThreshold) {
//            ACTIVE = false ;
//        }
//
//
//        // Trasfer needs to be done only once
//        if (migrationSignatures == votingThreshold) {
//            uint voteCount = 0 ;
//            address maxVotedAddress ;
//            for(uint i=0;i< migrationAddressList.length;i++){
//                if (migrationAddressMap[migrationAddressList[i]] > voteCount){
//                    voteCount = migrationAddressMap[migrationAddressList[i]];
//                    maxVotedAddress = migrationAddressList[i];
//                }
//            }
//            if (voteCount >= votingThreshold) {
//            (bool success, ) = maxVotedAddress.call.value(address(this).balance)("");
//            }
//            require(success, "Transfer failed");
//        }
//    }

    // function called by user
    function lock() payable public isActive {
        require(msg.value >= 0, "Must pay a balance more than 0");
        emit Lock(msg.sender,msg.value);
    }

    // function called by user
    function redeem(uint256 amount_)  payable public isActive {
        require(redeemRequests[msg.sender].amount == uint256(0));
        require(msg.value >= redeem_gas_charge,"Redeem Fees not supplied");
        require(amount_ > 0, "amount should be bigger than 0");
        require(redeemRequests[msg.sender].until < block.number, "request is locked, not available");

        redeemRequests[msg.sender].isCompleted = false;
        redeemRequests[msg.sender].signature_count = uint256(0);
        redeemRequests[msg.sender].recipient = msg.sender;
        redeemRequests[msg.sender].amount = amount_ ;
        redeemRequests[msg.sender].redeemFee += msg.value;
        redeemRequests[msg.sender].until = block.number + LOCK_PERIOD;

        emit RedeemRequest(redeemRequests[msg.sender].recipient,redeemRequests[msg.sender].amount);
    }

    //function called by Validators using protocol to sign on RedeemRequests
    function sign(uint amount_, address payable recipient_) public isActive onlyValidator {
        uint startGas = gasleft();
        require(!redeemRequests[recipient_].isCompleted, "redeem request is completed");
        require(redeemRequests[recipient_].amount == amount_,"redeem amount is different" );
        require(!redeemRequests[recipient_].votes[msg.sender]);
        require(redeemRequests[recipient_].until > block.number, "request has expired");
        // update votes
        redeemRequests[recipient_].votes[msg.sender] = true;
        redeemRequests[recipient_].signature_count += 1;
        //validators[msg.sender].validatorFee = validators[msg.sender].validatorFee + tx.gasprice * avg_gas;

        // if threshold is reached then transfer the amount
        if (redeemRequests[recipient_].signature_count >= votingThreshold ) {
            (bool success, ) = redeemRequests[recipient_].recipient.call.value((redeemRequests[recipient_].amount))("");
            require(success, "Transfer failed.");
            redeemRequests[recipient_].amount = 0;
            redeemRequests[recipient_].isCompleted = true;
        }

        // Trasaction Cost is a an average trasaction cost .
        // additionalGasCost is the extra gas used by the lines of code after gas calculation is done.
        // This is an approximate calcualtion and actuall cost might vary sightly .
        uint gasUsed = startGas - gasleft() + trasactionCost + additionalGasCost ;
        validators[msg.sender].validatorFee = gasUsed;
        uint gasFee = gasUsed * tx.gasprice;
        (bool success, ) = msg.sender.call.value(gasFee)("");
        require(success, "Transfer back to validator failed");
        redeemRequests[recipient_].redeemFee -= gasFee;
        emit ValidatorSignedRedeem(recipient_, msg.sender, amount_,gasUsed);
    }

    function validatorFee() public view isActive onlyValidator returns(uint) {
        return validators[msg.sender].validatorFee;
    }

    function collectUserFee() public isActive{
        require(redeemRequests[msg.sender].isCompleted = true, "request signing is still in progress");
        (bool success, ) = msg.sender.call.value(redeemRequests[msg.sender].redeemFee)("");
        require(success, "Transfer failed.");

    }

    // Function used by protocol and wallet
    function isredeemAvailable (address recepient_) public view returns (bool){
        return redeemRequests[recepient_].until == 0 || redeemRequests[recepient_].until < block.number;
    }
    // Function used by protocol
    function hasValidatorSigned(address recipient_) public view returns(bool) {
        return redeemRequests[recipient_].votes[msg.sender];
    }
    // Function used by protocol
    function verifyRedeem(address recipient_) public view returns(bool){
        return redeemRequests[recipient_].isCompleted;
    }
    // Function used by protocol
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
        validators[v].votingPower = DEFAULT_VALIDATOR_POWER;
        validators[v].validatorFee = 0;
        numValidators += 1;
        emit AddValidator(v);
    }


}
