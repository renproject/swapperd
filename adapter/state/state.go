package state

import (
	"github.com/republicprotocol/renex-swapper-go/adapter/store"
	"github.com/republicprotocol/renex-swapper-go/service/logger"
	"github.com/republicprotocol/renex-swapper-go/service/state"
)

type stateAdapter struct {
	store.Store
	logger.Logger
}

func New(store store.Store, logger logger.Logger) state.Adapter {
	return &stateAdapter{
		store,
		logger,
	}
}
