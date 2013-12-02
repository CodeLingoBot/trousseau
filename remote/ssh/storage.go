package ssh

import (
    "strings"
	"code.google.com/p/go.crypto/ssh"
)

type ScpStorage struct {
	host      string
	port      string
	connexion *ssh.ClientConn

	Keychain *Keychain
	Password string
	User     string
	Endpoint string
}

func NewScpStorage(host, port, user, password string, keychain *Keychain) *ScpStorage {
	return &ScpStorage{
		Keychain: keychain,
        Password: password,
		User:     user,
		Endpoint: strings.Join([]string{host, port}, ":"),
		host:     host,
		port:     port,
	}
}
