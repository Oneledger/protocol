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

func GetFundersForProposalID(pf *ProposalFundStore, id ProposalID, fn func(proposalID ProposalID, fundingAddr keys.Address, amt *ProposalAmount) Funder) []Funder {
	var foundProposals []Funder
	pf.Iterate(func(proposalID ProposalID, fundingAddr keys.Address, amt *ProposalAmount) bool {
		if proposalID == id {
			foundProposals = append(foundProposals, fn(proposalID, fundingAddr, amt))
		}
		return false
	})
	return foundProposals
}
func GetProposalsForFunder(pf *ProposalFundStore, funderAddress keys.Address, fn func(proposalID ProposalID, fundingAddr keys.Address, amt *ProposalAmount) Funder) []Funder {
	var foundProposals []Funder
	pf.Iterate(func(proposalID ProposalID, fundingAddr keys.Address, amt *ProposalAmount) bool {
		if bytes.Equal(keys.Address(funderAddress.String()), fundingAddr) {
			foundProposals = append(foundProposals, fn(proposalID, fundingAddr, amt))
		}
		return false
	})
	return foundProposals
}
