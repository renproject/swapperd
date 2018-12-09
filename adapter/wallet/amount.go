package wallet

import (
	"fmt"
	"math/big"

	"github.com/republicprotocol/swapperd/foundation"
)

func (wallet *wallet) VerifyBalance(token foundation.Token, amount *big.Int) error {
	switch token.Blockchain {
	case foundation.Ethereum:
		return wallet.verifyEthereumBalance()
	case foundation.Bitcoin:

		return wallet.verifyBitcoinBalance(amount)
	default:
		return foundation.NewErrUnsupportedToken("unsupported blockchain")
	}
}

func (wallet *wallet) verifyEthereumBalance() error {
	balance, err := wallet.balance(foundation.ETH)
	if err != nil {
		return err
	}

	balanceAmount, ok := big.NewInt(0).SetString(balance.Amount, 10)
	if !ok {
		return fmt.Errorf("invalid balance amount %s", balance.Amount)
	}

	minVal, ok := big.NewInt(0).SetString("5000000000000000", 10) // 0.005 eth
	if !ok {
		return fmt.Errorf("invalid minimum value")
	}
	if balanceAmount.Cmp(minVal) < 0 {
		return fmt.Errorf("minimum balance required to start an atomic swap on ethereum blockchain is 0.005 eth (to cover the transaction fees) leftover: %v", balanceAmount)
	}
	return nil
}

func (wallet *wallet) verifyBitcoinBalance(amount *big.Int) error {
	balance, err := wallet.balance(foundation.BTC)
	if err != nil {
		return err
	}

	balanceAmount, ok := big.NewInt(0).SetString(balance.Amount, 10)
	if !ok {
		return fmt.Errorf("invalid balance amount %s", balance.Amount)
	}

	leftover := balanceAmount.Sub(balanceAmount, amount)
	if leftover.Cmp(big.NewInt(10000)) < 0 {
		return fmt.Errorf("minimum send amount for bitcoin is 10000 sat")
	}
	return nil
}
