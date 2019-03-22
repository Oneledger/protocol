package data

import (
	"os"
	"testing"

	"github.com/Oneledger/protocol/node/config"
	"github.com/Oneledger/protocol/node/global"
)

func setup() {
	global.Current.RootDir = "./"
	global.Current.Config = config.DefaultServerConfig()
}

func teardown() {
	os.RemoveAll(global.Current.DatabaseDir())
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
