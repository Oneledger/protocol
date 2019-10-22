package query

import (
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	codes "github.com/Oneledger/protocol/status_codes"
)

type Service struct {
	name       string
	ext        client.ExtServiceContext
	balances   *balance.Store
	currencies *balance.CurrencySet
	validators *identity.ValidatorStore
	ons        *ons.DomainStore
	logger     *log.Logger
}

func Name() string {
	return "query"
}

func NewService(ctx client.ExtServiceContext, balances *balance.Store, currencies *balance.CurrencySet, validators *identity.ValidatorStore,
	domains *ons.DomainStore, logger *log.Logger) *Service {
	return &Service{
		name:       "query",
		ext:        ctx,
		currencies: currencies,
		balances:   balances,
		validators: validators,
		ons:        domains,
		logger:     logger,
	}
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

// ListValidator returns a list of all validator
func (svc *Service) ListValidators(_ client.ListValidatorsRequest, reply *client.ListValidatorsReply) error {
	validators, err := svc.validators.GetValidatorSet()
	if err != nil {
		svc.logger.Error("error listing validators")
		return codes.ErrListValidators
	}

	*reply = client.ListValidatorsReply{
		Validators: validators,
		Height:     svc.balances.State.Version(),
	}

	return nil
}

func (svc *Service) ListCurrencies(_ client.ListCurrenciesRequest, reply *client.ListCurrenciesReply) error {
	reply.Currencies = svc.currencies.GetCurrencies()
	return nil
}
