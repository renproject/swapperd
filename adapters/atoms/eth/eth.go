package eth

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"time"

	"github.com/republicprotocol/atom-go/services/store"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	bindings "github.com/republicprotocol/atom-go/adapters/blockchain/bindings/eth"
	ethclient "github.com/republicprotocol/atom-go/adapters/blockchain/clients/eth"
	"github.com/republicprotocol/atom-go/services/swap"
)

// EthereumData
type EthereumData struct {
	SwapID   [32]byte `json:"swap_id"`
	HashLock [32]byte `json:"hash_lock"`
}

type EthereumAtom struct {
	orderID [32]byte
	context context.Context
	client  ethclient.Conn
	key     swap.Key
	binding *bindings.AtomicSwap
	data    EthereumData
}

// NewEthereumAtom returns a new Ethereum RequestAtom instance
func NewEthereumAtom(client ethclient.Conn, key swap.Key, orderID [32]byte) (swap.Atom, error) {
	contract, err := bindings.NewAtomicSwap(client.AtomAddress(), bind.ContractBackend(client.Client()))
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
		data: EthereumData{
			SwapID: swapID,
		},
	}, nil
}

// Initiate a new Atom swap by calling a function on ethereum
func (atom *EthereumAtom) Initiate(to []byte, hash [32]byte, value *big.Int, expiry int64) error {
	auth := bind.NewKeyedTransactor(atom.key.GetKey())
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
	auth := bind.NewKeyedTransactor(atom.key.GetKey())
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
	auth := bind.NewKeyedTransactor(atom.key.GetKey())
	auth.GasLimit = 3000000
	tx, err := atom.binding.Refund(auth, atom.data.SwapID)
	if err == nil {
		_, err = atom.client.PatchedWaitMined(atom.context, tx)
	}
	return err
}

// Audit an Atom swap by calling a function on ethereum
func (atom *EthereumAtom) Audit() ([32]byte, []byte, *big.Int, int64, error) {
	auditReport, err := atom.binding.Audit(&bind.CallOpts{}, atom.data.SwapID)
	if err != nil {
		return [32]byte{}, nil, nil, 0, err
	}
	return auditReport.SecretLock, auditReport.From.Bytes(), auditReport.Value, auditReport.Timelock.Int64(), nil
}

// AuditSecret audits the secret of an Atom swap by calling a function on ethereum
func (atom *EthereumAtom) AuditSecret() ([32]byte, error) {
	return atom.binding.AuditSecret(&bind.CallOpts{}, atom.data.SwapID)
}

// Store stores the atom details
func (atom *EthereumAtom) Store(state store.State) error {
	b, err := json.Marshal(atom.data)
	if err != nil {
		return err
	}
	return state.PutAtomDetails(atom.orderID, b)
}

// Restore restores the atom details
func (atom *EthereumAtom) Restore(state store.State) error {
	b, err := state.AtomDetails(atom.orderID)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &atom.data)
}

// From returns the address of the sender
func (atom *EthereumAtom) From() ([]byte, error) {
	return atom.key.GetAddress()
}

// PriorityCode returns the priority code of the currency.
func (atom *EthereumAtom) PriorityCode() uint32 {
	return atom.key.PriorityCode()
}

// GetSecretHash returns the Secret Hash of the atom.
func (atom *EthereumAtom) GetSecretHash() [32]byte {
	return atom.data.HashLock
}

// GetKey returns the key of the atom.
func (atom *EthereumAtom) GetKey() swap.Key {
	return atom.key
}
