package wallet

import (
	"math/big"

	"github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/libbtc-go"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
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

type Wallet interface {
	SupportedTokens() []blockchain.Token
	SupportedBlockchains() []blockchain.Blockchain
	Balances() (map[blockchain.TokenName]blockchain.Balance, error)
	Transfer(password string, token blockchain.Token, to string, amount *big.Int) (string, error)
	GetAddress(blockchain blockchain.BlockchainName) (string, error)
	Addresses() (map[blockchain.TokenName]string, error)
	VerifyAddress(blockchain blockchain.BlockchainName, address string) error
	VerifyBalance(token blockchain.Token, balance *big.Int) error
	DefaultFee(blockchainName blockchain.BlockchainName) (*big.Int, error)

	EthereumAccount(password string) (beth.Account, error)
	BitcoinAccount(password string) (libbtc.Account, error)
}

type wallet struct {
	config Config
}

func New(config Config) Wallet {
	return &wallet{
		config: config,
	}
}
