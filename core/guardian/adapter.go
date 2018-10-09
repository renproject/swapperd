package guardian

import (
	"github.com/republicprotocol/swapperd/service/logger"
	"github.com/republicprotocol/swapperd/service/state"
)

type Adapter interface {
	logger.Logger
	state.State
	Refund([32]byte) error
}
