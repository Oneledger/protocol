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

package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterator(t *testing.T) {
	items := []IterItem{
		{[]byte("a"), []byte("a123")},
		{[]byte("b"), []byte("a123")},
		{[]byte("c"), []byte("a123")},
		{[]byte("d"), []byte("a123")},
		{[]byte("e"), []byte("a123")},
		{[]byte("f"), []byte("a123")},
		{[]byte("g"), []byte("a123")},
		{[]byte("h"), []byte("a123")},
	}

	iter := newIterator(items)

	assert.Equal(t, []byte("a"), iter.Key())
	assert.Equal(t, []byte("a123"), iter.Value())
	assert.True(t, iter.Next())

	assert.Equal(t, []byte("b"), iter.Key())
	assert.Equal(t, []byte("a123"), iter.Value())
	assert.True(t, iter.Next())

	assert.Equal(t, []byte("c"), iter.Key())
	assert.Equal(t, []byte("a123"), iter.Value())
	assert.True(t, iter.Next())

	assert.Equal(t, []byte("d"), iter.Key())
	assert.Equal(t, []byte("a123"), iter.Value())
	assert.True(t, iter.Next())

	assert.Equal(t, []byte("e"), iter.Key())
	assert.Equal(t, []byte("a123"), iter.Value())
	assert.True(t, iter.Next())

	assert.Equal(t, []byte("f"), iter.Key())
	assert.Equal(t, []byte("a123"), iter.Value())
	assert.True(t, iter.Next())

	assert.Equal(t, []byte("g"), iter.Key())
	assert.Equal(t, []byte("a123"), iter.Value())
	assert.True(t, iter.Next())

	assert.Equal(t, []byte("h"), iter.Key())
	assert.Equal(t, []byte("a123"), iter.Value())
	assert.False(t, iter.Next())

}
