package query

import (
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/rpc"
)

type Service struct {
	name       string
	balances   *balance.Store
	currencies *balance.CurrencyList
	validators *identity.ValidatorStore
	logger     *log.Logger
}

func Name() string {
	return "query"
}

func NewService(balances *balance.Store, currencies *balance.CurrencyList, validators *identity.ValidatorStore, logger *log.Logger) *Service {
	return &Service{
		name:       "query",
		currencies: currencies,
		balances:   balances,
		validators: validators,
		logger:     logger,
	}
}

func (svc *Service) Balance(req client.BalanceRequest, resp *client.BalanceReply) error {
	addr := req.Address
	bal, err := svc.balances.Get(addr, true)

	if err != nil && err == balance.ErrNoBalanceFoundForThisAddress {
		// Return a zero for balance if the account doesn't exist
		// TODO: Zero in the balances
		bal = balance.NewBalance()
	} else if err != nil {
		svc.logger.Error("error getting balance", err)
		return rpc.InternalError("error getting balance")
	}

	*resp = client.BalanceReply{
		Balance: *bal,
		Height:  svc.balances.Version,
	}
	return nil
}

// ListValidator returns a list of all validator
func (svc *Service) ListValidators(_ client.ListValidatorsRequest, reply *client.ListValidatorsReply) error {
	validators, err := svc.validators.GetValidatorSet()
	if err != nil {
		return rpc.InternalError("err while retrieving validators info")
	}

	*reply = client.ListValidatorsReply{
		Validators: validators,
		Height:     svc.balances.Version,
	}

	return nil
}
