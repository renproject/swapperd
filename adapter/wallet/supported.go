package wallet

import "github.com/renproject/swapperd/foundation/blockchain"

func (wallet *wallet) SupportedTokens() []blockchain.Token {
	return blockchain.SupportedTokens
}
