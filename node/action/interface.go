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

func Noop(app interface{}, chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing Noop Command", "chain", chain, "context", context)
	return true, nil
}

func PrepareTransaction(app interface{}, chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing PrepareTransaction Command", "chain", chain, "context", context)
	return true, nil
}

func SubmitTransaction(app interface{}, chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing SubmitTransaction Command", "chain", chain, "context", context, "sequence", context[COUNT])
	return SubmitTransactionOLT(context, chain)
}

func Initiate(app interface{}, chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing Initiate Command", "chain", chain, "context", context)
	switch chain {

	case data.BITCOIN:
		return CreateContractBTC(app, context)
	case data.ETHEREUM:
		return CreateContractETH(app, context)
	case data.ONELEDGER:
		return CreateContractOLT(app, context)
	default:
		return false, nil
	}
}

func Participate(app interface{}, chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing Participate Command", "chain", chain, "context", context)
	switch chain {

	case data.BITCOIN:
		return ParticipateBTC(app, context)
	case data.ETHEREUM:
		return ParticipateETH(app, context)
	case data.ONELEDGER:
		return ParticipateOLT(app, context)
	default:
		return false, nil
	}
}

func Redeem(app interface{}, chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing Redeem Command", "chain", chain, "context", context)

	switch chain {

	case data.BITCOIN:
		return RedeemBTC(app, context)
	case data.ETHEREUM:
		return RedeemETH(app, context)
	case data.ONELEDGER:
		return RedeemOLT(app, context)
	default:
		return false, nil
	}
}

func Refund(app interface{}, chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing Refund Command", "chain", chain, "context", context)
	switch chain {

	case data.BITCOIN:
		return RefundBTC(app, context)
	case data.ETHEREUM:
		return RefundETH(app, context)
	case data.ONELEDGER:
		return RefundOLT(app, context)
	default:
		return false, nil
	}
}

func ExtractSecret(app interface{}, chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing ExtractSecret Command", "chain", chain, "context", context)
	switch chain {

	case data.BITCOIN:
		return ExtractSecretBTC(app, context)
	case data.ETHEREUM:
		return ExtractSecretETH(app, context)
	case data.ONELEDGER:
		return ExtractSecretOLT(app, context)
	default:
		return false, nil
	}
}

func AuditContract(app interface{}, chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing AuditContract Command", "chain", chain, "context", context)
	switch chain {

	case data.BITCOIN:
		return AuditContractBTC(app, context)
	case data.ETHEREUM:
		return AuditContractETH(app, context)
	case data.ONELEDGER:
		return AuditContractOLT(app, context)
	default:
		return false, nil
	}
}

func WaitForChain(app interface{}, chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing WaitForChain Command", "chain", chain, "context", context)
	//todo : make this to check finish status, and then rollback if necessary
	// Make sure it is pushed forward first...
	global.Current.Sequence += 32

	signers := []id.PublicKey(nil)
    owner := GetParty(context[MY_ACCOUNT])
    target := GetParty(context[THEM_ACCOUNT])
    eventType := GetType(context[EVENTTYPE])
    nonce := GetInt64(context[NONCE])
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
		    Nonce:  nonce,
        },
	}
	DelayedTransaction(VERIFY, Transaction(verify), 3 * lockPeriod)

	return true, nil
}
