package wallet

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
)

func (wallet *wallet) Addresses(password string) (map[blockchain.TokenName]string, error) {
	addresses := map[blockchain.TokenName]string{}
	for _, token := range wallet.SupportedTokens() {
		addr, err := wallet.GetAddress(password, token.Blockchain)
		if err != nil {
			return nil, err
		}
		addresses[token.Name] = addr
	}
	return addresses, nil
}

func (wallet *wallet) GetAddress(password string, blockchainName blockchain.BlockchainName) (string, error) {
	switch blockchainName {
	case blockchain.Ethereum:
		return wallet.getEthereumAddress(password)
	case blockchain.Bitcoin:
		return wallet.getBitcoinAddress(password)
	default:
		return "", blockchain.NewErrUnsupportedToken("unsupported blockchain")
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

func (wallet *wallet) VerifyAddress(blockchainName blockchain.BlockchainName, address string) error {
	switch blockchainName {
	case blockchain.Ethereum:
		return wallet.verifyEthereumAddress(address)
	case blockchain.Bitcoin:
		return wallet.verifyBitcoinAddress(address)
	default:
		return blockchain.NewErrUnsupportedToken("unsupported blockchain")
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
		if _, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams); err != nil {
			return fmt.Errorf("Invalid %s bitcoin address: %s", network, address)
		}

	case "testnet":
		if _, err := btcutil.DecodeAddress(address, &chaincfg.TestNet3Params); err != nil {
			return fmt.Errorf("Invalid %s bitcoin address: %s", network, address)
		}
	}
	return nil
}
