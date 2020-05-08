package governance

import (
	"fmt"
	"math/big"

	"github.com/Oneledger/protocol/data/keys"
)

type (
	ProposalAmount  big.Int
	ProposalID      string
	ProposalType    int
	ProposalStatus  int
	ProposalOutcome int
	ProposalState   int
)

type ProposalFund struct {
	id            ProposalID
	address       keys.Address
	fundingAmount ProposalAmount
}

func (fund *ProposalFund) Print() {
	fmt.Printf("Proposal ID : %s | Funding Address : %s  | Funding Amount  : %s \n", fund.id, fund.address.String(), fund.fundingAmount.String())
}
