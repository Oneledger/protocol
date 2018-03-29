package core

import (
  "../utils"
  "bytes"
)

type TxOutput struct {
  Value int
  PubKeyHash []byte
}

func (out *TxOutput) AssignPubKeyHash(address []byte) {
  pubKeyHash := utils.Base58Decode(address)
  pubKeyHash = pubKeyHash[1:len(pubKeyHash) - 4]
  out.PubKeyHash = pubKeyHash
}

func (out *TxOutput) VerifyTxOutput(pubKeyHash []byte) bool {
  return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func NewTxOutput(value int, address string) *TxOutput {
  output := TxOutput{value, nil}
  output.AssignPubKeyHash([]byte(address))
  return &output
}
