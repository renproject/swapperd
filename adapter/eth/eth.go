package eth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	swapDomain "github.com/republicprotocol/renex-swapper-go/domain/swap"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/service/logger"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type ethereumAtom struct {
	id     [32]byte
	client Conn
	key    keystore.EthereumKey
	req    swapDomain.Request
	logger logger.Logger
	binder *RenExAtomicSwapper
}

// NewEthereumAtom returns a new Ethereum RequestAtom instance
func NewEthereumAtom(conf config.EthereumNetwork, key keystore.EthereumKey, req swapDomain.Request) (swap.Atom, error) {
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

	fmt.Printf("Swap ID: %s\n", base64.StdEncoding.EncodeToString(id[:]))
	// logger.LogDebug(req.UID, fmt.Sprintf("Swap ID: %s", base64.StdEncoding.EncodeToString(id[:])))
	req.TimeLock = expiry

	return &ethereumAtom{
		client: conn,
		key:    key,
		binder: contract,
		// logger:  logger,
		req: req,
		id:  id,
	}, nil
}

// Initiate a new Atom swap by calling a function on ethereum
func (atom *ethereumAtom) Initiate() error {
	// atom.logger.LogInfo(atom.req.UID, "Initiating on Ethereum blockchain")
	fmt.Println("Initiating on Ethereum blockchain")
	initiatable, err := atom.binder.Initiatable(&bind.CallOpts{}, atom.id)
	if err != nil {
		return err
	}
	if !initiatable {
		//		atom.logger.LogInfo(atom.req.UID, "Skipping initiation as it is already initiated")
		return swap.ErrSwapAlreadyInitiated
	}

	prevValue := atom.key.TransactOpts.Value
	prevGasLimit := atom.key.TransactOpts.GasLimit
	atom.key.TransactOpts.Value = atom.req.SendValue
	atom.key.TransactOpts.GasLimit = 3000000
	_, err = atom.binder.Initiate(atom.key.TransactOpts, atom.id, common.HexToAddress(atom.req.SendToAddress), atom.req.SecretHash, big.NewInt(atom.req.TimeLock))
	atom.key.TransactOpts.Value = prevValue
	atom.key.TransactOpts.GasLimit = prevGasLimit
	if err != nil {
		return fmt.Errorf("Failed to initiate on the Ethereum blockchain: %v", err)
	}
	fmt.Println("Initiated the atomic swap on ethereum blockchain")
	// atom.logger.LogInfo(atom.req.UID,
	// 	fmt.Sprintf("Initiated the atomic swap on ethereum blockchain (TX Hash): %s", tx.Hash().String()))
	return nil
}

// Refund an Atom swap by calling a function on ethereum
func (atom *ethereumAtom) Refund() error {
	fmt.Println("Refunding the atomic swap on ethereum blockchain")
	refundable, err := atom.binder.Refundable(&bind.CallOpts{}, atom.id)
	if err != nil {
		return err
	}
	if !refundable {
		return swap.ErrSwapAlreadyRedeemedOrRefunded
	}
	prevGasLimit := atom.key.TransactOpts.GasLimit
	atom.key.TransactOpts.GasLimit = 3000000
	_, err = atom.binder.Refund(atom.key.TransactOpts, atom.id)
	atom.key.TransactOpts.GasLimit = prevGasLimit
	if err != nil {
		return err
	}
	fmt.Println("Refunded the atomic swap on ethereum blockchain")
	return nil
}

// AuditSecret audits the secret of an Atom swap by calling a function on ethereum
func (atom *ethereumAtom) AuditSecret() ([32]byte, error) {
	for {
		fmt.Println("Waiting for counter party redemption ......  on ethereum blockchain")
		secret, err := atom.binder.AuditSecret(&bind.CallOpts{}, atom.id)
		if err != nil {
			return [32]byte{}, err
		}
		if secret != [32]byte{} {
			return secret, nil
		}
		if time.Now().Unix() > atom.req.TimeLock {
			break
		}
		time.Sleep(1 * time.Minute)
	}
	return [32]byte{}, errors.New("Timed Out")
}

// Audit an Atom swap by calling a function on ethereum
func (atom *ethereumAtom) Audit() error {
	for {
		fmt.Println("Waiting for counter party initiation ......  on ethereum blockchain")
		auditReport, err := atom.binder.Audit(&bind.CallOpts{}, atom.id)
		if err != nil {
			return err
		}
		if auditReport.To.String() != auditReport.From.String() {
			break
		}
		if time.Now().Unix() > atom.req.TimeLock {
			return errors.New("Timed Out")
		}
		time.Sleep(1 * time.Minute)
	}

	auditReport, err := atom.binder.Audit(&bind.CallOpts{}, atom.id)
	if err != nil {
		return err
	}

	if auditReport.From.String() != atom.req.ReceiveFromAddress {
		return errors.New("From Address Mismatch")
	}

	if auditReport.SecretLock != atom.req.SecretHash {
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
func (atom *ethereumAtom) Redeem(secret [32]byte) error {
	fmt.Println("redeeming the atomic swap on ethereum blockchain")
	redeemable, err := atom.binder.Redeemable(&bind.CallOpts{}, atom.id)
	if err != nil {
		return err
	}
	if !redeemable {
		return swap.ErrSwapAlreadyRedeemedOrRefunded
	}
	prevGasLimit := atom.key.TransactOpts.GasLimit
	atom.key.TransactOpts.GasLimit = 3000000
	_, err = atom.binder.Redeem(atom.key.TransactOpts, atom.id, secret)
	atom.key.TransactOpts.GasLimit = prevGasLimit
	if err != nil {
		return err
	}
	fmt.Println("redeemed the atomic swap on ethereum blockchain")
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