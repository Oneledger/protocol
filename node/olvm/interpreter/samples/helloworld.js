var HellowWorldContract = function (context) {
  this.context = context;
}

HellowWorldContract.prototype.default__ = function() {
  this.context.set('word', 'hello, by default');
  return 'hello, by default';
}

HellowWorldContract.prototype.setWord = function(word) {
  this.context.set('word', word);
  return word;
}
HellowWorldContract.prototype.setWord.func_type = "write";

HellowWorldContract.prototype.getWord = function () {
  return this.context.get('word');
}
module.Contract = HellowWorldContract;
