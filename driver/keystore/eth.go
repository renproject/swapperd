package keystore

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"

	ethKeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pborman/uuid"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
)

type plainEthereumKeyJSON struct {
	PrivateKey string `json:"privateKey"`
}

// GenerateRandomEthereumKey creates a new encrypted ethereum key. If an empty
// passphrase is given this function generates an unencrypted ethereum key.
func GenerateRandomEthereumKey(passphrase string) ([]byte, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	if passphrase != "" {
		return encryptEthereumKey(privKey, passphrase)
	}
	return plainEthereumKey(privKey)
}

// DecodeEthereumKey decrypts an ethereum key using the given keystore. If an
// empty passphrase is given this function generates an unencrypted ethereum
// key.
func DecodeEthereumKey(data []byte, ethNetwork string, passphrase string) (keystore.EthereumKey, error) {
	if passphrase != "" {
		return decodeEncryptedEthereumKey(data, ethNetwork, passphrase)
	}
	return decodePlainEthereumKey(data, ethNetwork)
}

func decodePlainEthereumKey(key []byte, ethNetwork string) (keystore.EthereumKey, error) {
	plainEthKey := plainEthereumKeyJSON{}
	json.Unmarshal(key, &plainEthKey)
	privKey, err := crypto.HexToECDSA(plainEthKey.PrivateKey)
	if err != nil {
		return keystore.EthereumKey{}, err
	}
	return keystore.NewEthereumKey(privKey, ethNetwork)
}

func decodeEncryptedEthereumKey(key []byte, ethNetwork string, passphrase string) (keystore.EthereumKey, error) {
	privKey, err := ethKeystore.DecryptKey(key, passphrase)
	if err != nil {
		return keystore.EthereumKey{}, nil
	}
	return keystore.NewEthereumKey(privKey.PrivateKey, ethNetwork)
}

func encryptEthereumKey(privKey *ecdsa.PrivateKey, passphrase string) ([]byte, error) {
	key := ethKeystore.Key{
		Id:         uuid.NewRandom(),
		PrivateKey: privKey,
		Address:    crypto.PubkeyToAddress(privKey.PublicKey),
	}
	encryptedKey, err := ethKeystore.EncryptKey(&key, passphrase, ethKeystore.StandardScryptN, ethKeystore.StandardScryptP)
	if err != nil {
		return nil, err
	}
	return encryptedKey, nil
}

func plainEthereumKey(privKey *ecdsa.PrivateKey) ([]byte, error) {
	ks := plainEthereumKeyJSON{
		hex.EncodeToString(crypto.FromECDSA(privKey)),
	}
	data, err := json.Marshal(ks)
	if err != nil {
		return nil, err
	}
	return data, nil
}
