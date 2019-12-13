package blockchain

import "github.com/renproject/tokens"

type Blockchain struct {
	Name    tokens.BlockchainName `json:"name"`
	Address string                `json:"address"`
}

type Balance struct {
	Address    string `json:"address"`
	Decimals   int    `json:"decimals"`
	FullAmount string `json:"fullBalance"`
	Amount     string `json:"balance"`
}
