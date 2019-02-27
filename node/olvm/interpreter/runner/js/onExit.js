var functionName = __callString__.substring(0, __callString__.indexOf("("));

var __func_type__ = contract[functionName].func_type;

var list = context.getUpdateIndexList();
var storage = context.getStorage();
var out = "__NO_ACTION__";
if (__func_type__ == "write") {
  var transaction = {};
  for (var i = 0; i< list.length; i ++) {
    var key = list[i];
    transaction[key] = storage[key];
  }
  transaction.__from__ = __from__;
  transaction.__olt__ = __olt__;
  transaction.__runtime__ = "es5"
  transaction.__version__ = "0.5.4"
  transaction.__line_data__ =  context.getLineData();
  out = JSON.stringify(transaction);
}

var ret = JSON.stringify(retValue);
