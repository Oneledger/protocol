package main

import (
	"os"
)

func setEnvVariables() {
	os.Setenv("API_KEY", "de5e96cbb6284d5ea1341bf6cb7fa401")
	os.Setenv("ETHPKPATH", "/home/tanmay/Codebase/pkdata")
	os.Setenv("WALLETADDR", "/home/tanmay/Codebase/walletAddr")
	ethBlockConfirmation = 0
}
