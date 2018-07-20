package watch

import (
	"fmt"
	"log"
	"time"

	"github.com/republicprotocol/atom-go/domains/match"
	"github.com/republicprotocol/atom-go/services/store"
	"github.com/republicprotocol/atom-go/services/swap"
	"github.com/republicprotocol/co-go"
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
	Run(<-chan struct{}, <-chan struct{}) <-chan error
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

// Run runs the watch object on the given order id
func (watch *watch) Run(done <-chan struct{}, notification <-chan struct{}) <-chan error {
	errs := make(chan error)
	log.Println("Starting the watcher......")
	go func() {
		for {
			select {
			case <-done:
				return
			case <-notification:
				swaps, err := watch.state.PendingSwaps()
				fmt.Println("Getting Pending Swaps")
				if err != nil {
					fmt.Println(err.Error())
					errs <- err
					return
				}
				co.ParForAll(swaps, func(i int) {
					fmt.Println("Inside Par for all")
					if err := watch.Swap(swaps[i]); err != nil {
						errs <- err
						return
					}
					watch.state.DeleteSwap(swaps[i])
				})
			}
			time.Sleep(10 * time.Second)
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

func (watch *watch) Swap(orderID [32]byte) error {
	if watch.state.Status(orderID) == "UNKNOWN" {
		err := watch.state.PutStatus(orderID, "PENDING")
		if err != nil {
			return err
		}
		log.Println("Initiated the Atomic Swap")
	}

	var match match.Match
	if watch.state.Status(orderID) == "PENDING" {
		var err error
		match, err = watch.wallet.GetMatch(orderID)
		if err != nil {
			return err
		}

		err = watch.state.PutMatch(orderID, match)
		if err != nil {
			return err
		}

		err = watch.state.PutStatus(orderID, "MATCHED")
		if err != nil {
			return err
		}

		log.Println("Match found for ", order.ID(orderID))
	} else {
		var err error
		match, err = watch.state.Match(orderID)
		if err != nil {
			return err
		}
	}

	personalAtom, foreignAtom, err := watch.builder.BuildAtoms(watch.state, match)
	if err != nil {
		return err
	}

	if watch.state.Status(orderID) == "MATCHED" {
		addr, err := foreignAtom.GetKey().GetAddress()
		if err != nil {
			return err
		}
		log.Println("Setting owner address for ", order.ID(orderID))
		if err := watch.info.SetOwnerAddress(orderID, addr); err != nil {
			fmt.Println(err)
			return err
		}
		log.Println("...done", order.ID(orderID))

		log.Println("Put status for ", order.ID(orderID))
		if err := watch.state.PutStatus(orderID, "INFO_SUBMITTED"); err != nil {
			return err
		}
		log.Println("...done", order.ID(orderID))

		log.Println("Info Submitted for ", order.ID(orderID))
	} else {
		log.Println("Skipping Info Submission for ", order.ID(orderID))
	}

	fmt.Printf("Personal Code -> %d Foreign Code -> %d", personalAtom.PriorityCode(), foreignAtom.PriorityCode())
	if watch.state.Status(orderID) != "REDEEMED" && watch.state.Status(orderID) != "REFUNDED" {
		atomicSwap := swap.NewSwap(personalAtom, foreignAtom, watch.info, match, watch.network, watch.state)
		err := atomicSwap.Execute()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}
