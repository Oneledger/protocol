package config

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Oneledger/protocol/node/cmd/shared"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/log"
	"github.com/spf13/viper"
)

// Going to abandon servers
// Configure Server sets viper to use the specified configuration filename
func ConfigureServer() {
	viper.SetConfigName(global.Current.ConfigName)

	viper.AddConfigPath(global.Current.RootDir) // Special user overrides
	viper.AddConfigPath(".")                    // Local directory override

	err := viper.ReadInConfig()
	if err != nil {
		log.Info("Not using config file", "err", err)
	}
}

// Struct for holding the configuration details for the node
type ServerConfig struct {
	Node      *NodeConfig      `toml:"node"`
	Network   *NetworkConfig   `toml:"network"`
	P2P       *P2PConfig       `toml:"p2p"`
	Mempool   *MempoolConfig   `toml:"mempool"`
	Consensus *ConsensusConfig `toml:"consensus"`
}

func (sc *ServerConfig) ReadFile(filepath string) error {
	bz, err := shared.ReadFile(filepath)
	if err != nil {
		return err
	}
	return sc.Unmarshal(bz)
}

func (sc *ServerConfig) Unmarshal(text []byte) error {
	_, err := toml.Decode(string(text), sc)
	if err != nil {
		return err
	}
	return nil
}

func (sc *ServerConfig) Marshal() (text []byte, err error) {
	var buf bytes.Buffer
	err = toml.NewEncoder(&buf).Encode(sc)
	text = buf.Bytes()
	return
}

func (sc *ServerConfig) SaveFile(filepath string) error {
	bz, err := sc.Marshal()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath, bz, 0644)
}

func (sc *ServerConfig) OpenFile(filepath string) error {
	bz, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return sc.Unmarshal(bz)
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Node:      DefaultNodeConfig(),
		Network:   DefaultNetworkConfig(),
		P2P:       DefaultP2PConfig(),
		Mempool:   DefaultMempoolConfig(),
		Consensus: DefaultConsensusConfig(),
	}
}

// NodeConfig handles general configuration settings for the node
type NodeConfig struct {
	NodeName string `toml:"node_name"`
	FastSync bool   `toml:"fast_sync"`
	// Specify what backend database to use: (goleveldb|cleveldb)
	// It is recommended to use cleveldb for production environments
	DB string `toml:"db"`
	// List of transaction tags to index in the db, allows them to be searched
	// by this parameter
	IndexTags []string `toml:"index_tags"`
	// Tells the indexer to index all available tags, IndexTags has precedence
	// over IndexAllTAgs
	IndexAllTags bool `toml:"index_all_tags"`
}

func DefaultNodeConfig() *NodeConfig {
	return &NodeConfig{
		NodeName:     "Newton-Node",
		FastSync:     true,
		DB:           "goleveldb",
		IndexTags:    []string{"foo", "bar", "baz"},
		IndexAllTags: false,
	}
}

// NetworkConfig exposes configuration files for the current
type NetworkConfig struct {
	RPCAddress string `toml:"rpc_address"`
	// Specify an address and port for incoming clients to connect to
	P2PAddress string `toml:"p2p_address"`

	// Address to advertise for incoming peers to connect
	ExternalP2PAddress string `toml:"external_p2p_address"`

	SDKAddress string `toml:"sdk_address"`

	// Point to a bitcoin node, can be empty
	BTCAddress string `toml:"btc_address"`

	// Point to an ethereum node, can be empty
	ETHAddress string `toml:"eth_address"`

	OLVMAddress  string `toml:"olvm_address"`
	OLVMProtocol string `toml:"olvm_protocol"`
}

func DefaultNetworkConfig() *NetworkConfig {
	// TODO: Fill in default values
	return &NetworkConfig{}
}

