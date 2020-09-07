package common

import abci "github.com/tendermint/tendermint/abci/types"

type HandlerList []func(*ExtAppData, abci.Header)

var Handlers HandlerList

var FunctionRouter ControllerRouter

func init() {
	Handlers = HandlerList{}
	FunctionRouter = ControllerRouter{}
}
