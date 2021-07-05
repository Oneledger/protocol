package utils

import (
	"errors"

	rpcclient "github.com/Oneledger/protocol/client"
	"github.com/ethereum/go-ethereum/rpc"
)

func StateAndHeaderByNumberOrHash(client rpcclient.Client, blockNrOrHash rpc.BlockNumberOrHash) (int64, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return blockNr.Int64(), nil
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header, err := client.BlockByHash(hash.Bytes())
		if err != nil {
			return 0, err
		}
		if header == nil {
			return 0, errors.New("header for hash not found")
		}
		if blockNrOrHash.RequireCanonical && b.eth.blockchain.GetCanonicalHash(header.Number.Uint64()) != hash {
			return 0, errors.New("hash is not currently canonical")
		}
		return header.Block.Height, nil
	}
	return 0, errors.New("invalid arguments; neither block nor hash specified")
}
