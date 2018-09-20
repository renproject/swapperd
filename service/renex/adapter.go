package renex

import (
	"github.com/republicprotocol/renex-swapper-go/domain/match"
	"github.com/republicprotocol/renex-swapper-go/domain/order"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/service/logger"
	"github.com/republicprotocol/renex-swapper-go/service/state"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type Adapter interface {
	state.State
	logger.Logger
	swap.Swapper
	swap.Network
	GetOrderMatch(orderID order.ID, waitTill int64) (match.Match, error)
	GetAddress(token.Token) []byte
}
