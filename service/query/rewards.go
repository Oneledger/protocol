package query

import (
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
)

func (svc *Service) ListRewardsForValidator(req client.ListRewardsRequest, resp *client.ListRewardsReply) error {
	validatorAddr := keys.Address{}
	err := validatorAddr.UnmarshalText([]byte(req.Validator))
	if err != nil {
		return err
	}
	var rewards []balance.Amount
	svc.rewardStore.Iterate(validatorAddr, func(c string, amt balance.Amount) bool {
		rewards = append(rewards, amt)
		return false
	})

	*resp = client.ListRewardsReply{
		Validator: validatorAddr,
		Rewards:   rewards,
	}

	return nil
}
