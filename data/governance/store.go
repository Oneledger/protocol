package governance

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"

	"github.com/Oneledger/protocol/data/network_delegation"

	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/evidence"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/data/rewards"
	"github.com/Oneledger/protocol/log"

	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/chains/bitcoin"
	ethchain "github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

const (
	ADMIN_INITIAL_KEY    string = "initial"
	ADMIN_CURRENCY_KEY   string = "currency"
	ADMIN_FEE_OPTION_KEY string = "feeopt"

	ADMIN_EPOCH_BLOCK_INTERVAL string = "epoch"

	ADMIN_ETH_CHAINDRIVER_OPTION string = "ethcdopt"

	ADMIN_BTC_CHAINDRIVER_OPTION string = "btccdopt"
	ADMIN_ONS_OPTION             string = "onsopt"

	ADMIN_PROPOSAL_OPTION string = "proposal"

	ADMIN_STAKING_OPTION string = "stakingopt"

	ADMIN_REWARD_OPTION string = "reward"

	ADMIN_EVIDENCE_OPTION string = "evidenceopt"

	ADMIN_NETWK_DELEG_OPTION string = "networkdelegopt"

	TOTAL_FUNDS_PREFIX string = "t"

	INDIVIDUAL_FUNDS_PREFIX string = "i"

	LAST_UPDATE_HEIGHT             string = "defaultOptions"
	LAST_UPDATE_HEIGHT_CURRENCY    string = "currencyOptions"
	LAST_UPDATE_HEIGHT_FEE         string = "feeOptions"
	LAST_UPDATE_HEIGHT_ETH         string = "ethOptions"
	LAST_UPDATE_HEIGHT_BTC         string = "btcOptions"
	LAST_UPDATE_HEIGHT_REWARDS     string = "rewardsOptions"
	LAST_UPDATE_HEIGHT_STAKING     string = "stakingOptions"
	LAST_UPDATE_HEIGHT_NETWK_DELEG string = "delegOptions"
	LAST_UPDATE_HEIGHT_ONS         string = "onsOptions"
	LAST_UPDATE_HEIGHT_PROPOSAL    string = "proposalOptions"
	LAST_UPDATE_HEIGHT_EVIDENCE    string = "evidenceOptions"
	HEIGHT_INDEPENDENT_VALUE       string = "heightindependent"

	// Pool names
	POOL_BOUNTY     = "BountyPool"
	POOL_FEE        = "FeePool"
	POOL_REWARDS    = "RewardsPool"
	POOL_DELEGATION = "DelegationPool"
)

type Store struct {
	state  *storage.State
	prefix []byte
	height int64
	logger *log.Logger
	mux    sync.RWMutex
}

func NewStore(prefix string, state *storage.State) *Store {
	return &Store{
		state:  state,
		prefix: storage.Prefix(prefix),
		height: 0,
		logger: log.NewDefaultLogger(os.Stdout).WithPrefix("governanceStore"),
	}
}

func (st *Store) WithState(state *storage.State) *Store {
	st.state = state
	return st
}

func (st *Store) WithHeight(height int64) *Store {
	st.height = height
	return st
}

func (st *Store) Get(key string, optKey string) ([]byte, error) {
	// Get the last update height for the present height
	// LUH is unversioned and specific for each option type
	luh, err := st.GetLUH(optKey)
	if err != nil {
		panic(errors.Wrap(err, "Unable to get Last Update Height"))
	}
	// Get the Options from the last update Height
	versionedKey := storage.StoreKey(string(rune(luh)) + storage.DB_PREFIX + key)
	prefixKey := append(st.prefix, versionedKey...)
	return st.state.Get(prefixKey)
}

func (st *Store) Set(key string, value []byte) error {
	versionedKey := storage.StoreKey(string(rune(st.height)) + storage.DB_PREFIX + key)
	prefixKey := append(st.prefix, versionedKey...)
	err := st.state.Set(prefixKey, value)
	return err
}

func (st *Store) GetUnversioned(key string, optKey string) ([]byte, error) {
	optionedLuh := storage.StoreKey(optKey + storage.DB_PREFIX + key)
	prefixKey := append(st.prefix, optionedLuh...)
	return st.state.Get(prefixKey)
}

