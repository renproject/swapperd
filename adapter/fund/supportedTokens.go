package fund

import (
	"fmt"

	"github.com/republicprotocol/swapperd/foundation"
)

func decodeSupportedTokens(config Config) ([]foundation.Token, error) {
	supportedTokens := []foundation.Token{}
	for _, blockchain := range config.Blockchains {
		for _, token := range blockchain.Tokens {
			supportedToken, err := foundation.PathToken(token.Name)
			if err != nil {
				return nil, fmt.Errorf("corrupted config file: %v", err)
			}
			supportedTokens = append(supportedTokens, supportedToken)
		}
	}
	return supportedTokens, nil
}

func (manager *manager) SupportedTokens() []foundation.Token {
	return manager.supportedTokens
}
