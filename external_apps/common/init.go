package common

type HandlerList []func(*ExtAppData)

var Handlers HandlerList

func init() {
	Handlers = HandlerList{}
}
