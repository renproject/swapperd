package callback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/renproject/swapperd/core/wallet/swapper/delayed"
	"github.com/renproject/swapperd/foundation/swap"
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

	switch resp.StatusCode {
	case http.StatusOK:
		filledSwap := swap.SwapBlob{}
		if err := json.Unmarshal(respBytes, &filledSwap); err != nil {
			return partialSwap, err
		}
		return verifyDelaySwap(partialSwap, filledSwap)
	case http.StatusNoContent:
		return partialSwap, delayed.ErrSwapDetailsUnavailable
	case http.StatusGone:
		return partialSwap, delayed.ErrSwapCancelled
	default:
		return partialSwap, fmt.Errorf("unexpected status code=%v: %v", resp.StatusCode, string(respBytes))
	}
}

func verifyDelaySwap(partialSwap, filledSwap swap.SwapBlob) (swap.SwapBlob, error) {
	initialMinReceiveValue, ok := new(big.Int).SetString(partialSwap.MinimumReceiveAmount, 10)
	if !ok {
		initialMinReceiveValue = big.NewInt(0)
	}

	initialSendValue, ok := new(big.Int).SetString(partialSwap.SendAmount, 10)
	if !ok {
		return partialSwap, fmt.Errorf("corrupted send value")
	}

	initialRecvValue, ok := big.NewInt(0).SetString(partialSwap.ReceiveAmount, 10)
	if !ok {
		return partialSwap, fmt.Errorf("corrupted receive value")
	}

	filledSendValue, ok := big.NewInt(0).SetString(filledSwap.SendAmount, 10)
	if !ok {
		return partialSwap, fmt.Errorf("corrupted filled send value")
	}

	filledReceiveValue, ok := big.NewInt(0).SetString(filledSwap.ReceiveAmount, 10)
	if !ok {
		return partialSwap, fmt.Errorf("corrupted filled receive value")
	}

	if filledReceiveValue.Cmp(initialMinReceiveValue) < 0 || initialSendValue.Cmp(filledSendValue) < 0 {
		return partialSwap, fmt.Errorf("invalid filled swap receive value too low or send value too high %v %v %v %v", initialMinReceiveValue, filledReceiveValue, initialSendValue, filledSendValue)
	}

	if validatePrice(filledReceiveValue, filledSendValue, initialRecvValue, initialSendValue, partialSwap.DelayRange) {
		return partialSwap, fmt.Errorf("invalid filled swap unfavorable price")
	}

	filledSwap.Delay = false
	return filledSwap, nil
}

func validatePrice(filledReceiveValue, filledSendValue, initialRecvValue, initialSendValue *big.Int, delayRange int64) bool {
	return new(big.Int).Mul(new(big.Int).Mul(filledReceiveValue, initialSendValue), big.NewInt(10000)).Cmp(new(big.Int).Mul(new(big.Int).Mul(initialRecvValue, filledSendValue), big.NewInt(10000-delayRange))) >= 0
}