// P2PConfig defines the options for P2P networking layer
type P2PConfig struct {
	// List of seed nodes to connect to
	Seeds []string `toml:"seeds"`

	// Enables seed mode, this node will constantly crawl the network looking for peers
	SeedMode bool `toml:"seed_mode"`

	// List of nodes to keep persistent connections to
	PersistentPeers []string `toml:"persistent_peers"`

	// Enable UPNP port forwarding
	UPNP bool `toml:"upnp"`

	// Set true for strict address routability rules
	// If true, the node will fail to start if the given P2P address isn't routable
	AddrBookStrict bool `toml:"addr_book_strict"`

	// Maximum number of inbound peers
	MaxNumInboundPeers int `toml:"max_num_inbound_peers"`

	// Maximum number of outbound peers to connect to, excluding persistent peers
	MaxNumOutboundPeers int `toml:"max_num_outbound_peers"`

	// Time to wait before flushing messages out on the connection
	FlushThrottleTimeout time.Duration `toml:"flush_throttle_timeout"`

	// Maximum size of a message packet payload, in bytes
	MaxPacketMsgPayloadSize int `toml:"max_packet_msg_payload_size"`

	// Rate at which packets can be sent, in bytes/second
	SendRate int64 `toml:"send_rate"`

	// Rate at which packets can be received, in bytes/second
	RecvRate int64 `toml:"recv_rate"`

	// Set true to enable the peer-exchange reactor
	PexReactor bool `toml:"pex"`

	// Comma separated list of peer IDs to keep private (will not be gossiped to
	// other peers)
	PrivatePeerIDs string `toml:"private_peer_ids"`

	// Toggle to disable guard against peers connecting from the same ip.
	AllowDuplicateIP bool `toml:"allow_duplicate_ip"`

	// Peer connection configuration.
	HandshakeTimeout time.Duration `toml:"handshake_timeout"`
	DialTimeout      time.Duration `toml:"dial_timeout"`
}

func DefaultP2PConfig() *P2PConfig {
	// TODO: Fill in default values
	return &P2PConfig{}
}

// MempoolConfig defines configuration options for the mempool
type MempoolConfig struct {
	Recheck   bool `toml:"recheck"`
	Broadcast bool `toml:"broadcast"`
	Size      int  `toml:"size"`
	CacheSize int  `toml:"cache_size"`
}

func DefaultMempoolConfig() *MempoolConfig {
	return &MempoolConfig{}
}

// ConsensusConfig handles consensus-specific options
type ConsensusConfig struct {
	TimeoutPropose        time.Duration `toml:"timeout_propose"`
	TimeoutProposeDelta   time.Duration `toml:"timeout_propose_delta"`
	TimeoutPrevote        time.Duration `toml:"timeout_prevote"`
	TimeoutPrevoteDelta   time.Duration `toml:"timeout_prevote_delta"`
	TimeoutPrecommit      time.Duration `toml:"timeout_precommit"`
	TimeoutPrecommitDelta time.Duration `toml:"timeout_precommit_delta"`
	TimeoutCommit         time.Duration `toml:"timeout_commit"`

	// Make progress as soon as we have all the precommits (as if TimeoutCommit = 0)
	SkipTimeoutCommit bool `toml:"skip_timeout_commit"`

	// EmptyBlocks mode and possible interval between empty blocks
	CreateEmptyBlocks         bool          `toml:"create_empty_blocks"`
	CreateEmptyBlocksInterval time.Duration `toml:"create_empty_blocks_interval"`

	// Reactor sleep duration parameters
	PeerGossipSleepDuration     time.Duration `toml:"peer_gossip_sleep_duration"`
	PeerQueryMaj23SleepDuration time.Duration `toml:"peer_query_maj23_sleep_duration"`

	// Block time parameters. Corresponds to the minimum time increment between consecutive blocks.
	BlockTimeIota time.Duration `toml:"blocktime_iota"`
}

func DefaultConsensusConfig() *ConsensusConfig {
	return &ConsensusConfig{}
}
