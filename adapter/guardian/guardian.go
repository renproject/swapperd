package guardian

import (
	"fmt"

	"github.com/republicprotocol/renex-swapper-go/adapter/btc"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/eth"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	swapDomain "github.com/republicprotocol/renex-swapper-go/domain/swap"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/service/guardian"
	"github.com/republicprotocol/renex-swapper-go/service/logger"
	"github.com/republicprotocol/renex-swapper-go/service/state"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type guardianAdapter struct {
	config.Config
	keystore.Keystore
	logger.Logger
	state.State
}

func New(conf config.Config, ks keystore.Keystore, state state.State, logger logger.Logger) guardian.Adapter {
	return &guardianAdapter{
		Config:   conf,
		Keystore: ks,
		State:    state,
		Logger:   logger,
	}
}

// TODO: Check whether the atom is initiated before building the atom.
func (adapter *guardianAdapter) Refund(orderID [32]byte) error {
	req, err := adapter.SwapRequest(orderID)
	if err != nil {
		return err
	}

	personalAtom, err := buildAtom(adapter.Keystore, adapter.Config, adapter.Logger, req.SendToken, req)
	if err != nil {
		return err
	}

	return personalAtom.Refund()
}

func buildAtom(key keystore.Keystore, config config.Config, logger logger.Logger, t token.Token, req swapDomain.Request) (swap.Atom, error) {
	switch t {
	case token.BTC:
		btcKey := key.GetKey(t).(keystore.BitcoinKey)
		return btc.NewBitcoinAtom(config.Bitcoin, btcKey, logger, req)
	case token.ETH:
		ethKey := key.GetKey(t).(keystore.EthereumKey)
		return eth.NewEthereumAtom(config.Ethereum, ethKey, logger, req)
	}
	return nil, fmt.Errorf("Atom Build Failed")
}
