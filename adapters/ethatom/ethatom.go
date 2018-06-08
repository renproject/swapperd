package ethatom

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/atom-go/adapters/ethclient"
	"github.com/republicprotocol/atom-go/services/atom"
)

// EthereumData
type EthereumData struct {
	SwapID [32]byte `json:"swap_id"`
	Owner  []byte   `json:"owner"`
}

type EthereumAtom struct {
	context      context.Context
	client       ethclient.Connection
	auth         *bind.TransactOpts
	binding      *Atom
	personalData EthereumData
	foreignData  EthereumData
}

// NewEthereumAtom returns a new EthereumAtom instance
func NewEthereumAtom(context context.Context, client ethclient.Connection, auth *bind.TransactOpts, swapID [32]byte) (atom.Atom, error) {
	contract, err := NewAtom(client.EthAddress, bind.ContractBackend(client.Client))
	if err != nil {
		return &EthereumAtom{}, err
	}

	return &EthereumAtom{
		context: context,
		client:  client,
		auth:    auth,
		binding: contract,
		personalData: EthereumData{
			SwapID: swapID,
			Owner:  auth.From.Bytes(),
		},
	}, nil
}

// Initiate a new Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) Initiate(hash [32]byte, value *big.Int, expiry int64) error {
	if bytes.Compare(ethAtom.foreignData.Owner, []byte{}) == 0 {
		return errors.New("Please run deserialize before running initiate")
	}
	ethAtom.auth.Value = value
	tx, err := ethAtom.binding.Initiate(ethAtom.auth, ethAtom.personalData.SwapID, common.BytesToAddress(ethAtom.foreignData.Owner), hash, big.NewInt(expiry))
	ethAtom.auth.Value = big.NewInt(0)
	if err != nil {
		return err
	}
	_, err = ethAtom.client.PatchedWaitMined(ethAtom.context, tx)
	return err
}

// Redeem an Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) Redeem(secret [32]byte) error {
	tx, err := ethAtom.binding.Redeem(ethAtom.auth, ethAtom.foreignData.SwapID, secret)
	if err == nil {
		_, err = ethAtom.client.PatchedWaitMined(ethAtom.context, tx)
	}
	return err
}

// Refund an Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) Refund() error {
	tx, err := ethAtom.binding.Refund(ethAtom.auth, ethAtom.personalData.SwapID)
	if err == nil {
		_, err = ethAtom.client.PatchedWaitMined(ethAtom.context, tx)
	}
	return err
}

// Audit an Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) Audit() (hash [32]byte, to, from []byte, value *big.Int, expiry int64, err error) {
	auditReport, err := ethAtom.binding.Audit(&bind.CallOpts{}, ethAtom.foreignData.SwapID)
	if err != nil {
		return [32]byte{}, nil, nil, nil, 0, err
	}
	return auditReport.SecretLock, auditReport.From.Bytes(), auditReport.To.Bytes(), auditReport.Value, auditReport.Timelock.Int64(), err
}

// AuditSecret audits the secret of an Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) AuditSecret() ([32]byte, error) {
	return ethAtom.binding.AuditSecret(&bind.CallOpts{}, ethAtom.personalData.SwapID)
}

// Serialize serializes the atom details into a bytes array
func (ethAtom *EthereumAtom) Serialize() ([]byte, error) {
	b, err := json.Marshal(ethAtom.personalData)
	return b, err
}

// Deserialize deserializes the atom details from a bytes array
func (ethAtom *EthereumAtom) Deserialize(b []byte) error {
	return json.Unmarshal(b, &ethAtom.foreignData)
}

// From returns the address of the sender
func (ethAtom *EthereumAtom) From() []byte {
	return ethAtom.auth.From.Bytes()
}

// PriorityCode returns the priority code of the currency.
func (ethAtom *EthereumAtom) PriorityCode() int64 {
	return 1
}
