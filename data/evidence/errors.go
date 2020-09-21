package evidence

import (
	codes "github.com/Oneledger/protocol/status_codes"
)

var (
	ErrCreateAllegationFailed = codes.ProtocolError{codes.TxErrEvidenceError, "failed to create allegation request"}
	ErrHandleReleaseFailed    = codes.ProtocolError{codes.TxErrEvidenceError, "failed to handle release"}
	ErrFrozenValidator        = codes.ProtocolError{codes.TxErrEvidenceError, "error frozen validator"}
	ErrNonFrozenValidator     = codes.ProtocolError{codes.TxErrEvidenceError, "error non frozen validator"}
	ErrNonActiveValidator     = codes.ProtocolError{codes.TxErrEvidenceError, "non active validator"}
)
