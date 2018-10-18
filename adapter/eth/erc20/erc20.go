package erc20

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/republicprotocol/beth-go"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/swapperd/core"
	"github.com/republicprotocol/swapperd/foundation"
)

type erc20SwapContractBinder struct {
	id             [32]byte
	account        beth.Account
	req            foundation.Swap
	logger         core.Logger
	swapperAddress common.Address
	tokenAddress   common.Address
	swapperBinder  *RenExAtomicSwapper
	tokenBinder    *CompatibleERC20
}

// NewERC20SwapContractBinder returns a new ERC20 Atom instance
func NewERC20SwapContractBinder(account beth.Account, req foundation.Swap, logger core.Logger) (core.SwapContractBinder, error) {
	token, expiry, err := getSwapDetails(req)
	if err != nil {
		return nil, err
	}

	tokenAddress, err := account.ReadAddress(fmt.Sprintf("ERC20:%s", token.Name))
	if err != nil {
		return nil, err
	}

	swapperAddress, err := account.ReadAddress(fmt.Sprintf("SWAPPER:%s", token.Name))
	if err != nil {
		return nil, err
	}

	client := account.EthClient()

	tokenBinder, err := NewCompatibleERC20(tokenAddress, bind.ContractBackend(client.EthClient()))
	if err != nil {
		return nil, err
	}

	swapperBinder, err := NewRenExAtomicSwapper(swapperAddress, bind.ContractBackend(client.EthClient()))
	if err != nil {
		return nil, err
	}

	id, err := swapperBinder.SwapID(&bind.CallOpts{}, req.SecretHash, big.NewInt(expiry))
	if err != nil {
		return nil, err
	}

	logger.LogInfo(req.ID, fmt.Sprintf("ERC20 Atomic Swap ID: %s", base64.StdEncoding.EncodeToString(id[:])))
	req.TimeLock = expiry

	return &erc20SwapContractBinder{
		account:        account,
		swapperAddress: swapperAddress,
		tokenAddress:   tokenAddress,
		swapperBinder:  swapperBinder,
		tokenBinder:    tokenBinder,
		logger:         logger,
		req:            req,
		id:             id,
	}, nil
}

// Initiate a new Atom swap by calling a function on ethereum
func (atom *erc20SwapContractBinder) Initiate() error {
	atom.logger.LogInfo(atom.req.ID, fmt.Sprintf("Initiating on Ethereum blockchain for ERC20 (%s)", atom.req.SendToken))

	// Approve the contract to transfer tokens
	if err := atom.account.Transact(
		context.Background(),
		nil,
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := atom.tokenBinder.Approve(tops, atom.swapperAddress, atom.req.SendValue)
			if err != nil {
				return tx, err
			}
			return tx, nil
		},
		nil,
		1,
	); err != nil {
		return err
	}

	// Initiate the Atomic Swap
	return atom.account.Transact(
		context.Background(),
		func() bool {
			initiatable, err := atom.swapperBinder.Initiatable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return initiatable
		},
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := atom.swapperBinder.Initiate(tops, atom.id, common.HexToAddress(atom.req.SendToAddress), atom.req.SecretHash, big.NewInt(atom.req.TimeLock), atom.req.SendValue)
			if err != nil {
				return tx, err
			}
			msg, _ := atom.account.FormatTransactionView("Initiated the atomic swap on Ethereum blockchain", tx.Hash().String())
			atom.logger.LogInfo(atom.req.ID, msg)
			return tx, nil
		},
		func() bool {
			initiatable, err := atom.swapperBinder.Initiatable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return !initiatable
		},
		1,
	)
}

// Refund an Atom swap by calling a function on ethereum
func (atom *erc20SwapContractBinder) Refund() error {
	atom.logger.LogInfo(atom.req.ID, "Refunding the atomic swap on ERC20 blockchain")
	return atom.account.Transact(
		context.Background(),
		func() bool {
			refundable, err := atom.swapperBinder.Refundable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return refundable
		},
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := atom.swapperBinder.Refund(tops, atom.id)
			if err != nil {
				return nil, err
			}
			msg, _ := atom.account.FormatTransactionView("Refunded the atomic swap on Ethereum blockchain", tx.Hash().String())
			atom.logger.LogInfo(atom.req.ID, msg)
			return tx, nil
		},
		func() bool {
			refundable, err := atom.swapperBinder.Refundable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return !refundable
		},
		1,
	)
}

// AuditSecret audits the secret of an Atom swap by calling a function on ethereum
func (atom *erc20SwapContractBinder) AuditSecret() ([32]byte, error) {
	for {
		atom.logger.LogInfo(atom.req.ID, "Auditing secret on ethereum blockchain")
		redeemable, err := atom.swapperBinder.Redeemable(&bind.CallOpts{}, atom.id)
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
	secret, err := atom.swapperBinder.AuditSecret(&bind.CallOpts{}, atom.id)
	if err != nil {
		return [32]byte{}, err
	}
	atom.logger.LogInfo(atom.req.ID, fmt.Sprintf("Audit success on ethereum blockchain secret=%s", base64.StdEncoding.EncodeToString(secret[:])))
	return secret, nil
}

// Audit an Atom swap by calling a function on ethereum
func (atom *erc20SwapContractBinder) Audit() error {
	atom.logger.LogInfo(atom.req.ID, fmt.Sprintf("Waiting for initiation on ethereum blockchain"))
	for {
		initiatable, err := atom.swapperBinder.Initiatable(&bind.CallOpts{}, atom.id)
		if err != nil {
			return err
		}
		if !initiatable {
			break
		}
		time.Sleep(15 * time.Second)
	}
	auditReport, err := atom.swapperBinder.Audit(&bind.CallOpts{}, atom.id)
	if err != nil {
		return err
	}
	if auditReport.Value.Cmp(atom.req.ReceiveValue) != 0 {
		return fmt.Errorf("Receive Value Mismatch Expected: %v Actual: %v", atom.req.ReceiveValue, auditReport.Value)
	}
	atom.logger.LogInfo(atom.req.ID, fmt.Sprintf("Audit successful on Ethereum blockchain"))
	return nil
}

// Redeem an Atom swap by calling a function on ethereum
func (atom *erc20SwapContractBinder) Redeem(secret [32]byte) error {
	atom.logger.LogInfo(atom.req.ID, "Redeeming the atomic swap on Ethereum blockchain")
	return atom.account.Transact(
		context.Background(),
		func() bool {
			redeemable, err := atom.swapperBinder.Redeemable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return redeemable
		},
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := atom.swapperBinder.Redeem(tops, atom.id, secret)
			if err != nil {
				return nil, err
			}
			msg, _ := atom.account.FormatTransactionView("Redeemed the atomic swap on Ethereum blockchain", tx.Hash().String())
			atom.logger.LogInfo(atom.req.ID, msg)
			return tx, nil
		},
		func() bool {
			refundable, err := atom.swapperBinder.Refundable(&bind.CallOpts{}, atom.id)
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
