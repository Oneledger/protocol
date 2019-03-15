package config

type Client struct {
	Node NodeConfig `toml:"node"`
}

// ClientConfig loads the configuration for the client onto the global Context
func ClientConfig() {
	// viper.SetConfigName(global.Current.ConfigName)

	// NOTE: Directories need the trailing slash
}
