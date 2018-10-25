package server

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/republicprotocol/swapperd/foundation"
)

var ErrInvalidAmount = errors.New("invalid amount")
var ErrInvalidLength = errors.New("invalid length")

// GetWhoAmIResponse data object contains the Swapper's internal information.
type GetWhoAmIResponse struct {
	Version         string             `json:"version"`
	SupportedTokens []foundation.Token `json:"supportedTokens"`
}

type GetSwapsResponse struct {
	Swaps []SwapStatus `json:"swaps"`
}

type SwapStatus struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
}

type GetBalanceResponse struct {
	Balances []Balance `json:"balances"`
}

type Balance struct {
	Token   string `json:"token"`
	Address string `json:"address"`
	Amount  string `json:"amount"`
}

type PostWithdrawalsRequest struct {
	To     string `json:"to"`
	Token  string `json:"token"`
	Amount string `json:"amount"`
}

type PostWithdrawalsResponse struct {
	TxHash string `txHash`
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

func toBytes32(data []byte) ([32]byte, error) {
	bytes32 := [32]byte{}
	if len(data) != 32 {
		return bytes32, ErrInvalidLength
	}
	copy(bytes32[:], data)
	return bytes32, nil
}
