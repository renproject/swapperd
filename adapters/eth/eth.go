package eth

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/atom-go/services/arc"
)

type EthereumArc struct {
	context context.Context
	client  Connection
	auth    *bind.TransactOpts
	binding *Arc
	swapID  [32]byte
}

// NewEthereumArc returns a new EthereumArc instance
func NewEthereumArc(context context.Context, client Connection, auth *bind.TransactOpts, swapID [32]byte) (arc.Arc, error) {
	contract, err := NewArc(client.EthAddress, bind.ContractBackend(client.Client))
	if err != nil {
		return &EthereumArc{}, err
	}

	return &EthereumArc{
		context: context,
		client:  client,
		auth:    auth,
		binding: contract,
		swapID:  swapID,
	}, nil
}

// Initiate a new Arc swap by calling a function on ethereum
func (ethArc *EthereumArc) Initiate(hash [32]byte, from []byte, to []byte, value *big.Int, expiry int64) error {
	ethArc.auth.Value = value
	tx, err := ethArc.binding.Initiate(ethArc.auth, ethArc.swapID, common.BytesToAddress(to), hash, big.NewInt(expiry))
	ethArc.auth.Value = big.NewInt(0)
	if err != nil {
		return err
	}
	_, err = ethArc.client.PatchedWaitMined(ethArc.context, tx)
	return err
}

// Redeem an Arc swap by calling a function on ethereum
func (ethArc *EthereumArc) Redeem(secret [32]byte) error {
	tx, err := ethArc.binding.Redeem(ethArc.auth, ethArc.swapID, secret)
	if err == nil {
		_, err = ethArc.client.PatchedWaitMined(ethArc.context, tx)
	}
	return err
}

// Refund an Arc swap by calling a function on ethereum
func (ethArc *EthereumArc) Refund() error {
	tx, err := ethArc.binding.Refund(ethArc.auth, ethArc.swapID)
	if err == nil {
		_, err = ethArc.client.PatchedWaitMined(ethArc.context, tx)
	}
	return err
}

// Audit an Arc swap by calling a function on ethereum
func (ethArc *EthereumArc) Audit() (hash [32]byte, to, from []byte, value *big.Int, expiry int64, err error) {
	auditReport, err := ethArc.binding.Audit(&bind.CallOpts{}, ethArc.swapID)
	if err != nil {
		return [32]byte{}, nil, nil, nil, 0, err
	}
	return auditReport.SecretLock, auditReport.From.Bytes(), auditReport.To.Bytes(), auditReport.Value, auditReport.Timelock.Int64(), err
}

// AuditSecret audits the secret of an Arc swap by calling a function on ethereum
func (ethArc *EthereumArc) AuditSecret() ([32]byte, error) {
	return ethArc.binding.AuditSecret(&bind.CallOpts{}, ethArc.swapID)
}

// Serialize serializes an Ethereum Arc object
func (ethArc *EthereumArc) Serialize() ([]byte, error) {
	return ethArc.swapID[:], nil
}

// Deserialize deserializes an Ethereum Arc object
func (ethArc *EthereumArc) Deserialize(b []byte) error {
	bytes32 := [32]byte{}
	if len(b) != 32 {
		return errors.New("Deserialization failed due to malformed input")
	}
	for i := range b {
		bytes32[i] = b[i]
	}
	ethArc.swapID = bytes32
	return nil
}
