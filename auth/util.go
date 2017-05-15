package auth

import (
	"crypto/md5"
	"fmt"
	"io"
	"net"

	"github.com/10gen/mongo-go-driver/model"
)

const defaultAuthDB = "admin"

func mongoPasswordDigest(username, password string) string {
	h := md5.New()
	io.WriteString(h, username)
	io.WriteString(h, ":mongo:")
	io.WriteString(h, password)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Hostname returns just the hostname of the address.
func getHostname(a model.Addr) string {
	if a.Network() != "unix" {
		if host, _, err := net.SplitHostPort(string(a)); err == nil {
			return host
		}
	}
	return string(a)
}
