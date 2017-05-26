package auth_test

import (
	"context"
	"strings"
	"testing"

	"github.com/10gen/mongo-go-driver/bson"
	"github.com/craiggwilson/go-sasl"
	"github.com/craiggwilson/go-sasl/scramsha1"

	"reflect"

	"encoding/base64"

	. "github.com/10gen/mongo-go-driver/auth"
	"github.com/10gen/mongo-go-driver/internal/conntest"
	"github.com/10gen/mongo-go-driver/internal/msgtest"
	"github.com/10gen/mongo-go-driver/msg"
)

type constReader []byte

func (cr constReader) Read(dst []byte) (n int, err error) {
	copy(dst, []byte(cr))
	return len(dst), nil
}

func TestScramSHA1Authenticator_Fails(t *testing.T) {
	t.Parallel()

	authenticator := &SaslAuthenticator{
		MechName: scramsha1.MechName,
		Source:   "source",
		Factory: func() sasl.ClientMech {
			return scramsha1.NewClientMech("", "user", "1c33006ec1ffd90f9cadcbcc0e118200", 24, constReader("fyko+d2lbbFgONRv9qkxdawL"))
		},
	}

	saslStartReply := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", []byte{}},
		{"code", 143},
		{"done", true},
	})

	conn := &conntest.MockConnection{
		ResponseQ: []*msg.Reply{saslStartReply},
	}

	err := authenticator.Auth(context.Background(), conn)
	if err == nil {
		t.Fatalf("expected an error but got none")
	}

	errPrefix := "unable to authenticate using mechanism \"SCRAM-SHA-1\""
	if !strings.HasPrefix(err.Error(), errPrefix) {
		t.Fatalf("expected an err starting with \"%s\" but got \"%s\"", errPrefix, err)
	}
}

func TestScramSHA1Authenticator_Missing_challenge_fields(t *testing.T) {
	t.Parallel()

	authenticator := &SaslAuthenticator{
		MechName: scramsha1.MechName,
		Source:   "source",
		Factory: func() sasl.ClientMech {
			return scramsha1.NewClientMech("", "user", "1c33006ec1ffd90f9cadcbcc0e118200", 24, constReader("fyko+d2lbbFgONRv9qkxdawL"))
		},
	}

	payload, _ := base64.StdEncoding.DecodeString("cz1yUTlaWTNNbnRCZXVQM0UxVERWQzR3PT0saT0xMDAwMA===")
	saslStartReply := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", payload},
		{"done", false},
	})

	conn := &conntest.MockConnection{
		ResponseQ: []*msg.Reply{saslStartReply},
	}

	err := authenticator.Auth(context.Background(), conn)
	if err == nil {
		t.Fatalf("expected an error but got none")
	}

	errPrefix := "sasl mechanism SCRAM-SHA-1: client failed to provide response: invalid challenge"
	if !strings.HasPrefix(err.Error(), errPrefix) {
		t.Fatalf("expected an err starting with \"%s\" but got \"%s\"", errPrefix, err)
	}
}

func TestScramSHA1Authenticator_Invalid_server_nonce1(t *testing.T) {
	t.Parallel()

	authenticator := &SaslAuthenticator{
		MechName: scramsha1.MechName,
		Source:   "source",
		Factory: func() sasl.ClientMech {
			return scramsha1.NewClientMech("", "user", "1c33006ec1ffd90f9cadcbcc0e118200", 24, constReader("fyko+d2lbbFgONRv9qkxdawL"))
		},
	}

	payload, _ := base64.StdEncoding.DecodeString("bD0yMzJnLHM9clE5WlkzTW50QmV1UDNFMVREVkM0dz09LGk9MTAwMDA=")
	saslStartReply := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", payload},
		{"done", false},
	})

	conn := &conntest.MockConnection{
		ResponseQ: []*msg.Reply{saslStartReply},
	}

	err := authenticator.Auth(context.Background(), conn)
	if err == nil {
		t.Fatalf("expected an error but got none")
	}

	errPrefix := "sasl mechanism SCRAM-SHA-1: client failed to provide response: invalid challenge: expected nonce"
	if !strings.HasPrefix(err.Error(), errPrefix) {
		t.Fatalf("expected an err starting with \"%s\" but got \"%s\"", errPrefix, err)
	}
}

