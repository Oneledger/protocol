package identity

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/abci/types"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/Oneledger/protocol/utils"
)

type ValidatorStore struct {
	prefix      []byte
	prefixPurge []byte
	store       *storage.State
	proposer    keys.Address
	queue       ValidatorQueue
	byzantine   []Validator
	lastActive  map[string]int64
	totalPower  int64
	isValidator bool
}

func NewValidatorStore(prefix string, prefixPurge string, state *storage.State) *ValidatorStore {
	// TODO: get the genesis validators when start the node
	return &ValidatorStore{
		prefix:      storage.Prefix(prefix),
		prefixPurge: storage.Prefix(prefixPurge),
		store:       state,
		proposer:    []byte(nil),
		queue:       ValidatorQueue{PriorityQueue: make(utils.PriorityQueue, 0, 100)},
		byzantine:   make([]Validator, 0),
		lastActive:  make(map[string]int64),
		totalPower:  0,
	}
}

func (vs *ValidatorStore) WithState(state *storage.State) *ValidatorStore {
	vs.store = state
	return vs
}

func (vs *ValidatorStore) Get(addr keys.Address) (*Validator, error) {
	key := append(vs.prefix, addr...)
	value, _ := vs.store.Get(key)
	if value == nil {
		return nil, errors.New("failed to get validator from store")
	}
	validator := &Validator{}
	validator, err := validator.FromBytes(value)
	if err != nil {
		return nil, errors.Wrap(err, "error deserialize validator")
	}
	return validator, nil
}

func (vs *ValidatorStore) Exists(addr keys.Address) bool {
	key := append(vs.prefix, addr...)
	return vs.store.Exists(key)
}

func (vs *ValidatorStore) Set(validator Validator) error {
	return vs.set(validator)
}

func (vs *ValidatorStore) set(validator Validator) error {
	value := (validator).Bytes()
	vkey := append(vs.prefix, validator.Address.Bytes()...)
	err := vs.store.Set(vkey, value)
	if err != nil {
		return errors.Wrap(err, "failed to set validator for stake")
	}
	return nil
}

func (vs *ValidatorStore) Iterate(fn func(addr keys.Address, validator *Validator) bool) (stopped bool) {
	return vs.store.IterateRange(
		vs.prefix,
		storage.Rangefix(string(vs.prefix)),
		true,
		func(key, value []byte) bool {
			validator, err := (&Validator{}).FromBytes(value)
			if err != nil {
				logger.Error("failed to deserialize validator")
				return false
			}
			addr := key[len(vs.prefix):]
			return fn(addr, validator)
		},
	)
}

func (vs *ValidatorStore) Init(req types.RequestInitChain, currencies *balance.CurrencySet) (types.ValidatorUpdates, error) {
	_, ok := currencies.GetCurrencyByName("OLT")
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

		validator, err := vs.Get(h.Address().Bytes())
		if err != nil {
			return validatorUpdates, errors.Wrap(err, "failed to get the validator")
		}
		//todo: add more check for initial validators if needed.
		_ = validator
		validatorUpdates = append(validatorUpdates, v)
	}
	return validatorUpdates, nil
}

// setup the validators according to begin block
func (vs *ValidatorStore) Setup(req types.RequestBeginBlock, nodeValidatorAddress keys.Address) error {
	vs.proposer = req.Header.GetProposerAddress()
	var def error
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
			def = err
		}
		logger.Infof("slashed byzantine validator, address = %s", remove.Address)
	}

	vs.InitValidatorQueue(nodeValidatorAddress, req.GetHeader().Height)
	vs.cacheActiveValidators(req.LastCommitInfo)

	return def
}

func (vs *ValidatorStore) InitValidatorQueue(nodeValidatorAddress keys.Address, height int64) {
	// initialize the queue for validators
	vs.queue.PriorityQueue = make(utils.PriorityQueue, 0, 100)
	i := 0
	vs.totalPower = 0
	vs.Iterate(func(addr keys.Address, validator *Validator) bool {
		key := append(vs.prefix, addr...)
		data := vs.store.GetVersioned(height-1, key)
		if len(data) == 0 {
			logger.Errorf("Previous state data not found for address: %s", keys.Address(addr).Humanize())
			return false
		}
		validator = &Validator{}
		if err := serialize.GetSerializer(serialize.JSON).Deserialize(data, validator); err != nil {
			logger.Errorf("Validator: %s not found", keys.Address(addr).Humanize())
			return false
		}

		queued := utils.NewQueued(addr, validator.Power, i)
		vs.queue.Push(queued)
		vs.totalPower += validator.Power
		if bytes.Equal(addr, nodeValidatorAddress) {
			vs.isValidator = true
		}
		i++
		return false
	})
	vs.queue.Init()
}

// cache last active validators in tendermint
func (vs *ValidatorStore) cacheActiveValidators(lastCommit types.LastCommitInfo) {
	vs.lastActive = make(map[string]int64)
	for _, vote := range lastCommit.Votes {
		addr := keys.Address(vote.Validator.Address)
		vs.lastActive[string(addr)] = vote.Validator.Power
	}
}

