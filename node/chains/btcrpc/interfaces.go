package btcrpc

import (
	"encoding/json"
	"errors"
	"strconv"
)


// Bitcoind - represents a Bitcoind client
type Bitcoind struct {
	client *BTCRpcClient
}

// New - return a new bitcoind
func New(host string, port int, user, passwd string, useSSL bool) (*Bitcoind, error) {
	BTCRpcClient, err := newClient(host, port, user, passwd, useSSL)
	if err != nil {
		return nil, err
	}
	return &Bitcoind{BTCRpcClient}, nil
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

// GetRawChangeAddress Returns a new Bitcoin address, for receiving change.
// This is for use with raw transactions, NOT normal use.
func (b *Bitcoind) GetRawChangeAddress(account ...string) (rawAddress string, err error) {
	// 0 or 1 account
	if len(account) > 1 {
		err = errors.New("Bad parameters for GetRawChangeAddress: you can set 0 or 1 account")
		return
	}
	r, err := b.client.call("getrawchangeaddress", account)
	if err = handleError(err, &r); err != nil {
		return
	}
	err = json.Unmarshal(r.Result, &rawAddress)
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