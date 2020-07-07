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

func (svc *Service) GetValidatorMaturityAmount(req client.RewardsRequest, resp *client.MatureRewardsReply) error {
	validatorAddr := keys.Address{}
	err := validatorAddr.UnmarshalText([]byte(req.Validator))
	if err != nil {
		return err
	}

	amount, err := svc.rewardMaster.RewardCm.GetMaturedBalance(validatorAddr)
	if err != nil {
		return err
	}

	*resp = client.MatureRewardsReply{
		Validator:     validatorAddr,
		MatureRewards: *amount,
	}

	return nil
}
