/*

 */

package btc

import (
	"bytes"
	"encoding/hex"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/btc"
	"github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/serialize"
	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	codes "github.com/Oneledger/protocol/status_codes"
)

const (
	MINIMUM_CONFIRMATIONS_REQ = 1
)

func (s *Service) PrepareLock(args client.BTCLockPrepareRequest, reply *client.BTCLockPrepareResponse) error {
	cd := bitcoin.NewChainDriver(s.blockCypherToken)

	btc := gobcy.API{s.blockCypherToken, "btc", s.btcChainType}
	tx, err := btc.GetTX(args.Hash, nil)
	if err != nil {
		s.logger.Error("error in getting txn from bitcoin network", err)
		return err
	}

	if tx.Confirmations < MINIMUM_CONFIRMATIONS_REQ {

		s.logger.Error("not enough txn confirmations", err)
		return errors.New("source transaction doesn't have enough confirmations")
	}

	hashh, _ := chainhash.NewHashFromStr(tx.Hash)
	inputAmount := int64(tx.Outputs[args.Index].Value)

	//tracker, err := s.trackerStore.Get("tracker_1")
	tracker, err := s.trackerStore.GetTrackerForLock()
	if err != nil {
		s.logger.Error("error getting tracker for lock", err)
		return errors.Wrap(err, "error getting tracker for lock")
	}

	txnBytes := cd.PrepareLockNew(tracker.ProcessTxId, 0, tracker.CurrentBalance,
		hashh, args.Index, inputAmount, tracker.ProcessLockScriptAddress)

	reply.Txn = hex.EncodeToString(txnBytes)
	reply.TrackerName = tracker.Name

	return nil
}

func (s *Service) AddUserSignatureAndProcessLock(args client.BTCLockRequest, reply *client.SendTxReply) error {

	tracker, err := s.trackerStore.Get(args.TrackerName)
	if err != nil {
		// tracker of that name not found
		return err
	}
	if tracker.IsBusy() {
		// tracker not available anymore, try another tracker
		return err
	}

	// initialize a chain driver for adding signature
	cd := bitcoin.NewChainDriver("")

	// add the users' btc signature to the lock txn in the appropriate place

	s.logger.Debug("----", hex.EncodeToString(args.Txn), hex.EncodeToString(args.Signature))

	newBTCTx := cd.AddUserLockSignature(args.Txn, args.Signature)

	totalLockAmount := newBTCTx.TxOut[0].Value

	if len(newBTCTx.TxIn) == 1 { // if new tracker

		if tracker.CurrentTxId != nil {
			// incorrect txn
			return err
		}
	} else if len(newBTCTx.TxIn) == 2 { // if not a new tracker

		if *tracker.CurrentTxId != newBTCTx.TxIn[0].PreviousOutPoint.Hash ||
			newBTCTx.TxIn[0].PreviousOutPoint.Index != 0 {

			// incorrect txn
			return err
		}
	} else {
		// incorrect txn
		return err
	}

	var txBytes []byte
	buf := bytes.NewBuffer(txBytes)
	err = newBTCTx.Serialize(buf)
	if err != nil {
		return err
	}
	txBytes = buf.Bytes()

	s.logger.Debug("-----", hex.EncodeToString(txBytes))

	lock := btc.Lock{
		Locker:      args.Address,
		TrackerName: args.TrackerName,
		BTCTxn:      txBytes,
		LockAmount:  totalLockAmount,
	}

	data, err := lock.Marshal()
	if err != nil {
		return err
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.GasPrice, args.Gas}
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
