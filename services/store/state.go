package store

import (
	"github.com/republicprotocol/atom-go/domains/match"
	"github.com/republicprotocol/atom-go/utils"
)

type State interface {
	UpdateStatus([32]byte, string) error
	ReadStatus([32]byte) string

	Read([]byte) ([]byte, error)
	Write([]byte, []byte) error

	SetMatch([32]byte, match.Match) error
	GetMatch([32]byte) (match.Match, error)
}

type state struct {
	Store
}

func NewSwapStore(store Store) State {
	return &state{
		store,
	}
}

func (str *state) UpdateStatus(orderID [32]byte, status string) error {
	return str.Write(append([]byte("status:"), orderID[:]...), []byte(status))
}

func (str *state) ReadStatus(orderID [32]byte) string {
	status, err := str.Read(append([]byte("status:"), orderID[:]...))
	if err != nil {
		return "UNKNOWN"
	}
	return string(status)
}

func (str *state) SetMatch(orderID [32]byte, m match.Match) error {
	data, err := m.Serialize()
	if err != nil {
		return err
	}
	return str.Write(append([]byte("match:"), orderID[:]...), data)
}

func (str *state) GetMatch(orderID [32]byte) (match.Match, error) {
	data, err := str.Read(append([]byte("match:"), orderID[:]...))
	if err != nil {
		return nil, err
	}
	return match.NewMatchFromBytes(data)
}

func (str *state) SetSecret(orderID [32]byte, sec [32]byte) error {
	return str.Write(append([]byte("secret:"), orderID[:]...), sec[:])
}

func (str *state) GetSecret(orderID [32]byte) ([32]byte, error) {
	secret, err := str.Read(append([]byte("secret:"), orderID[:]...))
	if err != nil {
		return [32]byte{}, err
	}
	return utils.ToBytes32(secret)
}
