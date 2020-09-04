package common

import "fmt"

type HandlerList []func(*ExtAppData)

var Handlers HandlerList

func init() {
	//todo remove
	fmt.Println("init from common/init")
	Handlers = HandlerList{}
}
