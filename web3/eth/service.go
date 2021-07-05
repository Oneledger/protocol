package eth

import (
	"errors"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/log"
	"github.com/ethereum/go-ethereum/rpc"

	rpcclient "github.com/Oneledger/protocol/client"
)

type Service struct {
	log     *log.Logger
	ext     client.ExtServiceContext
	stateDB *action.CommitStateDB
}

func NewService(logger *log.Logger, ext client.ExtServiceContext, stateDB *action.CommitStateDB) *Service {
	return &Service{
		log:     logger,
		ext:     ext,
		stateDB: stateDB,
	}
}

func StateAndHeaderByNumberOrHash(client rpcclient.Client, blockNrOrHash rpc.BlockNumberOrHash) (int64, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return blockNr.Int64(), nil
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header, err := client.BlockByHash(hash.Bytes())
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
