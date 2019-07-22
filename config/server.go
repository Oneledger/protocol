package config

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/Oneledger/toml"
	"github.com/pkg/errors"
	tmconfig "github.com/tendermint/tendermint/config"
)

const (
	// Default permissions for writing files
	// These are based on default umask settings
	// User+Group: rw, Other: r
	FilePerms = 0664
	// User+Group: rwx, Other: rx
	DirPerms   = 0775
	FileName   = "config.toml"
	DefaultDir = ".olfullnode"
)

// Duration is a time.Duration that marshals and unmarshals with millisecond values
type Duration int64

// Returns a nanosecond duration
func (d Duration) Nanoseconds() time.Duration {
	return time.Duration(d) * time.Millisecond
}

func toConfigDuration(d time.Duration) Duration {
	return Duration(d / time.Millisecond)
}

// Struct for holding the configuration details for the node
type Server struct {
	Node      *NodeConfig      `toml:"node"`
	Network   *NetworkConfig   `toml:"network"`
	P2P       *P2PConfig       `toml:"p2p"`
	Mempool   *MempoolConfig   `toml:"mempool"`
	Consensus *ConsensusConfig `toml:"consensus"`

	chainID string
	rootDir string
}

func (cfg *Server) RootDir() string {
	return cfg.rootDir
}

func (cfg *Server) ChainID() string {
	return cfg.chainID
}

func (cfg *Server) setChainID(doc GenesisDoc) {
	cfg.chainID = doc.ChainID
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

	csConfig := cfg.Consensus.TMConfig()
	csConfig.WalPath = filepath.Join(baseConfig.DBPath, "cs.wal", "wal")

	rpcConfig := tmconfig.DefaultRPCConfig()
	rpcConfig.ListenAddress = cfg.Network.RPCAddress

	nilMetricsConfig := tmconfig.InstrumentationConfig{Namespace: "metrics"}

	tmcfg := &tmconfig.Config{
		BaseConfig: baseConfig,
		RPC:        rpcConfig,
		P2P:        p2pConfig,
		Mempool:    cfg.Mempool.TMConfig(),
		Consensus:  csConfig,
		TxIndex: &tmconfig.TxIndexConfig{
			Indexer:      "kv",
			IndexTags:    strings.Join(cfg.Node.IndexTags, ","),
			IndexAllTags: cfg.Node.IndexAllTags,
		},
		Instrumentation: &nilMetricsConfig,
	}

	tmcfg.SetRoot(filepath.Join(cfg.rootDir, "consensus"))

	return *tmcfg
}

// ReadFile accepts a filepath and returns the
func (cfg *Server) ReadFile(path string) error {
	bz, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "cfg.ReadFile error")
	}
	err = cfg.Unmarshal(bz)
	if err != nil {
		return errors.Wrap(err, "cfg.ReadFile error unmarshaling JSON")
	}

	// Set internal root directory
	cfg.rootDir = filepath.Dir(path)

	return nil
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
	FastSync bool   `toml:"fast_sync" desc:"Fast sync allows a block to catch up quickly to the chain by downloading blocks in parallel and verifying their commits"`
	DB       string `toml:"db" desc:"Specify what backend database to use (goleveldb|cleveldb)"`
	DBDir    string `toml:"db_dir" desc:"Specify the application database directory. This is always relative to the root directory of the app."`
	// List of transaction tags to index in the db, allows them to be searched
	// by this parameter
	IndexTags []string `toml:"index_tags" desc:"List of transaction tags to index in the database, allows them to be searched by the specified tags"`
	// Tells the indexer to index all available tags, IndexTags has precedence
	// over IndexAllTAgs
	IndexAllTags bool `toml:"index_all_tags" desc:"Tells the indexer to index all available tags, IndexTags has precedence over IndexAllTags"`
}

func DefaultNodeConfig() *NodeConfig {
	return &NodeConfig{
		NodeName:     "Newton-Node",
		FastSync:     true,
		DB:           "goleveldb",
		DBDir:        "nodedata",
		IndexTags:    []string{"tx.owner", "tx.type", "tx.swapkey"},
		IndexAllTags: false,
	}
}

// NetworkConfig exposes configuration files for the current
type NetworkConfig struct {
	RPCAddress string `toml:"rpc_address"`
	P2PAddress string `toml:"p2p_address" desc:"Main address for P2P connections"`

	ExternalP2PAddress string `toml:"external_p2p_address" desc:"Address to advertise for incoming peers to connect to"`

	SDKAddress string `toml:"sdk_address"`

	BTCAddress string `toml:"btc_address"`
	ETHAddress string `toml:"eth_address"`

	OLVMAddress  string `toml:"olvm_address"`
	OLVMProtocol string `toml:"olvm_protocol"`
}

