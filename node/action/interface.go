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

func Noop(chain data.ChainType, context map[string]string) bool {
	log.Info("Executing Noop Command", "chain", chain, "context", context)
	return true
}

func SubmitTransaction(chain data.ChainType, context map[string]string) bool {
	log.Info("Executing SubmitTransaction Command", "chain", chain, "context", context)
	return true
}

func CreateLockbox(chain data.ChainType, context map[string]string) bool {
	log.Info("Executing CreateLockbox Command", "chain", chain, "context", context)
	return true
}

func SignLockbox(chain data.ChainType, context map[string]string) bool {
	log.Info("Executing SignLockbox Command", "chain", chain, "context", context)
	return true
}

func VerifyLockbox(chain data.ChainType, context map[string]string) bool {
	log.Info("Executing VerifyLockbox Command", "chain", chain, "context", context)
	return true
}

func SendKey(chain data.ChainType, context map[string]string) bool {
	log.Info("Executing SendKey Command", "chain", chain, "context", context)
	return true
}

func ReadChain(chain data.ChainType, context map[string]string) bool {
	log.Info("Executing ReadChain Command", "chain", chain, "context", context)
	return true
}

func OpenLockbox(chain data.ChainType, context map[string]string) bool {
	log.Info("Executing OpenLockbox Command", "chain", chain, "context", context)
	return true
}

func WaitForChain(chain data.ChainType, context map[string]string) bool {
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

	return true
}
