package keystore

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/crypto"
)

type ethereumKey key

func (key *ethereumKey) GetKey() (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(key.PrivateKey)
}

func (key *ethereumKey) GetKeyString() string {
	return key.PrivateKey
}

func (key *ethereumKey) GetAddress() ([]byte, error) {
	privKey, err := key.GetKey()
	if err != nil {
		return nil, err
	}
	return bind.NewKeyedTransactor(privKey).From.Bytes(), nil
}

func (key *ethereumKey) PriorityCode() uint32 {
	return key.Code
}

func (key *ethereumKey) Chain() string {
	return key.Network
}

func RandomEthereumKeyString() (string, error) {
	priv, err := crypto.GenerateKey()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(crypto.FromECDSA(priv)), nil
}
