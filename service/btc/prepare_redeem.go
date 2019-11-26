/*

 */

package btc

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/action/btc"
	"github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/serialize"
	codes "github.com/Oneledger/protocol/status_codes"
)

func (s *Service) PrepareRedeem(args client.BTCRedeemRequest, reply *client.BTCRedeemPrepareResponse) error {
	cd := bitcoin.NewChainDriver(s.blockCypherToken)

	//tracker, err := s.trackerStore.Get("tracker_1")
	tracker, err := s.trackerStore.GetTrackerForRedeem()
	if err != nil {

		s.logger.Error("error getting tracker for lock", err)
		return errors.Wrap(err, "error getting tracker for lock")
	}

	if tracker.CurrentBalance < (args.Amount + args.FeesBTC) {
		return errors.New("not tracker with enough balance")
	}

	btcAddr := base58.Decode(args.BTCAddress)
	if len(btcAddr) == 0 {
		return errors.New("redeem address not base58")
	}

	fmt.Printf("%#v \n", tracker)
	txnBytes := cd.PrepareRedeemNew(tracker.CurrentTxId, 0, tracker.CurrentBalance,
		btcAddr, args.Amount, args.FeesBTC, tracker.ProcessLockScriptAddress)

	fmt.Println(hex.EncodeToString(txnBytes))

	redeem := btc.Redeem{
		Redeemer:     args.Address,
		TrackerName:  tracker.Name,
		BTCTxn:       txnBytes,
		RedeemAmount: args.Amount,
	}

	data, err := redeem.Marshal()
	if err != nil {
		return codes.ErrSerialization
	}

	uuidNew, _ := uuid.NewUUID()
	fee := action.Fee{args.GasPrice, args.Gas}
	tx := &action.RawTx{
		Type: action.BTC_REDEEM,
		Data: data,
		Fee:  fee,
		Memo: uuidNew.String(),
	}

	packet, err := serialize.GetSerializer(serialize.NETWORK).Serialize(tx)
	if err != nil {
		return codes.ErrSerialization
	}

	reply.RawTx = packet
	reply.TrackerName = tracker.Name
	return nil
}
