package states

import (
	"github.com/Oneledger/prototype/consensus/types"
)

type BasicState interface {
	Enter (event types.Event) *BasicState
	Exit (event types.Event) *BasicState
	Process ()
}
