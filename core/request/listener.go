package request

import "github.com/republicprotocol/swapperd/foundation"

type Listener interface {
	Run(done <-chan struct{}, swapRequests chan<- foundation.SwapRequest, statusQueries chan<- foundation.StatusQuery)
}
