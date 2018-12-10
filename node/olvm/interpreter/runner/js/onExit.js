var list = context.getUpdateIndexList();
var storage = context.getStorage();
var transaction = {};
for (var i = 0; i< list.length; i ++) {
  var key = list[i];
  transaction[key] = storage[key];
}
transaction.__from__ = __from__;
transaction.__olt__ = __olt__;
var out = JSON.stringify(transaction);
var ret = JSON.stringify(retValue);
