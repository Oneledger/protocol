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
	*storage.ChainState
	proposer  keys.Address
	queue     ValidatorQueue
	byzantine []Validator
}

func NewValidatorStore(cfg config.Server, dbPath string, dbType string) *ValidatorStore {
	store := storage.NewChainState("validators", dbPath, dbType, storage.PERSISTENT)
	// TODO: get the genesis validators when start the node
	return &ValidatorStore{
		ChainState: store,
		proposer:   []byte(nil),
		queue:      ValidatorQueue{PriorityQueue: make(utils.PriorityQueue, 0, 100)},
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
			StakeAddress: h.Address(), // TODO : put the validator address for validator in genesis for now. should be a different address from who the validator pay the stake
			PubKey:       vpubkey,
			Power:        v.Power,
			Name:         "",
			// TODO: this should be change with @TODO99
			Staking: currency.NewCoinFromInt(v.Power),
		}
		err = vs.ChainState.Set(h.Address().Bytes(), validator.Bytes())
		if err != nil {
			return req.Validators, errors.New("failed to add initial validators")
		}
		validatorUpdates = append(validatorUpdates, v)
	}
	return validatorUpdates, nil
}

// setup the validators according to begin block
func (vs *ValidatorStore) Set(req types.RequestBeginBlock) error {
	vs.proposer = req.Header.GetProposerAddress()
	err := updateValidatorSet(vs.ChainState, req.LastCommitInfo.Votes)
	if err != nil {
		return errors.Wrapf(err, "height=%d", req.Header.Height)
	}

	// update the byzantine node that need to be slashed
	// this should happened before initialize the queue.
	vs.byzantine = make([]Validator, 0)
	vs.byzantine = makingslash(vs, req.ByzantineValidators)
	for _, remove := range vs.byzantine {

		err := vs.ChainState.Set(remove.Address.Bytes(), remove.Bytes())
		if err != nil {
			logger.Error("failed to set byzantine validator power")
			// TODO: add fatal status for here after we decide what to do with byzantine validator
		}
	}

	// initialize the queue for validators
	vs.queue.PriorityQueue = make(utils.PriorityQueue, 0, 100)
	i := 0
	vs.ChainState.Iterate(func(key, value []byte) bool {
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
	vs.ChainState.Iterate(func(key, value []byte) bool {
		validator, err := (&Validator{}).FromBytes(value)
		if err != nil {
			return false
		}
		validatorSet = append(validatorSet, *validator)
		return false
	})
	return validatorSet, nil
}

func updateValidatorSet(store *storage.ChainState, votes []types.VoteInfo) error {

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
			ECDSAPubKey:  apply.ECDSAPubKey,
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

	err := vs.ChainState.Set(validator.Address.Bytes(), value)
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
		if vs.ChainState.Exists(evidence.Validator.Address) {
			value := vs.ChainState.Get(evidence.Validator.GetAddress(), false)
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

	if !vs.ChainState.Exists(unstake.Address.Bytes()) {

		return errors.New("address not exist in validator set")
	}

	value := vs.ChainState.Get(unstake.Address.Bytes(), false)
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

	err = vs.ChainState.Set(validator.Address.Bytes(), validator.Bytes())
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
			result := vs.ChainState.Get(queued.Value(), true)
			validator, err := (&Validator{}).FromBytes(result)
			if err != nil {
				logger.Error(err, "error deserialize validator")
				continue
			}
			// purge validator who's power is 0
			if validator.Power <= 0 {
				_, ok := vs.ChainState.Remove(validator.Address)
				if !ok {
					logger.Error("Failed to remove invalid validator from chainstate")
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
	return vs.ChainState.Commit()
}
