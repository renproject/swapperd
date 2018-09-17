package swap

import (
	"encoding/json"
	"math/big"
	"sync"

	"github.com/republicprotocol/renex-swapper-go/adapter/store"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
	"github.com/republicprotocol/republic-go/logger"
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
	logger.Logger
	store.Store
}

type State interface {
	Status([32]byte) swap.Status
	PutStatus([32]byte, swap.Status) error
	InitiateDetails([32]byte) ([32]byte, int64, error)
	PutInitiateDetails([32]byte, [32]byte, int64) error
	RedeemDetails([32]byte) ([32]byte, error)
	PutRedeemDetails([32]byte, [32]byte) error
	AtomDetails([32]byte) ([]byte, error)
	PutAtomDetails([32]byte, []byte) error
	Redeemed([32]byte) error
	PutRedeemable([32]byte) error
}

func NewState(store store.Store, logger logger.Logger) State {
	return &state{
		Store:  store,
		Logger: logger,
		swapMu: new(sync.RWMutex),
	}
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

func (state *state) PutAtomDetails(orderID [32]byte, data []byte) error {
	return state.Write(append([]byte("Atom Details:"), orderID[:]...), data)
}

func (state *state) AtomDetails(orderID [32]byte) ([]byte, error) {
	return state.Read(append([]byte("Atom Details:"), orderID[:]...))
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
