package wallet

import (
	"math/big"

	"github.com/renproject/swapperd/core/wallet/transfer"
	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/libbtc-go"
)

type Config struct {
	Mnemonic string           `json:"mnemonic"`
	Ethereum BlockchainConfig `json:"ethereum"`
	Bitcoin  BlockchainConfig `json:"bitcoin"`
}

type BlockchainConfig struct {
	Network Network `json:"network"`
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
	ID(password, idType string) (string, error)
	SupportedTokens() []blockchain.Token
	Balances(password string) (map[blockchain.TokenName]blockchain.Balance, error)
	Balance(password string, token blockchain.Token) (blockchain.Balance, error)
	Lookup(token blockchain.Token, txHash string) (transfer.UpdateReceipt, error)
	Transfer(password string, token blockchain.Token, to string, amount, fee *big.Int, senAll bool) (string, blockchain.Cost, error)
	GetAddress(password string, blockchainName blockchain.BlockchainName) (string, error)
	Addresses(password string) (map[blockchain.TokenName]string, error)
	VerifyAddress(blockchain blockchain.BlockchainName, address string) error
	VerifyBalance(password string, token blockchain.Token, balance *big.Int) error

	EthereumAccount(password string) (beth.Account, error)
	BitcoinAccount(password string) (libbtc.Account, error)
	ECDSASigner(password string) (ECDSASigner, error)
}

type wallet struct {
	config Config
}

func New(config Config) Wallet {
	return &wallet{
		config: config,
	}
}
