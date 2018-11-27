package fund

import (
	"github.com/republicprotocol/swapperd/foundation"
)

type Blockchain struct {
	Name    foundation.Blockchain
	Address string
}

func (manager *manager) SupportedTokens() []foundation.Token {
	return []foundation.Token{
		foundation.TokenBTC,
		foundation.TokenETH,
		foundation.TokenWBTC,
	}
}

func (manager *manager) SupportedBlockchains() []foundation.Blockchain {
	return []foundation.Blockchain{
		foundation.Blockchain{
			foundation.Bitcoin,
			manager.config.Bitcoin.Address,
		},
		foundation.Blockchain{
			foundation.Ethereum,
			manager.config.Ethereum.Address,
		},
	}
}