func TestScramSHA1Authenticator_Invalid_server_nonce2(t *testing.T) {
	t.Parallel()

	authenticator := &SaslAuthenticator{
		MechName: scramsha1.MechName,
		Source:   "source",
		Factory: func() sasl.ClientMech {
			return scramsha1.NewClientMech("", "user", "1c33006ec1ffd90f9cadcbcc0e118200", 24, constReader("fyko+d2lbbFgONRv9qkxdawL"))
		},
	}

	payload, _ := base64.StdEncoding.DecodeString("cj1meWtvLWQybGJiRmdPTlJ2OXFreGRhd0xIbytWZ2s3cXZVT0tVd3VXTElXZzRsLzlTcmFHTUhFRSxzPXJROVpZM01udEJldVAzRTFURFZDNHc9PSxpPTEwMDAw")
	saslStartReply := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", payload},
		{"done", false},
	})

	conn := &conntest.MockConnection{
		ResponseQ: []*msg.Reply{saslStartReply},
	}

	err := authenticator.Auth(context.Background(), conn)
	if err == nil {
		t.Fatalf("expected an error but got none")
	}

	errPrefix := "sasl mechanism SCRAM-SHA-1: client failed to provide response: invalid challenge: nonce mismatch"
	if !strings.HasPrefix(err.Error(), errPrefix) {
		t.Fatalf("expected an err starting with \"%s\" but got \"%s\"", errPrefix, err)
	}
}

func TestScramSHA1Authenticator_No_salt(t *testing.T) {
	t.Parallel()

	authenticator := &SaslAuthenticator{
		MechName: scramsha1.MechName,
		Source:   "source",
		Factory: func() sasl.ClientMech {
			return scramsha1.NewClientMech("", "user", "1c33006ec1ffd90f9cadcbcc0e118200", 24, constReader("fyko+d2lbbFgONRv9qkxdawL"))
		},
	}

	payload, _ := base64.StdEncoding.DecodeString("cj1meWtvK2QybGJiRmdPTlJ2OXFreGRhd0xIbytWZ2s3cXZVT0tVd3VXTElXZzRsLzlTcmFHTUhFRSxrPXJROVpZM01udEJldVAzRTFURFZDNHc9PSxpPTEwMDAw======")
	saslStartReply := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", payload},
		{"done", false},
	})

	conn := &conntest.MockConnection{
		ResponseQ: []*msg.Reply{saslStartReply},
	}

	err := authenticator.Auth(context.Background(), conn)
	if err == nil {
		t.Fatalf("expected an error but got none")
	}

	errPrefix := "sasl mechanism SCRAM-SHA-1: client failed to provide response: invalid challenge: expected salt"
	if !strings.HasPrefix(err.Error(), errPrefix) {
		t.Fatalf("expected an err starting with \"%s\" but got \"%s\"", errPrefix, err)
	}
}

func TestScramSHA1Authenticator_No_iteration_count(t *testing.T) {
	t.Parallel()

	authenticator := &SaslAuthenticator{
		MechName: scramsha1.MechName,
		Source:   "source",
		Factory: func() sasl.ClientMech {
			return scramsha1.NewClientMech("", "user", "1c33006ec1ffd90f9cadcbcc0e118200", 24, constReader("fyko+d2lbbFgONRv9qkxdawL"))
		},
	}

	payload, _ := base64.StdEncoding.DecodeString("cj1meWtvK2QybGJiRmdPTlJ2OXFreGRhd0xIbytWZ2s3cXZVT0tVd3VXTElXZzRsLzlTcmFHTUhFRSxzPXJROVpZM01udEJldVAzRTFURFZDNHc9PSxrPXNkZg======")
	saslStartReply := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", payload},
		{"done", false},
	})

	conn := &conntest.MockConnection{
		ResponseQ: []*msg.Reply{saslStartReply},
	}

	err := authenticator.Auth(context.Background(), conn)
	if err == nil {
		t.Fatalf("expected an error but got none")
	}

	errPrefix := "sasl mechanism SCRAM-SHA-1: client failed to provide response: invalid challenge: expected iteration-count"
	if !strings.HasPrefix(err.Error(), errPrefix) {
		t.Fatalf("expected an err starting with \"%s\" but got \"%s\"", errPrefix, err)
	}
}

