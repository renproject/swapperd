package auth

import (
	"crypto/subtle"

	"github.com/ethereum/go-ethereum/crypto/sha3"
)

type Authenticator interface {
	VerifyPassword(password string) bool
}

type authenticator struct {
	passwordHash [32]byte
}

func NewAuthenticator(passwordHash [32]byte) Authenticator {
	return &authenticator{passwordHash}
}

func (authenticator *authenticator) VerifyPassword(password string) bool {
	passwordHash := sha3.Sum256([]byte(password))
	passwordOk := subtle.ConstantTimeCompare(passwordHash[:], authenticator.passwordHash[:]) == 1
	return passwordOk
}
