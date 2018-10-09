package eth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/swapperd/adapter/config"
	"github.com/republicprotocol/swapperd/adapter/keystore"
	swapDomain "github.com/republicprotocol/swapperd/domain/swap"
	"github.com/republicprotocol/swapperd/domain/token"
	"github.com/republicprotocol/swapperd/service/guardian"
	"github.com/republicprotocol/swapperd/service/logger"
	"github.com/republicprotocol/swapperd/service/swap"
)

type erc20Atom struct {
	id     [32]byte
	client Conn
	key    keystore.EthereumKey
	req    swapDomain.Request
	logger logger.Logger
	binder *RenExAtomicSwapperERC20
}

// NewERC20Atom returns a new ERC20 Atom instance
func NewERC20Atom(conf config.EthereumNetwork, key keystore.EthereumKey, logger logger.Logger, req swapDomain.Request) (swap.Atom, error) {
	conn, err := NewConnWithConfig(conf)
	if err != nil {
		return nil, err
	}

	contract, err := NewRenExAtomicSwapper(conn.RenExAtomicSwapper, bind.ContractBackend(conn.Client))
	if err != nil {
		return nil, err
	}

	addr, expiry, err := buildValues(req, key.Address.String())
	if err != nil {
		return nil, err
	}

	id, err := contract.SwapID(&bind.CallOpts{}, common.HexToAddress(addr), req.SecretHash, big.NewInt(expiry))
	if err != nil {
		return nil, err
	}

	logger.LogInfo(req.UID, fmt.Sprintf("ERC20 Atomic Swap ID: %s", base64.StdEncoding.EncodeToString(id[:])))
	req.TimeLock = expiry

	return &erc20Atom{
		client: conn,
		key:    key,
		binder: contract,
		logger: logger,
		req:    req,
		id:     id,
	}, nil
}

// Initiate a new Atom swap by calling a function on ethereum
func (atom *ethereumAtom) Initiate() error {
	atom.logger.LogInfo(atom.req.UID, "Initiating on Ethereum blockchain for ERC20(%s)", atom.req.Token)
	initiatable, err := atom.binder.Initiatable(&bind.CallOpts{}, atom.id)
	if err != nil {
		return err
	}
	if !initiatable {
		return swap.ErrSwapAlreadyInitiated
	}
	if err := atom.key.SubmitTx(
		func(tops *bind.TransactOpts) error {
			val, ok := big.NewInt(0).SetString(atom.req.SendValue, 10)
			if !ok {
				return fmt.Errorf("Invalid Send Value: %s", atom.req.SendValue)
			}
			prevValue := tops.Value
			tops.Value = val
			_, err := atom.binder.Initiate(atom.key.TransactOpts, atom.id, common.HexToAddress(atom.req.SendToAddress), atom.req.SecretHash, big.NewInt(atom.req.TimeLock))
			atom.key.TransactOpts.Value = prevValue
			return err
		},
		func() bool {
			initiatable, _ := atom.binder.Initiatable(&bind.CallOpts{}, atom.id)
			return !initiatable
		},
	); err != nil {
		return err
	}
	atom.logger.LogInfo(atom.req.UID, fmt.Sprintf("Initiated the atomic swap on ERC20 blockchain"))
	return nil
}

// Refund an Atom swap by calling a function on ethereum
func (atom *ethereumAtom) Refund() error {
	atom.logger.LogInfo(atom.req.UID, "Refunding the atomic swap on ERC20 blockchain")
	refundable, err := atom.binder.Refundable(&bind.CallOpts{}, atom.id)
	if err != nil {
		return err
	}
	if !refundable {
		return guardian.ErrNotRefundable
	}
	if err := atom.key.SubmitTx(
		func(tops *bind.TransactOpts) error {
			_, err = atom.binder.Refund(atom.key.TransactOpts, atom.id)
			return err
		},
		func() bool {
			refundable, _ := atom.binder.Refundable(&bind.CallOpts{}, atom.id)
			return !refundable
		},
	); err != nil {
		return err
	}
	atom.logger.LogInfo(atom.req.UID, "Refunded the atomic swap on ERC20 blockchain")
	return nil
}

// AuditSecret audits the secret of an Atom swap by calling a function on ethereum
func (atom *ethereumAtom) AuditSecret() ([32]byte, error) {
	for {
		atom.logger.LogInfo(atom.req.UID, "Auditing secret on ethereum blockchain")
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
	atom.logger.LogInfo(atom.req.UID, fmt.Sprintf("Audit success on ethereum blockchain secret=%s", base64.StdEncoding.EncodeToString(secret[:])))
	return secret, nil
}

// Audit an Atom swap by calling a function on ethereum
func (atom *ethereumAtom) Audit() error {
	atom.logger.LogInfo(atom.req.UID, fmt.Sprintf("Waiting for initiation on ethereum blockchain"))
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
	recvValue, ok := big.NewInt(0).SetString(atom.req.ReceiveValue, 10)
	if !ok {
		return fmt.Errorf("Invalid Receive Value %s", recvValue)
	}
	if auditReport.Value.Cmp(recvValue) != 0 {
		return fmt.Errorf("Receive Value Mismatch Expected: %v Actual: %v", atom.req.ReceiveValue, auditReport.Value)
	}
	atom.logger.LogInfo(atom.req.UID, fmt.Sprintf("Audit successful on ERC20 blockchain"))
	return nil
}

// Redeem an Atom swap by calling a function on ethereum
func (atom *ethereumAtom) Redeem(secret [32]byte) error {
	atom.logger.LogInfo(atom.req.UID, "Redeeming the atomic swap on ERC20 blockchain")
	redeemable, err := atom.binder.Redeemable(&bind.CallOpts{}, atom.id)
	if err != nil {
		return err
	}
	if !redeemable {
		return swap.ErrSwapAlreadyRedeemedOrRefunded
	}

	if err := atom.key.SubmitTx(
		func(tops *bind.TransactOpts) error {
			_, err = atom.binder.Redeem(atom.key.TransactOpts, atom.id, secret)
			return err
		},
		func() bool {
			refundable, _ := atom.binder.Redeemable(&bind.CallOpts{}, atom.id)
			return !refundable
		},
	); err != nil {
		return err
	}
	atom.logger.LogInfo(atom.req.UID, "Redeemed the atomic swap on ERC20 blockchain")
	return nil
}

func buildValues(req swapDomain.Request, personalAddr string) (string, int64, error) {
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
		addr = personalAddr
	}
	return addr, expiry, nil
}
