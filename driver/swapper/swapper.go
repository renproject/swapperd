package swapper

import (
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	guardianAdapter "github.com/republicprotocol/renex-swapper-go/adapter/guardian"
	"github.com/republicprotocol/renex-swapper-go/adapter/http"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	renexAdapter "github.com/republicprotocol/renex-swapper-go/adapter/renex"
	stateAdapter "github.com/republicprotocol/renex-swapper-go/adapter/state"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	loggerDriver "github.com/republicprotocol/renex-swapper-go/driver/logger"
	"github.com/republicprotocol/renex-swapper-go/driver/network"
	storeDriver "github.com/republicprotocol/renex-swapper-go/driver/store"
	"github.com/republicprotocol/renex-swapper-go/service/guardian"
	"github.com/republicprotocol/renex-swapper-go/service/renex"
	"github.com/republicprotocol/renex-swapper-go/service/state"
)

type Swapper interface {
	Http(port int64)
	Withdraw(tk, to string, value, fee float64) error
}

type swapper struct {
	httpAdapter  http.Adapter
	renexSwapper renex.RenEx
	guardian     guardian.Guardian
	conf         config.Config
	keys         keystore.Keystore
}

func NewSwapper(conf config.Config, keys keystore.Keystore) Swapper {
	db, err := storeDriver.NewLevelDB(conf.HomeDir + "/db")
	if err != nil {
		panic(err)
	}
	logger := loggerDriver.NewStdOut()
	state := state.NewState(stateAdapter.New(db, logger))
	ingressNet := network.NewIngress(conf.RenEx.Ingress, keys.GetKey(token.ETH).(keystore.EthereumKey))
	binder, err := renexAdapter.NewBinder(conf, logger)
	if err != nil {
		panic(err)
	}
	renexSwapper := renex.NewRenEx(renexAdapter.New(conf, keys, ingressNet, state, logger, binder))
	return &swapper{
		httpAdapter:  http.NewAdapter(conf, keys, renexSwapper),
		renexSwapper: renexSwapper,
		guardian:     guardian.NewGuardian(guardianAdapter.New(conf, keys, state, logger)),
		conf:         conf,
		keys:         keys,
	}
}
