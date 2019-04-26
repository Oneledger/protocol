package runner

import (
	"github.com/Oneledger/protocol/data"
	"github.com/robertkrimen/otto"
)

func (runner Runner) initialContext(from string, olt int, callString string, context data.OLVMContext) {

	err := runner.vm.Set("__GetContextValue__", func(call otto.FunctionCall) otto.Value {
		key := call.Argument(0).String()
		log.Debug("OLVM get value from context", "key", key)
		value, _ := runner.vm.ToValue(context.GetValue(key))

		return value
	})
	logIfError("error in setting context", err)

	sourceCode := getCodeFromJsLibs("onEnter")
	err = runner.vm.Set("__callString__", callString)
	logIfError("error in setting callstring ", err)

	_, err = runner.vm.Run(sourceCode)
	logIfError("error in running sourceCode ", err)

	err = runner.vm.Set("__from__", from)
	logIfError("error in setting __from__ ", err)

	err = runner.vm.Set("__olt__", olt)
	logIfError("error in setting __olt__ ", err)
}

func (runner Runner) initialAnalyzeContext(from string, olt int, callString string, context data.OLVMContext) {

	sourceCode := getCodeFromJsLibs("onAnalyzeEnter")
	err := runner.vm.Set("__callString__", callString)
	logIfError("error in setting callstring ", err)

	_, err = runner.vm.Run(sourceCode)
	logIfError("error in running sourceCode ", err)

	err = runner.vm.Set("__from__", from)
	logIfError("error in setting __from__ ", err)

	err = runner.vm.Set("__olt__", olt)
	logIfError("error in setting __olt__ ", err)
}
