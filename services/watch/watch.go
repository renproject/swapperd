package watch

import (
	"fmt"
	"log"

	"github.com/republicprotocol/atom-go/services/store"
	"github.com/republicprotocol/atom-go/services/swap"
	"github.com/republicprotocol/republic-go/order"
)

type watch struct {
	network swap.Network
	info    swap.Info
	builder swap.AtomBuilder
	wallet  Wallet
	state   store.SwapState
}

type Watch interface {
	// Run(<-chan struct{}, <-chan struct{}) <-chan error
	Add([32]byte) error
	Status([32]byte) string
	Swap([32]byte) error
}

func NewWatch(network swap.Network, info swap.Info, wallet Wallet, builder swap.AtomBuilder, state store.SwapState) Watch {
	return &watch{
		network: network,
		info:    info,
		builder: builder,
		wallet:  wallet,
		state:   state,
	}
}

// // Run runs the watch object on the given order id
// func (watch *watch) Run(done <-chan struct{}, notification <-chan struct{}) <-chan error {
// 	errs := make(chan error)
// 	log.Println("Starting the watcher......")
// 	go func() {
// 		for {
// 			select {
// 			case <-done:
// 				return
// 			case <-notification:
// 				swaps, err := watch.state.PendingSwaps()
// 				fmt.Println("Getting Pending Swaps")
// 				if err != nil {
// 					fmt.Println(err.Error())
// 					errs <- err
// 					return
// 				}
// 				co.ParForAll(swaps, func(i int) {
// 					fmt.Println("Inside Par for all")
// 					if err := watch.Swap(swaps[i]); err != nil {
// 						errs <- err
// 						return
// 					}
// 					watch.state.DeleteSwap(swaps[i])
// 				})
// 			}
// 			time.Sleep(10 * time.Second)
// 		}
// 	}()
// 	return errs
// }

func (watch *watch) Add(orderID [32]byte) error {
	return watch.state.AddSwap(orderID)
}

func (watch *watch) Status(orderID [32]byte) string {
	return watch.state.Status(orderID)
}

func (watch *watch) Swap(orderID [32]byte) error {
	if watch.state.Status(orderID) == "UNKNOWN" {
		if err := watch.initiate(orderID); err != nil {
			return err
		}
	} else {
		log.Println("Skipping swap initiation for ", order.ID(orderID))
	}

	if watch.state.Status(orderID) == "PENDING" {
		if err := watch.getMatch(orderID); err != nil {
			return err
		}
	} else {
		log.Println("Skipping get match for ", order.ID(orderID))
	}

	if watch.state.Status(orderID) == "MATCHED" {
		if err := watch.setInfo(orderID); err != nil {
			return err
		}
	} else {
		log.Println("Skipping Info Submission for ", order.ID(orderID))
	}

	if watch.state.Status(orderID) != "REDEEMED" && watch.state.Status(orderID) != "REFUNDED" {
		if err := watch.execute(orderID); err != nil {
			return err
		}
	} else {
		log.Println("Skipping Execute for ", order.ID(orderID))
	}

	return nil
}

func (w *watch) setInfo(orderID [32]byte) error {
	log.Println("Submitting info for", order.ID(orderID))

	m, err := w.state.Match(orderID)
	if err != nil {
		return err
	}

	_, foreignAtom, err := w.builder.BuildAtoms(w.state, m)
	if err != nil {
		return err
	}

	addr, err := foreignAtom.GetKey().GetAddress()
	if err != nil {
		return err
	}

	if err := w.info.SetOwnerAddress(orderID, addr); err != nil {
		fmt.Println(err)
		return err
	}

	if err := w.state.PutStatus(orderID, swap.StatusInfoSubmitted); err != nil {
		return err
	}

	log.Println("Info Submitted for ", order.ID(orderID))
	return nil
}

func (w *watch) execute(orderID [32]byte) error {
	m, err := w.state.Match(orderID)
	if err != nil {
		return err
	}

	personalAtom, foreignAtom, err := w.builder.BuildAtoms(w.state, m)
	if err != nil {
		return err
	}

	atomicSwap := swap.NewSwap(personalAtom, foreignAtom, w.info, m, w.network, w.state)
	if err := atomicSwap.Execute(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (w *watch) initiate(orderID [32]byte) error {
	log.Println("Initiating the Atomic Swap")
	err := w.state.PutStatus(orderID, "PENDING")
	if err != nil {
		return err
	}
	log.Println("Initiated the Atomic Swap")
	return nil
}

func (w *watch) getMatch(orderID [32]byte) error {
	log.Println("Waiting for the match to be found for ", order.ID(orderID))
	match, err := w.wallet.GetMatch(orderID)
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

	log.Println("Match found for ", order.ID(orderID))
	return nil
}
