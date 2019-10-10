/*

 */

package bitcoin

import (
	"bytes"
	"crypto/rand"
	"io"

	"github.com/Oneledger/protocol/data/bitcoin"

	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type TxHash [chainhash.HashSize]byte

type ChainDriver interface {
	PrepareLock(*UTXO, *UTXO, []byte) []byte
	AddUserLockSignature([]byte, []byte) *wire.MsgTx
	AddLockSignature([]byte, []byte) *wire.MsgTx
	BroadcastTx(*wire.MsgTx, *rpcclient.Client) (*chainhash.Hash, error)

	CheckFinality(chainhash.Hash) (bool, error)

	PrepareRedeem(UTXO, []byte, int64, []byte) []byte
}

type chainDriver struct {
	blockCypherToken string
}

var _ ChainDriver = &chainDriver{}

func NewChainDriver(token string) ChainDriver {

	return &chainDriver{token}
}

func (c *chainDriver) PrepareLock(prevLock, input *UTXO, lockScriptAddress []byte) (txBytes []byte) {
	tx := wire.NewMsgTx(wire.TxVersion)

	if prevLock.TxID != bitcoin.NilTxHash {
		prevLockOP := wire.NewOutPoint(prevLock.TxID, prevLock.Index)
		vin1 := wire.NewTxIn(prevLockOP, nil, nil)
		tx.AddTxIn(vin1)
	}

	userOP := wire.NewOutPoint(input.TxID, input.Index)
	vin2 := wire.NewTxIn(userOP, nil, nil)
	tx.AddTxIn(vin2)

	// will adjust fees later
	balance := prevLock.Balance + input.Balance

	out := wire.NewTxOut(balance, lockScriptAddress)
	tx.AddTxOut(out)

	tempBuf := bytes.NewBuffer([]byte{})
	tx.Serialize(tempBuf)
	size := len(tempBuf.Bytes()) * 2
	fees := int64(70 * size)

	tx.TxOut[0].Value = tx.TxOut[0].Value - fees

	buf := bytes.NewBuffer(txBytes)
	tx.Serialize(buf)
	txBytes = buf.Bytes()

	return
}

func (c *chainDriver) AddUserLockSignature(txBytes []byte, sigScript []byte) *wire.MsgTx {
	tx := wire.NewMsgTx(wire.TxVersion)

	buf := bytes.NewBuffer(txBytes)
	tx.Deserialize(buf)

	tx.TxIn[0].SignatureScript = sigScript

	return tx
}

func (c *chainDriver) AddLockSignature(txBytes []byte, sigScript []byte) *wire.MsgTx {
	tx := wire.NewMsgTx(wire.TxVersion)

	buf := bytes.NewBuffer(txBytes)
	tx.Deserialize(buf)

	tx.TxIn[0].SignatureScript = sigScript

	return tx
}

func (c *chainDriver) BroadcastTx(tx *wire.MsgTx, clt *rpcclient.Client) (*chainhash.Hash, error) {

	hash, err := clt.SendRawTransaction(tx, false)
	if err != nil {
		return &chainhash.Hash{}, err
	}

	return hash, nil
}

func (c *chainDriver) CheckFinality(hash chainhash.Hash) (bool, error) {

	btc := gobcy.API{"dd53aae66b83431ca57a1f656af8ed69", "btc", "main"}
	tx, err := btc.GetTX(hash.String(), nil)
	if err != nil {
		return false, err
	}

	if tx.Confirmations > 10 {
		return true, nil
	}

	return false, nil
}

func (c *chainDriver) PrepareRedeem(prevLock UTXO,
	userAddress []byte, amount int64,
	lockAddress []byte) (txBytes []byte) {

	tx := wire.NewMsgTx(wire.TxVersion)

	prevLockOP := wire.NewOutPoint(prevLock.TxID, prevLock.Index)
	vin1 := wire.NewTxIn(prevLockOP, nil, nil)
	tx.AddTxIn(vin1)

	balance := prevLock.Balance - amount

	out := wire.NewTxOut(balance, lockAddress)
	tx.AddTxOut(out)

	userOP := wire.NewTxOut(amount, userAddress)
	tx.AddTxOut(userOP)

	tempBuf := bytes.NewBuffer([]byte{})
	tx.Serialize(tempBuf)
	size := len(tempBuf.Bytes()) * 2
	fees := int64(70 * size)

	tx.TxOut[0].Value = tx.TxOut[0].Value - fees

	buf := bytes.NewBuffer(txBytes)
	tx.Serialize(buf)
	txBytes = buf.Bytes()

	return
}

func CreateMultiSigAddress(m int, publicKeys []*btcutil.AddressPubKey) (script, address []byte, err error) {
	// ideally m should be
	//	m = len(publicKeys) * 2 /3 ) + 1

	builder := txscript.NewScriptBuilder().AddInt64(int64(m))
	for _, key := range publicKeys {
		builder.AddData(key.ScriptAddress())
	}
	builder.AddInt64(int64(len(publicKeys)))
	builder.AddOp(txscript.OP_CHECKMULTISIG)

	// add randomness
	builder.AddOp(txscript.OP_RETURN)

	data := [4]byte{}
	_, err = io.ReadFull(rand.Reader, data[:])
	if err != nil {
		return
	}

	builder.AddData(data[:])

	script, err = builder.Script()
	if err != nil {
		return
	}

	addressObj, err := btcutil.NewAddressScriptHash(script, &chaincfg.RegressionNetParams)
	if err != nil {
		return
	}

	address = addressObj.ScriptAddress()

	return
}
