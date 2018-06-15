/*
	Copyright 2017 - 2018 OneLedger
*/
package convert

import (
	"strconv"
)

func GetInt(value string, defaultValue int) int {
	// TODO: Should be ParseInt and should specific 64 or 32
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}
