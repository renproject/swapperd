package state

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/republicprotocol/renex-swapper-go/domain/match"
	"github.com/republicprotocol/renex-swapper-go/domain/swap"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
)

type SwapDetails struct {
	TimeStamp       int64       `json:"timeStamp"`
	Status          swap.Status `json:"status"`
	Expiry          int64       `json:"expiry"`
	HashLock        [32]byte    `json:"hashLock"`
	Secret          [32]byte    `json:"secret"`
	ForeignOrderID  [32]byte    `json:"foreignOrderID"`
	SendValue       *big.Int    `json:"sendValue"`
	ReceiveValue    *big.Int    `json:"receiveValue"`
	SendCurrency    token.Token `json:"sendCurrency"`
	ReceiveCurrency token.Token `json:"receiveCurrency"`
	PersonalAtom    []byte      `json:"personalAtom"`
	ForeignAtom     []byte      `json:"foreignAtom"`
}

type state struct {
	swapMu    *sync.RWMutex
	swapCache map[[32]byte]SwapDetails
	swapList  [][32]byte
	Adapter
}

type State interface {
	ActiveSwapList
	InitiateDetails([32]byte) ([32]byte, int64, error)
	PutInitiateDetails([32]byte, [32]byte, int64) error

	RedeemDetails([32]byte) ([32]byte, error)
	PutRedeemDetails([32]byte, [32]byte) error

	Status([32]byte) swap.Status
	PutStatus([32]byte, swap.Status) error

	Match([32]byte) (match.Match, error)
	PutMatch([32]byte, match.Match) error

	SwapTimestamp([32]byte) (int64, error)
	PutSwapTimestamp([32]byte) error

	PersonalAtom([32]byte) ([]byte, error)
	PutPersonalAtom([32]byte, []byte) error

	ForeignAtom([32]byte) ([]byte, error)
	PutForeignAtom([32]byte, []byte) error

	AtomsExist([32]byte) bool
}

// NewState creates a new state interface
func NewState(adapter Adapter) State {
	return &state{
		Adapter:   adapter,
		swapCache: map[[32]byte]SwapDetails{},
		swapMu:    new(sync.RWMutex),
	}
}

// PutInitiateDetails into both persistent storage and in memory cache.
func (state *state) PutInitiateDetails(orderID [32]byte, hashLock [32]byte, expiry int64) error {
	swap := state.swapCache[orderID]
	swap.HashLock = hashLock
	swap.Expiry = expiry
	if err := state.WriteSwapDetails(orderID, swap); err != nil {
		return err
	}
	state.swapCache[orderID] = swap
	return nil
}

// InitiateDetails tries to get the initiate details from in memory cache if
// they do not exist it tries to read from the persistent storage.
func (state *state) InitiateDetails(orderID [32]byte) ([32]byte, int64, error) {
	swap := state.swapCache[orderID]
	if swap.HashLock != [32]byte{} && swap.Expiry != 0 {
		return swap.HashLock, swap.Expiry, nil
	}
	swap, err := state.ReadSwapDetails(orderID)
	if err != nil {
		return [32]byte{}, 0, err
	}
	state.swapCache[orderID] = swap
	return swap.HashLock, swap.Expiry, nil
}

// PutRedeemDetails into both persistent storage and in memory cache.
func (state *state) PutRedeemDetails(orderID [32]byte, secret [32]byte) error {
	swap := state.swapCache[orderID]
	swap.Secret = secret
	if err := state.WriteSwapDetails(orderID, swap); err != nil {
		return err
	}
	state.swapCache[orderID] = swap
	return nil
}

// RedeemDetails tries to get the redeem details from in memory cache if
// they do not exist it tries to read from the persistent storage.
func (state *state) RedeemDetails(orderID [32]byte) ([32]byte, error) {
	swap := state.swapCache[orderID]
	if swap.Secret != [32]byte{} {
		return swap.Secret, nil
	}
	swap, err := state.ReadSwapDetails(orderID)
	if err != nil {
		return [32]byte{}, err
	}
	state.swapCache[orderID] = swap
	return swap.Secret, nil
}

