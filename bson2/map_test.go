package bson2

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapCodec(t *testing.T) {
	input, err := hex.DecodeString("160000000378000E0000000261000200000062000000")
	require.NoError(t, err)

	target := make(map[string]interface{})

	err = Unmarshal(input, &target)
	require.NoError(t, err)

	expected := map[string]interface{}{
		"x": map[string]interface{}{
			"a": "b",
		},
	}

	require.Equal(t, expected, target)
}
