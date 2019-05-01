package identity

import (
	"github.com/Oneledger/protocol/log"
	"os"
)

var logger *log.Logger

func init() {

	logger = log.NewDefaultLogger(os.Stdout).WithPrefix("identity")

}
