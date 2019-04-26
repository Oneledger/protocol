package consensus

import (
	"os"

	"github.com/Oneledger/protocol/log"
)

var logger *log.Logger

func init() {
	logger = log.NewLoggerWithPrefix(os.Stdout, "consensus:")
}
