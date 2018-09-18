package swap

import (
	"fmt"

	"github.com/republicprotocol/renex-swapper-go/domain/match"

	"github.com/republicprotocol/renex-swapper-go/adapter/atoms/btc"
	"github.com/republicprotocol/renex-swapper-go/adapter/atoms/eth"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/adapter/network"
	"github.com/republicprotocol/renex-swapper-go/adapter/watchdog"
	"github.com/republicprotocol/renex-swapper-go/domain/order"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/service/logger"
	"github.com/republicprotocol/renex-swapper-go/service/state"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type swapperAdapter struct {
	network.Network
	watchdog.Watchdog
	state.State
	logger.Logger
	config   config.Config
	keystore keystore.Keystore
}

func New(cfg config.Config, ks keystore.Keystore, network network.Network, watchdog watchdog.Watchdog, state state.State, logger logger.Logger) swap.SwapperAdapter {
	return &swapperAdapter{
		config:   cfg,
		keystore: ks,
		Network:  network,
		Watchdog: watchdog,
		State:    state,
		Logger:   logger,
	}
}

func (swapper *swapperAdapter) NewSwap(orderID order.ID) (swap.Atom, swap.Atom, match.Match, swap.Adapter, error) {
	var personalAtom, foreignAtom swap.Atom
	match, err := swapper.Match(orderID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	personalAtom, err = buildAtom(swapper.Network, swapper.keystore, swapper.config, match.SendCurrency(), match.PersonalOrderID())
	if err != nil {
		return nil, nil, nil, nil, err
	}

	foreignAtom, err = buildAtom(swapper.Network, swapper.keystore, swapper.config, match.ReceiveCurrency(), match.ForeignOrderID())
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if swapper.AtomExists(match.PersonalOrderID()) {
		details, err := swapper.AtomDetails(match.PersonalOrderID())
		if err != nil {
			return nil, nil, nil, nil, err
		}
		if err := personalAtom.Deserialize(details); err != nil {
			return nil, nil, nil, nil, err
		}
	}

	if swapper.AtomExists(match.ForeignOrderID()) {
		details, err := swapper.AtomDetails(match.ForeignOrderID())
		if err != nil {
			return nil, nil, nil, nil, err
		}
		if err := foreignAtom.Deserialize(details); err != nil {
			return nil, nil, nil, nil, err
		}
	}

	return personalAtom, foreignAtom, match, swapper, nil
}

func buildAtom(network network.Network, key keystore.Keystore, config config.Config, t token.Token, orderID [32]byte) (swap.Atom, error) {
	switch t {
	case token.BTC:
		btcKey := key.GetKey(t).(keystore.BitcoinKey)
		return btc.NewBitcoinAtom(network, config.Bitcoin, btcKey, orderID)
	case token.ETH:
		ethKey := key.GetKey(t).(keystore.EthereumKey)
		return eth.NewEthereumAtom(network, config.Ethereum, ethKey, orderID)
	}
	return nil, fmt.Errorf("Atom Build Failed")
}
