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

func SubmitTransaction(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing SubmitTransaction Command", "chain", chain, "context", context)
	return true, nil
}

func Initiate(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing Initiate Command", "chain", chain, "context", context)
	switch chain {

	case data.BITCOIN:
		return true, nil
	case data.ETHEREUM:
		return true, nil
	case data.ONELEDGER:
		return true, nil
	default:
		return false, nil
	}
}

func Participate(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing Participate Command", "chain", chain, "context", context)
	return true, nil
}

func Redeem(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing Redeem Command", "chain", chain, "context", context)
	return true, nil
}

func Refund(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing Refund Command", "chain", chain, "context", context)
	return true, nil
}

func ExtractSecret(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing ExtractSecret Command", "chain", chain, "context", context)
	return true, nil
}

func AuditContract(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing AuditContract Command", "chain", chain, "context", context)
	return true, nil
}

func WaitForChain(chain data.ChainType, context map[Parameter]FunctionValue) (bool, map[Parameter]FunctionValue) {
	log.Info("Executing WaitForChain Command", "chain", chain, "context", context)

	// Make sure it is pushed forward first...
	global.Current.Sequence += 32

	signers := []id.PublicKey(nil)

	verify := Verify{
		Base: Base{
			Type:     VERIFY,
			ChainId:  "OneLedger-Root",
			Signers:  signers,
			Sequence: global.Current.Sequence,
		},
	}
	BroadcastTransaction(VERIFY, Transaction(verify))

	return true, nil
}
