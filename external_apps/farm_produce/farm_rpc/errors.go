package farm_rpc

import (
	"github.com/Oneledger/protocol/external_apps/farm_produce/farm_error"
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	ErrGettingProduceBatchInQuery = codes.ProtocolError{farm_error.ErrGettingProduceBatchInQuery, "failed to get produce batch in query"}
)
