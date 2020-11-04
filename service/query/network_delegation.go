package query

import (
	"math/big"

	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/network_delegation"
	"github.com/pkg/errors"
)

func (svc *Service) ListDelegation(req client.ListDelegationRequest, reply *client.ListDelegationReply) error {
	zeroAmount := balance.NewAmount(0)
	zeroCoin := balance.Coin{
		Currency: balance.Currency{Id: 0, Name: "OLT", Chain: chain.ONELEDGER, Decimal: 18, Unit: "nue"},
		Amount:   zeroAmount,
	}
	fullDelgStatsList := make([]client.FullDelegStats, 0)
	delegStore := svc.netwkDelegators.Deleg
	delegRewardsStore := svc.netwkDelegators.Rewards
	reply.AllDelegStats = fullDelgStatsList

	// if input is a non-empty list of addresses, get info for them
	for _, address := range req.DelegationAddresses {
		// get delegation stats
		activeDelegation, err := delegStore.WithPrefix(network_delegation.ActiveType).Get(address)
		if err != nil {
			activeDelegation = &zeroCoin
		}
		pendingDelegation := zeroCoin
		delegStore.WithPrefix(network_delegation.PendingType)
		delegStore.IterateAllPendingAmounts(func(height int64, addr *keys.Address, coin *balance.Coin) bool {
			if addr.Equal(address) {
				pendingDelegation = pendingDelegation.Plus(*coin)
			}
			return false
		})
		delegStats := client.DelegStats{
			Active: *activeDelegation.Amount,
			Pending: *pendingDelegation.Amount,
		}

		// get delegation rewards stats
		activeRewards, err := delegRewardsStore.GetRewardsBalance(address)
		if err != nil {
			activeRewards = zeroAmount
		}
		pendingRewards := zeroAmount
		delegRewardsStore.IterateAllPD(func(height int64, addr keys.Address, amt *balance.Amount) bool {
			if addr.Equal(address) {
				pendingRewards = pendingRewards.Plus(*amt)
			}
			return false
		})
		delegRewardsStats := client.DelegRewardsStats{
			Active: *activeRewards,
			Pending: *pendingRewards,
		}

		// combine info into full delegation stats
		fullDelegStats := client.FullDelegStats{
			DelegAddress:      address,
			DelegStats:        delegStats,
			DelegRewardsStats: delegRewardsStats,
		}

		// add into the list
		fullDelgStatsList = append(fullDelgStatsList, fullDelegStats)
	}

	if len(req.DelegationAddresses) > 0 {
		return nil
	}

	// if input is an empty list, get info for all delegators
	// use a map to store all results, key is address, value is the pointer of fullDelegStats object
	delegationMap := make(map[string]*client.FullDelegStats)

	// get delegation stats
	// active amounts
	delegStore.IterateActiveAmounts(func(addr *keys.Address, coin *balance.Coin) bool {
		// load each result into the reply, and put the pointer to the map
		delegStats := client.DelegStats{
			Active: *coin.Amount,
			Pending: *zeroAmount,
		}
		fullDelgStats := client.FullDelegStats{
			DelegAddress: *addr,
			DelegStats: delegStats,
			DelegRewardsStats: client.DelegRewardsStats{},
		}
		reply.AllDelegStats = append(reply.AllDelegStats, fullDelgStats)

		delegationMap[addr.String()] = &fullDelgStats
		return false
	})

	// pending amounts
	// first check if there is already one data entry to this address
	//, if so, add to existing structure, if not, add new data entry with zero active amount
	delegStore.IterateAllPendingAmounts(func(height int64, addr *keys.Address, coin *balance.Coin) bool {
		if fullDelgStats, ok := delegationMap[addr.String()]; ok {
			pendingAmount := coin.Amount
			fullDelgStats.DelegStats.Pending = *fullDelgStats.DelegStats.Pending.Plus(*pendingAmount)
		} else {
			delegStats := client.DelegStats{
				Active: *zeroAmount,
				Pending: *coin.Amount,
			}
			fullDelgStats := client.FullDelegStats{
				DelegAddress: *addr,
				DelegStats: delegStats,
				DelegRewardsStats: client.DelegRewardsStats{},
			}
			delegationMap[addr.String()] = &fullDelgStats
		}
		return false
	})

	// get delegation rewards stats
	// active rewards
	delegRewardsStore.IterateActiveRewards(func(addr *keys.Address, amt *balance.Amount) bool {
		// check if the address is already in the map, load each result into the reply, and put the pointer to the map
		if fullDelgStats, ok := delegationMap[addr.String()]; ok {
			fullDelgStats.DelegRewardsStats.Active = *amt
		} else {
			delegRewardsStats := client.DelegRewardsStats{
				Active: *amt,
				Pending: *zeroAmount,
			}
			fullDelgStats := client.FullDelegStats{
				DelegAddress: *addr,
				DelegStats: client.DelegStats{},
				DelegRewardsStats: delegRewardsStats,
			}
			reply.AllDelegStats = append(reply.AllDelegStats, fullDelgStats)

			delegationMap[addr.String()] = &fullDelgStats
		}
		return false
	})

	// pending withdraw rewards
	delegRewardsStore.IterateAllPD(func(height int64, addr keys.Address, amt *balance.Amount) bool {
		// check if the address is already in the map, load each result into the reply, and put the pointer to the map
		if fullDelgStats, ok := delegationMap[addr.String()]; ok {
			fullDelgStats.DelegRewardsStats.Pending = *fullDelgStats.DelegRewardsStats.Pending.Plus(*amt)
		} else {
			delegRewardsStats := client.DelegRewardsStats{
				Active: *zeroAmount,
				Pending: *amt,
			}
			fullDelgStats := client.FullDelegStats{
				DelegAddress: addr,
				DelegStats: client.DelegStats{},
				DelegRewardsStats: delegRewardsStats,
			}
			reply.AllDelegStats = append(reply.AllDelegStats, fullDelgStats)

			delegationMap[addr.String()] = &fullDelgStats
		}
		return false
	})
	return nil
}

