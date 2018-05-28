package bson2

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDCodec(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		input, err := hex.DecodeString("160000000378000E0000000261000200000062000000")
		require.NoError(t, err)

		target := D{}

		err = Unmarshal(input, &target)
		require.NoError(t, err)

		expected := D{
			{"x", D{
				{"a", "b"},
			}}}

		require.Equal(t, expected, target)
	})
	t.Run("Encode", func(t *testing.T) {
		input := D{
			{"x", D{
				{"a", "b"},
			}}}

		var target bytes.Buffer

		err := Marshal(input, &target)
		require.NoError(t, err)

		expected, err := hex.DecodeString("160000000378000E0000000261000200000062000000")
		require.NoError(t, err)

		require.Equal(t, expected, target.Bytes())
	})

}
