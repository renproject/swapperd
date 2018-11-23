package fund

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"

	"github.com/republicprotocol/swapperd/foundation"
)

func (manager *manager) VerifyAddress(blockchain foundation.BlockchainName, address string) error {
	switch blockchain {
	case foundation.Ethereum:
		return manager.verifyEthereumAddress(address)
	case foundation.Bitcoin:
		return manager.verifyBitcoinAddress(address)
	default:
		return foundation.NewErrUnsupportedToken("unsupported blockchain")
	}
}

func (manager *manager) verifyEthereumAddress(address string) error {
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

func (manager *manager) verifyBitcoinAddress(address string) error {
	network := manager.config.Bitcoin.Network.Name
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
