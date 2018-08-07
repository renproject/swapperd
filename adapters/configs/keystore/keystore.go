package keystore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

var ErrPrefix = "Config: Keystore: %s"

type keystore struct {
	Keys []key `json:"keys"`

	keyMap map[uint32][]key
	mu     *sync.RWMutex
	path   string
}
type Keystore interface {
	GetKey(uint32, uint32) (Key, error)
	AppendKey(uint32, uint32, Key)
}

func NewKeystore(priorityCodes []uint32, chains []string, path string) (Keystore, error) {
	if len(priorityCodes) != len(chains) {
		return nil, fmt.Errorf(ErrPrefix, "Invalid Parameters")
	}

	var keystore keystore
	keystore.mu = new(sync.RWMutex)
	keystore.keyMap = make(map[uint32][]key)
	keystore.path = path

	for i := range priorityCodes {
		key, err := generateKey(priorityCodes[i], chains[i])
		if err != nil {
			return nil, err
		}
		keystore.Keys = append(keystore.Keys, key)
	}

	if err := keystore.update(); err != nil {
		return nil, err
	}

	for _, key := range keystore.Keys {
		keystore.keyMap[key.Code] = append(keystore.keyMap[key.Code], key)
	}

	return &keystore, nil
}

func Load(path string) (Keystore, error) {
	var keystore keystore
	keystore.path = path
	keystore.mu = new(sync.RWMutex)
	keystore.keyMap = make(map[uint32][]key)
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return &keystore, err
	}
	json.Unmarshal(raw, &keystore)
	for _, key := range keystore.Keys {
		keystore.keyMap[key.Code] = append(keystore.keyMap[key.Code], key)
	}
	return &keystore, nil
}

func (keystore *keystore) GetKey(priorityCode, index uint32) (Key, error) {
	key := keystore.keyMap[priorityCode][index]
	return NewKey(key.PrivateKey, key.Code, key.Network)
}

func (keystore *keystore) AppendKey(priorityCode, index uint32, keyInf Key) {
	keyStruct := key{
		keyInf.GetKeyString(),
		keyInf.PriorityCode(),
		keyInf.Chain(),
	}
	keystore.Keys = append(keystore.Keys, keyStruct)
	keystore.keyMap[keyStruct.Code] = append(keystore.keyMap[keyStruct.Code], keyStruct)
}

func generateKey(priorityCode uint32, chain string) (key, error) {
	switch priorityCode {
	case 0:
		btcKey, err := RandomBitcoinKeyString(chain)
		return key{
			Code:       0,
			PrivateKey: btcKey,
			Network:    chain,
		}, err
	case 1:
		ethKey, err := RandomEthereumKeyString()
		return key{
			Code:       1,
			PrivateKey: ethKey,
			Network:    chain,
		}, err
	}
	return key{}, fmt.Errorf(ErrPrefix, "Unknown Priority Code")
}

func (keystore *keystore) update() error {
	data, err := json.Marshal(keystore)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(keystore.path, data, 0600)
}
