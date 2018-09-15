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
	"github.com/republicprotocol/renex-swapper-go/adapter/configs/keystore"
	"github.com/republicprotocol/renex-swapper-go/domains/order"
	"github.com/republicprotocol/renex-swapper-go/services/swap"
)

type Adapter interface {
	ReceiveSwapDetails(order.ID, int64) ([]byte, error)
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
	key     keystore.Key
	binding *bindings.AtomicSwap
	adapter Adapter
	data    EthereumData
}

// NewEthereumAtom returns a new Ethereum RequestAtom instance
func NewEthereumAtom(adapter Adapter, client ethclient.Conn, key keystore.Key, orderID [32]byte) (swap.Atom, error) {
	contract, err := bindings.NewAtomicSwap(client.RenExAtomicSwapperAddress(), bind.ContractBackend(client.Client()))
	if err != nil {
		return &EthereumAtom{}, err
	}

	swapID := [32]byte{}
	rand.Read(swapID[:])

	return &EthereumAtom{
		context: context.Background(),
		client:  client,
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
	key, err := atom.key.GetKey()
	if err != nil {
		return err
	}
	auth := bind.NewKeyedTransactor(key)
	auth.Value = value
	auth.GasLimit = 3000000
	atom.data.HashLock = hash

	tx, err := atom.binding.Initiate(auth, atom.data.SwapID, common.BytesToAddress(to), hash, big.NewInt(expiry))
	if err != nil {
		return err
	}
	_, err = atom.client.PatchedWaitMined(atom.context, tx)
	return err
}

// Redeem an Atom swap by calling a function on ethereum
func (atom *EthereumAtom) Redeem(secret [32]byte) error {
	key, err := atom.key.GetKey()
	if err != nil {
		return err
	}
	auth := bind.NewKeyedTransactor(key)
	auth.GasLimit = 3000000
	tx, err := atom.binding.Redeem(auth, atom.data.SwapID, secret)
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
	key, err := atom.key.GetKey()
	if err != nil {
		return err
	}
	auth := bind.NewKeyedTransactor(key)
	auth.GasLimit = 3000000
	tx, err := atom.binding.Refund(auth, atom.data.SwapID)
	if err == nil {
		_, err = atom.client.PatchedWaitMined(atom.context, tx)
	}
	return err
}

// Audit an Atom swap by calling a function on ethereum
func (atom *EthereumAtom) Audit() ([32]byte, []byte, *big.Int, int64, error) {
	details, err := atom.adapter.ReceiveSwapDetails(atom.orderID, time.Now().Add(15*time.Minute).Unix())
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
	details, err := atom.adapter.ReceiveSwapDetails(atom.orderID, time.Now().Add(15*time.Minute).Unix())
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
	details, err := atom.adapter.ReceiveSwapDetails(atom.orderID, time.Now().Add(15*time.Minute).Unix())
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
	return atom.key.GetAddress()
}

// PriorityCode returns the priority code of the currency.
func (atom *EthereumAtom) PriorityCode() uint32 {
	return atom.key.PriorityCode()
}
