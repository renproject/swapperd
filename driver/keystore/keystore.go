package keystore

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/republicprotocol/swapperd/adapter/wallet"
	"golang.org/x/crypto/bcrypt"
)

// Testnet is the Swapperd's testnet config object
var Testnet = wallet.Config{
	Bitcoin: wallet.BlockchainConfig{
		Network: wallet.Network{
			Name: "testnet",
		},
		Tokens: []string{"BTC"},
	},
	Ethereum: wallet.BlockchainConfig{
		Network: wallet.Network{
			Name: "kovan",
			URL:  "https://kovan.infura.io",
		},
		Tokens: []string{"ETH", "WBTC"},
	},
}

// Mainnet is the Swapperd's mainnet config object
var Mainnet = wallet.Config{
	Bitcoin: wallet.BlockchainConfig{
		Network: wallet.Network{
			Name: "mainnet",
		},
		Tokens: []string{"BTC"},
	},
	Ethereum: wallet.BlockchainConfig{
		Network: wallet.Network{
			Name: "mainnet",
			URL:  "https://mainnet.infura.io",
		},
		Tokens: []string{"ETH", "WBTC"},
	},
}

type Keystore struct {
	Username     string        `json:"username"`
	PasswordHash string        `json:"passwordHash"`
	Config       wallet.Config `json:"config"`
}

type Address struct {
	Blockchain string `json:"token"`
	Address    string `json:"address"`
}

func Wallet(homeDir, network string) (wallet.Wallet, error) {
	path := keystorePath(homeDir, network)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	keystore := Keystore{}
	if err := json.Unmarshal(data, &keystore); err != nil {
		return nil, err
	}
	return wallet.New(keystore.Config), nil
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

func generateConfig(network, password, mnemonic string) (wallet.Config, error) {
	var config wallet.Config
	switch network {
	case "testnet":
		config = Testnet
	case "mainnet":
		config = Mainnet
	default:
		return wallet.Config{}, fmt.Errorf("Invalid Network %s", network)
	}
	config.Mnemonic = mnemonic
	wallet := wallet.New(config)
	ethAccount, err := wallet.EthereumAccount(password)
	if err != nil {
		return wallet.Config{}, err
	}
	config.Ethereum.Address = ethAccount.Address().String()
	btcAccount, err := wallet.BitcoinAccount(password)
	if err != nil {
		return wallet.Config{}, err
	}
	btcAddress, err := btcAccount.Address()
	if err != nil {
		return wallet.Config{}, err
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
