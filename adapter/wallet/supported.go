package wallet

import "github.com/republicprotocol/swapperd/foundation/blockchain"

func (wallet *wallet) SupportedTokens() []blockchain.Token {
	return blockchain.SupportedTokens
}