func DefaultNetworkConfig() *NetworkConfig {
	return &NetworkConfig{
		RPCAddress:         "http://127.0.0.1:26601",
		P2PAddress:         "tcp://127.0.0.1:26611",
		ExternalP2PAddress: "",
		SDKAddress:         "tcp://127.0.0.1:26631",
		OLVMAddress:        "tcp://127.0.0.1:26641",
		OLVMProtocol:       "tcp",
		BTCAddress:         "tcp://127.0.0.1:NONE",
		ETHAddress:         "NONE",
	}
}

// P2PConfig defines the options for P2P networking layer
type P2PConfig struct {
	Seeds []string `toml:"seeds" desc:"List of seed nodes to connect to"`

	SeedMode bool `toml:"seed_mode" desc:"Enables seed mode, which will make the node crawl the network looking for peers"`

	PersistentPeers []string `toml:"persistent_peers" desc:"List of peers to maintain a persistent connection to"`

	UPNP bool `toml:"upnp" desc:"Enable UPNP port forwarding"`

	AddrBookStrict bool `toml:"addr_book_strict" desc:"Set true for strict address routability rules. If true, the node will fail to start if the given P2P address isn't routable'"`

	MaxNumInboundPeers int `toml:"max_num_inbound_peers" desc:"Max number of inbound peers"`

	MaxNumOutboundPeers int `toml:"max_num_outbound_peers" desc:"Max number of outbound peers to connect to, excluding persistent peers"`

	FlushThrottleTimeout Duration `toml:"flush_throttle_timeout" desc:"Time to wait before flushing messages out on the connection in milliseconds"`

	MaxPacketMsgPayloadSize int `toml:"max_packet_msg_payload_size" desc:"Max size of a message packet payload, in bytes"`

	SendRate int64 `toml:"send_rate" desc:"Rate at which packets can be sent, in bytes/second"`

	// Rate at which packets can be received, in bytes/second
	RecvRate int64 `toml:"recv_rate" desc:"Rate at which packets can be received, in bytes/second"`

	PexReactor bool `toml:"pex" desc:"Set true to enable the peer-exchange reactor"`

	PrivatePeerIDs []string `toml:"private_peer_ids" desc:"List of peer IDs to keep private (will not be gossiped to other peers)"`

	AllowDuplicateIP bool `toml:"allow_duplicate_ip" desc:"Toggle to disable guard against peers connecting from the same IP"`

	HandshakeTimeout Duration `toml:"handshake_timeout" desc:"In milliseconds"`
	DialTimeout      Duration `toml:"dial_timeout" desc:"In milliseconds"`
}

func (cfg *P2PConfig) TMConfig() *tmconfig.P2PConfig {
	return &tmconfig.P2PConfig{
		Seeds:                   strings.Join(cfg.Seeds, ","),
		PersistentPeers:         strings.Join(cfg.PersistentPeers, ","),
		UPNP:                    cfg.UPNP,
		AddrBook:                filepath.Join("config", "addrbook.json"),
		AddrBookStrict:          cfg.AddrBookStrict,
		MaxNumInboundPeers:      cfg.MaxNumInboundPeers,
		MaxNumOutboundPeers:     cfg.MaxNumOutboundPeers,
		MaxPacketMsgPayloadSize: cfg.MaxPacketMsgPayloadSize,
		SendRate:                cfg.SendRate,
		RecvRate:                cfg.RecvRate,
		PexReactor:              cfg.PexReactor,
		SeedMode:                cfg.SeedMode,
		PrivatePeerIDs:          strings.Join(cfg.PrivatePeerIDs, ","),
		AllowDuplicateIP:        cfg.AllowDuplicateIP,
		FlushThrottleTimeout:    cfg.FlushThrottleTimeout.Nanoseconds(),
		HandshakeTimeout:        cfg.HandshakeTimeout.Nanoseconds(),
		DialTimeout:             cfg.DialTimeout.Nanoseconds(),
	}
}

func DefaultP2PConfig() *P2PConfig {
	var cfg P2PConfig
	tmDefaults := tmconfig.DefaultP2PConfig()
	cfg.Seeds = make([]string, 0)
	cfg.PersistentPeers = make([]string, 0)
	cfg.PrivatePeerIDs = make([]string, 0)
	cfg.UPNP = tmDefaults.UPNP
	cfg.AddrBookStrict = false
	cfg.AllowDuplicateIP = true
	cfg.MaxNumInboundPeers = tmDefaults.MaxNumInboundPeers
	cfg.MaxNumOutboundPeers = tmDefaults.MaxNumOutboundPeers
	cfg.MaxPacketMsgPayloadSize = tmDefaults.MaxPacketMsgPayloadSize
	cfg.SendRate = tmDefaults.SendRate
	cfg.RecvRate = tmDefaults.RecvRate
	cfg.PexReactor = tmDefaults.PexReactor
	cfg.SeedMode = tmDefaults.SeedMode
	cfg.FlushThrottleTimeout = toConfigDuration(tmDefaults.FlushThrottleTimeout)
	cfg.HandshakeTimeout = toConfigDuration(tmDefaults.HandshakeTimeout)
	cfg.DialTimeout = toConfigDuration(tmDefaults.DialTimeout)
	return &cfg
}

