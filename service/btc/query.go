/*

 */

package btc

import (
	"encoding/json"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/status_codes"
)

func (s *Service) GetTracker(args client.BTCGetTrackerRequest, reply *client.BTCGetTrackerReply) error {

	tracker, err := s.trackerStore.Get(args.Name)
	if err != nil {
		return status_codes.ErrSerialization
	}

	b, _ := json.MarshalIndent(tracker, "", "	")
	reply.TrackerData = string(b)

	return nil
}

func (s *Service) GetMaxTrackerBalance(args client.EmptyRequest, reply *client.MaxTrackerBalanceReply) error {

	reply.MaxBalance = s.trackerStore.GetMaxAvailableBalance()

	return nil
}
