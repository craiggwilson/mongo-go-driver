package bson2

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocumentCodec(t *testing.T) {
	input, err := hex.DecodeString("160000000378000E0000000261000200000062000000")
	require.NoError(t, err)

	target := &Document{}

	err = Unmarshal(input, target)
	require.NoError(t, err)

	expected := &Document{
		elems: []*Element{&Element{&Value{
			start:  0,
			offset: 3,
			data:   input[4:21],
		}}},
		index: []uint32{0},
	}

	require.Equal(t, expected, target)
}
