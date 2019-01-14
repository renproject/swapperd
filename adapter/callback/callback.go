package callback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/republicprotocol/swapperd/core/swapper/delayed"
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

	// initialMinReceiveValue, ok := new(big.Int).SetString(partialSwap.MinimumReceiveAmount, 10)
	// if !ok {
	// 	initialMinReceiveValue = big.NewInt(0)
	// }

	initialSendValue, ok := new(big.Int).SetString(partialSwap.SendAmount, 10)
	if !ok {
		return partialSwap, fmt.Errorf("corrupted send value")
	}

	initialRecvValue, ok := big.NewInt(0).SetString(partialSwap.ReceiveAmount, 10)
	if !ok {
		return partialSwap, fmt.Errorf("corrupted receive value")
	}

	sendBrokerFee := new(big.Int).Div(new(big.Int).Mul(initialSendValue, big.NewInt(partialSwap.BrokerFee)), big.NewInt(10000))
	recvBrokerFee := new(big.Int).Div(new(big.Int).Mul(initialRecvValue, big.NewInt(partialSwap.BrokerFee)), big.NewInt(10000))
	// minRecvBrokerFee := new(big.Int).Div(new(big.Int).Mul(initialMinReceiveValue, big.NewInt(partialSwap.BrokerFee)), big.NewInt(10000))

	actualSendValue := new(big.Int).Sub(initialSendValue, sendBrokerFee)
	actualRecvValue := new(big.Int).Sub(initialRecvValue, recvBrokerFee)
	// actualMinRecvValue := new(big.Int).Sub(initialMinReceiveValue, minRecvBrokerFee)

	filledSendValue, ok := big.NewInt(0).SetString(filledSwap.SendAmount, 10)
	if !ok {
		return partialSwap, fmt.Errorf("corrupted filled send value")
	}

	filledReceiveValue, ok := big.NewInt(0).SetString(filledSwap.ReceiveAmount, 10)
	if !ok {
		return partialSwap, fmt.Errorf("corrupted filled receive value")
	}

	// if filledReceiveValue.Cmp(actualMinRecvValue) < 0 || actualSendValue.Cmp(filledSendValue) > 0 {
	// 	return partialSwap, fmt.Errorf("invalid filled swap receive value too low or send value too high %v %v %v %v", actualMinRecvValue, filledReceiveValue, actualSendValue, filledSendValue)
	// }

	if filledReceiveValue.Mul(filledReceiveValue, actualSendValue).Cmp(actualRecvValue.Mul(actualRecvValue, filledSendValue)) < 0 {
		return partialSwap, fmt.Errorf("invalid filled swap unfavorable price")
	}

	filledSwap.Delay = false
	return filledSwap, nil
}
