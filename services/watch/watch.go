package watch

import (
	"github.com/republicprotocol/atom-go/services/swap"
)

type Watch struct {
	network swap.Network
	info    swap.Info
	reqAtom swap.AtomRequester
	resAtom swap.AtomResponder
	wallet  Wallet
}

func NewWatch(network swap.Network, info swap.Info, wallet Wallet, reqAtom swap.AtomRequester, resAtom swap.AtomResponder) Watch {
	return Watch{
		network: network,
		info:    info,
		wallet:  wallet,
		reqAtom: reqAtom,
		resAtom: resAtom,
	}
}

// Run runs the watch object on the given order id
func (watch *Watch) Run(orderID [32]byte) error {
	match, err := watch.wallet.GetMatch(orderID)
	if err != nil {
		return err
	}
	// watch.info.SetOwnerAddress()
	atomicSwap := swap.NewSwap(watch.reqAtom, watch.resAtom, watch.info, match, watch.network)
	err = atomicSwap.Execute()
	return err
}
