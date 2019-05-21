/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package data

import (
	"testing"

	"github.com/Oneledger/protocol/serialize"

	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	f := map[string]interface{}{"a": "b"}
	r, err := NewRequest("Query", f)
	assert.NoError(t, err)
	assert.Equal(t, r.Query, "Query")
	assert.Equal(t, r.Params, f)
	assert.NotNil(t, r.Data)

	r, err = NewRequest("Query", nil)
	assert.NoError(t, err)
	assert.Equal(t, r.Query, "Query")
	assert.Nil(t, r.Params)

	r = NewRequestFromData("asdf", []byte("asdf"))
	assert.Equal(t, r.Query, "asdf")
	assert.Equal(t, r.Data, []byte("asdf"))

	r = NewRequestFromData("1234", []byte("asdf"))
	assert.Equal(t, r.Query, "1234")
	assert.Equal(t, r.Data, []byte("asdf"))

	r = NewRequestFromData("1234", nil)
	assert.Equal(t, r.Query, "1234")
	assert.Nil(t, r.Data)

	r, err = NewRequestFromObj("poiuy", 90)
	assert.Nil(t, err)
	assert.Equal(t, r.Query, "poiuy")
	r.Parse()

	var i int
	err = clSzlr.Deserialize(r.Data, &i)
	assert.Nil(t, err)
	assert.Equal(t, i, 90)

}

func TestRequest(t *testing.T) {

	bytes := []byte("asdlfkjasldkjf")
	f := map[string]interface{}{"a": "b", "str": "str", "int": 1, "bool": false,
		"float": 67.89, "bytes": bytes}
	r, err := NewRequest("Query", f)
	assert.NoError(t, err)
	assert.Equal(t, r.Query, "Query")

	assert.Equal(t, r.GetString("a"), "b")
	assert.Equal(t, r.GetString("str"), "str")
	assert.Equal(t, r.GetString("bad_key"), "")
	assert.Equal(t, r.GetString("int"), "")

	assert.Equal(t, r.GetInt("int"), 1)
	assert.Equal(t, r.GetInt("bad_key"), 0)
	assert.Equal(t, r.GetInt("a"), 0)

	b, err := r.GetBool("bool")
	assert.NoError(t, err)
	assert.False(t, b)
	assert.NotEqual(t, b, true)

	b, err = r.GetBool("bad_key")
	assert.Error(t, err)
	assert.Equal(t, err, ErrParamNotFound)
	assert.False(t, b)

	b, err = r.GetBool("a")
	assert.Error(t, err)
	assert.Equal(t, err, ErrWrongParamType)
	assert.False(t, b)

	assert.Equal(t, r.GetFloat64("float"), 67.89)
	assert.Equal(t, r.GetFloat64("bad_key"), 0.0)
	assert.Equal(t, r.GetFloat64("a"), 0.0)

	assert.Equal(t, r.GetBytes("bytes"), bytes)
	assert.Len(t, r.GetBytes("bad_key"), 0)
	assert.Len(t, r.GetBytes("a"), 0)

	cd := &ContractData{[]byte("akshqioeroinkjlasnqbnkjfbnalksd"), map[string]interface{}{"asdfasdf": 12341234}}
	rq, err := NewRequestFromObj("query", cd)
	assert.NoError(t, err)
	assert.Equal(t, rq.Query, "query")

	cdr := &ContractData{}
	err = rq.ClientDeserialize("data", cdr)
	assert.NoError(t, err)
	assert.NotEqual(t, cdr, cd)
}

func TestResponse_Error(t *testing.T) {
	resp := &Response{}
	resp.Error("some error msg")

	assert.Equal(t, resp.ErrorMsg, "some error msg")
	assert.False(t, resp.Success)
	assert.Nil(t, resp.Data)
}

func TestResponse_SetData(t *testing.T) {
	d := []byte("some binary data kjlwsehrh8923uj4tjiuhc89340ujtikfna;[092uj4ljkldjflkdnfljojjml;akdf;lk")

	resp := &Response{}
	resp.SetData(d)

	assert.Equal(t, resp.Data, d)
	assert.True(t, resp.Success)
	assert.Equal(t, resp.ErrorMsg, "")
}

func TestResponse_SetDataObj(t *testing.T) {
	cd := &ContractData{[]byte("lkhj930o4ji90jf9023uj1239hef0-293jraase"),
		map[string]interface{}{"lll": "9108741098374"}}

	resp := &Response{}
	err := resp.SetDataObj(cd)

	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, resp.ErrorMsg, "")
	assert.NotNil(t, resp.Data)

	cdr := &ContractData{}
	err = clSzlr.Deserialize(resp.Data, cdr)
	assert.NoError(t, err)
	assert.Equal(t, cdr, cd)
}

func TestResponse_JSON(t *testing.T) {

	cd := &ContractData{[]byte("lkhj930o4ji90jf9023uj1239hef0-293jraase"),
		map[string]interface{}{"lll": "9108741098374"}}

	resp := &Response{}
	err := resp.JSON(cd)

	assert.NoError(t, err)
	assert.True(t, resp.Success)

	cdJson, err := serialize.GetSerializer(serialize.JSON).Serialize(cd)
	assert.NoError(t, err)
	assert.Equal(t, resp.Data, cdJson)
}
