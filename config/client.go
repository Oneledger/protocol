package config

import "github.com/Oneledger/protocol-temp"

type Client struct {
	Node protocol_temp.NodeConfig `toml:"node"`

	BroadcastMode string `toml:"async"`
	Proof         bool   `toml:"proof"`
}
