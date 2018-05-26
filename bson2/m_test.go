package bson2

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMCodec(t *testing.T) {
	input, err := hex.DecodeString("160000000378000E0000000261000200000062000000")
	require.NoError(t, err)

	target := M{}

	err = Unmarshal(input, &target)
	require.NoError(t, err)

	expected := M{
		"x": M{
			"a": "b",
		},
	}

	require.Equal(t, expected, target)
}
