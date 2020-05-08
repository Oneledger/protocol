package governance

import (
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type ProposalStore struct {
	state *storage.State
	szlr  serialize.Serializer

	prefix []byte //Current Store Prefix

	prefixActive []byte
	prefixPassed []byte
	prefixFailed []byte

	proposalOptions *ProposalOptions
}

func (ps *ProposalStore) Set(proposal *Proposal) error {
	prefixed := append(ps.prefix, proposal.ProposalID...)
	data, err := ps.szlr.Serialize(proposal)
	if err != nil {
		return err
	}

	err = ps.state.Set(prefixed, data)

	return err
}

func (ps *ProposalStore) Get(proposalID ProposalID) (*Proposal, error) {
	proposal := &Proposal{}
	prefixed := append(ps.prefix, proposalID...)
	data, err := ps.state.Get(prefixed)
	if err != nil {
		return nil, err
	}
	err = ps.szlr.Deserialize(data, proposal)
	if err != nil {
		return nil, err
	}

	return proposal, nil
}

func (ps *ProposalStore) Exists(key ProposalID) bool {
	prefixed := append(ps.prefix, key...)
	return ps.state.Exists(prefixed)
}

func (ps *ProposalStore) Delete(key ProposalID) (bool, error) {
	prefixed := append(ps.prefix, key...)
	return ps.state.Delete(prefixed)
}

func (ps *ProposalStore) GetIterable() storage.Iterable {
	return ps.state.GetIterable()
}

func (ps *ProposalStore) Iterate(fn func(id ProposalID, proposal *Proposal) bool) (stopped bool) {
	return ps.state.IterateRange(
		ps.prefix,
		storage.Rangefix(string(ps.prefix)),
		true,
		func(key, value []byte) bool {
			proposalID := ProposalID(key)
			proposal := &Proposal{}

			err := ps.szlr.Deserialize(value, proposal)
			if err != nil {
				return true
			}
			return fn(proposalID, proposal)
		},
	)
}

func (ps *ProposalStore) IterateProposer(fn func(id ProposalID, proposal *Proposal) bool, proposer keys.Address) (stopped bool) {
	return ps.Iterate(func(id ProposalID, proposal *Proposal) bool {
		if proposal.Proposer.Equal(proposer) {
			return fn(id, proposal)
		}
		return false
	})
}

func (ps *ProposalStore) IterateProposalType(fn func(id ProposalID, proposal *Proposal) bool, proposalType ProposalType) (stopped bool) {
	return ps.Iterate(func(id ProposalID, proposal *Proposal) bool {
		if proposal.Type == proposalType {
			return fn(id, proposal)
		}
		return false
	})
}

func (ps *ProposalStore) WithState(state *storage.State) *ProposalStore {
	ps.state = state
	return ps
}

func (ps *ProposalStore) WithPrefix(prefix []byte) *ProposalStore {
	ps.prefix = prefix
	return ps
}

func (ps *ProposalStore) WithPrefixType(prefixType ProposalState) *ProposalStore {
	switch prefixType {
	case ProposalStateActive:
		ps.prefix = ps.prefixActive
	case ProposalStatePassed:
		ps.prefix = ps.prefixPassed
	case ProposalStateFailed:
		ps.prefix = ps.prefixFailed
	}
	return ps
}

func (ps *ProposalStore) QueryAllStores(key ProposalID) (*Proposal, ProposalState, error) {
	proposal, err := ps.WithPrefixType(ProposalStateActive).Get(key)
	if err == nil {
		return proposal, ProposalStateActive, nil
	}
	proposal, err = ps.WithPrefixType(ProposalStatePassed).Get(key)
	if err == nil {
		return proposal, ProposalStatePassed, nil
	}
	proposal, err = ps.WithPrefixType(ProposalStateFailed).Get(key)
	if err == nil {
		return proposal, ProposalStateFailed, nil
	}
	return nil, -1, err
}

func (ps *ProposalStore) SetOptions(pOpt *ProposalOptions) {
	ps.proposalOptions = pOpt
}

func (ps *ProposalStore) GetOptions() *ProposalOptions {
	return ps.proposalOptions
}

func NewProposalStore(prefixActive string, prefixPassed string, prefixFailed string, state *storage.State) *ProposalStore {
	return &ProposalStore{
		state:           state,
		szlr:            serialize.GetSerializer(serialize.PERSISTENT),
		prefix:          []byte(prefixActive),
		prefixActive:    []byte(prefixActive),
		prefixPassed:    []byte(prefixPassed),
		prefixFailed:    []byte(prefixFailed),
		proposalOptions: &ProposalOptions{},
	}
}
