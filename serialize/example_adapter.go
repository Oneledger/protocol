package serialize

import (
	"errors"
	"strconv"
)

type testStuffAd struct {
	F string
	A int
	B int64
	C []byte
	H float64
}

type testStuffAdData struct {
	X string
	Y int
	Z int64
	J []byte
	K string
}

func (t *testStuffAd) NewDataInstance() Data {
	return &testStuffAdData{}
}

func (t *testStuffAd) Data() Data {
	return &testStuffAdData{t.F, t.A, t.B, t.C, strconv.FormatFloat(t.H, 'f', -1, 64)}
}

func (t *testStuffAd) SetData(a interface{}) error {
	ad, ok := a.(*testStuffAdData)
	if !ok {
		return errors.New("Wrong data")
	}

	var e error
	t.F = ad.X
	t.A = ad.Y
	t.B = ad.Z
	t.C = ad.J
	t.H, e = strconv.ParseFloat(ad.K, 64)

	return e
}

func (ad *testStuffAdData) SerialTag() string {
	return ""
}
