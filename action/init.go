package action

import (
	"os"

	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
)

type Type int

type TxTypeMap map[int]string

var txTypeMap TxTypeMap

const (
	SEND     Type = 0x01
	SENDPOOL Type = 0x02

	//staking related transaction
	STAKE    Type = 0x11
	UNSTAKE  Type = 0x12
	WITHDRAW Type = 0x13

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
	PROPOSAL_CREATE         Type = 0x30
	PROPOSAL_CANCEL         Type = 0x31
	PROPOSAL_FUND           Type = 0x32
	PROPOSAL_VOTE           Type = 0x33
	PROPOSAL_FINALIZE       Type = 0x34
	EXPIRE_VOTES            Type = 0x35
	PROPOSAL_WITHDRAW_FUNDS Type = 0x36

	//Rewards
	WITHDRAW_REWARD Type = 0x41

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
	txTypeMap  = TxTypeMap{}
}
//todo in all txs, register tx type use this func(proposal is done)
func RegisterTxType(value int, name string, ) {
	if dupName, ok := txTypeMap[value]; ok {
		logger.Errorf("Trying to register tx type %s failed, type value conflicts with existing type: %d: %s", value, dupName)
		return
	}
	txTypeMap[value] = name
}

func (t Type) String() string {
	if name, ok := txTypeMap[int(t)]; ok {
		return name
	}
	return "UNKNOWN"
}
