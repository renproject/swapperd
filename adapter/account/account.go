package account

import (
	"github.com/republicprotocol/swapperd/foundation"
)

type Keystore interface {
	GetKey(token foundation.Token) interface{}
}

type KeyMap map[foundation.Token]interface{}

type keystore struct {
	keyMap KeyMap
}

func New(accounts ...interface{}) Keystore {
	keyMap := KeyMap{}

	for _, key := range keys {
		keyMap[key.Token()] = key
	}

	return &keystore{
		keyMap,
	}
}

// GetKey returns the key object of the given token
func (keystore *keystore) GetKey(token foundation.Token) Key {
	return keystore.keyMap[token]
}
