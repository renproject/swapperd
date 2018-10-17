package erc20

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/swapperd/core"
	"github.com/republicprotocol/swapperd/foundation"
)

type erc20Atom struct {
	id          [32]byte
	client      Client
	req         foundation.Swap
	logger      core.Logger
	tokenBinder *CompatibleERC20
	binder      *RenExAtomicSwapper
}

type Client interface {
	GetSwapperAddress(foundation.Token) common.Address
	GetTokenAddress(foundation.Token) common.Address
	Conn() *ethclient.Client
	Transact(
		ctx context.Context,
		preConditionCheck func() bool,
		f func(*bind.TransactOpts) (*types.Transaction, error),
		postConditionCheck func() bool,
		waitForBlocks int64,
	) error
	Address() common.Address
	FormatTransactionView(string, string) string
}

// NewERC20Atom returns a new ERC20 Atom instance
func NewERC20Atom(client Client, logger core.Logger, req foundation.Swap) (core.SwapContractBinder, error) {
	token, expiry, err := getSwapDetails(req)
	if err != nil {
		return nil, err
	}

	tokenContract, err := NewCompatibleERC20(client.GetTokenAddress(token), bind.ContractBackend(client.Conn()))
	if err != nil {
		return nil, err
	}

	contract, err := NewRenExAtomicSwapper(client.GetSwapperAddress(token), bind.ContractBackend(client.Conn()))
	if err != nil {
		return nil, err
	}

	id, err := contract.SwapID(&bind.CallOpts{}, req.SecretHash, big.NewInt(expiry))
	if err != nil {
		return nil, err
	}

	logger.LogInfo(req.ID, fmt.Sprintf("ERC20 Atomic Swap ID: %s", base64.StdEncoding.EncodeToString(id[:])))
	req.TimeLock = expiry

	return &erc20Atom{
		client:      client,
		tokenBinder: tokenContract,
		binder:      contract,
		logger:      logger,
		req:         req,
		id:          id,
	}, nil
}

// Initiate a new Atom swap by calling a function on ethereum
func (atom *erc20Atom) Initiate() error {
	atom.logger.LogInfo(atom.req.ID, fmt.Sprintf("Initiating on Ethereum blockchain for ERC20 (%s)", atom.req.SendToken))

	sendValue, ok := big.NewInt(0).SetString(hex.EncodeToString(atom.req.SendValue[:]), 16)
	if !ok {
		return fmt.Errorf("Invalid Send Value: %s", atom.req.SendValue)
	}

	// Approve the contract to transfer tokens
	if err := atom.client.Transact(
		context.Background(),
		nil,
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := atom.tokenBinder.Approve(tops, atom.client.GetSwapperAddress(atom.req.SendToken), sendValue)
			if err != nil {
				return tx, err
			}
			atom.logger.LogInfo(atom.req.ID, atom.client.FormatTransactionView(fmt.Sprintf("Approved %f %s on ethereum blockchain", float64(sendValue.Int64())/100000000, atom.req.SendToken), tx.Hash().String()))
			return tx, nil
		},
		nil,
		1,
	); err != nil {
		return err
	}

	// Initiate the Atomic Swap
	return atom.client.Transact(
		context.Background(),
		func() bool {
			initiatable, err := atom.binder.Initiatable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return initiatable
		},
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := atom.binder.Initiate(tops, atom.id, common.HexToAddress(atom.req.SendToAddress), atom.req.SecretHash, big.NewInt(atom.req.TimeLock), sendValue)
			if err != nil {
				return tx, err
			}
			atom.logger.LogInfo(atom.req.ID, atom.client.FormatTransactionView("Initiated the atomic swap on Ethereum blockchain", tx.Hash().String()))
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
	)
}

// Refund an Atom swap by calling a function on ethereum
func (atom *erc20Atom) Refund() error {
	atom.logger.LogInfo(atom.req.ID, "Refunding the atomic swap on ERC20 blockchain")
	return atom.client.Transact(
		context.Background(),
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
			atom.logger.LogInfo(atom.req.ID, atom.client.FormatTransactionView("Refunded the atomic swap on Ethereum blockchain", tx.Hash().String()))
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
	)
}

// AuditSecret audits the secret of an Atom swap by calling a function on ethereum
func (atom *erc20Atom) AuditSecret() ([32]byte, error) {
	for {
		atom.logger.LogInfo(atom.req.ID, "Auditing secret on ethereum blockchain")
		redeemable, err := atom.binder.Redeemable(&bind.CallOpts{}, atom.id)
		if err != nil {
			atom.logger.LogError(atom.req.ID, err)
			return [32]byte{}, err
		}
		if !redeemable {
			break
		}
		if time.Now().Unix() > atom.req.TimeLock {
			return [32]byte{}, fmt.Errorf("Timed Out")
		}
		time.Sleep(15 * time.Second)
	}
	secret, err := atom.binder.AuditSecret(&bind.CallOpts{}, atom.id)
	if err != nil {
		return [32]byte{}, err
	}
	atom.logger.LogInfo(atom.req.ID, fmt.Sprintf("Audit success on ethereum blockchain secret=%s", base64.StdEncoding.EncodeToString(secret[:])))
	return secret, nil
}

// Audit an Atom swap by calling a function on ethereum
func (atom *erc20Atom) Audit() error {
	atom.logger.LogInfo(atom.req.ID, fmt.Sprintf("Waiting for initiation on ethereum blockchain"))
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
	recvValue, ok := big.NewInt(0).SetString(hex.EncodeToString(atom.req.ReceiveValue[:]), 16)
	if !ok {
		return fmt.Errorf("Invalid Receive Value %s", recvValue)
	}
	if auditReport.Value.Cmp(recvValue) != 0 {
		return fmt.Errorf("Receive Value Mismatch Expected: %v Actual: %v", atom.req.ReceiveValue, auditReport.Value)
	}
	atom.logger.LogInfo(atom.req.ID, fmt.Sprintf("Audit successful on Ethereum blockchain"))
	return nil
}

// Redeem an Atom swap by calling a function on ethereum
func (atom *erc20Atom) Redeem(secret [32]byte) error {
	atom.logger.LogInfo(atom.req.ID, "Redeeming the atomic swap on Ethereum blockchain")
	return atom.client.Transact(
		context.Background(),
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
			atom.logger.LogInfo(atom.req.ID, atom.client.FormatTransactionView("Redeemed the atomic swap on ERC20 blockchain", tx.Hash().String()))
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
	)
}

func getSwapDetails(req foundation.Swap) (foundation.Token, int64, error) {
	if req.SendToken.Blockchain != "Ethereum" && req.ReceiveToken.Blockchain != "Ethereum" {
		return foundation.Token{}, 0, errors.New("Expected one of the tokens to be ethereum")
	}
	if (req.IsFirst && req.SendToken.Blockchain == "Ethereum") || (!req.IsFirst && req.ReceiveToken.Blockchain == "Ethereum") {
		if req.SendToken.Blockchain == "Ethereum" {
			return req.SendToken, req.TimeLock, nil
		}
		return req.ReceiveToken, req.TimeLock, nil
	}

	if req.SendToken.Blockchain == "Ethereum" {
		return req.SendToken, req.TimeLock - 24*60*60, nil
	}
	return req.ReceiveToken, req.TimeLock - 24*60*60, nil
}
