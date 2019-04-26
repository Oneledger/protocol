
var SimpleContract = function (context) {
  this.context = context;
}

SimpleContract.prototype.setWord = function(word) {
  this.context.set('word', word);
}

SimpleContract.prototype.getWord = function () {
  return this.context.get('word');
}

module.exports =  SimpleContract;
