/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package client

import (
	"github.com/pkg/errors"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var ErrEmptyTransaction = errors.New("empty transaction")

func (ctx ExtServiceContext) BroadcastTxSync(packet []byte) (res *ctypes.ResultBroadcastTx, err error) {

	if len(packet) < 1 {
		return res, ErrEmptyTransaction
	}

	client := ctx.rpcClient
	result, err := client.BroadcastTxSync(packet)
	if err != nil {
		return res, err
	}

	return result, nil
}

func (ctx ExtServiceContext) BroadcastTxAsync(packet []byte) (res *ctypes.ResultBroadcastTx, err error) {

	if len(packet) < 1 {
		return nil, ErrEmptyTransaction
	}

	client := ctx.rpcClient
	result, err := client.BroadcastTxAsync(packet)
	if err != nil {
		return res, err
	}

	return result, nil
}

func (ctx ExtServiceContext) BroadcastTxCommit(packet []byte) (res *ctypes.ResultBroadcastTxCommit, err error) {

	if len(packet) < 1 {
		return nil, ErrEmptyTransaction
	}

	client := ctx.rpcClient
	result, err := client.BroadcastTxCommit(packet)
	if err != nil {
		return nil, err
	}

	return result, nil
}
