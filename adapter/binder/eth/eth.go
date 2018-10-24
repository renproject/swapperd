package eth

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/republicprotocol/beth-go"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/swapperd/core/swapper"
	"github.com/republicprotocol/swapperd/foundation"
)

type ethSwapContractBinder struct {
	id      [32]byte
	account beth.Account
	swap    foundation.Swap
	logger  swapper.Logger
	binder  *RenExAtomicSwapper
}

// NewETHSwapContractBinder returns a new Ethereum RequestAtom instance
func NewETHSwapContractBinder(account beth.Account, swap foundation.Swap, logger swapper.Logger) (swapper.Contract, error) {
	client := account.EthClient()
	swapperAddr, err := account.ReadAddress(fmt.Sprintf("SWAPPER:%s", swap.Token.Name))
	if err != nil {
		return nil, err
	}

	contract, err := NewRenExAtomicSwapper(swapperAddr, bind.ContractBackend(client.EthClient()))
	if err != nil {
		return nil, err
	}

	id, err := contract.SwapID(&bind.CallOpts{}, swapperAddr, swap.SecretHash, big.NewInt(swap.TimeLock))
	if err != nil {
		return nil, err
	}

	logger.LogInfo(swap.ID, fmt.Sprintf("Ethereum Atomic Swap ID: %s", base64.StdEncoding.EncodeToString(id[:])))

	return &ethSwapContractBinder{
		account: account,
		binder:  contract,
		logger:  logger,
		swap:    swap,
		id:      id,
	}, nil
}

// Initiate a new Atom swap by calling a function on ethereum
func (atom *ethSwapContractBinder) Initiate() error {
	atom.logger.LogInfo(atom.swap.ID, "Initiating on Ethereum blockchain")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Initiate the Atomic Swap
	if err := atom.account.Transact(
		ctx,
		func() bool {
			initiatable, err := atom.binder.Initiatable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return initiatable
		},
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tops.Value = atom.swap.Value
			tx, err := atom.binder.Initiate(tops, atom.id, common.HexToAddress(atom.swap.SpendingAddress), atom.swap.SecretHash, big.NewInt(atom.swap.TimeLock))
			if err != nil {
				return tx, err
			}
			tops.Value = big.NewInt(0)
			msg, _ := atom.account.FormatTransactionView("Initiated the atomic swap on Ethereum blockchain", tx.Hash().String())
			atom.logger.LogInfo(atom.swap.ID, msg)
			return tx, nil
		},
		func() bool {
			initiatable, err := atom.binder.Initiatable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return !initiatable
		},
		1,
	); err != nil && err != beth.ErrPreConditionCheckFailed {
		return err
	}
	return nil
}

// Refund an Atom swap by calling a function on ethereum
func (atom *ethSwapContractBinder) Refund() error {
	atom.logger.LogInfo(atom.swap.ID, "Refunding the atomic swap on Ethereum blockchain")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := atom.account.Transact(
		ctx,
		func() bool {
			refundable, err := atom.binder.Refundable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return refundable
		},
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := atom.binder.Refund(tops, atom.id)
			if err != nil {
				return nil, err
			}
			msg, _ := atom.account.FormatTransactionView("Refunded the atomic swap on Ethereum blockchain", tx.Hash().String())
			atom.logger.LogInfo(atom.swap.ID, msg)
			return tx, nil
		},
		func() bool {
			refundable, err := atom.binder.Refundable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return !refundable
		},
		1,
	); err != nil && err != beth.ErrPreConditionCheckFailed {
		return err
	}
	return nil
}

// AuditSecret audits the secret of an Atom swap by calling a function on ethereum
func (atom *ethSwapContractBinder) AuditSecret() ([32]byte, error) {
	for {
		atom.logger.LogInfo(atom.swap.ID, "Auditing secret on ethereum blockchain")
		redeemable, err := atom.binder.Redeemable(&bind.CallOpts{}, atom.id)
		if err != nil {
			atom.logger.LogError(atom.swap.ID, err)
			return [32]byte{}, err
		}
		if !redeemable {
			break
		}
		if time.Now().Unix() > atom.swap.TimeLock {
			return [32]byte{}, fmt.Errorf("Timed Out")
		}
		time.Sleep(15 * time.Second)
	}
	secret, err := atom.binder.AuditSecret(&bind.CallOpts{}, atom.id)
	if err != nil {
		return [32]byte{}, err
	}
	atom.logger.LogInfo(atom.swap.ID, fmt.Sprintf("Audit success on ethereum blockchain secret=%s", base64.StdEncoding.EncodeToString(secret[:])))
	return secret, nil
}

// Audit an Atom swap by calling a function on ethereum
func (atom *ethSwapContractBinder) Audit() error {
	atom.logger.LogInfo(atom.swap.ID, fmt.Sprintf("Waiting for initiation on ethereum blockchain"))
	for {
		initiatable, err := atom.binder.Initiatable(&bind.CallOpts{}, atom.id)
		if err != nil {
			return err
		}
		if !initiatable {
			break
		}
		time.Sleep(15 * time.Second)
	}
	auditReport, err := atom.binder.Audit(&bind.CallOpts{}, atom.id)
	if err != nil {
		return err
	}
	if auditReport.Value.Cmp(atom.swap.Value) != 0 {
		return fmt.Errorf("Receive Value Mismatch Expected: %v Actual: %v", atom.swap.Value, auditReport.Value)
	}
	atom.logger.LogInfo(atom.swap.ID, fmt.Sprintf("Audit successful on Ethereum blockchain"))
	return nil
}

// Redeem an Atom swap by calling a function on ethereum
func (atom *ethSwapContractBinder) Redeem(secret [32]byte) error {
	atom.logger.LogInfo(atom.swap.ID, "Redeeming the atomic swap on Ethereum blockchain")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := atom.account.Transact(
		ctx,
		func() bool {
			redeemable, err := atom.binder.Redeemable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return redeemable
		},
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := atom.binder.Redeem(tops, atom.id, secret)
			if err != nil {
				return nil, err
			}
			msg, _ := atom.account.FormatTransactionView("Redeemed the atomic swap on ERC20 blockchain", tx.Hash().String())
			atom.logger.LogInfo(atom.swap.ID, msg)
			return tx, nil
		},
		func() bool {
			refundable, err := atom.binder.Refundable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return !refundable
		},
		1,
	); err != nil && err != beth.ErrPreConditionCheckFailed {
		return err
	}
	return nil
}