func (svc *Service) GetUndelegatedAmount(req client.GetUndelegatedRequest, reply *client.GetUndelegatedReply) error {
	pendingAmounts := make([]client.SinglePendingAmount, 0)
	// get all pending amount
	nd := svc.netwkDelegators.Deleg
	nd.WithPrefix(network_delegation.PendingType)
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

	// get total amount
	totalAmount := balance.NewAmountFromBigInt(big.NewInt(0))
	for _, amount := range pendingAmounts {
		totalAmount = totalAmount.Plus(amount.Amount)
	}

	*reply = client.GetUndelegatedReply{
		PendingAmounts: pendingAmounts,
		TotalAmount:    *totalAmount,
		Height:         nd.GetState().Version(),
	}
	return nil
}

func (svc *Service) GetTotalNetwkDelegation(req client.GetTotalNetwkDelegation, reply *client.GetTotalNetwkDelgReply) error {
	nd := svc.netwkDelegators.Deleg
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

	if req.OnlyActive == 1 {
		*reply = client.GetTotalNetwkDelgReply{
			ActiveAmount:  *activeCoin.Amount,
			PendingAmount: *balance.NewAmountFromBigInt(big.NewInt(0)),
			TotalAmount:   *activeCoin.Amount,
			Height:        nd.GetState().Version(),
		}
		return nil
	}

	// get pending delegation amount
	pendingCoin := balance.Coin{Currency: currencyOLT, Amount: balance.NewAmountFromBigInt(big.NewInt(0))}
	nd.WithPrefix(network_delegation.PendingType)
	nd.IterateAllPendingAmounts(func(height int64, addr *keys.Address, coin *balance.Coin) bool {
		pendingCoin = pendingCoin.Plus(*coin)
		return false
	})

	// get total delegation amount
	totalCoin := activeCoin.Plus(pendingCoin)

	*reply = client.GetTotalNetwkDelgReply{
		ActiveAmount:  *activeCoin.Amount,
		PendingAmount: *pendingCoin.Amount,
		TotalAmount:   *totalCoin.Amount,
		Height:        nd.GetState().Version(),
	}
	return nil
}

func (svc *Service) GetDelegRewards(req client.GetDelegRewardsRequest, resp *client.GetDelegRewardsReply) error {
	height := svc.netwkDelegators.Rewards.GetState().Version()
	options, err := svc.govern.GetNetworkDelegOptions()
	if err != nil {
		return network_delegation.ErrGettingDelgOption
	}

	balance, err := svc.netwkDelegators.Rewards.GetRewardsBalance(req.Delegator)
	if err != nil {
		return err
	}
	//matured, err := svc.netwkDelegators.Rewards.GetMaturedRewards(req.Delegator)
	//if err != nil {
	//	return err
	//}

	pending := &network_delegation.DelegPendingRewards{Rewards: []*network_delegation.PendingRewards{}}
	if req.InclPending {
		pending, err = svc.netwkDelegators.Rewards.GetPendingRewards(req.Delegator, height, options.RewardsMaturityTime+1)
		if err != nil {
			return err
		}
	}

	*resp = client.GetDelegRewardsReply{
		Balance: *balance,
		Pending: pending.Rewards,
		//Matured: *matured,
		Height:  height,
	}
	return nil
}

func (svc *Service) GetTotalDelegRewards(req client.GetTotalDelegRewardsRequest, resp *client.GetTotalDelegRewardsReply) error {
	height := svc.netwkDelegators.Rewards.GetState().Version()

	total, err := svc.netwkDelegators.Rewards.GetTotalRewards()
	if err != nil {
		return err
	}

	*resp = client.GetTotalDelegRewardsReply{
		TotalRewards: *total,
		Height:       height,
	}
	return nil
}