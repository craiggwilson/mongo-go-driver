package auth

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/10gen/mongo-go-driver/bson"
	"github.com/10gen/mongo-go-driver/conn"
	"github.com/10gen/mongo-go-driver/msg"
	"github.com/craiggwilson/go-sasl"
	"github.com/craiggwilson/go-sasl/gssapi"
	"github.com/craiggwilson/go-sasl/plain"
	"github.com/craiggwilson/go-sasl/scramsha1"
)

func newGSSAPIAuthenticator(cred *Cred) (Authenticator, error) {
	return &SaslAuthenticator{
		MechName: gssapi.MechName,
		Source:   cred.Source,
		Factory: func(c conn.Connection) sasl.ClientMech {
			cfg := &gssapi.ClientConfig{
				Username:    cred.Username,
				Password:    cred.Password,
				ServiceName: "mongodb",
				ServiceFQDN: getHostname(c.Model().Addr),
			}

			return gssapi.NewClientMech(cfg)
		},
	}, nil
}

func newPlainAuthenticator(cred *Cred) (Authenticator, error) {
	return &SaslAuthenticator{
		MechName: plain.MechName,
		Source:   cred.Source,
		Factory: func(_ conn.Connection) sasl.ClientMech {
			return plain.NewClientMech("", cred.Username, cred.Password)
		},
	}, nil
}

func newScramSHA1Authenticator(cred *Cred) (Authenticator, error) {
	return &SaslAuthenticator{
		MechName: scramsha1.MechName,
		Source:   cred.Source,
		Factory: func(_ conn.Connection) sasl.ClientMech {
			return scramsha1.NewClientMech("", cred.Username, mongoPasswordDigest(cred.Username, cred.Password), 16, rand.Reader)
		},
	}, nil
}

// ClientMechFactory creates a sasl.ClientMech.
type ClientMechFactory func(conn.Connection) sasl.ClientMech

// SaslAuthenticator uses the SASL protocol to authenticate a connection using a provided mechanism.
type SaslAuthenticator struct {
	MechName string
	Factory  ClientMechFactory
	Source   string
}

// Auth authenticates the connection.
func (a *SaslAuthenticator) Auth(ctx context.Context, c conn.Connection) error {

	incoming := make(chan []byte, 1)
	outgoing := make(chan []byte, 1)
	errs := make(chan error, 1)

	mech := a.Factory(c)
	if closer, ok := mech.(sasl.MechCloser); ok {
		defer closer.Close()
	}
	go func() {
		errs <- sasl.ConverseAsClient(ctx, mech, incoming, outgoing)
	}()

	var payload []byte
	var err error
	var cid int
	for {
		select {
		case payload = <-outgoing:
		case err = <-errs:
			return err
		}

		var cmd msg.Request
		if cid == 0 {
			cmd = msg.NewCommand(
				msg.NextRequestID(),
				a.Source,
				true,
				bson.D{
					{"saslStart", 1},
					{"mechanism", a.MechName},
					{"payload", payload},
				},
			)
		} else {
			cmd = msg.NewCommand(
				msg.NextRequestID(),
				a.Source,
				true,
				bson.D{
					{"saslContinue", 1},
					{"conversationId", cid},
					{"payload", payload},
				},
			)
		}

		saslResp := struct {
			ConversationID int    `bson:"conversationId"`
			Code           int    `bson:"code"`
			Done           bool   `bson:"done"`
			Payload        []byte `bson:"payload"`
		}{}

		err = conn.ExecuteCommand(ctx, c, cmd, &saslResp)
		if err != nil {
			return err
		}

		cid = saslResp.ConversationID
		if saslResp.Code != 0 {
			return newError(fmt.Errorf("server failed with code %d", saslResp.Code), a.MechName)
		}

		incoming <- saslResp.Payload
	}
}
