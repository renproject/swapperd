package regtest

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/btcsuite/btcutil"
	btcclient "github.com/republicprotocol/atom-go/adapters/clients/btc"
)

// Start a local Ganache instance.
func Start() *exec.Cmd {
	cmd := exec.Command("bitcoind", "--regtest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
	return cmd
}

func Stop(cmd *exec.Cmd) {
	cmd.Process.Signal(syscall.SIGTERM)
	cmd.Wait()
}

func Mine(connection btcclient.Conn) error {
	_, err := connection.Client.Generate(100)
	if err != nil {
		return err
	}
	tick := time.NewTicker(2 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			_, err := connection.Client.Generate(1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewAccount(connection btcclient.Conn, name string, value btcutil.Amount) (btcutil.Address, error) {
	addr, err := connection.Client.GetAccountAddress(name)
	if err != nil {
		return nil, err
	}

	if value > 0 {
		_, err = connection.Client.SendToAddress(addr, value)
		if err != nil {
			return nil, err
		}

		_, err = connection.Client.Generate(1)
		if err != nil {
			return nil, err
		}
	}

	return addr, nil
}

func GetAddressForAccount(connection btcclient.Conn, name string) (btcutil.Address, error) {
	addresses, err := connection.Client.GetAddressesByAccount(name)
	newAddress := addresses[2]
	balance, err := connection.Client.GetReceivedByAddress(newAddress)

	fmt.Println("Balance in regtest :", balance)
	if balance == 0 {
		amt := btcutil.Amount(1000000000)
		log.Println(amt.ToBTC())
		tx, err := connection.Client.SendFrom("", newAddress, amt)
		if err != nil {
			return nil, err
		}

		_, err = connection.Client.Generate(1)
		log.Println("Mined a block")
		if err != nil {
			return nil, err
		}

		err = connection.WaitForConfirmations(tx, 10)
	}

	return addresses[2], err
}

// func NewAccount(conn client.btcclient.Conn, eth *big.Int) (*bind.TransactOpts, common.Address, error) {
// 	ethereumPair, err := crypto.GenerateKey()
// 	if err != nil {
// 		return nil, common.Address{}, err
// 	}
// 	addr := crypto.PubkeyToAddress(ethereumPair.PublicKey)
// 	account := bind.NewKeyedTransactor(ethereumPair)
// 	if eth.Cmp(big.NewInt(0)) > 0 {
// 		if err := DistributeEth(conn, addr); err != nil {
// 			return nil, common.Address{}, err
// 		}
// 	}

// 	return account, addr, nil
// }
