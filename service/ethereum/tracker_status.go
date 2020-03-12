package ethereum

import (
	codes "github.com/Oneledger/protocol/status_codes"
)

func (svc *Service) GetOngoingTrackerStatus(req TrackerStatusRequest, out *TrackerStatusReply) error {
	tracker, err := svc.trackersOngoing.Get(req.TrackerName)
	if err != nil {
		tracker, err = svc.trackersFailed.Get(req.TrackerName)
		if err != nil {
			svc.logger.Error(err, codes.ErrGettingTrackerStatus.ErrorMsg())
			return codes.ErrGettingTrackerStatus
		}
	}
	*out = TrackerStatusReply{
		Status: tracker.State.String(),
	}
	return nil
}

func (svc *Service) GetFailedTrackerStatus(req TrackerStatusRequest, out *TrackerStatusReply) error {
	tracker, err := svc.trackersFailed.Get(req.TrackerName)
	if err != nil {
		svc.logger.Error(err, codes.ErrGettingTrackerStatus.ErrorMsg())
		return codes.ErrGettingTrackerStatus
	}
	*out = TrackerStatusReply{
		Status: tracker.State.String(),
	}
	return nil
}
