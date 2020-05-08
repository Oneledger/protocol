package governance

import (
	"bytes"
	"math/big"

	"github.com/Oneledger/protocol/data/keys"
)

func NewAmount(x int64) *ProposalAmount {
	return NewAmountFromInt(x)
}

func NewAmountFromInt(x int64) *ProposalAmount {
	return (*ProposalAmount)(big.NewInt(x))
}

var foundProposals []Funder

func GetFundersForProposalID(pf *ProposalFundStore, id ProposalID, fn func(proposalID ProposalID, fundingAddr keys.Address, amt *ProposalAmount) Funder) bool {
	return pf.Iterate(func(proposalID ProposalID, fundingAddr keys.Address, amt *ProposalAmount) bool {
		if proposalID == id {
			foundProposals = append(foundProposals, fn(proposalID, fundingAddr, amt))
		}
		return false
	})
}
func GetProposalsForFunder(pf *ProposalFundStore, funderAddress keys.Address, fn func(proposalID ProposalID, fundingAddr keys.Address, amt *ProposalAmount) Funder) bool {
	return pf.Iterate(func(proposalID ProposalID, fundingAddr keys.Address, amt *ProposalAmount) bool {
		if bytes.Equal(funderAddress, fundingAddr) {
			foundProposals = append(foundProposals, fn(proposalID, fundingAddr, amt))
		}
		return false
	})
}
