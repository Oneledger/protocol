package query

import (
	"fmt"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/pkg/errors"
	"math/big"
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
	fmt.Println("pendingAmounts", pendingAmounts)

	// get matured amount
	//maturedCoin, err := nd.WithPrefix(net_delg.MatureType).Get(req.Delegator)
	//if err != nil {
	//	return errors.New("failed to get matured undelegation amount")
	//}
	//maturedAmount := maturedCoin.Amount
	maturedAmount := balance.NewAmountFromBigInt(big.NewInt(0))
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

func (svc *Service) GetTotalNetwkDelegation(reply *client.GetTotalNetwkDelgReply) error {
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
	svc.netwkDelegators.Deleg.IterateAllPendingAmounts(func(height int64, addr *keys.Address, coin *balance.Coin) bool {
		pendingCoin = pendingCoin.Plus(*coin)
		return false
	})

	// get matured delegation amount
	maturedCoin := balance.Coin{Currency: currencyOLT, Amount: balance.NewAmountFromBigInt(big.NewInt(0))}
	svc.netwkDelegators.Deleg.IterateMatureAmounts(func(addr *keys.Address, coin *balance.Coin) bool {
		maturedCoin = maturedCoin.Plus(*coin)
		return false
	})

	// get total delegation amount
	totalCoin := activeCoin.Plus(pendingCoin).Plus(maturedCoin)

	*reply = client.GetTotalNetwkDelgReply{
		ActiveAmount: *activeCoin.Amount,
		PendingAmount: *pendingCoin.Amount,
		MaturedAmount: *maturedCoin.Amount,
		TotalAmount: *totalCoin.Amount,
	}
	return nil
}