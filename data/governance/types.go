package governance

import "github.com/pkg/errors"

type (
	ProposalID      string
	ProposalType    int
	ProposalStatus  int
	ProposalOutcome int
	ProposalState   int
	VoteOpinion     int
	VoteResult      int
)

func IDFromString(s string) ProposalID {
	return ProposalID(s)
}

func (opinion VoteOpinion) String() string {
	switch opinion {
	case OPIN_UNKNOWN:
		return "Unknown"
	case OPIN_POSITIVE:
		return "Positive"
	case OPIN_NEGATIVE:
		return "Negative"
	case OPIN_GIVEUP:
		return "Giveup"
	default:
		return "Invalid opinion"
	}
}

func (opinion VoteOpinion) Err() error {
	opName := opinion.String()
	if opName == "" {
		return errors.New("vote opinion must be one of [UNKNOWN, POSITIVE, NEGATIVE, GIVEUP]")
	}
	return nil
}

//func (p ProposalAmount) MarshalJSON() ([]byte, error) {
//	v := p.BigInt().String()
//	return json.Marshal(v)
//}
//
//func (p *ProposalAmount) UnmarshalJSON(b []byte) error {
//	v := ""
//	err := json.Unmarshal(b, &v)
//	if err != nil {
//		return err
//	}
//	i, ok := big.NewInt(0).SetString(v, 0)
//	if !ok {
//		return errors.New("failed to unmarshal amount" + v)
//	}
//	*p = *(*ProposalAmount)(i)
//	return nil
//}
//
//func (p *ProposalAmount) BigInt() *big.Int {
//	return (*big.Int)(p)
//}
//
//func (p ProposalAmount) String() string {
//	return p.BigInt().String()
//}
//
//func (p ProposalAmount) Plus(add *ProposalAmount) *ProposalAmount {
//	base := big.NewInt(0)
//	return (*ProposalAmount)(base.Add(p.BigInt(), add.BigInt()))
//}
//func NewAmount(x int64) *ProposalAmount {
//	return NewAmountFromInt(x)
//}
//
//func NewAmountFromInt(x int64) *ProposalAmount {
//	return (*ProposalAmount)(big.NewInt(x))
//}
//
//func NewAmountFromBigInt(x *big.Int) *ProposalAmount {
//	return (*ProposalAmount)(x)
//}
