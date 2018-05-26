package bson2

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRawDCodec(t *testing.T) {
	input, err := hex.DecodeString("160000000378000E0000000261000200000062000000")
	require.NoError(t, err)

	target := RawD{}

	err = Unmarshal(input, &target)
	require.NoError(t, err)

	expected := RawD{
		{"x", Raw{Kind: 3, Data: []byte{0xe, 0x0, 0x0, 0x0, 0x2, 0x61, 0x0, 0x2, 0x0, 0x0, 0x0, 0x62, 0x0, 0x0}, reg: globalRegistry}},
	}

	require.Equal(t, expected, target)
}
