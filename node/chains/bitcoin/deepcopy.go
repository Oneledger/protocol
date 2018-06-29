package bitcoin

import "github.com/btcsuite/btcd/wire"

func copyArray(from []byte) []byte {
	to := make([]byte, len(from))
	for i := 0; i < len(from); i++ {
		to[i] = from[i]
	}
	return to
}

func copyMsgTx(from *wire.MsgTx) *wire.MsgTx {
	var to wire.MsgTx

	to.Version = from.Version
	to.TxIn = copyTxIn(from.TxIn)
	to.TxOut = copyTxOut(from.TxOut)
	to.LockTime = from.LockTime

	return &to
}

func copyTxIn(from []*wire.TxIn) []*wire.TxIn {
	var to []*wire.TxIn
	to = make([]*wire.TxIn, len(from))
	for i := 0; i < len(from); i++ {
		to[i] = &wire.TxIn{}
		to[i].PreviousOutPoint = from[i].PreviousOutPoint
		to[i].SignatureScript = copyArray(from[i].SignatureScript)
		to[i].Witness = from[i].Witness
		to[i].Sequence = from[i].Sequence
	}
	return to
}

func copyTxOut(from []*wire.TxOut) []*wire.TxOut {
	var to []*wire.TxOut
	to = make([]*wire.TxOut, len(from))
	for i := 0; i < len(from); i++ {
		to[i] = &wire.TxOut{}
		to[i].Value = from[i].Value
		to[i].PkScript = copyArray(from[i].PkScript)
	}
	return to
}
