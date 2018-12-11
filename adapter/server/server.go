package server

import (
	"github.com/republicprotocol/swapperd/core/balance"
	"github.com/republicprotocol/swapperd/foundation/swap"
)

type Server interface {
	Run(doneCh <-chan struct{}, swaps chan<- swap.SwapBlob, receipts chan<- swap.SwapReceipt, statusQueries chan<- swap.ReceiptQuery, balanceQueries chan<- balance.BalanceQuery)
}
