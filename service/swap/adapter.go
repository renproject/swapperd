package swap

import (
	"math/big"

	"github.com/republicprotocol/renex-swapper-go/domain/match"
	"github.com/republicprotocol/renex-swapper-go/domain/order"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/service/logger"
	"github.com/republicprotocol/renex-swapper-go/service/state"
)

type Request struct {
	UID                [32]byte    `json:"uid"`
	TimeLock           int64       `json:"timeLock"`
	Secret             [32]byte    `json:"secret"`
	SecretHash         [32]byte    `json:"secretHash"`
	SendToAddress      string      `json:"sendToAddress"`
	ReceiveFromAddress string      `json:"receiveFromAddress"`
	SendValue          *big.Int    `json:"sendValue"`
	ReceiveValue       *big.Int    `json:"sendValue"`
	SendToken          token.Token `json:"sendToken"`
	ReceiveToken       token.Token `json:"receiveToken"`
	GoesFirst          bool        `json:"goesFirst"`
}

type SwapperAdapter interface {
	NewSwap(order.ID, Request) (Atom, Atom, match.Match, Adapter, error)
}

type Adapter interface {
	logger.Logger
	state.State
	Network
	Watchdog
}

type Atom interface {
	Initiate() error
	Refund() error
	AuditSecret() (secret [32]byte, err error)
	Redeem(secret [32]byte) error
	Audit() error
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
