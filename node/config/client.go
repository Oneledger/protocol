package config

import (
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/viper"
)

func ClientConfig() {
	viper.SetConfigName(global.Current.ConfigName)

	// NOTE: Directories need the trailing slash
	viper.AddConfigPath(".")                    // Local directory override
	viper.AddConfigPath(global.Current.RootDir) // Special user overrides

	err := viper.ReadInConfig()
	if err != nil {
		log.Info("Not using config file", "err", err)
	}
}
