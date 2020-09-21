package vpart_rpc

import (
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/data/ons"
	"github.com/Oneledger/protocol/external_apps/vpart_tracking/vpart_data"
	"github.com/Oneledger/protocol/log"
)

type Service struct {
	balances   *balance.Store
	currencies *balance.CurrencySet
	ons        *ons.DomainStore
	logger     *log.Logger
	vPartStore *vpart_data.VPartStore
}

func Name() string {
	return "vpart_query"
}

func NewService(balances *balance.Store, currencies *balance.CurrencySet, logger *log.Logger, vPartStore *vpart_data.VPartStore) *Service {
	return &Service{
		currencies: currencies,
		balances:   balances,
		logger:     logger,
		vPartStore: vPartStore,
	}
}

func (svc *Service) GetVPart(req GetVehiclePartRequest, reply *GetVehiclePartReply) error {
	vPart, err := svc.vPartStore.Get(req.VIN, req.PartType)
	if err != nil {
		return ErrGettingVPartInQuery.Wrap(err)
	}

	*reply = GetVehiclePartReply{
		VehiclePart: *vPart,
		Height:      svc.vPartStore.GetState().Version(),
	}
	return nil
}
