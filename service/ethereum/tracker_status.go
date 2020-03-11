package ethereum

import (
	codes "github.com/Oneledger/protocol/status_codes"
)

func (svc *Service) GetTrackerStatus(req TrackerStatusRequest, out *TrackerStatusReply) error {
	tracker, err := svc.trackerStore.Get(req.TrackerName)
	if err != nil {
		svc.logger.Error(err, codes.ErrGettingTrackerStatus.ErrorMsg())
		return codes.ErrGettingTrackerStatus
	}
	*out = TrackerStatusReply{
		Status: tracker.State.String(),
	}
	return nil
}
