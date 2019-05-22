/*
	Copyright 2017-2018 OneLedger
*/
package runner

import (
	"errors"
	"fmt"
	"time"

	"github.com/Oneledger/protocol/data"
	"github.com/robertkrimen/otto"
)

func (runner Runner) analyze(callString string) (string, string) {

	if callString == "" {
		callString = "default__()"
		log.Info("callString is empty, use default", "callString", callString)
	}

	// Pretext to set up the execution
	_, err := runner.vm.Run(`
    var contract = new module.Contract(context);
    `)
	if err != nil {
		panic(err)
	}

	sourceCode := getCodeFromJsLibs("onAnalyzeExit")
	// Set the transaction parameters
	runner.vm.Run(sourceCode)
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

func (runner Runner) exec(callString string) (string, string) {

	if callString == "" {
		callString = "default__()"
		log.Info("callString is empty, use default", "callString", callString)
	}

	// Pretext to set up the execution
	_, err := runner.vm.Run(`
    var contract = new module.Contract(context);
    var retValue = contract.` + callString)
	if err != nil {
		panic(err)
	}

	sourceCode := getCodeFromJsLibs("onExit")
	// Set the transaction parameters
	// TODO: use value returned below
	_, err = runner.vm.Run(sourceCode)
	if err != nil {
		log.Error("error in running vm", "err", err)
	}

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

func (runner Runner) Analyze(request *data.OLVMRequest, result *data.OLVMResult) (err error) {
	log.Debug("Analyzing the Script")

	defer func() {
		if r := recover(); r != nil {
			log.Debug("HALTING")
			err = errors.New(fmt.Sprintf("Runtime Error: %v", r))
			result.Out = r.(error).Error()
			result.Ret = "HALT"
			result.Elapsed = "Timed out after 3 secs"
		}
	}()

	log.Debug("Setup the Context")
	runner.initialAnalyzeContext(request.From, request.Value, request.CallString, request.Context)

	log.Debug("Setup the SourceCode")
	runner.setupContract(request)

	done := make(chan string, 1) // The buffer prevents blocking

	// Setup a go routine to timeout and interrup processing in otto
	runner.vm.Interrupt = make(chan func(), 1) // The buffer prevents blocking
	go func() {
		for {
			select {
			case <-time.After(3 * time.Second):
				runner.vm.Interrupt <- func() {
					log.Debug("Halting due to timeout")
					panic(errors.New("Halting Execution"))
				}
			case status := <-done:
				log.Debug("Finished timer Cleanly", "status", status)
				return
			}
		}
	}()

	log.Debug("Exec the Smart Contract")

	start := time.Now()
	out, ret := runner.analyze(request.CallString)
	elapsed := time.Since(start)

	done <- "Finished"

	log.Debug("Smart Contract is finished")

	result.Out = out
	result.Ret = ret
	result.Elapsed = elapsed.String()
	return
}

//func (runner Runner) Call(from string, address string, callString string, olt int) (transaction string, returnValue string, err error) {
func (runner Runner) Call(request *data.OLVMRequest, result *data.OLVMResult) (err error) {
	log.Debug("Calling the Script")
	defer func() {
		if r := recover(); r != nil {
			log.Debug("HALTING")
			err = errors.New(fmt.Sprintf("Runtime Error: %v", r))
			result.Out = r.(error).Error()
			result.Ret = "HALT"
			result.Elapsed = "Timed out after 3 secs"
		}
	}()

	log.Debug("Setup the Context")
	runner.initialContext(request.From, request.Value, request.CallString, request.Context)

	log.Debug("Setup the SourceCode")
	runner.setupContract(request)

	done := make(chan string, 1) // The buffer prevents blocking

	// Setup a go routine to timeout and interrup processing in otto
	runner.vm.Interrupt = make(chan func(), 1) // The buffer prevents blocking
	go func() {
		for {
			select {
			case <-time.After(3 * time.Second):
				runner.vm.Interrupt <- func() {
					log.Debug("Halting due to timeout")
					panic(errors.New("Halting Execution"))
				}
			case status := <-done:
				log.Debug("Finished timer Cleanly", "status", status)
				return
			}
		}
	}()

	log.Debug("Exec the Smart Contract")

	start := time.Now()
	out, ret := runner.exec(request.CallString)
	elapsed := time.Since(start)

	done <- "Finished"

	log.Debug("Smart Contract is finished")

	result.Out = out
	result.Ret = ret
	result.Elapsed = elapsed.String()
	return
}

func CreateRunner() Runner {
	vm := otto.New()
	err := vm.Set("version", "OVM v0.1.0 TEST")
	if err != nil {
		log.Errorf("error in setting version", "err", err)
	}
	return Runner{vm}
}
