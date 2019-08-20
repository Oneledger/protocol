package identity

import (
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/Oneledger/protocol/utils"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/abci/types"
)

type ValidatorStore struct {
	prefix    []byte
	store     *storage.State
	proposer  keys.Address
	queue     ValidatorQueue
	byzantine []Validator
}

func NewValidatorStore(prefix string, cfg config.Server, state *storage.State) *ValidatorStore {
	// TODO: get the genesis validators when start the node
	return &ValidatorStore{
		prefix:    []byte(prefix + storage.DB_PREFIX),
		store:     state,
		proposer:  []byte(nil),
		queue:     ValidatorQueue{PriorityQueue: make(utils.PriorityQueue, 0, 100)},
		byzantine: make([]Validator, 0),
	}
}

func (vs *ValidatorStore) WithGas(gc storage.GasCalculator) *ValidatorStore {
	vs.store = vs.store.WithGas(gc)
	return vs
}

func (vs *ValidatorStore) Init(req types.RequestInitChain, currencies *balance.CurrencyList) (types.ValidatorUpdates, error) {
	currency, ok := currencies.GetCurrencyByName("VT")
	if !ok {
		return req.Validators, errors.New("stake token not registered")
	}
	validatorUpdates := make([]types.ValidatorUpdate, 0)

	for _, v := range req.Validators {

		vpubkey, err := keys.GetPublicKeyFromBytes(v.PubKey.Data, keys.GetAlgorithmFromTmKeyName(v.PubKey.Type))
		if err != nil {
			return validatorUpdates, errors.Wrap(err, "invalid pubkey type")
		}

		h, err := vpubkey.GetHandler()
		if err != nil {
			return validatorUpdates, errors.Wrap(err, "invalid pubkey type")
		}

		validator := Validator{
			Address:      h.Address(),
			StakeAddress: h.Address(), // TODO : put the validator address for validator in genesis for now. should be a different address from who the validator pay the stake
			PubKey:       vpubkey,
			Power:        v.Power,
			Name:         "",
			// TODO: this should be change with @TODO99
			Staking: currency.NewCoinFromInt(v.Power),
		}
		key := append(vs.prefix, h.Address().Bytes()...)
		err = vs.store.Set(key, validator.Bytes())
		if err != nil {
			return req.Validators, errors.New("failed to add initial validators")
		}
		validatorUpdates = append(validatorUpdates, v)
	}
	return validatorUpdates, nil
}

// setup the validators according to begin block
func (vs *ValidatorStore) Setup(req types.RequestBeginBlock) error {
	vs.proposer = req.Header.GetProposerAddress()
	err := updateValidatorSet(vs.prefix, vs.store, req.LastCommitInfo.Votes)
	if err != nil {
		return errors.Wrapf(err, "height=%d", req.Header.Height)
	}

	// update the byzantine node that need to be slashed
	// this should happened before initialize the queue.
	vs.byzantine = make([]Validator, 0)
	vs.byzantine = makingslash(vs, req.ByzantineValidators)
	for _, remove := range vs.byzantine {
		key := append(vs.prefix, remove.Address.Bytes()...)
		err := vs.store.Set(key, remove.Bytes())
		if err != nil {
			logger.Error("failed to set byzantine validator power")
			// TODO: add fatal status for here after we decide what to do with byzantine validator
		}
	}

	// initialize the queue for validators
	vs.queue.PriorityQueue = make(utils.PriorityQueue, 0, 100)
	i := 0
	vs.Iterate(func(key, value []byte) bool {
		validator, err := (&Validator{}).FromBytes(value)
		if err != nil {
			return false
		}
		queued := utils.NewQueued(key, validator.Power, i)
		vs.queue.append(queued)
		// 		vs.queue.Push(queued)
		i++
		return false
	})
	vs.queue.Init()

	return err
}

// get validators set
func (vs *ValidatorStore) GetValidatorSet() ([]Validator, error) {

	validatorSet := make([]Validator, 0)
	vs.Iterate(func(key, value []byte) bool {
		validator, err := (&Validator{}).FromBytes(value)
		if err != nil {
			return false
		}
		validatorSet = append(validatorSet, *validator)
		return false
	})
	return validatorSet, nil
}

func updateValidatorSet(prefix []byte, store *storage.State, votes []types.VoteInfo) error {

	for _, v := range votes {
		addr := v.Validator.GetAddress()
		key := append(prefix, addr...)
		if !store.Exists(key) {
			return errors.New("validator set not match to last commit")
		}
	}
	return nil
}

