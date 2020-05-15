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

	//Governance Action
	PROPOSAL_CREATE Type = 0x96

	//Governanace Actions
	PROPOSAL_FUND Type = 0x31
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

	case BTC_LOCK:
		return "BTC_LOCK"
	case BTC_ADD_SIGNATURE:
		return "BTC_ADD_SIGNATURE"
	case BTC_BROADCAST_SUCCESS:
		return "BTC_BROADCAST_SUCCESS"
	case BTC_REPORT_FINALITY_MINT:
		return "BTC_REPORT_FINALITY_MINT"
	case BTC_EXT_MINT:
		return "BTC_EXT_MINT"
	case BTC_REDEEM:
		return "BTC_REDEEM"
	case BTC_FAILED_BROADCAST_RESET:
		return "BTC_FAILED_BROADCAST_RESET"

	case ETH_LOCK:
		return "ETH_LOCK"
	case ETH_REPORT_FINALITY_MINT:
		return "ETH_REPORT_FINALITY_MINT"
	case ETH_REDEEM:
		return "ETH_REDEEM"
	case ERC20_LOCK:
		return "ERC20_LOCK"
	case ERC20_REDEEM:
		return "ERC20_REDEEM"

	default:
		return "UNKNOWN"
	}
}
