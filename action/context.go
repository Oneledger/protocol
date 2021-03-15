package action

import (
	"github.com/Oneledger/protocol/data"
	"github.com/Oneledger/protocol/data/rewards"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Oneledger/protocol/data/accounts"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/bitcoin"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/ethereum"
	"github.com/Oneledger/protocol/data/evidence"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/jobs"
	netwkDeleg "github.com/Oneledger/protocol/data/network_delegation"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
)

type Context struct {
	Router              Router
	State               *storage.State
	Header              *abci.Header
	Accounts            accounts.Wallet
	Balances            *balance.Store
	Domains             *ons.DomainStore
	Delegators          *delegation.DelegationStore
	NetwkDelegators     *netwkDeleg.MasterStore
	EvidenceStore       *evidence.EvidenceStore
	FeePool             *fees.Store
	Currencies          *balance.CurrencySet
	FeeOpt              *fees.FeeOption
	Validators          *identity.ValidatorStore
	Witnesses           *identity.WitnessStore
	BTCTrackers         *bitcoin.TrackerStore
	ETHTrackers         *ethereum.TrackerStore
	Logger              *log.Logger
	JobStore            *jobs.JobStore
	LockScriptStore     *bitcoin.LockScriptStore
	ProposalMasterStore *governance.ProposalMasterStore
	RewardMasterStore   *rewards.RewardMasterStore
	GovernanceStore     *governance.Store
	ExtStores           data.Router
	GovUpdate           *GovernaceUpdateAndValidate

	// evm
	StateDB *CommitStateDB
}

func NewContext(r Router, header *abci.Header, state *storage.State,
	wallet accounts.Wallet, balances *balance.Store,
	currencies *balance.CurrencySet, feePool *fees.Store,
	validators *identity.ValidatorStore, witnesses *identity.WitnessStore,
	domains *ons.DomainStore, delegators *delegation.DelegationStore, netwkDelegators *netwkDeleg.MasterStore, evidenceStore *evidence.EvidenceStore,
	btcTrackers *bitcoin.TrackerStore, ethTrackers *ethereum.TrackerStore, jobStore *jobs.JobStore,
	lockScriptStore *bitcoin.LockScriptStore, logger *log.Logger, proposalmaster *governance.ProposalMasterStore,
	rewardmaster *rewards.RewardMasterStore, govern *governance.Store, extStores data.Router, govUpdate *GovernaceUpdateAndValidate,
	stateDB *CommitStateDB,
) *Context {
	return &Context{
		Router:              r,
		State:               state,
		Header:              header,
		Accounts:            wallet,
		Balances:            balances,
		Domains:             domains,
		Delegators:          delegators,
		EvidenceStore:       evidenceStore,
		NetwkDelegators:     netwkDelegators,
		FeePool:             feePool,
		Currencies:          currencies,
		Validators:          validators,
		Witnesses:           witnesses,
		BTCTrackers:         btcTrackers,
		ETHTrackers:         ethTrackers,
		Logger:              logger,
		JobStore:            jobStore,
		LockScriptStore:     lockScriptStore,
		ProposalMasterStore: proposalmaster,
		RewardMasterStore:   rewardmaster,
		GovernanceStore:     govern,
		ExtStores:           extStores,
		GovUpdate:           govUpdate,
		StateDB:             stateDB,
	}
}
