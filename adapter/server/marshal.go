package server

import (
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/republicprotocol/swapperd/foundation"
)

// GetPingResponse data object contains the Swapper's internal information.
type GetPingResponse struct {
	Version         string             `json:"version"`
	SupportedTokens []foundation.Token `json:"supportedTokens"`
}

type PostSwapMessage struct {
	ID                  string `json:"id"`
	SendToken           string `json:"sendToken"`
	ReceiveToken        string `json:"receiveToken"`
	SendAmount          string `json:"sendAmount"`    // hex
	ReceiveAmount       string `json:"receiveAmount"` //hex
	SendTo              string `json:"sendTo"`
	ReceiveFrom         string `json:"receiveFrom"`
	TimeLock            int64  `json:"timeLock"`
	SecretHash          string `json:"secretHash"`
	ShouldInitiateFirst bool   `json:"shouldInitiateFirst"`
}

// GetSwapResponse
type GetSwapResponse struct {
	Swaps []SwapStatus `json:"swaps"`
}

// SwapStatus
type SwapStatus struct {
}

type GetBalanceResponse struct {
	Balances []Balance `json:"balances"`
}

type Balance struct {
	TokenName string `json:"name"`
	Address   string `json:"address"`
	Amount    string `json:"amount"`
}

func MarshalSwapID(swapID foundation.SwapID) string {
	return base64.StdEncoding.EncodeToString(swapID[:])
}

func UnmarshalSwapID(swapID string) (foundation.SwapID, error) {
	swapIDBytes, err := base64.StdEncoding.DecodeString(swapID)
	if err != nil {
		return foundation.SwapID([32]byte{}), err
	}
	id, err := ToBytes32(swapIDBytes)
	return foundation.SwapID(id), err
}

func MarshalToken(token foundation.Token) string {
	return token.Name
}

func UnmarshalToken(token string) (foundation.Token, error) {
	token = strings.ToLower(token)
	switch token {
	case "btc", "bitcoin", "xbt":
		return foundation.TokenBTC, nil
	case "eth", "ethereum", "ether":
		return foundation.TokenETH, nil
	case "wbtc", "wrappedbtc", "wrappedbitcoin":
		return foundation.TokenWBTC, nil
	default:
		return foundation.Token{}, foundation.NewErrUnsupportedToken(token)
	}
}

func MarshalAmount(amount *big.Int) string {
	return hex.EncodeToString(amount.Bytes())
}

func UnmarshalAmount(amount string) (*big.Int, error) {
	val, success := big.NewInt(0).SetString(amount, 16)
	if !success {
		return nil, ErrInvalidAmount
	}
	return val, nil
}

func MarshalSecretHash(hash [32]byte) string {
	return base64.StdEncoding.EncodeToString(hash[:])
}

func UnmarshalSecretHash(hash string) ([32]byte, error) {
	hashBytes, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return [32]byte{}, err
	}
	return ToBytes32(hashBytes)
}

func ToBytes32(data []byte) ([32]byte, error) {
	bytes32 := [32]byte{}
	if len(data) != 32 {
		return bytes32, ErrInvalidLength
	}
	copy(bytes32[:], data)
	return bytes32, nil
}
