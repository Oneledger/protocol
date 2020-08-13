package bidding

import "github.com/Oneledger/protocol/data/ons"

type TestAsset struct {
	Lalala ons.Name `json:"lalala"`
}

func NewTestAsset(lalala string) *TestAsset {
	return &TestAsset{
		Lalala: ons.Name(lalala),
	}
}

func (da *TestAsset) ToString() string {
	return string(da.Lalala)
}
