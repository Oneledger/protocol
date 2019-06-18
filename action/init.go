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

	//ons related transaction
	DOMAIN_CREATE
	DOMAIN_UPDATE
	DOMAIN_SELL
	DOMAIN_BUY
	DOMAIN_SEND
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
	case DOMAIN_CREATE:
		return "DOMAIN_CREATE"
	case DOMAIN_UPDATE:
		return "DOMAIN_UPDATE"
	case DOMAIN_SELL:
		return "DOMAIN_SELL"
	case DOMAIN_BUY:
		return "DOMAIN_BUY"
	case DOMAIN_SEND:
		return "DOMAIN_SEND"
	default:
		return "UNKNOWN"
	}
}
