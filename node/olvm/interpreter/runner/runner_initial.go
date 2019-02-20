package runner

import (
	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/log"
	"github.com/robertkrimen/otto"
)

func (runner Runner) initialContext(from string, olt int, callString string, context action.OLVMContext) {

	runner.vm.Set("__GetContextValue__", func(call otto.FunctionCall) otto.Value {
		key := call.Argument(0).String()
		log.Debug("OLVM get value from context", "key", key)
		value, _ := runner.vm.ToValue(context.GetValue(key))

		return value
	})
	sourceCode := getCodeFromJsLibs("onEnter")
	runner.vm.Set("__callString__", callString)
	runner.vm.Run(sourceCode)
	runner.vm.Set("__from__", from)
	runner.vm.Set("__olt__", olt)
}

func (runner Runner) initialAnalyzeContext(from string, olt int, callString string, context action.OLVMContext) {

	sourceCode := getCodeFromJsLibs("onAnalyzeEnter")
	runner.vm.Set("__callString__", callString)
	runner.vm.Run(sourceCode)
	runner.vm.Set("__from__", from)
	runner.vm.Set("__olt__", olt)
}
