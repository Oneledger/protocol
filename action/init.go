package action

import (
	"os"

	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
)

type Type int

const (
	SEND Type = 0x01

	//staking related transaction
	APPLYVALIDATOR Type = 0x11
	WITHDRAW       Type = 0x12

	//ons related transaction
	DOMAIN_CREATE   Type = 0x21
	DOMAIN_UPDATE   Type = 0x22
	DOMAIN_SELL     Type = 0x23
	DOMAIN_PURCHASE Type = 0x24
	DOMAIN_SEND     Type = 0x25
)

var logger *log.Logger

func init() {

	serialize.RegisterInterface(new(Msg))
	logger = log.NewLoggerWithPrefix(os.Stdout, "action")
}

func (t Type) String() string {
	switch t {
	case SEND:
		return "SEND"
	case APPLYVALIDATOR:
		return "APPLY_VALIDATOR"
	case WITHDRAW:
		return "WITHDRAW"
	case DOMAIN_CREATE:
		return "DOMAIN_CREATE"
	case DOMAIN_UPDATE:
		return "DOMAIN_UPDATE"
	case DOMAIN_SELL:
		return "DOMAIN_SELL"
	case DOMAIN_PURCHASE:
		return "DOMAIN_PURCHASE"
	case DOMAIN_SEND:
		return "DOMAIN_SEND"
	default:
		return "UNKNOWN"
	}
}
