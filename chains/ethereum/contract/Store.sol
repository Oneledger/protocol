pragma solidity >=0.4.25 <0.6.0;

contract Store {
  event ItemSet(bytes32 key, bytes32 value);

  uint public balance;
  string public version;
  mapping (bytes32 => bytes32) public items;
  mapping (address => uint) public balances;

  constructor(string memory _version) public {
    version = _version;
  }


  function payItem() public payable {
    require(msg.value >= 30);
    balances[msg.sender] = msg.value;
  }


  function setItem(bytes32 key, bytes32 value) external {
    items[key] = value;
    emit ItemSet(key, value);
  }

}

