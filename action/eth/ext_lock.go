package eth

import "github.com/Oneledger/protocol/action"

type Lock struct {
	Locker      action.Address
	TrackerName string
	ETHTxn      []byte
	LockAmount  int64
}
