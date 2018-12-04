package shared

import (
	"github.com/Oneledger/protocol/node/log"
	"io/ioutil"
	"os"
)

func ReadFile(filePath string) []byte {
	textFile, err := os.Open(filePath)
	if err != nil {
		log.Debug("ReadFile", "err", err)
		defer textFile.Close()
		return nil
	}
	defer textFile.Close()

	byteValue, _ := ioutil.ReadAll(textFile)

	return byteValue
}
