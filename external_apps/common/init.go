package common

type HandlerList []func(*ExtAppData)

var Handlers HandlerList

var FunctionRouter ControllerRouter

func init() {
	Handlers = HandlerList{}
	FunctionRouter = ControllerRouter{}
}
