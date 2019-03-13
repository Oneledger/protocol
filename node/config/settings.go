/*
	Copyright 2017-2018 OneLedger
*/
package config

// Log all of the global settings
import (
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
)

// Print the current settings to the log file
func LogSettings() {
	log.Info("Diagnostics", "Debug", global.Current.Debug, "DisablePasswords", global.Current.DisablePasswords)

	log.Info("Ownership", "NodeName", global.Current.NodeName, "NodeAccountName", global.Current.NodeAccountName,
		"NodeIdentity", global.Current.NodeIdentity)

	log.Info("Locations", "RootDir", global.Current.RootDir)
	log.Info("Addresses", "RpcAddress", global.Current.RpcAddress)

	log.Info("Tendermint", "TendermintAddress", global.Current.TendermintAddress,
		"TendermintPubKey", global.Current.TendermintPubKey)
	log.Info("SDK", "SDKAddress", global.Current.SDKAddress)
	//log.Info("OLVM", "OLVMAddress", global.Current.OLVMAddress)

	log.Info("Bitcoin", "BTCAddress", global.Current.BTCAddress)
	log.Info("Ethereum", "ETHAddress", global.Current.ETHAddress)
}
