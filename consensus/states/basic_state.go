package states

import (
	"github.com/Oneledger/prototype/consensus/types"
	"sync"
)

type BasicState interface {
	Enter (event types.Event) *BasicState
	Exit (event types.Event) *BasicState
	Process ()
}



type Propose struct {

}



var instance *Propose
var once sync.Once

func GetPropose() *Propose {
	once.Do(func() {
		instance = &Propose{}
	})
	return instance
}

type NewHeight struct {

}
