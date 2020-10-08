package query

import (
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/network_delegation"
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
