pragma solidity >=0.5.0 <0.6.0;



/**
 * @dev Interface of the ERC20 standard as defined in the EIP. Does not include
 * the optional functions; to access them see {ERC20Detailed}.
 */
interface IERC20 {
    /**
     * @dev Returns the amount of tokens in existence.
     */
    function totalSupply() external view returns (uint256);

    /**
     * @dev Returns the amount of tokens owned by `account`.
     */
    function balanceOf(address account) external view returns (uint256);

    /**
     * @dev Moves `amount` tokens from the caller's account to `recipient`.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transfer(address recipient, uint256 amount) external returns (bool);

    /**
     * @dev Returns the remaining number of tokens that `spender` will be
     * allowed to spend on behalf of `owner` through {transferFrom}. This is
     * zero by default.
     *
     * This value changes when {approve} or {transferFrom} are called.
     */
    function allowance(address owner, address spender) external view returns (uint256);

    /**
     * @dev Sets `amount` as the allowance of `spender` over the caller's tokens.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * IMPORTANT: Beware that changing an allowance with this method brings the risk
     * that someone may use both the old and the new allowance by unfortunate
     * transaction ordering. One possible solution to mitigate this race
     * condition is to first reduce the spender's allowance to 0 and set the
     * desired value afterwards:
     * https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729
     *
     * Emits an {Approval} event.
     */
    function approve(address spender, uint256 amount) external returns (bool);

    /**
     * @dev Moves `amount` tokens from `sender` to `recipient` using the
     * allowance mechanism. `amount` is then deducted from the caller's
     * allowance.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);

    /**
     * @dev Emitted when `value` tokens are moved from one account (`from`) to
     * another (`to`).
     *
     * Note that `value` may be zero.
     */
    event Transfer(address indexed from, address indexed to, uint256 value);

    /**
     * @dev Emitted when the allowance of a `spender` for an `owner` is set by
     * a call to {approve}. `value` is the new allowance.
     */
    event Approval(address indexed owner, address indexed spender, uint256 value);
}



/**
 * @dev Wrappers over Solidity's arithmetic operations with added overflow
 * checks.
 *
 * Arithmetic operations in Solidity wrap on overflow. This can easily result
 * in bugs, because programmers usually assume that an overflow raises an
 * error, which is the standard behavior in high level programming languages.
 * `SafeMath` restores this intuition by reverting the transaction when an
 * operation overflows.
 *
 * Using this library instead of the unchecked operations eliminates an entire
 * class of bugs, so it's recommended to use it always.
 */
