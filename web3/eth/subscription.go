package eth

import (
	"context"

	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
)

func (svc *Service) NewHeads(ctx context.Context) (*rpc.Subscription, error) {
	return svc.filterAPI.NewHeads(ctx)
}

func (svc *Service) NewPendingTransactions(ctx context.Context) (*rpc.Subscription, error) {
	return svc.filterAPI.NewPendingTransactions(ctx)
}

func (svc *Service) Logs(ctx context.Context, crit filters.FilterCriteria) (*rpc.Subscription, error) {
	return svc.filterAPI.Logs(ctx, crit)
}
