package config

type Client struct {
	Node NodeConfig `toml:"node"`

	BroadcastMode string `toml:"async"`
	Proof         bool   `toml:"proof"`
}