// LUH FORMAT
// KEY :LAST_UPDATE_HEIGHT_LAST_UPDATE_HEIGHT_FEE (Key for each option)
// VALUE : 0 (Height)
func (st *Store) SetUnversioned(key string, optKey string, value []byte) error {
	optionedLuh := storage.StoreKey(optKey + storage.DB_PREFIX + key)
	prefixKey := append(st.prefix, optionedLuh...)
	err := st.state.Set(prefixKey, value)
	return err
}

func (st *Store) Exists(key []byte) bool {
	prefixKey := append(st.prefix, storage.StoreKey(key)...)
	return st.state.Exists(prefixKey)
}

func (st *Store) SetLUH(optKey string) error {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(st.height))
	// Always gives last update Height for present block
	// Replaying tx ,will keep updating this value at every change
	err := st.SetUnversioned(LAST_UPDATE_HEIGHT, optKey, b)
	if err != nil {
		return err
	}
	st.logger.Debugf("Setting new update height : %d | For : %s ", st.height, optKey)
	return nil
}

func (st *Store) SetAllLUH() error {
	err := st.SetLUH(LAST_UPDATE_HEIGHT_PROPOSAL)
	if err != nil {
		return err
	}
	err = st.SetLUH(LAST_UPDATE_HEIGHT_ONS)
	if err != nil {
		return err
	}
	err = st.SetLUH(LAST_UPDATE_HEIGHT_FEE)
	if err != nil {
		return err
	}
	err = st.SetLUH(LAST_UPDATE_HEIGHT_ETH)
	if err != nil {
		return err
	}
	err = st.SetLUH(LAST_UPDATE_HEIGHT_BTC)
	if err != nil {
		return err
	}
	err = st.SetLUH(LAST_UPDATE_HEIGHT_REWARDS)
	if err != nil {
		return err
	}
	err = st.SetLUH(LAST_UPDATE_HEIGHT_STAKING)
	if err != nil {
		return err
	}
	err = st.SetLUH(LAST_UPDATE_HEIGHT_CURRENCY)
	if err != nil {
		return err
	}
	err = st.SetLUH(LAST_UPDATE_HEIGHT_EVIDENCE)
	if err != nil {
		return err
	}
	err = st.SetLUH(LAST_UPDATE_HEIGHT_NETWK_DELEG)
	if err != nil {
		return err
	}
	err = st.SetLUH(LAST_UPDATE_HEIGHT)
	if err != nil {
		return err
	}
	err = st.SetLUH(LAST_UPDATE_HEIGHT_NETWK_DELEG)
	if err != nil {
		return err
	}
	return nil
}

// LUH -> LAST_UPDATE_HEIGHT_LAST_UPDATE_HEIGHT_FEE
// FEEOPTION ->(CurrentHeight) + storage.DB_PREFIX + optionKey + storage.DB_PREFIX + FeeOption)
func (st *Store) GetLUH(optKey string) (int64, error) {
	data, err := st.GetUnversioned(LAST_UPDATE_HEIGHT, optKey)
	if err != nil {
		return 0, err
	}
	height := int64(binary.LittleEndian.Uint64(data))

	return height, nil
}

func (st *Store) GetCurrencies() (balance.Currencies, error) {
	result, err := st.Get(ADMIN_CURRENCY_KEY, LAST_UPDATE_HEIGHT_CURRENCY)
	currencies := make(balance.Currencies, 0, 10)
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(result, &currencies)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the currencies")
	}
	return currencies, nil
}

func (st *Store) SetCurrencies(currencies balance.Currencies) error {

	currenciesBytes, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(currencies)
	if err != nil {
		return errors.Wrap(err, "failed to serialize currencies")
	}
	err = st.Set(ADMIN_CURRENCY_KEY, currenciesBytes)
	if err != nil {
		return errors.Wrap(err, "failed to set the currencies")
	}
	return nil
}

func (st *Store) GetFeeOption() (*fees.FeeOption, error) {
	feeOpt := &fees.FeeOption{}
	bytes, err := st.Get(ADMIN_FEE_OPTION_KEY, LAST_UPDATE_HEIGHT_FEE)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get FeeOption")
	}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(bytes, feeOpt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize FeeOption stored")
	}

	return feeOpt, nil
}

func (st *Store) SetFeeOption(feeOpt fees.FeeOption) error {

	bytes, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(feeOpt)
	if err != nil {
		return errors.Wrap(err, "failed to serialize FeeOption")
	}
	err = st.Set(ADMIN_FEE_OPTION_KEY, bytes)
	if err != nil {
		return errors.Wrap(err, "failed to set the FeeOption")
	}
	return nil
}

