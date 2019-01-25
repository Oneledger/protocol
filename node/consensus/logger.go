package consensus

import (
	"os"

	"github.com/Oneledger/protocol/node/log"
	tmconfig "github.com/tendermint/tendermint/config"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/common"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

func NewLogger(logPath string, cfg tmconfig.Config) tmlog.Logger {
	var file *os.File
	var err error
	if !common.FileExists(logPath) {
		log.Info("Creating consensus log file at", "path", logPath)
		file, err = os.Create(logPath)
		if err != nil {
			log.Fatal("Failed to create logging file", "location", logPath, "err", err)
		}
	} else {
		file, err = os.OpenFile(logPath, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Failed to open logging file", "location", logPath, "err", err)
		}
	}
	tmLogger := tmlog.NewTMLogger(tmlog.NewSyncWriter(file))
	logger, err := tmflags.ParseLogLevel(cfg.LogLevel, tmLogger, tmconfig.DefaultLogLevel())
	if err != nil {
		log.Fatal("Failed to configure loglevel for logger", "loglevel", cfg.LogLevel)
	}
	return logger
}
