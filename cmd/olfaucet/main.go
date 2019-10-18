package main

import (
	"fmt"
	"github.com/Oneledger/protocol/data/chain"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Oneledger/protocol/version"

	"github.com/Oneledger/protocol/data/balance"

	"github.com/Oneledger/protocol/action"
	"github.com/Oneledger/protocol/app/node"
	"github.com/Oneledger/protocol/client"
	"github.com/Oneledger/protocol/config"
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/log"
	"github.com/Oneledger/protocol/rpc"
	"github.com/pkg/errors"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
)

var logger = log.NewLoggerWithPrefix(os.Stdout, "faucet")
var args faucetArgs

type faucetArgs struct {
	rootDir         string
	fullnodeConnStr string
	listenAddr      string
	lockTime        int
	maxReqAmount    int
	dbDir           string
}

func (f faucetArgs) MaxLockTime() time.Duration {
	return time.Minute * time.Duration(f.lockTime)
}

var faucetCmd = &cobra.Command{
	Use:   "olfaucet",
	Short: "a faucet service for testnets",
	Long:  "a faucet service for testnets",
	RunE:  runFaucet,
}

var apiRoutes = make(map[string]http.HandlerFunc)

func init() {
	faucetCmd.Flags().StringVarP(&args.rootDir, "root", "r", "./", "Set root directory containing olfullnode data files")
	faucetCmd.Flags().StringVar(&args.dbDir, "db_dir", "./", "Set directory of keyvalue store")
	faucetCmd.Flags().StringVar(&args.fullnodeConnStr, "fullnode_connection", "", "Specify connection string for fullnode. If not present, uses the root directory's configuration to make the connection")
	faucetCmd.Flags().StringVarP(&args.listenAddr, "listen_address", "l", "http://0.0.0.0:65432", "Specify what address to listen on")
	faucetCmd.Flags().IntVar(&args.lockTime, "lock_time", 60, "Max lock time for individual accounts, specified in terms of minutes")
	faucetCmd.Flags().IntVar(&args.maxReqAmount, "max_req_amount", 50, "Maximum number of OLT to request per account")

	addRestfulAPIEndpoint()
}

// addRestfulAPIEndpoint collects all restful API router and function mapping
// update this function to extend more restful API calls
func addRestfulAPIEndpoint() {
	apiRoutes["/"] = restfulAPIRoot
	apiRoutes["/health"] = health
}

