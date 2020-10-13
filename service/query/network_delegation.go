package query

import (
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/network_delegation"
	"github.com/pkg/errors"
	"math/big"
)

func (svc *Service) ListDelegation(req client.ListDelegationRequest, reply *client.ListDelegationReply) error {
	zeroAmount := balance.Coin{
		Currency: balance.Currency{Id: 0, Name: "OLT", Chain: chain.ONELEDGER, Decimal: 18, Unit: "nue"},
		Amount:   balance.NewAmount(0),
	}
	active, err := svc.networkDelegation.Deleg.WithPrefix(network_delegation.ActiveType).Get(req.DelegationAddress)
	if err != nil {
		active = &zeroAmount
	}
	pending := zeroAmount
	svc.networkDelegation.Deleg.WithPrefix(network_delegation.PendingType)
	svc.networkDelegation.Deleg.IterateAllPendingAmounts(func(height int64, addr *keys.Address, coin *balance.Coin) bool {
		if addr.Equal(req.DelegationAddress) {
			pending = pending.Plus(*coin)
		}
		return false
	})
	mature, err := svc.networkDelegation.Deleg.WithPrefix(network_delegation.MatureType).Get(req.DelegationAddress)
	if err != nil {
		mature = &zeroAmount
		//return err
	}
	ds := client.DelegationStats{
		Active:  active.String(),
		Pending: pending.String(),
		Matured: mature.String(),
	}
	*reply = client.ListDelegationReply{
		DelegationStats: ds,
		Height:          svc.proposalMaster.Proposal.GetState().Version(),
	}
	return nil
}

func (svc *Service) GetUndelegatedAmount(req client.GetUndelegatedRequest, reply *client.GetUndelegatedReply) error {
	pendingAmounts := make([]client.SinglePendingAmount, 0)
	// get all pending amount
	nd := svc.networkDelegation.Deleg
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

	// get matured amount
	nd.WithPrefix(network_delegation.MatureType)
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

func (svc *Service) GetTotalNetwkDelegation(req client.GetTotalNetwkDelegation, reply *client.GetTotalNetwkDelgReply) error {
	nd := svc.networkDelegation.Deleg
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
			MaturedAmount: *balance.NewAmountFromBigInt(big.NewInt(0)),
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

	// get matured delegation amount
	nd.WithPrefix(network_delegation.MatureType)
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
