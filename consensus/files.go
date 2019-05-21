package consensus

import (
	"path/filepath"
)

// Common file names & directories
const (
	GenesisFilename            = "genesis.json"
	NodeKeyFilename            = "node_key.json"
	PrivValidatorKeyFilename   = "priv_validator_key.json"
	PrivValidatorStateFilename = "priv_validator_state.json"
	AddrBookFilename           = "addrbook.json"

	RootDirName   = "consensus"
	ConfigDirName = "config"
	DataDirName   = "data"
)

// Dir returns the root folder for consensus files
func Dir(rootDir string) string {
	if !filepath.IsAbs(rootDir) {
		// Please don't pass an invalid rootdir
		rootDir, _ = filepath.Abs(rootDir)
	}
	return filepath.Join(rootDir, RootDirName)
}
