package config

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/jbsmith7741/toml"
	tmconfig "github.com/tendermint/tendermint/config"
)

// Default permissions for writing files
// These are based on default umask settings
// User+Group: rw, Other: r
const FilePerms = 0664

// User+Group: rwx, Other: rx
const DirPerms = 0775
const FileName = "config.toml"
const DefaultDir = ".olfullnode"

// Struct for holding the configuration details for the node
type Server struct {
	Node      *NodeConfig      `toml:"node"`
	Network   *NetworkConfig   `toml:"network"`
	P2P       *P2PConfig       `toml:"p2p"`
	Mempool   *MempoolConfig   `toml:"mempool"`
	Consensus *ConsensusConfig `toml:"consensus"`
}

func (cfg *Server) TMConfig() tmconfig.Config {
	leveldb := cfg.Node.DB
	if cfg.Node.DB == "goleveldb" {
		leveldb = "leveldb"
	}

	baseConfig := tmconfig.DefaultBaseConfig()
	baseConfig.ProxyApp = "OneLedgerProtocol"
	baseConfig.Moniker = cfg.Node.NodeName
	baseConfig.FastSync = cfg.Node.FastSync
	baseConfig.DBBackend = leveldb
	baseConfig.DBPath = "data"
	baseConfig.LogLevel = cfg.Consensus.LogLevel

	p2pConfig := cfg.P2P.TMConfig()
	p2pConfig.ListenAddress = cfg.Network.P2PAddress
	p2pConfig.ExternalAddress = cfg.Network.ExternalP2PAddress
	if cfg.Network.ExternalP2PAddress == "" {
		p2pConfig.ExternalAddress = cfg.Network.P2PAddress
	}

	rpcConfig := tmconfig.DefaultRPCConfig()
	rpcConfig.ListenAddress = cfg.Network.RPCAddress

	nilMetricsConfig := tmconfig.InstrumentationConfig{Namespace: "metrics"}

	return tmconfig.Config{
		BaseConfig: baseConfig,
		RPC:        rpcConfig,
		P2P:        p2pConfig,
		Mempool:    cfg.Mempool.TMConfig(),
		Consensus:  cfg.Consensus.TMConfig(),
		TxIndex: &tmconfig.TxIndexConfig{
			Indexer:      "kv",
			IndexTags:    strings.Join(cfg.Node.IndexTags, ","),
			IndexAllTags: cfg.Node.IndexAllTags,
		},
		Instrumentation: &nilMetricsConfig,
	}
}

// ReadFile accepts a filepath and returns the
func (cfg *Server) ReadFile(filepath string) error {
	bz, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return cfg.Unmarshal(bz)
}

// Marshal accepts the text form of a TOML configuration file and fills Server with its values
func (cfg *Server) Unmarshal(text []byte) error {
	_, err := toml.Decode(string(text), cfg)
	if err != nil {
		return err
	}
	// TODO: Some basic validation for string fields, should return an error for invalid values
	return nil
}

// Marshal returns the text form of Server as TOML
func (cfg *Server) Marshal() (text []byte, err error) {
	var buf bytes.Buffer
	err = toml.NewEncoder(&buf).Encode(cfg)
	text = buf.Bytes()
	return
}

// SaveFile saves the current config to a file at the specified path
func (cfg *Server) SaveFile(filepath string) error {
	bz, err := cfg.Marshal()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath, bz, FilePerms)
}

// OpenFile opens the file at a given path and injects the
func (cfg *Server) OpenFile(filepath string) error {
	bz, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return cfg.Unmarshal(bz)
}

