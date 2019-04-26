package runner

import (
	"bytes"
	"github.com/Oneledger/protocol/node/log"
	"os"
	"path/filepath"
)

func getCodeFromJsLibs(address string) string {

	jsFilePath := filepath.Join(os.Getenv("OLROOT"), "/protocol/node/olvm/interpreter/runner/js", address+".js")
	log.Debug("get source code from local file system", "path", jsFilePath)

	file, err := os.Open(jsFilePath)
	if err != nil {
		log.Fatal("cannot get source code", "err", err)
		return ""
	}

	defer file.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	contents := buf.String()

	return contents
}
