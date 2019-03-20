package blockchain

import (
	"math/big"

	"github.com/renproject/tokens"
)

type TxExecutionSpeed uint8

const (
	Nil = TxExecutionSpeed(iota)
	Slow
	Standard
	Fast
)

// Cost of an atomic swap
type Cost map[tokens.Name]*big.Int

// TODO: Why this is in foundation?
// CostBlob is the json representation of cost.
type CostBlob map[tokens.Name]string

// CostToCostBlob converts cost to cost blob
func CostToCostBlob(cost Cost) CostBlob {
	costBlob := CostBlob{}
	for tokenName, tokenCost := range cost {
		costBlob[tokenName] = tokenCost.String()
	}
	return costBlob
}

// CostBlobToCost converts cost blob to cost
func CostBlobToCost(costBlob CostBlob) Cost {
	cost := Cost{}
	for tokenName, tokenCost := range costBlob {
		cost[tokenName], _ = new(big.Int).SetString(tokenCost, 10)
	}
	return cost
}
