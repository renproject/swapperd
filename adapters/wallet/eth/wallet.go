package ethwallet

import (
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	bindings "github.com/republicprotocol/atom-go/adapters/bindings/eth"
	ethclient "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/services/swap"
	"github.com/republicprotocol/atom-go/services/watch"
)

type EthWallet struct {
	wallet bindings.AtomWallet
	auth   bind.TransactOpts
}

func NewEthereumWallet(conn ethclient.Conn, auth bind.TransactOpts) watch.Wallet {
	wallet, err := bindings.NewAtomWallet(conn.WalletAddress, conn.Client)

	return EthWallet{
		wallet: wallet,
		auth:   auth,
	}
}

func (wallet *EthWallet) WaitForMatch(orderID [32]byte) ([32]byte, error) {
	now := time.Now()
	for time.Now().After(now.Add(48 * time.Hour)) {
		foreignOrderID, err := wallet.wallet.Matches(&bind.CallOpts{}, orderID)
		if foreignOrderID != [32]byte{} {
			return foreignOrderID, nil
		}
		time.Sleep(5 * time.Second)
	}
	return [32]byte{}, errors.New("Order Expired")
}

func (wallet *EthWallet) GetMatch(personalOrderID [32]byte, foreignOrderID [32]byte) (swap.OrderMatch, error) {
	personalOrder, err := wallet.wallet.Orders(&bind.CallOpts{}, personalOrderID)
	foreignOrder, err := wallet.wallet.Orders(&bind.CallOpts{}, foreignOrderID)

	personalPrice, personalVolume, 
}
