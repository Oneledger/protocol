package governance

import (
	"encoding/json"
	"math/big"

	"github.com/pkg/errors"
)

type (
	ProposalAmount  big.Int
	ProposalID      string
	ProposalType    int
	ProposalStatus  int
	ProposalOutcome int
	ProposalState   int
)

func (p ProposalAmount) MarshalJSON() ([]byte, error) {
	v := p.BigInt().String()
	return json.Marshal(v)
}

func (p *ProposalAmount) UnmarshalJSON(b []byte) error {
	v := ""
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	i, ok := big.NewInt(0).SetString(v, 0)
	if !ok {
		return errors.New("failed to unmarshal amount" + v)
	}
	*p = *(*ProposalAmount)(i)
	return nil
}

func (p *ProposalAmount) BigInt() *big.Int {
	return (*big.Int)(p)
}

func (p ProposalAmount) String() string {
	return p.BigInt().String()
}

func (p ProposalAmount) Plus(add *ProposalAmount) *ProposalAmount {
	base := big.NewInt(0)
	return (*ProposalAmount)(base.Add(p.BigInt(), add.BigInt()))
}

type Funder struct {
	id            ProposalID
	address       string
	fundingAmount string
}
