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
	cl := NewCurrencySet()

	curr := Currency{Id: 0, Name: "OLT", Chain: 0, Decimal: 18, Unit: "nue"}
	err := cl.Register(curr)
	assert.NoError(t, err)

	currNew, ok := cl.GetCurrencyByName("OLT")
	assert.True(t, ok)

	assert.Equal(t, currNew, curr)

	id := curr.Id
	curr3, ok := cl.GetCurrencyById(id)
	assert.True(t, ok)
	assert.Equal(t, curr3, curr)

	assert.Equal(t, cl.Len(), 1)

	curr4 := Currency{Id: 1, Name: "Bitcoin", Chain: 1, Decimal: 18, Unit: "satoshi"}
	err = cl.Register(curr4)
	assert.NoError(t, err)

	assert.Equal(t, cl.Len(), 2)
}
