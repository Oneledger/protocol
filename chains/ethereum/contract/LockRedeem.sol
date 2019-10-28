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
    uint constant MIN_VALIDATORS = 4;

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

    mapping (address => int) public validators;
    //mapping (address => uint) public balances;

    event AddValidator(
        address indexed _address,
        int _power
    );

    //Redeem
   // mapping (address => bool) public isSigned;
    struct RedeemTX {
        address payable recipient;
        uint amount;
        uint signature_count;
    }
    mapping (uint => RedeemTX) redeemRequests;

    event RedeemSuccessful(
        address indexed recepient,
        uint amount_trafered
    );

    event ValidatorSignedRedeem(
        address indexed validator_addresss
    );

    event DeleteValidator(address indexed _address);

    event NewEpoch(uint epochHeight);

    event Lock(
        address sender,
        uint amunt_received,
        uint updated_balance
    );

    event NewThreshold(uint _prevThreshold, uint _newThreshold);

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

        // Set the initial epochBlockHeight
        declareNewEpoch(block.number);
    }

    function isValidator(address addr) public view returns(bool) {
        return validators[addr] > 0;
    }

    function lock() payable public {                                        
        require(msg.value >= 0, "Must pay a balance more than 0");
        //balances[oltAddress] += msg.value;
        emit Lock(msg.sender,msg.value,getTotalEthBalance());
    }

    // getLockedBalance returns the amount of ether locked by a specific account
//    function getLockedBalance(address addr) public view returns (uint) {
//        return balances[addr];
//    }
    function sign(uint redeemID_,uint amount_,address payable recipient_) public  {
        require(isValidator(msg.sender),"validator not pressent in list");
        if(redeemRequests[redeemID_].signature_count > 0 )
        {
            require(redeemRequests[redeemID_].amount == amount_,"ValidatorCompromised" );
            require(redeemRequests[redeemID_].recipient == recipient_, "ValidatorCompromised");
            redeemRequests[redeemID_].signature_count = redeemRequests[redeemID_].signature_count + 1;
        }
        else
        {
            redeemRequests[redeemID_].amount = amount_;
            redeemRequests[redeemID_].recipient = recipient_;
            redeemRequests[redeemID_].signature_count = 1;
        }
        emit ValidatorSignedRedeem(msg.sender);
    }
    function redeem (uint redeemID_)  public  {
        require(redeemRequests[redeemID_].recipient == msg.sender);
        require(redeemRequests[redeemID_].signature_count >= votingThreshold);
        redeemRequests[redeemID_].recipient.transfer(redeemRequests[redeemID_].amount);
        emit RedeemSuccessful(redeemRequests[redeemID_].recipient,redeemRequests[redeemID_].amount);
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

    function declareNewEpoch(uint nextEpochHeight) internal onlyValidator {
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
        validators[v] = DEFAULT_VALIDATOR_POWER;
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
