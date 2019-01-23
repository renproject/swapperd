package wallet

import (
	"crypto/ecdsa"
	"encoding/base64"
	"strings"

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

func (wallet *wallet) ID(password, idType string) (string, error) {
	signer, err := wallet.ECDSASigner(password)
	if err != nil {
		return "", err
	}
	pubKey := signer.PublicKey()
	idType = strings.ToLower(idType)
	switch idType {
	case "ethereum", "eth":
		return crypto.PubkeyToAddress(pubKey).Hex(), nil
	default:
		return base64.StdEncoding.EncodeToString(crypto.FromECDSAPub(&pubKey)), nil
	}
}
