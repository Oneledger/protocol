package action

import (
	"os"

	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
)

type Type int

const (
	SEND Type = iota
	APPLYVALIDATOR
)

var logger *log.Logger

func init() {

	serialize.RegisterInterface(new(Msg))
	serialize.RegisterConcrete(new(Send), "action_send")
	serialize.RegisterConcrete(new(ApplyValidator), "action_av")
	logger = log.NewLoggerWithPrefix(os.Stdout, "action")
}

func (t Type) String() string {
	switch t {
	case SEND:
		return "SEND"
	case APPLYVALIDATOR:
		return "APPLY_VALIDATOR"
	default:
		return "UNKNOWN"
	}
}
