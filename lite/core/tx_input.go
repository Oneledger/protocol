package core

import (
	"../utils"
	"bytes"
)

type TxInput struct {
	Id          []byte
	OutputIndex int
	Signature   []byte
	PubKey      []byte
}

func (in *TxInput) isOwnedByPubKeyHash(pubKeyHash []byte) bool {
	return bytes.Compare(utils.HashPubKey(in.PubKey), pubKeyHash) == 0
}
