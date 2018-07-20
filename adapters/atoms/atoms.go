package atoms

import (
	"fmt"

	"github.com/republicprotocol/atom-go/domains/match"

	"github.com/republicprotocol/atom-go/services/store"
	"github.com/republicprotocol/atom-go/services/swap"

	"github.com/republicprotocol/atom-go/adapters/atoms/btc"
	"github.com/republicprotocol/atom-go/adapters/atoms/eth"
	"github.com/republicprotocol/atom-go/adapters/config"

	btcClient "github.com/republicprotocol/atom-go/adapters/clients/btc"
	ethClient "github.com/republicprotocol/atom-go/adapters/clients/eth"
)

type atomBuilder struct {
	keys   []swap.Key
	config config.Config
}

func NewAtomBuilder(config config.Config, keys []swap.Key) swap.AtomBuilder {
	return &atomBuilder{
		keys:   keys,
		config: config,
	}
}

func (ab *atomBuilder) BuildAtoms(state store.SwapState, m match.Match) (swap.Atom, swap.Atom, error) {

	if len(ab.keys) != 2 {
		return nil, nil, fmt.Errorf("This software does not support more than two keys at the moment")
	}

	var personalAtom, foreignAtom swap.Atom
	var err error

	if ab.keys[0].PriorityCode() == m.SendCurrency() {
		personalAtom, err = buildAtom(ab.keys[0], ab.config, m.PersonalOrderID())
		if err != nil {
			return nil, nil, err
		}

		foreignAtom, err = buildAtom(ab.keys[1], ab.config, m.ForeignOrderID())
		if err != nil {
			return nil, nil, err
		}
	} else {
		personalAtom, err = buildAtom(ab.keys[1], ab.config, m.PersonalOrderID())
		if err != nil {
			return nil, nil, err
		}

		foreignAtom, err = buildAtom(ab.keys[0], ab.config, m.ForeignOrderID())
		if err != nil {
			return nil, nil, err
		}
	}

	if state.AtomExists(m.PersonalOrderID()) {
		if err := personalAtom.Restore(state); err != nil {
			return nil, nil, err
		}
	}

	if state.AtomExists(m.ForeignOrderID()) {
		if err := foreignAtom.Restore(state); err != nil {
			return nil, nil, err
		}
	}

	return personalAtom, foreignAtom, nil
}

func buildAtom(key swap.Key, config config.Config, orderID [32]byte) (swap.Atom, error) {
	switch key.PriorityCode() {
	case 0:
		conn, err := btcClient.Connect(config)
		if err != nil {
			return nil, err
		}
		return btc.NewBitcoinAtom(conn, key, orderID), nil
	case 1:
		conn, err := ethClient.Connect(config)
		if err != nil {
			return nil, err
		}
		return eth.NewEthereumAtom(conn, key, orderID)
	}
	return nil, fmt.Errorf("Atom Build Failed")
}
