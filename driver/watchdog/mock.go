package watchdog

import "github.com/republicprotocol/renex-swapper-go/adapter/watchdog"

type mock struct {
}

func NewMock() watchdog.Watchdog {
	return &mock{}
}

func (mock *mock) ComplainDelayedAddressSubmission(orderID [32]byte) error {
	return nil
}

func (mock *mock) ComplainDelayedRequestorInitiation(orderID [32]byte) error {
	return nil
}

func (mock *mock) ComplainWrongRequestorInitiation(orderID [32]byte) error {
	return nil
}

func (mock *mock) ComplainDelayedResponderInitiation(orderID [32]byte) error {
	return nil
}

func (mock *mock) ComplainWrongResponderInitiation(orderID [32]byte) error {
	return nil
}

func (mock *mock) ComplainDelayedRequestorRedemption(orderID [32]byte) error {
	return nil
}
