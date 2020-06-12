package governance

import (
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"

	"github.com/pkg/errors"
)

type ProposalStore struct {
	state *storage.State
	szlr  serialize.Serializer

	prefix []byte //Current Store Prefix

	prefixActive         []byte
	prefixPassed         []byte
	prefixFailed         []byte
	prefixFinalized      []byte
	prefixFinalizeFailed []byte

	proposalOptions *ProposalOptionSet
}

func (ps *ProposalStore) Set(proposal *Proposal) error {
	prefixed := append(ps.prefix, proposal.ProposalID...)
	data, err := ps.szlr.Serialize(proposal)
	if err != nil {
		return errors.Wrap(err, errorSerialization)
	}

	err = ps.state.Set(prefixed, data)

	return errors.Wrap(err, errorSettingRecord)
}

func (ps *ProposalStore) Get(proposalID ProposalID) (*Proposal, error) {
	proposal := &Proposal{}
	prefixed := append(ps.prefix, proposalID...)
	data, err := ps.state.Get(prefixed)
	if err != nil {
		return nil, errors.Wrap(err, errorGettingRecord)
	}
	err = ps.szlr.Deserialize(data, proposal)
	if err != nil {
		return nil, errors.Wrap(err, errorDeSerialization)
	}

	return proposal, nil
}

func (ps *ProposalStore) Exists(key ProposalID) bool {
	active := append(ps.prefixActive, key...)
	passed := append(ps.prefixPassed, key...)
	failed := append(ps.prefixFailed, key...)
	finalized := append(ps.prefixFinalized, key...)
	finalizeFailed := append(ps.prefixFinalizeFailed, key...)
	return ps.state.Exists(active) || ps.state.Exists(passed) || ps.state.Exists(failed) || ps.state.Exists(finalized) || ps.state.Exists(finalizeFailed)
}

func (ps *ProposalStore) Delete(key ProposalID) (bool, error) {
	prefixed := append(ps.prefix, key...)
	res, err := ps.state.Delete(prefixed)
	if err != nil {
		return false, errors.Wrap(err, errorDeletingRecord)
	}
	return res, err
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

func (ps *ProposalStore) GetState() *storage.State {
	return ps.state
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
	case ProposalStateFinalized:
		ps.prefix = ps.prefixFinalized
	case ProposalStateFinalizeFailed:
		ps.prefix = ps.prefixFinalizeFailed

	}
	return ps
}

func (ps *ProposalStore) QueryAllStores(key ProposalID) (*Proposal, ProposalState, error) {
	prefix := ps.prefix
	defer func() { ps.prefix = prefix }()

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
	proposal, err = ps.WithPrefixType(ProposalStateFinalized).Get(key)
	if err == nil {
		return proposal, ProposalStateFinalized, nil
	}
	proposal, err = ps.WithPrefixType(ProposalStateFinalizeFailed).Get(key)
	if err == nil {
		return proposal, ProposalStateFinalizeFailed, nil
	}
	return nil, ProposalStateInvalid, errors.Wrap(err, errorGettingRecord)
}

// Filter proposals by Proposer and Type in a specified store
func (ps *ProposalStore) FilterProposals(state ProposalState, proposer keys.Address, pType ProposalType) []Proposal {
	prefix := ps.prefix
	defer func() { ps.prefix = prefix }()

	proposals := make([]Proposal, 0)
	ps.WithPrefixType(state).Iterate(func(id ProposalID, proposal *Proposal) bool {
		if pType != ProposalTypeInvalid && proposal.Type != pType {
			return false
		}
		if len(proposer) != 0 && !proposal.Proposer.Equal(proposer) {
			return false
		}
		proposals = append(proposals, *proposal)
		return false
	})
	return proposals
}

func (ps *ProposalStore) SetOptions(pOpt *ProposalOptionSet) {
	ps.proposalOptions = pOpt
}

func (ps *ProposalStore) GetOptions() *ProposalOptionSet {
	return ps.proposalOptions
}

func (ps *ProposalStore) GetOptionsByType(typ ProposalType) *ProposalOption {
	switch typ {
	case ProposalTypeConfigUpdate:
		return &ps.proposalOptions.ConfigUpdate
	case ProposalTypeCodeChange:
		return &ps.proposalOptions.CodeChange
	case ProposalTypeGeneral:
		return &ps.proposalOptions.General
	}
	return nil
}

func NewProposalStore(prefixActive string, prefixPassed string, prefixFailed string, prefixFinalized string, prefixFinalizeFailed string, state *storage.State) *ProposalStore {
	return &ProposalStore{
		state:                state,
		szlr:                 serialize.GetSerializer(serialize.PERSISTENT),
		prefix:               []byte(prefixActive),
		prefixActive:         []byte(prefixActive),
		prefixPassed:         []byte(prefixPassed),
		prefixFailed:         []byte(prefixFailed),
		prefixFinalized:      []byte(prefixFinalized),
		prefixFinalizeFailed: []byte(prefixFinalizeFailed),
		proposalOptions:      &ProposalOptionSet{},
	}
}
