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

func Noop(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {
	log.Info("Executing Noop Command", "chain", chain, "context", context)
	return true, nil
}

func PrepareTransaction(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {
	log.Info("Executing PrepareTransaction Command", "chain", chain, "context", context)
	return true, nil
}

func SubmitTransaction(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {
	log.Info("Executing SubmitTransaction Command", "chain", chain, "context", context, "sequence", context[COUNT])
	return SubmitTransactionOLT(context, chain)
}

func WaitForChain(app interface{}, chain data.ChainType, context FunctionValues) (bool, FunctionValues) {
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
