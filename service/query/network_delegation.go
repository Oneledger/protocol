package query

import (
	"fmt"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/network_delegation"
)

func (svc *Service) ListDelegation(req client.ListDelegationRequest, reply *client.ListDelegationReply) error {

	active, err := svc.networkDelegation.Deleg.WithPrefix(network_delegation.ActiveType).Get(req.DelegationAddress)
	if err != nil {
		fmt.Println("1", req.DelegationAddress.String())
		return err
	}
	pending, err := svc.networkDelegation.Deleg.WithPrefix(network_delegation.PendingType).Get(req.DelegationAddress)
	if err != nil {
		fmt.Println("2")
		return err
	}
	mature, err := svc.networkDelegation.Deleg.WithPrefix(network_delegation.MatureType).Get(req.DelegationAddress)
	if err != nil {
		fmt.Println("3")
		return err
	}
	ds := client.DelegationStats{
		Active:  active.String(),
		Pending: pending.String(),
		Matured: mature.String(),
	}
	fmt.Println(ds)
	*reply = client.ListDelegationReply{
		DelegationStats: ds,
		Height:          svc.proposalMaster.Proposal.GetState().Version(),
	}
	return nil
}
