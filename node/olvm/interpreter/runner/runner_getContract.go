/*
	Copyright 2017-2018 OneLedger
*/
package runner

import (
	"bytes"
	"os"
	"strings"

	"github.com/Oneledger/protocol/node/log"
)

func (runner Runner) setupContract(request *OLVMRequest) bool {
	address := request.Address
	sourceCode := ""

	switch {
	case strings.HasPrefix(address, "samples://"):
		sourceCode = getSourceCodeFromSamples(address)
	case address == "embed://":
		sourceCode = request.SourceCode
	default:
		sourceCode = getSourceCodeFromBlockChain(address)
	}

	// TODO: Needs better error handling
	if sourceCode == "" {
		return false
	}
	log.Debug("get source code", "sourceCode", sourceCode)
	_, error := runner.vm.Run(`var module = {};(function(module){` + sourceCode + `})(module)`)
	if error == nil {
		return true
	} else {
		return false
	}
}

func getSourceCodeFromSamples(address string) string {

	prefix := "samples://"
	sampleCodeName := address[len(prefix):]

	file, err := os.Open("./samples/" + sampleCodeName + ".js")
	if err != nil {

		// TODO: Needs better error handling
		return ""
		//log.Fatal(err)
	}

	defer file.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	contents := buf.String()

	return contents
}

func getSourceCodeFromBlockChain(address string) string {
	log.Fatal("SourceCodeFrom BlockChain is Unimplemented")
	return ""
}
