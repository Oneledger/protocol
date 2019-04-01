package serialize

import (
	"github.com/Oneledger/protocol/node/log"
	"github.com/tendermint/go-amino"
)

type aminoStrategy struct {
	metaStrategy
	codec *amino.Codec
}


//NewAminoStrategy generates a new object for amino serialization with amino codec
func NewAminoStrategy(cdc *amino.Codec) *aminoStrategy {
	return &aminoStrategy{codec: cdc}
}


func (a *aminoStrategy) Serialize(obj interface{}) ([]byte, error) {

	if _, ok := obj.(DataAdapter); ok {

		log.Warn("amino strategy does not support adapters")
	}
	bz, err := a.codec.MarshalBinaryLengthPrefixed(obj)

	return bz, err
}

//Deserialize
func (a *aminoStrategy) Deserialize(src []byte, dest interface{}) error {

	/*
	if apr, ok := dest.(adapter.DataAdapter); ok {
		return a.wrapDataAdapter(src, dest, a.deserialize, apr)
	}

	return a.deserialize(src, dest)
	*/

	err := a.deserialize(src, dest)
	return err
}

//deserialize
func (a *aminoStrategy) deserialize(src []byte, dest interface{}) error {

	err :=  a.codec.UnmarshalBinaryLengthPrefixed(src, dest)

	return err
}

