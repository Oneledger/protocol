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
    bytes private scr;

    constructor(
        address _sender

    ) public {
        require(_sender != address(0));
        sender = _sender;
        startFromTime = now;
        lockPeriod = 0;
    }

    event Release(address sender, address receiver, uint256 value);
    event Rollback(address sender, address receiver, uint256 value);

    function funds(uint256 _lockPeriod, address _receiver, bytes32 _scrHash) public payable {
        require(msg.sender == sender);
        require(msg.value > 0);
        require(balance == 0);
        require(_receiver != address(0));
        require(_lockPeriod >= 4 minutes );
        require(startFromTime + lockPeriod < now );

        balance = msg.value;
        receiver = _receiver;
        lockPeriod = _lockPeriod;
        scrHash = _scrHash;

        startFromTime = now;
    }


    function audit(address receiver_, uint256 balance_, bytes32 scrHash_) public view returns (bool) {
        require(receiver == receiver_);
        require(balance == balance_);
        require(lockPeriod+startFromTime > now + 2 minutes);
        require(scrHash == scrHash_);
        return true;
    }

    function validate(bytes scr_) public view {
        require(sha256((abi.encodePacked(scr_))) == scrHash);
        require(balance > 0);
    }

    function redeem(bytes scr_) public returns (bool) {
        validate(scr_);
        scr = scr_;
        address(receiver).transfer(balance);
        balance = 0;
        emit Release(this, receiver, balance);
        return true;
    }

    function refund() public returns (bool) {
        require((startFromTime + lockPeriod) > now);
        require(msg.sender == sender);

        address(sender).transfer(balance);
        balance = 0;
        emit Rollback(this, sender, balance);
        return true;
    }

    function extractMsg() public view returns (bytes) {
        return scr;
    }

}
