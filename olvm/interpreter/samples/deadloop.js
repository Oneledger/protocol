var DeadLoopContract = function (context) {
  this.context = context;
}


DeadLoopContract.prototype.default__ = function () {
  while(true){
    console.log('loop...');
  }
}
module.Contract = DeadLoopContract;
