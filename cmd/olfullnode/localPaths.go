package main

import (
	"os"
)

func setEnvVariablesInfura() {
	os.Setenv("API_KEY", "de5e96cbb6284d5ea1341bf6cb7fa401")
	os.Setenv("ETHPKPATH", "/tmp/pkdata")
	os.Setenv("WALLETADDR", "/tmp/walletAddr")
	ethBlockConfirmation = 12
	totalETHSupply = "20000000000000000000"
}

func setEnvVariablesGanache() {
	os.Setenv("API_KEY", "de5e96cbb6284d5ea1341bf6cb7fa401")
	os.Setenv("ETHPKPATH", "/tmp/pkdataGanache")
	os.Setenv("WALLETADDR", "/tmp/walletAddr")
	ethBlockConfirmation = 12
	totalETHSupply = "20000000000000000000"
}
