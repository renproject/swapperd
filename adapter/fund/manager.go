package fund

import (
	"math/big"

	"github.com/republicprotocol/swapperd/adapter/binder"
	"github.com/republicprotocol/swapperd/core/request"
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
	request.FundManager
	binder.Accounts
}

type manager struct {
	config Config
}

func New(config Config) Manager {
	return &manager{
		config: config,
	}
}
