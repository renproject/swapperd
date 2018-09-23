package renex

import (
	"github.com/republicprotocol/renex-swapper-go/domain/swap"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/service/logger"
	"github.com/republicprotocol/renex-swapper-go/service/state"
	swapService "github.com/republicprotocol/renex-swapper-go/service/swap"
)

type Adapter interface {
	state.State
	logger.Logger
	swapService.Swapper
	Network
	GetOrderMatch(orderID [32]byte, waitTill int64) (swap.Match, error)
	GetAddresses(token.Token, token.Token) (string, string)
}

type Network interface {
	SendSwapDetails([32]byte, SwapDetails) error
	ReceiveSwapDetails([32]byte, int64) (SwapDetails, error)
}

type SwapDetails struct {
	SecretHash         [32]byte `json:"secretHash"`
	TimeLock           int64    `json:"timeLock"`
	SendToAddress      string   `json:"sendToAddress"`
	ReceiveFromAddress string   `json:"receiveFromAddress"`
}
