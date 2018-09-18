package guardian

import (
	"github.com/republicprotocol/renex-swapper-go/service/logger"
	"github.com/republicprotocol/renex-swapper-go/service/state"
)

type Adapter interface {
	logger.Logger
	state.State
	Refund([32]byte) error
}
