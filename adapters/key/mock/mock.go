package mock

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/atom-go/services/swap"
)

type mockKey struct {
	privateKey *ecdsa.PrivateKey
}

func NewMockKey() (swap.Key, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	return &mockKey{
		privateKey: privateKey,
	}, nil
}

func (key *mockKey) GetKey() *ecdsa.PrivateKey {
	return key.privateKey
}

func (key *mockKey) GetKeyString() (string, error) {
	return hex.EncodeToString(crypto.FromECDSA(key.privateKey)), nil
}

func (key *mockKey) GetAddress() ([]byte, error) {
	return crypto.Keccak256(crypto.FromECDSA(key.privateKey)), nil
}

func (key *mockKey) PriorityCode() uint32 {
	return 4294967295
}
