package wallet

import (
	"github.com/republicprotocol/swapperd/foundation"
)

type Blockchain struct {
	Name    foundation.Blockchain
	Address string
}

func (wallet *wallet) SupportedTokens() []foundation.Token {
	return []foundation.Token{
		foundation.TokenBTC,
		foundation.TokenETH,
		foundation.TokenWBTC,
	}
}

func (wallet *wallet) SupportedBlockchains() []foundation.Blockchain {
	return []foundation.Blockchain{
		foundation.Blockchain{
			foundation.Bitcoin,
			wallet.config.Bitcoin.Address,
		},
		foundation.Blockchain{
			foundation.Ethereum,
			wallet.config.Ethereum.Address,
		},
	}
}
