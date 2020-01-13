pragma solidity >=0.5.0 <0.6.0;

library SafeMath {
function add(uint256 a, uint256 b) internal pure returns (uint256) {
uint256 c = a + b;
require(c >= a, "SafeMath: addition overflow");

return c;
}

function sub(uint256 a, uint256 b) internal pure returns (uint256) {
return sub(a, b, "SafeMath: subtraction overflow");
}

function sub(uint256 a, uint256 b, string memory errorMessage) internal pure returns (uint256) {
require(b <= a, errorMessage);
uint256 c = a - b;

return c;
}

function mul(uint256 a, uint256 b) internal pure returns (uint256) {

if (a == 0) {
return 0;
}

uint256 c = a * b;
require(c / a == b, "SafeMath: multiplication overflow");

return c;
}


function div(uint256 a, uint256 b) internal pure returns (uint256) {
return div(a, b, "SafeMath: division by zero");
}

function div(uint256 a, uint256 b, string memory errorMessage) internal pure returns (uint256) {
// Solidity only automatically asserts when dividing by 0
require(b > 0, errorMessage);
uint256 c = a / b;

return c;
}


function mod(uint256 a, uint256 b) internal pure returns (uint256) {
return mod(a, b, "SafeMath: modulo by zero");
}

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

// ERC20 public ERC20Interface;
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
votingThreshold = 1; // Take as input
// Set the initial epochBlockHeight
//declareNewEpoch(block.number);
}

function isValidator(address addr) public view returns(bool) {
return validators[addr] > 0;
}



//function called by user andcopy of trasaction sent to protocol
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
require(redeemRequests[msg.sender].recipient == msg.sender);
require(redeemRequests[msg.sender].signature_count >= votingThreshold);
require(redeemRequests[msg.sender].tokenAddress == tokenAddress_);
ERC20Basic tesToken  =  ERC20Basic(tokenAddress_);
tesToken.transfer(redeemRequests[msg.sender].recipient,redeemRequests[msg.sender].amount);
redeemRequests[msg.sender].amount = 0;
redeemRequests[msg.sender].isCompleted = true;
emit RedeemSuccessful(redeemRequests[msg.sender].recipient,redeemRequests[msg.sender].amount);
}

function getTotalErcBalance(address tokenAddress_) public view returns (uint) {
ERC20Basic tesToken =  ERC20Basic(tokenAddress_);
return tesToken.balanceOf(address(this));
}

function getOLTErcAddress() public view returns(address){
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
