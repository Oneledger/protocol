package types

import (
	"net/url"

	"bytes"
	"errors"
	"github.com/tendermint/tendermint/crypto/merkle"
	cmn "github.com/tendermint/tendermint/libs/common"
)

var (
	ErrRefInvalidType = errors.New("Invalid Reference Type")
	ErrRefInvalidHash = errors.New("Invalid Reference Hash")
)

type Reference struct {
	Type          string       `json:"type"`
	Url           url.URL      `json:"url"`
	ReferenceHash cmn.HexBytes `json:"referenceHash"`
}

func (reference *Reference) verify(referenceType string, hash cmn.HexBytes) error {
	if reference.Type != referenceType {
		return ErrRefInvalidType
	} else if bytes.Equal(reference.ReferenceHash.Bytes(), hash.Bytes()) {
		return ErrRefInvalidHash
	}
	return nil
}

func (reference *Reference) Hash() cmn.HexBytes {
	return merkle.SimpleHashFromBytes([]byte(reference.Type + reference.Url.String()))
}
