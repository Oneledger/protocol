/*

 */

package btc

import (
	"encoding/hex"

	"github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/Oneledger/protocol/client"
	"github.com/pkg/errors"
)

func (s *Service) PrepareRedeem(args client.BTCLockRedeemRequest, reply *client.BTCRedeemPrepareResponse) error {
	cd := bitcoin.NewChainDriver(s.blockCypherToken)

	//tracker, err := s.trackerStore.Get("tracker_1")
	tracker, err := s.trackerStore.GetTrackerForLock()
	if err != nil {
		s.logger.Error("error getting tracker for lock", err)
		return errors.Wrap(err, "error getting tracker for lock")
	}

	addr, err := hex.DecodeString(args.Address)

	txnBytes := cd.PrepareRedeemNew(tracker.ProcessTxId, 0, tracker.CurrentBalance,
		addr, args.Amount, tracker.ProcessLockScriptAddress)

	reply.Txn = hex.EncodeToString(txnBytes)
	reply.TrackerName = tracker.Name

	return nil
}
