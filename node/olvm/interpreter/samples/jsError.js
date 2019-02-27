var jsErrorContract = function (context) {
  this.context = context;
}


jsErrorContract.prototype.default__ = function () {
  return undefine_value.length
}
jsErrorContract.prototype.post = function () {
  var len = undefine_value.length
}
jsErrorContract.prototype.post.func_type = "write"

module.Contract = jsErrorContract;
