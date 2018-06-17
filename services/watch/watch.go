package watch

import "github.com/republicprotocol/atom-go/services/swap"

type Watch struct {
	network swap.Network
	axc     swap.Contract
	wallet  Wallet
}

// Run runs the watch object on the given order id
func (watch *Watch) Run(orderID [32]byte) error {
	orderID2, err := watch.wallet.WaitForMatch(orderID)
	if err != nil {
		return err
	}

	_, err = watch.wallet.GetMatch(orderID, orderID2)
	if err != nil {
		return err
	}

	// atomicSwap := swap.NewSwap(personalAtom, foreignAtom, watch.axc, match, watch.network)

	// err = atomicSwap.Execute()
	// if err != nil {
	// 	return err
	// }

	return nil
}
