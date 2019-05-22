package runner

import (
	"bytes"
	"os"
	"path/filepath"
)

func getCodeFromJsLibs(address string) string {

	jsFilePath := filepath.Join(os.Getenv("OLROOT"), "/protocol/olvm/interpreter/runner/js", address+".js")
	log.Debug("get source code from local file system", "path", jsFilePath)

	file, err := os.Open(jsFilePath)
	if err != nil {
		log.Fatal("cannot get source code", "err", err)
		return ""
	}

	defer file.Close()

	buf := new(bytes.Buffer)

	// TODO : maybe implement a limit on file size?
	_, err = buf.ReadFrom(file)
	if err != nil {
		log.Error("error in reading js file", "err", err)
	}

	contents := buf.String()

	return contents
}
