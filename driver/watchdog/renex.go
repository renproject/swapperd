package watchdog

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/watchdog"
)

type renexWatchdog struct {
	ipAddress string
}

// NewRenEx creates a new RenExWatchdog object, that interacts with the
// RenEx watchdog over http.
func NewRenEx(config config.Config) watchdog.Watchdog {
	return &renexWatchdog{
		ipAddress: config.RenEx.Watchdog,
	}
}

func (client *renexWatchdog) ComplainDelayedAddressSubmission(orderID [32]byte) error {
	return client.watch(orderID)
}

func (client *renexWatchdog) ComplainDelayedRequestorInitiation(orderID [32]byte) error {
	return client.watch(orderID)
}

func (client *renexWatchdog) ComplainWrongRequestorInitiation(orderID [32]byte) error {
	return client.watch(orderID)
}

func (client *renexWatchdog) ComplainDelayedResponderInitiation(orderID [32]byte) error {
	return client.watch(orderID)
}

func (client *renexWatchdog) ComplainWrongResponderInitiation(orderID [32]byte) error {
	return client.watch(orderID)
}

func (client *renexWatchdog) ComplainDelayedRequestorRedemption(orderID [32]byte) error {
	return client.watch(orderID)
}

func (client *renexWatchdog) watch(orderID [32]byte) error {
	resp, err := http.Post(fmt.Sprintf("https://"+client.ipAddress+"/watch?orderID="+hex.EncodeToString(orderID[:])), "text", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		return nil
	}
	return fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
}
