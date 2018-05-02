package core

import (
  "testing"
  "../utils"

  "github.com/stretchr/testify/assert"
)

func TestAssignPubKeyHash(t *testing.T) {
  tx := TxOutput{100, nil}
  address := []byte("this is a test")
  tx.AssignPubKeyHash(address)
  assert.Equal(t, tx.PubKeyHash, utils.HashPubKey(address))
}

func TestVerifyTxOutput(t *testing.T) {
  tx := TxOutput{100, nil}
  address := []byte("this is a test")
  tx.AssignPubKeyHash(address)
  pubKeyHash := utils.HashPubKey(address)
  assert.Equal(t,tx.VerifyTxOutput(pubKeyHash), true)
}
