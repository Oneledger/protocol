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

package balance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCurrencyList(t *testing.T) {
	cl := NewCurrencyList()

	curr := Currency{"OLT", 0, 18}
	err := cl.Register(curr)
	assert.NoError(t, err)

	currNew, ok := cl.GetCurrencyByName("OLT")
	assert.True(t, ok)

	assert.Equal(t, currNew, curr)

	key := curr.StringKey()
	curr3, ok := cl.GetCurrencyByStringKey(key)
	assert.True(t, ok)
	assert.Equal(t, curr3, curr)

	assert.Equal(t, cl.Len(), 1)

	curr4 := Currency{"Bitcoin", 1, 18}
	err = cl.Register(curr4)
	assert.NoError(t, err)

	assert.Equal(t, cl.Len(), 2)
}
