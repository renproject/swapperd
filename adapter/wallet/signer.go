package wallet

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

type ECDSASigner interface {
	PublicKey() ecdsa.PublicKey
	Sign(hash []byte) ([]byte, error)
}

type ecdsaSigner struct {
	privateKey *ecdsa.PrivateKey
}

func (signer *ecdsaSigner) Sign(hash []byte) ([]byte, error) {

	return crypto.Sign(hash, signer.privateKey)
}

func (signer *ecdsaSigner) PublicKey() ecdsa.PublicKey {
	return signer.privateKey.PublicKey
}
