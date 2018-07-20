package errors

import "fmt"

var (
	ErrNonRefundable = fmt.Errorf("Trying to refund a non refundable order")
)

func ErrAtomBuildFailed(err error) error {
	return fmt.Errorf("Failed to build the atom: %v", err)
}

func ErrRefundAfterRedeem(err error) error {
	return fmt.Errorf("Trying to an atomic swap refund after redeem: %v", err)
}
