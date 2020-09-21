package identity

import (
	"math/big"
	"os"
	"testing"
	"time"

	db "github.com/tendermint/tm-db"

	"github.com/Oneledger/protocol/storage"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/chain"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/evidence"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
)

func generateKeyPair() (crypto.Address, crypto.PubKey, ed25519.PrivKeyEd25519) {

	prikey := ed25519.GenPrivKey()
	pubkey := prikey.PubKey()
	addr := pubkey.Address()

	return addr, pubkey, prikey
}

var (
	from, fromPubkey, _       = generateKeyPair()
	fromMal, fromPubkeyMal, _ = generateKeyPair()
)

func tearDownCtx() {
	os.Remove("test_dbpath")
}

func setUpCtx() *ValidatorContext {
	db := db.NewDB("test", db.MemDBBackend, "test_dbpath")
	cs := storage.NewState(storage.NewChainState("balance", db))
	balances := balance.NewStore("tb", cs)

	currencies := balance.NewCurrencySet()
	currency := balance.Currency{
		Id:      0,
		Name:    "OLT",
		Chain:   chain.ONELEDGER,
		Decimal: 18,
		Unit:    "nue",
	}
	currencies.Register(currency)

	feeOpt := &fees.FeeOption{
		FeeCurrency:   currency,
		MinFeeDecimal: 9,
	}
	feePool := &fees.Store{}
	feePool.SetupOpt(feeOpt)

	delegators := delegation.NewDelegationStore("tst", cs)
	evidenceStore := evidence.NewEvidenceStore("tes", cs)
	govern := governance.NewStore("tg", cs)
	validators := NewValidatorStore("tv", "purged", cs)

	evidenceOption := evidence.Options{
		MinVotesRequired: 2,
		BlockVotesDiff:   3,

		PenaltyBasePercentage: 30,
		PenaltyBaseDecimals:   100,

		PenaltyBountyPercentage: 50,
		PenaltyBountyDecimals:   100,

		PenaltyBurnPercentage: 50,
		PenaltyBurnDecimals:   100,

		ValidatorReleaseTime:    5,
		ValidatorVotePercentage: 100,
		ValidatorVoteDecimals:   100,

		AllegationPercentage: 50,
		AllegationDecimals:   100,
	}
	pOpt := governance.ProposalOptionSet{
		ConfigUpdate:      governance.ProposalOption{},
		CodeChange:        governance.ProposalOption{},
		General:           governance.ProposalOption{},
		BountyProgramAddr: "0lt581891d8411faaa15bfd3020b8bd78942cb74ecb",
	}
	govern.SetFeeOption(*feeOpt)
	govern.SetEvidenceOptions(evidenceOption)
	govern.WithHeight(0).SetAllLUH()
	govern.SetProposalOptions(pOpt)

	validator := NewValidator(
		from.Bytes(),
		from.Bytes(),
		keys.PublicKey{keys.ED25519, fromPubkey.Bytes()[5:]},
		keys.PublicKey{keys.ED25519, fromPubkey.Bytes()[5:]},
		*balance.NewAmountFromInt(0),
		"test_node",
	)
	_ = validators.Set(*validator)
	_ = evidenceStore.SetValidatorStatus(validator.Address, true, 1)

	malVal := NewValidator(
		fromMal.Bytes(),
		fromMal.Bytes(),
		keys.PublicKey{keys.ED25519, fromPubkeyMal.Bytes()[5:]},
		keys.PublicKey{keys.ED25519, fromPubkeyMal.Bytes()[5:]},
		*balance.NewAmountFromInt(0),
		"test_node_2",
	)
	_ = validators.Set(*malVal)
	_ = evidenceStore.SetValidatorStatus(malVal.Address, true, 1)

	ctx := NewValidatorContext(
		balances, feePool, delegators, evidenceStore,
		govern, currencies, validators,
	)
	now := time.Now()
	ctx.Validators.lastBlockTime = &now
	ctx.Validators.lastHeight = 4
	return ctx
}