// MempoolConfig defines configuration options for the mempool
type MempoolConfig struct {
	Recheck   bool `toml:"recheck"`
	Broadcast bool `toml:"broadcast"`
	Size      int  `toml:"size" desc:"Size of the mempool"`
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
	LogOutput             string   `toml:"log_output" desc:"Determines where consensus is logged (stdout|<filename>)"`
	LogLevel              string   `toml:"log_level" desc:"Determines the verbosity of consensus logs"`
	TimeoutPropose        Duration `toml:"timeout_propose" desc:"All timeouts are in milliseconds"`
	TimeoutProposeDelta   Duration `toml:"timeout_propose_delta"`
	TimeoutPrevote        Duration `toml:"timeout_prevote"`
	TimeoutPrevoteDelta   Duration `toml:"timeout_prevote_delta"`
	TimeoutPrecommit      Duration `toml:"timeout_precommit"`
	TimeoutPrecommitDelta Duration `toml:"timeout_precommit_delta"`
	TimeoutCommit         Duration `toml:"timeout_commit"`

	SkipTimeoutCommit bool `toml:"skip_timeout_commit" desc:"Make progress as soon as we have all precommits (as if TimeoutCommit = 0)"`

	CreateEmptyBlocks           bool     `toml:"create_empty_blocks" desc:"Should this node create empty blocks"`
	CreateEmptyBlocksInterval   Duration `toml:"create_empty_blocks_interval" desc:"Interval between empty block creation in milliseconds"`
	PeerGossipSleepDuration     Duration `toml:"peer_gossip_sleep_duration" desc:"Duration values in milliseconds"`
	PeerQueryMaj23SleepDuration Duration `toml:"peer_query_maj23_sleep_duration"`
	BlockTimeIota               Duration `toml:"blocktime_iota" desc:"Block time parameter, corresponds to the minimum time increment between consecutive blocks"`
}

func (cfg *ConsensusConfig) TMConfig() *tmconfig.ConsensusConfig {
	c := tmconfig.DefaultConsensusConfig()
	c.TimeoutPropose = cfg.TimeoutPropose.Nanoseconds()
	c.TimeoutProposeDelta = cfg.TimeoutProposeDelta.Nanoseconds()
	c.TimeoutPrevote = cfg.TimeoutPrevote.Nanoseconds()
	c.TimeoutPrevoteDelta = cfg.TimeoutPrevoteDelta.Nanoseconds()
	c.TimeoutPrecommit = cfg.TimeoutPrecommit.Nanoseconds()
	c.TimeoutPrecommitDelta = cfg.TimeoutPrecommitDelta.Nanoseconds()
	c.TimeoutCommit = cfg.TimeoutCommit.Nanoseconds()
	c.SkipTimeoutCommit = cfg.SkipTimeoutCommit
	c.CreateEmptyBlocks = cfg.CreateEmptyBlocks
	c.CreateEmptyBlocksInterval = cfg.CreateEmptyBlocksInterval.Nanoseconds()
	c.PeerGossipSleepDuration = cfg.PeerGossipSleepDuration.Nanoseconds()
	c.PeerQueryMaj23SleepDuration = cfg.PeerQueryMaj23SleepDuration.Nanoseconds()
	return c
}

func DefaultConsensusConfig() *ConsensusConfig {
	var cfg ConsensusConfig
	tmDefault := tmconfig.DefaultConsensusConfig()
	cfg.LogOutput = "consensus.log"
	cfg.LogLevel = tmconfig.DefaultPackageLogLevels()
	cfg.TimeoutPropose = toConfigDuration(tmDefault.TimeoutPropose)
	cfg.TimeoutProposeDelta = toConfigDuration(tmDefault.TimeoutProposeDelta)
	cfg.TimeoutPrevote = toConfigDuration(tmDefault.TimeoutPrevote)
	cfg.TimeoutPrevoteDelta = toConfigDuration(tmDefault.TimeoutPrevoteDelta)
	cfg.TimeoutPrecommit = toConfigDuration(tmDefault.TimeoutPrecommit)
	cfg.TimeoutPrecommitDelta = toConfigDuration(tmDefault.TimeoutPrecommitDelta)
	cfg.TimeoutCommit = toConfigDuration(tmDefault.TimeoutCommit)
	cfg.SkipTimeoutCommit = tmDefault.SkipTimeoutCommit
	cfg.CreateEmptyBlocks = tmDefault.CreateEmptyBlocks
	cfg.CreateEmptyBlocksInterval = toConfigDuration(tmDefault.CreateEmptyBlocksInterval)
	cfg.PeerGossipSleepDuration = toConfigDuration(tmDefault.PeerGossipSleepDuration)
	cfg.PeerQueryMaj23SleepDuration = toConfigDuration(tmDefault.PeerQueryMaj23SleepDuration)
	return &cfg
}
