package watch

import (
	"github.com/republicprotocol/atom-go/domains/match"
	"github.com/republicprotocol/atom-go/services/store"
	"github.com/republicprotocol/atom-go/services/swap"
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
	Run([32]byte) error
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
func (watch *watch) Run(orderID [32]byte) error {
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

	if watch.str.ReadStatus(orderID) == "MATCHED" {
		if watch.reqAtom.PriorityCode() == match.ReceiveCurrency() {
			addr, err := watch.reqAtom.GetKey().GetAddress()
			if err != nil {
				return err
			}
			if err = watch.info.SetOwnerAddress(orderID, addr); err != nil {
				return err
			}
		} else {
			addr, err := watch.resAtom.GetKey().GetAddress()
			if err != nil {
				return err
			}
			if err = watch.info.SetOwnerAddress(orderID, addr); err != nil {
				return err
			}
		}
		err := watch.str.UpdateStatus(orderID, "INFO_SUBMITTED")
		if err != nil {
			return err
		}
	}

	atomicSwap := swap.NewSwap(watch.reqAtom, watch.resAtom, watch.info, match, watch.network, watch.str)
	err := atomicSwap.Execute()
	return err
}

func (watch *watch) Status(orderID [32]byte) string {
	return watch.str.ReadStatus(orderID)
}
