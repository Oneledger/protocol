var functionName = __callString__.substring(0, __callString__.indexOf("("));

var __func_type__ = contract[functionName].func_type;
var out = "__NO_ACTION__";
if (__func_type__ == "write") {
  out = "__WRITE__"
}

var ret = 0
