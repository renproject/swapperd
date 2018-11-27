package fund

import (
	"fmt"
	"math/big"

	"github.com/republicprotocol/swapperd/foundation"
)

func (manager *manager) VerifyBalance(token foundation.Token, amount *big.Int) error {
	balance, err := manager.balance(token.Name)
	if err != nil {
		return err
	}

	balanceAmount, ok := big.NewInt(0).SetString(balance.Amount, 10)
	if !ok {
		return fmt.Errorf("invalid balance amount %s", balance.Amount)
	}

	leftover := amount.Sub(balanceAmount, amount)
	// leftover amount should be greater than the min tx fees
	switch token.Blockchain {
	case foundation.Ethereum:
		return manager.verifyEthereumBalance(leftover)
	case foundation.Bitcoin:
		return manager.verifyBitcoinBalance(leftover)
	default:
		return foundation.NewErrUnsupportedToken("unsupported blockchain")
	}
}

func (manager *manager) verifyEthereumBalance(leftover *big.Int) error {
	minVal, ok := big.NewInt(0).SetString("5000000000000000", 10) // 0.005 eth
	if !ok {
		return fmt.Errorf("invalid minimum value")
	}
	if leftover.Cmp(minVal) < 0 {
		return fmt.Errorf("minimum balance required to start an atomic swap on ethereum blockchain is 0.005 eth (to cover the transaction fees)")
	}
	return nil
}

func (manager *manager) verifyBitcoinBalance(leftover *big.Int) error {
	if leftover.Cmp(big.NewInt(10000)) < 0 {
		return fmt.Errorf("minimum send amount for bitcoin is 10000 sat")
	}
	return nil
}
