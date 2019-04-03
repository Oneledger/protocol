package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Oneledger/protocol/node/serialize"
	"reflect"
	"strconv"

	"github.com/Oneledger/protocol/node/action"
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/sdk"
	"github.com/Oneledger/protocol/node/sdk/pb"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

type SDKServer struct {
	App *Application
}

// Ensure SDKServer implements pb.SDKServer
var _ pb.SDKServer = SDKServer{}

func NewSDKServer(app *Application, addr string) (*sdk.Server, error) {
	server, err := sdk.NewServer(addr, &SDKServer{app})
	if err != nil {
		log.Fatal("SDK Server failed to start", "err", err.Error())
	}
	return server, nil
}

type SDKQuery struct {
	Path      string
	Arguments map[string]interface{}
}

type SDKSet struct {
	Path      string
	Arguments map[string]interface{}
}

func init() {
	serial.Register(SDKQuery{})
	serial.Register(SDKSet{})

	serialize.RegisterConcrete(new(SDKQuery), "app_sdk_query")
	serialize.RegisterConcrete(new(SDKSet), "app_sdk_set")
}

func (server SDKServer) Request(ctx context.Context, request *pb.SDKRequest) (*pb.SDKReply, error) {
	var prototype interface{}
	var err error
	var result interface{}

	if string(request.Parameters) == "" {
		return &pb.SDKReply{Results: []byte("Empty Data")}, nil
	}

	err = clSerializer.Deserialize(request.Parameters, &prototype)
	if err != nil {
		log.Fatal("Deserialized Failed", "err", err, "data", string(request.Parameters))
	}

	switch base := prototype.(type) {
	case *SDKQuery:
		result = HandleQuery(*server.App, base.Path, base.Arguments)

	case *SDKSet:
		result = HandleSet(*server.App, base.Path, base.Arguments)

	default:
		result = HandleTypeError("Unknown Parameter")
	}

	//buffer, err := serial.Serialize(result, serial.CLIENT)
	switch base := result.(type) {
	case []byte:
		return &pb.SDKReply{
			Results: base,
		}, nil

	case string:
		return &pb.SDKReply{
			Results: []byte(base),
		}, nil
	}

	return nil, fmt.Errorf("Invalid return type: %s", reflect.TypeOf(result).String())
}

func HandleTypeError(message string) []byte {
	result, err := clSerializer.Serialize(message)
	if err != nil {
		log.Fatal("Failed to Serialize", "err", err)
	}

	return result
}

func (server SDKServer) Status(ctx context.Context, request *pb.StatusRequest) (*pb.StatusReply, error) {
	return &pb.StatusReply{Ok: server.App.SDK.IsRunning()}, nil
}

// Given the name of the identity and the chain type, generate new keys and broadcast this new identity
func (server SDKServer) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterReply, error) {
	name := request.Identity

	//fee := request.Fee
	var fee float64
	var err error
	if fee, err = strconv.ParseFloat("0.01", 64); err == nil {
		// TODO: Minimum fee
		fee = 0.01
	}

	chain := parseChainType(request.Chain)

	// First check if this account already exists, return with error if not
	_, stat := server.App.Accounts.FindNameOnChain(name, chain)
	// If this account already existsm don't let this go through
	if stat != status.MISSING_VALUE {
		return nil, gstatus.Errorf(codes.AlreadyExists, "Identity %s already exists", name)
	}

	privKey, pubKey := id.GenerateKeys(secret(name+chain.String()), true)

	packet := action.SignAndPack(
		action.CreateRegisterRequest(
			name,
			chain.String(),
			global.Current.Sequence,
			global.Current.NodeName,
			nil,
			pubKey.Address(),
			fee))
	comm.Broadcast(packet)

	// TODO: Use proper secret for key generation
	return &pb.RegisterReply{
		Ok:         true,
		Identity:   name,
		PublicKey:  pubKey.Bytes(),
		PrivateKey: privKey.Bytes(),
	}, nil
}

// CheckAccount returns the balance of a given account ID
func (server SDKServer) CheckAccount(ctx context.Context, request *pb.CheckAccountRequest) (*pb.CheckAccountReply, error) {
	accountName := request.Name
	if accountName == "" {
		return nil, gstatus.Errorf(codes.InvalidArgument, "Account name can't be empty")
	}
	account, err := server.App.Accounts.FindName(accountName)
	if err != status.SUCCESS {
		return nil, errAccountNotFound(accountName)
	}

	// Get balance
	balance := server.App.Balances.Get(data.DatabaseKey(account.AccountKey()), false)

	var result *pb.Balance
	if balance != nil {
		result = &pb.Balance{
			Amount:   balance.GetAmountByName("OLT").Float64(),
			Currency: currencyProtobuf(balance.GetAmountByName("OLT").Currency),
		}
	} else {
		result = &pb.Balance{Amount: 0, Currency: pb.Currency_OLT}
	}

	return &pb.CheckAccountReply{
		Name:       account.Name(),
		Chain:      account.Chain().String(),
		AccountKey: account.AccountKey().Bytes(),
		Balance:    result,
	}, nil
}

