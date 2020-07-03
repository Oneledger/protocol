package governance

import (
	"fmt"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/pkg/errors"
)

type ProposalFund struct {
	id            ProposalID
	address       keys.Address
	fundingAmount *balance.Amount
}

func (fund *ProposalFund) Print() {
	fmt.Printf("Proposal ID : %s | Funding Address : %s  | Funding Amount  : %s \n", fund.id, fund.address.String(), fund.fundingAmount.String())
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

	// set total funds record to 0
	keyTotal := assembleTotalFundsKey(id)
	zeroFunds := balance.NewAmount(0)
	err := pf.set(keyTotal, *zeroFunds)
	if err != nil {
		return errors.Wrap(err, errorSettingRecord)
	}
	return nil
}