// handle stake action
func (vs *ValidatorStore) HandleStake(apply Stake) error {
	validator := &Validator{}
	key := append(vs.prefix, apply.ValidatorAddress.Bytes()...)
	if !vs.store.Exists(key) {

		validator = &Validator{
			Address:      apply.ValidatorAddress,
			StakeAddress: apply.StakeAddress,
			PubKey:       apply.Pubkey,
			Power:        calculatePower(apply.Amount),
			Name:         apply.Name,
			Staking:      apply.Amount,
		}
		// push the new validator to queue
	} else {
		value, _ := vs.store.Get(key)
		if value == nil {
			return errors.New("failed to get validator from store")
		}
		validator, err := validator.FromBytes(value)
		if err != nil {
			return errors.Wrap(err, "error deserialize validator")
		}
		amt, err := validator.Staking.Plus(apply.Amount)
		if err != nil {
			return errors.Wrap(err, "error adding staking amount")
		}
		validator.Staking = amt
		validator.Power = calculatePower(amt)

	}

	value := (validator).Bytes()
	vkey := append(vs.prefix, validator.Address.Bytes()...)
	err := vs.store.Set(vkey, value)
	if err != nil {
		return errors.Wrap(err, "failed to set validator for stake")
	}

	return nil
}

func calculatePower(stake balance.Coin) int64 {
	// TODO: change to correct power function @TODO99
	return stake.Amount.Int.Int64()
}

// TODO: implement the proper slashing
func makingslash(vs *ValidatorStore, evidences []types.Evidence) []Validator {
	remove := make([]Validator, 0)
	for _, evidence := range evidences {
		if vs.store.Exists(evidence.Validator.Address) {
			value, _ := vs.store.Get(evidence.Validator.GetAddress())
			if value == nil {
				logger.Error("failed to get validator from store", evidence.Validator.Address)
			}
			validator, err := (&Validator{}).FromBytes(value)
			if err != nil {
				logger.Error("error deserialize validator")
			}
			validator.Power = 0
			remove = append(remove, *validator)
		}
	}
	return remove
}

func (vs *ValidatorStore) HandleUnstake(unstake Unstake) error {
	validator := &Validator{}
	unstakeKey := append(vs.prefix, unstake.Address.Bytes()...)
	if !vs.store.Exists(unstakeKey) {

		return errors.New("address not exist in validator set")
	}

	value, _ := vs.store.Get(unstakeKey)
	if value == nil {
		return errors.New("failed to get validator from store")
	}
	validator, err := validator.FromBytes(value)
	if err != nil {
		return errors.Wrap(err, "error deserialize validator")
	}

	amt, err := validator.Staking.Minus(unstake.Amount)
	if err != nil {
		return errors.Wrap(err, "minus staking amount")
	}
	validator.Staking = amt
	validator.Power = calculatePower(amt)
	vKey := append(vs.prefix, validator.Address.Bytes()...)

	err = vs.store.Set(vKey, validator.Bytes())
	if err != nil {
		return errors.Wrap(err, "failed to set validator for unstake")
	}

	return nil
}

func (vs *ValidatorStore) GetEndBlockUpdate(ctx *ValidatorContext, req types.RequestEndBlock) []types.ValidatorUpdate {

	validatorUpdates := make([]types.ValidatorUpdate, 0)

	if req.Height > 1 || (len(vs.byzantine) > 0) {
		cnt := 0
		for vs.queue.Len() > 0 && cnt < 64 {
			queued := vs.queue.Pop()
			cqKey := append(vs.prefix, queued.Value()...)

			result, _ := vs.store.Get(cqKey)
			validator, err := (&Validator{}).FromBytes(result)
			if err != nil {
				logger.Error(err, "error deserialize validator")
				continue
			}
			// purge validator who's power is 0
			if validator.Power <= 0 {
				vKey := append(vs.prefix, validator.Address.Bytes()...)

				ok, err := vs.store.Delete(vKey)
				if !ok {
					logger.Error(err.Error())
				}
			}
			validatorUpdates = append(validatorUpdates, types.ValidatorUpdate{
				PubKey: validator.PubKey.GetABCIPubKey(),
				Power:  validator.Power,
			})
			cnt++
		}
	}

	// TODO : get the final updates from vs.cached
	return validatorUpdates
}

func (vs *ValidatorStore) Commit() ([]byte, int64) {
	return vs.Commit()
}