// Issue a send transacgion
func (server SDKServer) Send(ctx context.Context, request *pb.SendRequest) (*pb.SendReply, error) {
	findAccount := server.App.Accounts.FindName

	currency := currencyString(request.Currency)
	fee := data.NewCoinFromFloat(request.Fee, "OLT")
	gas := data.NewCoinFromUnits(request.Gas, "OLT")
	sendAmount := data.NewCoinFromFloat(request.Amount, currency)

	// Get party & counterparty accounts
	partyAccount, err := findAccount(request.Party)
	if err != status.SUCCESS {
		return nil, errAccountNotFound(request.Party)
	}

	counterPartyAccount, err := findAccount(request.CounterParty)
	if err != status.SUCCESS {
		return nil, errAccountNotFound(request.CounterParty)
	}

	send, errr := prepareSend(partyAccount, counterPartyAccount, sendAmount, fee, gas, server.App)
	if errr != nil {
		return &pb.SendReply{Ok: false, Reason: errr.Error()}, nil
	}

	packet := action.SignAndPack(send)

	// TODO: Include ResultBroadcastTxCommit in SendReply
	comm.Broadcast(packet)

	return &pb.SendReply{
		Ok:     true,
		Reason: "",
	}, nil
}

func (server SDKServer) Tx(ctx context.Context, request *pb.TxRequest) (*pb.SDKReply, error) {
	result := comm.Tx(request.Hash, request.Proof)

	buff, err := serialize.JSONSzr.Serialize(result)
	if err != nil {
		return nil, gstatus.Error(codes.Internal, err.Error())
	}

	return &pb.SDKReply{Results: buff}, nil
}

func (server SDKServer) TxSearch(ctx context.Context, request *pb.TxSearchRequest) (*pb.SDKReply, error) {
	result := comm.Search(request.Query, request.Proof, int(request.Page), int(request.PerPage))

	buff, err := serialize.JSONSzr.Serialize(result)
	if err != nil {
		return nil, gstatus.Error(codes.Internal, err.Error())
	}

	return &pb.SDKReply{Results: buff}, nil
}

func (server SDKServer) Block(ctx context.Context, request *pb.BlockRequest) (*pb.SDKReply, error) {
	result := comm.Block(int64(request.Height))
	if result == nil {
		return nil, gstatus.Errorf(codes.Internal, "Lookup failed")
	}

	bz, err := json.Marshal(result)
	if err != nil {
		return nil, gstatus.Errorf(codes.Internal, "Serializing result failed %s", err)
	}

	return &pb.SDKReply{Results: bz}, nil
}

func prepareSend(
	party id.Account,
	counterParty id.Account,
	sendAmount data.Coin,
	fee data.Coin,
	gas data.Coin,
	app *Application,
) (*action.Send_Absolute, error) {
	// TODO: Use functions in shared package after resolving import cycles
	findBalance := func(key id.AccountKey) (*data.Balance, error) {
		balance := app.Balances.Get(data.DatabaseKey(key), false)
		if balance == nil {
			return nil, fmt.Errorf("Balance not found for key %x", key)
		}
		return balance, nil
	}

	pKey := party.AccountKey()
	pBalance, errr := findBalance(pKey)
	if errr != nil {
		return nil, fmt.Errorf("Account %s has insufficient balance", party.Name())
	}

	cpKey := counterParty.AccountKey()
	cpBalance, errr := findBalance(cpKey)
	if errr != nil {
		return nil, fmt.Errorf("Account %s has insufficient balance", counterParty.Name())
	}

	inputs := []action.SendInput{
		action.NewSendInput(pKey, pBalance.GetAmountByName("OLT")),
		action.NewSendInput(cpKey, cpBalance.GetAmountByName("OLT")),
	}

	outputs := []action.SendOutput{
		action.NewSendOutput(pKey, pBalance.GetAmountByName("OLT").Minus(sendAmount)),
		action.NewSendOutput(cpKey, cpBalance.GetAmountByName("OLT").Plus(sendAmount)),
	}

	return &action.Send_Absolute{
		Base: action.Base{
			Type:    action.SEND,
			ChainId: ChainId,
			// GetSigners not implemneted
			Signers:  nil,
			Sequence: global.Current.Sequence,
		},
		Inputs:  inputs,
		Outputs: outputs,
		Fee:     fee,
		Gas:     gas,
	}, nil
}

func currencyProtobuf(c data.Currency) pb.Currency {
	switch c.Chain {
	case data.ONELEDGER:
		return pb.Currency_OLT
	case data.BITCOIN:
		return pb.Currency_BTC
	case data.ETHEREUM:
		return pb.Currency_ETH
	}
	// Shouldn't ever reach here
	log.Error("Converting nonexistent currency!!!")
	return pb.Currency_OLT
}

// Functions for returning gRPC errors
func errAccountNotFound(account string) error {
	return gstatus.Errorf(codes.NotFound, "Account %s not found", account)
}

func currencyString(currency pb.Currency) string {
	switch currency {
	case pb.Currency_OLT:
		return "OLT"
	case pb.Currency_ETH:
		return "ETH"
	case pb.Currency_BTC:
		return "BTC"
	}

	return "OLT"
}

func parseChainType(chain pb.ChainType) data.ChainType {
	switch chain {
	case pb.ChainType_ONELEDGER:
		return data.ONELEDGER
	case pb.ChainType_BITCOIN:
		return data.BITCOIN
	case pb.ChainType_ETHEREUM:
		return data.ETHEREUM
	}
	return data.UNKNOWN
}

func secret(str string) []byte {
	// TODO: proper secret for key generation
	return []byte(str + global.Current.NodeName)
}
