package wallet

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
)

func (wallet *wallet) Addresses() (map[blockchain.TokenName]string, error) {
	addresses := map[blockchain.TokenName]string{}
	for _, token := range wallet.SupportedTokens() {
		addr, err := wallet.GetAddress(token.Blockchain)
		if err != nil {
			return nil, err
		}
		addresses[token.Name] = addr
	}
	return addresses, nil
}

func (wallet *wallet) GetAddress(blockchainName blockchain.BlockchainName) (string, error) {
	switch blockchainName {
	case blockchain.Ethereum:
		return wallet.config.Ethereum.Address, nil
	case blockchain.Bitcoin:
		return wallet.config.Bitcoin.Address, nil
	default:
		return "", blockchain.NewErrUnsupportedToken("unsupported blockchain")
	}
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
	address = strings.ToLower(address)
	if address[:2] == "0x" {
		address = address[2:]
	}
	addrBytes, err := hex.DecodeString(address)
	if err != nil || len(addrBytes) != 20 {
		return fmt.Errorf("invalid ethereum address: %s", address)
	}
	return nil
}

func (wallet *wallet) verifyBitcoinAddress(address string) error {
	network := wallet.config.Bitcoin.Network.Name
	switch network {
	case "mainnet":
		if _, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams); err != nil {
			return fmt.Errorf("invalid %s bitcoin address: %s", network, address)
		}

	case "testnet":
		if _, err := btcutil.DecodeAddress(address, &chaincfg.TestNet3Params); err != nil {
			return fmt.Errorf("invalid %s bitcoin address: %s", network, address)
		}
	}
	return nil
}
