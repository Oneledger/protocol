package common_action

import (
	"errors"

	"github.com/Oneledger/protocol/log"

	"github.com/tendermint/tendermint/libs/kv"

	"github.com/Oneledger/protocol/status_codes"
)

//Helper Function to log error
func LogAndReturnFalse(logger *log.Logger, error status_codes.ProtocolError, tags kv.Pairs, err error) (bool, Response) {
	if err == nil {
		err = errors.New("No Err String")
	}
	logger.Error(error)
	result := Response{
		Events: GetEvent(tags, error.Msg),
		Log:    error.Wrap(err).Marshal(),
	}
	return false, result
}

func LogAndReturnTrue(logger *log.Logger, tags kv.Pairs, eventType string) (bool, Response) {
	logger.Detail(eventType)
	result := Response{
		Events: GetEvent(tags, eventType),
	}
	return true, result
}
