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
	wallet *bindings.RenExAtomicSettlement
	auth   bind.TransactOpts
	conn   client.Conn
}

func NewEthereumWallet(conn client.Conn, auth bind.TransactOpts) (watch.Wallet, error) {
	wallet, err := bindings.NewRenExAtomicSettlement(conn.WalletAddress(), bind.ContractBackend(conn.Client()))

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
		matchDetails, err := wallet.wallet.GetMatchDetails(&bind.CallOpts{}, personalOrderID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if matchDetails.PersonalOrder == [32]byte{} {
			continue
		}
		return match.NewMatch(matchDetails.PersonalOrder, matchDetails.ForeignOrder, matchDetails.SendValue, matchDetails.RecieveValue, matchDetails.SendCurrency, matchDetails.RecieveCurrency), nil
	}
	// return nil, errors.New("Failed to get match")
}

func (wallet *ethereumWallet) SetMatch(match match.Match) error {
	// wallet.auth.GasLimit = 3000000
	// fmt.Println(match.PersonalOrderID(), match.ForeignOrderID())
	// tx, err := wallet.wallet.SetMatchDetails(&wallet.auth, match.PersonalOrderID(), match.ForeignOrderID(), match.RecieveCurrency(), match.SendCurrency(), match.RecieveValue(), match.SendValue())
	// if err != nil {
	// 	return err
	// }
	// _, err = wallet.conn.PatchedWaitMined(context.Background(), tx)
	// return err
	return nil
}
