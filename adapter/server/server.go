package server

import (
	"github.com/republicprotocol/swapperd/core/balance"
	"github.com/republicprotocol/swapperd/foundation"
)

type Server interface {
	Run(done <-chan struct{}, swapRequests chan<- foundation.SwapRequest, statusQueries chan<- foundation.StatusQuery, balanceQueries chan<- balance.BalanceQuery)
}
