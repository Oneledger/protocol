package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/balance"
	"github.com/Oneledger/protocol/log"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	genesis, _ := consensus.NewGenesisDoc("bazbar", []balance.Currency{}, []consensus.StateInput{})
	input := &initContext{
		genesis:  genesis,
		logger:   log.NewDefaultLogger(os.Stdout).WithPrefix("fooe"),
		rootDir:  dir,
		nodeName: "newton-Node",
	}
	err := initNode(input)
	assert.Nil(t, err, "Running initNode shouldn't return an error")

	// These are the files that should exist and be nonempty
	fileTree := map[string]bool{
		filepath.Join(dir, "config.toml"):                                    true,
		filepath.Join(dir, "consensus", "config", "genesis.json"):            true,
		filepath.Join(dir, "consensus", "config", "node_key.json"):           true,
		filepath.Join(dir, "consensus", "config", "priv_validator_key.json"): true,
		filepath.Join(dir, "consensus", "data", "priv_validator_state.json"): true,
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if path == dir {
			return nil
		}

		_, ok := fileTree[path]
		if ok {
			delete(fileTree, path)
		}
		assert.NotZero(t, info.Size(), "Got empty file "+info.Name())

		return nil
	})

	assert.Nil(t, err, "Walk should succeed")
	assert.Empty(t, fileTree)
}
