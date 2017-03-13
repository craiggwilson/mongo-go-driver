package msg_test

import (
	"testing"

	"encoding/hex"

	"github.com/10gen/mongo-go-driver/bson"
	. "github.com/10gen/mongo-go-driver/msg"
	"github.com/stretchr/testify/require"
)

func TestWrapWithMeta(t *testing.T) {
	req := NewCommand(10, "admin", true, bson.M{"a": 1}).(*Query)

	buf, err := bson.Marshal(req.Query)
	require.NoError(t, err)
	require.Equal(t, "0c0000001061000100000000", hex.EncodeToString(buf))

	WrapWithMeta(req, map[string]interface{}{
		"$readPreference": bson.M{
			"mode": "secondary",
		},
	})

	buf, err = bson.Marshal(req.Query)
	require.NoError(t, err)
	require.Equal(t, `43000000032472656164507265666572656e63650019000000026d6f6465000a0000007365636f6e64617279000003247175657279000c000000106100010000000000`, hex.EncodeToString(buf))
}
