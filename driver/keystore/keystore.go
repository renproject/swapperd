package keystore

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
)

// ErrKeyFileExists is returned when the keystore file exists, and the user is
// trying to overwrite it.
var ErrKeyFileExists = errors.New("Keystore file exists")

// ErrKeyFileDoesNotExist is returned when the keystore file doesnot exist, and
// the user is trying to read from it.
var ErrKeyFileDoesNotExist = errors.New("Keystore file doesnot exist")

// LoadFromFile
func LoadFromFile(conf config.Config, passphrase string) (keystore.Keystore, error) {
	keys := []keystore.Key{}
	for _, token := range conf.SupportedCurrencies {
		var loc string
		if passphrase == "" {
			loc = fmt.Sprintf("%s/%s-%s-unsafe.json", conf.HomeDir, token, conf.RenEx.Network)
		} else {
			loc = fmt.Sprintf("%s/%s-%s.json", conf.HomeDir, token, conf.RenEx.Network)
		}
		key, err := LoadKeyFromFile(loc, passphrase, conf, token)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keystore.New(keys...), nil
}

func GenerateRandomKeystore(conf config.Config) (keystore.Keystore, error) {
	keys := []keystore.Key{}
	for _, token := range conf.SupportedCurrencies {
		key, err := randomKey(conf, token)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keystore.New(keys...), nil
}

// GenerateFile
func GenerateFile(conf config.Config, passphrase string) error {
	for _, token := range conf.SupportedCurrencies {
		var loc string
		if passphrase == "" {
			loc = fmt.Sprintf("%s/%s-%s-unsafe.json", conf.HomeDir, token, conf.RenEx.Network)
		} else {
			loc = fmt.Sprintf("%s/%s-%s.json", conf.HomeDir, token, conf.RenEx.Network)
		}
		if err := StoreKeyToFile(loc, passphrase, conf, token); err != nil {
			return err
		}
	}
	return nil
}

// LoadKeyFromFile loads a key from a file and tries to decrypt it using the
// given passphrase. If the passphrase is empty, then it tries to load an
// unencrypted key.
func LoadKeyFromFile(loc, passphrase string, conf config.Config, tok token.Token) (keystore.Key, error) {
	data, err := ioutil.ReadFile(loc)
	if err != nil {
		return nil, ErrKeyFileDoesNotExist
	}
	return decodeKey(data, passphrase, conf, tok)
}

// StoreKeyToFile stores a key to a file after encrypting it using the given
// passphrase. If the passphrase is empty, then it tries to load an unencrypted
// key.
func StoreKeyToFile(loc, passphrase string, conf config.Config, tok token.Token) error {
	if _, err := ioutil.ReadFile(loc); err == nil {
		return ErrKeyFileExists
	}
	generatedKey, err := generateRandomKey(passphrase, conf, tok)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(loc, generatedKey, 0400)
}

func randomKey(conf config.Config, tok token.Token) (keystore.Key, error) {
	switch tok {
	case token.BTC:
		return keystore.RandomBitcoinKey(conf.Bitcoin.Network)
	case token.ETH:
		return keystore.RandomEthereumKey(conf.Ethereum.Network)
	default:
		return nil, token.ErrUnsupportedToken
	}
}

func generateRandomKey(passphrase string, conf config.Config, tok token.Token) ([]byte, error) {
	switch tok {
	case token.BTC:
		return GenerateRandomBitcoinKey(conf.Bitcoin.Network, passphrase)
	case token.ETH:
		return GenerateRandomEthereumKey(passphrase)
	default:
		return nil, token.ErrUnsupportedToken
	}
}

func decodeKey(data []byte, passphrase string, conf config.Config, tok token.Token) (keystore.Key, error) {
	switch tok {
	case token.BTC:
		return DecodeBitcoinKey(data, conf.Bitcoin.Network, passphrase)
	case token.ETH:
		return DecodeEthereumKey(data, conf.Ethereum.Network, passphrase)
	default:
		return nil, token.ErrUnsupportedToken
	}
}

// // LoadKeyFromNet loads a key from the network and tries to decrypt it using
// // the given passphrase. If the  passphrase is empty, then it tries to load an
// // unencrypted key.
// func LoadKeyFromNet(url, chain, network, passphrase string) (keystore.Key, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, ErrKeyFileExists
// 	}
// 	if resp.StatusCode == 200 {
// 		data, err := ioutil.ReadAll(resp.Body)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return decodeKey(data, passphrase, network, )
// 	}
// 	return nil, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
// }
