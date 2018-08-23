package errors

import (
	"fmt"
)

const errPrefix = "Persistent Storage Error: "

func ErrFailedPendingSwaps(err error) error {
	return fmt.Errorf("%sFailed To get pending swaps: %v", errPrefix, err)
}

func ErrFailedInitiateDetails(err error) error {
	return fmt.Errorf("%sFailed To get initiate details: %v", errPrefix, err)
}
