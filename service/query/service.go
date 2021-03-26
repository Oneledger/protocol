package query

import (
	"encoding/hex"
	"strings"

	"github.com/Oneledger/protocol/data/evm"
	netwkDeleg "github.com/Oneledger/protocol/data/network_delegation"
	"github.com/Oneledger/protocol/data/rewards"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/delegation"
	"github.com/Oneledger/protocol/data/evidence"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	codes "github.com/Oneledger/protocol/status_codes"
	"github.com/Oneledger/protocol/utils"
)

type Service struct {
	name            string
	ext             client.ExtServiceContext
	balances        *balance.Store
	currencies      *balance.CurrencySet
	validators      *identity.ValidatorStore
	witnesses       *identity.WitnessStore
	delegators      *delegation.DelegationStore
	netwkDelegators *netwkDeleg.MasterStore
	evidenceStore   *evidence.EvidenceStore
	govern          *governance.Store
	ons             *ons.DomainStore
	feePool         *fees.Store
	proposalMaster  *governance.ProposalMasterStore
	rewardMaster    *rewards.RewardMasterStore
	governance      *governance.Store
	logger          *log.Logger
	txTypes         *[]action.TxTypeDescribe
	contracts       *evm.ContractStore
	accountKeeper   balance.AccountKeeper
}

func Name() string {
	return "query"
}

func NewService(ctx client.ExtServiceContext, balances *balance.Store, currencies *balance.CurrencySet, validators *identity.ValidatorStore, witnesses *identity.WitnessStore,
	domains *ons.DomainStore, delegators *delegation.DelegationStore, netwkDelegators *netwkDeleg.MasterStore, evidenceStore *evidence.EvidenceStore, govern *governance.Store, feePool *fees.Store, proposalMaster *governance.ProposalMasterStore, rewardMaster *rewards.RewardMasterStore, logger *log.Logger, txTypes *[]action.TxTypeDescribe,
	contracts *evm.ContractStore, accountKeeper balance.AccountKeeper,
) *Service {
	service := &Service{
		name:            "query",
		ext:             ctx,
		currencies:      currencies,
		balances:        balances,
		validators:      validators,
		witnesses:       witnesses,
		delegators:      delegators,
		netwkDelegators: netwkDelegators,
		evidenceStore:   evidenceStore,
		govern:          govern,
		ons:             domains,
		feePool:         feePool,
		proposalMaster:  proposalMaster,
		rewardMaster:    rewardMaster,
		logger:          logger,
		txTypes:         txTypes,
		governance:      govern,
		contracts:       contracts,
		accountKeeper:   accountKeeper,
	}
	return service
}

func (svc *Service) Balance(req client.BalanceRequest, resp *client.BalanceReply) error {
	err := req.Address.Err()
	if err != nil {
		return codes.ErrBadAddress
	}

	addr := req.Address
	bal, err := svc.balances.GetBalance(addr, svc.currencies)

	if err != nil {
		svc.logger.Error("error getting balance", err)
		return codes.ErrGettingBalance
	}

	*resp = client.BalanceReply{
		Balance: bal.String(),
		Height:  svc.balances.State.Version(),
	}
	return nil
}

func (svc *Service) BalancePool(req client.BalancePoolRequest, resp *client.BalanceReply) error {

	poolname := req.Poolname
	poolList, err := svc.governance.GetPoolList()
	if err != nil {
		return err
	}
	bal := &balance.Balance{}
	if pool, ok := poolList[poolname]; ok {
		bal, err = svc.balances.GetBalance(pool, svc.currencies)
		if err != nil {
			svc.logger.Error("error getting balance", err)
			return codes.ErrGettingBalance
		}
	}

	*resp = client.BalanceReply{
		Balance: bal.String(),
		Height:  svc.balances.State.Version(),
	}
	return nil
}

// ListValidator returns a list of all validator
func (svc *Service) ListValidators(_ client.ListValidatorsRequest, reply *client.ListValidatorsReply) error {
	validators, err := svc.validators.GetValidatorSet()
	if err != nil {
		svc.logger.Error("error listing validators")
		return codes.ErrListValidators
	}

	vMap := svc.evidenceStore.GetValidatorMap()
	fMap := svc.evidenceStore.GetFrozenMap()

	*reply = client.ListValidatorsReply{
		Validators: validators,
		Height:     svc.balances.State.Version(),
		VMap:       vMap,
		FMap:       fMap,
	}
	return nil
}

// ListWitnesses returns a list of all witness
func (svc *Service) ListWitnesses(req client.ListWitnessesRequest, reply *client.ListWitnessesReply) error {
	witnesses, err := svc.witnesses.GetWitnessAddresses(req.ChainType)
	if err != nil {
		svc.logger.Error("error listing witnesses")
		return codes.ErrListWitnesses
	}

	*reply = client.ListWitnessesReply{
		Witnesses: witnesses,
		Height:    svc.balances.State.Version(),
	}
	return nil
}

func (svc *Service) ListCurrencies(_ client.ListCurrenciesRequest, reply *client.ListCurrenciesReply) error {
	reply.Currencies = svc.currencies.GetCurrencies()
	return nil
}