library SafeMath {
    /**
     * @dev Returns the addition of two unsigned integers, reverting on
     * overflow.
     *
     * Counterpart to Solidity's `+` operator.
     *
     * Requirements:
     * - Addition cannot overflow.
     */
    function add(uint256 a, uint256 b) internal pure returns (uint256) {
        uint256 c = a + b;
        require(c >= a, "SafeMath: addition overflow");

        return c;
    }

    /**
     * @dev Returns the subtraction of two unsigned integers, reverting on
     * overflow (when the result is negative).
     *
     * Counterpart to Solidity's `-` operator.
     *
     * Requirements:
     * - Subtraction cannot overflow.
     */
    function sub(uint256 a, uint256 b) internal pure returns (uint256) {
        return sub(a, b, "SafeMath: subtraction overflow");
    }

    /**
     * @dev Returns the subtraction of two unsigned integers, reverting with custom message on
     * overflow (when the result is negative).
     *
     * Counterpart to Solidity's `-` operator.
     *
     * Requirements:
     * - Subtraction cannot overflow.
     *
     * NOTE: This is a feature of the next version of OpenZeppelin Contracts.
     * @dev Get it via `npm install @openzeppelin/contracts@next`.
     */
    function sub(uint256 a, uint256 b, string memory errorMessage) internal pure returns (uint256) {
        require(b <= a, errorMessage);
        uint256 c = a - b;

        return c;
    }

    /**
     * @dev Returns the multiplication of two unsigned integers, reverting on
     * overflow.
     *
     * Counterpart to Solidity's `*` operator.
     *
     * Requirements:
     * - Multiplication cannot overflow.
     */
    function mul(uint256 a, uint256 b) internal pure returns (uint256) {
        // Gas optimization: this is cheaper than requiring 'a' not being zero, but the
        // benefit is lost if 'b' is also tested.
        // See: https://github.com/OpenZeppelin/openzeppelin-contracts/pull/522
        if (a == 0) {
            return 0;
        }

        uint256 c = a * b;
        require(c / a == b, "SafeMath: multiplication overflow");

        return c;
    }

    /**
     * @dev Returns the integer division of two unsigned integers. Reverts on
     * division by zero. The result is rounded towards zero.
     *
     * Counterpart to Solidity's `/` operator. Note: this function uses a
     * `revert` opcode (which leaves remaining gas untouched) while Solidity
     * uses an invalid opcode to revert (consuming all remaining gas).
     *
     * Requirements:
     * - The divisor cannot be zero.
     */
    function div(uint256 a, uint256 b) internal pure returns (uint256) {
        return div(a, b, "SafeMath: division by zero");
    }

    /**
     * @dev Returns the integer division of two unsigned integers. Reverts with custom message on
     * division by zero. The result is rounded towards zero.
     *
     * Counterpart to Solidity's `/` operator. Note: this function uses a
     * `revert` opcode (which leaves remaining gas untouched) while Solidity
     * uses an invalid opcode to revert (consuming all remaining gas).
     *
     * Requirements:
     * - The divisor cannot be zero.

     * NOTE: This is a feature of the next version of OpenZeppelin Contracts.
     * @dev Get it via `npm install @openzeppelin/contracts@next`.
     */
    function div(uint256 a, uint256 b, string memory errorMessage) internal pure returns (uint256) {
        // Solidity only automatically asserts when dividing by 0
        require(b > 0, errorMessage);
        uint256 c = a / b;
        // assert(a == b * c + a % b); // There is no case in which this doesn't hold

        return c;
    }

    /**
     * @dev Returns the remainder of dividing two unsigned integers. (unsigned integer modulo),
     * Reverts when dividing by zero.
     *
     * Counterpart to Solidity's `%` operator. This function uses a `revert`
     * opcode (which leaves remaining gas untouched) while Solidity uses an
     * invalid opcode to revert (consuming all remaining gas).
     *
     * Requirements:
     * - The divisor cannot be zero.
     */
    function mod(uint256 a, uint256 b) internal pure returns (uint256) {
        return mod(a, b, "SafeMath: modulo by zero");
    }

    /**
     * @dev Returns the remainder of dividing two unsigned integers. (unsigned integer modulo),
     * Reverts with custom message when dividing by zero.
     *
     * Counterpart to Solidity's `%` operator. This function uses a `revert`
     * opcode (which leaves remaining gas untouched) while Solidity uses an
     * invalid opcode to revert (consuming all remaining gas).
     *
     * Requirements:
     * - The divisor cannot be zero.
     *
     * NOTE: This is a feature of the next version of OpenZeppelin Contracts.
     * @dev Get it via `npm install @openzeppelin/contracts@next`.
     */
    function mod(uint256 a, uint256 b, string memory errorMessage) internal pure returns (uint256) {
        require(b != 0, errorMessage);
        return a % b;
    }
}

contract ERC20Basic {

    string public constant name = "TestToken";
    string public constant symbol = "TTC";
    uint8 public constant decimals = 18;


    event Approval(address indexed tokenOwner, address indexed spender, uint tokens);
    event Transfer(address indexed from, address indexed to, uint tokens);


    mapping(address => uint256) balances;

    mapping(address => mapping (address => uint256)) allowed;

    uint256 totalSupply_;

    using SafeMath for uint256;


    constructor(uint256 total) public {
        totalSupply_ = total;
        balances[msg.sender] = totalSupply_;
    }

    function totalSupply() public view returns (uint256) {
        return totalSupply_;
    }

    function balanceOf(address tokenOwner) public view returns (uint) {
        return balances[tokenOwner];
    }

    function transfer(address receiver, uint numTokens) public returns (bool) {
        require(numTokens <= balances[msg.sender]);
        balances[msg.sender] = balances[msg.sender].sub(numTokens);
        balances[receiver] = balances[receiver].add(numTokens);
        emit Transfer(msg.sender, receiver, numTokens);
        return true;
    }

    function approve(address delegate, uint numTokens) public returns (bool) {
        allowed[msg.sender][delegate] = numTokens;
        emit Approval(msg.sender, delegate, numTokens);
        return true;
    }

    function allowance(address owner, address delegate) public view returns (uint) {
        return allowed[owner][delegate];
    }

    function transferFrom(address owner, address buyer, uint numTokens) public returns (bool) {
        require(numTokens <= balances[owner]);
        require(numTokens <= allowed[owner][msg.sender]);

        balances[owner] = balances[owner].sub(numTokens);
        allowed[owner][msg.sender] = allowed[owner][msg.sender].sub(numTokens);
        balances[buyer] = balances[buyer].add(numTokens);
        emit Transfer(owner, buyer, numTokens);
        return true;
    }
}

