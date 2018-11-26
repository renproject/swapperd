package fund

import (
	"math/big"

	"github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/libbtc-go"
	"github.com/republicprotocol/swapperd/foundation"
)

type Config struct {
	Mnemonic string           `json:"mnemonic"`
	Ethereum BlockchainConfig `json:"ethereum"`
	Bitcoin  BlockchainConfig `json:"bitcoin"`
}

type BlockchainConfig struct {
	Network Network  `json:"network"`
	Address string   `json:"address"`
	Tokens  []string `json:"tokens"`
}

type Network struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Balance struct {
	Address string
	Amount  *big.Int
}

type Manager interface {
	SupportedTokens() []foundation.Token
	SupportedBlockchains() []foundation.Blockchain
	Balances() (map[foundation.TokenName]foundation.Balance, error)
	Transfer(password string, token foundation.Token, to string, amount *big.Int) (string, error)
	VerifyAddress(blockchain foundation.BlockchainName, address string) error
	VerifyBalance(token foundation.Token, balance *big.Int) error
	EthereumAccount(password string) (beth.Account, error)
	BitcoinAccount(password string) (libbtc.Account, error)
}

type manager struct {
	config Config
}

func New(config Config) Manager {
	return &manager{
		config: config,
	}
}
