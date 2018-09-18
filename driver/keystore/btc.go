package keystore

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"

	"github.com/btcsuite/btcutil"

	"github.com/btcsuite/btcd/btcec"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/utils"
)

// ErrMalformedPadding is returned when padding cannot be stripped during
// decryption.
var ErrMalformedPadding = errors.New("malformed padding")

// ErrMalformedCipherText is returned when a cipher text is not a multiple of
// the block size.
var ErrMalformedCipherText = errors.New("malformed cipher text")

type EncryptedBitcoinKey struct {
	CipherText string `json:"ciphertext"`
}

type PlainBitcoinKey struct {
	PrivateKey string `json:"privateKey"`
}

// GenerateRandomBitcoinKey creates a new encrypted Bitcoin keystore. If an
// empty passphrase is given this function generates an unencrypted Bitcoin
// keystore.
func GenerateRandomBitcoinKey(network, passphrase string) ([]byte, error) {
	btcNetwork := utils.GetChainParams(network)
	key, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}
	wif, err := btcutil.NewWIF(key, btcNetwork, true)
	if err != nil {
		return nil, err
	}
	if passphrase != "" {
		return encryptBitcoinKey(wif.String(), passphrase)
	}
	return plainBitcoinKey(wif.String())
}

// DecodeBitcoinKey decrypts a bitcoin key using the given passphrase. If an
// empty passphrase is given this function decodes an unencrypted bitcoin key.
func DecodeBitcoinKey(key []byte, btcNetwork string, passphrase string) (keystore.BitcoinKey, error) {
	return keystore.BitcoinKey{}, nil
}

func decodePlainBitcoinKey(key []byte, btcNetwork string) (keystore.BitcoinKey, error) {
	plainBtcKey := PlainBitcoinKey{}
	json.Unmarshal(key, &plainBtcKey)
	return keystore.NewBitcoinKey(plainBtcKey.PrivateKey, btcNetwork)
}

func decodeEncryptedBitcoinKey(encryptedKey []byte, btcNetwork string, passphrase string) (keystore.BitcoinKey, error) {
	encryptedBtcKey := EncryptedBitcoinKey{}
	json.Unmarshal(encryptedKey, &encryptedBtcKey)

	hash := sha256.Sum256([]byte(passphrase))
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return keystore.BitcoinKey{}, err
	}

	cipherText, err := base64.StdEncoding.DecodeString(encryptedBtcKey.CipherText)
	if err != nil {
		return keystore.BitcoinKey{}, err
	}

	if (len(cipherText) % aes.BlockSize) != 0 {
		return keystore.BitcoinKey{}, ErrMalformedCipherText
	}
	iv := cipherText[:aes.BlockSize]
	message := cipherText[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(message, message)

	strippedMsg, err := strip(message)
	if err != nil {
		return keystore.BitcoinKey{}, err
	}

	return keystore.NewBitcoinKey(string(strippedMsg), btcNetwork)
}

func encryptBitcoinKey(privKey, passphrase string) ([]byte, error) {
	hash := sha256.Sum256([]byte(passphrase))
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}

	paddedPlainText := pad([]byte(privKey))
	cipherText := make([]byte, aes.BlockSize+len(paddedPlainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], []byte(paddedPlainText))

	encryptBitcoinKey := EncryptedBitcoinKey{
		CipherText: base64.StdEncoding.EncodeToString(cipherText),
	}

	return json.Marshal(encryptBitcoinKey)
}

func plainBitcoinKey(privKey string) ([]byte, error) {
	btcKey := PlainBitcoinKey{
		PrivateKey: privKey,
	}
	return json.Marshal(&btcKey)
}

func pad(src []byte) []byte {
	p := aes.BlockSize - len(src)%aes.BlockSize
	padding := bytes.Repeat([]byte{byte(p)}, p)
	return append(src, padding...)
}

func strip(src []byte) ([]byte, error) {
	length := len(src)
	p := int(src[length-1])
	if p > length {
		return nil, ErrMalformedPadding
	}
	return src[:(length - p)], nil
}
