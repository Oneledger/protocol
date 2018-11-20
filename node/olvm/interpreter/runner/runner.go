package runner

import (
	"errors"
	"fmt"
	"github.com/robertkrimen/otto"
)

func (runner Runner) exec(callString string) (string, string) {
	_, error := runner.vm.Run(`
    var contract = new module.Contract(context);
    var retValue = contract.` + callString)
	if error != nil {
		panic(error)
	}
	runner.vm.Run(`
    var list = context.getUpdateIndexList();
    var storage = context.getStorage();
    var transaction = {};
    for (var i = 0; i< list.length; i ++) {
      var key = list[i];
      transaction[key] = storage[key];
    }
    transaction.__from__ = __from__;
    transaction.__olt__ = __olt__;
    `)
	runner.vm.Run(`
    transaction = JSON.stringify(transaction);
    retValue = JSON.stringify(retValue);
    `)
	output := ""
	returnValue := ""

	if value, err := runner.vm.Get("transaction"); err == nil {
		output, _ = value.ToString()
	}

	if value, err := runner.vm.Get("retValue"); err == nil {
		returnValue, _ = value.ToString()
	}
	return output, returnValue
}

func (runner Runner) Call(from string, address string, callString string, olt int) (transaction string, returnValue string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Runtime Error: %v", r))
		}
	}()
	runner.initialContext(from, olt)
	runner.getContract(address)
	transaction, returnValue = runner.exec(callString)
	return
}

func CreateRunner() Runner {
	vm := otto.New()
	vm.Set("version", "OVM 0.1 TEST")
	return Runner{vm}
}
