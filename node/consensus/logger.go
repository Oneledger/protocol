package consensus

import (
	"io"
	"os"

	"github.com/Oneledger/protocol/node/log"
	tmconfig "github.com/tendermint/tendermint/config"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/common"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

func newLogger(w io.Writer, cfg tmconfig.Config) (tmlog.Logger, error) {
	tmLogger := tmlog.NewTMLogger(w)
	return tmflags.ParseLogLevel(cfg.LogLevel, tmLogger, tmconfig.DefaultLogLevel())
}

func newStdOutLogger(cfg tmconfig.Config) (tmlog.Logger, error) {
	return newLogger(os.Stdout, cfg)
}

func newFileLogger(logPath string, cfg tmconfig.Config) (tmlog.Logger, error) {
	var file *os.File
	var err error
	if !common.FileExists(logPath) {
		log.Info("Creating consensus log file at", "path", logPath)
		file, err = os.Create(logPath)
		if err != nil {
			return nil, err
		}
	} else {
		file, err = os.OpenFile(logPath, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
	}
	return newLogger(file, cfg)
}
