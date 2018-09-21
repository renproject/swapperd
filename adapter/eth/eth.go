package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethclient "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/clients/eth"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type ethereum struct {
	id      [32]byte
	context context.Context
	client  ethclient.Conn
	key     keystore.EthereumKey
	req     swap.Request
	binder  *RenExAtomicSwapper
}

// NewEthereumAtom returns a new Ethereum RequestAtom instance
func NewEthereumAtom(conf config.EthereumNetwork, key keystore.EthereumKey, req swap.Request) (swap.Atom, error) {
	conn, err := ethclient.NewConnWithConfig(conf)
	if err != nil {
		return &ethereum{}, err
	}

	contract, err := NewRenExAtomicSwapper(conn.RenExAtomicSwapper, bind.ContractBackend(conn.Client))
	if err != nil {
		return &ethereum{}, err
	}

	addr, expiry, err := buildValues(req)
	if err != nil {
		return &ethereum{}, err
	}

	id, err := contract.SwapID(&bind.CallOpts{}, common.HexToAddress(addr), req.HashLock, big.NewInt(expiry))
	if err != nil {
		return nil, err
	}

	return &ethereum{
		context: context.Background(),
		client:  conn,
		key:     key,
		binder:  contract,
		req:     req,
		id:      id,
	}, nil
}

// Initiate a new Atom swap by calling a function on ethereum
func (atom *ethereum) Initiate() error {
	initiatable, err := atom.binder.Initiatable(&bind.CallOpts{}, atom.id)
	if err != nil {
		return err
	}
	if !initiatable {
		return swap.ErrSwapAlreadyInitiated
	}

	prevValue := atom.key.TransactOpts.Value
	prevGasLimit := atom.key.TransactOpts.GasLimit
	atom.key.TransactOpts.Value = atom.req.SendValue
	atom.key.TransactOpts.GasLimit = 3000000
	tx, err := atom.binder.Initiate(atom.key.TransactOpts, atom.id, common.HexToAddress(atom.req.SendToAddress), atom.req.HashLock, big.NewInt(atom.req.TimeLock))
	atom.key.TransactOpts.Value = prevValue
	atom.key.TransactOpts.GasLimit = prevGasLimit

	if err != nil {
		return fmt.Errorf("Failed to initiate on the Ethereum blockchain: %v", err)
	}

	return nil
}

// Refund an Atom swap by calling a function on ethereum
func (atom *ethereum) Refund() error {
	refundable, err := atom.binder.Refundable(&bind.CallOpts{}, atom.id)
	if err != nil {
		return err
	}
	if !refundable {
		return swap.ErrNotRefundable
	}
	prevGasLimit := atom.key.TransactOpts.GasLimit
	atom.key.TransactOpts.GasLimit = 3000000
	tx, err := atom.binder.Refund(atom.key.TransactOpts, atom.id)
	atom.key.TransactOpts.GasLimit = prevGasLimit
	if err != nil {
		return err
	}
	return nil
}

// AuditSecret audits the secret of an Atom swap by calling a function on ethereum
func (atom *ethereum) AuditSecret() ([32]byte, error) {
	return atom.binder.AuditSecret(&bind.CallOpts{}, atom.id)
}

// Audit an Atom swap by calling a function on ethereum
func (atom *ethereum) Audit() error {
	auditReport, err := atom.binder.Audit(&bind.CallOpts{}, atom.id)
	if err != nil {
		return err
	}

	if auditReport.From.String() != atom.req.ReceiveFromAddress {
		return errors.New("From Address Mismatch")
	}

	if auditReport.SecretLock != atom.req.HashLock {
		return errors.New("Secret Hash Mismatch")
	}

	if auditReport.Timelock.Int64() != atom.req.TimeLock {
		return errors.New("Time Locks Mismatch")
	}

	if auditReport.To.String() != atom.key.Address.String() {
		return errors.New("To Address Mismatch")
	}

	if auditReport.Value.Cmp(atom.req.ReceiveValue) != 0 {
		return errors.New("Receive Value Mismatch")
	}

	return nil
}

// Redeem an Atom swap by calling a function on ethereum
func (atom *ethereum) Redeem(secret [32]byte) error {
	prevGasLimit := atom.key.TransactOpts.GasLimit
	atom.key.TransactOpts.GasLimit = 3000000
	tx, err := atom.binder.Redeem(atom.key.TransactOpts, atom.id, secret)
	atom.key.TransactOpts.GasLimit = prevGasLimit
	if err == nil {
		if _, err := atom.client.PatchedWaitMined(atom.context, tx); err != nil {
			return err
		}
		return nil
	}
	return err
}

func buildValues(req swap.Request) (string, int64, error) {
	var addr string
	var expiry int64
	if req.SendToken != token.ETH && req.ReceiveToken != token.ETH {
		return "", 0, errors.New("Expected one of the tokens to be ethereum")
	}

	if (req.GoesFirst && req.SendToken == token.ETH) || (!req.GoesFirst && req.ReceiveToken == token.ETH) {
		expiry = req.TimeLock
	} else {
		expiry = req.TimeLock - 24*60*60
	}

	if req.SendToken == token.ETH {
		addr = req.SendToAddress
	} else {
		addr = req.ReceiveFromAddress
	}

	return addr, expiry, nil
}
