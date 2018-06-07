/*
	Copyright 2017-2018 OneLedger

	Common errors across the entire system
*/
package err

type Code = uint32 // Matches Tendermint status

const (
	SUCCESS         Code = 0
	PARSE_ERROR     Code = 101
	NOT_IMPLEMENTED Code = 201
	MISSING_VALUE   Code = 301
	EXPAND_ERROR    Code = 401
	DUPLICATE       Code = 501
)
