package swap

import "errors"

var ErrSwapAlreadyInitiated = errors.New("Duplicate swap initiation")
var ErrSwapAlreadyRedeemedOrRefunded = errors.New("The swap is already redeemed or refunded")
