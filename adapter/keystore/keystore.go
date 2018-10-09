package keystore

import (
	"github.com/republicprotocol/swapperd/foundation"
)

type Key interface {
	Token() foundation.Token
}

type Keystore interface {
	GetKey(token foundation.Token) Key
}

type KeyMap map[foundation.Token]Key

type keystore struct {
	keyMap KeyMap
}

func New(keys ...Key) Keystore {
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
