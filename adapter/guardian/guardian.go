package guardian

import (
	"fmt"

	"github.com/republicprotocol/renex-swapper-go/adapter/atoms/btc"
	"github.com/republicprotocol/renex-swapper-go/adapter/atoms/eth"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/adapter/network"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/service/guardian"
	"github.com/republicprotocol/renex-swapper-go/service/logger"
	"github.com/republicprotocol/renex-swapper-go/service/state"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type guardianAdapter struct {
	config.Config
	keystore.Keystore
	network.Network
	logger.Logger
	state.State
}

func New(conf config.Config, ks keystore.Keystore, net network.Network, state state.State, logger logger.Logger) guardian.Adapter {
	return &guardianAdapter{
		Config:   conf,
		Keystore: ks,
		Network:  net,
		State:    state,
		Logger:   logger,
	}
}

func (guardian *guardianAdapter) Refund(orderID [32]byte) error {
	match, err := guardian.Match(orderID)
	if err != nil {
		return err
	}

	personalAtom, err := buildAtom(guardian.Network, guardian.Keystore, guardian.Config, match.SendCurrency(), match.PersonalOrderID())
	if err != nil {
		return err
	}

	details, err := guardian.AtomDetails(match.PersonalOrderID())
	if err != nil {
		return err
	}

	if err := personalAtom.Deserialize(details); err != nil {
		return err
	}

	return personalAtom.Refund()
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
