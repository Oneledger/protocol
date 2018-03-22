package types

import (
	"net/url"

	"bytes"
	"errors"
	cmn "github.com/tendermint/tmlibs/common"
)

var (
	ErrRefInvalidType	            = errors.New("Invalid Reference Type")
	ErrRefInvalidHash			    = errors.New("Invalid Reference Hash")
	ErrRefInvalidUrl				= errors.New("Invalid Reference Url")
)


type Reference struct {
	Type string `json:"type"`
	Url url.URL `json:"type"`
	ReferenceHash cmn.HexBytes `json:"referenceHash"`
}

func (referenceA *Reference) verify(referenceType string, hash cmn.HexBytes) error {
	if referenceA.Type != referenceType{
		return ErrRefInvalidType
	}else if bytes.Equal(referenceA.ReferenceHash.Bytes(), hash.Bytes()){
		return ErrRefInvalidHash
	}
	return nil
}