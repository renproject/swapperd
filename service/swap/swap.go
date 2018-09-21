package swap

import (
	"fmt"
)

type Swap interface {
	Execute() error
}

type swap struct {
	req          Request
	personalAtom Atom
	foreignAtom  Atom
	Adapter
}

func (swap *swap) Execute() error {
	if swap.req.GoesFirst {
		return swap.request(swap.req.UID, swap.req.Secret)
	}
	return swap.respond(swap.req.UID)
}

func (swap *swap) request(id [32]byte, secret [32]byte) error {
	if err := swap.personalAtom.Initiate(); err != nil {
		if err != ErrSwapAlreadyInitiated {
			swap.LogError(id, fmt.Sprintf("failed to initiate details: %v", err))
			return fmt.Errorf("failed to initiate details: %v", err)
		}
		swap.LogInfo(id, "swap already initiated")
	}
	if err := swap.foreignAtom.Audit(); err != nil {
		if err := swap.ComplainWrongResponderInitiation(id); err != nil {
			swap.LogError(id, fmt.Sprintf("failed to complain to the watch dog: %v", err))
			return fmt.Errorf("failed to complain to the watch dog: %v", err)
		}
		return fmt.Errorf("incorrect swap complained to the watchdog: %v", err)
	}
	if err := swap.foreignAtom.Redeem(secret); err != nil {
		swap.LogError(id, fmt.Sprintf("failed to redeem: %v", err))
		return fmt.Errorf("failed to redeem: %v", err)
	}
	return nil
}

func (swap *swap) respond(id [32]byte) error {
	if err := swap.foreignAtom.Audit(); err != nil {
		if err := swap.ComplainWrongRequestorInitiation(id); err != nil {
			swap.LogError(id, fmt.Sprintf("failed to complain to the watch dog: %v", err))
			return fmt.Errorf("failed to complain to the watch dog: %v", err)
		}
		return fmt.Errorf("incorrect swap complained to the watchdog: %v", err)
	}
	if err := swap.personalAtom.Initiate(); err != nil {
		if err != ErrSwapAlreadyInitiated {
			swap.LogError(id, fmt.Sprintf("failed to initiate atomic swap: %v", err))
			return fmt.Errorf("failed to initiate atomic swap: %v", err)
		}
		swap.LogInfo(id, "swap already initiated")
	}
	secret, err := swap.personalAtom.AuditSecret()
	if err != nil {
		if err := swap.ComplainDelayedRequestorRedemption(id); err != nil {
			swap.LogError(id, fmt.Sprintf("failed to complain to the watch dog: %v", err))
			return fmt.Errorf("failed to complain to the watch dog: %v", err)
		}
		return fmt.Errorf("incorrect swap complained to the watchdog: %v", err)
	}
	if err := swap.foreignAtom.Redeem(secret); err != nil {
		swap.LogError(id, fmt.Sprintf("failed to redeem: %v", err))
		return fmt.Errorf("failed to redeem: %v", err)
	}
	return nil
}
