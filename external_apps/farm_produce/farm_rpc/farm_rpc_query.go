package farm_rpc

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/external_apps/farm_produce/farm_data"
	"github.com/Oneledger/protocol/log"
)

type Service struct {
	balances     *balance.Store
	currencies   *balance.CurrencySet
	ons          *ons.DomainStore
	logger       *log.Logger
	productStore *farm_data.ProductStore
}

func Name() string {
	return "farm_query"
}

func NewService(balances *balance.Store, currencies *balance.CurrencySet,
	domains *ons.DomainStore, logger *log.Logger, productStore *farm_data.ProductStore) *Service {
	return &Service{
		currencies:   currencies,
		balances:     balances,
		ons:          domains,
		logger:       logger,
		productStore: productStore,
	}
}

func (svc *Service) GetBatchByID(req GetBatchByIDRequest, reply *GetBatchByIDReply) error {
	batch, err := svc.productStore.Get(req.BatchID)
	if err != nil {
		return ErrGettingProductBatchInQuery.Wrap(err)
	}

	*reply = GetBatchByIDReply{
		ProductBatch: *batch,
		Height:       svc.productStore.GetState().Version(),
	}
	return nil
}
