package keystore

import (
	"fmt"
	"io/ioutil"

	"github.com/republicprotocol/renex-swapper-go/utils"

	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
)

// NewErrKeyFileExists is returned when the keystore file exists, and the user is
// trying to overwrite it.
func NewErrKeyFileExists(loc string) error {
	return fmt.Errorf("Keystore file exists at %s", loc)
}

// NewErrKeyFileDoesNotExist is returned when the keystore file doesnot exist, and
// the user is trying to read from it.
func NewErrKeyFileDoesNotExist(loc string) error {
	return fmt.Errorf("Keystore file not found at %s", loc)
}

// LoadFromFile
func LoadFromFile(conf config.Config, passphrase string) (keystore.Keystore, error) {
	keys := []keystore.Key{}
	for _, token := range conf.SupportedCurrencies {
		var loc string
		if passphrase == "" {
			loc = utils.BuildKeystorePath(conf.HomeDir, string(token), conf.RenEx.Network, true)
		} else {
			loc = utils.BuildKeystorePath(conf.HomeDir, string(token), conf.RenEx.Network, false)
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
			loc = utils.BuildKeystorePath(conf.HomeDir, string(token), conf.RenEx.Network, true)
		} else {
			loc = utils.BuildKeystorePath(conf.HomeDir, string(token), conf.RenEx.Network, false)
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
		return nil, NewErrKeyFileDoesNotExist(loc)
	}
	return decodeKey(data, passphrase, conf, tok)
}

// StoreKeyToFile stores a key to a file after encrypting it using the given
// passphrase. If the passphrase is empty, then it tries to load an unencrypted
// key.
func StoreKeyToFile(loc, passphrase string, conf config.Config, tok token.Token) error {
	if _, err := ioutil.ReadFile(loc); err == nil {
		return NewErrKeyFileExists(loc)
	}
	generatedKey, err := generateRandomKey(passphrase, conf, tok)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(loc, generatedKey, 0444)
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
