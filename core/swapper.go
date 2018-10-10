package core

import (
	"fmt"

	"github.com/republicprotocol/swapperd/foundation"
)

type SwapContractBinder interface {
	Initiate() error
	Audit() error
	Redeem([32]byte) error
	AuditSecret() ([32]byte, error)
	Refund() error
}

type Logger interface {
	LogInfo(foundation.SwapID, string)
	LogDebug(foundation.SwapID, string)
	LogError(foundation.SwapID, string)
}

type Result struct {
	ID    foundation.SwapID
	Retry bool
}

func Swap(native, foreign SwapContractBinder, logger Logger, req foundation.Swap, done chan<- Result) {
	if req.IsFirst {
		if err := native.Initiate(); err != nil {
			logger.LogError(req.ID, fmt.Sprintf("Initiate failed: %v", err))
			done <- Result{req.ID, false}
			return
		}
		if err := foreign.Audit(); err != nil {
			if err := native.Refund(); err != nil {
				logger.LogError(req.ID, fmt.Sprintf("Refund failed: %v", err))
				done <- Result{req.ID, false}
				return
			}
			done <- Result{req.ID, true}
			return
		}
		if err := foreign.Redeem(req.Secret); err != nil {
			logger.LogError(req.ID, fmt.Sprintf("Redeem failed: %v", err))
			done <- Result{req.ID, false}
			return
		}
		done <- Result{req.ID, true}
		return
	}
	if err := foreign.Audit(); err != nil {
		done <- Result{req.ID, true}
		return
	}
	if err := native.Initiate(); err != nil {
		logger.LogError(req.ID, fmt.Sprintf("Initiate failed: %v", err))
		done <- Result{req.ID, false}
		return
	}
	secret, err := native.AuditSecret()
	if err != nil {
		if err := native.Refund(); err != nil {
			logger.LogError(req.ID, fmt.Sprintf("Refund failed: %v", err))
			done <- Result{req.ID, false}
			return
		}
		done <- Result{req.ID, true}
	}
	if err := foreign.Redeem(secret); err != nil {
		logger.LogError(req.ID, fmt.Sprintf("Redeem failed: %v", err))
		done <- Result{req.ID, false}
		return
	}
	done <- Result{req.ID, true}
}
