/*
	Copyright 2017-2018 OneLedger

	Convert strings (cli input) into various types.
*/
package convert

import (
	"errors"
	"strconv"
	"strings"

	crypto "github.com/tendermint/go-crypto"
)

const CURRENCY = 0x123

var Domain map[string]bool

func init() {
	Domain = make(map[string]bool)

	Domain["bitcoin"] = true
	Domain["oneledger"] = true
	Domain["etheruem"] = true
	Domain["btc"] = true
	Domain["olt"] = true
	Domain["eth"] = true
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

func (convert *Convert) GetPublicKey(value string) PublicKey {
	// TODO: Is this a file reference? If so, read in the key
	// TODO: Is this actionally the key
	return PublicKey{}
}

func (convert *Convert) GetPrivateKey(value string) PrivateKey {
	return PrivateKey{}
}

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

func (convert *Convert) GetCurrency(value string) string {
	key := strings.ToLower(value)
	if Domain[key] {
		return key
	}
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
