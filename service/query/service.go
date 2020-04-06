package query

import (
	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/identity"
	"github.com/Oneledger/protocol/log"
	codes "github.com/Oneledger/protocol/status_codes"
	"strings"
)

type Service struct {
	name       string
	ext        client.ExtServiceContext
	balances   *balance.Store
	currencies *balance.CurrencySet
	validators *identity.ValidatorStore
	ons        *ons.DomainStore
	feePool    *fees.Store
	logger     *log.Logger
	txTypes    *[]action.TxTypeDescribe
}

func Name() string {
	return "query"
}

func NewService(ctx client.ExtServiceContext, balances *balance.Store, currencies *balance.CurrencySet, validators *identity.ValidatorStore,
	domains *ons.DomainStore, feePool *fees.Store, logger *log.Logger, txTypes *[]action.TxTypeDescribe) *Service {
	return &Service{
		name:       "query",
		ext:        ctx,
		currencies: currencies,
		balances:   balances,
		validators: validators,
		ons:        domains,
		feePool:    feePool,
		logger:     logger,
		txTypes:    txTypes,
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
	for i := 0; i < int(action.EOF); i++{
		if strings.Compare(action.Type(i).String(), "UNKNOWN") != 0 {
			txTypeDescribe := action.TxTypeDescribe{
				TxTypeNum:       action.Type(i),
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