package netwk_delegation

import (
	"os"

	"github.com/Oneledger/protocol/log"
)

var logger *log.Logger

func init() {
	logger = log.NewDefaultLogger(os.Stdout).WithPrefix("netwk_delegation")
}