func TestScramSHA1Authenticator_Invalid_iteration_count(t *testing.T) {
	t.Parallel()

	authenticator := &SaslAuthenticator{
		MechName: scramsha1.MechName,
		Source:   "source",
		Factory: func() sasl.ClientMech {
			return scramsha1.NewClientMech("", "user", "1c33006ec1ffd90f9cadcbcc0e118200", 24, constReader("fyko+d2lbbFgONRv9qkxdawL"))
		},
	}

	payload, _ := base64.StdEncoding.DecodeString("cj1meWtvK2QybGJiRmdPTlJ2OXFreGRhd0xIbytWZ2s3cXZVT0tVd3VXTElXZzRsLzlTcmFHTUhFRSxzPXJROVpZM01udEJldVAzRTFURFZDNHc9PSxpPWFiYw====")
	saslStartReply := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", payload},
		{"done", false},
	})

	conn := &conntest.MockConnection{
		ResponseQ: []*msg.Reply{saslStartReply},
	}

	err := authenticator.Auth(context.Background(), conn)
	if err == nil {
		t.Fatalf("expected an error but got none")
	}

	errPrefix := "sasl mechanism SCRAM-SHA-1: client failed to provide response: invalid challenge: invalid iteration-count"
	if !strings.HasPrefix(err.Error(), errPrefix) {
		t.Fatalf("expected an err starting with \"%s\" but got \"%s\"", errPrefix, err)
	}
}

func TestScramSHA1Authenticator_Invalid_server_signature(t *testing.T) {
	t.Parallel()

	authenticator := &SaslAuthenticator{
		MechName: scramsha1.MechName,
		Source:   "source",
		Factory: func() sasl.ClientMech {
			return scramsha1.NewClientMech("", "user", "1c33006ec1ffd90f9cadcbcc0e118200", 24, constReader("fyko+d2lbbFgONRv9qkxdawL"))
		},
	}

	payload, _ := base64.StdEncoding.DecodeString("cj1meWtvK2QybGJiRmdPTlJ2OXFreGRhd0xIbytWZ2s3cXZVT0tVd3VXTElXZzRsLzlTcmFHTUhFRSxzPXJROVpZM01udEJldVAzRTFURFZDNHc9PSxpPTEwMDAw")
	saslStartReply := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", payload},
		{"done", false},
	})
	payload, _ = base64.StdEncoding.DecodeString("dj1VTVdlSTI1SkQxeU5ZWlJNcFo0Vkh2aFo5ZTBh")
	saslContinueReply := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", payload},
		{"done", false},
	})

	conn := &conntest.MockConnection{
		ResponseQ: []*msg.Reply{saslStartReply, saslContinueReply},
	}

	err := authenticator.Auth(context.Background(), conn)
	if err == nil {
		t.Fatalf("expected an error but got none")
	}

	errPrefix := "sasl mechanism SCRAM-SHA-1: client failed to provide response: invalid challenge: server signature mismatch"
	if !strings.HasPrefix(err.Error(), errPrefix) {
		t.Fatalf("expected an err starting with \"%s\" but got \"%s\"", errPrefix, err)
	}
}

