package query

import (
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
)

func (svc *Service) ListRewardsForValidator(req client.RewardsRequest, resp *client.ListRewardsReply) error {
	validatorAddr := keys.Address{}
	err := validatorAddr.UnmarshalText([]byte(req.Validator))
	if err != nil {
		return err
	}
	var rewards []balance.Amount
	svc.rewardMaster.Reward.Iterate(validatorAddr, func(c string, amt balance.Amount) bool {
		rewards = append(rewards, amt)
		return false
	})

	*resp = client.ListRewardsReply{
		Validator: validatorAddr,
		Rewards:   rewards,
	}

	return nil
}

func (svc *Service) GetTotalRewardsForValidator(req client.RewardsRequest, reply *client.ValidatorRewardStats) error {
	*reply = client.ValidatorRewardStats{}
	validatorAddr := keys.Address{}
	err := validatorAddr.UnmarshalText([]byte(req.Validator))
	if err != nil {
		return err
	}

	matureBalance, err := svc.rewardMaster.RewardCm.GetMaturedBalance(validatorAddr)
	if err != nil || matureBalance == nil {
		return err
	}
	//Get Withdrawn Amount
	withdrawnAmount, err := svc.rewardMaster.RewardCm.GetWithdrawnRewards(validatorAddr)
	if err != nil || withdrawnAmount == nil {
		return err
	}

	//Get Total Mature rewards. Mature balance + withdrawn amount
	matureAmount, err := svc.rewardMaster.RewardCm.GetMaturedRewards(validatorAddr)
	if err != nil || matureAmount == nil {
		return err
	}
	//Get Amount pending in the last 2 chunks
	pendingAmount, err := svc.rewardMaster.Reward.GetLastTwoChunks(validatorAddr)
	if err != nil || pendingAmount == nil {
		return err
	}

	total := matureAmount.Plus(*pendingAmount)

	*reply = client.ValidatorRewardStats{
		Address:         validatorAddr,
		PendingAmount:   *pendingAmount,
		WithdrawnAmount: *withdrawnAmount,
		MatureBalance:   *matureBalance,
		TotalAmount:     *total,
	}

	return nil
}

func (svc *Service) GetTotalRewards(_ client.RewardsRequest, reply *client.RewardStat) error {
	totalRewards := balance.NewAmount(0)
	validatorList := make([]client.ValidatorRewardStats, 0, 64)

	svc.rewardMaster.Reward.IterateAddrList(func(addr keys.Address) bool {
		validatorStats := &client.ValidatorRewardStats{}
		valRewardReq := client.RewardsRequest{
			Validator: addr.String(),
		}
		err := svc.GetTotalRewardsForValidator(valRewardReq, validatorStats)
		if err != nil {
			return true
		}
		validatorList = append(validatorList, *validatorStats)
		totalRewards = totalRewards.Plus(validatorStats.TotalAmount)
		return false
	})

	*reply = client.RewardStat{
		Validators:   validatorList,
		TotalRewards: *totalRewards,
	}

	return nil
}
