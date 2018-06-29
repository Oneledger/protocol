package rpc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/Oneledger/protocol/node/log"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

// Bitcoind - represents a Bitcoind client
type Bitcoind struct {
	client      *BTCRpcClient
	ChainParams *chaincfg.Params
}

// New - return a new bitcoind
func New(host string, port int, user, passwd string, useSSL bool, chainParams *chaincfg.Params) (*Bitcoind, error) {
	BTCRpcClient, err := newClient(host, port, user, passwd, useSSL)
	if err != nil {
		return nil, err
	}
	return &Bitcoind{client: BTCRpcClient, ChainParams: chainParams}, nil
}

// GetAccount returns the account associated with the given address.
func (b *Bitcoind) GetAccount(address string) (account string, err error) {
	r, err := b.client.call("getaccount", []string{address})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &account)
	return
}

// GetAccountAddress Returns the current bitcoin address for receiving
// payments to this account.
// If account does not exist, it will be created along with an
// associated new address that will be returned.
func (b *Bitcoind) GetAccountAddress(account string) (address string, err error) {
	r, err := b.client.call("getaccountaddress", []string{account})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &address)
	return
}

// GetAddressesByAccount return addresses associated with account <account>
func (b *Bitcoind) GetAddressesByAccount(account string) (addresses []string, err error) {
	r, err := b.client.call("getaddressesbyaccount", []string{account})
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &addresses)
	return
}

// GetBalance return the balance of the server or of a specific account
//If [account] is "", returns the server's total available balance.
//If [account] is specified, returns the balance in the account
func (b *Bitcoind) GetBalance(account string, minconf uint64) (balance float64, err error) {
	r, err := b.client.call("getbalance", []interface{}{account, minconf})
	if err = handleError(err, &r); err != nil {
		return
	}

	balance, err = strconv.ParseFloat(string(r.Result), 64)
	return
}

// GetBestBlockhash returns the hash of the best (tip) block in the longest block chain.
func (b *Bitcoind) GetBestBlockhash() (bestBlockHash string, err error) {
	r, err := b.client.call("getbestblockhash", nil)
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &bestBlockHash)
	return
}

// GetBlock returns information about the block with the given hash.
func (b *Bitcoind) GetBlock(blockHash string) (block Block, err error) {
	r, err := b.client.call("getblock", []string{blockHash})
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &block)
	return
}

// GetBlockCount returns the number of blocks in the longest block chain.
func (b *Bitcoind) GetBlockCount() (count uint64, err error) {
	r, err := b.client.call("getblockcount", nil)
	if err = handleError(err, &r); err != nil {
		return
	}

	count, err = strconv.ParseUint(string(r.Result), 10, 64)
	return
}

// GetBlockHash returns hash of block in best-block-chain at <index>
func (b *Bitcoind) GetBlockHash(index uint64) (hash string, err error) {
	r, err := b.client.call("getblockhash", []uint64{index})
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &hash)
	return
}

// GetDifficulty returns the proof-of-work difficulty as a multiple of
// the minimum difficulty.
func (b *Bitcoind) GetDifficulty() (difficulty float64, err error) {
	r, err := b.client.call("getdifficulty", nil)
	if err = handleError(err, &r); err != nil {
		return
	}

	difficulty, err = strconv.ParseFloat(string(r.Result), 64)
	return
}

// GetGenerate returns true or false whether bitcoind is currently generating hashes
func (b *Bitcoind) GetGenerate() (generate bool, err error) {
	r, err := b.client.call("getgenerate", nil)
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &generate)
	return
}

// GetHashesPerSec returns a recent hashes per second performance measurement while generating.
func (b *Bitcoind) GetHashesPerSec() (hashpersec float64, err error) {
	r, err := b.client.call("gethashespersec", nil)
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &hashpersec)
	return
}

// GetNewAddress return a new address for account [account].
func (b *Bitcoind) GetNewAddress(account ...string) (addr string, err error) {
	// 0 or 1 account
	if len(account) > 1 {
		err = errors.New("Bad parameters for GetNewAddress: you can set 0 or 1 account")
		return
	}

	r, err := b.client.call("getnewaddress", account)
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &addr)
	return
}

// GetRawMempool returns all transaction ids in memory pool
func (b *Bitcoind) GetRawMempool() (txId []string, err error) {
	r, err := b.client.call("getrawmempool", nil)
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &txId)
	return
}

// GetRawTransaction returns raw transaction representation for given transaction id.
func (b *Bitcoind) GetRawTransaction(txId string, verbose bool) (rawTx interface{}, err error) {
	r, err := b.client.call("getrawtransaction", []interface{}{txId, verbose})
	if err = handleError(err, &r); err != nil {
		return
	}

	if !verbose {
		err = json.Unmarshal(r.Result, &rawTx)
	} else {
		var t RawTransaction
		err = json.Unmarshal(r.Result, &t)
		rawTx = t
	}

	return
}

