package wallet

import "github.com/republicprotocol/swapperd/foundation/blockchain"

type Blockchain struct {
	Name    blockchain.Blockchain
	Address string
}

func (wallet *wallet) SupportedTokens() []blockchain.Token {
	return []blockchain.Token{
		blockchain.TokenBTC,
		blockchain.TokenETH,
		blockchain.TokenWBTC,
	}
}

func (wallet *wallet) SupportedBlockchains() []blockchain.Blockchain {
	return []blockchain.Blockchain{
		blockchain.Blockchain{
			blockchain.Bitcoin,
			wallet.config.Bitcoin.Address,
		},
		blockchain.Blockchain{
			blockchain.Ethereum,
			wallet.config.Ethereum.Address,
		},
	}
}
