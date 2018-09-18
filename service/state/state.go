package state

import (
	"bytes"
	"encoding/json"
	"math/big"
	"sync"
	"time"

	"github.com/republicprotocol/renex-swapper-go/domain/match"
	"github.com/republicprotocol/renex-swapper-go/domain/swap"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
)

// SwapStatus stores the swap status
type SwapStatus struct {
	Status swap.Status `json:"status"`
}

// OrderTimeStamp stores the timestamp when the order is submitted to the
// swapper
type OrderTimeStamp struct {
	TimeStamp int64 `json:"timeStamp"`
}

// SwapInitiateDetails stores the initiate details
type SwapInitiateDetails struct {
	Expiry   int64    `json:"expiry"`
	HashLock [32]byte `json:"hashLock"`
}

// SwapRedeemDetails stores the redeem details
type SwapRedeemDetails struct {
	Secret [32]byte `json:"secret"`
}

// SwapMatch stores the swap status
type SwapMatch struct {
	PersonalOrderID [32]byte    `json:"personalOrderID"`
	ForeignOrderID  [32]byte    `json:"foreignOrderID"`
	SendValue       *big.Int    `json:"sendValue"`
	ReceiveValue    *big.Int    `json:"receiveValue"`
	SendCurrency    token.Token `json:"sendCurrency"`
	ReceiveCurrency token.Token `json:"receiveCurrency"`
}

// PendingSwaps stores all the swaps that are pending
type PendingSwaps struct {
	Swaps [][32]byte `json:"pendingSwaps"`
}

type state struct {
	swapMu *sync.RWMutex
	Adapter
}

type State interface {
	AddSwap([32]byte) error
	DeleteSwap([32]byte) error
	ExecutableSwaps(bool) ([][32]byte, error)
	RefundableSwaps() ([][32]byte, error)

	InitiateDetails([32]byte) ([32]byte, int64, error)
	PutInitiateDetails([32]byte, [32]byte, int64) error

	RedeemDetails([32]byte) ([32]byte, error)
	PutRedeemDetails([32]byte, [32]byte) error

	Status([32]byte) swap.Status
	PutStatus([32]byte, swap.Status) error

	Match([32]byte) (match.Match, error)
	PutMatch([32]byte, match.Match) error

	OrderTimeStamp([32]byte) (int64, error)
	PutOrderTimeStamp([32]byte) error

	AtomDetails([32]byte) ([]byte, error)
	PutAtomDetails([32]byte, []byte) error
	AtomExists([32]byte) bool

	PutRedeemable([32]byte) error
	IsRedeemable([32]byte) bool
	Complained([32]byte) bool
	Redeemed([32]byte) error
}

func NewState(adapter Adapter) State {
	return &state{
		Adapter: adapter,
		swapMu:  new(sync.RWMutex),
	}
}

func (state *state) AddSwap(orderID [32]byte) error {
	state.swapMu.Lock()
	defer state.swapMu.Unlock()
	pendingSwapsRawBytes, err := state.Read([]byte("Pending Swaps:"))
	pendingSwaps := PendingSwaps{}

	if err == nil {
		if err := json.Unmarshal(pendingSwapsRawBytes, &pendingSwaps); err != nil {
			return err
		}
	}

	pendingSwaps.Swaps = append(pendingSwaps.Swaps, orderID)

	pendingSwapsProcessedBytes, err := json.Marshal(pendingSwaps)
	if err != nil {
		return err
	}

	return state.Write([]byte("Pending Swaps:"), pendingSwapsProcessedBytes)
}

func (state *state) DeleteSwap(orderID [32]byte) error {
	state.swapMu.Lock()
	defer state.swapMu.Unlock()
	defer state.LogInfo(orderID, "deleted the swap and its details")

	pendingSwapsRawBytes, err := state.Read([]byte("Pending Swaps:"))
	if err != nil {
		return err
	}

	pendingSwaps := PendingSwaps{}
	if err := json.Unmarshal(pendingSwapsRawBytes, &pendingSwaps); err != nil {
		return err
	}

	for i, swap := range pendingSwaps.Swaps {
		if swap == orderID {
			if len(pendingSwaps.Swaps) == 1 {
				pendingSwaps.Swaps = [][32]byte{}
				break
			}

			if i == 0 {
				pendingSwaps.Swaps = pendingSwaps.Swaps[1:]
				break
			}

			pendingSwaps.Swaps = append(pendingSwaps.Swaps[:i-1], pendingSwaps.Swaps[i:]...)
			break
		}
	}

	pendingSwapsProcessedBytes, err := json.Marshal(pendingSwaps)
	if err != nil {
		return err
	}

	return state.Write([]byte("Pending Swaps:"), pendingSwapsProcessedBytes)
}

