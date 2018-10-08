package swapper

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	netHttp "net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/renex-swapper-go/adapter/btc"
	"github.com/republicprotocol/renex-swapper-go/adapter/eth"

	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	guardianAdapter "github.com/republicprotocol/renex-swapper-go/adapter/guardian"
	"github.com/republicprotocol/renex-swapper-go/adapter/http"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	renexAdapter "github.com/republicprotocol/renex-swapper-go/adapter/renex"
	stateAdapter "github.com/republicprotocol/renex-swapper-go/adapter/state"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	httpDriver "github.com/republicprotocol/renex-swapper-go/driver/http"
	loggerDriver "github.com/republicprotocol/renex-swapper-go/driver/logger"
	"github.com/republicprotocol/renex-swapper-go/driver/network"
	storeDriver "github.com/republicprotocol/renex-swapper-go/driver/store"
	"github.com/republicprotocol/renex-swapper-go/service/guardian"
	"github.com/republicprotocol/renex-swapper-go/service/renex"
	"github.com/republicprotocol/renex-swapper-go/service/state"
)

type Swapper interface {
	Http(port int64)
	Withdraw(tk, to string, value float64) error
}

type swapper struct {
	httpAdapter  http.Adapter
	renexSwapper renex.RenEx
	guardian     guardian.Guardian
	conf         config.Config
	keys         keystore.Keystore
}

func NewSwapper(conf config.Config, keys keystore.Keystore) Swapper {
	return &swapper{
		conf: conf,
		keys: keys,
	}
}

func (swapper *swapper) Http(port int64) {
	db, err := storeDriver.NewLevelDB(swapper.conf.HomeDir)
	if err != nil {
		panic(err)
	}
	logger := loggerDriver.NewStdOut()
	state := state.NewState(stateAdapter.New(db, logger))
	ingressNet := network.NewIngress(swapper.conf.RenEx.Ingress, swapper.keys.GetKey(token.ETH).(keystore.EthereumKey))
	binder, err := renexAdapter.NewBinder(swapper.conf, logger)
	if err != nil {
		panic(err)
	}
	renexSwapper := renex.NewRenEx(renexAdapter.New(swapper.conf, swapper.keys, ingressNet, state, logger, binder))
	guardian := guardian.NewGuardian(guardianAdapter.New(swapper.conf, swapper.keys, state, logger))
	errCh := make(chan error, 1)
	go renexSwapper.Run(errCh)
	go guardian.Run(errCh)
	go func() {
		for err := range errCh {
			fmt.Println("Swapper Error: ", err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		defer close(errCh)
		_ = <-c
		log.Println("Stopping the atom box safely")
		os.Exit(1)
	}()
	log.Fatal(netHttp.ListenAndServe(fmt.Sprintf(":%d", port), httpDriver.NewServer(http.NewAdapter(swapper.conf, swapper.keys, renexSwapper))))
}

func (swapper *swapper) Withdraw(tokenStr, to string, value float64) error {
	// Parse and validate the token
	tokenStr = strings.ToLower(strings.TrimSpace(tokenStr))
	switch tokenStr {
	case "btc", "bitcoin", "xbt":
		valueBig, _ := big.NewFloat(value * math.Pow10(8)).Int(nil)
		return swapper.withdrawBitcoin(to, valueBig.Int64())
	case "eth", "ethereum", "ether":
		valueBig, _ := big.NewFloat(value * math.Pow10(18)).Int(nil)
		return swapper.withdrawEthereum(to, valueBig)
	default:
		return errors.New("unknown token")
	}
}

func (swapper *swapper) withdrawBitcoin(to string, value int64) error {
	conn := btc.NewConnWithConfig(swapper.conf.Bitcoin)
	return conn.Withdraw(to, swapper.keys.GetKey(token.BTC).(keystore.BitcoinKey), value, 3000)
}

func (swapper *swapper) withdrawEthereum(to string, value *big.Int) error {
	conn, err := eth.NewConnWithConfig(swapper.conf.Ethereum)
	if err != nil {
		return err
	}
	return conn.Transfer(common.HexToAddress(to), swapper.keys.GetKey(token.ETH).(keystore.EthereumKey), value)
}
