/*

 */

package btc

import (
	"bytes"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/btc"
	"github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/Oneledger/protocol/client"
	btc_data "github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/serialize"
	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	codes "github.com/Oneledger/protocol/status_codes"
)

func (s *Service) PrepareLock(hash string, index uint32, net *chaincfg.Params) ([]byte, error) {
	cd := bitcoin.NewChainDriver("")

	btc := gobcy.API{"dd53aae66b83431ca57a1f656af8ed69", "btc", "main"}
	tx, err := btc.GetTX(hash, nil)
	if err != nil {
		return nil, err
	}

	if tx.Confirmations > 10 {
		return nil, errors.New("source transaction doesn't have enough confirmations")
	}

	hashh, _ := chainhash.NewHashFromStr(tx.Hash)
	balance := int64(tx.Outputs[index].Value)

	tracker, err := s.trackerStore.GetTrackerForLock()
	if err != nil {
		return nil, errors.Wrap(err, "error getting tracker for lock")
	}

	source := btc_data.NewUTXO(hashh, index, balance)

	utxoInit := btc_data.NewUTXO(&chainhash.Hash{}, 0, 0)

	txn := cd.PrepareLock(utxoInit, source, tracker.ProcessLockScriptAddress)

	return txn, nil
}

func (s *Service) AddUserSignatureAndBroadcast(args *client.BTCLockRequest, reply *client.SendTxReply) error {
	cd := bitcoin.NewChainDriver("")

	newTx := cd.AddUserLockSignature(args.Txn, args.Signature)
	lockAmount := newTx.TxOut[0].Value

	var txBytes []byte
	buf := bytes.NewBuffer(txBytes)
	newTx.Serialize(buf)
	txBytes = buf.Bytes()

	lock := btc.Lock{
		Locker:      args.Address,
		TrackerName: args.TrackerName,
		BTCTxn:      txBytes,
		LockAmount:  lockAmount,
	}

	data, err := lock.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.Fee, args.Gas}
	tx := &action.RawTx{
		Type: action.BTC_LOCK,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return codes.ErrSerialization
	}

	*reply = client.SendTxReply{
		RawTx: packet,
	}
	return nil
}
