package blockchain

import "fmt"

func ErrUnsupportedBlockchain(blockchain BlockchainName) error {
	return fmt.Errorf("unsupported blockchain: %s", blockchain)
}

type BlockchainName string

var (
	Bitcoin  = BlockchainName("bitcoin")
	Ethereum = BlockchainName("ethereum")
)

type Blockchain struct {
	Name    BlockchainName `json:"name"`
	Address string         `json:"address"`
}

type Balance struct {
	Address string `json:"address"`
	Amount  string `json:"balance"`
}
