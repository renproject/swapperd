package swap

import (
	"github.com/republicprotocol/renex-swapper-go/domains/order"
	"github.com/republicprotocol/renex-swapper-go/services/logger"
	"github.com/republicprotocol/renex-swapper-go/services/renguardClient"
)

type SwapAdapter interface {
	SendOwnerAddress(order.ID, []byte) error
	ReceiveOwnerAddress(order.ID, int64) ([]byte, error)
	ReceiveSwapDetails(order.ID, int64) ([]byte, error)
	SendSwapDetails(order.ID, []byte) error
	renguardClient.RenguardClient
	logger.Logger
}
