package action

import "github.com/Oneledger/protocol/serialize"


type Type int

const (
	SEND Type = iota
)

func init() {

	serialize.RegisterInterface(new(Msg))
	serialize.RegisterConcrete(new(Send), "action_send")

}

func (t Type) String() string {
	switch t {
	case SEND:
		return "SEND"
	default:
		return "UNKNOWN"
	}
}