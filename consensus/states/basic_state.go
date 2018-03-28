package states

import (
	"github.com/Oneledger/prototype/consensus/types"
)

type BasicState interface {
	Enter (previous BasicState, event types.Event)
	Exit (next BasicState, event types.Event )
	Process ()
}
