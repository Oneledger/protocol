pragma solidity ^0.4.0;

//import "openzeppelin-solidity/contracts/math/SafeMath.sol";

contract HTLC {
    //using SafeMath for uint256;

    address public sender;
    address public receiver;
    uint256 public balance;
    uint256 public lockPeriod;
    uint256 public startFromTime;
    bytes32 public scrHash;
    bytes32 private scr;

    constructor(
        address _sender
    ) public {
        require(_sender != address(0));
        sender = _sender;
    }

    event Release(address sender, address receiver, uint256 value);
    event Rollback(address sender, address receiver, uint256 value);

    function funds() public payable {
        require(msg.sender == sender);
        require(msg.value > 0);
        require(balance == 0);
        balance = msg.value;
    }

    function setup(uint256 _lockPeriod, address _receiver, bytes32 _scrHash) public returns (bool) {
        require(msg.sender == sender);
        require(_receiver != address(0));
        require(_lockPeriod >= 24 hours );
        require(balance > 0);

        receiver = _receiver;
        lockPeriod = _lockPeriod;
        scrHash = _scrHash;

        startFromTime = now;

    }

    function audit(address receiver_, uint256 balance_, bytes32 scrHash_) public view returns (bool) {
        require(receiver == receiver_);
        require(balance == balance_);
        require(lockPeriod+startFromTime > now + 12 hours);
        require(scrHash == scrHash_);
        return true;
    }

    function validate(bytes32 scr_) private view returns (bool) {
        return keccak256(abi.encodePacked(scr_)) == scrHash;
    }

    function redeem(bytes32 scr_) public returns (bool) {
        require(validate(scr_));
        scr = scr_;
        address(receiver).transfer(balance);
        balance = 0;
        emit Release(this, receiver, balance);
        return true;
    }

    function refund(bytes32 scr_) public returns (bool) {
        require((startFromTime + lockPeriod) > now);
        require(validate(scr_));

        address(sender).transfer(balance);
        balance = 0;
        emit Rollback(this, sender, balance);
        return true;
    }

    function extractMsg() public view returns (bytes32) {
        require(msg.sender == receiver);
        require(validate(scr));
        return scr;
    }

}
