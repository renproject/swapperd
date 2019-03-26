package wallet

import (
	"fmt"
	"math/big"

	"github.com/renproject/tokens"
)

func (wallet *wallet) VerifyBalance(password string, token tokens.Token, amount *big.Int) error {
	switch token.Blockchain {
	case tokens.ETHEREUM:
		return wallet.verifyEthereumBalance(password, amount)
	case tokens.ERC20:
		return wallet.verifyERC20Balance(password, token, amount)
	case tokens.BITCOIN:
		return wallet.verifyBitcoinBalance(password, amount)
	case tokens.ZCASH:
		return wallet.verifyZCashBalance(password, amount)
	default:
		return tokens.NewErrUnsupportedBlockchain(token.Blockchain)
	}
}

func (wallet *wallet) verifyEthereumBalance(password string, amount *big.Int) error {
	balanceAmount, err := wallet.AvailableBalance(password, tokens.ETH)
	if err != nil {
		return err
	}

	if amount != nil {
		balanceAmount = new(big.Int).Sub(balanceAmount, amount)
	}
	fee := big.NewInt(1200000000000000)
	if balanceAmount.Cmp(fee) < 0 {
		return fmt.Errorf("You must have at least %d WEI remaining in your wallet to cover transaction fees. You have %v WEI", fee, balanceAmount)
	}
	return nil
}

func (wallet *wallet) verifyERC20Balance(password string, token tokens.Token, amount *big.Int) error {
	ethAmount, err := wallet.AvailableBalance(password, tokens.ETH)
	if err != nil {
		return err
	}

	fee := big.NewInt(1200000000000000)
	if amount != nil {
		erc20Amount, err := wallet.AvailableBalance(password, token)
		if err != nil {
			return err
		}

		expectedAmount := amount
		if extraFee := token.AdditionalTransactionFee(erc20Amount); extraFee != nil {
			expectedAmount = new(big.Int).Add(amount, extraFee)
		}

		if erc20Amount.Cmp(expectedAmount) < 0 {
			return fmt.Errorf("You must have at least %s %s remaining in your wallet to execute the swap. You have %s %s", expectedAmount, token.Name, erc20Amount, token.Name)
		}
	}

	if ethAmount.Cmp(fee) < 0 {
		return fmt.Errorf("You must have at least %d WEI remaining in your wallet to cover transaction fees. You have %v WEI", fee, ethAmount)
	}

	return nil
}

func (wallet *wallet) verifyBitcoinBalance(password string, amount *big.Int) error {
	if amount == nil {
		return nil
	}

	if amount.Cmp(big.NewInt(20000)) < 0 {
		return fmt.Errorf("invalid bitcoin amount: minimum swappable bitcoin amount 20000 SAT (or 0.0002 BTC)")
	}

	fee := big.NewInt(10000)

	balanceAmount, err := wallet.AvailableBalance(password, tokens.BTC)
	if err != nil {
		return err
	}

	leftover := new(big.Int).Sub(balanceAmount, amount)
	if leftover.Cmp(new(big.Int).Add(fee, big.NewInt(600))) < 0 {
		return fmt.Errorf("You need at least 10600 SAT (or 0.000106 BTC) remaining in your wallet to cover transaction fees. You have: %v", leftover)
	}
	return nil
}

func (wallet *wallet) verifyZCashBalance(password string, amount *big.Int) error {
	if amount == nil {
		return nil
	}

	if amount.Cmp(big.NewInt(20000)) < 0 {
		return fmt.Errorf("invalid ZCash amount: minimum swappable ZCash amount 20000 ZAT (or 0.0002 ZEC)")
	}

	fee := big.NewInt(10000)

	balanceAmount, err := wallet.AvailableBalance(password, tokens.ZEC)
	if err != nil {
		return err
	}

	leftover := new(big.Int).Sub(balanceAmount, amount)
	if leftover.Cmp(new(big.Int).Add(fee, big.NewInt(600))) < 0 {
		return fmt.Errorf("You need at least 10600 ZAT (or 0.000106 ZEC) remaining in your wallet to cover transaction fees. You have: %v", leftover)
	}
	return nil
}
