/*

 */

package bitcoin

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"

	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
)

type TxHash [chainhash.HashSize]byte

const (
	BLOCK_CONFIRMATIONS      = 6
	BLOCK_CONFIRMATIONS_TEST = 1
)

type ChainDriver interface {
	//	PrepareLock(prevLock, input *bitcoin.UTXO, lockScriptAddress []byte) (txBytes []byte)

	PrepareLockNew(prevLockTxID *chainhash.Hash, prevLockIndex uint32, prevLockBalance int64,
		inputs []InputTransaction, feesInSatoshi int64, lockAmount int64, returnAddress []byte,
		lockScriptAddress []byte) (txBytes []byte, err error)

	AddLockSignature([]byte, []byte, bool) *wire.MsgTx

	BroadcastTx(*wire.MsgTx, *rpcclient.Client) (*chainhash.Hash, error)

	CheckFinality(*chainhash.Hash, string, string) (bool, error)

	PrepareRedeemNew(prevLockTxID *chainhash.Hash, prevLockIndex uint32, prevLockBalance int64,
		userAddress []byte, redeemAmount int64, feesInSatoshi int64,
		lockScriptAddress []byte) (txBytes []byte)
}

type chainDriver struct {
	blockCypherToken string
}

type InputTransaction struct {
	TxID    *chainhash.Hash
	Index   uint32
	Balance int64
}

var _ ChainDriver = &chainDriver{}

func NewChainDriver(token string) ChainDriver {

	return &chainDriver{token}
}

func (c *chainDriver) PrepareLockNew(prevLockTxID *chainhash.Hash, prevLockIndex uint32, prevLockBalance int64,
	inputs []InputTransaction, feesRate int64, lockAmount int64, returnAddress []byte,
	lockScriptAddress []byte) (txBytes []byte, err error) {

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

	var totalInputBalance int64

	for _, input := range inputs {
		// create users' txin & add to txn
		userOP := wire.NewOutPoint(input.TxID, input.Index)
		vin2 := wire.NewTxIn(userOP, nil, nil)
		tx.AddTxIn(vin2)

		totalInputBalance += input.Balance
	}

	var balance = lockAmount

	// add a txout to the lockscriptaddress
	out := wire.NewTxOut(balance, lockScriptAddress)
	tx.AddTxOut(out)

	feesInSatoshi := int64(EstimateTxSizeBeforeUserSign(tx, true)+20) * feesRate

	var returnBalance = prevLockBalance + totalInputBalance - lockAmount - feesInSatoshi
	if returnBalance < 0 {
		err = errors.New("not enough balance to lock")
		return
	}

	if returnBalance > 0 {
		returnOut := wire.NewTxOut(returnBalance, returnAddress)
		tx.AddTxOut(returnOut)
	}

	buf := bytes.NewBuffer(txBytes)
	tx.Serialize(buf)
	txBytes = buf.Bytes()

	fmt.Println(hex.EncodeToString(txBytes))
	return
}

func (c *chainDriver) AddLockSignature(txBytes []byte, sigScript []byte, isFirstLock bool) *wire.MsgTx {
	tx := wire.NewMsgTx(wire.TxVersion)

	buf := bytes.NewBuffer(txBytes)
	tx.Deserialize(buf)

	// if first lock & not redeem
	if isFirstLock {
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

	if tx.Confirmations >= BLOCK_CONFIRMATIONS {
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

func CreateMultiSigAddress(m int, publicKeys []*btcutil.AddressPubKey, randomBytes []byte,
	params *chaincfg.Params) (

	script, address []byte, btcAddressList []string, err error) {

	if len(publicKeys) >= 20 {
		err = errors.New("signers should be less than 20")
	}
	// ideally m should be
	//	m = len(publicKeys) * 2 /3 ) + 1
	btcAddressList = make([]string, len(publicKeys))

	btcPubKeyMap := make(map[string][]byte)

	for i := range publicKeys {
		address := base58.Encode(publicKeys[i].ScriptAddress())
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
	// builder.AddOp(txscript.OP_RETURN)
	// builder.AddData(randomBytes)

	script, err = builder.Script()
	if err != nil {
		return
	}

	addressObj, err := btcutil.NewAddressScriptHash(script, params)
	if err != nil {
		return
	}

	address, err = txscript.PayToAddrScript(addressObj)

	return
}
