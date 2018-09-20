package swap

import (
	"math/big"

	"github.com/republicprotocol/renex-swapper-go/domain/match"

	"github.com/republicprotocol/renex-swapper-go/domain/order"
	"github.com/republicprotocol/renex-swapper-go/service/logger"
	"github.com/republicprotocol/renex-swapper-go/service/state"
)

type SwapperAdapter interface {
	NewSwap(order.ID) (Atom, Atom, match.Match, Adapter, error)
}

type Adapter interface {
	logger.Logger
	state.State
	Network
	Watchdog
}

type Atom interface {
	Initiate(to []byte, hash [32]byte, value *big.Int, expiry int64) error
	Refund() error
	AuditSecret() (secret [32]byte, err error)
	Redeem(secret [32]byte) error
	Audit() ([32]byte, []byte, *big.Int, int64, error)
	WaitForCounterRedemption() error
	Serialize() ([]byte, error)
	Deserialize([]byte) error
	GetFromAddress() ([]byte, error)
	PriorityCode() uint32
	RedeemedAt() (int64, error)
}

type Network interface {
	SendOwnerAddress(order.ID, []byte) error
	SendSwapDetails(order.ID, []byte) error
	ReceiveOwnerAddress(order.ID, int64) ([]byte, error)
	ReceiveSwapDetails(order.ID, int64) ([]byte, error)
}

type Watchdog interface {
	ComplainDelayedAddressSubmission([32]byte) error
	ComplainDelayedRequestorInitiation([32]byte) error
	ComplainWrongRequestorInitiation([32]byte) error
	ComplainDelayedResponderInitiation([32]byte) error
	ComplainWrongResponderInitiation([32]byte) error
	ComplainDelayedRequestorRedemption([32]byte) error
}
