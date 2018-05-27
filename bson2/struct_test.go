package bson2

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStructCodec(t *testing.T) {
	input, err := hex.DecodeString("160000000378000E0000000261000200000062000000")
	require.NoError(t, err)

	target := testStruct{}

	err = Unmarshal(input, &target)
	require.NoError(t, err)

	expected := testStruct{
		X: testStructEmbedded{A: "b"},
	}

	require.Equal(t, expected, target)
}

type testStruct struct {
	X testStructEmbedded
}

type testStructEmbedded struct {
	A string
}
