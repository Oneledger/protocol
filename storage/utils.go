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

package storage

import (
	"os"
	"path/filepath"
)

const (
	OLDATA   = "OLDATA"
	NODEDATA = "nodedata"
)

func dbDir() string {

	rootDir := os.Getenv(OLDATA)
	result, _ := filepath.Abs(filepath.Join(rootDir, NODEDATA))

	return result
}

func Prefix(prefix string) []byte {
	return []byte(prefix + DB_PREFIX)
}

func Rangefix(prefix string) []byte {
	if string(prefix[len(prefix)-1]) == DB_PREFIX {
		prefix = prefix[:len(prefix)-1]
	}
	a := []byte(prefix + DB_RANGEFIX)
	return a
}
