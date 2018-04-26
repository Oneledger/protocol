/*
	Copyright 2017-2018 OneLedger
*/
package app

import "strings"

func HandleQuery(path string, message []byte) []byte {
	switch path {
	case "/account":
		return HandleAccountQuery(message)
	}
	return HandleError("Unknown Path", path, message)
}

func HandleAccountQuery(message []byte) []byte {
	text := string(message)

	parts := strings.Split(text, "=")
	if len(parts) == 2 {
		if parts[0] == "User" {
			return []byte("User Information")
		}
	}
	return []byte("Unknown Query")
}

func HandleError(text string, path string, massage []byte) []byte {
	return []byte("Invalid Query")
}
