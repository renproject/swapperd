package wallet

import (
	"encoding/hex"
	"fmt"
	"strings"
	"sync"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/renproject/libzec-go"
	"github.com/renproject/tokens"
	"github.com/republicprotocol/co-go"
)

func (wallet *wallet) Addresses(password string) (map[tokens.Name]string, error) {
	addresses := map[tokens.Name]string{}
	mu := new(sync.RWMutex)
	tokens := wallet.SupportedTokens()
	co.ParForAll(tokens, func(i int) {
		token := tokens[i]
		addr, err := wallet.GetAddress(password, token.Blockchain)
		if err != nil {
			return
		}
		mu.Lock()
		defer mu.Unlock()
		addresses[token.Name] = addr
	})
	return addresses, nil
}

func (wallet *wallet) GetAddress(password string, blockchainName tokens.BlockchainName) (string, error) {
	switch blockchainName {
	case tokens.ETHEREUM, tokens.ERC20:
		return wallet.getEthereumAddress(password)
	case tokens.BITCOIN:
		return wallet.getBitcoinAddress(password)
	case tokens.ZCASH:
		return wallet.getZCashAddress(password)
	default:
		return "", tokens.NewErrUnsupportedBlockchain(blockchainName)
	}
}

func (wallet *wallet) getEthereumAddress(password string) (string, error) {
	ethAccount, err := wallet.EthereumAccount(password)
	if err != nil {
		return "", err
	}
	return ethAccount.Address().String(), nil
}

func (wallet *wallet) getBitcoinAddress(password string) (string, error) {
	btcAccount, err := wallet.BitcoinAccount(password)
	if err != nil {
		return "", err
	}
	btcAddr, err := btcAccount.Address()
	if err != nil {
		return "", err
	}
	return btcAddr.String(), nil
}

func (wallet *wallet) getZCashAddress(password string) (string, error) {
	zecAccount, err := wallet.ZCashAccount(password)
	if err != nil {
		return "", err
	}
	zecAddr, err := zecAccount.Address()
	if err != nil {
		return "", err
	}
	return zecAddr.String(), nil
}

func (wallet *wallet) VerifyAddress(blockchainName tokens.BlockchainName, address string) error {
	switch blockchainName {
	case tokens.ETHEREUM, tokens.ERC20:
		return wallet.verifyEthereumAddress(address)
	case tokens.BITCOIN:
		return wallet.verifyBitcoinAddress(address)
	case tokens.ZCASH:
		return wallet.verifyZCashAddress(address)
	default:
		return tokens.NewErrUnsupportedBlockchain(blockchainName)
	}
}

func (wallet *wallet) verifyEthereumAddress(address string) error {
	if address == "" {
		return fmt.Errorf("Empty ethereum address")
	}

	address = strings.ToLower(address)
	if len(address) > 2 && address[:2] == "0x" {
		address = address[2:]
	}

	addrBytes, err := hex.DecodeString(address)
	if err != nil || len(addrBytes) != 20 {
		return fmt.Errorf("Invalid ethereum address: %s", address)
	}
	return nil
}

func (wallet *wallet) verifyBitcoinAddress(address string) error {
	if address == "" {
		return fmt.Errorf("Empty bitcoin address")
	}

	network := wallet.config.Bitcoin.Network.Name
	switch network {
	case "mainnet":
		if _, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams); err == nil {
			return nil
		}

	case "testnet3", "testnet", "regtest":
		if _, err := btcutil.DecodeAddress(address, &chaincfg.TestNet3Params); err == nil {
			return nil
		}
	}
	return fmt.Errorf("Invalid %s bitcoin address: %s", network, address)
}

func (wallet *wallet) verifyZCashAddress(address string) error {
	if address == "" {
		return fmt.Errorf("Empty ZCash address")
	}

	network := wallet.config.ZCash.Network.Name
	switch network {
	case "mainnet":
		if _, err := libzec.DecodeAddress(address, &chaincfg.MainNetParams); err == nil {
			return nil
		}

	case "testnet3", "testnet", "regtest":
		if _, err := libzec.DecodeAddress(address, &chaincfg.TestNet3Params); err == nil {
			return nil
		}
	}
	return fmt.Errorf("Invalid %s ZCash address: %s", network, address)
}
