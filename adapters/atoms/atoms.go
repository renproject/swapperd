package atoms

import (
	"fmt"

	"github.com/republicprotocol/atom-go/domains/match"

	"github.com/republicprotocol/atom-go/services/store"
	"github.com/republicprotocol/atom-go/services/swap"

	"github.com/republicprotocol/atom-go/adapters/atoms/btc"
	"github.com/republicprotocol/atom-go/adapters/atoms/eth"
	"github.com/republicprotocol/atom-go/adapters/configs/keystore"
	"github.com/republicprotocol/atom-go/adapters/configs/network"

	"github.com/republicprotocol/atom-go/adapters/blockchain/binder"
	btcClient "github.com/republicprotocol/atom-go/adapters/blockchain/clients/btc"
	ethClient "github.com/republicprotocol/atom-go/adapters/blockchain/clients/eth"
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
	b, err := binder.NewBinder(keystore.EthereumKey.GetKey(), ethConn)
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
		return btc.NewBitcoinAtom(&binder, conn, &key.BitcoinKey, orderID), nil
	case 1:
		conn, err := ethClient.Connect(config)
		if err != nil {
			return nil, err
		}
		return eth.NewEthereumAtom(&binder, conn, &key.EthereumKey, orderID)
	}
	return nil, fmt.Errorf("Atom Build Failed")
}
