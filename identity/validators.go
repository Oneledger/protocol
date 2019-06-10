package identity

import (
	"bytes"
	"fmt"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/Oneledger/protocol/storage"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/abci/types"
)

type Validator struct {
	Address keys.Address
	PubKey  keys.PublicKey
	Power   int64
	Name    string
	Staking balance.Coin
}

func (v Validator) Bytes() []byte {
	value, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(v)
	if err != nil {
		logger.Error("validator not serializable", err)
		return nil
	}
	return value
}

func (v *Validator) FromBytes(msg []byte) *Validator {
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(msg, v)
	if err != nil {
		logger.Error("failed to deserialize account from bytes", err)
		return &Validator{}
	}
	return v
}

type Stake struct {
	Address keys.Address
	Pubkey  keys.PublicKey
	Name    string
	Amount  balance.Coin
}

type Unstake struct {
	Address keys.Address
	Amount  balance.Coin
}

type Validators struct {
	cached        storage.Store
	proposer      keys.Address
	newValidators []Validator
	toBeRemoved   []Validator
}

func NewValidators() *Validators {
	cache := storage.NewStorage(storage.CACHE, "validators")
	//todo: get the genesis validators when start the node
	return &Validators{
		cached:        cache,
		proposer:      []byte(nil),
		newValidators: make([]Validator, 0),
		toBeRemoved:   make([]Validator, 0),
	}
}

//setup the validators according to begin block
func (vs *Validators) Set(req types.RequestBeginBlock) error {
	vs.proposer = req.Header.GetProposerAddress()
	err := updateValidiatorSet(vs.cached, req.LastCommitInfo.Votes)
	if err != nil {
		return err
	}
	vs.newValidators = make([]Validator, 0)
	vs.toBeRemoved = make([]Validator, 0)
	vs.toBeRemoved = makingslash(vs, req.ByzantineValidators)
	return err
}

func updateValidiatorSet(cached storage.Store, votes []types.VoteInfo) error {

	fmt.Println("q", votes)
	for _, v := range votes {
		addr := v.Validator.GetAddress()
		if !cached.Exists(addr) {
			return errors.New("validator set not match to last commit")
		}
	}
	return nil
}

func (vs *Validators) HandleStake(apply Stake) *Validators {
	if !checkPubkeyAddress(apply.Pubkey, apply.Address) {
		return vs
	}

	validator := &Validator{}
	if vs.cached.Exists(apply.Address.Bytes()) {

		validator = &Validator{
			Address: apply.Address,
			PubKey:  apply.Pubkey,
			Power:   calculatePower(apply.Amount),
			Name:    apply.Name,
			Staking: apply.Amount,
		}
	}

	value, err := vs.cached.Get(apply.Address.Bytes())
	if err != nil {
		logger.Error("failed to get validator from cache even it exist", err)
	}
	validator = validator.FromBytes(value)

	amt, err := validator.Staking.Plus(apply.Amount)
	if err != nil {
		logger.Error("error adding staking amount", err)
		return vs
	}
	validator.Staking = amt

	vs.newValidators = append(vs.newValidators, *validator)
	return vs
}

func calculatePower(stake balance.Coin) int64 {
	//todo: change to correct power function
	return stake.Amount.Int64()
}

func checkPubkeyAddress(pubkey keys.PublicKey, address keys.Address) bool {
	handler, err := pubkey.GetHandler()
	if err != nil {
		return false
	}
	if bytes.Equal(address, handler.Address()) {
		return true
	}
	return false
}

//todo: implement the proper slashing
func makingslash(vs *Validators, evidences []types.Evidence) []Validator {
	remove := make([]Validator, 0)
	for _, evidence := range evidences {
		if vs.cached.Exists(evidence.Validator.Address) {
			value, err := vs.cached.Get(evidence.Validator.GetAddress())
			if err != nil {
				logger.Error("failed to get validator from cache even it exist", err)
			}
			validator := (&Validator{}).FromBytes(value)
			validator.Power = 0
			remove = append(remove, *validator)
		}
	}
	return remove
}

func (vs *Validators) HandleUnstake(purge Unstake) *Validators {

	if vs.cached.Exists(purge.Address.Bytes()) {
		return vs
	}
	value, err := vs.cached.Get(purge.Address.Bytes())
	if err != nil {
		logger.Error("failed to get validator from cache even it exist", err)
	}
	validator := (&Validator{}).FromBytes(value)
	validator.Staking = purge.Amount
	vs.toBeRemoved = append(vs.toBeRemoved, *validator)
	return vs
}

func (vs *Validators) GetEndBlockUpdate(ctx *ValidatorContext, req types.RequestEndBlock) []types.ValidatorUpdate {

	validatorUpdates := make([]types.ValidatorUpdate, 0)

	if req.Height > 1 && (len(vs.newValidators) > 0 || len(vs.toBeRemoved) > 0) {

		for _, add := range vs.newValidators {
			if !transferVT(*ctx, add) {
				logger.Error("failed to transfer the vt token for validator", add)
			}
			_ = vs.cached.Set(add.Address.Bytes(), add.Bytes())
		}
		for _, remove := range vs.toBeRemoved {
			if !transferVT(*ctx, remove) {
				logger.Error("failed to transfer the vt token for validator", remove)
			}
			_ = vs.cached.Set(remove.Address.Bytes(), remove.Bytes())
		}

	}

	//todo : get the final updates from vs.cached
	return validatorUpdates
}

func (vs *Validators) GetValidator(addr keys.Address) *Validator {
	return nil
}

type ValidatorContext struct {
	Balances *balance.Store
	//todo: add necessary config
}

func NewValidatorContext(balances *balance.Store) *ValidatorContext {
	return &ValidatorContext{
		Balances: balances,
	}
}

func transferVT(ctx ValidatorContext, validator Validator) bool {
	logger.Debug("Processing Transfer of VT to Payment Account")

	//todo: implement transfer vt from account balance to some where else

	return true
}
