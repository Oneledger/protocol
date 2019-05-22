var GetContextValueContract = function (context) {
  this.context = context;
}

GetContextValueContract.prototype.default__ = function() {
  return this.context.get("test_key")
}
module.Contract = GetContextValueContract;
