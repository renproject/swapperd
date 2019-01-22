package blockchain

import "math/big"

// Cost of an atomic swap
type Cost map[TokenName]*big.Int

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

var (
	EthereumTransactionCost = map[TokenName]*big.Int{
		ETH: big.NewInt(12000000000000000),
	}

	BitcoinTransactionCost = map[TokenName]*big.Int{
		BTC: big.NewInt(10000),
	}
)

func (token Token) TransactionCost(amount *big.Int) (Cost, error) {
	switch token.Blockchain {
	case ERC20:
		return erc20TransactionCost(token, amount), nil
	case Ethereum:
		return EthereumTransactionCost, nil
	case Bitcoin:
		return BitcoinTransactionCost, nil
	default:
		return nil, NewErrUnsupportedToken(token.Name)
	}
}

func (token Token) BlockchainTxFees() (*big.Int, error) {
	switch token.Blockchain {
	case Ethereum, ERC20:
		return big.NewInt(12000000000), nil
	case Bitcoin:
		return big.NewInt(10000), nil
	default:
		return nil, NewErrUnsupportedToken(token.Name)
	}
}

func erc20TransactionCost(token Token, amount *big.Int) Cost {
	switch token {
	case TokenDGX:
		cost := EthereumTransactionCost
		cost[DGX] = calculateFeesFromBips(amount, 13)
		return cost
	default:
		return EthereumTransactionCost
	}
}

func calculateFeesFromBips(value *big.Int, bips int64) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(value, big.NewInt(bips)), new(big.Int).Sub(big.NewInt(10000), big.NewInt(bips)))
}
