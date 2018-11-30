/*
	Copyright 2017-2018 OneLedger

	Define aliases to hide Tendermint types and enhance readability of the code
*/
package app

import "github.com/tendermint/tendermint/abci/types"

type RequestInitChain = types.RequestInitChain
type ResponseInitChain = types.ResponseInitChain

type RequestInfo = types.RequestInfo
type ResponseInfo = types.ResponseInfo

type RequestQuery = types.RequestQuery
type ResponseQuery = types.ResponseQuery

type RequestSetOption = types.RequestSetOption
type ResponseSetOption = types.ResponseSetOption

type RequestCheckTx = types.RequestCheckTx
type ResponseCheckTx = types.ResponseCheckTx

type RequestBeginBlock = types.RequestBeginBlock
type ResponseBeginBlock = types.ResponseBeginBlock

type RequestDeliverTx = types.RequestDeliverTx
type ResponseDeliverTx = types.ResponseDeliverTx

type RequestEndBlock = types.RequestEndBlock
type ResponseEndBlock = types.ResponseEndBlock

type ResponseCommit = types.ResponseCommit
