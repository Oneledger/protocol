/*

 */

package btc

import (
	"bytes"
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/btc"
	"github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"
)

func (s *Service) PrepareRedeem(args client.BTCLockRedeemRequest, reply *client.BTCRedeemPrepareResponse) error {
	cd := bitcoin.NewChainDriver(s.blockCypherToken)

	//tracker, err := s.trackerStore.Get("tracker_1")
	tracker, err := s.trackerStore.GetTrackerForRedeem()
	if err != nil {

		s.logger.Error("error getting tracker for lock", err)
		return errors.Wrap(err, "error getting tracker for lock")
	}

	addr, err := hex.DecodeString(args.Address)
	if err != nil {

	}

	txnBytes := cd.PrepareRedeemNew(tracker.ProcessTxId, 0, tracker.CurrentBalance,
		addr, args.Amount, args.FeesBTC, tracker.ProcessLockScriptAddress)

	reply.Txn = hex.EncodeToString(txnBytes)
	reply.TrackerName = tracker.Name

	return nil
}

func (s *Service) AddUserSignatureAndProcessRedeem(args client.BTCLockRequest, reply *client.SendTxReply) error {

	tracker, err := s.trackerStore.Get(args.TrackerName)
	if err != nil {
		// tracker of that name not found
		return codes.ErrTrackerNotFound
	}

	if tracker.IsBusy() {
		// tracker not available anymore, try another tracker
		return codes.ErrTrackerBusy
	}

	// initialize a chain driver for adding signature
	cd := bitcoin.NewChainDriver(s.blockCypherToken)

	// add the users' btc signature to the lock txn in the appropriate place
	s.logger.Debug("----", hex.EncodeToString(args.Txn), hex.EncodeToString(args.Signature))

	newBTCTx := cd.AddUserLockSignature(args.Txn, args.Signature)

	totalRedeemAmount := tracker.CurrentBalance - newBTCTx.TxOut[0].Value

	if len(newBTCTx.TxIn) != 1 { // if new tracker

		return codes.ErrBadBTCTxn
	}

	if len(newBTCTx.TxOut) != 2 { // if not a new tracker

		// incorrect txn
		return codes.ErrBadBTCTxn
	}

	if *tracker.CurrentTxId != newBTCTx.TxIn[0].PreviousOutPoint.Hash ||
		newBTCTx.TxIn[0].PreviousOutPoint.Index != 0 {
		// incorrect txn
		return codes.ErrBadBTCTxn
	}

	var txBytes []byte
	buf := bytes.NewBuffer(txBytes)
	err = newBTCTx.Serialize(buf)
	if err != nil {
		return codes.ErrSerialization
	}
	txBytes = buf.Bytes()

	s.logger.Debug("-----", hex.EncodeToString(txBytes))

	redeem := btc.Redeem{
		Redeemer:     args.Address,
		TrackerName:  args.TrackerName,
		BTCTxn:       txBytes,
		RedeemAmount: totalRedeemAmount,
	}

	data, err := redeem.Marshal()
	if err != nil {
		return codes.ErrSerialization
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
