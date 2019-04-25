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

type ChainType int

type Chain struct {
	ChainType   ChainType
	Description string
	Features    []string
}

var chainTypes = map[string]ChainType{}
var chainTypeNames = map[ChainType]string{}

func RegisterChainType(name string, id int) {
	chainTypes[name] = ChainType(id)
	chainTypeNames[ChainType(id)] = name
}

func (ctype ChainType) String() string {

	name, ok := chainTypeNames[ctype]
	if !ok {
		return "INVALID"
	}

	return name
}
