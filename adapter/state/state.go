package state

import (
	"github.com/republicprotocol/renex-swapper-go/service/logger"
	"github.com/republicprotocol/renex-swapper-go/service/state"
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