func TestValidatorStore_CheckMaliciousValidators(t *testing.T) {
	baddr1 := keys.Address{}
	_ = baddr1.UnmarshalText([]byte(from.String()))

	baddr2 := keys.Address{}
	_ = baddr2.UnmarshalText([]byte(fromMal.String()))

	t.Run("check malicious validators and no new entries", func(t *testing.T) {
		ctx := setUpCtx()
		defer tearDownCtx()

		votes := []types.VoteInfo{
			{
				Validator: types.Validator{
					Address: baddr1,
				},
				SignedLastBlock: true,
			},
			{
				Validator: types.Validator{
					Address: baddr2,
				},
				SignedLastBlock: true,
			},
		}

		for i := int64(1); i <= 3; i++ {
			eopts, _ := ctx.Govern.GetEvidenceOptions()
			_ = ctx.EvidenceStore.SetVoteBlock(i, votes)
			cv, _ := ctx.EvidenceStore.GetCumulativeVote()
			_ = ctx.EvidenceStore.SetCumulativeVote(cv, i, eopts.BlockVotesDiff)
		}

		err := ctx.Validators.CheckMaliciousValidators(ctx.EvidenceStore, ctx.Govern)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(ctx.Validators.maliciousValidators))
	})

	t.Run("check malicious validators and get malicious validators as not passed criteria", func(t *testing.T) {
		ctx := setUpCtx()
		defer tearDownCtx()

		vMap := make(map[int64][]types.VoteInfo)
		vMap[1] = []types.VoteInfo{
			{
				Validator: types.Validator{
					Address: baddr1,
				},
				SignedLastBlock: false,
			},
			{
				Validator: types.Validator{
					Address: baddr2,
				},
				SignedLastBlock: false,
			},
		}
		vMap[2] = []types.VoteInfo{
			{
				Validator: types.Validator{
					Address: baddr1,
				},
				SignedLastBlock: false,
			},
			{
				Validator: types.Validator{
					Address: baddr2,
				},
				SignedLastBlock: false,
			},
		}
		vMap[3] = []types.VoteInfo{
			{
				Validator: types.Validator{
					Address: baddr1,
				},
				SignedLastBlock: true,
			},
			{
				Validator: types.Validator{
					Address: baddr2,
				},
				SignedLastBlock: false,
			},
		}
		vMap[4] = []types.VoteInfo{
			{
				Validator: types.Validator{
					Address: baddr1,
				},
				SignedLastBlock: false,
			},
			{
				Validator: types.Validator{
					Address: baddr2,
				},
				SignedLastBlock: true,
			},
		}

		for i := int64(1); i <= 3; i++ {
			eopts, _ := ctx.Govern.GetEvidenceOptions()
			_ = ctx.EvidenceStore.SetVoteBlock(i, vMap[i])
			cv, _ := ctx.EvidenceStore.GetCumulativeVote()
			_ = ctx.EvidenceStore.SetCumulativeVote(cv, i, eopts.BlockVotesDiff)
			result, _ := ctx.Validators.store.Commit()
			assert.NotEmpty(t, result)
		}

		cv, _ := ctx.EvidenceStore.GetCumulativeVote()
		assert.Equal(t, 1, len(cv.Addresses))

		err := ctx.Validators.CheckMaliciousValidators(ctx.EvidenceStore, ctx.Govern)
		assert.NoError(t, err)

		assert.Equal(t, 1, len(ctx.Validators.maliciousValidators))

		// test at next block
		ctx.Validators.maliciousValidators = make(map[string]*evidence.LastValidatorHistory)
		i := int64(4)
		eopts, _ := ctx.Govern.GetEvidenceOptions()
		_ = ctx.EvidenceStore.SetVoteBlock(i, vMap[i])
		cv, _ = ctx.EvidenceStore.GetCumulativeVote()
		_ = ctx.EvidenceStore.SetCumulativeVote(cv, i, eopts.BlockVotesDiff)
		result, _ := ctx.Validators.store.Commit()
		assert.NotEmpty(t, result)

		cv, _ = ctx.EvidenceStore.GetCumulativeVote()
		assert.Equal(t, 2, len(cv.Addresses))

		now := time.Now()
		ctx.Validators.lastBlockTime = &now
		ctx.Validators.lastHeight = 5

		err = ctx.Validators.CheckMaliciousValidators(ctx.EvidenceStore, ctx.Govern)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(ctx.Validators.maliciousValidators))
	})
}

