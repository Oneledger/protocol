package identity

import (
	"fmt"
	"math"
	"math/big"
	"sort"

	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/evidence"
	govern "github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/kv"
)

// TODO: Reduce power in allegation tracker for GUILTY
// TODO: Add staking recalculation

// CheckMaliciousValidators that not pass criterias and will be marked as malicious
func (vs *ValidatorStore) CheckMaliciousValidators(es *evidence.EvidenceStore, govern *govern.Store) error {
	vs.maliciousValidators = make(map[string]*evidence.LastValidatorHistory)
	evidenceOptions, err := govern.GetEvidenceOptions()
	if err != nil {
		logger.Fatal("failed to get the evidence options")
	}

	// skip checks if does not met height criteria
	if vs.lastHeight <= evidenceOptions.BlockVotesDiff {
		logger.Infof("Height must be more than equal %d for check. Current: %d\n", evidenceOptions.BlockVotesDiff, vs.lastHeight)
		return nil
	}

	if vs.lastBlockTime == nil {
		logger.Info("Last block time is nil")
		return nil
	}

	cv, err := es.GetCumulativeVote()
	if err != nil {
		return err
	}

	// fetch previous suspicious validators
	es.IterateSuspiciousValidators(func(lvh *evidence.LastValidatorHistory) bool {
		vs.maliciousValidators[lvh.Address.String()] = lvh
		return false
	})

	addresses := make([]string, 0, len(cv.Addresses))
	for addr := range cv.Addresses {
		addresses = append(addresses, addr)
	}
	sort.Strings(addresses)

	// update found addresses with missed votes
	for _, addr := range addresses {
		votes := cv.Addresses[addr]
		baddr := keys.Address{}
		err = baddr.UnmarshalText([]byte(addr))
		if err != nil {
			continue
		}
		if votes < evidenceOptions.MinVotesRequired {
			key := append(vs.prefix, baddr...)
			data := vs.store.GetVersioned(vs.lastHeight-1, key)
			if len(data) == 0 {
				logger.Errorf("Previous state data not found for address: %s", baddr)
				continue
			}
			validator := &Validator{}
			if err := serialize.GetSerializer(serialize.JSON).Deserialize(data, validator); err != nil {
				logger.Errorf("Validator: %s not found", baddr)
				continue
			}
			validatorStatus, _ := es.GetValidatorStatus(validator.Address)
			if validatorStatus != nil && validatorStatus.IsActive {
				// get active from starting height as when it incomed it could be malicious, this will prevent this
				plusDiff := validatorStatus.Height + evidenceOptions.BlockVotesDiff
				if plusDiff > vs.lastHeight {
					continue
				}
				logger.Detailf("Found validator with missed required votes: %s\n", validator.Address)
				lvh, err := es.CreateSuspiciousValidator(
					baddr, evidence.MISSED_REQUIRED_VOTES,
					vs.lastHeight, vs.lastBlockTime)
				if err != nil {
					continue
				}
				vs.maliciousValidators[addr] = lvh
			}
		}
	}
	return nil
}

