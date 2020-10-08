package query

import (
	"github.com/Oneledger/protocol/client"
	netwkDeleg "github.com/Oneledger/protocol/data/network_delegation"
)

func (svc *Service) GetDelegRewards(req client.GetDelegRewardsRequest, resp *client.GetDelegRewardsReply) error {
	height := svc.netwkDelegators.Rewards.GetState().Version()
	options, err := svc.govern.GetNetworkDelegOptions()
	if err != nil {
		return netwkDeleg.ErrGettingDelgOption
	}

	balance, err := svc.netwkDelegators.Rewards.GetRewardsBalance(req.Delegator)
	if err != nil {
		return err
	}
	matured, err := svc.netwkDelegators.Rewards.GetMaturedRewards(req.Delegator)
	if err != nil {
		return err
	}

	var pending *netwkDeleg.DelegPendingRewards
	if req.InclPending {
		pending, err = svc.netwkDelegators.Rewards.GetPendingRewards(req.Delegator, height, options.RewardsMaturityTime+1)
		if err != nil {
			return err
		}
	}

	*resp = client.GetDelegRewardsReply{
		Balance: *balance,
		Pending: pending.Rewards,
		Matured: *matured,
		Height:  height,
	}

	return nil
}
