package node

import (
	"testing"

	"github.com/Oneledger/protocol/data/keys"

	"github.com/stretchr/testify/assert"

	"github.com/Oneledger/protocol/config"
)

func setup() *config.Server {
	networkConfig := &config.NetworkConfig{
		P2PAddress: "tcp://127.0.0.1:26601",
		RPCAddress: "tcp://127.0.0.1:26600",
		SDKAddress: "http://127.0.0.1:26603",
	}
	consensus := &config.ConsensusConfig{}
	p2pConfig := &config.P2PConfig{}
	mempool := &config.MempoolConfig{}
	nodeConfig := &config.NodeConfig{
		NodeName: "test_node",
		FastSync: true,
		DBDir:    "test_dbpath",
		DB:       "goleveldb",
	}
	config := &config.Server{
		Node:      nodeConfig,
		Network:   networkConfig,
		Consensus: consensus,
		P2P:       p2pConfig,
		Mempool:   mempool,
	}
	return config
}

func TestNewNodeContext(t *testing.T) {
	t.Run("should return error when given nil as config input", func(t *testing.T) {
		_, err := NewNodeContext(nil)
		assert.Error(t, err)
	})
	t.Run("should return a valid context given a valid config", func(t *testing.T) {
		cfg := setup()
		_, err := NewNodeContext(cfg)
		assert.Error(t, err)
	})
}

func TestContext_PrivVal(t *testing.T) {
	_, privateKey, err := keys.NewKeyPairFromTendermint()
	assert.Nil(t, err)
	ctx := Context{
		NodeName:   "test_node",
		privateKey: privateKey,
		privval:    privateKey,
	}

	t.Run("PrivVal test : returned private key should be equal to the generated one", func(t *testing.T) {
		prikey := ctx.PrivVal()
		assert.Equal(t, privateKey, prikey)
	})
	t.Run("PrivKey test : returned private key should be equal to the generated one", func(t *testing.T) {
		prikey := ctx.PrivKey()
		assert.Equal(t, privateKey, prikey)
	})
}

func TestContext_PubKey(t *testing.T) {
	t.Run("should return an empty keys.publicKey when given a invalid context", func(t *testing.T) {
		ctx := Context{
			NodeName: "test_node",
		}
		returnedPubKey := ctx.PubKey()
		assert.Equal(t, keys.PublicKey{}, returnedPubKey)
	})
	t.Run("should return the generated public key when given a valid private key", func(t *testing.T) {
		publicKey, privateKey, err := keys.NewKeyPairFromTendermint()
		assert.Nil(t, err)
		ctx := Context{
			NodeName:   "test_node",
			privateKey: privateKey,
			privval:    privateKey,
		}
		returnedPubKey := ctx.PubKey()
		assert.NotEqual(t, keys.PublicKey{}, returnedPubKey)
		assert.Equal(t, publicKey, returnedPubKey)
	})
}

func TestContext_ValidatorPubKey(t *testing.T) {
	t.Run("should return an empty keys.publicKey when given a invalid context", func(t *testing.T) {
		ctx := Context{
			NodeName: "test_node",
		}
		returnedPubKey := ctx.ValidatorPubKey()
		assert.Equal(t, keys.PublicKey{}, returnedPubKey)
	})
	t.Run("should return the generated public key when given a valid private key", func(t *testing.T) {
		publicKey, privateKey, err := keys.NewKeyPairFromTendermint()
		assert.Nil(t, err)
		ctx := Context{
			NodeName:   "test_node",
			privateKey: privateKey,
			privval:    privateKey,
		}
		returnedPubKey := ctx.ValidatorPubKey()
		assert.NotEqual(t, keys.PublicKey{}, returnedPubKey)
		assert.Equal(t, publicKey, returnedPubKey)
	})
}

func TestContext_Address(t *testing.T) {
	t.Run("should return nil when given an invalid context", func(t *testing.T) {
		ctx := Context{
			NodeName: "test_node",
		}
		returnedAddress := ctx.Address()
		assert.Nil(t, returnedAddress)
	})
	t.Run("should return valid address when given a valid context", func(t *testing.T) {
		_, privateKey, err := keys.NewKeyPairFromTendermint()
		assert.Nil(t, err)
		ctx := Context{
			NodeName:   "test_node",
			privateKey: privateKey,
			privval:    privateKey,
		}
		returnedAddress := ctx.Address()
		assert.IsType(t, keys.Address{}, returnedAddress)
		assert.NotEmpty(t, returnedAddress)
	})
}

func TestContext_ValidatorAddress(t *testing.T) {
	t.Run("should return nil when given an invalid context", func(t *testing.T) {
		ctx := Context{
			NodeName: "test_node",
		}
		returnedAddress := ctx.ValidatorAddress()
		assert.Nil(t, returnedAddress)
	})
	t.Run("should return valid address when given a valid context", func(t *testing.T) {
		_, privateKey, err := keys.NewKeyPairFromTendermint()
		assert.Nil(t, err)
		ctx := Context{
			NodeName:   "test_node",
			privateKey: privateKey,
			privval:    privateKey,
		}
		returnedAddress := ctx.ValidatorAddress()
		assert.IsType(t, keys.Address{}, returnedAddress)
		assert.NotEmpty(t, returnedAddress)
	})
}
