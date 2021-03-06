package wallet

import (
	"math/big"

	"github.com/renproject/libbtc-go"
	"github.com/renproject/libeth-go"
	"github.com/renproject/swapperd/core/wallet/transfer"
	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/renproject/tokens"
	"github.com/sirupsen/logrus"
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
	SupportedTokens() []tokens.Token
	Balances(password string) (map[tokens.Name]blockchain.Balance, error)
	Balance(password string, token tokens.Token) (blockchain.Balance, error)
	Lookup(token tokens.Token, txHash string) (transfer.UpdateReceipt, error)
	Transfer(password string, token tokens.Token, to string, amount *big.Int, speed blockchain.TxExecutionSpeed, senAll bool) (string, blockchain.Cost, error)
	GetAddress(password string, blockchainName tokens.BlockchainName) (string, error)
	Addresses(password string) (map[tokens.Name]string, error)
	VerifyAddress(blockchain tokens.BlockchainName, address string) error
	VerifyBalance(password string, token tokens.Token, balance *big.Int) error

	EthereumAccount(password string) (libeth.Account, error)
	BitcoinAccount(password string) (libbtc.Account, error)
	ECDSASigner(password string) (ECDSASigner, error)
}

type wallet struct {
	config Config
	logger logrus.FieldLogger
}

func New(config Config, logger logrus.FieldLogger) Wallet {
	return &wallet{
		config: config,
		logger: logger,
	}
}
