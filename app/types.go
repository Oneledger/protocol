package app

import (
	abci "github.com/tendermint/tendermint/abci/types"
)

type RequestInitChain = abci.RequestInitChain
type ResponseInitChain = abci.ResponseInitChain

type RequestInfo = abci.RequestInfo
type ResponseInfo = abci.ResponseInfo

type RequestQuery = abci.RequestQuery
type ResponseQuery = abci.ResponseQuery

type RequestSetOption = abci.RequestSetOption
type ResponseSetOption = abci.ResponseSetOption

type RequestCheckTx = abci.RequestCheckTx
type ResponseCheckTx = abci.ResponseCheckTx

type RequestBeginBlock = abci.RequestBeginBlock
type ResponseBeginBlock = abci.ResponseBeginBlock

type RequestDeliverTx = abci.RequestDeliverTx
type ResponseDeliverTx = abci.ResponseDeliverTx

type RequestEndBlock = abci.RequestEndBlock
type ResponseEndBlock = abci.ResponseEndBlock

type ResponseCommit = abci.ResponseCommit

type Header = abci.Header

type Validator = abci.Validator

type ABCIApp = abci.Application
