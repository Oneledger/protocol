package vpart_rpc

import (
	"github.com/Oneledger/protocol/external_apps/vpart_tracking/vpart_error"
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	ErrGettingVPartInQuery   = codes.ProtocolError{vpart_error.ErrGettingVPartInQuery, "failed to get vehicle part in query"}
)
