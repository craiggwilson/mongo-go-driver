package bson2

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMCodec(t *testing.T) {
	b, err := hex.DecodeString("160000000378000E0000000261000200000062000000")
	require.NoError(t, err)

	target := M{}

	err = Unmarshal(bytes.NewReader(b), &target)
	require.NoError(t, err)

	expected := M{
		"x": M{
			"a": "b",
		},
	}

	require.Equal(t, expected, target)
}
