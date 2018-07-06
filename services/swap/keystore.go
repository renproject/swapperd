package swap

import (
	"crypto/ecdsa"
	"errors"
)

type Keystore interface {
	LoadKeys() ([]Key, error)
}

type Key interface {
	GetKey() *ecdsa.PrivateKey
	GetKeyString() (string, error)
	GetAddress() ([]byte, error)
	PriorityCode() uint32
}

func GetAddress(kstrs []Key, cc uint32) ([]byte, error) {
	for _, kstr := range kstrs {
		if kstr.PriorityCode() == cc {
			return kstr.GetAddress()
		}
	}
	return []byte{}, errors.New("Unknown Currency Code")
}