func main() {
	err := faucetCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func runFaucet(_ *cobra.Command, _ []string) error {
	cfg := new(config.Server)
	err := cfg.ReadFile(filepath.Join(args.rootDir, config.FileName))
	if err != nil {
		logger.Error("failed to read configuration file")
		return errors.Wrap(err, "failed to read configuration file")
	}

	if args.fullnodeConnStr == "" {
		args.fullnodeConnStr = cfg.Network.SDKAddress
	}

	faucet, err := NewFaucet(cfg)
	if err != nil {
		return err
	}
	defer faucet.db.Close()

	srv := rpc.NewServer(os.Stdout)

	logger.Info(args.listenAddr)
	u, err := url.Parse(args.listenAddr)
	if err != nil {
		return errors.New("Invalid listen address")
	}

	err = srv.Prepare(u)
	if err != nil {
		return err
	}

	err = srv.Register("faucet", faucet)
	if err != nil {
		return err
	}

	srv.RegisterRestfulMap(apiRoutes)

	err = srv.Start()
	if err != nil {
		return err
	}

	HandleSigTerm(func() {
		faucet.db.Close()
		srv.Close()
	})

	select {}
}

type Faucet struct {
	nodeCtx  *node.Context
	fullnode *client.ServiceClient
	db       *leveldb.DB
}

type ParamsReply struct {
	MaxAmount   int    `json:"maxAmount"`
	MinWaitTime int    `json:"minWaitTimeMinutes"`
	Version     string `json:"version"`
}

type Request struct {
	Address keys.Address `json:"address"`
	Amount  int          `json:"amount"`
	jsonrpc2.Ctx
}

func (req *Request) HTTPRequest() *http.Request {
	return jsonrpc2.HTTPRequestFromContext(req.Ctx.Context())
}

type Reply struct {
	OK   bool          `json:"ok"`
	Sent action.Amount `json:"sent"`
}

// Return the time this IP last made a request
func (f *Faucet) get(key string) (time.Time, bool) {
	val, err := f.db.Get([]byte(key), nil)
	if err != nil {
		return time.Time{}, false
	}
	t := new(time.Time)
	err = t.UnmarshalText(val)
	if err != nil {
		logger.Error("failed to unmarshal stored time", val)
		return time.Time{}, false
	}

	return *t, true
}

// Mark a time this was requested
func (f *Faucet) set(key string, value time.Time) bool {
	t, err := value.MarshalText()
	if err != nil {
		logger.Error("set err", err)
		return false
	}
	err = f.db.Put([]byte(key), t, nil)
	if err != nil {
		logger.Error("set err", err)
		return false
	}

	return true
}

func timeSinceLastReq(f *Faucet, req Request, now time.Time) time.Duration {
	//host, _, _ := net.SplitHostPort(req.HTTPRequest().RemoteAddr)
	t2, ok := f.get(req.Address.Humanize())
	if !ok {
		return args.MaxLockTime() + 50000
	}

	return now.Sub(t2)
}

func ip(rawAddr string) string {
	host, _, _ := net.SplitHostPort(rawAddr)
	return host
}

func (f *Faucet) RequestOLT(req Request, reply *Reply) error {
	requestIP := ip(req.HTTPRequest().RemoteAddr)
	logger.Info("Incoming request from", requestIP)
	olt := balance.Currency{
		Id:      0,
		Name:    "OLT",
		Chain:   chain.Type(0),
		Decimal: 18,
		Unit:    "nue",
	}

	err := req.Address.Err()
	if err != nil {
		return rpc.InvalidParamsError(err.Error())
	}
	logger.Infof("Request \n\t Address: %s\n\t Amount: %d", req.Address.String(), req.Amount)

	sinceLastReqTime := timeSinceLastReq(f, req, time.Now())
	if sinceLastReqTime < args.MaxLockTime() {
		waitTimeMins := args.MaxLockTime().Minutes() - sinceLastReqTime.Minutes()
		waitTime := fmt.Sprintf("%.1f minutes", waitTimeMins)
		return rpc.NotAllowedError("This address is locked from making requests for " + waitTime)
	}

	if req.Amount <= 0 {
		return rpc.NotAllowedError("Request needs an amount greater than 0")
	} else if req.Amount > args.maxReqAmount {
		req.Amount = args.maxReqAmount
	}

	if req.Address == nil {
		return rpc.InvalidRequestError("address must not be nil")
	}
	amt := olt.NewCoinFromInt(int64(req.Amount))
	toSend := action.Amount{
		Currency: olt.Name,
		Value:    *amt.Amount,
	}

	sendTxResults, err := f.fullnode.CreateRawSend(client.SendTxRequest{
		From:   f.nodeCtx.Address(),
		To:     req.Address,
		Amount: toSend,
		Fee:    action.Amount{Currency: olt.Name, Value: *balance.NewAmount(1000000000)},
		Gas:    40000,
	})
	if err != nil {
		logger.Error("failed to sendTx", err)
		return rpc.InternalError(err.Error())
	}

	rawTx := sendTxResults.RawTx
	h, err := f.nodeCtx.PrivKey().GetHandler()
	if err != nil {
		panic("invalid nodeCtx private key")
	}

	sig, err := h.Sign(rawTx)
	if err != nil {
		panic("failed to sign raw tx " + err.Error())
	}

	broadcastResult, err := f.fullnode.TxSync(client.BroadcastRequest{
		RawTx:     sendTxResults.RawTx,
		Signature: sig,
		PublicKey: f.nodeCtx.PubKey(),
	})
	if err != nil {
		logger.Error("failed to sendTx", err)

		return rpc.InternalError(err.Error())
	}

	*reply = Reply{
		OK:   broadcastResult.OK,
		Sent: toSend,
	}

	// Set the time this request was made
	f.set(req.Address.Humanize(), time.Now())
	return nil
}

func (f *Faucet) GetParams(_ struct{}, reply *ParamsReply) error {
	*reply = ParamsReply{MaxAmount: args.maxReqAmount, MinWaitTime: args.lockTime, Version: version.Fullnode.String()}
	return nil
}

// restful API functions
func restfulAPIRoot(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "Available endpoints: ")
	if err != nil {
		logger.Error("failed to display available endpoints info")
	}
	for path := range apiRoutes {
		_, err = fmt.Fprintln(w, r.Host+path)
		if err != nil {
			logger.Error("failed to display available endpoints info")
		}
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	healthCheck := ParamsReply{MaxAmount: args.maxReqAmount, MinWaitTime: args.lockTime, Version: version.Fullnode.String()}
	_, err := fmt.Fprintf(w, "MaxAmount : %v, MinWaitTime : %d, version : %v\n", healthCheck.MaxAmount, healthCheck.MinWaitTime, healthCheck.Version)
	if err != nil {
		logger.Error("failed to display health check info")
	}
}

func NewFaucet(cfg *config.Server) (*Faucet, error) {
	// Use the node's account
	nodeCtx, err := node.NewNodeContext(cfg)
	if err != nil {
		logger.Error("failed to load nodeContext from configuration")
		return nil, err
	}

	logger.Info("Loaded account", nodeCtx.Address().String())

	dbPath := filepath.Join(args.dbDir, "faucet.db")

	// Create the database
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new db")
	}

	// Create the fullnodeclient
	fullnode, err := client.NewServiceClient(args.fullnodeConnStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start fullnode client")
	}

	addr, err := fullnode.NodeAddress()
	if err != nil {
		return nil, errors.Wrap(err, "failed to start node conn")
	}

	logger.Info("Connected to node run by", addr.String())

	balReply, err := fullnode.Balance(nodeCtx.Address())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get my balance")
	}

	logger.Info("Faucet balance:", balReply.Balance)

	return &Faucet{
		fullnode: fullnode,
		db:       db,
		nodeCtx:  nodeCtx,
	}, nil
}

func HandleSigTerm(cancel func()) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-signals
		logger.Info("Got quit signal...", sig)
		cancel()
		os.Exit(1)
	}()

}
