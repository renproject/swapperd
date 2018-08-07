package keystore

import (
	"crypto/ecdsa"
	"fmt"
)

type key struct {
	PrivateKey string `json:"privateKey"`
	Code       uint32 `json:"priorityCode"`
	Network    string `json:"network"`
}
type Key interface {
	GetKey() (*ecdsa.PrivateKey, error)
	GetKeyString() string
	GetAddress() ([]byte, error)
	PriorityCode() uint32
	Chain() string
}

func NewKey(privKey string, priCode uint32, network string) (Key, error) {
	key := key{
		privKey,
		priCode,
		network,
	}
	switch priCode {
	case 0:
		btcKey := bitcoinKey(key)
		return &btcKey, nil
	case 1:
		ethKey := ethereumKey(key)
		return &ethKey, nil
	}
	return nil, fmt.Errorf(ErrPrefix, "Unknown Priority Code")
}
