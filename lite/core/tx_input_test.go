package core

import (
	"../utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsOwnedByAddress(t *testing.T) {
	txInput := &TxInput{[]byte{}, 0, []byte{}, []byte("thisismyaddress")}
	assert.Equal(t, txInput.isOwnedByPubKeyHash(utils.HashPubKey([]byte("thisismyaddress"))), true)
}
