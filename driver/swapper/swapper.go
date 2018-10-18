package swapper

import (
	"github.com/republicprotocol/co-go"
	"fmt"
	"log"
	netHttp "net/http"

	"github.com/republicprotocol/swapperd/adapter/server"
	"github.com/republicprotocol/swapperd/foundation"
)

type Swapper interface {
	Http(port int64)
}

type swapper struct {
}

func NewSwapper() Swapper {
	return &swapper{}
}

func (swapper *swapper) Run(port int64) {
	swapCh := make(chan foundation.Swap)
	co.ParBegin(
		func() {
			defer close(swapCh)
			swapper.Http(port, swapCh)
		}
	)
}

func (swapper *swapper) Http(port int64, swapCh chan<- foundation.Swap) {
	log.Fatal(netHttp.ListenAndServe(fmt.Sprintf(":%d", port), server.NewHandler(swapCh)))
}
