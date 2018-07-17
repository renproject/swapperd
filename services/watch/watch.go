package watch

import (
	"log"

	"github.com/republicprotocol/atom-go/domains/match"
	"github.com/republicprotocol/atom-go/services/store"
	"github.com/republicprotocol/atom-go/services/swap"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/co-go"
)

type watch struct {
	network swap.Network
	info    swap.Info
	reqAtom swap.Atom
	resAtom swap.Atom
	wallet  Wallet
	str     store.State
}

type Watch interface {
	Run(<-chan struct{}, <-chan struct{}) <-chan error
	Add([32]byte) error
	Status([32]byte) string
}

func NewWatch(network swap.Network, info swap.Info, wallet Wallet, reqAtom swap.Atom, resAtom swap.Atom, str store.State) Watch {
	return &watch{
		network: network,
		info:    info,
		wallet:  wallet,
		reqAtom: reqAtom,
		resAtom: resAtom,
		str:     str,
	}
}

// Run runs the watch object on the given order id
func (watch *watch) Run(done <-chan struct{}, notification <-chan struct{}) <-chan error {
	errs := make(chan error)
	go func() {
		// TODO: Remove spin lock.
		for {
			select {
			case <-done:
				return
			case <-notification:
				swaps := watch.str.GetSwaps()
				co.ParForAll(swaps, func(i int) {
					if err := watch.doSwap(swaps[i]); err != nil {
						errs <- err
						return
					}
					watch.str.DeleteSwap(swaps[i])
				})
			}
		}
	}()
	return errs
}

func (watch *watch) Add(orderID [32]byte) error {
	return watch.str.AddSwap(orderID)
}

func (watch *watch) Status(orderID [32]byte) string {
	return watch.str.ReadStatus(orderID)
}

func (watch *watch) doSwap(orderID [32]byte) error {
	// TODO: All statuses should be defined as enumerated constants.
	if watch.str.ReadStatus(orderID) == "UNKNOWN" {
		err := watch.str.UpdateStatus(orderID, "PENDING")
		if err != nil {
			return err
		}
	}

	var match match.Match
	if watch.str.ReadStatus(orderID) == "PENDING" {
		var err error
		match, err = watch.wallet.GetMatch(orderID)
		if err != nil {
			return err
		}

		err = watch.str.SetMatch(orderID, match)
		if err != nil {
			return err
		}

		err = watch.str.UpdateStatus(orderID, "MATCHED")
		if err != nil {
			return err
		}
	} else {
		var err error
		match, err = watch.str.GetMatch(orderID)
		if err != nil {
			return err
		}
	}

	log.Println("Match found for ", order.ID(orderID))

	if watch.str.ReadStatus(orderID) == "MATCHED" {
		if watch.reqAtom.PriorityCode() == match.ReceiveCurrency() {
			addr, err := watch.reqAtom.GetKey().GetAddress()
			if err != nil {
				return err
			}
			if err := watch.info.SetOwnerAddress(orderID, addr); err != nil {
				return err
			}
		} else {
			addr, err := watch.resAtom.GetKey().GetAddress()
			if err != nil {
				return err
			}
			if err := watch.info.SetOwnerAddress(orderID, addr); err != nil {
				return err
			}
		}
		if err := watch.str.UpdateStatus(orderID, "INFO_SUBMITTED"); err != nil {
			return err
		}
	}

	if watch.str.ReadStatus(orderID) != "REDEEMED" && watch.str.ReadStatus(orderID) != "REFUNDED" {
		atomicSwap := swap.NewSwap(watch.reqAtom, watch.resAtom, watch.info, match, watch.network, watch.str)
		err := atomicSwap.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}
