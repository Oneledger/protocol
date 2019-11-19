/*

 */

package bitcoin

import "github.com/btcsuite/btcd/chaincfg"

func GetChainParams(typeString string) *chaincfg.Params {

	var params *chaincfg.Params
	switch typeString {
	case "mainnet":
		params = &chaincfg.MainNetParams
	case "testnet3":
		params = &chaincfg.TestNet3Params
	case "regtest":
		params = &chaincfg.RegressionNetParams
	case "simnet":
		params = &chaincfg.SimNetParams
	default:
		params = &chaincfg.TestNet3Params
	}

	return params
}
