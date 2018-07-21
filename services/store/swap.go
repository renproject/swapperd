package store

import (
	"bytes"
	"encoding/json"
	"math/big"
	"sync"

	"github.com/republicprotocol/atom-go/domains/match"
)

// SwapStatus stores the swap status
type SwapStatus struct {
	Status string `json:"status"`
}

// SwapInitiateDetails stores the swap status
type SwapInitiateDetails struct {
	Expiry   int64    `json:"expiry"`
	HashLock [32]byte `json:"hashLock"`
}

type SwapRedeemDetails struct {
	Secret [32]byte `json:"secret"`
}

// SwapMatch stores the swap status
type SwapMatch struct {
	PersonalOrderID [32]byte `json:"personalOrderID"`
	ForeignOrderID  [32]byte `json:"foreignOrderID"`
	SendValue       *big.Int `json:"sendValue"`
	ReceiveValue    *big.Int `json:"receiveValue"`
	SendCurrency    uint32   `json:"sendCurrency"`
	ReceiveCurrency uint32   `json:"receiveCurrency"`
}

// PendingSwaps stores all the swaps that are pending
type PendingSwaps struct {
	Swaps [][32]byte `json:"pendingSwaps"`
}

type state struct {
	Store
	swapMu *sync.RWMutex
}

type SwapState interface {
	AddSwap([32]byte) error
	DeleteSwap([32]byte) error
	PendingSwaps() ([][32]byte, error)

	InitiateDetails([32]byte) (int64, [32]byte, error)
	PutInitiateDetails([32]byte, int64, [32]byte) error

	RedeemDetails([32]byte) ([32]byte, error)
	PutRedeemDetails([32]byte, [32]byte) error

	Status([32]byte) string
	PutStatus([32]byte, string) error

	Match([32]byte) (match.Match, error)
	PutMatch([32]byte, match.Match) error

	AtomDetails([32]byte) ([]byte, error)
	PutAtomDetails([32]byte, []byte) error
	AtomExists([32]byte) bool
}

func NewSwapState(store Store) SwapState {
	return &state{
		Store:  store,
		swapMu: new(sync.RWMutex),
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
			if i == 0 {
				pendingSwaps.Swaps = pendingSwaps.Swaps[i:]
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

func (state *state) PendingSwaps() ([][32]byte, error) {
	state.swapMu.Lock()
	defer state.swapMu.Unlock()

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

func (state *state) PutInitiateDetails(orderID [32]byte, expiry int64, hashLock [32]byte) error {
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

func (state *state) InitiateDetails(orderID [32]byte) (int64, [32]byte, error) {
	initiateDetailsBytes, err := state.Read(append([]byte("Initiate Details:"), orderID[:]...))
	if err != nil {
		return 0, [32]byte{}, err
	}
	swapInitiateDetails := SwapInitiateDetails{}

	if err := json.Unmarshal(initiateDetailsBytes, &swapInitiateDetails); err != nil {
		return 0, [32]byte{}, err
	}

	return swapInitiateDetails.Expiry, swapInitiateDetails.HashLock, nil
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

func (state *state) PutStatus(orderID [32]byte, status string) error {
	swapStatus := SwapStatus{
		Status: status,
	}
	statusBytes, err := json.Marshal(swapStatus)
	if err != nil {
		return err
	}
	return state.Write(append([]byte("Status:"), orderID[:]...), statusBytes)
}

func (state *state) Status(orderID [32]byte) string {
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
