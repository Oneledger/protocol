package governance

import (
	"github.com/Oneledger/protocol/log"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/status_codes"
)

//Helper Function to log error
func logAndReturnFalse(logger *log.Logger, error status_codes.ProtocolError, tags kv.Pairs, err error) (bool, action.Response) {
	logger.Error(error)
	result := action.Response{
		Events: action.GetEvent(tags, error.Msg),
		Log:    error.Wrap(err).Marshal(),
	}
	return false, result
}

func logAndReturnTrue(logger *log.Logger, tags kv.Pairs, eventType string) (bool, action.Response) {
	logger.Detail(eventType)
	result := action.Response{
		Events: action.GetEvent(tags, eventType),
	}
	return true, result
}
