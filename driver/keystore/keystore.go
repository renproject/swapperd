package keystore

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/republicprotocol/swapperd/adapter/fund"
	"golang.org/x/crypto/bcrypt"
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
			Name: "mainnet",
			URL:  "https://mainnet.infura.io",
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

func FundManager(homeDir, network string) (fund.Manager, error) {
	path := keystorePath(homeDir, network)
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

func LoadPasswordHash(homeDir, network string) ([]byte, error) {
	path := keystorePath(homeDir, network)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	keystore := Keystore{}
	if err := json.Unmarshal(data, &keystore); err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(keystore.PasswordHash)
}

func Generate(homeDir, network, username, password, mnemonic string) error {
	network = strings.ToLower(network)
	path := keystorePath(homeDir, network)
	keystore := Keystore{}
	keystore.Username = username

	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	keystore.PasswordHash = base64.StdEncoding.EncodeToString(passwordHashBytes)
	config, err := generateConfig(network, password, mnemonic)
	if err != nil {
		return err
	}
	keystore.Config = config
	data, err := json.Marshal(keystore)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

func generateConfig(network, password, mnemonic string) (fund.Config, error) {
	var config fund.Config
	switch network {
	case "testnet":
		config = Testnet
	case "mainnet":
		config = Mainnet
	default:
		return fund.Config{}, fmt.Errorf("Invalid Network %s", network)
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

func keystorePath(homeDir, network string) string {
	return path.Join(homeDir, fmt.Sprintf("%s.json", network))
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
