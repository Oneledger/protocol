package runner

import (
	"bytes"
	"log"
	"os"
	"strings"
)

func (runner Runner) getContract(address string) bool {
	sourceCode := ""
	switch {
	case strings.HasPrefix(address, "samples://"):
		sourceCode = getSourceCodeFromSamples(address)
	default:
		sourceCode = getSourceCodeFromBlockChain(address)
	}
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
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	contents := buf.String()

	return contents

}

func getSourceCodeFromBlockChain(address string) string {
	log.Fatal("Unimplemented")
	return ""
}
