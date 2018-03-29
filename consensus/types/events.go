package types


type Event interface {

}

type EventDataNewBlock struct {
	Block *Block `json:"block"`
}

type EventDataNewProposal struct {
	Proposal *Proposal `json:"proposal"`
}

type EventDataNewPrevote struct {
	Vote *Vote `json:"vote"`
}

