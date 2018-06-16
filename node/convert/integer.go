/*
	Copyright 2017 - 2018 OneLedger
*/
package convert

import (
	"fmt"
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

func GetString(value int) string {
	return fmt.Sprintf("%d", value)
}

func GetString64(value int64) string {
	return fmt.Sprintf("%d", value)
}
