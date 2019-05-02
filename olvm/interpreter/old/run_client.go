package main

import (
	"github.com/Oneledger/protocol/olvm/interpreter/runner"
	"github.com/Oneledger/protocol/olvm/interpreter/vm"
)

var sourceCode = `var HellowWorldContract = function (context) {
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

HellowWorldContract.prototype.getWord = function () {
  return this.context.get('word');
}
module.Contract = HellowWorldContract;
`

func main() {
	request := runner.OLVMRequest{"0x0", "embed://", "", 0, sourceCode}

	reply, err := vm.AutoRun(&request)
	if err != nil {
		log.Fatal("Failed", err)
	}

	log.Info("get the result", "reply", reply)
}
