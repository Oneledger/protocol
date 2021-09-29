package eth

import (
	"context"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethfilters "github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
)

func (svc *Service) NewFilter(crit ethfilters.FilterCriteria) (rpc.ID, error) {
	return svc.filterAPI.NewFilter(crit)
}

func (svc *Service) NewBlockFilter() rpc.ID {
	return svc.filterAPI.NewBlockFilter()
}

func (svc *Service) NewPendingTransactionFilter() rpc.ID {
	return svc.filterAPI.NewPendingTransactionFilter()
}

func (svc *Service) UninstallFilter(id rpc.ID) bool {
	return svc.filterAPI.UninstallFilter(id)
}

func (svc *Service) GetFilterChanges(id rpc.ID) (interface{}, error) {
	return svc.filterAPI.GetFilterChanges(id)
}

func (svc *Service) GetFilterLogs(id rpc.ID) ([]*ethtypes.Log, error) {
	return svc.filterAPI.GetFilterLogs(context.Background(), id)
}

func (svc *Service) GetLogs(crit ethfilters.FilterCriteria) ([]*ethtypes.Log, error) {
	return svc.filterAPI.GetLogs(context.Background(), crit)
}
