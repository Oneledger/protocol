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

func (svc *Service) GetValidatorMaturityAmount(req client.RewardsRequest, resp *client.RewardsReply) error {
	validatorAddr := keys.Address{}
	err := validatorAddr.UnmarshalText([]byte(req.Validator))
	if err != nil {
		return err
	}

	amount, err := svc.rewardMaster.RewardCumula.GetMaturedBalance(validatorAddr)
	if err != nil {
		return err
	}

	*resp = client.RewardsReply{
		Validator: validatorAddr,
		Amount:    *amount,
	}

	return nil
}

func (svc *Service) GetTotalRewardsForValidator(req client.RewardsRequest, reply *client.RewardsReply) error {
	validatorAddr := keys.Address{}
	err := validatorAddr.UnmarshalText([]byte(req.Validator))
	if err != nil {
		return err
	}
	//Get Total Mature rewards. Mature balance + withdrawn amount
	matureAmount, err := svc.rewardMaster.RewardCumula.GetMaturedRewards(validatorAddr)
	if err != nil {
		return err
	}

	pendingAmount, err := svc.rewardMaster.Reward.GetLastTwoChunks(validatorAddr)
	if err != nil {
		return err
	}

	total := matureAmount.Plus(*pendingAmount)

	*reply = client.RewardsReply{
		Validator: validatorAddr,
		Amount:    *total,
	}

	return nil
}