func TestScramSHA1Authenticator_Server_provided_error(t *testing.T) {
	t.Parallel()

	authenticator := &SaslAuthenticator{
		MechName: scramsha1.MechName,
		Source:   "source",
		Factory: func() sasl.ClientMech {
			return scramsha1.NewClientMech("", "user", "1c33006ec1ffd90f9cadcbcc0e118200", 24, constReader("fyko+d2lbbFgONRv9qkxdawL"))
		},
	}

	payload, _ := base64.StdEncoding.DecodeString("cj1meWtvK2QybGJiRmdPTlJ2OXFreGRhd0xIbytWZ2s3cXZVT0tVd3VXTElXZzRsLzlTcmFHTUhFRSxzPXJROVpZM01udEJldVAzRTFURFZDNHc9PSxpPTEwMDAw")
	saslStartReply := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", payload},
		{"done", false},
	})
	payload, _ = base64.StdEncoding.DecodeString("ZT1zZXJ2ZXIgcGFzc2VkIGVycm9y")
	saslContinueReply := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", payload},
		{"done", false},
	})

	conn := &conntest.MockConnection{
		ResponseQ: []*msg.Reply{saslStartReply, saslContinueReply},
	}

	err := authenticator.Auth(context.Background(), conn)
	if err == nil {
		t.Fatalf("expected an error but got none")
	}

	errPrefix := "sasl mechanism SCRAM-SHA-1: client failed to provide response: server passed error"
	if !strings.HasPrefix(err.Error(), errPrefix) {
		t.Fatalf("expected an err starting with \"%s\" but got \"%s\"", errPrefix, err)
	}
}

func TestScramSHA1Authenticator_Invalid_final_message(t *testing.T) {
	t.Parallel()

	authenticator := &SaslAuthenticator{
		MechName: scramsha1.MechName,
		Source:   "source",
		Factory: func() sasl.ClientMech {
			return scramsha1.NewClientMech("", "user", "1c33006ec1ffd90f9cadcbcc0e118200", 24, constReader("fyko+d2lbbFgONRv9qkxdawL"))
		},
	}

	payload, _ := base64.StdEncoding.DecodeString("cj1meWtvK2QybGJiRmdPTlJ2OXFreGRhd0xIbytWZ2s3cXZVT0tVd3VXTElXZzRsLzlTcmFHTUhFRSxzPXJROVpZM01udEJldVAzRTFURFZDNHc9PSxpPTEwMDAw")
	saslStartReply := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", payload},
		{"done", false},
	})
	payload, _ = base64.StdEncoding.DecodeString("Zj1VTVdlSTI1SkQxeU5ZWlJNcFo0Vkh2aFo5ZTBh")
	saslContinueReply := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", payload},
		{"done", false},
	})

	conn := &conntest.MockConnection{
		ResponseQ: []*msg.Reply{saslStartReply, saslContinueReply},
	}

	err := authenticator.Auth(context.Background(), conn)
	if err == nil {
		t.Fatalf("expected an error but got none")
	}

	errPrefix := "sasl mechanism SCRAM-SHA-1: client failed to provide response: invalid challenge: expected server signature"
	if !strings.HasPrefix(err.Error(), errPrefix) {
		t.Fatalf("expected an err starting with \"%s\" but got \"%s\"", errPrefix, err)
	}
}

func TestScramSHA1Authenticator_Extra_message(t *testing.T) {
	t.Parallel()

	authenticator := &SaslAuthenticator{
		MechName: scramsha1.MechName,
		Source:   "source",
		Factory: func() sasl.ClientMech {
			return scramsha1.NewClientMech("", "user", "1c33006ec1ffd90f9cadcbcc0e118200", 24, constReader("fyko+d2lbbFgONRv9qkxdawL"))
		},
	}

	payload, _ := base64.StdEncoding.DecodeString("cj1meWtvK2QybGJiRmdPTlJ2OXFreGRhd0xIbytWZ2s3cXZVT0tVd3VXTElXZzRsLzlTcmFHTUhFRSxzPXJROVpZM01udEJldVAzRTFURFZDNHc9PSxpPTEwMDAw")
	saslStartReply := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", payload},
		{"done", false},
	})
	payload, _ = base64.StdEncoding.DecodeString("dj1VTVdlSTI1SkQxeU5ZWlJNcFo0Vkh2aFo5ZTA9")
	saslContinueReply := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", payload},
		{"done", false},
	})
	saslContinueReply2 := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", []byte{}},
		{"done", false},
	})

	conn := &conntest.MockConnection{
		ResponseQ: []*msg.Reply{saslStartReply, saslContinueReply, saslContinueReply2},
	}

	err := authenticator.Auth(context.Background(), conn)
	if err == nil {
		t.Fatalf("expected an error but got none")
	}

	errPrefix := "unexpected server challenge"
	if !strings.HasPrefix(err.Error(), errPrefix) {
		t.Fatalf("expected an err starting with \"%s\" but got \"%s\"", errPrefix, err)
	}
}

