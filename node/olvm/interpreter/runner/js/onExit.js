var list = context.getUpdateIndexList();
var storage = context.getStorage();
var transaction = {};
for (var i = 0; i< list.length; i ++) {
  var key = list[i];
  transaction[key] = storage[key];
}
transaction.__from__ = __from__;
transaction.__olt__ = __olt__;
transaction.__runtime__ = "es5"
transaction.__version__ = "0.5.2"
var out = JSON.stringify(transaction);
var ret = JSON.stringify(retValue);
