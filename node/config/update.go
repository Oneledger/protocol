/*
	Copyright 2017-2018 OneLedger
*/

package config

import (
	"reflect"

	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/viper"
)

type Parameter struct {
	Name     string
	DataType string
}

var updateParameters = []Parameter{
	Parameter{"NodeName", "string"},
	Parameter{"RpcAddress", "string"},
	Parameter{"AppAddress", "string"},
	Parameter{"P2PAddress", "string"},
	Parameter{"SDKAddress", "string"},
	Parameter{"BTCAddress", "string"},
	Parameter{"ETHAddress", "string"},
}

func UpdateContext() {
	//global.Current.SDKAddress = viper.Get("SDKAddress").(string)
	valueOf := reflect.ValueOf(global.Current).Elem()

	for _, parameter := range updateParameters {
		switch parameter.DataType {
		case "string":
			param := viper.Get(parameter.Name).(string)
			if param != "" {
				field := valueOf.FieldByName(parameter.Name)
				if field.IsValid() {
					field.SetString(param)
				}
			}
		default:
			log.Warn("Unknown Config Parameter", "name", parameter.Name)
		}
	}
}
