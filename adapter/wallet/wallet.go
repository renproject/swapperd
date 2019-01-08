package wallet

import (
	"math/big"

	"github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/libbtc-go"
	"github.com/republicprotocol/swapperd/core/wallet/balance"
	"github.com/republicprotocol/swapperd/core/wallet/transfer"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
)

type Config struct {
	Mnemonic    string           `json:"mnemonic"`
	IDPublicKey string           `json:"idPublicKey"`
	Ethereum    BlockchainConfig `json:"ethereum"`
	Bitcoin     BlockchainConfig `json:"bitcoin"`
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
	ID() string
	SupportedTokens() []blockchain.Token
	SupportedBlockchains() []blockchain.Blockchain
	Balances() (balance.BalanceMap, error)
	BalancesWithPassword(password string) (balance.BalanceMap, error)
	Lookup(token blockchain.Token, txHash string) (transfer.UpdateReceipt, error)
	Transfer(password string, token blockchain.Token, to string, amount *big.Int) (string, error)
	GetAddress(blockchain blockchain.BlockchainName) (string, error)
	GetAddressWithPassword(blockchainName blockchain.BlockchainName, password string) (string, error)
	Addresses() (map[blockchain.TokenName]string, error)
	AddressesWithPassword(password string) (map[blockchain.TokenName]string, error)
	VerifyAddress(blockchain blockchain.BlockchainName, address string) error
	VerifyBalance(token blockchain.Token, balance *big.Int) error
	DefaultFee(blockchainName blockchain.BlockchainName) (*big.Int, error)

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
