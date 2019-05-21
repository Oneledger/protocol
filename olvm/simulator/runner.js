var context = require('./context');
var Contract = require('./simple_contract');

var contract = new Contract(context);

contract.setWord('myaddress','hello,world')
var ret = contract.getWord('myaddress');


console.log(ret);
var list = context.getUpdateIndexList();
var storage = context.getStorage();
var transaction = {};
for (var i = 0; i< list.length; i ++) {
  var key = list[i];
  transaction[key] = storage[key];
}

console.log(JSON.stringify(transaction));
