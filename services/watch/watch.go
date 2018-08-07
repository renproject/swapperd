package watch

import (
	"fmt"
	"log"

	co "github.com/republicprotocol/co-go"
	"github.com/republicprotocol/renex-swapper-go/services/store"
	"github.com/republicprotocol/renex-swapper-go/services/swap"
)

type watch struct {
	adapter  Adapter
	state    store.State
	notifyCh chan struct{}
	doneCh   chan struct{}
}

type Watch interface {
	Start() <-chan error
	Add([32]byte) error
	Status([32]byte) string
	Notify()
	Stop()
}

func NewWatch(adapter Adapter, state store.State) Watch {
	return &watch{
		adapter:  adapter,
		state:    state,
		notifyCh: make(chan struct{}, 1),
		doneCh:   make(chan struct{}, 1),
	}
}

// Run runs the watch object on the given order id
func (watch *watch) Start() <-chan error {
	errs := make(chan error)
	log.Println("Starting the watcher......")
	go func() {
		defer close(errs)
		defer log.Println("Stopping the watcher......")
		for {
			select {
			case <-watch.doneCh:
				return
			case <-watch.notifyCh:
				swaps, err := watch.state.PendingSwaps()
				if err != nil {
					errs <- err
					continue
				}
				co.ParForAll(swaps, func(i int) {
					if err := watch.Swap(swaps[i]); err != nil {
						errs <- err
						return
					}
					watch.state.DeleteSwap(swaps[i])
				})
			}
		}
	}()
	return errs
}

func (watch *watch) Add(orderID [32]byte) error {
	return watch.state.AddSwap(orderID)
}

func (watch *watch) Status(orderID [32]byte) string {
	return watch.state.Status(orderID)
}

func (watch *watch) Notify() {
	watch.notifyCh <- struct{}{}
}

func (watch *watch) Stop() {
	watch.doneCh <- struct{}{}
}

func (watch *watch) Swap(orderID [32]byte) error {
	if watch.state.Status(orderID) == "UNKNOWN" {
		if err := watch.initiate(orderID); err != nil {
			return err
		}
	} else {
		// log.Println("Skipping swap initiation for ", order.ID(orderID))
	}

	if watch.state.Status(orderID) == "PENDING" {
		if err := watch.getMatch(orderID); err != nil {
			return err
		}
	} else {
		// log.Println("Skipping get match for ", order.ID(orderID))
	}

	if watch.state.Status(orderID) == "MATCHED" {
		if err := watch.setInfo(orderID); err != nil {
			return err
		}
	} else {
		// log.Println("Skipping Info Submission for ", order.ID(orderID))
	}

	if watch.state.Status(orderID) != "REDEEMED" && watch.state.Status(orderID) != "REFUNDED" {
		if err := watch.execute(orderID); err != nil {
			return err
		}
	} else {
		// log.Println("Skipping Execute for ", order.ID(orderID))
	}

	return nil
}

func (w *watch) setInfo(orderID [32]byte) error {
	// log.Println("Submitting info for", order.ID(orderID))

	m, err := w.state.Match(orderID)
	if err != nil {
		return err
	}

	_, foreignAtom, err := w.adapter.BuildAtoms(w.state, m)
	if err != nil {
		return err
	}

	addr, err := foreignAtom.GetFromAddress()
	if err != nil {
		return err
	}

	if err := w.adapter.SendOwnerAddress(orderID, addr); err != nil {
		fmt.Println(err)
		return err
	}

	if err := w.state.PutStatus(orderID, swap.StatusInfoSubmitted); err != nil {
		return err
	}

	// log.Println("Info Submitted for ", order.ID(orderID))
	return nil
}

func (w *watch) execute(orderID [32]byte) error {
	m, err := w.state.Match(orderID)
	if err != nil {
		return err
	}

	personalAtom, foreignAtom, err := w.adapter.BuildAtoms(w.state, m)
	if err != nil {
		return err
	}

	atomicSwap := swap.NewSwap(personalAtom, foreignAtom, m, w.adapter, w.state)
	if err := atomicSwap.Execute(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (w *watch) initiate(orderID [32]byte) error {
	// log.Println("Starting the Atomic Swap for", order.ID(orderID))
	err := w.state.PutStatus(orderID, "PENDING")
	if err != nil {
		return err
	}
	// log.Println("Started the Atomic Swap for", order.ID(orderID))
	return nil
}

func (w *watch) getMatch(orderID [32]byte) error {
	// log.Println("Waiting for the match to be found for ", order.ID(orderID))
	match, err := w.adapter.CheckForMatch(orderID, true)
	if err != nil {
		return err
	}

	err = w.state.PutMatch(orderID, match)
	if err != nil {
		return err
	}

	err = w.state.PutStatus(orderID, "MATCHED")
	if err != nil {
		return err
	}

	// log.Println("Match found :", order.ID(orderID), " <---->", order.ID(match.ForeignOrderID()))
	return nil
}
