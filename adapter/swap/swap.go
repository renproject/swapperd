package swap

import (
	"github.com/republicprotocol/renex-swapper-go/adapter/atoms"
	"github.com/republicprotocol/renex-swapper-go/adapter/logger"
	"github.com/republicprotocol/renex-swapper-go/adapter/watchdog"
	"github.com/republicprotocol/renex-swapper-go/domain/match"
	"github.com/republicprotocol/renex-swapper-go/driver/network"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type builder struct {
	network.Network
	watchdog.Watchdog
	State
	logger.Logger
	atoms.Builder
}

func New(network swap.Network, watchdog swap.Watchdog, state State, logger swap.Logger, atomBuilder swap.AtomBuilder) swap.Builder {
	return &builder{
		Network:     network,
		Watchdog:    watchdog,
		State:       state,
		Logger:      logger,
		AtomBuilder: atomBuilder,
	}
}

type adapter struct {
	swap.Network
	swap.Watchdog
	swap.State
	swap.Logger
	personalAtom swap.Atom
	foreignAtom  swap.Atom
	match        match.Match
}

func (builder *builder) New(match match.Match) (swap.Adapter, error) {
	personal, foreign, err := builder.BuildAtoms(builder.State, req)
	if err != nil {
		return nil, err
	}
	return &adapter{
		Network:      builder.Network,
		Watchdog:     builder.Watchdog,
		State:        builder.State,
		Logger:       builder.Logger,
		personalAtom: personal,
		foreignAtom:  foreign,
		match:        match,
	}, nil
}

func (adapter *adapter) PersonalAtom() swap.Atom {
	return adapter.personalAtom
}

func (adapter *adapter) ForeignAtom() swap.Atom {
	return adapter.foreignAtom
}

func (adapter *adapter) Request() swap.Request {
	return adapter.request
}
