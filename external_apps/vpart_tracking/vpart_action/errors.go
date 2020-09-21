package vpart_action

import (
	"github.com/Oneledger/protocol/external_apps/vpart_tracking/vpart_error"
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	ErrFailedToUnmarshal = codes.ProtocolError{vpart_error.ErrFailedToUnmarshal, "failed to unmarshal"}
	ErrGettingVPartStore = codes.ProtocolError{vpart_error.ErrGettingVPartStore, "failed to get vPart store"}

)
