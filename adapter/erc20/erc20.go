package erc20

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/renex-swapper-go/service/guardian"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
	"github.com/republicprotocol/swapperd/adapter/config"
	"github.com/republicprotocol/swapperd/adapter/keystore"
	"github.com/republicprotocol/swapperd/core"
	"github.com/republicprotocol/swapperd/foundation"
)

type erc20Atom struct {
	id          [32]byte
	client      Conn
	key         keystore.EthereumKey
	req         foundation.Swap
	logger      core.Logger
	tokenBinder *CompatibleERC20
	binder      *RenExAtomicSwapper
}

// NewERC20Atom returns a new ERC20 Atom instance
func NewERC20Atom(conf config.EthereumNetwork, key keystore.EthereumKey, logger core.Logger, req foundation.Swap) (core.SwapContractBinder, error) {
	conn, err := NewConnWithConfig(conf)
	if err != nil {
		return nil, err
	}

	tokenContract, err := NewCompatibleERC20(common.HexToAddress("0xA1D3EEcb76285B4435550E4D963B8042A8bffbF0"), bind.ContractBackend(conn.Client))
	if err != nil {
		return nil, err
	}

	contract, err := NewRenExAtomicSwapper(conn.SwapperAddresses[foundation.TokenWBTC], bind.ContractBackend(conn.Client))
	if err != nil {
		return nil, err
	}

	expiry, err := getExpiry(req)
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
		client:      conn,
		key:         key,
		tokenBinder: tokenContract,
		binder:      contract,
		logger:      logger,
		req:         req,
		id:          id,
	}, nil
}

// Initiate a new Atom swap by calling a function on ethereum
func (atom *erc20Atom) Initiate() error {
	atom.logger.LogInfo(atom.req.ID, fmt.Sprintf("Initiating on Ethereum blockchain for ERC20(%s)", atom.req.SendToken))
	initiatable, err := atom.binder.Initiatable(&bind.CallOpts{}, atom.id)
	if err != nil {
		return err
	}
	if !initiatable {
		return nil
	}

	sendValue, ok := big.NewInt(0).SetString(hex.EncodeToString(atom.req.SendValue[:]), 16)
	if !ok {
		return fmt.Errorf("Invalid Send Value: %s", atom.req.SendValue)
	}

	// Approve the contract to transfer tokens
	if err := atom.key.SubmitTx(
		func(tops *bind.TransactOpts) error {
			tx, err := atom.tokenBinder.Approve(tops, atom.client.SwapperAddresses[atom.req.SendToken], sendValue)
			if err != nil {
				return err
			}
			atom.logger.LogInfo(atom.req.ID, atom.client.FormatTransactionView(fmt.Sprintf("Approved %f %s on ethereum blockchain", float64(sendValue.Int64())/100000000, atom.req.SendToken), tx.Hash().String()))
			return nil
		},
		func() bool {
			allowance, err := atom.tokenBinder.Allowance(&bind.CallOpts{}, atom.key.Address, atom.client.SwapperAddresses[atom.req.SendToken])
			if err != nil {
				atom.logger.LogError(atom.req.ID, fmt.Sprintf("Error: %v", err))
			}
			return sendValue.Cmp(allowance) == 0
		},
	); err != nil {
		return err
	}

	// Initiate the swap
	if err := atom.key.SubmitTx(
		func(tops *bind.TransactOpts) error {
			tx, err := atom.binder.Initiate(atom.key.TransactOpts, atom.id, common.HexToAddress(atom.req.SendToAddress), atom.req.SecretHash, big.NewInt(atom.req.TimeLock), sendValue)
			if err != nil {
				return err
			}
			atom.logger.LogInfo(atom.req.ID, atom.client.FormatTransactionView("Initiated the atomic swap on Ethereum blockchain", tx.Hash().String()))
			return nil
		},
		func() bool {
			initiatable, err := atom.binder.Initiatable(&bind.CallOpts{}, atom.id)
			if err != nil {
				atom.logger.LogError(atom.req.ID, fmt.Sprintf("Error: %v", err))
			}
			return !initiatable
		},
	); err != nil {
		return err
	}
	return nil
}

// Refund an Atom swap by calling a function on ethereum
func (atom *erc20Atom) Refund() error {
	atom.logger.LogInfo(atom.req.ID, "Refunding the atomic swap on ERC20 blockchain")
	refundable, err := atom.binder.Refundable(&bind.CallOpts{}, atom.id)
	if err != nil {
		return err
	}
	if !refundable {
		return guardian.ErrNotRefundable
	}
	if err := atom.key.SubmitTx(
		func(tops *bind.TransactOpts) error {
			tx, err := atom.binder.Refund(atom.key.TransactOpts, atom.id)
			if err != nil {
				return err
			}
			atom.logger.LogInfo(atom.req.ID, atom.client.FormatTransactionView("Refunded the atomic swap on Ethereum blockchain", tx.Hash().String()))
			return nil
		},
		func() bool {
			refundable, _ := atom.binder.Refundable(&bind.CallOpts{}, atom.id)
			return !refundable
		},
	); err != nil {
		return err
	}
	return nil
}

// AuditSecret audits the secret of an Atom swap by calling a function on ethereum
func (atom *erc20Atom) AuditSecret() ([32]byte, error) {
	for {
		atom.logger.LogInfo(atom.req.ID, "Auditing secret on ethereum blockchain")
		redeemable, err := atom.binder.Redeemable(&bind.CallOpts{}, atom.id)
		if err != nil {
			return [32]byte{}, err
		}
		if !redeemable {
			break
		}
		if time.Now().Unix() > atom.req.TimeLock {
			return [32]byte{}, swap.ErrTimedOut
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
	atom.logger.LogInfo(atom.req.ID, fmt.Sprintf("Audit successful on ERC20 blockchain"))
	return nil
}

// Redeem an Atom swap by calling a function on ethereum
func (atom *erc20Atom) Redeem(secret [32]byte) error {
	atom.logger.LogInfo(atom.req.ID, "Redeeming the atomic swap on ERC20 blockchain")
	redeemable, err := atom.binder.Redeemable(&bind.CallOpts{}, atom.id)
	if err != nil {
		return err
	}
	if !redeemable {
		return swap.ErrSwapAlreadyRedeemedOrRefunded
	}

	if err := atom.key.SubmitTx(
		func(tops *bind.TransactOpts) error {
			tx, err := atom.binder.Redeem(atom.key.TransactOpts, atom.id, secret)
			if err != nil {
				return err
			}
			atom.logger.LogInfo(atom.req.ID, atom.client.FormatTransactionView("Redeemed the atomic swap on ERC20 blockchain", tx.Hash().String()))
			return nil
		},
		func() bool {
			refundable, _ := atom.binder.Redeemable(&bind.CallOpts{}, atom.id)
			return !refundable
		},
	); err != nil {
		return err
	}
	return nil
}

func getExpiry(req foundation.Swap) (int64, error) {
	if req.SendToken != foundation.TokenWBTC && req.ReceiveToken != foundation.TokenWBTC {
		return 0, errors.New("Expected one of the tokens to be ethereum")
	}
	if (req.IsFirst && req.SendToken == foundation.TokenWBTC) || (!req.IsFirst && req.ReceiveToken == foundation.TokenWBTC) {
		return req.TimeLock, nil
	}
	return req.TimeLock - 24*60*60, nil
}
