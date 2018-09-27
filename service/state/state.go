package state

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/republicprotocol/renex-swapper-go/domain/swap"
)

type ProtectedSwapMap struct {
	mu    *sync.RWMutex
	cache map[[32]byte]SwapDetails
}

func (swapMap *ProtectedSwapMap) Read(id [32]byte) SwapDetails {
	swapMap.mu.RLock()
	defer swapMap.mu.RUnlock()
	return swapMap.cache[id]
}

func (swapMap *ProtectedSwapMap) Write(id [32]byte, det SwapDetails) {
	swapMap.mu.Lock()
	defer swapMap.mu.Unlock()
	swapMap.cache[id] = det
}

type SwapDetails struct {
	TimeStamp      int64        `json:"timeStamp"`
	Status         swap.Status  `json:"status"`
	ForeignOrderID [32]byte     `json:"foreignOrderID"`
	Request        swap.Request `json:"request"`
}

type state struct {
	swapCache ProtectedSwapMap
	Adapter
}

type State interface {
	ActiveSwapList

	Status([32]byte) swap.Status
	PutStatus([32]byte, swap.Status) error

	SwapRequest([32]byte) (swap.Request, error)
	PutSwapRequest([32]byte, swap.Request) error

	AddTimestamp([32]byte) (int64, error)
	PutAddTimestamp([32]byte) error

	PrintSwapRequest(swap.Request)
}

// NewState creates a new state interface
func NewState(adapter Adapter) State {
	return &state{
		Adapter: adapter,
		swapCache: ProtectedSwapMap{
			cache: map[[32]byte]SwapDetails{},
			mu:    new(sync.RWMutex),
		},
	}
}

// PutStatus into both persistent storage and in memory cache.
func (state *state) PutStatus(orderID [32]byte, status swap.Status) error {
	swap := state.swapCache.Read(orderID)
	swap.Status = status
	if err := state.WriteSwapDetails(orderID, swap); err != nil {
		return err
	}
	state.swapCache.Write(orderID, swap)
	return nil
}

// Status tries to get the swap status from in memory cache if it does not exist
// it tries to read from the persistent storage.
func (state *state) Status(orderID [32]byte) swap.Status {
	swapDetails := state.swapCache.Read(orderID)
	if swapDetails.Status != swap.Status("") {
		return swapDetails.Status
	}
	swapDetails = state.ReadSwapDetails(orderID)
	state.swapCache.Write(orderID, swapDetails)
	if swapDetails.Status == swap.Status("") {
		return swap.StatusOpen
	}
	return swapDetails.Status
}

// PutAddTimestamp into both persistent storage and in memory cache.
func (state *state) PutAddTimestamp(orderID [32]byte) error {
	swap := state.swapCache.Read(orderID)
	swap.TimeStamp = time.Now().Unix()
	if err := state.WriteSwapDetails(orderID, swap); err != nil {
		return err
	}
	state.swapCache.Write(orderID, swap)
	return nil
}

// AddTimestamp tries to get the redeem details from in memory cache if
// they do not exist it tries to read from the persistent storage.
func (state *state) AddTimestamp(orderID [32]byte) (int64, error) {
	swap := state.swapCache.Read(orderID)
	if swap.TimeStamp != 0 {
		return swap.TimeStamp, nil
	}
	swap = state.ReadSwapDetails(orderID)
	state.swapCache.Write(orderID, swap)
	return swap.TimeStamp, nil
}

// PutSwapRequest into persistent storage
func (state *state) PutSwapRequest(orderID [32]byte, req swap.Request) error {
	swap := state.swapCache.Read(orderID)
	swap.Request = req
	if err := state.WriteSwapDetails(orderID, swap); err != nil {
		return err
	}
	state.swapCache.Write(orderID, swap)
	return nil
}

// ReadSwapDetails from persistent storage
func (state *state) SwapRequest(orderID [32]byte) (swap.Request, error) {
	req := swap.Request{}
	swapDetails := state.swapCache.Read(orderID)
	if swapDetails.Request != req {
		return swapDetails.Request, nil
	}
	swapDetails = state.ReadSwapDetails(orderID)
	state.swapCache.Write(orderID, swapDetails)
	return swapDetails.Request, nil
}

// WriteSwapDetails to persistent storage
func (state *state) WriteSwapDetails(orderID [32]byte, swapDetails SwapDetails) error {
	data, err := json.Marshal(swapDetails)
	if err != nil {
		return err
	}
	return state.Write(append([]byte("Swap Details:"), orderID[:]...), data)
}

// ReadSwapDetails from persistent storage
func (state *state) ReadSwapDetails(orderID [32]byte) SwapDetails {
	data, err := state.Read(append([]byte("Swap Details:"), orderID[:]...))
	if err != nil {
		return SwapDetails{}
	}
	swapDetails := SwapDetails{}
	if err = json.Unmarshal(data, &swapDetails); err != nil {
		return SwapDetails{}
	}
	return swapDetails
}

// PrintSwapRequest to Std Out
func (state *state) PrintSwapRequest(swapRequest swap.Request) {
	state.LogInfo(swapRequest.UID, fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s%s",
		fmt.Sprintf("\nSWAP REQUEST\n\n"),
		fmt.Sprintf("UID: %s\n", hex.EncodeToString(swapRequest.UID[:])),
		fmt.Sprintf("Expiry: %d\n", swapRequest.TimeLock),
		fmt.Sprintf("Secret Hash: %s\n", base64.StdEncoding.EncodeToString(swapRequest.SecretHash[:])),
		fmt.Sprintf("Secret: %s\n", base64.StdEncoding.EncodeToString(swapRequest.Secret[:])),
		fmt.Sprintf("Send To Address: %v\n", swapRequest.SendToAddress),
		fmt.Sprintf("Receive From Address: %v\n", swapRequest.ReceiveFromAddress),
		fmt.Sprintf("Send Value: %v\n", swapRequest.SendValue),
		fmt.Sprintf("Receive Value: %v\n", swapRequest.ReceiveValue),
		fmt.Sprintf("Send Token: %s\n", swapRequest.SendToken),
		fmt.Sprintf("Receive Token: %s\n", swapRequest.ReceiveToken),
		fmt.Sprintf("Goes First: %v\n", swapRequest.GoesFirst),
		fmt.Sprintf("-------------------------------------------------------------\n"),
	))
}
