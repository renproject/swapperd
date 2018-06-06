package eth

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/atom-go/services/atom"
)

type EthereumAtom struct {
	context context.Context
	client  Connection
	auth    *bind.TransactOpts
	binding *Atom
	swapID  [32]byte
}

// NewEthereumAtom returns a new EthereumAtom instance
func NewEthereumAtom(context context.Context, client Connection, auth *bind.TransactOpts, swapID [32]byte) (atom.Atom, error) {
	contract, err := NewAtom(client.EthAddress, bind.ContractBackend(client.Client))
	if err != nil {
		return &EthereumAtom{}, err
	}

	return &EthereumAtom{
		context: context,
		client:  client,
		auth:    auth,
		binding: contract,
		swapID:  swapID,
	}, nil
}

// Initiate a new Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) Initiate(hash [32]byte, from []byte, to []byte, value *big.Int, expiry int64) error {
	ethAtom.auth.Value = value
	tx, err := ethAtom.binding.Initiate(ethAtom.auth, ethAtom.swapID, common.BytesToAddress(to), hash, big.NewInt(expiry))
	ethAtom.auth.Value = big.NewInt(0)
	if err != nil {
		return err
	}
	_, err = ethAtom.client.PatchedWaitMined(ethAtom.context, tx)
	return err
}

// Redeem an Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) Redeem(secret [32]byte) error {
	tx, err := ethAtom.binding.Redeem(ethAtom.auth, ethAtom.swapID, secret)
	if err == nil {
		_, err = ethAtom.client.PatchedWaitMined(ethAtom.context, tx)
	}
	return err
}

// Refund an Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) Refund() error {
	tx, err := ethAtom.binding.Refund(ethAtom.auth, ethAtom.swapID)
	if err == nil {
		_, err = ethAtom.client.PatchedWaitMined(ethAtom.context, tx)
	}
	return err
}

// Audit an Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) Audit() (hash [32]byte, to, from []byte, value *big.Int, expiry int64, err error) {
	auditReport, err := ethAtom.binding.Audit(&bind.CallOpts{}, ethAtom.swapID)
	if err != nil {
		return [32]byte{}, nil, nil, nil, 0, err
	}
	return auditReport.SecretLock, auditReport.From.Bytes(), auditReport.To.Bytes(), auditReport.Value, auditReport.Timelock.Int64(), err
}

// AuditSecret audits the secret of an Atom swap by calling a function on ethereum
func (ethAtom *EthereumAtom) AuditSecret() ([32]byte, error) {
	return ethAtom.binding.AuditSecret(&bind.CallOpts{}, ethAtom.swapID)
}

// Store stores the swap details on Ethereum
func (ethAtom *EthereumAtom) Store(orderID [32]byte) error {
	tx, err := ethAtom.binding.SubmitDetails(ethAtom.auth, orderID, ethAtom.swapID[:])
	if err != nil {
		return err
	}
	ethAtom.client.PatchedWaitMined(ethAtom.context, tx)
	return nil
}

// Retrieve retrieves the swap details from Ethereum
func (ethAtom *EthereumAtom) Retrieve(orderID [32]byte) error {
	b, err := ethAtom.binding.SwapDetails(&bind.CallOpts{}, orderID)
	if err != nil {
		return err
	}
	bytes32 := [32]byte{}
	if len(b) != 32 {
		return errors.New("Deserialization failed due to malformed input")
	}
	for i := range b {
		bytes32[i] = b[i]
	}
	ethAtom.swapID = bytes32
	return nil
}
