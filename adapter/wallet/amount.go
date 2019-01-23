package wallet

import (
	"fmt"
	"math/big"

	"github.com/republicprotocol/swapperd/foundation/blockchain"
)

func (wallet *wallet) VerifyBalance(password string, token blockchain.Token, amount *big.Int) error {
	switch token.Blockchain {
	case blockchain.Ethereum:
		return wallet.verifyEthereumBalance(password, amount)
	case blockchain.ERC20:
		return wallet.verifyERC20Balance(password, token, amount)
	case blockchain.Bitcoin:
		return wallet.verifyBitcoinBalance(password, amount)
	default:
		return blockchain.NewErrUnsupportedToken("unsupported blockchain")
	}
}

func (wallet *wallet) verifyEthereumBalance(password string, amount *big.Int) error {
	balance, err := wallet.Balance(password, blockchain.TokenETH)
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

	fee, err := blockchain.TokenETH.TransactionCost(amount)
	if err != nil {
		return err
	}

	if balanceAmount.Cmp(fee[blockchain.ETH]) < 0 {
		return fmt.Errorf("You must have at least %d WEI remaining in your wallet to cover transaction fees. You have %v WEI", fee[blockchain.ETH], balanceAmount)
	}
	return nil
}

func (wallet *wallet) verifyERC20Balance(password string, token blockchain.Token, amount *big.Int) error {
	ethBalance, err := wallet.Balance(password, blockchain.TokenETH)
	if err != nil {
		return err
	}

	ethAmount, ok := big.NewInt(0).SetString(ethBalance.Amount, 10)
	if !ok {
		return fmt.Errorf("Invalid balance amount: %s", ethBalance.Amount)
	}

	fee, err := blockchain.TokenETH.TransactionCost(amount)
	if err != nil {
		return err
	}

	if amount != nil {
		erc20Balance, err := wallet.Balance(password, token)
		if err != nil {
			return err
		}

		erc20Amount, ok := big.NewInt(0).SetString(erc20Balance.Amount, 10)
		if !ok {
			return fmt.Errorf("Invalid balance amount: %s", erc20Balance.Amount)
		}

		feeValue, ok := fee[token.Name]
		if !ok {
			feeValue = big.NewInt(0)
		}

		expectedAmount := new(big.Int).Add(amount, feeValue)
		if erc20Amount.Cmp(expectedAmount) < 0 {
			return fmt.Errorf("You must have at least %s %s remaining in your wallet to execute the swap. You have %s %s", expectedAmount, token.Name, erc20Amount, token.Name)
		}
	}

	if ethAmount.Cmp(fee[blockchain.ETH]) < 0 {
		return fmt.Errorf("You must have at least %d WEI remaining in your wallet to cover transaction fees. You have %v WEI", fee[blockchain.ETH], ethAmount)
	}

	return nil
}

func (wallet *wallet) verifyBitcoinBalance(password string, amount *big.Int) error {
	if amount == nil {
		return nil
	}

	fee, err := blockchain.TokenBTC.TransactionCost(amount)
	if err != nil {
		return err
	}

	balance, err := wallet.Balance(password, blockchain.TokenBTC)
	if err != nil {
		return err
	}

	balanceAmount, ok := big.NewInt(0).SetString(balance.Amount, 10)
	if !ok {
		return fmt.Errorf("Invalid balance amount: %s", balance.Amount)
	}

	leftover := balanceAmount.Sub(balanceAmount, amount)
	if leftover.Cmp(fee[blockchain.BTC]) < 0 {
		return fmt.Errorf("You need at least 10000 SAT (or 0.0001 BTC) remaining in your wallet to cover transaction fees. You have: %v", balanceAmount)
	}
	return nil
}