func (vs *ValidatorStore) ExecuteAllegationTracker(ctx *ValidatorContext, activeCount int64) error {
	if vs.lastBlockTime == nil {
		return fmt.Errorf("Last block time not set")
	}

	if activeCount == 0 {
		return fmt.Errorf("Active count could not be zero")
	}

	cs := ctx.Currencies.GetCurrencies().GetCurrencySet()
	currency, _ := cs.GetCurrencyByName("OLT")
	decimal := new(big.Int).Exp(big.NewInt(10), big.NewInt(currency.Decimal), nil)

	at, err := ctx.EvidenceStore.GetAllegationTracker()
	if err != nil {
		return err
	}

	options, err := ctx.Govern.GetEvidenceOptions()
	if err != nil {
		return err
	}

	requiredVotesCount := int(math.Ceil(float64(activeCount) * float64(options.ValidatorVotePercentage) / float64(options.ValidatorVoteDecimals)))

	popt, err := ctx.Govern.GetProposalOptions()
	if err != nil {
		return err
	}

	bountyAddress := keys.Address(popt.BountyProgramAddr)

	addrToDelete := make([]string, 0)
	//processedValidators := make(map[string]bool)
	ctx.EvidenceStore.CleanTracker()
	for requestID := range at.Requests {
		ar, err := ctx.EvidenceStore.GetAllegationRequest(requestID)
		decisionMade := false
		if err != nil {
			//logger.Errorf("Failed to retrieve allegation request: %s\n", err)
			continue
		}

		//_, ok := processedValidators[ar.MaliciousAddress.Humanize()]
		//if ok {
		//	addrToDelete = append(addrToDelete, requestID)
		//	continue
		//}
		yesCount := 0
		noCount := 0

		for i := range ar.Votes {
			vote := ar.Votes[i]
			switch vote.Choice {
			case evidence.YES:
				yesCount++
			case evidence.NO:
				noCount++
			}
		}

		yesP := float64(yesCount) / float64(requiredVotesCount)
		noP := float64(noCount) / float64(requiredVotesCount)
		percentage := float64(options.AllegationPercentage) / float64(options.AllegationDecimals)
		arToUpdate := false

		logger.Detailf("Request ID: %s, yes votes count: %d, no votes count: %d, total count: %d \n", requestID, yesCount, noCount, requiredVotesCount)
		if yesP > percentage {
			decisionMade = true
			ar.Status = evidence.GUILTY
			sv, err := ctx.EvidenceStore.CreateSuspiciousValidator(
				ar.MaliciousAddress, evidence.BYZANTINE_FAULT,
				vs.lastHeight, vs.lastBlockTime)
			if err != nil {
				logger.Errorf("Failed to create suspicious validator: %s\n", err)
				continue
			}
			logger.Detailf("Suspicious validator created: %s\n", sv.Address)

			// taking malicoius validator data
			addrHuman := ar.MaliciousAddress.Humanize()
			key := append(vs.prefix, ar.MaliciousAddress...)
			data := vs.store.GetVersioned(vs.lastHeight-1, key)
			if len(data) == 0 {
				logger.Errorf("Previous state data not found for address: %s\n", addrHuman)
				continue
			}
			validator := &Validator{}
			if err := serialize.GetSerializer(serialize.JSON).Deserialize(data, validator); err != nil {
				logger.Errorf("Validator: %s not found\n", addrHuman)
				continue
			}
			// retrieving balance
			amt, err := ctx.Delegators.GetValidatorAmount(validator.Address)
			if err != nil {
				logger.Errorf("Failed to get balance from delegators: %s\n", err)
				continue
			}

			// calculate evidence percent
			penalizationAmt := new(big.Float).Mul(
				amt.BigFloat(),
				big.NewFloat(float64(options.PenaltyBasePercentage)),
			)
			penalizationAmt = new(big.Float).Quo(
				penalizationAmt,
				big.NewFloat(float64(options.PenaltyBaseDecimals)),
			)
			penalizationAmt.Add(penalizationAmt, new(big.Float).SetFloat64(0.5))
			pAmt, _ := penalizationAmt.Int(nil)

			// calculate bounty percent
			bountyAmt := new(big.Int).Mul(pAmt, decimal)
			bountyAmt = new(big.Int).Mul(
				bountyAmt,
				balance.NewAmount(options.PenaltyBountyPercentage).BigInt(),
			)
			bountyAmt = new(big.Int).Div(
				bountyAmt,
				balance.NewAmount(options.PenaltyBountyDecimals).BigInt(),
			)
			// starting minusing the bounty
			err = ctx.Delegators.MinusFromAddress(validator.Address, validator.StakeAddress, *balance.NewAmountFromBigInt(pAmt))
			if err == nil {
				bountyCoin := balance.Coin{
					Currency: currency,
					Amount:   balance.NewAmountFromBigInt(bountyAmt),
				}
				err = ctx.Balances.AddToAddress(bountyAddress, bountyCoin)
				if err != nil {
					logger.Errorf("Failed to add balance on bounty program: %s\n", err)
					continue
				}
				logger.Detailf("Successfully added bounty coin to bounty address: %s\n", bountyAddress.Humanize())
			} else {
				logger.Detailf("Nothing to withdraw from addr on bounty program: %s\n", bountyAddress.Humanize())
			}
			//processedValidators[ar.MaliciousAddress.Humanize()] = true
			addrToDelete = append(addrToDelete, requestID)
			arToUpdate = true
			vs.createAllegationEvent(ar)

			// refreshing balance
			amt, err = ctx.Delegators.GetValidatorAmount(validator.Address)
			if err != nil {
				logger.Errorf("Failed to get balance from delegators: %s\n", err)
				continue
			}
			// postpone the update for next block
			err = vs.delayHandleUnstake(validator.Address, *balance.NewAmountFromBigInt(pAmt))
			if err != nil {
				logger.Errorf("Failed to update postponed: %s\n", err)
				continue
			}
		} else if noP > 1-percentage {
			decisionMade = true
			//processedValidators[ar.MaliciousAddress.Humanize()] = true
			ar.Status = evidence.INNOCENT
			addrToDelete = append(addrToDelete, requestID)
			arToUpdate = true
			vs.createAllegationEvent(ar)
		}

		if arToUpdate {
			err := ctx.EvidenceStore.SetAllegationRequest(ar)
			if err != nil {
				return err
			}
		}
		if decisionMade {
			logger.Infof("Decision made on Validator, Deleting Allegation Request :%s", ar.String())
			ctx.EvidenceStore.DeleteAllegationRequest(ar.ID)
		}
	}
	update := false
	for i := range addrToDelete {
		if !update {
			update = true
		}
		delete(at.Requests, addrToDelete[i])
	}
	if update {
		return ctx.EvidenceStore.SetAllegationTracker(at)
	}

	return nil
}

