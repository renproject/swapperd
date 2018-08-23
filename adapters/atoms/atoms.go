package atoms

import (
	"fmt"

	"github.com/republicprotocol/renex-swapper-go/domains/match"

	"github.com/republicprotocol/renex-swapper-go/services/store"
	"github.com/republicprotocol/renex-swapper-go/services/swap"

	"github.com/republicprotocol/renex-swapper-go/adapters/atoms/btc"
	"github.com/republicprotocol/renex-swapper-go/adapters/atoms/eth"
	"github.com/republicprotocol/renex-swapper-go/adapters/configs/keystore"
	"github.com/republicprotocol/renex-swapper-go/adapters/configs/network"

	"github.com/republicprotocol/renex-swapper-go/adapters/blockchain/binder"
	btcClient "github.com/republicprotocol/renex-swapper-go/adapters/blockchain/clients/btc"
	ethClient "github.com/republicprotocol/renex-swapper-go/adapters/blockchain/clients/eth"
)

type atomBuilder struct {
	binder   binder.Binder
	keystore keystore.Keystore
	config   network.Config
}

type AtomBuilder interface {
	BuildAtoms(state store.State, m match.Match) (swap.Atom, swap.Atom, error)
}

func NewAtomBuilder(config network.Config, keystore keystore.Keystore) (AtomBuilder, error) {
	ethConn, err := ethClient.Connect(config)
	if err != nil {
		return nil, err
	}
	ethKey, err := keystore.GetKey(1, 0)
	if err != nil {
		return nil, err
	}
	privKey, err := ethKey.GetKey()
	if err != nil {
		return nil, err
	}
	b, err := binder.NewBinder(privKey, ethConn)
	if err != nil {
		return nil, err
	}
	return &atomBuilder{
		binder:   b,
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

func buildAtom(binder binder.Binder, key keystore.Keystore, config network.Config, cc uint32, orderID [32]byte) (swap.Atom, error) {
	switch cc {
	case 0:
		conn, err := btcClient.Connect(config)
		if err != nil {
			return nil, err
		}
		btcKey, err := key.GetKey(0, 0)
		if err != nil {
			return nil, err
		}
		return btc.NewBitcoinAtom(&binder, conn, btcKey, orderID), nil
	case 1:
		conn, err := ethClient.Connect(config)
		if err != nil {
			return nil, err
		}
		ethKey, err := key.GetKey(1, 0)
		if err != nil {
			return nil, err
		}
		return eth.NewEthereumAtom(&binder, conn, ethKey, orderID)
	}
	return nil, fmt.Errorf("Atom Build Failed")
}
