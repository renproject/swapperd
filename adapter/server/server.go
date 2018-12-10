package server

import (
	"github.com/republicprotocol/swapperd/core/balance"
	"github.com/republicprotocol/swapperd/foundation/swap"
)

type Server interface {
	Run(done <-chan struct{}, swapRequests chan<- swap.SwapRequest, statusQueries chan<- swap.ReceiptQuery, balanceQueries chan<- balance.BalanceQuery)
}