func (svc *Service) CurrencyBalance(req client.CurrencyBalanceRequest, resp *client.CurrencyBalanceReply) error {
	err := req.Address.Err()
	if err != nil {
		return codes.ErrBadAddress
	}

	currency, ok := svc.currencies.GetCurrencyByName(req.Currency)
	if !ok {
		return codes.ErrFindingCurrency
	}

	addr := req.Address
	bal, err := svc.balances.GetBalance(addr, svc.currencies)

	if err != nil {
		svc.logger.Error("error getting balance", err)
		return codes.ErrGettingBalance
	}

	coin := bal.GetCoin(currency)

	*resp = client.CurrencyBalanceReply{
		Currency: currency.Name,
		Balance:  coin.Humanize(),
		Height:   svc.balances.State.Version(),
	}
	return nil
}

func (svc *Service) VoteRequests(req client.VoteRequestRequest, resp *client.VoteRequestReply) error {
	requests := make([]evidence.AllegationRequest, 0)

	// TODO: Add filter for address
	svc.evidenceStore.IterateRequests(func(ar *evidence.AllegationRequest) bool {
		if len(req.Address) == 0 || ar.MaliciousAddress.Equal(req.Address) {
			requests = append(requests, *ar)
		}
		return false
	})

	*resp = client.VoteRequestReply{
		Requests: requests,
	}

	return nil
}

func (svc *Service) ValidatorStatus(req client.ValidatorStatusRequest, resp *client.ValidatorStatusReply) error {
	err := req.Address.Err()
	if err != nil {
		return codes.ErrBadAddress
	}

	exists := false
	validator, err := svc.validators.Get(req.Address)
	if err != nil {
		*resp = client.ValidatorStatusReply{
			Exists:                exists,
			Height:                svc.balances.State.Version(),
			Staking:               "0",
			TotalDelegationAmount: "0",
			SelfDelegationAmount:  "0",
			DelegationAmount:      "0",
		}
		return nil
	}

	svc.logger.Infof("Validator - %s, delegator - %s\n", validator.Address.Humanize(), validator.StakeAddress.Humanize())

	totalDelegationAmount, _ := svc.delegators.GetValidatorAmount(validator.Address)
	selfDelegationAmount, _ := svc.delegators.GetValidatorDelegationAmount(validator.Address, validator.StakeAddress)
	delegationAmount, _ := totalDelegationAmount.Minus(*selfDelegationAmount)

	if validator.Power > 0 {
		exists = true
	}

	*resp = client.ValidatorStatusReply{
		Exists:                exists,
		Height:                svc.balances.State.Version(),
		Power:                 validator.Power,
		Staking:               validator.Staking.String(),
		TotalDelegationAmount: totalDelegationAmount.String(),
		SelfDelegationAmount:  selfDelegationAmount.String(),
		DelegationAmount:      delegationAmount.String(),
	}

	return nil
}

func (svc *Service) DelegationStatus(req client.DelegationStatusRequest, resp *client.DelegationStatusReply) error {
	err := req.Address.Err()
	if err != nil {
		return codes.ErrBadAddress
	}

	options, _ := svc.govern.GetStakingOptions()
	if options == nil {
		return codes.ErrFlagNotSet
	}

	effectiveDelegationAmount, _ := svc.delegators.GetDelegatorEffectiveAmount(req.Address)
	withdrawableAmount, _ := svc.delegators.GetDelegatorBoundedAmount(req.Address)

	height := svc.balances.State.Version()
	maturedAmounts := svc.delegators.GetMaturedPendingAmount(req.Address, height, options.MaturityTime+1)

	bal, err := svc.balances.GetBalance(req.Address, svc.currencies)

	if err != nil {
		svc.logger.Error("error getting balance", err)
		return codes.ErrGettingBalance
	}

	*resp = client.DelegationStatusReply{
		Balance:                   bal.String(),
		EffectiveDelegationAmount: effectiveDelegationAmount.String(),
		WithdrawableAmount:        withdrawableAmount.String(),
		MaturedAmounts:            maturedAmounts,
	}

	return nil
}

func (svc *Service) FeeOptions(_ struct{}, reply *client.FeeOptionsReply) error {
	*reply = client.FeeOptionsReply{
		FeeOption: *svc.feePool.GetOpt(),
	}
	return nil
}

func (svc *Service) ListTxTypes(_ client.ListTxTypesRequest, reply *client.ListTxTypesReply) error {
	var txTypes []action.TxTypeDescribe
	//find all const types that less than EOF marker
	//and not "UNKNOWN"(this also prevents the potential future const that not be the type: "Type"
	//from showing up)
	for i := 0; i < int(action.EOF); i++ {
		if strings.Compare(action.Type(i).String(), "UNKNOWN") != 0 {
			txTypeDescribe := action.TxTypeDescribe{
				TxTypeNum:    action.Type(i),
				TxTypeString: action.Type(i).String(),
			}
			txTypes = append(txTypes, txTypeDescribe)
		}
	}

	*reply = client.ListTxTypesReply{
		TxTypes: txTypes,
	}
	return nil
}

func (svc *Service) Tx(req client.TxRequest, reply *client.TxResponse) error {
	hash, err := hex.DecodeString(utils.TrimHex(req.Hash))
	res, err := svc.ext.Tx(hash, req.Prove)
	if err != nil {
		return codes.ErrGetTx.Wrap(err)
	}

	*reply = client.TxResponse{
		Result: *res,
	}

	return nil
}
