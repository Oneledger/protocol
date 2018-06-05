pragma solidity ^0.4.0;

contract HTLCChannel{
    using SafeMath for uint256;

    address public sender;
    address public receiver;
    uint256 public cap;
    uint256 public balance;
    uint256 public lockPeriod;
    uint256 public startFromTime;
    bytes32 public msgHash;

    //TODO: gas for release/rollback
    constructor(
        uint256 _lockPeriod,
        address _sender,
        address _receiver,
        uint256 _cap,
        bytes32 _msgHash
    ) public {
        require(_sender != address(0));
        require(_receiver != address(0));
        require(_lockPeriod >= 24 hours );
        require(_cap > 0);

        sender = _sender;
        receiver = _receiver;
        cap = _cap;
        balance = 0;
        lockPeriod = _lockPeriod;
        msgHash = _msgHash;

    }

    event Release(address sender, address receiver, uint256 value);
    event Rollback(address sender, address receiver, uint256 value);

    function() external payable{
        funds();
    }

    function funds() public payable {
        require(msg.sender == sender);
        require(msg.value == cap);
        require(balance == 0);
        balance = msg.value;
        startFromTime = now();
    }

    function validate(bytes32 msg) public returns (bool) {
        return sha3(msg) == msgHash;
    }

    function release(bytes32 msg) public returns (bool) {
        require(validate(msg));
        this.transfer(receiver, balance);
        balance = 0;
        emit Release(this, receiver, balance);
        //todo: selfdestruct()?
        return true;
    }

    function rollback(bytes32 msg) public returns (bool) {
        require((startFromTime + lockPeriod) > now());
        require(balance == cap);
        require(validate(msg));

        this.transfer(sender, balance);
        balance = 0;
        emit Rollback(this, sender, balance);
        //todo: selfdestruct()?
        return true;
    }

}