func DefaultServerConfig() *Server {
	return &Server{
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
		IndexTags:    []string{"tx.owner", "tx.type", "tx.swapkey"},
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
	return &NetworkConfig{
		RPCAddress:         "tcp://127.0.0.1:26601",
		P2PAddress:         "tcp://127.0.0.1:26611",
		ExternalP2PAddress: "",
		SDKAddress:         "tcp://127.0.0.1:26631",
		OLVMAddress:        "tcp://127.0.0.1:26641",
		OLVMProtocol:       "tcp",
		BTCAddress:         "127.0.0.1:NONE",
		ETHAddress:         "NONE",
	}
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

func (cfg *P2PConfig) TMConfig() *tmconfig.P2PConfig {
	return &tmconfig.P2PConfig{
		Seeds:                   strings.Join(cfg.Seeds, ","),
		PersistentPeers:         strings.Join(cfg.PersistentPeers, ","),
		UPNP:                    cfg.UPNP,
		AddrBook:                filepath.Join("consensus", "config", "addrbook.json"),
		AddrBookStrict:          cfg.AddrBookStrict,
		MaxNumInboundPeers:      cfg.MaxNumInboundPeers,
		MaxNumOutboundPeers:     cfg.MaxNumOutboundPeers,
		FlushThrottleTimeout:    cfg.FlushThrottleTimeout,
		MaxPacketMsgPayloadSize: cfg.MaxPacketMsgPayloadSize,
		SendRate:                cfg.SendRate,
		RecvRate:                cfg.RecvRate,
		PexReactor:              cfg.PexReactor,
		SeedMode:                cfg.SeedMode,
		AllowDuplicateIP:        cfg.AllowDuplicateIP,
		HandshakeTimeout:        cfg.HandshakeTimeout,
		DialTimeout:             cfg.DialTimeout,
	}
}

func DefaultP2PConfig() *P2PConfig {
	var cfg P2PConfig
	tmDefaults := tmconfig.DefaultP2PConfig()
	cfg.Seeds = make([]string, 0)
	cfg.PersistentPeers = make([]string, 0)
	cfg.UPNP = tmDefaults.UPNP
	cfg.AddrBookStrict = false
	cfg.AllowDuplicateIP = true
	cfg.MaxNumInboundPeers = tmDefaults.MaxNumInboundPeers
	cfg.MaxNumOutboundPeers = tmDefaults.MaxNumOutboundPeers
	cfg.FlushThrottleTimeout = tmDefaults.FlushThrottleTimeout
	cfg.MaxPacketMsgPayloadSize = tmDefaults.MaxPacketMsgPayloadSize
	cfg.SendRate = tmDefaults.SendRate
	cfg.RecvRate = tmDefaults.RecvRate
	cfg.PexReactor = tmDefaults.PexReactor
	cfg.SeedMode = tmDefaults.SeedMode
	cfg.HandshakeTimeout = tmDefaults.HandshakeTimeout
	cfg.DialTimeout = tmDefaults.DialTimeout
	return &cfg
}

// MempoolConfig defines configuration options for the mempool
type MempoolConfig struct {
	Recheck   bool `toml:"recheck"`
	Broadcast bool `toml:"broadcast"`
	Size      int  `toml:"size"`
	CacheSize int  `toml:"cache_size"`
}

func (cfg *MempoolConfig) TMConfig() *tmconfig.MempoolConfig {
	c := tmconfig.DefaultMempoolConfig()
	c.Recheck = cfg.Recheck
	c.Broadcast = cfg.Broadcast
	c.Size = cfg.Size
	c.CacheSize = cfg.CacheSize
	return c
}

func DefaultMempoolConfig() *MempoolConfig {
	var cfg MempoolConfig
	tmDefault := tmconfig.DefaultMempoolConfig()
	cfg.Recheck = tmDefault.Recheck
	cfg.Broadcast = tmDefault.Broadcast
	cfg.Size = tmDefault.Size
	cfg.CacheSize = tmDefault.CacheSize
	return &cfg
}

// ConsensusConfig handles consensus-specific options
type ConsensusConfig struct {
	LogOutput             string        `toml:"log_output" desc:"Determines where consensus is logged (stdout|<filename>)"`
	LogLevel              string        `toml:"log_level"`
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

func (cfg *ConsensusConfig) TMConfig() *tmconfig.ConsensusConfig {
	c := tmconfig.DefaultConsensusConfig()
	c.TimeoutPropose = cfg.TimeoutPropose
	c.TimeoutProposeDelta = cfg.TimeoutProposeDelta
	c.TimeoutPrevote = cfg.TimeoutPrevote
	c.TimeoutPrevoteDelta = cfg.TimeoutPrevoteDelta
	c.TimeoutPrecommit = cfg.TimeoutPrecommit
	c.TimeoutPrecommitDelta = cfg.TimeoutPrecommitDelta
	c.TimeoutCommit = cfg.TimeoutCommit
	c.SkipTimeoutCommit = cfg.SkipTimeoutCommit
	c.CreateEmptyBlocks = cfg.CreateEmptyBlocks
	c.CreateEmptyBlocksInterval = cfg.CreateEmptyBlocksInterval
	c.PeerGossipSleepDuration = cfg.PeerGossipSleepDuration
	c.PeerQueryMaj23SleepDuration = cfg.PeerQueryMaj23SleepDuration
	c.BlockTimeIota = cfg.BlockTimeIota
	return c
}

func DefaultConsensusConfig() *ConsensusConfig {
	var cfg ConsensusConfig
	tmDefault := tmconfig.DefaultConsensusConfig()
	cfg.LogOutput = "consensus.log"
	cfg.LogLevel = tmconfig.DefaultPackageLogLevels()
	cfg.TimeoutPropose = tmDefault.TimeoutPropose
	cfg.TimeoutProposeDelta = tmDefault.TimeoutProposeDelta
	cfg.TimeoutPrevote = tmDefault.TimeoutPrevote
	cfg.TimeoutPrevoteDelta = tmDefault.TimeoutPrevoteDelta
	cfg.TimeoutPrecommit = tmDefault.TimeoutPrecommit
	cfg.TimeoutPrecommitDelta = tmDefault.TimeoutPrecommitDelta
	cfg.TimeoutCommit = tmDefault.TimeoutCommit
	cfg.SkipTimeoutCommit = tmDefault.SkipTimeoutCommit
	cfg.CreateEmptyBlocks = tmDefault.CreateEmptyBlocks
	cfg.CreateEmptyBlocksInterval = tmDefault.CreateEmptyBlocksInterval
	cfg.PeerGossipSleepDuration = tmDefault.PeerGossipSleepDuration
	cfg.PeerQueryMaj23SleepDuration = tmDefault.PeerQueryMaj23SleepDuration
	cfg.BlockTimeIota = tmDefault.BlockTimeIota
	return &cfg
}
