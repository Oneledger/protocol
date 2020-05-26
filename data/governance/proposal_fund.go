package governance

import (
	"fmt"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
)

type ProposalFund struct {
	id            ProposalID
	address       keys.Address
	fundingAmount *balance.Amount
}

func (fund *ProposalFund) Print() {
	fmt.Printf("Proposal ID : %s | Funding Address : %s  | Funding Amount  : %s \n", fund.id, fund.address.String(), fund.fundingAmount.String())
}

func GetCurrentFunds(id ProposalID, store *ProposalFundStore) *balance.Amount {
	funds := store.GetFundersForProposalID(id, func(proposalID ProposalID, fundingAddr keys.Address, amt *balance.Amount) ProposalFund {
		return ProposalFund{
			id:            proposalID,
			address:       fundingAddr,
			fundingAmount: amt,
		}
	})
	totalBalance := balance.NewAmountFromInt(0)
	for _, fund := range funds {
		totalBalance = totalBalance.Plus(fund.fundingAmount)
	}
	return totalBalance
}

func DeleteAllFunds(id ProposalID, store *ProposalFundStore) error {
	panic("Implement function to delete all funds")
	return nil
}
