package wallet

import (
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
)

func (wallet *wallet) VerifyBalance(password string, token blockchain.Token, amount *big.Int) error {
	switch token.Name {
	case blockchain.ETH:
		return wallet.verifyEthereumBalance(password, amount)
	case blockchain.WBTC:
		return wallet.verifyERC20Balance(password, token, amount)
	case blockchain.BTC:
		return wallet.verifyBitcoinBalance(password, amount)
	default:
		return blockchain.NewErrUnsupportedToken("unsupported blockchain")
	}
}

func (wallet *wallet) verifyEthereumBalance(password string, amount *big.Int) error {
	balance, err := wallet.balance(password, blockchain.TokenETH)
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

func (wallet *wallet) verifyERC20Balance(password string, token blockchain.Token, amount *big.Int) error {
	ethBalance, err := wallet.balance(password, blockchain.TokenETH)
	if err != nil {
		return err
	}

	ethAmount, ok := big.NewInt(0).SetString(ethBalance.Amount, 10)
	if !ok {
		return fmt.Errorf("Invalid balance amount: %s", ethBalance.Amount)
	}

	if amount != nil {
		erc20Balance, err := wallet.balance(password, token)
		if err != nil {
			return err
		}

		erc20Amount, ok := big.NewInt(0).SetString(erc20Balance.Amount, 10)
		if !ok {
			return fmt.Errorf("Invalid balance amount: %s", erc20Balance.Amount)
		}

		if erc20Amount.Cmp(amount) < 0 {
			return fmt.Errorf("You must have at least %s %s remaining in your wallet to execute the swap. You have %s %s", amount, token.Name, erc20Amount, token.Name)
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

func (wallet *wallet) verifyBitcoinBalance(password string, amount *big.Int) error {
	if amount == nil {
		return nil
	}

	balance, err := wallet.balance(password, blockchain.TokenBTC)
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

func (wallet *wallet) ID(password string) (string, error) {
	signer, err := wallet.ECDSASigner(password)
	if err != nil {
		return "", nil
	}
	pubKey := signer.PublicKey()
	pubKeyBytes := crypto.FromECDSAPub(&pubKey)
	return base64.StdEncoding.EncodeToString(pubKeyBytes), nil
}
