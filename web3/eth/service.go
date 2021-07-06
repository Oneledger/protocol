package eth

import (
	"errors"
	"sync"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/storage"
	"github.com/ethereum/go-ethereum/rpc"

	rpcclient "github.com/Oneledger/protocol/client"
)

type Service struct {
	log     *log.Logger
	ext     client.ExtServiceContext
	stateDB *action.CommitStateDB

	mu sync.Mutex
}

func NewService(logger *log.Logger, ext client.ExtServiceContext, stateDB *action.CommitStateDB) *Service {
	return &Service{
		log:     logger,
		ext:     ext,
		stateDB: stateDB,
	}
}

func (svc *Service) getTMClient() rpcclient.Client {
	return svc.ext.RPCClient()
}

func (svc *Service) getState() *storage.State {
	return svc.stateDB.GetContractStore().State
}

func (svc *Service) getStateHeight(height int64) int64 {
	switch height {
	case -1, -2:
		return svc.getState().Version()
	}
	return height
}

func (svc *Service) stateAndHeaderByNumberOrHash(blockNrOrHash rpc.BlockNumberOrHash) (int64, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return svc.getStateHeight(blockNr.Int64()), nil
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header, err := svc.getTMClient().BlockByHash(hash.Bytes())
		if err != nil {
			return 0, err
		}
		if header == nil || header.Block == nil {
			return 0, errors.New("header for hash not found")
		}
		return header.Block.Header.Height, nil
	}
	return 0, errors.New("invalid arguments; neither block nor hash specified")
}
