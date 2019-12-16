package ons

import (
	"github.com/Oneledger/protocol/action"
	"testing"
)

var (
	testCases           map[int]Case
	transactionHandlers map[string]action.Tx
	ctx                 *action.Context
)

type Case struct {
	//input values
	ctx      *action.Context
	tx       *action.RawTx
	signedTx action.SignedTx
	startGas action.Gas
	endGas   action.Gas
	txType   string

	//expected output
	validateResp bool
	checkResp    bool
	deliverResp  bool
	feeResp      bool
}

func setup() {
	testCases = make(map[int]Case)
	transactionHandlers = make(map[string]action.Tx)

	//create and initialize new context

}

func init() {

}

func TestONSTx(t *testing.T) {

}
