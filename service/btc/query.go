/*

 */

package btc

import (
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/status_codes"
)

func (s *Service) GetTracker(args client.BTCGetTrackerRequest, reply *client.BTCGetTrackerReply) error {

	tracker, err := s.trackerStore.Get(args.Name)
	if err != nil {
		return status_codes.ErrSerialization
	}

	reply.Tracker = *tracker

	return nil
}
