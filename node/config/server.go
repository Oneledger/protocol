package config

import (
	"github.com/Oneledger/protocol/node/global"
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
)

func ServerConfig() {
	viper.SetConfigName(global.Current.ConfigName)

	viper.AddConfigPath("$OLSCRIPT/data/") // Common script configs
	viper.AddConfigPath("~/.olfullnode")   // Special user overrides
	viper.AddConfigPath(".")               // Local directory override

	err := viper.ReadInConfig()
	if err != nil {
		log.Info("Not using config file", "err", err)
	}
}
