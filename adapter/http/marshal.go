package http

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/republicprotocol/renex-swapper-go/adapter/config"
)

// WhoAmI data object contains the swapper's internal information.
type WhoAmI struct {
	Challenge           string   `json:"challenge"`
	Version             string   `json:"version"`
	AuthorizedAddresses []string `json:"authorizedAddresses"`
	SupportedCurrencies []string `json:"supportedCurrencies"`
}

// WhoAmISigned data object contains the WhoAmI object, and the validating
// signature of the atomic swapper.
type WhoAmISigned struct {
	WhoAmI    WhoAmI `json:"whoAmI"`
	Signature string `json:"signature"`
}

// Status data object contains an order ID's atomic swap status.
type Status struct {
	OrderID string `json:"orderID"`
	Status  string `json:"status"`
}

// PostOrder data obje
type PostOrder struct {
	OrderID   string `json:"orderID"`
	Signature string `json:"signature"`
}

type Balance struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
}

type Balances struct {
	Ethereum Balance `json:"ethereum"`
	Bitcoin  Balance `json:"bitcoin"`
}

func MarshalSignature(signatureIn []byte) string {
	return hex.EncodeToString(signatureIn)
}

func UnmarshalSignature(signatureIn string) ([]byte, error) {
	return hex.DecodeString(signatureIn)
}

func NewWhoAmI(challenge string, conf config.Config) WhoAmI {
	return WhoAmI{
		Challenge:           challenge,
		Version:             conf.Version,
		AuthorizedAddresses: conf.AuthorizedAddresses,
		SupportedCurrencies: conf.SupportedCurrencies,
	}
}

func MarshalWhoAmI(whoAmI WhoAmI) ([]byte, error) {
	return json.Marshal(whoAmI)
}

func MarshalOrderID(orderIDIn [32]byte) string {
	return hex.EncodeToString(orderIDIn[:])
}

func UnmarshalOrderID(orderIDIn string) ([32]byte, error) {
	orderID := [32]byte{}
	orderIDBytes, err := hex.DecodeString(orderIDIn)
	if err != nil {
		return orderID, fmt.Errorf("cannot decode order id %v: %v", orderIDIn, err)
	}
	if len(orderIDBytes) != 32 {
		return orderID, ErrInvalidOrderIDLength
	}
	copy(orderID[:], orderIDBytes)
	return orderID, nil
}
