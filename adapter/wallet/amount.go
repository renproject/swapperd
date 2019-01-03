package wallet

import (
	"fmt"
	"math/big"

	"github.com/republicprotocol/swapperd/foundation/blockchain"
)

func (wallet *wallet) VerifyBalance(token blockchain.Token, amount *big.Int) error {
	switch token.Name {
	case blockchain.ETH:
		return wallet.verifyEthereumBalance(amount)
	case blockchain.REN, blockchain.DGX, blockchain.TUSD, blockchain.OMG,
		blockchain.ZRX, blockchain.WBTC:
		return wallet.verifyERC20Balance(token.Name, amount)
	case blockchain.BTC:
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

func (wallet *wallet) verifyERC20Balance(tokenName blockchain.TokenName, amount *big.Int) error {

	ethBalance, err := wallet.balance(blockchain.ETH)
	if err != nil {
		return err
	}

	ethAmount, ok := big.NewInt(0).SetString(ethBalance.Amount, 10)
	if !ok {
		return fmt.Errorf("Invalid balance amount: %s", ethBalance.Amount)
	}

	if amount != nil {
		erc20Balance, err := wallet.balance(tokenName)
		if err != nil {
			return err
		}

		erc20Amount, ok := big.NewInt(0).SetString(erc20Balance.Amount, 10)
		if !ok {
			return fmt.Errorf("Invalid balance amount: %s", erc20Balance.Amount)
		}

		if erc20Amount.Cmp(amount) < 0 {
			return fmt.Errorf("You must have at least %s %s remaining in your wallet to execute the swap. You have %s %s", amount, tokenName, erc20Amount, tokenName)
		}
	}

	minVal, ok := big.NewInt(0).SetString("5000000000000000", 10) // 0.005 eth
	if !ok {
		return fmt.Errorf("Invalid minimum value")
	}

	if ethAmount.Cmp(minVal) < 0 {
		return fmt.Errorf("You must have at least 0.005 ETH remaining in your wallet to cover transaction fees. You have %v ETH", ethAmount)
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
		return fmt.Errorf("You need at least 10000 SAT (or 0.0001 BTC) remaining in your wallet to cover transaction fees. You have: %v", balanceAmount)
	}
	return nil
}

func (wallet *wallet) ID() string {
	return wallet.config.IDPublicKey
}
