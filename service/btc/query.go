/*

 */

package btc

import "github.com/Oneledger/protocol/client"

func (s *Service) GetTracker(args client.BTCGetTrackerRequest, reply *client.BTCGetTrackerReply) error {

	tracker, err := s.trackerStore.Get(args.Name)
	if err != nil {
		return err
	}

	reply.Tracker = *tracker

	return nil
}
