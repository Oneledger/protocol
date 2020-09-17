package farm_rpc

import (
	"github.com/Oneledger/protocol/external_apps/farm_produce/farm_error"
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	ErrGettingProductBatchInQuery       = codes.ProtocolError{farm_error.ErrGettingProductBatchInQuery, "failed to get product batch in query"}
)
