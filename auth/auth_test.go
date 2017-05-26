package auth_test

import (
	"testing"

	. "github.com/10gen/mongo-go-driver/auth"
	"github.com/stretchr/testify/require"
)

func TestCreateAuthenticator(t *testing.T) {

	tests := []struct {
		name   string
		source string
		auther Authenticator
	}{
		{name: "", auther: &DefaultAuthenticator{}},
		{name: "SCRAM-SHA-1", auther: &SaslAuthenticator{}},
		{name: "MONGODB-CR", auther: &MongoDBCRAuthenticator{}},
		{name: "PLAIN", auther: &SaslAuthenticator{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cred := &Cred{
				Username:    "user",
				Password:    "pencil",
				PasswordSet: true,
			}

			a, err := CreateAuthenticator(test.name, cred)
			require.NoError(t, err)
			require.IsType(t, test.auther, a)
		})
	}
}
