package wallet

import (
	"fmt"
	"sync"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/renproject/libeth-go"
	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/republicprotocol/co-go"
)

func (wallet *wallet) Addresses(password string) (map[blockchain.TokenName]string, error) {
	addresses := map[blockchain.TokenName]string{}
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

func (wallet *wallet) GetAddress(password string, blockchainName blockchain.BlockchainName) (string, error) {
	switch blockchainName {
	case blockchain.Ethereum, blockchain.ERC20:
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

func (wallet *wallet) ResolveAddress(bcName blockchain.BlockchainName, address string) (string, error) {
	switch bcName {
	case blockchain.Ethereum, blockchain.ERC20:
		return wallet.resolveEthereumAddress(address)
	case blockchain.Bitcoin:
		return wallet.resolveBitcoinAddress(address)
	default:
		return "", blockchain.NewErrUnsupportedBlockchain(bcName)
	}
}

func (wallet *wallet) resolveEthereumAddress(address string) (string, error) {
	conn, err := libeth.Connect(wallet.config.Ethereum.Network.URL)
	if err != nil {
		return "", err
	}
	addr, err := conn.Resolve(address)
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

func (wallet *wallet) resolveBitcoinAddress(address string) (string, error) {
	if address == "" {
		return "", fmt.Errorf("Empty bitcoin address")
	}

	network := wallet.config.Bitcoin.Network.Name
	switch network {
	case "mainnet":
		if _, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams); err != nil {
			return "", fmt.Errorf("Invalid %s bitcoin address: %s", network, address)
		}

	case "testnet":
		if _, err := btcutil.DecodeAddress(address, &chaincfg.TestNet3Params); err != nil {
			return "", fmt.Errorf("Invalid %s bitcoin address: %s", network, address)
		}
	}
	return address, nil
}
