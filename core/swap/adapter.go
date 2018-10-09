package swap

import (
	"github.com/republicprotocol/swapperd/core/logger"
	"github.com/republicprotocol/swapperd/foundation"
)

type SwapperAdapter interface {
	NewSwap(foundation.Swap) (Atom, Atom, Adapter, error)
}

type Adapter interface {
	logger.Logger
	Complain([32]byte) error
}

type Atom interface {
	Initiate() error
	Refund() error
	AuditSecret() (secret [32]byte, err error)
	Redeem(secret [32]byte) error
	Audit() error
}