// GetReceivedByAccount Returns the total amount received by addresses with [account] in
// transactions with at least [minconf] confirmations. If [account] is set to all return
// will include all transactions to all accounts
func (b *Bitcoind) GetReceivedByAccount(account string, minconf uint32) (amount float64, err error) {
	if account == "all" {
		account = ""
	}

	r, err := b.client.call("getreceivedbyaccount", []interface{}{account, minconf})
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &amount)
	return
}

// Returns the amount received by <address> in transactions with at least [minconf] confirmations.
// It correctly handles the case where someone has sent to the address in multiple transactions.
// Keep in mind that addresses are only ever used for receiving transactions. Works only for addresses
// in the local wallet, external addresses will always show 0.
func (b *Bitcoind) GetReceivedByAddress(address string, minconf uint32) (amount float64, err error) {
	r, err := b.client.call("getreceivedbyaddress", []interface{}{address, minconf})
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &amount)
	return
}

// GetTransaction returns a Bitcoind.Transation struct about the given transaction
func (b *Bitcoind) GetTransaction(txid string) (transaction Transaction, err error) {
	r, err := b.client.call("gettransaction", []interface{}{txid})
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &transaction)
	return
}

// GetTxOut returns details about an unspent transaction output (UTXO)
func (b *Bitcoind) GetTxOut(txid string, n uint32, includeMempool bool) (transactionOut UTransactionOut, err error) {
	r, err := b.client.call("gettxout", []interface{}{txid, n, includeMempool})
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &transactionOut)
	return
}

// GetTxOutsetInfo returns statistics about the unspent transaction output (UTXO) set
func (b *Bitcoind) GetTxOutsetInfo() (txOutSet TransactionOutSet, err error) {
	r, err := b.client.call("gettxoutsetinfo", nil)
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &txOutSet)
	return
}

// ListAccounts returns Object that has account names as keys, account balances as values.
func (b *Bitcoind) ListAccounts(minconf int32) (accounts map[string]float64, err error) {
	r, err := b.client.call("listaccounts", []int32{minconf})
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &accounts)
	return
}

// ListSinceBlock
func (b *Bitcoind) ListSinceBlock(blockHash string, targetConfirmations uint32) (transaction []Transaction, err error) {
	r, err := b.client.call("listsinceblock", []interface{}{blockHash, targetConfirmations})
	if err = handleError(err, &r); err != nil {
		return
	}

	type ts struct {
		Transactions []Transaction
	}

	var result ts
	if err = json.Unmarshal(r.Result, &result); err != nil {
		return
	}
	transaction = result.Transactions
	return
}

// ListTransactions returns up to [count] most recent transactions skipping the first
// [from] transactions for account [account]. If [account] not provided it'll return
// recent transactions from all accounts.
func (b *Bitcoind) ListTransactions(account string, count, from uint32) (transaction []Transaction, err error) {
	r, err := b.client.call("listtransactions", []interface{}{account, count, from})
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &transaction)
	return
}

// ListUnspent returns array of unspent transaction inputs in the wallet.
func (b *Bitcoind) ListUnspent(minconf, maxconf uint32) (transactions []Transaction, err error) {
	if maxconf > 999999 {
		maxconf = 999999
	}

	r, err := b.client.call("listunspent", []interface{}{minconf, maxconf})
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &transactions)
	return
}

// ListLockUnspent returns list of temporarily unspendable outputs
func (b *Bitcoind) ListLockUnspent() (unspendableOutputs []UnspendableOutput, err error) {
	r, err := b.client.call("listlockunspent", nil)
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &unspendableOutputs)
	return
}

// LockUnspent updates(lock/unlock) list of temporarily unspendable outputs
func (b *Bitcoind) LockUnspent(lock bool, outputs []UnspendableOutput) (success bool, err error) {
	r, err := b.client.call("lockunspent", []interface{}{lock, outputs})
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &success)
	return
}

// Move from one account in your wallet to another
func (b *Bitcoind) Move(formAccount, toAccount string, amount float64, minconf uint32, comment string) (success bool, err error) {
	r, err := b.client.call("move", []interface{}{formAccount, toAccount, amount, minconf, comment})
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &success)
	return

}

// SendFrom send amount from fromAccount to toAddress
//  amount is a real and is rounded to 8 decimal places.
//  Will send the given amount to the given address, ensuring the account has a valid balance using [minconf] confirmations.
func (b *Bitcoind) SendFrom(fromAccount, toAddress string, amount float64, minconf uint32, comment, commentTo string) (txID string, err error) {
	r, err := b.client.call("sendfrom", []interface{}{fromAccount, toAddress, amount, minconf, comment, commentTo})
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &txID)
	return
}

