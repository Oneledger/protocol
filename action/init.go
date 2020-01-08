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
	DOMAIN_CREATE           Type = 0x21
	DOMAIN_UPDATE           Type = 0x22
	DOMAIN_SELL             Type = 0x23
	DOMAIN_PURCHASE         Type = 0x24
	DOMAIN_SEND             Type = 0x25
	DOMAIN_EXPIRED_PURCHASE Type = 0x26
	DOMAIN_CREATE_SUB       Type = 0x27
	DOMAIN_RENEW            Type = 0x28

	BTC_LOCK                   Type = 0x81
	BTC_ADD_SIGNATURE          Type = 0x82
	BTC_BROADCAST_SUCCESS      Type = 0x83
	BTC_REPORT_FINALITY_MINT   Type = 0x84
	BTC_EXT_MINT               Type = 0x85
	BTC_REDEEM                 Type = 0x86
	BTC_FAILED_BROADCAST_RESET Type = 0x87

	//Ethereum Actions
	ETH_LOCK                 Type = 0x91
	ETH_SIGN                 Type = 0x92
	ETH_FINALITY             Type = 0x93
	ETH_MINT                 Type = 0x94
	ETH_REPORT_FINALITY_MINT Type = 0x95
	ETH_REDEEM               Type = 0x96
	ERC20_LOCK				 Type = 0x97
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

	case BTC_LOCK:
		return "BTC_LOCK"
	case BTC_ADD_SIGNATURE:
		return "BTC_ADD_SIGNATURE"
	case BTC_REPORT_FINALITY_MINT:
		return "BTC_REPORT_FINALITY_MINT"
	case BTC_EXT_MINT:
		return "BTC_EXT_MINT"
	case ETH_LOCK:
		return "ETH_LOCK"
	case ETH_REPORT_FINALITY_MINT:
		return "ETH_REPORT_FINALITY_MINT"
	case ERC20_LOCK:
		return "ERC20_LOCK"
	default:
		return "UNKNOWN"
	}
}
