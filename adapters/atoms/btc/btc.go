package btc

import (
	"encoding/json"
	"errors"
	"math/big"

	bindings "github.com/republicprotocol/atom-go/adapters/bindings/btc"
	"github.com/republicprotocol/atom-go/adapters/clients/btc"
	"github.com/republicprotocol/atom-go/services/store"
	"github.com/republicprotocol/atom-go/services/swap"
)

type BitcoinData struct {
	ContractHash   string   `json:"contract_hash"`
	Contract       []byte   `json:"contract"`
	ContractTxHash []byte   `json:"contract_tx_hash"`
	ContractTx     []byte   `json:"contract_tx"`
	RefundTxHash   [32]byte `json:"refund_tx_hash"`
	RefundTx       []byte   `json:"refund_tx"`
	RedeemTxHash   [32]byte `json:"redeem_tx_hash"`
	RedeemTx       []byte   `json:"redeem_tx"`
	SecretHash     [32]byte `json:"secret_hash"`
}

// BitcoinAtom is a struct for Bitcoin Atom
type BitcoinAtom struct {
	key        swap.Key
	orderID    [32]byte
	connection btc.Conn
	data       BitcoinData
}

// NewBitcoinAtom returns a new Bitcoin Atom instance
func NewBitcoinAtom(connection btc.Conn, key swap.Key, orderID [32]byte) swap.Atom {
	return &BitcoinAtom{
		orderID:    orderID,
		key:        key,
		connection: connection,
	}
}

// Initiate a new Atom swap by calling Bitcoin
func (atom *BitcoinAtom) Initiate(to []byte, hash [32]byte, value *big.Int, expiry int64) error {
	from, err := atom.From()
	if err != nil {
		return err
	}
	result, err := bindings.Initiate(atom.connection, string(from), string(to), value.Int64(), hash[:], expiry)
	if err != nil {
		return err
	}

	atom.data.Contract = result.Contract
	atom.data.ContractHash = result.ContractHash
	atom.data.ContractTx = result.ContractTx
	atom.data.ContractTxHash = result.ContractTxHash
	atom.data.RefundTx = result.RefundTx
	atom.data.RefundTxHash = result.RefundTxHash
	atom.data.SecretHash = hash
	return nil
}

// Redeem an Atom swap by calling a function on Bitcoin
func (atom *BitcoinAtom) Redeem(secret [32]byte) error {
	from, err := atom.From()
	if err != nil {
		return err
	}

	result, err := bindings.Redeem(atom.connection, string(from), atom.data.Contract, atom.data.ContractTx, secret)
	if err != nil {
		return err
	}
	atom.data.RedeemTx = result.RedeemTx
	atom.data.RedeemTxHash = result.RedeemTxHash
	return nil
}

// WaitForCounterRedemption waits for the counter party to initiate
func (atom *BitcoinAtom) WaitForCounterRedemption() error {
	panic("unimplemented")
	return nil
}

// Refund an Atom swap by calling Bitcoin
func (atom *BitcoinAtom) Refund() error {
	from, err := atom.From()
	if err != nil {
		return err
	}
	return bindings.Refund(atom.connection, string(from), atom.data.Contract, atom.data.ContractTx)
}

// Audit an Atom swap by calling a function on Bitcoin
func (atom *BitcoinAtom) Audit() ([32]byte, []byte, *big.Int, int64, error) {
	result, err := bindings.Audit(atom.connection, atom.data.Contract, atom.data.ContractTx)
	if err != nil {
		return [32]byte{}, nil, nil, 0, err
	}
	return result.SecretHash, result.RecipientAddress, big.NewInt(result.Amount), result.LockTime, nil
}

// AuditSecret audits the secret of an Atom swap by calling Bitcoin
func (atom *BitcoinAtom) AuditSecret() ([32]byte, error) {
	result, err := bindings.AuditSecret(atom.connection, atom.data.RedeemTx, atom.data.SecretHash[:])
	if err != nil {
		return [32]byte{}, errors.New("Cannot read the secret")
	}
	return result, nil
}

// Store stores the atom details
func (atom *BitcoinAtom) Store(state store.State) error {
	b, err := json.Marshal(atom.data)
	if err != nil {
		return err
	}
	return state.PutAtomDetails(atom.orderID, b)
}

// Restore restores the atom details
func (atom *BitcoinAtom) Restore(state store.State) error {
	b, err := state.AtomDetails(atom.orderID)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &atom.data)
}

// PriorityCode returns the priority code of the currency.
func (atom *BitcoinAtom) PriorityCode() uint32 {
	return atom.key.PriorityCode()
}

// GetSecretHash returns the Secret Hash of the atom.
func (atom *BitcoinAtom) GetSecretHash() [32]byte {
	return atom.data.SecretHash
}

// From returns the address of the sender
func (atom *BitcoinAtom) From() ([]byte, error) {
	return atom.key.GetAddress()
}

// GetKey returns the key of the atom.
func (atom *BitcoinAtom) GetKey() swap.Key {
	return atom.key
}
