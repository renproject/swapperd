package keystore

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/republicprotocol/swapperd/adapter/fund"
	"github.com/republicprotocol/swapperd/core/auth"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/sha3"
)

// Testnet is the Swapperd's testnet config object
var Testnet = fund.Config{
	Bitcoin: fund.BlockchainConfig{
		Network: fund.Network{
			Name: "testnet",
		},
		Tokens: []string{"BTC"},
	},
	Ethereum: fund.BlockchainConfig{
		Network: fund.Network{
			Name: "kovan",
			URL:  "https://kovan.infura.io",
		},
		Tokens: []string{"ETH", "WBTC"},
	},
}

// Mainnet is the Swapperd's mainnet config object
var Mainnet = fund.Config{
	Bitcoin: fund.BlockchainConfig{
		Network: fund.Network{
			Name: "mainnet",
		},
		Tokens: []string{"BTC"},
	},
	Ethereum: fund.BlockchainConfig{
		Network: fund.Network{
			Name: "kovan",
			URL:  "https://kovan.infura.io",
		},
		Tokens: []string{"ETH", "WBTC"},
	},
}

type Keystore struct {
	Username     string      `json:"username"`
	PasswordHash string      `json:"passwordHash"`
	Config       fund.Config `json:"config"`
}

type Address struct {
	Blockchain string `json:"token"`
	Address    string `json:"address"`
}

func FundManager(network string) (fund.Manager, error) {
	path := keystorePath(network)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	keystore := Keystore{}
	if err := json.Unmarshal(data, &keystore); err != nil {
		return nil, err
	}
	return fund.New(keystore.Config), nil
}

func LoadAuthenticator(network string) (auth.Authenticator, error) {
	path := keystorePath(network)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	keystore := Keystore{}
	if err := json.Unmarshal(data, &keystore); err != nil {
		return nil, err
	}
	passwordHash, err := toBytes32(keystore.PasswordHash)
	return auth.NewAuthenticator(keystore.Username, passwordHash), nil
}

func Generate(network, username, password string) (string, error) {
	network = strings.ToLower(network)
	path := keystorePath(network)
	keystore := Keystore{}
	keystore.Username = username
	passwordHashBytes := sha3.Sum256([]byte(password))
	keystore.PasswordHash = base64.StdEncoding.EncodeToString(passwordHashBytes[:])
	config, err := generateConfig(network, password)
	if err != nil {
		return "", err
	}
	keystore.Config = config
	data, err := json.Marshal(keystore)
	if err != nil {
		return "", err
	}
	return config.Mnemonic, ioutil.WriteFile(path, data, 0644)
}

func generateConfig(network string, password string) (fund.Config, error) {
	var config fund.Config
	switch network {
	case "testnet":
		config = Testnet
	case "mainnet":
		config = Mainnet
	default:
		return fund.Config{}, fmt.Errorf("Invalid Network %s", network)
	}
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return fund.Config{}, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return fund.Config{}, err
	}
	config.Mnemonic = mnemonic
	manager := fund.New(config)
	ethAccount, err := manager.EthereumAccount(password)
	if err != nil {
		return fund.Config{}, err
	}
	config.Ethereum.Address = ethAccount.Address().String()
	btcAccount, err := manager.BitcoinAccount(password)
	if err != nil {
		return fund.Config{}, err
	}
	btcAddress, err := btcAccount.Address()
	if err != nil {
		return fund.Config{}, err
	}
	config.Bitcoin.Address = btcAddress.String()
	return config, nil
}

func keystorePath(network string) string {
	network = strings.ToLower(network)
	unix := os.Getenv("HOME")
	if unix != "" {
		return fmt.Sprintf("%s/%s.json", unix+"/.swapperd", network)
	}
	windows := os.Getenv("userprofile")
	if windows != "" {
		return fmt.Sprintf("%s\\%s.json", strings.Join(strings.Split(windows, "\\"), "\\\\")+"\\swapperd", network)
	}
	panic(fmt.Sprintf("unknown operating system"))
}

func toBytes32(data string) ([32]byte, error) {
	bytes32 := [32]byte{}
	dataBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil || len(dataBytes) != 32 {
		return bytes32, fmt.Errorf("Invalid data")
	}
	copy(bytes32[:], dataBytes)
	return bytes32, nil
}
