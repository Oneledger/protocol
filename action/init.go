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
	PURGE          Type = 0x13

	//ons related transaction
	DOMAIN_CREATE     Type = 0x21
	DOMAIN_UPDATE     Type = 0x22
	DOMAIN_SELL       Type = 0x23
	DOMAIN_PURCHASE   Type = 0x24
	DOMAIN_SEND       Type = 0x25
	DOMAIN_DELETE_SUB Type = 0x26
	DOMAIN_RENEW      Type = 0x27

	BTC_LOCK                   Type = 0x81
	BTC_ADD_SIGNATURE          Type = 0x82
	BTC_BROADCAST_SUCCESS      Type = 0x83
	BTC_REPORT_FINALITY_MINT   Type = 0x84
	BTC_EXT_MINT               Type = 0x85
	BTC_REDEEM                 Type = 0x86
	BTC_FAILED_BROADCAST_RESET Type = 0x87

	//Ethereum Actions
	ETH_LOCK                 Type = 0x91
	ETH_REPORT_FINALITY_MINT Type = 0x92
	ETH_REDEEM               Type = 0x93
	ERC20_LOCK               Type = 0x94
	ERC20_REDEEM             Type = 0x95

	//Passport Actions
	PASSPORT_HOSP_ADMIN  Type = 0x61
	PASSPORT_SCR_ADMIN   Type = 0x62
	PASSPORT_UPLOAD_TEST Type = 0x63
	PASSPORT_READ_TEST   Type = 0x64
	PASSPORT_UPDATE_TEST Type = 0x65

	//EOF here Only used as a marker to mark the end of Type list
	//So that the query for Types can return all Types dynamically
	//, when there is a change made in Type list
	//This value should be manually set as the largest among the list
	EOF Type = 0xFF
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
	case DOMAIN_DELETE_SUB:
		return "DOMAIN_DELETE_SUB"
	case DOMAIN_RENEW:
		return "DOMAIN_RENEW"

	case PASSPORT_HOSP_ADMIN:
		return "PASSPORT_HOSP_ADMIN"
	case PASSPORT_UPLOAD_TEST:
		return "PASSPORT_UPLOAD_TEST"
	case PASSPORT_UPDATE_TEST:
		return "PASSPORT_UPDATE_TEST"

	default:
		return "UNKNOWN"
	}
}
