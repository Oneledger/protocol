package consensus

import (
	"os"

	"github.com/Oneledger/protocol/node/log"
	"github.com/tendermint/tendermint/libs/common"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

func NewLogger(logPath string) tmlog.Logger {
	var file *os.File
	var err error
	if !common.FileExists(logPath) {
		log.Info("Creating consensus log file at", "path", logPath)
		file, err = os.Create(logPath)
		if err != nil {
			log.Fatal("Failed to create logging file", "location", logPath, "err", err)
		}
	} else {
		file, err = os.Open(logPath)
		if err != nil {
			log.Fatal("Failed to open logging file", "location", logPath, "err", err)
		}
	}
	return tmlog.NewTMLogger(file)
}
