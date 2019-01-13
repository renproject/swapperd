package wallet

import "github.com/republicprotocol/swapperd/foundation/blockchain"

func (wallet *wallet) SupportedTokens() []blockchain.Token {
	return []blockchain.Token{
		blockchain.TokenBTC,
		blockchain.TokenETH,
		blockchain.TokenWBTC,
		blockchain.TokenREN,
		blockchain.TokenZRX,
		blockchain.TokenOMG,
		blockchain.TokenTUSD,
		blockchain.TokenDGX,
		blockchain.TokenGUSD,
		blockchain.TokenDAI,
		blockchain.TokenUSDC,
	}
}
