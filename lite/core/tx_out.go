package core

import (
	"../utils"
	"bytes"
)

type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

// TXOutputs collects TXOutput
type TxOutputs struct {
	Outputs []TxOutput
}

func (out *TxOutput) AssignPubKeyHash(address []byte) {
	out.PubKeyHash = utils.HashPubKey(address)
}

func (out *TxOutput) VerifyTxOutput(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func NewTxOutput(value int, address string) *TxOutput {
	output := TxOutput{value, nil}
	output.AssignPubKeyHash([]byte(address))
	return &output
}
