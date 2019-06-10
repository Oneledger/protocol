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

func TestOLVMContext_GetValue(t *testing.T) {
	o := OLVMContext{}

	o.Data = map[string]interface{}{
		"key": []byte("value"),
	}

	assert.NotNil(t, o.Data)

	val := o.GetValue("key")
	assert.Equal(t, val, "value")

	val = o.GetValue("badKey")
	assert.Nil(t, val)
}

func TestNewOLVMRequest(t *testing.T) {
	oc := OLVMContext{}

	or := NewOLVMRequest([]byte("var a = 'a';"), oc)

	assert.NotNil(t, or)
	assert.Equal(t, or.CallString, "")
	assert.Equal(t, or.SourceCode, []byte("var a = 'a';"))
	assert.Equal(t, or.Context, oc)
}

func TestNewOLVMResult(t *testing.T) {
	or := NewOLVMResult()

	assert.NotNil(t, or)
	assert.Equal(t, or.Status, "MISSING DATA")
}