func (st *Store) Initiated() bool {
	_ = st.SetUnversioned(ADMIN_INITIAL_KEY, HEIGHT_INDEPENDENT_VALUE, []byte("initialed"))
	return true
}

func (st *Store) InitialChain() bool {
	data, err := st.GetUnversioned(ADMIN_INITIAL_KEY, HEIGHT_INDEPENDENT_VALUE)
	if err != nil {
		return true
	}
	if data == nil {
		return true
	}
	return false
}

func (st *Store) GetEpoch() (int64, error) {
	result, err := st.Get(ADMIN_EPOCH_BLOCK_INTERVAL, LAST_UPDATE_HEIGHT)
	if err != nil {
		return 0, err
	}

	epoch := int64(binary.LittleEndian.Uint64(result))

	return epoch, nil
}

func (st *Store) SetEpoch(epoch int64) error {

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(epoch))

	err := st.Set(ADMIN_EPOCH_BLOCK_INTERVAL, b)
	if err != nil {
		return errors.Wrap(err, "failed to set the currencies")
	}
	return nil
}

func (st *Store) GetStakingOptions() (*delegation.Options, error) {

	bytes, err := st.Get(ADMIN_STAKING_OPTION, LAST_UPDATE_HEIGHT_STAKING)
	if err != nil {
		return nil, err
	}

	r := &delegation.Options{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(bytes, r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize staking options")
	}

	return r, nil
}

func (st *Store) SetEvidenceOptions(opt evidence.Options) error {

	bytes, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(opt)
	if err != nil {
		return errors.Wrap(err, "failed to serialize evidence options")
	}

	err = st.Set(ADMIN_EVIDENCE_OPTION, bytes)
	if err != nil {
		return errors.Wrap(err, "failed to set the evidence options")
	}

	return nil
}

func (st *Store) GetEvidenceOptions() (*evidence.Options, error) {

	bytes, err := st.Get(ADMIN_EVIDENCE_OPTION, LAST_UPDATE_HEIGHT_EVIDENCE)
	if err != nil {
		return nil, err
	}

	r := &evidence.Options{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(bytes, r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize evidence options")
	}

	return r, nil
}

func (st *Store) SetStakingOptions(opt delegation.Options) error {

	bytes, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(opt)
	if err != nil {
		return errors.Wrap(err, "failed to serialize staking options")
	}

	err = st.Set(ADMIN_STAKING_OPTION, bytes)
	if err != nil {
		return errors.Wrap(err, "failed to set the staking options")
	}

	return nil
}

func (st *Store) GetETHChainDriverOption() (*ethchain.ChainDriverOption, error) {

	bytes, err := st.Get(ADMIN_ETH_CHAINDRIVER_OPTION, LAST_UPDATE_HEIGHT_ETH)
	if err != nil {
		return nil, err
	}

	r := &ethchain.ChainDriverOption{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(bytes, r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize eth chaindriver option stored")
	}

	return r, nil
}

func (st *Store) SetETHChainDriverOption(opt ethchain.ChainDriverOption) error {

	bytes, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(opt)
	if err != nil {
		return errors.Wrap(err, "failed to serialize eth chaindriver option")
	}

	err = st.Set(ADMIN_ETH_CHAINDRIVER_OPTION, bytes)
	if err != nil {
		return errors.Wrap(err, "failed to set the eth chaindriver option")
	}

	return nil
}

func (st *Store) GetBTCChainDriverOption() (*bitcoin.ChainDriverOption, error) {

	bytes, err := st.Get(ADMIN_BTC_CHAINDRIVER_OPTION, LAST_UPDATE_HEIGHT_BTC)
	if err != nil {
		return nil, err
	}

	r := &bitcoin.ChainDriverOption{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(bytes, r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize btc chaindriver option stored")
	}

	return r, nil
}

func (st *Store) SetBTCChainDriverOption(opt bitcoin.ChainDriverOption) error {

	bytes, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(opt)
	if err != nil {
		return errors.Wrap(err, "failed to serialize btc chaindriver option")
	}

	err = st.Set(ADMIN_BTC_CHAINDRIVER_OPTION, bytes)
	if err != nil {
		return errors.Wrap(err, "failed to set the btc chaindriver option")
	}

	return nil
}

func (st *Store) SetONSOptions(onsOpt ons.Options) error {
	bytes, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(onsOpt)
	if err != nil {
		return errors.Wrap(err, "failed to serialize ons options")
	}
	err = st.Set(ADMIN_ONS_OPTION, bytes)
	if err != nil {
		return errors.Wrap(err, "failed to set the ons options")
	}
	return nil
}

func (st *Store) GetONSOptions() (*ons.Options, error) {
	bytes, err := st.Get(ADMIN_ONS_OPTION, LAST_UPDATE_HEIGHT_ONS)
	if err != nil {
		return nil, err
	}
	r := &ons.Options{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(bytes, r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize ons options")
	}
	return r, nil
}

func (st *Store) SetProposalOptions(propOpt ProposalOptionSet) error {
	bytes, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(propOpt)
	if err != nil {
		return errors.Wrap(err, "failed to serialize proposal options")
	}
	err = st.Set(ADMIN_PROPOSAL_OPTION, bytes)
	if err != nil {
		return errors.Wrap(err, "failed to set the proposal options")
	}
	return nil
}

func (st *Store) GetProposalOptions() (*ProposalOptionSet, error) {
	bytes, err := st.Get(ADMIN_PROPOSAL_OPTION, LAST_UPDATE_HEIGHT_PROPOSAL)
	if err != nil {
		return nil, err
	}
	propOpt := &ProposalOptionSet{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(bytes, propOpt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize proposal options")
	}

	return propOpt, nil
}

func (st *Store) GetProposalOptionsByType(ptype ProposalType) (*ProposalOption, error) {

	pOpts, err := st.GetProposalOptions()
	if err != nil {
		return nil, err
	}
	switch ptype {
	case ProposalTypeConfigUpdate:
		return &pOpts.ConfigUpdate, nil
	case ProposalTypeCodeChange:
		return &pOpts.CodeChange, nil
	case ProposalTypeGeneral:
		return &pOpts.General, nil
	}
	return nil, errors.New(fmt.Sprintf("Options of Type %s not found", ptype))
}

func (st *Store) SetRewardOptions(rewardOptions rewards.Options) error {
	bytes, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(rewardOptions)
	if err != nil {
		return errors.Wrap(err, "failed to serialize reward options")
	}
	err = st.Set(ADMIN_REWARD_OPTION, bytes)
	if err != nil {
		return errors.Wrap(err, "failed to set the reward options")
	}
	return nil
}

func (st *Store) GetRewardOptions() (*rewards.Options, error) {
	bytes, err := st.Get(ADMIN_REWARD_OPTION, LAST_UPDATE_HEIGHT_REWARDS)
	if err != nil {
		return nil, err
	}

	rewardOptions := &rewards.Options{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(bytes, rewardOptions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize reward options")
	}
	return rewardOptions, nil
}

func (st *Store) SetNetworkDelegOptions(delegOptions network_delegation.Options) error {
	bytes, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(delegOptions)
	if err != nil {
		return errors.Wrap(err, "failed to serialize network delegation options")
	}
	err = st.Set(ADMIN_NETWK_DELEG_OPTION, bytes)
	if err != nil {
		return errors.Wrap(err, "failed to set network delegation options")
	}
	return nil
}

func (st *Store) GetNetworkDelegOptions() (*network_delegation.Options, error) {
	bytes, err := st.Get(ADMIN_NETWK_DELEG_OPTION, LAST_UPDATE_HEIGHT_NETWK_DELEG)
	if err != nil {
		return nil, err
	}

	delegOptions := &network_delegation.Options{}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(bytes, delegOptions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to deserialize network delegation options")
	}
	return delegOptions, nil
}

func (st *Store) GetPoolList() (map[string]keys.Address, error) {
	poolList := map[string]keys.Address{}
	propOpt, err := st.GetProposalOptions()
	if err != nil {
		return nil, err
	}
	rewardOpt, err := st.GetRewardOptions()
	if err != nil {
		return nil, err
	}
	poolList[POOL_BOUNTY] = keys.Address(propOpt.BountyProgramAddr)
	poolList[POOL_FEE] = keys.Address(fees.POOL_KEY)
	poolList[POOL_REWARDS] = keys.Address(rewardOpt.RewardPoolAddress)
	poolList[POOL_DELEGATION] = keys.Address(network_delegation.DELEGATION_POOL_KEY)
	return poolList, nil
}

func (st *Store) GetPoolByName(poolName string) (address keys.Address, err error) {
	poolList, err := st.GetPoolList()
	if err != nil {
		return
	}
	address, ok := poolList[poolName]
	if !ok {
		err = errors.New("Pool not found: " + poolName)
	}
	return
}
