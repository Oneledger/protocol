/*
	Copyright 2017-2018 OneLedger
*/
package runner

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/Oneledger/protocol/data"
)

func (runner Runner) setupContract(request *data.OLVMRequest) bool {
	address := request.Address
	sourceCode := ""

	switch {
	case strings.HasPrefix(address, "samples://"):
		sourceCode = getSourceCodeFromSamples(address)
	case address == "embed://":
		// TODO: Should preserve byte array, to support UTF8?
		sourceCode = string(request.SourceCode)
	default:
		sourceCode = getSourceCodeFromBlockChain(address)
	}

	if sourceCode == "" {
		log.Error("error in getting source code of contract")
		return false
	}
	log.Debug("get source code", "sourceCode", sourceCode)

	_, err := runner.vm.Run(`var module = {};(function(module){` + sourceCode + `})(module)`)
	if err == nil {
		return true
	}

	return false
}

func getSourceCodeFromSamples(address string) string {

	prefix := "samples://"
	sampleCodeName := address[len(prefix):]

	jsFilePath := filepath.Join(os.Getenv("OLROOT"), "/protocol/olvm/interpreter/samples/", sampleCodeName+".js")
	log.Debug("get source code from local file system", "path", jsFilePath)

	file, err := os.Open(jsFilePath)
	if err != nil {
		log.Error("cannot get source code of contract", "err", err)
		return ""
		//log.Fatal(err)
	}

	defer file.Close()

	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(file)
	if err != nil {
		log.Error("error reading file in getSourceCodeFromSamples", "err", err)
	}
	contents := buf.String()

	return contents
}

func getSourceCodeFromBlockChain(address string) string {

	log.Fatal("SourceCodeFrom BlockChain is Unimplemented")
	return ""
}
