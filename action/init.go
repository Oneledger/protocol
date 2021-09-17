package action

import (
	"os"

	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/serialize"
)

type Type int

type TxTypeMap map[Type]string

var txTypeMap TxTypeMap

const (
	SEND     Type = 0x01
	SENDPOOL Type = 0x02

	//staking related transaction
	STAKE    Type = 0x11
	UNSTAKE  Type = 0x12
	WITHDRAW Type = 0x13

	//network network_delegation
	ADD_NETWORK_DELEGATE              Type = 0x51
	NETWORK_UNDELEGATE                Type = 0x52
	REWARDS_WITHDRAW_NETWORK_DELEGATE Type = 0x53
	REWARDS_REINVEST_NETWORK_DELEGATE Type = 0x54

	//Evidence
	ALLEGATION      Type = 0x61
	ALLEGATION_VOTE Type = 0x62
	RELEASE         Type = 0x63

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

	// OLVM transactions (new sends + evm)
	OLVM Type = 0x101

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
	txTypeMap = TxTypeMap{}
	RegisterTxType(SEND, "SEND")
	RegisterTxType(SENDPOOL, "SENDPOOL")

	RegisterTxType(STAKE, "STAKE")
	RegisterTxType(UNSTAKE, "UNSTAKE")
	RegisterTxType(WITHDRAW, "WITHDRAW")

	RegisterTxType(DOMAIN_CREATE, "DOMAIN_CREATE")
	RegisterTxType(DOMAIN_UPDATE, "DOMAIN_UPDATE")
	RegisterTxType(DOMAIN_SELL, "DOMAIN_SELL")
	RegisterTxType(DOMAIN_PURCHASE, "DOMAIN_PURCHASE")
	RegisterTxType(DOMAIN_SEND, "DOMAIN_SEND")
	RegisterTxType(DOMAIN_DELETE_SUB, "DOMAIN_DELETE_SUB")
	RegisterTxType(DOMAIN_RENEW, "DOMAIN_RENEW")

	RegisterTxType(BTC_LOCK, "BTC_LOCK")
	RegisterTxType(BTC_ADD_SIGNATURE, "BTC_ADD_SIGNATURE")
	RegisterTxType(BTC_BROADCAST_SUCCESS, "BTC_BROADCAST_SUCCESS")
	RegisterTxType(BTC_REPORT_FINALITY_MINT, "BTC_REPORT_FINALITY_MINT")
	RegisterTxType(BTC_EXT_MINT, "BTC_EXT_MINT")
	RegisterTxType(BTC_REDEEM, "BTC_REDEEM")
	RegisterTxType(BTC_FAILED_BROADCAST_RESET, "BTC_FAILED_BROADCAST_RESET")

	RegisterTxType(ETH_LOCK, "ETH_LOCK")
	RegisterTxType(ETH_REPORT_FINALITY_MINT, "ETH_REPORT_FINALITY_MINT")
	RegisterTxType(ETH_REDEEM, "ETH_REDEEM")
	RegisterTxType(ERC20_LOCK, "ERC20_LOCK")
	RegisterTxType(ERC20_REDEEM, "ERC20_REDEEM")

	RegisterTxType(PROPOSAL_CREATE, "PROPOSAL_CREATE")
	RegisterTxType(PROPOSAL_CANCEL, "PROPOSAL_CANCEL")
	RegisterTxType(PROPOSAL_FUND, "PROPOSAL_FUND")
	RegisterTxType(PROPOSAL_VOTE, "PROPOSAL_VOTE")
	RegisterTxType(PROPOSAL_FINALIZE, "PROPOSAL_FINALIZE")
	RegisterTxType(EXPIRE_VOTES, "EXPIRE_VOTES")
	RegisterTxType(PROPOSAL_WITHDRAW_FUNDS, "PROPOSAL_WITHDRAW_FUNDS")

	RegisterTxType(WITHDRAW_REWARD, "WITHDRAW_REWARD")

	RegisterTxType(ADD_NETWORK_DELEGATE, "ADD_NETWORK_DELEGATION")
	RegisterTxType(NETWORK_UNDELEGATE, "NETWORK_UNDELEGATE")
	RegisterTxType(REWARDS_WITHDRAW_NETWORK_DELEGATE, "REWARDS_WITHDRAW_NETWORK_DELEGATE")
	RegisterTxType(REWARDS_REINVEST_NETWORK_DELEGATE, "REWARDS_REINVEST_NETWORK_DELEGATE")

	RegisterTxType(ALLEGATION, "ALLEGATION")
	RegisterTxType(ALLEGATION_VOTE, "ALLEGATION_VOTE")
	RegisterTxType(RELEASE, "RELEASE")

	RegisterTxType(OLVM, "OLVM")
}

func RegisterTxType(value Type, name string) {
	if dupName, ok := txTypeMap[value]; ok {
		logger.Errorf("Trying to register tx type %s failed, type value conflicts with existing type: %s", value, dupName)
		return
	}
	txTypeMap[value] = name
}

func (t Type) String() string {
	if name, ok := txTypeMap[t]; ok {
		return name
	}
	return "UNKNOWN"
}
