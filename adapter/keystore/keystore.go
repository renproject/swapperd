package keystore

import "github.com/republicprotocol/renex-swapper-go/domain/token"

type Key interface {
	Token() token.Token
}

type Keystore interface {
	GetKey(token token.Token) Key
}

type KeyMap map[token.Token]Key

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
func (keystore *keystore) GetKey(token token.Token) Key {
	return keystore.keyMap[token]
}
