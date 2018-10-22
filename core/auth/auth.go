package auth

import (
	"crypto/subtle"

	"github.com/ethereum/go-ethereum/crypto/sha3"
)

type Authenticator interface {
	VerifyUsernameAndPassword(username, password string) bool
}

type authenticator struct {
	username     string
	passwordHash [32]byte
}

func NewAuthenticator(username string, passwordHash [32]byte) Authenticator {
	return &authenticator{username, passwordHash}
}

func (authenticator *authenticator) VerifyUsernameAndPassword(username string, password string) bool {
	passwordHash := sha3.Sum256([]byte(password))
	usernameOk := subtle.ConstantTimeCompare([]byte(username), []byte(authenticator.username)) == 1
	passwordOk := subtle.ConstantTimeCompare(passwordHash[:], authenticator.passwordHash[:]) == 1
	return usernameOk && passwordOk
}
