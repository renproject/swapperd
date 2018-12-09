package status

import (
	"github.com/republicprotocol/swapperd/foundation"
)

type Book interface {
	Run(done <-chan struct{}, swaps <-chan foundation.SwapStatus, updates <-chan foundation.StatusUpdate, queries <-chan foundation.StatusQuery)
}

type book struct {
	monitor *monitor
}

func New() Book {
	return &book{newMonitor()}
}

func (book *book) Run(done <-chan struct{}, statuses <-chan foundation.SwapStatus, updates <-chan foundation.StatusUpdate, queries <-chan foundation.StatusQuery) {
	for {
		select {
		case <-done:
			return
		case status, ok := <-statuses:
			if !ok {
				return
			}
			book.monitor.set(status)
		case update, ok := <-updates:
			if !ok {
				return
			}
			book.monitor.update(update)
		case query, ok := <-queries:
			if !ok {
				return
			}
			go func() {
				query.Responder <- book.monitor.get()
			}()
		}
	}
}
