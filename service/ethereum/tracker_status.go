package ethereum

import (
	"github.com/Oneledger/protocol/data/ethereum"
	codes "github.com/Oneledger/protocol/status_codes"
)

func (svc *Service) GetTrackerStatus(req TrackerStatusRequest, out *TrackerStatusReply) error {
	tracker, err := svc.trackers.QueryAllStores(req.TrackerName)
	if err != nil {
		return codes.ErrGettingTrackerStatusSuccess.Wrap(codes.ErrGettingTrackerStatusFailed).Wrap(codes.ErrGettingTrackerStatusOngoing)
	}

	*out = TrackerStatusReply{
		Status: tracker.State.String(),
	}
	return nil
}

func (svc *Service) GetFailedTrackerStatus(req TrackerStatusRequest, out *TrackerStatusReply) error {
	tracker, err := svc.trackers.WithPrefixType(ethereum.PrefixFailed).Get(req.TrackerName)
	if err != nil {
		//svc.logger.Error(err, codes.ErrGettingTrackerStatusFailed.ErrorMsg())
		return codes.ErrGettingTrackerStatusFailed
	}
	*out = TrackerStatusReply{
		Status: tracker.State.String(),
	}
	return nil
}

func (svc *Service) GetSuccessTrackerStatus(req TrackerStatusRequest, out *TrackerStatusReply) error {
	tracker, err := svc.trackers.WithPrefixType(ethereum.PrefixPassed).Get(req.TrackerName)
	if err != nil {
		//svc.logger.Error(err, codes.ErrGettingTrackerStatusSuccess.ErrorMsg())
		return codes.ErrGettingTrackerStatusSuccess
	}
	*out = TrackerStatusReply{
		Status: tracker.State.String(),
	}
	return nil
}
