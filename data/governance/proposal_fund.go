package governance

import (
	"fmt"

	"github.com/Oneledger/protocol/data/keys"
)

type ProposalFund struct {
	id            ProposalID
	address       keys.Address
	fundingAmount ProposalAmount
}

func (fund *ProposalFund) Print() {
	fmt.Printf("Proposal ID : %s | Funding Address : %s  | Funding Amount  : %s \n", fund.id, fund.address.String(), fund.fundingAmount.String())
}
