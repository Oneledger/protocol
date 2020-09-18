package farm_data

import (
	"github.com/Oneledger/protocol/external_apps/farm_produce/farm_error"
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	ErrFailedInSerialization   = codes.ProtocolError{farm_error.ErrFailedInSerialization, "failed to serialize"}
	ErrFailedInDeserialization = codes.ProtocolError{farm_error.ErrFailedInDeserialization, "failed to deserialize"}
	ErrSettingRecord           = codes.ProtocolError{farm_error.ErrSettingRecord, "failed to set record"}
	ErrGettingRecord           = codes.ProtocolError{farm_error.ErrGettingRecord, "failed to get record"}
	ErrDeletingRecord          = codes.ProtocolError{farm_error.ErrDeletingRecord, "failed to delete record"}
	ErrInvalidBatchID          = codes.ProtocolError{farm_error.ErrInvalidBatchID, "invalid batch ID"}
	ErrInvalidFarmID           = codes.ProtocolError{farm_error.ErrInvalidFarmID, "invalid farm ID"}
	ErrGettingProduceStore     = codes.ProtocolError{farm_error.ErrGettingProduceStore, "failed to get produce store"}
	ErrBatchIDAlreadyExists    = codes.ProtocolError{farm_error.ErrBatchIDAlreadyExists, "this batch ID already exists in store"}
	ErrInsertingProduct        = codes.ProtocolError{farm_error.ErrInsertingProduct, "failed to insert product batch into the store"}
)
