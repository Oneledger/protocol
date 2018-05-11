package core

import (
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestIsCointbase(t *testing.T) {
  tx := Tx{nil, []TxInput{TxInput{[]byte{}, -1, nil, nil}},[]TxOutput{}}
  assert.Equal(t, tx.IsCoinbase(), true)
}
