package query

import (
	"math/big"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	net_delg "github.com/Oneledger/protocol/data/network_delegation"
	"github.com/pkg/errors"
)

func (svc *Service) GetUndelegatedAmount(req client.GetUndelegatedRequest, reply *client.GetUndelegatedReply) error {
	pendingAmounts := make([]client.SinglePendingAmount, 0)
	// get all pending amount
	nd := svc.netwkDelegators.Deleg
	nd.WithPrefix(net_delg.PendingType)
	nd.IterateAllPendingAmounts(func(height int64, addr *keys.Address, coin *balance.Coin) bool {
		if addr.Equal(req.Delegator) {
			pending := client.SinglePendingAmount{
				Amount:       *coin.Amount,
				MatureHeight: height,
			}
			pendingAmounts = append(pendingAmounts, pending)
		}
		return false
	})

	// get matured amount
	nd.WithPrefix(net_delg.MatureType)
	maturedCoin, err := nd.Get(req.Delegator)
	if err != nil {
		return err
	}
	maturedAmount := maturedCoin.Amount
	// get total amount
	totalAmount := balance.NewAmountFromBigInt(big.NewInt(0))
	for _, amount := range pendingAmounts {
		totalAmount = totalAmount.Plus(amount.Amount)
	}
	totalAmount = totalAmount.Plus(*maturedAmount)

	*reply = client.GetUndelegatedReply{
		PendingAmounts: pendingAmounts,
		MaturedAmount:  *maturedAmount,
		TotalAmount:    *totalAmount,
		Height:         nd.GetState().Version(),
	}
	return nil
}

func (svc *Service) GetTotalNetwkDelegation(_ client.GetTotalNetwkDelegation, reply *client.GetTotalNetwkDelgReply) error {
	// get active delegation amount
	poolList, err := svc.governance.GetPoolList()
	if err != nil {
		return err
	}
	if _, ok := poolList["DelegationPool"]; !ok {
		return errors.New("failed to get network delegation pool")
	}
	delagationPool := poolList["DelegationPool"]

	activeBalance, err := svc.balances.GetBalance(delagationPool, svc.currencies)
	if err != nil {
		return err
	}
	currencyOLT, ok := svc.currencies.GetCurrencyByName("OLT")
	if ok != true {
		return errors.New("failed to get OLT from currency set")
	}
	activeCoin := activeBalance.GetCoin(currencyOLT)

	// get pending delegation amount
	pendingCoin := balance.Coin{Currency: currencyOLT, Amount: balance.NewAmountFromBigInt(big.NewInt(0))}
	nd := svc.netwkDelegators.Deleg
	nd.WithPrefix(net_delg.PendingType)
	nd.IterateAllPendingAmounts(func(height int64, addr *keys.Address, coin *balance.Coin) bool {
		pendingCoin = pendingCoin.Plus(*coin)
		return false
	})

	// get matured delegation amount
	nd.WithPrefix(net_delg.MatureType)
	maturedCoin := balance.Coin{Currency: currencyOLT, Amount: balance.NewAmountFromBigInt(big.NewInt(0))}
	nd.IterateMatureAmounts(func(addr *keys.Address, coin *balance.Coin) bool {
		maturedCoin = maturedCoin.Plus(*coin)
		return false
	})

	// get total delegation amount
	totalCoin := activeCoin.Plus(pendingCoin).Plus(maturedCoin)

	*reply = client.GetTotalNetwkDelgReply{
		ActiveAmount:  *activeCoin.Amount,
		PendingAmount: *pendingCoin.Amount,
		MaturedAmount: *maturedCoin.Amount,
		TotalAmount:   *totalCoin.Amount,
		Height:        nd.GetState().Version(),
	}
	return nil
}

func (svc *Service) GetDelegRewards(req client.GetDelegRewardsRequest, resp *client.GetDelegRewardsReply) error {
	height := svc.netwkDelegators.Rewards.GetState().Version()
	options, err := svc.govern.GetNetworkDelegOptions()
	if err != nil {
		return net_delg.ErrGettingDelgOption
	}

	balance, err := svc.netwkDelegators.Rewards.GetRewardsBalance(req.Delegator)
	if err != nil {
		return err
	}
	matured, err := svc.netwkDelegators.Rewards.GetMaturedRewards(req.Delegator)
	if err != nil {
		return err
	}

	var pending *net_delg.DelegPendingRewards
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
