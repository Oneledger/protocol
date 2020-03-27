/*

 */

package btc

import (
	"bytes"
	"encoding/hex"

	"github.com/blockcypher/gobcy"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/google/uuid"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/btc"
	"github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"
)

const (
	MINIMUM_CONFIRMATIONS_REQ      = 6
	MINIMUM_CONFIRMATIONS_REQ_TEST = 1
)

// PrepareLock
func (s *Service) PrepareLock(args client.BTCLockPrepareRequest, reply *client.BTCLockPrepareResponse) error {

	cfg := s.trackerStore.GetConfig()

	cd := bitcoin.NewChainDriver(cfg.BlockCypherToken)

	btc := gobcy.API{cfg.BlockCypherToken, "btc", cfg.BlockCypherChainType}

	cdInput := make([]bitcoin.InputTransaction, 0, len(args.Inputs))
	var totalInput int64 = 0

	for _, input := range args.Inputs {
		tx, err := btc.GetTX(input.Hash, nil)
		if err != nil {
			s.logger.Error("error in getting txn from bitcoin network", err)
			return codes.ErrBTCReadingTxn
		}

		if tx.Confirmations < MINIMUM_CONFIRMATIONS_REQ {

			s.logger.Error("not enough txn confirmations", err)
			return codes.ErrBTCNotEnoughConfirmations
		}

		if tx.Outputs[input.Index].SpentBy != "" {

			s.logger.Error("source is not spendable", err)
			return codes.ErrBTCNotSpendable
		}

		hashh, _ := chainhash.NewHashFromStr(tx.Hash)
		inputAmount := int64(tx.Outputs[input.Index].Value)
		totalInput += inputAmount

		cdInput = append(cdInput, bitcoin.InputTransaction{hashh, input.Index, inputAmount})
	}

	//tracker, err := s.trackerStore.Get("tracker_1")
	tracker, err := s.trackerStore.GetTrackerForLock()
	if err != nil {
		s.logger.Error("error getting tracker for lock", err)
		return codes.ErrGettingTracker
	}

	s.logger.Infof("%#v \n", tracker)

	returnAddress, err := btcutil.DecodeAddress(args.ReturnAddressStr, cfg.BTCParams)
	if err != nil {
		return codes.ErrBadBTCAddress.Wrap(err)
	}

	returnAddressBytes, err := txscript.PayToAddrScript(returnAddress)
	if err != nil {
		return codes.ErrBadBTCAddress.Wrap(err)
	}

	txnBytes, err := cd.PrepareLockNew(tracker.CurrentTxId, 0, tracker.CurrentBalance,
		cdInput, args.FeeRate, args.AmountSatoshi, returnAddressBytes, tracker.ProcessLockScriptAddress)
	if err != nil {
		return codes.ErrBadBTCTxn.Wrap(err)
	}

	reply.Txn = hex.EncodeToString(txnBytes)
	reply.TrackerName = tracker.Name

	return nil
}

func (s *Service) AddUserSignatureAndProcessLock(args client.BTCLockRequest, reply *client.CreateTxReply) error {

	tracker, err := s.trackerStore.Get(args.TrackerName)
	if err != nil {
		// tracker of that name not found
		return codes.ErrTrackerNotFound
	}
	if tracker.IsBusy() {
		// tracker not available anymore, try another tracker
		return codes.ErrTrackerBusy
	}

	cfg := s.trackerStore.GetConfig()

	// add the users' btc signature to the redeem txn in the appropriate place
	s.logger.Debug("----", hex.EncodeToString(args.Txn))

	newBTCTx := wire.NewMsgTx(wire.TxVersion)

	buf := bytes.NewBuffer(args.Txn)
	newBTCTx.Deserialize(buf)

	totalLockAmount := newBTCTx.TxOut[0].Value - tracker.CurrentBalance

	isFirstLock := tracker.CurrentTxId == nil
	if isFirstLock {
		// if this is first lock for tracker, then all inputs must be signed

		for i := range newBTCTx.TxIn {
			if len(newBTCTx.TxIn[i].SignatureScript) == 0 {

				s.logger.Error("all user sources for lock are not signed")
				return codes.ErrBadBTCTxn
			}
		}
	} else {

		// if not the first tracker txn then the first input should be
		// tracker previous txn
		if *tracker.CurrentTxId != newBTCTx.TxIn[0].PreviousOutPoint.Hash ||
			newBTCTx.TxIn[0].PreviousOutPoint.Index != 0 {

			// incorrect txn
			s.logger.Error("btc txn doesn;t match tracker")
			return codes.ErrBadBTCTxn
		}

		for i := range newBTCTx.TxIn {
			if i == 0 {
				continue
			}

			if len(newBTCTx.TxIn[i].SignatureScript) == 0 {

				s.logger.Error("all user sources for lock are not signed")
				return codes.ErrBadBTCTxn
			}
		}
	}

	if !bitcoin.ValidateLock(newBTCTx, cfg.BlockCypherToken, cfg.BlockCypherChainType,
		tracker.ProcessLockScriptAddress, tracker.CurrentBalance, totalLockAmount, isFirstLock) {

		return codes.ErrBadBTCTxn
	}

	var txBytes []byte
	buf = bytes.NewBuffer(txBytes)
	err = newBTCTx.Serialize(buf)
	if err != nil {
		return codes.ErrSerialization
	}
	txBytes = buf.Bytes()

	lock := btc.Lock{
		Locker:      args.Address,
		TrackerName: args.TrackerName,
		BTCTxn:      txBytes,
		LockAmount:  totalLockAmount,
	}

	data, err := lock.Marshal()
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

	*reply = client.CreateTxReply{
		RawTx: packet,
	}
	return nil
}
