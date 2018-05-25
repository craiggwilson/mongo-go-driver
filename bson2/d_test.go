package bson2

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDCodec(t *testing.T) {
	b, err := hex.DecodeString("160000000378000E0000000261000200000062000000")
	require.NoError(t, err)

	target := D{}

	err = Unmarshal(bytes.NewReader(b), &target)
	require.NoError(t, err)

	expected := D{
		{"x", D{
			{"a", "b"},
		}}}

	require.Equal(t, expected, target)
}
