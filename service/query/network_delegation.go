package query

import (
	"errors"
	"github.com/Oneledger/protocol/client"
)

func (svc *Service) GetTotalNetwkDelg(reply *client.GetTotalNetwkDelgReply) error {
	poolList, err := svc.governance.GetPoolList()
	if err != nil {
		return err
	}
	if _, ok := poolList["DelegationPool"]; !ok {
		return errors.New("failed to get network delegation pool")
	}
	delagationPool := poolList["DelegationPool"]

	balance, err := svc.balances.GetBalance(delagationPool, svc.currencies)
	if err != nil {
		return err
	}
	*reply = client.GetTotalNetwkDelgReply{
		Amount: *balance.Amounts["OLT"].Amount,
	}
	return nil
}