// PutStatus into both persistent storage and in memory cache.
func (state *state) PutStatus(orderID [32]byte, status swap.Status) error {
	swap := state.swapCache[orderID]
	swap.Status = status
	if err := state.WriteSwapDetails(orderID, swap); err != nil {
		return err
	}
	state.swapCache[orderID] = swap
	return nil
}

// Status tries to get the swap status from in memory cache if it does not exist
// it tries to read from the persistent storage.
func (state *state) Status(orderID [32]byte) swap.Status {
	swapDetails := state.swapCache[orderID]
	if swapDetails.Status != swap.Status("") {
		return swapDetails.Status
	}
	swapDetails, err := state.ReadSwapDetails(orderID)
	if err != nil {
		return swap.StatusUnknown
	}
	state.swapCache[orderID] = swapDetails
	return swapDetails.Status
}

// PutMatch into both persistent storage and in memory cache.
func (state *state) PutMatch(orderID [32]byte, match match.Match) error {
	swap := state.swapCache[orderID]
	swap.ForeignOrderID = match.ForeignOrderID()
	swap.SendValue = match.SendValue()
	swap.ReceiveValue = match.ReceiveValue()
	swap.SendCurrency = match.SendCurrency()
	swap.ReceiveCurrency = match.ReceiveCurrency()
	if err := state.WriteSwapDetails(orderID, swap); err != nil {
		return err
	}
	state.swapCache[orderID] = swap
	return nil
}

// Match tries to get the match details from in memory cache if they do not
// exist it tries to read from the persistent storage.
func (state *state) Match(orderID [32]byte) (match.Match, error) {
	swap := state.swapCache[orderID]
	if swap.ForeignOrderID != [32]byte{} {
		return match.NewMatch(orderID, swap.ForeignOrderID, swap.SendValue, swap.ReceiveValue, swap.SendCurrency, swap.ReceiveCurrency), nil
	}
	swap, err := state.ReadSwapDetails(orderID)
	if err != nil {
		return nil, err
	}
	state.swapCache[orderID] = swap
	return match.NewMatch(orderID, swap.ForeignOrderID, swap.SendValue, swap.ReceiveValue, swap.SendCurrency, swap.ReceiveCurrency), nil
}

// PutPersonalAtom into both persistent storage and in memory cache.
func (state *state) PutPersonalAtom(orderID [32]byte, atomDetails []byte) error {
	swap := state.swapCache[orderID]
	swap.PersonalAtom = atomDetails
	if err := state.WriteSwapDetails(orderID, swap); err != nil {
		return err
	}
	state.swapCache[orderID] = swap
	return nil
}

// PersonalAtom tries to get the personal atom details from in memory cache if
// they do not exist it tries to read from the persistent storage.
func (state *state) PersonalAtom(orderID [32]byte) ([]byte, error) {
	swap := state.swapCache[orderID]
	if bytes.Compare(swap.PersonalAtom, []byte{}) != 0 {
		return swap.PersonalAtom, nil
	}
	swap, err := state.ReadSwapDetails(orderID)
	if err != nil {
		return []byte{}, err
	}
	state.swapCache[orderID] = swap
	return swap.PersonalAtom, nil
}

// PutForeignAtom into both persistent storage and in memory cache.
func (state *state) PutForeignAtom(orderID [32]byte, atomDetails []byte) error {
	swap := state.swapCache[orderID]
	swap.ForeignAtom = atomDetails
	if err := state.WriteSwapDetails(orderID, swap); err != nil {
		return err
	}
	state.swapCache[orderID] = swap
	return nil
}