func (state *state) pendingSwaps() ([][32]byte, error) {
	pendingSwapsBytes, err := state.Read([]byte("Pending Swaps:"))
	if err != nil {
		return [][32]byte{}, nil
	}

	pendingSwaps := PendingSwaps{}
	if err := json.Unmarshal(pendingSwapsBytes, &pendingSwaps); err != nil {
		return nil, err
	}

	return pendingSwaps.Swaps, nil
}

func (state *state) ExecutableSwaps(fullsync bool) ([][32]byte, error) {
	state.swapMu.RLock()
	pendingSwaps, err := state.pendingSwaps()
	state.swapMu.RUnlock()
	state.swapMu.Lock()
	defer state.swapMu.Unlock()
	if err != nil {
		return nil, err
	}
	if fullsync {
		return state.executableFullSync(pendingSwaps)
	}
	return state.executablePartialSync(pendingSwaps)
}

func (state *state) executableFullSync(pendingSwaps [][32]byte) ([][32]byte, error) {
	executableSwaps := [][32]byte{}
	for _, pendingSwap := range pendingSwaps {
		if state.Status(pendingSwap) != "COMPLAINED" {
			executableSwaps = append(executableSwaps, pendingSwap)
		}
	}
	return executableSwaps, nil
}

func (state *state) executablePartialSync(pendingSwaps [][32]byte) ([][32]byte, error) {
	executableSwaps := [][32]byte{}
	for _, pendingSwap := range pendingSwaps {
		if state.Status(pendingSwap) == "UNKNOWN" {
			executableSwaps = append(executableSwaps, pendingSwap)
		}
	}
	return executableSwaps, nil
}

func (state *state) RefundableSwaps() ([][32]byte, error) {
	state.swapMu.RLock()
	pendingSwaps, err := state.pendingSwaps()
	state.swapMu.RUnlock()
	if err != nil {
		return nil, err
	}
	refundableSwaps := [][32]byte{}
	for _, pendingSwap := range pendingSwaps {
		if state.Status(pendingSwap) == "COMPLAINED" {
			refundableSwaps = append(refundableSwaps, pendingSwap)
		}
	}
	return refundableSwaps, nil
}

func (state *state) PutInitiateDetails(orderID [32]byte, hashLock [32]byte, expiry int64) error {
	swapInitiateDetails := SwapInitiateDetails{
		Expiry:   expiry,
		HashLock: hashLock,
	}
	initiateDetailsBytes, err := json.Marshal(swapInitiateDetails)
	if err != nil {
		return err
	}
	return state.Write(append([]byte("Initiate Details:"), orderID[:]...), initiateDetailsBytes)
}

func (state *state) InitiateDetails(orderID [32]byte) ([32]byte, int64, error) {
	initiateDetailsBytes, err := state.Read(append([]byte("Initiate Details:"), orderID[:]...))
	if err != nil {
		return [32]byte{}, 0, err
	}
	swapInitiateDetails := SwapInitiateDetails{}

	if err := json.Unmarshal(initiateDetailsBytes, &swapInitiateDetails); err != nil {
		return [32]byte{}, 0, err
	}

	return swapInitiateDetails.HashLock, swapInitiateDetails.Expiry, nil
}

func (state *state) PutRedeemDetails(orderID [32]byte, secret [32]byte) error {
	swapRedeemDetails := SwapRedeemDetails{
		Secret: secret,
	}
	redeemDetailsBytes, err := json.Marshal(swapRedeemDetails)
	if err != nil {
		return err
	}
	return state.Write(append([]byte("Redeem Details:"), orderID[:]...), redeemDetailsBytes)
}

func (state *state) RedeemDetails(orderID [32]byte) ([32]byte, error) {
	redeemDetailsBytes, err := state.Read(append([]byte("Redeem Details:"), orderID[:]...))
	if err != nil {
		return [32]byte{}, err
	}
	swapRedeemDetails := SwapRedeemDetails{}

	if err := json.Unmarshal(redeemDetailsBytes, &swapRedeemDetails); err != nil {
		return [32]byte{}, err
	}

	return swapRedeemDetails.Secret, nil
}

