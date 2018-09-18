package btc

import (
	"encoding/json"
	"errors"
	"math/big"

	bindings "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/bindings/btc"
	"github.com/republicprotocol/renex-swapper-go/adapter/blockchain/clients/btc"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/adapter/network"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
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
	key        keystore.BitcoinKey
	orderID    [32]byte
	connection btc.Conn
	network    network.Network
	data       BitcoinData
}

// NewBitcoinAtom returns a new Bitcoin Atom instance
func NewBitcoinAtom(network network.Network, conf config.BitcoinNetwork, key keystore.BitcoinKey, orderID [32]byte) (swap.Atom, error) {
	conn, err := btc.NewConnWithConfig(conf)
	if err != nil {
		return nil, err
	}
	return &BitcoinAtom{
		orderID:    orderID,
		key:        key,
		network:    network,
		connection: conn,
	}, nil
}

// Initiate a new Atom swap by calling Bitcoin
func (atom *BitcoinAtom) Initiate(to []byte, hash [32]byte, value *big.Int, expiry int64) error {
	result, err := bindings.Initiate(atom.connection, atom.key, string(to), value.Int64(), hash[:], expiry)
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
	result, err := bindings.Redeem(atom.connection, atom.key, atom.data.Contract, atom.data.ContractTx, secret)
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
}

// RedeemedAt returns the timestamp at which the atom is redeemed
func (atom *BitcoinAtom) RedeemedAt() (int64, error) {
	panic("unimplemented")
}

// Refund an Atom swap by calling Bitcoin
func (atom *BitcoinAtom) Refund() error {
	return bindings.Refund(atom.connection, atom.key, atom.data.Contract, atom.data.ContractTx)
}

// Audit an Atom swap by calling a function on Bitcoin
func (atom *BitcoinAtom) Audit() ([32]byte, []byte, *big.Int, int64, error) {
	details, err := atom.network.ReceiveSwapDetails(atom.orderID, 0)
	if err != nil {
		return [32]byte{}, nil, nil, 0, err
	}
	if err := atom.Deserialize(details); err != nil {
		return [32]byte{}, nil, nil, 0, err
	}
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

// Serialize serializes the atom details
func (atom *BitcoinAtom) Serialize() ([]byte, error) {
	return json.Marshal(atom.data)
}

// Deserialize deserializes the atom details
func (atom *BitcoinAtom) Deserialize(data []byte) error {
	return json.Unmarshal(data, &atom.data)
}

// PriorityCode returns the priority code of the currency.
func (atom *BitcoinAtom) PriorityCode() uint32 {
	return uint32(token.BTC)
}

// GetFromAddress returns the address of the sender
func (atom *BitcoinAtom) GetFromAddress() ([]byte, error) {
	return []byte(atom.key.AddressString), nil
}
