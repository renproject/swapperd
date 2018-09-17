package eth

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	bindings "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/bindings/eth"
	ethclient "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/clients/eth"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/domain/order"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type Adapter interface {
	RecieveSwapDetails(order.ID, int64) ([]byte, error)
}

// EthereumData
type EthereumData struct {
	SwapID   [32]byte `json:"swap_id"`
	HashLock [32]byte `json:"hash_lock"`
}

type EthereumAtom struct {
	orderID [32]byte
	context context.Context
	client  ethclient.Conn
	key     keystore.EthereumKey
	binding *bindings.AtomicSwap
	adapter Adapter
	data    EthereumData
}

// NewEthereumAtom returns a new Ethereum RequestAtom instance
func NewEthereumAtom(adapter Adapter, conf config.EthereumNetwork, key keystore.EthereumKey, orderID [32]byte) (swap.Atom, error) {
	conn, err := ethclient.Connect(conf)
	if err != nil {
		return &EthereumAtom{}, err
	}

	contract, err := bindings.NewAtomicSwap(conn.RenExAtomicSwapper, bind.ContractBackend(conn.Client))
	if err != nil {
		return &EthereumAtom{}, err
	}

	swapID := [32]byte{}
	rand.Read(swapID[:])

	return &EthereumAtom{
		context: context.Background(),
		client:  conn,
		key:     key,
		binding: contract,
		orderID: orderID,
		adapter: adapter,
		data: EthereumData{
			SwapID: swapID,
		},
	}, nil
}

// Initiate a new Atom swap by calling a function on ethereum
func (atom *EthereumAtom) Initiate(to []byte, hash [32]byte, value *big.Int, expiry int64) error {
	prevValue := atom.key.TransactOpts.Value
	prevGasLimit := atom.key.TransactOpts.GasLimit
	atom.key.TransactOpts.Value = value
	atom.key.TransactOpts.GasLimit = 3000000
	atom.data.HashLock = hash

	tx, err := atom.binding.Initiate(atom.key.TransactOpts, atom.data.SwapID, common.BytesToAddress(to), hash, big.NewInt(expiry))
	atom.key.TransactOpts.Value = prevValue
	atom.key.TransactOpts.GasLimit = prevGasLimit
	if err != nil {
		return err
	}

	_, err = atom.client.PatchedWaitMined(atom.context, tx)
	return err
}

// Redeem an Atom swap by calling a function on ethereum
func (atom *EthereumAtom) Redeem(secret [32]byte) error {
	prevGasLimit := atom.key.TransactOpts.GasLimit
	atom.key.TransactOpts.GasLimit = 3000000
	tx, err := atom.binding.Redeem(atom.key.TransactOpts, atom.data.SwapID, secret)
	atom.key.TransactOpts.GasLimit = prevGasLimit
	if err == nil {
		_, err = atom.client.PatchedWaitMined(atom.context, tx)
	}

	return err
}

// WaitForCounterRedemption waits for the counter-party to initiate.
func (atom *EthereumAtom) WaitForCounterRedemption() error {
	for {
		secret, err := atom.binding.AuditSecret(&bind.CallOpts{}, atom.data.SwapID)
		if err != nil || secret == [32]byte{} {
			time.Sleep(1 * time.Second)
			continue
		}
		return nil
	}
}

// Refund an Atom swap by calling a function on ethereum
func (atom *EthereumAtom) Refund() error {
	prevGasLimit := atom.key.TransactOpts.GasLimit
	atom.key.TransactOpts.GasLimit = 3000000
	tx, err := atom.binding.Refund(atom.key.TransactOpts, atom.data.SwapID)
	atom.key.TransactOpts.GasLimit = prevGasLimit
	if err == nil {
		_, err = atom.client.PatchedWaitMined(atom.context, tx)
	}
	return err
}

// Audit an Atom swap by calling a function on ethereum
func (atom *EthereumAtom) Audit() ([32]byte, []byte, *big.Int, int64, error) {
	details, err := atom.adapter.RecieveSwapDetails(atom.orderID, time.Now().Add(15*time.Minute).Unix())
	if err != nil {
		return [32]byte{}, nil, nil, 0, err
	}

	if err := atom.Deserialize(details); err != nil {
		return [32]byte{}, nil, nil, 0, err
	}
	auditReport, err := atom.binding.Audit(&bind.CallOpts{}, atom.data.SwapID)
	if err != nil {
		return [32]byte{}, nil, nil, 0, err
	}
	return auditReport.SecretLock, auditReport.To.Bytes(), auditReport.Value, auditReport.Timelock.Int64(), nil
}

// AuditSecret audits the secret of an Atom swap by calling a function on ethereum
func (atom *EthereumAtom) AuditSecret() ([32]byte, error) {
	details, err := atom.adapter.RecieveSwapDetails(atom.orderID, time.Now().Add(15*time.Minute).Unix())
	if err != nil {
		return [32]byte{}, err
	}

	if err := atom.Deserialize(details); err != nil {
		return [32]byte{}, err
	}
	return atom.binding.AuditSecret(&bind.CallOpts{}, atom.data.SwapID)
}

// RedeemedAt returns the timestamp at which the atom is redeemed
func (atom *EthereumAtom) RedeemedAt() (int64, error) {
	details, err := atom.adapter.RecieveSwapDetails(atom.orderID, time.Now().Add(15*time.Minute).Unix())
	if err != nil {
		return 0, err
	}

	if err := atom.Deserialize(details); err != nil {
		return 0, err
	}

	redeemedAt, err := atom.binding.RedeemedAt(&bind.CallOpts{}, atom.data.SwapID)
	if err != nil {
		return 0, err
	}

	return redeemedAt.Int64(), nil
}

// Serialize serializes the atom details
func (atom *EthereumAtom) Serialize() ([]byte, error) {
	return json.Marshal(atom.data)
}

// Deserialize deserializes the atom details
func (atom *EthereumAtom) Deserialize(data []byte) error {
	return json.Unmarshal(data, &atom.data)
}

// GetFromAddress returns the address of the sender
func (atom *EthereumAtom) GetFromAddress() ([]byte, error) {
	return []byte(atom.key.Address.String()), nil
}

// PriorityCode returns the priority code of the currency.
func (atom *EthereumAtom) PriorityCode() uint32 {
	return uint32(token.ETH)
}
