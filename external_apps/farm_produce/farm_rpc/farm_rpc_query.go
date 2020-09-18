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
	produceStore *farm_data.ProduceStore
}

func Name() string {
	return "farm_query"
}

func NewService(balances *balance.Store, currencies *balance.CurrencySet,
	domains *ons.DomainStore, logger *log.Logger, productStore *farm_data.ProduceStore) *Service {
	return &Service{
		currencies:   currencies,
		balances:     balances,
		ons:          domains,
		logger:       logger,
		produceStore: productStore,
	}
}

func (svc *Service) GetBatchByID(req GetBatchByIDRequest, reply *GetBatchByIDReply) error {
	batch, err := svc.produceStore.Get(req.BatchID)
	if err != nil {
		return ErrGettingProduceBatchInQuery.Wrap(err)
	}

	*reply = GetBatchByIDReply{
		ProduceBatch: *batch,
		Height:       svc.produceStore.GetState().Version(),
	}
	return nil
}
