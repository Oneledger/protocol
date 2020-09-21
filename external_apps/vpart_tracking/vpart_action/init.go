package vpart_action

import "github.com/Oneledger/protocol/action"

const (
	VPART_INSERT action.Type = 990201
)

func init() {
	action.RegisterTxType(VPART_INSERT, "VPART_INSERT")
}

