/*
	Copyright 2017-2018 OneLedger

	Common errors across the entire system
*/
package status

type Code = uint32 // Alias to match Tendermint status

const (
	SUCCESS         Code = 0
	INVALID         Code = 101
	PARSE_ERROR     Code = 201
	NOT_IMPLEMENTED Code = 301
	MISSING_VALUE   Code = 401
	EXPAND_ERROR    Code = 501
	DUPLICATE       Code = 601
	MISSING_DATA    Code = 701
	BAD_VALUE       Code = 801
	EXECUTE_ERROR   Code = 901
)

// TODO: Code should be a real type, so that String really works?
//func (code Code) String() string {
func String(code Code) string {
	switch code {
	case SUCCESS:
		return "SUCCESS"
	case INVALID:
		return "INVALID"
	case PARSE_ERROR:
		return "PARSE_ERROR"
	case NOT_IMPLEMENTED:
		return "NOT_IMPLEMENTED"
	case MISSING_VALUE:
		return "MISSING_VALUE"
	case EXPAND_ERROR:
		return "EXPAND_ERROR"
	case DUPLICATE:
		return "DUPLICATE"
	case MISSING_DATA:
		return "MISSING_DATA"
	case BAD_VALUE:
		return "BAD_VALUE"
	case EXECUTE_ERROR:
		return "EXECUTE_ERROR"
	}
	return "UNKNOWN ERR"
}
