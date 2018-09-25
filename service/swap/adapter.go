package swap

import (
	"github.com/republicprotocol/renex-swapper-go/domain/swap"
	"github.com/republicprotocol/renex-swapper-go/service/logger"
)

type SwapperAdapter interface {
	NewSwap(swap.Request) (Atom, Atom, Adapter, error)
}

type Adapter interface {
	logger.Logger
	Watchdog
}

type Atom interface {
	Initiate() error
	Refund() error
	AuditSecret() (secret [32]byte, err error)
	Redeem(secret [32]byte) error
	Audit() error
}

// TODO: Remove watchdog for Complain([32]byte) error
type Watchdog interface {
	ComplainDelayedRequestorInitiation([32]byte) error
	ComplainWrongRequestorInitiation([32]byte) error
	ComplainDelayedResponderInitiation([32]byte) error
	ComplainWrongResponderInitiation([32]byte) error
	ComplainDelayedRequestorRedemption([32]byte) error
}
