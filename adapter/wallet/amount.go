package wallet

import (
	"fmt"
	"math/big"

	"github.com/republicprotocol/swapperd/foundation/blockchain"
)

func (wallet *wallet) VerifyBalance(token blockchain.Token, amount *big.Int) error {
	switch token.Blockchain {
	case blockchain.Ethereum:
		return wallet.verifyEthereumBalance()
	case blockchain.Bitcoin:

		return wallet.verifyBitcoinBalance(amount)
	default:
		return blockchain.NewErrUnsupportedToken("unsupported blockchain")
	}
}

func (wallet *wallet) verifyEthereumBalance() error {
	balance, err := wallet.balance(blockchain.ETH)
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
	balance, err := wallet.balance(blockchain.BTC)
	if err != nil {
		return err
	}

	balanceAmount, ok := big.NewInt(0).SetString(balance.Amount, 10)
	if !ok {
		return fmt.Errorf("invalid balance amount %s", balance.Amount)
	}

	leftover := balanceAmount.Sub(balanceAmount, amount)
	if leftover.Cmp(big.NewInt(10000)) < 0 {
		return fmt.Errorf("not enough bitcoin to complete the swap have: %v need: %v", balanceAmount, amount)
	}
	return nil
}

func (wallet *wallet) ID() string {
	return wallet.config.IDPublicKey
}
