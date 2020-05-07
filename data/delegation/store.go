package delegation

import (
	"fmt"
	"sync"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
)

type DelegationStore struct {
	state  *storage.State
	szlr   serialize.Serializer
	prefix []byte
	mux    sync.Mutex
}

type MatureBlock struct {
	Height int64
	Data   []*MatureData
}

type MatureData struct {
	Address keys.Address
	Amount  balance.Amount
	Height  int64
}

func NewDelegationStore(prefix string, state *storage.State) *DelegationStore {
	return &DelegationStore{
		state:  state,
		prefix: storage.Prefix(prefix),
		szlr:   serialize.GetSerializer(serialize.PERSISTENT),
	}
}

func (st *DelegationStore) WithState(state *storage.State) *DelegationStore {
	st.state = state
	return st
}

func (st *DelegationStore) Get(key []byte) (amt *balance.Amount, err error) {
	prefixKey := append(st.prefix, key...)

	dat, _ := st.state.Get(storage.StoreKey(prefixKey))
	amt = balance.NewAmount(0)
	if len(dat) == 0 {
		return
	}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, amt)
	return
}

func (st *DelegationStore) Set(key []byte, amt *balance.Amount) error {
	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(amt)
	if err != nil {
		return err
	}

	prefixKey := append(st.prefix, key...)
	err = st.state.Set(storage.StoreKey(prefixKey), dat)
	return err
}

// Validator - Delegator data

func (st *DelegationStore) getVDKey(validatorAddress keys.Address, delegatorAddress keys.Address) []byte {
	key := []byte(fmt.Sprintf("_e_%s_%s", validatorAddress, delegatorAddress))
	return key
}

func (st *DelegationStore) GetValidatorDelegationAmount(validatorAddress keys.Address, delegatorAddress keys.Address) (amount *balance.Amount, err error) {
	key := st.getVDKey(validatorAddress, delegatorAddress)
	amount, err = st.Get(key)
	return
}

func (st *DelegationStore) SetValidatorDelegationAmount(validatorAddress keys.Address, delegatorAddress keys.Address, amt balance.Amount) (err error) {
	key := st.getVDKey(validatorAddress, delegatorAddress)
	err = st.Set(key, &amt)
	return
}

// Validator data

func (st *DelegationStore) getVKey(validatorAddress keys.Address) []byte {
	key := []byte(fmt.Sprintf("_t_%s", validatorAddress))
	return key
}

func (st *DelegationStore) GetValidatorAmount(validatorAddress keys.Address) (amount *balance.Amount, err error) {
	key := st.getVKey(validatorAddress)
	amount, err = st.Get(key)
	return
}

func (st *DelegationStore) SetValidatorAmount(validatorAddress keys.Address, amt balance.Amount) (err error) {
	key := st.getVKey(validatorAddress)
	err = st.Set(key, &amt)
	return
}

// Delegator effective data

func (st *DelegationStore) getDEKey(delegatorAddress keys.Address) []byte {
	key := []byte(fmt.Sprintf("_d_e_%s", delegatorAddress))
	return key
}

func (st *DelegationStore) GetDelegatorEffectiveAmount(delegatorAddress keys.Address) (amount *balance.Amount, err error) {
	key := st.getDEKey(delegatorAddress)
	amount, err = st.Get(key)
	return
}

func (st *DelegationStore) SetDelegatorEffectiveAmount(delegatorAddress keys.Address, amt balance.Amount) (err error) {
	key := st.getDEKey(delegatorAddress)
	err = st.Set(key, &amt)
	return
}

// Delegator bounded data

func (st *DelegationStore) getDBKey(delegatorAddress keys.Address) []byte {
	key := []byte(fmt.Sprintf("_d_b_%s", delegatorAddress))
	return key
}

func (st *DelegationStore) GetDelegatorBoundedAmount(delegatorAddress keys.Address) (amount *balance.Amount, err error) {
	key := st.getDBKey(delegatorAddress)
	amount, err = st.Get(key)
	return
}

func (st *DelegationStore) SetDelegatorBoundedAmount(delegatorAddress keys.Address, amt balance.Amount) (err error) {
	key := st.getDBKey(delegatorAddress)
	err = st.Set(key, &amt)
	return
}

// mature

func (st *DelegationStore) getMatureKey(version int64) []byte {
	key := []byte(fmt.Sprintf("_m_%d", version))
	return key
}

func (st *DelegationStore) GetMatureAmounts(version int64) (mature *MatureBlock, err error) {
	key := st.getMatureKey(version)
	prefixKey := append(st.prefix, key...)

	dat, _ := st.state.Get(storage.StoreKey(prefixKey))
	mature = &MatureBlock{
		Height: version,
		Data:   make([]*MatureData, 0),
	}
	if len(dat) == 0 {
		return
	}
	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, mature)
	if err != nil {
		return
	}
	return
}

func (st *DelegationStore) SetMatureAmounts(version int64, mature *MatureBlock) (err error) {
	key := st.getMatureKey(version)

	dat, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(mature)
	if err != nil {
		return err
	}

	prefixKey := append(st.prefix, key...)
	err = st.state.Set(storage.StoreKey(prefixKey), dat)

	return
}

