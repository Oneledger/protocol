package rpc

// Block - Represents a block
type Block struct {
	// The block hash
	Hash string `json:"hash"`
	// The number of confirmations
	Confirmations uint64 `json:"confirmations"`
	// The block size
	Size uint64 `json:"size"`
	// The block height or index
	Height uint64 `json:"height"`
	// The block version
	Version uint32 `json:"version"`
	// The merkle root
	Merkleroot string `json:"merkleroot"`
	// The block time in seconds since epoch (Jan 1 1970 GMT)
	Time int64 `json:"time"`
	// The nonce
	Nonce uint64 `json:"nonce"`
	// The bits
	Bits string `json:"bits"`
	// The difficulty
	Difficulty float64 `json:"difficulty"`
	// Total amount of work in active chain, in hexadecimal
	Chainwork string `json:"chainwork,omitempty"`
	// The hash of the previous block
	Previousblockhash string `json:"previousblockhash"`
	// The hash of the next block
	Nextblockhash string `json:"nextblockhash"`
	// Slice on transaction ids
	Tx []string `json:"tx"`
}

// ScriptSig - represents a scriptsyg
type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

// ScriptPubKey - represents a scriptpubkey
type ScriptPubKey struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	ReqSigs   int      `json:"reqSigs,omitempty"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses,omitempty"`
}

// Vin - represent an IN value
type Vin struct {
	Coinbase  string    `json:"coinbase"`
	Txid      string    `json:"txid"`
	Vout      int       `json:"vout"`
	ScriptSig ScriptSig `json:"scriptSig"`
	Sequence  uint32    `json:"sequence"`
}

// Vout - represent an OUT value
type Vout struct {
	Value        float64      `json:"value"`
	N            int          `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

// UnspendableOutput - represents a unspendable (locked) output
type UnspendableOutput struct {
	TxId string `json:"txid"`
	Vout uint64 `json:"vout"`
}

// ValidateAddressResponse - represents a response to "validateaddress" call
type ValidateAddressResponse struct {
	IsValid      bool   `json:"isvalid"`
	Address      string `json:"address"`
	IsMine       bool   `json:"ismine"`
	IsScript     bool   `json:"isscript"`
	PubKey       string `json:"pubkey"`
	IsCompressed bool   `json:"iscompressed"`
	Account      string `json:"account"`
}

// Transaction - represents a transaction
type Transaction struct {
	Amount          float64              `json:"amount"`
	Account         string               `json:"account,omitempty"`
	Address         string               `json:"address,omitempty"`
	Category        string               `json:"category,omitempty"`
	Fee             float64              `json:"fee,omitempty"`
	Confirmations   int64                `json:"confirmations"`
	BlockHash       string               `json:"blockhash"`
	BlockIndex      int64                `json:"blockindex"`
	BlockTime       int64                `json:"blocktime"`
	TxID            string               `json:"txid"`
	WalletConflicts []string             `json:"walletconflicts"`
	Time            int64                `json:"time"`
	TimeReceived    int64                `json:"timereceived"`
	Details         []TransactionDetails `json:"details,omitempty"`
	Hex             string               `json:"hex,omitempty"`
}

// RawTx - represents a raw transaction
type RawTransaction struct {
	Hex           string `json:"hex"`
	Txid          string `json:"txid"`
	Version       uint32 `json:"version"`
	LockTime      uint32 `json:"locktime"`
	Vin           []Vin  `json:"vin"`
	Vout          []Vout `json:"vout"`
	BlockHash     string `json:"blockhash,omitempty"`
	Confirmations uint64 `json:"confirmations,omitempty"`
	Time          int64  `json:"time,omitempty"`
	Blocktime     int64  `json:"blocktime,omitempty"`
}

// TransactionDetails - represents details about a transaction
type TransactionDetails struct {
	Account  string  `json:"account"`
	Address  string  `json:"address,omitempty"`
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
	Fee      float64 `json:"fee,omitempty"`
}

// UTransactionOut - represents a unspent transaction out (UTXO)
type UTransactionOut struct {
	Bestblock     string       `json:"bestblock"`
	Confirmations uint32       `json:"confirmations"`
	Value         float64      `json:"value"`
	ScriptPubKey  ScriptPubKey `json:"scriptPubKey"`
	Version       uint32       `json:"version"`
	Coinbase      bool         `json:"coinbase"`
}

// TransactionOutSet - represents statistics about the unspent transaction output database
type TransactionOutSet struct {
	Height          uint32  `json:"height"`
	Bestblock       string  `json:"bestblock"`
	Transactions    float64 `json:"transactions"`
	TxOuts          float64 `json:"txouts"`
	BytesSerialized float64 `json:"bytes_serialized"`
	HashSerialized  string  `json:"hash_serialized"`
	TotalAmount     float64 `json:"total_amount"`
}