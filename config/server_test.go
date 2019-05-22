package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Oneledger/protocol/config"
	"github.com/stretchr/testify/assert"
	tmconfig "github.com/tendermint/tendermint/config"
)

var path string

func handleErr(err error) {
	panic(err)
}

func init() {
	// Declared here to suppress IDE error about unused 'path' variable
	var err error
	path, err = filepath.Abs(filepath.Join(".", "config_test"))
	if err != nil {
		handleErr(err)
	}

}

func setup() {
	err := os.MkdirAll(path, config.DirPerms)
	if err != nil {
		handleErr(err)
	}
}

func teardown() {
	_ = os.RemoveAll(path)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestServer_FileHandling(t *testing.T) {
	cfg := config.DefaultServerConfig()
	cfgPath := filepath.Join(path, config.FileName)

	t.Run("SaveFile should fail on saving to a nonexistent path", func(t *testing.T) {
		nonexistingPath := filepath.Join("nonexistent", "path", "file")
		err := cfg.SaveFile(nonexistingPath)
		assert.Error(t, err)
	})

	t.Run("SaveFile should save the file successfully", func(t *testing.T) {
		err := cfg.SaveFile(cfgPath)
		assert.Nil(t, err)
		assert.FileExists(t, cfgPath, "Should successfully create the file")
	})

	t.Run("ReadFile should fail when reading a nonexistent file", func(t *testing.T) {
		err := cfg.ReadFile(filepath.Join(cfgPath, "nonexistent.toml"))
		assert.Error(t, err, "Reading nonexistent file should return an error")
	})

	t.Run("ReadFile should correctly read an existing file", func(t *testing.T) {
		cfg := new(config.Server)
		defaultCFG := config.DefaultServerConfig()
		err := cfg.ReadFile(cfgPath)
		assert.Nil(t, err)
		assert.Equal(t, *cfg, *defaultCFG, "Should read the default config")
	})
}

func TestTMConfig_Durations(t *testing.T) {
	cfg := config.DefaultServerConfig()
	tmcfg := cfg.TMConfig()
	tmcfgDefault := tmconfig.DefaultConfig()

	int3Array := func(d config.Duration, td time.Duration, tdDefault time.Duration) [3]int64 {
		return [3]int64{int64(d), int64(td), int64(tdDefault)}
	}

	// Assert all of the millisecond values get converted to nanosecond values
	fieldDurations := map[string][3]int64{
		"FlushThrottleTimeout":        int3Array(cfg.P2P.FlushThrottleTimeout, tmcfg.P2P.FlushThrottleTimeout, tmcfgDefault.P2P.FlushThrottleTimeout),
		"HandshakeTimeout":            int3Array(cfg.P2P.HandshakeTimeout, tmcfg.P2P.HandshakeTimeout, tmcfgDefault.P2P.HandshakeTimeout),
		"DialTimeout":                 int3Array(cfg.P2P.DialTimeout, tmcfg.P2P.DialTimeout, tmcfgDefault.P2P.DialTimeout),
		"TimeoutPropose":              int3Array(cfg.Consensus.TimeoutPropose, tmcfg.Consensus.TimeoutPropose, tmcfgDefault.Consensus.TimeoutPropose),
		"TimeoutProposeDelta":         int3Array(cfg.Consensus.TimeoutProposeDelta, tmcfg.Consensus.TimeoutProposeDelta, tmcfgDefault.Consensus.TimeoutProposeDelta),
		"TimeoutPrevote":              int3Array(cfg.Consensus.TimeoutPrevote, tmcfg.Consensus.TimeoutPrevote, tmcfgDefault.Consensus.TimeoutPrevote),
		"TimeoutPrevoteDelta":         int3Array(cfg.Consensus.TimeoutPrevoteDelta, tmcfg.Consensus.TimeoutPrevoteDelta, tmcfgDefault.Consensus.TimeoutPrevoteDelta),
		"TimeoutPrecommit":            int3Array(cfg.Consensus.TimeoutPrecommit, tmcfg.Consensus.TimeoutPrecommit, tmcfgDefault.Consensus.TimeoutPrecommit),
		"TimeoutPrecommitDelta":       int3Array(cfg.Consensus.TimeoutPrecommitDelta, tmcfg.Consensus.TimeoutPrecommitDelta, tmcfgDefault.Consensus.TimeoutPrecommitDelta),
		"TimeoutCommit":               int3Array(cfg.Consensus.TimeoutCommit, tmcfg.Consensus.TimeoutCommit, tmcfgDefault.Consensus.TimeoutCommit),
		"CreateEmptyBlocksInterval":   int3Array(cfg.Consensus.CreateEmptyBlocksInterval, tmcfg.Consensus.CreateEmptyBlocksInterval, tmcfgDefault.Consensus.CreateEmptyBlocksInterval),
		"PeerGossipSleepDuration":     int3Array(cfg.Consensus.PeerGossipSleepDuration, tmcfg.Consensus.PeerGossipSleepDuration, tmcfgDefault.Consensus.PeerGossipSleepDuration),
		"PeerQueryMaj23SleepDuration": int3Array(cfg.Consensus.PeerQueryMaj23SleepDuration, tmcfg.Consensus.PeerQueryMaj23SleepDuration, tmcfgDefault.Consensus.PeerQueryMaj23SleepDuration),
	}

	for field, duration := range fieldDurations {
		olDefault, nano, defaultNano := duration[0], duration[1], duration[2]

		// First check the times were converted to milliseconds properly
		defaultWant := defaultNano / int64(time.Millisecond)
		assert.Equalf(t, defaultWant, olDefault, "Failed to generate proper default duration value in %s, want: %d, got: %d", field, defaultWant, olDefault)

		// Next check
		assert.Equalf(t, nano, defaultNano, "Failed to convert TMConfig Duration values for field %s, want: %d, got: %d", field, nano, defaultNano)
	}
}
