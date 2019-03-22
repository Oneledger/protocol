package log

import (
	"flag"
	"os"
	"testing"

	"github.com/Oneledger/protocol/node/global"
)

// Control the execution
func TestMain(m *testing.M) {
	flag.Parse()

	// Set the debug flags according to whether the -v flag is set in go test
	if testing.Verbose() {
		Debug("DEBUG TURNED ON")
		global.Current.Debug = true
	} else {
		Debug("DEBUG TURNED OFF")
		global.Current.Debug = false
	}

	// Run it all.
	code := m.Run()

	os.Exit(code)
}