// SendToAddress send an amount to a given address
func (b *Bitcoind) SendToAddress(toAddress string, amount float64, comment, commentTo string) (txID string, err error) {
	r, err := b.client.call("sendtoaddress", []interface{}{toAddress, amount, comment, commentTo})
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &txID)
	return
}

// Stop stop bitcoin server.
func (b *Bitcoind) Stop() error {
	r, err := b.client.call("stop", nil)
	return handleError(err, &r)
}

// ValidateAddress return information about <bitcoinaddress>.
func (b *Bitcoind) ValidateAddress(address string) (va ValidateAddressResponse, err error) {
	r, err := b.client.call("validateaddress", []interface{}{address})
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &va)
	return
}

//generate blockNumber of block on regtest network
func (b *Bitcoind) Generate(blockNumber uint64) (bh []string, err error) {
	r, err := b.client.call("generate", []uint64{blockNumber})
	if err = handleError(err, &r); err != nil {
		return
	}

	err = json.Unmarshal(r.Result, &bh)
	return
}

// DumpPrivKey return private key as string associated to public <address>
func (b *Bitcoind) DumpPrivKey(address btcutil.Address) (*btcutil.WIF, error) {
	addr := address.EncodeAddress()
	res, err := b.client.call("dumpprivkey", []string{addr})
	log.Debug("Call dumpprivkey", "res", res, "err", err, "address", address, "addr", addr)

	// Unmarshal result as a string.
	var privKeyWIF string
	err = json.Unmarshal(res.Result, &privKeyWIF)
	if err != nil {
		log.Debug("Unmarshal", "res", res)
		return nil, err
	}

	return btcutil.DecodeWIF(privKeyWIF)
}

func (b *Bitcoind) FundRawTransaction(tx *wire.MsgTx, feePerKb btcutil.Amount) (fundedTx *wire.MsgTx, fee btcutil.Amount, err error) {
	log.Debug("FundRawTransaction", "fee", feePerKb)

	var buf bytes.Buffer
	buf.Grow(tx.SerializeSize())
	tx.Serialize(&buf)

	param0, err := json.Marshal(hex.EncodeToString(buf.Bytes()))
	if err != nil {
		return nil, 0, err
	}

	param1, err := json.Marshal(struct {
		FeeRate float64 `json:"feeRate"`
		//Change_type string  `json:"change_type"`
	}{
		FeeRate: feePerKb.ToBTC(),
		//FeeRate:     0.00000000001,
		//Change_type: "legacy",
	})
	if err != nil {
		return nil, 0, err
	}

	params := []json.RawMessage{param0, param1}
	rawResp, err := b.client.call("fundrawtransaction", params)

	var resp struct {
		Hex       string  `json:"hex"`
		Fee       float64 `json:"fee"`
		ChangePos float64 `json:"changepos"`
	}

	if err != nil {
		log.Debug("Failed to Fund", "feePerKb", feePerKb, "params", params, "err", err, "rawResp", rawResp)
		resp.Hex = hex.EncodeToString(buf.Bytes())
		resp.Fee = 0.00001
		resp.ChangePos = -1

		//return nil, 0, err
	} else {
		err = json.Unmarshal(rawResp.Result, &resp)
		if err != nil {
			return nil, 0, err
		}
	}

	fundedTxBytes, err := hex.DecodeString(resp.Hex)
	if err != nil {
		return nil, 0, err
	}

	fundedTx = &wire.MsgTx{}
	err = fundedTx.Deserialize(bytes.NewReader(fundedTxBytes))
	if err != nil {
		return nil, 0, err
	}

	feeAmount, err := btcutil.NewAmount(resp.Fee)
	if err != nil {
		return nil, 0, err
	}
	return fundedTx, feeAmount, nil
}

