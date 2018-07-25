package keystore

import (
	"crypto/ecdsa"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/crypto"
)

type EthereumKey struct {
	privateKey   *ecdsa.PrivateKey `json:"private_key"`
	priorityCode uint32            `json:"priority_code"`
	network      string            `json:"network"`
}

func GetEthereumKey(key string, network string) (EthereumKey, error) {
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
	return key.privateKey
}

func (key *EthereumKey) GetKeyString() (string, error) {
	return hex.EncodeToString(crypto.FromECDSA(key.privateKey)), nil
}

func (key *EthereumKey) GetAddress() ([]byte, error) {
	return bind.NewKeyedTransactor(key.privateKey).From.Bytes(), nil
}

func (key *EthereumKey) PriorityCode() uint32 {
	return 1
}
