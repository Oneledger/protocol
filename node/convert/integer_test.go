package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert_GetInt(t *testing.T) {
	lookup := map[string]int{
		"1":       1,
		"333":     333,
		"-20":     -20,
		"--09989": 0,
		"-0987":   -987,
	}

	for s, i := range lookup {
		assert.Equal(t, i, GetInt(s, 0), "error converting string: "+s)
	}

	assert.Equal(t, 0, GetInt("asdf", 0))
	assert.Equal(t, 666, GetInt("oiut", 666))

}

func TestGetString(t *testing.T) {
	lookup := map[int]string{
		1:        "1",
		786:      "786",
		335:      "335",
		12341234: "12341234",
	}

	for i, s := range lookup {
		assert.Equal(t, s, GetString(i))
	}
}