func TestValidatorStore_ExecuteAllegationTracker(t *testing.T) {
	t.Run("execute allegation tracker without requests and no changes", func(t *testing.T) {
		ctx := setUpCtx()
		defer tearDownCtx()

		at, err := ctx.EvidenceStore.GetAllegationTracker()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(at.Requests))

		err = ctx.Validators.ExecuteAllegationTracker(ctx, 4)
		assert.NoError(t, err)

		at, err = ctx.EvidenceStore.GetAllegationTracker()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(at.Requests))
	})

	t.Run("execulte allegation tracker with 1 no request and no changes as not pass percentage criteria", func(t *testing.T) {
		ctx := setUpCtx()
		defer tearDownCtx()

		requestID := "test"

		at, err := ctx.EvidenceStore.GetAllegationTracker()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(at.Requests))

		err = ctx.EvidenceStore.PerformAllegation(from.Bytes(), fromMal.Bytes(), requestID, 4, "test")
		assert.NoError(t, err)
		err = ctx.EvidenceStore.Vote(requestID, from.Bytes(), evidence.NO)
		assert.NoError(t, err)

		at, err = ctx.EvidenceStore.GetAllegationTracker()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(at.Requests))

		ar, err := ctx.EvidenceStore.GetAllegationRequest(requestID)
		assert.NoError(t, err)
		assert.Equal(t, evidence.VOTING, ar.Status)

		// not enough percentage to change a status
		err = ctx.Validators.ExecuteAllegationTracker(ctx, 4)
		assert.NoError(t, err)

		at, err = ctx.EvidenceStore.GetAllegationTracker()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(at.Requests))

		ar, err = ctx.EvidenceStore.GetAllegationRequest(requestID)
		assert.NoError(t, err)
		assert.Equal(t, evidence.VOTING, ar.Status)
	})

	t.Run("execulte allegation tracker with 1 yes request and no changes as not pass percentage criteria", func(t *testing.T) {
		ctx := setUpCtx()
		defer tearDownCtx()

		requestID := "test"

		at, err := ctx.EvidenceStore.GetAllegationTracker()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(at.Requests))

		err = ctx.EvidenceStore.PerformAllegation(from.Bytes(), fromMal.Bytes(), requestID, 4, "test")
		assert.NoError(t, err)
		err = ctx.EvidenceStore.Vote(requestID, from.Bytes(), evidence.YES)
		assert.NoError(t, err)

		at, err = ctx.EvidenceStore.GetAllegationTracker()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(at.Requests))

		ar, err := ctx.EvidenceStore.GetAllegationRequest(requestID)
		assert.NoError(t, err)
		assert.Equal(t, evidence.VOTING, ar.Status)

		// not enough percentage to change a status
		err = ctx.Validators.ExecuteAllegationTracker(ctx, 4)
		assert.NoError(t, err)

		at, err = ctx.EvidenceStore.GetAllegationTracker()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(at.Requests))

		ar, err = ctx.EvidenceStore.GetAllegationRequest(requestID)
		assert.NoError(t, err)
		assert.Equal(t, evidence.VOTING, ar.Status)
	})

	t.Run("execulte allegation tracker with 1 no request and have changes as it passes percentage criteria", func(t *testing.T) {
		ctx := setUpCtx()
		defer tearDownCtx()

		at, err := ctx.EvidenceStore.GetAllegationTracker()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(at.Requests))

		requestID := "test"

		err = ctx.EvidenceStore.PerformAllegation(from.Bytes(), fromMal.Bytes(), requestID, 4, "test")
		assert.NoError(t, err)
		err = ctx.EvidenceStore.Vote(requestID, from.Bytes(), evidence.NO)
		assert.NoError(t, err)

		at, err = ctx.EvidenceStore.GetAllegationTracker()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(at.Requests))

		ar, err := ctx.EvidenceStore.GetAllegationRequest(requestID)
		assert.NoError(t, err)
		assert.Equal(t, evidence.VOTING, ar.Status)

		// will meet 100% criteria so the status will be changed as INNOCENT
		err = ctx.Validators.ExecuteAllegationTracker(ctx, 1)
		assert.NoError(t, err)

		// removed as it meet final status
		at, err = ctx.EvidenceStore.GetAllegationTracker()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(at.Requests))

		ar, err = ctx.EvidenceStore.GetAllegationRequest(requestID)
		assert.NoError(t, err)
		assert.Equal(t, evidence.INNOCENT, ar.Status)
	})

	t.Run("execulte allegation tracker with 1 yes request and have changes as it passes percentage criteria", func(t *testing.T) {
		ctx := setUpCtx()
		defer tearDownCtx()

		requestID := "test"

		currency := balance.Currency{
			Id:      0,
			Name:    "OLT",
			Chain:   chain.ONELEDGER,
			Decimal: 18,
			Unit:    "nue",
		}
		decimal := new(big.Int).Exp(big.NewInt(10), big.NewInt(currency.Decimal), nil)

		at, err := ctx.EvidenceStore.GetAllegationTracker()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(at.Requests))

		err = ctx.EvidenceStore.PerformAllegation(from.Bytes(), fromMal.Bytes(), requestID, 4, "test")
		assert.NoError(t, err)
		err = ctx.EvidenceStore.Vote(requestID, from.Bytes(), evidence.YES)
		assert.NoError(t, err)

		at, err = ctx.EvidenceStore.GetAllegationTracker()
		assert.NoError(t, err)
		assert.Equal(t, 1, len(at.Requests))

		ar, err := ctx.EvidenceStore.GetAllegationRequest(requestID)
		assert.NoError(t, err)
		assert.Equal(t, evidence.VOTING, ar.Status)

		err = ctx.Delegators.AddToAddress(fromMal.Bytes(), fromMal.Bytes(), *balance.NewAmount(137))
		assert.NoError(t, err)

		amt, err := ctx.Delegators.GetDelegatorEffectiveAmount(fromMal.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, balance.NewAmount(137).BigInt(), amt.BigInt())

		assert.Equal(t, 0, len(ctx.Validators.GetEvents()))

		// moving to next block
		result, _ := ctx.Validators.store.Commit()
		assert.NotEmpty(t, result)
		ctx.Validators.lastHeight = 2

		// will meet 100% criteria so the status will be changed as GUILTY and bounty and burn program will be executed
		err = ctx.Validators.ExecuteAllegationTracker(ctx, 1)
		assert.NoError(t, err)

		// removed as it meet final status
		at, err = ctx.EvidenceStore.GetAllegationTracker()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(at.Requests))

		ar, err = ctx.EvidenceStore.GetAllegationRequest(requestID)
		assert.NoError(t, err)
		assert.Equal(t, evidence.GUILTY, ar.Status)

		amt, err = ctx.Delegators.GetDelegatorEffectiveAmount(fromMal.Bytes())
		assert.NoError(t, err)
		// 30 percent from 137 amount was withdrawed (96)
		assert.Equal(t, balance.NewAmount(96).BigInt(), amt.BigInt())

		bountyAddress := keys.Address("0lt581891d8411faaa15bfd3020b8bd78942cb74ecb")
		// 50% from 20.5 OLT coins were withdrawed
		bountyCoin, err := ctx.Balances.GetBalanceForCurr(bountyAddress, &currency)
		assert.NoError(t, err)

		rAmt := new(big.Int).Mul(big.NewInt(205), decimal)
		rAmt = new(big.Int).Div(rAmt, big.NewInt(10))
		assert.Equal(t, rAmt, bountyCoin.Amount.BigInt())

		assert.Equal(t, 1, len(ctx.Validators.GetEvents()))

		// moving to next block
		result, _ = ctx.Validators.store.Commit()
		assert.NotEmpty(t, result)
		ctx.Validators.lastHeight = 3

		// also this amount will be applyied (41) to withdraw and recalculate of power
		unstake, err := ctx.Validators.GetDelayUnstake(fromMal.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, balance.NewAmount(41).BigInt(), unstake.Amount.BigInt())
	})
}
