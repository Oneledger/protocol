/*
	Copyright 2017-2018 OneLedger

	Convert strings (cli input) into various types.
*/
package convert

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	crypto "github.com/tendermint/go-crypto"
)

const CURRENCY = 0x123

var Domain map[string]string

func init() {
	Domain = make(map[string]string)

	// keys are lowercase, maps back to proper string
	Domain["oneledger"] = "OneLedger"
	Domain["bitcoin"] = "Bitcoin"
	Domain["etheruem"] = "Ethereum"
	Domain["olt"] = "OLT"
	Domain["btc"] = "BTC"
	Domain["eth"] = "ETH"
}

// Build up a list of error messages
type Convert struct {
	Errors map[string]error
	Index  map[string]int
	Next   int
}

func NewConvert() *Convert {
	return &Convert{
		Errors: make(map[string]error),
		Index:  make(map[string]int),
		Next:   0,
	}
}

type PublicKey = crypto.PubKey
type PrivateKey = crypto.PrivKey

func (convert *Convert) HasErrors() bool {
	if len(convert.Errors) < 1 {
		return false
	}
	return true
}

func (convert *Convert) GetErrors() string {
	buffer := ""
	for _, value := range convert.Errors {
		buffer += value.Error()
	}
	return buffer
}

func (convert *Convert) AddError(value string, err error) {
	convert.Errors[value] = err
	convert.Index[value] = convert.Next
	convert.Next++
}

// TODO: Should go through a set of different possibilties and settle on the best option
func (convert *Convert) GetAccountKey(value string) id.AccountKey {
	// TODO: See if this is an identity?
	// TODO: See if this is in utxo?
	return id.AccountKey(value)
}

// TODO: Needs to have real values
func (convert *Convert) GetPublicKey(value string) PublicKey {
	// TODO: Is this a file reference? If so, read in the key
	// TODO: Is this actionally the key
	return PublicKey{}
}

// TODO: Needs to have real values
func (convert *Convert) GetPrivateKey(value string) PrivateKey {
	return PrivateKey{}
}

// TODO: Need to be ripeMd?
func (convert *Convert) HashKey(key PublicKey) []byte {
	return nil
}

func (convert *Convert) GetHash(value string) []byte {
	result := convert.GetPublicKey(value)
	if convert.HasErrors() {
		return nil
	}
	return convert.HashKey(result)
}

func (convert *Convert) GetCoin(amountStr string, currencyStr string) data.Coin {
	currency := convert.GetCurrency(currencyStr)
	amountInt64 := convert.GetInt64(amountStr)
	return data.NewCoin(amountInt64, currency)
}

func (convert *Convert) GetCurrency(value string) string {
	key := strings.ToLower(value)
	if result, ok := Domain[key]; ok {
		return result
	}
	log.Error("MISSING Currency", "value", value)

	convert.AddError(value, errors.New("Invalid Currency"))
	return ""
}

func (convert *Convert) GetInt64(value string) int64 {
	result, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return result
	}

	convert.AddError(value, err)

	return 0
}

func (convert *Convert) GetInt(value string) int {
	// TODO: Not portable, ints match cpu arch (32 or 64)
	result, err := strconv.ParseInt(value, 10, 0)
	if err == nil {
		return int(result)
	}

	convert.AddError(value, err)

	return 0
}