func TestScramSHA1Authenticator_Succeeds(t *testing.T) {
	t.Parallel()

	authenticator := &SaslAuthenticator{
		MechName: scramsha1.MechName,
		Source:   "source",
		Username: "user",
		Password: "pencil",
		Factory: func() sasl.ClientMech {
			return scramsha1.NewClientMech("", "user", "1c33006ec1ffd90f9cadcbcc0e118200", 24, constReader("fyko+d2lbbFgONRv9qkxdawL"))
		},
	}

	payload, _ := base64.StdEncoding.DecodeString("cj1meWtvK2QybGJiRmdPTlJ2OXFreGRhd0xIbytWZ2s3cXZVT0tVd3VXTElXZzRsLzlTcmFHTUhFRSxzPXJROVpZM01udEJldVAzRTFURFZDNHc9PSxpPTEwMDAw")
	saslStartReply := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", payload},
		{"done", false},
	})
	payload, _ = base64.StdEncoding.DecodeString("dj1VTVdlSTI1SkQxeU5ZWlJNcFo0Vkh2aFo5ZTA9")
	saslContinueReply := msgtest.CreateCommandReply(bson.D{
		{"ok", 1},
		{"conversationId", 1},
		{"payload", payload},
		{"done", true},
	})

	conn := &conntest.MockConnection{
		ResponseQ: []*msg.Reply{saslStartReply, saslContinueReply},
	}

	err := authenticator.Auth(context.Background(), conn)
	if err != nil {
		t.Fatalf("expected no error but got \"%s\"", err)
	}

	if len(conn.Sent) != 2 {
		t.Fatalf("expected 2 messages to be sent but had %d", len(conn.Sent))
	}

	saslStartRequest := conn.Sent[0].(*msg.Query)
	payload, _ = base64.RawStdEncoding.DecodeString("biwsbj11c2VyLHI9ZnlrbytkMmxiYkZnT05Sdjlxa3hkYXdM")
	expectedCmd := bson.D{
		{"saslStart", 1},
		{"mechanism", "SCRAM-SHA-1"},
		{"payload", payload},
	}
	if !reflect.DeepEqual(saslStartRequest.Query, expectedCmd) {
		t.Fatalf("saslStart command was incorrect:\n  expected: %v\n    actual: %v", expectedCmd, saslStartRequest.Query)
	}

	saslContinueRequest := conn.Sent[1].(*msg.Query)
	payload, _ = base64.RawStdEncoding.DecodeString("Yz1iaXdzLHI9ZnlrbytkMmxiYkZnT05Sdjlxa3hkYXdMSG8rVmdrN3F2VU9LVXd1V0xJV2c0bC85U3JhR01IRUUscD1NQzJUOEJ2Ym1XUmNrRHc4b1dsNUlWZ2h3Q1k9")
	expectedCmd = bson.D{
		{"saslContinue", 1},
		{"conversationId", 1},
		{"payload", payload},
	}
	if !reflect.DeepEqual(saslContinueRequest.Query, expectedCmd) {
		t.Fatalf("saslContinue command was incorrect:\n  expected: %v\n    actual: %v", expectedCmd, saslContinueRequest.Query)
	}
}
