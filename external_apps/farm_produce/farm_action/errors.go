package farm_action

import (
	"github.com/Oneledger/protocol/external_apps/farm_produce/farm_error"
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	ErrFailedToUnmarshal = codes.ProtocolError{farm_error.ErrFailedToUnmarshal, "failed to unmarshal"}
)