func (vs *ValidatorStore) getDelayUnstakeKey(height int64, address keys.Address) []byte {
	key := append(vs.prefixPurge, []byte(fmt.Sprintf("unstake_%d", height))...)
	key = append(key, address...)
	return key
}

func (vs *ValidatorStore) SetDelayUnstake(unstake *Unstake) error {
	key := vs.getDelayUnstakeKey(vs.lastHeight, unstake.Address)
	value, err := serialize.GetSerializer(serialize.PERSISTENT).Serialize(unstake)
	if err != nil {
		return errors.Wrap(err, "failed to serialize unstake")
	}

	err = vs.store.Set(key, value)
	if err != nil {
		return errors.Wrap(err, "failed to set unstake")
	}
	return nil
}

func (vs *ValidatorStore) GetDelayUnstake(addr keys.Address) (*Unstake, error) {
	key := vs.getDelayUnstakeKey(vs.lastHeight-1, addr)
	dat, _ := vs.store.Get(key)
	if len(dat) == 0 {
		return nil, fmt.Errorf("Postponed unstake not found")
	}

	apply := &Unstake{}
	err := serialize.GetSerializer(serialize.PERSISTENT).Deserialize(dat, apply)
	if err != nil {
		return nil, err
	}
	return apply, nil
}

func (vs *ValidatorStore) fetchPostponedUnstakes() error {
	vs.Iterate(func(addr keys.Address, validator *Validator) bool {
		unstake, err := vs.GetDelayUnstake(validator.Address)
		if err != nil {
			return false
		}
		err = vs.HandleUnstake(*unstake, vs.lastHeight)
		if err != nil {
			logger.Errorf("Handle unstake for validator: %s failed, %s\n", validator.Address, err)
			return false
		}
		logger.Infof("Unstake %s for validator: %s was applied!\n", unstake.Amount, validator.Address)
		return false
	})
	return nil
}

func (vs *ValidatorStore) delayHandleUnstake(addr keys.Address, amt balance.Amount) error {
	apply := &Unstake{
		Address: addr,
		Amount:  amt,
	}
	err := vs.SetDelayUnstake(apply)
	if err != nil {
		return err
	}
	return nil
}

func (vs *ValidatorStore) createAllegationEvent(ar *evidence.AllegationRequest) {
	tags := make([]kv.Pair, 3)
	statusBytes := make([]byte, 1)
	statusBytes[0] = byte(ar.Status)
	tags = append(
		tags,
		kv.Pair{
			Key:   []byte("block.reporter"),
			Value: ar.ReporterAddress.Bytes(),
		},
		kv.Pair{
			Key:   []byte("block.malicious"),
			Value: ar.MaliciousAddress.Bytes(),
		},
		kv.Pair{
			Key:   []byte("block.status"),
			Value: statusBytes,
		},
	)
	vs.PushEvent("allegation_tracker", tags)
}
