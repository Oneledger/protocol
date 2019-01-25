var CounterContract = function (context) {
  this.context = context;
}

CounterContract.prototype.default__ = function() {
  return this.context.get('counter');
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

module.Contract = CounterContract;
