package ethclient

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// KeyStore are the parameters required to deploy contracts and create an ethereum client
type KeyStore struct {
	Chain      string `json:"chain"`
	URL        string `json:"url"`
	PrivateKey string `json:"private_key"`
}

// KeyStoreFile has the KeyStores of different ethereum chains
type KeyStoreFile struct {
	Ganache KeyStore `json:"ganache"`
	Ropsten KeyStore `json:"ropsten"`
	Kovan   KeyStore `json:"kovan"`
	Mainnet KeyStore `json:"mainnet"`
}

var keyStorePath = path.Join(os.Getenv("GOPATH"), "/src/github.com/republicprotocol/atom-go/adapters/clients/secrets/ethKeys.json")

func readKeyStore(chain string) (KeyStore, error) {
	var keyStore KeyStore
	var keyStoreFile KeyStoreFile

	raw, err := ioutil.ReadFile(keyStorePath)
	if err != nil {
		return keyStore, err
	}
	json.Unmarshal(raw, &keyStoreFile)

	switch chain {
	case "ganache":
		keyStore = keyStoreFile.Ganache
	case "ropsten":
		keyStore = keyStoreFile.Ropsten
	case "kovan":
		keyStore = keyStoreFile.Kovan
	case "mainnet":
		keyStore = keyStoreFile.Mainnet
	}
	return keyStore, nil
}

func writeKeyStore(KeyStore KeyStore) error {
	var KeyStoreFile KeyStoreFile

	raw, err := ioutil.ReadFile(keyStorePath)
	if err != nil {
		return err
	}
	json.Unmarshal(raw, &KeyStoreFile)

	switch KeyStore.Chain {
	case "ganache":
		KeyStoreFile.Ganache = KeyStore
	case "ropsten":
		KeyStoreFile.Ropsten = KeyStore
	case "kovan":
		KeyStoreFile.Kovan = KeyStore
	case "mainnet":
		KeyStoreFile.Mainnet = KeyStore
	}

	data, err := json.Marshal(KeyStoreFile)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(keyStorePath, data, 700)
}
