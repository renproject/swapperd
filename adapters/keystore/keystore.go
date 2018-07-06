package keystore

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"

	"github.com/republicprotocol/atom-go/adapters/key/btc"
	"github.com/republicprotocol/atom-go/adapters/key/eth"
	"github.com/republicprotocol/atom-go/services/swap"
)

var ErrUnknownChain = errors.New("unknown chain")

type keyJSON struct {
	PrivateKey   string `json:"privateKey"`
	Chain        string `json:"chain"`
	PriorityCode uint32 `json:"priorityCode"`
}
type keystore struct {
	Keys []keyJSON `json:"keys"`
	path string
	mu   *sync.RWMutex
}

func NewKeystore(path string) swap.Keystore {
	return &keystore{
		path: path,
		mu:   new(sync.RWMutex),
	}
}

func (kstr *keystore) LoadKeys() ([]swap.Key, error) {
	err := kstr.loadFromFile()
	if err != nil {
		return nil, err
	}
	kstr.mu.RLock()
	defer kstr.mu.RUnlock()

	keys := []swap.Key{}

	for _, key := range kstr.Keys {
		switch key.PriorityCode {
		case uint32(0):
			btcKey, err := btc.NewBitcoinKey(key.PrivateKey, key.Chain)
			if err != nil {
				return nil, errors.New("Failed to initialize bitcoin key")
			}
			keys = append(keys, btcKey)
		case uint32(1):
			ethKey, err := eth.NewEthereumKey(key.PrivateKey, key.Chain)
			if err != nil {
				return nil, errors.New("Failed to initialize ethereum key")
			}
			keys = append(keys, ethKey)
		default:
			return nil, errors.New("Unknown blockchain")
		}
	}
	return keys, nil
}

func (kstr *keystore) loadFromFile() error {
	kstr.mu.Lock()
	defer kstr.mu.Unlock()

	raw, err := ioutil.ReadFile(kstr.path)
	if err != nil {
		return nil
	}
	return json.Unmarshal(raw, &kstr)
}