// get validators set
func (vs *ValidatorStore) GetValidatorSet() ([]Validator, error) {

	validatorSet := make([]Validator, 0)
	vs.Iterate(func(addr keys.Address, validator *Validator) bool {
		validatorSet = append(validatorSet, *validator)
		return false
	})
	return validatorSet, nil
}

func (vs *ValidatorStore) GetValidatorMap(stakingOptions *delegation.Options) map[string]bool {
	vMap := make(map[string]bool)
	cnt := int64(0)

	i := 0
	queue := &ValidatorQueue{PriorityQueue: make(utils.PriorityQueue, 0, 100)}
	vs.Iterate(func(addr keys.Address, validator *Validator) bool {
		queued := utils.NewQueued(addr, validator.Power, i)
		queue.Push(queued)
		i++
		return false
	})
	queue.Init()

	for queue.Len() > 0 && cnt < stakingOptions.TopValidatorCount {
		queued := queue.Pop()
		addr := queued.Value()
		validator, err := vs.Get(addr)
		if err != nil {
			continue
		}
		if validator.Power < stakingOptions.MinSelfDelegationAmount.BigInt().Int64() {
			break
		}
		vMap[validator.Address.String()] = true
		cnt++
	}
	return vMap
}

func (vs *ValidatorStore) GetActiveValidatorList(stakingOptions *delegation.Options) ([]Validator, error) {
	activeList := []Validator{}
	validators, err := vs.GetValidatorSet()
	if err != nil {
		return activeList, err
	}
	vMap := vs.GetValidatorMap(stakingOptions)
	for _, v := range validators {
		isActive := vMap[v.Address.String()]
		if isActive {
			activeList = append(activeList, v)
		}
	}
	return activeList, nil
}

// get validators set
func (vs *ValidatorStore) GetValidatorsAddress() ([]keys.Address, error) {

	validatorAddress := make([]keys.Address, 0)
	vs.Iterate(func(addr keys.Address, validator *Validator) bool {
		validatorAddress = append(validatorAddress, addr)
		return false
	})
	return validatorAddress, nil
}

func (vs *ValidatorStore) IsValidator() bool {
	return vs.isValidator
}

func (vs *ValidatorStore) IsValidatorAddress(addr keys.Address) bool {
	v, err := vs.Get(addr)
	if err != nil {
		return false
	}
	if v.Power > 0 {
		return true
	}
	return false
}

// handle stake action
func (vs *ValidatorStore) HandleStake(apply Stake, updateStakeAddress bool) error {
	validator := &Validator{}
	if !vs.Exists(apply.ValidatorAddress) {
		validator = NewValidator(
			apply.ValidatorAddress,
			apply.StakeAddress,
			apply.Pubkey,
			apply.ECDSAPubKey,
			apply.Amount,
			apply.Name,
		)
		// push the new validator to queue
	} else {
		v, err := vs.Get(apply.ValidatorAddress)
		if err != nil {
			return errors.Wrap(err, "error deserialize validator")
		}
		validator = v
		amt := big.NewInt(0).Add(validator.Staking.BigInt(), apply.Amount.BigInt())

		validator.Staking = *balance.NewAmountFromBigInt(amt)
		validator.Power = calculatePower(validator.Staking)
		if updateStakeAddress {
			logger.Infof("update stake address: old= %s, new= %s", validator.StakeAddress.Humanize(), apply.StakeAddress.Humanize())
			validator.StakeAddress = apply.StakeAddress
		}
	}

	err := vs.set(*validator)
	if err != nil {
		return errors.Wrap(err, "failed to set validator for stake")
	}

	return nil
}

func calculatePower(stake balance.Amount) int64 {
	// TODO: change to correct power function @TODO99
	return stake.BigInt().Int64()
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

	validator, err := vs.Get(unstake.Address)
	if err != nil {
		return errors.Wrap(err, "error deserialize validator")
	}

	amt := big.NewInt(0).Sub(validator.Staking.BigInt(), unstake.Amount.BigInt())

	validator.Staking = *balance.NewAmountFromBigInt(amt)
	validator.Power = calculatePower(validator.Staking)
	err = vs.set(*validator)
	if err != nil {
		return errors.Wrap(err, "failed to set validator for unstake")
	}

	return nil
}

// Set validator last purge height
func (vs *ValidatorStore) SetLastPurgeHeight(validator keys.Address, height int64) error {
	key := append(vs.prefixPurge, validator...)
	value, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(height)
	if err != nil {
		return errors.Wrap(err, "failed to serialize purged height")
	}

	err = vs.store.Set(key, value)
	if err != nil {
		return errors.Wrap(err, "failed to set purged height")
	}
	return nil
}

// Get validator last purge height
func (vs *ValidatorStore) GetLastPurgeHeight(validator keys.Address) (height int64, err error) {
	key := append(vs.prefixPurge, validator...)
	dat, _ := vs.store.Get(key)
	height = 0
	if len(dat) == 0 {
		return
	}

	err = serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, &height)
	if err != nil {
		logger.Error("failed to deserialize purged height", err)
		return
	}
	return
}

