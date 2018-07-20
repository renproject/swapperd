package guardian

import (
	"fmt"
	"log"
	"time"

	"github.com/republicprotocol/atom-go/services/errors"
	"github.com/republicprotocol/atom-go/services/store"
	"github.com/republicprotocol/atom-go/services/swap"
	co "github.com/republicprotocol/co-go"
)

type Guardian interface {
	Run(<-chan struct{}, <-chan struct{}) <-chan error
}

type guardian struct {
	builder swap.AtomBuilder
	state   store.SwapState
}

func (g *guardian) Run(done <-chan struct{}, notification <-chan struct{}) {
	errs := make(chan error)
	log.Println("Starting the guardian service......")
	go func() {
		defer log.Println("Ending the guardian service......")
		for {
			select {
			case <-done:
				return
			case <-notification:
				swaps, err := g.expiredSwaps()
				if err != nil {
					fmt.Println(err.Error())
					errs <- err
					return
				}
				co.ParForAll(swaps, func(i int) {
					if err := g.refund(swaps[i]); err != nil {
						errs <- err
						return
					}
					g.state.DeleteSwap(swaps[i])
				})
			}
			time.Sleep(10 * time.Second)
		}
	}()
}

func (g *guardian) refund(orderID [32]byte) error {
	atom, err := g.buildAtom(orderID)
	if err != nil {
		return errors.ErrAtomBuildFailed(err)
	}

	if err := atom.Refund(); err != nil {
		return errors.ErrRefundAfterRedeem(err)
	}
	return nil
}

func (g *guardian) expiredSwaps() ([][32]byte, error) {
	pendingSwaps, err := g.state.PendingSwaps()
	if err != nil {
		return nil, errors.ErrFailedPendingSwaps(err)
	}
	expiredSwaps := [][32]byte{}
	for _, swap := range pendingSwaps {
		expiry, _, err := g.state.InitiateDetails(swap)
		if err != nil {
			return nil, errors.ErrFailedInitiateDetails(err)
		}
		if expiry <= time.Now().Unix() {
			expiredSwaps = append(expiredSwaps, swap)
		}
	}
	return expiredSwaps, nil
}

func (g *guardian) buildAtom(orderID [32]byte) (swap.Atom, error) {
	m, err := g.state.Match(orderID)
	if err != nil {
		return nil, err
	}
	atom, _, err := g.builder.BuildAtoms(g.state, m)
	return atom, err
}
