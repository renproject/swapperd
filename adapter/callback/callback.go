package callback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation/swap"
)

type callback struct {
}

func New() swapper.DelayCallback {
	return &callback{}
}

func (callback *callback) DelayCallback(partialSwap swap.SwapBlob) (swap.SwapBlob, error) {
	for {
		data, err := json.MarshalIndent(partialSwap, "", "  ")
		if err != nil {
			return partialSwap, err
		}
		buf := bytes.NewBuffer(data)

		resp, err := http.Post(fmt.Sprintf(partialSwap.DelayCallbackURL), "application/json", buf)
		if err != nil {
			return partialSwap, err
		}

		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return partialSwap, err
		}

		filledSwap := swap.SwapBlob{}
		if resp.StatusCode == 200 {
			if err := json.Unmarshal(respBytes, &filledSwap); err != nil {
				return partialSwap, err
			}
			return verifyDelaySwap(partialSwap, filledSwap)
		}

		if resp.StatusCode == 404 {
			time.Sleep(30 * time.Second)
			continue
		}

		return partialSwap, fmt.Errorf("unexpected error %d: %s", resp.StatusCode, respBytes)
	}
}

func verifyDelaySwap(partialSwap, filledSwap swap.SwapBlob) (swap.SwapBlob, error) {
	initialMinReceiveValue, ok := big.NewInt(0).SetString(partialSwap.MinimumReceiveAmount, 10)
	if !ok {
		return partialSwap, fmt.Errorf("corrupted minimum receive value")
	}

	initialMaxSendValue, ok := big.NewInt(0).SetString(partialSwap.SendAmount, 10)
	if !ok {
		return partialSwap, fmt.Errorf("corrupted send value")
	}

	initialReceiveValue, ok := big.NewInt(0).SetString(partialSwap.ReceiveAmount, 10)
	if !ok {
		return partialSwap, fmt.Errorf("corrupted receive value")
	}

	filledReceiveValue, ok := big.NewInt(0).SetString(filledSwap.ReceiveAmount, 10)
	if !ok {
		return partialSwap, fmt.Errorf("corrupted filled receive value")
	}

	filledSendValue, ok := big.NewInt(0).SetString(filledSwap.SendAmount, 10)
	if !ok {
		return partialSwap, fmt.Errorf("corrupted filled send value")
	}

	if initialMinReceiveValue.Cmp(filledReceiveValue) > 0 || initialMaxSendValue.Cmp(filledSendValue) < 0 {
		return partialSwap, fmt.Errorf("invalid filled swap receive value too low or send value too high")
	}

	if filledReceiveValue.Mul(filledReceiveValue, initialMaxSendValue).Cmp(initialReceiveValue.Mul(initialReceiveValue, filledSendValue)) < 0 {
		return partialSwap, fmt.Errorf("invalid filled swap unfavorable price")
	}

	filledSwap.Delay = false

	return filledSwap, nil
}