func (vs *ValidatorStore) GetEndBlockUpdate(ctx *ValidatorContext, req types.RequestEndBlock) []types.ValidatorUpdate {
	height := req.GetHeight()
	logger.Detailf("GetEndBlockUpdate started at block: %d\n", height)
	validatorUpdates := make([]types.ValidatorUpdate, 0)
	distribute := false
	total, err := ctx.FeePool.Get([]byte(fees.POOL_KEY))
	if err != nil {
		logger.Fatal("failed to get the total fee pool")
	} else if ctx.FeePool.GetOpt().MinFee().LessThanCoin(total) {
		distribute = true
	}
	stakingOptions, err := ctx.Govern.GetStakingOptions()
	if err != nil {
		logger.Fatal("failed to get the staking options")
	}

	minSelfDelegationAmount := stakingOptions.MinSelfDelegationAmount.BigInt().Int64()

	if height > 1 || (len(vs.byzantine) > 0) {
		// map for non top-power validators
		nonTopValidators := make(map[string]types.PubKey)

		// collect top-power validators
		cnt := int64(0)
		for vs.queue.Len() > 0 {
			// pop element to test
			queued := vs.queue.Pop()
			addr := queued.Value()
			addrHuman := keys.Address(addr).Humanize()
			key := append(vs.prefix, addr...)
			data := vs.store.GetVersioned(height-1, key)
			if len(data) == 0 {
				logger.Errorf("Previous state data not found for address: %s", addrHuman)
				continue
			}
			validator := &Validator{}
			if err := serialize.GetSerializer(serialize.JSON).Deserialize(data, validator); err != nil {
				logger.Errorf("Validator: %s not found", addrHuman)
				continue
			}

			updateTendermint := false

			if validator.Power >= minSelfDelegationAmount && cnt < stakingOptions.TopValidatorCount {
				updateTendermint = true
				cnt++
			}

			// delete validator who's power is 0
			if validator.Power <= 0 {
				vKey := append(vs.prefix, validator.Address.Bytes()...)
				//TODO: validator delete will not properly delete the item because of state implementation
				ok, err := vs.store.Delete(vKey)
				if !ok {
					logger.Error(err.Error())
				}
			}

			// append to update list
			if updateTendermint {
				logger.Detailf("Validator for update ready: %s - with power: %d\n", addrHuman, validator.Power)
				validatorUpdates = append(validatorUpdates, types.ValidatorUpdate{
					PubKey: validator.PubKey.GetABCIPubKey(),
					Power:  validator.Power,
				})
			} else {
				nonTopValidators[addrHuman] = validator.PubKey.GetABCIPubKey()
			}

			//distribute the fee for validators
			if distribute {
				feeShare := total.MultiplyInt64(queued.Priority()).DivideInt64(vs.totalPower)
				err = ctx.FeePool.MinusFromPool(feeShare)
				if err != nil {
					logger.Fatal("failed to minus from fee pool")
				}
				err = ctx.FeePool.AddToAddress(validator.StakeAddress, feeShare)
				if err != nil {
					logger.Fatal("failed to distribute fee")
				}
			}
		}

		// purge all other active validators if not among the top
		for addr, power := range vs.lastActive {
			addrHuman := keys.Address(addr).Humanize()
			pub, ok := nonTopValidators[addrHuman]
			if !ok {
				continue
			}
			// get last purge height
			purgeHeight, err := vs.GetLastPurgeHeight(keys.Address(addr))
			if err != nil {
				logger.Errorf("failed to get last purge height, validator: %s", addrHuman)
				continue
			}
			// validator purged in block H persists in block H+1, H+2, so we can't purge it again
			if purgeHeight > 0 && height <= purgeHeight+2 {
				continue
			}

			// add to validator update list
			logger.Infof("Validator for purge ready: %s - with power: %d\n", addrHuman, power)
			validatorUpdates = append(validatorUpdates, types.ValidatorUpdate{
				PubKey: pub,
				Power:  0,
			})

			// update last purge height
			err = vs.SetLastPurgeHeight(keys.Address(addr), height)
			if err != nil {
				logger.Fatalf("failed  set last purge height, validator: %s", addrHuman)
			}
		}

		// delegation
		ctx.Delegators.UpdateWithdrawReward(height)
	}
	logger.Detailf("GetEndBlockUpdate end at block: %d\n", height)

	// TODO : get the final updates from vs.cached
	return validatorUpdates
}

func (vs *ValidatorStore) GetBitcoinKeys(net *chaincfg.Params) (list []*btcutil.AddressPubKey, err error) {

	list = make([]*btcutil.AddressPubKey, 0)
	vs.Iterate(func(key keys.Address, validator *Validator) bool {

		var pubKey *btcutil.AddressPubKey
		h, err := validator.ECDSAPubKey.GetHandler()
		if err != nil {
			fmt.Println("GetBitcoinKeys", err)
			return true
		}

		pubKey, err = btcutil.NewAddressPubKey(h.Bytes(), net)
		if err != nil {
			return true
		}

		list = append(list, pubKey)
		return false
	})

	return
}
