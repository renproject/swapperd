package renex

import (
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	swapAdapter "github.com/republicprotocol/renex-swapper-go/adapter/swap"
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
	renex.Network
	Binder
}

func New(config config.Config, keystore keystore.Keystore, network renex.Network, state state.State, logger logger.Logger, binder Binder) renex.Adapter {
	return &renexAdapter{
		Keystore: keystore,
		State:    state,
		Logger:   logger,
		Network:  network,
		Binder:   binder,
		Swapper:  swap.NewSwapper(swapAdapter.New(config, keystore, logger)),
	}
}

// GetAddress corresponding to the given token.
func (renex *renexAdapter) GetAddresses(tok1, tok2 token.Token) (string, string) {
	return renex.getAddress(tok1), renex.getAddress(tok2)
}

func (renex *renexAdapter) getAddress(tok token.Token) string {
	switch tok {
	case token.BTC:
		return renex.Keystore.GetKey(tok).(keystore.BitcoinKey).AddressString
	case token.ETH:
		return renex.Keystore.GetKey(tok).(keystore.EthereumKey).Address.String()
	default:
		panic("Unexpected token")
	}
}
