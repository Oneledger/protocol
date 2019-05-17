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
	var i int
	err = clSzlr.Deserialize(r.Data, &i)
	assert.Nil(t, err)
	assert.Equal(t, i, 90)

}

func TestRequest(t *testing.T) {

	f := map[string]interface{}{"a": "b", "str": "str", "int": 1, "bool": false}
	r, err := NewRequest("Query", f)
	assert.NoError(t, err)
	assert.Equal(t, r.Query, "Query")

	assert.Equal(t, r.GetString("a"), "b")
	assert.Equal(t, r.GetString("str"), "str")
	assert.Equal(t, r.GetInt("int"), 1)

	b, err := r.GetBool("bool")
	assert.NoError(t, err)
	assert.False(t, b)
	assert.NotEqual(t, b, true)

}