contract LockRedeemERC {
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
    uint constant MIN_VALIDATORS = 1; // 4 for production
    uint256 LOCK_PERIOD = 28800;

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


    //string tokenName;
    // address tokenAddress; // Address of tokenName smartContract

    //Redeem
    //mapping (address => bool) public isSigned;
    struct RedeemTX {
        address tokenAddress;
        address payable recipient;
        mapping (address => bool) votes;
        uint256 amount;
        uint256 signature_count;
        bool isCompleted ;
        uint256 until;
    }
    mapping (address => RedeemTX) redeemRequests;


    event AddValidator(
        address indexed _address,
        int _power
    );

    event DeleteValidator(address indexed _address);

    event NewEpoch(uint epochHeight);

    event RedeemSuccessful(
        address indexed recepient,
        uint amount_trafered
    );

    event ValidatorSignedRedeem(
        address indexed recipient,
        address validator_addresss,
        uint256 amount
    );
    event RedeemRequest(
        address indexed recepient,
        uint256 amount_requested
    );



    event NewThreshold(uint _prevThreshold, uint _newThreshold);

    constructor(address[] memory initialValidators) public {
        // Require at least 4 validators
        require(initialValidators.length >= MIN_VALIDATORS, "insufficient validators passed to constructor");
        // require( tokenAddress_ != '0x0');

        // Add the initial validators
        for (uint i = 0; i < initialValidators.length; i++) {
            // Ensure these validators are unique
            address v = initialValidators[i];
            require(validators[v] == 0, "found non-unique validator in initialValidators");
            addValidator(v);
        }
        votingThreshold = (initialValidators.length * 2 / 3) + 1;
        // Set the initial epochBlockHeight
        //  declareNewEpoch(block.number);
    }

    function isValidator(address addr) public view returns(bool) {
        return validators[addr] > 0;
    }



    //function called by user and dcopy of trasaction sent to protocol
    function redeem(uint256 amount_,address tokenAddress_)  public  {
        require(redeemRequests[msg.sender].amount == uint256(0));
        require(amount_ > 0, "amount should be bigger than 0");
        //require(redeemRequests[msg.sender].isCompleted == true, "earlier redeem has not been executed yet");
        require(redeemRequests[msg.sender].until < block.number, "request is locked, not available");

        redeemRequests[msg.sender].isCompleted = false;
        redeemRequests[msg.sender].tokenAddress = tokenAddress_;
        redeemRequests[msg.sender].signature_count = uint256(0);
        redeemRequests[msg.sender].recipient = msg.sender;
        redeemRequests[msg.sender].amount = amount_ ;
        redeemRequests[msg.sender].until = block.number + LOCK_PERIOD;
        emit RedeemRequest(redeemRequests[msg.sender].recipient,redeemRequests[msg.sender].amount);
    }

    function sign(uint amount_, address payable recipient_) public  {
        require(isValidator(msg.sender),"validator not present in list");
        require(!redeemRequests[recipient_].isCompleted, "redeem request is completed");
        require(redeemRequests[recipient_].amount == amount_,"redeem amount Compromised" );
        require(!redeemRequests[recipient_].votes[msg.sender]);

        // update votes
        redeemRequests[recipient_].votes[msg.sender] = true;
        redeemRequests[recipient_].signature_count += 1;
        emit ValidatorSignedRedeem(recipient_, msg.sender, amount_);
    }

    //function called by user
    function executeredeem (address tokenAddress_) public  {
        require(redeemRequests[msg.sender].recipient == msg.sender,"Request not called from token owner");
        require(redeemRequests[msg.sender].signature_count >= votingThreshold ,"Not enough Validator votes");
        require(redeemRequests[msg.sender].tokenAddress == tokenAddress_,"Wrong token selected");
        require(redeemRequests[msg.sender].isCompleted != true,"Redeem has already been executed");
        IERC20 testToken  =   IERC20(tokenAddress_);
        testToken.transfer(redeemRequests[msg.sender].recipient,redeemRequests[msg.sender].amount);
        redeemRequests[msg.sender].amount = 0;
        redeemRequests[msg.sender].isCompleted = true;
        emit RedeemSuccessful(redeemRequests[msg.sender].recipient,redeemRequests[msg.sender].amount);
    }

    function getTotalErcBalance(address tokenAddress_) public view returns (uint) {
        IERC20 testToken  =   IERC20(tokenAddress_);
        return testToken.balanceOf(address(this));
    }

    function getOLTErcAddress() public view returns(address){
        return address(this);
    }

    function hasValidatorSigned(address recipient_) public view returns(bool) {
        return redeemRequests[recipient_].votes[msg.sender];
    }

    function verifyRedeem(address recipient_) public view returns(bool){
        return redeemRequests[recipient_].signature_count >= votingThreshold;
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