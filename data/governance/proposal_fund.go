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

func (pf *ProposalFundStore) GetCurrentFundsForProposal(id ProposalID) *balance.Amount {
	totalBalance := balance.NewAmountFromInt(0)
	pf.GetFundsForProposalID(id, func(proposalID ProposalID, fundingAddr keys.Address, amt *balance.Amount) ProposalFund {
		totalBalance = totalBalance.Plus(amt)
		return ProposalFund{}
	})
	return totalBalance
}

func GetCurrentFundsByFunder(id ProposalID, funder keys.Address, store *ProposalFundStore) (*balance.Amount, error) {
	funds := store.GetFundsForProposalID(id, func(proposalID ProposalID, fundingAddr keys.Address, amt *balance.Amount) ProposalFund {
		return ProposalFund{
			id:            proposalID,
			address:       fundingAddr,
			fundingAmount: amt,
		}
	})
	funderBalance := balance.NewAmountFromInt(0)
	for _, fund := range funds {
		if fund.address.Equal(funder) {
			funderBalance = funderBalance.Plus(fund.fundingAmount)
			return funderBalance, nil
		}
	}
	return nil, ErrWithdrawCheckFundsFailed
}

func (pf *ProposalFundStore) DeleteAllFunds(id ProposalID) error {
	e := error(nil)
	pf.GetFundsForProposalID(id, func(proposalID ProposalID, fundingAddr keys.Address, amt *balance.Amount) ProposalFund {
		ok, err := pf.DeleteFunds(proposalID, fundingAddr)
		if err != nil {
			e = err
			return ProposalFund{}
		}
		if !ok {
			e = ErrDeductFunding
		}
		return ProposalFund{}
	})
	if e != nil {
		return e
	}
	return nil
}