func (state *state) PutStatus(orderID [32]byte, status swap.Status) error {
	swapStatus := SwapStatus{
		Status: status,
	}
	statusBytes, err := json.Marshal(swapStatus)
	if err != nil {
		return err
	}
	return state.Write(append([]byte("Status:"), orderID[:]...), statusBytes)
}

func (state *state) Status(orderID [32]byte) swap.Status {
	statusBytes, err := state.Read(append([]byte("Status:"), orderID[:]...))
	if err != nil {
		return "UNKNOWN"
	}
	swapStatus := SwapStatus{}

	if err := json.Unmarshal(statusBytes, &swapStatus); err != nil {
		return "UNKNOWN"
	}
	return swapStatus.Status
}

func (state *state) PutMatch(orderID [32]byte, m match.Match) error {
	match := SwapMatch{
		PersonalOrderID: m.PersonalOrderID(),
		ForeignOrderID:  m.ForeignOrderID(),
		SendValue:       m.SendValue(),
		ReceiveValue:    m.ReceiveValue(),
		SendCurrency:    m.SendCurrency(),
		ReceiveCurrency: m.ReceiveCurrency(),
	}

	matchBytes, err := json.Marshal(match)
	if err != nil {
		return err
	}
	return state.Write(append([]byte("Match:"), orderID[:]...), matchBytes)
}

func (state *state) Match(orderID [32]byte) (match.Match, error) {
	matchBytes, err := state.Read(append([]byte("Match:"), orderID[:]...))
	if err != nil {
		return nil, err
	}
	swapMatch := SwapMatch{}

	if err := json.Unmarshal(matchBytes, &swapMatch); err != nil {
		return nil, err
	}
	return match.NewMatch(swapMatch.PersonalOrderID, swapMatch.ForeignOrderID, swapMatch.SendValue, swapMatch.ReceiveValue, swapMatch.SendCurrency, swapMatch.ReceiveCurrency), nil
}

func (state *state) PutAtomDetails(orderID [32]byte, data []byte) error {
	return state.Write(append([]byte("Atom Details:"), orderID[:]...), data)
}

func (state *state) AtomDetails(orderID [32]byte) ([]byte, error) {
	return state.Read(append([]byte("Atom Details:"), orderID[:]...))
}

func (state *state) AtomExists(orderID [32]byte) bool {
	atomDerails, err := state.AtomDetails(orderID)
	if err != nil || bytes.Compare(atomDerails, []byte{}) == 0 {
		return false
	}
	return true
}

func (state *state) IsRedeemable(orderID [32]byte) bool {
	_, err := state.Read(append([]byte("Redeemable"), orderID[:]...))
	if err != nil {
		return false
	}
	return true
}

func (state *state) Complained(orderID [32]byte) bool {
	statusBytes, err := state.Read(append([]byte("Status:"), orderID[:]...))
	if err != nil {
		return false
	}
	if string(statusBytes) != "COMPLAINED" {
		return false
	}
	return true
}

func (state *state) PutRedeemable(orderID [32]byte) error {
	return state.Write(append([]byte("Redeemable"), orderID[:]...), orderID[:])
}

func (state *state) Redeemed(orderID [32]byte) error {
	err := state.Delete(append([]byte("Redeemable"), orderID[:]...))
	if err != nil {
		return err
	}
	return nil
}

func (state *state) PutOrderTimeStamp(orderID [32]byte) error {
	data, err := json.Marshal(
		OrderTimeStamp{
			TimeStamp: time.Now().Unix(),
		},
	)
	if err != nil {
		return err
	}
	return state.Write(append([]byte("OrderTimeStamp:"), orderID[:]...), data)
}

func (state *state) OrderTimeStamp(orderID [32]byte) (int64, error) {
	tsBytes, err := state.Read(append([]byte("OrderTimeStamp:"), orderID[:]...))
	if err != nil {
		return 0, err
	}
	timeStamp := OrderTimeStamp{}

	if err := json.Unmarshal(tsBytes, &timeStamp); err != nil {
		return 0, err
	}
	return timeStamp.TimeStamp, nil
}
