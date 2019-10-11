/*

 */

package internal

import (
	"bytes"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/btc"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/google/uuid"
)

type signBTCTxnData struct {
	txnData          []byte
	validatorPrivKey *btcec.PrivateKey
	lockscript       []byte
	pubKey           keys.PublicKey
	validatorAddress action.Address

	trackerName string
}

func signBTCTxn(data interface{}) {
	signingData := data.(signBTCTxnData)

	lockTx := wire.NewMsgTx(wire.TxVersion)
	lockTx.Deserialize(bytes.NewReader(signingData.txnData))

	sig, err := txscript.RawTxInSignature(lockTx, 0, signingData.lockscript, txscript.SigHashAll, signingData.validatorPrivKey)
	if err != nil {

	}

	PubKey := signingData.validatorPrivKey.PubKey().SerializeCompressed()
	valPubKey, err := btcutil.NewAddressPubKey(PubKey, &chaincfg.RegressionNetParams)

	addSigData := btc.AddSignature{
		TrackerName:      signingData.trackerName,
		ValidatorPubKey:  valPubKey,
		PubKey:           signingData.pubKey,
		BTCSignature:     sig,
		ValidatorAddress: signingData.validatorAddress,
	}

	dat, err := addSigData.Marshal()

	uuidNew, _ := uuid.NewUUID()
	tx := action.RawTx{
		Type: action.BTC_ADD_SIGNATURE,
		Data: dat,
		Fee:  action.Fee{},
		Memo: uuidNew.String(),
	}
}
