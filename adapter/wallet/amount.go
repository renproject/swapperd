package wallet

import (
	"fmt"
	"math/big"

	"github.com/republicprotocol/swapperd/foundation/blockchain"
)

func (wallet *wallet) VerifyBalance(token blockchain.Token, amount *big.Int) error {
	switch token.Blockchain {
	case blockchain.Ethereum:
		return wallet.verifyEthereumBalance(amount)
	case blockchain.Bitcoin:
		return wallet.verifyBitcoinBalance(amount)
	default:
		return blockchain.NewErrUnsupportedToken("unsupported blockchain")
	}
}

func (wallet *wallet) verifyEthereumBalance(amount *big.Int) error {
	balance, err := wallet.balance(blockchain.ETH)
	if err != nil {
		return err
	}

	balanceAmount, ok := big.NewInt(0).SetString(balance.Amount, 10)
	if !ok {
		return fmt.Errorf("Invalid balance amount: %s", balance.Amount)
	}

	if amount != nil {
		balanceAmount = new(big.Int).Sub(balanceAmount, amount)
	}

	minVal, ok := big.NewInt(0).SetString("5000000000000000", 10) // 0.005 eth
	if !ok {
		return fmt.Errorf("Invalid minimum value")
	}
	if balanceAmount.Cmp(minVal) < 0 {
		return fmt.Errorf("You must have at least 0.005 ETH remaining in your wallet to cover transaction fees. You have %v ETH", balanceAmount)
	}
	return nil
}

func (wallet *wallet) verifyBitcoinBalance(amount *big.Int) error {
	if amount == nil {
		return nil
	}

	balance, err := wallet.balance(blockchain.BTC)
	if err != nil {
		return err
	}

	balanceAmount, ok := big.NewInt(0).SetString(balance.Amount, 10)
	if !ok {
		return fmt.Errorf("Invalid balance amount: %s", balance.Amount)
	}

	leftover := balanceAmount.Sub(balanceAmount, amount)
	if leftover.Cmp(big.NewInt(10000)) < 0 {
		return fmt.Errorf("You need at least %v BTC remaining in your wallet to cover transaction fees. You have: %v", amount, balanceAmount)
	}
	return nil
}

func (wallet *wallet) ID() string {
	return wallet.config.IDPublicKey
}
