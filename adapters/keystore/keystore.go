package keystore

import (
	"crypto/ecdsa"
)

type Keystore interface {
	LoadECDSA() ecdsa.PrivateKey
	RandomECDSA() ecdsa.PrivateKey
	SaveECDSA(ecdsa.PrivateKey) error
}

type keyStore struct {
	PrivateKeys []ecdsa.PrivateKey `json:"privateKeys"`
}
