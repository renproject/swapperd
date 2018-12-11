package callback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/republicprotocol/swapperd/core/delayed"
	"github.com/republicprotocol/swapperd/foundation/swap"
)

type cb struct {
}

func New() delayed.DelayCallback {
	return &cb{}
}

func (cb *cb) DelayCallback(partialSwap swap.SwapBlob) (swap.SwapBlob, error) {
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

	if resp.StatusCode == http.StatusNoContent {
		return partialSwap, delayed.ErrSwapDetailsUnavailable
	}

	if resp.StatusCode == http.StatusGone {
		return partialSwap, delayed.ErrSwapCancelled
	}

	return partialSwap, fmt.Errorf("unexpected error %d: %s", resp.StatusCode, respBytes)
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

	if filledSwap.BrokerFee > partialSwap.BrokerFee {
		return partialSwap, fmt.Errorf("invalid filled swap unfavourable broker fee")
	}

	filledSwap.Delay = false
	return filledSwap, nil
}
