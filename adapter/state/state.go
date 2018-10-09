package state

import (
	"github.com/republicprotocol/swapperd/service/logger"
	"github.com/republicprotocol/swapperd/service/state"
)

type stateAdapter struct {
	state.Store
	logger.Logger
}

func New(store state.Store, logger logger.Logger) state.Adapter {
	return &stateAdapter{
		store,
		logger,
	}
}
