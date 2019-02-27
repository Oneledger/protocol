var CounterContract = function (context) {
  this.context = context;
}

CounterContract.prototype.default__ = function() {
  var counter = this.context.get('counter');
  if(counter === undefined) {
    return 0
  }else{
    return counter
  }
}

CounterContract.prototype.increase = function() {
  var counter = this.context.get('counter');
  if(counter === undefined) {
    counter  = 0;
  }
  counter = counter + 1;
  this.context.set('counter', counter);
  return counter;
}
CounterContract.prototype.increase.func_type = "write";

module.Contract = CounterContract;
