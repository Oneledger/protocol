package node

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/btcsuite/btcd/btcec"

	hdwallet "github.com/Oneledger/hdkeychain"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
)

// Context holds key information about the running node. This is generally used to
// to access its address and perform signing functions
type Context struct {
	NodeName string

	// Node keys
	privateKey keys.PrivateKey

	// Validator key
	privval keys.PrivateKey

	// Validator hd keywords
	keywords []string
}

// PrivVal returns the private validator file
func (n Context) PrivVal() keys.PrivateKey {
	return n.privval
}

// private key of the nodes
func (n Context) PrivKey() keys.PrivateKey {
	return n.privateKey
}

// PubKey returns the public key of the node's NodeKey
func (n Context) PubKey() keys.PublicKey {
	h, err := n.privateKey.GetHandler()
	if err != nil {
		return keys.PublicKey{}
	}

	return h.PubKey()
}

// Address returns the address of the node's public key (the key's hash)
func (n Context) Address() keys.Address {
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

func (n Context) ValidatorPubKey() keys.PublicKey {
	priv, err := n.privval.GetHandler()
	if err != nil {
		return keys.PublicKey{}
	}

	return priv.PubKey()
}

func (n Context) ValidatorAddress() keys.Address {
	priv, err := n.privval.GetHandler()
	if err != nil {
		return nil
	}

	pub, err := priv.PubKey().GetHandler()
	if err != nil {
		return nil
	}

	return pub.Address()
}

func (n Context) ValidatorBTCPubKey(i uint32) keys.PublicKey {

	hdw, err := hdwallet.NewHDWalletFromKeywords(n.keywords, "")
	if err != nil {
		return keys.PublicKey{}
	}

	pubKey := &btcec.PublicKey{}

	_, pubKey, err = hdw.GetBTCExternalKeyPair(i)

	pkey, err := keys.GetPublicKeyFromBytes(pubKey.SerializeCompressed(), keys.BTCECSECP)
	return pkey
}

func (n Context) ValidatorBTCPrivateKey(i uint32) *keys.PrivateKey {
	hdw, err := hdwallet.NewHDWalletFromKeywords(n.keywords, "")
	if err != nil {
		return &keys.PrivateKey{}
	}

	privKey := &btcec.PrivateKey{}

	privKey, _, err = hdw.GetBTCExternalKeyPair(i)

	pkey, err := keys.GetPrivateKeyFromBytes(privKey.Serialize(), keys.BTCECSECP)
	return &pkey
}

func (n Context) ValidatorHDWallet() (hdwallet.HDWallet, error) {

	return hdwallet.NewHDWalletFromKeywords(n.keywords, "")
}

func (n Context) isValid() bool {
	if n.ValidatorAddress() == nil || n.Address() == nil {
		return false
	} //else if n.NodeName == "" { return false }
	return true
}

// NewNodeContext returns a Context by reading from the specified configuration files.
// This function WILL exit if the private validator key files (priv_validator_state, and priv_validator_key) don't
// exist in the configured location
func NewNodeContext(cfg *config.Server) (*Context, error) {
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

func readKeyFiles(cfg *consensus.Config) (*Context, error) {
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

	pvkey, err := keys.PVKeyFromTendermint(&filePV.Key)
	if err != nil {
		return nil, err
	}

	hdwalletFile := strings.Replace(consensus.PrivValidatorKeyFilename, ".json", "hdkeychain.json", 1)

	keywordsBytes, err := ioutil.ReadFile(hdwalletFile)
	if err != nil {
		return nil, err
	}

	keywords := strings.Split(string(keywordsBytes), " ")

	return &Context{
		privateKey: priv,
		privval:    pvkey,
		keywords:   keywords,
	}, nil
}
