/*
	Copyright 2017-2018 OneLedger

	Implement all of the query mechanics for the node and the chain
*/
package app

import "strings"

// Top-level list of all query types
func HandleQuery(path string, message []byte) []byte {
	switch path {
	case "/account":
		return HandleAccountQuery(message)
	}
	return HandleError("Unknown Path", path, message)
}

// Get the account information for a given user
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

// Return a nicely formatted error message
func HandleError(text string, path string, massage []byte) []byte {
	return []byte("Invalid Query")
}
