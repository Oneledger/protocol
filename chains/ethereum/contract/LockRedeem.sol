pragma solidity >=0.5.0 <0.6.0;

contract LockRedeem {
    // Proposal represents the set of validators who have voted for a particular
    // issue (proposeAddValidator, proposeRemoveValidator, ...).
    // Only people who have vote
    struct Proposal {
        uint voteCount;
        mapping (address => bool) voters;
    }


    // numValidators holds the total number of validators
    uint public numValidators;
    // Require this amount of signatures to push a proposal through
    uint public votingThreshold;

    // Default Voting power should be updated at one point
    int constant DEFAULT_VALIDATOR_POWER = 50;
    uint constant MIN_VALIDATORS = 0;


    //Approx
    // 270 blocks per hour
    // 5 blocks per min
    // 28800 old value
    uint256 LOCK_PERIOD;

    // This is the height at which the current epoch was started
    uint public epochBlockHeight;


    // Mapping of validators in the set, from address to its voting power
    // These are epoch-dependent proposals
    mapping (address => Proposal) public addValidatorProposals;
    mapping (address => Proposal) public removeValidatorProposals;
    mapping (uint => Proposal) public newThresholdProposals;

    // Keep track of every validator to add. When the new epoch begins, add every epoch validator
    // and remove every validator from this array. It is emptied on every NewEpoch.
    address[] validatorsToAdd;
    address[] validatorsToRemove;

    mapping (address => uint8) public validators;
    address[] validatorList;
    event AddValidator(
        address indexed _address,
        uint8 index
    );

    struct RedeemTX {
        address payable recipient;
        uint256 votes;
        uint256 amount;
        uint256 signature_count;
        bool isProcessing ;
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

    event DeleteValidator(address indexed _address);

    event NewEpoch(uint epochHeight);

    event Lock(
        address sender,
        uint256 amount_received
    //uint updated_balance,
    // string oltEthAdress_User
    );

    event NewThreshold(uint _prevThreshold, uint _newThreshold);

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
        votingThreshold = (initialValidators.length * 2 / 3) + 1;
        // Set the initial epochBlockHeight
        declareNewEpoch(block.number);
        LOCK_PERIOD = _lock_period;
    }

    function isValidator(address addr) public view returns(bool) {
        return validators[addr] > 0;
    }
    // function called by user
    function lock() payable public {
        require(msg.value >= 0, "Must pay a balance more than 0");
        emit Lock(msg.sender,msg.value);
    }
    function isredeemAvailable (address recepient_) public view returns (bool){
        return redeemRequests[recepient_].until < block.number;
    }

    // function called by user
    function redeem(uint256 amount_)  public  {
        require(isredeemAvailable(msg.sender) ,"redeem to this address is not available yet");
        require(amount_ > 0, "amount should be bigger than 0");

        // redeemRequests[msg.sender].isProcessing = true;
        redeemRequests[msg.sender].signature_count = uint256(0);
        redeemRequests[msg.sender].recipient = msg.sender;
        redeemRequests[msg.sender].amount = amount_ ;
        redeemRequests[msg.sender].until = block.number + LOCK_PERIOD;
        redeemRequests[msg.sender].votes = 0;
        emit RedeemRequest(redeemRequests[msg.sender].recipient,redeemRequests[msg.sender].amount);
    }

    //function called by go
    // todo: validator sign and if enough vote, transfer directly
    function sign(uint amount_, address payable recipient_) public  {
        require(isValidator(msg.sender),"validator not present in list");
        require(redeemRequests[recipient_].until > block.number, "redeem request is not available");
        require(redeemRequests[recipient_].amount == amount_,"redeem amount is different" );
        require((redeemRequests[recipient_].votes >> validators[msg.sender])% 2 == uint256(0), "validator already vote");

        // update votes
        redeemRequests[recipient_].votes = redeemRequests[recipient_].votes + (uint256(1) << validators[msg.sender]);
        redeemRequests[recipient_].signature_count += 1;

        // if reach threshold, transfer
        if (redeemRequests[recipient_].signature_count >= votingThreshold ) {
            redeemRequests[recipient_].recipient.transfer(redeemRequests[recipient_].amount);
            redeemRequests[recipient_].amount = 0;
            // redeemRequests[recipient_].isProcessing = false;
            redeemRequests[recipient_].until = block.number;
            // for (uint i =0 ;i < validatorList.length ;i++)
            // {
            //     redeemRequests[recipient_].votes[validatorList[i]] = false;
            // }
            redeemRequests[recipient_].votes = 0;
        }
        emit ValidatorSignedRedeem(recipient_, msg.sender, amount_);
    }

    function hasValidatorSigned(address recipient_) public view returns(bool) {
        return redeemRequests[recipient_].votes >> validators[msg.sender] % 2 == uint256(1);
    }

    //BeforeRedeem : -1  (amount == 0 and until = 0)
    //Ongoing : 0  (amount > 0 and until > block.number)
    //Success : 1  (amount = 0 and until < block.number)
    //Expired : 2  (amount > 0 and until < block.number)
    function verifyRedeem(address recipient_) public view returns(int8){
        return (redeemRequests[recipient_].until > 0 && redeemRequests[recipient_].until < block.number) ? (redeemRequests[recipient_].amount > 0? int8(2) : int8(1) ): (redeemRequests[recipient_].amount == 0 ? int8(-1) : int8(0));
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

    // Proposals
    function proposeAddValidator(address v) public onlyValidator {
        Proposal storage proposal = addValidatorProposals[v];
        require(!proposal.voters[msg.sender], "sender has already voted to add this address");
        // Mark this voter as added and increment the vote count
        proposal.voters[msg.sender] = true;
        proposal.voteCount += 1;

        addValidator(v);
    }

    function proposeRemoveValidator(address v) public view onlyValidator {
        Proposal storage proposal = removeValidatorProposals[v];
        require(!proposal.voters[msg.sender], "sender has already voted to add this to proposal");
    }

    function proposeNewThreshold(uint threshold) public onlyValidator {
        require(threshold < numValidators, "New thresholds (m) must be less than the number of validators (n)");
        Proposal storage proposal = newThresholdProposals[threshold];
        require(!proposal.voters[msg.sender], "sender has already voted for this proposal");
        proposal.voters[msg.sender] = true;
        proposal.voteCount += 1;

    }

    function declareNewEpoch(uint nextEpochHeight) internal {
        epochBlockHeight = nextEpochHeight;
        emit NewEpoch(epochBlockHeight);

        // Add validators and remove them from the set.
        for (uint i = 0; i < validatorsToAdd.length; i++) {
            address v = validatorsToAdd[i];
            addValidator(v);
            delete addValidatorProposals[v];
        }

        delete validatorsToAdd;

        // Delete validators from the proposal mappings
        for (uint i = 0; i < validatorsToRemove.length; i++) {
            address v = validatorsToRemove[i];
            removeValidator(v);
            delete removeValidatorProposals[v];
        }
        delete validatorsToRemove;
    }

    modifier onlyValidator() {
        require(validators[msg.sender] > 0);
        _; // Continues control flow after this is validates
    }

    // Adds a validator to our current store
    function addValidator(address v) internal {

        uint8 index = uint8(validatorList.push(v));
        validators[v] = index;
        numValidators += 1;
        emit AddValidator(v, validators[v]);
    }

    // Deletes a validator from our store
    function removeValidator(address v) internal {
        delete validators[v];
        numValidators -= 1;
        emit DeleteValidator(v);
    }
}
