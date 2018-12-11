package blockchain

type BlockchainName string

var (
	Bitcoin  = BlockchainName("bitcoin")
	Ethereum = BlockchainName("ethereum")
)

type Blockchain struct {
	Name    BlockchainName `json:"name"`
	Address string         `json:"address"`
}
