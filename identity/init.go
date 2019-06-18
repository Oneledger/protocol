package identity

import (
	"os"

	"github.com/Oneledger/protocol/log"
)

var logger *log.Logger

func init() {

	logger = log.NewDefaultLogger(os.Stdout).WithPrefix("identity")

}
