package keystore

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/republicprotocol/swapperd/adapter/account"
	"github.com/republicprotocol/swapperd/core/auth"
	"github.com/tyler-smith/go-bip39"
)

// Testnet is the Swapperd's testnet config object
var Testnet = account.Config{
	Bitcoin: "testnet",
	Ethereum: account.EthereumConfig{
		URL:     "https://kovan.infura.io",
		Network: "kovan",
		Swapper: "0x2218fa20c33765e7e01671ee6aaca75fbaf3a974",
		Tokens: []account.EthereumToken{
			account.EthereumToken{
				Name:    "WBTC",
				ERC20:   "0xA1D3EEcb76285B4435550E4D963B8042A8bffbF0",
				Swapper: "0x2218fa20c33765e7e01671ee6aaca75fbaf3a974",
			},
		},
	},
}

// Mainnet is the Swapperd's mainnet config object
var Mainnet = account.Config{
	Bitcoin: "mainnet",
	Ethereum: account.EthereumConfig{
		URL:     "https://kovan.infura.io",
		Network: "kovan",
		Swapper: "0x2218fa20c33765e7e01671ee6aaca75fbaf3a974",
		Tokens: []account.EthereumToken{
			account.EthereumToken{
				Name:    "WBTC",
				ERC20:   "0xA1D3EEcb76285B4435550E4D963B8042A8bffbF0",
				Swapper: "0x2218fa20c33765e7e01671ee6aaca75fbaf3a974",
			},
		},
	},
}

type Keystore struct {
	Username     string         `json:"username"`
	PasswordHash string         `json:"passwordHash"`
	Mnemonic     string         `json:"mnemonic"`
	Config       account.Config `json:"config"`
}

func LoadAccounts(network string) (account.Accounts, error) {
	path := keystorePath(network)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	keystore := Keystore{}
	if err := json.Unmarshal(data, &keystore); err != nil {
		return nil, err
	}
	return account.New(keystore.Mnemonic, keystore.Config), nil
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

func Generate(network, username, passwordHash string) (string, error) {
	network = strings.ToLower(network)
	path := keystorePath(network)

	keystore := Keystore{}
	switch network {
	case "testnet":
		keystore.Config = Testnet
	case "mainnet":
		keystore.Config = Mainnet
	default:
		return "", fmt.Errorf("Invalid Network %s", network)
	}
	keystore.Username = username
	keystore.PasswordHash = passwordHash

	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	keystore.Mnemonic = mnemonic

	data, err := json.Marshal(keystore)
	if err != nil {
		return "", err
	}
	return mnemonic, ioutil.WriteFile(path, data, 0644)
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
	panic(fmt.Sprintf("unknown Operating System: unix: %s windows: %s", os.Getenv("HOME"), os.Getenv("userprofile")))
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
