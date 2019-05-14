package app

import (
	"errors"
	"os"

	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
)

// NodeContext holds key information about the running node. This is generally used to
// to access its address and perform signing functions
type NodeContext struct {
	NodeName string

	// Node keys
	privateKey keys.PrivateKey

	// Validator key
	privval *privval.FilePV
}

// PrivVal returns the private validator file
func (n NodeContext) PrivVal() *privval.FilePV {
	return n.privval
}

// Address returns the address of the node's public key (the key's hash)
func (n NodeContext) Address() keys.Address {
	priv, err := n.privateKey.GetHandler()
	if err != nil {
		return nil
	}

	pub, err := priv.PubKey().GetHandler()
	if err != nil {
		return nil
	}

	return pub.Address()
}

func (n NodeContext) ValidatorAddress() keys.Address {
	return keys.Address(n.privval.GetPubKey().Address())
}

func (n NodeContext) isValid() bool {
	if n.privval == nil || n.Address() == nil {
		return false
	} //else if n.NodeName == "" { return false }
	return true
}

// NewNodeContext returns a NodeContext by reading from the specified configuration files.
// This function WILL exit if the private validator key files (priv_validator_state, and priv_validator_key) don't
// exist in the configured location
func NewNodeContext(cfg *config.Server) (*NodeContext, error) {
	if cfg == nil {
		return nil, errors.New("got nil argument")
	}

	// (1) Read the node_keys
	// Read the public key
	tmcfg := consensus.Config(cfg.TMConfig())
	ctx, err := readKeyFiles(&tmcfg)
	if err != nil {
		return nil, err
	}
	ctx.NodeName = cfg.Node.NodeName

	if !ctx.isValid() {
		return nil, errors.New("failed to read keys properly")
	}

	return ctx, nil
}

func readKeyFiles(cfg *consensus.Config) (*NodeContext, error) {
	nodeKeyF := cfg.NodeKeyFile()

	nodeKey, err := p2p.LoadNodeKey(nodeKeyF)
	if err != nil {
		return nil, err
	}

	priv, err := keys.NodeKeyFromTendermint(nodeKey)
	if err != nil {
		return nil, err
	}

	pvKeyF := cfg.PrivValidatorKeyFile()
	pvStateF := cfg.PrivValidatorStateFile()
	paths := []string{pvKeyF, pvStateF}
	for _, path := range paths {
		_, err := os.Stat(path)
		if err != nil {
			return nil, errors.New("failed to find " + path)
		}
	}

	// This function quits the process if either of these files don't exist
	filePV := privval.LoadFilePV(pvKeyF, pvStateF)

	return &NodeContext{
		privateKey: priv,
		privval:    filePV,
	}, nil
}
