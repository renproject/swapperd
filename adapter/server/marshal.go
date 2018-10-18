package server

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/republicprotocol/swapperd/foundation"
)

var ErrInvalidAmount = errors.New("invalid amount")
var ErrInvalidLength = errors.New("invalid length")

// GetPingResponse data object contains the Swapper's internal information.
type GetPingResponse struct {
	Version         string             `json:"version"`
	SupportedTokens []foundation.Token `json:"supportedTokens"`
}

type PostSwapRequestResponse struct {
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

type GetSwapResponse struct {
	Swaps []SwapStatus `json:"swaps"`
}

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
	id, err := toBytes32(swapIDBytes)
	return foundation.SwapID(id), err
}

func MarshalToken(token foundation.Token) string {
	return token.String()
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
	return toBytes32(hashBytes)
}

func UnmarshalSwapRequestResponse(swapReqRes PostSwapRequestResponse) (foundation.Swap, error) {
	secret := [32]byte{}
	if swapReqRes.ShouldInitiateFirst {
		rand.Read(secret[:])
		hash := sha256.Sum256(secret[:])
		swapReqRes.SecretHash = MarshalSecretHash(hash)
	}
	swapID, err := UnmarshalSwapID(swapReqRes.ID)
	if err != nil {
		return foundation.Swap{}, nil
	}
	sendToken, err := UnmarshalToken(swapReqRes.SendToken)
	if err != nil {
		return foundation.Swap{}, nil
	}
	receiveToken, err := UnmarshalToken(swapReqRes.ReceiveToken)
	if err != nil {
		return foundation.Swap{}, nil
	}
	sendValue, err := UnmarshalAmount(swapReqRes.SendAmount)
	if err != nil {
		return foundation.Swap{}, nil
	}
	receiveValue, err := UnmarshalAmount(swapReqRes.ReceiveAmount)
	if err != nil {
		return foundation.Swap{}, nil
	}
	secretHash, err := UnmarshalSecretHash(swapReqRes.SecretHash)
	if err != nil {
		return foundation.Swap{}, nil
	}
	return foundation.Swap{
		ID:                 swapID,
		Secret:             secret,
		SecretHash:         secretHash,
		TimeLock:           swapReqRes.TimeLock,
		SendToAddress:      swapReqRes.SendTo,
		ReceiveFromAddress: swapReqRes.ReceiveFrom,
		SendValue:          sendValue,
		ReceiveValue:       receiveValue,
		SendToken:          sendToken,
		ReceiveToken:       receiveToken,
		IsFirst:            swapReqRes.ShouldInitiateFirst,
	}, nil
}

func toBytes32(data []byte) ([32]byte, error) {
	bytes32 := [32]byte{}
	if len(data) != 32 {
		return bytes32, ErrInvalidLength
	}
	copy(bytes32[:], data)
	return bytes32, nil
}
