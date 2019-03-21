package shared

import (
	"io/ioutil"
	"os"

	"github.com/Oneledger/protocol/node/log"
)

func MustReadFile(filePath string) []byte {
	bz, err := ReadFile(filePath)
	if err != nil {
		log.Error("Failed to ReadFile", "filepath", filePath)
		return nil
	}

	return bz
}

func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}
