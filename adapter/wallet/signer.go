package wallet

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

type Signer interface {
	PublicKey() []byte
	Sign(hash []byte) ([]byte, error)
}

type ecdsaSigner struct {
	privateKey *ecdsa.PrivateKey
}

func (signer *ecdsaSigner) Sign(hash []byte) ([]byte, error) {
	return crypto.Sign(hash, signer.privateKey)
}

func (signer *ecdsaSigner) PublicKey() []byte {
	return crypto.CompressPubkey(&signer.privateKey.PublicKey)
}
