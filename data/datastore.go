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

package data

import "encoding/hex"

type StoreKey []byte

func (sk StoreKey) Bytes() []byte {
	return sk
}

func (sk StoreKey) String() string {
	return hex.EncodeToString(sk)
}

type Store interface {
	Get(StoreKey) ([]byte, error)
	Set(StoreKey, []byte) error
	Exists(StoreKey) (bool, error)
	Delete(StoreKey) (bool, error)
}
