/*
	Copyright 2017-2018 OneLedger
*/
package runner

import (
	"errors"
	"fmt"

	"github.com/Oneledger/protocol/node/log"
	"github.com/robertkrimen/otto"
)

func (runner Runner) exec(callString string) (string, string) {
	if callString == "" {
		callString = "default__()"
		log.Info("callString is empty, use default", "callString", callString)
	}
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
    out = JSON.stringify(transaction);
    ret = JSON.stringify(retValue);
    `)

	output := ""
	returnValue := ""

	if value, err := runner.vm.Get("out"); err == nil {
		output, _ = value.ToString()
	}

	if value, err := runner.vm.Get("ret"); err == nil {
		returnValue, _ = value.ToString()
	}
	return output, returnValue
}

//func (runner Runner) Call(from string, address string, callString string, olt int) (transaction string, returnValue string, err error) {
func (runner Runner) Call(request *OLVMRequest, result *OLVMResult) (err error) {
	log.Debug("Calling the Script")

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Runtime Error: %v", r))
		}
	}()

	log.Debug("Setup the Context")
	runner.initialContext(request.From, request.Value)

	log.Debug("Setup the SourceCode")
	runner.setupContract(request)

	log.Debug("Exec the Smart Contract")
	out, ret := runner.exec(request.CallString)
	result.Out = out
	result.Ret = ret
	return
}

func CreateRunner() Runner {
	vm := otto.New()
	vm.Set("version", "OVM 0.1 TEST")
	return Runner{vm}
}
