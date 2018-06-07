package htlc

import (
	"strconv"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	rpc "github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/wire"
	"encoding/hex"
	"bytes"
)



// There are two directions that the atomic swap can be performed, as the
// initiator can be on either chain.  This tool only deals with creating the
// Bitcoin transactions for these swaps.  A second tool should be used for the
// transaction on the other chain.  Any chain can be used so long as it supports
// OP_SHA256 and OP_CHECKLOCKTIMEVERIFY.
//
// Example scenerios using bitcoin as the second chain:
//
// Scenerio 1:
//   cp1 initiates (olt)
//   cp2 participates with cp1 H(S) (btc)
//   cp1 redeems btc revealing S
//     - must verify H(S) in contract is hash of known secret
//   cp2 redeems olt with S
//
// Scenerio 2:
//   cp1 initiates (btc)
//   cp2 participates with cp1 H(S) (olt)
//   cp1 redeems olt revealing S
//     - must verify H(S) in contract is hash of known secret
//   cp2 redeems btc with S


// TestInitiateCommand
func TestInitiateCommand(t *testing.T) {
	chainParams = &chaincfg.TestNet3Params
	cp2Addr, err := btcutil.DecodeAddress("", chainParams)
	if err != nil {
		t.Errorf("failed to decode participant address: %v", err)
	}
	if !cp2Addr.IsForNet(chainParams) {
		t.Errorf("participant address is not intended for use on %v", chainParams.Name)
	}
	cp2AddrP2PKH, ok := cp2Addr.(*btcutil.AddressPubKeyHash)
	if !ok {
		t.Errorf("participant address is not P2PKH")
	}

	amountF64, err := strconv.ParseFloat("", 64)
	if err != nil {
		t.Errorf("failed to decode amount: %v", err)
	}
	amount, err := btcutil.NewAmount(amountF64)
	if err != nil {
		t.Errorf("failed to decode amount: %v", err)
	}

	var cmd command
	cmd = &initiateCmd{cp2Addr: cp2AddrP2PKH, amount: amount}
	runCommand(cmd)
}


func runCommand(cmd command) {
	// Offline commands don't need to talk to the wallet.
	if cmd, ok := cmd.(offlineCommand); ok {
		cmd.runOfflineCommand()
	}

	connect, err := normalizeAddress(*connectFlag, walletPort(chainParams))
	if err != nil {
		t.Errorf("wallet server address: %v", err)
	}

	connConfig := &rpc.ConnConfig{
		Host:         connect,
		User:         *rpcuserFlag,
		Pass:         *rpcpassFlag,
		DisableTLS:   true,
		HTTPPostMode: true,
	}
	client, err := rpc.New(connConfig, nil)
	if err != nil {
		t.Errorf("rpc connect: %v", err)
	}
	defer func() {
		client.Shutdown()
		client.WaitForShutdown()
	}()

	err = cmd.runCommand(client)
}


func TestRedeemCommand(t *testing.T) {
	contract, err := hex.DecodeString("")
	if err != nil {
		t.Errorf("failed to decode contract: %v", err)
	}

	contractTxBytes, err := hex.DecodeString("")
	if err != nil {
		t.Errorf("failed to decode contract transaction: %v", err)
	}
	var contractTx wire.MsgTx
	err = contractTx.Deserialize(bytes.NewReader(contractTxBytes))
	if err != nil {
		t.Errorf("failed to decode contract transaction: %v", err)
	}

	secret, err := hex.DecodeString("")
	if err != nil {
		t.Errorf("failed to decode secret: %v", err)
	}

	var cmd command
	cmd = &redeemCmd{contract: contract, contractTx: &contractTx, secret: secret}
	runCommand(cmd)
}