package status

import (
	"github.com/republicprotocol/swapperd/foundation"
)

type Statuses interface {
	Run(done <-chan struct{}, swaps <-chan foundation.SwapStatus, updates <-chan foundation.StatusUpdate, queries <-chan foundation.StatusQuery)
}

type statuses struct {
	monitor *monitor
}

func New() Statuses {
	return &statuses{newMonitor()}
}

func (statuses *statuses) Run(done <-chan struct{}, receipts <-chan foundation.SwapStatus, updates <-chan foundation.StatusUpdate, queries <-chan foundation.StatusQuery) {
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
