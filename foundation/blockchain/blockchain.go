package blockchain

import "fmt"

type BlockchainName string

var (
	Bitcoin  = BlockchainName("bitcoin")
	Ethereum = BlockchainName("ethereum")
)

type Blockchain struct {
	Name    BlockchainName `json:"name"`
	Address string         `json:"address"`
}

type ErrUnsupportedBlockchain string

func NewErrUnsupportedBlockchain(blockchain BlockchainName) error {
	return ErrUnsupportedBlockchain(fmt.Sprintf("unsupported blockchain: %s", blockchain))
}

func (err ErrUnsupportedBlockchain) Error() string {
	return string(err)
}
