package status

import "github.com/republicprotocol/swapperd/foundation/swap"

type Statuses interface {
	Run(done <-chan struct{}, swaps <-chan swap.SwapReceipt, updates <-chan swap.StatusUpdate, queries <-chan swap.ReceiptQuery)
}

type statuses struct {
	monitor *monitor
}

func New() Statuses {
	return &statuses{newMonitor()}
}

func (statuses *statuses) Run(done <-chan struct{}, receipts <-chan swap.SwapReceipt, updates <-chan swap.StatusUpdate, queries <-chan swap.ReceiptQuery) {
	for {
		select {
		case <-done:
			return
		case receipt, ok := <-receipts:
			if !ok {
				return
			}
			statuses.monitor.set(receipt)
		case update, ok := <-updates:
			if !ok {
				return
			}
			statuses.monitor.update(update)
		case query, ok := <-queries:
			if !ok {
				return
			}
			go func() {
				query.Responder <- statuses.monitor.get()
			}()
		}
	}
}
