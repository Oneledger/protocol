package config

import (
	"errors"
	"path/filepath"
	"strings"
)

var setup map[string]func(cfg *Server, value string) error

func init() {
	setup = make(map[string]func(cfg *Server, value string) error)
	setup["persistent_peers"] = persistent

}

func Setup(cfg *Server, key string, value string) error {
	if fn, ok := setup[key]; ok {
		return fn(cfg, value)
	}
	return errors.New("setup key not available")
}

func persistent(cfg *Server, value string) error {
	peers := strings.Split(value, ",")
	cfg.P2P.SetPersistentPeers(peers)
	err := cfg.SaveFile(filepath.Join(cfg.rootDir, FileName))
	return err
}
