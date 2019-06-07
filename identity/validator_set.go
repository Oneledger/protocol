package identity

import (
	"container/heap"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/abci/types"
)

type ValidatorStore struct {
	*storage.ChainState
	proposer  keys.Address
	queue     ValidatorQueue
	byzantine []Validator
}

func NewValidatorStore(cfg config.Server, dbPath string, dbType string) *ValidatorStore {
	store := storage.NewChainState("validators", dbPath, dbType, storage.PERSISTENT)
	//todo: get the genesis validators when start the node
	return &ValidatorStore{
		ChainState: store,
		proposer:   []byte(nil),
		queue:      make(ValidatorQueue, 0, 100),
		byzantine:  make([]Validator, 0),
	}
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
			StakeAddress: h.Address(), //todo : put the validator address for validator in genesis for now. should be a different address from who the validator pay the stake
			PubKey:       vpubkey,
			Power:        v.Power,
			Name:         "",
			//todo: this should be change with @todo99
			Staking: currency.NewCoinFromInt(v.Power),
		}
		err = vs.ChainState.Set(h.Address().Bytes(), validator.Bytes())
		if err != nil {
			return req.Validators, errors.New("failed to add initial validators")
		}
		//fmt.Println("address in init ", hex.EncodeToString(h.Address().Bytes()))
		validatorUpdates = append(validatorUpdates, v)
	}
	return validatorUpdates, nil
}

//setup the validators according to begin block
func (vs *ValidatorStore) Set(req types.RequestBeginBlock) error {
	vs.proposer = req.Header.GetProposerAddress()
	err := updateValidiatorSet(vs.ChainState, req.LastCommitInfo.Votes)
	if err != nil {
		return err
	}

	//initialize the queue for validators
	vs.queue = make(ValidatorQueue, 0, 100)
	i := 0
	vs.ChainState.Iterate(func(key, value []byte) bool {
		validator := (&Validator{}).FromBytes(value)
		queued := &Queued{
			value:    key,
			priority: validator.Power,
			index:    i,
		}
		vs.queue = append(vs.queue, queued)
		//		vs.queue.Push(queued)
		i++
		return false
	})
	heap.Init(&vs.queue)

	logger.Debug("validator_update queue", vs.queue.Len(), "votes", len(req.LastCommitInfo.Votes))
	//update the byzantine node that need to be slashed
	vs.byzantine = make([]Validator, 0)
	vs.byzantine = makingslash(vs, req.ByzantineValidators)
	return err
}

func updateValidiatorSet(store *storage.ChainState, votes []types.VoteInfo) error {

	for _, v := range votes {
		addr := v.Validator.GetAddress()
		if !store.Exists(addr) {
			return errors.New("validator set not match to last commit")
		}
	}
	return nil
}

// handle stake action
func (vs *ValidatorStore) HandleStake(apply Stake) error {
	validator := &Validator{}

	if !vs.ChainState.Exists(apply.ValidatorAddress.Bytes()) {

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
		value := vs.ChainState.Get(apply.ValidatorAddress.Bytes(), false)
		if value == nil {
			return errors.New("failed to get validator from store")
		}
		validator = validator.FromBytes(value)

		amt, err := validator.Staking.Plus(apply.Amount)
		if err != nil {
			return errors.Wrap(err, "error adding staking amount")
		}
		validator.Staking = amt
		validator.Power = calculatePower(amt)

		queued := &Queued{}
		queued.value = validator.Address
		vs.queue.update(queued, queued.value, validator.Power)
	}

	value := (validator).Bytes()
	logger.Debug("validator stake ", validator, value)

	err := vs.ChainState.Set(validator.Address.Bytes(), value)
	if err != nil {
		return errors.Wrap(err, "failed to set validator for stake")
	}

	return nil
}

func calculatePower(stake balance.Coin) int64 {
	//todo: change to correct power function @todo99
	return stake.Amount.Int64()
}

//todo: implement the proper slashing
func makingslash(vs *ValidatorStore, evidences []types.Evidence) []Validator {
	remove := make([]Validator, 0)
	for _, evidence := range evidences {
		if vs.ChainState.Exists(evidence.Validator.Address) {
			value := vs.ChainState.Get(evidence.Validator.GetAddress(), false)
			if value == nil {
				logger.Error("failed to get validator from store", evidence.Validator.Address)
			}
			validator := (&Validator{}).FromBytes(value)
			validator.Power = 0
			remove = append(remove, *validator)
		}
	}
	return remove
}

func (vs *ValidatorStore) HandleUnstake(unstake Unstake) error {
	validator := &Validator{}
	queued := &Queued{}
	if !vs.ChainState.Exists(unstake.Address.Bytes()) {

		return errors.New("address not exist in validator set")
	}

	value := vs.ChainState.Get(unstake.Address.Bytes(), false)
	if value == nil {
		return errors.New("failed to get validator from store")
	}
	validator = validator.FromBytes(value)
	queued.value = validator.Address

	amt, err := validator.Staking.Minus(unstake.Amount)
	if err != nil {
		return errors.Wrap(err, "minus staking amount")
	}
	validator.Staking = amt
	validator.Power = calculatePower(amt)

	err = vs.ChainState.Set(validator.Address.Bytes(), validator.Bytes())
	if err != nil {
		return errors.Wrap(err, "failed to set validator for unstake")
	}

	vs.queue.update(queued, queued.value, validator.Power)
	return nil
}

func (vs *ValidatorStore) GetEndBlockUpdate(ctx *ValidatorContext, req types.RequestEndBlock) []types.ValidatorUpdate {

	validatorUpdates := make([]types.ValidatorUpdate, 0)

	if req.Height > 1 || (len(vs.byzantine) > 0) {

		for _, remove := range vs.byzantine {

			err := vs.ChainState.Set(remove.Address.Bytes(), remove.Bytes())
			if err != nil {
				logger.Error("failed to set byzantine validator at end block")
			}
		}

		length := vs.queue.Len()
		for i := 0; i < 64 && i < length; i++ {
			queued := heap.Pop(&vs.queue).(*Queued)
			result := vs.ChainState.Get(queued.value, true)
			validator := (&Validator{}).FromBytes(result)
			logger.Debugf("validator_update id: %d validator: %#v", i, validator)
			validatorUpdates = append(validatorUpdates, types.ValidatorUpdate{
				PubKey: validator.PubKey.GetABCIPubKey(),
				Power:  validator.Power,
			})
		}
	}

	//todo : get the final updates from vs.cached
	return validatorUpdates
}
