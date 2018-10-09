package swapper

import "github.com/republicprotocol/swapperd/foundation"

type SwapContractBinder interface {
	Initiate()
	Audit() error
	Redeem([32]byte)
	AuditSecret() ([32]byte, error)
	Refund()
}

func Swap(native, foreign SwapContractBinder, req foundation.Swap, done chan<- foundation.SwapID) {
	if req.IsFirst {
		native.Initiate()
		if err := foreign.Audit(); err != nil {
			native.Refund()
			done <- req.ID
		}
		foreign.Redeem(req.Secret)
		done <- req.ID
		return
	}
	if err := foreign.Audit(); err != nil {
		done <- req.ID
	}
	native.Initiate()
	secret, err := native.AuditSecret()
	if err != nil {
		native.Refund()
		done <- req.ID
	}
	foreign.Redeem(secret)
	done <- req.ID
}
