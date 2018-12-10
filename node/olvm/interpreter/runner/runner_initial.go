package runner

func (runner Runner) initialContext(from string, olt int) {
  sourceCode := getCodeFromJsLibs("onEnter")
  runner.vm.Run(sourceCode)
	runner.vm.Set("__from__", from)
	runner.vm.Set("__olt__", olt)
}
