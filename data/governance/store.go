package governance

import (
	"encoding/binary"

	"github.com/Oneledger/protocol/data/rewards"

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
	ADMIN_REWARD_OPTION   string = "reward"
	LAST_UPDATE_HEIGHT    string = "lastupdateheight"
)

type Store struct {
	state  *storage.State
	prefix []byte
	height int64
}

func NewStore(prefix string, state *storage.State) *Store {
	return &Store{
		state:  state,
		prefix: storage.Prefix(prefix),
		height: 0,
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

func (st *Store) Get(key string) ([]byte, error) {
	// Get the last update height for the present height
	luh, err := st.GetLUH()
	if err != nil {
		panic(errors.Wrap(err, "Unable to get Last Update Height"))
	}
	// Get the Options from the last update Height
	versionedKey := storage.StoreKey(string(luh) + storage.DB_PREFIX + key)
	prefixKey := append(st.prefix, versionedKey...)
	return st.state.Get(prefixKey)
}

func (st *Store) Set(key string, value []byte) error {
	versionedKey := storage.StoreKey(string(st.height) + storage.DB_PREFIX + key)
	prefixKey := append(st.prefix, versionedKey...)
	err := st.state.Set(prefixKey, value)
	return err
}

func (st *Store) GetUnversioned(key string) ([]byte, error) {
	prefixKey := append(st.prefix, key...)
	return st.state.Get(prefixKey)
}

func (st *Store) SetUnversioned(key string, value []byte) error {
	prefixKey := append(st.prefix, key...)
	err := st.state.Set(prefixKey, value)
	return err
}

func (st *Store) Exists(key []byte) bool {
	prefixKey := append(st.prefix, storage.StoreKey(key)...)
	return st.state.Exists(prefixKey)
}

func (st *Store) SetLUH() error {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(st.height))
	err := st.SetUnversioned(LAST_UPDATE_HEIGHT, b)
	if err != nil {
		return err
	}
	return nil
}

func (st *Store) GetLUH() (int64, error) {
	data, err := st.GetUnversioned(LAST_UPDATE_HEIGHT)
	if err != nil {
		return 0, err
	}
	height := int64(binary.LittleEndian.Uint64(data))

	return height, nil
}

func (st *Store) GetCurrencies() (balance.Currencies, error) {
	result, err := st.Get(ADMIN_CURRENCY_KEY)
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
	bytes, err := st.Get(ADMIN_FEE_OPTION_KEY)
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
	_ = st.SetUnversioned(ADMIN_INITIAL_KEY, []byte("initialed"))
	return true
}

func (st *Store) InitialChain() bool {
	data, err := st.GetUnversioned(ADMIN_INITIAL_KEY)
	if err != nil {
		return true
	}
	if data == nil {
		return true
	}
	return false
}

func (st *Store) GetEpoch() (int64, error) {
	result, err := st.Get(ADMIN_EPOCH_BLOCK_INTERVAL)
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

func (st *Store) GetETHChainDriverOption() (*ethchain.ChainDriverOption, error) {

	bytes, err := st.Get(ADMIN_ETH_CHAINDRIVER_OPTION)
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

	bytes, err := st.Get(ADMIN_BTC_CHAINDRIVER_OPTION)
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
	bytes, err := st.Get(ADMIN_ONS_OPTION)
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
	//TODO :Add Validations for proposal options . EX : Sum of proposal fund distribution must be 100
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
	bytes, err := st.Get(ADMIN_PROPOSAL_OPTION)
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
	bytes, err := st.Get(ADMIN_REWARD_OPTION)
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
