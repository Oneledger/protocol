pragma solidity >=0.5.0 <0.6.0;

contract LockRedeem {
    // numValidators holds the total number of validators
    uint totalWitness;
    uint activeWitness;
    uint min_stake;
    // Default Voting power should be updated at one point
    uint256 constant DEFAULT_VALIDATOR_POWER = 50;
    uint constant MIN_WITNESSES = 0;
    uint256 LOCK_PERIOD = 28800;
    bool ACTIVE = false;


    // Keep track of every validator to add. When the new epoch begins, add every epoch validator
    // and remove every validator from this array. It is emptied on every NewEpoch.

    mapping (address => uint256) public witnesses;
    uint totalWitnessStake = 0  ;
    struct RedeemTX {
        address payable recipient;
        mapping (address => bool) votes;
        uint256 amount;
        uint256 signature_count;
        bool isCompleted ;
        uint256 until;
    }
    mapping (address => RedeemTX) redeemRequests;

    mapping (address => uint) migrationCount;
    address [] migrationAddress;

    event RedeemRequest(
        address indexed recepient,
        uint256 amount_requested
    );
    event WitnessSignedRedeem(
        address indexed recipient,
        address validator_addresss,
        uint256 amount
    );
    event Lock(
        address sender,
        uint256 amount_received
    );

    event AddedWitness(
        address indexed _address,
        uint256 stake
    );

    event DeletedWitness(
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


    constructor(uint _witness_count,uint _min_stake) public {
        // Require at least 4 validators
        require(_witness_count >= MIN_WITNESSES, "insufficient witnesses to create contract");
        totalWitness = _witness_count;
        min_stake = _min_stake;
    }

    function migrate (address newSmartContractAddress) public {
        require(isWitness(msg.sender),"Witness not present in list");
        uint transfer_amount = witnesses[msg.sender];
        activeWitness = activeWitness - 1;
        witnesses[msg.sender] = 0 ;

        newSmartContractAddress.delegatecall(abi.encodePacked(bytes4(keccak256("someExternalFunction()"))));
        (bool success, ) = newSmartContractAddress.call.value(transfer_amount)("");
        require(success, "Transfer failed.");
        totalWitnessStake = totalWitnessStake - transfer_amount ;
        if (activeWitness < (totalWitness * 2 / 3) + 1){
            ACTIVE = false ;
        }
        if (activeWitness < (totalWitness * 1 / 3) + 1){
            uint userLockedAmount = address(this).balance - totalWitnessStake;
            uint voteCount ;
            address maxVotedAddress ;
            for(uint i=0;i<migrationAddress.length;i++){
                if (migrationCount[migrationAddress[i]] > voteCount){
                    voteCount = migrationCount[migrationAddress[i]];
                    maxVotedAddress =migrationAddress[i];
                }
            }
            (bool userEthsuccess, ) = maxVotedAddress.call.value(userLockedAmount)("");
            require(userEthsuccess, "Transfer failed.");
        }
    }


    function isWitness(address addr) public view returns(bool) {
        return witnesses[addr] > 0;
    }
    // function called by user
    function lock() payable public isActive {
        require(msg.value >= 0, "Must pay a balance more than 0");
        emit Lock(msg.sender,msg.value);
    }
    function isredeemAvailable (address recepient_) public view returns (bool)  {
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

    //function called by go
    // todo: validator sign and if enough vote, transfer directly
    function sign(uint amount_, address payable recipient_) public isActive {
        require(isWitness(msg.sender),"Witness not present in list");
        require(!redeemRequests[recipient_].isCompleted, "redeem request is completed");
        require(redeemRequests[recipient_].amount == amount_,"redeem amount Compromised" );
        require(!redeemRequests[recipient_].votes[msg.sender]);

        // update votes
        redeemRequests[recipient_].votes[msg.sender] = true;
        redeemRequests[recipient_].signature_count += 1;

        // if reach threshold, transfer
        if (redeemRequests[recipient_].signature_count >=  _getThreshold() ) {
            //  redeemRequests[recipient_].recipient.transfer(redeemRequests[recipient_].amount);
            (bool success, ) = redeemRequests[recipient_].recipient.call.value((redeemRequests[recipient_].amount))("");
            require(success, "Transfer failed.");
            redeemRequests[recipient_].amount = 0;
            redeemRequests[recipient_].isCompleted = true;
        }
        emit WitnessSignedRedeem(recipient_, msg.sender, amount_);
    }

    function hasWitnessSigned(address recipient_) public view returns(bool) {
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

    function _getThreshold() internal view returns (uint){
        return ((totalWitness * 2 / 3) + 1);
    }
    // Adds a Witness to our current store , smart contract address given to Validators offline
    function AddWitness() public payable {
        require(msg.value >= min_stake,"Not Enough Stake Provided");
        require(witnesses[msg.sender] == 0,"Witness already present in list");
        witnesses[msg.sender] = msg.value;
        totalWitnessStake = totalWitnessStake + msg.value;
        activeWitness = activeWitness + 1;
        if(activeWitness==totalWitness){
            ACTIVE = true ;
        }
        emit AddedWitness(msg.sender, witnesses[msg.sender]);
    }

    // Deletes a Witness from our store
    function RemoveWitness() public {
        require(witnesses[msg.sender] != 0,"Witness not present in list");
        delete witnesses[msg.sender];
        totalWitnessStake = totalWitnessStake -1 ;
        activeWitness = activeWitness - 1;
        emit DeletedWitness(msg.sender);
    }

}

