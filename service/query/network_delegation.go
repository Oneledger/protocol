package query

import (
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
)

func (svc *Service) GetUndelegatedAmount(req client.GetUndelegatedRequest, reply *client.GetUndelegatedReply) error {
	pendingAmount := make([]client.SinglePendingAmount, 0)
	// iterate every pending amount entry

	// get total amount for all the pending amount
	totalAmount := &balance.Amount{}
	for _, amount := range pendingAmount {
		totalAmount = totalAmount.Plus(amount.Amount)
	}

	*reply = client.GetUndelegatedReply{
		PendingAmount: pendingAmount,
		TotalPendingAmount: *totalAmount,
		Height:        svc.store.GetState().Version(),
	}
	return nil
}
