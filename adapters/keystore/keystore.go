package keystore

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

var ErrUnknownChain = errors.New("unknown chain")

type Keystore interface {
	LoadKeypair(string) (*ecdsa.PrivateKey, error)
	UpdateEtherumKey(*ecdsa.PrivateKey) error
	UpdateEtherumKeyString(string) error
	RandomEtherumKey() error
	UpdateBitcoinKey(*ecdsa.PrivateKey) error
	UpdateBitcoinKeyString(string) error
	RandomBitcoinKey() error
}

type keystore struct {
	Ethereum string `json:"ethereum"`
	Bitcoin  string `json:"bitcoin"`

	path string
	mu   *sync.RWMutex
}

func NewKeystore(path string) Keystore {
	return &keystore{
		path: path,
		mu:   new(sync.RWMutex),
	}
}

func (keystore *keystore) LoadKeypair(chain string) (*ecdsa.PrivateKey, error) {
	err := keystore.loadFromFile()
	if err != nil {
		return nil, err
	}
	keystore.mu.RLock()
	defer keystore.mu.RUnlock()

	switch chain {
	case "ethereum":
		return ethCrypto.HexToECDSA(keystore.Ethereum)
	case "bitcoin":
		return ethCrypto.HexToECDSA(keystore.Bitcoin)
	default:
		return nil, errors.New("Unknown blockchain")
	}
}

func (keystore *keystore) UpdateEtherumKey(pk *ecdsa.PrivateKey) error {
	return keystore.UpdateEtherumKeyString(hex.EncodeToString(ethCrypto.FromECDSA(pk)))
}

func (keystore *keystore) UpdateEtherumKeyString(pk string) error {
	err := keystore.loadFromFile()
	if err != nil {
		return err
	}
	keystore.mu.Lock()
	defer keystore.mu.Unlock()

	keystore.Ethereum = pk

	println(keystore.Ethereum)
	data, err := json.Marshal(&keystore)
	if err != nil {
		return err
	}
	fmt.Println("Updated ethereum key string")
	return ioutil.WriteFile(keystore.path, data, 0777)
}

func (keystore *keystore) RandomEtherumKey() error {
	keyPair, err := ethCrypto.GenerateKey()
	if err != nil {
		return err
	}
	return keystore.UpdateEtherumKey(keyPair)
}

func (keystore *keystore) UpdateBitcoinKey(pk *ecdsa.PrivateKey) error {
	return keystore.UpdateBitcoinKeyString(hex.EncodeToString(ethCrypto.FromECDSA(pk)))
}

func (keystore *keystore) UpdateBitcoinKeyString(pk string) error {
	err := keystore.loadFromFile()
	if err != nil {
		return err
	}
	keystore.mu.Lock()
	defer keystore.mu.Unlock()

	keystore.Bitcoin = pk
	data, err := json.Marshal(keystore)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(keystore.path, data, 0777)
}

func (keystore *keystore) RandomBitcoinKey() error {
	keyPair, err := ethCrypto.GenerateKey()
	if err != nil {
		return err
	}
	return keystore.UpdateBitcoinKey(keyPair)
}

func (keystore *keystore) loadFromFile() error {
	keystore.mu.Lock()
	defer keystore.mu.Unlock()

	raw, err := ioutil.ReadFile(keystore.path)
	if err != nil {
		return nil
	}
	return json.Unmarshal(raw, keystore)
}
