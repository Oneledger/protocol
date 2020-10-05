package query

import (
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
	maturedCoin, err := nd.WithPrefix(net_delg.MatureType).Get(req.Delegator)
	if err != nil {
		return errors.New("failed to get matured undelegation amount")
	}
	maturedAmount := maturedCoin.Amount
	// get total amount
	totalAmount := &balance.Amount{}
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
