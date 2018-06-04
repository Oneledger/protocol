pragma solidity ^0.4.0;

contract HTLCChannel {
    using SafeMath for uint256;

    address public sender;
    address public receiver;
    uint256 public balance;
    uint256 public cap;
    uint256 public lockPeriod;
    uint256 public startFromTime;
    //if we need lock based on height
    //uint256 public lockHeight;
    //uint256 public startFromHeight;

    //receiver pay gas to create this channel
    constructor(
        uint256 _cap,
        uint256 _lockPeriod,
        //uint256 _lockHeight,
        address _sender,
        address _receiver
    ) public {
        require(_sender != address(0));
        require(_receiver != address(0));
        require(lockPeriod > 48 hours );
        //require(lockHeight > 2800);
        require(cap > 0);

        sender = _sender;
        receiver = _receiver;

        cap = _cap;
        balance = 0;
        lockPeriod = _lockPeriod;
        //lockHeight = _lockHeight;
    }

    function() external payable{

    }

    function funds() public payable{
        require(msg.sender == sender);
        require(balance == 0);
        require(msg.value == cap);

        startFromTime = now();
        balance = msg.value;

    }

    function release() public returns (bool) {
        require(msg.sender = sender);
        require(balance == cap);
        balance = 0;
        emit Transfer(this, receiver, balance);
        //todo: selfdestruct()?
        return true;
    }

    function rollback() public returns (bool) {
        require(msg.sender == receiver || ((lockPeriod + startFromTime) > now() && msg.sender == sender) );
        require(balance == cap);
        balance = 0;
        emit Transfer(this, sender, balance);
        //todo: selfdestruct()?
        return true;
    }

}
