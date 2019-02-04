package keystore

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/republicprotocol/swapperd/adapter/wallet"
)

func Wallet(homeDir, network string) (wallet.Wallet, error) {
	path := keystorePath(homeDir, network)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := wallet.Config{}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return wallet.New(config), nil
}

func Generate(homeDir, network, mnemonic string) error {
	network = strings.ToLower(network)
	path := keystorePath(homeDir, network)
	config, err := generateConfig(network, mnemonic)
	if err != nil {
		return err
	}
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

func generateConfig(network, mnemonic string) (wallet.Config, error) {
	var config wallet.Config
	switch network {
	case "testnet":
		config = wallet.Testnet
	case "mainnet":
		config = wallet.Mainnet
	default:
		return wallet.Config{}, fmt.Errorf("invalid network %s", network)
	}
	config.Mnemonic = mnemonic
	return config, nil
}

func keystorePath(homeDir, network string) string {
	return path.Join(homeDir, fmt.Sprintf("%s.json", network))
}

func toBytes32(data string) ([32]byte, error) {
	bytes32 := [32]byte{}
	dataBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil || len(dataBytes) != 32 {
		return bytes32, fmt.Errorf("invalid data")
	}
	copy(bytes32[:], dataBytes)
	return bytes32, nil
}
