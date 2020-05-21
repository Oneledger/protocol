package governance

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/serialize"
	"github.com/magiconair/properties/assert"
	"testing"
)

var proposal = Proposal{
	ProposalID:      "c24e6520cbd66adcd17fa9c300b5070b",
	Type:            33,
	Status:          35,
	Outcome:         38,
	Description:     "DESCRIPTION A",
	Proposer:        []byte("0lt9b7ae9be50f7b378b591da41e99843c3d8220004"),
	FundingDeadline: 12,
	FundingGoal:     balance.NewAmount(10000000000),
	VotingDeadline:  12,
	PassPercentage:  51,
}

func TestProposalSerialization(t *testing.T) {
	szlr := serialize.GetSerializer(serialize.PERSISTENT)

	propBytes, err := szlr.Serialize(proposal)
	assert.Equal(t, err, nil)

	prop := &Proposal{}
	err = szlr.Deserialize(propBytes, prop)
	assert.Equal(t, err, nil)

	assert.Equal(t, *prop, proposal)
}
