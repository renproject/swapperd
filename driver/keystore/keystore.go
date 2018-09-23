package keystore

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
)

// ErrKeyFileExists is returned when the keystore file exists, and the user is
// trying to overwrite it.
var ErrKeyFileExists = errors.New("Keystore file exists")

// ErrKeyFileDoesNotExist is returned when the keystore file doesnot exist, and
// the user is trying to read from it.
var ErrKeyFileDoesNotExist = errors.New("Keystore file doesnot exist")

// LoadFromFile
func LoadFromFile(repNetwork, loc, passphrase string) (keystore.Keystore, error) {
	var ethLoc, btcLoc string
	if passphrase == "" {
		ethLoc = fmt.Sprintf("%s/ethereum-%s-unsafe.json", loc, repNetwork)
		btcLoc = fmt.Sprintf("%s/bitcoin-%s-unsafe.json", loc, repNetwork)
	} else {
		ethLoc = fmt.Sprintf("%s/ethereum-%s.json", loc, repNetwork)
		btcLoc = fmt.Sprintf("%s/bitcoin-%s.json", loc, repNetwork)
	}

	ethNet, btcNet, err := getSpecificNetworks(repNetwork)
	if err != nil {
		return nil, err
	}
	ethKey, err := LoadKeyFromFile(ethLoc, "ethereum", ethNet, passphrase)
	if err != nil {
		return nil, err
	}
	btcKey, err := LoadKeyFromFile(btcLoc, "bitcoin", btcNet, passphrase)
	if err != nil {
		return nil, err
	}
	return keystore.New(btcKey, ethKey), nil
}

func GenerateRandom(repNetwork string) (keystore.Keystore, error) {
	ethNet, btcNet, err := getSpecificNetworks(repNetwork)
	if err != nil {
		return nil, err
	}
	ethKey, err := keystore.RandomEthereumKey(ethNet)
	if err != nil {
		return nil, err
	}
	btcKey, err := keystore.RandomBitcoinKey(btcNet)
	if err != nil {
		return nil, err
	}
	return keystore.New(ethKey, btcKey), nil
}

// GenerateFile
func GenerateFile(loc string, repNetwork string, passphrase string) error {
	var ethLoc, btcLoc string
	if passphrase == "" {
		ethLoc = fmt.Sprintf("%s/ethereum-%s-unsafe.json", loc, repNetwork)
		btcLoc = fmt.Sprintf("%s/bitcoin-%s-unsafe.json", loc, repNetwork)
	} else {
		ethLoc = fmt.Sprintf("%s/ethereum-%s.json", loc, repNetwork)
		btcLoc = fmt.Sprintf("%s/bitcoin-%s.json", loc, repNetwork)
	}
	ethNet, btcNet, err := getSpecificNetworks(repNetwork)
	if err != nil {
		return err
	}
	if err := StoreKeyToFile(ethLoc, "ethereum", ethNet, passphrase); err != nil {
		return err
	}
	return StoreKeyToFile(btcLoc, "bitcoin", btcNet, passphrase)
}

// LoadKeyFromFile loads a key from a file and tries to decrypt it using the
// given passphrase. If the passphrase is empty, then it tries to load an
// unencrypted key.
func LoadKeyFromFile(loc, chain, network, passphrase string) (keystore.Key, error) {
	data, err := ioutil.ReadFile(loc)
	if err != nil {
		return nil, ErrKeyFileDoesNotExist
	}
	return decodeKey(data, chain, network, passphrase)
}

// StoreKeyToFile stores a key to a file after encrypting it using the given
// passphrase. If the passphrase is empty, then it tries to load an unencrypted
// key.
func StoreKeyToFile(loc, chain, network, passphrase string) error {
	if _, err := ioutil.ReadFile(loc); err == nil {
		return ErrKeyFileExists
	}
	generatedKey, err := generateRandomKey(chain, network, passphrase)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(loc, generatedKey, 0400)
}

// LoadKeyFromNet loads a key from the network and tries to decrypt it using
// the given passphrase. If the  passphrase is empty, then it tries to load an
// unencrypted key.
func LoadKeyFromNet(url, chain, network, passphrase string) (keystore.Key, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, ErrKeyFileExists
	}
	if resp.StatusCode == 200 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return decodeKey(data, chain, network, passphrase)
	}
	return nil, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
}

func generateRandomKey(chain, network, passphrase string) ([]byte, error) {
	switch chain {
	case "bitcoin":
		return GenerateRandomBitcoinKey(network, passphrase)
	case "ethereum":
		return GenerateRandomEthereumKey(passphrase)
	default:
		return nil, fmt.Errorf("Unsupported blockchain %s", chain)
	}
}

func decodeKey(data []byte, chain, network, passphrase string) (keystore.Key, error) {
	switch chain {
	case "bitcoin":
		return DecodeBitcoinKey(data, network, passphrase)
	case "ethereum":
		return DecodeEthereumKey(data, network, passphrase)
	default:
		return nil, fmt.Errorf("Unsupported blockchain %s", chain)
	}
}

func getSpecificNetworks(repNetwork string) (string, string, error) {
	switch repNetwork {
	case "nightly":
		return "kovan", "testnet", nil
	case "falcon":
		return "kovan", "testnet", nil
	case "testnet":
		return "kovan", "testnet", nil
	default:
		return "", "", fmt.Errorf("Unknown RenEx network %s", repNetwork)
	}
}
