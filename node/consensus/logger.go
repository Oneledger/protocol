package consensus

import (
	"io"
	"os"

	"github.com/Oneledger/protocol/node/config"
	"github.com/Oneledger/protocol/node/log"
	tmconfig "github.com/tendermint/tendermint/config"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
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
	file, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, config.FilePerms)
	if err != nil {
		log.Info("Failed to open new file", "err", err)
		return nil, err
	}
	return newLogger(file, cfg)
}
