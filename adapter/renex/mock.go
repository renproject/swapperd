package renex

import (
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/adapter/network"
	swapAdapter "github.com/republicprotocol/renex-swapper-go/adapter/swap"
	"github.com/republicprotocol/renex-swapper-go/adapter/watchdog"
	"github.com/republicprotocol/renex-swapper-go/domain/match"
	"github.com/republicprotocol/renex-swapper-go/domain/order"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/service/logger"
	"github.com/republicprotocol/renex-swapper-go/service/renex"
	"github.com/republicprotocol/renex-swapper-go/service/state"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type Mock struct {
	state.State
	keystore.Keystore
	logger.Logger
	network.Network
	swap.Swapper
	orderbook map[order.ID]match.Match
}

func NewMock(config config.Config, keystore keystore.Keystore, network network.Network, watchdog watchdog.Watchdog, state state.State, logger logger.Logger) (renex.Adapter, error) {
	return &Mock{
		Keystore:  keystore,
		State:     state,
		Logger:    logger,
		Network:   network,
		Swapper:   swap.NewSwapper(swapAdapter.New(config, keystore, network, watchdog, state, logger)),
		orderbook: map[order.ID]match.Match{},
	}, nil
}

// GetOrderMatch gets the order match from the mock renex adapter.
func (renex *Mock) GetOrderMatch(orderID order.ID, waitTill int64) (match.Match, error) {
	return renex.orderbook[orderID], nil
}

// SubmitOrderMatch subbmits an order match to the mock renex adapter.
func (renex *Mock) SubmitOrderMatch(match match.Match) error {
	renex.orderbook[match.PersonalOrderID()] = match
	return nil
}

// GetAddress corresponding to the given token.
func (renex *Mock) GetAddress(tok token.Token) []byte {
	switch tok {
	case token.BTC:
		return []byte(renex.Keystore.GetKey(tok).(keystore.BitcoinKey).AddressString)
	case token.ETH:
		return []byte(renex.Keystore.GetKey(tok).(keystore.EthereumKey).Address.String())
	default:
		panic("Unexpected token")
	}
}
