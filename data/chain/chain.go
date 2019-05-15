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

package chain

import (
	"github.com/pkg/errors"
)

type Type int

type Chain struct {
	ChainType   Type
	Description string
	Features    []string
}

var chainTypes = map[string]Type{}
var chainTypeNames = map[Type]string{}

func init() {
	RegisterChainType("OneLedger", 0)
}

func RegisterChainType(name string, id int) {
	chainTypes[name] = Type(id)
	chainTypeNames[Type(id)] = name
}

func (ctype Type) String() string {

	name, ok := chainTypeNames[ctype]
	if !ok {
		return "INVALID"
	}

	return name
}

func TypeFromName(chainName string) (Type, error) {
	typ, ok := chainTypes[chainName]
	if !ok {
		return Type(-1), errors.New("wrong chain name")
	}

	return typ, nil
}
