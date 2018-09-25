package swap

import (
	"fmt"

	"github.com/republicprotocol/renex-swapper-go/adapter/btc"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/eth"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	swapDomain "github.com/republicprotocol/renex-swapper-go/domain/swap"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/service/logger"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type swapperAdapter struct {
	logger.Logger
	config   config.Config
	keystore keystore.Keystore
}

func New(cfg config.Config, ks keystore.Keystore, logger logger.Logger) swap.SwapperAdapter {
	return &swapperAdapter{
		config:   cfg,
		keystore: ks,
		Logger:   logger,
	}
}

func (swapper *swapperAdapter) NewSwap(req swapDomain.Request) (swap.Atom, swap.Atom, swap.Adapter, error) {
	personalAtom, err := buildAtom(swapper.keystore, swapper.config, req.SendToken, req)
	if err != nil {
		return nil, nil, nil, err
	}
	foreignAtom, err := buildAtom(swapper.keystore, swapper.config, req.ReceiveToken, req)
	if err != nil {
		return nil, nil, nil, err
	}
	return personalAtom, foreignAtom, swapper, nil
}

func (swapper *swapperAdapter) Complain(UID [32]byte) error {
	return nil
}

func buildAtom(key keystore.Keystore, config config.Config, t token.Token, req swapDomain.Request) (swap.Atom, error) {
	switch t {
	case token.BTC:
		btcKey := key.GetKey(t).(keystore.BitcoinKey)
		return btc.NewBitcoinAtom(config.Bitcoin, btcKey, req)
	case token.ETH:
		ethKey := key.GetKey(t).(keystore.EthereumKey)
		return eth.NewEthereumAtom(config.Ethereum, ethKey, req)
	}
	return nil, fmt.Errorf("Atom Build Failed")
}
