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

package utils_test

import (
	"testing"

	"github.com/Oneledger/protocol/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetPort(t *testing.T) {
	port, err := utils.GetPort("https://1.1.1.1:789")
	assert.Equal(t, port, "789")
	assert.NoError(t, err)

	port, err = utils.GetPort("http://google.com:1234")
	assert.Equal(t, port, "1234")
	assert.NoError(t, err)

	port, err = utils.GetPort("google.com:1234")
	assert.Equal(t, port, "")
	assert.Equal(t, err, utils.ErrParsingAddress)
}
