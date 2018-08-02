/*
	Copyright 2017 - 2018 OneLedger

	Interface to specific chain functions
*/

package action

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
)

func Noop(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing Noop Command", "chain", chain, "context", context)
	return true, nil
}

func PrepareTransaction(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing PrepareTransaction Command", "chain", chain, "context", context)
	return true, nil
}

func SubmitTransaction(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing SubmitTransaction Command", "chain", chain, "context", context, "sequence", context[SEQUENCE])
	return SubmitTransactionOLT(context, chain)
}

func Initiate(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing Initiate Command", "chain", chain, "context", context)
	switch chain {

	case data.BITCOIN:
		return CreateContractBTC(context)
	case data.ETHEREUM:
		return CreateContractETH(context)
	case data.ONELEDGER:
		return CreateContractOLT(context)
	default:
		return false, nil
	}
}

func Participate(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing Participate Command", "chain", chain, "context", context)
	switch chain {

	case data.BITCOIN:
		return ParticipateBTC(context)
	case data.ETHEREUM:
		return ParticipateETH(context)
	case data.ONELEDGER:
		return ParticipateOLT(context)
	default:
		return false, nil
	}
}

func Redeem(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing Redeem Command", "chain", chain, "context", context)



	switch chain {

	case data.BITCOIN:
		return RedeemBTC(context)
	case data.ETHEREUM:
		return RedeemETH(context)
	case data.ONELEDGER:
		return RedeemOLT(context)
	default:
		return false, nil
	}
}

func Refund(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing Refund Command", "chain", chain, "context", context)
	switch chain {

	case data.BITCOIN:
		return RefundBTC(context)
	case data.ETHEREUM:
		return RefundETH(context)
	case data.ONELEDGER:
		return RefundOLT(context)
	default:
		return false, nil
	}
}

func ExtractSecret(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing ExtractSecret Command", "chain", chain, "context", context)
	switch chain {

	case data.BITCOIN:
		return ExtractSecretBTC(context)
	case data.ETHEREUM:
		return ExtractSecretETH(context)
	case data.ONELEDGER:
		return ExtractSecretOLT(context)
	default:
		return false, nil
	}
}

func AuditContract(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing AuditContract Command", "chain", chain, "context", context)
	switch chain {

	case data.BITCOIN:
		return AuditContractBTC(context)
	case data.ETHEREUM:
		return AuditContractETH(context)
	case data.ONELEDGER:
		return AuditContractOLT(context)
	default:
		return false, nil
	}
}

func WaitForChain(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing WaitForChain Command", "chain", chain, "context", context)
	//todo : make this to check finish status, and then rollback if necessary
	// Make sure it is pushed forward first...
	global.Current.Sequence += 32

	signers := []id.PublicKey(nil)
    owner := GetParty(context[MY_ACCOUNT])
    target := GetParty(context[THEM_ACCOUNT])
    eventType := GetType(context[EVENTTYPE])
	verify := Verify{
		Base: Base{
			Type:     VERIFY,
			ChainId:  "OneLedger-Root",
			Owner:    owner.Key,
			Signers:  signers,
			Sequence: global.Current.Sequence,
		},
		Target: owner.Key,
		Event:  Event{
		    Type:   eventType,
		    Key:    target.Key,
		    Nonce:  global.Current.Sequence,
        },
	}
	DelayedTransaction(VERIFY, Transaction(verify), 3*lockPeriod)

	return true, nil
}
