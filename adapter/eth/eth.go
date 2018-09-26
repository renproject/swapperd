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
	fmt.Println(req)
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

	// TODO: Fix data race on transact opts
	prevValue := atom.key.TransactOpts.Value
	prevGasLimit := atom.key.TransactOpts.GasLimit
	var ok bool
	atom.key.TransactOpts.Value, ok = big.NewInt(0).SetString(atom.req.SendValue, 10)
	if !ok {
		return fmt.Errorf("Invalid Send Value: %s", atom.req.SendValue)
	}

	atom.key.TransactOpts.GasLimit = 3000000
	_, err = atom.binder.Initiate(atom.key.TransactOpts, atom.id, common.HexToAddress(atom.req.SendToAddress), atom.req.SecretHash, big.NewInt(atom.req.TimeLock))
	atom.key.TransactOpts.Value = prevValue
	atom.key.TransactOpts.GasLimit = prevGasLimit
	if err != nil {
		return fmt.Errorf("Failed to initiate on the Ethereum blockchain: %v", err)
	}

	if err := atom.waitForInitiation(); err != nil {
		return err
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
	if err := atom.waitForRedemption(); err != nil {
		return [32]byte{}, err
	}
	return atom.binder.AuditSecret(&bind.CallOpts{}, atom.id)
}

// Audit an Atom swap by calling a function on ethereum
func (atom *ethereumAtom) Audit() error {
	if err := atom.waitForInitiation(); err != nil {
		return err
	}
	auditReport, err := atom.binder.Audit(&bind.CallOpts{}, atom.id)
	if err != nil {
		return err
	}
	fmt.Println(auditReport)
	recvValue, ok := big.NewInt(0).SetString(atom.req.ReceiveValue, 10)
	if !ok {
		return fmt.Errorf("Invalid Receive Value %s", recvValue)
	}
	if auditReport.Value.Cmp(recvValue) != 0 {
		return fmt.Errorf("Receive Value Mismatch Expected: %v Actual: %v", atom.req.ReceiveValue, auditReport.Value)
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
	if err := atom.waitForRedemption(); err != nil {
		return err
	}
	fmt.Println("redeemed the atomic swap on ethereum blockchain")
	return nil
}

// TODO: change req to personalReq
func buildValues(req swapDomain.Request, personalAddr string) (string, int64, error) {
	var addr string
	var expiry int64
	if req.SendToken != token.ETH && req.ReceiveToken != token.ETH {
		return "", 0, errors.New("Expected one of the tokens to be ethereum")
	}

	if (req.GoesFirst && req.SendToken == token.ETH) || (!req.GoesFirst && req.ReceiveToken == token.ETH) {
		expiry = req.TimeLock
	} else {
		// TODO: Document times
		expiry = req.TimeLock - 24*60*60
	}

	if req.SendToken == token.ETH {
		addr = req.SendToAddress
	} else {
		addr = personalAddr
	}

	return addr, expiry, nil
}

func (atom *ethereumAtom) waitForInitiation() error {
	for {
		fmt.Println("Waiting for initiation ......  on ethereum blockchain")
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
	return nil
}

func (atom *ethereumAtom) waitForRedemption() error {
	for {
		fmt.Println("Waiting for counter party redemption ......  on ethereum blockchain")
		secret, err := atom.binder.AuditSecret(&bind.CallOpts{}, atom.id)
		if err != nil {
			return err
		}
		if secret != [32]byte{} {
			return nil
		}
		if time.Now().Unix() > atom.req.TimeLock {
			break
		}
		time.Sleep(1 * time.Minute)
	}
	return errors.New("Timed Out")
}