// ForeignAtom tries to get the foreign atom details from in memory cache if
// they do not exist it tries to read from the persistent storage.
func (state *state) ForeignAtom(orderID [32]byte) ([]byte, error) {
	swap := state.swapCache[orderID]
	if bytes.Compare(swap.ForeignAtom, []byte{}) != 0 {
		return swap.ForeignAtom, nil
	}
	swap, err := state.ReadSwapDetails(orderID)
	if err != nil {
		return []byte{}, err
	}
	state.swapCache[orderID] = swap
	return swap.ForeignAtom, nil
}

// AtomsExist checks whether the atoms are created beforehand.
func (state *state) AtomsExist(orderID [32]byte) bool {
	swapAtomDetails, err := state.ReadSwapDetails(orderID)
	if err != nil || (bytes.Compare(swapAtomDetails.PersonalAtom, []byte{}) == 0 && bytes.Compare(swapAtomDetails.ForeignAtom, []byte{}) == 0) {
		return false
	}
	return true
}

// PutSwapTimestamp into both persistent storage and in memory cache.
func (state *state) PutSwapTimestamp(orderID [32]byte) error {
	swap := state.swapCache[orderID]
	swap.TimeStamp = time.Now().Unix()
	if err := state.WriteSwapDetails(orderID, swap); err != nil {
		return err
	}
	state.swapCache[orderID] = swap
	return nil
}

// SwapTimestamp tries to get the redeem details from in memory cache if
// they do not exist it tries to read from the persistent storage.
func (state *state) SwapTimestamp(orderID [32]byte) (int64, error) {
	swap := state.swapCache[orderID]
	if swap.TimeStamp != 0 {
		return swap.TimeStamp, nil
	}
	swap, err := state.ReadSwapDetails(orderID)
	if err != nil {
		return 0, err
	}
	state.swapCache[orderID] = swap
	return swap.TimeStamp, nil
}

// WriteSwapDetails to persistent storage
func (state *state) WriteSwapDetails(orderID [32]byte, swapDetails SwapDetails) error {
	state.PrintSwapDetails(swapDetails)
	data, err := json.Marshal(swapDetails)
	if err != nil {
		return err
	}
	return state.Write(append([]byte("Swap Details:"), orderID[:]...), data)
}

// ReadSwapDetails from persistent storage
func (state *state) ReadSwapDetails(orderID [32]byte) (SwapDetails, error) {
	data, err := state.Read(append([]byte("Swap Details:"), orderID[:]...))
	if err != nil {
		return SwapDetails{}, err
	}
	swapDetails := SwapDetails{}
	if err = json.Unmarshal(data, &swapDetails); err != nil {
		return SwapDetails{}, err
	}
	return swapDetails, nil
}

// PrintSwapDetails to Std Out
func (state *state) PrintSwapDetails(swapDetails SwapDetails) {
	fmt.Printf("\n\n\n\t\tSTATE UPDATED\n\n")
	fmt.Printf("Timestamp: %d\n", swapDetails.TimeStamp)
	fmt.Printf("Status: %s\n", swapDetails.Status)
	fmt.Printf("Expiry: %d\n", swapDetails.Expiry)
	fmt.Printf("Hash Lock: %s\n", base64.StdEncoding.EncodeToString(swapDetails.HashLock[:]))
	fmt.Printf("Secret: %s\n", base64.StdEncoding.EncodeToString(swapDetails.Secret[:]))
	fmt.Printf("Foreign Order ID: %s\n", base64.StdEncoding.EncodeToString(swapDetails.ForeignOrderID[:]))
	fmt.Printf("Send Value: %v\n", swapDetails.SendValue)
	fmt.Printf("Receive Value: %v\n", swapDetails.ReceiveValue)
	fmt.Printf("Send Currency: %d\n", swapDetails.SendCurrency)
	fmt.Printf("Receive Currency: %d\n", swapDetails.ReceiveCurrency)
	fmt.Printf("Personal Atom: %v\n", hex.EncodeToString(swapDetails.PersonalAtom[:]))
	fmt.Printf("Foreign Atom: %v\n", hex.EncodeToString(swapDetails.ForeignAtom[:]))
	fmt.Println("-------------------------------------------------------------")
}
