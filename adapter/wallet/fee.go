package wallet

import (
	"math/big"

	"github.com/republicprotocol/swapperd/foundation/blockchain"
)

func (wallet *wallet) DefaultFee(blockchainName blockchain.BlockchainName) (*big.Int, error) {
	switch blockchainName {
	case blockchain.Ethereum:
		return big.NewInt(12000000000), nil
	case blockchain.Bitcoin:
		return big.NewInt(10000), nil
	default:
		return nil, blockchain.NewErrUnsupportedBlockchain(blockchainName)
	}
}
