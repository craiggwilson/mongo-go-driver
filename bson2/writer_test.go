package bson2

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriter(t *testing.T) {
	w := NewValueWriter()

	dw, err := w.WriteDocument()
	require.NoError(t, err)

	vr, err := dw.WriteElement("x")
	require.NoError(t, err)

	xDoc, err := vr.WriteDocument()
	require.NoError(t, err)

	xvr, err := xDoc.WriteElement("a")
	require.NoError(t, err)
	xvr.WriteString("b")

	err = xDoc.WriteEndDocument()
	require.NoError(t, err)

	err = dw.WriteEndDocument()
	require.NoError(t, err)

	var actual bytes.Buffer
	_, err = w.WriteTo(&actual)
	require.NoError(t, err)

	expected, err := hex.DecodeString("160000000378000E0000000261000200000062000000")
	require.NoError(t, err)

	require.Equal(t, expected, actual.Bytes())
}
