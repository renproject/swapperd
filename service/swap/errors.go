package swap

import "errors"

var ErrSwapAlreadyInitiated = errors.New("Duplicate swap initiation")
var ErrNotRefundable = errors.New("Swap not refundable")
