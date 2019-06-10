package serialize

import (
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
	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic in amino", r)
			// debug.PrintStack()
		}
	}()

	if apr, ok := obj.(DataAdapter); ok {
		obj = apr.Data()
	}
	bz, err := a.codec.MarshalBinaryLengthPrefixed(obj)

	return bz, err
}

//Deserialize
func (a *aminoStrategy) Deserialize(src []byte, dest interface{}) error {

	if apr, ok := dest.(DataAdapter); ok {
		err := a.wrapDataAdapter(src, dest, a.deserialize, apr)
		return err
	}

	err := a.deserialize(src, dest)
	return err
}

//deserialize
func (a *aminoStrategy) deserialize(src []byte, dest interface{}) error {

	err := a.codec.UnmarshalBinaryLengthPrefixed(src, dest)

	return err
}
