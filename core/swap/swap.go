package swap

import (
	"fmt"

	"github.com/republicprotocol/swapperd/domain/swap"
)

type Swap interface {
	Execute() error
}

type swapExec struct {
	req          swap.Request
	personalAtom Atom
	foreignAtom  Atom
	Adapter
}

func (swap *swapExec) Execute() error {
	swap.LogInfo(swap.req.UID, "starting the atomic swap ")
	if swap.req.GoesFirst {
		return swap.request(swap.req.UID, swap.req.Secret)
	}
	return swap.respond(swap.req.UID)
}

func (swap *swapExec) request(id [32]byte, secret [32]byte) error {
	swap.LogInfo(swap.req.UID, "Initiating the atomic swap ")
	if err := swap.personalAtom.Initiate(); err != nil {
		if err != ErrSwapAlreadyInitiated {
			swap.LogError(id, fmt.Sprintf("failed to initiate details: %v", err))
			return fmt.Errorf("failed to initiate details: %v", err)
		}
		swap.LogInfo(id, "swap already initiated")
	}
	swap.LogInfo(swap.req.UID, "Initiated the atomic swap ")

	swap.LogInfo(swap.req.UID, "Auditing the atomic swap ")
	if err := swap.foreignAtom.Audit(); err != nil {
		if err := swap.Complain(id); err != nil {
			swap.LogError(id, fmt.Sprintf("failed to complain to the watch dog: %v", err))
			return fmt.Errorf("failed to complain to the watch dog: %v", err)
		}
		return fmt.Errorf("incorrect swap complaint sent: %v", err)
	}
	swap.LogInfo(swap.req.UID, "Audited the atomic swap ")

	swap.LogInfo(swap.req.UID, "Redeeming the atomic swap ")
	if err := swap.foreignAtom.Redeem(secret); err != nil {
		if err != ErrSwapAlreadyRedeemedOrRefunded {
			swap.LogError(id, fmt.Sprintf("failed to redeem: %v", err))
			return fmt.Errorf("failed to redeem: %v", err)
		}
		swap.LogInfo(id, "swap already redeemed or refunded")
	}
	swap.LogInfo(swap.req.UID, "Redeemed the atomic swap ")

	return nil
}

func (swap *swapExec) respond(id [32]byte) error {
	swap.LogInfo(swap.req.UID, "Auditing the atomic swap ")
	if err := swap.foreignAtom.Audit(); err != nil {
		if err := swap.Complain(id); err != nil {
			swap.LogError(id, fmt.Sprintf("failed to complain to the watch dog: %v", err))
			return fmt.Errorf("failed to complain to the watch dog: %v", err)
		}
		return fmt.Errorf("incorrect swap complaint sent: %v", err)
	}
	swap.LogInfo(swap.req.UID, "Audited the atomic swap ")

	swap.LogInfo(swap.req.UID, "Initiating the atomic swap ")
	if err := swap.personalAtom.Initiate(); err != nil {
		if err != ErrSwapAlreadyInitiated {
			swap.LogError(id, fmt.Sprintf("failed to initiate atomic swap: %v", err))
			return fmt.Errorf("failed to initiate atomic swap: %v", err)
		}
		swap.LogInfo(id, "swap already initiated")
	}
	swap.LogInfo(swap.req.UID, "Initiated the atomic swap ")

	swap.LogInfo(swap.req.UID, "Auditing Secret ")
	secret, err := swap.personalAtom.AuditSecret()
	if err != nil {
		if err := swap.Complain(id); err != nil {
			swap.LogError(id, fmt.Sprintf("failed to complain to the watch dog: %v", err))
			return fmt.Errorf("failed to complain to the watch dog: %v", err)
		}
		return fmt.Errorf("incorrect swap complaint sent: %v", err)
	}
	swap.LogInfo(swap.req.UID, "Secret Audit success")

	swap.LogInfo(swap.req.UID, "Redeeming the atomic swap ")
	if err := swap.foreignAtom.Redeem(secret); err != nil {
		if err != ErrSwapAlreadyRedeemedOrRefunded {
			swap.LogError(id, fmt.Sprintf("failed to redeem: %v", err))
			return fmt.Errorf("failed to redeem: %v", err)
		}
		swap.LogInfo(id, "swap already redeemed or refunded")
	}

	swap.LogInfo(swap.req.UID, "Redeemed the atomic swap ")
	return nil
}
