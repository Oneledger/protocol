package btcrpc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlockUnmarshal(t *testing.T) {
	block := new(Block)
	err := json.Unmarshal([]byte("222"), block)
	require.NotNil(t, err)
}