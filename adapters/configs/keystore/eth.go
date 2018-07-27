package keystore

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/crypto"
)

type EthereumKey struct {
	PrivateKey   *ecdsa.PrivateKey `json:"private_key"`
	CurrencyCode uint32            `json:"currency_code"`
	Network      string            `json:"network"`
}

func NewEthereumKey(key string, network string) (EthereumKey, error) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return EthereumKey{}, err
	}
	return EthereumKey{
		privateKey,
		1,
		network,
	}, nil
}

func (key *EthereumKey) GetKey() *ecdsa.PrivateKey {
	return key.PrivateKey
}

func (key *EthereumKey) GetKeyString() (string, error) {
	return hex.EncodeToString(crypto.FromECDSA(key.PrivateKey)), nil
}

func (key *EthereumKey) GetAddress() ([]byte, error) {
	return bind.NewKeyedTransactor(key.PrivateKey).From.Bytes(), nil
}

func (key *EthereumKey) PriorityCode() uint32 {
	return key.CurrencyCode
}
