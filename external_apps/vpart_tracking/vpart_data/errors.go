package vpart_data

import (
	"github.com/Oneledger/protocol/external_apps/vpart_tracking/vpart_error"
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	ErrFailedInSerialization   = codes.ProtocolError{vpart_error.ErrFailedInSerialization, "failed to serialize"}
	ErrFailedInDeserialization = codes.ProtocolError{vpart_error.ErrFailedInDeserialization, "failed to deserialize"}
	ErrSettingRecord           = codes.ProtocolError{vpart_error.ErrSettingRecord, "failed to set record"}
	ErrGettingRecord           = codes.ProtocolError{vpart_error.ErrGettingRecord, "failed to get record"}
	ErrDeletingRecord          = codes.ProtocolError{vpart_error.ErrDeletingRecord, "failed to delete record"}
)
