package ethwallet

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	bindings "github.com/republicprotocol/atom-go/adapters/bindings/eth"
	client "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/domains/match"
	"github.com/republicprotocol/atom-go/services/watch"
)

type ethereumWallet struct {
	wallet *bindings.RenExSettlement
	auth   bind.TransactOpts
	conn   client.Conn
}

func NewEthereumWallet(conn client.Conn, auth bind.TransactOpts) (watch.Wallet, error) {
	wallet, err := bindings.NewRenExSettlement(conn.WalletAddress(), bind.ContractBackend(conn.Client()))

	if err != nil {
		return nil, err
	}
	return &ethereumWallet{
		wallet: wallet,
		auth:   auth,
		conn:   conn,
	}, nil
}

func (wallet *ethereumWallet) GetMatch(personalOrderID [32]byte) (match.Match, error) {
	for {
		time.Sleep(2 * time.Second)
		PersonalOrder, ForeignOrder, ReceiveValue, SendValue, ReceiveCurrency, SendCurrency, err := wallet.wallet.GetMatchDetails(&bind.CallOpts{}, personalOrderID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if PersonalOrder == [32]byte{} {
			continue
		}
		return match.NewMatch(PersonalOrder, ForeignOrder, SendValue, ReceiveValue, SendCurrency, ReceiveCurrency), nil
	}
	// return nil, errors.New("Failed to get match")
}

func (wallet *ethereumWallet) SetMatch(match match.Match) error {
	// wallet.auth.GasLimit = 3000000
	// tx, err := wallet.wallet.SetMatchDetails(&wallet.auth, match.PersonalOrderID(), match.ForeignOrderID(), match.ReceiveCurrency(), match.SendCurrency(), match.ReceiveValue(), match.SendValue())
	// if err != nil {
	// 	return err
	// }
	// _, err = wallet.conn.PatchedWaitMined(context.Background(), tx)
	return nil
}

// func (wallet *ethereumWallet) SubmitOrder(order order.Order) error {
// 	wallet.auth.GasLimit = 3000000

// 	tx, err := wallet.wallet.SubmitOrder(&wallet.auth, order.Type(), order.Parity(), order.Expiry(), order.Tokens(), order.PriceC(), order.PriceQ(), order.VolumeC(), order.VolumeQ(), order.MinimumVolumeC(), order.MinimumVolumeQ(), order.NonceHash())

// }
