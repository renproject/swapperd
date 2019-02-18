package blockchain

import "math/big"

// Cost of an atomic swap
type Cost map[TokenName]*big.Int

// TODO: Why this is in foundation?
// CostBlob is the json representation of cost.
type CostBlob map[TokenName]string

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

func (token Token) AdditionalTransactionFee(amount *big.Int) *big.Int {
	switch token {
	case TokenDGX:
		return calculateFeesFromBips(amount, 13)
	default:
		return nil
	}
}

func calculateFeesFromBips(value *big.Int, bips int64) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(value, big.NewInt(bips)), new(big.Int).Sub(big.NewInt(10000), big.NewInt(bips)))
}
