package eth

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/atom-go/services/swap"
)

type ethereumKey struct {
	privateKey *ecdsa.PrivateKey
	network    string
}

func NewEthereumKey(key string, network string) (swap.Key, error) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, err
	}
	return &ethereumKey{
		privateKey: privateKey,
		network:    network,
	}, nil
}

func (key *ethereumKey) GetKey() *ecdsa.PrivateKey {
	return key.privateKey
}

func (key *ethereumKey) GetKeyString() (string, error) {
	return hex.EncodeToString(crypto.FromECDSA(key.privateKey)), nil
}

func (key *ethereumKey) GetAddress() ([]byte, error) {
	return bind.NewKeyedTransactor(key.privateKey).From.Bytes(), nil
}

func (key *ethereumKey) PriorityCode() uint32 {
	return 1
}
