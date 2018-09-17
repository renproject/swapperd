package atoms

import (
	"fmt"

	"github.com/republicprotocol/renex-swapper-go/domain/match"
	"github.com/republicprotocol/renex-swapper-go/domain/token"

	"github.com/republicprotocol/renex-swapper-go/service/store"
	"github.com/republicprotocol/renex-swapper-go/service/swap"

	"github.com/republicprotocol/renex-swapper-go/adapter/atoms/btc"
	"github.com/republicprotocol/renex-swapper-go/adapter/atoms/eth"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"

	btcClient "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/clients/btc"
	ethClient "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/clients/eth"
)

type Builder interface {
	BuildAtoms(state store.State, m match.Match) (swap.Atom, swap.Atom, error)
}
type atomBuilder struct {
	binder   eth.Adapter
	keystore keystore.Keystore
	config   config.Config
}

func NewAtomBuilder(network eth.Adapter, config config.Config, keystore keystore.Keystore) (Builder, error) {
	return &atomBuilder{
		binder:   network,
		keystore: keystore,
		config:   config,
	}, nil
}

func (ab *atomBuilder) BuildAtoms(state store.State, m match.Match) (swap.Atom, swap.Atom, error) {
	var personalAtom, foreignAtom swap.Atom
	var err error

	personalAtom, err = buildAtom(ab.binder, ab.keystore, ab.config, m.SendCurrency(), m.PersonalOrderID())
	if err != nil {
		return nil, nil, err
	}

	foreignAtom, err = buildAtom(ab.binder, ab.keystore, ab.config, m.ReceiveCurrency(), m.ForeignOrderID())
	if err != nil {
		return nil, nil, err
	}

	if state.AtomExists(m.PersonalOrderID()) {
		details, err := state.AtomDetails(m.PersonalOrderID())
		if err != nil {
			return nil, nil, err
		}
		if err := personalAtom.Deserialize(details); err != nil {
			return nil, nil, err
		}
	}

	if state.AtomExists(m.ForeignOrderID()) {
		details, err := state.AtomDetails(m.ForeignOrderID())
		if err != nil {
			return nil, nil, err
		}
		if err := foreignAtom.Deserialize(details); err != nil {
			return nil, nil, err
		}
	}

	return personalAtom, foreignAtom, nil
}

func buildAtom(network eth.Adapter, key keystore.Keystore, config config.Config, t uint32, orderID [32]byte) (swap.Atom, error) {
	switch t {
	case 0:
		conn, err := btcClient.NewConnWithConfig(config)
		if err != nil {
			return nil, err
		}
		btcKey := key.GetKey(token.BTC).(keystore.BitcoinKey)
		return btc.NewBitcoinAtom(network, conn, btcKey, orderID), nil
	case 1:
		conn, err := ethClient.Connect(config)
		if err != nil {
			return nil, err
		}
		ethKey := key.GetKey(token.ETH).(keystore.EthereumKey)
		return eth.NewEthereumAtom(network, conn, ethKey, orderID)
	}
	return nil, fmt.Errorf("Atom Build Failed")
}
