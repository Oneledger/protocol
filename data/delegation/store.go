package delegation

import (
	"fmt"
	"strings"
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
	err = st.szlr.Deserialize(dat, amt)
	return
}

func (st *DelegationStore) Set(key []byte, amt *balance.Amount) error {
	dat, err := st.szlr.Serialize(amt)
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
	err = st.szlr.Deserialize(dat, mature)
	if err != nil {
		return
	}
	return
}

func (st *DelegationStore) SetMatureAmounts(version int64, mature *MatureBlock) (err error) {
	key := st.getMatureKey(version)
	dat, err := st.szlr.Serialize(mature)
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
			if len(delegatorAddress) == 0 || m.Address.Equal(delegatorAddress) {
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
			fmt.Println("failed to get GetDelegatorBoundedAmount!!!")
			continue
		}

		// update current bound delegator amount with unlocked unstaked amount
		err = st.SetDelegatorBoundedAmount(m.Address, *delegatorBoundCoin.Plus(m.Amount))
		if err != nil {
			fmt.Println("failed to get SetDelegatorBoundedAmount!!!")
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

// iterate
func (st *DelegationStore) iterate(subkey string, fn func(addr keys.Address, amt *balance.Amount) bool) (stopped bool) {
	prefix := append(st.prefix, subkey...)
	return st.state.IterateRange(
		prefix,
		storage.Rangefix(string(prefix)),
		true,
		func(key, value []byte) bool {
			amt := balance.NewAmount(0)
			err := st.szlr.Deserialize(value, amt)
			if err != nil {
				fmt.Printf("failed to deserialize delegation amount")
				return false
			}
			addr := keys.Address{}
			err = addr.UnmarshalText(key[len(prefix):])
			if err != nil {
				fmt.Printf("failed to deserialize delegator address")
				return false
			}
			return fn(addr, amt)
		},
	)
}

// iterate
func (st *DelegationStore) iterateVD(subkey string, fn func(validator keys.Address, delegator keys.Address, amt *balance.Amount) bool) (stopped bool) {
	prefix := append(st.prefix, subkey...)
	return st.state.IterateRange(
		prefix,
		storage.Rangefix(string(prefix)),
		true,
		func(key, value []byte) bool {
			amt := balance.NewAmount(0)
			err := st.szlr.Deserialize(value, amt)
			if err != nil {
				fmt.Printf("failed to deserialize delegation amount")
				return false
			}
			// key in format "%validator_%delegator"
			vd := strings.Split(string(key[len(prefix):]), "_")
			if len(vd) != 2 {
				fmt.Printf("failed to deserialize validator delegator addresses")
				return false
			}

			// parse validator and delegator addresses
			validator := keys.Address{}
			err = validator.UnmarshalText([]byte(vd[0]))
			if err != nil {
				fmt.Printf("failed to deserialize validator address")
				return false
			}
			delegator := keys.Address{}
			err = delegator.UnmarshalText([]byte(vd[1]))
			if err != nil {
				fmt.Printf("failed to deserialize delegator address")
				return false
			}
			return fn(validator, delegator, amt)
		},
	)
}

//-----------------------------Dump/Load chain state
//
type DelegationAmount struct {
	Address keys.Address    `json:"address"`
	Amount  *balance.Amount `json:"amount"`
}

type ValidatorDelegationAmount struct {
	Validator keys.Address    `json:"validator"`
	Delegator keys.Address    `json:"delegator"`
	Amount    *balance.Amount `json:"amount"`
}

type DelegationState struct {
	ValidatorAmounts           []*DelegationAmount          `json:"validatorAmounts"`
	ValidatorDelegationAmounts []*ValidatorDelegationAmount `json:"validatorDelegationAmounts"`
	DelegatorEffectiveAmounts  []*DelegationAmount          `json:"delegatorEffectiveAmounts"`
	DelegatorBoundedAmounts    []*DelegationAmount          `json:"delegatorBoundedAmounts"`
	MatureAmounts              []*MatureData                `json:"matureAmounts"`
}

func NewDelegationState() *DelegationState {
	return &DelegationState{
		ValidatorAmounts:           []*DelegationAmount{},
		ValidatorDelegationAmounts: []*ValidatorDelegationAmount{},
		DelegatorEffectiveAmounts:  []*DelegationAmount{},
		DelegatorBoundedAmounts:    []*DelegationAmount{},
		MatureAmounts:              []*MatureData{},
	}
}

func (st *DelegationStore) DumpState(options *Options) (state *DelegationState, succeed bool) {
	state = NewDelegationState()

	// dump each validator's total amount
	st.iterate("_t_", func(addr keys.Address, amt *balance.Amount) bool {
		dm := &DelegationAmount{
			Address: addr,
			Amount:  amt,
		}
		state.ValidatorAmounts = append(state.ValidatorAmounts, dm)
		return false
	})
	// dump each validator_delegator amount
	st.iterateVD("_e_", func(validator keys.Address, delegator keys.Address, amt *balance.Amount) bool {
		vdm := &ValidatorDelegationAmount{
			Validator: validator,
			Delegator: delegator,
			Amount:    amt,
		}
		state.ValidatorDelegationAmounts = append(state.ValidatorDelegationAmounts, vdm)
		return false
	})
	// dump each delegator effective amount
	st.iterate("_d_e_", func(addr keys.Address, amt *balance.Amount) bool {
		dm := &DelegationAmount{
			Address: addr,
			Amount:  amt,
		}
		state.DelegatorEffectiveAmounts = append(state.DelegatorEffectiveAmounts, dm)
		return false
	})
	// dump each delegator bounded amount
	st.iterate("_d_b_", func(addr keys.Address, amt *balance.Amount) bool {
		dm := &DelegationAmount{
			Address: addr,
			Amount:  amt,
		}
		state.DelegatorBoundedAmounts = append(state.DelegatorBoundedAmounts, dm)
		return false
	})
	// dump pending mature amount
	version := st.state.Version()
	matureAmounts := st.GetMaturedPendingAmount(keys.Address{}, version, options.MaturityTime+1)
	state.MatureAmounts = append(state.MatureAmounts, matureAmounts...)

	succeed = true
	return
}

func (st *DelegationStore) LoadState(state DelegationState) (succeed bool) {
	// load each validator's total amount
	for _, dm := range state.ValidatorAmounts {
		err := st.SetValidatorAmount(dm.Address, *dm.Amount)
		if err != nil {
			return
		}
	}
	// load each validator_delegator amount
	for _, vdm := range state.ValidatorDelegationAmounts {
		err := st.SetValidatorDelegationAmount(vdm.Validator, vdm.Delegator, *vdm.Amount)
		if err != nil {
			return
		}
	}
	// load each delegator effective amount
	for _, dm := range state.DelegatorEffectiveAmounts {
		err := st.SetDelegatorEffectiveAmount(dm.Address, *dm.Amount)
		if err != nil {
			return
		}
	}
	// load each delegator bounded amount
	for _, dm := range state.DelegatorBoundedAmounts {
		err := st.SetDelegatorBoundedAmount(dm.Address, *dm.Amount)
		if err != nil {
			return
		}
	}
	// load pending mature amount
	blocks := make(map[int64]*MatureBlock)
	for _, data := range state.MatureAmounts {
		blk, ok := blocks[data.Height]
		if !ok {
			blocks[data.Height] = &MatureBlock{
				Height: data.Height,
				Data:   []*MatureData{data},
			}
		} else {
			blk.Data = append(blk.Data, data)
		}
	}

	succeed = true
	return
}