func (st *DelegationStore) GetMaturedPendingAmount(delegatorAddress keys.Address, version int64, count int64) []*MatureData {
	amts := make([]*MatureData, 0)
	for i := int64(0); i < count; i++ {
		height := version + i
		matureBlock, err := st.GetMatureAmounts(height)
		if err != nil {
			continue
		}
		for i, m := range matureBlock.Data {
			if m.Amount.Equals(*balance.NewAmountFromInt(0)) {
				continue
			}
			if m.Address.Equal(delegatorAddress) {
				amts = append(amts, matureBlock.Data[i])
			}
		}
	}
	return amts
}

// Staking

func (st *DelegationStore) Stake(validatorAddress keys.Address, delegatorAddress keys.Address, amount balance.Amount) error {
	st.mux.Lock()
	defer st.mux.Unlock()

	lockedAmt, err := st.GetValidatorAmount(validatorAddress)
	if err != nil {
		return err
	}

	err = st.SetValidatorAmount(validatorAddress, *lockedAmt.Plus(amount))
	if err != nil {
		return err
	}

	lockedAmt, err = st.GetValidatorDelegationAmount(validatorAddress, delegatorAddress)
	if err != nil {
		return err
	}

	err = st.SetValidatorDelegationAmount(validatorAddress, delegatorAddress, *lockedAmt.Plus(amount))
	if err != nil {
		return err
	}

	lockedAmt, err = st.GetDelegatorEffectiveAmount(delegatorAddress)
	if err != nil {
		return err
	}

	err = st.SetDelegatorEffectiveAmount(delegatorAddress, *lockedAmt.Plus(amount))
	if err != nil {
		return err
	}
	return nil
}

func (st *DelegationStore) Unstake(validatorAddress keys.Address, delegatorAddress keys.Address, coin balance.Amount, height int64) error {
	st.mux.Lock()
	defer st.mux.Unlock()
	// st_v_ operation

	// take current total effective amount from total
	totalEffectiveCoin, err := st.GetValidatorAmount(validatorAddress)
	if err != nil {
		return err
	}

	// withdraw from total
	newTotalEffectiveCoin, err := totalEffectiveCoin.Minus(coin)
	if err != nil {
		return err
	}

	// update a new total amount
	err = st.SetValidatorAmount(validatorAddress, *newTotalEffectiveCoin)
	if err != nil {
		return err
	}

	// st_e_ operation

	// take current validator-delegator effective amount
	validatorDelegatedCoin, err := st.GetValidatorDelegationAmount(validatorAddress, delegatorAddress)
	if err != nil {
		return err
	}

	// withdraw from total
	newvalidatorDelegatedCoin, err := validatorDelegatedCoin.Minus(coin)
	if err != nil {
		return err
	}

	// update a new vd effective amount
	err = st.SetValidatorDelegationAmount(validatorAddress, delegatorAddress, *newvalidatorDelegatedCoin)
	if err != nil {
		return err
	}

	// st_d_e_ operation

	// take current delegated effective amount
	delegatedEffectiveCoin, err := st.GetDelegatorEffectiveAmount(delegatorAddress)
	if err != nil {
		return err
	}

	// withdraw from total
	newDelegatedEffectiveCoin, err := delegatedEffectiveCoin.Minus(coin)
	if err != nil {
		return err
	}

	// update a new vd effective amount
	err = st.SetDelegatorEffectiveAmount(delegatorAddress, *newDelegatedEffectiveCoin)
	if err != nil {
		return err
	}

	// st_m_ operation

	// get pending mature coins at block height
	mature, err := st.GetMatureAmounts(height)
	if err != nil {
		return err
	}
	fmt.Printf("Mature got: %+v\n", mature)
	mature.Data = append(mature.Data, &MatureData{
		Address: delegatorAddress,
		Amount:  coin,
		Height:  height,
	})
	fmt.Printf("Mature added: %+v\n", mature)
	// update a new vd effective amount
	err = st.SetMatureAmounts(height, mature)
	if err != nil {
		return err
	}

	return nil
}

func (st *DelegationStore) Withdraw(validatorAddress keys.Address, delegatorAddress keys.Address, coin balance.Amount) error {
	st.mux.Lock()
	defer st.mux.Unlock()

	// taking into apply delegator bound amount
	delegatorBoundCoin, err := st.GetDelegatorBoundedAmount(delegatorAddress)
	if err != nil {
		return err
	}

	// withdraw amount for unstake from bound amount
	resultCoin, err := delegatorBoundCoin.Minus(coin)
	if err != nil {
		return err
	}

	// update current bound delegator amount with unstake amount
	err = st.SetDelegatorBoundedAmount(delegatorAddress, *resultCoin)
	if err != nil {
		return err
	}

	return nil
}

func (st *DelegationStore) UpdateWithdrawReward(height int64) {
	st.mux.Lock()
	defer st.mux.Unlock()

	// st_m_ operation

	// get pending mature coins at block height
	mature, err := st.GetMatureAmounts(height)
	if err != nil {
		return
	}

	for i := range mature.Data {
		m := mature.Data[i]
		if m.Amount.Equals(*balance.NewAmountFromInt(0)) {
			continue
		}

		// taking into apply delegator bound amount
		delegatorBoundCoin, err := st.GetDelegatorBoundedAmount(m.Address)
		if err != nil {
			continue
		}

		// update current bound delegator amount with unlocked unstaked amount
		err = st.SetDelegatorBoundedAmount(m.Address, *delegatorBoundCoin.Plus(m.Amount))
		if err != nil {
			continue
		}
	}

	if len(mature.Data) != 0 {
		// reset info
		mature = &MatureBlock{
			Height: height,
			Data:   make([]*MatureData, 0),
		}
		st.SetMatureAmounts(height, mature)
	}
}
