/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package utils

import (
	"errors"
	"net/url"
)

// Hex prefix
const HexPrefix = "0x"

func PrefixHex(in string) string {
	return HexPrefix + in
}

var ErrParsingAddress = errors.New("failed to parse network address")

// Pick out the port from a full address
func GetPort(addr string) (string, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return "", err
	} else if u.Port() == "" {
		return "", ErrParsingAddress
	}

	return u.Port(), nil
}
