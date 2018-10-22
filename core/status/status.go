package status

import (
	"github.com/republicprotocol/swapperd/foundation"
)

type Query struct {
	Responder chan<- map[foundation.SwapID]foundation.SwapStatus
}

type Book interface {
	Run(done <-chan struct{}, statuses <-chan foundation.SwapStatus, queries <-chan Query)
}

type book struct {
	monitor *monitor
}

func New() Book {
	return &book{newMonitor()}
}

func (book *book) Run(done <-chan struct{}, statuses <-chan foundation.SwapStatus, queries <-chan Query) {
	for {
		select {
		case <-done:
			return
		case status, ok := <-statuses:
			if !ok {
				return
			}
			book.monitor.set(status.ID, status)
		case query, ok := <-queries:
			if !ok {
				return
			}
			query.Responder <- book.monitor.get()
		}
	}
}
