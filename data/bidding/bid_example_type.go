package bidding

import (
	"github.com/Oneledger/protocol/action"
	"strconv"
)

type ExampleAsset struct {
	ExampleField int `json:"exampleField"`
}

func NewExampleAsset(exampleValue int) *ExampleAsset {
	return &ExampleAsset{
		ExampleField: exampleValue,
	}
}

func (ta *ExampleAsset) ToString() string {
	return strconv.Itoa(ta.ExampleField)
}

func (ta *ExampleAsset) ValidateAsset(ctx *action.Context) (bool, error) {
	return true, nil
}
