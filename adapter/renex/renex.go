package renex

import (
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/adapter/network"
	swapAdapter "github.com/republicprotocol/renex-swapper-go/adapter/swap"
	"github.com/republicprotocol/renex-swapper-go/adapter/watchdog"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/service/logger"
	"github.com/republicprotocol/renex-swapper-go/service/renex"
	"github.com/republicprotocol/renex-swapper-go/service/state"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type renexAdapter struct {
	state.State
	keystore.Keystore
	logger.Logger
	swap.Swapper
	network.Network
	Binder
}

func New(config config.Config, keystore keystore.Keystore, network network.Network, watchdog watchdog.Watchdog, state state.State, logger logger.Logger, binder Binder) renex.Adapter {
	return &renexAdapter{
		Keystore: keystore,
		State:    state,
		Logger:   logger,
		Network:  network,
		Binder:   binder,
		Swapper:  swap.NewSwapper(swapAdapter.New(config, keystore, network, watchdog, state, logger)),
	}
}

// GetAddress corresponding to the given token.
func (renex *renexAdapter) GetAddress(tok token.Token) []byte {
	switch tok {
	case token.BTC:
		return []byte(renex.Keystore.GetKey(tok).(keystore.BitcoinKey).AddressString)
	case token.ETH:
		return []byte(renex.Keystore.GetKey(tok).(keystore.EthereumKey).Address.String())
	default:
		panic("Unexpected token")
	}
}