// getFeePerKb queries the wallet for the transaction relay fee/kB to use and
// the minimum mempool relay fee.  It first tries to get the user-set fee in the
// wallet.  If unset, it attempts to find an estimate using estimatefee 6.  If
// both of these fail, it falls back to mempool relay fee policy.
func (b *Bitcoind) GetFeePerKb() (useFee, relayFee btcutil.Amount, err error) {
	var netInfoResp struct {
		RelayFee float64 `json:"relayfee"`
	}
	var walletInfoResp struct {
		PayTxFee float64 `json:"paytxfee"`
	}
	var estimateResp struct {
		FeeRate float64 `json:"feerate"`
	}

	netInfoRawResp, err := b.client.call("getnetworkinfo", nil)
	if err == nil {
		err = json.Unmarshal(netInfoRawResp.Result, &netInfoResp)
		if err != nil {
			return 0, 0, err
		}
	}

	walletInfoRawResp, err := b.client.call("getwalletinfo", nil)
	if err == nil {
		err = json.Unmarshal(walletInfoRawResp.Result, &walletInfoResp)
		if err != nil {
			return 0, 0, err
		}
	}

	relayFee, err = btcutil.NewAmount(netInfoResp.RelayFee)
	if err != nil {
		return 0, 0, err
	}

	payTxFee, err := btcutil.NewAmount(walletInfoResp.PayTxFee)
	if err != nil {
		return 0, 0, err
	}

	// Use user-set wallet fee when set and not lower than the network relay
	// fee.
	if payTxFee != 0 {
		maxFee := payTxFee
		if relayFee > maxFee {
			maxFee = relayFee
		}
		return maxFee, relayFee, nil
	}

	params := []json.RawMessage{[]byte("6")}
	estimateRawResp, err := b.client.call("estimatesmartfee", params)
	if err != nil {
		return 0, 0, err
	}

	err = json.Unmarshal(estimateRawResp.Result, &estimateResp)
	if err == nil && estimateResp.FeeRate > 0 {
		useFee, err = btcutil.NewAmount(estimateResp.FeeRate)
		if relayFee > useFee {
			useFee = relayFee
		}
		return useFee, relayFee, err
	}

	fmt.Println("warning: falling back to mempool relay fee policy")
	return relayFee, relayFee, nil
}

// getRawChangeAddress calls the getrawchangeaddress JSON-RPC method.  It is
// implemented manually as the rpcclient implementation always passes the
// account parameter which was removed in Bitcoin Core 0.15.
func (b *Bitcoind) GetRawChangeAddress() (btcutil.Address, error) {
	params := []json.RawMessage{[]byte(`"legacy"`)}

	rawResp, err := b.client.call("getrawchangeaddress", params)
	if err != nil {
		return nil, err
	}

	var addrStr string
	err = json.Unmarshal(rawResp.Result, &addrStr)
	if err != nil {
		return nil, err
	}
	addr, err := btcutil.DecodeAddress(addrStr, b.ChainParams)
	if err != nil {
		return nil, err
	}

	/*
		// TODO: Failing
		if !addr.IsForNet(chainParams) {
			log.Debug("Failed for Net")
			return nil, fmt.Errorf("address %v is not intended for use on %v",
				addrStr, chainParams.Name)
		}
	*/

	if _, ok := addr.(*btcutil.AddressPubKeyHash); !ok {
		log.Debug("Failed type conversion")
		return nil, fmt.Errorf("getrawchangeaddress: address %v is not P2PKH",
			addr)
	}
	return addr, nil
}

func (b *Bitcoind) PublishTx(tx *wire.MsgTx, name string) (*chainhash.Hash, error) {
	txHex := ""
	if tx != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		if err := tx.Serialize(buf); err != nil {

		}
		txHex = hex.EncodeToString(buf.Bytes())
	}

	txHash, err := b.client.call("sendrawtransaction", []interface{}{txHex, false})
	if err != nil {
		return nil, fmt.Errorf("sendrawtransaction: %v", err)
	}
	fmt.Printf("Published %s transaction (%v)\n", name, txHash)

	// Unmarshal result as a string.
	var txHashStr string
	err = json.Unmarshal(txHash.Result, &txHashStr)
	if err != nil {
		return nil, err
	}

	return chainhash.NewHashFromStr(txHashStr)
}

// SignRawTransaction signs inputs for the passed transaction and returns the
// signed transaction as well as whether or not all inputs are now signed.
//
// This function assumes the RPC server already knows the input transactions and
// private keys for the passed transaction which needs to be signed and uses the
// default signature hash type.  Use one of the SignRawTransaction# variants to
// specify that information if needed.
func (b *Bitcoind) SignRawTransaction(tx *wire.MsgTx) (*wire.MsgTx, bool, error) {
	txHex := ""
	if tx != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		if err := tx.Serialize(buf); err != nil {

		}
		txHex = hex.EncodeToString(buf.Bytes())
	}

	res, err := b.client.call("signrawtransaction", []interface{}{txHex})
	if err != nil {
		return nil, false, err
	}

	// Unmarshal as a signrawtransaction result.
	var signRawTxResult btcjson.SignRawTransactionResult
	err = json.Unmarshal(res.Result, &signRawTxResult)
	if err != nil {
		return nil, false, err
	}

	// Decode the serialized transaction hex to raw bytes.
	serializedTx, err := hex.DecodeString(signRawTxResult.Hex)
	if err != nil {
		return nil, false, err
	}

	// Deserialize the transaction and return it.
	var msgTx wire.MsgTx
	if err := msgTx.Deserialize(bytes.NewReader(serializedTx)); err != nil {
		return nil, false, err
	}

	return &msgTx, signRawTxResult.Complete, nil
}
