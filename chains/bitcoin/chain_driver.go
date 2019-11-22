/*

 */

package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sort"

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
	//	PrepareLock(prevLock, input *bitcoin.UTXO, lockScriptAddress []byte) (txBytes []byte)

	PrepareLockNew(*chainhash.Hash, uint32, int64, *chainhash.Hash, uint32, int64, int64, []byte) (txBytes []byte)

	AddUserLockSignature([]byte, []byte) *wire.MsgTx
	AddLockSignature([]byte, []byte) *wire.MsgTx

	BroadcastTx(*wire.MsgTx, *rpcclient.Client) (*chainhash.Hash, error)

	CheckFinality(*chainhash.Hash, string, string) (bool, error)

	PrepareRedeemNew(prevLockTxID *chainhash.Hash, prevLockIndex uint32, prevLockBalance int64,
		userAddress []byte, redeemAmount int64, feesInSatoshi int64,
		lockScriptAddress []byte) (txBytes []byte)
}

type chainDriver struct {
	blockCypherToken string
}

var _ ChainDriver = &chainDriver{}

func NewChainDriver(token string) ChainDriver {

	return &chainDriver{token}
}

func (c *chainDriver) PrepareLockNew(prevLockTxID *chainhash.Hash, prevLockIndex uint32, prevLockBalance int64,
	inputTxID *chainhash.Hash, inputIndex uint32, inputBalance int64, feesInSatoshi int64,
	lockScriptAddress []byte) (txBytes []byte) {

	// start a new empty txn
	tx := wire.NewMsgTx(wire.TxVersion)

	// create tracker txin & add to txn
	// if this is a newly initialized tracker the prev lock transaction id is nil,
	// in that case there is no balance in the tracker, so just go ahead without carrying balance
	if prevLockTxID != nil {

		prevLockOP := wire.NewOutPoint(prevLockTxID, prevLockIndex)
		vin1 := wire.NewTxIn(prevLockOP, nil, nil)
		tx.AddTxIn(vin1)
	}

	// create users' txin & add to txn
	userOP := wire.NewOutPoint(inputTxID, inputIndex)
	vin2 := wire.NewTxIn(userOP, nil, nil)
	tx.AddTxIn(vin2)

	// will adjust fees later
	balance := prevLockBalance + inputBalance - feesInSatoshi

	// add a txout to the lockscriptaddress
	out := wire.NewTxOut(balance, lockScriptAddress)
	tx.AddTxOut(out)

	// serialize again with fees accounted for
	buf := bytes.NewBuffer(txBytes)
	tx.Serialize(buf)
	txBytes = buf.Bytes()

	fmt.Println(hex.EncodeToString(txBytes))
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

	// if first lock
	if len(tx.TxIn) == 1 {
		return tx
	}

	tx.TxIn[0].SignatureScript = sigScript

	return tx
}

func (c *chainDriver) BroadcastTx(tx *wire.MsgTx, clt *rpcclient.Client) (*chainhash.Hash, error) {

	// temp

	hash, err := clt.SendRawTransaction(tx, false)
	if err != nil {
		return &chainhash.Hash{}, err
	}

	return hash, nil
}

func (c *chainDriver) CheckFinality(hash *chainhash.Hash, token, chain string) (bool, error) {

	btc := gobcy.API{token, "btc", chain}

	tx, err := btc.GetTX(hash.String(), nil)
	if err != nil {
		return false, err
	}

	if tx.Confirmations > 10 {
		return true, nil
	}

	return false, nil
}

func (c *chainDriver) PrepareRedeemNew(prevLockTxID *chainhash.Hash, prevLockIndex uint32,
	prevLockBalance int64, userAddress []byte, redeemAmount int64, feesInSatoshi int64,
	lockScriptAddress []byte) (txBytes []byte) {

	tx := wire.NewMsgTx(wire.TxVersion)

	prevLockOP := wire.NewOutPoint(prevLockTxID, prevLockIndex)
	vin1 := wire.NewTxIn(prevLockOP, nil, nil)
	tx.AddTxIn(vin1)

	balance := prevLockBalance - redeemAmount

	out := wire.NewTxOut(balance, lockScriptAddress)
	tx.AddTxOut(out)

	userOP := wire.NewTxOut(redeemAmount, userAddress)
	tx.AddTxOut(userOP)

	tempBuf := bytes.NewBuffer([]byte{})
	tx.Serialize(tempBuf)
	size := len(tempBuf.Bytes()) * 2
	fees := int64(40 * size)

	tx.TxOut[1].Value = tx.TxOut[1].Value - fees

	buf := bytes.NewBuffer(txBytes)
	tx.Serialize(buf)
	txBytes = buf.Bytes()

	return
}

func CreateMultiSigAddress(m int, publicKeys []*btcutil.AddressPubKey, randomBytes []byte) (script, address []byte,
	btcAddressList []string, err error) {

	// ideally m should be
	//	m = len(publicKeys) * 2 /3 ) + 1
	btcAddressList = make([]string, len(publicKeys))
	btcPubKeyMap := make(map[string][]byte)

	for i := range publicKeys {
		address := publicKeys[i].EncodeAddress()

		btcAddressList[i] = address
		btcPubKeyMap[address] = publicKeys[i].ScriptAddress()
	}

	sort.Strings(btcAddressList)

	builder := txscript.NewScriptBuilder().AddInt64(int64(m))
	for _, addr := range btcAddressList {
		builder.AddData(btcPubKeyMap[addr])
	}
	builder.AddInt64(int64(len(publicKeys)))
	builder.AddOp(txscript.OP_CHECKMULTISIG)

	// add randomness
	builder.AddOp(txscript.OP_RETURN)
	builder.AddData(randomBytes)

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
